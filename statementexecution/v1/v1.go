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

// Requests that an executing statement be canceled. Callers must poll for
// status to see the terminal state. Cancel response is empty; receiving
// response indicates successful receipt.
func (c *Client) CancelStatement(ctx context.Context, req *CancelStatementRequest, opts ...api.Option) (*CancelStatementResponse, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/sql/statements/%v/cancel", *req.StatementId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CancelStatementResponse{}

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

// Execute a SQL statement and optionally await its results for a specified
// time.
//
// **Use case: small result sets with INLINE + JSON_ARRAY**
//
// For flows that generate small and predictable result sets (<= 25 MiB),
// `INLINE` responses of `JSON_ARRAY` result data are typically the simplest way
// to execute and fetch result data.
//
// **Use case: large result sets with EXTERNAL_LINKS**
//
// Using `EXTERNAL_LINKS` to fetch result data allows you to fetch large result
// sets efficiently. The main differences from using `INLINE` disposition are
// that the result data is accessed with URLs, and that there are 3 supported
// formats: `JSON_ARRAY`, `ARROW_STREAM` and `CSV` compared to only `JSON_ARRAY`
// with `INLINE`.
//
// ** URLs**
//
// External links point to data stored within your workspace's internal storage,
// in the form of a URL. The URLs are valid for only a short period, <= 15
// minutes. Alongside each `external_link` is an expiration field indicating the
// time at which the URL is no longer valid. In `EXTERNAL_LINKS` mode, chunks
// can be resolved and fetched multiple times and in parallel.
//
// ----
//
// ### **Warning: Databricks strongly recommends that you protect the URLs that
// are returned by the `EXTERNAL_LINKS` disposition.**
//
// When you use the `EXTERNAL_LINKS` disposition, a short-lived, URL is
// generated, which can be used to download the results directly from . As a
// short-lived is embedded in this URL, you should protect the URL.
//
// Because URLs are already generated with embedded temporary s, you must not
// set an `Authorization` header in the download requests.
//
// The `EXTERNAL_LINKS` disposition can be disabled upon request by creating a
// support case.
//
// See also [Security best
// practices](/sql/admin/sql-execution-tutorial.html#security-best-practices).
//
// ----
//
// StatementResponse contains `statement_id` and `status`; other fields might be
// absent or present depending on context. If the SQL warehouse fails to execute
// the provided statement, a 200 response is returned with `status.state` set to
// `FAILED` (in contrast to a failure when accepting the request, which results
// in a non-200 response). Details of the error can be found at `status.error`
// in case of execution failures.
func (c *Client) ExecuteStatement(ctx context.Context, req *ExecuteStatementRequest, opts ...api.Option) (*StatementResponse, error) {
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
	baseURL.Path = "/api/2.0/sql/statements/"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StatementResponse{}

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

// After the statement execution has `SUCCEEDED`, this request can be used to
// fetch any chunk by index. Whereas the first chunk with `chunk_index=0` is
// typically fetched with :method:statementexecution/executeStatement or
// :method:statementexecution/getStatement, this request can be used to fetch
// subsequent chunks. The response structure is identical to the nested `result`
// element described in the :method:statementexecution/getStatement request, and
// similarly includes the `next_chunk_index` and `next_chunk_internal_link`
// fields for simple iteration through the result set. Depending on
// `disposition`, the response returns chunks of data either inline, or as
// links.
func (c *Client) GetResultData(ctx context.Context, req *GetResultDataRequest, opts ...api.Option) (*ResultData, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/sql/statements/%v/result/chunks/%v", *req.StatementId, *req.ChunkIndex)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResultData{}

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

// This request can be used to poll for the statement's status.
// StatementResponse contains `statement_id` and `status`; other fields might be
// absent or present depending on context. When the `status.state` field is
// `SUCCEEDED` it will also return the result manifest and the first chunk of
// the result data. When the statement is in the terminal states `CANCELED`,
// `CLOSED` or `FAILED`, it returns HTTP 200 with the state set. After at least
// 12 hours in terminal state, the statement is removed from the warehouse and
// further calls will receive an HTTP 404 response.
//
// **NOTE** This call currently might take up to 5 seconds to get the latest
// status and result.
func (c *Client) GetStatementResult(ctx context.Context, req *GetStatementResultRequest, opts ...api.Option) (*StatementResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/sql/statements/%v", *req.StatementId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StatementResponse{}

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
