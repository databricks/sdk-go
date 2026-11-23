package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/ops"
)

const testLastModified = "Thu, 28 Nov 2024 16:39:14 GMT"

// fileServer serves one file over the Files API on the workspace host, with
// Range support, and lets a test cut or stall a response.
type fileServer struct {
	t    *testing.T
	data []byte
	// lastModified is sent as Last-Modified unless empty.
	lastModified string
	// respond, when set, may write the response for the given 1-based call
	// itself; it returns false to let the file be served.
	respond func(call int, w http.ResponseWriter, r *http.Request, body []byte) bool

	mu      sync.Mutex
	calls   int
	headers []http.Header
}

func newFileServer(t *testing.T, data []byte) (*fileServer, *engine) {
	t.Helper()
	fs := &fileServer{t: t, data: data, lastModified: testLastModified}
	srv := httptest.NewServer(requireCredentials(t, fs.serve))
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
	fast := e.tun.cpBackoff
	fast.Initial, fast.Maximum, fast.Factor = time.Millisecond, time.Millisecond, 1
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	return fs, e
}

func (fs *fileServer) serve(w http.ResponseWriter, r *http.Request) {
	fs.mu.Lock()
	fs.calls++
	call := fs.calls
	fs.headers = append(fs.headers, r.Header.Clone())
	fs.mu.Unlock()
	body, status := fs.data, http.StatusOK
	if spec, ok := strings.CutPrefix(r.Header.Get("Range"), "bytes="); ok {
		startText, endText, _ := strings.Cut(spec, "-")
		start, _ := strconv.Atoi(startText)
		end := len(fs.data) - 1
		if endText != "" {
			end, _ = strconv.Atoi(endText)
		}
		body, status = fs.data[start:end+1], http.StatusPartialContent
	}
	if fs.lastModified != "" {
		w.Header().Set("Last-Modified", fs.lastModified)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if fs.respond != nil && fs.respond(call, w, r, body) {
		return
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func (fs *fileServer) requests() []http.Header {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.headers
}

// cutAfter writes the response's first n bytes and then drops the
// connection, as a network break does.
func cutAfter(w http.ResponseWriter, body []byte, n int) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body[:n])
	w.(http.Flusher).Flush()
	conn, _, err := w.(http.Hijacker).Hijack()
	if err == nil {
		_ = conn.Close()
	}
}

// workspaceDownload opens the whole file from the workspace's Files API, as
// the last endpoint.
func workspaceDownload(t *testing.T, e *engine) (*resilientDownload, error) {
	t.Helper()
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: filesAPIEndpoint, baseURL: e.c.host}, true, 1)
	return e.openDownload(t.Context(), dc, byteRange{start: 0, end: -1}, "")
}

func payload(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i * 7)
	}
	return b
}

func TestDownloadResumesAfterABreak(t *testing.T) {
	// A break mid-body resumes from the current offset, guarded by the first
	// response's Last-Modified.
	data := payload(5000)
	fs, e := newFileServer(t, data)
	fs.respond = func(call int, w http.ResponseWriter, _ *http.Request, body []byte) bool {
		if call == 1 {
			cutAfter(w, body, 1000)
			return true
		}
		return false
	}
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	got, err := io.ReadAll(d)
	noErr(t, err)
	eqBytes(t, got, data)
	reqs := fs.requests()
	if len(reqs) != 2 {
		t.Fatalf("requests = %d, want 2", len(reqs))
	}
	if got := reqs[1].Get("Range"); got != "bytes=1000-" {
		t.Errorf("resume Range = %q, want bytes=1000-", got)
	}
	if got := reqs[1].Get("If-Unmodified-Since"); got != testLastModified {
		t.Errorf("resume If-Unmodified-Since = %q, want the first response's Last-Modified", got)
	}
}

