package apierr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/databricks/sdk-go/databricks/apierr/codes"
)

// APIError is a transport-agnostic error representing a Databricks API error.
type APIError struct {
	// The canonical error code of the error.
	code codes.Code

	// The message received in the error response. This is left empty if the
	// error response is not a standard Databricks API error.
	message string

	// The structured error details of the error. This is left empty if the
	// error response is not a standard Databricks API error.
	details ErrorDetails

	// The "raw" HTTP error details.
	httpErr *httpError

	wrapped error // convenience to implement Unwrap()
}

// httpError holds the raw HTTP response details for debugging.
type httpError struct {
	statusCode int
	body       []byte
	header     http.Header
}

func (e *APIError) Error() string {
	// TODO: Implement error printing logic.
	return fmt.Sprintf("databricks-api error: %s", e.Message())
}

func (e *APIError) Unwrap() error {
	return e.wrapped
}

// Code returns the APIError's error code.
func (e *APIError) Code() codes.Code {
	return e.code
}

// Message returns the APIError's error message.
func (e *APIError) Message() string {
	return e.message
}

// Details returns the APIError's error details.
func (e *APIError) Details() ErrorDetails {
	return e.details
}

// HTTPStatusCode returns the APIError's HTTP status code. If the APIError
// is not an HTTP error, it returns -1.
func (e *APIError) HTTPStatusCode() int {
	if e.httpErr == nil {
		return -1
	}
	return e.httpErr.statusCode
}

// HTTPHeader returns the APIError's HTTP headers. If the APIError is not
// an HTTP error, it returns nil.
func (e *APIError) HTTPHeader() http.Header {
	if e.httpErr == nil {
		return nil
	}
	return e.httpErr.header
}

// HTTPBody returns the APIError's HTTP body. If the APIError is not an HTTP
// error, it returns nil.
func (e *APIError) HTTPBody() []byte {
	if e.httpErr == nil {
		return nil
	}
	return e.httpErr.body
}

// Code returns the error code from err if it's an APIError (or wraps one).
// It returns codes.Unknown for non-APIError errors, codes.OK for nil.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	var e *APIError
	if errors.As(err, &e) {
		return e.Code()
	}
	return codes.Unknown
}

// FromHTTPError parses an HTTP error response into an APIError; it returns
// nil if the status code is 2xx.
func FromHTTPError(statusCode int, header http.Header, body []byte) *APIError {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	apiErr := &APIError{
		httpErr: &httpError{
			statusCode: statusCode,
			header:     header,
			body:       body,
		},
	}

	if len(body) == 0 {
		apiErr.code = toCode(statusCode)
		return apiErr
	}

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// The JSON error is simply swallowed, this typically happens when the
		// error does not come directly from a Databricks API. A typical example
		// is when the error is returned by a proxy.
		apiErr.code = toCode(statusCode)
		apiErr.wrapped = err // retain the JSON parsing error
		return apiErr
	}

	// Error codes may be missing or be an integer (legacy APIs). In such cases,
	// defer to the HTTP status code to infer the closest canonical error code.
	switch v := errResp.ErrorCode.(type) {
	case string:
		apiErr.code = codes.FromString(v)
	default:
		apiErr.code = toCode(statusCode)
	}

	switch {
	case errResp.Message != "":
		apiErr.message = errResp.Message // standard
	case errResp.API12Error != "":
		apiErr.message = errResp.API12Error // legacy
	case errResp.ScimDetail != "":
		apiErr.message = errResp.ScimDetail // scim
	case errResp.ScimType != "":
		apiErr.message = errResp.ScimType // scim fallback
	}

	apiErr.details = parseErrorDetails(errResp.RawDetails)

	return apiErr
}

type errorResponse struct {
	Message    string            `json:"message,omitempty"`
	RawDetails []json.RawMessage `json:"details,omitempty"`

	// Some Databricks APIs incorrectly return the HTTP status code as an
	// integer rather than the actual error code as a string.
	ErrorCode any `json:"error_code,omitempty"` // int or string

	// Legacy Databricks APIs (e.g. version 1.2 and earlier) used "error"
	// instead of "message".
	API12Error string `json:"error,omitempty"`

	// SCIM error fields (RFC7644 section 3.7.3).
	// The "status" field is intentionally omitted; it duplicates HTTP status.
	ScimDetail string `json:"detail,omitempty"`
	ScimType   string `json:"scimType,omitempty"`
}

// toCode maps an HTTP status code to the closest canonical error code.
func toCode(httpCode int) codes.Code {
	// Canonical mapping.
	switch httpCode {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusRequestedRangeNotSatisfiable:
		return codes.OutOfRange
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	}

	// Fallback for status codes without a direct canonical mapping.
	switch {
	case httpCode >= 200 && httpCode < 300:
		return codes.OK
	case httpCode >= 400 && httpCode < 500:
		// Most non-canonical 4xx status codes are state related and map
		// to the definition of FailedPrecondition.
		return codes.FailedPrecondition
	case httpCode >= 500 && httpCode < 600:
		return codes.Internal
	}

	return codes.Unknown
}
