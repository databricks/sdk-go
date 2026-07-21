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
	lastAuthUser string
	lastAuthPass string

	// Response overrides.
	accessToken  string
	tokenType    string
	expiresIn    int
	tokenStatus  int // 0 means 200
	discoveryBad bool
}

func newFakeOIDCServer() *fakeOIDCServer {
	f := &fakeOIDCServer{
		accessToken: "access-token",
		tokenType:   "Bearer",
		expiresIn:   3600,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/oidc/.well-known/oauth-authorization-server", func(w http.ResponseWriter, r *http.Request) {
		if f.discoveryBad {
			http.Error(w, "nope", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"token_endpoint": f.server.URL + "/oidc/v1/token",
		})
	})
	mux.HandleFunc("/oidc/v1/token", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.tokenCalls, 1)
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.lastRequest = r.PostForm
		f.lastAuthUser, f.lastAuthPass, _ = r.BasicAuth()

		if f.tokenStatus != 0 {
			http.Error(w, "failure", f.tokenStatus)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": f.accessToken,
			"token_type":   f.tokenType,
			"expires_in":   f.expiresIn,
		})
	})
	f.server = httptest.NewServer(mux)
	return f
}

func (f *fakeOIDCServer) Close() { f.server.Close() }

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
