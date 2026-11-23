package files

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

// FuzzEndpointEntry checks that any entry of the endpoints API either is
// skipped or gives an endpoint a transfer can safely use: a bare http or https
// origin, with plain http only where credentials may travel in the clear.
func FuzzEndpointEntry(f *testing.F) {
	for _, seed := range []string{
		`{"api":"IO_OFFLOAD","scheme":"https","hostname":"ws.example.com","multipart_upload_supported":true}`,
		`{"api":"FILES_API","api_version":"2.0","scheme":"http","hostname":"proxy.example.com:8080"}`,
		`{"api":"FILES_API","api_version":"3.0","scheme":"https","hostname":"[::1]:8443"}`,
		`{"api":"FILES_API","api_version":"2.0","scheme":"https","hostname":"user@host/path?q#f"}`,
		`{"api":"FILES_API","api_version":"3.0","scheme":"http","hostname":"rig.example.com"}`,
		`{"api":"IO_OFFLOAD","api_version":null,"scheme":"http","hostname":"localhost"}`,
		`{"api":["x"],"api_version":2,"scheme":1,"hostname":{}}`,
		`[]`, `null`, `"x"`,
	} {
		f.Add(seed, "https://ws.example.com")
		f.Add(seed, "http://localhost")
	}
	f.Fuzz(func(t *testing.T, raw, workspaceHost string) {
		var entry map[string]any
		if json.Unmarshal([]byte(raw), &entry) != nil {
			entry = nil
		}
		kind, ok := endpointKindOf(entry)
		if !ok {
			return
		}
		e := newEngine(&Client{internalClient: internalClient{host: workspaceHost}})
		ep, ok := e.parseEndpoint(entry, kind)
		if !ok {
			return
		}
		u, err := url.Parse(ep.baseURL)
		if err != nil {
			t.Fatalf("accepted %q, which doesn't parse: %v", ep.baseURL, err)
		}
		if u.Scheme != "http" && u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
			t.Fatalf("accepted %q, which is not a bare origin", ep.baseURL)
		}
		if u.Scheme == "http" {
			workspace, err := originOf(workspaceHost)
			toWorkspace := err == nil && strings.EqualFold(u.Hostname(), workspace.host)
			switch {
			case toWorkspace && workspace.scheme != "http":
				t.Fatalf("accepted plain http to a workspace reached over %s: %q", workspace.scheme, ep.baseURL)
			case !toWorkspace && kind != filesAPIEndpoint:
				t.Fatalf("accepted plain http to %q for %v", ep.baseURL, kind)
			}
		}
		if kind == filesAPIv3Endpoint && u.Scheme != "https" {
			t.Fatalf("accepted Files API 3.0 over %s", u.Scheme)
		}
	})
}

// FuzzSplitIntoRanges checks that the ranges of a parallel download cover the
// file exactly once, in order, each at most a part long.
func FuzzSplitIntoRanges(f *testing.F) {
	f.Add(int64(0), int64(1))
	f.Add(int64(4196), int64(1024))
	f.Add(int64(2048), int64(1024))
	f.Add(int64(1), int64(10<<20))
	f.Fuzz(func(t *testing.T, size, partSize int64) {
		if size < 0 || partSize < 1 || size/partSize > 1<<16 {
			t.Skip()
		}
		next := int64(0)
		for _, r := range splitIntoRanges(size, partSize) {
			if r.start != next || r.end < r.start || r.end-r.start+1 > partSize {
				t.Fatalf("range %+v after byte %d, part size %d", r, next, partSize)
			}
			next = r.end + 1
		}
		if next != size {
			t.Fatalf("ranges end at %d, want %d", next, size)
		}
	})
}
