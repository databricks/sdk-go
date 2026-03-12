package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/databricks/sdk-go/databricks/api"
	"github.com/databricks/sdk-go/databricks/options"
	"github.com/databricks/sdk-go/databricks/options/unstable"
	"github.com/databricks/sdk-go/databricks/transport"
)

type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
	host       string
}

func NewClient(ctx context.Context, opts ...options.ClientOption) (*Client, error) {
	resolved, err := unstable.Resolve(opts...)
	if err != nil {
		return nil, err
	}
	httpClient, err := transport.NewHTTPClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: httpClient,
		logger:     resolved.Logger,
		host:       resolved.Host,
	}, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Cancels an already-initiated refresh job.
func (c *Client) CancelRefresh(ctx context.Context, req *CancelRefresh, opts ...api.Option) (*CancelRefresh_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor/refreshes/%v/cancel", *req.FullTableNameArg, *req.RefreshId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelRefresh_Response{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Creates a new monitor for the specified
// table.
//
// The caller must either: 1. be an owner of the table's parent catalog, have
// **USE_SCHEMA** on the table's parent schema, and have **SELECT** access on
// the table 2. have **USE_CATALOG** on the table's parent catalog, be an owner
// of the table's parent schema, and have **SELECT** access on the table. 3.
// have the following permissions: - **USE_CATALOG** on the table's parent
// catalog - **USE_SCHEMA** on the table's parent schema - be an owner of the
// table.
//
// Workspace assets, such as the dashboard, will be created in the workspace
// where this call was made.
func (c *Client) CreateMonitor(ctx context.Context, req *CreateMonitor, opts ...api.Option) (*DataMonitorInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DataMonitorInfo{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Deletes a monitor for the specified table.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema - be an
// owner of the table.
//
// Additionally, the call must be made from the workspace where the monitor was
// created.
//
// Note that the metric tables and dashboard will not be deleted as part of this
// call; those assets must be manually cleaned up (if desired).
func (c *Client) DeleteMonitor(ctx context.Context, req *DeleteMonitor, opts ...api.Option) (*DeleteMonitor_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteMonitor_Response{}

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
		if err := json.Unmarshal(respBody, resp); err != nil {
			return err
		}
		return nil
	}

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Gets a monitor for the specified table.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema. 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema -
// **SELECT** privilege on the table.
//
// The returned information includes configuration values, as well as
// information on assets created by the monitor. Some information (e.g.,
// dashboard) may be filtered out if the caller is in a different workspace than
// where the monitor was created.
func (c *Client) GetMonitor(ctx context.Context, req *GetMonitor, opts ...api.Option) (*DataMonitorInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DataMonitorInfo{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Gets info about a specific monitor refresh
// using the given refresh ID.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema -
// **SELECT** privilege on the table.
//
// Additionally, the call must be made from the workspace where the monitor was
// created.
func (c *Client) GetRefresh(ctx context.Context, req *GetRefresh, opts ...api.Option) (*RefreshInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor/refreshes/%v", *req.FullTableNameArg, *req.RefreshId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RefreshInfo{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Gets an array containing the history of the
// most recent refreshes (up to 25) for this table.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema -
// **SELECT** privilege on the table.
//
// Additionally, the call must be made from the workspace where the monitor was
// created.
func (c *Client) ListRefreshes(ctx context.Context, req *ListRefreshes, opts ...api.Option) (*ListRefreshes_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor/refreshes", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListRefreshes_Response{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Regenerates the monitoring dashboard for the
// specified table.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema - be an
// owner of the table
//
// The call must be made from the workspace where the monitor was created. The
// dashboard will be regenerated in the assets directory that was specified when
// the monitor was created.
func (c *Client) RegenerateDashboard(ctx context.Context, req *RegenerateDashboard, opts ...api.Option) (*RegenerateDashboard_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/quality-monitoring/tables/%v/monitor/dashboard", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RegenerateDashboard_Response{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Queues a metric refresh on the monitor for
// the specified table. The refresh will execute in the background.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema - be an
// owner of the table
//
// Additionally, the call must be made from the workspace where the monitor was
// created.
func (c *Client) RunRefresh(ctx context.Context, req *RunRefresh, opts ...api.Option) (*RefreshInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor/refreshes", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RefreshInfo{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deprecated: Use Data Quality Monitors API instead
// (/api/data-quality/v1/monitors). Updates a monitor for the specified table.
//
// The caller must either: 1. be an owner of the table's parent catalog 2. have
// **USE_CATALOG** on the table's parent catalog and be an owner of the table's
// parent schema 3. have the following permissions: - **USE_CATALOG** on the
// table's parent catalog - **USE_SCHEMA** on the table's parent schema - be an
// owner of the table.
//
// Additionally, the call must be made from the workspace where the monitor was
// created, and the caller must be the original creator of the monitor.
//
// Certain configuration fields, such as output asset identifiers, cannot be
// updated.
func (c *Client) UpdateMonitor(ctx context.Context, req *UpdateMonitor, opts ...api.Option) (*DataMonitorInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/monitor", *req.FullTableNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DataMonitorInfo{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("PUT", baseURL.String(), bytes.NewBuffer(body))
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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}
