// Package api provides utilities to make API calls against the Databricks API.
//
// This package draws inspiration from the AWS and GCP SDKs.
package api

import (
	"context"
	"time"
)

// Call represents a call to a Databricks API.
type Call func(context.Context) error

// Execute makes a call to a Databricks API using the given options.
func Execute(ctx context.Context, call Call, opts ...Option) error {
	options := Options{}
	for _, opt := range opts {
		if err := opt.Apply(&options); err != nil {
			return err
		}
	}
	return execute(ctx, call, options, sleep)
}

// WaitConfig configures the behavior of [ExecuteWait].
type WaitConfig struct {
	// Timeout for the entire wait operation. If zero, no timeout is applied
	// and the wait relies solely on the context's deadline.
	Timeout time.Duration

	// Backoff controls the delay between poll attempts. See [BackoffPolicy]
	// for the default values applied to zero-valued fields.
	Backoff BackoffPolicy
}

// WithOverrides returns a copy of the WaitConfig with non-zero fields from
// override applied on top. If override is nil, a copy of the original is
// returned unchanged.
func (c WaitConfig) WithOverrides(override *WaitConfig) WaitConfig {
	if override == nil {
		return c
	}
	if override.Timeout != 0 {
		c.Timeout = override.Timeout
	}
	if override.Backoff.Initial != 0 {
		c.Backoff.Initial = override.Backoff.Initial
	}
	if override.Backoff.Maximum != 0 {
		c.Backoff.Maximum = override.Backoff.Maximum
	}
	if override.Backoff.Factor != 0 {
		c.Backoff.Factor = override.Backoff.Factor
	}
	return c
}

// ExecuteWait polls a Databricks API until the call succeeds or a
// non-retriable error is encountered. The isRetriable predicate determines
// whether a given error should trigger another poll attempt.
//
// If cfg is nil, default values are used.
func ExecuteWait(ctx context.Context, call Call, isRetriable func(error) bool, cfg *WaitConfig) error {
	if cfg == nil {
		cfg = &WaitConfig{}
	}
	opts := Options{
		timeout: cfg.Timeout,
		retrier: func() Retrier {
			return RetryOn(cfg.Backoff, isRetriable)
		},
	}
	return execute(ctx, call, opts, sleep)
}

// sleep sleeps for the given duration. It is mostly equivalent to time.Sleep,
// but can be interrupted by ctx.Done() if the context completes before the
// duration elapses.
func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// sleeper is a convenience type for readability.
type sleeper func(ctx context.Context, d time.Duration) error

// execute is the actual implementation of Execute. Its purpose is to ease
// testing by providing a convenient way to mock the sleeping logic.
func execute(ctx context.Context, apiCall Call, opts Options, sleep sleeper) error {
	// Optionally update the context with the timeout. If the context already
	// has a deadline, that deadline is updated to the minimum of the context's
	// deadline and the timeout.
	if opts.timeout != 0 {
		c, f := context.WithTimeout(ctx, opts.timeout)
		defer f()
		ctx = c
	}

	// Get a new retrier for this specific execution. This is instantiated
	// lazilly if and if the first call execution returns an error.
	var retrier Retrier

	for {
		if opts.rateLimiter != nil { // no limiter == no wait
			if err := opts.rateLimiter.Wait(ctx); err != nil {
				return err
			}
		}

		err := apiCall(ctx)
		if err == nil {
			return nil // nothing to retry
		}

		if retrier == nil {
			if opts.retrier != nil {
				retrier = opts.retrier() // lazily instantiate the retrier
			}
			if retrier == nil {
				return err // no retrier == no retry
			}
		}
		delay, ok := retrier.IsRetriable(err)
		if !ok {
			return err // not retriable
		}
		if err := sleep(ctx, delay); err != nil {
			return err
		}
	}
}
