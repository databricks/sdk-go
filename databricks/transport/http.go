package transport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/databricks/internal"
	"github.com/databricks/sdk-go/databricks/options"
)

// NewHTTPClient creates a new HTTP client with the given options.
func NewHTTPClient(ctx context.Context, opts ...options.ClientOption) (*http.Client, error) {
	copts := &internal.ClientOptions{}
	for _, opt := range opts {
		if err := opt(copts); err != nil {
			return nil, err
		}
	}
	if err := copts.Initialize(); err != nil {
		return nil, err
	}

	// If an HTTP client is provided, use it as is without any additional
	// configuration.
	if copts.HTTPClient != nil {
		return copts.HTTPClient, nil
	}

	if copts.Credentials == nil {
		// TODO: Load default credentials from profile
		return nil, errors.New("no credentials provided")
	}

	transport := &authTransport{
		base:   http.DefaultTransport,
		creds:  copts.Credentials,
		logger: copts.Logger,
	}

	return &http.Client{
		Timeout:   copts.Timeout,
		Transport: transport,
	}, nil
}

type authTransport struct {
	base   http.RoundTripper // base transport to wrap
	creds  auth.Credentials  // credentials to use for authentication
	logger *slog.Logger
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	headers, err := t.creds.AuthHeaders(req.Context())
	if err != nil {
		// RoundTripper must always close the request, including on errors.
		if req.Body != nil {
			closeErr := req.Body.Close()
			// Swallow the cleanup error; the credentials error is the primary
			// failure and the one that is the most actionable for callers.
			if closeErr != nil {
				t.logger.ErrorContext(req.Context(), "error closing request body", "error", closeErr)
			}
		}
		return nil, err
	}
	// RoundTripper must not modify the request, except for consuming and
	// closing the Request's Body.
	clone := req.Clone(req.Context())
	for _, header := range headers {
		clone.Header.Add(header.Key, header.Value)
	}
	return t.base.RoundTrip(clone)
}
