package credentials

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/databricks/sdk-go/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	errClientIDRequired     = errors.New("client ID is required")
	errClientSecretRequired = errors.New("client secret is required")
	errHostRequired         = errors.New("host is required")
)

// M2MOptions configures a machine-to-machine (OAuth 2.0 client credentials)
// [auth.TokenProvider].
type M2MOptions struct {
	// ClientID is the OAuth client ID (service principal application ID).
	// Required.
	ClientID string

	// ClientSecret is the OAuth client secret. Required.
	ClientSecret string

	// Host is the Databricks workspace or account URL (for example,
	// "https://example.cloud.databricks.com"). Required.
	Host string

	// GroupID is the ID of the group whose role is assumed by the issued token.
	// When empty, no group role is assumed.
	GroupID string

	// Scopes overrides the OAuth scopes requested for the token. If empty,
	// defaults to ["all-apis"].
	Scopes []string

	// HTTPClient is the HTTP client used for OIDC endpoint discovery and
	// for the OAuth token exchange. If nil, [http.DefaultClient] is used.
	HTTPClient *http.Client
}

// NewM2MCredentials returns an [auth.TokenProvider] that fetches OAuth 2.0
// access tokens using the client credentials grant.
//
// The returned provider does not cache tokens or retry on failure. Wrap it
// with [auth.NewCachedTokenProvider] and [retrying.NewTokenProvider] as
// needed.
func NewM2MCredentials(opts M2MOptions) (auth.TokenProvider, error) {
	if opts.ClientID == "" {
		return nil, errClientIDRequired
	}
	if opts.ClientSecret == "" {
		return nil, errClientSecretRequired
	}
	if opts.Host == "" {
		return nil, errHostRequired
	}

	scopes := opts.Scopes
	if len(scopes) == 0 {
		scopes = []string{"all-apis"}
	}

	return &m2mTokenProvider{
		clientID:     opts.ClientID,
		clientSecret: opts.ClientSecret,
		host:         opts.Host,
		groupID:      opts.GroupID,
		scopes:       scopes,
		httpClient:   opts.HTTPClient,
	}, nil
}

type m2mTokenProvider struct {
	clientID     string
	clientSecret string
	host         string
	groupID      string
	scopes       []string
	httpClient   *http.Client

	// server holds the discovered authorization server metadata. Discovery
	// runs at most once per provider; concurrent first-discoveries are benign
	// because the result is immutable for a given host.
	server atomic.Pointer[oauthAuthorizationServer]
}

func (m *m2mTokenProvider) Token(ctx context.Context) (*auth.Token, error) {
	server, err := m.authorizationServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("oidc discovery: %w", err)
	}

	cfg := &clientcredentials.Config{
		ClientID:     m.clientID,
		ClientSecret: m.clientSecret,
		TokenURL:     server.TokenEndpoint,
		Scopes:       m.scopes,
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	if m.groupID != "" {
		cfg.EndpointParams = url.Values{"assume_group": {m.groupID}}
	}
	ot, err := fetchToken(ctx, m.httpClient, cfg)
	if err != nil {
		return nil, fmt.Errorf("fetch token: %w", err)
	}

	return &auth.Token{
		Value:  ot.AccessToken,
		Type:   ot.Type(),
		Expiry: ot.Expiry,
	}, nil
}

func (m *m2mTokenProvider) authorizationServer(ctx context.Context) (*oauthAuthorizationServer, error) {
	if s := m.server.Load(); s != nil {
		return s, nil
	}
	s, err := discoverAuthorizationServer(ctx, m.httpClient, m.host)
	if err != nil {
		return nil, err
	}
	m.server.Store(s)
	return s, nil
}
