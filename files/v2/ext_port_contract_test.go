//go:build files_port_harness

package files

// Temporary harness for porting the Python SDK's Files transfer layer to this
// package; delete it with ext_port_fake_test.go and the transfer contract in
// testdata/ once the port is complete.
//
// The transfer contract is the scenario matrix of the Python SDK's
// tests/test_files_walk.py, recorded with the outcome and requests the Python
// SDK produced. This test runs every scenario against this package through
// portFake. The recording is not distributed with this module; the test skips
// without it.
//
// The ratchet: a scenario on neither list must match the contract, a scenario
// in known_gaps.txt must not (remove it once it does), and a scenario in
// deviations.txt must have its recorded outcome instead of the contract's.

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/options/call"
)

const (
	contractHeaderFile = "testdata/transfer_contract.json"
	knownGapsFile      = "testdata/known_gaps.txt"
	deviationsFile     = "testdata/deviations.txt"

	// writeGapsEnv, when set, records the scenarios that fail instead of
	// failing on them, in known_gaps.txt under TEST_UNDECLARED_OUTPUTS_DIR (or
	// the directory the variable names). It is for the harness PR only: after
	// it, known_gaps.txt only shrinks.
	writeGapsEnv = "DATABRICKS_FILES_CONTRACT_WRITE_GAPS"
)

var contractScenarioFiles = []string{
	"testdata/transfer_contract_upload.jsonl",
	"testdata/transfer_contract_download.jsonl",
}

// errScenarioHung is the cause of a scenario's context running out, so that a
// hang is not mistaken for the SDK's own timeout.
var errScenarioHung = errors.New("scenario did not finish")

type contractHeader struct {
	Version   int `json:"version"`
	Constants struct {
		WorkspaceHost    string `json:"workspace_host"`
		StorageProxyHost string `json:"storage_proxy_host"`
		RIGHost          string `json:"rig_host"`
		CloudHost        string `json:"cloud_host"`
		FilePath         string `json:"file_path"`
		LastModified     string `json:"last_modified"`
		PartSize         int64  `json:"part_size"`
		FaultCalls       []int  `json:"fault_calls"`
	} `json:"constants"`
	Sizes         map[string]int `json:"sizes"`
	EndpointLists map[string]struct {
		Endpoints    []string         `json:"endpoints"`
		EndpointsAPI []map[string]any `json:"endpoints_api"`
	} `json:"endpoint_lists"`
}

type contractScenario struct {
	ID           string                     `json:"id"`
	Transfer     string                     `json:"transfer"`
	EndpointList string                     `json:"endpoint_list"`
	Size         string                     `json:"size"`
	Source       string                     `json:"source"`
	Mode         string                     `json:"mode"`
	Faults       map[string]json.RawMessage `json:"faults"`
	Expect       struct {
		Outcome       string         `json:"outcome"`
		EndpointsUsed []string       `json:"endpoints_used"`
		Requests      map[string]int `json:"requests"`
		RequestOps    []string       `json:"request_ops"`
	} `json:"expect"`
}

type encodedFault struct {
	Status     int    `json:"status"`
	Error      string `json:"error"`
	BreakAfter *int64 `json:"break_after"`
}

func TestTransferContract(t *testing.T) {
	header, scenarios := loadContract(t)
	gaps := loadList(t, knownGapsFile, 2)
	deviations := loadList(t, deviationsFile, 3)
	ids := map[string]bool{}
	for _, s := range scenarios {
		ids[s.ID] = true
	}
	for _, list := range []map[string][]string{gaps, deviations} {
		for id := range list {
			if !ids[id] {
				t.Errorf("%q is listed but is not a contract scenario; remove it", id)
			}
		}
	}

	_, recording := os.LookupEnv(writeGapsEnv)
	var mu sync.Mutex
	var failing []contractScenario
	t.Cleanup(func() {
		if recording {
			writeGaps(t, failing)
		}
	})

	for _, s := range scenarios {
		t.Run(s.ID, func(t *testing.T) {
			t.Parallel()
			want := s.Expect.Outcome
			if d, ok := deviations[s.ID]; ok {
				want = d[0]
			}
			mismatches, broken := runScenario(t, header, s, want)
			if len(broken) > 0 {
				t.Errorf("breaks a rule that holds whatever the Python SDK does, gap or not:\n%s", strings.Join(broken, "\n"))
			}
			_, gap := gaps[s.ID]
			switch {
			case recording:
				if len(mismatches) > 0 {
					mu.Lock()
					failing = append(failing, s)
					mu.Unlock()
				}
			case gap && len(mismatches) == 0:
				t.Errorf("passes now: remove it from %s", knownGapsFile)
			case !gap && len(mismatches) > 0:
				t.Errorf("does not match the contract:\n%s", strings.Join(mismatches, "\n"))
			}
		})
	}
}

