package fixture

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/databricks/sdk-go/databricks/apierr/codes"
)

func TestSequence_Match_Success(t *testing.T) {
	seq := NewSequence(
		Fixture{
			Method:   "GET",
			Path:     "/api/test",
			Response: []byte(`{"id":1}`),
		},
	)

	resp, err := seq.Match(Request{
		Method: "GET",
		Path:   "/api/test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("StatusCode: got %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if body["id"] != float64(1) {
		t.Errorf("body id: got %v, want 1", body["id"])
	}
}

func TestSequence_Match_OrderedSequence(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "POST", Path: "/api/create", Response: []byte(`{"id":1}`)},
		Fixture{Method: "GET", Path: "/api/get", Response: []byte(`{"id":1}`)},
		Fixture{Method: "DELETE", Path: "/api/delete", StatusCode: 204},
	)

	// First request.
	resp, err := seq.Match(Request{Method: "POST", Path: "/api/create"})
	if err != nil {
		t.Fatalf("request 1: unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("request 1: StatusCode: got %d, want 200", resp.StatusCode)
	}

	// Second request.
	resp, err = seq.Match(Request{Method: "GET", Path: "/api/get"})
	if err != nil {
		t.Fatalf("request 2: unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("request 2: StatusCode: got %d, want 200", resp.StatusCode)
	}

	// Third request.
	resp, err = seq.Match(Request{Method: "DELETE", Path: "/api/delete"})
	if err != nil {
		t.Fatalf("request 3: unexpected error: %v", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("request 3: StatusCode: got %d, want 204", resp.StatusCode)
	}

	// Verify all fixtures consumed.
	if err := seq.Verify(); err != nil {
		t.Errorf("Verify: unexpected error: %v", err)
	}
}

func TestSequence_Match_MethodMismatch(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/test"},
	)

	_, err := seq.Match(Request{Method: "POST", Path: "/api/test"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var mismatch *MismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("expected MismatchError, got %T", err)
	}
	if mismatch.Index != 0 {
		t.Errorf("Index: got %d, want 0", mismatch.Index)
	}
}

func TestSequence_Match_PathMismatch(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/test"},
	)

	_, err := seq.Match(Request{Method: "GET", Path: "/api/other"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var mismatch *MismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("expected MismatchError, got %T", err)
	}
}

func TestSequence_Match_OutOfOrderRequest(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "POST", Path: "/api/create"},
		Fixture{Method: "GET", Path: "/api/get"},
	)

	// Try to make the second request first.
	_, err := seq.Match(Request{Method: "GET", Path: "/api/get"})
	if err == nil {
		t.Fatal("expected error for out-of-order request, got nil")
	}

	var mismatch *MismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("expected MismatchError, got %T", err)
	}
}

func TestSequence_Match_ExtraRequest(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/test"},
	)

	// First request succeeds.
	_, err := seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err != nil {
		t.Fatalf("request 1: unexpected error: %v", err)
	}

	// Extra request should fail.
	_, err = seq.Match(Request{Method: "GET", Path: "/api/extra"})
	if err == nil {
		t.Fatal("expected error for extra request, got nil")
	}

	var unexpected *UnexpectedRequestError
	if !errors.As(err, &unexpected) {
		t.Fatalf("expected UnexpectedRequestError, got %T", err)
	}
}

func TestSequence_Match_ExpectedBody(t *testing.T) {
	seq := NewSequence(
		Fixture{
			Method:       "POST",
			Path:         "/api/create",
			ExpectedBody: []byte(`{"name":"test"}`),
			Response:     []byte(`{"id":1}`),
		},
	)

	resp, err := seq.Match(Request{
		Method: "POST",
		Path:   "/api/create",
		Body:   []byte(`{"name":"test"}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode: got %d, want 200", resp.StatusCode)
	}
}

func TestSequence_Match_ExpectedBodyMismatch(t *testing.T) {
	seq := NewSequence(
		Fixture{
			Method:       "POST",
			Path:         "/api/create",
			ExpectedBody: []byte(`{"name":"expected"}`),
		},
	)

	_, err := seq.Match(Request{
		Method: "POST",
		Path:   "/api/create",
		Body:   []byte(`{"name":"actual"}`),
	})
	if err == nil {
		t.Fatal("expected error for body mismatch, got nil")
	}

	var mismatch *MismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("expected MismatchError, got %T", err)
	}
}

func TestSequence_Match_ErrorResponse(t *testing.T) {
	seq := NewSequence(
		Fixture{
			Method: "GET",
			Path:   "/api/test",
			Error: &ErrorResponse{
				Code:    codes.NotFound,
				Message: "Resource not found",
			},
		},
	)

	resp, err := seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("StatusCode: got %d, want 404", resp.StatusCode)
	}

	var body map[string]any
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal error body: %v", err)
	}
	if body["error_code"] != "NOT_FOUND" {
		t.Errorf("error_code: got %v, want NOT_FOUND", body["error_code"])
	}
}

func TestSequence_Match_ErrorResponseWithCustomStatus(t *testing.T) {
	seq := NewSequence(
		Fixture{
			Method:     "GET",
			Path:       "/api/test",
			StatusCode: 418,
			Error: &ErrorResponse{
				Code:    codes.Internal,
				Message: "I'm a teapot",
			},
		},
	)

	resp, err := seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != 418 {
		t.Errorf("StatusCode: got %d, want 418", resp.StatusCode)
	}
}

func TestSequence_Verify_Unconsumed(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/first"},
		Fixture{Method: "GET", Path: "/api/second"},
		Fixture{Method: "GET", Path: "/api/third"},
	)

	// Only consume the first fixture.
	_, err := seq.Match(Request{Method: "GET", Path: "/api/first"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = seq.Verify()
	if err == nil {
		t.Fatal("expected error for unconsumed fixtures, got nil")
	}

	var unconsumed *UnconsumedFixturesError
	if !errors.As(err, &unconsumed) {
		t.Fatalf("expected UnconsumedFixturesError, got %T", err)
	}
	if unconsumed.NextIndex != 1 {
		t.Errorf("NextIndex: got %d, want 1", unconsumed.NextIndex)
	}
}

func TestSequence_Reset(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/test", Response: []byte(`{"id":1}`)},
	)

	// First use.
	_, err := seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err != nil {
		t.Fatalf("first use: unexpected error: %v", err)
	}
	if err := seq.Verify(); err != nil {
		t.Fatalf("first use verify: unexpected error: %v", err)
	}

	// Reset.
	seq.Reset()

	// Second use.
	_, err = seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err != nil {
		t.Fatalf("second use: unexpected error: %v", err)
	}
	if err := seq.Verify(); err != nil {
		t.Fatalf("second use verify: unexpected error: %v", err)
	}
}

func TestSequence_Calls(t *testing.T) {
	seq := NewSequence(
		Fixture{Method: "GET", Path: "/api/test"},
	)

	req := Request{Method: "GET", Path: "/api/test", Body: []byte("test")}
	_, err := seq.Match(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := seq.Calls()
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

func TestSequence_EmptySequence(t *testing.T) {
	seq := NewSequence()

	// Any request should fail.
	_, err := seq.Match(Request{Method: "GET", Path: "/api/test"})
	if err == nil {
		t.Fatal("expected error for empty sequence, got nil")
	}

	var unexpected *UnexpectedRequestError
	if !errors.As(err, &unexpected) {
		t.Fatalf("expected UnexpectedRequestError, got %T", err)
	}

	// Verify should pass on empty sequence.
	if err := seq.Verify(); err != nil {
		t.Errorf("Verify: unexpected error on empty sequence: %v", err)
	}
}
