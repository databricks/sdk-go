// Package cloudstorage issues unauthenticated, idempotent HTTP requests to cloud
// object storage (S3, Azure Blob, GCS) using short-lived presigned URLs. The
// URLs are self-authenticating, so the client deliberately carries no Databricks
// credentials. It owns the inactivity-stall and exact-length guards and a capped
// retry-with-backoff policy; callers supply the bytes via a [Body] and interpret
// the returned [Response].
package cloudstorage

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/ops"
)

// defaultIdleTimeout cuts a transfer that stops making progress: if the request
// body is not read for this long, no bytes are moving and the attempt is cancelled.
var defaultIdleTimeout = 60 * time.Second

var (
	errStalled   = errors.New("transfer stalled (no progress within the idle timeout)")
	errShortRead = errors.New("body produced fewer bytes than expected (was the source modified mid-transfer?)")
)

// Response is the captured result of a presigned request. Unlike a Databricks
// API call, cloud requests expose the raw status and headers (ETag, Range, 308)
// and are never mapped to a Databricks API error.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// Client issues requests to cloud object storage over presigned URLs. It holds
// only an *http.Client: no Databricks credentials, no config.
type Client struct {
	http        *http.Client
	backoff     ops.BackoffPolicy
	idleTimeout time.Duration
	logger      *slog.Logger
}

// Option configures a [Client].
type Option func(*Client)

// WithLogger sets the logger used for debug-level transfer diagnostics. When
// unset, transfer diagnostics are discarded.
func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithBackoff sets the backoff between the retries of [Client.Send]. The zero
// policy uses the [ops.BackoffPolicy] defaults, as when it is not set.
func WithBackoff(bp ops.BackoffPolicy) Option {
	return func(c *Client) {
		c.backoff = bp
	}
}

