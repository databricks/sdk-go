package apierr

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/google/go-cmp/cmp"
)

func TestAPIError_nonHTTPGetters(t *testing.T) {
	testCases := []struct {
		apiErr      *APIError
		wantCode    codes.Code
		wantMessage string
		wantDetails ErrorDetails
	}{
		{
			apiErr:      &APIError{},
			wantCode:    codes.Unknown,
			wantMessage: "",
			wantDetails: ErrorDetails{},
		},
		{
			apiErr: &APIError{
				code:    codes.InvalidArgument,
				message: "Invalid request",
				details: ErrorDetails{
					ErrorInfo: &ErrorInfo{
						Reason: "bad_param",
						Domain: "databricks.com",
					},
				},
			},
			wantCode:    codes.InvalidArgument,
			wantMessage: "Invalid request",
			wantDetails: ErrorDetails{
				ErrorInfo: &ErrorInfo{
					Reason: "bad_param",
					Domain: "databricks.com",
				},
			},
		},
	}

	for _, tc := range testCases {
		if got := tc.apiErr.Code(); got != tc.wantCode {
			t.Errorf("Code: want %v, got %v", tc.wantCode, got)
		}
		if got := tc.apiErr.Message(); got != tc.wantMessage {
			t.Errorf("Message: want %q, got %q", tc.wantMessage, got)
		}
		if diff := cmp.Diff(tc.apiErr.Details(), tc.wantDetails); diff != "" {
			t.Errorf("details mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestAPIError_httpGetters(t *testing.T) {
	testCases := []struct {
		apiErr         *APIError
		wantStatusCode int
		wantHeader     http.Header
		wantBody       []byte
	}{
		{
			apiErr:         &APIError{},
			wantStatusCode: -1,
			wantHeader:     nil,
			wantBody:       nil,
		},
		{
			apiErr: &APIError{
				httpErr: &httpError{
					statusCode: http.StatusBadRequest,
					header:     http.Header{"Content-Type": {"application/json"}},
					body:       []byte(`{"error_code": "INVALID_ARGUMENT", "message": "Invalid request"}`),
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantHeader:     http.Header{"Content-Type": {"application/json"}},
			wantBody:       []byte(`{"error_code": "INVALID_ARGUMENT", "message": "Invalid request"}`),
		},
	}

	for _, tc := range testCases {
		if got := tc.apiErr.HTTPStatusCode(); got != tc.wantStatusCode {
			t.Errorf("HTTPStatusCode: want %d, got %d", tc.wantStatusCode, got)
		}
		if !cmp.Equal(tc.apiErr.HTTPHeader(), tc.wantHeader) {
			t.Errorf("HTTPHeader: want %v, got %v", tc.wantHeader, tc.apiErr.HTTPHeader())
		}
		if !cmp.Equal(tc.apiErr.HTTPBody(), tc.wantBody) {
			t.Errorf("HTTPBody: want %v, got %v", tc.wantBody, tc.apiErr.HTTPBody())
		}
	}
}

func TestFromHTTPError(t *testing.T) {
	testCases := []struct {
		desc        string
		statusCode  int
		header      http.Header
		body        []byte
		want        *APIError // httpErr set automatically from inputs
		wantWrapped bool      // whether APIError should have a wrapped error
	}{
		{
			desc:       "200 returns nil",
			statusCode: http.StatusOK,
		},
		{
			desc:       "201 returns nil",
			statusCode: http.StatusCreated,
		},
		{
			desc:       "204 returns nil",
			statusCode: http.StatusNoContent,
		},
		{
			desc:       "300 returns nil",
			statusCode: http.StatusMultipleChoices,
		},
		{
			desc:       "empty body with status",
			statusCode: http.StatusBadRequest,
			want: &APIError{
				code: codes.InvalidArgument,
			},
		},
		{
			desc:       "empty body with status and headers",
			statusCode: http.StatusNotFound,
			header:     http.Header{"Content-Type": {"application/json"}},
			want: &APIError{
				code: codes.NotFound,
			},
		},
		{
			desc:       "HTML body",
			statusCode: http.StatusBadGateway,
			body:       []byte(`<html><body>Bad Gateway</body></html>`),
			want: &APIError{
				code: codes.Internal,
			},
			wantWrapped: true,
		},
		{
			desc:       "malformed JSON",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{not valid json`),
			want: &APIError{
				code: codes.InvalidArgument,
			},
			wantWrapped: true,
		},
		{
			desc:       "standard error no details",
			statusCode: http.StatusNotFound,
			body:       []byte(`{"error_code": "NOT_FOUND", "message": "Job 123 not found"}`),
			want: &APIError{
				code:    codes.NotFound,
				message: "Job 123 not found",
			},
		},
		{
			desc:       "standard error with details",
			statusCode: http.StatusNotFound,
			body:       []byte(`{"error_code": "NOT_FOUND", "message": "Job 123 not found", "details": [{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": "bad_param", "domain": "databricks.com"}]}`),
			want: &APIError{
				code:    codes.NotFound,
				message: "Job 123 not found",
				details: ErrorDetails{
					ErrorInfo: &ErrorInfo{
						Reason: "bad_param",
						Domain: "databricks.com",
					},
				},
			},
		},
		{
			desc:       "standard error with unknown error_code",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"error_code": "SOME_UNKNOWN_CODE", "message": "Something went wrong"}`),
			want: &APIError{
				code:    codes.Unknown,
				message: "Something went wrong",
			},
		},
		{
			desc:       "standard error with missing error_code",
			statusCode: http.StatusForbidden,
			body:       []byte(`{"message": "Access denied"}`),
			want: &APIError{
				code:    codes.PermissionDenied,
				message: "Access denied",
			},
		},
		{
			desc:       "standard error with integer error_code",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"error_code": 42, "message": "Invalid request"}`),
			want: &APIError{
				code:    codes.InvalidArgument,
				message: "Invalid request",
			},
		},
		{
			desc:       "legacy API 1.2 error field",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"error": "Invalid parameter"}`),
			want: &APIError{
				code:    codes.InvalidArgument,
				message: "Invalid parameter",
			},
		},
		{
			desc:       "message takes precedence over error field",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"message": "New message", "error": "Old error"}`),
			want: &APIError{
				code:    codes.InvalidArgument,
				message: "New message",
			},
		},
		{
			desc:       "SCIM error with detail",
			statusCode: http.StatusNotFound,
			body:       []byte(`{"detail": "User not found", "scimType": "invalidValue"}`),
			want: &APIError{
				code:    codes.NotFound,
				message: "User not found",
			},
		},
		{
			desc:       "SCIM error with only scimType",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"scimType": "uniqueness"}`),
			want: &APIError{
				code:    codes.InvalidArgument,
				message: "uniqueness",
			},
		},
		{
			desc:       "message takes precedence over SCIM detail",
			statusCode: http.StatusBadRequest,
			body:       []byte(`{"message": "Standard message", "detail": "SCIM detail"}`),
			want: &APIError{
				code:    codes.InvalidArgument,
				message: "Standard message",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := FromHTTPError(tc.statusCode, tc.header, tc.body)

			// Early return for nil cases.
			if got == nil && tc.want == nil {
				return // success
			}
			if (got == nil) != (tc.want == nil) {
				t.Fatalf("want %v, got %v", tc.want, got)
			}

			// Input HTTP fields are meant to be preserved in httpErr.
			tc.want.httpErr = &httpError{
				statusCode: tc.statusCode,
				header:     tc.header,
				body:       tc.body,
			}

			if tc.want.HTTPStatusCode() != got.HTTPStatusCode() {
				t.Errorf("HTTPStatusCode: want %d, got %d", tc.want.HTTPStatusCode(), got.HTTPStatusCode())
			}
			if !cmp.Equal(tc.want.HTTPHeader(), got.HTTPHeader()) {
				t.Errorf("HTTPHeader: want %v, got %v", tc.want.HTTPHeader(), got.HTTPHeader())
			}
			if !cmp.Equal(tc.want.HTTPBody(), got.HTTPBody()) {
				t.Errorf("HTTPBody: want %v, got %v", tc.want.HTTPBody(), got.HTTPBody())
			}
			if tc.want.Code() != got.Code() {
				t.Errorf("Code: want %v, got %v", tc.want.Code(), got.Code())
			}
			if tc.want.Message() != got.Message() {
				t.Errorf("Message: want %q, got %q", tc.want.Message(), got.Message())
			}
			if diff := cmp.Diff(tc.want.Details(), got.Details()); diff != "" {
				t.Errorf("Details: mismatch (-want +got):\n%s", diff)
			}
			if tc.wantWrapped != (got.wrapped != nil) {
				t.Errorf("wrapped: want %v, got %v", tc.wantWrapped, got.wrapped != nil)
			}
		})
	}
}

func TestCode(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want codes.Code
	}{
		{
			name: "nil error returns OK",
			err:  nil,
			want: codes.OK,
		},
		{
			name: "non-APIError returns Unknown",
			err:  errors.New("some error"),
			want: codes.Unknown,
		},
		{
			name: "APIError returns its code",
			err:  &APIError{code: codes.NotFound},
			want: codes.NotFound,
		},
		{
			name: "wrapped APIError returns its code",
			err:  fmt.Errorf("wrapped: %w", &APIError{code: codes.PermissionDenied}),
			want: codes.PermissionDenied,
		},
		{
			name: "double wrapped APIError returns its code",
			err:  fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", &APIError{code: codes.InvalidArgument})),
			want: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Code(tc.err)
			if got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestToCode(t *testing.T) {
	testCases := []struct {
		httpCode int
		want     codes.Code
	}{
		// Direct mappings.
		{http.StatusOK, codes.OK},
		{http.StatusBadRequest, codes.InvalidArgument},
		{http.StatusUnauthorized, codes.Unauthenticated},
		{http.StatusForbidden, codes.PermissionDenied},
		{http.StatusNotFound, codes.NotFound},
		{http.StatusConflict, codes.Aborted},
		{http.StatusRequestedRangeNotSatisfiable, codes.OutOfRange},
		{http.StatusTooManyRequests, codes.ResourceExhausted},
		{http.StatusGatewayTimeout, codes.DeadlineExceeded},
		{http.StatusNotImplemented, codes.Unimplemented},
		{http.StatusServiceUnavailable, codes.Unavailable},

		// Fallback ranges.
		{http.StatusCreated, codes.OK},                   // 2xx
		{http.StatusNoContent, codes.OK},                 // 2xx
		{http.StatusTeapot, codes.FailedPrecondition},    // 4xx
		{http.StatusInternalServerError, codes.Internal}, // 5xx
		{599, codes.Internal},                            // 5xx

		// Unknown (valid).
		{100, codes.Unknown}, // 1xx
		{300, codes.Unknown}, // 3xx

		// Unknown (invalid).
		{-1, codes.Unknown},
		{0, codes.Unknown},
		{42, codes.Unknown},
		{600, codes.Unknown},
		{1337, codes.Unknown},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("status_%d", tc.httpCode), func(t *testing.T) {
			if got := toCode(tc.httpCode); got != tc.want {
				t.Errorf("toCode(%d): want %v, got %v", tc.httpCode, tc.want, got)
			}
		})
	}
}

// Ensure APIError implements error interface and Unwrap.
func TestAPIErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	apiErr := &APIError{
		code:    codes.Internal,
		message: "outer error",
		wrapped: innerErr,
	}

	if got := apiErr.Unwrap(); got != innerErr {
		t.Errorf("Unwrap: want %v, got %v", innerErr, got)
	}
	if !errors.Is(apiErr, innerErr) {
		t.Error("errors.Is should find the wrapped error")
	}
}
