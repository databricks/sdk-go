package options

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/databricks/internal"
)

// ClientOption is an option to configure Databricks API clients.
type ClientOption func(*internal.ClientOptions) error

// WithHost returns a ClientOption to set the host for the client.
func WithHost(h string) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Host = h
		return nil
	}
}

// WithHTTPClient returns a ClientOption to use a specific HTTP Client when
// making HTTP requests.
//
// Important: When set, this option ignores all other options.
func WithHTTPClient(c *http.Client) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.HTTPClient = c
		return nil
	}
}

// WithCredentials returns a ClientOption to set a specific credentials.
func WithCredentials(c auth.Credentials) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Credentials = c
		return nil
	}
}

// WithTimeout returns a ClientOption to set the overall API call timeout to
// the given duration by default.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Timeout = d
		return nil
	}
}

// WithLogger returns a ClientOption to use the provided logger. Log messages
// are only logged if the logger is enabled.
func WithLogger(l *slog.Logger) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Logger = l
		return nil
	}
}

// WithoutResolution returns a ClientOption to disable profile resolution.
func WithoutResolution() ClientOption {
	return func(o *internal.ClientOptions) error {
		o.NoResolution = true
		return nil
	}
}

// WithProfile returns a ClientOption to set the profile name to use for
// default resolution.
//
// This option is ignored if resolution is disabled by WithoutResolution.
func WithProfile(profile string) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.Profile = profile
		return nil
	}
}

// WithProfileFile returns a ClientOption to set the path to the profile
// file to use for default resolution.
//
// This option is ignored if resolution is disabled by WithoutResolution.
func WithProfileFile(file string) ClientOption {
	return func(o *internal.ClientOptions) error {
		o.ProfileFile = file
		return nil
	}
}

// WithoutEnvResolution returns a ClientOption to disable environment
// variable resolution.
//
// This option is ignored if resolution is disabled by WithoutResolution.
func WithoutEnvResolution() ClientOption {
	return func(o *internal.ClientOptions) error {
		o.NoEnv = true
		return nil
	}
}
