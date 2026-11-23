package files

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
)

// downloadServer serves a file through presigned URLs it mints itself and
// through the workspace's Files API, with Range support on both, and counts
// which of them served each request.
type downloadServer struct {
	data []byte
	// missing answers 404 everywhere; refusePresigned answers presigned GETs
	// with a 403. holdFirstRange holds the first range after its first bytes
	// until another presigned GET arrives, for up to 5 s, so that a parallel
	// download shows it runs ranges at once.
	missing, refusePresigned, holdFirstRange bool

	mu             sync.Mutex
	mints          int
	presignedGets  int
	workspaceGets  int
	heads          int
	endpointsAsked int
	presignedRange []string
	inFlight       int
	maxInFlight    int
}

func newDownloadServer(t *testing.T, data []byte) (*downloadServer, *Client) {
	t.Helper()
	ds := &downloadServer{data: data}
	var srv *httptest.Server
	notFound := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error_code":"NOT_FOUND","message":"The file does not exist."}`)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/2.0/fs/endpoints", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		ds.count(&ds.endpointsAsked)
		_, _ = io.WriteString(w, `{"endpoints":[]}`)
	}))
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		ds.count(&ds.mints)
		if ds.missing {
			notFound(w)
			return
		}
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("credentials were sent to a presigned URL")
		}
		ds.mu.Lock()
		ds.presignedGets++
		ds.presignedRange = append(ds.presignedRange, r.Header.Get("Range"))
		ds.inFlight++
		ds.maxInFlight = max(ds.maxInFlight, ds.inFlight)
		ds.mu.Unlock()
		defer func() {
			ds.mu.Lock()
			ds.inFlight--
			ds.mu.Unlock()
		}()
		if ds.refusePresigned {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `<Error><Code>AccessDenied</Code></Error>`)
			return
		}
		ds.serve(w, r)
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		if ds.missing {
			notFound(w)
			return
		}
		if r.Method == http.MethodHead {
			ds.count(&ds.heads)
			w.Header().Set("Last-Modified", testLastModified)
			w.Header().Set("Content-Length", strconv.Itoa(len(ds.data)))
			return
		}
		ds.count(&ds.workspaceGets)
		ds.serve(w, r)
	}))
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return ds, buildUploadClient(t, srv.URL, srv.Client(), "")
}

func (ds *downloadServer) count(n *int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	*n++
}

// serve writes the file, or the range the request asks for.
func (ds *downloadServer) serve(w http.ResponseWriter, r *http.Request) {
	body, status := ds.data, http.StatusOK
	if spec, ok := strings.CutPrefix(r.Header.Get("Range"), "bytes="); ok {
		startText, endText, _ := strings.Cut(spec, "-")
		start, _ := strconv.Atoi(startText)
		end := len(ds.data) - 1
		if endText != "" {
			end, _ = strconv.Atoi(endText)
		}
		body, status = ds.data[start:end+1], http.StatusPartialContent
	}
	w.Header().Set("Last-Modified", testLastModified)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.WriteHeader(status)
	if ds.holdFirstRange && status == http.StatusPartialContent && strings.HasPrefix(r.Header.Get("Range"), "bytes=0-") && len(body) > 1024 {
		_, _ = w.Write(body[:1024])
		w.(http.Flusher).Flush()
		for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(time.Millisecond) {
			ds.mu.Lock()
			others := ds.maxInFlight > 1
			ds.mu.Unlock()
			if others {
				break
			}
		}
		body = body[1024:]
	}
	_, _ = w.Write(body)
}

// servedBy names the endpoint that served the file's bytes.
func (ds *downloadServer) servedBy(t *testing.T) string {
	t.Helper()
	ds.mu.Lock()
	defer ds.mu.Unlock()
	switch {
	case ds.workspaceGets > 0:
		return "workspace"
	case ds.presignedGets > 0 && !ds.refusePresigned:
		return "presigned"
	}
	return "none"
}