func TestDownloadDoesNotResumeWithoutLastModified(t *testing.T) {
	// Without a Last-Modified to guard it, a resumed download could splice two
	// versions of a file that changed in between.
	fs, e := newFileServer(t, payload(5000))
	fs.lastModified = ""
	fs.respond = func(_ int, w http.ResponseWriter, _ *http.Request, body []byte) bool {
		cutAfter(w, body, 1000)
		return true
	}
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	if _, err := io.ReadAll(d); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("err = %v, want the break", err)
	}
	if n := len(fs.requests()); n != 1 {
		t.Errorf("requests = %d, want 1", n)
	}
}

func TestDownloadResumesOnceWithoutProgress(t *testing.T) {
	// A resume that brings no byte is not tried again.
	fs, e := newFileServer(t, payload(5000))
	fs.respond = func(call int, w http.ResponseWriter, _ *http.Request, body []byte) bool {
		n := 0
		if call == 1 {
			n = 1000
		}
		cutAfter(w, body, n)
		return true
	}
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	if _, err := io.ReadAll(d); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("err = %v, want the break", err)
	}
	if n := len(fs.requests()); n != 2 {
		t.Errorf("requests = %d, want 2", n)
	}
}

func TestDownloadResumesAfterAStall(t *testing.T) {
	// A connection that stops delivering bytes is cut after the idle timeout,
	// and the download resumes on a fresh one.
	data := payload(5000)
	fs, e := newFileServer(t, data)
	e.tun.transferIdleTimeout = 50 * time.Millisecond
	fs.respond = func(call int, w http.ResponseWriter, r *http.Request, body []byte) bool {
		if call == 1 {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body[:1000])
			w.(http.Flusher).Flush()
			<-r.Context().Done()
			return true
		}
		return false
	}
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	got, err := io.ReadAll(d)
	noErr(t, err)
	eqBytes(t, got, data)
}

func TestDownloadCancelledByTheCallerIsNotResumed(t *testing.T) {
	fs, e := newFileServer(t, payload(5000))
	fs.respond = func(_ int, w http.ResponseWriter, r *http.Request, body []byte) bool {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body[:1000])
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		return true
	}
	ctx, cancel := context.WithCancel(t.Context())
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: filesAPIEndpoint, baseURL: e.c.host}, true, 1)
	d, err := e.openDownload(ctx, dc, byteRange{start: 0, end: -1}, "")
	noErr(t, err)
	if _, err := io.ReadFull(d, make([]byte, 500)); err != nil {
		t.Fatal(err)
	}
	cancel()
	if _, err := io.ReadAll(d); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want the caller's cancellation", err)
	}
	if n := len(fs.requests()); n != 1 {
		t.Errorf("requests = %d, want 1", n)
	}
}

func TestCheckRange(t *testing.T) {
	resp := func(status int, length int64) *http.Response {
		return &http.Response{StatusCode: status, ContentLength: length, Body: http.NoBody}
	}
	tests := []struct {
		name      string
		requested byteRange
		resp      *http.Response
		expected  int64
		wantOK    bool
	}{
		{"the whole file as a 200", byteRange{0, -1}, resp(200, 10), -1, true},
		{"a range as a 206", byteRange{5, 9}, resp(206, 5), 5, true},
		{"a range as a 200 of its length", byteRange{5, 9}, resp(200, 5), 5, true},
		{"a range as a 200 of the whole file", byteRange{5, 9}, resp(200, 10), 5, false},
		{"a range as a 200 of unknown length", byteRange{5, -1}, resp(200, -1), -1, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkRange(tc.requested, tc.resp, tc.expected)
			if (err == nil) != tc.wantOK {
				t.Errorf("checkRange err = %v, want ok %v", err, tc.wantOK)
			}
		})
	}
	if got := (byteRange{5, -1}).header(); got != "bytes=5-" {
		t.Errorf("header = %q", got)
	}
	if got := (byteRange{5, 9}).header(); got != "bytes=5-9" {
		t.Errorf("header = %q", got)
	}
}

