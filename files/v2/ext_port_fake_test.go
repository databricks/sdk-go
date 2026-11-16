//go:build files_port_harness

package files

// Temporary harness for porting the Python SDK's Files transfer layer to this
// package. Delete this file, ext_port_contract_test.go and the transfer
// contract in testdata/ together once the port is complete.
//
// portFake is an in-process fake of the endpoints a Files transfer talks to:
// the workspace (Files API 2.0, the upload control plane and the endpoints
// API), the storage proxy, a Files API 3.0 host, and a cloud store reached
// through presigned URLs. It mirrors FakeFilesServer in the Python SDK's
// tests/test_files_walk.py operation for operation, so that the scenarios
// recorded from that file run unchanged here. It is an http.RoundTripper:
// requests never reach the network.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
)

// errPortNetwork is the error a request fails with when its operation has a
// connection fault, as if the host could not be reached.
var errPortNetwork = errors.New("port fake: simulated network failure")

// portConnectionError is what a request with a connection fault fails with: a
// refused connection, as a real dial would report it, that also matches
// errPortNetwork. Python's requests.ConnectionError, which the contract was
// recorded with, is raised for the same failure.
type portConnectionError struct{ op *net.OpError }

func newPortConnectionError() error {
	return &portConnectionError{op: &net.OpError{
		Op:  "dial",
		Net: "tcp",
		Err: os.NewSyscallError("connect", syscall.ECONNREFUSED),
	}}
}

func (e *portConnectionError) Error() string   { return e.op.Error() }
func (e *portConnectionError) Unwrap() []error { return []error{errPortNetwork, e.op} }

// portFault describes how one call of an operation fails. The zero portFault, and a
// nil *portFault, mean the call succeeds.
type portFault struct {
	// Status, when not zero, answers the call with this HTTP status and a
	// Databricks error body.
	Status int
	// Connection fails the call with a refused connection before any response.
	Connection bool
	// Break, for a download, cuts the response body after BreakAfter bytes:
	// later reads fail with [io.ErrUnexpectedEOF].
	Break      bool
	BreakAfter int64
}

// portFaultFunc returns the fault for the call-th call of an operation, counting
// from 1, or nil for none.
type portFaultFunc func(call int) *portFault

// Always returns f on every call.
func portAlways(f *portFault) portFaultFunc { return func(int) *portFault { return f } }

// OnCalls returns f on the listed calls only.
func portOnCalls(f *portFault, calls ...int) portFaultFunc {
	return func(call int) *portFault {
		if slices.Contains(calls, call) {
			return f
		}
		return nil
	}
}

// portHosts names the hosts the fake answers for. The workspace host serves Files
// API 2.0 and the control plane; StorageProxy serves Files API 2.0; RIG serves
// Files API 3.0; Cloud serves presigned part uploads and downloads.
type portHosts struct {
	Workspace    string
	StorageProxy string
	RIG          string
	Cloud        string
}

// portFakeConfig configures a [portFake].
type portFakeConfig struct {
	portHosts
	// FilePath is the one remote file the fake stores and serves.
	FilePath string
	// LastModified is the Last-Modified header of downloads and metadata.
	LastModified string
	// Endpoints is the list GET /api/2.0/fs/endpoints answers with. When nil,
	// the workspace does not serve that API and answers 404.
	Endpoints []map[string]any
	// Faults maps an operation name to the faults of its calls.
	Faults map[string]portFaultFunc
	// DownloadContent is the content of the file downloads read.
	DownloadContent []byte
}

// portFake is the fake. Its methods are safe for concurrent use.
type portFake struct {
	cfg portFakeConfig

	mu          sync.Mutex
	calls       []string
	counts      map[string]int
	parts       map[string]map[int][]byte
	stored      map[string][]byte
	aborted     map[string]bool
	credentials map[string][]bool
	ranges      map[string][]string
	violations  []string
	unknown     []string
}

// New returns a portFake for cfg.
func newPortFake(cfg portFakeConfig) *portFake {
	return &portFake{
		cfg:         cfg,
		counts:      map[string]int{},
		parts:       map[string]map[int][]byte{},
		stored:      map[string][]byte{},
		aborted:     map[string]bool{},
		credentials: map[string][]bool{},
		ranges:      map[string][]string{},
	}
}

