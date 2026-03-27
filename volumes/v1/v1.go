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

// Creates a new volume.
//
// The user could create either an external volume or a managed volume. An
// external volume will be created in the specified external location, while a
// managed volume will be located in the default location which is specified by
// the parent schema, or the parent catalog, or the Metastore.
//
// For the volume creation to succeed, the user must satisfy following
// conditions: - The caller must be a metastore admin, or be the owner of the
// parent catalog and schema, or have the **USE_CATALOG** privilege on the
// parent catalog and the **USE_SCHEMA** privilege on the parent schema. - The
// caller must have **CREATE VOLUME** privilege on the parent schema.
//
// For an external volume, following conditions also need to satisfy - The
// caller must have **CREATE EXTERNAL VOLUME** privilege on the external
// location. - There are no other tables, nor volumes existing in the specified
// storage location. - The specified storage location is not under the location
// of other tables, nor volumes, or catalogs or schemas.
func (c *Client) CreateVolume(ctx context.Context, req *CreateVolume, opts ...api.Option) (*VolumeInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/volumes"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &VolumeInfo{}

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

// Deletes a volume from the specified parent catalog and schema.
//
// The caller must be a metastore admin or an owner of the volume. For the
// latter case, the caller must also be the owner or have the **USE_CATALOG**
// privilege on the parent catalog and the **USE_SCHEMA** privilege on the
// parent schema.
func (c *Client) DeleteVolume(ctx context.Context, req *DeleteVolume, opts ...api.Option) (*DeleteVolume_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/volumes/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteVolume_Response{}

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

// Gets a volume from the metastore for a specific catalog and schema.
//
// The caller must be a metastore admin or an owner of (or have the **READ
// VOLUME** privilege on) the volume. For the latter case, the caller must also
// be the owner or have the **USE_CATALOG** privilege on the parent catalog and
// the **USE_SCHEMA** privilege on the parent schema.
func (c *Client) GetVolume(ctx context.Context, req *GetVolume, opts ...api.Option) (*VolumeInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/volumes/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &VolumeInfo{}

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

// Gets an array of volumes for the current metastore under the parent catalog
// and schema.
//
// The returned volumes are filtered based on the privileges of the calling
// user. For example, the metastore admin is able to list all the volumes. A
// regular user needs to be the owner or have the **READ VOLUME** privilege on
// the volume to receive the volumes in the response. For the latter case, the
// caller must also be the owner or have the **USE_CATALOG** privilege on the
// parent catalog and the **USE_SCHEMA** privilege on the parent schema.
//
// There is no guarantee of a specific ordering of the elements in the array.
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListVolumes(ctx context.Context, req *ListVolumes, opts ...api.Option) (*ListVolumes_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/volumes"
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

	resp := &ListVolumes_Response{}

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

func (c *Client) ListVolumesIter(ctx context.Context, req *ListVolumes, opts ...api.Option) iter.Seq2[*VolumeInfo, error] {
	return func(yield func(*VolumeInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListVolumes{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListVolumes(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Volumes {
				if !yield(&resp.Volumes[i], nil) {
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

// Updates the specified volume under the specified parent catalog and schema.
//
// The caller must be a metastore admin or an owner of the volume. For the
// latter case, the caller must also be the owner or have the **USE_CATALOG**
// privilege on the parent catalog and the **USE_SCHEMA** privilege on the
// parent schema.
//
// Currently only the name, the owner or the comment of the volume could be
// updated.
func (c *Client) UpdateVolume(ctx context.Context, req *UpdateVolume, opts ...api.Option) (*VolumeInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/volumes/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &VolumeInfo{}

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