const awsExpiredBody = `<?xml version="1.0" encoding="UTF-8"?><Error><Code>AccessDenied</Code><Message>Request has expired</Message></Error>`

// presignedServer serves the workspace's create-download-url and the
// presigned GET it mints, which answers with get for each 1-based call.
func presignedServer(t *testing.T, data []byte, get func(call int, w http.ResponseWriter) bool) (*engine, *atomic.Int32, *atomic.Int32) {
	t.Helper()
	var mints, gets atomic.Int32
	mux := http.NewServeMux()
	var srv *httptest.Server
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		mints.Add(1)
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("credentials were sent to a presigned URL")
		}
		call := int(gets.Add(1))
		if get != nil && get(call, w) {
			return
		}
		w.Header().Set("Last-Modified", testLastModified)
		_, _ = w.Write(data)
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Last-Modified", testLastModified)
		_, _ = w.Write(data)
	}))
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
	fast := e.tun.cpBackoff
	fast.Initial, fast.Maximum, fast.Factor = time.Millisecond, time.Millisecond, 1
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	return e, &mints, &gets
}

func TestPresignedDownloadRefetchesAnExpiredURL(t *testing.T) {
	data := payload(100)
	e, mints, _ := presignedServer(t, data, func(call int, w http.ResponseWriter) bool {
		if call == 1 {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, awsExpiredBody)
			return true
		}
		return false
	})
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: presignedEndpoint, baseURL: e.c.host, multipart: true}, true, 1)
	d, err := e.openDownload(t.Context(), dc, byteRange{start: 0, end: -1}, "")
	noErr(t, err)
	got, err := io.ReadAll(d)
	noErr(t, err)
	eqBytes(t, got, data)
	if int(mints.Load()) != 2 {
		t.Errorf("URLs minted = %d, want 2", mints.Load())
	}
}

func TestPresignedDownloadURLThatKeepsExpiring(t *testing.T) {
	e, mints, _ := presignedServer(t, payload(100), func(_ int, w http.ResponseWriter) bool {
		w.WriteHeader(http.StatusForbidden)
		_, _ = io.WriteString(w, awsExpiredBody)
		return true
	})
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: presignedEndpoint, baseURL: e.c.host, multipart: true}, true, 1)
	if _, err := e.openDownload(t.Context(), dc, byteRange{start: 0, end: -1}, ""); !errors.Is(err, errDownloadURLExpired) {
		t.Errorf("err = %v, want errDownloadURLExpired", err)
	}
	if want := e.tun.maxRetries + 1; int(mints.Load()) != want {
		t.Errorf("URLs minted = %d, want %d", mints.Load(), want)
	}
}

func TestPresignedDownloadURLMintFailureIsKept(t *testing.T) {
	// Every range of a download shares one URL: a failure to mint it is kept, so
	// none commits an endpoint that is failing.
	u := &presignedDownloadURL{}
	boom := errors.New("boom")
	calls := 0
	mint := func(context.Context) (*presignedURL, error) {
		calls++
		return nil, boom
	}
	for range 2 {
		if _, _, err := u.get(t.Context(), false, mint); !errors.Is(err, boom) {
			t.Errorf("err = %v, want %v", err, boom)
		}
	}
	if calls != 1 {
		t.Errorf("mints = %d, want 1", calls)
	}
	// Once the endpoint is committed, or for a cancelled request, a failure
	// is not kept: a later range or resume mints again.
	for name, get := range map[string]func(*presignedDownloadURL) error{
		"committed": func(u *presignedDownloadURL) error {
			_, _, err := u.get(t.Context(), true, mint)
			return err
		},
		"cancelled": func(u *presignedDownloadURL) error {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			_, _, err := u.get(ctx, false, mint)
			return err
		},
	} {
		u := &presignedDownloadURL{}
		calls = 0
		_ = get(u)
		_ = get(u)
		if calls != 2 {
			t.Errorf("%s: mints = %d, want 2", name, calls)
		}
	}
	// Only the version that expired is replaced.
	u = &presignedDownloadURL{}
	ok := func(context.Context) (*presignedURL, error) {
		calls++
		return &presignedURL{URL: strconv.Itoa(calls)}, nil
	}
	first, v, _ := u.get(t.Context(), false, ok)
	u.invalidate(v)
	u.invalidate(v)
	second, v2, _ := u.get(t.Context(), false, ok)
	if first.URL == second.URL || v2 != v+1 {
		t.Errorf("after invalidate: %s (v%d), then %s (v%d)", first.URL, v, second.URL, v2)
	}
}

