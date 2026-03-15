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

// Creates access requests for Unity Catalog permissions for a specified
// principal on a securable object. This Batch API can take in multiple
// principals, securable objects, and permissions as the input and returns the
// access request destinations for each. Principals must be unique across the
// API call.
//
// The supported securable types are: "metastore", "catalog", "schema", "table",
// "external_location", "connection", "credential", "function",
// "registered_model", and "volume".
func (c *Client) BatchCreateAccessRequests(ctx context.Context, req *BatchCreateAccessRequestsRequest, opts ...api.Option) (*BatchCreateAccessRequestsResponse, error) {
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
	baseURL.Path = "/api/3.0/rfa/requests"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &BatchCreateAccessRequestsResponse{}

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

// Gets an array of access request destinations for the specified securable. Any
// caller can see URL destinations or the destinations on the metastore.
// Otherwise, only those with **BROWSE** permissions on the securable can see
// destinations.
//
// The supported securable types are: "metastore", "catalog", "schema", "table",
// "external_location", "connection", "credential", "function",
// "registered_model", and "volume".
func (c *Client) GetAccessRequestDestinations(ctx context.Context, req *GetAccessRequestDestinationsRequest, opts ...api.Option) (*AccessRequestDestinations, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/3.0/rfa/destinations/%v/%v", *req.SecurableType, *req.FullName)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AccessRequestDestinations{}

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

// Updates the access request destinations for the given securable. The caller
// must be a metastore admin, the owner of the securable, or a user that has the
// **MANAGE** privilege on the securable in order to assign destinations.
// Destinations cannot be updated for securables underneath schemas (tables,
// volumes, functions, and models). For these securable types, destinations are
// inherited from the parent securable. A maximum of 5 emails and 5 external
// notification destinations (Slack, Microsoft Teams, and Generic Webhook
// destinations) can be assigned to a securable. If a URL destination is
// assigned, no other destinations can be set.
//
// The supported securable types are: "metastore", "catalog", "schema", "table",
// "external_location", "connection", "credential", "function",
// "registered_model", and "volume".
func (c *Client) UpdateAccessRequestDestinations(ctx context.Context, req *UpdateAccessRequestDestinationsRequest, opts ...api.Option) (*AccessRequestDestinations, error) {
	body, err := json.Marshal(req.AccessRequestDestinations)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/3.0/rfa/destinations"
	queryParams := url.Values{}
	if req.UpdateMask != nil {
		queryParams.Add("update_mask", fmt.Sprintf("%v", *req.UpdateMask))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AccessRequestDestinations{}

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
