// Package client defines the options used to configure Databricks API
// clients.
//
// Databricks API clients can resolve their options from a variety of sources:
//
// - Profile from file or default.
// - Environment variables override.
// - Explicit value override.
//
// If no credentials are provided, the credentials are automatically resolved
// from the client configuration after the above chain of resolution.
//
// In practice, we recommend users to stick to a single resolution level.
// That is either rely purely on automatic resolution from the environment or
// explicitly set all options in code. Mixing the two approaches is not recommended.
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

// WithAccountID returns an Option that sets the account ID for the client.
func WithAccountID(id string) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.AccountID = id
		return nil
	}
}

// WithWorkspaceID returns an Option that sets the workspace ID for the client.
func WithWorkspaceID(id string) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.WorkspaceID = id
		return nil
	}
}

// WithoutConfigFile returns an Option that disables loading client
// configuration from a config file. It cannot be combined with
// [WithConfigFile] or [WithProfile].
func WithoutConfigFile() Option {
	return func(c *internaloptions.ClientOptions) error {
		c.DisableConfigFile = true
		return nil
	}
}

// WithoutEnv returns an Option that disables overlaying client configuration
// from environment variables. It does not affect DATABRICKS_CONFIG_FILE or
// DATABRICKS_CONFIG_PROFILE, which select the profile file and name.
func WithoutEnv() Option {
	return func(c *internaloptions.ClientOptions) error {
		c.DisableEnv = true
		return nil
	}
}

// WithConfigFile returns an Option that sets the config file to use for
// profile resolution. By default, the config file is resolved from the
// environment variable $DATABRICKS_CONFIG_FILE.
func WithConfigFile(file string) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.ConfigFile = file
		return nil
	}
}

// WithProfile returns an Option that sets the profile name to use for
// profile resolution. By default, the profile name is resolved from the
// environment variable $DATABRICKS_CONFIG_PROFILE.
func WithProfile(name string) Option {
	return func(c *internaloptions.ClientOptions) error {
		c.ProfileName = name
		return nil
	}
}
