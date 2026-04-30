// Package client defines the options used to configure Databricks API
// clients.
package client

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/options/internaloptions"
)

// Option configures a Databricks API client.
type Option func(*internaloptions.ClientOptions) error

// WithHost returns an Option that sets the host for the client.
func WithHost(h string) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.Host = h
		return nil
	}
}

// WithHTTPClient returns an Option that uses a specific HTTP client when
// making HTTP requests.
//
// Important: When set, this option ignores all other options.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.HTTPClient = hc
		return nil
	}
}

// WithCredentials returns an Option that sets a specific credentials.
func WithCredentials(creds auth.Credentials) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.Credentials = creds
		return nil
	}
}

// WithTimeout returns an Option that sets the overall API call timeout to
// the given duration by default.
func WithTimeout(d time.Duration) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.Timeout = d
		return nil
	}
}

// WithLogger returns an Option that uses the provided logger. Log messages
// are only logged if the logger is enabled.
func WithLogger(l *slog.Logger) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.Logger = l
		return nil
	}
}
