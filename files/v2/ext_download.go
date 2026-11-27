package files

// Downloads: a byte range of a file from one endpoint, resumed after a break,
// and the endpoint walk that moves on only while no byte has arrived.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apiretry"
	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/files/v2/internal/cloudstorage"
)

// errDownloadStalled cancels a download attempt that received no bytes for the
// download idle timeout, so that it resumes on a fresh connection.
var errDownloadStalled = errors.New("download stalled")

// errDownloadURLExpired is returned when a presigned download URL keeps
// expiring across refetches.
var errDownloadURLExpired = errors.New("download URL expired after retries")

// errDownloadClosed is a Read of a download that was closed, including one
// that was blocked when Close ran.
var errDownloadClosed = errors.New("read of a closed download")

// byteRange is a range of a file's bytes. end is inclusive, or -1 for the end
// of the file.
type byteRange struct{ start, end int64 }

func (r byteRange) wholeFile() bool { return r.start == 0 && r.end < 0 }

func (r byteRange) header() string {
	if r.end < 0 {
		return fmt.Sprintf("bytes=%d-", r.start)
	}
	return fmt.Sprintf("bytes=%d-%d", r.start, r.end)
}

// downloadContext is one download's state on one endpoint, shared by the
// goroutines of a parallel download.
type downloadContext struct {
	tun  tunables
	path string
	ep   endpoint
	// last is whether no endpoint is left after ep.
	last      bool
	committed atomic.Bool
	// cloud sends presigned GETs, without credentials; files sends GETs to a
	// Files API host other than the workspace, with them.
	cloud, files *http.Client
	// url is the presigned URL every range of the download shares.
	url *presignedDownloadURL

	// breaks are the breaks resumeAfter gave up on. A break is a passing
	// fault, so it does not count against the endpoint.
	breaksMu sync.Mutex
	breaks   []error
}

// gaveUpOn records a break the download gave up resuming after.
func (dc *downloadContext) gaveUpOn(breakErr error) {
	dc.breaksMu.Lock()
	defer dc.breaksMu.Unlock()
	dc.breaks = append(dc.breaks, breakErr)
}

// countsAgainst reports whether a failure the walk moved on from counts
// against the endpoint (see countsAgainstEndpoint). A break the download gave
// up resuming after doesn't.
func (dc *downloadContext) countsAgainst(ctx context.Context, err error) bool {
	if !countsAgainstEndpoint(ctx, err) {
		return false
	}
	dc.breaksMu.Lock()
	defer dc.breaksMu.Unlock()
	return !slices.ContainsFunc(dc.breaks, func(b error) bool { return errors.Is(err, b) })
}

// retries is how often a request to the endpoint may be retried, or -1 for no
// cap but the retry timeout. Before the first byte arrives, with an endpoint
// left, a download moves on without retrying. Otherwise the workspace host
// keeps the API client's retries, and other hosts get cloudRetries.
func (dc *downloadContext) retries(onWorkspace bool) int {
	switch {
	case !dc.committed.Load() && !dc.last:
		return 0
	case onWorkspace:
		return -1
	}
	return dc.tun.cloudRetries
}

// downloadConfig holds the options of a download.
type downloadConfig struct {
	// transferClient sends presigned GETs, without credentials. When nil, a
	// client sized to the download's parallelism is used.
	transferClient *http.Client

	// The options of a download to a file. noOverwrite fails it if the
	// destination exists, so that the zero value overwrites. parallel
	// downloads a large file in ranges on up to parallelism goroutines, or
	// tunables.downloadParallelism when it is nil.
	noOverwrite bool
	parallel    bool
	parallelism *int
}

// resolvedParallelism is how many ranges a parallel download runs at once.
func (cfg downloadConfig) resolvedParallelism(tun tunables) int {
	if cfg.parallelism != nil {
		return *cfg.parallelism
	}
	return tun.downloadParallelism
}

func (e *engine) newDownloadContext(cfg downloadConfig, path string, ep endpoint, last bool, parallelism int) *downloadContext {
	dc := &downloadContext{tun: e.tun, path: path, ep: ep, last: last, url: &presignedDownloadURL{}}
	switch {
	case ep.kind == presignedEndpoint && cfg.transferClient != nil:
		dc.cloud = cfg.transferClient
	case ep.kind == presignedEndpoint:
		dc.cloud = newTransferClient(parallelism, e.tun)
		http1Only(dc.cloud)
	case !e.onWorkspace(ep):
		// Not limited to HTTP/1.1: a Files API 3.0 host closes HTTP/1.1
		// downloads without answering.
		dc.files = e.authenticatedHTTPClient(parallelism)
	}
	return dc
}

