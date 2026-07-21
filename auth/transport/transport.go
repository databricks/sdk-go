// Package transport provides a HTTP transport that automatically adds authentication
// headers to outgoing requests. It is meant to provide a convenient way to make
// authenticated requests against Databricks APIs that are not part of the SDKs.
package transport

import (
	"net/http"

	"github.com/databricks/sdk-go/auth"
)

// NewAuthTransport returns a new HTTP transport that wraps the base transport
// to automatically add authentication headers to outgoing requests.
//
// The returned transport is safe for concurrent use by multiple goroutines.
// If base is nil, the default transport is used. The function assumes that
// the given credentials are non-nil.
func NewAuthTransport(base http.RoundTripper, creds auth.Credentials) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &authTransport{base: base, creds: creds}
}

// authTransport is the implementation of the HTTP transport that adds
// authentication headers to outgoing requests.
type authTransport struct {
	base  http.RoundTripper // base transport to wrap
	creds auth.Credentials  // credentials to use for authentication
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	headers, err := t.creds.AuthHeaders(req.Context())
	if err != nil {
		// RoundTripper must always close the request, including on errors.
		if req.Body != nil {
			// Swallow the cleanup error; the credentials error is the primary
			// failure and the one that is the most actionable for callers.
			_ = req.Body.Close()
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
