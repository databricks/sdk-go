package files

// Endpoint discovery: the workspace lists the endpoints a transfer should use,
// best first (GET /api/2.0/fs/endpoints).

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
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

// endpointsAPIEnv set to false turns discovery off for Client calls, where it
// is on by default. It is an experimental, unsupported setting that may be
// removed. client.WithoutEnv does not cover it, because the Client does not
// keep that choice.
const endpointsAPIEnv = "DATABRICKS_EXPERIMENTAL_FILES_EXT_ENABLE_ENDPOINTS_API"

// endpointsAPIEnabled reports whether v, the value of endpointsAPIEnv, leaves
// discovery on. A value strconv.ParseBool rejects leaves it on, with a
// warning, so that a typo does not pass for turning it off.
func endpointsAPIEnabled(ctx context.Context, logger *slog.Logger, v string) bool {
	if v == "" {
		return true
	}
	on, err := strconv.ParseBool(v)
	if err != nil {
		logger.WarnContext(ctx, "ignoring an invalid "+endpointsAPIEnv+", endpoint discovery stays on", "value", v)
		return true
	}
	return on
}

// endpointsResolveTimeout bounds discovery, which is a single attempt: it
// precedes a transfer, so it must fall back to the built-in order quickly
// rather than retry and stall the upload.
const endpointsResolveTimeout = 10 * time.Second

// How long a client keeps the workspace's answer before asking again, so that
// a change on the server reaches long-running processes. A refresh that fails
// transiently keeps the workspace's last list, but only for endpointsMaxStale
// after the workspace gave it, so discovery that keeps failing still ends on
// the built-in order.
const (
	endpointsCacheTTL        = 30 * time.Minute
	endpointsFailureCacheTTL = 30 * time.Minute
	endpointsMaxStale        = 2 * time.Hour
)

// endpointsCacheKey keys the client's discovered endpoints in its extensions.
type endpointsCacheKey struct{}

// endpointsCache holds a client's endpoint list and when to ask again.
type endpointsCache struct {
	// mu guards cached, and is held only to read or replace it.
	mu     sync.Mutex
	cached *cachedEndpoints
	// refresh is held by the one caller that asks the workspace, so that
	// concurrent transfers don't each ask.
	refresh chan struct{}
	// httpClient is what discovery sends with. It is built and used only
	// while refresh is held.
	httpClient *http.Client
}

func (c *endpointsCache) load() *cachedEndpoints {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cached
}

func (c *endpointsCache) store(cached *cachedEndpoints) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cached = cached
}

// cachedEndpoints is an endpoint list and when to ask the workspace again.
type cachedEndpoints struct {
	endpoints []endpoint
	expiresAt time.Time
	// fromServer is whether endpoints is the workspace's list rather than the
	// built-in order, and fetchedAt when the workspace gave it.
	fromServer bool
	fetchedAt  time.Time
}

// pastStaleLimit reports whether the workspace's list is too old to keep
// through refreshes that fail.
func (c *cachedEndpoints) pastStaleLimit(now time.Time) bool {
	return c.fromServer && !now.Before(c.fetchedAt.Add(endpointsMaxStale))
}

