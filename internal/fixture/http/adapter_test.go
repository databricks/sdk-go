package httpfixture

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/databricks/sdk-go/databricks/apierr/codes"
	"github.com/databricks/sdk-go/internal/fixture"
)

func TestTransport_RoundTrip_Success(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method:   "GET",
			Path:     "/api/test",
			Response: map[string]any{"id": 1, "name": "test"},
		},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	req, _ := http.NewRequest("GET", "http://example.com/api/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode: got %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["id"] != float64(1) {
		t.Errorf("body id: got %v, want 1", body["id"])
	}
	if body["name"] != "test" {
		t.Errorf("body name: got %v, want test", body["name"])
	}
}

func TestTransport_RoundTrip_WithRequestBody(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method:       "POST",
			Path:         "/api/create",
			ExpectedBody: map[string]any{"name": "new-item"},
			Response:     map[string]any{"id": 42},
		},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	reqBody, _ := json.Marshal(map[string]any{"name": "new-item"})
	req, _ := http.NewRequest("POST", "http://example.com/api/create", bytes.NewReader(reqBody))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode: got %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["id"] != float64(42) {
		t.Errorf("body id: got %v, want 42", body["id"])
	}
}

func TestTransport_RoundTrip_ErrorResponse(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method: "GET",
			Path:   "/api/notfound",
			Error: &fixture.ErrorResponse{
				Code:    codes.NotFound,
				Message: "Resource not found",
			},
		},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	req, _ := http.NewRequest("GET", "http://example.com/api/notfound", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Errorf("StatusCode: got %d, want 404", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error body: %v", err)
	}
	if body["error_code"] != "NOT_FOUND" {
		t.Errorf("error_code: got %v, want NOT_FOUND", body["error_code"])
	}
}

func TestTransport_RoundTrip_Mismatch(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method: "GET",
			Path:   "/api/expected",
		},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	req, _ := http.NewRequest("GET", "http://example.com/api/actual", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Fatal("expected error for mismatched request, got nil")
	}

	var mismatch *fixture.MismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("expected MismatchError, got %T", err)
	}
}

func TestTransport_RoundTrip_OrderedSequence(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{Method: "POST", Path: "/api/create", Response: map[string]any{"id": 1}},
		fixture.Fixture{Method: "GET", Path: "/api/get/1", Response: map[string]any{"id": 1, "status": "created"}},
		fixture.Fixture{Method: "DELETE", Path: "/api/delete/1", StatusCode: 204},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	// Request 1: POST /api/create.
	req1, _ := http.NewRequest("POST", "http://example.com/api/create", nil)
	resp1, err := client.Do(req1)
	if err != nil {
		t.Fatalf("request 1: unexpected error: %v", err)
	}
	resp1.Body.Close()
	if resp1.StatusCode != 200 {
		t.Errorf("request 1: StatusCode: got %d, want 200", resp1.StatusCode)
	}

	// Request 2: GET /api/get/1.
	req2, _ := http.NewRequest("GET", "http://example.com/api/get/1", nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("request 2: unexpected error: %v", err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != 200 {
		t.Errorf("request 2: StatusCode: got %d, want 200", resp2.StatusCode)
	}

	// Request 3: DELETE /api/delete/1.
	req3, _ := http.NewRequest("DELETE", "http://example.com/api/delete/1", nil)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Fatalf("request 3: unexpected error: %v", err)
	}
	resp3.Body.Close()
	if resp3.StatusCode != 204 {
		t.Errorf("request 3: StatusCode: got %d, want 204", resp3.StatusCode)
	}

	// Verify all fixtures consumed.
	if err := transport.Verify(); err != nil {
		t.Errorf("Verify: unexpected error: %v", err)
	}
}

func TestTransport_Verify_Unconsumed(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{Method: "GET", Path: "/api/first"},
		fixture.Fixture{Method: "GET", Path: "/api/second"},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	// Only make the first request.
	req, _ := http.NewRequest("GET", "http://example.com/api/first", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	// Verify should fail.
	err = transport.Verify()
	if err == nil {
		t.Fatal("expected error for unconsumed fixtures, got nil")
	}

	var unconsumed *fixture.UnconsumedFixturesError
	if !errors.As(err, &unconsumed) {
		t.Fatalf("expected UnconsumedFixturesError, got %T", err)
	}
}

func TestTransport_Reset(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{Method: "GET", Path: "/api/test", Response: map[string]any{"id": 1}},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	// First use.
	req, _ := http.NewRequest("GET", "http://example.com/api/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("first use: unexpected error: %v", err)
	}
	resp.Body.Close()
	if err := transport.Verify(); err != nil {
		t.Fatalf("first use verify: unexpected error: %v", err)
	}

	// Reset.
	transport.Reset()

	// Second use.
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("second use: unexpected error: %v", err)
	}
	resp.Body.Close()
	if err := transport.Verify(); err != nil {
		t.Fatalf("second use verify: unexpected error: %v", err)
	}
}

func TestTransport_Calls(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{Method: "GET", Path: "/api/test"},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	req, _ := http.NewRequest("GET", "http://example.com/api/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()

	calls := transport.Calls()
	if len(calls) != 1 {
		t.Fatalf("Calls: got %d, want 1", len(calls))
	}
	if calls[0].Method != "GET" {
		t.Errorf("Calls[0].Method: got %s, want GET", calls[0].Method)
	}
	if calls[0].Path != "/api/test" {
		t.Errorf("Calls[0].Path: got %s, want /api/test", calls[0].Path)
	}
}

func TestTransport_ResponseHeaders(t *testing.T) {
	seq := fixture.NewSequence(
		fixture.Fixture{
			Method:   "GET",
			Path:     "/api/test",
			Response: map[string]any{"id": 1},
		},
	)

	transport := NewTransport(seq)
	client := transport.Client()

	req, _ := http.NewRequest("GET", "http://example.com/api/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	// Verify the response body can be read completely.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if len(body) == 0 {
		t.Error("response body is empty")
	}
}
