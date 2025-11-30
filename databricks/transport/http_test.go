package transport

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/databricks/sdk-go/auth"
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
	if req.Body != nil {
		defer req.Body.Close()
	}
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
		credHeaders   []auth.Header  // headers returned by the credentials
		credErr       error          // error returned by the credentials
		transportResp *http.Response // response returned by the base transport
		transportErr  error          // error returned by the base transport
		bodycloseErr  error          // error returned by the body close
		wantErr       error          // error returned by the transport
		wantLog       string         // message that should be logged
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
			desc: "propagates transport error",
			credHeaders: []auth.Header{
				{Key: "Authorization", Value: "Bearer token123"},
			},
			transportErr: errTransport,
			wantErr:      errTransport,
		},
		{
			desc:    "credentials error with no close error",
			credErr: errCredentials,
			wantErr: errCredentials,
		},
		{
			desc:         "credentials error with close error",
			bodycloseErr: errClose,
			credErr:      errCredentials,
			wantErr:      errCredentials,
			wantLog:      "level=ERROR msg=\"error closing request body\" error=\"close error\"\n",
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
			buf := bytes.Buffer{} // collect logs
			transport := &authTransport{
				base:   base,
				creds:  creds,
				logger: slog.New(slog.NewTextHandler(&buf, nil)),
			}

			body := &mockBody{closeErr: tc.bodycloseErr}
			req, err := http.NewRequest("GET", "https://example.com/api", body)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			gotResp, gotErr := transport.RoundTrip(req)

			// The body should always be closed by the transport, even in case
			// of errors.
			if !body.closed {
				t.Error("request body was not closed")
			}
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
			// Validate that the expected logs were logged.
			if gotLog := buf.String(); !strings.Contains(gotLog, tc.wantLog) {
				t.Errorf("got log %q, want %q", gotLog, tc.wantLog)
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
