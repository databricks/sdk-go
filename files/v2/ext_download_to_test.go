package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// rangeServer serves one file through the built-in endpoints: presigned URLs
// it mints itself, with Range support, and the workspace's Files API (GET and
// HEAD).
type rangeServer struct {
	data []byte
	// get, when set, may answer the given 1-based presigned GET itself; it
	// returns false to let the bytes be served.
	get func(call int, w http.ResponseWriter, body []byte) bool

	mu     sync.Mutex
	gets   int
	ranges []string
	// pins are the If-Unmodified-Since of each presigned GET.
	pins  []string
	heads int
}

func newRangeServer(t *testing.T, data []byte) (*rangeServer, *engine) {
	t.Helper()
	rs := &rangeServer{data: data}
	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/2.0/fs/create-download-url", requireCredentials(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `{"url":%q,"headers":[]}`, srv.URL+"/presigned")
	}))
	mux.HandleFunc("/presigned", func(w http.ResponseWriter, r *http.Request) {
		rs.mu.Lock()
		rs.gets++
		call := rs.gets
		rs.ranges = append(rs.ranges, r.Header.Get("Range"))
		rs.pins = append(rs.pins, r.Header.Get("If-Unmodified-Since"))
		rs.mu.Unlock()
		body, status := rs.data, http.StatusOK
		if spec, ok := strings.CutPrefix(r.Header.Get("Range"), "bytes="); ok {
			startText, endText, _ := strings.Cut(spec, "-")
			start, _ := strconv.Atoi(startText)
			end := len(rs.data) - 1
			if endText != "" {
				end, _ = strconv.Atoi(endText)
			}
			body, status = rs.data[start:end+1], http.StatusPartialContent
		}
		w.Header().Set("Last-Modified", testLastModified)
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		if rs.get != nil && rs.get(call, w, body) {
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write(body)
	})
	mux.HandleFunc("/api/2.0/fs/files/", requireCredentials(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", testLastModified)
		w.Header().Set("Content-Length", strconv.Itoa(len(rs.data)))
		if r.Method == http.MethodHead {
			rs.mu.Lock()
			rs.heads++
			rs.mu.Unlock()
			return
		}
		_, _ = w.Write(rs.data)
	}))
	srv = httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	e := newEngine(buildUploadClient(t, srv.URL, srv.Client(), ""))
	fast := e.tun.cpBackoff
	fast.Initial, fast.Maximum, fast.Factor = time.Millisecond, time.Millisecond, 1
	e.tun.cpBackoff, e.tun.cloudBackoff = fast, fast
	e.tun.downloadPartSize, e.tun.downloadMinParallelSize = 1024, 1024
	return rs, e
}

// leftovers lists the files in dir other than want.
func leftovers(t *testing.T, dir string, want ...string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	noErr(t, err)
	var extra []string
	for _, entry := range entries {
		if !slices.Contains(want, entry.Name()) {
			extra = append(extra, entry.Name())
		}
	}
	return extra
}

func TestDownloadToANewFile(t *testing.T) {
	for _, parallel := range []bool{false, true} {
		t.Run(fmt.Sprintf("parallel=%v", parallel), func(t *testing.T) {
			data := payload(5000)
			rs, e := newRangeServer(t, data)
			dir := t.TempDir()
			dest := filepath.Join(dir, "out.bin")
			noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{parallel: parallel, parallelism: new(4)}))
			got, err := os.ReadFile(dest)
			noErr(t, err)
			eqBytes(t, got, data)
			if extra := leftovers(t, dir, "out.bin"); len(extra) != 0 {
				t.Errorf("left behind %v", extra)
			}
			rs.mu.Lock()
			defer rs.mu.Unlock()
			if parallel {
				want := []string{"bytes=0-1023", "bytes=1024-2047", "bytes=2048-3071", "bytes=3072-4095", "bytes=4096-4999"}
				if got := slices.Sorted(slices.Values(rs.ranges)); !slices.Equal(got, want) || rs.heads != 1 {
					t.Errorf("ranges = %v, HEADs = %d; want %v and 1", got, rs.heads, want)
				}
				// Every range is pinned to the metadata's Last-Modified, so a file
				// that changes meanwhile fails rather than mixing two versions.
				for i, pin := range rs.pins {
					if pin != testLastModified {
						t.Errorf("range %s pinned to %q, want %q", rs.ranges[i], pin, testLastModified)
					}
				}
			} else if !slices.Equal(rs.ranges, []string{""}) || rs.heads != 0 {
				t.Errorf("ranges = %v, HEADs = %d; want one whole-file GET and none", rs.ranges, rs.heads)
			}
		})
	}
}

