package files

// Transfers through the endpoint walk, end to end against hostsFake: which
// endpoint serves an upload or a download as hosts go down, reject data or
// break mid-transfer, and what the caller gets back. Each case checks only
// what a user sees: the outcome, the bytes, roughly which endpoint served the
// transfer, and that no credentials reached a presigned URL and every request
// followed the protocol.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apierr/codes"
	"github.com/databricks/sdk-go/core/ops"
)

const (
	transferWorkspace    = "localhost"
	transferProxy        = "storage-proxy.databricks.com"
	transferRIG          = "rig.example.com"
	transferCloud        = "cloud.example.com"
	transferPath         = "/Volumes/catalog/schema/volume/file.bin"
	transferLastModified = "Thu, 28 Nov 2024 16:39:14 GMT"
	transferPartSize     = 1024
)

// transferSizes has one size for each way a file is routed: an empty file and
// a small one go in a single request, a larger one in several parts.
var transferSizes = map[string]int{"empty": 0, "small": 100, "several parts": 4196}

// transferLists are the endpoint lists the workspace answers with. Built-in
// means it doesn't serve the endpoints API, so a transfer walks presigned URLs
// and then the workspace.
var transferLists = map[string][]map[string]any{
	"built-in": nil,
	"proxy": {
		{"api": "FILES_API", "api_version": "2.0", "scheme": "http", "hostname": transferProxy, "multipart_upload_supported": true},
		transferPresignedEntry,
		transferWorkspaceEntry,
	},
	"rig": {
		{"api": "FILES_API", "api_version": "3.0", "scheme": "https", "hostname": transferRIG, "multipart_upload_supported": true},
		transferPresignedEntry,
		transferWorkspaceEntry,
	},
	"presigned only": {transferPresignedEntry},
}

var (
	transferPresignedEntry = map[string]any{"api": "IO_OFFLOAD", "scheme": "http", "hostname": transferWorkspace, "multipart_upload_supported": true}
	transferWorkspaceEntry = map[string]any{"api": "FILES_API", "api_version": "2.0", "scheme": "http", "hostname": transferWorkspace, "multipart_upload_supported": false}
)

// transferCase is one transfer. hosts lists the hosts that aren't up, each
// with its state: down (no connection), rejects (a 403 on every data
// request), breaks (the first data request goes through, or for a stream or a
// file download breaks after 1536 bytes, and later ones get no connection),
// fails once (only the first data request gets no connection), or breaks once
// (a download's first request breaks after 1536 bytes, later ones go
// through). servedBy is the last endpoint that saw data: "offload" for
// presigned URLs, "floor" for the workspace Files API, "dp" for the storage
// proxy and "rig" for the Files API 3.0 host.
type transferCase struct {
	list     string
	hosts    string
	size     string
	source   string // file, seekable or non-seekable; uploads only
	mode     string // sequential or parallel for uploads; stream, to-file or parallel for downloads
	want     string // ok, deadline, connection, broken or permission denied
	servedBy string
}

func (c transferCase) name() string {
	hosts := c.hosts
	if hosts == "" {
		hosts = "all up"
	}
	what := c.size
	if c.source != "" {
		what += ", " + c.source
	}
	return c.list + " | " + hosts + " | " + what + ", " + c.mode
}

