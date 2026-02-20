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

// Create a new task and start its execution. This method returns immediately,
// but the task continues running. Use GetTask to poll for completion.
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest, opts ...api.Option) (*Task, error) {
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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateTaskWaiter starts the operation and returns a waiter to track its completion.
func (c *Client) CreateTaskWaiter(ctx context.Context, req *CreateTaskRequest, opts ...api.Option) (*CreateTaskWaiter, error) {
	resp, err := c.CreateTask(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if resp.TaskId == nil {
		return nil, fmt.Errorf("response field TaskId required for polling is nil")
	}
	return &CreateTaskWaiter{
		rawResponse: resp,
		service:     c,
		taskId:      resp.TaskId,
	}, nil
}

type CreateTaskWaiter struct {
	rawResponse *Task
	service     *Client
	taskId      *string
}
