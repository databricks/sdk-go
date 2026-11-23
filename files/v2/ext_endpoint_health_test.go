package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/ops"
)

// ======================== Endpoint health ========================

var (
	healthA = endpoint{kind: presignedEndpoint, baseURL: "https://a.example.com", multipart: true}
	healthB = endpoint{kind: filesAPIEndpoint, baseURL: "https://b.example.com"}
)

// testClock is a clock the test moves.
type testClock struct {
	mu  sync.Mutex
	now time.Time
}

func newTestClock() *testClock { return &testClock{now: time.Unix(1_700_000_000, 0)} }

func (c *testClock) time() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *testClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// recordedLogs keeps every log record, for tests that check what was logged.
type recordedLogs struct {
	mu      sync.Mutex
	records []slog.Record
}

func (l *recordedLogs) Enabled(context.Context, slog.Level) bool { return true }
func (l *recordedLogs) Handle(_ context.Context, r slog.Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, r.Clone())
	return nil
}
func (l *recordedLogs) WithAttrs([]slog.Attr) slog.Handler { return l }
func (l *recordedLogs) WithGroup(string) slog.Handler      { return l }

// levels returns the level of each record whose message contains msg.
func (l *recordedLogs) levels(msg string) []slog.Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []slog.Level
	for _, r := range l.records {
		if strings.Contains(r.Message, msg) {
			out = append(out, r.Level)
		}
	}
	return out
}

func testHealth(skipAfter int, logs *recordedLogs) *endpointHealth {
	var h slog.Handler = slog.NewTextHandler(io.Discard, nil)
	if logs != nil {
		h = logs
	}
	return newEndpointHealth(tunables{endpointSkipAfter: skipAfter, endpointSkipCoolOff: 10 * time.Second, endpointSkipMaxCoolOff: 25 * time.Second}, slog.New(h))
}

func eqEndpoints(t *testing.T, what string, got []endpoint, want ...endpoint) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}

func probed(order endpointOrder) []endpoint {
	var out []endpoint
	for ep := range order.probes {
		out = append(out, ep)
	}
	return out
}

