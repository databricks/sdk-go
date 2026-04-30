package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/databricks/transport"
	"github.com/databricks/sdk-go/options/call"
	"github.com/databricks/sdk-go/options/client"
	"github.com/databricks/sdk-go/options/internaloptions"
)

type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
	host       string
}

func NewClient(ctx context.Context, opts ...client.Option) (*Client, error) {
	cfg := internaloptions.ClientOptions{}
	for _, opt := range opts {
		if err := opt.Apply(&cfg); err != nil {
			return nil, err
		}
	}
	httpClient, err := transport.NewHTTPClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: httpClient,
		logger:     cfg.Logger,
		host:       cfg.Host,
	}, nil
}

// executeCall resolves call.Option values to ops.Option values and invokes
// ops.Execute. Lives at the package level to keep call sites concise.
func executeCall(ctx context.Context, op func(context.Context) error, opts []call.Option) error {
	cfg := internaloptions.CallOptions{}
	for _, opt := range opts {
		if err := opt.Apply(&cfg); err != nil {
			return err
		}
	}
	var opsOpts []ops.Option
	if cfg.Retrier != nil {
		opsOpts = append(opsOpts, ops.WithRetrier(cfg.Retrier))
	}
	if cfg.RateLimiter != nil {
		opsOpts = append(opsOpts, ops.WithLimiter(cfg.RateLimiter))
	}
	if cfg.Timeout != 0 {
		opsOpts = append(opsOpts, ops.WithTimeout(cfg.Timeout))
	}
	return ops.Execute(ctx, op, opsOpts...)
}