// New returns a Client that sends over the given HTTP client. The caller owns
// the client's transport, timeouts, and connection-pool sizing; it must not
// attach Databricks credentials, and should not set a whole-request
// http.Client.Timeout (it would abort a legitimately long transfer; bound the
// operation with the request context instead).
func New(httpClient *http.Client, opts ...Option) *Client {
	c := &Client{
		http:        httpClient,
		idleTimeout: defaultIdleTimeout,
		logger:      slog.New(slog.DiscardHandler),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Attempt performs a single request and returns the raw response. body may be
// nil for a bodyless request (a status query or an abort).
func (c *Client) Attempt(ctx context.Context, method, url string, headers map[string]string, body Body) (*Response, error) {
	var rdr io.Reader
	var bodyLen int64
	if body != nil {
		var err error
		if rdr, err = body.Reader(); err != nil {
			return nil, err
		}
		bodyLen = body.Size()
	}

	// Cancel the attempt on two conditions, each with a distinct cause: the
	// transfer stalls (a dead connection that accepts no more bytes but never
	// errors), or the source comes up short (truncated mid-transfer). Both are
	// told apart from the caller cancelling ctx; the stall is retried, the short
	// read is not.
	attemptCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	req, err := http.NewRequestWithContext(attemptCtx, method, url, rdr)
	if err != nil {
		return nil, RedactURLError(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if bodyLen > 0 && req.Body != nil {
		// A non-bytes source (a file section) is not length-detected by NewRequest,
		// so set Content-Length explicitly. Enforce the exact length and guard
		// against a stall; the body is read as bytes move to the socket, so a gap
		// longer than IdleTimeout means the transfer has stalled. The response phase
		// is bounded by the transport's ResponseHeaderTimeout.
		req.ContentLength = bodyLen
		req.Body = stallGuardedBody(exactBody(req.Body, bodyLen, cancel), c.idleTimeout, cancel)
	}

	// Trace connection reuse so a slow request can be attributed to connection
	// behavior (e.g. HTTP/2 multiplexing onto one connection vs a fresh one).
	var connReused, connWasIdle bool
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			connReused, connWasIdle = info.Reused, info.WasIdle
		},
	}))

	start := time.Now()
	resp, err := c.http.Do(req)
	if err != nil {
		err = RedactURLError(err)
		// Surface our own cancellations as distinct errors; a caller ctx
		// cancellation stays as-is. errShortRead carries its detail message.
		if ctx.Err() == nil {
			switch cause := context.Cause(attemptCtx); {
			case errors.Is(cause, errShortRead):
				err = cause
			case errors.Is(cause, errStalled):
				err = errStalled
			}
		}
		c.logger.DebugContext(ctx, "request failed",
			"method", method, "host", req.URL.Host, "req_bytes", bodyLen,
			"duration_ms", time.Since(start).Milliseconds(), "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}
	c.logger.DebugContext(ctx, "response",
		"method", method, "host", req.URL.Host, "status", resp.StatusCode, "proto", resp.Proto,
		"conn_reused", connReused, "conn_was_idle", connWasIdle, "req_bytes", bodyLen,
		"duration_ms", duration.Milliseconds())
	return &Response{StatusCode: resp.StatusCode, Header: resp.Header, Body: respBody}, nil
}

// Send performs Attempt and retries it up to retries times: transient network
// errors, the retryable statuses in retryStatusCodes (honoring Retry-After), and
// stalled attempts are retried with exponential backoff. The caller sets the
// count, as it depends on where the transfer stands. A fresh request and body reader
// are built per attempt, so no explicit rewind is needed. Exhausting retries on
// a retryable status returns that status as an error.
func (c *Client) Send(ctx context.Context, method, url string, headers map[string]string, body Body, retries int) (*Response, error) {
	var result *Response
	op := func(ctx context.Context) error {
		resp, err := c.Attempt(ctx, method, url, headers, body)
		if err != nil {
			return err // network error or inactivity cut; the retrier decides
		}
		// Surface a failure status as an APIError and let the retrier decide: it
		// retries the retryable statuses (honoring Retry-After) and transient
		// transport failures, and returns everything else. A non-retryable status
		// (403 URL expired, 412 already exists, ...) therefore reaches the caller as
		// this error, which carries the status, headers, and body for inspection.
		// FromHTTPError returns nil for 2xx and 3xx, so both pass through as a
		// successful response (the http.Client already follows real redirects).
		if apiErr := apierr.FromHTTPError(resp.StatusCode, resp.Header, resp.Body); apiErr != nil {
			return apiErr
		}
		result = resp
		return nil
	}
	newRetrier := func() ops.Retrier { return &retrier{maxRetries: retries, backoff: c.backoff} }
	if err := ops.Execute(ctx, op, ops.WithRetrier(newRetrier)); err != nil {
		return nil, err
	}
	return result, nil
}

// RedactURLError returns err, as [http.Client.Do] or [http.NewRequestWithContext]
// returns it, without the query, user info and fragment of any URL in it, since a
// presigned URL's query is a credential. It changes no error in place, because a
// transport may share the errors it returns.
func RedactURLError(err error) error {
	uerr, ok := err.(*url.Error)
	if !ok {
		return redactMessage(err)
	}
	redacted := *uerr
	redacted.URL = redactedURL(uerr.URL)
	redacted.Err = redactMessage(uerr.Err)
	return &redacted
}

// redactMessage returns err, or, when err's message holds a URL with a query,
// user info or fragment, an error that wraps err with those removed from its
// message. Errors can carry a URL in their text alone, as when net/http copies
// a transport error's message on a client timeout.
func redactMessage(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	redacted := urlQuery.ReplaceAllString(urlUserInfo.ReplaceAllString(msg, "$1"), "$1")
	if redacted == msg {
		return err
	}
	return &redactedError{msg: redacted, err: err}
}

// urlUserInfo and urlQuery match the user info, and the query and fragment, of
// an http or https URL in an error message, where %q may have escaped its bytes.
var (
	urlUserInfo = regexp.MustCompile(`(?i)(https?://)(?:\\.|[^\s"\\/?#])*@`)
	urlQuery    = regexp.MustCompile(`(?i)(https?://(?:\\.|[^\s"\\?#])*)[?#](?:\\.|[^\s"\\])*`)
)

// redactedURL returns raw without its query, user info and fragment, or
// "<unparsable URL>" when raw doesn't parse.
func redactedURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<unparsable URL>"
	}
	u.RawQuery, u.ForceQuery, u.Fragment, u.RawFragment, u.User = "", false, "", "", nil
	return u.String()
}

// redactedError has the message of the error it wraps, with its URLs redacted.
type redactedError struct {
	msg string
	err error
}

func (e *redactedError) Error() string { return e.msg }

func (e *redactedError) Unwrap() error { return e.err }

// Timeout and Temporary answer for the wrapped error, as [url.Error]'s do, so
// that a [*url.Error] around a redactedError reports the same timeouts.
func (e *redactedError) Timeout() bool {
	t, ok := e.err.(interface{ Timeout() bool })
	return ok && t.Timeout()
}

func (e *redactedError) Temporary() bool {
	t, ok := e.err.(interface{ Temporary() bool })
	return ok && t.Temporary()
}

