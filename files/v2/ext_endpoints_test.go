package files

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/databricks/sdk-go/auth"
)

func TestPlanUpload(t *testing.T) {
	e := newEngine(&Client{internalClient: internalClient{host: "https://ws.example.com"}})
	e.tun.minStreamSize = 1024
	presigned := endpoint{kind: presignedEndpoint, baseURL: "https://ws.example.com", multipart: true}
	workspace := endpoint{kind: filesAPIEndpoint, baseURL: "https://ws.example.com"}
	proxy := endpoint{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com", multipart: true}
	tests := []struct {
		name          string
		ep            endpoint
		size          int64
		wantOK        bool
		wantMultipart bool
	}{
		{"presigned takes a large upload in parts", presigned, 4096, true, true},
		{"presigned takes a stream of unknown size in parts", presigned, -1, true, true},
		{"presigned has no single shot for a small upload", presigned, 100, false, false},
		{"presigned has no single shot for an empty upload", presigned, 0, false, false},
		{"the workspace takes a large upload in a single shot", workspace, 4096, true, false},
		{"the workspace takes a stream of unknown size in a single shot", workspace, -1, true, false},
		{"the workspace takes up to 5 GiB in a single shot", workspace, workspaceMaxSingleShotSize, true, false},
		{"the workspace takes no more than 5 GiB in a single shot", workspace, workspaceMaxSingleShotSize + 1, false, false},
		{"another host has no single-shot limit", proxy.withoutMultipart(), workspaceMaxSingleShotSize + 1, true, false},
		{"a small upload goes in a single shot where multipart is offered", proxy, 100, true, false},
		{"an upload at the threshold goes in parts", proxy, 1024, true, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plan, ok := e.planUpload(tc.ep, tc.size)
			if ok != tc.wantOK || plan.multipart != tc.wantMultipart {
				t.Errorf("planUpload = (multipart %v, %v), want (multipart %v, %v)", plan.multipart, ok, tc.wantMultipart, tc.wantOK)
			}
		})
	}
}

func (ep endpoint) withoutMultipart() endpoint {
	ep.multipart = false
	return ep
}

