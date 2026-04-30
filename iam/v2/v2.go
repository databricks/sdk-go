package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/databricks/sdk-go/core/ops"
	"github.com/databricks/sdk-go/databricks/transport"
	"github.com/databricks/sdk-go/options/call"
	"github.com/databricks/sdk-go/options/client"
	"github.com/databricks/sdk-go/options/internaloptions"
)

type Client struct {
	httpClient *http.Client
	logger     *slog.Logger
	host       string
}

func NewClient(ctx context.Context, opts ...client.Option) (*Client, error) {
	cfg := internaloptions.ClientOptions{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}
	httpClient, err := transport.NewHTTPClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: httpClient,
		logger:     cfg.Logger,
		host:       cfg.Host,
	}, nil
}

// executeCall resolves call.Option values to ops.Option values and invokes
// ops.Execute. Lives at the package level to keep call sites concise.
func executeCall(ctx context.Context, op func(context.Context) error, opts []call.Option) error {
	cfg := internaloptions.CallOptions{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return err
		}
	}
	var opsOpts []ops.Option
	if cfg.Retrier != nil {
		opsOpts = append(opsOpts, ops.WithRetrier(cfg.Retrier))
	}
	if cfg.RateLimiter != nil {
		opsOpts = append(opsOpts, ops.WithLimiter(cfg.RateLimiter))
	}
	if cfg.Timeout != 0 {
		opsOpts = append(opsOpts, ops.WithTimeout(cfg.Timeout))
	}
	return ops.Execute(ctx, op, opsOpts...)
}