// transferEndpoints returns the endpoints a transfer walks, best first: the
// built-in order, or with discovery on, the workspace's list. The answer, a
// failure included, is kept until it expires. Then one caller asks again while
// the others keep the expired list; a client with no list yet waits for the
// first answer. Anything short of a usable list falls back to the built-in
// order, so discovery never fails a transfer.
func (e *engine) transferEndpoints(ctx context.Context) []endpoint {
	if !e.tun.endpointsAPI {
		return e.builtinEndpoints()
	}
	cache := e.endpointsCache()
	cached := cache.load()
	if cached != nil && e.clock().Before(cached.expiresAt) {
		return slices.Clone(cached.endpoints)
	}
	if cached == nil {
		// Nothing to use yet: wait for whichever caller is asking.
		select {
		case cache.refresh <- struct{}{}:
		case <-ctx.Done():
			return e.builtinEndpoints()
		}
	} else {
		select {
		case cache.refresh <- struct{}{}:
		default:
			// Another caller is asking. Keep transferring on the expired list
			// rather than wait, unless it is past the stale limit, which holds
			// however long that caller takes.
			if cached.pastStaleLimit(e.clock()) {
				return e.builtinEndpoints()
			}
			return slices.Clone(cached.endpoints)
		}
	}
	defer func() { <-cache.refresh }()
	if cached = cache.load(); cached != nil && e.clock().Before(cached.expiresAt) {
		// Another caller asked while this one waited.
		return slices.Clone(cached.endpoints)
	}
	if cache.httpClient == nil {
		cache.httpClient = e.discoveryClient()
	}
	fetched, transient := e.fetchEndpoints(ctx, cache.httpClient)
	if ctx.Err() != nil {
		// The caller gave up; that is not the workspace's answer.
		return e.builtinEndpoints()
	}
	now := e.clock()
	switch {
	case transient && cached != nil && cached.fromServer && !cached.pastStaleLimit(now):
		// A blip in discovery must not move a working client off the
		// workspace's list. A definitive answer, such as discovery turned off,
		// still does.
		expiresAt := now.Add(endpointsFailureCacheTTL)
		if limit := cached.fetchedAt.Add(endpointsMaxStale); limit.Before(expiresAt) {
			expiresAt = limit
		}
		cached = &cachedEndpoints{endpoints: cached.endpoints, expiresAt: expiresAt, fromServer: true, fetchedAt: cached.fetchedAt}
	case len(fetched) > 0:
		names := make([]string, len(fetched))
		for i, ep := range fetched {
			names[i] = ep.String()
		}
		e.c.logger.DebugContext(ctx, "using the workspace's file endpoints", "endpoints", names)
		cached = &cachedEndpoints{endpoints: fetched, expiresAt: now.Add(endpointsCacheTTL), fromServer: true, fetchedAt: now}
	default:
		cached = &cachedEndpoints{endpoints: e.builtinEndpoints(), expiresAt: now.Add(endpointsFailureCacheTTL), fetchedAt: now}
	}
	cache.store(cached)
	return slices.Clone(cached.endpoints)
}

// endpointsCache returns the client's endpoint cache.
func (e *engine) endpointsCache() *endpointsCache {
	if v, ok := e.c.extensions.Load(endpointsCacheKey{}); ok {
		return v.(*endpointsCache)
	}
	v, _ := e.c.extensions.LoadOrStore(endpointsCacheKey{}, &endpointsCache{refresh: make(chan struct{}, 1)})
	return v.(*endpointsCache)
}

// discoveryClient returns the client discovery sends with: the client's, on a
// copy of its transport. A refresh runs while transfers may hold every
// connection in that transport's pool, and must not wait for them before its
// timeout even starts. The copy keeps no idle connections: it sends one request
// every half hour, and nothing closes it. A transport other than
// *http.Transport can't be copied, so it is shared.
func (e *engine) discoveryClient() *http.Client {
	hc := e.c.httpClient
	if hc == nil {
		hc = &http.Client{}
	}
	base := hc.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	t, ok := base.(*http.Transport)
	if !ok {
		return hc
	}
	t = t.Clone()
	t.DisableKeepAlives = true
	c := *hc
	c.Transport = t
	return &c
}

