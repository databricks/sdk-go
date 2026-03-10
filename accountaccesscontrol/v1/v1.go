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

// Gets all the roles that can be granted on an account level resource. A role
// is grantable if the rule set on the resource can contain an access rule of
// the role.
func (c *Client) GetAssignableRolesForResource(ctx context.Context, req *GetAssignableRolesForResourceRequest, opts ...api.Option) (*GetAssignableRolesForResourceResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/preview/accounts/{account_id}/access-control/assignable-roles"
	queryParams := url.Values{}
	if req.Resource != nil {
		queryParams.Add("resource", fmt.Sprintf("%v", *req.Resource))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetAssignableRolesForResourceResponse{}

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

// Get a rule set by its name. A rule set is always attached to a resource and
// contains a list of access rules on the said resource. Currently only a
// default rule set for each resource is supported.
func (c *Client) GetRuleSet(ctx context.Context, req *GetRuleSetRequest, opts ...api.Option) (*RuleSet, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/preview/accounts/{account_id}/access-control/rule-sets"
	queryParams := url.Values{}
	if req.Name != nil {
		queryParams.Add("name", fmt.Sprintf("%v", *req.Name))
	}
	if req.Etag != nil {
		queryParams.Add("etag", fmt.Sprintf("%v", *req.Etag))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RuleSet{}

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

// Replace the rules of a rule set. First, use get to read the current version
// of the rule set before modifying it. This pattern helps prevent conflicts
// between concurrent updates.
func (c *Client) UpdateRuleSet(ctx context.Context, req *UpdateRuleSetRequest, opts ...api.Option) (*RuleSet, error) {
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
	baseURL.Path = "/api/2.0/preview/accounts/{account_id}/access-control/rule-sets"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &RuleSet{}

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