func TestEndpointHealth_skipsAfterFailuresInARowThenProbes(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(2, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	eqEndpoints(t, "after one failure", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	eqEndpoints(t, "after two failures", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)

	clock.advance(10 * time.Second)
	eqEndpoints(t, "the probe", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
	// Only one transfer probes; the others keep skipping until it reports.
	eqEndpoints(t, "while the probe is out", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)
	h.recordSuccess(ctx, healthA)
	eqEndpoints(t, "after a success", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
}

func TestEndpointHealth_doublesTheCoolOffAfterAFailedProbeUpToTheMaximum(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	var skippedFor []time.Duration
	for range 3 {
		start := clock.time()
		order := h.order(clock.time(), []endpoint{healthA, healthB})
		for len(order.endpoints) == 1 {
			clock.advance(time.Second)
			order = h.order(clock.time(), []endpoint{healthA, healthB})
		}
		skippedFor = append(skippedFor, clock.time().Sub(start))
		h.recordFailure(ctx, clock.time(), healthA, order.probes[healthA])
	}
	if want := []time.Duration{10 * time.Second, 20 * time.Second, 25 * time.Second}; !slices.Equal(skippedFor, want) {
		t.Errorf("skipped for %v, want %v", skippedFor, want)
	}
}

func TestEndpointHealth_neverSkipsTheLastEndpoint(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	h.recordFailure(ctx, clock.time(), healthB, 0)
	eqEndpoints(t, "two endpoints", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)
	eqEndpoints(t, "one endpoint", h.order(clock.time(), []endpoint{healthA}).endpoints, healthA)
}

func TestEndpointHealth_handsAReturnedProbeToTheNextTransfer(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	clock.advance(10 * time.Second)

	first := h.order(clock.time(), []endpoint{healthA, healthB})
	eqEndpoints(t, "the first transfer's endpoints", first.endpoints, healthA, healthB)
	eqEndpoints(t, "the first transfer's probes", probed(first), healthA)
	// The transfer never tried a, for example because it planned no request
	// there.
	h.returnProbe(clock.time(), healthA, first.probes[healthA])
	second := h.order(clock.time(), []endpoint{healthA, healthB})
	eqEndpoints(t, "the second transfer's endpoints", second.endpoints, healthA, healthB)
	eqEndpoints(t, "the second transfer's probes", probed(second), healthA)
	// Without the probe back, a third transfer keeps skipping a until the
	// second reports.
	eqEndpoints(t, "the third transfer's endpoints", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)
}

func TestEndpointHealth_ignoresAFailedProbeReportedAfterItWasIssuedAgain(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	clock.advance(10 * time.Second)
	slow := h.order(clock.time(), []endpoint{healthA, healthB})
	// The slow transfer hasn't reported within a cool-off, so the probe goes
	// to another transfer.
	clock.advance(10 * time.Second)
	fresh := h.order(clock.time(), []endpoint{healthA, healthB})
	eqEndpoints(t, "the fresh transfer's probes", probed(fresh), healthA)

	h.recordFailure(ctx, clock.time(), healthA, slow.probes[healthA])

	// The late failure ended neither the fresh probe nor the cool-off it runs
	// in.
	eqEndpoints(t, "after the late failure", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)
	clock.advance(10 * time.Second)
	eqEndpoints(t, "a cool-off later", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
}

func TestEndpointHealth_ignoresAProbeReturnedAfterItWasIssuedAgain(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	clock.advance(10 * time.Second)
	slow := h.order(clock.time(), []endpoint{healthA, healthB})
	clock.advance(10 * time.Second)
	h.order(clock.time(), []endpoint{healthA, healthB})

	h.returnProbe(clock.time(), healthA, slow.probes[healthA])

	// The fresh probe is still out, so no third transfer probes at once.
	eqEndpoints(t, "after the late return", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthB)
}

func TestEndpointHealth_ignoresFailuresOfTransfersThatStartedBeforeTheSkipOrHoldNoProbe(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(1, nil)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	// A transfer that started before the skip: it doesn't start the cool-off
	// again.
	clock.advance(5 * time.Second)
	h.recordFailure(ctx, clock.time(), healthA, 0)
	clock.advance(5 * time.Second)
	eqEndpoints(t, "the probe", probed(h.order(clock.time(), []endpoint{healthA, healthB})), healthA)
	// A transfer that holds no probe: it doesn't end the probe or double the
	// cool-off.
	h.recordFailure(ctx, clock.time(), healthA, 0)
	clock.advance(10 * time.Second)
	eqEndpoints(t, "a cool-off later", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
}

func TestEndpointHealth_withSkippingOffNeverSkips(t *testing.T) {
	ctx, clock, h := t.Context(), newTestClock(), testHealth(0, nil)
	for range 5 {
		h.recordFailure(ctx, clock.time(), healthA, 0)
	}
	eqEndpoints(t, "endpoints", h.order(clock.time(), []endpoint{healthA, healthB}).endpoints, healthA, healthB)
}

func TestEndpointHealth_concurrentTransfersMakeOneSkipOneLogLineAndOneProbe(t *testing.T) {
	var logs recordedLogs
	ctx, clock, h := t.Context(), newTestClock(), testHealth(50, &logs)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 8 {
		wg.Go(func() {
			<-start
			for range 100 {
				h.recordFailure(ctx, clock.time(), healthA, 0)
			}
		})
	}
	close(start)
	wg.Wait()
	if got := logs.levels("skipping it"); !slices.Equal(got, []slog.Level{slog.LevelWarn}) {
		t.Errorf("skip log levels = %v, want one warning", got)
	}

	clock.advance(10 * time.Second)
	var mu sync.Mutex
	probes := 0
	start = make(chan struct{})
	for range 8 {
		wg.Go(func() {
			<-start
			if order := h.order(clock.time(), []endpoint{healthA, healthB}); order.endpoints[0] == healthA {
				mu.Lock()
				probes++
				mu.Unlock()
			}
		})
	}
	close(start)
	wg.Wait()
	if probes != 1 {
		t.Errorf("transfers that probed = %d, want 1", probes)
	}
}

func TestCountsAgainstEndpoint(t *testing.T) {
	apiErr := func(status int, header http.Header) error {
		return apierr.FromHTTPError(status, header, []byte(`{"error_code":"X","message":"m"}`))
	}
	cancelled, cancel := context.WithCancel(t.Context())
	cancel()
	cases := []struct {
		name   string
		ctx    context.Context
		err    error
		counts bool
	}{
		{"403", t.Context(), apiErr(http.StatusForbidden, nil), true},
		{"404", t.Context(), apiErr(http.StatusNotFound, nil), true},
		{"500", t.Context(), apiErr(http.StatusInternalServerError, nil), true},
		{"an unexpected response", t.Context(), errUnexpectedServerResponse, true},
		{"a URL that keeps expiring", t.Context(), errUploadURLExpired, true},
		{"408", t.Context(), apiErr(http.StatusRequestTimeout, nil), false},
		{"429", t.Context(), apiErr(http.StatusTooManyRequests, nil), false},
		{"502", t.Context(), apiErr(http.StatusBadGateway, nil), false},
		{"503", t.Context(), apiErr(http.StatusServiceUnavailable, nil), false},
		{"504", t.Context(), apiErr(http.StatusGatewayTimeout, nil), false},
		{"a 500 with Retry-After", t.Context(), apiErr(http.StatusInternalServerError, http.Header{"Retry-After": {"7"}}), false},
		{"a wrapped 503", t.Context(), fmt.Errorf("part 3: %w", apiErr(http.StatusServiceUnavailable, nil)), false},
		{"a refused connection", t.Context(), &net.OpError{Op: "dial", Err: os.NewSyscallError("connect", syscall.ECONNREFUSED)}, false},
		{"a dropped request", t.Context(), &url.Error{Op: "Get", URL: "https://x", Err: io.ErrUnexpectedEOF}, false},
		{"a stalled single shot", t.Context(), &stalledError{idle: time.Minute}, false},
		{"a spent retry budget", t.Context(), fmt.Errorf("initiate: %w", context.DeadlineExceeded), false},
		{"the caller's cancellation", cancelled, apiErr(http.StatusForbidden, nil), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := countsAgainstEndpoint(tc.ctx, tc.err); got != tc.counts {
				t.Errorf("countsAgainstEndpoint(%v) = %v, want %v", tc.err, got, tc.counts)
			}
		})
	}
}

// ======================== Downloads ========================

// skipFake serves a download from the workspace's presigned URLs and its Files
// API, with faults per operation: "mint" (create-download-url), "presigned"
// (the presigned GET) and "workspace" (the Files API GET).
type skipFake struct {
	t    *testing.T
	srv  *httptest.Server
	data []byte

	mu sync.Mutex
	// faults answers the operation's 1-based call with a status, or with 0 to
	// serve it. A status of -1 drops the connection after the headers.
	faults map[string]func(call int) int
	counts map[string]int
}

func newSkipFake(t *testing.T, data []byte) *skipFake {
	f := &skipFake{t: t, data: data, faults: map[string]func(int) int{}, counts: map[string]int{}}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		if f.fault(w, r, "mint") {
			return
		}
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, f.srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("credentials were sent to a presigned URL")
		}
		if !f.fault(w, r, "presigned") {
			f.serve(w, r)
		}
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set("Last-Modified", testLastModified)
			w.Header().Set("Content-Length", strconv.Itoa(len(f.data)))
			return
		}
		if !f.fault(w, r, "workspace") {
			f.serve(w, r)
		}
	}))
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

// fault counts a call to op and answers it if a fault applies.
func (f *skipFake) fault(w http.ResponseWriter, r *http.Request, op string) bool {
	f.mu.Lock()
	f.counts[op]++
	call, fault := f.counts[op], f.faults[op]
	f.mu.Unlock()
	if fault == nil {
		return false
	}
	switch status := fault(call); status {
	case 0:
		return false
	case -1:
		// The response's headers, then the connection drops before a byte.
		body, status := f.rangeOf(r)
		w.Header().Set("Last-Modified", testLastModified)
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(status)
		w.(http.Flusher).Flush()
		if conn, _, err := w.(http.Hijacker).Hijack(); err == nil {
			_ = conn.Close()
		}
	default:
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, `{"error_code":"E%d","message":"fault"}`, status)
	}
	return true
}