// downloadAll downloads the file with Download and reads it whole.
func downloadAll(t *testing.T, c *Client, opts ...DownloadOption) ([]byte, error) {
	t.Helper()
	resp, err := c.Download(t.Context(), "/Volumes/c/s/v/f.bin", opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Contents.Close()
	return io.ReadAll(resp.Contents)
}

func TestDownload_streamsTheFile(t *testing.T) {
	t.Parallel()
	data := payload(3000)
	ds, c := newDownloadServer(t, data)
	resp, err := c.Download(t.Context(), "/Volumes/c/s/v/f.bin")
	noErr(t, err)
	got, err := io.ReadAll(resp.Contents)
	noErr(t, err)
	noErr(t, resp.Contents.Close())
	eqBytes(t, got, data)
	if resp.ContentLength == nil || *resp.ContentLength != int64(len(data)) {
		t.Errorf("ContentLength = %v, want %d", resp.ContentLength, len(data))
	}
	if resp.LastModified == nil || *resp.LastModified != testLastModified {
		t.Errorf("LastModified = %v, want %q", resp.LastModified, testLastModified)
	}
	if got := ds.servedBy(t); got != "presigned" {
		t.Errorf("served by %s, want presigned", got)
	}
}

func TestDownloadTo_writesTheFile(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		size       int
		opts       []DownloadOption
		wantRanges bool
	}{
		{"small", 3000, nil, false},
		{"small, parallel", 3000, []DownloadOption{WithDownloadParallel(true)}, false},
		// Over 50 MiB, a parallel download goes in ranges of 10 MiB, several at
		// once.
		{"large, parallel", 52 << 20, []DownloadOption{WithDownloadParallel(true)}, true},
		{"large, sequential", 52 << 20, nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := payload(tc.size)
			ds, c := newDownloadServer(t, data)
			ds.holdFirstRange = true
			dest := filepath.Join(t.TempDir(), "out.bin")
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, tc.opts...)
			noErr(t, err)
			got, err := os.ReadFile(dest)
			noErr(t, err)
			eqBytes(t, got, data)
			if got := ds.servedBy(t); got != "presigned" {
				t.Errorf("served by %s, want presigned", got)
			}
			ds.mu.Lock()
			defer ds.mu.Unlock()
			if ranged := ds.presignedGets > 1; ranged != tc.wantRanges {
				t.Errorf("presigned GETs = %d (ranges %q), want ranges %v", ds.presignedGets, ds.presignedRange, tc.wantRanges)
			}
			if tc.wantRanges && ds.maxInFlight < 2 {
				t.Errorf("at most %d range at once, want several", ds.maxInFlight)
			}
		})
	}
}

func TestDownload_fallsBackToTheWorkspace(t *testing.T) {
	t.Parallel()
	// Presigned URLs refused before the first byte, as behind a storage
	// firewall: the workspace serves the file.
	data := payload(3000)
	t.Run("Download", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, data)
		ds.refusePresigned = true
		got, err := downloadAll(t, c)
		noErr(t, err)
		eqBytes(t, got, data)
		if got := ds.servedBy(t); got != "workspace" {
			t.Errorf("served by %s, want workspace", got)
		}
	})
	t.Run("DownloadTo", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, data)
		ds.refusePresigned = true
		dest := filepath.Join(t.TempDir(), "out.bin")
		_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadParallel(true))
		noErr(t, err)
		got, err := os.ReadFile(dest)
		noErr(t, err)
		eqBytes(t, got, data)
		if got := ds.servedBy(t); got != "workspace" {
			t.Errorf("served by %s, want workspace", got)
		}
	})
}

func TestDownload_usesTheTransferClient(t *testing.T) {
	t.Parallel()
	// A transfer client that refuses every request: presigned URLs fail through
	// it, so the workspace, reached through the client's own transport, serves
	// the file.
	data := payload(3000)
	refusing := func() (*http.Client, *int, *sync.Mutex) {
		var mu sync.Mutex
		calls := 0
		return &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			mu.Lock()
			calls++
			mu.Unlock()
			return nil, errors.New("refused by the transfer client")
		})}, &calls, &mu
	}
	check := func(t *testing.T, ds *downloadServer, calls *int, mu *sync.Mutex) {
		t.Helper()
		mu.Lock()
		defer mu.Unlock()
		if *calls == 0 {
			t.Error("the transfer client sent nothing")
		}
		if got := ds.servedBy(t); got != "workspace" {
			t.Errorf("served by %s, want workspace", got)
		}
	}
	t.Run("Download", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, data)
		hc, calls, mu := refusing()
		got, err := downloadAll(t, c, WithDownloadTransferClient(hc))
		noErr(t, err)
		eqBytes(t, got, data)
		check(t, ds, calls, mu)
	})
	t.Run("DownloadTo", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, data)
		hc, calls, mu := refusing()
		dest := filepath.Join(t.TempDir(), "out.bin")
		_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadTransferClient(hc))
		noErr(t, err)
		got, err := os.ReadFile(dest)
		noErr(t, err)
		eqBytes(t, got, data)
		check(t, ds, calls, mu)
	})
}

