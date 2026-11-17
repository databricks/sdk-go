package files

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync/atomic"
	"testing"
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
