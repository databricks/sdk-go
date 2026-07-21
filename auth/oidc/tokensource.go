package oidc

import (
	"context"
	"errors"
	"net/url"

	"github.com/databricks/databricks-sdk-go/credentials/u2m"
	"github.com/databricks/sdk-go/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// DatabricksOIDCTokenProviderConfig is the configuration for a Databricks OIDC
// TokenProvider.
type DatabricksOIDCTokenProviderConfig struct {
	// ClientID of the Databricks OIDC application. It corresponds to the
	// Application ID of the Databricks Service Principal.
	//
	// This field is only required for Workload Identity Federation and should
	// be empty for Account-wide token federation.
	ClientID string

	// AccountID is the account ID of the Databricks Account. This field is
	// only required for Account-wide token federation.
	AccountID string

	// Host is the host of the Databricks account or workspace.
	Host string

	// TokenEndpointProvider returns the token endpoint for the Databricks OIDC
	// application.
	TokenEndpointProvider func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error)

	// Audience is the audience of the Databricks OIDC application.
	// This is only used for Workspace level tokens.
	Audience string

	// IDTokenProvider returns the IDToken to be used for the token exchange.
	IDTokenProvider IDTokenProvider
}

// NewDatabricksOIDCTokenProvider returns a new Databricks OIDC TokenProvider.
func NewDatabricksOIDCTokenProvider(cfg DatabricksOIDCTokenProviderConfig) auth.TokenProvider {
	return &databricksOIDCTokenProvider{cfg: cfg}
}

// databricksOIDCTokenProvider is a auth.TokenProvider which exchanges a token using
// Workload Identity Federation.
type databricksOIDCTokenProvider struct {
	cfg DatabricksOIDCTokenProviderConfig
}

// Token implements [TokenProvider.Token]
func (w *databricksOIDCTokenProvider) Token(ctx context.Context) (*auth.Token, error) {
	if w.cfg.Host == "" {
		return nil, errors.New("missing Host")
	}
	endpoints, err := w.cfg.TokenEndpointProvider(ctx)
	if err != nil {
		return nil, err
	}

	// if w.cfg.ClientID == "" {
	// 	logger.Debugf(ctx, "No ClientID provided, authenticating with Account-wide token federation")
	// } else {
	// 	logger.Debugf(ctx, "Client ID provided, authenticating with Workload Identity Federation")
	// }

	// TODO: The audience is a concept of the IDToken that should likely be
	// configured when the IDTokenProvider is created.
	audience := w.determineAudience(endpoints)
	idToken, err := w.cfg.IDTokenProvider.IDToken(ctx, audience)
	if err != nil {
		return nil, err
	}

	c := &clientcredentials.Config{
		ClientID:  w.cfg.ClientID,
		AuthStyle: oauth2.AuthStyleInParams,
		TokenURL:  endpoints.TokenEndpoint,
		Scopes:    []string{"all-apis"},
		EndpointParams: url.Values{
			"subject_token_type": {"urn:ietf:params:oauth:token-type:jwt"},
			"subject_token":      {idToken.Value},
			"grant_type":         {"urn:ietf:params:oauth:grant-type:token-exchange"},
		},
	}
	token, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}
	return &auth.Token{
		Value:  token.AccessToken,
		Type:   token.Type(),
		Expiry: token.Expiry,
	}, nil
}

func (w *databricksOIDCTokenProvider) determineAudience(endpoints *u2m.OAuthAuthorizationServer) string {
	if w.cfg.Audience != "" {
		return w.cfg.Audience
	}
	// For Databricks Accounts, the account id is the default audience.
	if w.cfg.AccountID != "" {
		return w.cfg.AccountID
	}
	// For Databricks Workspaces, the auth endpoint is the default audience.
	return endpoints.TokenEndpoint
}
