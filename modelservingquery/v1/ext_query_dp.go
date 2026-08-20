package modelservingquery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/ops"
)

// dpStateKey is the private key under which the route-optimization state is
// stored in the client's extension map.
type dpStateKey struct{}

// dpState holds the route-optimization state for a Client, built once per
// client and cached in the extension map.
type dpState struct {
	client *Client

	// cpTokens is the control-plane token provider used to mint data-plane
	// tokens. It is nil when the credentials cannot mint OAuth tokens (for
	// example a personal access token).
	cpTokens auth.TokenProvider

	// endpoints caches the per-endpoint route-optimization state keyed by
	// endpoint name (string -> *endpointState).
	//
	// An entry is evicted on a data-plane call failure so neither a stale
	// endpoint URL / authorization detail nor a bad-but-unexpired data-plane
	// token can wedge an endpoint permanently.
	endpoints sync.Map
}

// endpointState bundles a route-optimized endpoint's discovered info with the
// credentials scoped to it. The two are cached and evicted together: a
// data-plane failure drops the possibly bad token alongside the metadata, so
// the next call re-discovers and re-mints instead of replaying a stale token.
type endpointState struct {
	info  *dataPlaneInfo
	creds auth.Credentials
}

// dpState returns the route-optimization state for the client, building it once
// and caching it in the extension map.
func (c *Client) dpState() *dpState {
	if v, ok := c.extensions.Load(dpStateKey{}); ok {
		return v.(*dpState)
	}
	s := &dpState{client: c}
	if tp, ok := c.credentials.(auth.TokenProvider); ok {
		s.cpTokens = auth.NewCachedTokenProvider(tp)
	}
	actual, _ := c.extensions.LoadOrStore(dpStateKey{}, s)
	return actual.(*dpState)
}

// dataPlaneInfo is the minimal projection of a serving endpoint's data-plane
// query info needed to route a query directly to the data plane.
type dataPlaneInfo struct {
	endpointURL          string
	authorizationDetails string
}

// ErrRouteOptimizationUnavailable is returned by QueryOptimized when the query
// cannot be routed to the data plane: the client is not configured with
// OAuth-capable credentials, the request has no endpoint name, or the endpoint
// does not advertise data-plane query info. A caller that wants best-effort
// behavior can detect this with errors.Is and fall back to Query:
//
//	resp, err := c.QueryOptimized(ctx, req)
//	if errors.Is(err, ErrRouteOptimizationUnavailable) {
//	    resp, err = c.Query(ctx, req)
//	}
var ErrRouteOptimizationUnavailable = errors.New("modelservingquery: route optimization unavailable for this endpoint")

