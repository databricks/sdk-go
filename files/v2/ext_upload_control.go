package files

// Files API control-plane calls for large uploads: initiate a multipart upload
// session, mint presigned URLs for its parts, complete it, and abort it for
// cleanup. It also owns the single-shot octet-stream PUT used for small
// files and as the multipart fallback.
//
// These calls attach the client's credentials per request (via newHTTPRequest,
// the same helper the generated methods use); the parts/chunks then transfer
// directly to cloud storage over the URLs minted here (see the cloudstorage
// subpackage), carrying no Databricks credentials.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/apiretry"
	"github.com/databricks/sdk-go/core/ops"
)

// Sentinel errors surfaced by the large-file upload API.
var (
	// errUnexpectedServerResponse is returned when the Files API control plane
	// returns a well-formed but semantically incomplete response (for example,
	// no presigned URLs). It is an internal signal, not a caller-facing sentinel.
	errUnexpectedServerResponse = errors.New("unexpected server response")

	// ErrAlreadyExists is returned when an upload targets an existing path with
	// overwrite disabled. It is surfaced as a sentinel so callers detect
	// "already exists" with errors.Is regardless of which protocol produced it.
	ErrAlreadyExists = errors.New("the file being created already exists")
)

// controlPlaneRetryStatusCodes are the HTTP statuses retried for an
// authenticated control-plane call.
var controlPlaneRetryStatusCodes = []int{
	http.StatusRequestTimeout,     // 408
	http.StatusTooManyRequests,    // 429
	http.StatusBadGateway,         // 502
	http.StatusServiceUnavailable, // 503
	http.StatusGatewayTimeout,     // 504
}

// --- Result types ---

// initiateResult reports which large-file protocol the server selected. Exactly
// multipartUpload is non-nil on a well-formed response; the pointer
// distinguishes "not offered" from "offered with an empty token".
type initiateResult struct {
	MultipartUpload *uploadSession `json:"multipart_upload"`
	// UploadID names the session on Files API 3.0.
	UploadID string `json:"upload_id"`
}

type uploadSession struct {
	SessionToken string `json:"session_token"`
}

// presignedURL is a resolved cloud-storage URL with its associated request
// headers, as minted by the control plane and transferred over by the caller.
type presignedURL struct {
	URL        string
	PartNumber int
	Headers    map[string]string
}

// --- Wire types for the Files API multipart coordination ---

type nameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type createPartURLsRequest struct {
	Path            string `json:"path"`
	SessionToken    string `json:"session_token"`
	StartPartNumber int    `json:"start_part_number"`
	Count           int    `json:"count"`
	ExpireTime      string `json:"expire_time"`
}

type createPartURLsResponse struct {
	UploadPartURLs []struct {
		URL        string      `json:"url"`
		PartNumber int         `json:"part_number"`
		Headers    []nameValue `json:"headers"`
	} `json:"upload_part_urls"`
}

type completePart struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

type completeRequest struct {
	Parts []completePart `json:"parts"`
}

