package v1

import (
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

// Returns billable usage logs in CSV format for the specified account and date
// range. For the data schema, see:
//
// - AWS: [CSV file schema]. - GCP: [CSV file schema].
//
// Note that this method might take multiple minutes to complete.
//
// **Warning**: Depending on the queried date range, the number of workspaces in
// the account, the size of the response and the internet speed of the caller,
// this API may hit a timeout after a few minutes. If you experience this, try
// to mitigate by calling the API with narrower date ranges.
//
// [CSV file schema]: https://docs.gcp.databricks.com/administration-guide/account-settings/usage-analysis.html#csv-file-schema
func (c *Client) Download(ctx context.Context, req *DownloadRequest, opts ...api.Option) (*DownloadResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/accounts/{account_id}/usage/download"
	queryParams := url.Values{}
	if req.StartMonth != nil {
		queryParams.Add("start_month", fmt.Sprintf("%v", *req.StartMonth))
	}
	if req.EndMonth != nil {
		queryParams.Add("end_month", fmt.Sprintf("%v", *req.EndMonth))
	}
	if req.PersonalData != nil {
		queryParams.Add("personal_data", fmt.Sprintf("%v", *req.PersonalData))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DownloadResponse{}

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