// Client returns an HTTP client whose every request goes to s.
func (s *portFake) Client() *http.Client { return &http.Client{Transport: s} }

// Operations on presigned cloud URLs. They must never carry workspace
// credentials.
var portPresignedOps = map[string]bool{"offload.part": true, "offload.get": true}

// Requests that pick the endpoint list or read file metadata rather than
// transfer data.
var portNotEndpoints = map[string]bool{"endpoints": true, "scim": true, "probe": true, "workspace": true}

// Error codes of the error bodies, by HTTP status.
var portErrorCodes = map[int]string{
	400: "INVALID_PARAMETER_VALUE",
	403: "PERMISSION_DENIED",
	404: "NOT_FOUND",
	409: "ALREADY_EXISTS",
	500: "INTERNAL_ERROR",
}

// RoundTrip implements [http.RoundTripper].
func (s *portFake) RoundTrip(req *http.Request) (*http.Response, error) {
	op, err := s.operation(req)
	if err != nil {
		return s.reject(req, err, &s.unknown), nil
	}
	s.mu.Lock()
	s.calls = append(s.calls, op)
	s.counts[op]++
	call := s.counts[op]
	s.credentials[op] = append(s.credentials[op], req.Header.Get("Authorization") != "")
	s.mu.Unlock()

	var fault *portFault
	if fn := s.cfg.Faults[op]; fn != nil {
		fault = fn(call)
	}
	if fault != nil && fault.Connection {
		portDrain(req)
		return nil, newPortConnectionError()
	}
	if fault != nil && fault.Status != 0 {
		portDrain(req)
		code := portErrorCodes[fault.Status]
		if code == "" {
			code = "UNKNOWN"
		}
		return portJSONResponse(req, fault.Status, map[string]any{
			"error_code": code,
			"message":    fmt.Sprintf("simulated %d on %s", fault.Status, op),
		}), nil
	}
	resp, err := s.respond(op, req, fault)
	if err != nil {
		return s.violation(req, err), nil
	}
	return resp, nil
}

// violation records a protocol check that failed.
func (s *portFake) violation(req *http.Request, err error) *http.Response {
	return s.reject(req, err, &s.violations)
}

// reject records err in list and answers 400, so that a client that swallows
// its request errors cannot hide it from the test.
func (s *portFake) reject(req *http.Request, err error, list *[]string) *http.Response {
	portDrain(req)
	s.mu.Lock()
	*list = append(*list, err.Error())
	s.mu.Unlock()
	return portJSONResponse(req, http.StatusBadRequest, map[string]any{"message": err.Error()})
}

func (s *portFake) operation(req *http.Request) (string, error) {
	host, method, path := req.URL.Hostname(), req.Method, req.URL.Path
	query := req.URL.Query()
	h := s.cfg.portHosts
	switch {
	case host == h.Workspace && path == "/api/2.0/fs/endpoints" && method == http.MethodGet:
		return "endpoints", nil
	case host == h.Workspace && path == "/api/2.0/preview/scim/v2/Me" && method == http.MethodGet:
		return "scim", nil
	case host == h.StorageProxy && path == "/api/2.0/fs/files/DatabricksInternal/Probe/fullstack/wis":
		return "probe", nil
	}
	if path == "/api/3.0/fs/files"+s.cfg.FilePath {
		if host != h.RIG {
			return "", fmt.Errorf("files API 3.0 request to %s", host)
		}
		if query.Has("upload_type") || query.Has("session_token") {
			return "", fmt.Errorf("2.0 parameters on %s", req.URL)
		}
		action := query.Get("action")
		switch {
		case method == http.MethodPost && (action == "initiate-upload" || action == "complete-upload" || action == "abort-upload"):
			return "rig." + strings.SplitN(action, "-", 2)[0], nil
		case method == http.MethodPut && query.Has("part_number"):
			return "rig.part", nil
		case method == http.MethodPut:
			return "rig.single", nil
		case method == http.MethodGet:
			return "rig.get", nil
		}
	}
	if path == "/api/2.0/fs/files"+s.cfg.FilePath {
		if host != h.Workspace && host != h.StorageProxy {
			return "", fmt.Errorf("files API 2.0 request to %s", host)
		}
		// The workspace's control plane is the "offload" endpoint and its
		// single-shot Files API the "floor"; the storage proxy is "dp" for both.
		control, transfer := "offload", "floor"
		if host == h.StorageProxy {
			control, transfer = "dp", "dp"
		}
		action := query.Get("action")
		switch {
		case method == http.MethodPost && action == "initiate-upload":
			return control + ".initiate", nil
		case method == http.MethodPost && action == "complete-upload":
			return control + ".complete", nil
		case method == http.MethodPost && action == "abort-upload":
			if query.Get("upload_type") != "multipart" {
				return "", fmt.Errorf("abort-upload without upload_type: %s", req.URL)
			}
			return control + ".abort", nil
		case method == http.MethodPut && query.Get("upload_type") == "multipart":
			if host != h.StorageProxy {
				return "", fmt.Errorf("multipart part sent to the workspace Files API: %s", req.URL)
			}
			return "dp.part", nil
		case method == http.MethodPut:
			return transfer + ".single", nil
		case method == http.MethodGet:
			return transfer + ".get", nil
		case method == http.MethodHead && host == h.Workspace:
			return "workspace.head", nil
		}
	}
	switch {
	case host == h.Workspace && path == "/api/2.0/fs/create-upload-part-urls":
		return "offload.part-urls", nil
	case host == h.Workspace && path == "/api/2.0/fs/create-download-url":
		return "offload.download-url", nil
	case host == h.Cloud && strings.HasPrefix(path, "/part/") && method == http.MethodPut:
		return "offload.part", nil
	case host == h.Cloud && path == "/download" && method == http.MethodGet:
		return "offload.get", nil
	}
	return "", fmt.Errorf("unexpected request %s %s", method, req.URL)
}

