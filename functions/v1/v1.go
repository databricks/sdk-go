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

// **WARNING: This API is experimental and will change in future versions**
//
// # Creates a new function
//
// The user must have the following permissions in order for the function to be
// created: - **USE_CATALOG** on the function's parent catalog - **USE_SCHEMA**
// and **CREATE_FUNCTION** on the function's parent schema
func (c *Client) CreateFunction(ctx context.Context, req *CreateFunctionRequest, opts ...api.Option) (*FunctionInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/functions"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &FunctionInfo{}

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

// Deletes the function that matches the supplied name. For the deletion to
// succeed, the user must satisfy one of the following conditions: - Is the
// owner of the function's parent catalog - Is the owner of the function's
// parent schema and have the **USE_CATALOG** privilege on its parent catalog -
// Is the owner of the function itself and have both the **USE_CATALOG**
// privilege on its parent catalog and the **USE_SCHEMA** privilege on its
// parent schema
func (c *Client) DeleteFunction(ctx context.Context, req *DeleteFunction, opts ...api.Option) (*DeleteFunction_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/functions/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.Force != nil {
		queryParams.Add("force", fmt.Sprintf("%v", *req.Force))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteFunction_Response{}

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

// Gets a function from within a parent catalog and schema. For the fetch to
// succeed, the user must satisfy one of the following requirements: - Is a
// metastore admin - Is an owner of the function's parent catalog - Have the
// **USE_CATALOG** privilege on the function's parent catalog and be the owner
// of the function - Have the **USE_CATALOG** privilege on the function's parent
// catalog, the **USE_SCHEMA** privilege on the function's parent schema, and
// the **EXECUTE** privilege on the function itself
func (c *Client) GetFunction(ctx context.Context, req *GetFunction, opts ...api.Option) (*FunctionInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/functions/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &FunctionInfo{}

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

// List functions within the specified parent catalog and schema. If the user is
// a metastore admin, all functions are returned in the output list. Otherwise,
// the user must have the **USE_CATALOG** privilege on the catalog and the
// **USE_SCHEMA** privilege on the schema, and the output list contains only
// functions for which either the user has the **EXECUTE** privilege or the user
// is the owner. There is no guarantee of a specific ordering of the elements in
// the array.
//
// NOTE: we recommend using max_results=0 to use the paginated version of this
// API. Unpaginated calls will be deprecated soon.
//
// PAGINATION BEHAVIOR: When using pagination (max_results >= 0), a page may
// contain zero results while still providing a next_page_token. Clients must
// continue reading pages until next_page_token is absent, which is the only
// indication that the end of results has been reached.
func (c *Client) ListFunctions(ctx context.Context, req *ListFunctions, opts ...api.Option) (*ListFunctions_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/functions"
	queryParams := url.Values{}
	if req.CatalogName != nil {
		queryParams.Add("catalog_name", fmt.Sprintf("%v", *req.CatalogName))
	}
	if req.SchemaName != nil {
		queryParams.Add("schema_name", fmt.Sprintf("%v", *req.SchemaName))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListFunctions_Response{}

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

func (c *Client) ListFunctionsIter(ctx context.Context, req *ListFunctions, opts ...api.Option) iter.Seq2[*FunctionInfo, error] {
	return func(yield func(*FunctionInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListFunctions{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListFunctions(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Functions {
				if !yield(&resp.Functions[i], nil) {
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

// Updates the function that matches the supplied name. Only the owner of the
// function can be updated. If the user is not a metastore admin, the user must
// be a member of the group that is the new function owner. - Is a metastore
// admin - Is the owner of the function's parent catalog - Is the owner of the
// function's parent schema and has the **USE_CATALOG** privilege on its parent
// catalog - Is the owner of the function itself and has the **USE_CATALOG**
// privilege on its parent catalog as well as the **USE_SCHEMA** privilege on
// the function's parent schema.
func (c *Client) UpdateFunction(ctx context.Context, req *UpdateFunction, opts ...api.Option) (*FunctionInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/functions/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &FunctionInfo{}

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
