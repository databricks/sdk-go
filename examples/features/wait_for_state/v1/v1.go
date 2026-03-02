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

var errMissingStatus = errors.New("response missing required status field")

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

// Create a new task.
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

// Create a new Task.
//
// Returns a waiter.
func (c *Client) CreateTaskWaiter(ctx context.Context, req *CreateTaskRequest, opts ...api.Option) (*CreateTaskWaiter, error) {
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

// Wait blocks until the task reaches a terminal state.
//
// Success states ([TaskStateCompleted], [TaskStateCancelled]) return the final
// [Task]. Failure states ([TaskStateFailed], [TaskStateInternalError]) return
// an error.
func (w *CreateTaskWaiter) Wait(ctx context.Context, opts ...api.Option) (*Task, error) {
	errStillRunning := errors.New("waiting for completion")
	var result *Task

	// Wait uses two levels of Execute: an outer one for the poll loop, and
	// an inner one (inside GetTask) for each HTTP call. The user's timeout
	// applies to the poll loop, not each individual call.
	pollOpts := append([]api.Option(nil), opts...)
	pollOpts = append(pollOpts, api.WithTimeout(0))

	call := func(ctx context.Context) error {
		pollReq := &GetTaskRequest{}
		pollReq.TaskId = w.taskId

		pollResp, err := w.client.GetTask(ctx, pollReq, pollOpts...)
		if err != nil {
			return err
		}

		if pollResp == nil || pollResp.Status == nil || pollResp.Status.State == nil {
			return errMissingStatus
		}
		state := *pollResp.Status.State

		switch state {
		case TaskStateCompleted, TaskStateCancelled:
			result = pollResp
			return nil
		case TaskStateFailed, TaskStateInternalError:
			// TODO: use a typed error so callers can distinguish task failure
			// from transport errors via errors.As without parsing the message.
			msg := "(no message)"
			if pollResp.Status.Message != nil {
				msg = *pollResp.Status.Message
			}
			return fmt.Errorf("terminal state %s: %s", state, msg)
		default:
			return errStillRunning
		}
	}

	// Rate limiting belongs at the HTTP call level, not the poll loop. The
	// retrier here determines when to stop polling, rather than when to retry
	// HTTP errors.
	waitOpts := append([]api.Option(nil), opts...)
	waitOpts = append(waitOpts, api.WithLimiter(nil), api.WithRetrier(func() api.Retrier {
		return api.RetryOn(api.BackoffPolicy{}, func(err error) bool {
			return errors.Is(err, errStillRunning)
		})
	}))

	if err := api.Execute(ctx, call, waitOpts...); err != nil {
		return nil, err
	}
	return result, nil
}

// Done returns true if a terminal state is reached (success or failure).
func (w *CreateTaskWaiter) Done(ctx context.Context, opts ...api.Option) (bool, error) {
	pollReq := &GetTaskRequest{}
	pollReq.TaskId = w.taskId

	pollResp, err := w.client.GetTask(ctx, pollReq, opts...)
	if err != nil {
		return false, err
	}

	if pollResp == nil || pollResp.Status == nil || pollResp.Status.State == nil {
		return false, errMissingStatus
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

// TaskId returns the task ID.
func (w *CreateTaskWaiter) TaskId() string {
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