// Cancels a data quality monitor refresh. Currently only supported for the
// `table` `object_type`. The call must be made in the same workspace as where
// the monitor was created.
//
// The caller must have either of the following sets of permissions: 1.
// **MANAGE** and **USE_CATALOG** on the table's parent catalog. 2.
// **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **MANAGE** on the table.
func (c *Client) CancelRefresh(ctx context.Context, req *CancelRefreshRequest, opts ...call.Option) (*CancelRefreshResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes/%v/cancel", *req.ObjectType, *req.ObjectId, *req.RefreshId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelRefreshResponse{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Create a data quality monitor on a Unity Catalog object. The caller must
// provide either `anomaly_detection_config` for a schema monitor or
// `data_profiling_config` for a table monitor.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the table's parent
// catalog, **USE_SCHEMA** on the table's parent schema, and **SELECT** on the
// table 2. **USE_CATALOG** on the table's parent catalog, **MANAGE** and
// **USE_SCHEMA** on the table's parent schema, and **SELECT** on the table. 3.
// **USE_CATALOG** on the table's parent catalog, **USE_SCHEMA** on the table's
// parent schema, and **MANAGE** and **SELECT** on the table.
//
// Workspace assets, such as the dashboard, will be created in the workspace
// where this call was made.
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and **MANAGE**
// and **USE_SCHEMA** on the schema.
func (c *Client) CreateMonitor(ctx context.Context, req *CreateMonitorRequest, opts ...call.Option) (*Monitor, error) {
	body, err := json.Marshal(req.Monitor)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/data-quality/v1/monitors"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Monitor{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Creates a refresh. Currently only supported for the `table` `object_type`.
// The call must be made in the same workspace as where the monitor was created.
//
// The caller must have either of the following sets of permissions: 1.
// **MANAGE** and **USE_CATALOG** on the table's parent catalog. 2.
// **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **MANAGE** on the table.
func (c *Client) CreateRefresh(ctx context.Context, req *CreateRefreshRequest, opts ...call.Option) (*Refresh, error) {
	body, err := json.Marshal(req.Refresh)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes", *req.Refresh.ObjectType, *req.Refresh.ObjectId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Refresh{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("POST", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Delete a data quality monitor on Unity Catalog object.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: **MANAGE** and **USE_CATALOG** on the table's parent
// catalog. **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. **USE_CATALOG** on the table's
// parent catalog, **USE_SCHEMA** on the table's parent schema, and **MANAGE**
// on the table.
//
// Note that the metric tables and dashboard will not be deleted as part of this
// call; those assets must be manually cleaned up (if desired).
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and **MANAGE**
// and **USE_SCHEMA** on the schema.
func (c *Client) DeleteMonitor(ctx context.Context, req *DeleteMonitorRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("DELETE", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		_ = respBody
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// (Unimplemented) Delete a refresh
func (c *Client) DeleteRefresh(ctx context.Context, req *DeleteRefreshRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes/%v", *req.ObjectType, *req.ObjectId, *req.RefreshId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("DELETE", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		_ = respBody
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// Read a data quality monitor on a Unity Catalog object.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the table's parent
// catalog. 2. **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **SELECT** on the table.
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and
// **USE_SCHEMA** on the schema.
//
// The returned information includes configuration values on the entity and
// parent entity as well as information on assets created by the monitor. Some
// information (e.g. dashboard) may be filtered out if the caller is in a
// different workspace than where the monitor was created.
func (c *Client) GetMonitor(ctx context.Context, req *GetMonitorRequest, opts ...call.Option) (*Monitor, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Monitor{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get data quality monitor refresh. The call must be made in the same workspace
// as where the monitor was created.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the table's parent
// catalog. 2. **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **SELECT** on the table.
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and
// **USE_SCHEMA** on the schema.
func (c *Client) GetRefresh(ctx context.Context, req *GetRefreshRequest, opts ...call.Option) (*Refresh, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes/%v", *req.ObjectType, *req.ObjectId, *req.RefreshId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Refresh{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// (Unimplemented) List data quality monitors.
func (c *Client) ListMonitor(ctx context.Context, req *ListMonitorRequest, opts ...call.Option) (*ListMonitorResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/data-quality/v1/monitors"
	queryParams := url.Values{}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListMonitorResponse{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListMonitorIter(ctx context.Context, req *ListMonitorRequest, opts ...call.Option) iter.Seq2[*Monitor, error] {
	return func(yield func(*Monitor, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListMonitorRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListMonitor(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Monitors {
				if !yield(&resp.Monitors[i], nil) {
					return
				}
			}
			if resp.NextPageToken == nil || *resp.NextPageToken == "" {
				return
			}
			pageReq.PageToken = resp.NextPageToken
		}
	}
}

// List data quality monitor refreshes. The call must be made in the same
// workspace as where the monitor was created.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the table's parent
// catalog. 2. **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **SELECT** on the table.
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and
// **USE_SCHEMA** on the schema.
func (c *Client) ListRefresh(ctx context.Context, req *ListRefreshRequest, opts ...call.Option) (*ListRefreshResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListRefreshResponse{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("GET", baseURL.String(), nil)
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListRefreshIter(ctx context.Context, req *ListRefreshRequest, opts ...call.Option) iter.Seq2[*Refresh, error] {
	return func(yield func(*Refresh, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListRefreshRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListRefresh(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Refreshes {
				if !yield(&resp.Refreshes[i], nil) {
					return
				}
			}
			if resp.NextPageToken == nil || *resp.NextPageToken == "" {
				return
			}
			pageReq.PageToken = resp.NextPageToken
		}
	}
}

// Update a data quality monitor on Unity Catalog object.
//
// For the `table` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the table's parent
// catalog. 2. **USE_CATALOG** on the table's parent catalog, and **MANAGE** and
// **USE_SCHEMA** on the table's parent schema. 3. **USE_CATALOG** on the
// table's parent catalog, **USE_SCHEMA** on the table's parent schema, and
// **MANAGE** on the table.
//
// For the `schema` `object_type`, the caller must have either of the following
// sets of permissions: 1. **MANAGE** and **USE_CATALOG** on the schema's parent
// catalog. 2. **USE_CATALOG** on the schema's parent catalog, and **MANAGE**
// and **USE_SCHEMA** on the schema.
func (c *Client) UpdateMonitor(ctx context.Context, req *UpdateMonitorRequest, opts ...call.Option) (*Monitor, error) {
	body, err := json.Marshal(req.Monitor)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Monitor{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("PATCH", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// (Unimplemented) Update a refresh
func (c *Client) UpdateRefresh(ctx context.Context, req *UpdateRefreshRequest, opts ...call.Option) (*Refresh, error) {
	body, err := json.Marshal(req.Refresh)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/data-quality/v1/monitors/%v/%v/refreshes/%v", *req.ObjectType, *req.ObjectId, *req.RefreshId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Refresh{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("PATCH", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		httpReq = httpReq.WithContext(ctx)
		httpReq.Header = headers

		respBody, err := executeHTTPCall(httpCallOptions{
			req:    httpReq,
			client: c.httpClient,
			logger: c.logger,
		})
		if err != nil {
			return err
		}
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}
