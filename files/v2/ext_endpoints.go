package files

// The endpoints an upload can go to, and the walk that tries them in order:
// the next endpoint gets the upload only if the previous one failed before any
// part of it landed.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
)

// endpointKind is a transfer protocol this package speaks.
type endpointKind int

const (
	// presignedEndpoint (IO_OFFLOAD) sends parts to presigned cloud URLs,
	// without credentials. The upload session and the URLs come from Files API
	// 2.0 on the endpoint's host. It has no single-shot upload.
	presignedEndpoint endpointKind = iota
	// filesAPIEndpoint (FILES_API 2.0) sends the data to the endpoint's host,
	// with credentials: the workspace, or the storage proxy.
	filesAPIEndpoint
	// filesAPIv3Endpoint (FILES_API 3.0) is the 2.0 protocol under /api/3.0,
	// with an upload_id for the session token and parts of at most
	// filesAPIv3MaxPartSize. It is reached with credentials over the internet,
	// so only over https.
	filesAPIv3Endpoint
)

// filesAPIv3MaxPartSize is the largest part Files API 3.0 takes, below the
// largest part size an upload may pick.
const filesAPIv3MaxPartSize int64 = 4000 << 20

func (k endpointKind) String() string {
	switch k {
	case presignedEndpoint:
		return "IO_OFFLOAD"
	case filesAPIEndpoint:
		return "FILES_API 2.0"
	case filesAPIv3Endpoint:
		return "FILES_API 3.0"
	}
	return fmt.Sprintf("endpointKind(%d)", int(k))
}

// workspaceMaxSingleShotSize is the largest single-shot upload the workspace's
// Files API takes.
const workspaceMaxSingleShotSize int64 = 5 << 30

// endpoint is one entry of the ordered list an upload walks, best first.
type endpoint struct {
	kind endpointKind
	// baseURL is the scheme and host, such as https://example.cloud.databricks.com.
	baseURL   string
	multipart bool
}

func (ep endpoint) String() string { return ep.kind.String() + " at " + ep.baseURL }

// filesPath is the path of a file in the endpoint's Files API. A presigned
// endpoint opens and completes its sessions through Files API 2.0.
func (ep endpoint) filesPath(absPath string) string {
	if ep.kind == filesAPIv3Endpoint {
		return "/api/3.0/fs/files" + absPath
	}
	return filesAPIPath(absPath)
}

// maxPartSize caps the part size of a multipart upload to the endpoint, or is
// -1 for no cap.
func (ep endpoint) maxPartSize() int64 {
	if ep.kind == filesAPIv3Endpoint {
		return filesAPIv3MaxPartSize
	}
	return -1
}

// sessionQuery returns the query that names a multipart session on the
// endpoint: Files API 3.0 calls it an upload_id.
func (ep endpoint) sessionQuery(token string) url.Values {
	if ep.kind == filesAPIv3Endpoint {
		return url.Values{"upload_id": {token}}
	}
	return url.Values{"upload_type": {"multipart"}, "session_token": {token}}
}

// builtinEndpoints is the order used without endpoint discovery: presigned
// URLs, then the workspace's Files API, which takes uploads only in a single
// shot.
func (e *engine) builtinEndpoints() []endpoint {
	return []endpoint{
		{kind: presignedEndpoint, baseURL: e.c.host, multipart: true},
		{kind: filesAPIEndpoint, baseURL: e.c.host},
	}
}

// uploadEndpoints returns the endpoints an upload walks, in order. A list of
// only presigned URLs gets the workspace's Files API, which every workspace
// serves: small and empty files need a single shot, which presigned URLs
// can't take.
func (e *engine) uploadEndpoints(ctx context.Context) []endpoint {
	endpoints := e.transferEndpoints(ctx)
	if !slices.ContainsFunc(endpoints, func(ep endpoint) bool { return ep.kind != presignedEndpoint }) {
		endpoints = append(endpoints, endpoint{kind: filesAPIEndpoint, baseURL: e.c.host})
	}
	return endpoints
}

// onWorkspace reports whether ep is the workspace host the client is
// configured with. Requests there keep the API client's retries, until the
// retry timeout; other hosts get bounded retries, so that a failing host gives
// way to the next endpoint quickly.
func (e *engine) onWorkspace(ep endpoint) bool {
	return sameOrigin(ep.baseURL, e.c.host)
}

// sameOrigin reports whether two URLs have the same scheme, host and port, with
// the scheme's default port filled in.
func sameOrigin(a, b string) bool {
	oa, erra := originOf(a)
	ob, errb := originOf(b)
	return erra == nil && errb == nil && oa == ob
}

type origin struct{ scheme, host, port string }

func originOf(raw string) (origin, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return origin{}, err
	}
	scheme := strings.ToLower(u.Scheme)
	port := u.Port()
	if port == "" {
		port = map[string]string{"http": "80", "https": "443"}[scheme]
	}
	return origin{scheme: scheme, host: strings.ToLower(u.Hostname()), port: port}, nil
}

// endpointPlan is how an upload goes to one endpoint.
type endpointPlan struct {
	ep endpoint
	// multipart is a multipart upload session; otherwise one single-shot PUT.
	multipart bool
}

// planUpload returns how the data goes to ep, or false if ep can't take it:
// presigned URLs have no single-shot upload, and the workspace takes a single
// shot of at most workspaceMaxSingleShotSize. Data below the multipart
// threshold, or to an endpoint without multipart, goes in a single shot.
func (e *engine) planUpload(ep endpoint, size int64) (endpointPlan, bool) {
	small := size >= 0 && size < max(e.tun.minStreamSize, 1)
	if !small && ep.multipart {
		return endpointPlan{ep: ep, multipart: true}, true
	}
	if ep.kind == presignedEndpoint {
		return endpointPlan{}, false
	}
	if e.onWorkspace(ep) && size > workspaceMaxSingleShotSize {
		return endpointPlan{}, false
	}
	return endpointPlan{ep: ep}, true
}

