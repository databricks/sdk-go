package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/log/logtesting"
)

// Sentinel test errors for use with errors.Is.
var (
	errCredentials = errors.New("credentials error")
	errTransport   = errors.New("transport error")
	errClose       = errors.New("close error")
)

// mockCredentials implements auth.Credentials for testing.
type mockCredentials struct {
	headers     []auth.Header
	err         error
	capturedCtx context.Context
}

func (m *mockCredentials) AuthHeaders(ctx context.Context) ([]auth.Header, error) {
	m.capturedCtx = ctx
	return m.headers, m.err
}

// mockTransport implements http.RoundTripper for testing.
type mockTransport struct {
	response *http.Response
	err      error

	// capturedReq stores the request passed to RoundTrip for inspection.
	capturedReq *http.Request
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.capturedReq = req
	return m.response, m.err
}

// mockBody implements io.ReadCloser for testing body close behavior.
type mockBody struct {
	closed   bool
	closeErr error
}

func (m *mockBody) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (m *mockBody) Close() error {
	m.closed = true
	return m.closeErr
}

func TestAuthTransport_RoundTrip(t *testing.T) {
	testCases := []struct {
		desc          string
		credHeaders   []auth.Header
		credErr       error
		transportResp *http.Response
		transportErr  error
		wantErr       error
	}{
		{
			desc: "adds single auth header",
			credHeaders: []auth.Header{
				{Key: "Authorization", Value: "Bearer token123"},
			},
			transportResp: &http.Response{StatusCode: 200},
		},
		{
			desc: "adds multiple auth headers",
			credHeaders: []auth.Header{
				{Key: "Authorization", Value: "Bearer token123"},
				{Key: "X-Custom-Auth", Value: "custom-value"},
			},
			transportResp: &http.Response{StatusCode: 200},
		},
		{
			desc:    "returns credentials error",
			credErr: errCredentials,
			wantErr: errCredentials,
		},
		{
			desc: "propagates transport error",
			credHeaders: []auth.Header{
				{Key: "Authorization", Value: "Bearer token123"},
			},
			transportErr: errTransport,
			wantErr:      errTransport,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			creds := &mockCredentials{
				headers: tc.credHeaders,
				err:     tc.credErr,
			}
			base := &mockTransport{
				response: tc.transportResp,
				err:      tc.transportErr,
			}
			transport := &authTransport{
				base:  base,
				creds: creds,
			}

			req, err := http.NewRequest("GET", "https://example.com/api", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			gotResp, gotErr := transport.RoundTrip(req)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Fatalf("got error %v, want %v", gotErr, tc.wantErr)
			}
			if gotResp != tc.transportResp {
				t.Errorf("response: got %v, want %v", gotResp, tc.transportResp)
			}

			// Check that the auth headers were added to the request sent to
			// the base transport.
			for _, h := range tc.credHeaders {
				if got := base.capturedReq.Header.Get(h.Key); got != h.Value {
					t.Errorf("%s header: got %q, want %q", h.Key, got, h.Value)
				}
			}
		})
	}
}

func TestAuthTransport_RoundTrip_DoesNotModifyOriginalRequest(t *testing.T) {
	creds := &mockCredentials{
		headers: []auth.Header{{Key: "Authorization", Value: "Bearer token123"}},
	}
	base := &mockTransport{
		response: &http.Response{StatusCode: 200},
	}
	transport := &authTransport{
		base:  base,
		creds: creds,
	}

	req, err := http.NewRequest("GET", "https://example.com/api", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Keep track of the original header state to verify that is was not
	// modified by the transport.
	originalAuthHeader := req.Header.Get("Authorization")

	_, err = transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	if req == base.capturedReq {
		t.Error("base transport received original request instead of clone")
	}
	if got := req.Header.Get("Authorization"); got != originalAuthHeader {
		t.Errorf("Authorization header from original request was modified: got %q, want %q", got, originalAuthHeader)
	}
	if got := base.capturedReq.Header.Get("Authorization"); got != "Bearer token123" {
		t.Errorf("cloned request missing auth header: got %q", got)
	}
}

func TestAuthTransport_RoundTrip_ClosesBodyOnCredentialsError(t *testing.T) {
	creds := &mockCredentials{
		err: errCredentials,
	}
	transport := &authTransport{
		base:  &mockTransport{},
		creds: creds,
	}

	body := &mockBody{
		closeErr: errClose,
	}
	req, err := http.NewRequest("POST", "https://example.com/api", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = transport.RoundTrip(req)

	if !body.closed {
		t.Error("request body was not closed")
	}
	// Should return credentials error, not close error
	if !errors.Is(err, errCredentials) {
		t.Errorf("got error %v, want %v", err, errCredentials)
	}
	if errors.Is(err, errClose) {
		t.Errorf("should not return close error, got: %v", err)
	}
}

func TestAuthTransport_RoundTrip_PassesContextToCredentials(t *testing.T) {
	type ctxKey string
	key := ctxKey("test-key")
	wantValue := "test-value"

	creds := &mockCredentials{
		headers: []auth.Header{{Key: "Auth", Value: "token"}},
	}
	base := &mockTransport{
		response: &http.Response{StatusCode: 200},
	}
	transport := &authTransport{
		base:  base,
		creds: creds,
	}

	ctx := context.WithValue(context.Background(), key, wantValue)
	req, err := http.NewRequestWithContext(ctx, "GET", "https://example.com/api", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = transport.RoundTrip(req)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	if creds.capturedCtx == nil {
		t.Fatal("context was not passed to credentials")
	}
	if got := creds.capturedCtx.Value(key); got != wantValue {
		t.Errorf("context value = %v, want %v", got, wantValue)
	}
}

func TestAuthTransport_RoundTrip_LogsCloseError(t *testing.T) {
	logger := &logtesting.Logger{}
	transport := &authTransport{
		base:   &mockTransport{},
		creds:  &mockCredentials{err: errCredentials},
		logger: logger,
	}

	body := &mockBody{closeErr: errClose}
	req, err := http.NewRequest("POST", "https://example.com/api", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = transport.RoundTrip(req)
	if !errors.Is(err, errCredentials) {
		t.Fatalf("got error %v, want %v", err, errCredentials)
	}

	records := logger.Records()
	if len(records) != 1 {
		t.Fatalf("got %d log records, want 1", len(records))
	}
	if len(records[0].Args) != 1 || !errors.Is(records[0].Args[0].(error), errClose) {
		t.Errorf("log record args = %v, want [%v]", records[0].Args, errClose)
	}
}
