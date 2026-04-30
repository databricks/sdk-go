// Package call defines the options used to configure individual calls
// against the Databricks API.
package call

import (
	"time"

	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/options/internaloptions"
)

// Option configures a single call against the Databricks API.
type Option interface {
	Apply(*internaloptions.CallOptions) error
}

type optionFn func(*internaloptions.CallOptions) error

func (fn optionFn) Apply(c *internaloptions.CallOptions) error { return fn(c) }

// WithRetrier returns an Option that uses the given Retrier provider. If no
// retrier is provided, the call is not retried.
//
// The provider function must be thread-safe.
func WithRetrier(provider func() ops.Retrier) Option {
	return optionFn(func(c *internaloptions.CallOptions) error {
		c.Retrier = provider
		return nil
	})
}

// WithDisableRetry is a convenience option that disables retries.
func WithDisableRetry() Option {
	return optionFn(func(c *internaloptions.CallOptions) error {
		c.Retrier = nil
		return nil
	})
}

// WithTimeout returns an Option that sets the timeout for the call. If the
// context already has a deadline, it is updated to the minimum of the
// context's deadline and the timeout.
//
// The timeout covers the entire execution, including retries.
func WithTimeout(t time.Duration) Option {
	return optionFn(func(c *internaloptions.CallOptions) error {
		c.Timeout = t
		return nil
	})
}

// WithLimiter returns an Option that uses the given Limiter. If no limiter is
// provided, the call is not rate limited.
func WithLimiter(l ops.Limiter) Option {
	return optionFn(func(c *internaloptions.CallOptions) error {
		c.RateLimiter = l
		return nil
	})
}
