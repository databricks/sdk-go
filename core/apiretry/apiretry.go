// Package apiretry provides retry utilities for Databricks API calls. It
// composes the generic retry machinery in [core/ops] with the error
// classification primitives in [core/apierr].
//
// This package is intended for internal use only.
package apiretry

import (
	"context"
	"errors"
	"net"
	"net/http"
	"slices"
	"strconv"
	"syscall"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/ops"
)

// RetrierConfig describes which errors [NewRetrier] retries.
//
// Error codes prevail over HTTP statuses. When ErrorCodes is non-empty and
// the error carries a canonical code that is not [codes.Unknown], ErrorCodes
// is consulted exclusively. Otherwise the HTTP status is examined against
// StatusCodes.
//
// Transient network errors are always retried regardless of rules. Retry-after
// hints (either passed in response header or via [apierr.RetryInfo]) are
// honored.
type RetrierConfig struct {
	// ErrorCodes are canonical codes considered retriable.
	ErrorCodes []codes.Code

	// StatusCodes are HTTP status codes considered retriable. Consulted only
	// when the error does not carry a canonical code, or when ErrorCodes is
	// empty.
	StatusCodes []int
}

// NewRetrier returns a [ops.Retrier] configured by rules. See [RetrierConfig]
// for the matching semantics.
//
// Important: the retrier has its own copy of the backoff policy which cannot
// be trivially reset by design. Users who need to reset the backoff policy
// should rather create a new retrier.
func NewRetrier(bp ops.BackoffPolicy, cfg RetrierConfig) ops.Retrier {
	return &retrier{backoff: bp, cfg: cfg}
}

type retrier struct {
	backoff ops.BackoffPolicy
	cfg     RetrierConfig
}

func (r *retrier) IsRetriable(err error) (time.Duration, bool) {
	if !r.isRetriable(err) {
		return 0, false
	}
	return max(r.backoff.Delay(), RetryDurationHint(err)), true
}

func (r *retrier) isRetriable(err error) bool {
	if IsTransientNetworkError(err) {
		return true
	}
	// Code-prevails, but only when ErrorCodes is configured. Otherwise callers
	// who pass only StatusCodes would have their known-code errors rejected
	// outright instead of falling through to status matching.
	if c := apierr.Code(err); len(r.cfg.ErrorCodes) > 0 && c != codes.Unknown {
		return slices.Contains(r.cfg.ErrorCodes, c)
	}
	if apiErr, ok := errors.AsType[*apierr.APIError](err); ok {
		return slices.Contains(r.cfg.StatusCodes, apiErr.HTTPStatusCode())
	}
	return false
}

// IsTransientNetworkError reports whether err represents a transient network
// condition that is likely to resolve on retry.
func IsTransientNetworkError(err error) bool {
	if err == nil {
		return false
	}
	// Context cancellation and deadline errors are intentionally treated as
	// non-retriable. The check must come before the net.Error check because
	// context.DeadlineExceeded satisfies net.Error with Timeout() == true.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNREFUSED) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

// RetryDurationHint returns a retry-duration hint extracted from err, if one
// is available. It examines, in order:
//
//   - the Retry-After header on an [*apierr.APIError]'s HTTP response;
//   - the [apierr.RetryInfo] detail on an [*apierr.APIError].
//
// Returns 0 when err does not carry a hint.
func RetryDurationHint(err error) time.Duration {
	apiErr, ok := errors.AsType[*apierr.APIError](err)
	if !ok {
		return 0
	}
	if d, ok := retryAfterFromHeader(apiErr.HTTPHeader()); ok {
		return d
	}
	if ri := apiErr.Details().RetryInfo; ri != nil && ri.RetryDelay > 0 {
		return ri.RetryDelay
	}
	return 0
}

// retryAfterFromHeader parses a Retry-After HTTP header value as either a
// delay in seconds or an HTTP-date. Returns (0, false) if the header is
// missing, malformed, or specifies a time in the past.
func retryAfterFromHeader(header http.Header) (time.Duration, bool) {
	if header == nil {
		return 0, false
	}
	v := header.Get("Retry-After")
	if v == "" {
		return 0, false
	}
	if seconds, parseErr := strconv.Atoi(v); parseErr == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second, true
	}
	if t, parseErr := http.ParseTime(v); parseErr == nil {
		if d := time.Until(t); d > 0 {
			return d, true
		}
	}
	return 0, false
}
