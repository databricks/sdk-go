package files

// Large-file upload engine for Databricks Unity Catalog Volumes. It chooses a
// single-shot or multipart protocol based on the stream size, and walks the
// upload endpoints in order (ext_endpoints.go) until one takes the upload.
//
// The Files API calls that run an upload (ext_upload_control.go) carry the
// client's credentials. Parts go either to presigned cloud-storage URLs, without
// credentials, through the cloudstorage subpackage, or to a Files API host with
// them.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/databricks/sdk-go/files/v2/internal/cloudstorage"
)

// errUploadURLExpired is returned when a part's presigned URL keeps expiring
// across re-mint retries. An already-exists condition is surfaced as the
// [ErrAlreadyExists] sentinel; a malformed server response is an internal error.
var errUploadURLExpired = errors.New("upload URL expired after retries")

// UploadResult holds the result of an upload. It is currently empty and exists
// for forward compatibility.
type UploadResult struct{}

// UploadOption configures an Upload or UploadFrom call.
type UploadOption func(*uploadConfig)

type uploadConfig struct {
	overwrite      *bool
	partSize       int64
	parallelism    *int
	progress       ProgressFunc
	transferClient *http.Client
	limiter        Limiter
}

// Progress reports the state of an in-flight upload. Fields may be added in
// future releases, so callers must not depend on the struct being comparable or
// on its exact size; always refer to fields by name.
type Progress struct {
	// Transferred is the cumulative number of bytes confirmed uploaded so far.
	Transferred int64
	// Total is the total size in bytes, or -1 if it is not known in advance (a
	// non-seekable stream).
	Total int64
}

// ProgressFunc is invoked as an upload makes progress. It is called from
// internal goroutines but never concurrently with itself, so it needs no
// locking of its own; it must return promptly.
type ProgressFunc func(Progress)

// WithOverwrite controls whether an existing file is overwritten. When this
// option is not supplied the parameter is omitted and the server applies its
// default.
func WithOverwrite(overwrite bool) UploadOption {
	return func(c *uploadConfig) { c.overwrite = &overwrite }
}

// WithPartSize sets the multipart part size in bytes. It must not exceed the
// cloud provider maximum. When not supplied an appropriate size is chosen from
// the content length.
func WithPartSize(partSize int64) UploadOption {
	return func(c *uploadConfig) { c.partSize = partSize }
}

// WithParallelism sets the number of concurrent upload workers used for a large
// file. A value of 1 uploads the parts sequentially on a single goroutine;
// higher values upload that many parts at once. It must be at least 1. When not
// set, a default is used.
func WithParallelism(n int) UploadOption {
	return func(c *uploadConfig) { c.parallelism = &n }
}

// WithProgress registers a callback invoked as the upload progresses, reporting
// the cumulative bytes uploaded and the total size (-1 if unknown). It is useful
// for rendering an upload progress bar.
func WithProgress(fn ProgressFunc) UploadOption {
	return func(c *uploadConfig) { c.progress = fn }
}

// WithTransferClient overrides the HTTP client that sends parts to presigned
// cloud-storage URLs, replacing the one the call would otherwise create
// internally. Use it to route those transfers through a custom transport, proxy,
// or CA. Requests to Databricks hosts, including parts sent to a Files API host,
// use the client's own transport.
//
// Because presigned URLs authenticate themselves, the client must not attach
// Databricks credentials; the storage provider rejects extra authentication. Do
// not set http.Client.Timeout, a whole-request deadline that would abort a
// legitimately long part transfer; bound the upload with the context passed to
// Upload instead.
func WithTransferClient(client *http.Client) UploadOption {
	return func(c *uploadConfig) { c.transferClient = client }
}

func resolveUploadConfig(opts []UploadOption, tun tunables) (*uploadConfig, error) {
	cfg := &uploadConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.parallelism != nil && *cfg.parallelism < 1 {
		return nil, fmt.Errorf("parallelism must be at least 1, got %d", *cfg.parallelism)
	}
	if cfg.partSize < 0 {
		return nil, fmt.Errorf("part size must be non-negative, got %d", cfg.partSize)
	}
	if cfg.partSize > tun.maxPartSize {
		return nil, fmt.Errorf("part size %d exceeds maximum %d", cfg.partSize, tun.maxPartSize)
	}
	return cfg, nil
}

// resolvedParallelism returns the worker count: the caller's explicit value when
// set, otherwise a default that depends on whether the source is buffered. A
// randomly-readable source streams parts from positioned reads and uses a high,
// bandwidth-oriented default; a buffered (non-seekable) source holds
// parallelism*partSize in memory, so its default is sized to a memory budget,
// capped at the file default so a tiny part size cannot translate the budget into
// an unbounded number of goroutines and sockets.
func (cfg *uploadConfig) resolvedParallelism(buffered bool, partSize int64, tun tunables) int {
	if cfg.parallelism != nil {
		return *cfg.parallelism
	}
	if !buffered {
		return tun.fileParallelism
	}
	workers := int(tun.streamMemoryBudget / partSize)
	return min(max(workers, 1), tun.fileParallelism)
}