// uploadSingleShot sends the whole body in one PUT to the current endpoint's
// Files API. size is the body's length, or -1 if unknown. The upload's overwrite
// is tri-state: a non-nil value is sent explicitly; nil lets the server apply
// its default. A 409 ALREADY_EXISTS is mapped to ErrAlreadyExists.
//
// The generated UploadFile is unsuitable here: it JSON-marshals the reader
// rather than streaming the raw octet stream the Files API expects. A seekable
// body is sent again on each retry (see retryBodies); a non-seekable body cannot
// be replayed, so it is attempted once.
func (e *engine) uploadSingleShot(ctx context.Context, uc *uploadContext, body io.Reader, size int64) error {
	query := url.Values{}
	if uc.overwrite != nil {
		query.Set("overwrite", strconv.FormatBool(*uc.overwrite))
	}
	urlStr, err := endpointURL(uc.ep.baseURL, uc.ep.filesPath(uc.targetPath), query)
	if err != nil {
		return err
	}

	bodies, seekable := newRetryBodies(body)

	call := func(ctx context.Context) error {
		attempt := body
		if seekable {
			next, err := bodies.next(ctx)
			if err != nil {
				return err
			}
			attempt = next
		} else if _, ok := body.(io.Closer); ok {
			// net/http closes a request body that is an io.Closer once the
			// attempt ends, and the body is the caller's.
			attempt = struct{ io.Reader }{body}
		}
		headers := http.Header{}
		headers.Set("Content-Type", "application/octet-stream")
		req, err := newHTTPRequest(ctx, httpRequestOptions{
			Method:      http.MethodPut,
			URL:         urlStr,
			Credentials: e.c.credentials,
			Headers:     headers,
			Body:        attempt,
		})
		if err != nil {
			return err
		}
		// Hosts require a Content-Length, which net/http knows only for a few
		// in-memory body types, and never for an attempt's body; an empty body
		// must be no body to send one.
		length := size
		if seekable {
			length, req.GetBody = bodies.size, bodies.getBody()
		}
		if length == 0 {
			_ = req.Body.Close()
			req.Body, req.GetBody = http.NoBody, nil
		}
		if length >= 0 {
			req.ContentLength = length
		}
		e.setWorkspaceHeader(req)
		_, _, err = executeHTTPCall(httpCallOptions{req: req, client: e.c.httpClient, logger: e.c.logger})
		return err
	}

	// Only retry when the body can be rewound; a consumed non-seekable stream
	// cannot be replayed, so it is attempted exactly once.
	opts := []ops.Option{ops.WithTimeout(e.tun.cpRetryTimeout)}
	if seekable {
		opts = append(opts, ops.WithRetrier(e.controlPlaneRetrier(uc)))
	}
	return e.asAlreadyExistsOn(uc, ops.Execute(ctx, call, opts...))
}

// retryBodies gives each attempt of a request its own body over the same
// seekable data. net/http may still read an attempt's body after the attempt
// returns, until it closes the body, so a retry that rewound a shared body
// sooner could lose bytes to that late read.
type retryBodies struct {
	body        io.Reader
	start, size int64
	// last is the latest attempt's body, for data without random access.
	last *closeSignal
}

// newRetryBodies returns the bodies over body from its current position to its
// end, or false if body cannot be sent again.
func newRetryBodies(body io.Reader) (*retryBodies, bool) {
	seeker, ok := body.(io.Seeker)
	if !ok {
		return nil, false
	}
	start, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, false
	}
	end, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, false
	}
	return &retryBodies{body: body, start: start, size: end - start}, true
}

// next returns the body of the next attempt. Data with random access gives each
// attempt its own section; other data is rewound only once the last attempt's
// body is closed.
func (b *retryBodies) next(ctx context.Context) (io.Reader, error) {
	if ra, ok := b.body.(io.ReaderAt); ok {
		return io.NewSectionReader(ra, b.start, b.size), nil
	}
	if b.last != nil {
		select {
		case <-b.last.closed:
		case <-ctx.Done():
			return nil, context.Cause(ctx)
		}
	}
	if _, err := b.body.(io.Seeker).Seek(b.start, io.SeekStart); err != nil {
		return nil, fmt.Errorf("rewinding upload body: %w", err)
	}
	b.last = &closeSignal{Reader: b.body, closed: make(chan struct{})}
	return b.last, nil
}

// getBody returns a function for http.Request.GetBody, which net/http calls to
// send the body again itself, or nil for data without random access.
func (b *retryBodies) getBody() func() (io.ReadCloser, error) {
	ra, ok := b.body.(io.ReaderAt)
	if !ok {
		return nil
	}
	return func() (io.ReadCloser, error) {
		return io.NopCloser(io.NewSectionReader(ra, b.start, b.size)), nil
	}
}

// closeSignal reports when net/http closes a request body, without closing
// the data under it.
type closeSignal struct {
	io.Reader
	once   sync.Once
	closed chan struct{}
}

func (c *closeSignal) Close() error {
	c.once.Do(func() { close(c.closed) })
	return nil
}

// initiateUpload opens a multipart upload session on the current endpoint.
func (e *engine) initiateUpload(ctx context.Context, uc *uploadContext) (*initiateResult, error) {
	query := url.Values{"action": {"initiate-upload"}}
	if uc.overwrite != nil {
		query.Set("overwrite", strconv.FormatBool(*uc.overwrite))
	}
	var out initiateResult
	if err := e.controlPlaneJSON(ctx, uc, http.MethodPost, uc.ep.filesPath(uc.targetPath), query, nil, &out); err != nil {
		return nil, e.asAlreadyExistsOn(uc, err)
	}
	return &out, nil
}

