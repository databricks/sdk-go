package cloudstorage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/ops"
)

// fastClient returns a Client whose retrier backs off in milliseconds, so retry
// tests do not sleep.
func fastClient(srv *httptest.Server) *Client {
	fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	return New(srv.Client(), WithBackoff(fast))
}

func TestSendRetriesRetryableStatus(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	resp, err := fastClient(srv).Send(t.Context(), http.MethodPut, srv.URL, nil, BytesBody([]byte("x")), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if calls.Load() != 2 {
		t.Errorf("calls = %d, want 2 (one retry)", calls.Load())
	}
}

func TestSendExhaustsRetryableStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	_, err := fastClient(srv).Send(t.Context(), http.MethodPut, srv.URL, nil, BytesBody([]byte("x")), 3)
	aerr, ok := errors.AsType[*apierr.APIError](err)
	if !ok || aerr.HTTPStatusCode() != http.StatusServiceUnavailable {
		t.Fatalf("err = %v, want a 503 APIError after exhausting retries", err)
	}
}

func TestSendRetriesAtMostTheGivenCount(t *testing.T) {
	for _, retries := range []int{0, 1, 2} {
		var calls atomic.Int32
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls.Add(1)
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		_, err := fastClient(srv).Send(t.Context(), http.MethodPut, srv.URL, nil, BytesBody([]byte("x")), retries)
		srv.Close()
		if err == nil {
			t.Fatalf("retries=%d: want an error", retries)
		}
		if got := int(calls.Load()); got != retries+1 {
			t.Errorf("retries=%d: calls = %d, want %d", retries, got, retries+1)
		}
	}
}

func TestSendReturnsNonRetryableStatusAsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("nope"))
	}))
	defer srv.Close()

	resp, err := New(srv.Client()).Send(t.Context(), http.MethodPut, srv.URL, nil, BytesBody([]byte("x")), 3)
	if resp != nil {
		t.Errorf("expected nil response on error, got %v", resp)
	}
	aerr, ok := errors.AsType[*apierr.APIError](err)
	if !ok || aerr.HTTPStatusCode() != http.StatusForbidden {
		t.Fatalf("err = %v, want a 403 APIError", err)
	}
	if string(aerr.HTTPBody()) != "nope" {
		t.Errorf("HTTPBody = %q, want %q", aerr.HTTPBody(), "nope")
	}
}

func TestAttemptReturnsRawResponse(t *testing.T) {
	// Attempt never maps a status to an error: a 308 and its Range header come
	// back as a Response.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Range", "bytes=0-41")
		w.WriteHeader(http.StatusPermanentRedirect)
	}))
	defer srv.Close()

	resp, err := New(srv.Client()).Attempt(t.Context(), http.MethodPut, srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusPermanentRedirect {
		t.Errorf("status = %d, want 308", resp.StatusCode)
	}
	if got := resp.Header.Get("Range"); got != "bytes=0-41" {
		t.Errorf("Range = %q, want bytes=0-41", got)
	}
}

