package files

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
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

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/options/client"
)

// testToken is the Authorization value a Databricks-authenticated request in
// these tests carries. Control-plane calls must attach it; presigned
// cloud-storage transfers must not.
const testToken = "Bearer test-token"

// testCredentials stamps testToken on every request the client builds. It emits
// a real value rather than nothing so that a request which drops the client's
// credentials is distinguishable from one that carries them.
type testCredentials struct{}

func (testCredentials) Name() string { return "test" }

func (testCredentials) AuthHeaders(context.Context) ([]auth.Header, error) {
	return []auth.Header{{Key: "Authorization", Value: testToken}}, nil
}

// buildUploadClient builds a v2 Client whose control plane is served by hc and
// addressed at host, with an optional workspace ID for the routing header.
func buildUploadClient(t *testing.T, host string, hc *http.Client, workspaceID string) *Client {
	t.Helper()
	opts := []client.Option{
		client.WithHost(host),
		client.WithHTTPClient(hc),
		client.WithCredentials(testCredentials{}),
		client.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))),
		// Keep tests hermetic: without these, an unset workspace ID could be
		// filled from the developer's local configuration.
		client.WithoutConfigFile(),
		client.WithoutEnv(),
	}
	if workspaceID != "" {
		opts = append(opts, client.WithWorkspaceID(workspaceID))
	}
	c, err := NewClient(context.Background(), opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// fakeServer emulates both the Files API control plane and the cloud storage
// provider: the presigned URLs it mints point back at itself, so a single
// httptest server backs the whole upload. It reassembles the uploaded object so
// tests can assert byte-for-byte equality with the input.
type fakeServer struct {
	srv *httptest.Server
	// onRequest, when set, sees every request first.
	onRequest func(*http.Request)
	// endpoints, when set, answers GET /api/2.0/fs/endpoints, with each
	// $HOST in it the server's host.
	endpoints string

	mu          sync.Mutex
	mode        string // "multipart", or "resumable" to offer a resumable session
	parts       map[int][]byte
	completed   []completePart
	singleBody  []byte
	overwriteQ  []string      // overwrite query values seen on initiate + single-shot
	partHeaders []http.Header // headers seen on cloud part PUTs
	partURLReqs []int         // Count of each create-upload-part-urls request
	aborts      []string      // session_token of each abort-upload request

	// abortStatus, when non-zero, is returned by abort-upload.
	abortStatus int

	// initiateStatus and initiateBody, when the status is non-zero, answer
	// initiate-upload.
	initiateStatus int
	initiateBody   string

	// partURLsStatus and partURLsBody, when the status is non-zero, answer
	// create-upload-part-urls.
	partURLsStatus int
	partURLsBody   string

	// partHook injects faults. It returns (statusCode, etag, body); a zero
	// statusCode means "use the default 200 + store".
	partHook func(n, attempt int) (int, string, []byte)
	partTry  map[int]int

	// singleShotStatus, when non-zero, is returned by the single-shot PUT.
	singleShotStatus int
	singleShotBody   string

	// singleShotHook injects per-attempt faults on the single-shot PUT. It
	// returns the status for the given 1-based attempt; a zero status means "use
	// the default 200 + store". singleTry counts single-shot PUTs seen.
	singleShotHook func(attempt int) int
	singleTry      int

	// singleShotRead, when set, reads the single-shot body in place of
	// io.ReadAll, to pace or stall it. It returns false to leave the request
	// unanswered. Tests count its calls themselves: an attempt the client gave
	// up on may return after the next one began.
	singleShotRead func(r *http.Request) ([]byte, bool)
}

func newFakeServer(t *testing.T, mode string) *fakeServer {
	f := &fakeServer{
		mode:    mode,
		parts:   map[int][]byte{},
		partTry: map[int]int{},
	}
	f.srv = httptest.NewServer(f)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeServer) base() string { return f.srv.URL }

// isControlPlane reports whether a path is a Files API control-plane endpoint,
// as opposed to a presigned cloud-storage URL this server also mints.
func isControlPlane(path string) bool { return strings.HasPrefix(path, "/api/2.0/fs/") }

func (f *fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.onRequest != nil {
		f.onRequest(r)
	}
	path := r.URL.Path

	// Enforce the credential split the real deployment enforces, so any upload
	// test fails when a call authenticates wrongly: the control plane rejects a
	// request without the client's credentials ("Credential was not sent"), and a
	// presigned cloud URL must never receive them.
	if isControlPlane(path) && r.Header.Get("Authorization") != testToken {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"error_code":"UNAUTHENTICATED","message":"Credential was not sent or was of an unsupported type"}`)
		return
	}
	if strings.HasPrefix(path, "/cloud/") && r.Header.Get("Authorization") != "" {
		http.Error(w, "presigned cloud URL received Databricks credentials", http.StatusInternalServerError)
		return
	}

	switch {
	case r.Method == http.MethodGet && path == endpointsAPIPath && f.endpoints != "":
		_, _ = io.WriteString(w, strings.ReplaceAll(f.endpoints, "$HOST", r.Host))
	case r.Method == http.MethodPost && strings.HasSuffix(path, "/create-upload-part-urls"):
		f.handleCreatePartURLs(w, r)
	case r.Method == http.MethodPut && strings.HasPrefix(path, "/cloud/part/"):
		f.handlePartPut(w, r)
	case strings.Contains(path, "/api/2.0/fs/files/"):
		f.handleFiles(w, r)
	default:
		http.Error(w, "unexpected request: "+r.Method+" "+path, http.StatusInternalServerError)
	}
}

func (f *fakeServer) handleFiles(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	switch {
	case r.Method == http.MethodPost && action == "initiate-upload":
		f.mu.Lock()
		f.overwriteQ = append(f.overwriteQ, r.URL.Query().Get("overwrite"))
		f.mu.Unlock()
		if f.initiateStatus != 0 {
			w.WriteHeader(f.initiateStatus)
			_, _ = io.WriteString(w, f.initiateBody)
			return
		}
		if f.mode == "resumable" {
			// A resumable session, which this package does not support.
			writeJSON(w, map[string]any{"resumable_upload": map[string]any{"session_token": "tok"}})
		} else {
			writeJSON(w, initiateResult{MultipartUpload: &uploadSession{SessionToken: "tok"}})
		}
	case r.Method == http.MethodPost && action == "abort-upload":
		if q := r.URL.Query(); q.Get("upload_type") != "multipart" {
			http.Error(w, "abort-upload without upload_type=multipart", http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		f.aborts = append(f.aborts, r.URL.Query().Get("session_token"))
		status := f.abortStatus
		f.mu.Unlock()
		if status != 0 {
			w.WriteHeader(status)
			return
		}
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodPost && action == "complete-upload":
		var req completeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		f.mu.Lock()
		f.completed = req.Parts
		f.mu.Unlock()
		w.WriteHeader(http.StatusOK)
	case r.Method == http.MethodPut: // single-shot
		f.mu.Lock()
		read := f.singleShotRead
		f.mu.Unlock()
		var body []byte
		if read != nil {
			var ok bool
			if body, ok = read(r); !ok {
				return
			}
		} else {
			body, _ = io.ReadAll(r.Body)
		}
		f.mu.Lock()
		f.singleBody = body
		f.overwriteQ = append(f.overwriteQ, r.URL.Query().Get("overwrite"))
		f.singleTry++
		status, respBody, hook, attempt := f.singleShotStatus, f.singleShotBody, f.singleShotHook, f.singleTry
		f.mu.Unlock()
		if hook != nil {
			if s := hook(attempt); s != 0 {
				w.WriteHeader(s)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		if status != 0 {
			w.WriteHeader(status)
			_, _ = io.WriteString(w, respBody)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "unexpected files request", http.StatusInternalServerError)
	}
}

func (f *fakeServer) handleCreatePartURLs(w http.ResponseWriter, r *http.Request) {
	var req createPartURLsRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	f.mu.Lock()
	f.partURLReqs = append(f.partURLReqs, req.Count)
	f.mu.Unlock()
	if f.partURLsStatus != 0 {
		w.WriteHeader(f.partURLsStatus)
		_, _ = io.WriteString(w, f.partURLsBody)
		return
	}
	var out createPartURLsResponse
	for i := range req.Count {
		n := req.StartPartNumber + i
		out.UploadPartURLs = append(out.UploadPartURLs, struct {
			URL        string      `json:"url"`
			PartNumber int         `json:"part_number"`
			Headers    []nameValue `json:"headers"`
		}{URL: f.base() + "/cloud/part/" + strconv.Itoa(n), PartNumber: n})
	}
	writeJSON(w, out)
}

func (f *fakeServer) handlePartPut(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/cloud/part/"))
	body, _ := io.ReadAll(r.Body)

	f.mu.Lock()
	f.partHeaders = append(f.partHeaders, r.Header.Clone())
	f.partTry[n]++
	attempt := f.partTry[n]
	hook := f.partHook
	f.mu.Unlock()

	if hook != nil {
		if status, etag, respBody := hook(n, attempt); status != 0 {
			if etag != "" {
				w.Header().Set("ETag", etag)
			}
			w.WriteHeader(status)
			_, _ = w.Write(respBody)
			return
		}
	}

	f.mu.Lock()
	f.parts[n] = body
	f.mu.Unlock()
	w.Header().Set("ETag", "etag-"+strconv.Itoa(n))
	w.WriteHeader(http.StatusOK)
}

// abortsSeen returns the session_token of each abort-upload request.
func (f *fakeServer) abortsSeen() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.aborts)
}

// assembled returns the bytes the server reconstructed for the upload.
func (f *fakeServer) assembled() []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.singleBody != nil {
		return f.singleBody
	}
	var out []byte
	for _, p := range f.completed {
		out = append(out, f.parts[p.PartNumber]...)
	}
	return out
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func newTestClient(t *testing.T, f *fakeServer) *Client {
	t.Helper()
	return buildUploadClient(t, f.base(), f.srv.Client(), "")
}

func TestSetWorkspaceHeader(t *testing.T) {
	// A workspace-scoped client stamps the routing header the Files API control
	// plane needs; a client without a workspace ID sets nothing.
	e := newEngine(buildUploadClient(t, "https://h.test", http.DefaultClient, "1234"))
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://h.test/api/2.0/fs/files/x", nil)
	noErr(t, err)
	e.setWorkspaceHeader(req)
	if got := req.Header.Get("X-Databricks-Workspace-Id"); got != "1234" {
		t.Errorf("X-Databricks-Workspace-Id = %q, want 1234", got)
	}

	e2 := newEngine(buildUploadClient(t, "https://h.test", http.DefaultClient, ""))
	req2, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "https://h.test/api/2.0/fs/files/x", nil)
	noErr(t, err)
	e2.setWorkspaceHeader(req2)
	if got := req2.Header.Get("X-Databricks-Workspace-Id"); got != "" {
		t.Errorf("X-Databricks-Workspace-Id = %q, want empty", got)
	}
}

// failingSeeker reports a size (Seek to the end succeeds) but cannot rewind
// (Seek to the start fails), modeling the unrecoverable probe case.
type failingSeeker struct{ io.ReadSeeker }

func (f failingSeeker) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		return 0, errors.New("seek start failed")
	}
	return f.ReadSeeker.Seek(offset, whence)
}

func TestUploadRewindFailureIsLoud(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)
	// A reader the engine can size but not rewind must fail the upload rather
	// than send a truncated object from the EOF position.
	r := failingSeeker{bytes.NewReader(make([]byte, 1024))}
	if _, err := e.upload(t.Context(), "/Volumes/c/s/v/f.bin", r); err == nil {
		t.Fatal("expected an error when the reader cannot be rewound after sizing")
	}
}

// testEngine returns an engine over a client pointed at f, with production
// tunables.
func testEngine(t *testing.T, f *fakeServer) *engine {
	t.Helper()
	return newEngine(newTestClient(t, f))
}

// shrunkEngine returns an engine whose size thresholds are tiny, so multipart
// boundaries are hit with small inputs. Tunables are injected per engine, so
// tests never mutate shared state and can run in parallel.
func shrunkEngine(t *testing.T, f *fakeServer) *engine {
	t.Helper()
	e := testEngine(t, f)
	e.tun.minStreamSize = 1024
	e.tun.defaultPartSize = 1024
	return e
}

// --- assertion helpers (standard library testing; no third-party framework) ---

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func eqBytes(t *testing.T, got, want []byte) {
	t.Helper()
	if !bytes.Equal(got, want) {
		t.Errorf("uploaded bytes differ: got %d bytes, want %d bytes", len(got), len(want))
	}
}

func data(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 251)
	}
	return b
}

func writeTemp(t *testing.T, b []byte) string {
	t.Helper()
	p := t.TempDir() + "/src.bin"
	noErr(t, os.WriteFile(p, b, 0o600))
	return p
}

func TestUploadSingleShot(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	in := data(500) // below the 1024 single-shot threshold
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/small.bin", bytes.NewReader(in))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

func TestUploadMultipartParallelMatchesSequential(t *testing.T) {
	in := data(20 * 1024)

	run := func(parallelism int) []byte {
		f := newFakeServer(t, "multipart")
		e := shrunkEngine(t, f)
		_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/big.bin", writeTemp(t, in),
			WithParallelism(parallelism))
		noErr(t, err)
		// Parts must be completed in sorted order with non-empty ETags.
		if len(f.completed) == 0 {
			t.Fatal("expected completed parts")
		}
		for i, p := range f.completed {
			if p.PartNumber != i+1 {
				t.Errorf("part %d: got part number %d, want %d", i, p.PartNumber, i+1)
			}
			if p.ETag == "" {
				t.Errorf("part %d: empty ETag", p.PartNumber)
			}
		}
		return f.assembled()
	}

	eqBytes(t, run(1), in)
	eqBytes(t, run(8), in)
}

func TestUploadFromReaderAtKnownSize(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	in := data(5000)
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/r.bin", bytes.NewReader(in), WithParallelism(4))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

func TestUploadNonSeekableStream(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	in := data(5000)
	// io.Reader wrapper that is neither Seeker nor ReaderAt.
	r := io.MultiReader(bytes.NewReader(in))
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/stream.bin", r, WithParallelism(4))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

func TestCloudPartsCarryNoDatabricksCredentials(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/big.bin", writeTemp(t, data(5000)), WithParallelism(4))
	noErr(t, err)

	if len(f.partHeaders) == 0 {
		t.Fatal("expected cloud part requests")
	}
	for _, h := range f.partHeaders {
		if got := h.Get("Authorization"); got != "" {
			t.Errorf("cloud part PUT must not carry Databricks auth, got Authorization=%q", got)
		}
		if got := h.Get("X-Databricks-Workspace-Id"); got != "" {
			t.Errorf("cloud part PUT must not carry workspace routing, got X-Databricks-Workspace-Id=%q", got)
		}
	}
}

func TestUploadOverwrite(t *testing.T) {
	cases := []struct {
		name string
		opts []UploadOption
		want string
	}{
		{"unset", nil, ""},
		{"true", []UploadOption{WithOverwrite(true)}, "true"},
		{"false", []UploadOption{WithOverwrite(false)}, "false"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			e := shrunkEngine(t, f)
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/x.bin", bytes.NewReader(data(100)), tc.opts...)
			noErr(t, err)
			if len(f.overwriteQ) == 0 {
				t.Fatal("expected an overwrite query value")
			}
			if f.overwriteQ[0] != tc.want {
				t.Errorf("overwrite query = %q, want %q", f.overwriteQ[0], tc.want)
			}
		})
	}
}

func TestFileUploadMintsURLsInBatches(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	in := data(16 * 1024) // 16 parts at the 1024 part size
	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/batch.bin", writeTemp(t, in), WithParallelism(8))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)

	// The file path pre-mints in batches, so at least one create-upload-part-urls
	// call requests more than one URL; per-part minting would never exceed 1.
	maxCount := 0
	for _, cnt := range f.partURLReqs {
		maxCount = max(maxCount, cnt)
	}
	if maxCount <= 1 {
		t.Errorf("expected batched URL minting (a request with count > 1), got max count %d across %d calls", maxCount, len(f.partURLReqs))
	}
}

func TestCloudRejectsAllPartsFallsBackToSingleShot(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// A cloud firewall forbids every part PUT on every attempt (a blanket block of
	// the storage host). Parts upload concurrently, so no single part is the canary;
	// because none lands, the upload falls back to single-shot.
	f.partHook = func(n, attempt int) (int, string, []byte) {
		return http.StatusForbidden, "", []byte("forbidden")
	}
	e := shrunkEngine(t, f)

	in := data(5000)
	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/fb.bin", writeTemp(t, in), WithParallelism(4))
	noErr(t, err)
	// The fallback single-shot PUT carries the whole object.
	eqBytes(t, f.assembled(), in)
	// The abandoned multipart session is aborted before falling back.
	if got, want := f.abortsSeen(), []string{"tok"}; !slices.Equal(got, want) {
		t.Errorf("aborts = %v, want %v", got, want)
	}
}

func TestUploadAlreadyExists(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.singleShotStatus = http.StatusConflict
	f.singleShotBody = `{"error_code":"ALREADY_EXISTS","message":"file already exists"}`
	e := shrunkEngine(t, f)

	_, err := e.upload(t.Context(), "/Volumes/c/s/v/exists.bin", bytes.NewReader(data(100)), WithOverwrite(false))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got error %v, want ErrAlreadyExists", err)
	}
}

// TestSingleShotProgressRewindsOnRetry: a retried single shot sends the body
// again from its start, so progress rewinds instead of counting the bytes
// twice, and Transferred never exceeds Total.
func TestSingleShotProgressRewindsOnRetry(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// The fake reads the whole body before it answers the first PUT with a 503.
	f.singleShotHook = func(attempt int) int {
		if attempt == 1 {
			return http.StatusServiceUnavailable
		}
		return 0
	}
	e := shrunkEngine(t, f)
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

	in := data(500) // below the 1024 single-shot threshold
	var mu sync.Mutex
	var reports []Progress
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", bytes.NewReader(in),
		WithProgress(func(p Progress) {
			mu.Lock()
			defer mu.Unlock()
			reports = append(reports, p)
		}))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if f.singleTry != 2 {
		t.Fatalf("single-shot attempts = %d, want 2 (one retry after the 503)", f.singleTry)
	}

	mu.Lock()
	defer mu.Unlock()
	rewound := false
	for i, p := range reports {
		if p.Transferred > p.Total {
			t.Errorf("report %d: Transferred %d exceeds Total %d", i, p.Transferred, p.Total)
		}
		if i > 0 && p.Transferred < reports[i-1].Transferred {
			rewound = true
		}
	}
	if !rewound {
		t.Error("progress never rewound for the retry")
	}
	if n := len(reports); n == 0 || reports[n-1].Transferred != int64(len(in)) {
		t.Errorf("last report = %v, want Transferred %d", reports, len(in))
	}
}

// TestSingleShotOutlastsTheRetryBudget: the retry budget stops new attempts
// and never cuts a running one short. A single shot whose bytes keep moving
// for longer than both the budget and the idle timeout lands in one attempt.
func TestSingleShotOutlastsTheRetryBudget(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.singleShotRead = func(r *http.Request) ([]byte, bool) {
		var body []byte
		buf := make([]byte, 32<<10)
		for {
			n, err := io.ReadFull(r.Body, buf)
			body = append(body, buf[:n]...)
			if err != nil {
				return body, true
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	// A small send buffer makes the client's writes wait on the server's
	// reads, so the client sees the bytes move at the server's pace instead of
	// handing them all to the kernel at once.
	transport := &http.Transport{DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		c, err := (&net.Dialer{}).DialContext(ctx, network, addr)
		if tc, ok := c.(*net.TCPConn); ok {
			_ = tc.SetWriteBuffer(16 << 10)
		}
		return c, err
	}}
	t.Cleanup(transport.CloseIdleConnections)
	e := newEngine(buildUploadClient(t, f.base(), &http.Client{Transport: transport}, ""))
	e.tun.minStreamSize = 8 << 20 // so 4 MiB is a single shot
	e.tun.cpRetryTimeout = 50 * time.Millisecond
	e.tun.transferIdleTimeout = 300 * time.Millisecond

	in := data(4 << 20)
	start := time.Now()
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/slow.bin", bytes.NewReader(in))
	took := time.Since(start)
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if f.singleTry != 1 {
		t.Errorf("single-shot attempts = %d, want 1", f.singleTry)
	}
	if took < e.tun.transferIdleTimeout {
		t.Fatalf("the upload took %v, under the idle timeout, so it shows nothing", took)
	}
}

// TestSingleShotStallIsCutAndRetried: an attempt that moves no bytes for the
// idle timeout fails as a timeout, which is retried as any transient failure
// is.
func TestSingleShotStallIsCutAndRetried(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// Registered after the server's Close, so it runs first and frees a held
	// answer.
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	var tries atomic.Int32
	f.singleShotRead = func(r *http.Request) ([]byte, bool) {
		// The server notices that the client gave up only once the body is read.
		body, _ := io.ReadAll(r.Body)
		if tries.Add(1) == 1 {
			// Hold the answer: no bytes move until the client gives up.
			select {
			case <-r.Context().Done():
			case <-release:
			}
			return nil, false
		}
		return body, true
	}
	e := shrunkEngine(t, f)
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.cpRetryTimeout = 5 * time.Second
	e.tun.transferIdleTimeout = 250 * time.Millisecond

	in := data(500)
	// Without the stall rule this would wait for the held answer; fail instead.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	start := time.Now()
	_, err := e.upload(ctx, "/Volumes/c/s/v/stall.bin", bytes.NewReader(in))
	took := time.Since(start)
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if n := tries.Load(); n != 2 {
		t.Errorf("single-shot attempts = %d, want 2 (the stall, then a retry)", n)
	}
	if took < e.tun.transferIdleTimeout || took > e.tun.cpRetryTimeout {
		t.Errorf("the upload took %v, want at least the %v idle timeout and well under the %v budget", took, e.tun.transferIdleTimeout, e.tun.cpRetryTimeout)
	}
}

// TestSingleShotRetriesStopAtTheBudget: no new attempt starts once the budget
// has passed. The upload then returns the deadline error, as ops.Execute does
// for any call, which the contract reads as RETRY_TIMEOUT.
func TestSingleShotRetriesStopAtTheBudget(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.singleShotHook = func(int) int { return http.StatusServiceUnavailable }
	e := shrunkEngine(t, f)
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: 5 * time.Millisecond, Maximum: 5 * time.Millisecond, Factor: 1}
	e.tun.cpRetryTimeout = 100 * time.Millisecond

	// Without the budget the workspace would be retried until ctx ends; fail
	// instead.
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	start := time.Now()
	_, err := e.upload(ctx, "/Volumes/c/s/v/p.bin", bytes.NewReader(data(500)))
	took := time.Since(start)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("got error %v, want the budget's deadline", err)
	}
	if f.singleTry < 2 {
		t.Errorf("single-shot attempts = %d, want retries within the budget", f.singleTry)
	}
	if took > time.Second {
		t.Errorf("the upload took %v, want it to stop at the %v budget", took, e.tun.cpRetryTimeout)
	}
}

// slowReader is a seekable reader whose every Read first waits, as a caller's
// slow source would. It has no ReadAt, so the upload reads through Read.
type slowReader struct {
	r     *bytes.Reader
	pause time.Duration
}

func (s *slowReader) Read(p []byte) (int, error) {
	time.Sleep(s.pause)
	return s.r.Read(p[:min(len(p), 100)])
}

func (s *slowReader) Seek(offset int64, whence int) (int64, error) {
	return s.r.Seek(offset, whence)
}

// TestSingleShotIdleTimeoutIgnoresTheCallersReader: time spent waiting on the
// caller's reader isn't a stall.
func TestSingleShotIdleTimeoutIgnoresTheCallersReader(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)
	e.tun.transferIdleTimeout = 200 * time.Millisecond
	e.tun.cpRetryTimeout = 10 * time.Millisecond

	in := data(500)
	start := time.Now()
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/slow-source.bin", &slowReader{r: bytes.NewReader(in), pause: 3 * e.tun.transferIdleTimeout / 2})
	took := time.Since(start)
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if f.singleTry != 1 {
		t.Errorf("single-shot attempts = %d, want 1", f.singleTry)
	}
	if took < 2*e.tun.transferIdleTimeout {
		t.Fatalf("the upload took %v, so the caller's reader never outlasted the idle timeout", took)
	}
}

func TestStalledSingleShotIsTransient(t *testing.T) {
	e := shrunkEngine(t, newFakeServer(t, "multipart"))
	err := &stalledError{idle: time.Minute}
	if _, ok := e.newControlPlaneRetrier().IsRetriable(err); !ok {
		t.Error("a stalled attempt isn't retried")
	}
	if !walksOn(t.Context(), err) {
		t.Error("a stalled attempt doesn't move on to the next endpoint")
	}
	if !errors.Is(err, errUploadStalled) {
		t.Error("a stalled attempt doesn't match errUploadStalled")
	}
}

// TestSingleShotRetriesWithProgress guards both MAJOR review findings: with
// WithProgress set, a seekable single-shot body must still be retried on a
// transient failure (the progress wrapper previously hid the io.Seeker, which
// silently disabled retry).
func TestSingleShotRetriesWithProgress(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// First single-shot PUT fails with a retryable 503; the second succeeds.
	f.singleShotHook = func(attempt int) int {
		if attempt == 1 {
			return http.StatusServiceUnavailable
		}
		return 0
	}
	e := shrunkEngine(t, f)
	// Shrink the control-plane backoff so the retry does not sleep a full second.
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

	var reported atomic.Int64
	in := data(500) // below the 1024 single-shot threshold
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", bytes.NewReader(in),
		WithProgress(func(p Progress) { reported.Store(p.Transferred) }))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)

	if f.singleTry != 2 {
		t.Errorf("single-shot attempts = %d, want 2 (one retry after the 503)", f.singleTry)
	}
	if reported.Load() == 0 {
		t.Error("progress callback never fired")
	}
}

func TestSingleShotRetriesAFileBody(t *testing.T) {
	// An *os.File body is an io.Closer, which net/http closes after each attempt.
	// The retry must still rewind it, and the caller's file must stay open.
	in := data(500) // below the 1024 single-shot threshold
	src := filepath.Join(t.TempDir(), "src.bin")
	if err := os.WriteFile(src, in, 0o600); err != nil {
		t.Fatal(err)
	}
	upload := map[string]func(*engine) error{
		"UploadFrom": func(e *engine) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/p.bin", src)
			return err
		},
		"Upload with the caller's file": func(e *engine) error {
			f, err := os.Open(src)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", f); err != nil {
				return err
			}
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("the caller's file was closed: %w", err)
			}
			return nil
		},
	}
	for name, run := range upload {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.singleShotHook = func(attempt int) int {
				if attempt == 1 {
					return http.StatusServiceUnavailable
				}
				return 0
			}
			e := shrunkEngine(t, f)
			e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

			noErr(t, run(e))
			eqBytes(t, f.assembled(), in)
			if f.singleTry != 2 {
				t.Errorf("single-shot attempts = %d, want 2 (one retry after the 503)", f.singleTry)
			}
		})
	}
}

func TestSingleShotRetryIsNotCutShortByALateRead(t *testing.T) {
	// net/http may still read a request body after the attempt returns, until
	// it closes it. This transport answers the first attempt with a 503 at once
	// and keeps reading that attempt's body while the retry runs, as a write
	// loop still sending it would. The retry must still send every byte, and
	// the late read mustn't count toward the retry's progress.
	in := data(500) // below the 1024 single-shot threshold
	src := filepath.Join(t.TempDir(), "src.bin")
	noErr(t, os.WriteFile(src, in, 0o600))
	uploads := map[string]struct {
		upload func(*engine, ...UploadOption) error
		// randomAccess is whether net/http can also send the body again itself.
		randomAccess bool
	}{
		"a file": {func(e *engine, opts ...UploadOption) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/p.bin", src, opts...)
			return err
		}, true},
		"a reader without random access": {func(e *engine, opts ...UploadOption) error {
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(in)}, opts...)
			return err
		}, false},
	}
	for name, tc := range uploads {
		t.Run(name, func(t *testing.T) {
			retrying, lateReadDone := make(chan struct{}), make(chan struct{})
			var attempts int
			var sent, resent []byte
			var sentLength int64
			rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
				attempts++
				if attempts == 1 {
					go func() {
						// Wait for the retry, but not forever: a retry that
						// waits for this attempt's body to close never starts.
						select {
						case <-retrying:
						case <-time.After(100 * time.Millisecond):
						}
						_, _ = io.CopyN(io.Discard, r.Body, 100)
						_ = r.Body.Close()
						close(lateReadDone)
					}()
					return cannedResponse(r, http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`), nil
				}
				if attempts == 2 {
					close(retrying)
				}
				<-lateReadDone
				var err error
				sent, err = io.ReadAll(r.Body)
				_ = r.Body.Close()
				sentLength = r.ContentLength
				if err != nil {
					return nil, err
				}
				if r.GetBody != nil {
					again, err := r.GetBody()
					if err != nil {
						return nil, err
					}
					resent, _ = io.ReadAll(again)
				}
				return cannedResponse(r, http.StatusNoContent, ""), nil
			})
			e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
			e.tun.minStreamSize = 1024
			e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

			var past int64 // the most Transferred ran past Total
			noErr(t, tc.upload(e, WithProgress(func(p Progress) { past = max(past, p.Transferred-p.Total) })))
			if attempts != 2 {
				t.Errorf("attempts = %d, want 2 (one retry after the 503)", attempts)
			}
			if past > 0 {
				t.Errorf("progress ran %d bytes past Total: the first attempt's late read counted", past)
			}
			eqBytes(t, sent, in)
			if sentLength != int64(len(in)) {
				t.Errorf("Content-Length = %d, want %d", sentLength, len(in))
			}
			if tc.randomAccess {
				eqBytes(t, resent, in)
			} else if resent != nil {
				t.Errorf("GetBody sent %d bytes, want no GetBody: the body cannot be read twice at once", len(resent))
			}
		})
	}
}

