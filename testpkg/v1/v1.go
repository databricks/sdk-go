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

// Cancel a running task. Returns the task with updated status once cancellation
// is confirmed.
func (c *Client) CancelTask(ctx context.Context, req *CancelTaskRequest, opts ...api.Option) (*CancelTaskWaiter, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/tasks/%v/cancel", *req.TaskId)
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
	return &CancelTaskWaiter{
		rawResponse: resp,
		service:     c,
		taskId:      req.TaskId,
	}, nil
}

type CancelTaskWaiter struct {
	rawResponse *Task
	service     *Client
	taskId      *string
}

// Wait polls the server until the operation reaches a terminal state or encounters an error.
// This method will return an error if a failure state is reached or an unknown state is encountered.
func (w *CancelTaskWaiter) Wait(ctx context.Context, opts ...api.Option) (*Task, error) {
	errOperationInProgress := errors.New("operation still in progress")
	var result *Task

	call := func(ctx context.Context) error {
		pollReq := &GetTaskRequest{}
		pollReq.TaskId = w.taskId

		pollResp, err := w.service.GetTask(ctx, pollReq)
		if err != nil {
			return err
		}

		state := *pollResp.Status.State

		switch state {
		case TaskStateCancelled:
			result = pollResp
			return nil
		case TaskStateFailed:
			return fmt.Errorf("operation failed with state %s: %s", state, *pollResp.Status.Message)
		case TaskStateTaskStateUnspecified, TaskStatePending, TaskStateRunning, TaskStateCompleted, TaskStateInternalError:
			return errOperationInProgress
		default:
			return fmt.Errorf("unknown state %q while waiting for operation - consider updating SDK", state)
		}
	}

	retrier := api.RetryOn(api.BackoffPolicy{}, func(err error) bool {
		return errors.Is(err, errOperationInProgress)
	})

	defaultOpts := []api.Option{
		api.WithRetrier(func() api.Retrier { return retrier }),
	}
	allOpts := append(defaultOpts, opts...)

	if err := api.Execute(ctx, call, allOpts...); err != nil {
		return nil, err
	}
	return result, nil
}

// Done reports whether the operation has completed.
// Returns an error if an unknown state is encountered.
func (w *CancelTaskWaiter) Done() (bool, error) {
	pollReq := &GetTaskRequest{}
	pollReq.TaskId = w.taskId

	pollResp, err := w.service.GetTask(context.Background(), pollReq)
	if err != nil {
		return false, err
	}

	state := *pollResp.Status.State

	switch state {
	case TaskStateCancelled:
		return true, nil
	case TaskStateFailed:
		return true, nil
	case TaskStateTaskStateUnspecified, TaskStatePending, TaskStateRunning, TaskStateCompleted, TaskStateInternalError:
		return false, nil
	default:
		return false, fmt.Errorf("unknown state %q while waiting for operation - consider updating SDK", state)
	}
}

// Create a new task and start its execution. This method returns immediately,
// but the task continues running. Use GetTask to poll for completion.
func (c *Client) CreateTask(ctx context.Context, req *CreateTaskRequest, opts ...api.Option) (*CreateTaskWaiter, error) {
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

// Wait polls the server until the operation reaches a terminal state or encounters an error.
// This method will return an error if a failure state is reached or an unknown state is encountered.
func (w *CreateTaskWaiter) Wait(ctx context.Context, opts ...api.Option) (*Task, error) {
	errOperationInProgress := errors.New("operation still in progress")
	var result *Task

	call := func(ctx context.Context) error {
		pollReq := &GetTaskRequest{}
		pollReq.TaskId = w.taskId

		pollResp, err := w.service.GetTask(ctx, pollReq)
		if err != nil {
			return err
		}

		state := *pollResp.Status.State

		switch state {
		case TaskStateCompleted, TaskStateCancelled:
			result = pollResp
			return nil
		case TaskStateFailed, TaskStateInternalError:
			return fmt.Errorf("operation failed with state %s: %s", state, *pollResp.Status.Message)
		case TaskStateTaskStateUnspecified, TaskStatePending, TaskStateRunning:
			return errOperationInProgress
		default:
			return fmt.Errorf("unknown state %q while waiting for operation - consider updating SDK", state)
		}
	}

	retrier := api.RetryOn(api.BackoffPolicy{}, func(err error) bool {
		return errors.Is(err, errOperationInProgress)
	})

	defaultOpts := []api.Option{
		api.WithRetrier(func() api.Retrier { return retrier }),
	}
	allOpts := append(defaultOpts, opts...)

	if err := api.Execute(ctx, call, allOpts...); err != nil {
		return nil, err
	}
	return result, nil
}

// Done reports whether the operation has completed.
// Returns an error if an unknown state is encountered.
func (w *CreateTaskWaiter) Done() (bool, error) {
	pollReq := &GetTaskRequest{}
	pollReq.TaskId = w.taskId

	pollResp, err := w.service.GetTask(context.Background(), pollReq)
	if err != nil {
		return false, err
	}

	state := *pollResp.Status.State

	switch state {
	case TaskStateCompleted, TaskStateCancelled:
		return true, nil
	case TaskStateFailed, TaskStateInternalError:
		return true, nil
	case TaskStateTaskStateUnspecified, TaskStatePending, TaskStateRunning:
		return false, nil
	default:
		return false, fmt.Errorf("unknown state %q while waiting for operation - consider updating SDK", state)
	}
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
