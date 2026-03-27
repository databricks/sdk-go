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

// Creates a tag assignment for an Unity Catalog entity.
//
// To add tags to Unity Catalog entities, you must own the entity or have the
// following privileges: - **APPLY TAG** on the entity - **USE SCHEMA** on the
// entity's parent schema - **USE CATALOG** on the entity's parent catalog
//
// To add a governed tag to Unity Catalog entities, you must also have the
// **ASSIGN** or **MANAGE** permission on the tag policy. See [Manage tag policy
// permissions].
//
// [Manage tag policy permissions]: https://docs.databricks.com/aws/en/admin/tag-policies/manage-permissions
func (c *Client) CreateEntityTagAssignment(ctx context.Context, req *CreateEntityTagAssignmentRequest, opts ...api.Option) (*EntityTagAssignment, error) {
	body, err := json.Marshal(req.TagAssignment)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/entity-tag-assignments"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &EntityTagAssignment{}

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

// Deletes a tag assignment for an Unity Catalog entity by its key.
//
// To delete tags from Unity Catalog entities, you must own the entity or have
// the following privileges: - **APPLY TAG** on the entity - **USE_SCHEMA** on
// the entity's parent schema - **USE_CATALOG** on the entity's parent catalog
//
// To delete a governed tag from Unity Catalog entities, you must also have the
// **ASSIGN** or **MANAGE** permission on the tag policy. See [Manage tag policy
// permissions].
//
// [Manage tag policy permissions]: https://docs.databricks.com/aws/en/admin/tag-policies/manage-permissions
func (c *Client) DeleteEntityTagAssignment(ctx context.Context, req *DeleteEntityTagAssignmentRequest, opts ...api.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/entity-tag-assignments/%v/%v/tags/%v", *req.EntityType, *req.EntityName, *req.TagKey)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

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
		_ = respBody
		return nil
	}

	if err := api.Execute(ctx, call, opts...); err != nil {
		return err
	}
	return nil
}

// Gets a tag assignment for an Unity Catalog entity by tag key.
func (c *Client) GetEntityTagAssignment(ctx context.Context, req *GetEntityTagAssignmentRequest, opts ...api.Option) (*EntityTagAssignment, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/entity-tag-assignments/%v/%v/tags/%v", *req.EntityType, *req.EntityName, *req.TagKey)
	queryParams := url.Values{}
	if req.IncludeInherited != nil {
		queryParams.Add("include_inherited", fmt.Sprintf("%v", *req.IncludeInherited))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &EntityTagAssignment{}

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

// List tag assignments for an Unity Catalog entity
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListEntityTagAssignments(ctx context.Context, req *ListEntityTagAssignmentsRequest, opts ...api.Option) (*ListEntityTagAssignmentsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/entity-tag-assignments/%v/%v/tags", *req.EntityType, *req.EntityName)
	queryParams := url.Values{}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.IncludeInherited != nil {
		queryParams.Add("include_inherited", fmt.Sprintf("%v", *req.IncludeInherited))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListEntityTagAssignmentsResponse{}

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

func (c *Client) ListEntityTagAssignmentsIter(ctx context.Context, req *ListEntityTagAssignmentsRequest, opts ...api.Option) iter.Seq2[*EntityTagAssignment, error] {
	return func(yield func(*EntityTagAssignment, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListEntityTagAssignmentsRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListEntityTagAssignments(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.TagAssignments {
				if !yield(&resp.TagAssignments[i], nil) {
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

// Updates an existing tag assignment for an Unity Catalog entity.
//
// To update tags to Unity Catalog entities, you must own the entity or have the
// following privileges: - **APPLY TAG** on the entity - **USE SCHEMA** on the
// entity's parent schema - **USE CATALOG** on the entity's parent catalog
//
// To update a governed tag to Unity Catalog entities, you must also have the
// **ASSIGN** or **MANAGE** permission on the tag policy. See [Manage tag policy
// permissions].
//
// [Manage tag policy permissions]: https://docs.databricks.com/aws/en/admin/tag-policies/manage-permissions
func (c *Client) UpdateEntityTagAssignment(ctx context.Context, req *UpdateEntityTagAssignmentRequest, opts ...api.Option) (*EntityTagAssignment, error) {
	body, err := json.Marshal(req.TagAssignment)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/entity-tag-assignments/%v/%v/tags/%v", *req.TagAssignment.EntityType, *req.TagAssignment.EntityName, *req.TagAssignment.TagKey)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &EntityTagAssignment{}

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
