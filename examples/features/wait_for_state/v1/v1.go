package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// Wait polls until a terminal state is reached or an error is encountered.
// It returns an error if a failure state is reached.
func (w *CreateTaskWaiter) Wait(ctx context.Context, opts ...api.Option) (*Task, error) {
	errStillRunning := errors.New("waiting for completion")
	var result *Task

	pollOpts := append(append([]api.Option{}, opts...), api.WithTimeout(0))

	call := func(ctx context.Context) error {
		pollReq := &GetTaskRequest{}
		pollReq.TaskId = w.taskId

		pollResp, err := w.service.GetTask(ctx, pollReq, pollOpts...)
		if err != nil {
			return err
		}

		if pollResp == nil || pollResp.Status == nil || pollResp.Status.State == nil {
			return fmt.Errorf("response missing required status field")
		}
		state := *pollResp.Status.State

		switch state {
		case TaskStateCompleted, TaskStateCancelled:
			result = pollResp
			return nil
		case TaskStateFailed, TaskStateInternalError:
			msg := "(no message)"
			if pollResp.Status.Message != nil {
				msg = *pollResp.Status.Message
			}
			return fmt.Errorf("terminal state %s: %s", state, msg)
		default:
			return errStillRunning
		}
	}

	waitOpts := append(append([]api.Option{}, opts...), api.WithLimiter(nil))
	waitOpts = append(waitOpts, api.WithRetrier(func() api.Retrier {
		return api.RetryOn(api.BackoffPolicy{}, func(err error) bool {
			return errors.Is(err, errStillRunning)
		})
	}))
	if err := api.Execute(ctx, call, waitOpts...); err != nil {
		return nil, err
	}
	return result, nil
}

// Done checks whether the operation has reached a terminal state (success or
// failure). It returns (false, error) only if the poll request itself fails.
// To wait for completion and get the result or failure error, use Wait instead.
func (w *CreateTaskWaiter) Done(ctx context.Context, opts ...api.Option) (bool, error) {
	pollReq := &GetTaskRequest{}
	pollReq.TaskId = w.taskId

	pollResp, err := w.service.GetTask(ctx, pollReq, opts...)
	if err != nil {
		return false, err
	}

	if pollResp == nil || pollResp.Status == nil || pollResp.Status.State == nil {
		return false, fmt.Errorf("response missing required status field")
	}
	state := *pollResp.Status.State

	switch state {
	case TaskStateCompleted, TaskStateCancelled:
		return true, nil
	case TaskStateFailed, TaskStateInternalError:
		return true, nil
	default:
		return false, nil
	}
}

// GetTaskId returns the taskId used for polling.
func (w *CreateTaskWaiter) GetTaskId() string {
	return *w.taskId
}

// Get the current state of a task.
func (c *Client) GetTask(ctx context.Context, req *GetTaskRequest, opts ...api.Option) (*Task, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/tasks/%v", *req.TaskId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Task{}

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
