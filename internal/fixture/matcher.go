package fixture

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// matchRequest checks if a request matches a fixture's criteria.
func matchRequest(req Request, f Fixture) error {
	// Match method exactly.
	if req.Method != f.Method {
		return fmt.Errorf("method mismatch: got %q, want %q", req.Method, f.Method)
	}

	// Match path exactly.
	if req.Path != f.Path {
		return fmt.Errorf("path mismatch: got %q, want %q", req.Path, f.Path)
	}

	// Match body if ExpectedBody is specified.
	if f.ExpectedBody != nil {
		if err := matchBody(req.Body, f.ExpectedBody); err != nil {
			return fmt.Errorf("body mismatch: %w", err)
		}
	}

	return nil
}

// matchBody compares the request body against the expected body.
// Both are compared as JSON-decoded values for semantic equality.
func matchBody(actual []byte, expected any) error {
	if len(actual) == 0 && expected == nil {
		return nil
	}

	var actualValue any
	if len(actual) > 0 {
		if err := json.Unmarshal(actual, &actualValue); err != nil {
			return fmt.Errorf("failed to unmarshal actual body: %w", err)
		}
	}

	// Marshal expected to JSON then unmarshal for consistent comparison.
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("failed to marshal expected body: %w", err)
	}
	var expectedValue any
	if err := json.Unmarshal(expectedBytes, &expectedValue); err != nil {
		return fmt.Errorf("failed to unmarshal expected body: %w", err)
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		return fmt.Errorf("body values differ:\n  got:  %v\n  want: %v", actualValue, expectedValue)
	}

	return nil
}

// buildResponse constructs a Response from a Fixture.
func buildResponse(req Request, f Fixture) Response {
	// If Error is set, build an error response.
	if f.Error != nil {
		statusCode := f.StatusCode
		if statusCode == 0 {
			statusCode = codeToHTTPStatus(f.Error.Code)
		}
		body, _ := json.Marshal(map[string]any{
			"error_code": f.Error.Code.String(),
			"message":    f.Error.Message,
		})
		return Response{
			StatusCode: statusCode,
			Body:       body,
			Error:      f.Error,
		}
	}

	// Build success response.
	statusCode := f.StatusCode
	if statusCode == 0 {
		statusCode = 200
	}

	var body []byte
	if f.Response != nil {
		body, _ = json.Marshal(f.Response)
	}

	return Response{
		StatusCode: statusCode,
		Body:       body,
	}
}

// codeToHTTPStatus maps canonical error codes to HTTP status codes.
func codeToHTTPStatus(c any) int {
	// Use type assertion to get the string representation.
	type stringer interface {
		String() string
	}
	s, ok := c.(stringer)
	if !ok {
		return 500
	}

	switch s.String() {
	case "CANCELLED":
		return 499
	case "INVALID_ARGUMENT":
		return 400
	case "DEADLINE_EXCEEDED":
		return 504
	case "NOT_FOUND":
		return 404
	case "ALREADY_EXISTS":
		return 409
	case "PERMISSION_DENIED":
		return 403
	case "RESOURCE_EXHAUSTED":
		return 429
	case "FAILED_PRECONDITION":
		return 400
	case "ABORTED":
		return 409
	case "OUT_OF_RANGE":
		return 400
	case "UNIMPLEMENTED":
		return 501
	case "INTERNAL":
		return 500
	case "UNAVAILABLE":
		return 503
	case "DATA_LOSS":
		return 500
	case "UNAUTHENTICATED":
		return 401
	default:
		return 500
	}
}