// walkUpload sends the data to each endpoint in turn until one takes it. It
// moves on only while no part has landed, the failure is one the next
// endpoint may not hit (see walksOn), and the data can be sent again from the
// start; otherwise it returns the failure.
func (e *engine) walkUpload(ctx context.Context, uc *uploadContext, data *uploadData) error {
	// A one-pass stream learns its size only once it is read: read up to the
	// multipart threshold, so a short one is planned as a single shot.
	if err := data.readAhead(e.tun.minStreamSize); err != nil {
		return err
	}
	var plans []endpointPlan
	var tooLarge []string
	for _, ep := range e.uploadEndpoints(ctx) {
		if plan, ok := e.planUpload(ep, data.size()); ok {
			plans = append(plans, plan)
		} else if ep.kind != presignedEndpoint {
			tooLarge = append(tooLarge, ep.String())
		}
	}
	if len(plans) == 0 {
		return fmt.Errorf("no endpoint can take this upload of %d bytes: %s takes a single shot of at most %d bytes",
			data.size(), strings.Join(tooLarge, ", "), workspaceMaxSingleShotSize)
	}
	for i, plan := range plans {
		uc.startEndpoint(plan.ep, i == len(plans)-1)
		err := e.uploadTo(ctx, uc, plan, data)
		if err == nil {
			return nil
		}
		if uc.committed() || uc.lastEndpoint || !walksOn(ctx, err) || !data.canSendAgain() {
			return err
		}
		e.c.logger.InfoContext(ctx, "upload failed before any part landed, trying the next endpoint",
			"endpoint", plan.ep.String(), "error", err)
		if err := data.restart(); err != nil {
			return err
		}
	}
	panic("unreachable: the last endpoint returns")
}

// walksOn reports whether a failure before any part has landed moves the
// upload to the next endpoint. Any failure does once its retries are spent,
// as the endpoint may be unreachable or blocked from here, except one the next
// endpoint would hit too (the file exists, or a precondition failed) and the
// caller's own cancellation.
func walksOn(ctx context.Context, err error) bool {
	if ctx.Err() != nil || errors.Is(err, ErrAlreadyExists) {
		return false
	}
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok {
		return aerr.Code() != codes.AlreadyExists && aerr.HTTPStatusCode() != http.StatusPreconditionFailed
	}
	return true
}

// uploadTo sends the data to one endpoint with its plan. A multipart session
// that fails is aborted once.
func (e *engine) uploadTo(ctx context.Context, uc *uploadContext, plan endpointPlan, data *uploadData) error {
	if !plan.multipart {
		return e.singleShotUpload(ctx, uc, data.whole(), data.size())
	}
	token, err := e.initiateMultipart(ctx, uc)
	if err != nil {
		return err
	}
	switch {
	case uc.parallel && data.ra != nil:
		// Parts stream from positioned reads of the source, without buffering.
		err = e.parallelMultipartFromReaderAt(ctx, uc, token, data.ra)
	case uc.parallel:
		err = e.parallelMultipartFromStream(ctx, uc, token, data.r)
	default:
		err = e.performMultipartUpload(ctx, uc, token, data.r)
	}
	if err != nil {
		e.abortMultipartBestEffort(ctx, uc, token)
	}
	return err
}

// filesAPIPartURLs returns the URLs of parts sent straight to a Files API
// host, with credentials: unlike presigned URLs, nothing is minted or expires.
func filesAPIPartURLs(uc *uploadContext, token string, startPart, count int) ([]presignedURL, error) {
	urls := make([]presignedURL, 0, count)
	for n := startPart; n < startPart+count; n++ {
		query := uc.ep.sessionQuery(token)
		query.Set("part_number", strconv.Itoa(n))
		u, err := endpointURL(uc.ep.baseURL, uc.ep.filesPath(uc.targetPath), query)
		if err != nil {
			return nil, err
		}
		urls = append(urls, presignedURL{URL: u, PartNumber: n})
	}
	return urls, nil
}

// authenticatedHTTPClient returns the HTTP client for parts sent to a Files API
// host, such as the storage proxy: the client's own transport, as for any
// Databricks host, sized like the transfer client for n concurrent parts, with
// the client's credentials on every request. A transport that isn't an
// *http.Transport is used as is. It sets no whole-request timeout, which would
// abort a long transfer.
func (e *engine) authenticatedHTTPClient(n int) *http.Client {
	base := http.DefaultTransport
	if e.c.httpClient != nil && e.c.httpClient.Transport != nil {
		base = e.c.httpClient.Transport
	}
	if t, ok := base.(*http.Transport); ok {
		base = sizedForTransfers(t, n, e.tun)
	}
	return &http.Client{Transport: &credentialTransport{base: base, creds: e.c.credentials}}
}

// credentialTransport adds credentials to every request it sends. Presigned
// URLs never go through it.
type credentialTransport struct {
	base  http.RoundTripper
	creds auth.Credentials
}

func (t *credentialTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.creds == nil {
		return t.base.RoundTrip(req)
	}
	headers, err := t.creds.AuthHeaders(req.Context())
	if err != nil {
		if req.Body != nil {
			_ = req.Body.Close()
		}
		return nil, err
	}
	// A RoundTripper must not modify the request it is given.
	req = req.Clone(req.Context())
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	return t.base.RoundTrip(req)
}