func TestDownloadWalksOnlyBeforeTheFirstByte(t *testing.T) {
	data := payload(100)
	cases := []struct {
		name     string
		status   int
		wantErr  bool
		wantGets int
	}{
		{"a server error moves on to the workspace", http.StatusInternalServerError, false, 1},
		{"a changed file does not", http.StatusPreconditionFailed, true, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, _, gets := presignedServer(t, data, func(_ int, w http.ResponseWriter) bool {
				w.WriteHeader(tc.status)
				return true
			})
			// The workspace's Files API serves the file after the presigned URL.
			resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", downloadConfig{})
			if tc.wantErr {
				if err == nil {
					t.Fatal("want an error")
				}
			} else {
				if err != nil {
					t.Fatalf("download: %v", err)
				}
				defer resp.Contents.Close()
				got, err := io.ReadAll(resp.Contents)
				noErr(t, err)
				if !bytes.Equal(got, data) {
					t.Error("downloaded bytes differ")
				}
			}
			if int(gets.Load()) != tc.wantGets {
				t.Errorf("presigned GETs = %d, want %d", gets.Load(), tc.wantGets)
			}
		})
	}
}

func TestDownloadStopsAtTheRangeEnd(t *testing.T) {
	// Bytes past a range's end would overwrite the next range of a parallel
	// download, so a response that runs past it fails.
	fs, e := newFileServer(t, payload(5000))
	fs.respond = func(_ int, w http.ResponseWriter, _ *http.Request, _ []byte) bool {
		w.Header().Set("Content-Length", "150")
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(payload(150))
		return true
	}
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: filesAPIEndpoint, baseURL: e.c.host}, true, 1)
	d, err := e.openDownload(t.Context(), dc, byteRange{start: 0, end: 99}, "")
	noErr(t, err)
	got, err := io.ReadAll(d)
	if !errors.Is(err, errUnexpectedServerResponse) || len(got) > 100 {
		t.Errorf("read %d bytes, err = %v; want at most 100 and errUnexpectedServerResponse", len(got), err)
	}
}

func TestDownloadIdleTimeoutIgnoresTheCallersPauses(t *testing.T) {
	// Only a Read that stays blocked is a stall: a caller that pauses between
	// reads keeps its connection.
	data := payload(5000)
	fs, e := newFileServer(t, data)
	e.tun.transferIdleTimeout = 20 * time.Millisecond
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	var got []byte
	buf := make([]byte, 1000)
	for {
		n, err := d.Read(buf)
		got = append(got, buf[:n]...)
		if err == io.EOF {
			break
		}
		noErr(t, err)
		time.Sleep(50 * time.Millisecond)
	}
	eqBytes(t, got, data)
	if n := len(fs.requests()); n != 1 {
		t.Errorf("requests = %d, want 1", n)
	}
}