func (s *portFake) respond(op string, req *http.Request, fault *portFault) (*http.Response, error) {
	switch op {
	case "endpoints":
		portDrain(req)
		if s.cfg.Endpoints == nil {
			return portJSONResponse(req, http.StatusNotFound, map[string]any{"error_code": "NOT_FOUND", "message": "no endpoints API"}), nil
		}
		return portJSONResponse(req, http.StatusOK, map[string]any{"endpoints": s.cfg.Endpoints}), nil
	case "scim":
		portDrain(req)
		resp := portJSONResponse(req, http.StatusOK, map[string]any{})
		resp.Header.Set("X-Databricks-Org-Id", "12345")
		return resp, nil
	case "probe":
		portDrain(req)
		return portResponse(req, http.StatusOK, nil, nil), nil
	}
	endpoint, action, _ := strings.Cut(op, ".")
	query := req.URL.Query()
	if endpoint == "rig" && (action == "part" || action == "complete" || action == "abort") && query.Get("upload_id") != "upload-rig" {
		return nil, fmt.Errorf("%s without the upload_id from initiate-upload", op)
	}
	if endpoint == "rig" && (action == "part" || action == "single") && !portSendsContentLength(req) {
		return nil, fmt.Errorf("files API 3.0 requires a Content-Length on %s", op)
	}
	body, err := portReadBody(req)
	if err != nil {
		return nil, err
	}
	switch {
	case op == "rig.initiate":
		return portJSONResponse(req, http.StatusOK, map[string]any{"upload_id": "upload-rig"}), nil
	case action == "initiate":
		return portJSONResponse(req, http.StatusOK, map[string]any{
			"multipart_upload": map[string]any{"session_token": "token-" + endpoint},
		}), nil
	case op == "offload.part-urls":
		return s.partURLs(req, body)
	case action == "part":
		return s.part(req, endpoint, body)
	case action == "complete":
		return s.complete(req, endpoint, body)
	case action == "abort":
		if endpoint != "rig" && query.Get("session_token") != "token-"+endpoint {
			return nil, fmt.Errorf("%s with session_token %q", op, query.Get("session_token"))
		}
		s.mu.Lock()
		s.aborted[endpoint] = true
		s.mu.Unlock()
		return portJSONResponse(req, http.StatusOK, map[string]any{}), nil
	case action == "single":
		s.mu.Lock()
		s.stored[endpoint] = body
		s.mu.Unlock()
		return portJSONResponse(req, http.StatusOK, map[string]any{}), nil
	case op == "offload.download-url":
		return portJSONResponse(req, http.StatusOK, map[string]any{
			"url":     "https://" + s.cfg.Cloud + "/download",
			"headers": []any{},
		}), nil
	case op == "workspace.head":
		header := http.Header{}
		header.Set("Content-Length", strconv.Itoa(len(s.cfg.DownloadContent)))
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Last-Modified", s.cfg.LastModified)
		resp := portResponse(req, http.StatusOK, header, nil)
		resp.ContentLength = int64(len(s.cfg.DownloadContent))
		return resp, nil
	case action == "get":
		return s.download(req, endpoint, fault)
	}
	return nil, fmt.Errorf("no response for %s", op)
}