func TestDownloadToAFileThatFitsOneRangeIsWhole(t *testing.T) {
	rs, e := newRangeServer(t, payload(1024))
	dest := filepath.Join(t.TempDir(), "out.bin")
	noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{parallel: true, parallelism: new(4)}))
	if !slices.Equal(rs.ranges, []string{""}) {
		t.Errorf("ranges = %v, want one whole-file GET", rs.ranges)
	}
}

func TestFailedDownloadToANewFileLeavesNothing(t *testing.T) {
	for _, parallel := range []bool{false, true} {
		t.Run(fmt.Sprintf("parallel=%v", parallel), func(t *testing.T) {
			rs, e := newRangeServer(t, payload(5000))
			// The first range lands, which commits the presigned endpoint; the
			// others keep failing.
			rs.get = func(call int, w http.ResponseWriter, body []byte) bool {
				if call == 1 && parallel {
					return false
				}
				if call == 1 {
					cutAfter(w, body, 1000)
					return true
				}
				w.WriteHeader(http.StatusInternalServerError)
				return true
			}
			dir := t.TempDir()
			dest := filepath.Join(dir, "out.bin")
			if err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{parallel: parallel, parallelism: new(4)}); err == nil {
				t.Fatal("want an error")
			}
			if extra := leftovers(t, dir); len(extra) != 0 {
				t.Errorf("left behind %v", extra)
			}
		})
	}
}

func TestDownloadToAnExistingFile(t *testing.T) {
	data := payload(5000)
	rs, e := newRangeServer(t, data)
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.bin")
	noErr(t, os.WriteFile(dest, []byte(strings.Repeat("x", 9000)), 0o600))

	// An empty download still empties it.
	rs.data = nil
	noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{}))
	got, err := os.ReadFile(dest)
	noErr(t, err)
	if len(got) != 0 {
		t.Errorf("an empty download left %d bytes", len(got))
	}

	// A download replaces its contents in place, keeping its mode.
	rs.data = data
	noErr(t, os.Chmod(dest, 0o640))
	noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{parallel: true, parallelism: new(4)}))
	got, err = os.ReadFile(dest)
	noErr(t, err)
	eqBytes(t, got, data)
	info, err := os.Stat(dest)
	noErr(t, err)
	if info.Mode().Perm() != 0o640 {
		t.Errorf("mode = %v, want the file's own 0640", info.Mode().Perm())
	}
	if extra := leftovers(t, dir, "out.bin"); len(extra) != 0 {
		t.Errorf("left behind %v", extra)
	}
}

func TestFailedDownloadBeforeAnyByteLeavesAnExistingFile(t *testing.T) {
	rs, e := newRangeServer(t, payload(5000))
	rs.get = func(_ int, w http.ResponseWriter, _ []byte) bool {
		w.WriteHeader(http.StatusPreconditionFailed)
		return true
	}
	dest := filepath.Join(t.TempDir(), "out.bin")
	noErr(t, os.WriteFile(dest, []byte("keep"), 0o600))
	if err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dest, downloadConfig{}); err == nil {
		t.Fatal("want an error")
	}
	got, err := os.ReadFile(dest)
	noErr(t, err)
	if string(got) != "keep" {
		t.Errorf("contents = %q, want them kept", got)
	}
}

func TestDownloadToThroughASymlink(t *testing.T) {
	data := payload(100)
	_, e := newRangeServer(t, data)
	dir := t.TempDir()
	target := filepath.Join(dir, "target.bin")
	link := filepath.Join(dir, "link.bin")
	noErr(t, os.WriteFile(target, nil, 0o600))
	noErr(t, os.Symlink(target, link))
	noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", link, downloadConfig{}))
	got, err := os.ReadFile(target)
	noErr(t, err)
	eqBytes(t, got, data)
	if info, err := os.Lstat(link); err != nil || info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("the symlink was replaced: %v, %v", info, err)
	}
}

func TestDownloadToRefusals(t *testing.T) {
	rs, e := newRangeServer(t, payload(100))
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing.bin")
	noErr(t, os.WriteFile(existing, []byte("keep"), 0o600))

	if err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", existing, downloadConfig{noOverwrite: true}); !errors.Is(err, fs.ErrExist) {
		t.Errorf("err = %v, want fs.ErrExist", err)
	}
	if err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", dir, downloadConfig{}); err == nil {
		t.Error("downloaded to a directory")
	}
	if err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", filepath.Join(dir, "p.bin"), downloadConfig{parallel: true, parallelism: new(65)}); err == nil {
		t.Error("accepted a parallelism of 65")
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if rs.gets != 0 || rs.heads != 0 {
		t.Errorf("GETs = %d, HEADs = %d; want none", rs.gets, rs.heads)
	}
}