// http1Only keeps the client for presigned downloads to HTTP/1.1, so that
// each range of a parallel download gets its own connection to cloud storage:
// over HTTP/2 they would share one, which made a parallel download of a
// 120 MiB file three times slower. A transport that isn't an *http.Transport,
// such as a test's, is left as is.
func http1Only(c *http.Client) {
	if t, ok := c.Transport.(*http.Transport); ok {
		t.Protocols = new(http.Protocols)
		t.Protocols.SetHTTP1(true)
	}
}

// presignedDownloadURL is the presigned URL of one download, created on first
// use and replaced once it expires. Until the endpoint is committed, a failure
// to create it is kept, so that every range sees the failure and none commits
// an endpoint that is failing; a cancelled request's failure is not kept.
type presignedDownloadURL struct {
	mu      sync.Mutex
	url     *presignedURL
	version int
	err     error
}

func (u *presignedDownloadURL) get(ctx context.Context, committed bool, mint func(context.Context) (*presignedURL, error)) (*presignedURL, int, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.err != nil {
		return nil, 0, u.err
	}
	if u.url == nil {
		url, err := mint(ctx)
		if err != nil {
			if !committed && ctx.Err() == nil {
				u.err = err
			}
			return nil, 0, err
		}
		u.url = url
	}
	return u.url, u.version, nil
}

// invalidate drops the URL if it is still the given version, so that only the
// first range to see it expire replaces it.
func (u *presignedDownloadURL) invalidate(version int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if version == u.version {
		u.url = nil
		u.version++
	}
}

// createDownloadURL mints a presigned URL for the file through the presigned
// endpoint's host.
func (e *engine) createDownloadURL(ctx context.Context, dc *downloadContext) (*presignedURL, error) {
	query := map[string][]string{"path": {dc.path}, "expire_time": {e.expireTime()}}
	var out struct {
		URL     string      `json:"url"`
		Headers []nameValue `json:"headers"`
	}
	call, err := e.controlPlaneCall(dc.ep.baseURL, http.MethodPost, "/api/2.0/fs/create-download-url", query, nil, &out)
	if err != nil {
		return nil, err
	}
	retrier := func() ops.Retrier {
		return &getRetrier{statuses: newStatusRetrier(e.tun.cpBackoff), backoff: e.tun.cpBackoff, left: dc.retries(e.onWorkspace(dc.ep))}
	}
	if err := ops.Execute(ctx, call, ops.WithRetrier(retrier), ops.WithTimeout(e.tun.cpRetryTimeout)); err != nil {
		return nil, err
	}
	if out.URL == "" {
		return nil, fmt.Errorf("%w: no download URL returned", errUnexpectedServerResponse)
	}
	return &presignedURL{URL: out.URL, Headers: headerMap(out.Headers)}, nil
}

func newStatusRetrier(bp ops.BackoffPolicy) ops.Retrier {
	return apiretry.NewRetrier(bp, apiretry.RetrierConfig{StatusCodes: controlPlaneRetryStatusCodes})
}

// getRetrier retries a download's requests, at most left more times, or
// without a cap for left < 0: what statuses retries (the retryable statuses
// and transient network failures), and any other failure of a request that got
// no answer, such as a connection the server closed (io.ErrUnexpectedEOF),
// which the API client's rules don't count as transient. These requests can
// always be sent again.
type getRetrier struct {
	statuses ops.Retrier
	backoff  ops.BackoffPolicy
	left     int
}

func (r *getRetrier) IsRetriable(err error) (time.Duration, bool) {
	if r.left == 0 {
		return 0, false
	}
	d, ok := r.statuses.IsRetriable(err)
	if !ok && unanswered(err) {
		d, ok = r.backoff.Delay(), true
	}
	if ok && r.left > 0 {
		r.left--
	}
	return d, ok
}

// unanswered reports whether err is a request's transport failure, other
// than its context ending.
func unanswered(err error) bool {
	var uerr *url.Error
	return errors.As(err, &uerr) && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
}