func TestWalkWaitsForTheLastBodyBeforeTheNextEndpoint(t *testing.T) {
	// RoundTrip may close a request body on another goroutine after it
	// returns. Here the storage proxy refuses the single shot at once and keeps
	// reading its body while the walk moves on, as a write loop still sending it
	// would. The workspace must still get exactly the source's bytes.
	in := data(500) // below the 1024 single-shot threshold
	src := filepath.Join(t.TempDir(), "src.bin")
	noErr(t, os.WriteFile(src, in, 0o600))
	sources := map[string]func(*engine) error{
		// It reads the source itself on every endpoint, so the walk must wait
		// for the proxy's body to close before it rewinds the source.
		"a reader that seeks but has no ReadAt": func(e *engine) error {
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(in)})
			return err
		},
		// Each attempt reads its own section, so a late read can't move the
		// next endpoint's position.
		"a file": func(e *engine) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/p.bin", src)
			return err
		},
	}
	for name, upload := range sources {
		t.Run(name, func(t *testing.T) {
			workspaceStarted, lateReadDone := make(chan struct{}), make(chan struct{})
			var proxyPuts, workspacePuts atomic.Int32
			var sent []byte
			rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch {
				case r.Method == http.MethodGet && r.URL.Path == endpointsAPIPath:
					return cannedResponse(r, http.StatusOK, `{"endpoints":[
						{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com"},
						{"api":"FILES_API","api_version":"2.0","scheme":"https","hostname":"ws.example.com"}]}`), nil
				case r.Method == http.MethodPut && r.URL.Host == "proxy.example.com":
					proxyPuts.Add(1)
					go func() {
						// Wait for the next endpoint, but not forever: a walk that
						// waits for this body to close never gets there.
						select {
						case <-workspaceStarted:
						case <-time.After(100 * time.Millisecond):
						}
						_, _ = io.CopyN(io.Discard, r.Body, 100)
						_ = r.Body.Close()
						close(lateReadDone)
					}()
					return cannedResponse(r, http.StatusBadRequest, `{"error_code":"BAD_REQUEST","message":"no"}`), nil
				case r.Method == http.MethodPut && r.URL.Host == "ws.example.com":
					if workspacePuts.Add(1) == 1 {
						close(workspaceStarted)
					}
					<-lateReadDone
					var err error
					sent, err = io.ReadAll(r.Body)
					_ = r.Body.Close()
					if err != nil {
						return nil, err
					}
					return cannedResponse(r, http.StatusNoContent, ""), nil
				}
				t.Errorf("unexpected request %s %s", r.Method, r.URL)
				return cannedResponse(r, http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`), nil
			})
			e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
			e.tun.endpointsAPI = true
			e.tun.minStreamSize = 1024
			e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

			noErr(t, upload(e))
			if p, w := proxyPuts.Load(), workspacePuts.Load(); p != 1 || w != 1 {
				t.Fatalf("PUTs: storage proxy %d, workspace %d; want 1 each, walking from the proxy to the workspace", p, w)
			}
			eqBytes(t, sent, in)
		})
	}
}

func TestWalkWaitingForTheLastBodyStopsWhenCancelled(t *testing.T) {
	// The proxy here never closes its body, so only the caller's cancellation
	// can end the walk's wait for it.
	var workspacePuts atomic.Int32
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == endpointsAPIPath:
			return cannedResponse(r, http.StatusOK, `{"endpoints":[
				{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com"},
				{"api":"FILES_API","api_version":"2.0","scheme":"https","hostname":"ws.example.com"}]}`), nil
		case r.Method == http.MethodPut && r.URL.Host == "proxy.example.com":
			return cannedResponse(r, http.StatusBadRequest, `{"error_code":"BAD_REQUEST","message":"no"}`), nil
		case r.Method == http.MethodPut && r.URL.Host == "ws.example.com":
			workspacePuts.Add(1)
			_ = r.Body.Close()
			return cannedResponse(r, http.StatusNoContent, ""), nil
		}
		t.Errorf("unexpected request %s %s", r.Method, r.URL)
		return cannedResponse(r, http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`), nil
	})
	e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
	e.tun.endpointsAPI = true
	e.tun.minStreamSize = 1024

	stop := errors.New("the caller gave up")
	ctx, cancel := context.WithCancelCause(t.Context())
	time.AfterFunc(50*time.Millisecond, func() { cancel(stop) })
	done := make(chan error, 1)
	go func() {
		_, err := e.upload(ctx, "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(data(500))})
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, stop) {
			t.Errorf("err = %v, want the cancellation cause %v", err, stop)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the walk still waits for a body that is never closed after its context was cancelled")
	}
	if n := workspacePuts.Load(); n != 0 {
		t.Errorf("workspace PUTs = %d, want none: the source was never free to rewind", n)
	}
}

// failOnCall is credentials whose nth AuthHeaders call fails.
type failOnCall struct {
	n     int32
	calls atomic.Int32
}

func (*failOnCall) Name() string { return "fail-on-call" }

func (c *failOnCall) AuthHeaders(context.Context) ([]auth.Header, error) {
	if c.calls.Add(1) == c.n {
		return nil, errors.New("no credentials this time")
	}
	return []auth.Header{{Key: "Authorization", Value: "Bearer test"}}, nil
}

func TestWalkDoesNotWaitForABodyThatWasNeverSent(t *testing.T) {
	// The proxy's request fails to build (its credentials fail), so net/http
	// never gets the body. newHTTPRequest closes it then, or the walk's wait
	// for it would only end with the caller's context.
	in := data(500)
	var sent []byte
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == endpointsAPIPath:
			return cannedResponse(r, http.StatusOK, `{"endpoints":[
				{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com"},
				{"api":"FILES_API","api_version":"2.0","scheme":"https","hostname":"ws.example.com"}]}`), nil
		case r.Method == http.MethodPut && r.URL.Host == "ws.example.com":
			var err error
			sent, err = io.ReadAll(r.Body)
			_ = r.Body.Close()
			if err != nil {
				return nil, err
			}
			return cannedResponse(r, http.StatusNoContent, ""), nil
		}
		t.Errorf("unexpected request %s %s", r.Method, r.URL)
		return cannedResponse(r, http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`), nil
	})
	e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
	e.c.credentials = &failOnCall{n: 2} // the first call is discovery's
	e.tun.endpointsAPI = true
	e.tun.minStreamSize = 1024

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	_, err := e.upload(ctx, "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(in)})
	noErr(t, err)
	eqBytes(t, sent, in)
}

func TestSingleShotRetryWaitingForTheLastBodyStopsWhenCancelled(t *testing.T) {
	// A retry of data without random access waits for net/http to close the
	// last attempt's body. This transport never closes it, so only the
	// caller's cancellation, or the retry budget, can end the wait.
	var attempts atomic.Int32
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts.Add(1)
		return cannedResponse(r, http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`), nil
	})
	e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
	e.tun.minStreamSize = 1024
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}

	stop := errors.New("the caller gave up")
	ctx, cancel := context.WithCancelCause(t.Context())
	time.AfterFunc(50*time.Millisecond, func() { cancel(stop) })
	done := make(chan error, 1)
	go func() {
		_, err := e.upload(ctx, "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(data(500))})
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, stop) {
			t.Errorf("err = %v, want the cancellation cause %v", err, stop)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the upload still waits for a body that is never closed after its context was cancelled")
	}
	if got := attempts.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1: the retry waits for the first attempt's body", got)
	}
}

