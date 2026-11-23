package files

// The large-file upload engine's policy knobs and the engine value that carries
// them. tunables is injected at construction (defaultTunables in production,
// custom values in tests) rather than read from package globals, so a test never
// mutates shared state and uploads with different policies cannot interfere.

import (
	"context"
	"os"
	"time"

	"github.com/databricks/sdk-go/core/ops"
)

// tunables holds every policy value the engine reads. It is copied by value into
// the engine, the per-upload context, and the straggler guard, so nothing shared
// is mutated after construction.
type tunables struct {
	// minStreamSize is the threshold below which a known-size upload is sent in a
	// single PUT instead of being split into parts.
	minStreamSize int64

	// defaultPartSize is the part size used when no part size is requested. A
	// fixed, modestly sized part (rather than one scaled to keep the part count
	// low) bounds the memory a buffering upload holds (about parallelism*partSize
	// for a non-seekable stream) and spreads stragglers across more parts. It is
	// only grown when a file is large enough to exceed maxUploadParts.
	defaultPartSize int64

	// maxPartSize is the largest part size a cloud provider accepts.
	maxPartSize int64

	// batchURLCount is the number of presigned URLs requested per batch when the
	// content length is unknown.
	batchURLCount int

	// fileParallelism is the default worker count when no parallelism is requested
	// and the source is randomly readable (a local file or other io.ReaderAt).
	// Each part streams from a positioned read, so resident memory stays bounded
	// to the transport's per-connection buffers regardless of the worker count
	// (not parallelism*partSize); a high default saturates a fast uplink, and the
	// driver caps the workers at the part count for small files.
	fileParallelism int

	// streamMemoryBudget bounds the bytes a buffered upload (a non-seekable
	// stream) holds in memory, where each worker keeps one partSize buffer for its
	// in-flight part. The default stream worker count is this budget divided by
	// the part size, so a larger part size lowers the worker count instead of
	// growing memory. At the 16 MiB default part size this is 16 workers. A
	// stream of unknown size also holds up to minStreamSize of read-ahead until
	// its first parts are cut from it.
	streamMemoryBudget int64

	// maxRetries caps how many times a single part is retried after a presigned
	// URL expires.
	maxRetries int

	// cloudRetries caps how many times a part PUT is retried once a part of the
	// upload has landed. cloudRetriesBeforeCommit applies before then: the
	// workspace is still left to fall back to, so failing over is cheaper than
	// waiting out retries on a storage host that may be blocked.
	cloudRetries             int
	cloudRetriesBeforeCommit int

	// maxUploadParts is the provider's maximum number of parts per multipart
	// upload (S3's 10,000 is the common floor). The default part size is grown for
	// files large enough to otherwise exceed it.
	maxUploadParts int64

	// responseHeaderTimeout bounds the response phase of a cloud transfer so a
	// connection that accepts the upload but never replies cannot hang a goroutine
	// forever. It is applied to the transfer client built in resolveCloudClient.
	responseHeaderTimeout time.Duration

	// controlAttemptTimeout bounds one attempt of a control-plane call
	// (initiate, complete, create-upload-part-urls, create-download-url, abort),
	// so that a host that accepts the connection and never answers gives way to
	// the next endpoint instead of holding the whole retry budget. It applies
	// only when the client's HTTP client sets no Timeout, which would already
	// bound each attempt.
	controlAttemptTimeout time.Duration

	// cleanupTimeout bounds a best-effort abort, which runs on a context detached
	// from the upload's (see cleanupContext).
	cleanupTimeout time.Duration

	// slowFactor cancels a part attempt that exceeds this multiple of the recent
	// p95 part duration. It is the main straggler knob: lower catches stragglers
	// sooner but re-issues more healthy parts.
	slowFactor int

	// slowWarmup is the number of completed parts required before the guard
	// switches from the cold-start deadline to the p95-relative one.
	slowWarmup int

	// slowWindow bounds the rolling sample set so the p95 tracks recent conditions
	// rather than the whole upload.
	slowWindow int

	// slowMinDeadline floors the soft deadline so a fast network (tiny p95) cannot
	// produce a deadline that clips normal variance.
	slowMinDeadline time.Duration

	// slowColdDeadline is the absolute soft deadline for a part in the opening
	// wave, before slowWarmup parts have completed and the p95-relative deadline is
	// armed. Because the deadline is re-evaluated in flight (see sendPart), this
	// only bites in the first seconds of an upload or for a file too small to warm
	// the guard; once warmed, an in-flight early part is caught at the tighter p95
	// deadline instead. It is set above the legitimate opening-burst latency (tens
	// of seconds at high concurrency).
	slowColdDeadline time.Duration

	// slowFirstPartDeadline is the soft deadline for the synchronous first part,
	// which runs alone before any worker (so it has no samples and, unlike the
	// opening wave, no burst contention). A healthy first PUT is a few seconds, so
	// a tight deadline re-issues a wedged first connection quickly without blocking
	// the whole upload at 0% on the cold-start floor.
	slowFirstPartDeadline time.Duration

	// slowCheckInterval is how often an in-flight attempt is re-checked against the
	// current deadline, so warmup and a falling p95 take effect while the attempt
	// is still running.
	slowCheckInterval time.Duration

	// slowMaxReissue caps soft re-issues per part; afterward the part rides out the
	// normal timeouts rather than being cancelled again, so it is never failed
	// solely for staying slow.
	slowMaxReissue int

	// cpBackoff is the exponential-backoff-with-jitter policy for control-plane
	// retries.
	cpBackoff ops.BackoffPolicy

	// cloudBackoff is the backoff between the retries of a cloud-storage
	// request. The zero policy uses the ops defaults.
	cloudBackoff ops.BackoffPolicy

	// cpRetryTimeout bounds the total time spent retrying a single control-plane
	// call. Without it, a persistently failing endpoint would be retried until the
	// caller's context expires (and forever if it has no deadline). For the
	// single-shot PUT, which may carry gigabytes, it only stops new attempts and
	// bounds the credential resolution before one; it never cuts a running
	// transfer short.
	cpRetryTimeout time.Duration

	// A download to a file larger than both downloadMinParallelSize and
	// downloadPartSize is split into ranges of downloadPartSize, written on up
	// to downloadParallelism goroutines by default.
	downloadMinParallelSize int64
	downloadPartSize        int64
	downloadParallelism     int

	// transferIdleTimeout breaks a transfer whose connection moves no bytes for
	// this long: a download resumes on a fresh connection, and a single-shot
	// upload's attempt fails as a timeout, to be retried or walked on. It never
	// bounds a transfer whose bytes keep moving.
	transferIdleTimeout time.Duration

	// downloadResumesWithoutProgress caps how often a download resumes in a row
	// without receiving a byte.
	downloadResumesWithoutProgress int

	// endpointsAPI is whether a transfer asks the workspace which endpoints to
	// use (GET /api/2.0/fs/endpoints); without it a transfer walks the built-in
	// order. Client calls turn it on, so that the workspace chooses the hosts,
	// unless endpointsAPIEnv is false. Engines built in tests leave it off.
	endpointsAPI bool

	// urlExpiry is how long requested presigned URLs are valid.
	urlExpiry time.Duration

	// After endpointSkipAfter transfers in a row that failed on an endpoint
	// and then completed on a later one, a client leaves that endpoint out of
	// new transfers for endpointSkipCoolOff (an upload only it can take still
	// uses it), then lets one transfer probe it. Each failed probe doubles the
	// cool-off, up to endpointSkipMaxCoolOff; a success clears it. 0 turns
	// skipping off.
	endpointSkipAfter      int
	endpointSkipCoolOff    time.Duration
	endpointSkipMaxCoolOff time.Duration
}

