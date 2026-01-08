// Package httpfixture provides HTTP transport adapters for the fixture package.
package httpfixture

import (
	"bytes"
	"io"
	"net/http"

	"github.com/databricks/sdk-go/internal/fixture"
)

// Transport implements http.RoundTripper using a fixture Sequence.
// This allows injecting fixtures into any HTTP client.
type Transport struct {
	seq *fixture.Sequence
}

// NewTransport creates a new Transport backed by the given Sequence.
func NewTransport(seq *fixture.Sequence) *Transport {
	return &Transport{seq: seq}
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Convert http.Request to fixture.Request.
	var body []byte
	if req.Body != nil {
		var err error
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()
		// Restore body for potential re-reads.
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	fixReq := fixture.Request{
		Method:  req.Method,
		Path:    req.URL.Path,
		Body:    body,
		Headers: req.Header,
	}

	// Match against sequence.
	resp, err := t.seq.Match(fixReq)
	if err != nil {
		// Return the mismatch as a clear test failure.
		// In tests, this will cause the test to fail with a descriptive message.
		return nil, err
	}

	// Convert fixture.Response to http.Response.
	httpResp := &http.Response{
		StatusCode: resp.StatusCode,
		Status:     http.StatusText(resp.StatusCode),
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(resp.Body)),
		Request:    req,
	}

	for k, v := range resp.Headers {
		httpResp.Header[k] = v
	}

	return httpResp, nil
}

// Client returns an http.Client using this transport.
// This is a convenience method for use with WithHTTPClient option.
func (t *Transport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// Verify delegates to the underlying Sequence's Verify method.
// Call this at the end of a test.
func (t *Transport) Verify() error {
	return t.seq.Verify()
}

// Calls returns all requests received by the underlying sequence.
// Useful for debugging test failures.
func (t *Transport) Calls() []fixture.Request {
	return t.seq.Calls()
}

// Reset resets the underlying sequence to its initial state.
// Useful for table-driven tests that reuse fixture definitions.
func (t *Transport) Reset() {
	t.seq.Reset()
}