func TestDownloadWalksWhenTheBodyFailsBeforeTheFirstByte(t *testing.T) {
	// The presigned URL answers and then breaks before any byte: the download
	// moves on to the workspace, pinned to the file the first answer described.
	for _, changed := range []bool{false, true} {
		t.Run(fmt.Sprintf("file changed=%v", changed), func(t *testing.T) {
			data := payload(100)
			var mu sync.Mutex
			var pinned string
			mux := http.NewServeMux()
			var srv *httptest.Server
			mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
			}))
			mux.HandleFunc("/presigned", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Last-Modified", testLastModified)
				w.Header().Set("Content-Length", strconv.Itoa(len(data)))
				cutAfter(w, data, 0)
			})
			mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				pinned = r.Header.Get("If-Unmodified-Since")
				mu.Unlock()
				if changed {
					w.WriteHeader(http.StatusPreconditionFailed)
					return
				}
				_, _ = w.Write(data)
			}))
			srv = httptest.NewServer(mux)
			t.Cleanup(srv.Close)
			e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))

			resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", downloadConfig{})
			noErr(t, err)
			defer resp.Contents.Close()
			got, err := io.ReadAll(resp.Contents)
			if changed {
				if err == nil {
					t.Error("served a file that changed since the download started")
				}
			} else {
				noErr(t, err)
				eqBytes(t, got, data)
			}
			mu.Lock()
			defer mu.Unlock()
			if pinned != testLastModified {
				t.Errorf("If-Unmodified-Since = %q, want the first answer's Last-Modified", pinned)
			}
		})
	}
}

func TestDownloadRetriesARequestTheServerDropped(t *testing.T) {
	// A connection the server closes before answering fails the request with
	// io.ErrUnexpectedEOF, which the API client's rules don't retry; a download
	// request can always be sent again.
	data := payload(5000)
	fs, e := newFileServer(t, data)
	fs.respond = func(call int, w http.ResponseWriter, _ *http.Request, _ []byte) bool {
		if call == 1 {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				_ = conn.Close()
			}
			return true
		}
		return false
	}
	d, err := workspaceDownload(t, e)
	noErr(t, err)
	got, err := io.ReadAll(d)
	noErr(t, err)
	eqBytes(t, got, data)
	if n := len(fs.requests()); n != 2 {
		t.Errorf("requests = %d, want 2", n)
	}
}

func TestPresignedDownloadsUseHTTP1(t *testing.T) {
	// Over HTTP/2 the ranges of a parallel download would share one connection
	// to cloud storage. A Files API host keeps the client's protocols: a Files
	// API 3.0 host refuses HTTP/1.1 downloads.
	e := newEngine(&Client{internalClient: internalClient{host: "https://ws.example.com", httpClient: &http.Client{Transport: &http.Transport{}}}})
	presigned := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: presignedEndpoint, baseURL: "https://ws.example.com", multipart: true}, true, 4)
	if p := presigned.cloud.Transport.(*http.Transport).Protocols; p == nil || p.HTTP2() || !p.HTTP1() {
		t.Errorf("presigned protocols = %v, want HTTP/1.1 only", p)
	}
	rig := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: filesAPIv3Endpoint, baseURL: "https://rig.example.com", multipart: true}, true, 4)
	if p := rig.files.Transport.(*credentialTransport).base.(*http.Transport).Protocols; p != nil {
		t.Errorf("Files API 3.0 protocols = %v, want the client's", p)
	}
}

// dropConnection closes the request's connection without answering.
func dropConnection(w http.ResponseWriter) {
	if conn, _, err := w.(http.Hijacker).Hijack(); err == nil {
		_ = conn.Close()
	}
}

