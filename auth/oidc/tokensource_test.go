package oidc

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/databricks/databricks-sdk-go/credentials/u2m"
	"github.com/databricks/databricks-sdk-go/httpclient/fixtures"
	"github.com/databricks/sdk-go/auth"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
)

func hasPrefix(err error, prefix string) bool {
	return strings.HasPrefix(err.Error(), prefix)
}

func TestDatabricksOidcTokenProvider(t *testing.T) {
	testCases := []struct {
		desc                 string
		clientID             string
		accountID            string
		host                 string
		tokenAudience        string
		httpTransport        http.RoundTripper
		oidcEndpointProvider func(context.Context) (*u2m.OAuthAuthorizationServer, error)
		idToken              string
		wantAudience         string
		tokenProviderError   error
		wantToken            string
		wantErrPrefix        *string
	}{
		{
			desc:          "missing host",
			clientID:      "client-id",
			tokenAudience: "token-audience",
			wantErrPrefix: new("missing Host"),
		},
		{
			desc: "token provider error",

			clientID:      "client-id",
			host:          "http://host.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			wantAudience:       "token-audience",
			tokenProviderError: errors.New("error getting id token"),
			wantErrPrefix:      new("error getting id token"),
		},
		{
			desc:          "databricks workspace server error",
			clientID:      "client-id",
			host:          "http://host.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {
					Status: http.StatusInternalServerError,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
				},
			},
			wantAudience:  "token-audience",
			idToken:       "id-token-42",
			wantErrPrefix: new("oauth2: cannot fetch token: Internal Server Error"),
		},
		{
			desc:          "invalid auth token",
			clientID:      "client-id",
			host:          "http://host.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {
					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Response: map[string]string{
						"foo": "bar",
					},
				},
			},
			wantAudience:  "token-audience",
			idToken:       "id-token-42",
			wantErrPrefix: new("oauth2: server response missing access_token"),
		},
		{
			desc:          "success WIF workspace",
			clientID:      "client-id",
			host:          "http://host.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {

					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					ExpectedRequest: url.Values{
						"client_id":          {"client-id"},
						"scope":              {"all-apis"},
						"subject_token_type": {"urn:ietf:params:oauth:token-type:jwt"},
						"subject_token":      {"id-token-42"},
						"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
					},
					Response: map[string]string{
						"token_type":    "access-token",
						"access_token":  "test-auth-token",
						"refresh_token": "refresh",
						"expires_on":    "0",
					},
				},
			},
			wantAudience: "token-audience",
			idToken:      "id-token-42",
			wantToken:    "test-auth-token",
		},
		{
			desc:          "success WIF account",
			clientID:      "client-id",
			accountID:     "ac123",
			host:          "https://accounts.databricks.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {
					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					ExpectedRequest: url.Values{
						"client_id":          {"client-id"},
						"scope":              {"all-apis"},
						"subject_token_type": {"urn:ietf:params:oauth:token-type:jwt"},
						"subject_token":      {"id-token-42"},
						"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
					},
					Response: map[string]string{
						"token_type":    "access-token",
						"access_token":  "test-auth-token",
						"refresh_token": "refresh",
						"expires_on":    "0",
					},
				},
			},
			wantAudience: "token-audience",
			idToken:      "id-token-42",
			wantToken:    "test-auth-token",
		},
		{
			desc:      "default token audience account",
			clientID:  "client-id",
			accountID: "ac123",
			host:      "https://accounts.databricks.com",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {
					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Response: map[string]string{
						"token_type":    "access-token",
						"access_token":  "test-auth-token",
						"refresh_token": "refresh",
						"expires_on":    "0",
					},
				},
			},
			wantAudience: "ac123",
			idToken:      "id-token-42",
			wantToken:    "test-auth-token",
		},
		{
			desc:     "default token audience workspace",
			clientID: "client-id",
			host:     "https://host.com",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {
					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Response: map[string]string{
						"token_type":    "access-token",
						"access_token":  "test-auth-token",
						"refresh_token": "refresh",
						"expires_on":    "0",
					},
				},
			},
			wantAudience: "https://host.com/oidc/v1/token",
			idToken:      "id-token-42",
			wantToken:    "test-auth-token",
		},
		{
			desc:          "success account-wide",
			host:          "http://host.com",
			tokenAudience: "token-audience",
			oidcEndpointProvider: func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
				return &u2m.OAuthAuthorizationServer{
					TokenEndpoint: "https://host.com/oidc/v1/token",
				}, nil
			},
			httpTransport: fixtures.MappingTransport{
				"POST /oidc/v1/token": {

					Status: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					ExpectedRequest: url.Values{
						"scope":              {"all-apis"},
						"subject_token_type": {"urn:ietf:params:oauth:token-type:jwt"},
						"subject_token":      {"id-token-42"},
						"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
					},
					Response: map[string]string{
						"token_type":    "access-token",
						"access_token":  "test-auth-token",
						"refresh_token": "refresh",
						"expires_on":    "0",
					},
				},
			},
			wantAudience: "token-audience",
			idToken:      "id-token-42",
			wantToken:    "test-auth-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var gotAudience string // set when IDTokenProvider is called
			cfg := DatabricksOIDCTokenProviderConfig{
				ClientID:              tc.clientID,
				AccountID:             tc.accountID,
				Host:                  tc.host,
				TokenEndpointProvider: tc.oidcEndpointProvider,
				Audience:              tc.tokenAudience,
				IDTokenProvider: IDTokenProviderFn(func(ctx context.Context, aud string) (*IDToken, error) {
					gotAudience = aud
					return &IDToken{Value: tc.idToken}, tc.tokenProviderError
				}),
			}

			ts := NewDatabricksOIDCTokenProvider(cfg)
			if tc.httpTransport != nil {
				ts.(*databricksOIDCTokenProvider).cfg.TokenEndpointProvider = func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
					return &u2m.OAuthAuthorizationServer{
						TokenEndpoint: "https://host.com/oidc/v1/token",
					}, nil
				}
			}

			ctx := context.Background()
			if tc.httpTransport != nil {
				ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
					Transport: tc.httpTransport,
				})
			}

			token, err := ts.Token(ctx)
			if tc.wantErrPrefix == nil && err != nil {
				t.Errorf("Token(ctx): got error %q, want none", err)
			}
			if tc.wantErrPrefix != nil && !hasPrefix(err, *tc.wantErrPrefix) {
				t.Errorf("Token(ctx): got error %q, want error with prefix %q", err, *tc.wantErrPrefix)
			}
			if tc.wantAudience != gotAudience {
				t.Errorf("mockTokenProvider: got audience %s, want %s", gotAudience, tc.wantAudience)
			}
			tokenValue := ""
			if token != nil {
				tokenValue = token.Value
			}
			if diff := cmp.Diff(tc.wantToken, tokenValue); diff != "" {
				t.Errorf("Authenticate(): mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type countingRoundTripper struct {
	transport http.RoundTripper
	calls     int
}

// RoundTrip records how many token exchanges reach the underlying transport.
func (t *countingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	t.calls++
	return t.transport.RoundTrip(r)
}

// newTokenExchangeRequest returns the standard WIF token exchange form with
// the optional client and group identifiers used by a test scenario.
func newTokenExchangeRequest(clientID, groupID string) url.Values {
	request := url.Values{
		"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
		"scope":              {"all-apis"},
		"subject_token":      {"id-token"},
		"subject_token_type": {"urn:ietf:params:oauth:token-type:jwt"},
	}
	if clientID != "" {
		request.Set("client_id", clientID)
	}
	if groupID != "" {
		request.Set("assume_group", groupID)
	}

	return request
}

// contextWithOAuthTransport routes OAuth requests through the supplied test
// transport instead of the default HTTP client.
func contextWithOAuthTransport(transport http.RoundTripper) context.Context {
	return context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: transport,
	})
}

