package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/options/client"
)

// discoveryEngine returns an engine with discovery on whose workspace answers
// GET /api/2.0/fs/endpoints with h, and counts the requests it gets.
func discoveryEngine(t *testing.T, h http.HandlerFunc) (*engine, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != endpointsAPIPath {
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
			return
		}
		calls.Add(1)
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
	e.tun.endpointsAPI = true
	return e, &calls
}

func respondJSON(status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}
}

func TestDiscoveryFallsBackToTheBuiltinOrder(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"not served":           respondJSON(http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`),
		"a server error":       respondJSON(http.StatusInternalServerError, `{"error_code":"INTERNAL_ERROR","message":"boom"}`),
		"not JSON":             respondJSON(http.StatusOK, `<html>`),
		"no endpoints":         respondJSON(http.StatusOK, `{}`),
		"endpoints not a list": respondJSON(http.StatusOK, `{"endpoints":{"api":"IO_OFFLOAD"}}`),
		"only unusable entries": respondJSON(http.StatusOK, `{"endpoints":[
			{"api":"NEW_API","scheme":"https","hostname":"new.example.com"},
			{"api":"FILES_API","api_version":"2.0","scheme":"ftp","hostname":"proxy.example.com"}]}`),
	}
	for name, h := range cases {
		t.Run(name, func(t *testing.T) {
			e, calls := discoveryEngine(t, h)
			for range 2 {
				if got, want := e.transferEndpoints(t.Context()), e.builtinEndpoints(); !slices.Equal(got, want) {
					t.Errorf("endpoints = %v, want the built-in %v", got, want)
				}
			}
			// The failure is cached too: the workspace is asked once.
			if calls.Load() != 1 {
				t.Errorf("endpoints requests = %d, want 1", calls.Load())
			}
		})
	}
}

func TestDiscoveryIsOffByDefault(t *testing.T) {
	e, calls := discoveryEngine(t, respondJSON(http.StatusOK, `{"endpoints":[]}`))
	e.tun = defaultTunables()
	if got, want := e.transferEndpoints(t.Context()), e.builtinEndpoints(); !slices.Equal(got, want) {
		t.Errorf("endpoints = %v, want the built-in %v", got, want)
	}
	if calls.Load() != 0 {
		t.Errorf("endpoints requests = %d, want none", calls.Load())
	}
}

func TestDiscoveryKeepsTheUsableEntriesInOrder(t *testing.T) {
	e, calls := discoveryEngine(t, respondJSON(http.StatusOK, `{"endpoints":[
		{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com","multipart_upload_supported":true},
		{"api":"NEW_API","scheme":"https","hostname":"new.example.com"},
		{"api":"IO_OFFLOAD","scheme":"https","hostname":"ws.example.com","multipart_upload_supported":true},
		{"api":"FILES_API","api_version":"2.0","scheme":"https","hostname":"ws.example.com"}]}`))
	want := []endpoint{
		{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com", multipart: true},
		{kind: presignedEndpoint, baseURL: "https://ws.example.com", multipart: true},
		{kind: filesAPIEndpoint, baseURL: "https://ws.example.com"},
	}
	for range 2 {
		if got := e.transferEndpoints(t.Context()); !slices.Equal(got, want) {
			t.Errorf("endpoints = %v, want %v", got, want)
		}
	}
	if calls.Load() != 1 {
		t.Errorf("endpoints requests = %d, want 1", calls.Load())
	}
}

func TestDiscoveryCancelledByTheCallerIsNotCached(t *testing.T) {
	e, calls := discoveryEngine(t, respondJSON(http.StatusOK, `{"endpoints":[
		{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com"}]}`))
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if got, want := e.transferEndpoints(ctx), e.builtinEndpoints(); !slices.Equal(got, want) {
		t.Errorf("endpoints = %v, want the built-in %v", got, want)
	}
	want := []endpoint{{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com"}}
	if got := e.transferEndpoints(t.Context()); !slices.Equal(got, want) {
		t.Errorf("endpoints = %v, want %v", got, want)
	}
	if calls.Load() != 1 {
		t.Errorf("endpoints requests = %d, want 1", calls.Load())
	}
}

func TestUploadEndpointsAddTheWorkspaceToPresignedOnly(t *testing.T) {
	e, _ := discoveryEngine(t, respondJSON(http.StatusOK, `{"endpoints":[
		{"api":"IO_OFFLOAD","scheme":"https","hostname":"ws.example.com","multipart_upload_supported":true}]}`))
	want := []endpoint{
		{kind: presignedEndpoint, baseURL: "https://ws.example.com", multipart: true},
		{kind: filesAPIEndpoint, baseURL: e.c.host},
	}
	if got := e.uploadEndpoints(t.Context()); !slices.Equal(got, want) {
		t.Errorf("upload endpoints = %v, want %v", got, want)
	}
}

func TestEndpointKindOf(t *testing.T) {
	tests := []struct {
		name   string
		entry  map[string]any
		want   endpointKind
		wantOK bool
	}{
		{"IO_OFFLOAD", map[string]any{"api": "IO_OFFLOAD"}, presignedEndpoint, true},
		{"IO_OFFLOAD with a null version", map[string]any{"api": "IO_OFFLOAD", "api_version": nil}, presignedEndpoint, true},
		{"a versioned IO_OFFLOAD is newer", map[string]any{"api": "IO_OFFLOAD", "api_version": "2.0"}, 0, false},
		{"FILES_API 2.0", map[string]any{"api": "FILES_API", "api_version": "2.0"}, filesAPIEndpoint, true},
		{"FILES_API 3.0", map[string]any{"api": "FILES_API", "api_version": "3.0"}, filesAPIv3Endpoint, true},
		{"FILES_API without a version", map[string]any{"api": "FILES_API"}, 0, false},
		{"a newer FILES_API", map[string]any{"api": "FILES_API", "api_version": "9.0"}, 0, false},
		{"a version that is not a string", map[string]any{"api": "FILES_API", "api_version": 2.0}, 0, false},
		{"an api that is not a string", map[string]any{"api": []any{"FILES_API"}, "api_version": "2.0"}, 0, false},
		{"not an object", nil, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := endpointKindOf(tc.entry)
			if got != tc.want || ok != tc.wantOK {
				t.Errorf("endpointKindOf = (%v, %v), want (%v, %v)", got, ok, tc.want, tc.wantOK)
			}
		})
	}
}

func TestParseEndpoint(t *testing.T) {
	entry := func(scheme, hostname string) map[string]any {
		return map[string]any{"scheme": scheme, "hostname": hostname, "multipart_upload_supported": true}
	}
	tests := []struct {
		name          string
		workspaceHost string
		kind          endpointKind
		entry         map[string]any
		wantBase      string // "" when the entry is unusable
	}{
		{"https", "https://ws.example.com", presignedEndpoint, entry("https", "ws.example.com"), "https://ws.example.com"},
		{"with a port", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com:8443"), "https://proxy.example.com:8443"},
		{"an IPv6 host", "https://ws.example.com", filesAPIEndpoint, entry("https", "[::1]:8443"), "https://[::1]:8443"},
		{"plain http to the storage proxy", "https://ws.example.com", filesAPIEndpoint, entry("http", "proxy.example.com"), "http://proxy.example.com"},
		{"plain http to a workspace reached that way", "http://localhost", presignedEndpoint, entry("http", "localhost"), "http://localhost"},
		{"plain http to a workspace reached over https", "https://ws.example.com", filesAPIEndpoint, entry("http", "WS.example.com"), ""},
		{"plain http presigned to another host", "https://ws.example.com", presignedEndpoint, entry("http", "other.example.com"), ""},
		{"Files API 3.0 over https", "https://ws.example.com", filesAPIv3Endpoint, entry("https", "rig.example.com"), "https://rig.example.com"},
		{"Files API 3.0 over plain http", "https://ws.example.com", filesAPIv3Endpoint, entry("http", "rig.example.com"), ""},
		{"Files API 3.0 over plain http to a workspace reached that way", "http://localhost", filesAPIv3Endpoint, entry("http", "localhost"), ""},
		{"another scheme", "https://ws.example.com", filesAPIEndpoint, entry("ftp", "proxy.example.com"), ""},
		{"no hostname", "https://ws.example.com", filesAPIEndpoint, entry("https", ""), ""},
		{"a hostname that is not a string", "https://ws.example.com", filesAPIEndpoint, map[string]any{"scheme": "https", "hostname": 1}, ""},
		{"a path", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com/files"), ""},
		{"a query", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com?x=1"), ""},
		{"user info", "https://ws.example.com", filesAPIEndpoint, entry("https", "user@proxy.example.com"), ""},
		{"an empty port", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com:"), ""},
		{"a port out of range", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com:70000"), ""},
		{"a port that is not a number", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy.example.com:https"), ""},
		{"a space", "https://ws.example.com", filesAPIEndpoint, entry("https", "proxy example.com"), ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEngine(&Client{internalClient: internalClient{host: tc.workspaceHost}})
			ep, ok := e.parseEndpoint(tc.entry, tc.kind)
			if tc.wantBase == "" {
				if ok {
					t.Errorf("parseEndpoint = %v, want the entry rejected", ep)
				}
				return
			}
			if !ok || ep.baseURL != tc.wantBase || ep.kind != tc.kind || !ep.multipart {
				t.Errorf("parseEndpoint = (%v, %v), want %s with multipart", ep, ok, tc.wantBase)
			}
		})
	}
}

func TestClient_uploadsReadEndpointsAPIEnv(t *testing.T) {
	// The workspace lists the storage proxy, which is this test server under
	// another name, so the host an upload lands on shows whether it asked.
	cases := []struct {
		value     string
		wantProxy bool
	}{
		{"", true},
		{"true", true},
		{"1", true},
		{"false", false},
		{"0", false},
		{"yes", true},
	}
	uploads := map[string]func(t *testing.T, c *Client) error{
		"Upload": func(t *testing.T, c *Client) error {
			_, err := c.Upload(t.Context(), "/Volumes/c/s/v/f", strings.NewReader("hello"))
			return err
		},
		"UploadFrom": func(t *testing.T, c *Client) error {
			src := filepath.Join(t.TempDir(), "f")
			noErr(t, os.WriteFile(src, []byte("hello"), 0o600))
			_, err := c.UploadFrom(t.Context(), "/Volumes/c/s/v/f", src)
			return err
		},
	}
	for name, upload := range uploads {
		for _, tc := range cases {
			t.Run(name+"/"+strconv.Quote(tc.value), func(t *testing.T) {
				t.Setenv(endpointsAPIEnv, tc.value)
				var srv *httptest.Server
				var listed atomic.Int32
				var putHost atomic.Value
				srv = httptest.NewServer(requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
					switch {
					case r.Method == http.MethodGet && r.URL.Path == endpointsAPIPath:
						listed.Add(1)
						proxy := strings.Replace(strings.TrimPrefix(srv.URL, "http://"), "127.0.0.1", "localhost", 1)
						respondJSON(http.StatusOK, fmt.Sprintf(`{"endpoints":[
						{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":%q}]}`, proxy))(w, r)
					case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/2.0/fs/files/"):
						_, _ = io.Copy(io.Discard, r.Body)
						putHost.Store(r.Host)
						w.WriteHeader(http.StatusNoContent)
					default:
						t.Errorf("unexpected request %s %s", r.Method, r.URL)
						w.WriteHeader(http.StatusInternalServerError)
					}
				}))
				t.Cleanup(srv.Close)
				var logs bytes.Buffer
				c, err := NewClient(t.Context(),
					client.WithHost(srv.URL),
					client.WithHTTPClient(srv.Client()),
					client.WithCredentials(testCredentials{}),
					client.WithLogger(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))),
					client.WithoutConfigFile(),
					client.WithoutEnv(),
				)
				noErr(t, err)

				noErr(t, upload(t, c))

				wantListed, wantHost := int32(0), strings.TrimPrefix(srv.URL, "http://")
				if tc.wantProxy {
					wantListed, wantHost = 1, strings.Replace(wantHost, "127.0.0.1", "localhost", 1)
				}
				if got := listed.Load(); got != wantListed {
					t.Errorf("endpoints requests = %d, want %d", got, wantListed)
				}
				if got, _ := putHost.Load().(string); got != wantHost {
					t.Errorf("upload went to %q, want %q", got, wantHost)
				}
				// The debug line is how an end-to-end run confirms that it used
				// the workspace's list.
				if got := strings.Contains(logs.String(), "using the workspace's file endpoints"); got != tc.wantProxy {
					t.Errorf("logged the workspace's endpoints = %v, want %v; logs:\n%s", got, tc.wantProxy, logs.String())
				}
			})
		}
	}
}

func TestClient_uploadsUseTheBuiltinOrderWhenTheWorkspaceCantList(t *testing.T) {
	// Discovery is on by default. When the workspace doesn't give a usable list,
	// an upload walks the built-in order: presigned URLs can't take a single
	// shot, so a small file goes to the workspace in one request.
	answers := map[string]http.HandlerFunc{
		"404":                       respondJSON(http.StatusNotFound, `{"error_code":"ENDPOINT_NOT_FOUND","message":"no"}`),
		"403":                       respondJSON(http.StatusForbidden, `{"error_code":"PERMISSION_DENIED","message":"no"}`),
		"503":                       respondJSON(http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`),
		"a body that doesn't parse": respondJSON(http.StatusOK, `{"endpoints":`),
		"a dropped connection": func(w http.ResponseWriter, r *http.Request) {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err != nil {
				t.Errorf("hijack: %v", err)
				return
			}
			_ = conn.Close()
		},
	}
	for name, answer := range answers {
		t.Run(name, func(t *testing.T) {
			t.Setenv(endpointsAPIEnv, "")
			var listed atomic.Int32
			var putHost atomic.Value
			srv := httptest.NewServer(requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.Method == http.MethodGet && r.URL.Path == endpointsAPIPath:
					listed.Add(1)
					answer(w, r)
				case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/2.0/fs/files/"):
					_, _ = io.Copy(io.Discard, r.Body)
					putHost.Store(r.Host)
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}))
			t.Cleanup(srv.Close)
			c, err := NewClient(t.Context(),
				client.WithHost(srv.URL),
				client.WithHTTPClient(srv.Client()),
				client.WithCredentials(testCredentials{}),
				client.WithoutConfigFile(),
				client.WithoutEnv(),
			)
			noErr(t, err)

			_, err = c.Upload(t.Context(), "/Volumes/c/s/v/f", strings.NewReader("hello"))
			noErr(t, err)
			if got := listed.Load(); got != 1 {
				t.Errorf("endpoints requests = %d, want 1", got)
			}
			if got, _ := putHost.Load().(string); got != strings.TrimPrefix(srv.URL, "http://") {
				t.Errorf("upload went to %q, want the workspace %q", got, strings.TrimPrefix(srv.URL, "http://"))
			}
		})
	}
}

func TestEndpointsAPIEnabled_invalidValueWarns(t *testing.T) {
	cases := []struct {
		value    string
		want     bool
		wantWarn bool
	}{
		{"", true, false},
		{"TRUE", true, false},
		{"0", false, false},
		{"false", false, false},
		{"on", true, true},
	}
	for _, tc := range cases {
		var logs bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&logs, nil))
		if got := endpointsAPIEnabled(t.Context(), logger, tc.value); got != tc.want {
			t.Errorf("endpointsAPIEnabled(%q) = %v, want %v", tc.value, got, tc.want)
		}
		if got := strings.Contains(logs.String(), endpointsAPIEnv); got != tc.wantWarn {
			t.Errorf("endpointsAPIEnabled(%q) warned = %v, want %v; logs:\n%s", tc.value, got, tc.wantWarn, logs.String())
		}
	}
}

// The endpoint-list TTL tests port the Python SDK's from
// tests/test_files_walk.py. Each scenario there is one transfer, so the
// contract can't show expiry; these drive transferEndpoints with a clock the
// test moves.

const (
	ttlProxyList = `{"endpoints":[{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com","multipart_upload_supported":true}]}`
	ttlRIGList   = `{"endpoints":[{"api":"FILES_API","api_version":"3.0","scheme":"https","hostname":"rig.example.com","multipart_upload_supported":true}]}`
)

var (
	ttlProxy = []endpoint{{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com", multipart: true}}
	ttlRIG   = []endpoint{{kind: filesAPIv3Endpoint, baseURL: "https://rig.example.com", multipart: true}}
)

// fakeClock is a clock that only the test moves.
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

// ttlDiscovery is a workspace whose endpoints answer the test can change, on an
// engine whose clock the test moves.
type ttlDiscovery struct {
	e     *engine
	clock *fakeClock
	calls *atomic.Int32
	mu    sync.Mutex
	h     http.HandlerFunc
}

func newTTLDiscovery(t *testing.T, h http.HandlerFunc) *ttlDiscovery {
	t.Helper()
	d := &ttlDiscovery{h: h, clock: &fakeClock{t: time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)}}
	d.e, d.calls = discoveryEngine(t, func(w http.ResponseWriter, r *http.Request) {
		d.mu.Lock()
		h := d.h
		d.mu.Unlock()
		h(w, r)
	})
	d.e.now = d.clock.Now
	return d
}

func (d *ttlDiscovery) answer(h http.HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.h = h
}

// expect checks that a transfer now walks want, and that the workspace has
// been asked wantCalls times in all.
func (d *ttlDiscovery) expect(t *testing.T, want []endpoint, wantCalls int32) {
	t.Helper()
	if got := d.e.transferEndpoints(t.Context()); !slices.Equal(got, want) {
		t.Errorf("endpoints = %v, want %v", got, want)
	}
	if got := d.calls.Load(); got != wantCalls {
		t.Errorf("endpoints requests = %d, want %d", got, wantCalls)
	}
}

// holdNext answers the next request with h only once release is called, and
// closes arrived when that request arrives. Later requests get h at once.
func holdNext(t *testing.T, h http.HandlerFunc) (handler http.HandlerFunc, arrived <-chan struct{}, release func()) {
	t.Helper()
	arrivedCh, releaseCh := make(chan struct{}), make(chan struct{})
	var first atomic.Bool
	var once sync.Once
	release = func() { once.Do(func() { close(releaseCh) }) }
	t.Cleanup(release)
	return func(w http.ResponseWriter, r *http.Request) {
		if first.CompareAndSwap(false, true) {
			close(arrivedCh)
			<-releaseCh
		}
		h(w, r)
	}, arrivedCh, release
}

// within returns what ch delivers, failing the test if that takes too long.
func within[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(10 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
		panic("unreachable")
	}
}

func TestTransferEndpoints_askedAgainOnceTheTTLPasses(t *testing.T) {
	d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlProxyList))
	d.expect(t, ttlProxy, 1)
	d.answer(respondJSON(http.StatusOK, ttlRIGList))
	d.clock.Advance(endpointsCacheTTL - time.Second)
	d.expect(t, ttlProxy, 1)

	d.clock.Advance(time.Second)
	d.expect(t, ttlRIG, 2)
}

func TestTransferEndpoints_failedAnswerKeptForTheFailureTTL(t *testing.T) {
	// The workspace doesn't serve the endpoints API yet.
	d := newTTLDiscovery(t, respondJSON(http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`))
	builtin := d.e.builtinEndpoints()
	d.expect(t, builtin, 1)
	d.answer(respondJSON(http.StatusOK, ttlRIGList))
	d.clock.Advance(endpointsFailureCacheTTL - time.Second)
	d.expect(t, builtin, 1)

	d.clock.Advance(time.Second)
	d.expect(t, ttlRIG, 2)
}

func TestTransferEndpoints_transientRefreshFailureKeepsTheServersList(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"408":                       respondJSON(http.StatusRequestTimeout, `{"error_code":"REQUEST_TIMEOUT","message":"slow"}`),
		"429":                       respondJSON(http.StatusTooManyRequests, `{"error_code":"TOO_MANY_REQUESTS","message":"slow down"}`),
		"500":                       respondJSON(http.StatusInternalServerError, `{"error_code":"INTERNAL_ERROR","message":"boom"}`),
		"503":                       respondJSON(http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`),
		"a dropped connection":      func(w http.ResponseWriter, _ *http.Request) { dropConnection(w) },
		"a body that doesn't parse": respondJSON(http.StatusOK, `<html>`),
	}
	for name, blip := range cases {
		t.Run(name, func(t *testing.T) {
			d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlRIGList))
			d.expect(t, ttlRIG, 1)
			d.answer(blip)
			d.clock.Advance(endpointsCacheTTL)
			// The blip doesn't move the client off the workspace's list.
			if got := d.e.transferEndpoints(t.Context()); !slices.Equal(got, ttlRIG) {
				t.Errorf("endpoints after the failed refresh = %v, want %v", got, ttlRIG)
			}
			afterBlip := d.calls.Load()
			if afterBlip < 2 {
				t.Fatalf("endpoints requests = %d, want a refresh", afterBlip)
			}

			// It asks again once the failure TTL passes.
			d.answer(respondJSON(http.StatusOK, ttlProxyList))
			d.clock.Advance(endpointsFailureCacheTTL - time.Second)
			d.expect(t, ttlRIG, afterBlip)
			d.clock.Advance(time.Second)
			d.expect(t, ttlProxy, afterBlip+1)
		})
	}
}

func TestTransferEndpoints_definitiveAnswerMovesToTheBuiltinOrder(t *testing.T) {
	// How the workspace rolls discovery back: the endpoints API answers 403, or
	// 404 where it is not served. Any other answer that is the workspace's own
	// also drops its list.
	cases := map[string]http.HandlerFunc{
		"403":                   respondJSON(http.StatusForbidden, `{"error_code":"PERMISSION_DENIED","message":"off"}`),
		"404":                   respondJSON(http.StatusNotFound, `{"error_code":"NOT_FOUND","message":"no"}`),
		"400":                   respondJSON(http.StatusBadRequest, `{"error_code":"INVALID_PARAMETER_VALUE","message":"no"}`),
		"no endpoints":          respondJSON(http.StatusOK, `{}`),
		"endpoints not a list":  respondJSON(http.StatusOK, `{"endpoints":{"api":"IO_OFFLOAD"}}`),
		"only unusable entries": respondJSON(http.StatusOK, `{"endpoints":[{"api":"NEW_API","scheme":"https","hostname":"new.example.com"}]}`),
		"204":                   respondJSON(http.StatusNoContent, ""),
		"a list with a 201":     respondJSON(http.StatusCreated, ttlProxyList),
	}
	for name, answer := range cases {
		t.Run(name, func(t *testing.T) {
			d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlRIGList))
			d.expect(t, ttlRIG, 1)
			d.answer(answer)
			d.clock.Advance(endpointsCacheTTL)
			d.expect(t, d.e.builtinEndpoints(), 2)
		})
	}
}

func TestTransferEndpoints_listKeptThroughFailingRefreshesExpiresAtTheStaleLimit(t *testing.T) {
	d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlRIGList))
	start := d.clock.Now()
	d.expect(t, ttlRIG, 1)
	// Discovery breaks for good, with an error that looks transient.
	d.answer(respondJSON(http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`))
	d.clock.Advance(endpointsCacheTTL)
	for d.clock.Now().Sub(start)+endpointsFailureCacheTTL < endpointsMaxStale {
		if got := d.e.transferEndpoints(t.Context()); !slices.Equal(got, ttlRIG) {
			t.Fatalf("endpoints %v after the list arrived = %v, want %v", d.clock.Now().Sub(start), got, ttlRIG)
		}
		d.clock.Advance(endpointsFailureCacheTTL)
	}
	d.clock.Advance(endpointsFailureCacheTTL)
	if got, want := d.e.transferEndpoints(t.Context()), d.e.builtinEndpoints(); !slices.Equal(got, want) {
		t.Errorf("endpoints past the stale limit = %v, want the built-in %v", got, want)
	}
}

func TestTransferEndpoints_failingRefreshJustBeforeTheStaleLimitKeepsTheListOnlyUntilTheLimit(t *testing.T) {
	d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlRIGList))
	d.expect(t, ttlRIG, 1)
	d.answer(respondJSON(http.StatusServiceUnavailable, `{"error_code":"TEMPORARILY_UNAVAILABLE","message":"busy"}`))
	// A refresh a minute before the limit fails and keeps the list, but only
	// until the limit, not a failure TTL more.
	d.clock.Advance(endpointsMaxStale - time.Minute)
	d.expect(t, ttlRIG, 2)
	d.clock.Advance(time.Minute)
	d.expect(t, d.e.builtinEndpoints(), 3)
}

func TestTransferEndpoints_keepTheExpiredListWhileAnotherCallerAsks(t *testing.T) {
	d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlProxyList))
	d.expect(t, ttlProxy, 1)
	d.clock.Advance(endpointsCacheTTL)
	held, arrived, release := holdNext(t, respondJSON(http.StatusOK, ttlRIGList))
	d.answer(held)

	asker := make(chan []endpoint, 1)
	go func() { asker <- d.e.transferEndpoints(t.Context()) }()
	within(t, arrived, "the refresh to reach the workspace")
	// The asker holds the refresh, so this caller neither asks nor waits.
	d.expect(t, ttlProxy, 2)

	release()
	if got := within(t, asker, "the asker"); !slices.Equal(got, ttlRIG) {
		t.Errorf("the asker's endpoints = %v, want the new %v", got, ttlRIG)
	}
	d.expect(t, ttlRIG, 2)
}

func TestTransferEndpoints_dropAListPastTheStaleLimitWhileAnotherCallerAsks(t *testing.T) {
	d := newTTLDiscovery(t, respondJSON(http.StatusOK, ttlRIGList))
	d.expect(t, ttlRIG, 1)
	d.clock.Advance(endpointsMaxStale)
	// Stands in for a refresh that takes long.
	held, arrived, release := holdNext(t, respondJSON(http.StatusOK, ttlRIGList))
	d.answer(held)

	asker := make(chan []endpoint, 1)
	go func() { asker <- d.e.transferEndpoints(t.Context()) }()
	within(t, arrived, "the refresh to reach the workspace")
	// The list is past the stale limit, so this caller uses the built-in order.
	d.expect(t, d.e.builtinEndpoints(), 2)

	release()
	if got := within(t, asker, "the asker"); !slices.Equal(got, ttlRIG) {
		t.Errorf("the asker's endpoints = %v, want %v", got, ttlRIG)
	}
}

func TestTransferEndpoints_withNoListYetWaitForTheFirstAnswer(t *testing.T) {
	held, arrived, release := holdNext(t, respondJSON(http.StatusOK, ttlProxyList))
	d := newTTLDiscovery(t, held)
	first := make(chan []endpoint, 1)
	go func() { first <- d.e.transferEndpoints(t.Context()) }()
	within(t, arrived, "the refresh to reach the workspace")

	cancelled, cancel := context.WithCancel(t.Context())
	gaveUp := make(chan []endpoint, 1)
	go func() { gaveUp <- d.e.transferEndpoints(cancelled) }()
	waiter := make(chan []endpoint, 1)
	go func() { waiter <- d.e.transferEndpoints(t.Context()) }()
	// A caller that gives up while waiting walks the built-in order.
	cancel()
	if got, want := within(t, gaveUp, "the cancelled waiter"), d.e.builtinEndpoints(); !slices.Equal(got, want) {
		t.Errorf("the cancelled waiter's endpoints = %v, want the built-in %v", got, want)
	}

	release()
	for name, ch := range map[string]chan []endpoint{"the asker": first, "the waiter": waiter} {
		if got := within(t, ch, name); !slices.Equal(got, ttlProxy) {
			t.Errorf("%s's endpoints = %v, want %v", name, got, ttlProxy)
		}
	}
	if got := d.calls.Load(); got != 1 {
		t.Errorf("endpoints requests = %d, want 1", got)
	}
}

func TestTransferEndpoints_notBlockedByTheClientsBusyConnections(t *testing.T) {
	// A refresh runs while transfers may hold every connection the client's
	// transport allows. Here the only one is held, and discovery still gets an
	// answer before its caller gives up.
	holding, release := make(chan struct{}), make(chan struct{})
	srv := httptest.NewServer(requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == endpointsAPIPath {
			respondJSON(http.StatusOK, ttlProxyList)(w, r)
			return
		}
		close(holding)
		<-release
	}))
	t.Cleanup(srv.Close)
	// Cleanups run last first: release the held request before Close waits for it.
	t.Cleanup(func() { close(release) })
	e := newEngine(buildUploadClient(t, srv.URL, &http.Client{Transport: &http.Transport{MaxConnsPerHost: 1}}, ""))
	e.tun.endpointsAPI = true

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/a-long-transfer", nil)
	noErr(t, err)
	req.Header.Set("Authorization", testToken)
	go func() {
		if resp, err := e.c.httpClient.Do(req); err == nil {
			_ = resp.Body.Close()
		}
	}()
	within(t, holding, "the long transfer to hold its connection")
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	if got := e.transferEndpoints(ctx); !slices.Equal(got, ttlProxy) {
		t.Errorf("endpoints = %v, want %v: discovery waited for the client's connection", got, ttlProxy)
	}
}

func TestTransferEndpoints_cancelledMidRefreshIsNotKept(t *testing.T) {
	// A caller that gives up while it refreshes leaves no answer behind, so
	// the next caller asks the workspace itself.
	held, arrived, release := holdNext(t, respondJSON(http.StatusOK, ttlProxyList))
	d := newTTLDiscovery(t, held)
	ctx, cancel := context.WithCancel(t.Context())
	gaveUp := make(chan []endpoint, 1)
	go func() { gaveUp <- d.e.transferEndpoints(ctx) }()
	within(t, arrived, "the refresh to reach the workspace")
	cancel()
	if got, want := within(t, gaveUp, "the cancelled caller"), d.e.builtinEndpoints(); !slices.Equal(got, want) {
		t.Errorf("the cancelled caller's endpoints = %v, want the built-in %v", got, want)
	}
	release()
	d.expect(t, ttlProxy, 2)
}

func TestDiscoveryClient_doesNotShareTheClientsPool(t *testing.T) {
	own := &http.Transport{MaxConnsPerHost: 4}
	custom := roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("unused") })
	cases := []struct {
		name      string
		transport http.RoundTripper
		shared    bool
	}{
		{"the default transport", nil, false},
		{"http.DefaultTransport", http.DefaultTransport, false},
		{"the caller's *http.Transport", own, false},
		// A RoundTripper of another type can't be copied.
		{"another RoundTripper", custom, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hc := &http.Client{Transport: tc.transport, Timeout: time.Minute}
			e := newEngine(&Client{internalClient: internalClient{httpClient: hc}})
			got := e.discoveryClient()
			if got.Timeout != hc.Timeout {
				t.Errorf("Timeout = %v, want the client's %v", got.Timeout, hc.Timeout)
			}
			if tc.shared {
				if got != hc {
					t.Error("discovery didn't use the client as is")
				}
				return
			}
			tr, ok := got.Transport.(*http.Transport)
			switch {
			case !ok:
				t.Fatalf("Transport = %T, want an *http.Transport", got.Transport)
			case tr == http.DefaultTransport || tr == own:
				t.Error("discovery shares the client's transport, and so its connection pool")
			case tc.transport == own && tr.MaxConnsPerHost != own.MaxConnsPerHost:
				t.Errorf("MaxConnsPerHost = %d, want the client's %d", tr.MaxConnsPerHost, own.MaxConnsPerHost)
			}
		})
	}
}

func TestLoggedBody_capsLongBodies(t *testing.T) {
	if got := loggedBody([]byte(`{"endpoints":[]}`)); got != `{"endpoints":[]}` {
		t.Errorf("loggedBody(short) = %q, want it unchanged", got)
	}
	if got := loggedBody(bytes.Repeat([]byte("x"), 10*maxLoggedBody)); len(got) != maxLoggedBody+len("...") {
		t.Errorf("len(loggedBody(10 KiB)) = %d, want %d", len(got), maxLoggedBody+len("..."))
	}
}

func TestTransferEndpoints_sharedTransportThatNeverAnswersGivesTheBuiltinOrderInTime(t *testing.T) {
	// A RoundTripper other than *http.Transport can't be copied, so discovery
	// shares it. If every connection it has is held, discovery
	// still ends at its timeout, on the built-in order.
	t.Parallel()
	blocked := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	e := newEngine(buildUploadClient(t, "https://ws.example.com", &http.Client{Transport: blocked}, ""))
	e.tun.endpointsAPI = true
	start := time.Now()
	done := make(chan []endpoint, 1)
	go func() { done <- e.transferEndpoints(t.Context()) }()
	select {
	case got := <-done:
		if want := e.builtinEndpoints(); !slices.Equal(got, want) {
			t.Errorf("endpoints = %v, want the built-in %v", got, want)
		}
		if elapsed := time.Since(start); elapsed > endpointsResolveTimeout+5*time.Second {
			t.Errorf("discovery took %v, want about %v", elapsed, endpointsResolveTimeout)
		}
	case <-time.After(endpointsResolveTimeout + 10*time.Second):
		t.Fatal("discovery waited past its timeout for a transport that never answers")
	}
}

func TestDiscovery_sendsTheSDKUserAgent(t *testing.T) {
	var got atomic.Value
	e, _ := discoveryEngine(t, func(w http.ResponseWriter, r *http.Request) {
		got.Store(r.UserAgent())
		respondJSON(http.StatusOK, `{"endpoints":[]}`)(w, r)
	})
	e.transferEndpoints(t.Context())
	want, err := e.c.userAgent()
	noErr(t, err)
	if want == "" || got.Load() != want {
		t.Errorf("discovery User-Agent = %q, want the SDK's %q", got.Load(), want)
	}
}