func TestUploadContextRetries(t *testing.T) {
	tun := defaultTunables()
	tests := []struct {
		name        string
		committed   bool
		last        bool
		onWorkspace bool
		want        int
	}{
		{"before commit with an endpoint left", false, false, false, tun.cloudRetriesBeforeCommit},
		{"before commit with an endpoint left, on the workspace", false, false, true, tun.cloudRetriesBeforeCommit},
		{"on the last endpoint", false, true, false, tun.cloudRetries},
		{"on the last endpoint, on the workspace", false, true, true, -1},
		{"once committed", true, false, false, tun.cloudRetries},
		{"once committed, on the workspace", true, false, true, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uc := &uploadContext{tun: tun}
			uc.startEndpoint(endpoint{}, tc.last)
			if tc.committed {
				uc.commit()
			}
			if got := uc.retries(tc.onWorkspace); got != tc.want {
				t.Errorf("retries = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestSameOrigin(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"https://ws.example.com", "https://ws.example.com", true},
		{"https://WS.example.com", "https://ws.example.com/", true},
		{"https://ws.example.com:443", "https://ws.example.com", true},
		{"http://localhost", "http://localhost:80", true},
		{"http://ws.example.com", "https://ws.example.com", false},
		{"https://ws.example.com:8443", "https://ws.example.com", false},
		{"https://proxy.example.com", "https://ws.example.com", false},
		{"://bad", "://bad", false},
	}
	for _, tc := range tests {
		if got := sameOrigin(tc.a, tc.b); got != tc.want {
			t.Errorf("sameOrigin(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func readAllString(t *testing.T, r io.Reader) string {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestReplayStreamShortStream(t *testing.T) {
	// A stream that ends within the read-ahead learns its size, and its single
	// shot is a bytes.Reader that can be sent again.
	d := newStreamData(io.MultiReader(bytes.NewReader([]byte("hello"))))
	noErr(t, d.readAhead(10))
	if d.size() != 5 {
		t.Fatalf("size = %d, want 5", d.size())
	}
	body := d.whole()
	if _, ok := body.(io.Seeker); !ok {
		t.Errorf("whole() of an ended stream is %T, want a seekable reader", body)
	}
	if got := readAllString(t, body); got != "hello" {
		t.Errorf("whole() = %q, want hello", got)
	}
	if !d.canSendAgain() {
		t.Fatal("an ended stream can't be sent again")
	}
	noErr(t, d.restart())
	if got := readAllString(t, d.whole()); got != "hello" {
		t.Errorf("whole() after restart = %q, want hello", got)
	}
}

func TestReplayStreamRestartBeforeCommit(t *testing.T) {
	// Until commit, every byte read is kept, so the stream starts again from its
	// first byte.
	d := newStreamData(io.MultiReader(bytes.NewReader([]byte("abcdefghij"))))
	noErr(t, d.readAhead(4))
	if d.size() != -1 {
		t.Errorf("size = %d, want -1 before the stream ends", d.size())
	}
	first := make([]byte, 6)
	if _, err := io.ReadFull(d.r, first); err != nil {
		t.Fatal(err)
	}
	if !d.canSendAgain() {
		t.Fatal("can't send again before commit")
	}
	noErr(t, d.restart())
	if got := readAllString(t, d.r); got != "abcdefghij" {
		t.Errorf("after restart = %q, want the whole stream", got)
	}
}

func TestReplayStreamCommit(t *testing.T) {
	// Commit drops what was handed out, keeps what was read ahead but not yet
	// handed out, and stops keeping: the stream can't start again.
	s := &replayStream{src: io.MultiReader(bytes.NewReader([]byte("abcdefghij"))), keeping: true}
	d := &uploadData{r: s, stream: s, knownSize: -1}
	noErr(t, d.readAhead(6))
	part := make([]byte, 3)
	if _, err := io.ReadFull(d.r, part); err != nil {
		t.Fatal(err)
	}
	d.commit()
	if d.canSendAgain() {
		t.Error("can send again after commit")
	}
	if got := len(s.kept); got != 3 {
		t.Errorf("kept %d bytes after commit, want the 3 read ahead", got)
	}
	if got := readAllString(t, d.r); got != "defghij" {
		t.Errorf("after commit = %q, want the rest", got)
	}
	if s.kept != nil {
		t.Errorf("kept %d bytes once handed out, want none", len(s.kept))
	}
}

func TestReplayStreamWholeOfALongStream(t *testing.T) {
	// A single shot of a stream that did not end within what was kept streams
	// the rest, so the stream can't be sent again.
	d := newStreamData(io.MultiReader(bytes.NewReader([]byte("abcdefghij"))))
	noErr(t, d.readAhead(4))
	if got := readAllString(t, d.whole()); got != "abcdefghij" {
		t.Errorf("whole() = %q, want the whole stream", got)
	}
	if d.canSendAgain() {
		t.Error("can send again after streaming the rest")
	}
}

func TestReplayStreamReadAheadError(t *testing.T) {
	boom := errors.New("boom")
	d := newStreamData(io.MultiReader(bytes.NewReader([]byte("ab")), errReader{boom}))
	if err := d.readAhead(10); !errors.Is(err, boom) {
		t.Errorf("readAhead err = %v, want %v", err, boom)
	}
}

type errReader struct{ err error }

func (r errReader) Read([]byte) (int, error) { return 0, r.err }

func TestStorageProxyPreconditionFailedIsAlreadyExists(t *testing.T) {
	// The storage proxy answers an existing file with a 412 where the workspace
	// sends 409 ALREADY_EXISTS. Only an upload that doesn't overwrite maps it,
	// whichever request gets it.
	no, yes := false, true
	cases := []struct {
		name      string
		otherHost bool
		overwrite *bool
		want      bool
	}{
		{"the storage proxy without overwrite", true, &no, true},
		{"the storage proxy with overwrite", true, &yes, false},
		{"the storage proxy with the server's default", true, nil, false},
		{"the workspace without overwrite", false, &no, false},
	}
	requests := map[string]func(*engine, *uploadContext) error{
		"single shot": func(e *engine, uc *uploadContext) error {
			return e.uploadSingleShot(t.Context(), uc, strings.NewReader("hello"), 5)
		},
		"initiate": func(e *engine, uc *uploadContext) error {
			_, err := e.initiateUpload(t.Context(), uc)
			return err
		},
		"complete": func(e *engine, uc *uploadContext) error {
			return e.completeMultipart(t.Context(), uc, "tok", map[int]string{1: "e1"})
		},
	}
	for _, tc := range cases {
		for name, send := range requests {
			t.Run(tc.name+", "+name, func(t *testing.T) {
				c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
					_, _ = io.Copy(io.Discard, r.Body)
					w.WriteHeader(http.StatusPreconditionFailed)
				})
				uc := controlContext(c, tc.overwrite)
				if tc.otherHost {
					// The same test server under another name is another origin.
					uc.ep = endpoint{kind: filesAPIEndpoint, baseURL: strings.Replace(c.c.host, "127.0.0.1", "localhost", 1), multipart: true}
				}
				err := send(c, uc)
				if err == nil {
					t.Fatal("want an error")
				}
				if got := errors.Is(err, ErrAlreadyExists); got != tc.want {
					t.Errorf("errors.Is(%v, ErrAlreadyExists) = %v, want %v", err, got, tc.want)
				}
			})
		}
	}
}

func TestAuthenticatedHTTPClientIsSizedForTheParts(t *testing.T) {
	// Parts to a Files API host go over the client's transport, sized like the
	// transfer client so that parallel parts reuse their connections.
	own := &http.Transport{MaxIdleConnsPerHost: 2}
	e := newEngine(&Client{internalClient: internalClient{httpClient: &http.Client{Transport: own}}})
	ct, ok := e.authenticatedHTTPClient(16).Transport.(*credentialTransport)
	if !ok {
		t.Fatalf("transport = %T, want a credentialTransport", e.authenticatedHTTPClient(16).Transport)
	}
	sized, ok := ct.base.(*http.Transport)
	if !ok || sized == own || sized.MaxIdleConnsPerHost != 16 || sized.ResponseHeaderTimeout != e.tun.responseHeaderTimeout {
		t.Errorf("base = %#v, want a copy of the client's transport sized for 16 parts", ct.base)
	}
	if own.MaxIdleConnsPerHost != 2 {
		t.Error("the client's own transport was changed")
	}

	// Any other RoundTripper, such as a test's, is used as is.
	custom := roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, errors.New("unused") })
	e = newEngine(&Client{internalClient: internalClient{httpClient: &http.Client{Transport: custom}}})
	if ct := e.authenticatedHTTPClient(16).Transport.(*credentialTransport); ct.base == nil {
		t.Error("want the client's RoundTripper")
	}
}

type staticCredentials struct {
	headers []auth.Header
	err     error
}

func (staticCredentials) Name() string { return "static" }

func (c staticCredentials) AuthHeaders(context.Context) ([]auth.Header, error) {
	return c.headers, c.err
}

func TestCredentialTransport(t *testing.T) {
	var seen http.Header
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		seen = r.Header.Clone()
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: r}, nil
	})
	rt := &credentialTransport{base: base, creds: staticCredentials{headers: []auth.Header{{Key: "Authorization", Value: "Bearer t"}}}}
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPut, "http://proxy.example.com/x", nil)
	noErr(t, err)
	_, err = rt.RoundTrip(req)
	noErr(t, err)
	if got := seen.Get("Authorization"); got != "Bearer t" {
		t.Errorf("sent Authorization = %q, want the credentials", got)
	}
	if req.Header.Get("Authorization") != "" {
		t.Error("the caller's request was modified")
	}

	// A failure to get credentials sends nothing, and closes the body as a
	// RoundTripper must.
	boom := errors.New("boom")
	body := &closeRecorder{Reader: strings.NewReader("x")}
	req, err = http.NewRequestWithContext(t.Context(), http.MethodPut, "http://proxy.example.com/x", body)
	noErr(t, err)
	seen = nil
	rt.creds = staticCredentials{err: boom}
	if _, err := rt.RoundTrip(req); !errors.Is(err, boom) {
		t.Errorf("err = %v, want %v", err, boom)
	}
	if seen != nil || !body.closed.Load() {
		t.Errorf("sent = %v, body closed = %v; want nothing sent and the body closed", seen != nil, body.closed.Load())
	}
}