func TestDownloadTo_overwrite(t *testing.T) {
	t.Parallel()
	data := payload(3000)
	old := []byte("the old contents")
	t.Run("off keeps the file", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, data)
		dest := filepath.Join(t.TempDir(), "out.bin")
		noErr(t, os.WriteFile(dest, old, 0o600))
		_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadOverwrite(false))
		if !errors.Is(err, fs.ErrExist) {
			t.Fatalf("err = %v, want fs.ErrExist", err)
		}
		got, err := os.ReadFile(dest)
		noErr(t, err)
		eqBytes(t, got, old)
		ds.mu.Lock()
		defer ds.mu.Unlock()
		if ds.mints+ds.presignedGets+ds.workspaceGets != 0 {
			t.Errorf("sent %d mints and %d GETs, want none", ds.mints, ds.presignedGets+ds.workspaceGets)
		}
	})
	for name, opts := range map[string][]DownloadOption{
		"by default": nil,
		"on":         {WithDownloadOverwrite(true)},
	} {
		t.Run(name+" replaces the file", func(t *testing.T) {
			t.Parallel()
			_, c := newDownloadServer(t, data)
			dest := filepath.Join(t.TempDir(), "out.bin")
			noErr(t, os.WriteFile(dest, old, 0o600))
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, opts...)
			noErr(t, err)
			got, err := os.ReadFile(dest)
			noErr(t, err)
			eqBytes(t, got, data)
		})
	}
}

func TestDownloadTo_parallelismAppliesOnlyWhenParallel(t *testing.T) {
	t.Parallel()
	data := payload(3000)
	t.Run("without parallel, any count is ignored", func(t *testing.T) {
		t.Parallel()
		for _, n := range []int{0, 100} {
			ds, c := newDownloadServer(t, data)
			dest := filepath.Join(t.TempDir(), "out.bin")
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadParallelism(n))
			noErr(t, err)
			got, err := os.ReadFile(dest)
			noErr(t, err)
			eqBytes(t, got, data)
			if ds.heads != 0 {
				t.Errorf("parallelism %d: the download asked for the file's metadata, as only a parallel one does", n)
			}
		}
	})
	t.Run("with parallel, a count outside 1 to 64 fails before any request", func(t *testing.T) {
		t.Parallel()
		for _, n := range []int{0, -1, 65, 100} {
			ds, c := newDownloadServer(t, data)
			dest := filepath.Join(t.TempDir(), "out.bin")
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadParallel(true), WithDownloadParallelism(n))
			if err == nil || !strings.Contains(err.Error(), "parallelism must be between 1 and 64") {
				t.Errorf("parallelism %d: err = %v, want the range error", n, err)
			}
			if ds.heads+ds.mints != 0 {
				t.Errorf("parallelism %d: sent requests before failing", n)
			}
			if _, err := os.Stat(dest); !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("parallelism %d: the destination exists: %v", n, err)
			}
		}
	})
}

func TestDownloadConfig_resolvedParallelism(t *testing.T) {
	tun := defaultTunables()
	if got := resolveDownloadConfig([]DownloadOption{WithDownloadParallel(true)}).resolvedParallelism(tun); got != 10 {
		t.Errorf("default parallelism = %d, want 10", got)
	}
	if got := resolveDownloadConfig([]DownloadOption{WithDownloadParallel(true), WithDownloadParallelism(3)}).resolvedParallelism(tun); got != 3 {
		t.Errorf("parallelism = %d, want the 3 asked for", got)
	}
}

