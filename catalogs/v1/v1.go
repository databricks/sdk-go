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

// Creates a new catalog instance in the parent metastore if the caller is a
// metastore admin or has the **CREATE_CATALOG** privilege.
func (c *Client) CreateCatalog(ctx context.Context, req *CreateCatalog, opts ...api.Option) (*CatalogInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/catalogs"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CatalogInfo{}

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

// Deletes the catalog that matches the supplied name. The caller must be a
// metastore admin or the owner of the catalog.
func (c *Client) DeleteCatalog(ctx context.Context, req *DeleteCatalog, opts ...api.Option) (*DeleteCatalog_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/catalogs/%v", *req.NameArg)
	queryParams := url.Values{}
	if req.Force != nil {
		queryParams.Add("force", fmt.Sprintf("%v", *req.Force))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteCatalog_Response{}

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

// Gets the specified catalog in a metastore. The caller must be a metastore
// admin, the owner of the catalog, or a user that has the **USE_CATALOG**
// privilege set for their account.
func (c *Client) GetCatalog(ctx context.Context, req *GetCatalog, opts ...api.Option) (*CatalogInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/catalogs/%v", *req.NameArg)
	queryParams := url.Values{}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CatalogInfo{}

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

// Gets an array of catalogs in the metastore. If the caller is the metastore
// admin, all catalogs will be retrieved. Otherwise, only catalogs owned by the
// caller (or for which the caller has the **USE_CATALOG** privilege) will be
// retrieved. There is no guarantee of a specific ordering of the elements in
// the array.
//
// NOTE: we recommend using max_results=0 to use the paginated version of this
// API. Unpaginated calls will be deprecated soon.
//
// PAGINATION BEHAVIOR: When using pagination (max_results >= 0), a page may
// contain zero results while still providing a next_page_token. Clients must
// continue reading pages until next_page_token is absent, which is the only
// indication that the end of results has been reached.
func (c *Client) ListCatalogs(ctx context.Context, req *ListCatalogs, opts ...api.Option) (*ListCatalogs_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/catalogs"
	queryParams := url.Values{}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.IncludeUnbound != nil {
		queryParams.Add("include_unbound", fmt.Sprintf("%v", *req.IncludeUnbound))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListCatalogs_Response{}

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

func (c *Client) ListCatalogsIter(ctx context.Context, req *ListCatalogs, opts ...api.Option) iter.Seq2[*CatalogInfo, error] {
	return func(yield func(*CatalogInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListCatalogs{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListCatalogs(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Catalogs {
				if !yield(&resp.Catalogs[i], nil) {
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

// Updates the catalog that matches the supplied name. The caller must be either
// the owner of the catalog, or a metastore admin (when changing the owner field
// of the catalog).
func (c *Client) UpdateCatalog(ctx context.Context, req *UpdateCatalog, opts ...api.Option) (*CatalogInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/catalogs/%v", *req.NameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CatalogInfo{}

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
