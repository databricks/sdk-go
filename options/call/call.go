// Package call defines the options used to configure individual calls
// against the Databricks API.
package call

import (
	"time"

	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/options/internaloptions"
)

// Option configures a single call against the Databricks API. If multiple
// options configure the same setting, the last option takes precedence.
type Option func(*internaloptions.CallOptions) error

// WithRetrier returns an Option that uses the given Retrier provider. If no
// retrier is provided, the call is not retried.
//
// The provider function must be thread-safe.
func WithRetrier(provider func() ops.Retrier) Option {
	return func(c *internaloptions.CallOptions) error {
		c.Retrier = provider
		return nil
	}
}

// WithDisableRetry is a convenience option that disables retries.
func WithDisableRetry() Option {
	return func(c *internaloptions.CallOptions) error {
		c.Retrier = nil
		return nil
	}
}

// WithTimeout returns an Option that limits how long a call can take,
// including retry delays and rate-limiter waits. The timeout applies until the
// method returns.
//
// An earlier deadline on the caller's context takes precedence.
// A zero duration disables the SDK timeout.
func WithTimeout(t time.Duration) Option {
	return func(c *internaloptions.CallOptions) error {
		c.Timeout = t
		return nil
	}
}

// WithLimiter returns an Option that uses the given Limiter. If no limiter is
// provided, the call is not rate limited.
func WithLimiter(l ops.Limiter) Option {
	return func(c *internaloptions.CallOptions) error {
		c.RateLimiter = l
		return nil
	}
}
