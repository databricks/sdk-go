package v2

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

// [DEPRECATED] Create a quality monitor on UC object. Use Data Quality
// Monitoring API instead.
func (c *Client) CreateQualityMonitor(ctx context.Context, req *CreateQualityMonitorRequest, opts ...ops.Option) (*QualityMonitor, error) {
	body, err := json.Marshal(req.QualityMonitor)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/quality-monitors"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &QualityMonitor{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// [DEPRECATED] Delete a quality monitor on UC object. Use Data Quality
// Monitoring API instead.
func (c *Client) DeleteQualityMonitor(ctx context.Context, req *DeleteQualityMonitorRequest, opts ...ops.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/quality-monitors/%v/%v", *req.ObjectType, *req.ObjectId)
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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return err
	}
	return nil
}

// [DEPRECATED] Read a quality monitor on UC object. Use Data Quality Monitoring
// API instead.
func (c *Client) GetQualityMonitor(ctx context.Context, req *GetQualityMonitorRequest, opts ...ops.Option) (*QualityMonitor, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/quality-monitors/%v/%v", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &QualityMonitor{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// [DEPRECATED] (Unimplemented) List quality monitors. Use Data Quality
// Monitoring API instead.
func (c *Client) ListQualityMonitor(ctx context.Context, req *ListQualityMonitorRequest, opts ...ops.Option) (*ListQualityMonitorResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/quality-monitors"
	queryParams := url.Values{}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListQualityMonitorResponse{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListQualityMonitorIter(ctx context.Context, req *ListQualityMonitorRequest, opts ...ops.Option) iter.Seq2[*QualityMonitor, error] {
	return func(yield func(*QualityMonitor, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListQualityMonitorRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListQualityMonitor(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.QualityMonitors {
				if !yield(&resp.QualityMonitors[i], nil) {
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

// [DEPRECATED] (Unimplemented) Update a quality monitor on UC object. Use Data
// Quality Monitoring API instead.
func (c *Client) UpdateQualityMonitor(ctx context.Context, req *UpdateQualityMonitorRequest, opts ...ops.Option) (*QualityMonitor, error) {
	body, err := json.Marshal(req.QualityMonitor)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/quality-monitors/%v/%v", *req.ObjectType, *req.ObjectId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &QualityMonitor{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}
