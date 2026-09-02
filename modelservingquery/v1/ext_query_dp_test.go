package modelservingquery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/options/client"
	"github.com/google/go-cmp/cmp"
)

// patCreds is a static credential that only implements auth.Credentials, like a
// personal access token. It cannot mint OAuth tokens, so route optimization is
// not possible and QueryOptimized must report it as unavailable.
type patCreds struct{}

func (patCreds) Name() string { return "pat" }

func (patCreds) AuthHeaders(context.Context) ([]auth.Header, error) {
	return []auth.Header{{Key: "Authorization", Value: "Bearer pat-token"}}, nil
}

// oauthCreds implements auth.TokenCredentials (both AuthHeaders and Token), like
// an OAuth-based credential. Its presence makes route optimization possible.
type oauthCreds struct{}

func (oauthCreds) Name() string { return "oauth" }

func (oauthCreds) AuthHeaders(context.Context) ([]auth.Header, error) {
	return []auth.Header{{Key: "Authorization", Value: "Bearer cp-token"}}, nil
}

func (oauthCreds) Token(context.Context) (*auth.Token, error) {
	return &auth.Token{Value: "cp-token", Type: "Bearer"}, nil
}

const (
	testEndpointName = "my-endpoint"
	dpAccessToken    = "dp-access-token"
	dpAuthDetails    = "auth-details-blob"
)

func newQueryTestClient(t *testing.T, server *httptest.Server, creds auth.Credentials) *Client {
	t.Helper()
	c, err := NewClient(context.Background(),
		client.WithHost(server.URL),
		client.WithHTTPClient(server.Client()),
		client.WithCredentials(creds),
		client.WithWorkspaceID("ws-123"),
		client.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))),
		client.WithoutConfigFile(),
		client.WithoutEnv(),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

// counters records how many times each server route was hit.
type counters struct {
	cpGet         atomic.Int64 // GET /api/2.0/serving-endpoints/{name}
	cpInvocations atomic.Int64 // POST /api/serving-endpoints/{name}/invocations
	tokenExchange atomic.Int64 // POST /oidc/v1/token
	dpInvoke      atomic.Int64 // POST /dp/invocations
}

// routeCounts is the plain-value expectation compared against counters.
type routeCounts struct {
	cpGet         int64
	cpInvocations int64
	tokenExchange int64
	dpInvoke      int64
}

// serverConfig configures the mock server behavior for a test case.
type serverConfig struct {
	// includeDataPlaneInfo controls whether the control-plane GET advertises a
	// data-plane endpoint URL (i.e. the endpoint is route-optimized).
	includeDataPlaneInfo bool
	// tokenStatus is the HTTP status the token-exchange /oidc/v1/token route
	// returns. Zero means 200 OK.
	tokenStatus int
	// dpStatus is the HTTP status the data-plane /dp/invocations route returns.
	// Zero means 200 OK.
	dpStatus int
}

// captured records request details the assertions care about.
type captured struct {
	dpAuthHeader        string
	dpWorkspaceIDHeader string
	tokenAssertion      string
	tokenAuthDetails    string
}