// fetchEndpoints asks the workspace for its endpoint list and keeps the
// entries this package can use. It returns nil if there are none, and whether
// that may be transient: no response, a 408, 429 or 5xx, or a body that does
// not parse. Any other answer is the workspace's: a 403 or 404 (discovery
// turned off or not served), another 4xx, or nothing usable in the list.
func (e *engine) fetchEndpoints(ctx context.Context, httpClient *http.Client) (endpoints []endpoint, transient bool) {
	ctx, cancel := context.WithTimeout(ctx, endpointsResolveTimeout)
	defer cancel()
	urlStr, err := endpointURL(e.c.host, endpointsAPIPath, nil)
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil, true
	}
	headers := http.Header{}
	headers.Set("Accept", "application/json")
	req, err := newHTTPRequest(ctx, httpRequestOptions{Method: http.MethodGet, URL: urlStr, Credentials: e.c.credentials, UserAgent: e.userAgent, Headers: headers})
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil, true
	}
	e.setWorkspaceHeader(req)
	resp, err := executeStreamingHTTPCall(httpCallOptions{req: req, client: httpClient, logger: e.c.logger})
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok {
		status, requestID := aerr.HTTPStatusCode(), aerr.HTTPHeader().Get("X-Request-Id")
		if status == http.StatusNotFound {
			// Not a warning: some workspaces don't serve the API yet.
			e.c.logger.InfoContext(ctx, "the endpoints API is not served for this workspace, using the default endpoints", "request_id", requestID)
			return nil, false
		}
		e.c.logger.WarnContext(ctx, "unexpected status resolving file endpoints, using the default endpoints",
			"status", status, "request_id", requestID)
		return nil, status >= 500 || status == http.StatusRequestTimeout || status == http.StatusTooManyRequests
	}
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to resolve file endpoints, using the default endpoints", "error", err)
		return nil, true
	}
	defer resp.Body.Close()
	requestID := resp.Header.Get("X-Request-Id")
	if resp.StatusCode != http.StatusOK {
		// Any other success or redirect is not a list: it's the workspace's
		// answer all the same, so it isn't retried sooner.
		e.c.logger.WarnContext(ctx, "unexpected status resolving file endpoints, using the default endpoints",
			"status", resp.StatusCode, "request_id", requestID)
		return nil, false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.c.logger.WarnContext(ctx, "failed to read the endpoints response, using the default endpoints", "request_id", requestID, "error", err)
		return nil, true
	}
	e.c.logger.DebugContext(ctx, "endpoints response", "request_id", requestID, "body", loggedBody(body))
	var doc any
	if err := json.Unmarshal(body, &doc); err != nil {
		e.c.logger.WarnContext(ctx, "failed to parse the endpoints response, using the default endpoints", "request_id", requestID, "error", err)
		return nil, true
	}
	obj, isObject := doc.(map[string]any)
	// Proto3 JSON leaves out an empty repeated field.
	entries, isList := obj["endpoints"].([]any)
	if _, present := obj["endpoints"]; !present {
		isList = true
	}
	if !isObject || !isList {
		e.c.logger.WarnContext(ctx, "unexpected endpoints response, using the default endpoints", "request_id", requestID, "body", loggedBody(body))
		return nil, false
	}
	for _, raw := range entries {
		entry, _ := raw.(map[string]any)
		kind, ok := endpointKindOf(entry)
		if !ok {
			// Not a warning: an older SDK skips the APIs added after it shipped.
			e.c.logger.InfoContext(ctx, "skipping an endpoint this SDK cannot use", "request_id", requestID, "entry", loggedEntry(raw))
			continue
		}
		ep, ok := e.parseEndpoint(entry, kind)
		if !ok {
			e.c.logger.WarnContext(ctx, "skipping an endpoint with an unusable address", "request_id", requestID, "entry", loggedEntry(raw))
			continue
		}
		endpoints = append(endpoints, ep)
	}
	if len(endpoints) == 0 {
		e.c.logger.WarnContext(ctx, "no usable endpoints in the response, using the default endpoints", "request_id", requestID)
	}
	return endpoints, false
}

// maxLoggedBody caps how much of a response discovery logs.
const maxLoggedBody = 1 << 10

// loggedBody returns body for a log line, cut to maxLoggedBody.
func loggedBody(body []byte) string {
	if len(body) > maxLoggedBody {
		return string(body[:maxLoggedBody]) + "..."
	}
	return string(body)
}

// loggedEntry returns a listed entry for a log line.
func loggedEntry(raw any) string {
	b, _ := json.Marshal(raw)
	return loggedBody(b)
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