func (s *portFake) partURLs(req *http.Request, body []byte) (*http.Response, error) {
	var in struct {
		SessionToken    string `json:"session_token"`
		StartPartNumber int    `json:"start_part_number"`
		Count           int    `json:"count"`
	}
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("create-upload-part-urls body: %w", err)
	}
	if in.SessionToken != "token-offload" {
		return nil, fmt.Errorf("create-upload-part-urls with session_token %q", in.SessionToken)
	}
	urls := make([]map[string]any, 0, in.Count)
	for n := in.StartPartNumber; n < in.StartPartNumber+in.Count; n++ {
		urls = append(urls, map[string]any{
			"part_number": n,
			"url":         fmt.Sprintf("https://%s/part/%d", s.cfg.Cloud, n),
			"headers":     []any{},
		})
	}
	return portJSONResponse(req, http.StatusOK, map[string]any{"upload_part_urls": urls}), nil
}

func (s *portFake) part(req *http.Request, endpoint string, body []byte) (*http.Response, error) {
	var raw string
	if endpoint == "offload" {
		raw = req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
	} else {
		raw = req.URL.Query().Get("part_number")
	}
	number, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("part number %q: %w", raw, err)
	}
	s.mu.Lock()
	if s.parts[endpoint] == nil {
		s.parts[endpoint] = map[int][]byte{}
	}
	s.parts[endpoint][number] = body
	s.mu.Unlock()
	header := http.Header{}
	header.Set("ETag", fmt.Sprintf("etag-%s-%d", endpoint, number))
	return portResponse(req, http.StatusOK, header, nil), nil
}

func (s *portFake) complete(req *http.Request, endpoint string, body []byte) (*http.Response, error) {
	var in struct {
		Parts []struct {
			PartNumber int    `json:"part_number"`
			ETag       string `json:"etag"`
		} `json:"parts"`
	}
	if err := json.Unmarshal(body, &in); err != nil {
		return nil, fmt.Errorf("complete-upload body: %w", err)
	}
	numbers := make([]int, 0, len(in.Parts))
	for _, p := range in.Parts {
		if want := fmt.Sprintf("etag-%s-%d", endpoint, p.PartNumber); p.ETag != want {
			return nil, fmt.Errorf("complete-upload part %d with etag %q, want %q", p.PartNumber, p.ETag, want)
		}
		numbers = append(numbers, p.PartNumber)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	uploaded := make([]int, 0, len(s.parts[endpoint]))
	for n := range s.parts[endpoint] {
		uploaded = append(uploaded, n)
	}
	sort.Ints(uploaded)
	if !slices.Equal(numbers, uploaded) {
		return nil, fmt.Errorf("complete-upload must list every uploaded part: got %v, uploaded %v", numbers, uploaded)
	}
	var file bytes.Buffer
	for _, n := range numbers {
		file.Write(s.parts[endpoint][n])
	}
	s.stored[endpoint] = file.Bytes()
	return portJSONResponse(req, http.StatusOK, map[string]any{}), nil
}

func (s *portFake) download(req *http.Request, endpoint string, fault *portFault) (*http.Response, error) {
	rangeHeader := req.Header.Get("Range")
	s.mu.Lock()
	s.ranges[endpoint] = append(s.ranges[endpoint], rangeHeader)
	s.mu.Unlock()
	data, status := s.cfg.DownloadContent, http.StatusOK
	if rangeHeader != "" {
		spec, ok := strings.CutPrefix(rangeHeader, "bytes=")
		startText, endText, found := strings.Cut(spec, "-")
		if !ok || !found {
			return nil, fmt.Errorf("unsupported Range %q", rangeHeader)
		}
		start, err := strconv.Atoi(startText)
		if err != nil {
			return nil, fmt.Errorf("range start %q: %w", rangeHeader, err)
		}
		end := len(data) - 1
		if endText != "" {
			if end, err = strconv.Atoi(endText); err != nil {
				return nil, fmt.Errorf("range end %q: %w", rangeHeader, err)
			}
		}
		end = min(end, len(data)-1)
		if start > end+1 {
			return nil, fmt.Errorf("range %q beyond the %d-byte file", rangeHeader, len(data))
		}
		data, status = data[start:end+1], http.StatusPartialContent
	}
	header := http.Header{}
	header.Set("Content-Length", strconv.Itoa(len(data)))
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Last-Modified", s.cfg.LastModified)
	var body io.Reader = bytes.NewReader(data)
	if fault != nil && fault.Break {
		body = &portBreakingReader{data: data, breakAt: fault.BreakAfter}
	}
	resp := portResponse(req, status, header, body)
	resp.ContentLength = int64(len(data))
	return resp, nil
}

// portBreakingReader serves the first breakAt bytes, then fails like a dropped
// connection.
type portBreakingReader struct {
	data    []byte
	breakAt int64
	pos     int64
}

func (r *portBreakingReader) Read(p []byte) (int, error) {
	if r.pos >= r.breakAt {
		return 0, io.ErrUnexpectedEOF
	}
	end := min(r.breakAt, int64(len(r.data)), r.pos+int64(len(p)))
	if end <= r.pos {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:end])
	r.pos += int64(n)
	return n, nil
}