// getRange opens requested on the download's endpoint, retrying within the
// endpoint's budget. ifUnmodifiedSince guards against a file that changed since
// the download started; expectedLength, or -1, is the length a 200 to a range
// request must have to be accepted in place of a 206.
func (e *engine) getRange(ctx context.Context, dc *downloadContext, requested byteRange, ifUnmodifiedSince string, expectedLength int64) (*http.Response, error) {
	onWorkspace := dc.ep.kind != presignedEndpoint && e.onWorkspace(dc.ep)
	backoff := e.tun.cloudBackoff
	if onWorkspace {
		backoff = e.tun.cpBackoff
	}
	for refreshes := 0; ; refreshes++ {
		var resp *http.Response
		// The attempts use ctx rather than opCtx, the one ops.Execute passes,
		// which it cancels on return, while the body is still to be read. A
		// presigned URL is minted under opCtx, so that the mint's own retries
		// end with the range's budget.
		attempt := func(opCtx context.Context) error {
			var err error
			resp, err = e.getRangeOnce(ctx, opCtx, dc, onWorkspace, requested, ifUnmodifiedSince, expectedLength)
			return err
		}
		retrier := func() ops.Retrier {
			return &getRetrier{statuses: newStatusRetrier(backoff), backoff: backoff, left: dc.retries(onWorkspace)}
		}
		err := ops.Execute(ctx, attempt, ops.WithRetrier(retrier), ops.WithTimeout(e.tun.cpRetryTimeout))
		if err == nil {
			return resp, nil
		}
		var expired *urlExpiredError
		if !errors.As(err, &expired) {
			return nil, err
		}
		if refreshes >= e.tun.maxRetries {
			return nil, errDownloadURLExpired
		}
		dc.url.invalidate(expired.version)
	}
}

// urlExpiredError is a presigned GET answered with the cloud's "URL expired".
type urlExpiredError struct {
	version int
	err     error
}

func (e *urlExpiredError) Error() string {
	return "the presigned download URL expired: " + e.err.Error()
}
func (e *urlExpiredError) Unwrap() error { return e.err }

func (e *engine) getRangeOnce(ctx, mintCtx context.Context, dc *downloadContext, onWorkspace bool, requested byteRange, ifUnmodifiedSince string, expectedLength int64) (*http.Response, error) {
	headers := http.Header{}
	headers.Set("Accept", "application/octet-stream")
	if !requested.wholeFile() {
		headers.Set("Range", requested.header())
	}
	if ifUnmodifiedSince != "" {
		headers.Set("If-Unmodified-Since", ifUnmodifiedSince)
	}
	var req *http.Request
	var client *http.Client
	version := 0
	switch {
	case dc.ep.kind == presignedEndpoint:
		url, v, err := dc.url.get(mintCtx, dc.committed.Load(), func(ctx context.Context) (*presignedURL, error) { return e.createDownloadURL(ctx, dc) })
		if err != nil {
			return nil, err
		}
		version = v
		for name, value := range url.Headers {
			if headers.Get(name) != "" {
				return nil, fmt.Errorf("%w: the presigned URL's headers overlap the request's: %s", errUnexpectedServerResponse, name)
			}
			headers.Set(name, value)
		}
		if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url.URL, nil); err != nil {
			return nil, cloudstorage.RedactURLError(err)
		}
		req.Header = headers
		client = dc.cloud
	default:
		urlStr, err := endpointURL(dc.ep.baseURL, dc.ep.filesPath(dc.path), nil)
		if err != nil {
			return nil, err
		}
		client = dc.files
		var creds = e.c.credentials
		if onWorkspace {
			client = e.c.httpClient
		} else {
			creds = nil // the files client adds them
		}
		if req, err = newHTTPRequest(ctx, httpRequestOptions{Method: http.MethodGet, URL: urlStr, Credentials: creds, UserAgent: e.userAgent, Headers: headers}); err != nil {
			return nil, err
		}
		e.setWorkspaceHeader(req)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, cloudstorage.RedactURLError(err)
	}
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusPartialContent {
		if err := checkRange(requested, resp, expectedLength); err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		return resp, nil
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	_ = resp.Body.Close()
	// 1xx and 3xx statuses aren't errors to FromHTTPError, nor a download.
	err = fmt.Errorf("%w: status %d", errUnexpectedServerResponse, resp.StatusCode)
	if aerr := apierr.FromHTTPError(resp.StatusCode, resp.Header, body); aerr != nil {
		err = aerr
	}
	if dc.ep.kind == presignedEndpoint && cloudstorage.IsURLExpired(err) {
		return nil, &urlExpiredError{version: version, err: err}
	}
	return nil, err
}

