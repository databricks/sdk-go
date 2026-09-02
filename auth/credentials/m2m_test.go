package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/apierr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewM2MCredentials_validation(t *testing.T) {
	testCases := []struct {
		desc    string
		opts    M2MOptions
		wantErr error
	}{
		{
			desc:    "missing client ID",
			opts:    M2MOptions{ClientSecret: "s", Host: "https://h"},
			wantErr: errClientIDRequired,
		},
		{
			desc:    "missing client secret",
			opts:    M2MOptions{ClientID: "c", Host: "https://h"},
			wantErr: errClientSecretRequired,
		},
		{
			desc:    "missing host",
			opts:    M2MOptions{ClientID: "c", ClientSecret: "s"},
			wantErr: errHostRequired,
		},
		{
			desc: "all required fields",
			opts: M2MOptions{ClientID: "c", ClientSecret: "s", Host: "https://h"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := NewM2MCredentials(tc.opts)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("NewM2MCredentials() err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewM2MCredentials() err = %v, want nil", err)
			}
			if got == nil {
				t.Fatal("NewM2MCredentials() returned nil")
			}
		})
	}
}

// fakeOIDCServer stands up an httptest server that speaks the discovery
// document at /oidc/.well-known/oauth-authorization-server and a token
// endpoint at /oidc/v1/token.
type fakeOIDCServer struct {
	server       *httptest.Server
	tokenCalls   int32
	lastRequest  url.Values
	requests     []url.Values
	lastAuthUser string
	lastAuthPass string

	// Response overrides.
	tokenType    string
	expiresIn    int
	tokenStatus  int // 0 means 200
	discoveryBad bool
	tokenPath    string
	// tokenValue makes the fake server's access token an explicit function of
	// the request, allowing tests to expose differences between token forms.
	tokenValue func(url.Values) string
}

func newFakeOIDCServer() *fakeOIDCServer {
	return newFakeOIDCServerWithTokenPath("/oidc/v1/token")
}

// newFakeOIDCServerWithTokenPath creates a discoverable fake OIDC server whose
// token endpoint uses the supplied path.
func newFakeOIDCServerWithTokenPath(tokenPath string) *fakeOIDCServer {
	f := &fakeOIDCServer{
		tokenType:  "Bearer",
		expiresIn:  3600,
		tokenPath:  tokenPath,
		tokenValue: func(url.Values) string { return "access-token" },
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/oidc/.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) {
		if f.discoveryBad {
			http.Error(w, "nope", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"token_endpoint": f.server.URL + f.tokenPath,
		})
	})

	mux.HandleFunc(tokenPath, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.tokenCalls, 1)

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		f.lastRequest = r.PostForm
		f.requests = append(f.requests, r.PostForm)
		f.lastAuthUser, f.lastAuthPass, _ = r.BasicAuth()

		if f.tokenStatus != 0 {
			http.Error(w, "failure", f.tokenStatus)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": f.tokenValue(r.PostForm),
			"token_type":   f.tokenType,
			"expires_in":   f.expiresIn,
		})
	})

	f.server = httptest.NewServer(mux)

	return f
}

func (f *fakeOIDCServer) Close() { f.server.Close() }