// Calls returns every operation in the order it was received.
func (s *portFake) Calls() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.calls)
}

// Counts returns how many times each operation was called.
func (s *portFake) Counts() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return maps.Clone(s.counts)
}

// EndpointsTouched returns the endpoints in the order they were first used,
// ignoring the probe, metadata and endpoint-list requests.
func (s *portFake) EndpointsTouched() []string {
	var seen []string
	for _, op := range s.Calls() {
		endpoint, _, _ := strings.Cut(op, ".")
		if !portNotEndpoints[endpoint] && !slices.Contains(seen, endpoint) {
			seen = append(seen, endpoint)
		}
	}
	return seen
}

// StoredFile returns the file as a caller would read it back, whichever
// endpoint stored it, and whether any endpoint did. It is an error for more
// than one endpoint to have stored it.
func (s *portFake) StoredFile() ([]byte, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.stored) > 1 {
		endpoints := make([]string, 0, len(s.stored))
		for e := range s.stored {
			endpoints = append(endpoints, e)
		}
		sort.Strings(endpoints)
		return nil, false, fmt.Errorf("the file was stored by more than one endpoint: %v", endpoints)
	}
	for _, data := range s.stored {
		return data, true, nil
	}
	return nil, false, nil
}

// Aborted reports whether endpoint received an abort-upload.
func (s *portFake) Aborted(endpoint string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.aborted[endpoint]
}

// Violations returns the protocol checks that failed: malformed requests to
// the contract's operations.
func (s *portFake) Violations() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.violations)
}

// UnknownRequests returns the requests that are none of the contract's
// operations. The Python SDK sends none, but a server may serve them.
func (s *portFake) UnknownRequests() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.unknown)
}

// CheckCredentials returns an error if workspace credentials reached a
// presigned URL, or were missing on a request to a Databricks host.
func (s *portFake) CheckCredentials() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ops := make([]string, 0, len(s.credentials))
	for op := range s.credentials {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	for _, op := range ops {
		sent := s.credentials[op]
		if portPresignedOps[op] && slices.Contains(sent, true) {
			return fmt.Errorf("workspace credentials were sent to a presigned URL (%s)", op)
		}
		if !portPresignedOps[op] && slices.Contains(sent, false) {
			return fmt.Errorf("workspace credentials were missing on %s", op)
		}
	}
	return nil
}

func portReadBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading the request body: %w", err)
	}
	return body, nil
}

// drain consumes and closes the request body, as a server that rejects a
// request still receives what the client sent.
func portDrain(req *http.Request) {
	if req.Body != nil {
		_, _ = io.Copy(io.Discard, req.Body)
		_ = req.Body.Close()
	}
}

// sendsContentLength reports whether the HTTP client would send a
// Content-Length header for req: it does when the length is known.
func portSendsContentLength(req *http.Request) bool {
	return req.ContentLength > 0 || (req.ContentLength == 0 && (req.Body == nil || req.Body == http.NoBody))
}

