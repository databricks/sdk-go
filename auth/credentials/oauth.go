package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/databricks/sdk-go/core/apierr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OAuthAuthorizationServer describes the OAuth 2.0 endpoints advertised by a
// Databricks host through its well-known metadata document. Consumers should
// only rely on endpoints that are relevant to the flow they implement.
type oauthAuthorizationServer struct {
	// TokenEndpoint is the URL used to obtain OAuth 2.0 access tokens.
	TokenEndpoint string `json:"token_endpoint"`

	// AuthorizationEndpoint is the URL used to start user-facing OAuth
	// authorization flows. May be empty for hosts that do not advertise it.
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}

// DiscoverAuthorizationServer fetches OAuth authorization server metadata
// from {host}/oidc/.well-known/oauth-authorization-server.
//
// If httpClient is nil, [http.DefaultClient] is used.
func discoverAuthorizationServer(ctx context.Context, httpClient *http.Client, host string) (*oauthAuthorizationServer, error) {
	if host == "" {
		return nil, errHostRequired
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	url := strings.TrimRight(host, "/") + "/oidc/.well-known/oauth-authorization-server"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if err != nil {
			return nil, fmt.Errorf("GET %s: reading response body: %w", url, err)
		}
		// Convert the HTTP error to an apierr.APIError so that it can easily
		// be introspected by upstream retriers.
		err = apierr.FromHTTPError(resp.StatusCode, resp.Header, body)
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}

	var server oauthAuthorizationServer
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("GET %s: decoding response: %w", url, err)
	}
	if server.TokenEndpoint == "" {
		return nil, fmt.Errorf("GET %s: missing token_endpoint", url)
	}
	return &server, nil
}

// FetchToken performs an OAuth 2.0 token request using the given
// [clientcredentials.Config] and returns the result as an [auth.Token].
//
// This is a thin wrapper around cfg.Token(ctx) that takes care of injecting
// httpClient via the oauth2 library's context idiom and converting the result
// to [auth.Token]. It is shared by credential providers that use the client
// credentials or token-exchange grants (for example, M2M and Workload
// Identity Federation).
//
// If httpClient is nil, [http.DefaultClient] is used.
func fetchToken(ctx context.Context, httpClient *http.Client, cfg *clientcredentials.Config) (*oauth2.Token, error) {
	if httpClient != nil {
		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	}
	tok, err := cfg.Token(ctx)
	if err != nil {
		return nil, bridgeOAuth2Error(err)
	}
	return tok, nil
}

// bridgeOAuth2Error converts an *oauth2.RetrieveError into an *apierr.APIError
// so downstream code (retry classification, error inspection) only needs to
// understand the apierr taxonomy. Returns err unchanged if it is not an
// *oauth2.RetrieveError with an HTTP response attached.
func bridgeOAuth2Error(err error) error {
	re, ok := errors.AsType[*oauth2.RetrieveError](err)
	if !ok || re.Response == nil {
		return err
	}
	return apierr.FromHTTPError(re.Response.StatusCode, re.Response.Header, re.Body)
}