func TestTakesRanges(t *testing.T) {
	dir := t.TempDir()
	regular := filepath.Join(dir, "f.bin")
	noErr(t, os.WriteFile(regular, nil, 0o600))
	for path, want := range map[string]bool{
		filepath.Join(dir, "new.bin"): true,
		regular:                       true,
		os.DevNull:                    false,
	} {
		if got := takesRanges(path); got != want {
			t.Errorf("takesRanges(%s) = %v, want %v", path, got, want)
		}
	}
}

func TestSplitIntoRanges(t *testing.T) {
	got := splitIntoRanges(2500, 1024)
	want := []byteRange{{0, 1023}, {1024, 2047}, {2048, 2499}}
	if !slices.Equal(got, want) {
		t.Errorf("splitIntoRanges = %v, want %v", got, want)
	}
	if got := splitIntoRanges(2048, 1024); !slices.Equal(got, []byteRange{{0, 1023}, {1024, 2047}}) {
		t.Errorf("splitIntoRanges(2048) = %v", got)
	}
}

func TestDownloadDestinationWritesRangesAtTheirOffsets(t *testing.T) {
	// Ranges arrive in any order. An in-order write goes through the file's own
	// offset, which only such writes move; any other is positioned.
	dir := t.TempDir()
	d, err := newDownloadDestination(filepath.Join(dir, "out.bin"), 3072)
	noErr(t, err)
	want := payload(3072)
	for _, w := range []struct{ start, end int }{{0, 512}, {1024, 2048}, {2048, 3072}, {512, 1024}} {
		noErr(t, d.writeAt(int64(w.start), want[w.start:w.end]))
	}
	noErr(t, d.finish())
	got, err := os.ReadFile(filepath.Join(dir, "out.bin"))
	noErr(t, err)
	eqBytes(t, got, want)
}

func TestParallelDownloadStartsWithTheFirstRangeAlone(t *testing.T) {
	// The other ranges start only once the first range's first bytes arrived,
	// not merely its headers, so that a failure before the endpoint is
	// committed is the only request in flight.
	rs, e := newRangeServer(t, payload(5000))
	var firstBytes atomic.Bool
	var early atomic.Int32
	rs.get = func(call int, w http.ResponseWriter, body []byte) bool {
		if call != 1 {
			if !firstBytes.Load() {
				early.Add(1)
			}
			return false
		}
		w.WriteHeader(http.StatusPartialContent)
		w.(http.Flusher).Flush()
		time.Sleep(100 * time.Millisecond)
		firstBytes.Store(true)
		_, _ = w.Write(body)
		return true
	}
	noErr(t, e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", filepath.Join(t.TempDir(), "out.bin"), downloadConfig{parallel: true, parallelism: new(4)}))
	if early.Load() != 0 {
		t.Errorf("%d ranges started before the first range's first bytes", early.Load())
	}
}

func TestParallelDownloadStopsAfterAFailure(t *testing.T) {
	// A range that fails, here because the file changed, stops the ranges not
	// yet started.
	rs, e := newRangeServer(t, payload(5000))
	rs.get = func(call int, w http.ResponseWriter, _ []byte) bool {
		if call == 2 {
			w.WriteHeader(http.StatusPreconditionFailed)
			return true
		}
		return false
	}
	err := e.downloadTo(t.Context(), "/Volumes/c/s/v/f.bin", filepath.Join(t.TempDir(), "out.bin"), downloadConfig{parallel: true, parallelism: new(1)})
	if err == nil {
		t.Fatal("want an error")
	}
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if rs.gets != 2 {
		t.Errorf("GETs = %d, want 2: the later ranges never start", rs.gets)
	}
}

// eofHook calls onEOF when its body is read to the end.
type eofHook struct {
	io.ReadCloser
	onEOF func()
}

func (b *eofHook) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if err == io.EOF {
		b.onEOF()
	}
	return n, err
}

func TestCancelledParallelDownloadFails(t *testing.T) {
	// The caller cancels between two ranges, while none is in flight: the
	// download fails rather than finishing with ranges never written.
	_, e := newRangeServer(t, payload(5000))
	ctx, cancel := context.WithCancel(t.Context())
	var gets atomic.Int32
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp, err := http.DefaultTransport.RoundTrip(r)
		if err == nil && gets.Add(1) == 1 {
			resp.Body = &eofHook{ReadCloser: resp.Body, onEOF: cancel}
		}
		return resp, err
	})}
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.bin")
	err := e.downloadTo(ctx, "/Volumes/c/s/v/f.bin", dest, downloadConfig{transferClient: client, parallel: true, parallelism: new(1)})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want the cancellation", err)
	}
	if extra := leftovers(t, dir); len(extra) != 0 {
		t.Errorf("left behind %v", extra)
	}
}