// Creates a new account access identity rule for a given account. This allows
// administrators to explicitly allow or deny specific principals from accessing
// the account.
func (c *Client) CreateAccountAccessIdentityRule(ctx context.Context, req *CreateAccountAccessIdentityRuleRequest, opts ...call.Option) (*AccountAccessIdentityRule, error) {
	body, err := json.Marshal(req.AccountAccessIdentityRule)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/accounts/{account_id}/aim-control-policy/account-access-identity-rules"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AccountAccessIdentityRule{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Deletes an account access identity rule for a given principal.
func (c *Client) DeleteAccountAccessIdentityRule(ctx context.Context, req *DeleteAccountAccessIdentityRuleRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/accounts//aim-control-policy/account-access-identity-rules/%v", *req.ExternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// Gets an account access identity rule for a given principal.
func (c *Client) GetAccountAccessIdentityRule(ctx context.Context, req *GetAccountAccessIdentityRuleRequest, opts ...call.Option) (*AccountAccessIdentityRule, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/accounts//aim-control-policy/account-access-identity-rules/%v", *req.ExternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AccountAccessIdentityRule{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Lists all account access identity rules for a given account. These rules
// control which principals (users, service principals, groups) from the
// customer's IdP are allowed or denied access to the <Databricks> account.
func (c *Client) ListAccountAccessIdentityRules(ctx context.Context, req *ListAccountAccessIdentityRulesRequest, opts ...call.Option) (*ListAccountAccessIdentityRulesResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/accounts/{account_id}/aim-control-policy/account-access-identity-rules"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListAccountAccessIdentityRulesResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateGroup(ctx context.Context, req *CreateGroupRequest, opts ...call.Option) (*Group, error) {
	body, err := json.Marshal(req.Group)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/groups"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateGroupProxy(ctx context.Context, req *CreateGroupProxyRequest, opts ...call.Option) (*Group, error) {
	body, err := json.Marshal(req.Group)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/groups"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteGroup(ctx context.Context, req *DeleteGroupRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//groups/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteGroupProxy(ctx context.Context, req *DeleteGroupProxyRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/groups/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetGroup(ctx context.Context, req *GetGroupRequest, opts ...call.Option) (*Group, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//groups/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetGroupProxy(ctx context.Context, req *GetGroupProxyRequest, opts ...call.Option) (*Group, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/groups/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListGroups(ctx context.Context, req *ListGroupsRequest, opts ...call.Option) (*ListGroupsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/groups"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListGroupsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListGroupsProxy(ctx context.Context, req *ListGroupsProxyRequest, opts ...call.Option) (*ListGroupsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/groups"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListGroupsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves a group with the given external ID from the customer's IdP. If the
// group does not exist, it will be created in the account. If the customer is
// not onboarded onto Automatic Identity Management (AIM), this will return an
// error.
func (c *Client) ResolveGroup(ctx context.Context, req *ResolveGroupRequest, opts ...call.Option) (*ResolveGroupResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/groups/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveGroupResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves a group with the given external ID from the customer's IdP. If the
// group does not exist, it will be created in the account. If the customer is
// not onboarded onto Automatic Identity Management (AIM), this will return an
// error.
func (c *Client) ResolveGroupProxy(ctx context.Context, req *ResolveGroupProxyRequest, opts ...call.Option) (*ResolveGroupResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/groups/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveGroupResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateGroup(ctx context.Context, req *UpdateGroupRequest, opts ...call.Option) (*Group, error) {
	body, err := json.Marshal(req.Group)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//groups/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateGroupProxy(ctx context.Context, req *UpdateGroupProxyRequest, opts ...call.Option) (*Group, error) {
	body, err := json.Marshal(req.Group)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/groups/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Group{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateServicePrincipal(ctx context.Context, req *CreateServicePrincipalRequest, opts ...call.Option) (*ServicePrincipal, error) {
	body, err := json.Marshal(req.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/servicePrincipals"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateServicePrincipalProxy(ctx context.Context, req *CreateServicePrincipalProxyRequest, opts ...call.Option) (*ServicePrincipal, error) {
	body, err := json.Marshal(req.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/servicePrincipals"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteServicePrincipal(ctx context.Context, req *DeleteServicePrincipalRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//servicePrincipals/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteServicePrincipalProxy(ctx context.Context, req *DeleteServicePrincipalProxyRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/servicePrincipals/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetServicePrincipal(ctx context.Context, req *GetServicePrincipalRequest, opts ...call.Option) (*ServicePrincipal, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//servicePrincipals/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetServicePrincipalProxy(ctx context.Context, req *GetServicePrincipalProxyRequest, opts ...call.Option) (*ServicePrincipal, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/servicePrincipals/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListServicePrincipals(ctx context.Context, req *ListServicePrincipalsRequest, opts ...call.Option) (*ListServicePrincipalsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/servicePrincipals"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListServicePrincipalsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListServicePrincipalsProxy(ctx context.Context, req *ListServicePrincipalsProxyRequest, opts ...call.Option) (*ListServicePrincipalsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/servicePrincipals"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListServicePrincipalsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves an SP with the given external ID from the customer's IdP. If the SP
// does not exist, it will be created. If the customer is not onboarded onto
// Automatic Identity Management (AIM), this will return an error.
func (c *Client) ResolveServicePrincipal(ctx context.Context, req *ResolveServicePrincipalRequest, opts ...call.Option) (*ResolveServicePrincipalResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/servicePrincipals/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveServicePrincipalResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves an SP with the given external ID from the customer's IdP. If the SP
// does not exist, it will be created. If the customer is not onboarded onto
// Automatic Identity Management (AIM), this will return an error.
func (c *Client) ResolveServicePrincipalProxy(ctx context.Context, req *ResolveServicePrincipalProxyRequest, opts ...call.Option) (*ResolveServicePrincipalResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/servicePrincipals/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveServicePrincipalResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateServicePrincipal(ctx context.Context, req *UpdateServicePrincipalRequest, opts ...call.Option) (*ServicePrincipal, error) {
	body, err := json.Marshal(req.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//servicePrincipals/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateServicePrincipalProxy(ctx context.Context, req *UpdateServicePrincipalProxyRequest, opts ...call.Option) (*ServicePrincipal, error) {
	body, err := json.Marshal(req.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/servicePrincipals/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ServicePrincipal{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateUser(ctx context.Context, req *CreateUserRequest, opts ...call.Option) (*User, error) {
	body, err := json.Marshal(req.User)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/users"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateUserProxy(ctx context.Context, req *CreateUserProxyRequest, opts ...call.Option) (*User, error) {
	body, err := json.Marshal(req.User)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/users"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteUser(ctx context.Context, req *DeleteUserRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//users/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteUserProxy(ctx context.Context, req *DeleteUserProxyRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/users/%v", *req.InternalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetUser(ctx context.Context, req *GetUserRequest, opts ...call.Option) (*User, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//users/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) GetUserProxy(ctx context.Context, req *GetUserProxyRequest, opts ...call.Option) (*User, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/users/%v", *req.InternalId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListUsers(ctx context.Context, req *ListUsersRequest, opts ...call.Option) (*ListUsersResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/users"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListUsersResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListUsersProxy(ctx context.Context, req *ListUsersProxyRequest, opts ...call.Option) (*ListUsersResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/users"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.Filter != nil {
		queryParams.Add("filter", fmt.Sprintf("%v", *req.Filter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListUsersResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves a user with the given external ID from the customer's IdP. If the
// user does not exist, it will be created. If the customer is not onboarded
// onto Automatic Identity Management (AIM), this will return an error.
func (c *Client) ResolveUser(ctx context.Context, req *ResolveUserRequest, opts ...call.Option) (*ResolveUserResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/accounts/{account_id}/users/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveUserResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resolves a user with the given external ID from the customer's IdP. If the
// user does not exist, it will be created. If the customer is not onboarded
// onto Automatic Identity Management (AIM), this will return an error.
func (c *Client) ResolveUserProxy(ctx context.Context, req *ResolveUserProxyRequest, opts ...call.Option) (*ResolveUserResponse, error) {
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
	baseURL.Path = "/api/2.0/identity/users/resolveByExternalId"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ResolveUserResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateUser(ctx context.Context, req *UpdateUserRequest, opts ...call.Option) (*User, error) {
	body, err := json.Marshal(req.User)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//users/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateUserProxy(ctx context.Context, req *UpdateUserProxyRequest, opts ...call.Option) (*User, error) {
	body, err := json.Marshal(req.User)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/users/%v", *req.InternalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &User{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateWorkspaceAccessDetail(ctx context.Context, req *CreateWorkspaceAccessDetailRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {
	body, err := json.Marshal(req.WorkspaceAccessDetail)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//workspaces/%v/workspaceAccessDetails", *req.Parent)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) CreateWorkspaceAccessDetailLocal(ctx context.Context, req *CreateWorkspaceAccessDetailLocalRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {
	body, err := json.Marshal(req.WorkspaceAccessDetail)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/workspaceAccessDetails"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteWorkspaceAccessDetail(ctx context.Context, req *DeleteWorkspaceAccessDetailRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//workspaces/%v/workspaceAccessDetails/%v", *req.WorkspaceId, *req.PrincipalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// TODO: Write description later when this method is implemented
func (c *Client) DeleteWorkspaceAccessDetailLocal(ctx context.Context, req *DeleteWorkspaceAccessDetailLocalRequest, opts ...call.Option) error {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/workspaceAccessDetails/%v", *req.PrincipalId)
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

	if err := executeCall(ctx, call, opts); err != nil {
		return err
	}
	return nil
}

// Returns the access details for a principal in a workspace. Allows for
// checking access details for any provisioned principal (user, service
// principal, or group) in a workspace. * Provisioned principal here refers to
// one that has been synced into <Databricks> from the customer's IdP or added
// explicitly to <Databricks> via SCIM/UI. Allows for passing in a "view"
// parameter to control what fields are returned (BASIC by default or FULL).
func (c *Client) GetWorkspaceAccessDetail(ctx context.Context, req *GetWorkspaceAccessDetailRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//workspaces/%v/workspaceAccessDetails/%v", *req.WorkspaceId, *req.PrincipalId)
	queryParams := url.Values{}
	if req.View != nil {
		queryParams.Add("view", fmt.Sprintf("%v", *req.View))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// Returns the access details for a principal in the current workspace. Allows
// for checking access details for any provisioned principal (user, service
// principal, or group) in the current workspace. * Provisioned principal here
// refers to one that has been synced into <Databricks> from the customer's IdP
// or added explicitly to <Databricks> via SCIM/UI. Allows for passing in a
// "view" parameter to control what fields are returned (BASIC by default or
// FULL).
func (c *Client) GetWorkspaceAccessDetailLocal(ctx context.Context, req *GetWorkspaceAccessDetailLocalRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/workspaceAccessDetails/%v", *req.PrincipalId)
	queryParams := url.Values{}
	if req.View != nil {
		queryParams.Add("view", fmt.Sprintf("%v", *req.View))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListWorkspaceAccessDetails(ctx context.Context, req *ListWorkspaceAccessDetailsRequest, opts ...call.Option) (*ListWorkspaceAccessDetailsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//workspaces/%v/workspaceAccessDetails", *req.WorkspaceId)
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListWorkspaceAccessDetailsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) ListWorkspaceAccessDetailsLocal(ctx context.Context, req *ListWorkspaceAccessDetailsLocalRequest, opts ...call.Option) (*ListWorkspaceAccessDetailsResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/identity/workspaceAccessDetails"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListWorkspaceAccessDetailsResponse{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateWorkspaceAccessDetail(ctx context.Context, req *UpdateWorkspaceAccessDetailRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {
	body, err := json.Marshal(req.WorkspaceAccessDetail)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/accounts//workspaces/%v/workspaceAccessDetails/%v", *req.WorkspaceId, *req.PrincipalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Write description later when this method is implemented
func (c *Client) UpdateWorkspaceAccessDetailLocal(ctx context.Context, req *UpdateWorkspaceAccessDetailLocalRequest, opts ...call.Option) (*WorkspaceAccessDetail, error) {
	body, err := json.Marshal(req.WorkspaceAccessDetail)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/identity/workspaceAccessDetails/%v", *req.PrincipalId)
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &WorkspaceAccessDetail{}

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

	if err := executeCall(ctx, call, opts); err != nil {
		return nil, err
	}
	return resp, nil
}