// newTestDatabricksOIDCProvider supplies fixed external ID token plumbing so
// tests can focus on the Databricks token exchange.
func newTestDatabricksOIDCProvider(tokenEndpoint, clientID, groupID string) auth.TokenProvider {
	return NewDatabricksOIDCTokenProvider(DatabricksOIDCTokenProviderConfig{
		ClientID: clientID,
		Host:     "https://host.com",
		GroupID:  groupID,
		TokenEndpointProvider: func(context.Context) (*u2m.OAuthAuthorizationServer, error) {
			return &u2m.OAuthAuthorizationServer{TokenEndpoint: tokenEndpoint}, nil
		},
		IDTokenProvider: IDTokenProviderFn(func(context.Context, string) (*IDToken, error) {
			return &IDToken{Value: "id-token"}, nil
		}),
	})
}

// TestDatabricksOIDCGroupCachesAreIsolated verifies that cached WIF providers
// configured for different groups retain the token issued for their own group.
func TestDatabricksOIDCGroupCachesAreIsolated(t *testing.T) {
	testCases := []struct {
		name      string
		groupID   string
		wantToken string
	}{
		{
			name:      "no group",
			wantToken: "normal-token",
		},
		{
			name:      "group A",
			groupID:   "group-a",
			wantToken: "group-a-token",
		},
		{
			name:      "group B",
			groupID:   "group-b",
			wantToken: "group-b-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transport := &countingRoundTripper{
				transport: fixtures.MappingTransport{
					"POST /oidc/v1/token": {
						Status:          http.StatusOK,
						ExpectedRequest: newTokenExchangeRequest("client-id", tc.groupID),
						Response: map[string]any{
							"access_token": tc.wantToken,
							"token_type":   "Bearer",
							"expires_in":   3600,
						},
					},
				},
			}
			ctx := contextWithOAuthTransport(transport)
			provider := newTestDatabricksOIDCProvider(
				"https://host.com/oidc/v1/token",
				"client-id",
				tc.groupID,
			)

			cached := auth.NewCachedTokenProvider(provider, auth.WithAsyncRefresh(false))
			for range 2 {
				token, err := cached.Token(ctx)
				if err != nil {
					t.Fatalf("Token() error = %v", err)
				}

				if token.Value != tc.wantToken {
					t.Errorf("Token().Value = %q, want %q", token.Value, tc.wantToken)
				}
			}

			if got, want := transport.calls, 1; got != want {
				t.Errorf("token exchanges = %d, want %d", got, want)
			}
		})
	}
}

