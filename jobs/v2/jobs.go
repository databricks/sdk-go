package jobs

import (
	"bytes"
	"context"
	"encoding/json"
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
	}, nil
}

func (c *Client) CreateJob(ctx context.Context, req *CreateJobRequest, opts ...api.Option) (*Job, error) {
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
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Job{}
	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL.String(), bytes.NewBuffer(body))
		if err != nil {
			return err
		}
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

func ListJobs(ctx context.Context, req *ListJobsRequest, opts ...api.Option) (*ListJobsResponse, error) {
	return nil, nil
}

func ListJobsIter(ctx context.Context, req *ListJobsRequest, opts ...api.Option) iter.Seq2[*Job, error] {
	return func(yield func(*Job, error) bool) {
		pageReq := &ListJobsRequest{
			PageSize:  req.PageSize,
			PageToken: req.PageToken,
		}
		for {
			resp, err := ListJobs(ctx, pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Jobs {
				if !yield(&resp.Jobs[i], nil) {
					return
				}
			}
			if resp.NextPageToken == "" {
				return
			}
			pageReq.PageToken = resp.NextPageToken
		}
	}
}
