// Package dataplane is an experimental package that provides a token source to
// directly access Databricks data plane.
package dataplane

import (
	"context"
	"sync"

	"github.com/databricks/sdk-go/auth"
)

// OAuthClient is an interface for Databricks OAuth client.
type OAuthClient interface {
	GetOAuthToken(ctx context.Context, authDetails string, t *auth.Token) (*auth.Token, error)
}

// EndpointTokenProvider is anything that returns tokens given a data plane
// endpoint and authentication details.
type EndpointTokenProvider interface {
	Token(ctx context.Context, endpoint string, authDetails string) (*auth.Token, error)
}

// NewEndpointTokenProvider returns a new EndpointTokenProvider that uses the given
// OAuthClient and control plane TokenProvider.
func NewEndpointTokenProvider(c OAuthClient, cpts auth.TokenProvider) *dataPlaneTokenProvider {
	return &dataPlaneTokenProvider{
		client: c,
		cpts: auth.NewCachedTokenProvider(
			cpts,
		),
	}
}

type tsKey struct {
	endpoint    string
	authDetails string
}

// dataPlaneTokenProvider implements the EndpointTokenProvider interface.
//
// For a given endpoint and authentication details, it uses the control plane
// TokenProvider to retrieve the control plane token, that is then used to
// retrieve the data plane token through the OAuthClient.
//
// Each token source is cached to avoid unnecessary token requests.
type dataPlaneTokenProvider struct {
	client  OAuthClient
	cpts    auth.TokenProvider
	sources sync.Map
}

// Token returns a token for the given endpoint and authentication details.
func (dpts *dataPlaneTokenProvider) Token(ctx context.Context, endpoint string, authDetails string) (*auth.Token, error) {
	key := tsKey{endpoint: endpoint, authDetails: authDetails}

	if a, ok := dpts.sources.Load(key); ok { // happy path
		return a.(auth.TokenProvider).Token(ctx)
	}

	ts := auth.NewCachedTokenProvider(
		&TokenProvider{
			client:      dpts.client,
			cpts:        dpts.cpts,
			authDetails: authDetails,
		},
	)
	dpts.sources.Store(key, ts)

	return ts.Token(ctx)
}

type TokenProvider struct {
	client      OAuthClient
	cpts        auth.TokenProvider //
	authDetails string
}

func (dpts *TokenProvider) Token(ctx context.Context) (*auth.Token, error) {
	innerToken, err := dpts.cpts.Token(ctx)
	if err != nil {
		return nil, err
	}
	return dpts.client.GetOAuthToken(ctx, dpts.authDetails, innerToken)
}