func TestDroppedRequestsKeepTheRetryBudget(t *testing.T) {
	data := payload(100)
	t.Run("before the first byte the download moves on", func(t *testing.T) {
		e, _, gets := presignedServer(t, data, func(_ int, w http.ResponseWriter) bool {
			dropConnection(w)
			return true
		})
		// A retry loop that ignored the budget would run until this timeout.
		e.tun.cpRetryTimeout = time.Second
		resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", downloadConfig{})
		noErr(t, err)
		defer resp.Contents.Close()
		got, err := io.ReadAll(resp.Contents)
		noErr(t, err)
		eqBytes(t, got, data)
		if int(gets.Load()) != 1 {
			t.Errorf("presigned GETs = %d, want 1: no retry before moving on", gets.Load())
		}
	})
	t.Run("once committed the download retries, then fails", func(t *testing.T) {
		e, _, gets := presignedServer(t, data, func(call int, w http.ResponseWriter) bool {
			if call == 1 {
				w.Header().Set("Last-Modified", testLastModified)
				w.Header().Set("Content-Length", strconv.Itoa(len(data)))
				cutAfter(w, data, 50)
				return true
			}
			dropConnection(w)
			return true
		})
		e.tun.cpRetryTimeout = time.Second
		resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", downloadConfig{})
		noErr(t, err)
		defer resp.Contents.Close()
		if _, err := io.ReadAll(resp.Contents); err == nil {
			t.Fatal("want the break")
		}
		if want := 2 + e.tun.cloudRetries; int(gets.Load()) != want {
			t.Errorf("presigned GETs = %d, want %d: the first, then the resume and its retries", gets.Load(), want)
		}
	})
}

func TestPresignedDownloadSendsTheURLsHeaders(t *testing.T) {
	// A cloud can sign a download URL with headers its GET must carry. One that
	// would replace a header the download sets itself is refused before the
	// GET goes out.
	data := payload(100)
	tests := []struct {
		name string
		// headers is what create-download-url returns with the URL.
		headers string
		rng     byteRange
		want    []byte
		wantErr error
		// wantGets is how many GETs the URL gets, the first with wantSig as
		// its X-Sig header.
		wantGets int
		wantSig  string
	}{
		{"a header the URL is signed with", `[{"name":"x-sig","value":"v1"}]`, byteRange{0, -1}, data, nil, 1, "v1"},
		{"a Range of its own on a range", `[{"name":"Range","value":"bytes=0-0"}]`, byteRange{10, 19}, nil, errUnexpectedServerResponse, 0, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gets atomic.Int32
			var sig atomic.Value
			mux := http.NewServeMux()
			var srv *httptest.Server
			mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(w, `{"url":%q,"headers":%s}`, srv.URL+"/presigned", tc.headers)
			}))
			mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
				if gets.Add(1) == 1 {
					sig.Store(r.Header.Get("X-Sig"))
				}
				w.Header().Set("Last-Modified", testLastModified)
				_, _ = w.Write(data)
			})
			srv = httptest.NewServer(mux)
			t.Cleanup(srv.Close)
			e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
			dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", endpoint{kind: presignedEndpoint, baseURL: e.c.host, multipart: true}, true, 1)

			var got []byte
			d, err := e.openDownload(t.Context(), dc, tc.rng, "")
			if err == nil {
				got, err = io.ReadAll(d)
				_ = d.Close()
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if !bytes.Equal(got, tc.want) {
				t.Errorf("read %d bytes, want %d", len(got), len(tc.want))
			}
			if int(gets.Load()) != tc.wantGets {
				t.Errorf("presigned GETs = %d, want %d", gets.Load(), tc.wantGets)
			}
			if got, _ := sig.Load().(string); got != tc.wantSig {
				t.Errorf("X-Sig = %q, want %q", got, tc.wantSig)
			}
		})
	}
}

// bodyTracker is a transport that counts the response bodies it hands out and
// how many of them were closed.
type bodyTracker struct {
	base http.RoundTripper
	// answered, when set, runs with each request once its answer arrived,
	// before the client gets it.
	answered func(r *http.Request)

	mu             sync.Mutex
	opened, closed int
}

func (b *bodyTracker) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := b.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	b.mu.Lock()
	b.opened++
	b.mu.Unlock()
	resp.Body = &trackedBody{ReadCloser: resp.Body, tracker: b}
	if b.answered != nil {
		b.answered(r)
	}
	return resp, nil
}

func (b *bodyTracker) counts() (opened, closed int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.opened, b.closed
}

type trackedBody struct {
	io.ReadCloser
	tracker *bodyTracker
	once    sync.Once
}