// rangeOf returns the bytes r asks for and the status that answers it.
func (f *skipFake) rangeOf(r *http.Request) ([]byte, int) {
	spec, ok := strings.CutPrefix(r.Header.Get("Range"), "bytes=")
	if !ok {
		return f.data, http.StatusOK
	}
	startText, endText, _ := strings.Cut(spec, "-")
	start, _ := strconv.Atoi(startText)
	end := len(f.data) - 1
	if endText != "" {
		end, _ = strconv.Atoi(endText)
	}
	return f.data[start : end+1], http.StatusPartialContent
}

func (f *skipFake) serve(w http.ResponseWriter, r *http.Request) {
	body, status := f.rangeOf(r)
	w.Header().Set("Last-Modified", testLastModified)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func (f *skipFake) setFault(op string, fault func(call int) int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if fault == nil {
		delete(f.faults, op)
	} else {
		f.faults[op] = fault
	}
}

func (f *skipFake) count(op string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.counts[op]
}

// httpStatus is the HTTP status of err's API error, or 0.
func httpStatus(err error) int {
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok {
		return aerr.HTTPStatusCode()
	}
	return 0
}

func always(status int) func(int) int { return func(int) int { return status } }

func onCalls(status int, calls ...int) func(int) int {
	return func(call int) int {
		if slices.Contains(calls, call) {
			return status
		}
		return 0
	}
}

// skipEngine is an engine over a fresh client of f, with retries that back off
// in milliseconds, ranges small enough to split f's data, and clock as its
// clock.
func skipEngine(t *testing.T, host string, hc *http.Client, clock *testClock) *engine {
	t.Helper()
	e := newEngine(buildUploadClient(t, host, hc, ""))
	fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	e.tun.downloadMinParallelSize, e.tun.downloadPartSize = 1000, 1000
	e.now = clock.time
	return e
}

const skipPath = "/Volumes/c/s/v/f.bin"

// downloadModes are the ways a test downloads: a stream, to a file, and to a
// file in parallel ranges.
var downloadModes = []string{"stream", "to-file", "parallel"}

// skipDownload downloads skipPath in the given mode.
func skipDownload(t *testing.T, e *engine, mode string) ([]byte, error) {
	t.Helper()
	if mode == "stream" {
		resp, err := e.download(t.Context(), skipPath, downloadConfig{})
		if err != nil {
			return nil, err
		}
		defer resp.Contents.Close()
		return io.ReadAll(resp.Contents)
	}
	dest := filepath.Join(t.TempDir(), "f.bin")
	if err := e.downloadTo(t.Context(), skipPath, dest, downloadConfig{parallel: mode == "parallel", parallelism: new(2)}); err != nil {
		return nil, err
	}
	return os.ReadFile(dest)
}

func TestDownload_anEndpointThatKeepsFailingIsSkippedThenProbed(t *testing.T) {
	for _, mode := range downloadModes {
		t.Run(mode, func(t *testing.T) {
			want := payload(4000)
			f := newSkipFake(t, want)
			f.setFault("presigned", always(http.StatusInternalServerError))
			clock := newTestClock()
			e := skipEngine(t, f.srv.URL, f.srv.Client(), clock)
			var logs recordedLogs
			e.c.logger = slog.New(&logs)
			skipAfter, coolOff := e.tun.endpointSkipAfter, e.tun.endpointSkipCoolOff

			// triesPresigned runs one download, which must complete, and tells
			// whether it tried the presigned URL.
			triesPresigned := func() bool {
				t.Helper()
				before := f.count("presigned")
				got, err := skipDownload(t, e, mode)
				noErr(t, err)
				eqBytes(t, got, want)
				return f.count("presigned") > before
			}
			var tried []bool
			for range skipAfter + 2 {
				tried = append(tried, triesPresigned())
			}
			if want := []bool{true, true, true, false, false}; !slices.Equal(tried, want) {
				t.Fatalf("tried presigned = %v, want %v", tried, want)
			}
			if got := logs.levels("skipping it"); !slices.Equal(got, []slog.Level{slog.LevelWarn}) {
				t.Errorf("skip log levels = %v, want one warning", got)
			}

			// After the cool-off, one transfer probes presigned URLs. They
			// still fail, so the next transfer skips them.
			clock.advance(coolOff)
			if got := []bool{triesPresigned(), triesPresigned()}; !slices.Equal(got, []bool{true, false}) {
				t.Errorf("after the cool-off, tried presigned = %v, want [true false]", got)
			}
			// The failed probe doubled the cool-off.
			clock.advance(coolOff)
			if triesPresigned() {
				t.Error("a cool-off after the failed probe, presigned URLs were tried")
			}
			clock.advance(coolOff)
			f.setFault("presigned", nil)
			// The probe succeeds, so every later transfer uses presigned URLs
			// and nothing else.
			workspaceBefore := f.count("workspace")
			for i := range 3 {
				if !triesPresigned() {
					t.Errorf("transfer %d after the successful probe skipped presigned URLs", i+1)
				}
			}
			if got := f.count("workspace"); got != workspaceBefore {
				t.Errorf("workspace GETs after the successful probe = %d, want none", got-workspaceBefore)
			}
		})
	}
}

func TestDownload_anEndpointIsNotSkippedWhenEveryEndpointFails(t *testing.T) {
	// A missing file fails everywhere: the request is at fault, not the
	// endpoints.
	for _, mode := range downloadModes {
		t.Run(mode, func(t *testing.T) {
			f := newSkipFake(t, payload(4000))
			f.setFault("presigned", always(http.StatusNotFound))
			f.setFault("workspace", always(http.StatusNotFound))
			e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
			n := e.tun.endpointSkipAfter + 2
			for range n {
				if _, err := skipDownload(t, e, mode); httpStatus(err) != http.StatusNotFound {
					t.Fatalf("err = %v, want a 404", err)
				}
			}
			if got := f.count("presigned"); got != n {
				t.Errorf("presigned GETs = %d, want %d", got, n)
			}
		})
	}
}

func TestDownload_anEndpointIsNotSkippedForFailuresThatAreRetried(t *testing.T) {
	// A throttled or unavailable endpoint may recover on its own, so it keeps
	// getting the first attempt. Before the first byte, with an endpoint left,
	// such a failure walks on without being retried.
	cases := map[string]struct {
		op     string
		status int
	}{
		"a 429 minting the download URL": {"mint", http.StatusTooManyRequests},
		"a 503 minting the download URL": {"mint", http.StatusServiceUnavailable},
		"a 408 on the presigned GET":     {"presigned", http.StatusRequestTimeout},
		"a 429 on the presigned GET":     {"presigned", http.StatusTooManyRequests},
		"a 502 on the presigned GET":     {"presigned", http.StatusBadGateway},
		"a 503 on the presigned GET":     {"presigned", http.StatusServiceUnavailable},
		"a 504 on the presigned GET":     {"presigned", http.StatusGatewayTimeout},
		"a break before the first byte":  {"presigned", -1},
	}
	for name, tc := range cases {
		for _, mode := range downloadModes {
			t.Run(name+"/"+mode, func(t *testing.T) {
				want := payload(4000)
				f := newSkipFake(t, want)
				f.setFault(tc.op, always(tc.status))
				e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
				n := e.tun.endpointSkipAfter + 2
				for range n {
					got, err := skipDownload(t, e, mode)
					noErr(t, err)
					eqBytes(t, got, want)
				}
				if got := f.count(tc.op); got < n {
					t.Errorf("%s calls = %d, want at least %d", tc.op, got, n)
				}
				if got := f.count("workspace"); got < n {
					t.Errorf("workspace GETs = %d, want at least %d", got, n)
				}
			})
		}
	}
}

func TestDownload_aStallBeforeTheFirstByteIsNotCountedAgainstTheEndpoint(t *testing.T) {
	want := payload(100)
	var srv *httptest.Server
	var mu sync.Mutex
	stalls := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		stalls++
		mu.Unlock()
		// The headers, then no byte until the client gives up.
		w.Header().Set("Last-Modified", testLastModified)
		w.Header().Set("Content-Length", strconv.Itoa(len(want)))
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Last-Modified", testLastModified)
		_, _ = w.Write(want)
	}))
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	e := skipEngine(t, srv.URL, srv.Client(), newTestClock())
	e.tun.transferIdleTimeout = 50 * time.Millisecond
	n := e.tun.endpointSkipAfter + 2
	for range n {
		got, err := skipDownload(t, e, "stream")
		noErr(t, err)
		eqBytes(t, got, want)
	}
	mu.Lock()
	defer mu.Unlock()
	if stalls != n {
		t.Errorf("presigned GETs = %d, want %d: a stall counted against presigned URLs", stalls, n)
	}
}