// runScenario runs one scenario against the fake. It returns how the result
// differs from the contract, and separately the rules it breaks that hold
// whatever the Python SDK does: a known gap excuses only the former.
func runScenario(t *testing.T, h contractHeader, s contractScenario, wantOutcome string) (mismatches, broken []string) {
	t.Helper()
	k := h.Constants
	data := contractPayload(h.Sizes[s.Size])
	faults, err := decodeFaults(s.Faults, k.FaultCalls)
	if err != nil {
		t.Fatalf("decoding the faults: %v", err)
	}
	cfg := portFakeConfig{
		portHosts:    portHosts{Workspace: k.WorkspaceHost, StorageProxy: k.StorageProxyHost, RIG: k.RIGHost, Cloud: k.CloudHost},
		FilePath:     k.FilePath,
		LastModified: k.LastModified,
		Endpoints:    h.EndpointLists[s.EndpointList].EndpointsAPI,
		Faults:       faults,
	}
	if s.Transfer == "download" {
		cfg.DownloadContent = data
	}
	srv := newPortFake(cfg)
	e := newEngine(buildUploadClient(t, "http://"+k.WorkspaceHost, srv.Client(), ""))
	e.tun = contractTunables(e.tun, k.PartSize)

	ctx, cancel := context.WithTimeoutCause(t.Context(), 30*time.Second, errScenarioHung)
	defer cancel()
	var downloaded []byte
	if s.Transfer == "upload" {
		err = contractUpload(ctx, t, e, srv, k.FilePath, k.PartSize, s, data)
	} else {
		downloaded, err = contractDownload(ctx, e, k.FilePath)
	}
	if errors.Is(context.Cause(ctx), errScenarioHung) {
		hung := []string{"the scenario did not finish within 30s"}
		return hung, hung
	}

	if got := contractOutcome(err); got != wantOutcome {
		mismatches = append(mismatches, fmt.Sprintf("outcome = %s (%v), want %s", got, err, wantOutcome))
	}
	if got := srv.EndpointsTouched(); !slices.Equal(got, s.Expect.EndpointsUsed) {
		mismatches = append(mismatches, fmt.Sprintf("endpoints used = %v, want %v", got, s.Expect.EndpointsUsed))
	}
	counts := srv.Counts()
	if s.Expect.Requests != nil && !maps.Equal(counts, s.Expect.Requests) {
		mismatches = append(mismatches, fmt.Sprintf("requests = %v, want %v", counts, s.Expect.Requests))
	}
	if s.Expect.RequestOps != nil {
		ops := slices.Sorted(maps.Keys(counts))
		if !slices.Equal(ops, s.Expect.RequestOps) {
			mismatches = append(mismatches, fmt.Sprintf("operations = %v, want %v", ops, s.Expect.RequestOps))
		}
	}
	switch s.Transfer {
	case "upload":
		stored, ok, serr := srv.StoredFile()
		switch {
		case serr != nil:
			mismatches = append(mismatches, serr.Error())
		case err == nil && (!ok || !bytes.Equal(stored, data)):
			broken = append(broken, fmt.Sprintf("the upload succeeded but stored %d bytes (stored: %v), want the %d uploaded", len(stored), ok, len(data)))
		case wantOutcome == "OK" && (!ok || !bytes.Equal(stored, data)):
			mismatches = append(mismatches, fmt.Sprintf("stored %d bytes (stored: %v), want the %d uploaded", len(stored), ok, len(data)))
		case wantOutcome != "OK" && ok:
			mismatches = append(mismatches, "a file was stored by a failed upload")
		}
	case "download":
		switch {
		case err == nil && !bytes.Equal(downloaded, data):
			broken = append(broken, fmt.Sprintf("the download succeeded but returned %d bytes, want the %d stored", len(downloaded), len(data)))
		case wantOutcome == "OK" && !bytes.Equal(downloaded, data):
			mismatches = append(mismatches, fmt.Sprintf("downloaded %d bytes, want the %d stored", len(downloaded), len(data)))
		}
	}
	if err := srv.CheckCredentials(); err != nil {
		broken = append(broken, err.Error())
	}
	if u := srv.UnknownRequests(); len(u) > 0 {
		mismatches = append(mismatches, "requests the contract has no operation for: "+strings.Join(u, "; "))
	}
	if v := srv.Violations(); len(v) > 0 {
		broken = append(broken, "protocol violations: "+strings.Join(v, "; "))
	}
	if counts["probe"] > 0 || counts["scim"] > 0 {
		broken = append(broken, "probe or scim requests were sent")
	}
	return mismatches, broken
}

