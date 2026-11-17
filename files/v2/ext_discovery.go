package files

// Endpoint discovery: the workspace lists the endpoints a transfer should use,
// best first (GET /api/2.0/fs/endpoints).

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
)

const endpointsAPIPath = "/api/2.0/fs/endpoints"

// endpointsResolveTimeout bounds discovery, which is a single attempt: it
// precedes a transfer, so it must fall back to the built-in order quickly
// rather than retry and stall the upload.
const endpointsResolveTimeout = 10 * time.Second

// endpointsCacheKey keys the client's discovered endpoints in its extensions.
type endpointsCacheKey struct{}

// endpointsCache holds the workspace's answer, a failure included, so a client
// asks once.
type endpointsCache struct {
	mu        sync.Mutex
	resolved  bool
	endpoints []endpoint
}

// transferEndpoints returns the endpoints a transfer walks, best first: the
// built-in order, or with discovery on, the workspace's list. Anything short of
// a usable list falls back to the built-in order, so discovery never fails a
// transfer.
func (e *engine) transferEndpoints(ctx context.Context) []endpoint {
	if !e.tun.endpointsAPI {
		return e.builtinEndpoints()
	}
	v, _ := e.c.extensions.LoadOrStore(endpointsCacheKey{}, &endpointsCache{})
	cache := v.(*endpointsCache)
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if !cache.resolved {
		endpoints := e.fetchEndpoints(ctx)
		if ctx.Err() != nil {
			// The caller gave up; that is not the workspace's answer.
			return e.builtinEndpoints()
		}
		if len(endpoints) == 0 {
			endpoints = e.builtinEndpoints()
		}
		cache.resolved, cache.endpoints = true, endpoints
	}
	return slices.Clone(cache.endpoints)
}

// fetchEndpoints asks the workspace for its endpoint list and keeps the
// entries this package can use. It returns nil if there are none.
func (e *engine) fetchEndpoints(ctx context.Context) []endpoint {
	ctx, cancel := context.WithTimeout(ctx, endpointsResolveTimeout)
	defer cancel()
	urlStr, err := endpointURL(e.c.host, endpointsAPIPath, nil)
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil
	}
	headers := http.Header{}
	headers.Set("Accept", "application/json")
	req, err := newHTTPRequest(ctx, httpRequestOptions{Method: http.MethodGet, URL: urlStr, Credentials: e.c.credentials, Headers: headers})
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil
	}
	e.setWorkspaceHeader(req)
	body, respHeader, err := executeHTTPCall(httpCallOptions{req: req, client: e.c.httpClient, logger: e.c.logger})
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok {
		requestID := aerr.HTTPHeader().Get("X-Request-Id")
		if aerr.HTTPStatusCode() == http.StatusNotFound {
			// Not a warning: some workspaces don't serve the API yet.
			e.c.logger.InfoContext(ctx, "the endpoints API is not served for this workspace, using the default endpoints", "request_id", requestID)
			return nil
		}
		e.c.logger.WarnContext(ctx, "unexpected status resolving file endpoints, using the default endpoints",
			"status", aerr.HTTPStatusCode(), "request_id", requestID)
		return nil
	}
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil
	}
	requestID := respHeader.Get("X-Request-Id")
	var resp struct {
		// Proto3 JSON leaves out an empty repeated field.
		Endpoints []json.RawMessage `json:"endpoints"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		e.c.logger.WarnContext(ctx, "unexpected endpoints response, using the default endpoints", "request_id", requestID, "error", err)
		return nil
	}
	var endpoints []endpoint
	for _, raw := range resp.Endpoints {
		var entry map[string]any
		_ = json.Unmarshal(raw, &entry)
		kind, ok := endpointKindOf(entry)
		if !ok {
			// Not a warning: an older SDK skips the APIs added after it shipped.
			e.c.logger.InfoContext(ctx, "skipping an endpoint this SDK cannot use", "request_id", requestID, "entry", string(raw))
			continue
		}
		ep, ok := e.parseEndpoint(entry, kind)
		if !ok {
			e.c.logger.WarnContext(ctx, "skipping an endpoint with an unusable address", "request_id", requestID, "entry", string(raw))
			continue
		}
		endpoints = append(endpoints, ep)
	}
	if len(endpoints) == 0 {
		e.c.logger.WarnContext(ctx, "no usable endpoints in the response, using the default endpoints", "request_id", requestID)
	}
	return endpoints
}

// endpointKindOf returns the protocol of a listed endpoint by its api and
// api_version, or false if this package doesn't speak it. Pairs not listed
// here are skipped, which lets the server list newer protocols without
// breaking older SDKs. IO_OFFLOAD is unversioned during the preview, so a
// versioned IO_OFFLOAD is newer than this package.
func endpointKindOf(entry map[string]any) (endpointKind, bool) {
	api, ok := entry["api"].(string)
	if !ok {
		return 0, false
	}
	version, hasVersion := entry["api_version"]
	if version == nil {
		hasVersion = false
	}
	versionStr, isString := version.(string)
	if hasVersion && !isString {
		return 0, false
	}
	switch {
	case api == "IO_OFFLOAD" && !hasVersion:
		return presignedEndpoint, true
	case api == "FILES_API" && versionStr == "2.0":
		return filesAPIEndpoint, true
	case api == "FILES_API" && versionStr == "3.0":
		return filesAPIv3Endpoint, true
	}
	return 0, false
}

// parseEndpoint returns the endpoint of a listed entry, or false if its address
// is unusable: anything but a bare http or https host and port, which requests
// go to as written.
func (e *engine) parseEndpoint(entry map[string]any, kind endpointKind) (endpoint, bool) {
	scheme, _ := entry["scheme"].(string)
	hostname, _ := entry["hostname"].(string)
	if (scheme != "http" && scheme != "https") || hostname == "" {
		return endpoint{}, false
	}
	if kind == filesAPIv3Endpoint && scheme != "https" {
		return endpoint{}, false
	}
	base := scheme + "://" + hostname
	u, err := url.Parse(base)
	if err != nil || u.Host != hostname || u.User != nil || strings.HasSuffix(u.Host, ":") {
		return endpoint{}, false
	}
	if p := u.Port(); p != "" {
		if n, err := strconv.Atoi(p); err != nil || n < 1 || n > 65535 {
			return endpoint{}, false
		}
	}
	// Plain HTTP carries workspace credentials in the clear. The workspace host
	// may use it only if the client reaches the workspace that way; another host
	// only for Files API 2.0, which is the in-cluster storage proxy.
	if scheme == "http" {
		workspace, err := originOf(e.c.host)
		if err != nil {
			return endpoint{}, false
		}
		if strings.EqualFold(u.Hostname(), workspace.host) {
			if workspace.scheme != "http" {
				return endpoint{}, false
			}
		} else if kind != filesAPIEndpoint {
			return endpoint{}, false
		}
	}
	// As the proto getter reads it, an unset bool is false.
	multipart, _ := entry["multipart_upload_supported"].(bool)
	return endpoint{kind: kind, baseURL: base, multipart: multipart}, true
}