// checkRange accepts a response to requested: a 206, or a 200 of exactly the
// requested length. A 200 with the whole file must not be written at the
// range's offset.
func checkRange(requested byteRange, resp *http.Response, expectedLength int64) error {
	if requested.wholeFile() || resp.StatusCode == http.StatusPartialContent {
		return nil
	}
	if resp.StatusCode == http.StatusOK && expectedLength >= 0 && resp.ContentLength == expectedLength {
		return nil
	}
	return fmt.Errorf("%w: status %d for %s", errUnexpectedServerResponse, resp.StatusCode, requested.header())
}

// resilientDownload reads one byte range from one endpoint. After a break it
// resumes from the current offset with a Range request guarded by
// If-Unmodified-Since (the caller's value, or else the first response's
// Last-Modified), within the resume budget.
type resilientDownload struct {
	e   *engine
	ctx context.Context
	dc  *downloadContext
	rng byteRange

	ifUnmodifiedSince string
	offset            int64
	// total is the end of the file, when the range runs to it and the first
	// response told its length; otherwise -1.
	total int64

	// mu guards body and closed: Close may run while a Read is blocked.
	mu     sync.Mutex
	body   io.ReadCloser
	closed bool
	// pending is a read error that came with bytes, for the next Read.
	pending error

	recovers, recoversWithoutProgress int

	// The first response's headers.
	contentLength           int64
	contentType, modifiedAt string
}

// openDownload opens rng on the download's endpoint.
func (e *engine) openDownload(ctx context.Context, dc *downloadContext, rng byteRange, ifUnmodifiedSince string) (*resilientDownload, error) {
	d := &resilientDownload{e: e, ctx: ctx, dc: dc, rng: rng, ifUnmodifiedSince: ifUnmodifiedSince, offset: rng.start, total: -1}
	resp, err := d.openAt(d.offset)
	if err != nil {
		return nil, err
	}
	d.contentLength = resp.ContentLength
	d.contentType = resp.Header.Get("Content-Type")
	d.modifiedAt = resp.Header.Get("Last-Modified")
	if d.ifUnmodifiedSince == "" {
		d.ifUnmodifiedSince = d.modifiedAt
	}
	if rng.end < 0 && resp.ContentLength >= 0 {
		d.total = rng.start + resp.ContentLength
	}
	return d, nil
}

// end is the offset just past the last byte the download expects, or -1.
func (d *resilientDownload) end() int64 {
	if d.rng.end >= 0 {
		return d.rng.end + 1
	}
	return d.total
}

func (d *resilientDownload) openAt(offset int64) (*http.Response, error) {
	requested := byteRange{start: offset, end: d.rng.end}
	expected := int64(-1)
	if end := d.end(); end >= 0 {
		expected = end - offset
	}
	// Each attempt gets its own context, cancelled if its body stalls.
	ctx, cancel := context.WithCancelCause(d.ctx)
	resp, err := d.e.getRange(ctx, d.dc, requested, d.ifUnmodifiedSince, expected)
	if err != nil {
		cancel(nil)
		return nil, err
	}
	body := newIdleTimeoutBody(resp.Body, d.dc.tun.transferIdleTimeout, cancel)
	if !d.setBody(body) {
		// Close ran while the request was out, so nothing else would close it.
		_ = body.Close()
		return nil, errDownloadClosed
	}
	return resp, nil
}

// setBody makes body the one Read reads, unless Close has run.
func (d *resilientDownload) setBody(body io.ReadCloser) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.closed {
		return false
	}
	d.body = body
	return true
}

func (d *resilientDownload) isClosed() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.closed
}

