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

// Creates a new tag policy, making the associated tag key governed. For
// Terraform usage, see the [Tag Policy Terraform documentation]. To manage
// permissions for tag policies, use the [Account Access Control Proxy API].
//
// [Account Access Control Proxy API]: https://docs.databricks.com/api/workspace/accountaccesscontrolproxy
// [Tag Policy Terraform documentation]: https://registry.terraform.io/providers/databricks/databricks/latest/docs/resources/tag_policy
func (c *Client) CreateTagPolicy(ctx context.Context, req *CreateTagPolicyRequest, opts ...api.Option) (*TagPolicy, error) {
	body, err := json.Marshal(req.TagPolicy)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/tag-policies"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TagPolicy{}

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

// Deletes a tag policy by its associated governed tag's key, leaving that tag
// key ungoverned. For Terraform usage, see the [Tag Policy Terraform
// documentation].
//
// [Tag Policy Terraform documentation]: https://registry.terraform.io/providers/databricks/databricks/latest/docs/resources/tag_policy
func (c *Client) DeleteTagPolicy(ctx context.Context, req *DeleteTagPolicyRequest, opts ...api.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/tag-policies/%v", *req.TagKey)
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

// Gets a single tag policy by its associated governed tag's key. For Terraform
// usage, see the [Tag Policy Terraform documentation]. To list granted
// permissions for tag policies, use the [Account Access Control Proxy API].
//
// [Account Access Control Proxy API]: https://docs.databricks.com/api/workspace/accountaccesscontrolproxy
// [Tag Policy Terraform documentation]: https://registry.terraform.io/providers/databricks/databricks/latest/docs/data-sources/tag_policy
func (c *Client) GetTagPolicy(ctx context.Context, req *GetTagPolicyRequest, opts ...api.Option) (*TagPolicy, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/tag-policies/%v", *req.TagKey)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TagPolicy{}

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

// Lists the tag policies for all governed tags in the account. For Terraform
// usage, see the [Tag Policy Terraform documentation]. To list granted
// permissions for tag policies, use the [Account Access Control Proxy API].
//
// [Account Access Control Proxy API]: https://docs.databricks.com/api/workspace/accountaccesscontrolproxy
// [Tag Policy Terraform documentation]: https://registry.terraform.io/providers/databricks/databricks/latest/docs/data-sources/tag_policies
func (c *Client) ListTagPolicies(ctx context.Context, req *ListTagPoliciesRequest, opts ...api.Option) (*ListTagPoliciesResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/tag-policies"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListTagPoliciesResponse{}

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

func (c *Client) ListTagPoliciesIter(ctx context.Context, req *ListTagPoliciesRequest, opts ...api.Option) iter.Seq2[*TagPolicy, error] {
	return func(yield func(*TagPolicy, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListTagPoliciesRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListTagPolicies(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.TagPolicies {
				if !yield(&resp.TagPolicies[i], nil) {
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

// Updates an existing tag policy for a single governed tag. For Terraform
// usage, see the [Tag Policy Terraform documentation]. To manage permissions
// for tag policies, use the [Account Access Control Proxy API].
//
// [Account Access Control Proxy API]: https://docs.databricks.com/api/workspace/accountaccesscontrolproxy
// [Tag Policy Terraform documentation]: https://registry.terraform.io/providers/databricks/databricks/latest/docs/resources/tag_policy
func (c *Client) UpdateTagPolicy(ctx context.Context, req *UpdateTagPolicyRequest, opts ...api.Option) (*TagPolicy, error) {
	body, err := json.Marshal(req.TagPolicy)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/tag-policies/%v", *req.TagPolicy.TagKey)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TagPolicy{}

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
