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

// Creates a new registered model in Unity Catalog.
//
// File storage for model versions in the registered model will be located in
// the default location which is specified by the parent schema, or the parent
// catalog, or the Metastore.
//
// For registered model creation to succeed, the user must satisfy the following
// conditions: - The caller must be a metastore admin, or be the owner of the
// parent catalog and schema, or have the **USE_CATALOG** privilege on the
// parent catalog and the **USE_SCHEMA** privilege on the parent schema. - The
// caller must have the **CREATE MODEL** or **CREATE FUNCTION** privilege on the
// parent schema.
func (c *Client) CreateRegisteredModel(ctx context.Context, req *CreateRegisteredModel, opts ...api.Option) (*RegisteredModelInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/models"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RegisteredModelInfo{}

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

// Deletes a model version from the specified registered model. Any aliases
// assigned to the model version will also be deleted.
//
// The caller must be a metastore admin or an owner of the parent registered
// model. For the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
func (c *Client) DeleteModelVersion(ctx context.Context, req *DeleteModelVersion, opts ...api.Option) (*DeleteModelVersion_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/versions/%v", *req.FullNameArg, *req.VersionArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteModelVersion_Response{}

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

// Deletes a registered model and all its model versions from the specified
// parent catalog and schema.
//
// The caller must be a metastore admin or an owner of the registered model. For
// the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
func (c *Client) DeleteRegisteredModel(ctx context.Context, req *DeleteRegisteredModel, opts ...api.Option) (*DeleteRegisteredModel_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteRegisteredModel_Response{}

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

// Deletes a registered model alias.
//
// The caller must be a metastore admin or an owner of the registered model. For
// the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
func (c *Client) DeleteRegisteredModelAlias(ctx context.Context, req *DeleteRegisteredModelAlias, opts ...api.Option) (*DeleteRegisteredModelAlias_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/aliases/%v", *req.FullNameArg, *req.AliasArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteRegisteredModelAlias_Response{}

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

// Get a model version.
//
// The caller must be a metastore admin or an owner of (or have the **EXECUTE**
// privilege on) the parent registered model. For the latter case, the caller
// must also be the owner or have the **USE_CATALOG** privilege on the parent
// catalog and the **USE_SCHEMA** privilege on the parent schema.
func (c *Client) GetModelVersion(ctx context.Context, req *GetModelVersion, opts ...api.Option) (*ModelVersionInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/versions/%v", *req.FullNameArg, *req.VersionArg)
	queryParams := url.Values{}
	if req.IncludeAliases != nil {
		queryParams.Add("include_aliases", fmt.Sprintf("%v", *req.IncludeAliases))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ModelVersionInfo{}

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

// Get a model version by alias.
//
// The caller must be a metastore admin or an owner of (or have the **EXECUTE**
// privilege on) the registered model. For the latter case, the caller must also
// be the owner or have the **USE_CATALOG** privilege on the parent catalog and
// the **USE_SCHEMA** privilege on the parent schema.
func (c *Client) GetModelVersionByAlias(ctx context.Context, req *GetModelVersionByAlias, opts ...api.Option) (*ModelVersionInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/aliases/%v", *req.FullNameArg, *req.AliasArg)
	queryParams := url.Values{}
	if req.IncludeAliases != nil {
		queryParams.Add("include_aliases", fmt.Sprintf("%v", *req.IncludeAliases))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ModelVersionInfo{}

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

// Get a registered model.
//
// The caller must be a metastore admin or an owner of (or have the **EXECUTE**
// privilege on) the registered model. For the latter case, the caller must also
// be the owner or have the **USE_CATALOG** privilege on the parent catalog and
// the **USE_SCHEMA** privilege on the parent schema.
func (c *Client) GetRegisteredModel(ctx context.Context, req *GetRegisteredModel, opts ...api.Option) (*RegisteredModelInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.IncludeAliases != nil {
		queryParams.Add("include_aliases", fmt.Sprintf("%v", *req.IncludeAliases))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RegisteredModelInfo{}

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

// List model versions. You can list model versions under a particular schema,
// or list all model versions in the current metastore.
//
// The returned models are filtered based on the privileges of the calling user.
// For example, the metastore admin is able to list all the model versions. A
// regular user needs to be the owner or have the **EXECUTE** privilege on the
// parent registered model to recieve the model versions in the response. For
// the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
//
// There is no guarantee of a specific ordering of the elements in the response.
// The elements in the response will not contain any aliases or tags.
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListModelVersions(ctx context.Context, req *ListModelVersions, opts ...api.Option) (*ListModelVersions_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/versions", *req.FullNameArg)
	queryParams := url.Values{}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListModelVersions_Response{}

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

func (c *Client) ListModelVersionsIter(ctx context.Context, req *ListModelVersions, opts ...api.Option) iter.Seq2[*ModelVersionInfo, error] {
	return func(yield func(*ModelVersionInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListModelVersions{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListModelVersions(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.ModelVersions {
				if !yield(&resp.ModelVersions[i], nil) {
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

// List registered models. You can list registered models under a particular
// schema, or list all registered models in the current metastore.
//
// The returned models are filtered based on the privileges of the calling user.
// For example, the metastore admin is able to list all the registered models. A
// regular user needs to be the owner or have the **EXECUTE** privilege on the
// registered model to recieve the registered models in the response. For the
// latter case, the caller must also be the owner or have the **USE_CATALOG**
// privilege on the parent catalog and the **USE_SCHEMA** privilege on the
// parent schema.
//
// There is no guarantee of a specific ordering of the elements in the response.
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListRegisteredModels(ctx context.Context, req *ListRegisteredModels, opts ...api.Option) (*ListRegisteredModels_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/models"
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

	resp := &ListRegisteredModels_Response{}

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

func (c *Client) ListRegisteredModelsIter(ctx context.Context, req *ListRegisteredModels, opts ...api.Option) iter.Seq2[*RegisteredModelInfo, error] {
	return func(yield func(*RegisteredModelInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListRegisteredModels{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListRegisteredModels(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.RegisteredModels {
				if !yield(&resp.RegisteredModels[i], nil) {
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

// Set an alias on the specified registered model.
//
// The caller must be a metastore admin or an owner of the registered model. For
// the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
func (c *Client) SetRegisteredModelAlias(ctx context.Context, req *SetRegisteredModelAlias, opts ...api.Option) (*RegisteredModelAliasInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/aliases/%v", *req.FullNameArg, *req.AliasArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RegisteredModelAliasInfo{}

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

// Updates the specified model version.
//
// The caller must be a metastore admin or an owner of the parent registered
// model. For the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
//
// Currently only the comment of the model version can be updated.
func (c *Client) UpdateModelVersion(ctx context.Context, req *UpdateModelVersion, opts ...api.Option) (*ModelVersionInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v/versions/%v", *req.FullNameArg, *req.VersionArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ModelVersionInfo{}

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

// Updates the specified registered model.
//
// The caller must be a metastore admin or an owner of the registered model. For
// the latter case, the caller must also be the owner or have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
//
// Currently only the name, the owner or the comment of the registered model can
// be updated.
func (c *Client) UpdateRegisteredModel(ctx context.Context, req *UpdateRegisteredModel, opts ...api.Option) (*RegisteredModelInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/models/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RegisteredModelInfo{}

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
