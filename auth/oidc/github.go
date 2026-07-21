package oidc

import (
	"context"
	"errors"
	"fmt"

	"github.com/databricks/databricks-sdk-go/httpclient"
	"github.com/databricks/databricks-sdk-go/logger"
)

// NewGithubIDTokenProvider returns a new IDTokenProvider that retrieves an IDToken
// from the Github Actions environment. This IDTokenProvider is only valid when
// running in Github Actions with OIDC enabled.
func NewGithubIDTokenProvider(client *httpclient.ApiClient, actionsIDTokenRequestURL, actionsIDTokenRequestToken string) IDTokenProvider {
	return &githubIDTokenProvider{
		actionsIDTokenRequestURL:   actionsIDTokenRequestURL,
		actionsIDTokenRequestToken: actionsIDTokenRequestToken,
		refreshClient:              client,
	}
}

// githubIDTokenProvider retrieves JWT Tokens from Github Actions.
type githubIDTokenProvider struct {
	actionsIDTokenRequestURL   string
	actionsIDTokenRequestToken string
	refreshClient              *httpclient.ApiClient
}

// IDToken returns a JWT Token for the specified audience. It will return
// an error if not running in GitHub Actions.
func (g *githubIDTokenProvider) IDToken(ctx context.Context, audience string) (*IDToken, error) {
	if g.actionsIDTokenRequestURL == "" {
		logger.Debugf(ctx, "Missing ActionsIDTokenRequestURL, likely not calling from a Github action")
		return nil, errors.New("missing ActionsIDTokenRequestURL")
	}
	if g.actionsIDTokenRequestToken == "" {
		logger.Debugf(ctx, "Missing ActionsIDTokenRequestToken, likely not calling from a Github action")
		return nil, errors.New("missing ActionsIDTokenRequestToken")
	}

	resp := &IDToken{}
	requestUrl := g.actionsIDTokenRequestURL
	if audience != "" {
		requestUrl = fmt.Sprintf("%s&audience=%s", requestUrl, audience)
	}
	err := g.refreshClient.Do(ctx, "GET", requestUrl,
		httpclient.WithRequestHeader("Authorization", fmt.Sprintf("Bearer %s", g.actionsIDTokenRequestToken)),
		httpclient.WithResponseUnmarshal(resp),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to request ID token from %s: %w", g.actionsIDTokenRequestURL, err)
	}

	return resp, nil
}