// createPartURLs mints count presigned URLs for multipart parts starting at
// startPart.
func (e *engine) createPartURLs(ctx context.Context, uc *uploadContext, token string, startPart, count int) ([]presignedURL, error) {
	if uc.ep.kind != presignedEndpoint {
		return filesAPIPartURLs(uc, token, startPart, count)
	}
	body := createPartURLsRequest{
		Path:            uc.targetPath,
		SessionToken:    token,
		StartPartNumber: startPart,
		Count:           count,
		ExpireTime:      e.expireTime(),
	}
	var out createPartURLsResponse
	if err := e.controlPlaneJSON(ctx, uc, http.MethodPost, "/api/2.0/fs/create-upload-part-urls", nil, body, &out); err != nil {
		return nil, err
	}
	if len(out.UploadPartURLs) == 0 {
		return nil, fmt.Errorf("%w: no upload part URLs returned", errUnexpectedServerResponse)
	}
	result := make([]presignedURL, 0, len(out.UploadPartURLs))
	for _, p := range out.UploadPartURLs {
		result = append(result, presignedURL{URL: p.URL, PartNumber: p.PartNumber, Headers: headerMap(p.Headers)})
	}
	return result, nil
}

// completeMultipart finalizes a multipart upload with the part ETags, which it
// sorts by part number before sending.
func (e *engine) completeMultipart(ctx context.Context, uc *uploadContext, token string, etags map[int]string) error {
	nums := slices.Sorted(maps.Keys(etags))
	parts := make([]completePart, 0, len(nums))
	for _, n := range nums {
		parts = append(parts, completePart{PartNumber: n, ETag: etags[n]})
	}
	query := uc.ep.sessionQuery(token)
	query.Set("action", "complete-upload")
	return e.asAlreadyExistsOn(uc, e.controlPlaneJSON(ctx, uc, http.MethodPost, uc.ep.filesPath(uc.targetPath), query, completeRequest{Parts: parts}, nil))
}

// abortMultipart cancels a multipart upload session with the Files API
// abort-upload action. It is best-effort cleanup, so it is sent once.
func (e *engine) abortMultipart(ctx context.Context, uc *uploadContext, token string) error {
	query := uc.ep.sessionQuery(token)
	query.Set("action", "abort-upload")
	call, err := e.controlPlaneCall(uc.ep.baseURL, http.MethodPost, uc.ep.filesPath(uc.targetPath), query, nil, nil)
	if err != nil {
		return err
	}
	return call(ctx)
}

// controlPlaneJSON performs an authenticated JSON request against the Files API
// on the current endpoint's host, retrying transient failures via core/ops.
// reqBody and out may be nil. A fresh request is built on each attempt so
// retries re-apply credentials (newHTTPRequest re-reads them, picking up any
// token refresh) and rewind the body.
func (e *engine) controlPlaneJSON(ctx context.Context, uc *uploadContext, method, path string, query url.Values, reqBody, out any) error {
	call, err := e.controlPlaneCall(uc.ep.baseURL, method, path, query, reqBody, out)
	if err != nil {
		return err
	}
	return ops.Execute(ctx, call, ops.WithRetrier(e.controlPlaneRetrier(uc)), ops.WithTimeout(e.tun.cpRetryTimeout))
}

// controlPlaneCall returns one attempt of an authenticated JSON request against
// the Files API on baseURL. Each call builds a fresh request, so a retry
// re-applies credentials and rewinds the body.
func (e *engine) controlPlaneCall(baseURL, method, path string, query url.Values, reqBody, out any) (func(context.Context) error, error) {
	urlStr, err := endpointURL(baseURL, path, query)
	if err != nil {
		return nil, err
	}
	var bodyBytes []byte
	if reqBody != nil {
		if bodyBytes, err = json.Marshal(reqBody); err != nil {
			return nil, err
		}
	}

	call := func(ctx context.Context) error {
		var rdr io.Reader
		if bodyBytes != nil {
			rdr = bytes.NewReader(bodyBytes)
		}
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")
		req, err := newHTTPRequest(ctx, httpRequestOptions{
			Method:      method,
			URL:         urlStr,
			Credentials: e.c.credentials,
			Headers:     headers,
			Body:        rdr,
		})
		if err != nil {
			return err
		}
		e.setWorkspaceHeader(req)
		respBody, _, err := executeHTTPCall(httpCallOptions{req: req, client: e.c.httpClient, logger: e.c.logger})
		if err != nil {
			return err
		}
		if out != nil && len(respBody) > 0 {
			return json.Unmarshal(respBody, out)
		}
		return nil
	}
	return call, nil
}