func (t *trackedBody) Close() error {
	t.once.Do(func() {
		t.tracker.mu.Lock()
		t.tracker.closed++
		t.tracker.mu.Unlock()
	})
	return t.ReadCloser.Close()
}

func TestDownloadClosedDuringARead(t *testing.T) {
	// Close may run while another goroutine's Read is blocked. The download
	// requests nothing after it, and closes every body it opened.
	data := payload(5000)
	firstBytes := func(w http.ResponseWriter, n int) {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data[:n])
		w.(http.Flusher).Flush()
	}
	serveFrom := func(w http.ResponseWriter, r *http.Request) {
		start, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.Header.Get("Range"), "bytes="), "-"))
		w.Header().Set("Content-Length", strconv.Itoa(len(data)-start))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(data[start:])
	}
	waitForClose := func(w http.ResponseWriter, r *http.Request) {
		firstBytes(w, 1000)
		<-r.Context().Done()
	}
	breakThenServe := func(call int, w http.ResponseWriter, r *http.Request) {
		if call == 1 {
			firstBytes(w, 1000)
			dropConnection(w)
			return
		}
		serveFrom(w, r)
	}
	tests := []struct {
		name string
		// presigned answers the presigned GET with the given 1-based number.
		presigned func(call int, w http.ResponseWriter, r *http.Request)
		// closeOn is the request on whose answer Close runs, before the
		// download takes it: "resume", the second presigned GET, or
		// "workspace". When empty, Close runs once the first bytes were read.
		closeOn string
		// rangeOnly reads one range from the presigned URL, without the walk.
		rangeOnly                    bool
		wantPresigned, wantWorkspace int
	}{
		{"a read that waits for bytes", func(_ int, w http.ResponseWriter, r *http.Request) { waitForClose(w, r) }, "", false, 1, 0},
		{"a range read that waits for bytes", func(_ int, w http.ResponseWriter, r *http.Request) { waitForClose(w, r) }, "", true, 1, 0},
		{"a resume answered", breakThenServe, "resume", false, 2, 0},
		{"a resume to be retried", func(call int, w http.ResponseWriter, r *http.Request) {
			if call == 2 {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			breakThenServe(call, w, r)
		}, "resume", false, 2, 0},
		{"the next endpoint answered", func(_ int, w http.ResponseWriter, _ *http.Request) {
			firstBytes(w, 0)
			dropConnection(w)
		}, "workspace", false, 1, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var presignedGets, workspaceGets atomic.Int32
			mux := http.NewServeMux()
			var srv *httptest.Server
			mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
			}))
			mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Last-Modified", testLastModified)
				tc.presigned(int(presignedGets.Add(1)), w, r)
			})
			mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
				workspaceGets.Add(1)
				w.Header().Set("Last-Modified", testLastModified)
				serveFrom(w, r)
			}))
			srv = httptest.NewServer(mux)
			t.Cleanup(srv.Close)

			answered, closed := make(chan struct{}), make(chan struct{})
			tracker := &bodyTracker{base: srv.Client().Transport, answered: func(r *http.Request) {
				switch {
				case tc.closeOn == "resume" && r.URL.Path == "/presigned" && presignedGets.Load() == 2,
					tc.closeOn == "workspace" && strings.HasPrefix(r.URL.Path, "/api/2.0/fs/files/"):
					close(answered)
					<-closed
				}
			}}
			client := &http.Client{Transport: tracker}
			e := newEngine(buildUploadClient(t, srv.URL, client, ""))
			fast := e.tun.cpBackoff
			fast.Initial, fast.Maximum, fast.Factor = time.Millisecond, time.Millisecond, 1
			e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
			cfg := downloadConfig{transferClient: client}

			var body io.ReadCloser
			if tc.rangeOnly {
				dc := e.newDownloadContext(cfg, "/Volumes/c/s/v/f.bin", endpoint{kind: presignedEndpoint, baseURL: e.c.host, multipart: true}, true, 1)
				d, err := e.openDownload(t.Context(), dc, byteRange{start: 0, end: -1}, "")
				noErr(t, err)
				body = d
			} else {
				resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", cfg)
				noErr(t, err)
				body = resp.Contents
			}

			readSome, done := make(chan struct{}), make(chan error, 1)
			go func() {
				if _, err := body.Read(make([]byte, 100)); err != nil {
					done <- err
					return
				}
				close(readSome)
				_, err := io.ReadAll(body)
				done <- err
			}()
			if tc.closeOn != "" {
				<-answered
			} else {
				<-readSome
				// Let the read block on the bytes that don't come.
				time.Sleep(50 * time.Millisecond)
			}
			noErr(t, body.Close())
			close(closed)

			select {
			case err := <-done:
				if !errors.Is(err, errDownloadClosed) {
					t.Errorf("read err = %v, want errDownloadClosed", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the read is still blocked after Close")
			}
			if int(presignedGets.Load()) != tc.wantPresigned || int(workspaceGets.Load()) != tc.wantWorkspace {
				t.Errorf("GETs: presigned %d, workspace %d; want %d, %d", presignedGets.Load(), workspaceGets.Load(), tc.wantPresigned, tc.wantWorkspace)
			}
			if opened, closedBodies := tracker.counts(); closedBodies != opened {
				t.Errorf("closed %d of the %d bodies opened", closedBodies, opened)
			}
		})
	}
}