func TestDownload_aSuccessClearsTheFailureCount(t *testing.T) {
	for _, mode := range downloadModes {
		t.Run(mode, func(t *testing.T) {
			want := payload(4000)
			f := newSkipFake(t, want)
			e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
			skipAfter := e.tun.endpointSkipAfter
			// Every transfer fails on presigned URLs but one in the middle, so
			// the failures in a row never reach the threshold. A parallel
			// download sends more than one presigned GET, so the fault goes by
			// transfer.
			var transfer atomic.Int32
			f.setFault("presigned", func(int) int {
				if int(transfer.Load()) == skipAfter {
					return 0
				}
				return http.StatusInternalServerError
			})
			for i := 1; i <= 2*skipAfter; i++ {
				transfer.Store(int32(i))
				before := f.count("presigned")
				got, err := skipDownload(t, e, mode)
				noErr(t, err)
				eqBytes(t, got, want)
				if f.count("presigned") == before {
					t.Fatalf("transfer %d skipped presigned URLs", i)
				}
			}
		})
	}
}

func TestDownload_skippingCanBeTurnedOff(t *testing.T) {
	want := payload(4000)
	f := newSkipFake(t, want)
	f.setFault("presigned", always(http.StatusInternalServerError))
	e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
	e.tun.endpointSkipAfter = 0
	for range 5 {
		got, err := skipDownload(t, e, "stream")
		noErr(t, err)
		eqBytes(t, got, want)
	}
	if got := f.count("presigned"); got != 5 {
		t.Errorf("presigned GETs = %d, want 5", got)
	}
}