// newTransferClient returns the default HTTP client for cloud-leg transfers,
// sized for n concurrent transfers: it keeps up to n idle connections to the
// storage host so they are reused rather than re-dialed (Go's default of 2 per
// host would otherwise force re-dialing). It attaches no Databricks credentials
// (presigned URLs are self-authenticating) and sets no whole-request timeout,
// which would abort a legitimately long transfer; the upload is bounded by the
// context instead. Callers who need to supply their own client use
// WithTransferClient.
func newTransferClient(n int, tun tunables) *http.Client {
	return &http.Client{Transport: sizedForTransfers(http.DefaultTransport.(*http.Transport), n, tun)}
}

// sizedForTransfers returns a copy of t that bounds the response phase of a
// transfer and keeps up to n idle connections to a host.
func sizedForTransfers(t *http.Transport, n int, tun tunables) *http.Transport {
	t = t.Clone()
	t.ResponseHeaderTimeout = tun.responseHeaderTimeout
	t.MaxIdleConnsPerHost = n
	t.MaxIdleConns = max(t.MaxIdleConns, n)
	return t
}

// resolveCloudClient returns the caller-supplied transfer client when set, and
// otherwise a fresh one sized to the upload parallelism.
func (cfg *uploadConfig) resolveCloudClient(parallelism int, tun tunables) *http.Client {
	if cfg.transferClient != nil {
		return cfg.transferClient
	}
	return newTransferClient(parallelism, tun)
}

// uploadContext is the resolved per-upload state. Its fields are set once at
// construction; the slowGuard and completed flag hold the only mutable state,
// both safe for concurrent use by the upload workers.
type uploadContext struct {
	tun           tunables
	targetPath    string
	overwrite     *bool
	partSize      int64 // of the current endpoint: basePartSize, capped by the endpoint
	basePartSize  int64
	batchSize     int
	contentLength int64 // -1 when unknown (non-seekable stream)
	parallel      bool
	parallelism   int
	cloud         *cloudstorage.Client // presigned URLs: no credentials
	// files sends parts to a Files API host with the client's credentials. It is
	// built only if the walk reaches such an endpoint.
	files     func() *cloudstorage.Client
	progress  *progressReporter
	slowGuard *slowAttemptGuard
	limiter   Limiter // bounds concurrent transfers; never nil (unlimited by default)

	// data is the upload's source; the walk sends it to each endpoint it tries.
	data *uploadData

	// ep is the endpoint the walk is trying, and lastEndpoint whether no
	// endpoint is left after it. They change only between endpoints, while no
	// request is in flight.
	ep           endpoint
	lastEndpoint bool

	// completed records whether any part has landed on the current endpoint.
	// The walk moves on only while none has; once one has, the upload is
	// committed to the endpoint and a later failure is returned.
	completed atomic.Bool
}

// startEndpoint sets the endpoint the walk tries next.
func (uc *uploadContext) startEndpoint(ep endpoint, last bool) {
	uc.ep, uc.lastEndpoint = ep, last
	uc.partSize = uc.basePartSize
	if limit := ep.maxPartSize(); limit >= 0 {
		uc.partSize = min(uc.partSize, limit)
	}
}

// commit records that a part landed, the first time on any goroutine.
func (uc *uploadContext) commit() {
	if !uc.completed.Swap(true) && uc.data != nil {
		uc.data.commit()
	}
}

func (uc *uploadContext) committed() bool { return uc.completed.Load() }

// retries is how often a request to the current endpoint may be retried, or -1
// for no cap but the retry timeout. Before any part lands, with an endpoint left
// to move on to, a failing request gets a single retry: moving on is cheaper than
// waiting out a host that may be unreachable from here. Otherwise the workspace
// host keeps the API client's retries, and other hosts get cloudRetries.
func (uc *uploadContext) retries(onWorkspace bool) int {
	switch {
	case !uc.completed.Load() && !uc.lastEndpoint:
		return uc.tun.cloudRetriesBeforeCommit
	case onWorkspace:
		return -1
	}
	return uc.tun.cloudRetries
}

// partRetries is how often a part PUT is retried. Parts never go to the
// workspace host.
func (uc *uploadContext) partRetries() int { return uc.retries(false) }

// partClient sends the current endpoint's part PUTs.
func (uc *uploadContext) partClient() *cloudstorage.Client {
	if uc.ep.kind == presignedEndpoint {
		return uc.cloud
	}
	return uc.files()
}