// IsURLExpired reports whether err (as returned by Send) is one of the known
// cloud-provider "presigned URL expired" errors: a 403 with the AWS or Azure
// signature-expiry body.
func IsURLExpired(err error) bool {
	aerr, ok := errors.AsType[*apierr.APIError](err)
	if !ok || aerr.HTTPStatusCode() != http.StatusForbidden {
		return false
	}
	var e struct {
		XMLName    xml.Name `xml:"Error"`
		Code       string   `xml:"Code"`
		Message    string   `xml:"Message"`
		AuthDetail string   `xml:"AuthenticationErrorDetail"`
	}
	if uerr := xml.Unmarshal(aerr.HTTPBody(), &e); uerr != nil {
		return false
	}
	switch e.Code {
	case "AuthenticationFailed": // Azure
		return strings.Contains(e.AuthDetail, "Signature not valid in the specified time frame")
	case "AccessDenied": // AWS
		return e.Message == "Request has expired"
	default:
		return false
	}
}

// stallReader wraps a request body so a transfer that stops making progress is
// cancelled. The transport reads the body as it writes bytes to the socket, so a
// gap longer than idle between reads means no bytes are moving -- a stalled or
// dead connection. A transfer making any progress keeps resetting the timer and
// so cannot be cut, unlike a whole-attempt or throughput-floor deadline.
type stallReader struct {
	r      io.Reader
	idle   time.Duration
	cancel context.CancelCauseFunc

	// mu guards timer and stopped: the transport may close the body, which stops
	// the timer, while another goroutine is in Read.
	mu      sync.Mutex
	timer   *time.Timer
	stopped bool
}

func newStallReader(r io.Reader, idle time.Duration, cancel context.CancelCauseFunc) *stallReader {
	return &stallReader{r: r, idle: idle, cancel: cancel}
}

func (s *stallReader) Read(p []byte) (int, error) {
	// Arm on the first read (after the connection is established, so dial time is
	// not charged against the idle budget) and reset on each later read.
	s.arm()
	n, err := s.r.Read(p)
	if err != nil {
		// EOF (transfer done) or a read error: stop arming so the response phase is
		// not subject to the stall timer.
		s.stop()
	}
	return n, err
}

// arm starts or resets the idle timer, unless stop has run.
func (s *stallReader) arm() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return
	}
	if s.timer == nil {
		s.timer = time.AfterFunc(s.idle, func() { s.cancel(errStalled) })
	} else {
		s.timer.Reset(s.idle)
	}
}

// stop disarms the idle timer for good, so a later Read doesn't arm it again.
// Safe to call before the first Read and more than once.
func (s *stallReader) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = true
	if s.timer != nil {
		s.timer.Stop()
	}
}

// stallGuardedBody wraps a request body's reads with stall detection while
// preserving its Close.
func stallGuardedBody(body io.ReadCloser, idle time.Duration, cancel context.CancelCauseFunc) io.ReadCloser {
	return &stallGuardedReadCloser{Closer: body, sr: newStallReader(body, idle, cancel)}
}

type stallGuardedReadCloser struct {
	io.Closer
	sr *stallReader
}

func (b *stallGuardedReadCloser) Read(p []byte) (int, error) {
	return b.sr.Read(p)
}

// Close disarms the stall timer. The transport closes the request body once it
// has finished writing it, which on a fixed-length request happens without a
// trailing EOF read (the writer stops at Content-Length). Without stopping here,
// the timer would stay armed into the response phase and could cancel a
// successful upload as errStalled while waiting on a slow response.
func (b *stallGuardedReadCloser) Close() error {
	b.sr.stop()
	return b.Closer.Close()
}

// exactReader enforces that a body delivers exactly the declared length. If the
// source ends early -- e.g. the file was truncated during the transfer -- it
// cancels the attempt with errShortRead rather than letting a short body be sent
// (silent corruption) or surfacing net/http's opaque ContentLength-mismatch
// error (which would then be pointlessly retried).
type exactReader struct {
	r      io.Reader
	left   int64
	cancel context.CancelCauseFunc
}

func (e *exactReader) Read(p []byte) (int, error) {
	n, err := e.r.Read(p)
	e.left -= int64(n)
	if err == io.EOF && e.left > 0 {
		cause := fmt.Errorf("%w: stream ended %d bytes short of the expected length", errShortRead, e.left)
		e.cancel(cause)
		return n, cause
	}
	return n, err
}

// exactBody wraps a request body's reads with the exact-length check while
// preserving its Close.
func exactBody(body io.ReadCloser, n int64, cancel context.CancelCauseFunc) io.ReadCloser {
	return &exactReadCloser{Closer: body, er: &exactReader{r: body, left: n, cancel: cancel}}
}

type exactReadCloser struct {
	io.Closer
	er *exactReader
}

func (b *exactReadCloser) Read(p []byte) (int, error) {
	return b.er.Read(p)
}