func TestDownload_aSkippedEndpointIsNotTriedEvenWhenTheOthersFail(t *testing.T) {
	f := newSkipFake(t, payload(4000))
	f.setFault("presigned", always(http.StatusInternalServerError))
	e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
	for range e.tun.endpointSkipAfter {
		_, err := skipDownload(t, e, "stream")
		noErr(t, err)
	}
	// Presigned URLs are skipped, then the workspace fails too.
	f.setFault("workspace", always(http.StatusInternalServerError))
	before := f.count("presigned")
	if _, err := skipDownload(t, e, "stream"); httpStatus(err) != http.StatusInternalServerError {
		t.Fatalf("err = %v, want a 500", err)
	}
	if got := f.count("presigned"); got != before {
		t.Errorf("presigned GETs = %d, want none", got-before)
	}
}

func TestDownload_theLastEndpointKeepsItsRetriesWhenAnEndpointIsSkipped(t *testing.T) {
	// Presigned downloads keep failing, as behind a storage firewall.
	for _, mode := range []string{"stream", "to-file"} {
		t.Run(mode, func(t *testing.T) {
			want := payload(4000)
			f := newSkipFake(t, want)
			f.setFault("presigned", always(http.StatusForbidden))
			e := skipEngine(t, f.srv.URL, f.srv.Client(), newTestClock())
			for range e.tun.endpointSkipAfter {
				_, err := skipDownload(t, e, mode)
				noErr(t, err)
			}
			presignedBefore := f.count("presigned")
			// The workspace throttles the next request once.
			f.setFault("workspace", onCalls(http.StatusServiceUnavailable, f.count("workspace")+1))

			got, err := skipDownload(t, e, mode)
			noErr(t, err)
			eqBytes(t, got, want)

			// The workspace is still the last endpoint, so it retried the 503
			// rather than handing the transfer to presigned URLs.
			if got := f.count("presigned"); got != presignedBefore {
				t.Errorf("presigned GETs = %d, want none", got-presignedBefore)
			}
		})
	}
}

