package options

import (
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/log"
	"github.com/databricks/sdk-go/databricks/internal"
)

// ClientOption is an option to configure Databricks API clients.
type ClientOption func(*internal.ClientOptions) error

// WithHTTPClient returns a ClientOption to use a specific
// HTTP Client when making HTTP requests.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.HTTPClient = c
		return nil
	}
}

// WithCredentials returns a ClientOption to set a
// specific credentials.
func WithCredentials(c auth.Credentials) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Credentials = c
		return nil
	}
}

// WithTimeout returns a ClientOption to set the overall
// API call timeout to the given duration by default.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Timeout = d
		return nil
	}
}

// WithLogger returns a ClientOption to use the provided
// logger instead of the default logger.
func WithLogger(l log.Logger) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Logger = l
		return nil
	}
}