func portJSONResponse(req *http.Request, status int, v any) *http.Response {
	body, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("port fake: encoding a response: %v", err))
	}
	header := http.Header{}
	header.Set("Content-Type", "application/json")
	resp := portResponse(req, status, header, bytes.NewReader(body))
	resp.ContentLength = int64(len(body))
	return resp
}

func portResponse(req *http.Request, status int, header http.Header, body io.Reader) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	rc := io.NopCloser(bytes.NewReader(nil))
	if body != nil && req.Method != http.MethodHead {
		rc = io.NopCloser(body)
	}
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", status, http.StatusText(status)),
		StatusCode:    status,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          rc,
		ContentLength: -1,
		Request:       req,
	}
}

// The fake's own tests.
const portFilePath = "/Volumes/catalog/schema/volume/file.bin"

func newTestPortFake(cfg portFakeConfig) *portFake {
	cfg.portHosts = portHosts{Workspace: "localhost", StorageProxy: "proxy.example.com", RIG: "rig.example.com", Cloud: "cloud.example.com"}
	cfg.FilePath = portFilePath
	cfg.LastModified = "Thu, 28 Nov 2024 16:39:14 GMT"
	return newPortFake(cfg)
}

func portDo(t *testing.T, s *portFake, method, url string, body []byte, header map[string]string) *http.Response {
	t.Helper()
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(t.Context(), method, url, r)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := s.Client().Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func portRead(t *testing.T, resp *http.Response) string {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

var portAuth = map[string]string{"Authorization": "Bearer token"}

func TestPortFakePresignedMultipartUploadIsStored(t *testing.T) {
	s := newTestPortFake(portFakeConfig{})
	files := "http://localhost/api/2.0/fs/files" + portFilePath

	initiate := portRead(t, portDo(t, s, http.MethodPost, files+"?action=initiate-upload", nil, portAuth))
	if !strings.Contains(initiate, "token-offload") {
		t.Fatalf("initiate-upload = %s", initiate)
	}
	urls := portRead(t, portDo(t, s, http.MethodPost, "http://localhost/api/2.0/fs/create-upload-part-urls",
		[]byte(`{"session_token":"token-offload","start_part_number":1,"count":2}`), portAuth))
	if !strings.Contains(urls, "https://cloud.example.com/part/2") {
		t.Fatalf("create-upload-part-urls = %s", urls)
	}
	for n, part := range []string{"hello ", "world"} {
		resp := portDo(t, s, http.MethodPut, "https://cloud.example.com/part/"+string(rune('1'+n)), []byte(part), nil)
		if got, want := resp.Header.Get("ETag"), "etag-offload-"+string(rune('1'+n)); got != want {
			t.Fatalf("ETag = %q, want %q", got, want)
		}
	}
	complete := `{"parts":[{"part_number":1,"etag":"etag-offload-1"},{"part_number":2,"etag":"etag-offload-2"}]}`
	if resp := portDo(t, s, http.MethodPost, files+"?action=complete-upload&upload_type=multipart&session_token=token-offload",
		[]byte(complete), portAuth); resp.StatusCode != http.StatusOK {
		t.Fatalf("complete-upload = %d %s", resp.StatusCode, portRead(t, resp))
	}

	got, ok, err := s.StoredFile()
	if err != nil || !ok || string(got) != "hello world" {
		t.Fatalf("StoredFile = %q, %v, %v", got, ok, err)
	}
	if want := []string{"offload"}; !slices.Equal(s.EndpointsTouched(), want) {
		t.Errorf("EndpointsTouched = %v, want %v", s.EndpointsTouched(), want)
	}
	if err := s.CheckCredentials(); err != nil {
		t.Error(err)
	}
	if v := s.Violations(); len(v) != 0 {
		t.Errorf("Violations = %v", v)
	}
}

func TestPortFakeCredentialsOnPresignedURLAreReported(t *testing.T) {
	s := newTestPortFake(portFakeConfig{})
	portDo(t, s, http.MethodPut, "https://cloud.example.com/part/1", []byte("x"), portAuth)
	if err := s.CheckCredentials(); err == nil {
		t.Fatal("CheckCredentials = nil, want an error for credentials on a presigned URL")
	}
}

func TestPortFakeFaults(t *testing.T) {
	s := newTestPortFake(portFakeConfig{Faults: map[string]portFaultFunc{
		"floor.single": portOnCalls(&portFault{Status: http.StatusForbidden}, 1),
		"floor.get":    portAlways(&portFault{Connection: true}),
	}})
	files := "http://localhost/api/2.0/fs/files" + portFilePath

	if resp := portDo(t, s, http.MethodPut, files, []byte("x"), portAuth); resp.StatusCode != http.StatusForbidden ||
		!strings.Contains(portRead(t, resp), "PERMISSION_DENIED") {
		t.Fatalf("first single shot = %d, want 403 PERMISSION_DENIED", resp.StatusCode)
	}
	if resp := portDo(t, s, http.MethodPut, files, []byte("x"), portAuth); resp.StatusCode != http.StatusOK {
		t.Fatalf("second single shot = %d, want 200", resp.StatusCode)
	}
	req, _ := http.NewRequestWithContext(t.Context(), http.MethodGet, files, nil)
	if _, err := s.Client().Do(req); !errors.Is(err, errPortNetwork) || !errors.Is(err, syscall.ECONNREFUSED) {
		t.Fatalf("get with a connection fault: err = %v, want a refused connection", err)
	}
	if got := s.Counts()["floor.single"]; got != 2 {
		t.Errorf("floor.single calls = %d, want 2", got)
	}
}

func TestPortFakeDownloadRangeAndBreak(t *testing.T) {
	content := []byte("0123456789")
	s := newTestPortFake(portFakeConfig{
		DownloadContent: content,
		Faults:          map[string]portFaultFunc{"offload.get": portOnCalls(&portFault{Break: true, BreakAfter: 4}, 2)},
	})

	resp := portDo(t, s, http.MethodGet, "https://cloud.example.com/download", nil, map[string]string{"Range": "bytes=2-5"})
	if resp.StatusCode != http.StatusPartialContent || portRead(t, resp) != "2345" {
		t.Fatalf("ranged get = %d", resp.StatusCode)
	}
	resp = portDo(t, s, http.MethodGet, "https://cloud.example.com/download", nil, nil)
	got, err := io.ReadAll(resp.Body)
	if !errors.Is(err, io.ErrUnexpectedEOF) || string(got) != "0123" {
		t.Fatalf("broken get read %q, err %v; want \"0123\" and io.ErrUnexpectedEOF", got, err)
	}
}

func TestPortFakeEndpointsAPI(t *testing.T) {
	if resp := portDo(t, newTestPortFake(portFakeConfig{}), http.MethodGet, "http://localhost/api/2.0/fs/endpoints", nil, portAuth); resp.StatusCode != http.StatusNotFound {
		t.Errorf("endpoints without a list = %d, want 404", resp.StatusCode)
	}
	s := newTestPortFake(portFakeConfig{Endpoints: []map[string]any{{"api": "IO_OFFLOAD", "scheme": "http", "hostname": "localhost"}}})
	if body := portRead(t, portDo(t, s, http.MethodGet, "http://localhost/api/2.0/fs/endpoints", nil, portAuth)); !strings.Contains(body, "IO_OFFLOAD") {
		t.Errorf("endpoints = %s", body)
	}
}

func TestPortFakeUnexpectedRequestIsUnknown(t *testing.T) {
	s := newTestPortFake(portFakeConfig{})
	if resp := portDo(t, s, http.MethodGet, "http://localhost/api/2.0/somewhere", nil, portAuth); resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected request = %d, want 400", resp.StatusCode)
	}
	if len(s.UnknownRequests()) != 1 || len(s.Violations()) != 0 {
		t.Fatalf("UnknownRequests = %v, Violations = %v, want one unknown request", s.UnknownRequests(), s.Violations())
	}
}

func TestPortFakeFilesAPI3RequiresUploadID(t *testing.T) {
	s := newTestPortFake(portFakeConfig{})
	portDo(t, s, http.MethodPut, "https://rig.example.com/api/3.0/fs/files"+portFilePath+"?part_number=1&upload_id=wrong", []byte("x"), portAuth)
	if len(s.Violations()) != 1 {
		t.Fatalf("Violations = %v, want the missing upload_id", s.Violations())
	}
}
