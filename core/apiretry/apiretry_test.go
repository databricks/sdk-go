package apiretry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/ops"
)

type fakeNetError struct {
	msg     string
	timeout bool
}

func (e *fakeNetError) Error() string   { return e.msg }
func (e *fakeNetError) Timeout() bool   { return e.timeout }
func (e *fakeNetError) Temporary() bool { return false }

func TestIsTransientNetworkError(t *testing.T) {
	testCases := []struct {
		desc string
		err  error
		want bool
	}{
		{desc: "nil", err: nil, want: false},
		{desc: "plain error", err: errors.New("boom"), want: false},
		{
			desc: "net timeout",
			err:  &fakeNetError{msg: "i/o timeout", timeout: true},
			want: true,
		},
		{
			desc: "net non-timeout",
			err:  &fakeNetError{msg: "connection refused", timeout: false},
			want: false,
		},
		{
			desc: "wrapped net timeout",
			err:  fmt.Errorf("dial: %w", &fakeNetError{msg: "i/o timeout", timeout: true}),
			want: true,
		},
		{desc: "context canceled", err: context.Canceled, want: false},
		{desc: "context deadline exceeded", err: context.DeadlineExceeded, want: false},
		{
			desc: "wrapped context deadline",
			err:  fmt.Errorf("op: %w", context.DeadlineExceeded),
			want: false,
		},
		{desc: "ECONNRESET", err: syscall.ECONNRESET, want: true},
		{desc: "ECONNREFUSED", err: syscall.ECONNREFUSED, want: true},
		{
			desc: "wrapped ECONNRESET",
			err:  fmt.Errorf("dial: %w", syscall.ECONNRESET),
			want: true,
		},
		{
			desc: "EPIPE (not in the allowlist)",
			err:  syscall.EPIPE,
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := IsTransientNetworkError(tc.err); got != tc.want {
				t.Errorf("isTransientNetworkError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestRetryAfterFromHeader(t *testing.T) {
	testCases := []struct {
		desc    string
		header  http.Header
		wantOK  bool
		wantMin time.Duration
		wantMax time.Duration
	}{
		{desc: "nil header", header: nil, wantOK: false},
		{desc: "missing header", header: http.Header{}, wantOK: false},
		{
			desc:    "seconds",
			header:  http.Header{"Retry-After": []string{"7"}},
			wantOK:  true,
			wantMin: 7 * time.Second,
			wantMax: 7 * time.Second,
		},
		{
			desc:   "zero seconds",
			header: http.Header{"Retry-After": []string{"0"}},
			wantOK: true,
		},
		{
			desc:   "negative seconds rejected",
			header: http.Header{"Retry-After": []string{"-1"}},
			wantOK: false,
		},
		{
			desc:    "http date in the future",
			header:  http.Header{"Retry-After": []string{time.Now().Add(30 * time.Second).UTC().Format(http.TimeFormat)}},
			wantOK:  true,
			wantMin: 28 * time.Second,
			wantMax: 31 * time.Second,
		},
		{
			desc:   "http date in the past rejected",
			header: http.Header{"Retry-After": []string{time.Now().Add(-1 * time.Hour).UTC().Format(http.TimeFormat)}},
			wantOK: false,
		},
		{
			desc:   "garbage rejected",
			header: http.Header{"Retry-After": []string{"not a number"}},
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, ok := retryAfterFromHeader(tc.header)
			if ok != tc.wantOK {
				t.Fatalf("retryAfterFromHeader ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got < tc.wantMin || got > tc.wantMax {
				t.Errorf("retryAfterFromHeader = %v, want in [%v, %v]", got, tc.wantMin, tc.wantMax)
			}
		})
	}
}

func TestRetryDurationHint(t *testing.T) {
	retryInfoBody := []byte(`{"message":"x","details":[{"@type":"type.googleapis.com/google.rpc.RetryInfo","retry_delay":"3s"}]}`)

	testCases := []struct {
		desc string
		err  error
		want time.Duration
	}{
		{
			desc: "nil",
			err:  nil,
			want: 0,
		},
		{
			desc: "plain error",
			err:  errors.New("boom"),
			want: 0,
		},
		{
			desc: "APIError with Retry-After header",
			err:  apierr.FromHTTPError(http.StatusServiceUnavailable, http.Header{"Retry-After": []string{"7"}}, nil),
			want: 7 * time.Second,
		},
		{
			desc: "APIError with RetryInfo details",
			err:  apierr.FromHTTPError(http.StatusServiceUnavailable, nil, retryInfoBody),
			want: 3 * time.Second,
		},
		{
			desc: "Retry-After header wins over RetryInfo",
			err:  apierr.FromHTTPError(http.StatusTooManyRequests, http.Header{"Retry-After": []string{"10"}}, retryInfoBody),
			want: 10 * time.Second,
		},
		{
			desc: "APIError with no hint",
			err:  apierr.FromHTTPError(http.StatusInternalServerError, nil, nil),
			want: 0,
		},
		{
			desc: "wrapped APIError",
			err:  fmt.Errorf("op: %w", apierr.FromHTTPError(http.StatusInternalServerError, nil, retryInfoBody)),
			want: 3 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := RetryDurationHint(tc.err); got != tc.want {
				t.Errorf("retryDurationHint = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewRetrier_Codes(t *testing.T) {
	testCases := []struct {
		name   string
		codes  []codes.Code
		err    error
		wantOK bool
	}{
		{
			name:   "not an API error",
			codes:  []codes.Code{codes.Internal},
			err:    errors.New("an error"),
			wantOK: false,
		},
		{
			name:   "API error with code not in list",
			codes:  []codes.Code{codes.Internal},
			err:    apierr.FromHTTPError(http.StatusNotFound, nil, nil),
			wantOK: false,
		},
		{
			name:   "API error with code in list",
			codes:  []codes.Code{codes.Internal},
			err:    apierr.FromHTTPError(http.StatusInternalServerError, nil, nil),
			wantOK: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRetrier(ops.BackoffPolicy{}, RetrierConfig{ErrorCodes: tc.codes})
			_, got := r.IsRetriable(tc.err)
			if got != tc.wantOK {
				t.Errorf("IsRetriable() = %v, want %v", got, tc.wantOK)
			}
		})
	}
}

func TestNewRetrier_Statuses(t *testing.T) {
	retriable := []int{
		http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	testCases := []struct {
		desc   string
		err    error
		wantOK bool
	}{
		{
			desc:   "429 retriable",
			err:    apierr.FromHTTPError(http.StatusTooManyRequests, nil, nil),
			wantOK: true,
		},
		{
			desc:   "503 retriable",
			err:    apierr.FromHTTPError(http.StatusServiceUnavailable, nil, nil),
			wantOK: true,
		},
		{
			desc:   "500 not in allowlist",
			err:    apierr.FromHTTPError(http.StatusInternalServerError, nil, nil),
			wantOK: false,
		},
		{
			desc:   "400 not retriable",
			err:    apierr.FromHTTPError(http.StatusBadRequest, nil, nil),
			wantOK: false,
		},
		{
			desc:   "plain error not retriable",
			err:    errors.New("boom"),
			wantOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			r := NewRetrier(ops.BackoffPolicy{}, RetrierConfig{StatusCodes: retriable})
			_, got := r.IsRetriable(tc.err)
			if got != tc.wantOK {
				t.Errorf("IsRetriable(%v) = %v, want %v", tc.err, got, tc.wantOK)
			}
		})
	}
}

func TestNewRetrier_CodePrevails(t *testing.T) {
	// An APIError with body that parses to a specific code; status 503.
	// Code: RESOURCE_EXHAUSTED (in allowlist). Should retry.
	retriableByCode := apierr.FromHTTPError(
		http.StatusServiceUnavailable,
		nil,
		[]byte(`{"error_code":"RESOURCE_EXHAUSTED","message":"rate limit"}`),
	)

	// An APIError with code NOT_FOUND but status 503. Code prevails — don't
	// retry even though 503 is in the status allowlist.
	terminalByCode := apierr.FromHTTPError(
		http.StatusServiceUnavailable,
		nil,
		[]byte(`{"error_code":"NOT_FOUND","message":"nope"}`),
	)

	cfg := RetrierConfig{
		ErrorCodes:  []codes.Code{codes.ResourceExhausted, codes.Unavailable},
		StatusCodes: []int{http.StatusServiceUnavailable, http.StatusTooManyRequests},
	}
	r := NewRetrier(ops.BackoffPolicy{}, cfg)

	if _, ok := r.IsRetriable(retriableByCode); !ok {
		t.Errorf("retriable-by-code: IsRetriable = false, want true")
	}
	if _, ok := r.IsRetriable(terminalByCode); ok {
		t.Errorf("terminal-by-code: IsRetriable = true, want false (code prevails over status)")
	}
}

func TestNewRetrier_RetryAfterHonored(t *testing.T) {
	err := apierr.FromHTTPError(
		http.StatusTooManyRequests,
		http.Header{"Retry-After": []string{"5"}},
		nil,
	)

	r := NewRetrier(ops.BackoffPolicy{}, RetrierConfig{StatusCodes: []int{http.StatusTooManyRequests}})
	delay, ok := r.IsRetriable(err)

	if !ok {
		t.Fatalf("IsRetriable ok = false, want true")
	}
	if want := 5 * time.Second; delay < want {
		t.Errorf("delay = %v, want >= %v (Retry-After honored)", delay, want)
	}
}