func (d *resilientDownload) Read(p []byte) (int, error) {
	for {
		var n int
		var err error
		if d.pending != nil {
			err, d.pending = d.pending, nil
		} else {
			n, err = d.body.Read(p)
		}
		if n > 0 {
			if end := d.end(); end >= 0 && d.offset+int64(n) > end {
				// Bytes past the end would overwrite the next range of a
				// parallel download.
				return 0, fmt.Errorf("%w: the download went past byte %d", errUnexpectedServerResponse, end)
			}
			d.offset += int64(n)
			d.recoversWithoutProgress = 0
			if err != nil && err != io.EOF {
				d.pending = err
			}
			return n, nil
		}
		switch {
		case err == nil:
			continue
		case err == io.EOF:
			end := d.end()
			if end < 0 || d.offset >= end {
				return 0, io.EOF
			}
			// A body cut short without an error is a break too.
			err = fmt.Errorf("the download ended at byte %d of %d: %w", d.offset, end, io.ErrUnexpectedEOF)
		case d.ctx.Err() != nil:
			// The caller's own cancellation is not a break.
			return 0, err
		}
		if rerr := d.resumeAfter(err); rerr != nil {
			return 0, rerr
		}
	}
}

// resumeAfter reopens the range at the current offset after a break. A failed
// reopen surfaces the break, as the caller saw it.
func (d *resilientDownload) resumeAfter(breakErr error) error {
	if d.isClosed() {
		// The break is Close's own, or came as it ran.
		return errDownloadClosed
	}
	if !(d.dc.committed.Load() || d.dc.last) || d.recoversWithoutProgress >= d.dc.tun.downloadResumesWithoutProgress {
		d.dc.gaveUpOn(breakErr)
		return breakErr
	}
	if d.ifUnmodifiedSince == "" {
		// Without it, a file changed in between would be spliced from two
		// versions.
		d.dc.gaveUpOn(breakErr)
		return breakErr
	}
	d.recovers++
	d.recoversWithoutProgress++
	_ = d.body.Close()
	d.e.c.logger.DebugContext(d.ctx, "resuming the download", "offset", d.offset, "error", breakErr)
	if _, err := d.openAt(d.offset); err != nil {
		if errors.Is(err, errDownloadClosed) {
			return err
		}
		d.setBody(io.NopCloser(errReaderOf{breakErr}))
		d.dc.gaveUpOn(breakErr)
		return breakErr
	}
	return nil
}

func (d *resilientDownload) Close() error {
	d.mu.Lock()
	d.closed = true
	body := d.body
	d.mu.Unlock()
	return body.Close()
}

type errReaderOf struct{ err error }

func (r errReaderOf) Read([]byte) (int, error) { return 0, r.err }

// idleTimeoutBody cancels its request, with errDownloadStalled, when a Read
// stays blocked for the timeout, so that a wedged connection breaks rather
// than hangs. The caller's own pauses between reads don't count.
type idleTimeoutBody struct {
	io.ReadCloser
	timer  *time.Timer
	idle   time.Duration
	cancel context.CancelCauseFunc
}

func newIdleTimeoutBody(body io.ReadCloser, idle time.Duration, cancel context.CancelCauseFunc) io.ReadCloser {
	b := &idleTimeoutBody{ReadCloser: body, idle: idle, cancel: cancel}
	if idle > 0 {
		b.timer = time.AfterFunc(idle, func() { cancel(errDownloadStalled) })
		b.timer.Stop()
	}
	return b
}

func (b *idleTimeoutBody) Read(p []byte) (int, error) {
	if b.timer != nil {
		b.timer.Reset(b.idle)
		defer b.timer.Stop()
	}
	return b.ReadCloser.Read(p)
}

func (b *idleTimeoutBody) Close() error {
	if b.timer != nil {
		b.timer.Stop()
	}
	err := b.ReadCloser.Close()
	b.cancel(nil)
	return err
}

// mayWalkDownload reports whether a failed download moves on to the next
// endpoint: only before any byte arrived, and not for a failure the next
// endpoint would hit too.
func mayWalkDownload(ctx context.Context, dc *downloadContext, err error) bool {
	return !dc.committed.Load() && !dc.last && walksOn(ctx, err)
}

