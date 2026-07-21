package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func TestDiscoverAuthorizationServer_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oidc/.well-known/oauth-authorization-server" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"token_endpoint":         "https://auth.example.com/oidc/v1/token",
			"authorization_endpoint": "https://auth.example.com/oidc/v1/authorize",
		})
	}))
	defer server.Close()

	got, err := discoverAuthorizationServer(context.Background(), http.DefaultClient, server.URL)
	if err != nil {
		t.Fatalf("discoverAuthorizationServer err = %v", err)
	}

	want := &oauthAuthorizationServer{
		TokenEndpoint:         "https://auth.example.com/oidc/v1/token",
		AuthorizationEndpoint: "https://auth.example.com/oidc/v1/authorize",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("discoverAuthorizationServer mismatch (-want +got):\n%s", diff)
	}
}

func TestDiscoverAuthorizationServer_TrimsTrailingSlash(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token_endpoint": "x"})
	}))
	defer server.Close()

	_, err := discoverAuthorizationServer(context.Background(), http.DefaultClient, server.URL+"/")
	if err != nil {
		t.Fatalf("discoverAuthorizationServer err = %v", err)
	}
	if gotPath != "/oidc/.well-known/oauth-authorization-server" {
		t.Errorf("request path = %q, want %q", gotPath, "/oidc/.well-known/oauth-authorization-server")
	}
}

func TestDiscoverAuthorizationServer_MissingHost(t *testing.T) {
	_, err := discoverAuthorizationServer(context.Background(), nil, "")
	if !errors.Is(err, errHostRequired) {
		t.Errorf("err = %v, want %v", err, errHostRequired)
	}
}

func TestDiscoverAuthorizationServer_Non2xx(t *testing.T) {
	testCases := []struct {
		desc     string
		status   int
		wantCode codes.Code
	}{
		{"503", http.StatusServiceUnavailable, codes.Unavailable},
		{"429", http.StatusTooManyRequests, codes.ResourceExhausted},
		{"500", http.StatusInternalServerError, codes.Internal},
		{"401", http.StatusUnauthorized, codes.Unauthenticated},
		{"404", http.StatusNotFound, codes.NotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "nope", tc.status)
			}))
			defer server.Close()

			_, err := discoverAuthorizationServer(context.Background(), nil, server.URL)
			if err == nil {
				t.Fatal("err = nil, want non-nil")
			}
			if got := apierr.Code(err); got != tc.wantCode {
				t.Errorf("apierr.Code(err) = %v, want %v", got, tc.wantCode)
			}
		})
	}
}

func TestDiscoverAuthorizationServer_MissingTokenEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	_, err := discoverAuthorizationServer(context.Background(), nil, server.URL)
	if err == nil || !strings.Contains(err.Error(), "missing token_endpoint") {
		t.Errorf("err = %v, want to contain \"missing token_endpoint\"", err)
	}
}

func TestDiscoverAuthorizationServer_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{not json"))
	}))
	defer server.Close()

	_, err := discoverAuthorizationServer(context.Background(), nil, server.URL)
	if err == nil || !strings.Contains(err.Error(), "decoding response") {
		t.Errorf("err = %v, want to contain \"decoding response\"", err)
	}
}

func TestDiscoverAuthorizationServer_NilHTTPClient(t *testing.T) {
	// Nil httpClient should be defaulted to http.DefaultClient and still work.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"token_endpoint": "x"})
	}))
	defer server.Close()

	_, err := discoverAuthorizationServer(context.Background(), nil, server.URL)
	if err != nil {
		t.Errorf("discoverAuthorizationServer err = %v, want nil", err)
	}
}

func TestFetchToken_Success(t *testing.T) {
	var gotAuthUser, gotAuthPass, gotGrantType, gotScope string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gotAuthUser, gotAuthPass, _ = r.BasicAuth()
		gotGrantType = r.PostForm.Get("grant_type")
		gotScope = r.PostForm.Get("scope")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer server.Close()

	cfg := &clientcredentials.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     server.URL,
		Scopes:       []string{"scope-a"},
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	got, err := fetchToken(context.Background(), http.DefaultClient, cfg)
	if err != nil {
		t.Fatalf("fetchToken err = %v", err)
	}
	if got.AccessToken != "tok" {
		t.Errorf("AccessToken = %q, want %q", got.AccessToken, "tok")
	}
	if got.TokenType != "Bearer" {
		t.Errorf("TokenType = %q, want %q", got.TokenType, "Bearer")
	}
	if got.Expiry.IsZero() || got.Expiry.Before(time.Now()) {
		t.Errorf("Expiry = %v, want future", got.Expiry)
	}
	if gotAuthUser != "id" || gotAuthPass != "secret" {
		t.Errorf("basic auth = (%q, %q), want (id, secret)", gotAuthUser, gotAuthPass)
	}
	if gotGrantType != "client_credentials" {
		t.Errorf("grant_type = %q, want client_credentials", gotGrantType)
	}
	if gotScope != "scope-a" {
		t.Errorf("scope = %q, want scope-a", gotScope)
	}
}

type countingTransport struct {
	inner http.RoundTripper
	calls int32
}

func (c *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&c.calls, 1)
	return c.inner.RoundTrip(req)
}

func TestFetchToken_HTTPClientInjection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer server.Close()

	ct := &countingTransport{inner: http.DefaultTransport}
	client := &http.Client{Transport: ct}

	cfg := &clientcredentials.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     server.URL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	if _, err := fetchToken(context.Background(), client, cfg); err != nil {
		t.Fatalf("fetchToken err = %v", err)
	}
	if got := atomic.LoadInt32(&ct.calls); got != 1 {
		t.Errorf("injected client calls = %d, want 1", got)
	}
}

func TestFetchToken_ErrorFromEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	cfg := &clientcredentials.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		TokenURL:     server.URL,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}

	_, err := fetchToken(context.Background(), nil, cfg)
	if err == nil {
		t.Fatal("err = nil, want non-nil")
	}
	// fetchToken bridges oauth2.RetrieveError to *apierr.APIError so that
	// downstream code (retry classification) sees a single error taxonomy.
	if code := apierr.Code(err); code != codes.Unavailable {
		t.Errorf("apierr.Code(err) = %v, want %v", code, codes.Unavailable)
	}
	var apiErr *apierr.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err is not *apierr.APIError: %T", err)
	}
	if got := apiErr.HTTPStatusCode(); got != http.StatusServiceUnavailable {
		t.Errorf("HTTPStatusCode = %d, want %d", got, http.StatusServiceUnavailable)
	}
}