// QueryOptimized queries a serving endpoint directly on the data plane,
// bypassing the control plane for lower latency.
//
// It requires OAuth-capable credentials and an endpoint that advertises
// data-plane query info. When either is missing, or the request has no endpoint
// name, it returns an error wrapping ErrRouteOptimizationUnavailable; callers
// that want to fall back can test for it with errors.Is and call Query.
//
// Unlike Query, it never falls back to the control plane once the data-plane
// call is made: an error from that call is returned as is, so a billed inference
// is not silently retried elsewhere. It otherwise behaves like Query.
func (c *Client) QueryOptimized(ctx context.Context, req *QueryEndpointRequest, opts ...ops.Option) (*QueryEndpointResponse, error) {
	dp := c.dpState()
	if dp.cpTokens == nil || req.Name == nil {
		return nil, ErrRouteOptimizationUnavailable
	}
	name := *req.Name

	ep, err := dp.endpointState(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	if ep == nil {
		// The endpoint is not route-optimized.
		return nil, fmt.Errorf("%w: endpoint %q", ErrRouteOptimizationUnavailable, name)
	}

	resp, err := dp.query(ctx, req, ep, opts...)
	if err != nil {
		// Evict info and token source together so neither a stale URL nor a
		// bad-but-unexpired token can wedge the endpoint.
		dp.endpoints.Delete(name)
		return nil, err
	}
	return resp, nil
}

// endpointState returns the cached route-optimization state for the endpoint,
// discovering it via the control plane on a cache miss. It returns (nil, nil)
// when the endpoint is not route-optimized (no data-plane query info). Negative
// results are not cached. Each freshly discovered endpoint gets its own token
// source so evicting the endpoint also discards its cached data-plane token.
func (dp *dpState) endpointState(ctx context.Context, name string, opts ...ops.Option) (*endpointState, error) {
	if v, ok := dp.endpoints.Load(name); ok {
		return v.(*endpointState), nil
	}

	info, err := dp.fetchDataPlaneInfo(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	ep := &endpointState{
		info:  info,
		creds: dp.dataPlaneCredentials(info.authorizationDetails),
	}
	actual, _ := dp.endpoints.LoadOrStore(name, ep)
	return actual.(*endpointState), nil
}

// servingEndpointDetailedWire is the minimal projection of the serving-endpoint
// GET response. Only the data-plane query info is decoded; everything else on
// the endpoint is ignored. This mirrors the wire shape of
// modelserving.ModelDataPlaneInfo without depending on that module, keeping
// modelservingquery self-contained.
type servingEndpointDetailedWire struct {
	DataPlaneInfo *struct {
		QueryInfo *struct {
			EndpointURL          *string `json:"endpoint_url,omitempty"`
			AuthorizationDetails *string `json:"authorization_details,omitempty"`
		} `json:"query_info,omitempty"`
	} `json:"data_plane_info,omitempty"`
}

// fetchDataPlaneInfo issues the control-plane GET that reveals whether the
// endpoint is route-optimized. It returns nil (no error) when the endpoint has
// no data-plane query info. The GET is signed with the client's control-plane
// credentials.
func (dp *dpState) fetchDataPlaneInfo(ctx context.Context, name string, opts ...ops.Option) (*dataPlaneInfo, error) {
	c := dp.client
	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	pb := pathBuilder{}
	pb.literal("/api/2.0/serving-endpoints/")
	pb.singleSegment(name)
	baseURL.Path, baseURL.RawPath = pb.build()
	urlStr := baseURL.String()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	var info *dataPlaneInfo
	call := func(ctx context.Context) error {
		httpReq, err := newHTTPRequest(ctx, httpRequestOptions{
			Method:      http.MethodGet,
			URL:         urlStr,
			Credentials: c.credentials,
			Headers:     headers,
		})
		if err != nil {
			return err
		}

		respBody, _, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		var wireResp servingEndpointDetailedWire
		if err := json.Unmarshal(respBody, &wireResp); err != nil {
			return err
		}
		info = dataPlaneInfoFromWire(&wireResp)
		return nil
	}

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return info, nil
}

// dataPlaneInfoFromWire projects the endpoint response to a dataPlaneInfo,
// returning nil when the endpoint does not advertise a data-plane endpoint URL.
func dataPlaneInfoFromWire(w *servingEndpointDetailedWire) *dataPlaneInfo {
	if w.DataPlaneInfo == nil || w.DataPlaneInfo.QueryInfo == nil {
		return nil
	}
	qi := w.DataPlaneInfo.QueryInfo
	if qi.EndpointURL == nil || *qi.EndpointURL == "" {
		return nil
	}
	info := &dataPlaneInfo{endpointURL: *qi.EndpointURL}
	if qi.AuthorizationDetails != nil {
		info.authorizationDetails = *qi.AuthorizationDetails
	}
	return info
}

// query posts the request directly to the data-plane endpoint URL, using the
// endpoint's data-plane credentials. It mirrors the generated control-plane
// Query (body/response wire conversion, served-model-name header) but targets
// the absolute data-plane URL and signs with the data-plane token instead of
// the control-plane credentials.
func (dp *dpState) query(ctx context.Context, req *QueryEndpointRequest, ep *endpointState, opts ...ops.Option) (*QueryEndpointResponse, error) {
	c := dp.client
	info := ep.info
	wireReq, err := queryEndpointRequestToWire(req)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(wireReq)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	if c.workspaceID != "" {
		headers.Set("X-Databricks-Workspace-Id", c.workspaceID)
	}

	var resp *QueryEndpointResponse
	call := func(ctx context.Context) error {
		// The token is minted per attempt via ep.creds.AuthHeaders, so a
		// refreshed token is used on retry.
		httpReq, err := newHTTPRequest(ctx, httpRequestOptions{
			Method:      http.MethodPost,
			URL:         info.endpointURL,
			Credentials: ep.creds,
			Headers:     headers,
			Body:        bytes.NewBuffer(body),
		})
		if err != nil {
			return err
		}

		respBody, respHeader, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		var wireResp queryEndpointResponseWire
		if err := json.Unmarshal(respBody, &wireResp); err != nil {
			return err
		}
		resp, err = queryEndpointResponseFromWire(&wireResp)
		if err != nil {
			return err
		}
		if v := respHeader.Get("served-model-name"); v != "" {
			h := v
			resp.ServedModelName = &h
		}
		return nil
	}

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// dataPlaneCredentials returns credentials that authenticate a request with a
// data-plane token minted for the given authorization details by exchanging a
// fresh control-plane token. The underlying token is cached and refreshed
// asynchronously before it expires, with concurrent refreshes coalesced.
func (dp *dpState) dataPlaneCredentials(authDetails string) auth.Credentials {
	tokens := auth.NewCachedTokenProvider(auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		cpToken, err := dp.cpTokens.Token(ctx)
		if err != nil {
			return nil, err
		}
		return dp.exchangeToken(ctx, authDetails, cpToken)
	}))
	return auth.NewTokenCredentials("dataplane", tokens)
}

const jwtBearerGrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"

// exchangeToken swaps a control-plane token for a data-plane-scoped token via
// the workspace OIDC token endpoint, using the JWT-bearer grant with RFC 9396
// authorization details. The request carries the assertion in its form body and
// is issued with nil credentials, so no control-plane signing is applied.
func (dp *dpState) exchangeToken(ctx context.Context, authDetails string, cpToken *auth.Token) (*auth.Token, error) {
	c := dp.client
	tokenURL := strings.TrimRight(c.host, "/") + "/oidc/v1/token"

	form := url.Values{}
	form.Set("grant_type", jwtBearerGrantType)
	form.Set("authorization_details", authDetails)
	form.Set("assertion", cpToken.Value)

	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Accept", "application/json")

	// Credentials is nil: the assertion in the form body is the credential, so
	// the request must not also carry control-plane signing.
	httpReq, err := newHTTPRequest(ctx, httpRequestOptions{
		Method:  http.MethodPost,
		URL:     tokenURL,
		Headers: headers,
		Body:    strings.NewReader(form.Encode()),
	})
	if err != nil {
		return nil, err
	}

	respBody, _, err := executeHTTPCall(httpCallOptions{
		req:    httpReq,
		client: c.httpClient,
		logger: c.logger,
	})
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("oidc token exchange: missing access_token (expires_in=%d)", tokenResp.ExpiresIn)
	}

	token := &auth.Token{
		Value: tokenResp.AccessToken,
		Type:  tokenResp.TokenType,
	}
	if tokenResp.ExpiresIn > 0 {
		token.Expiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	return token, nil
}