// newUploadContext resolves the per-upload state shared by Upload and UploadFrom.
// contentLength is -1 for a non-seekable stream. buffered reports whether the
// upload holds parts in memory (a non-seekable stream) rather than streaming
// them from positioned reads, which sets the default worker count.
func (e *engine) newUploadContext(cfg *uploadConfig, targetPath string, data *uploadData, buffered bool) *uploadContext {
	contentLength := data.size()
	partSize, batchSize := optimizeParams(contentLength, cfg.partSize, e.tun)
	parallelism := cfg.resolvedParallelism(buffered, partSize, e.tun)
	limiter := cfg.limiter
	if limiter == nil {
		limiter = unlimitedLimiter{}
	}
	return &uploadContext{
		tun:           e.tun,
		targetPath:    targetPath,
		overwrite:     cfg.overwrite,
		partSize:      partSize,
		basePartSize:  partSize,
		batchSize:     batchSize,
		contentLength: contentLength,
		data:          data,
		parallel:      parallelism > 1,
		parallelism:   parallelism,
		cloud:         cloudstorage.New(cfg.resolveCloudClient(parallelism, e.tun), cloudstorage.WithLogger(e.c.logger), cloudstorage.WithBackoff(e.tun.cloudBackoff)),
		files: sync.OnceValue(func() *cloudstorage.Client {
			return cloudstorage.New(e.authenticatedHTTPClient(parallelism), cloudstorage.WithLogger(e.c.logger), cloudstorage.WithBackoff(e.tun.cloudBackoff))
		}),
		progress:  newProgressReporter(cfg.progress, contentLength),
		slowGuard: &slowAttemptGuard{tun: e.tun},
		limiter:   limiter,
	}
}

// progressReporter accumulates confirmed-uploaded bytes and forwards the running
// total to a ProgressFunc. Its methods are safe for concurrent use (parallel
// uploads report from many goroutines) and serialize callbacks so the
// ProgressFunc never runs concurrently with itself.
type progressReporter struct {
	fn    ProgressFunc
	total int64

	mu          sync.Mutex
	transferred int64
	// attempt numbers the single-shot attempts; only the latest one's reads
	// count (see restart).
	attempt int64
}

// newProgressReporter returns nil when fn is nil so that reporter calls are
// no-ops without the callers having to check.
func newProgressReporter(fn ProgressFunc, total int64) *progressReporter {
	if fn == nil {
		return nil
	}
	return &progressReporter{fn: fn, total: total}
}

// add records n newly uploaded bytes and reports the new cumulative total. It is
// a no-op on a nil reporter, so callers need not check.
func (p *progressReporter) add(n int64) {
	if p == nil || n <= 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.transferred += n
	p.fn(Progress{Transferred: p.transferred, Total: p.total})
}