func TestIsURLExpired(t *testing.T) {
	azure := `<?xml version="1.0"?><Error><Code>AuthenticationFailed</Code><AuthenticationErrorDetail>Signature not valid in the specified time frame</AuthenticationErrorDetail></Error>`
	aws := `<?xml version="1.0"?><Error><Code>AccessDenied</Code><Message>Request has expired</Message></Error>`
	other := `<?xml version="1.0"?><Error><Code>AccessDenied</Code><Message>Access Denied</Message></Error>`

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"azure expired", apierr.FromHTTPError(http.StatusForbidden, nil, []byte(azure)), true},
		{"aws expired", apierr.FromHTTPError(http.StatusForbidden, nil, []byte(aws)), true},
		{"other 403", apierr.FromHTTPError(http.StatusForbidden, nil, []byte(other)), false},
		{"404", apierr.FromHTTPError(http.StatusNotFound, nil, []byte(azure)), false},
		{"nil", nil, false},
		{"non-apierr", errors.New("boom"), false},
	}
	for _, tc := range cases {
		if got := IsURLExpired(tc.err); got != tc.want {
			t.Errorf("%s: IsURLExpired = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestStallGuardStopsTimerOnClose verifies the idle timer is disarmed when the
// transport closes the request body, not only on a terminal Read. On a
// fixed-length request the transport normally issues a trailing EOF read that
// stops the timer, but Close is the reliable signal that body transmission is
// done: if the timer stayed armed into the response phase, a slow response could
// be cancelled as errStalled even though the upload succeeded. After Close, the
// cancel cause must never be set.
func TestStallGuardDoesNotArmTheTimerAfterClose(t *testing.T) {
	var cancelled atomic.Bool
	cancel := func(error) { cancelled.Store(true) }
	const idle = 20 * time.Millisecond
	body := stallGuardedBody(io.NopCloser(bytes.NewReader([]byte("payload"))), idle, cancel)

	// The transport can close the body before the goroutine that writes it reads
	// again, as when the server answers early or the stream is reset.
	if err := body.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := body.Read(make([]byte, 1)); err != nil {
		t.Fatalf("Read: %v", err)
	}

	time.Sleep(3 * idle)
	if cancelled.Load() {
		t.Fatal("stall timer fired after the body was closed")
	}
}

func TestStallGuardStopsTimerOnClose(t *testing.T) {
	var cancelled atomic.Bool
	cancel := func(error) { cancelled.Store(true) }

	// A tiny idle window so a leaked timer would fire quickly.
	const idle = 20 * time.Millisecond
	payload := []byte("payload")
	body := stallGuardedBody(io.NopCloser(bytes.NewReader(payload)), idle, cancel)

	// Simulate a transport that stops reading at Content-Length without a trailing
	// EOF read: read exactly len(payload) bytes (each Read returns n>0, nil, which
	// arms the timer), then Close. Only Close can disarm the timer here.
	buf := make([]byte, len(payload))
	if _, err := io.ReadFull(body, buf); err != nil {
		t.Fatalf("ReadFull: %v", err)
	}
	if err := body.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Wait past the idle window; a timer left armed by Close would fire here.
	time.Sleep(3 * idle)
	if cancelled.Load() {
		t.Fatal("stall timer fired after the body was fully read and closed")
	}
}

func TestExactReaderShortRead(t *testing.T) {
	var caught error
	er := &exactReader{r: bytes.NewReader(make([]byte, 10)), left: 100, cancel: func(c error) { caught = c }}

	_, err := io.ReadAll(er)
	if !errors.Is(err, errShortRead) {
		t.Fatalf("read error = %v, want errShortRead", err)
	}
	if !errors.Is(caught, errShortRead) {
		t.Errorf("attempt cancelled with %v, want errShortRead", caught)
	}
}

func TestAttemptKeepsThePresignedSignatureOutOfErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	presigned := srv.URL + "/bucket/key?X-Amz-Signature=secret"
	srv.Close() // so that the request fails before any response

	_, err := New(http.DefaultClient).Attempt(t.Context(), http.MethodPut, presigned, nil, BytesBody([]byte("x")))

	var uerr *url.Error
	if !errors.As(err, &uerr) {
		t.Fatalf("err = %v, want a *url.Error", err)
	}
	if want := srv.URL + "/bucket/key"; uerr.URL != want {
		t.Errorf("error URL = %q, want %q", uerr.URL, want)
	}
}

func TestAttemptKeepsTheSignatureOfAMalformedURLOutOfErrors(t *testing.T) {
	_, err := New(http.DefaultClient).Attempt(t.Context(), http.MethodPut, "https://bucket.test/key\x00?X-Amz-Signature=secret", nil, BytesBody([]byte("x")))

	if want := `parse "<unparsable URL>": net/url: invalid control character in URL`; err == nil || err.Error() != want {
		t.Errorf("err = %v, want %q", err, want)
	}
}

func TestRedactURLErrorKeepsATransportTimeout(t *testing.T) {
	signed := "https://bucket.test/key?sig=secret"
	inner := &url.Error{Op: "Put", URL: signed, Err: context.DeadlineExceeded}

	got := RedactURLError(&url.Error{Op: "Put", URL: signed, Err: inner})

	if netErr, ok := got.(net.Error); !ok || !netErr.Timeout() {
		t.Errorf("RedactURLError() = %v, want a net.Error that reports a timeout, as %v does", got, inner)
	}
}

func TestRedactURLError(t *testing.T) {
	plain := errors.New("boom")
	signed := "https://bucket.test/key?sig=secret"
	cases := []struct {
		name string
		err  func() error
		want string
	}{
		{"nil", func() error { return nil }, ""},
		{"not a url.Error", func() error { return plain }, "boom"},
		{
			"query and user info removed",
			func() error {
				return &url.Error{Op: "Put", URL: "https://u:p@bucket.test/key?sig=secret#f", Err: plain}
			},
			`Put "https://bucket.test/key": boom`,
		},
		{
			"an unparsable URL is replaced",
			func() error { return &url.Error{Op: "Put", URL: "https://bucket.test/%zz?sig=secret", Err: plain} },
			`Put "<unparsable URL>": boom`,
		},
		{
			"a url.Error that a transport returned, inside the client's",
			func() error {
				return &url.Error{Op: "Put", URL: signed, Err: &url.Error{Op: "Put", URL: signed, Err: plain}}
			},
			`Put "https://bucket.test/key": Put "https://bucket.test/key": boom`,
		},
		{
			"a transport's url.Errors joined",
			func() error {
				return &url.Error{Op: "Put", URL: signed, Err: errors.Join(&url.Error{Op: "Put", URL: signed, Err: plain}, &url.Error{Op: "Get", URL: signed, Err: plain})}
			},
			"Put \"https://bucket.test/key\": Put \"https://bucket.test/key\": boom\nGet \"https://bucket.test/key\": boom",
		},
		{
			"a transport's url.Error inside fmt.Errorf, whose message is fixed when it is built",
			func() error {
				return &url.Error{Op: "Put", URL: signed, Err: fmt.Errorf("proxy: %w", &url.Error{Op: "Put", URL: signed, Err: plain})}
			},
			`Put "https://bucket.test/key": proxy: Put "https://bucket.test/key": boom`,
		},
		{
			"a URL that begins another URL",
			func() error {
				return &url.Error{Op: "Get", URL: signed, Err: errors.Join(&url.Error{Op: "Get", URL: "https://bucket.test/key?a=1", Err: plain}, &url.Error{Op: "Get", URL: "https://bucket.test/key?a=1&sig=secret", Err: plain})}
			},
			"Get \"https://bucket.test/key\": Get \"https://bucket.test/key\": boom\nGet \"https://bucket.test/key\": boom",
		},
		{
			"a URL in the text of an error that wraps no url.Error, as on a client timeout",
			func() error {
				return &url.Error{Op: "Get", URL: signed, Err: fmt.Errorf("Get %q: context deadline exceeded (Client.Timeout exceeded while awaiting headers): %w", signed, plain)}
			},
			`Get "https://bucket.test/key": Get "https://bucket.test/key": context deadline exceeded (Client.Timeout exceeded while awaiting headers): boom`,
		},
		{
			"a redirect's unparsable Location, with user info and an escaped byte",
			func() error {
				return &url.Error{Op: "Get", URL: signed, Err: fmt.Errorf("failed to parse Location header %q: %w", "https://u:p@other.test/k\x00?sig=secret", plain)}
			},
			`Get "https://bucket.test/key": failed to parse Location header "https://other.test/k\x00": boom`,
		},
		{
			"a URL right after an escaped byte",
			func() error {
				return &url.Error{Op: "Get", URL: signed, Err: fmt.Errorf("proxy: %q: %w", "error:\nhttps://bucket.test/key?sig=secret", plain)}
			},
			`Get "https://bucket.test/key": proxy: "error:\nhttps://bucket.test/key": boom`,
		},
		{
			"a password with an @ in it",
			func() error {
				return &url.Error{Op: "Get", URL: signed, Err: fmt.Errorf("Get %q: %w", "https://u:p@secret@bucket.test/key", plain)}
			},
			`Get "https://bucket.test/key": Get "https://bucket.test/key": boom`,
		},
		{
			"a url.Error inside an error of another type",
			func() error { return fmt.Errorf("send: %w", &url.Error{Op: "Put", URL: signed, Err: plain}) },
			`send: Put "https://bucket.test/key": boom`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := tc.err()
			before := ""
			if in != nil {
				before = in.Error()
			}

			got := RedactURLError(in)

			gotMsg := ""
			if got != nil {
				gotMsg = got.Error()
			}
			if gotMsg != tc.want {
				t.Errorf("RedactURLError() = %q, want %q", gotMsg, tc.want)
			}
			if in != nil && in.Error() != before {
				t.Errorf("RedactURLError changed its input to %q, want %q", in.Error(), before)
			}
			if in != nil && !errors.Is(got, plain) {
				t.Errorf("errors.Is(RedactURLError(), plain) = false, want true")
			}
		})
	}
}