// TestDatabricksOIDCAssumeGroupForEveryEndpointType verifies that WIF sends
// group assumption to workspace, shared account or unified, and account-wide
// exchanges.
func TestDatabricksOIDCAssumeGroupForEveryEndpointType(t *testing.T) {
	testCases := []struct {
		name      string
		tokenPath string
		clientID  string
	}{
		{
			name:      "workspace",
			tokenPath: "/oidc/v1/token",
			clientID:  "client-id",
		},
		{
			name:      "account or unified",
			tokenPath: "/oidc/accounts/account-id/v1/token",
			clientID:  "client-id",
		},
		{
			name:      "account-wide",
			tokenPath: "/oidc/accounts/account-id/v1/token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transport := fixtures.MappingTransport{
				"POST " + tc.tokenPath: {
					Status:          http.StatusOK,
					ExpectedRequest: newTokenExchangeRequest(tc.clientID, "group-id"),
					Response: map[string]any{
						"access_token": "access-token",
						"token_type":   "Bearer",
						"expires_in":   3600,
					},
				},
			}
			ctx := contextWithOAuthTransport(transport)
			provider := newTestDatabricksOIDCProvider(
				"https://host.com"+tc.tokenPath,
				tc.clientID,
				"group-id",
			)

			if _, err := provider.Token(ctx); err != nil {
				t.Fatalf("Token() error = %v", err)
			}
		})
	}
}

// TestDatabricksOIDCGroupRejectionIsNotRetried verifies that a rejected grouped
// exchange is returned to the caller without another token request.
func TestDatabricksOIDCGroupRejectionIsNotRetried(t *testing.T) {
	transport := &countingRoundTripper{
		transport: fixtures.MappingTransport{
			"POST /oidc/v1/token": {
				Status:          http.StatusBadRequest,
				ExpectedRequest: newTokenExchangeRequest("client-id", "group-id"),
				Response: map[string]string{
					"error":             "invalid_request",
					"error_description": "assume_group is not supported",
				},
			},
		},
	}
	ctx := contextWithOAuthTransport(transport)
	provider := newTestDatabricksOIDCProvider(
		"https://host.com/oidc/v1/token",
		"client-id",
		"group-id",
	)

	_, err := provider.Token(ctx)
	if err == nil || !strings.Contains(err.Error(), "invalid_request") || !strings.Contains(err.Error(), "assume_group is not supported") {
		t.Fatalf("Token() error = %v", err)
	}

	if got, want := transport.calls, 1; got != want {
		t.Errorf("token requests = %d, want %d", got, want)
	}
}

// TestDatabricksOIDCExpiredTokenRetainsGroup verifies that every refresh of an
// expired WIF token retains the configured group.
func TestDatabricksOIDCExpiredTokenRetainsGroup(t *testing.T) {
	transport := &countingRoundTripper{
		transport: fixtures.MappingTransport{
			"POST /oidc/v1/token": {
				Status:          http.StatusOK,
				ExpectedRequest: newTokenExchangeRequest("client-id", "group-id"),
				Response: map[string]any{
					"access_token": "expired-token",
					"token_type":   "Bearer",
					"expires_in":   -1,
				},
			},
		},
	}
	ctx := contextWithOAuthTransport(transport)
	provider := newTestDatabricksOIDCProvider(
		"https://host.com/oidc/v1/token",
		"client-id",
		"group-id",
	)

	cached := auth.NewCachedTokenProvider(provider, auth.WithAsyncRefresh(false))
	for range 2 {
		if _, err := cached.Token(ctx); err != nil {
			t.Fatalf("Token() error = %v", err)
		}
	}

	if got, want := transport.calls, 2; got != want {
		t.Fatalf("token exchanges = %d, want %d", got, want)
	}
}