// contractTunables matches the thresholds the Python matrix sets, and makes
// retries immediate so that the scenarios run in milliseconds. Retry budgets
// stay finite, so that scenarios which exhaust them still do.
func contractTunables(tun tunables, partSize int64) tunables {
	tun.minStreamSize = partSize
	tun.defaultPartSize = partSize
	fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
	tun.cpBackoff = fast
	tun.cloudBackoff = fast
	tun.cpRetryTimeout = 200 * time.Millisecond
	return tun
}

func contractUpload(ctx context.Context, t *testing.T, e *engine, srv *portFake, path string, partSize int64, s contractScenario, data []byte) error {
	parallelism := 1
	if s.Mode == "parallel" {
		// The Python SDK's default parallelism.
		parallelism = 10
	}
	opts := []UploadOption{WithPartSize(partSize), WithParallelism(parallelism), WithTransferClient(srv.Client())}
	switch s.Source {
	case "file":
		local := filepath.Join(t.TempDir(), "source.bin")
		if err := os.WriteFile(local, data, 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := e.uploadFrom(ctx, path, local, opts...)
		return err
	case "seekable":
		_, err := e.upload(ctx, path, bytes.NewReader(data), opts...)
		return err
	case "non-seekable":
		_, err := e.upload(ctx, path, struct{ io.Reader }{bytes.NewReader(data)}, opts...)
		return err
	}
	t.Fatalf("unknown source %q", s.Source)
	return nil
}

// contractDownload downloads through the generated client, the only download
// this package has. Its retries are off because they back off for seconds.
func contractDownload(ctx context.Context, e *engine, path string) ([]byte, error) {
	resp, err := e.c.DownloadFile(ctx, DownloadFileRequest{FilePath: new(path)}, call.WithDisableRetry())
	if err != nil {
		return nil, err
	}
	defer resp.Contents.Close()
	return io.ReadAll(resp.Contents)
}

// contractOutcome maps an error to the contract's language-neutral outcome
// codes.
func contractOutcome(err error) string {
	switch {
	case err == nil:
		return "OK"
	case errors.Is(err, ErrAlreadyExists), apierr.Code(err) == codes.AlreadyExists:
		return "ALREADY_EXISTS"
	case apierr.Code(err) == codes.PermissionDenied:
		return "PERMISSION_DENIED"
	case apierr.Code(err) == codes.Internal:
		return "INTERNAL"
	case errors.Is(err, io.ErrUnexpectedEOF):
		return "TRANSFER_BROKEN"
	case errors.Is(err, context.DeadlineExceeded):
		return "RETRY_TIMEOUT"
	case errors.Is(err, errPortNetwork):
		return "CONNECTION_FAILED"
	case errors.Is(err, errUnexpectedServerResponse):
		return "UNEXPECTED_RESPONSE"
	}
	return "UNCLASSIFIED"
}

// contractPayload returns size bytes that are the same on every run.
func contractPayload(size int) []byte {
	data := make([]byte, size)
	r := rand.New(rand.NewChaCha8([32]byte{'f', 'i', 'l', 'e', 's'}))
	for i := range data {
		data[i] = byte(r.Uint32())
	}
	return data
}

// decodeFaults turns the contract's faults into fake faults. A fault is either
// one fault for every call, or one per call in calls, the last standing for
// every later call.
func decodeFaults(raw map[string]json.RawMessage, calls []int) (map[string]portFaultFunc, error) {
	faults := map[string]portFaultFunc{}
	for op, msg := range raw {
		var list []*encodedFault
		if err := json.Unmarshal(msg, &list); err != nil {
			var one *encodedFault
			if err := json.Unmarshal(msg, &one); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			list = []*encodedFault{one}
		}
		decoded := make([]*portFault, len(list))
		for i, f := range list {
			d, err := f.decode()
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			decoded[i] = d
		}
		faults[op] = func(call int) *portFault {
			if len(decoded) == 1 {
				return decoded[0]
			}
			i := sort.SearchInts(calls, call)
			return decoded[min(i, len(decoded)-1)]
		}
	}
	return faults, nil
}

func (f *encodedFault) decode() (*portFault, error) {
	switch {
	case f == nil:
		return nil, nil
	case f.Status != 0:
		return &portFault{Status: f.Status}, nil
	case f.Error == "connection":
		return &portFault{Connection: true}, nil
	case f.BreakAfter != nil:
		return &portFault{Break: true, BreakAfter: *f.BreakAfter}, nil
	}
	return nil, fmt.Errorf("unknown fault %+v", *f)
}

func loadContract(t *testing.T) (contractHeader, []contractScenario) {
	t.Helper()
	raw, err := os.ReadFile(contractHeaderFile)
	if errors.Is(err, fs.ErrNotExist) {
		// The public mirror does not ship the contract. Under Bazel it is the
		// test's data, so missing means the data is wired wrong.
		if _, bazel := os.LookupEnv("TEST_SRCDIR"); bazel {
			t.Fatalf("%s is missing from the test's data", contractHeaderFile)
		}
		t.Skip("the transfer contract is not available")
	}
	if err != nil {
		t.Fatal(err)
	}
	var h contractHeader
	if err := json.Unmarshal(raw, &h); err != nil {
		t.Fatalf("%s: %v", contractHeaderFile, err)
	}
	if h.Version != 1 {
		t.Fatalf("%s: version %d, want 1", contractHeaderFile, h.Version)
	}
	var scenarios []contractScenario
	for _, file := range contractScenarioFiles {
		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1<<20), 1<<20)
		for line := 1; scanner.Scan(); line++ {
			var s contractScenario
			if err := json.Unmarshal(scanner.Bytes(), &s); err != nil {
				t.Fatalf("%s:%d: %v", file, line, err)
			}
			scenarios = append(scenarios, s)
		}
		if err := scanner.Err(); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		f.Close()
	}
	return h, scenarios
}

