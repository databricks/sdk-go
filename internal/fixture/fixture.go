// Package fixture provides a transport-agnostic testing framework for mocking
// API requests and responses with support for ordered sequences.
package fixture

import (
	"github.com/databricks/sdk-go/databricks/apierr/codes"
)

// Request represents a transport-agnostic API request.
// It captures the essential elements needed for matching without
// being tied to HTTP or any other transport.
type Request struct {
	// Method is the operation type (e.g., "GET", "POST" for HTTP,
	// or RPC method names for gRPC).
	Method string

	// Path is the resource path (e.g., "/api/2.1/jobs/create").
	Path string

	// Body is the request body as bytes. Can be nil for requests
	// without a body.
	Body []byte

	// Headers is a map of header key to values. For HTTP this maps
	// directly; for gRPC this could represent metadata.
	Headers map[string][]string
}

// Response represents a transport-agnostic API response.
type Response struct {
	// StatusCode is the response status. For HTTP, this is the HTTP
	// status code. For gRPC, this could be mapped from gRPC status codes.
	StatusCode int

	// Body is the response body. For successful responses, this is
	// typically JSON-encoded. For errors, this may contain error details.
	Body []byte

	// Headers is a map of response headers.
	Headers map[string][]string

	// Error, if set, indicates the response is an error response.
	// This is transport-agnostic and uses the SDK's canonical error codes.
	Error *ErrorResponse
}

// ErrorResponse represents an error that should be returned.
type ErrorResponse struct {
	// Code is the canonical error code.
	Code codes.Code

	// Message is the error message.
	Message string
}

// Fixture defines a single expected request and its response.
// Fixtures are immutable after creation.
type Fixture struct {
	// Request matching criteria.
	Method string // Required: the method to match (e.g., "GET", "POST").
	Path   string // Required: the path to match exactly.

	// Optional request validation.
	ExpectedBody []byte // If set, request body must match exactly as bytes.

	// Response definition.
	Response   []byte         // Raw response body bytes.
	StatusCode int            // HTTP status code (default: 200 for Response, derived for Error).
	Error      *ErrorResponse // Error response to return.
}

// Sequence is an ordered collection of fixtures that must be consumed
// in order. It is the primary type used in tests.
//
// Sequence is NOT thread-safe; each test should have its own Sequence.
type Sequence struct {
	fixtures []Fixture
	index    int       // Current position (next fixture to match).
	consumed []bool    // Tracks which fixtures were consumed.
	calls    []Request // Records all requests received.
}

// NewSequence creates a new Sequence from the given fixtures.
// The fixtures will be matched in order.
func NewSequence(fixtures ...Fixture) *Sequence {
	return &Sequence{
		fixtures: fixtures,
		consumed: make([]bool, len(fixtures)),
		calls:    make([]Request, 0),
	}
}

// Match attempts to match a request against the current expected fixture.
// It returns the response for the matching fixture or an error if:
// - The request does not match the expected fixture.
// - All fixtures have been consumed.
//
// This method advances the sequence position on successful match.
func (s *Sequence) Match(req Request) (Response, error) {
	s.calls = append(s.calls, req)

	if s.index >= len(s.fixtures) {
		return Response{}, &UnexpectedRequestError{
			Request:  req,
			Message:  "all fixtures have been consumed",
			Expected: s.fixtures,
			Received: s.calls,
		}
	}

	fixture := s.fixtures[s.index]
	if err := matchRequest(req, fixture); err != nil {
		return Response{}, &MismatchError{
			Request:  req,
			Expected: fixture,
			Index:    s.index,
			Cause:    err,
		}
	}

	s.consumed[s.index] = true
	s.index++

	return buildResponse(req, fixture)
}

// Verify checks that all fixtures were consumed.
// Call this at the end of a test to ensure all expected calls were made.
func (s *Sequence) Verify() error {
	if s.index < len(s.fixtures) {
		return &UnconsumedFixturesError{
			Fixtures:  s.fixtures,
			Consumed:  s.consumed,
			NextIndex: s.index,
		}
	}
	return nil
}

// Reset resets the sequence to its initial state.
// Useful for table-driven tests that reuse fixture definitions.
func (s *Sequence) Reset() {
	s.index = 0
	s.consumed = make([]bool, len(s.fixtures))
	s.calls = make([]Request, 0)
}

// Calls returns all requests received by this sequence.
// Useful for debugging test failures.
func (s *Sequence) Calls() []Request {
	return s.calls
}
