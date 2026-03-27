package v1

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

// Creates a new metastore based on a provided name and optional storage root
// path. By default (if the __owner__ field is not set), the owner of the new
// metastore is the user calling the __createMetastore__ API. If the __owner__
// field is set to the empty string (**""**), the ownership is assigned to the
// System User instead.
func (c *Client) CreateMetastore(ctx context.Context, req *CreateMetastore, opts ...api.Option) (*MetastoreInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/metastores"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &MetastoreInfo{}

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

// Creates a new metastore assignment. If an assignment for the same
// __workspace_id__ exists, it will be overwritten by the new __metastore_id__
// and __default_catalog_name__. The caller must be an account admin.
func (c *Client) CreateMetastoreAssignment(ctx context.Context, req *CreateMetastoreAssignment, opts ...api.Option) (*CreateMetastoreAssignment_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/workspaces/%v/metastore", *req.WorkspaceId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CreateMetastoreAssignment_Response{}

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

	if err := api.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deletes a metastore. The caller must be a metastore admin.
func (c *Client) DeleteMetastore(ctx context.Context, req *DeleteMetastore, opts ...api.Option) (*DeleteMetastore_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/metastores/%v", *req.Id)
	queryParams := url.Values{}
	if req.Force != nil {
		queryParams.Add("force", fmt.Sprintf("%v", *req.Force))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteMetastore_Response{}

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

// Deletes a metastore assignment. The caller must be an account administrator.
func (c *Client) DeleteMetastoreAssignment(ctx context.Context, req *DeleteMetastoreAssignment, opts ...api.Option) (*DeleteMetastoreAssignment_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/workspaces/%v/metastore", *req.WorkspaceId)
	queryParams := url.Values{}
	if req.MetastoreId != nil {
		queryParams.Add("metastore_id", fmt.Sprintf("%v", *req.MetastoreId))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteMetastoreAssignment_Response{}

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

// Gets the metastore assignment for the workspace being accessed.
func (c *Client) GetCurrentMetastoreAssignment(ctx context.Context, req *GetCurrentMetastoreAssignment, opts ...api.Option) (*MetastoreAssignment, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/current-metastore-assignment"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &MetastoreAssignment{}

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

// Gets a metastore that matches the supplied ID. The caller must be a metastore
// admin to retrieve this info.
func (c *Client) GetMetastore(ctx context.Context, req *GetMetastore, opts ...api.Option) (*MetastoreInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/metastores/%v", *req.Id)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &MetastoreInfo{}

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

// Gets information about a metastore. This summary includes the storage
// credential, the cloud vendor, the cloud region, and the global metastore ID.
func (c *Client) GetMetastoreSummary(ctx context.Context, req *GetMetastoreSummary, opts ...api.Option) (*GetMetastoreSummary_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/metastore_summary"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetMetastoreSummary_Response{}

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

// Lists metastores across all shards in the account, federating one shard at a
// time.
func (c *Client) ListGlobalMetastores(ctx context.Context, req *ListGlobalMetastoresRequest, opts ...api.Option) (*ListGlobalMetastoresResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/metastores:globalList"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.ShardLocalOnly != nil {
		queryParams.Add("shard_local_only", fmt.Sprintf("%v", *req.ShardLocalOnly))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListGlobalMetastoresResponse{}

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

func (c *Client) ListGlobalMetastoresIter(ctx context.Context, req *ListGlobalMetastoresRequest, opts ...api.Option) iter.Seq2[*MetastoreInfo, error] {
	return func(yield func(*MetastoreInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListGlobalMetastoresRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListGlobalMetastores(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Metastores {
				if !yield(&resp.Metastores[i], nil) {
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

// Gets an array of the available metastores (as __MetastoreInfo__ objects). The
// caller must be an admin to retrieve this info. There is no guarantee of a
// specific ordering of the elements in the array.
//
// NOTE: we recommend using max_results=0 to use the paginated version of this
// API. Unpaginated calls will be deprecated soon.
//
// PAGINATION BEHAVIOR: When using pagination (max_results >= 0), a page may
// contain zero results while still providing a next_page_token. Clients must
// continue reading pages until next_page_token is absent, which is the only
// indication that the end of results has been reached.
func (c *Client) ListMetastores(ctx context.Context, req *ListMetastores, opts ...api.Option) (*ListMetastores_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/metastores"
	queryParams := url.Values{}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListMetastores_Response{}

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

func (c *Client) ListMetastoresIter(ctx context.Context, req *ListMetastores, opts ...api.Option) iter.Seq2[*MetastoreInfo, error] {
	return func(yield func(*MetastoreInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListMetastores{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListMetastores(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Metastores {
				if !yield(&resp.Metastores[i], nil) {
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

// Updates information for a specific metastore. The caller must be a metastore
// admin. If the __owner__ field is set to the empty string (**""**), the
// ownership is updated to the System User.
func (c *Client) UpdateMetastore(ctx context.Context, req *UpdateMetastore, opts ...api.Option) (*MetastoreInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/metastores/%v", *req.Id)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &MetastoreInfo{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("PATCH", baseURL.String(), bytes.NewBuffer(body))
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

// Updates a metastore assignment. This operation can be used to update
// __metastore_id__ or __default_catalog_name__ for a specified Workspace, if
// the Workspace is already assigned a metastore. The caller must be an account
// admin to update __metastore_id__; otherwise, the caller can be a Workspace
// admin.
func (c *Client) UpdateMetastoreAssignment(ctx context.Context, req *UpdateMetastoreAssignment, opts ...api.Option) (*UpdateMetastoreAssignment_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/workspaces/%v/metastore", *req.WorkspaceId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UpdateMetastoreAssignment_Response{}

	call := func(ctx context.Context) error {
		httpReq, err := http.NewRequest("PATCH", baseURL.String(), bytes.NewBuffer(body))
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