// restart begins a pass that sends the whole upload from its start, a
// single-shot attempt or the next endpoint of the walk: the count rewinds to
// zero, so Transferred never exceeds Total. It returns the attempt's number
// for addAttempt; reads of earlier attempts, on this endpoint or an earlier
// one, no longer count.
func (p *progressReporter) restart() int64 {
	if p == nil {
		return 0
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.attempt++
	if p.transferred != 0 {
		p.transferred = 0
		p.fn(Progress{Transferred: 0, Total: p.total})
	}
	return p.attempt
}

// addAttempt records n bytes that single-shot attempt read, unless a later
// attempt has begun. It is a no-op on a nil reporter.
func (p *progressReporter) addAttempt(attempt, n int64) {
	if p == nil || n <= 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if attempt != p.attempt {
		return
	}
	p.transferred += n
	p.fn(Progress{Transferred: p.transferred, Total: p.total})
}

// Upload uploads contents to the remote file at filePath, which must be an
// absolute path such as "/Volumes/catalog/schema/volume/file". It chooses a
// single-shot or multipart upload based on the size, which it learns by seeking
// a seekable reader, or by reading any other reader up to the multipart
// threshold.
//
// When contents implements io.ReaderAt (for example an *os.File or a
// *bytes.Reader) and more than one part is in flight, parts are read from it
// with concurrent positioned reads rather than buffered in memory. Otherwise
// each in-flight part is buffered as it is read.
//
// ctx bounds the entire operation, including all retries; pass a context with a
// deadline to cap the total upload time.
func (c *Client) Upload(ctx context.Context, filePath string, contents io.Reader, opts ...UploadOption) (*UploadResult, error) {
	return newClientEngine(ctx, c).upload(ctx, filePath, contents, opts...)
}

// UploadFrom uploads the local file at sourcePath to the remote file at
// filePath. The local file is opened as needed; its size is always known, so it
// never takes the non-seekable path.
//
// A multipart upload covers the file as of its size when the call begins:
// bytes appended afterward are not included, and truncating the file mid-upload
// fails the upload (it does not produce a partial object). A single-shot
// upload, including the fallback after a multipart upload fails before any part
// lands, sends the file as it is when that request starts. Do not modify the
// file while an upload is in progress.
//
// ctx bounds the entire operation, including all retries; pass a context with a
// deadline to cap the total upload time.
func (c *Client) UploadFrom(ctx context.Context, filePath, sourcePath string, opts ...UploadOption) (*UploadResult, error) {
	return newClientEngine(ctx, c).uploadFrom(ctx, filePath, sourcePath, opts...)
}

func (e *engine) upload(ctx context.Context, filePath string, contents io.Reader, opts ...UploadOption) (*UploadResult, error) {
	cfg, err := resolveUploadConfig(opts, e.tun)
	if err != nil {
		return nil, err
	}

	// A seekable reader lets us learn the size up front; otherwise it is unknown.
	// Seeking to the end moves the reader, so a failed rewind afterward is
	// unrecoverable: the reader is left at EOF and a later read would send zero
	// bytes. Fail loud in that case rather than silently truncating the object.
	contentLength := int64(-1)
	if seeker, ok := contents.(io.Seeker); ok {
		if end, serr := seeker.Seek(0, io.SeekEnd); serr == nil {
			if _, serr := seeker.Seek(0, io.SeekStart); serr != nil {
				return nil, fmt.Errorf("cannot rewind reader after sizing: %w", serr)
			}
			contentLength = end
		}
	}

	// A randomly-readable source large enough for the parallel path streams its
	// parts from positioned reads; any other source buffers each in-flight part.
	_, isReaderAt := contents.(io.ReaderAt)
	buffered := !isReaderAt || contentLength < e.tun.minStreamSize
	var data *uploadData
	if contentLength >= 0 {
		data = newSeekableData(contents.(io.ReadSeeker), contentLength)
	} else {
		data = newStreamData(contents)
	}
	uc := e.newUploadContext(cfg, filePath, data, buffered)
	return &UploadResult{}, e.walkUpload(ctx, uc, data)
}

func (e *engine) uploadFrom(ctx context.Context, filePath, sourcePath string, opts ...UploadOption) (*UploadResult, error) {
	cfg, err := resolveUploadConfig(opts, e.tun)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(sourcePath)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// A local file is randomly readable, so the parallel path streams parts from
	// positioned reads and never buffers them: buffered is false.
	data := newSeekableData(f, info.Size())
	uc := e.newUploadContext(cfg, filePath, data, false)
	return &UploadResult{}, e.walkUpload(ctx, uc, data)
}

// optimizeParams picks a part size and a presigned-URL batch size. contentLength
// is -1 when unknown; override is 0 when no part size was requested. The override
// is already range-checked in resolveUploadConfig.
func optimizeParams(contentLength, override int64, tun tunables) (partSize int64, batchSize int) {
	switch {
	case override > 0:
		partSize = override
	default:
		partSize = tun.defaultPartSize
		// Grow the part size only if the fixed default would exceed the
		// provider's maximum part count for this (known) content length.
		if contentLength >= 0 {
			if minPart := (contentLength + tun.maxUploadParts - 1) / tun.maxUploadParts; minPart > partSize {
				partSize = min(minPart, tun.maxPartSize)
			}
		}
	}

	batchSize = tun.batchURLCount
	if contentLength >= 0 && partSize > 0 {
		numParts := (contentLength + partSize - 1) / partSize
		// The square root of the part count is a heuristic that balances the
		// number of URL-minting round trips against batch size.
		batchSize = max(int(math.Ceil(math.Sqrt(float64(numParts)))), 1)
	}
	return partSize, batchSize
}

// initiateMultipart opens a multipart upload session on the current endpoint
// and returns its token. Only multipart sessions are supported; any other
// initiate response is an unexpected one.
func (e *engine) initiateMultipart(ctx context.Context, uc *uploadContext) (string, error) {
	resp, err := e.initiateUpload(ctx, uc)
	if err != nil {
		return "", err
	}
	if uc.ep.kind == filesAPIv3Endpoint {
		if resp.UploadID == "" {
			return "", fmt.Errorf("%w: missing upload_id", errUnexpectedServerResponse)
		}
		return resp.UploadID, nil
	}
	if resp.MultipartUpload == nil {
		return "", fmt.Errorf("%w: no multipart upload session offered", errUnexpectedServerResponse)
	}
	if resp.MultipartUpload.SessionToken == "" {
		return "", fmt.Errorf("%w: missing multipart session token", errUnexpectedServerResponse)
	}
	return resp.MultipartUpload.SessionToken, nil
}