// defaultTunables returns the production policy values.
func defaultTunables() tunables {
	return tunables{
		minStreamSize:      50 * 1024 * 1024,
		defaultPartSize:    16 << 20, // 16 MiB
		maxPartSize:        4 << 30,
		batchURLCount:      1,
		fileParallelism:    128,
		streamMemoryBudget: 256 << 20,
		maxRetries:         3,
		// A part PUT gets 3 attempts once a part has landed, and 2 before.
		cloudRetries:             2,
		cloudRetriesBeforeCommit: 1,
		maxUploadParts:           10_000,
		responseHeaderTimeout:    60 * time.Second,
		controlAttemptTimeout:    60 * time.Second,
		cleanupTimeout:           30 * time.Second,
		slowFactor:               3,
		slowWarmup:               32,
		slowWindow:               128,
		slowMinDeadline:          5 * time.Second,
		slowColdDeadline:         30 * time.Second,
		slowFirstPartDeadline:    10 * time.Second,
		slowCheckInterval:        1 * time.Second,
		slowMaxReissue:           2,
		cpBackoff:                ops.BackoffPolicy{},
		cpRetryTimeout:           5 * time.Minute,
		urlExpiry:                time.Hour,
		// Downloads.
		downloadMinParallelSize:        50 << 20,
		downloadPartSize:               10 << 20,
		downloadParallelism:            10,
		transferIdleTimeout:            60 * time.Second,
		downloadResumesWithoutProgress: 1,
		// Endpoint skipping.
		endpointSkipAfter:      3,
		endpointSkipCoolOff:    2 * time.Minute,
		endpointSkipMaxCoolOff: 30 * time.Minute,
	}
}

// engine performs a single or shared-limiter set of uploads for a Client. It
// pairs the Client (for its authenticated transport, host, logger, and workspace
// routing) with the tunables that govern part sizing, retries, and straggler
// control. All internal upload logic hangs off *engine so tests can drive it
// with custom tunables without touching package state.
type engine struct {
	c   *Client
	tun tunables
	// now is the clock that endpoint-list expiry reads; nil is time.Now.
	now func() time.Time
}

// clock returns the current time.
func (e *engine) clock() time.Time {
	if e.now != nil {
		return e.now()
	}
	return time.Now()
}

// newEngine builds an engine over c with the production tunables.
func newEngine(c *Client) *engine {
	return &engine{c: c, tun: defaultTunables()}
}

// clientTunables returns the tunables a Client method's engine starts from:
// defaultTunables. Tests replace it to shrink the limits, so that a transfer
// through the Client methods needs neither large files nor long retries.
var clientTunables = defaultTunables

// newClientEngine builds the engine for a Client method. Unlike newEngine, it
// reads endpointsAPIEnv, so that engines built in tests do not depend on the
// environment.
func newClientEngine(ctx context.Context, c *Client) *engine {
	e := &engine{c: c, tun: clientTunables()}
	e.tun.endpointsAPI = endpointsAPIEnabled(ctx, c.logger, os.Getenv(endpointsAPIEnv))
	return e
}
