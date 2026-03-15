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

// Cancels all active runs of a job. The runs are canceled asynchronously, so it
// doesn't prevent new runs from being started.
func (c *Client) CancelAllRuns(ctx context.Context, req *CancelAllRuns, opts ...api.Option) (*CancelAllRuns_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/runs/cancel-all"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelAllRuns_Response{}

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

// Cancels a job run or a task run. The run is canceled asynchronously, so it
// may still be running when this request completes.
func (c *Client) CancelRun(ctx context.Context, req *CancelRun, opts ...api.Option) (*CancelRun_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/runs/cancel"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelRun_Response{}

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

// Create a new job.
func (c *Client) CreateJob(ctx context.Context, req *CreateJob, opts ...api.Option) (*CreateJob_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/create"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CreateJob_Response{}

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

// Deletes a job.
func (c *Client) DeleteJob(ctx context.Context, req *DeleteJob, opts ...api.Option) (*DeleteJob_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteJob_Response{}

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

// Deletes a non-active run. Returns an error if the run is active.
func (c *Client) DeleteRun(ctx context.Context, req *DeleteRun, opts ...api.Option) (*DeleteRun_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/runs/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteRun_Response{}

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

// Export and retrieve the job run task.
func (c *Client) ExportRun(ctx context.Context, req *ExportRun, opts ...api.Option) (*ExportRun_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/runs/export"
	queryParams := url.Values{}
	if req.RunId != nil {
		queryParams.Add("run_id", fmt.Sprintf("%v", *req.RunId))
	}
	if req.ViewsToExport != nil {
		queryParams.Add("views_to_export", fmt.Sprintf("%v", *req.ViewsToExport))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ExportRun_Response{}

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

// Retrieves the details for a single job.
//
// Large arrays in the results will be paginated when they exceed 100 elements.
// A request for a single job will return all properties for that job, and the
// first 100 elements of array properties (`tasks`, `job_clusters`,
// `environments` and `parameters`). Use the `next_page_token` field to check
// for more results and pass its value as the `page_token` in subsequent
// requests. If any array properties have more than 100 elements, additional
// results will be returned on subsequent requests. Arrays without additional
// results will be empty on later pages.
func (c *Client) GetJob(ctx context.Context, req *GetJob, opts ...api.Option) (*GetJob_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/get"
	queryParams := url.Values{}
	if req.JobId != nil {
		queryParams.Add("job_id", fmt.Sprintf("%v", *req.JobId))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetJob_Response{}

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

// Retrieves the metadata of a run.
//
// Large arrays in the results will be paginated when they exceed 100 elements.
// A request for a single run will return all properties for that run, and the
// first 100 elements of array properties (`tasks`, `job_clusters`,
// `job_parameters` and `repair_history`). Use the next_page_token field to
// check for more results and pass its value as the page_token in subsequent
// requests. If any array properties have more than 100 elements, additional
// results will be returned on subsequent requests. Arrays without additional
// results will be empty on later pages.
func (c *Client) GetRun(ctx context.Context, req *GetRun, opts ...api.Option) (*GetRun_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/runs/get"
	queryParams := url.Values{}
	if req.RunId != nil {
		queryParams.Add("run_id", fmt.Sprintf("%v", *req.RunId))
	}
	if req.IncludeHistory != nil {
		queryParams.Add("include_history", fmt.Sprintf("%v", *req.IncludeHistory))
	}
	if req.IncludeResolvedValues != nil {
		queryParams.Add("include_resolved_values", fmt.Sprintf("%v", *req.IncludeResolvedValues))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetRun_Response{}

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

// Retrieve the output and metadata of a single task run. When a notebook task
// returns a value through the `dbutils.notebook.exit()` call, you can use this
// endpoint to retrieve that value. <Databricks> restricts this API to returning
// the first 5 MB of the output. To return a larger result, you can store job
// results in a cloud storage service.
//
// This endpoint validates that the __run_id__ parameter is valid and returns an
// HTTP status code 400 if the __run_id__ parameter is invalid. Runs are
// automatically removed after 60 days. If you to want to reference them beyond
// 60 days, you must save old run results before they expire.
func (c *Client) GetRunOutput(ctx context.Context, req *GetRunOutput, opts ...api.Option) (*GetRunOutput_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/runs/get-output"
	queryParams := url.Values{}
	if req.RunId != nil {
		queryParams.Add("run_id", fmt.Sprintf("%v", *req.RunId))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetRunOutput_Response{}

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

// Retrieves a list of jobs.
func (c *Client) ListJobs(ctx context.Context, req *ListJobs, opts ...api.Option) (*ListJobs_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/list"
	queryParams := url.Values{}
	if req.Offset != nil {
		queryParams.Add("offset", fmt.Sprintf("%v", *req.Offset))
	}
	if req.Limit != nil {
		queryParams.Add("limit", fmt.Sprintf("%v", *req.Limit))
	}
	if req.ExpandTasks != nil {
		queryParams.Add("expand_tasks", fmt.Sprintf("%v", *req.ExpandTasks))
	}
	if req.Name != nil {
		queryParams.Add("name", fmt.Sprintf("%v", *req.Name))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListJobs_Response{}

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

func (c *Client) ListJobsIter(ctx context.Context, req *ListJobs, opts ...api.Option) iter.Seq2[*BaseJob, error] {
	return func(yield func(*BaseJob, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListJobs{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListJobs(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Jobs {
				if !yield(&resp.Jobs[i], nil) {
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

// List runs in descending order by start time.
func (c *Client) ListRuns(ctx context.Context, req *ListRuns, opts ...api.Option) (*ListRuns_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.2/jobs/runs/list"
	queryParams := url.Values{}
	if req.JobId != nil {
		queryParams.Add("job_id", fmt.Sprintf("%v", *req.JobId))
	}
	if req.ActiveOnly != nil {
		queryParams.Add("active_only", fmt.Sprintf("%v", *req.ActiveOnly))
	}
	if req.CompletedOnly != nil {
		queryParams.Add("completed_only", fmt.Sprintf("%v", *req.CompletedOnly))
	}
	if req.Offset != nil {
		queryParams.Add("offset", fmt.Sprintf("%v", *req.Offset))
	}
	if req.Limit != nil {
		queryParams.Add("limit", fmt.Sprintf("%v", *req.Limit))
	}
	if req.RunType != nil {
		queryParams.Add("run_type", fmt.Sprintf("%v", *req.RunType))
	}
	if req.ExpandTasks != nil {
		queryParams.Add("expand_tasks", fmt.Sprintf("%v", *req.ExpandTasks))
	}
	if req.StartTimeFrom != nil {
		queryParams.Add("start_time_from", fmt.Sprintf("%v", *req.StartTimeFrom))
	}
	if req.StartTimeTo != nil {
		queryParams.Add("start_time_to", fmt.Sprintf("%v", *req.StartTimeTo))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListRuns_Response{}

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

func (c *Client) ListRunsIter(ctx context.Context, req *ListRuns, opts ...api.Option) iter.Seq2[*BaseRun, error] {
	return func(yield func(*BaseRun, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListRuns{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListRuns(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Runs {
				if !yield(&resp.Runs[i], nil) {
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

// Re-run one or more tasks. Tasks are re-run as part of the original job run.
// They use the current job and task settings, and can be viewed in the history
// for the original job run.
func (c *Client) Repair(ctx context.Context, req *RepairRun, opts ...api.Option) (*RepairRun_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/runs/repair"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RepairRun_Response{}

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

// Overwrite all settings for the given job. Use the [_Update_
// endpoint](:method:jobs/update) to update job settings partially.
func (c *Client) ResetJob(ctx context.Context, req *ResetJob, opts ...api.Option) (*ResetJob_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/reset"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResetJob_Response{}

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

// Run a job and return the `run_id` of the triggered run.
func (c *Client) RunNow(ctx context.Context, req *RunNow, opts ...api.Option) (*RunNow_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/run-now"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RunNow_Response{}

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

// Submit a one-time run. This endpoint allows you to submit a workload directly
// without creating a job. Runs submitted using this endpoint don’t display in
// the UI. Use the `jobs/runs/get` API to check the run state after the job is
// submitted.
//
// **Important:** Jobs submitted using this endpoint are not saved as a job.
// They do not show up in the Jobs UI, and do not retry when they fail. Because
// they are not saved, <Databricks> cannot auto-optimize serverless compute in
// case of failure. If your job fails, you may want to use classic compute to
// specify the compute needs for the job. Alternatively, use the `POST
// /jobs/create` and `POST /jobs/run-now` endpoints to create and run a saved
// job.
func (c *Client) SubmitRun(ctx context.Context, req *SubmitRun, opts ...api.Option) (*SubmitRun_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/runs/submit"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &SubmitRun_Response{}

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

// Add, update, or remove specific settings of an existing job. Use the [_Reset_
// endpoint](:method:jobs/reset) to overwrite all job settings.
func (c *Client) UpdateJob(ctx context.Context, req *UpdateJob, opts ...api.Option) (*UpdateJob_Response, error) {
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
	baseURL.Path = "/api/2.2/jobs/update"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UpdateJob_Response{}

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
