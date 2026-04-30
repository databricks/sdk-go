package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

// Create a new task.
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest, opts ...call.Option) (*Task, error) {
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
	baseURL.Path = "/api/2.0/tasks"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Task{}

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

// Create a new Task.
//
// Returns a waiter.
func (c *Client) CreateTaskWaiter(ctx context.Context, req *CreateTaskRequest, opts ...call.Option) (*CreateTaskWaiter, error) {
	resp, err := c.CreateTask(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if resp.TaskId == nil {
		return nil, fmt.Errorf("nil TaskId in response")
	}
	return &CreateTaskWaiter{
		rawResponse: resp,
		client:      c,
		taskId:      resp.TaskId,
	}, nil
}

type CreateTaskWaiter struct {
	rawResponse *Task
	client      *Client
	taskId      *string
}
