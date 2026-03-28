package ops

import (
	"time"
)

// Option is an option used by Execute to control the behavior of an operation.
// An Option essentially acts as a convenient way to configure Options.
type Option interface {
	Apply(*Options) error
}

type optionFn func(*Options)

func (fn optionFn) Apply(opts *Options) error {
	fn(opts)
	return nil
}

// WithRetrier configures to use the given Retrier provider. If no retrier is
// provided, the operation is not retried.
//
// The provider function must be thread-safe.
func WithRetrier(provider func() Retrier) Option {
	return optionFn(func(o *Options) {
		o.retrier = provider
	})
}

// WithDisableRetry is a convenience option to disable retries.
func WithDisableRetry() Option {
	return optionFn(func(o *Options) {
		o.retrier = nil
	})
}

// WithTimeout sets the timeout duration for the operation. If the context
// already has a deadline, it is updated to the minimum of the context's
// deadline and the timeout.
//
// The timeout covers the entire execution, including retries.
func WithTimeout(t time.Duration) Option {
	return optionFn(func(o *Options) {
		o.timeout = t
	})
}

// WithLimiter configures to use the given Limiter. If no limiter is provided,
// the operation is not rate limited.
func WithLimiter(l Limiter) Option {
	return optionFn(func(o *Options) {
		o.rateLimiter = l
	})
}

// Options controls the behavior of an operation.
type Options struct {
	// Provides a new Retrier for each Execute call. The function must be
	// thread-safe. The retrier must be fresh within the context of an
	// Execute call (e.g. no need to reset a BackoffStrategy).
	retrier func() Retrier

	rateLimiter Limiter

	timeout time.Duration
}