func TestSingleShotRetryWaitingForTheLastBodyStopsAtTheBudget(t *testing.T) {
	// The same wait as above, with no cancellation: the budget, which bounds
	// starting a new attempt, ends it.
	var attempts atomic.Int32
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts.Add(1)
		return cannedResponse(r, http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`), nil
	})
	e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: rt}, ""))
	e.tun.minStreamSize = 1024
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.cpRetryTimeout = 50 * time.Millisecond

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()
	start := time.Now()
	_, err := e.upload(ctx, "/Volumes/c/s/v/p.bin", struct{ io.ReadSeeker }{bytes.NewReader(data(500))})
	if took := time.Since(start); took > time.Second {
		t.Errorf("the upload took %v, want it to stop at the %v budget", took, e.tun.cpRetryTimeout)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want the budget's deadline", err)
	}
	if got := attempts.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1: the retry waits for the first attempt's body", got)
	}
}

// flakyCredentials fails its first AuthHeaders call with a connection reset,
// which the control-plane retrier retries, and gives testToken after that.
type flakyCredentials struct{ calls atomic.Int32 }

func (*flakyCredentials) Name() string { return "flaky" }

func (c *flakyCredentials) AuthHeaders(context.Context) ([]auth.Header, error) {
	if c.calls.Add(1) == 1 {
		return nil, fmt.Errorf("refreshing the token: %w", syscall.ECONNRESET)
	}
	return []auth.Header{{Key: "Authorization", Value: testToken}}, nil
}

// resolveOnceCredentials resolves on its first AuthHeaders call, which takes
// delay, and keeps that resolution's result, as the default credentials do:
// a resolution its context ended fails every later call too.
type resolveOnceCredentials struct {
	delay time.Duration
	once  sync.Once
	err   error
}

func (*resolveOnceCredentials) Name() string { return "resolve-once" }

func (c *resolveOnceCredentials) AuthHeaders(ctx context.Context) ([]auth.Header, error) {
	c.once.Do(func() {
		select {
		case <-time.After(c.delay):
		case <-ctx.Done():
			c.err = fmt.Errorf("resolving the credentials: %w", context.Cause(ctx))
		}
	})
	if c.err != nil {
		return nil, c.err
	}
	return []auth.Header{{Key: "Authorization", Value: testToken}}, nil
}

func TestSingleShotRetriesAfterTheCredentialsFail(t *testing.T) {
	// An attempt whose credentials fail sends nothing, but its body must still
	// be closed: the retry of data without random access waits for that
	// before it rewinds.
	var got []byte
	e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})
	creds := &flakyCredentials{}
	e.c.credentials = creds
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	// A retry left waiting for that body would fail at this budget.
	e.tun.cpRetryTimeout = 2 * time.Second

	in := data(500)
	noErr(t, e.uploadSingleShot(t.Context(), controlContext(e, nil), struct{ io.ReadSeeker }{bytes.NewReader(in)}, int64(len(in))))
	eqBytes(t, got, in)
	if n := creds.calls.Load(); n != 2 {
		t.Errorf("credential calls = %d, want 2: the failed one, then the retry's", n)
	}
}

func TestSingleShotResolvesCredentialsOutsideTheIdleTimeout(t *testing.T) {
	// Credentials that resolve once, as the default ones do, keep the first
	// resolution's result. A first resolution slower than the idle timeout
	// moves no bytes, but it isn't a stall: the attempt, sent only once since
	// its body can't be rewound, must land, and so must the next upload.
	var mu sync.Mutex
	var got [][]byte
	e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		got = append(got, b)
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	})
	e.c.credentials = &resolveOnceCredentials{delay: 300 * time.Millisecond}
	e.tun.transferIdleTimeout = 100 * time.Millisecond

	in := data(500)
	for upload := range 2 {
		if err := e.uploadSingleShot(t.Context(), controlContext(e, nil), struct{ io.Reader }{bytes.NewReader(in)}, int64(len(in))); err != nil {
			t.Fatalf("upload %d: %v", upload+1, err)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if want := [][]byte{in, in}; !slices.EqualFunc(got, want, bytes.Equal) {
		t.Errorf("the host got %d bodies, want both uploads' %d bytes", len(got), len(in))
	}
}

func TestSingleShotCredentialResolutionIsBoundedByTheRetryBudget(t *testing.T) {
	// Resolving credentials isn't a stall, so the idle timeout doesn't bound
	// it, but the retry budget does: credentials that never resolve fail the
	// upload once the budget runs out, instead of holding it for good.
	e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request sent without resolved credentials: %s %s", r.Method, r.URL)
		w.WriteHeader(http.StatusOK)
	})
	e.c.credentials = &resolveOnceCredentials{delay: time.Hour}
	e.tun.cpRetryTimeout = 200 * time.Millisecond

	in := data(500)
	done := make(chan error, 1)
	go func() {
		done <- e.uploadSingleShot(t.Context(), controlContext(e, nil), struct{ io.Reader }{bytes.NewReader(in)}, int64(len(in)))
	}()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v, want the retry budget's deadline", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("the upload still waits on its credentials 5 s into a %v retry budget", e.tun.cpRetryTimeout)
	}
}

// lateCredentials returns testToken after delay whatever its context, as
// credentials waiting on another call's resolution or token refresh do.
type lateCredentials struct{ delay time.Duration }

func (lateCredentials) Name() string { return "late" }

func (c lateCredentials) AuthHeaders(context.Context) ([]auth.Header, error) {
	time.Sleep(c.delay)
	return []auth.Header{{Key: "Authorization", Value: testToken}}, nil
}

func TestSingleShotStartsNoAttemptOnceTheRetryBudgetIsSpent(t *testing.T) {
	// The credentials return only after the retry budget ran out. The upload
	// fails with the budget's deadline and sends nothing, as the budget stops
	// new attempts.
	var puts atomic.Int32
	e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		puts.Add(1)
		w.WriteHeader(http.StatusOK)
	})
	e.c.credentials = lateCredentials{delay: 300 * time.Millisecond}
	e.tun.cpRetryTimeout = 100 * time.Millisecond

	in := data(500)
	err := e.uploadSingleShot(t.Context(), controlContext(e, nil), struct{ io.Reader }{bytes.NewReader(in)}, int64(len(in)))
	type result struct {
		puts     int32
		deadline bool
	}
	if got, want := (result{puts.Load(), errors.Is(err, context.DeadlineExceeded)}), (result{0, true}); got != want {
		t.Errorf("got %+v (err %v), want %+v", got, err, want)
	}
}

// TestSingleShotProgressStartsAtZeroOnTheNextEndpoint: an endpoint the walk
// moves on to sends the whole upload again, so progress starts from zero there,
// not from what the last endpoint's attempts read.
func TestSingleShotProgressStartsAtZeroOnTheNextEndpoint(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)
	in := data(500)
	var reports []Progress
	uc := controlContext(e, nil)
	uc.progress = newProgressReporter(func(p Progress) { reports = append(reports, p) }, int64(len(in)))
	// The last endpoint's attempt read the whole body before it failed.
	uc.progress.addAttempt(uc.progress.restart(), int64(len(in)))

	noErr(t, e.uploadSingleShot(t.Context(), uc, bytes.NewReader(in), int64(len(in))))
	eqBytes(t, f.assembled(), in)
	for i, p := range reports {
		if p.Transferred > p.Total {
			t.Errorf("report %d: Transferred %d exceeds Total %d", i, p.Transferred, p.Total)
		}
	}
	if last := reports[len(reports)-1]; last.Transferred != int64(len(in)) {
		t.Errorf("last report = %+v, want Transferred %d", last, len(in))
	}
}

func TestUploadWalksPastAnEndpointThatNeverAnswers(t *testing.T) {
	// A discovered Files API host accepts initiate-upload and never answers.
	// Each attempt there ends at the control-call bound, so the upload moves on
	// to the next endpoint long before the retry budget runs out.
	f := newFakeServer(t, "multipart")
	proxy := strings.Replace(strings.TrimPrefix(f.base(), "http://"), "127.0.0.1", "localhost", 1)
	f.endpoints = fmt.Sprintf(`{"endpoints":[
		{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":%q,"multipart_upload_supported":true},
		{"api":"IO_OFFLOAD","scheme":"http","hostname":"$HOST","multipart_upload_supported":true}]}`, proxy)
	var hung atomic.Int32
	f.onRequest = func(r *http.Request) {
		if r.Host == proxy && r.URL.Query().Get("action") == "initiate-upload" {
			hung.Add(1)
			<-r.Context().Done()
		}
	}
	e := shrunkEngine(t, f)
	e.tun.endpointsAPI = true
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.controlAttemptTimeout = 200 * time.Millisecond
	// Without a bound per attempt, the first attempt would hold all of this.
	e.tun.cpRetryTimeout = 10 * time.Second

	in := data(3000)
	start := time.Now()
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", bytes.NewReader(in))
	took := time.Since(start)
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if n := hung.Load(); n != 2 {
		t.Errorf("unanswered initiates = %d, want 2: the attempt and its one retry before any part lands", n)
	}
	if took > 3*time.Second {
		t.Errorf("the upload took %v, want it to move on after two %v attempts", took, e.tun.controlAttemptTimeout)
	}
}

func TestUploadSkipsAnEndpointThatNeverAnswers(t *testing.T) {
	// A discovered Files API host accepts initiate-upload and never answers.
	// Each upload that walks past it to the next endpoint counts against it,
	// so after endpointSkipAfter uploads the next one skips it, and after the
	// cool-off one upload probes it again.
	f := newFakeServer(t, "multipart")
	proxy := strings.Replace(strings.TrimPrefix(f.base(), "http://"), "127.0.0.1", "localhost", 1)
	f.endpoints = fmt.Sprintf(`{"endpoints":[
		{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":%q,"multipart_upload_supported":true},
		{"api":"IO_OFFLOAD","scheme":"http","hostname":"$HOST","multipart_upload_supported":true}]}`, proxy)
	var hung atomic.Int32
	f.onRequest = func(r *http.Request) {
		if r.Host == proxy && r.URL.Query().Get("action") == "initiate-upload" {
			hung.Add(1)
			<-r.Context().Done()
		}
	}
	clock := newTestClock()
	e := shrunkEngine(t, f)
	e.now = clock.time
	e.tun.endpointsAPI = true
	e.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.controlAttemptTimeout = 100 * time.Millisecond
	e.tun.cpRetryTimeout = 10 * time.Second

	// triesProxy runs one upload, which must land, and tells whether it tried
	// the host that never answers.
	in := data(3000)
	triesProxy := func() bool {
		t.Helper()
		before := hung.Load()
		_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", bytes.NewReader(in))
		noErr(t, err)
		eqBytes(t, f.assembled(), in)
		return hung.Load() > before
	}
	var tried, want []bool
	for i := range e.tun.endpointSkipAfter + 1 {
		tried = append(tried, triesProxy())
		want = append(want, i < e.tun.endpointSkipAfter)
	}
	if !slices.Equal(tried, want) {
		t.Fatalf("tried the host = %v, want %v", tried, want)
	}
	clock.advance(e.tun.endpointSkipCoolOff)
	if !triesProxy() {
		t.Error("after the cool-off, the upload didn't probe the host")
	}
}

// TestUploadProgressRewindsWhenTheWalkMovesOn: a single shot that read the
// whole upload and failed leaves its bytes in the count; the multipart upload
// on the next endpoint sends them all again, so progress starts from zero.
func TestUploadProgressRewindsWhenTheWalkMovesOn(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.endpoints = `{"endpoints":[
		{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"$HOST"},
		{"api":"IO_OFFLOAD","scheme":"http","hostname":"$HOST","multipart_upload_supported":true}]}`
	// The fake reads the whole body before it answers the single shot.
	f.singleShotHook = func(int) int { return http.StatusInternalServerError }
	e := shrunkEngine(t, f)
	e.tun.endpointsAPI = true

	in := data(3000) // at the 1024 multipart threshold or more
	var mu sync.Mutex
	var reports []Progress
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/p.bin", bytes.NewReader(in),
		WithProgress(func(p Progress) {
			mu.Lock()
			defer mu.Unlock()
			reports = append(reports, p)
		}))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if f.singleTry != 1 || len(f.completed) != 3 {
		t.Fatalf("single shots = %d, parts = %d; want 1 single shot, then 3 parts", f.singleTry, len(f.completed))
	}

	mu.Lock()
	defer mu.Unlock()
	var most int64
	for i, p := range reports {
		most = max(most, p.Transferred)
		if p.Transferred > p.Total {
			t.Errorf("report %d: Transferred %d exceeds Total %d", i, p.Transferred, p.Total)
		}
	}
	if most != int64(len(in)) {
		t.Errorf("most Transferred = %d, want %d: the single shot read the whole upload", most, len(in))
	}
	if last := reports[len(reports)-1]; last.Transferred != int64(len(in)) {
		t.Errorf("last report = %+v, want Transferred %d", last, len(in))
	}
}

// closeRecorder is a one-pass body that records whether it was closed.
type closeRecorder struct {
	io.Reader
	closed atomic.Bool
}

func (r *closeRecorder) Close() error {
	r.closed.Store(true)
	return nil
}

func TestSingleShotDoesNotCloseANonSeekableBody(t *testing.T) {
	// net/http closes a request body that is an io.Closer. A body without Seek
	// is sent once and stays the caller's, so the PUT must not close it.
	var got []byte
	e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})
	body := &closeRecorder{Reader: strings.NewReader("hello")}
	noErr(t, e.uploadSingleShot(t.Context(), controlContext(e, nil), body, -1))
	if string(got) != "hello" {
		t.Errorf("sent %q, want %q", got, "hello")
	}
	if body.closed.Load() {
		t.Error("the PUT closed the caller's body")
	}
}

func TestSingleShotSendsContentLength(t *testing.T) {
	// Hosts require a Content-Length, which net/http can't tell for a file or
	// for the body of an attempt; the single shot sets it from the body's size.
	src := filepath.Join(t.TempDir(), "f.bin")
	noErr(t, os.WriteFile(src, []byte("hello"), 0o600))
	cases := map[string]struct {
		upload func(*engine) error
		size   int64
	}{
		"a file": {func(e *engine) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/f.bin", src)
			return err
		}, 5},
		"an empty body": {func(e *engine) error {
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/f.bin", bytes.NewReader(nil))
			return err
		}, 0},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var gotLength int64
			var gotEncoding []string
			e := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
				gotLength, gotEncoding = r.ContentLength, r.TransferEncoding
				_, _ = io.Copy(io.Discard, r.Body)
				w.WriteHeader(http.StatusOK)
			})
			noErr(t, tc.upload(e))
			if gotLength != tc.size || len(gotEncoding) != 0 {
				t.Errorf("Content-Length = %d, Transfer-Encoding = %v, want %d and none", gotLength, gotEncoding, tc.size)
			}
		})
	}
}

func cannedResponse(r *http.Request, status int, body string) *http.Response {
	return &http.Response{
		Status:     http.StatusText(status),
		StatusCode: status,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    r,
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestUploadWithTransferClient(t *testing.T) {
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)

	var used atomic.Int32
	custom := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		used.Add(1)
		return http.DefaultTransport.RoundTrip(r)
	})}

	in := data(5000)
	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/tc.bin", writeTemp(t, in),
		WithParallelism(2), WithTransferClient(custom))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	if used.Load() == 0 {
		t.Error("the cloud transfer must go through the supplied client")
	}
}

func TestUploadPresignedURLExpiryRetried(t *testing.T) {
	f := newFakeServer(t, "multipart")
	// Every part's first presigned URL is reported expired; the second succeeds.
	azure := `<Error><Code>AuthenticationFailed</Code><AuthenticationErrorDetail>Signature not valid in the specified time frame</AuthenticationErrorDetail></Error>`
	f.partHook = func(n, attempt int) (int, string, []byte) {
		if attempt == 1 {
			return http.StatusForbidden, "", []byte(azure)
		}
		return 0, "", nil
	}
	e := shrunkEngine(t, f)

	in := data(5000)
	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/exp.bin", writeTemp(t, in), WithParallelism(4))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

// Retry-timing coverage (retryable-status retried/exhausted, with shrunk backoff)
// lives in the cloudstorage package, where Send's retry policy is defined.

// walkEngine is a shrunkEngine whose retries back off in milliseconds.
func walkEngine(t *testing.T, f *fakeServer) *engine {
	t.Helper()
	e := shrunkEngine(t, f)
	fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	return e
}

func TestFailureBeforeAnyPartLandsFallsBackToSingleShot(t *testing.T) {
	// Before any part lands, a part that keeps failing gets one retry, and then
	// the whole upload is sent in a single shot; a transient failure is no
	// different from a cloud rejection.
	for _, parallelism := range []int{1, 4} {
		t.Run(strconv.Itoa(parallelism), func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.partHook = func(n, attempt int) (int, string, []byte) {
				return http.StatusServiceUnavailable, "", nil
			}
			e := walkEngine(t, f)

			in := data(5000)
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/fb.bin", writeTemp(t, in), WithParallelism(parallelism))
			noErr(t, err)
			eqBytes(t, f.assembled(), in)
			f.mu.Lock()
			defer f.mu.Unlock()
			for n, tries := range f.partTry {
				if tries > 2 {
					t.Errorf("part %d was attempted %d times, want at most 2", n, tries)
				}
			}
			if f.singleTry != 1 {
				t.Errorf("single-shot PUTs = %d, want 1", f.singleTry)
			}
		})
	}
}

func TestFailureAfterAPartLandsIsRetriedThenReturned(t *testing.T) {
	f := newFakeServer(t, "multipart")
	f.partHook = func(n, attempt int) (int, string, []byte) {
		if n == 1 {
			return 0, "", nil
		}
		return http.StatusServiceUnavailable, "", nil
	}
	e := walkEngine(t, f)

	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/c.bin", writeTemp(t, data(5000)), WithParallelism(1))
	if aerr, ok := errors.AsType[*apierr.APIError](err); !ok || aerr.HTTPStatusCode() != http.StatusServiceUnavailable {
		t.Fatalf("err = %v, want the last failure, a 503", err)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if got, want := f.partTry[2], e.tun.cloudRetries+1; got != want {
		t.Errorf("part 2 attempts = %d, want %d", got, want)
	}
	if f.singleTry != 0 {
		t.Errorf("single-shot PUTs = %d, want none once a part has landed", f.singleTry)
	}
}

func TestFailureThatTheSingleShotWouldHitTooIsReturned(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*fakeServer)
		check func(*testing.T, error)
	}{
		{
			name: "a part fails its precondition",
			setup: func(f *fakeServer) {
				f.partHook = func(n, attempt int) (int, string, []byte) {
					return http.StatusPreconditionFailed, "", nil
				}
			},
			check: func(t *testing.T, err error) {
				if aerr, ok := errors.AsType[*apierr.APIError](err); !ok || aerr.HTTPStatusCode() != http.StatusPreconditionFailed {
					t.Errorf("err = %v, want a 412", err)
				}
			},
		},
		{
			name: "initiate finds the file exists",
			setup: func(f *fakeServer) {
				f.initiateStatus = http.StatusConflict
				f.initiateBody = `{"error_code":"ALREADY_EXISTS","message":"exists"}`
			},
			check: func(t *testing.T, err error) {
				if !errors.Is(err, ErrAlreadyExists) {
					t.Errorf("err = %v, want ErrAlreadyExists", err)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			tc.setup(f)
			e := walkEngine(t, f)
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/x.bin", writeTemp(t, data(5000)), WithParallelism(1))
			tc.check(t, err)
			if f.singleTry != 0 {
				t.Errorf("single-shot PUTs = %d, want none", f.singleTry)
			}
		})
	}
}

func TestFailedInitiateFallsBackToSingleShot(t *testing.T) {
	sources := map[string]func([]byte) io.Reader{
		"seekable":     func(b []byte) io.Reader { return bytes.NewReader(b) },
		"non-seekable": func(b []byte) io.Reader { return io.MultiReader(bytes.NewReader(b)) },
	}
	for name, source := range sources {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.initiateStatus = http.StatusInternalServerError
			f.initiateBody = `{"error_code":"INTERNAL_ERROR","message":"boom"}`
			e := walkEngine(t, f)

			// Larger than the single-shot threshold, so the read-ahead of the
			// non-seekable stream is followed by the rest of it.
			in := data(5000)
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/i.bin", source(in), WithParallelism(1))
			noErr(t, err)
			eqBytes(t, f.assembled(), in)
		})
	}
}

func TestFailedPartURLMintFallsBackToSingleShot(t *testing.T) {
	// Each driver mints URLs its own way: the sequential loop per batch, the file
	// producer ahead of its workers, and the stream's workers on demand.
	uploads := map[string]func(*engine, []byte) error{
		"sequential": func(e *engine, in []byte) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/m.bin", writeTemp(t, in), WithParallelism(1))
			return err
		},
		"parallel file": func(e *engine, in []byte) error {
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/m.bin", writeTemp(t, in), WithParallelism(4))
			return err
		},
		"parallel stream": func(e *engine, in []byte) error {
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/m.bin", io.MultiReader(bytes.NewReader(in)), WithParallelism(4))
			return err
		},
	}
	for name, upload := range uploads {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.partURLsStatus = http.StatusInternalServerError
			f.partURLsBody = `{"error_code":"INTERNAL_ERROR","message":"boom"}`
			e := walkEngine(t, f)

			in := data(5000)
			noErr(t, upload(e, in))
			eqBytes(t, f.assembled(), in)
			if f.singleTry != 1 {
				t.Errorf("single-shot PUTs = %d, want 1", f.singleTry)
			}
		})
	}

	t.Run("the file exists", func(t *testing.T) {
		f := newFakeServer(t, "multipart")
		f.partURLsStatus = http.StatusConflict
		f.partURLsBody = `{"error_code":"ALREADY_EXISTS","message":"exists"}`
		e := walkEngine(t, f)

		_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/m.bin", writeTemp(t, data(5000)), WithParallelism(1))
		if apierr.Code(err) != codes.AlreadyExists {
			t.Errorf("err = %v, want ALREADY_EXISTS", err)
		}
		if f.singleTry != 0 {
			t.Errorf("single-shot PUTs = %d, want none", f.singleTry)
		}
	})
}

func TestControlPlaneFailureMovesOnAfterOneRetry(t *testing.T) {
	// Before any part lands, a control-plane call that keeps failing gets one
	// retry, and then the upload moves on to the single-shot PUT: as on any host,
	// moving on is cheaper than waiting out the retry timeout.
	cases := map[string]struct {
		fault    func(*fakeServer)
		requests func(*fakeServer) int
	}{
		"initiate": {
			fault:    func(f *fakeServer) { f.initiateStatus = http.StatusServiceUnavailable },
			requests: func(f *fakeServer) int { return len(f.overwriteQ) - f.singleTry },
		},
		"part URLs": {
			fault:    func(f *fakeServer) { f.partURLsStatus = http.StatusServiceUnavailable },
			requests: func(f *fakeServer) int { return len(f.partURLReqs) },
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			tc.fault(f)
			e := walkEngine(t, f)

			in := data(5000)
			_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/cp.bin", writeTemp(t, in), WithParallelism(1))
			noErr(t, err)
			eqBytes(t, f.assembled(), in)
			f.mu.Lock()
			defer f.mu.Unlock()
			if got := tc.requests(f); got != 2 {
				t.Errorf("%s requests = %d, want 2", name, got)
			}
		})
	}
}

func TestStorageTimeoutFallsBackToSingleShot(t *testing.T) {
	// A blackholed storage host times out rather than refusing the connection;
	// that is the case the fallback exists for, so a deadline error from a part
	// PUT falls back like any other failure.
	f := newFakeServer(t, "multipart")
	e := walkEngine(t, f)
	blackholed := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("blackholed: %w", context.DeadlineExceeded)
	})}

	in := data(5000)
	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/t.bin", writeTemp(t, in), WithParallelism(1), WithTransferClient(blackholed))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

func TestSingleShotFallbackOfAPipe(t *testing.T) {
	// A pipe is an *os.File, so it has Seek, but Seek fails: the fallback sends
	// the bytes read so far and then the rest of the pipe.
	f := newFakeServer(t, "multipart")
	f.partHook = func(n, attempt int) (int, string, []byte) {
		return http.StatusForbidden, "", []byte("blocked")
	}
	e := walkEngine(t, f)

	in := data(5000)
	pr, pw, err := os.Pipe()
	noErr(t, err)
	t.Cleanup(func() { _ = pr.Close() })
	go func() {
		_, _ = pw.Write(in)
		_ = pw.Close()
	}()
	_, err = e.upload(t.Context(), "/Volumes/c/s/v/p.bin", pr, WithParallelism(1))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
}

func TestSingleShotFallbackRetriesOnlyASeekableSource(t *testing.T) {
	// The fallback rewinds a seekable source and sends it itself, so the PUT is
	// retried; a one-pass stream has been read past the bytes it holds, so the
	// PUT is sent once.
	cases := []struct {
		name       string
		source     func([]byte) io.Reader
		wantTries  int
		wantFailed bool
	}{
		{"seekable", func(b []byte) io.Reader { return bytes.NewReader(b) }, 2, false},
		{"non-seekable", func(b []byte) io.Reader { return io.MultiReader(bytes.NewReader(b)) }, 1, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeServer(t, "multipart")
			f.partHook = func(n, attempt int) (int, string, []byte) {
				return http.StatusForbidden, "", []byte("blocked")
			}
			f.singleShotHook = func(attempt int) int {
				if attempt == 1 {
					return http.StatusServiceUnavailable
				}
				return 0
			}
			e := walkEngine(t, f)

			in := data(5000)
			_, err := e.upload(t.Context(), "/Volumes/c/s/v/s.bin", tc.source(in), WithParallelism(1))
			if tc.wantFailed != (err != nil) {
				t.Fatalf("err = %v, want failed = %v", err, tc.wantFailed)
			}
			if !tc.wantFailed {
				eqBytes(t, f.assembled(), in)
			}
			if f.singleTry != tc.wantTries {
				t.Errorf("single-shot PUTs = %d, want %d", f.singleTry, tc.wantTries)
			}
		})
	}
}

func TestResolvedParallelism(t *testing.T) {
	n := func(v int) *int { return &v }
	const mib int64 = 1 << 20
	tun := defaultTunables()
	tests := []struct {
		name     string
		cfg      uploadConfig
		buffered bool
		partSize int64
		want     int
	}{
		{"explicit override wins for a file", uploadConfig{parallelism: n(7)}, false, 16 * mib, 7},
		{"explicit override wins for a stream", uploadConfig{parallelism: n(3)}, true, 16 * mib, 3},
		{"file default is bandwidth-oriented", uploadConfig{}, false, 16 * mib, tun.fileParallelism},
		{"stream default fits the memory budget", uploadConfig{}, true, 16 * mib, int(tun.streamMemoryBudget / (16 * mib))},
		{"larger stream part size lowers the worker count", uploadConfig{}, true, 64 * mib, int(tun.streamMemoryBudget / (64 * mib))},
		{"tiny stream part size is capped at the file default", uploadConfig{}, true, 1, tun.fileParallelism},
		{"huge stream part size floors at one worker", uploadConfig{}, true, tun.streamMemoryBudget * 2, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.resolvedParallelism(tc.buffered, tc.partSize, tun); got != tc.want {
				t.Errorf("resolvedParallelism(%v, %d) = %d, want %d", tc.buffered, tc.partSize, got, tc.want)
			}
		})
	}
}

func TestUploadRejectsResumableSession(t *testing.T) {
	// Only multipart sessions are supported. An initiate response with a
	// resumable session and no multipart one is unexpected: no data goes to the
	// session, and the upload moves on to the single-shot PUT, as for any other
	// failure before a part lands.
	f := newFakeServer(t, "resumable")
	e := shrunkEngine(t, f)
	if _, err := e.initiateMultipart(t.Context(), controlContext(e, nil)); !errors.Is(err, errUnexpectedServerResponse) {
		t.Fatalf("initiateMultipart err = %v, want errUnexpectedServerResponse", err)
	}

	in := data(4 * 1024)
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/r.bin", bytes.NewReader(in))
	noErr(t, err)
	eqBytes(t, f.assembled(), in)
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.partTry) != 0 {
		t.Errorf("parts sent = %v, want none", f.partTry)
	}
}

func TestFailedMultipartUploadIsAbortedOnce(t *testing.T) {
	// A part fails with a status that isn't retried after another part has
	// landed, so the upload fails. Its session is aborted once, and a failing
	// abort is not retried and does not replace the upload's own error.
	f := newFakeServer(t, "multipart")
	f.partHook = func(n, attempt int) (int, string, []byte) {
		if n == 1 {
			return 0, "", nil
		}
		return http.StatusInternalServerError, "", []byte("boom")
	}
	f.abortStatus = http.StatusServiceUnavailable
	e := shrunkEngine(t, f)

	_, err := e.uploadFrom(t.Context(), "/Volumes/c/s/v/fail.bin", writeTemp(t, data(5000)), WithParallelism(1))
	if err == nil {
		t.Fatal("expected the upload to fail")
	}
	if apierr.Code(err) != codes.Internal {
		t.Errorf("err = %v, want the part's 500", err)
	}
	if got, want := f.abortsSeen(), []string{"tok"}; !slices.Equal(got, want) {
		t.Errorf("aborts = %v, want %v", got, want)
	}
}

func TestParallelShortStreamIsSentInASingleShot(t *testing.T) {
	// A stream of unknown size that turns out shorter than the single-shot
	// threshold is sent in one PUT, without opening a multipart session, even
	// when the upload is parallel.
	for _, size := range []int{0, 100} {
		f := newFakeServer(t, "multipart")
		e := shrunkEngine(t, f)
		in := data(size)
		_, err := e.upload(t.Context(), "/Volumes/c/s/v/short.bin", struct{ io.Reader }{bytes.NewReader(in)}, WithParallelism(4))
		noErr(t, err)
		eqBytes(t, f.assembled(), in)
		f.mu.Lock()
		singleShots, initiates := f.singleTry, len(f.overwriteQ)-f.singleTry
		f.mu.Unlock()
		if singleShots != 1 || initiates != 0 {
			t.Errorf("size %d: %d single shots and %d initiates, want 1 and 0", size, singleShots, initiates)
		}
	}
}

func TestParallelStreamAroundTheSingleShotThreshold(t *testing.T) {
	// The read-ahead decides between one PUT and a multipart session at exactly
	// the threshold (1024 in the shrunk engine). With a part size that doesn't
	// divide it, a part also straddles the end of the read-ahead, so the stream
	// must be reassembled without losing or repeating its bytes.
	for _, tc := range []struct {
		size      int
		multipart bool
	}{
		{1023, false},
		{1024, true},
		{1025, true},
		{5000, true},
	} {
		f := newFakeServer(t, "multipart")
		e := shrunkEngine(t, f)
		in := data(tc.size)
		_, err := e.upload(t.Context(), "/Volumes/c/s/v/edge.bin", struct{ io.Reader }{bytes.NewReader(in)},
			WithParallelism(4), WithPartSize(700))
		noErr(t, err)
		eqBytes(t, f.assembled(), in)
		f.mu.Lock()
		singleShots, initiates := f.singleTry, len(f.overwriteQ)-f.singleTry
		f.mu.Unlock()
		if gotMultipart := initiates == 1 && singleShots == 0; gotMultipart != tc.multipart {
			t.Errorf("size %d: %d initiates and %d single shots, want multipart=%v", tc.size, initiates, singleShots, tc.multipart)
		}
	}
}

// failingAfterReader returns n bytes and then err.
type failingAfterReader struct {
	n   int
	err error
}

func (r *failingAfterReader) Read(p []byte) (int, error) {
	if r.n == 0 {
		return 0, r.err
	}
	k := min(len(p), r.n)
	r.n -= k
	return k, nil
}

func TestParallelStreamReadErrorDuringReadAhead(t *testing.T) {
	// A read error before the threshold is reached fails the upload before any
	// request, so there is no session to abort.
	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)
	readErr := errors.New("disk on fire")
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/err.bin", &failingAfterReader{n: 500, err: readErr}, WithParallelism(4))
	if !errors.Is(err, readErr) {
		t.Fatalf("err = %v, want the read error", err)
	}
	f.mu.Lock()
	requests := len(f.overwriteQ) + len(f.aborts)
	f.mu.Unlock()
	if requests != 0 {
		t.Errorf("%d initiate, single-shot or abort requests, want none", requests)
	}
}

// userAgents records the User-Agent of each request, by method and path.
type userAgents struct {
	mu     sync.Mutex
	byPath map[string]string
}

func (u *userAgents) record(r *http.Request) {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.byPath == nil {
		u.byPath = map[string]string{}
	}
	u.byPath[r.Method+" "+r.URL.Path] = r.UserAgent()
}

// check fails unless every Files API request sent c's User-Agent, as the
// generated client does, and no presigned one did: those keep the transport's.
func (u *userAgents) check(t *testing.T, c *Client, presigned func(path string) bool) {
	t.Helper()
	want, err := c.userAgent()
	noErr(t, err)
	if want == "" {
		t.Fatal("the client has no User-Agent")
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	var api, cloud int
	for req, got := range u.byPath {
		_, path, _ := strings.Cut(req, " ")
		if presigned(path) {
			cloud++
			if got == want {
				t.Errorf("%s: sent the SDK's User-Agent to a presigned URL", req)
			}
			continue
		}
		api++
		if got != want {
			t.Errorf("%s: User-Agent = %q, want the SDK's %q", req, got, want)
		}
	}
	if api == 0 || cloud == 0 {
		t.Fatalf("requests = %v, want Files API and presigned ones", u.byPath)
	}
}

func TestUpload_sendsTheSDKUserAgent(t *testing.T) {
	// Requests the generated client doesn't send still carry the SDK's
	// User-Agent.
	f := newFakeServer(t, "multipart")
	var uas userAgents
	f.onRequest = uas.record
	e := shrunkEngine(t, f)
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/multipart.bin", bytes.NewReader(data(3000)))
	noErr(t, err)
	_, err = e.upload(t.Context(), "/Volumes/c/s/v/single-shot.bin", bytes.NewReader(data(500)))
	noErr(t, err)
	uas.check(t, e.c, func(path string) bool { return strings.HasPrefix(path, "/cloud/") })
}

func TestRequestsGoWithoutAUserAgentThatCantBeBuilt(t *testing.T) {
	// The User-Agent is only for attribution, so a client that can't build
	// one still uploads and discovers endpoints, and says why at debug.
	noUserAgent := func() (string, error) { return "", errors.New("no version") }
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))

	f := newFakeServer(t, "multipart")
	e := shrunkEngine(t, f)
	e.c.logger, e.c.userAgent = logger, noUserAgent
	_, err := e.upload(t.Context(), "/Volumes/c/s/v/multipart.bin", bytes.NewReader(data(3000)))
	noErr(t, err)
	_, err = e.upload(t.Context(), "/Volumes/c/s/v/single-shot.bin", bytes.NewReader(data(500)))
	noErr(t, err)
	eqBytes(t, f.assembled(), data(500))

	d, _ := discoveryEngine(t, respondJSON(http.StatusOK, `{"endpoints":[{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com"}]}`))
	d.c.logger, d.c.userAgent = logger, noUserAgent
	if got, want := d.transferEndpoints(t.Context()), []endpoint{{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com"}}; !slices.Equal(got, want) {
		t.Errorf("endpoints = %v, want the workspace's %v", got, want)
	}
	if !strings.Contains(logs.String(), "sending without the SDK's User-Agent") {
		t.Errorf("no debug line for the missing User-Agent; logs:\n%s", logs.String())
	}
}
