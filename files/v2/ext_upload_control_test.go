package files

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/ops"
)

// newControlEngine builds an engine whose control plane points at an httptest
// server running h. h is only reached by a request that carries the client's
// credentials: like the real control plane, an unauthenticated call is rejected
// rather than served, so every test here also covers that invariant.
func newControlEngine(t *testing.T, h http.HandlerFunc) *engine {
	t.Helper()
	srv := httptest.NewServer(requireCredentials(t, h))
	t.Cleanup(srv.Close)
	return newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
}

// controlContext returns the context of an upload of /Volumes/c/s/v/f.bin to
// e's workspace host, as the walk's last endpoint, so that its requests keep
// the API client's retries.
func controlContext(e *engine, overwrite *bool) *uploadContext {
	uc := &uploadContext{tun: e.tun, targetPath: "/Volumes/c/s/v/f.bin", overwrite: overwrite}
	uc.startEndpoint(endpoint{kind: presignedEndpoint, baseURL: e.c.host, multipart: true}, true)
	return uc
}

// requireCredentials rejects a control-plane request that does not carry
// testToken, the way the Files API rejects one with "Credential was not sent".
func requireCredentials(t *testing.T, h http.HandlerFunc) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != testToken {
			t.Errorf("%s %s: Authorization = %q, want %q", r.Method, r.URL.Path, got, testToken)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error_code":"UNAUTHENTICATED","message":"Credential was not sent or was of an unsupported type"}`)
			return
		}
		h(w, r)
	}
}

func TestInitiate(t *testing.T) {
	c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/2.0/fs/files/") {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("action"); got != "initiate-upload" {
			t.Errorf("action = %q", got)
		}
		if got := r.URL.Query().Get("overwrite"); got != "true" {
			t.Errorf("overwrite = %q, want true", got)
		}
		_, _ = io.WriteString(w, `{"multipart_upload":{"session_token":"tok"}}`)
	})

	ow := true
	res, err := c.initiateUpload(t.Context(), controlContext(c, &ow))
	if err != nil {
		t.Fatal(err)
	}
	if res.MultipartUpload == nil || res.MultipartUpload.SessionToken != "tok" {
		t.Errorf("multipart = %+v", res.MultipartUpload)
	}
}

func TestCreatePartURLs(t *testing.T) {
	c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"upload_part_urls":[{"url":"https://cloud.test/part/1","part_number":1,"headers":[{"name":"x-h","value":"v"}]}]}`)
	})

	urls, err := c.createPartURLs(t.Context(), controlContext(c, nil), "tok", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Fatalf("got %d URLs, want 1", len(urls))
	}
	if urls[0].URL != "https://cloud.test/part/1" || urls[0].PartNumber != 1 || urls[0].Headers["x-h"] != "v" {
		t.Errorf("url = %+v", urls[0])
	}
}

func TestCreatePartURLsEmpty(t *testing.T) {
	c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"upload_part_urls":[]}`)
	})
	if _, err := c.createPartURLs(t.Context(), controlContext(c, nil), "tok", 1, 1); !errors.Is(err, errUnexpectedServerResponse) {
		t.Fatalf("err = %v, want errUnexpectedServerResponse", err)
	}
}

func TestAbortMultipart(t *testing.T) {
	// The abort is an authenticated POST with the abort-upload action, sent once
	// even when it fails with a retryable status.
	var requests atomic.Int32
	c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		q := r.URL.Query()
		if r.Method != http.MethodPost || r.URL.Path != "/api/2.0/fs/files/Volumes/c/s/v/f.bin" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if q.Get("action") != "abort-upload" || q.Get("upload_type") != "multipart" || q.Get("session_token") != "tok" {
			t.Errorf("query = %v", q)
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	// Fast retries, so that a retried abort shows up as a second request rather
	// than as a test that runs out its retry budget.
	c.tun.cpBackoff = ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	c.tun.cpRetryTimeout = time.Second

	if err := c.abortMultipart(t.Context(), controlContext(c, nil), "tok"); err == nil {
		t.Fatal("expected the 503 to be returned")
	}
	if got := requests.Load(); got != 1 {
		t.Errorf("requests = %d, want 1", got)
	}
}

func TestCompleteMultipartSortsParts(t *testing.T) {
	var got []int
	c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.URL.Query().Get("action"); action != "complete-upload" {
			t.Errorf("action = %q", action)
		}
		var body struct {
			Parts []struct {
				PartNumber int    `json:"part_number"`
				ETag       string `json:"etag"`
			} `json:"parts"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		for _, p := range body.Parts {
			got = append(got, p.PartNumber)
		}
		w.WriteHeader(http.StatusOK)
	})

	err := c.completeMultipart(t.Context(), controlContext(c, nil), "tok", map[int]string{3: "e3", 1: "e1", 2: "e2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Errorf("parts = %v, want sorted [1 2 3]", got)
	}
}

func TestControlPlaneAlreadyExists(t *testing.T) {
	conflict := func(errorCode string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"error_code":"`+errorCode+`","message":"the file being created already exists"}`)
		}
	}

	// complete-upload (the multipart finish) returning 409 ALREADY_EXISTS maps to
	// the sentinel so callers can skip-if-exists, matching single-shot.
	c := newControlEngine(t, conflict("ALREADY_EXISTS"))
	if err := c.completeMultipart(t.Context(), controlContext(c, nil), "tok", map[int]string{1: "e1"}); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("completeMultipart err = %v, want ErrAlreadyExists", err)
	}

	// initiate-upload maps the same way.
	c2 := newControlEngine(t, conflict("ALREADY_EXISTS"))
	if _, err := c2.initiateUpload(t.Context(), controlContext(c2, nil)); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("initiateUpload err = %v, want ErrAlreadyExists", err)
	}

	// A different 409 is not mistaken for the already-exists sentinel.
	c3 := newControlEngine(t, conflict("RESOURCE_CONFLICT"))
	if _, err := c3.initiateUpload(t.Context(), controlContext(c3, nil)); errors.Is(err, ErrAlreadyExists) {
		t.Errorf("a non-ALREADY_EXISTS 409 mapped to ErrAlreadyExists: %v", err)
	}
}

func TestSingleShotUploadOverwrite(t *testing.T) {
	// The single-shot octet-stream PUT streams the raw body and sends the
	// overwrite query only when explicitly set.
	yes, no := true, false
	cases := []struct {
		name      string
		overwrite *bool
		wantQuery string
	}{
		{"unset", nil, ""},
		{"true", &yes, "true"},
		{"false", &no, "false"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody, gotOverwrite, gotContentType string
			c := newControlEngine(t, func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				gotBody = string(b)
				gotOverwrite = r.URL.Query().Get("overwrite")
				gotContentType = r.Header.Get("Content-Type")
				w.WriteHeader(http.StatusOK)
			})
			err := c.uploadSingleShot(t.Context(), controlContext(c, tc.overwrite), strings.NewReader("hello"), 5)
			if err != nil {
				t.Fatal(err)
			}
			if gotBody != "hello" {
				t.Errorf("body = %q, want the raw octet stream %q", gotBody, "hello")
			}
			if gotContentType != "application/octet-stream" {
				t.Errorf("Content-Type = %q, want application/octet-stream", gotContentType)
			}
			if gotOverwrite != tc.wantQuery {
				t.Errorf("overwrite query = %q, want %q", gotOverwrite, tc.wantQuery)
			}
		})
	}
}