func (e *engine) newControlPlaneRetrier() ops.Retrier {
	return apiretry.NewRetrier(e.tun.cpBackoff, apiretry.RetrierConfig{StatusCodes: controlPlaneRetryStatusCodes})
}

// controlPlaneRetrier returns retriers for a request to the current endpoint,
// capped at the retries the walk allows it now (see uploadContext.retries).
func (e *engine) controlPlaneRetrier(uc *uploadContext) func() ops.Retrier {
	n := uc.retries(e.onWorkspace(uc.ep))
	return func() ops.Retrier {
		if n < 0 {
			return e.newControlPlaneRetrier()
		}
		return &cappedRetrier{Retrier: e.newControlPlaneRetrier(), left: n}
	}
}

// cappedRetrier retries what its Retrier does, at most left more times.
type cappedRetrier struct {
	ops.Retrier
	left int
}

func (r *cappedRetrier) IsRetriable(err error) (time.Duration, bool) {
	if r.left <= 0 {
		return 0, false
	}
	d, ok := r.Retrier.IsRetriable(err)
	if ok {
		r.left--
	}
	return d, ok
}

// setWorkspaceHeader applies the workspace routing header when the client is
// workspace-scoped. The Files API control plane routes a request to the right
// workspace by this header; the credentials are attached separately by
// newHTTPRequest.
func (e *engine) setWorkspaceHeader(req *http.Request) {
	if e.c.workspaceID != "" {
		req.Header.Set("X-Databricks-Workspace-Id", e.c.workspaceID)
	}
}

// endpointURL joins an endpoint's base URL with a path and optional query,
// mirroring the generated client's URL construction.
func endpointURL(base, path string, query url.Values) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	u.Path = path
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	return u.String(), nil
}

// asAlreadyExistsOn is asAlreadyExists for a request to the current endpoint.
// The storage proxy answers an existing file with a 412 rather than the
// workspace's 409 ALREADY_EXISTS, which an upload that doesn't overwrite maps
// to ErrAlreadyExists too.
func (e *engine) asAlreadyExistsOn(uc *uploadContext, err error) error {
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok && aerr.HTTPStatusCode() == http.StatusPreconditionFailed &&
		!e.onWorkspace(uc.ep) && uc.overwrite != nil && !*uc.overwrite {
		return ErrAlreadyExists
	}
	return asAlreadyExists(err)
}

// asAlreadyExists maps an "already exists" API error to the ErrAlreadyExists
// sentinel, so callers detect it with errors.Is regardless of which path
// surfaced it. The Files API returns HTTP 409 with error_code ALREADY_EXISTS
// when an upload targets an existing path with overwrite=false; for multipart
// this surfaces from complete-upload (after the parts are sent). Every path runs
// over core/apierr, so a single code check suffices. Other errors pass through
// unchanged.
func asAlreadyExists(err error) error {
	if apierr.Code(err) == codes.AlreadyExists {
		return ErrAlreadyExists
	}
	return err
}

// filesAPIPath builds the Files API URL path for an absolute volume path. It
// matches the single-shot UploadFile path construction so the control-plane
// initiate/complete calls escape the path identically.
func filesAPIPath(absPath string) string {
	return fmt.Sprintf("/api/2.0/fs/files%v", absPath)
}

// expireTime is the presigned-URL expiry timestamp sent to the control plane,
// urlExpiry into the future.
func (e *engine) expireTime() string {
	return time.Now().UTC().Add(e.tun.urlExpiry).Format("2006-01-02T15:04:05Z")
}

func headerMap(headers []nameValue) map[string]string {
	out := make(map[string]string, len(headers))
	for _, h := range headers {
		out[h.Name] = h.Value
	}
	return out
}
