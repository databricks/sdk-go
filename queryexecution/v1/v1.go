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

// Cancel the results for the a query for a published, embedded dashboard.
func (c *Client) CancelPublishedQueryExecution(ctx context.Context, req *CancelPublishedQueryExecutionRequest, opts ...api.Option) (*CancelQueryExecutionResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/lakeview-query/query/published"
	queryParams := url.Values{}
	if len(req.Tokens) > 0 {
		queryParams.Add("tokens", fmt.Sprintf("%v", req.Tokens))
	}
	if req.DashboardName != nil {
		queryParams.Add("dashboard_name", fmt.Sprintf("%v", *req.DashboardName))
	}
	if req.DashboardRevisionId != nil {
		queryParams.Add("dashboard_revision_id", fmt.Sprintf("%v", *req.DashboardRevisionId))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelQueryExecutionResponse{}

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

// Execute a query for a published dashboard.
func (c *Client) ExecutePublishedDashboardQuery(ctx context.Context, req *ExecutePublishedDashboardQueryRequest, opts ...api.Option) (*ExecuteQueryResponse, error) {
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
	baseURL.Path = "/api/2.0/lakeview-query/query/published"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ExecuteQueryResponse{}

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

// Poll the results for the a query for a published, embedded dashboard.
// Supports both GET and POST methods. POST is recommended for polling many
// tokens to avoid URL length limitations.
func (c *Client) PollPublishedQueryStatus(ctx context.Context, req *PollPublishedQueryStatusRequest, opts ...api.Option) (*PollQueryStatusResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/lakeview-query/query/published"
	queryParams := url.Values{}
	if len(req.Tokens) > 0 {
		queryParams.Add("tokens", fmt.Sprintf("%v", req.Tokens))
	}
	if req.DashboardName != nil {
		queryParams.Add("dashboard_name", fmt.Sprintf("%v", *req.DashboardName))
	}
	if req.DashboardRevisionId != nil {
		queryParams.Add("dashboard_revision_id", fmt.Sprintf("%v", *req.DashboardRevisionId))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PollQueryStatusResponse{}

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