var uploadCases = []transferCase{
	{list: "built-in", hosts: "", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "workspace down", size: "small", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "workspace down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "empty", source: "seekable", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "small", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "several parts", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "empty", source: "seekable", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned breaks", size: "several parts", source: "file", mode: "sequential", want: "connection", servedBy: "offload"},
	{list: "built-in", hosts: "presigned breaks", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "offload"},
	{list: "proxy", hosts: "", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "proxy down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy rejects", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy rejects", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "small", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "several parts", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "empty", source: "seekable", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy breaks", size: "several parts", source: "file", mode: "sequential", want: "connection", servedBy: "dp"},
	{list: "proxy", hosts: "proxy breaks", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "dp"},
	{list: "rig", hosts: "", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "rig down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig rejects", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig rejects", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "small", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "several parts", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "empty", source: "seekable", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig breaks", size: "several parts", source: "file", mode: "sequential", want: "connection", servedBy: "rig"},
	{list: "rig", hosts: "rig breaks", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "rig"},
	{list: "presigned only", hosts: "", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned rejects", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned rejects", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned rejects", size: "several parts", source: "non-seekable", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned rejects", size: "empty", source: "seekable", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down, workspace down", size: "small", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down, workspace down", size: "several parts", source: "file", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down, workspace down", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned down, workspace down", size: "empty", source: "seekable", mode: "sequential", want: "deadline", servedBy: "floor"},
	{list: "presigned only", hosts: "presigned breaks", size: "several parts", source: "file", mode: "sequential", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned breaks", size: "several parts", source: "non-seekable", mode: "parallel", want: "connection", servedBy: "offload"},
	// A host that fails once is retried: on the workspace with its own retries,
	// and on another host before any part lands with the one retry it gets then.
	{list: "built-in", hosts: "workspace fails once", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned fails once", size: "several parts", source: "file", mode: "sequential", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy fails once", size: "several parts", source: "seekable", mode: "parallel", want: "ok", servedBy: "dp"},
	{list: "rig", hosts: "rig fails once", size: "small", source: "file", mode: "sequential", want: "ok", servedBy: "rig"},
}

var downloadCases = []transferCase{
	{list: "built-in", hosts: "", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "presigned down", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "several parts", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "several parts", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "small", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down", size: "empty", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "several parts", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "several parts", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "small", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "presigned rejects", size: "empty", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "built-in", hosts: "workspace down", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "workspace down", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "several parts", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "several parts", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "several parts", mode: "parallel", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "small", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned down, workspace down", size: "empty", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "built-in", hosts: "presigned breaks", size: "several parts", mode: "stream", want: "broken", servedBy: "offload"},
	{list: "built-in", hosts: "presigned breaks", size: "several parts", mode: "to-file", want: "broken", servedBy: "offload"},
	{list: "built-in", hosts: "presigned breaks", size: "several parts", mode: "parallel", want: "connection", servedBy: "offload"},
	{list: "proxy", hosts: "", size: "several parts", mode: "stream", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "several parts", mode: "to-file", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "several parts", mode: "parallel", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "small", mode: "stream", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "", size: "empty", mode: "to-file", want: "ok", servedBy: "dp"},
	{list: "proxy", hosts: "proxy down", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy rejects", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "several parts", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "several parts", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "small", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down", size: "empty", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "several parts", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "several parts", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "several parts", mode: "parallel", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "small", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy down, presigned down, workspace down", size: "empty", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "proxy", hosts: "proxy breaks", size: "several parts", mode: "stream", want: "broken", servedBy: "dp"},
	{list: "proxy", hosts: "proxy breaks", size: "several parts", mode: "to-file", want: "broken", servedBy: "dp"},
	{list: "proxy", hosts: "proxy breaks", size: "several parts", mode: "parallel", want: "connection", servedBy: "dp"},
	{list: "rig", hosts: "", size: "several parts", mode: "stream", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "several parts", mode: "to-file", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "several parts", mode: "parallel", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "small", mode: "stream", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "", size: "empty", mode: "to-file", want: "ok", servedBy: "rig"},
	{list: "rig", hosts: "rig down", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig rejects", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig down, presigned down", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "several parts", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "several parts", mode: "parallel", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "small", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down", size: "empty", mode: "to-file", want: "ok", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "several parts", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "several parts", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "several parts", mode: "parallel", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "small", mode: "stream", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig down, presigned down, workspace down", size: "empty", mode: "to-file", want: "deadline", servedBy: "floor"},
	{list: "rig", hosts: "rig breaks", size: "several parts", mode: "stream", want: "broken", servedBy: "rig"},
	{list: "rig", hosts: "rig breaks", size: "several parts", mode: "to-file", want: "broken", servedBy: "rig"},
	{list: "rig", hosts: "rig breaks", size: "several parts", mode: "parallel", want: "connection", servedBy: "rig"},
	{list: "presigned only", hosts: "", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "several parts", mode: "parallel", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "small", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "", size: "empty", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned down", size: "several parts", mode: "stream", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned down", size: "several parts", mode: "to-file", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned down", size: "several parts", mode: "parallel", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned down", size: "small", mode: "stream", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned down", size: "empty", mode: "to-file", want: "connection", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned rejects", size: "several parts", mode: "stream", want: "permission denied", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned rejects", size: "several parts", mode: "to-file", want: "permission denied", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned rejects", size: "several parts", mode: "parallel", want: "permission denied", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned rejects", size: "small", mode: "stream", want: "permission denied", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned rejects", size: "empty", mode: "to-file", want: "permission denied", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned breaks", size: "several parts", mode: "stream", want: "broken", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned breaks", size: "several parts", mode: "to-file", want: "broken", servedBy: "offload"},
	{list: "presigned only", hosts: "presigned breaks", size: "several parts", mode: "parallel", want: "connection", servedBy: "offload"},
	// Before the first byte, a host other than the last one gets no retry, so
	// a single failure moves the download on; the workspace retries.
	{list: "built-in", hosts: "presigned fails once", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	{list: "proxy", hosts: "proxy fails once", size: "several parts", mode: "to-file", want: "ok", servedBy: "offload"},
	{list: "built-in", hosts: "presigned down, workspace fails once", size: "several parts", mode: "stream", want: "ok", servedBy: "floor"},
	// A break once bytes have arrived resumes on the same endpoint.
	{list: "built-in", hosts: "presigned breaks once", size: "several parts", mode: "stream", want: "ok", servedBy: "offload"},
	{list: "rig", hosts: "rig breaks once", size: "several parts", mode: "to-file", want: "ok", servedBy: "rig"},
}

func TestTransfer_upload(t *testing.T) {
	for _, c := range uploadCases {
		t.Run(c.name(), func(t *testing.T) {
			content := transferPayload(transferSizes[c.size])
			srv, client := newTransferHosts(t, c, nil)
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			parallelism := 1
			if c.mode == "parallel" {
				parallelism = 10
			}
			opts := []UploadOption{WithPartSize(transferPartSize), WithParallelism(parallelism), WithTransferClient(srv.Client())}
			var err error
			switch c.source {
			case "file":
				_, err = client.UploadFrom(ctx, transferPath, writeTemp(t, content), opts...)
			case "seekable":
				_, err = client.Upload(ctx, transferPath, bytes.NewReader(content), opts...)
			case "non-seekable":
				_, err = client.Upload(ctx, transferPath, struct{ io.Reader }{bytes.NewReader(content)}, opts...)
			default:
				t.Fatalf("unknown source %q", c.source)
			}
			checkTransferOutcome(t, ctx, c, err)

			stored, ok, serr := srv.StoredFile()
			switch {
			case serr != nil:
				t.Error(serr)
			case c.want == "ok" && (!ok || !bytes.Equal(stored, content)):
				t.Errorf("stored %d bytes (stored: %v), want the %d uploaded", len(stored), ok, len(content))
			case c.want != "ok" && ok:
				t.Error("a failed upload stored a file")
			}
			checkTransferHosts(t, srv, c)
		})
	}
}

func TestTransfer_download(t *testing.T) {
	for _, c := range downloadCases {
		t.Run(c.name(), func(t *testing.T) {
			content := transferPayload(transferSizes[c.size])
			srv, client := newTransferHosts(t, c, content)
			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			opts := []DownloadOption{WithDownloadTransferClient(srv.Client())}
			var got []byte
			var err error
			switch c.mode {
			case "stream":
				var resp *DownloadFileResponse
				if resp, err = client.Download(ctx, transferPath, opts...); err == nil {
					got, err = io.ReadAll(resp.Contents)
					_ = resp.Contents.Close()
				}
			case "to-file", "parallel":
				if c.mode == "parallel" {
					opts = append(opts, WithDownloadParallel(true))
				}
				dir := t.TempDir()
				dest := filepath.Join(dir, "file.bin")
				if _, err = client.DownloadTo(ctx, transferPath, dest, opts...); err == nil {
					got, err = os.ReadFile(dest)
				} else if left, _ := os.ReadDir(dir); len(left) > 0 {
					t.Errorf("the failed download left %d files behind", len(left))
				}
			default:
				t.Fatalf("unknown mode %q", c.mode)
			}
			checkTransferOutcome(t, ctx, c, err)
			if c.want == "ok" && !bytes.Equal(got, content) {
				t.Errorf("downloaded %d bytes, want the %d stored", len(got), len(content))
			}
			checkTransferHosts(t, srv, c)
		})
	}
}

// newTransferHosts returns the fake for c and a client of it, with the
// client's discovery switch set from the case's list. Until the test ends,
// Client methods build their engine with 1 KiB parts and quick retries, so
// that several parts and a spent retry budget take milliseconds rather than
// 50 MiB and minutes.
func newTransferHosts(t *testing.T, c transferCase, content []byte) (*hostsFake, *Client) {
	t.Helper()
	endpoints, ok := transferLists[c.list]
	if !ok {
		t.Fatalf("unknown list %q", c.list)
	}
	srv := newHostsFake(hostsFakeConfig{
		hostNames:       hostNames{Workspace: transferWorkspace, StorageProxy: transferProxy, RIG: transferRIG, Cloud: transferCloud},
		FilePath:        transferPath,
		LastModified:    transferLastModified,
		Endpoints:       endpoints,
		Faults:          transferFaults(t, c),
		DownloadContent: content,
	})
	t.Setenv(endpointsAPIEnv, strconv.FormatBool(endpoints != nil))
	shrinkClientTunables(t)
	return srv, buildUploadClient(t, "http://"+transferWorkspace, srv.Client(), "")
}

// shrinkClientTunables makes Client methods use small parts and quick, finite
// retries until the test ends. The tests that call it don't run in parallel.
func shrinkClientTunables(t *testing.T) {
	t.Helper()
	original := clientTunables
	t.Cleanup(func() { clientTunables = original })
	clientTunables = func() tunables {
		tun := original()
		tun.minStreamSize = transferPartSize
		tun.defaultPartSize = transferPartSize
		tun.downloadMinParallelSize = transferPartSize
		tun.downloadPartSize = transferPartSize
		fast := ops.BackoffPolicy{Initial: time.Millisecond, Maximum: time.Millisecond, Factor: 1}
		tun.cpBackoff, tun.cloudBackoff = fast, fast
		tun.cpRetryTimeout = 200 * time.Millisecond
		return tun
	}
}

func TestClientTunables_defaultIsTheProductionTunables(t *testing.T) {
	if got, want := clientTunables(), defaultTunables(); !reflect.DeepEqual(got, want) {
		t.Errorf("clientTunables() = %+v, want defaultTunables() = %+v", got, want)
	}
}

// transferHostOps names each host's endpoint in the fake and the operations a
// host that is down fails: all of them, the storage proxy's probe included.
var transferHostOps = map[string]struct {
	endpoint string
	down     []string
}{
	"presigned": {"offload", []string{"get", "part"}},
	"workspace": {"floor", []string{"get", "single"}},
	"proxy":     {"dp", []string{"abort", "complete", "get", "initiate", "part", "single"}},
	"rig":       {"rig", []string{"abort", "complete", "get", "initiate", "part", "single"}},
}

// transferFaults turns c.hosts into the fake's faults.
func transferFaults(t *testing.T, c transferCase) map[string]hostsFaultFunc {
	t.Helper()
	faults := map[string]hostsFaultFunc{}
	if c.hosts == "" {
		return faults
	}
	refused := &hostsFault{Connection: true}
	partial := &hostsFault{Break: true, BreakAfter: 1536}
	download := c.source == ""
	for h := range strings.SplitSeq(c.hosts, ", ") {
		name, state, _ := strings.Cut(h, " ")
		host, ok := transferHostOps[name]
		if !ok {
			t.Fatalf("unknown host %q", name)
		}
		set := func(fault hostsFaultFunc, actions ...string) {
			for _, a := range actions {
				faults[host.endpoint+"."+a] = fault
			}
		}
		data := "part"
		if download {
			data = "get"
		}
		switch state {
		case "down":
			set(hostsAlways(refused), host.down...)
			if name == "proxy" {
				faults["probe"] = hostsAlways(refused)
			}
		case "rejects":
			set(hostsAlways(&hostsFault{Status: 403}), "get", "part", "single")
		case "breaks":
			first := (*hostsFault)(nil)
			if download && c.mode != "parallel" {
				first = partial
			}
			set(firstThen(first, refused), data)
		case "fails once":
			set(hostsOnCalls(refused, 1), "get", "part", "single")
		case "breaks once":
			set(hostsOnCalls(partial, 1), data)
		default:
			t.Fatalf("unknown state %q of %s", state, name)
		}
	}
	return faults
}

// firstThen answers the first call with first and every later one with later.
func firstThen(first, later *hostsFault) hostsFaultFunc {
	return func(call int) *hostsFault {
		if call == 1 {
			return first
		}
		return later
	}
}

// checkTransferOutcome checks err against c.want, as a caller would tell the
// outcomes apart.
func checkTransferOutcome(t *testing.T, ctx context.Context, c transferCase, err error) {
	t.Helper()
	if ctx.Err() != nil {
		t.Fatalf("the transfer did not finish in time: %v", err)
	}
	if got := transferOutcome(err); got != c.want {
		t.Errorf("outcome = %s (%v), want %s", got, err, c.want)
	}
}

func transferOutcome(err error) string {
	var opErr *net.OpError
	switch {
	case err == nil:
		return "ok"
	case apierr.Code(err) == codes.PermissionDenied:
		return "permission denied"
	case errors.Is(err, io.ErrUnexpectedEOF):
		return "broken"
	case errors.Is(err, context.DeadlineExceeded):
		return "deadline"
	case errors.As(err, &opErr):
		return "connection"
	}
	return fmt.Sprintf("unexpected %T", err)
}

// checkTransferHosts checks which endpoint served the transfer, and the
// fake's own checks: no credentials on presigned URLs, credentials on every
// Databricks request, and nothing outside the protocol.
func checkTransferHosts(t *testing.T, srv *hostsFake, c transferCase) {
	t.Helper()
	touched := srv.EndpointsTouched()
	if len(touched) == 0 || touched[len(touched)-1] != c.servedBy {
		t.Errorf("endpoints used = %v, want the last to be %s", touched, c.servedBy)
	}
	if err := srv.CheckCredentials(); err != nil {
		t.Error(err)
	}
	if v := srv.Violations(); len(v) > 0 {
		t.Errorf("protocol violations: %v", v)
	}
	if u := srv.UnknownRequests(); len(u) > 0 {
		t.Errorf("requests the fake doesn't serve: %v", u)
	}
}

// transferPayload returns n bytes that differ from part to part.
func transferPayload(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i*31 + i/transferPartSize)
	}
	return b
}
