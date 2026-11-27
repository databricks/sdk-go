// Package ops provides utilities to execute operations with retry, timeout,
// and rate limiting.
//
// This package draws inspiration from the AWS and GCP SDKs.
package ops

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Execute executes operation op with the given options. The operation is
// expected to honor cancellation and deadlines on the context passed to it.
//
// When the retrier stops, Execute returns the error of the last attempt. When
// the context ends, or the limiter fails, after a failed attempt, Execute
// returns that cause, and the message also names the error of the last
// attempt. [LastError] returns the error of the last attempt.
func Execute(ctx context.Context, op func(context.Context) error, opts ...Option) error {
	options := Options{}
	for _, opt := range opts {
		if err := opt.Apply(&options); err != nil {
			return err
		}
	}
	return execute(ctx, op, options, sleep)
}

// ExecuteWithResult returns the result of calling op with the given options.
// It is a convenience wrapper around Execute for operations that return a
// value. In case of error, the zero value of T is returned.
func ExecuteWithResult[T any](ctx context.Context, op func(context.Context) (T, error), opts ...Option) (T, error) {
	var result T
	err := Execute(ctx, func(ctx context.Context) error {
		var err error
		result, err = op(ctx)
		return err
	}, opts...)
	if err != nil {
		var zero T // guarantee zero value on error
		return zero, err
	}
	return result, nil
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
func execute(ctx context.Context, fn func(context.Context) error, opts Options, sleep sleeper) error {
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

	// Error of the last failed attempt, if any.
	var lastErr error

	for {
		if opts.rateLimiter != nil { // no limiter == no wait
			if err := opts.rateLimiter.Wait(ctx); err != nil {
				return withLastAttempt(err, lastErr)
			}
		}

		err := fn(ctx)
		if err == nil {
			return nil // nothing to retry
		}
		lastErr = err

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
			return withLastAttempt(err, lastErr)
		}
	}
}

// withLastAttempt returns cause, which ended the execution before its next
// attempt, together with lastErr, the error of the last attempt, if any. A
// bare context error would hide what failed before the time ran out.
func withLastAttempt(cause, lastErr error) error {
	if lastErr == nil {
		return cause
	}
	return &interruptedError{cause: cause, lastErr: lastErr}
}

// interruptedError is an execution that ended before its next attempt.
type interruptedError struct {
	cause   error // what ended the execution: a context or a limiter error
	lastErr error // error of the last attempt
}

func (e *interruptedError) Error() string {
	return fmt.Sprintf("%v; last attempt: %v", e.cause, e.lastErr)
}

// Unwrap returns the cause only. An execution that ran out of time must not
// match as the error of its last attempt: a caller, or the retrier of an outer
// execution, would otherwise handle a timeout as, for example, an HTTP 503.
func (e *interruptedError) Unwrap() error {
	return e.cause
}

// LastError returns the error of the last attempt of an execution that ended
// because its context ended, or its limiter failed, after that attempt. It
// returns nil for every other error.
//
// errors.Is and errors.As do not find the error of the last attempt in err;
// LastError is the only way to read it.
func LastError(err error) error {
	if e, ok := errors.AsType[*interruptedError](err); ok {
		return e.lastErr
	}
	return nil
}