func TestFilesAPIv3Protocol(t *testing.T) {
	// Files API 3.0 is the 2.0 protocol under /api/3.0, naming the session by
	// upload_id, with parts of at most 4000 MiB.
	v2 := endpoint{kind: filesAPIEndpoint, baseURL: "http://proxy.example.com", multipart: true}
	v3 := endpoint{kind: filesAPIv3Endpoint, baseURL: "https://rig.example.com", multipart: true}
	if got := v3.filesPath("/Volumes/c/s/v/f.bin"); got != "/api/3.0/fs/files/Volumes/c/s/v/f.bin" {
		t.Errorf("3.0 path = %q", got)
	}
	if got := v2.filesPath("/Volumes/c/s/v/f.bin"); got != "/api/2.0/fs/files/Volumes/c/s/v/f.bin" {
		t.Errorf("2.0 path = %q", got)
	}
	if got := v3.sessionQuery("u1").Encode(); got != "upload_id=u1" {
		t.Errorf("3.0 session query = %q", got)
	}
	if got := v2.sessionQuery("t1").Encode(); got != "session_token=t1&upload_type=multipart" {
		t.Errorf("2.0 session query = %q", got)
	}

	uc := &uploadContext{tun: defaultTunables(), targetPath: "/Volumes/c/s/v/f.bin", basePartSize: 4 << 30}
	uc.startEndpoint(v3, false)
	if uc.partSize != filesAPIv3MaxPartSize {
		t.Errorf("3.0 part size = %d, want %d", uc.partSize, filesAPIv3MaxPartSize)
	}
	urls, err := filesAPIPartURLs(uc, "u1", 2, 2)
	noErr(t, err)
	if len(urls) != 2 || urls[0].URL != "https://rig.example.com/api/3.0/fs/files/Volumes/c/s/v/f.bin?part_number=2&upload_id=u1" || urls[1].PartNumber != 3 {
		t.Errorf("3.0 part URLs = %+v", urls)
	}
	uc.startEndpoint(v2, false)
	if uc.partSize != 4<<30 {
		t.Errorf("2.0 part size = %d, want the upload's %d", uc.partSize, int64(4<<30))
	}
}

func TestFilesAPIv3InitiateNeedsAnUploadID(t *testing.T) {
	for name, body := range map[string]string{
		"an upload_id":  `{"upload_id":"u1"}`,
		"no upload_id":  `{}`,
		"a 2.0 session": `{"multipart_upload":{"session_token":"t1"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/3.0/fs/files/Volumes/c/s/v/f.bin" || r.URL.Query().Get("action") != "initiate-upload" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL)
				}
				_, _ = io.WriteString(w, body)
			})
			uc := controlContext(c, nil)
			uc.startEndpoint(endpoint{kind: filesAPIv3Endpoint, baseURL: c.c.host, multipart: true}, true)
			token, err := c.initiateMultipart(t.Context(), uc)
			if name == "an upload_id" {
				noErr(t, err)
				if token != "u1" {
					t.Errorf("token = %q, want u1", token)
				}
				return
			}
			if !errors.Is(err, errUnexpectedServerResponse) {
				t.Errorf("err = %v, want errUnexpectedServerResponse", err)
			}
		})
	}
}