// walkingDownload is the endpoint walk of download, driven by the caller's
// reads. Close may run while a Read is blocked: it ends the walk's context, so
// that a request in flight stops, and no endpoint is opened after it.
type walkingDownload struct {
	e   *engine
	cfg downloadConfig
	// ctx is the caller's context, ended by Close too.
	ctx       context.Context
	stop      context.CancelCauseFunc
	path      string
	endpoints []endpoint
	index     int
	dc        *downloadContext
	// modifiedAt pins the file once download returned the first endpoint's
	// headers: a file changed since must fail, not be served under them.
	modifiedAt string

	// The endpoint health the walk was ordered by, the endpoints it moved on
	// from on a failure that counts against them, and whether it reported.
	health     *endpointHealth
	order      endpointOrder
	walkedPast []endpoint
	reported   bool

	// mu guards d and closed.
	mu     sync.Mutex
	d      *resilientDownload
	closed bool
}

// download opens the file for reading from the first endpoint that answers.
func (e *engine) download(ctx context.Context, path string, cfg downloadConfig) (*DownloadFileResponse, error) {
	health := e.endpointHealth(false)
	order := health.order(e.clock(), e.transferEndpoints(ctx))
	walkCtx, stop := context.WithCancelCause(ctx)
	w := &walkingDownload{e: e, cfg: cfg, ctx: walkCtx, stop: stop, path: path, endpoints: order.endpoints, index: -1, health: health, order: order}
	if err := w.openNext(); err != nil {
		stop(nil)
		return nil, err
	}
	w.modifiedAt = w.d.modifiedAt
	resp := &DownloadFileResponse{Contents: w}
	if w.d.contentLength >= 0 {
		resp.ContentLength = new(w.d.contentLength)
	}
	if w.d.contentType != "" {
		resp.ContentType = new(w.d.contentType)
	}
	if w.d.modifiedAt != "" {
		resp.LastModified = new(w.d.modifiedAt)
	}
	return resp, nil
}

// openNext opens the next endpoint that answers, from byte 0.
func (w *walkingDownload) openNext() error {
	for {
		w.index++
		ep := w.endpoints[w.index]
		w.dc = w.e.newDownloadContext(w.cfg, w.path, ep, w.index == len(w.endpoints)-1, 1)
		d, err := w.e.openDownload(w.ctx, w.dc, byteRange{start: 0, end: -1}, w.modifiedAt)
		if err == nil {
			return w.setDownload(d)
		}
		if !mayWalkDownload(w.ctx, w.dc, err) {
			return err
		}
		w.e.c.logger.InfoContext(w.ctx, "download failed before any data arrived, trying the next endpoint", "endpoint", ep.String(), "error", err)
		if w.dc.countsAgainst(w.ctx, err) {
			w.walkedPast = append(w.walkedPast, ep)
		}
	}
}

// reportDone reports the download to the endpoint health once: at its first
// byte, or at the end of an empty file.
func (w *walkingDownload) reportDone() {
	if !w.reported {
		w.reported = true
		w.e.transferDone(w.ctx, w.health, w.order, w.walkedPast, w.endpoints[:w.index+1])
	}
}

// setDownload makes d the download Read reads, unless Close has run: then it
// closes d, since nothing else would.
func (w *walkingDownload) setDownload(d *resilientDownload) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		_ = d.Close()
		return errDownloadClosed
	}
	w.d = d
	return nil
}

func (w *walkingDownload) isClosed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closed
}

func (w *walkingDownload) Read(p []byte) (int, error) {
	if w.isClosed() {
		return 0, errDownloadClosed
	}
	for {
		n, err := w.d.Read(p)
		if n > 0 {
			w.dc.committed.Store(true)
			w.reportDone()
			return n, nil
		}
		if err == io.EOF {
			w.reportDone()
		}
		if err == nil || err == io.EOF {
			return 0, err
		}
		if w.isClosed() {
			return 0, errDownloadClosed
		}
		if !mayWalkDownload(w.ctx, w.dc, err) {
			return 0, err
		}
		w.e.c.logger.InfoContext(w.ctx, "download failed before any data arrived, trying the next endpoint", "endpoint", w.dc.ep.String(), "error", err)
		if w.dc.countsAgainst(w.ctx, err) {
			w.walkedPast = append(w.walkedPast, w.dc.ep)
		}
		_ = w.d.Close()
		if err := w.openNext(); err != nil {
			return 0, err
		}
	}
}

func (w *walkingDownload) Close() error {
	w.mu.Lock()
	w.closed = true
	d := w.d
	w.mu.Unlock()
	w.stop(errDownloadClosed)
	return d.Close()
}