// newTestM2MProvider creates an M2M provider with the standard credentials
// shared by group-assumption token exchange tests.
func newTestM2MProvider(t *testing.T, server *fakeOIDCServer, groupID string) auth.TokenProvider {
	t.Helper()

	provider, err := NewM2MCredentials(M2MOptions{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Host:         server.server.URL,
		GroupID:      groupID,
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials() error = %v", err)
	}

	return provider
}

func TestM2M_Token_Success(t *testing.T) {
	s := newFakeOIDCServer()
	defer s.Close()

	tp, err := NewM2MCredentials(M2MOptions{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Host:         s.server.URL,
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials() error = %v", err)
	}

	got, err := tp.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	want := &auth.Token{Value: "access-token", Type: "Bearer"}
	if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(auth.Token{}, "Expiry")); diff != "" {
		t.Errorf("Token() mismatch (-want +got):\n%s", diff)
	}
	if got.Expiry.IsZero() || got.Expiry.Before(time.Now()) {
		t.Errorf("Token().Expiry = %v, want non-zero future time", got.Expiry)
	}
	if got, want := s.lastAuthUser, "client-id"; got != want {
		t.Errorf("basic auth user = %q, want %q", got, want)
	}
	if got, want := s.lastAuthPass, "client-secret"; got != want {
		t.Errorf("basic auth pass = %q, want %q", got, want)
	}
	if got, want := s.lastRequest.Get("grant_type"), "client_credentials"; got != want {
		t.Errorf("grant_type = %q, want %q", got, want)
	}
	if got, want := s.lastRequest.Get("scope"), "all-apis"; got != want {
		t.Errorf("default scope = %q, want %q", got, want)
	}
}

// TestM2M_Token_AssumeGroup verifies that M2M token requests include the
// configured group and omit assume_group when no group is configured.
func TestM2M_Token_AssumeGroup(t *testing.T) {
	testCases := []struct {
		name    string
		groupID string
	}{
		{
			name: "normal",
		},
		{
			name:    "group A",
			groupID: "group-a",
		},
		{
			name:    "group B",
			groupID: "group-b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newFakeOIDCServer()
			defer server.Close()

			provider := newTestM2MProvider(t, server, tc.groupID)

			// Repeated exchanges verify that the fixed provider configuration is
			// retained after discovery has been cached.
			for range 2 {
				if _, err := provider.Token(context.Background()); err != nil {
					t.Fatalf("Token() error = %v", err)
				}

				wantForm := url.Values{
					"grant_type": {"client_credentials"},
					"scope":      {"all-apis"},
				}
				if tc.groupID != "" {
					wantForm.Set("assume_group", tc.groupID)
				}

				if diff := cmp.Diff(wantForm, server.lastRequest); diff != "" {
					t.Errorf("token form mismatch (-want +got):\n%s", diff)
				}
			}

			if got, want := atomic.LoadInt32(&server.tokenCalls), int32(2); got != want {
				t.Errorf("token calls = %d, want %d", got, want)
			}
		})
	}
}

// TestM2M_Token_AssumeGroupForEveryEndpointType verifies that group assumption
// is sent to workspace and the shared account or unified token endpoint shape.
func TestM2M_Token_AssumeGroupForEveryEndpointType(t *testing.T) {
	testCases := []struct {
		name      string
		tokenPath string
	}{
		{
			name:      "workspace",
			tokenPath: "/oidc/v1/token",
		},
		{
			name:      "account or unified",
			tokenPath: "/oidc/accounts/account-id/v1/token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := newFakeOIDCServerWithTokenPath(tc.tokenPath)
			defer server.Close()

			provider := newTestM2MProvider(t, server, "group-id")

			if _, err := provider.Token(context.Background()); err != nil {
				t.Fatalf("Token() error = %v", err)
			}

			if got, want := server.lastRequest.Get("assume_group"), "group-id"; got != want {
				t.Errorf("assume_group = %q, want %q", got, want)
			}
		})
	}
}

// TestM2M_GroupCachesAreIsolated verifies that cached providers configured for
// different groups retain the token issued for their own group.
func TestM2M_GroupCachesAreIsolated(t *testing.T) {
	server := newFakeOIDCServer()
	defer server.Close()

	server.tokenValue = func(form url.Values) string {
		if groupID := form.Get("assume_group"); groupID != "" {
			return groupID + "-token"
		}
		return "normal-token"
	}

	for _, groupID := range []string{"", "group-a", "group-b"} {
		provider := newTestM2MProvider(t, server, groupID)

		cached := auth.NewCachedTokenProvider(provider, auth.WithAsyncRefresh(false))
		for range 2 {
			token, err := cached.Token(context.Background())
			if err != nil {
				t.Fatalf("Token() error = %v", err)
			}

			want := "normal-token"
			if groupID != "" {
				want = groupID + "-token"
			}

			if token.Value != want {
				t.Errorf("Token().Value = %q, want %q", token.Value, want)
			}
		}
	}

	if got, want := atomic.LoadInt32(&server.tokenCalls), int32(3); got != want {
		t.Errorf("token calls = %d, want %d", got, want)
	}
}

// TestM2M_ExpiredCachedTokenRetainsGroup verifies that every refresh of an
// expired M2M token retains the configured group.
func TestM2M_ExpiredCachedTokenRetainsGroup(t *testing.T) {
	server := newFakeOIDCServer()
	defer server.Close()

	server.expiresIn = -1

	provider := newTestM2MProvider(t, server, "group-id")

	cached := auth.NewCachedTokenProvider(provider, auth.WithAsyncRefresh(false))
	for range 2 {
		if _, err := cached.Token(context.Background()); err != nil {
			t.Fatalf("Token() error = %v", err)
		}
	}

	if got, want := atomic.LoadInt32(&server.tokenCalls), int32(2); got != want {
		t.Fatalf("token calls = %d, want %d", got, want)
	}

	for i, form := range server.requests {
		if got, want := form.Get("assume_group"), "group-id"; got != want {
			t.Errorf("request %d assume_group = %q, want %q", i, got, want)
		}
	}
}