// loadList reads a tab-separated list keyed by scenario ID. Blank lines and
// lines starting with # are ignored.
func loadList(t *testing.T, file string, fields int) map[string][]string {
	t.Helper()
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	list := map[string][]string{}
	for i, line := range strings.Split(string(raw), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != fields {
			t.Fatalf("%s:%d: %d tab-separated fields, want %d", file, i+1, len(parts), fields)
		}
		if _, dup := list[parts[0]]; dup {
			t.Fatalf("%s:%d: %q is listed twice", file, i+1, parts[0])
		}
		list[parts[0]] = parts[1:]
	}
	return list
}

// writeGaps writes the failing scenarios as known_gaps.txt, each with the
// phase of the porting plan expected to close it.
func writeGaps(t *testing.T, failing []contractScenario) {
	dir := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if dir == "" {
		dir = os.Getenv(writeGapsEnv)
	}
	slices.SortFunc(failing, func(a, b contractScenario) int { return strings.Compare(a.ID, b.ID) })
	var b strings.Builder
	b.WriteString("# Contract scenarios this SDK does not pass yet, with the porting phase\n")
	b.WriteString("# expected to close each. Entries are only ever removed.\n")
	for _, s := range failing {
		fmt.Fprintf(&b, "%s\t%s\n", s.ID, gapPhase(s))
	}
	if err := os.WriteFile(filepath.Join(dir, "known_gaps.txt"), []byte(b.String()), 0o644); err != nil {
		t.Errorf("writing known_gaps.txt: %v", err)
	}
}

func gapPhase(s contractScenario) string {
	switch {
	case s.Transfer == "download":
		return "P5"
	case strings.HasPrefix(s.EndpointList, "server:"):
		return "P3-P4"
	}
	return "P1-P2"
}