func newMockServer(t *testing.T, cfg serverConfig, c *counters, cap *captured) *httptest.Server {
	t.Helper()

	queryResp := mustMarshalJSON(t, &queryEndpointResponseWire{Model: new("served-model")})

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/2.0/serving-endpoints/"+testEndpointName:
			c.cpGet.Add(1)
			body := servingEndpointGetPayload(t, cfg.includeDataPlaneInfo, srv.URL+"/dp/invocations")
			_, _ = w.Write(body)

		case r.Method == http.MethodPost && r.URL.Path == "/oidc/v1/token":
			c.tokenExchange.Add(1)
			_ = r.ParseForm()
			cap.tokenAssertion = r.Form.Get("assertion")
			cap.tokenAuthDetails = r.Form.Get("authorization_details")
			if cfg.tokenStatus != 0 && cfg.tokenStatus != http.StatusOK {
				http.Error(w, `{"error":"invalid_grant"}`, cfg.tokenStatus)
				return
			}
			_, _ = w.Write(mustMarshalJSON(t, map[string]any{
				"access_token": dpAccessToken,
				"token_type":   "Bearer",
				"expires_in":   3600,
			}))

		case r.Method == http.MethodPost && r.URL.Path == "/dp/invocations":
			c.dpInvoke.Add(1)
			cap.dpAuthHeader = r.Header.Get("Authorization")
			cap.dpWorkspaceIDHeader = r.Header.Get("X-Databricks-Workspace-Id")
			if cfg.dpStatus != 0 && cfg.dpStatus != http.StatusOK {
				http.Error(w, `{"error_code":"INTERNAL","message":"boom"}`, cfg.dpStatus)
				return
			}
			w.Header().Set("served-model-name", "served-model")
			_, _ = w.Write(queryResp)

		case r.Method == http.MethodPost && r.URL.Path == "/api/serving-endpoints/"+testEndpointName+"/invocations":
			c.cpInvocations.Add(1)
			w.Header().Set("served-model-name", "served-model")
			_, _ = w.Write(queryResp)

		default:
			http.Error(w, fmt.Sprintf(`{"error":"not found: %s %s"}`, r.Method, r.URL.Path), http.StatusNotFound)
		}
	}))
	return srv
}

func servingEndpointGetPayload(t *testing.T, includeDataPlaneInfo bool, dpURL string) []byte {
	t.Helper()
	payload := map[string]any{"name": testEndpointName}
	if includeDataPlaneInfo {
		payload["data_plane_info"] = map[string]any{
			"query_info": map[string]any{
				"endpoint_url":          dpURL,
				"authorization_details": dpAuthDetails,
			},
		}
	}
	return mustMarshalJSON(t, payload)
}

func mustMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestQueryOptimized(t *testing.T) {
	testCases := []struct {
		name  string
		creds auth.Credentials
		cfg   serverConfig
		calls int
		// wantErr is the sentinel every call must match with errors.Is: nil for
		// success or ErrRouteOptimizationUnavailable.
		wantErr error
		// wantAPIError is set when the call instead fails with an opaque
		// *apierr.APIError (data-plane or token-exchange failures), which carries
		// no sentinel to match against.
		wantAPIError bool
		want         routeCounts
	}{
		{
			name:    "non-OAuth credentials are unavailable without any request",
			creds:   patCreds{},
			cfg:     serverConfig{includeDataPlaneInfo: true},
			calls:   1,
			wantErr: ErrRouteOptimizationUnavailable,
			want:    routeCounts{},
		},
		{
			name:    "endpoint without data-plane info is unavailable after detection",
			creds:   oauthCreds{},
			cfg:     serverConfig{includeDataPlaneInfo: false},
			calls:   1,
			wantErr: ErrRouteOptimizationUnavailable,
			want:    routeCounts{cpGet: 1},
		},
		{
			name:  "route-optimized endpoint caches info and token across calls",
			creds: oauthCreds{},
			cfg:   serverConfig{includeDataPlaneInfo: true},
			calls: 2,
			// Detection GET and token exchange fire once and are reused; only
			// the billed data-plane invocation repeats.
			want: routeCounts{cpGet: 1, tokenExchange: 1, dpInvoke: 2},
		},
		{
			name:         "data-plane failure returns the error and evicts the endpoint",
			creds:        oauthCreds{},
			cfg:          serverConfig{includeDataPlaneInfo: true, dpStatus: http.StatusInternalServerError},
			calls:        2,
			wantAPIError: true,
			// Eviction re-runs detection and token exchange on the second call;
			// the control plane is never used as a fallback.
			want: routeCounts{cpGet: 2, tokenExchange: 2, dpInvoke: 2},
		},
		{
			name:         "token-exchange failure returns the error and evicts the endpoint",
			creds:        oauthCreds{},
			cfg:          serverConfig{includeDataPlaneInfo: true, tokenStatus: http.StatusUnauthorized},
			calls:        2,
			wantAPIError: true,
			// The data plane is never reached when minting fails; eviction
			// re-runs detection and exchange on the second call.
			want: routeCounts{cpGet: 2, tokenExchange: 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var c counters
			var cap captured
			srv := newMockServer(t, tc.cfg, &c, &cap)
			defer srv.Close()

			cl := newQueryTestClient(t, srv, tc.creds)

			for i := range tc.calls {
				resp, gotErr := cl.QueryOptimized(context.Background(), QueryEndpointRequest{Name: new(testEndpointName)})

				if tc.wantAPIError {
					if _, ok := errors.AsType[*apierr.APIError](gotErr); !ok {
						t.Fatalf("call %d: error = %v, want an *apierr.APIError", i, gotErr)
					}
					continue
				}

				if !errors.Is(gotErr, tc.wantErr) {
					t.Fatalf("call %d: error = %v, want %v", i, gotErr, tc.wantErr)
				}
				if tc.wantErr == nil && (resp.ServedModelName == nil || *resp.ServedModelName != "served-model") {
					t.Errorf("call %d: served-model-name header not applied: %+v", i, resp.ServedModelName)
				}
			}

			got := routeCounts{
				cpGet:         c.cpGet.Load(),
				cpInvocations: c.cpInvocations.Load(),
				tokenExchange: c.tokenExchange.Load(),
				dpInvoke:      c.dpInvoke.Load(),
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(routeCounts{})); diff != "" {
				t.Errorf("route hit counts mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestQueryOptimized_DataPlaneRequest checks the data-plane call carries the
// exchanged token and workspace id, and that the exchange forwards the
// control-plane token and the endpoint's authorization details.
func TestQueryOptimized_DataPlaneRequest(t *testing.T) {
	var c counters
	var cap captured
	srv := newMockServer(t, serverConfig{includeDataPlaneInfo: true}, &c, &cap)
	defer srv.Close()

	cl := newQueryTestClient(t, srv, oauthCreds{})
	if _, err := cl.QueryOptimized(context.Background(), QueryEndpointRequest{Name: new(testEndpointName)}); err != nil {
		t.Fatalf("QueryOptimized: %v", err)
	}

	if cap.dpAuthHeader != "Bearer "+dpAccessToken {
		t.Errorf("data-plane Authorization = %q, want %q", cap.dpAuthHeader, "Bearer "+dpAccessToken)
	}
	if cap.dpWorkspaceIDHeader != "ws-123" {
		t.Errorf("data-plane workspace id = %q, want %q", cap.dpWorkspaceIDHeader, "ws-123")
	}
	if cap.tokenAssertion != "cp-token" {
		t.Errorf("token exchange assertion = %q, want %q", cap.tokenAssertion, "cp-token")
	}
	if cap.tokenAuthDetails != dpAuthDetails {
		t.Errorf("token exchange authorization_details = %q, want %q", cap.tokenAuthDetails, dpAuthDetails)
	}
}

// TestServingEndpointWireProjection guards the minimal wire projection used for
// detection against a representative response body.
func TestServingEndpointWireProjection(t *testing.T) {
	body := servingEndpointGetPayload(t, true, "https://dp.example/invocations")
	var w servingEndpointDetailedWire
	if err := json.Unmarshal(body, &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := dataPlaneInfoFromWire(&w)
	want := &dataPlaneInfo{
		endpointURL:          "https://dp.example/invocations",
		authorizationDetails: dpAuthDetails,
	}
	if diff := cmp.Diff(want, got, cmp.AllowUnexported(dataPlaneInfo{})); diff != "" {
		t.Errorf("dataPlaneInfoFromWire mismatch (-want +got):\n%s", diff)
	}

	// Missing data-plane info projects to nil (not route-optimized).
	empty := servingEndpointGetPayload(t, false, "")
	var we servingEndpointDetailedWire
	if err := json.Unmarshal(empty, &we); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := dataPlaneInfoFromWire(&we); got != nil {
		t.Errorf("dataPlaneInfoFromWire on endpoint without data-plane info = %+v, want nil", got)
	}
}