// TestM2M_Token_GroupRejectionIsNotRetried verifies that a rejected grouped
// token request is returned to the caller without another token request.
func TestM2M_Token_GroupRejectionIsNotRetried(t *testing.T) {
	server := newFakeOIDCServer()
	defer server.Close()

	server.tokenStatus = http.StatusForbidden

	provider := newTestM2MProvider(t, server, "group-id")

	_, err := provider.Token(context.Background())

	var apiError *apierr.APIError
	if !errors.As(err, &apiError) || apiError.HTTPStatusCode() != http.StatusForbidden {
		t.Fatalf("Token() error = %v, want wrapped APIError with status %d", err, http.StatusForbidden)
	}

	if got, want := atomic.LoadInt32(&server.tokenCalls), int32(1); got != want {
		t.Errorf("token calls = %d, want %d", got, want)
	}

	if got, want := server.lastRequest.Get("assume_group"), "group-id"; got != want {
		t.Errorf("assume_group = %q, want %q", got, want)
	}
}

func TestM2M_Token_CustomScopes(t *testing.T) {
	s := newFakeOIDCServer()
	defer s.Close()

	tp, err := NewM2MCredentials(M2MOptions{
		ClientID:     "c",
		ClientSecret: "s",
		Host:         s.server.URL,
		Scopes:       []string{"scope-a", "scope-b"},
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials() error = %v", err)
	}

	if _, err := tp.Token(context.Background()); err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if got, want := s.lastRequest.Get("scope"), "scope-a scope-b"; got != want {
		t.Errorf("scope = %q, want %q", got, want)
	}
}

func TestM2M_Token_DiscoveryFailure(t *testing.T) {
	s := newFakeOIDCServer()
	s.discoveryBad = true
	defer s.Close()

	tp, err := NewM2MCredentials(M2MOptions{
		ClientID:     "c",
		ClientSecret: "s",
		Host:         s.server.URL,
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials() error = %v", err)
	}

	_, err = tp.Token(context.Background())
	if err == nil || !strings.Contains(err.Error(), "oidc discovery") {
		t.Errorf("Token() error = %v, want oidc discovery error", err)
	}
}

func TestM2M_Token_DiscoveryIsMemoized(t *testing.T) {
	var discoveryCalls int32
	s := newFakeOIDCServer()
	defer s.Close()
	s.server.Config.Handler = wrapMux(s.server.Config.Handler, func(path string) {
		if path == "/oidc/.well-known/oauth-authorization-server" {
			atomic.AddInt32(&discoveryCalls, 1)
		}
	})

	tp, err := NewM2MCredentials(M2MOptions{
		ClientID:     "c",
		ClientSecret: "s",
		Host:         s.server.URL,
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials err = %v", err)
	}

	for range 3 {
		if _, err := tp.Token(context.Background()); err != nil {
			t.Fatalf("Token() err = %v", err)
		}
	}

	if got := atomic.LoadInt32(&discoveryCalls); got != 1 {
		t.Errorf("discovery calls = %d, want 1 (should be memoized)", got)
	}
}

// wrapMux returns an http.Handler that runs observe for every request before
// delegating to inner.
func wrapMux(inner http.Handler, observe func(path string)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		observe(r.URL.Path)
		inner.ServeHTTP(w, r)
	})
}

func TestM2M_Token_HTTPClientInjection(t *testing.T) {
	s := newFakeOIDCServer()
	defer s.Close()

	rt := &countingTransport{inner: http.DefaultTransport}
	client := &http.Client{Transport: rt}

	tp, err := NewM2MCredentials(M2MOptions{
		ClientID:     "c",
		ClientSecret: "s",
		Host:         s.server.URL,
		HTTPClient:   client,
	})
	if err != nil {
		t.Fatalf("NewM2MCredentials() error = %v", err)
	}

	if _, err := tp.Token(context.Background()); err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	// At least two calls: one for discovery, one for the token exchange.
	if got := atomic.LoadInt32(&rt.calls); got < 2 {
		t.Errorf("injected HTTP client calls = %d, want >= 2", got)
	}
}