// ======================== Uploads ========================

// uploadOps counts f's requests by operation: "initiate", "part-urls", "part"
// and "single".
func uploadOps(f *fakeServer) func(op string) int {
	var mu sync.Mutex
	counts := map[string]int{}
	f.onRequest = func(r *http.Request) {
		op := "other"
		switch {
		case r.URL.Query().Get("action") == "initiate-upload":
			op = "initiate"
		case strings.HasSuffix(r.URL.Path, "/create-upload-part-urls"):
			op = "part-urls"
		case strings.HasPrefix(r.URL.Path, "/cloud/part/"):
			op = "part"
		case r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/api/2.0/fs/files/"):
			op = "single"
		}
		mu.Lock()
		counts[op]++
		mu.Unlock()
	}
	return func(op string) int {
		mu.Lock()
		defer mu.Unlock()
		return counts[op]
	}
}

// uploadSkipEngine is a walkEngine over f with clock as its clock.
func uploadSkipEngine(t *testing.T, f *fakeServer, clock *testClock) *engine {
	t.Helper()
	e := walkEngine(t, f)
	e.now = clock.time
	return e
}

// sentByPresigned reports whether the last upload went to presigned URLs:
// completed as a multipart upload, with no single shot since reset.
func sentByPresigned(f *fakeServer) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.singleBody == nil && len(f.completed) > 0
}