func TestDownload_sendsTheSDKUserAgent(t *testing.T) {
	// The presigned URL refuses, so the download also walks to the workspace's
	// Files API.
	data := payload(100)
	var uas userAgents
	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		uas.record(r)
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		uas.record(r)
		w.WriteHeader(http.StatusForbidden)
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		uas.record(r)
		w.Header().Set("Last-Modified", testLastModified)
		_, _ = w.Write(data)
	}))
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))

	resp, err := e.download(t.Context(), "/Volumes/c/s/v/f.bin", downloadConfig{transferClient: srv.Client()})
	noErr(t, err)
	got, err := io.ReadAll(resp.Contents)
	_ = resp.Contents.Close()
	noErr(t, err)
	eqBytes(t, got, data)
	uas.check(t, e.c, func(path string) bool { return path == "/presigned" })
}

func TestDownloadURLMintEndsWithTheRangeBudget(t *testing.T) {
	// A committed range on a presigned endpoint whose URL-mint host never
	// answers. The first mint gives up after its retries, each ending at the
	// control-call bound, and the range retries. The second mint gets only
	// what is left of the range's budget, not a budget of its own.
	var mints atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/2.0/fs/create-download-url" {
			mints.Add(1)
		}
		<-r.Context().Done()
	}))
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
	fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	e.tun.cloudRetries = 1
	e.tun.controlAttemptTimeout = 200 * time.Millisecond
	e.tun.cpRetryTimeout = 500 * time.Millisecond
	// Not the workspace host, so the mint's retries are capped like the range's.
	ep := endpoint{kind: presignedEndpoint, baseURL: strings.Replace(srv.URL, "127.0.0.1", "localhost", 1), multipart: true}
	dc := e.newDownloadContext(downloadConfig{}, "/Volumes/c/s/v/f.bin", ep, true, 1)
	dc.committed.Store(true)

	start := time.Now()
	_, err := e.getRange(t.Context(), dc, byteRange{start: 0, end: -1}, "", -1)
	took := time.Since(start)
	if err == nil {
		t.Fatal("the range opened without a URL")
	}
	// Two mints for the first range attempt, and one cut at the range's budget
	// for the second; a mint with its own budget would make four.
	if n := mints.Load(); n != 3 {
		t.Errorf("mint requests = %d, want 3 (err %v)", n, err)
	}
	if took > 700*time.Millisecond {
		t.Errorf("the range took %v, want it to end near its %v budget", took, e.tun.cpRetryTimeout)
	}
}