func TestDownload_missingFile(t *testing.T) {
	t.Parallel()
	isNotFound := func(t *testing.T, err error) {
		t.Helper()
		var aerr *apierr.APIError
		if !errors.As(err, &aerr) || aerr.Code() != codes.NotFound {
			t.Errorf("err = %v, want a NOT_FOUND API error", err)
		}
	}
	t.Run("Download", func(t *testing.T) {
		t.Parallel()
		ds, c := newDownloadServer(t, nil)
		ds.missing = true
		_, err := downloadAll(t, c)
		isNotFound(t, err)
	})
	for name, opts := range map[string][]DownloadOption{
		"DownloadTo":           nil,
		"DownloadTo, parallel": {WithDownloadParallel(true)},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ds, c := newDownloadServer(t, nil)
			ds.missing = true
			dir := t.TempDir()
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", filepath.Join(dir, "out.bin"), opts...)
			isNotFound(t, err)
			if extra := leftovers(t, dir); len(extra) > 0 {
				t.Errorf("left behind %v", extra)
			}
		})
	}
}

func TestDownloadTo_aDirectoryIsAnError(t *testing.T) {
	t.Parallel()
	ds, c := newDownloadServer(t, payload(3000))
	dir := t.TempDir()
	if _, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dir); err == nil {
		t.Fatal("a download to a directory succeeded")
	}
	if extra := leftovers(t, dir); len(extra) > 0 {
		t.Errorf("left behind %v", extra)
	}
	if ds.mints != 0 {
		t.Error("the download sent requests for a directory destination")
	}
}

func TestDownloadTo_parallelOnWindowsIsSequential(t *testing.T) {
	// Not parallel: it changes downloadGOOS. Only a parallel download asks for
	// the file's metadata, so the HEAD count shows which one ran.
	for _, tc := range []struct {
		goos      string
		wantHeads int
		wantWarn  bool
	}{
		{"linux", 1, false},
		{"windows", 0, true},
	} {
		t.Run(tc.goos, func(t *testing.T) {
			prev := downloadGOOS
			downloadGOOS = tc.goos
			t.Cleanup(func() { downloadGOOS = prev })
			data := payload(3000)
			ds, c := newDownloadServer(t, data)
			var logs bytes.Buffer
			c.logger = slog.New(slog.NewTextHandler(&logs, nil))
			dest := filepath.Join(t.TempDir(), "out.bin")
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, WithDownloadParallel(true))
			noErr(t, err)
			got, err := os.ReadFile(dest)
			noErr(t, err)
			eqBytes(t, got, data)
			if ds.heads != tc.wantHeads {
				t.Errorf("metadata requests = %d, want %d", ds.heads, tc.wantHeads)
			}
			warned := strings.Contains(logs.String(), "level=WARN") && strings.Contains(logs.String(), "Windows")
			if warned != tc.wantWarn {
				t.Errorf("warned = %v, want %v; logs:\n%s", warned, tc.wantWarn, logs.String())
			}
		})
	}
}

func TestDownload_ignoresTheFileOnlyOptions(t *testing.T) {
	t.Parallel()
	data := payload(3000)
	ds, c := newDownloadServer(t, data)
	got, err := downloadAll(t, c, WithDownloadOverwrite(false), WithDownloadParallel(true), WithDownloadParallelism(100))
	noErr(t, err)
	eqBytes(t, got, data)
	if got := ds.servedBy(t); got != "presigned" {
		t.Errorf("served by %s, want presigned", got)
	}
}

func TestDownload_readsTheEndpointsAPIEnv(t *testing.T) {
	// Not parallel: it sets the environment.
	downloads := map[string]func(t *testing.T, c *Client) error{
		"Download": func(t *testing.T, c *Client) error {
			_, err := downloadAll(t, c)
			return err
		},
		"DownloadTo": func(t *testing.T, c *Client) error {
			_, err := c.DownloadTo(t.Context(), "/Volumes/c/s/v/f.bin", filepath.Join(t.TempDir(), "out.bin"))
			return err
		},
	}
	for name, download := range downloads {
		for value, wantAsked := range map[string]bool{"": true, "true": true, "false": false} {
			t.Run(name+"/"+strconv.Quote(value), func(t *testing.T) {
				t.Setenv(endpointsAPIEnv, value)
				ds, c := newDownloadServer(t, payload(100))
				noErr(t, download(t, c))
				if asked := ds.endpointsAsked > 0; asked != wantAsked {
					t.Errorf("asked for the endpoint list = %v, want %v", asked, wantAsked)
				}
			})
		}
	}
}