// resetStored forgets what f stored, so that the next upload's route shows.
func resetStored(f *fakeServer) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.singleBody, f.completed = nil, nil
}

func TestUpload_anEndpointThatKeepsFailingIsSkipped(t *testing.T) {
	sources := map[string]func(t *testing.T, e *engine, b []byte) error{
		"file": func(t *testing.T, e *engine, b []byte) error {
			_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, b), WithParallelism(1))
			return err
		},
		"seekable": func(t *testing.T, e *engine, b []byte) error {
			_, err := e.upload(t.Context(), skipPath, bytes.NewReader(b), WithParallelism(1))
			return err
		},
	}
	for name, upload := range sources {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.initiateStatus = http.StatusInternalServerError
			ops := uploadOps(f)
			e := uploadSkipEngine(t, f, newTestClock())
			in := data(5000)
			for range e.tun.endpointSkipAfter + 2 {
				resetStored(f)
				noErr(t, upload(t, e, in))
				eqBytes(t, f.assembled(), in)
			}
			if got := ops("initiate"); got != e.tun.endpointSkipAfter {
				t.Errorf("presigned initiates = %d, want %d", got, e.tun.endpointSkipAfter)
			}
		})
	}
}

func TestUpload_anEndpointIsNotSkippedForFailuresThatAreRetried(t *testing.T) {
	cases := map[string]func(f *fakeServer){
		// Minting part URLs keeps getting 429 until the retries before commit
		// run out: a passing fault, not the endpoint's.
		"a 429 on the part URLs": func(f *fakeServer) {
			f.partURLsStatus, f.partURLsBody = http.StatusTooManyRequests, `{"error_code":"TOO_MANY_REQUESTS","message":"slow down"}`
		},
		"a 503 on every part": func(f *fakeServer) {
			f.partHook = func(int, int) (int, string, []byte) { return http.StatusServiceUnavailable, "", nil }
		},
	}
	for name, fault := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			fault(f)
			ops := uploadOps(f)
			e := uploadSkipEngine(t, f, newTestClock())
			in := data(5000)
			n := e.tun.endpointSkipAfter + 2
			for i := range n {
				resetStored(f)
				before := ops("initiate")
				_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, in), WithParallelism(1))
				noErr(t, err)
				eqBytes(t, f.assembled(), in)
				if ops("initiate") == before {
					t.Fatalf("upload %d skipped presigned URLs", i+1)
				}
			}
		})
	}
}

func TestUpload_a503ThatARetryRecoversNeverSkips(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// Every part's first attempt of each upload gets a 503; its one retry
	// before commit lands it.
	f.partHook = func(_, attempt int) (int, string, []byte) {
		if attempt%2 == 1 {
			return http.StatusServiceUnavailable, "", nil
		}
		return 0, "", nil
	}
	e := uploadSkipEngine(t, f, newTestClock())
	in := data(5000)
	for i := range e.tun.endpointSkipAfter + 2 {
		resetStored(f)
		_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, in), WithParallelism(1))
		noErr(t, err)
		eqBytes(t, f.assembled(), in)
		if !sentByPresigned(f) {
			t.Fatalf("upload %d didn't go to presigned URLs", i+1)
		}
	}
}

func TestUpload_aProbeATransferNeverTriesGoesToTheNextTransfer(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.partURLsStatus, f.partURLsBody = http.StatusInternalServerError, `{"error_code":"INTERNAL_ERROR","message":"boom"}`
	clock := newTestClock()
	e := uploadSkipEngine(t, f, clock)
	large := data(5000)
	for range e.tun.endpointSkipAfter {
		resetStored(f)
		_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, large), WithParallelism(1))
		noErr(t, err)
		if sentByPresigned(f) {
			t.Fatal("an upload went to presigned URLs while minting their part URLs failed")
		}
	}
	f.partURLsStatus = 0
	clock.advance(e.tun.endpointSkipCoolOff)
	// A small upload is ordered while presigned URLs are due a probe, but
	// takes no presigned request.
	resetStored(f)
	_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, data(100)), WithParallelism(1))
	noErr(t, err)
	if sentByPresigned(f) {
		t.Fatal("a small upload went to presigned URLs")
	}

	resetStored(f)
	_, err = e.uploadFrom(t.Context(), skipPath, writeTemp(t, large), WithParallelism(1))
	noErr(t, err)

	// The probe went to this upload, which tried presigned URLs first.
	if !sentByPresigned(f) {
		t.Error("the upload after the small one didn't probe presigned URLs")
	}
}

func TestUpload_aSkippedEndpointIsNotTriedEvenWhenTheOthersFail(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.initiateStatus = http.StatusInternalServerError
	ops := uploadOps(f)
	e := uploadSkipEngine(t, f, newTestClock())
	in := data(5000)
	for range e.tun.endpointSkipAfter {
		_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, in), WithParallelism(1))
		noErr(t, err)
	}
	f.singleShotStatus, f.singleShotBody = http.StatusInternalServerError, `{"error_code":"INTERNAL_ERROR","message":"boom"}`
	before := ops("initiate")
	_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, in), WithParallelism(1))
	if httpStatus(err) != http.StatusInternalServerError {
		t.Fatalf("err = %v, want a 500", err)
	}
	if got := ops("initiate"); got != before {
		t.Errorf("presigned initiates = %d, want none", got-before)
	}
}

func TestUpload_anUploadOnlyASkippedEndpointCanTakeStillUsesIt(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := uploadSkipEngine(t, f, newTestClock())
	health := e.endpointHealth(true)
	every := e.builtinEndpoints()
	presigned, workspace := every[0], every[1]
	for range e.tun.endpointSkipAfter {
		health.recordFailure(t.Context(), e.clock(), presigned, 0)
	}
	plannedTo := func(size int64) []endpoint {
		_, plans, _ := e.planUploads(health, every, size)
		var out []endpoint
		for _, p := range plans {
			out = append(out, p.ep)
		}
		return out
	}
	// Presigned URLs are skipped for what the workspace takes in a single
	// shot.
	eqEndpoints(t, "a small upload", plannedTo(100), workspace)
	eqEndpoints(t, "an upload of 5 GiB", plannedTo(workspaceMaxSingleShotSize), workspace)
	// Over the workspace's single-shot limit, or of unknown size, which a
	// single shot can't send again, only the skipped multipart endpoint can
	// take it.
	eqEndpoints(t, "an upload over 5 GiB", plannedTo(workspaceMaxSingleShotSize+1), presigned)
	eqEndpoints(t, "an upload of unknown size", plannedTo(-1), presigned, workspace)

	// End to end, a stream of unknown size goes to presigned URLs.
	in := data(5000)
	_, err := e.upload(t.Context(), skipPath, io.MultiReader(bytes.NewReader(in)), WithParallelism(1))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if !sentByPresigned(f) {
		t.Error("the stream of unknown size didn't go to presigned URLs")
	}
}

func TestUpload_downloadFailuresDoNotChangeHowUploadsSendData(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := uploadSkipEngine(t, f, newTestClock())
	presigned := e.builtinEndpoints()[0]
	for range e.tun.endpointSkipAfter {
		e.endpointHealth(false).recordFailure(t.Context(), e.clock(), presigned, 0)
	}
	if got := e.endpointHealth(false).order(e.clock(), e.builtinEndpoints()).endpoints; got[0] == presigned {
		t.Fatal("downloads didn't skip presigned URLs")
	}
	in := data(5000)
	_, err := e.uploadFrom(t.Context(), skipPath, writeTemp(t, in), WithParallelism(1))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if !sentByPresigned(f) {
		t.Error("the upload didn't go to presigned URLs")
	}
}
