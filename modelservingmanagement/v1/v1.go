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

// Create a new serving endpoint.
func (c *Client) CreateInferenceEndpoint(ctx context.Context, req *CreateInferenceEndpoint, opts ...api.Option) (*InferenceEndpointDetailed, error) {
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
	baseURL.Path = "/api/2.0/preview/serving-endpoints"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &InferenceEndpointDetailed{}

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

// Create a new PT serving endpoint.
func (c *Client) CreateProvisionedThroughputInferenceEndpoint(ctx context.Context, req *CreatePtEndpoint, opts ...api.Option) (*InferenceEndpointDetailed, error) {
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
	baseURL.Path = "/api/2.0/serving-endpoints/pt"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &InferenceEndpointDetailed{}

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

// Delete a serving endpoint.
func (c *Client) DeleteInferenceEndpoint(ctx context.Context, req *DeleteInferenceEndpoint, opts ...api.Option) (*DeleteInferenceEndpoint_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteInferenceEndpoint_Response{}

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

// Retrieves the details for a single serving endpoint.
func (c *Client) GetInferenceEndpoint(ctx context.Context, req *GetInferenceEndpoint, opts ...api.Option) (*InferenceEndpointDetailed, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &InferenceEndpointDetailed{}

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

// Get the query schema of the serving endpoint in OpenAPI format. The schema
// contains information for the supported paths, input and output format and
// datatypes.
func (c *Client) GetInferenceEndpointSchema(ctx context.Context, req *GetInferenceEndpointSchema, opts ...api.Option) (*GetOpenApiResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/openapi", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetOpenApiResponse{}

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

// List only foundation model endpoints
func (c *Client) ListFoundationModelEndpoints(ctx context.Context, req *ListFoundationModelEndpoints, opts ...api.Option) (*ListFoundationModelEndpoints_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/preview/serving-endpoints:foundation-models"
	queryParams := url.Values{}
	if req.Product != nil {
		queryParams.Add("product", fmt.Sprintf("%v", *req.Product))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListFoundationModelEndpoints_Response{}

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

// Get all serving endpoints.
func (c *Client) ListInferenceEndpoints(ctx context.Context, req *ListInferenceEndpoints, opts ...api.Option) (*ListInferenceEndpoints_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/preview/serving-endpoints"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListInferenceEndpoints_Response{}

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

// Updates the budget policy of a serving endpoint.
func (c *Client) PatchInferenceEndpointBudgetPolicy(ctx context.Context, req *PatchInferenceEndpointBudgetPolicy, opts ...api.Option) (*PatchInferenceEndpointBudgetPolicy_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/budget-policy", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PatchInferenceEndpointBudgetPolicy_Response{}

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

func (c *Client) PatchInferenceEndpointDescription(ctx context.Context, req *PatchInferenceEndpointDescription, opts ...api.Option) (*PatchInferenceEndpointDescription_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/description", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PatchInferenceEndpointDescription_Response{}

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

// Used to batch add and delete tags from a serving endpoint with a single API
// call.
func (c *Client) PatchInferenceEndpointTags(ctx context.Context, req *PatchInferenceEndpointTags, opts ...api.Option) (*PatchInferenceEndpointTags_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/tags", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PatchInferenceEndpointTags_Response{}

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

// Updates the usage policy of a serving endpoint.
func (c *Client) PatchInferenceEndpointUsagePolicy(ctx context.Context, req *PatchInferenceEndpointUsagePolicy, opts ...api.Option) (*PatchInferenceEndpointUsagePolicy_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/usage-policy", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PatchInferenceEndpointUsagePolicy_Response{}

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

// Used to update the AI Gateway of a serving endpoint. NOTE: External model,
// provisioned throughput, and pay-per-token endpoints are fully supported;
// agent endpoints currently only support inference tables.
func (c *Client) PutInferenceEndpointAiGateway(ctx context.Context, req *PutInferenceEndpointAiGateway, opts ...api.Option) (*PutInferenceEndpointAiGateway_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/ai-gateway", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PutInferenceEndpointAiGateway_Response{}

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

// Updates any combination of the serving endpoint's served entities, the
// compute configuration of those served entities, and the endpoint's traffic
// config. An endpoint that already has an update in progress can not be updated
// until the current update completes or fails.
func (c *Client) PutInferenceEndpointConfig(ctx context.Context, req *PutInferenceEndpointConfig, opts ...api.Option) (*InferenceEndpointDetailed, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/config", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &InferenceEndpointDetailed{}

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

// Deprecated: Please use AI Gateway to manage rate limits instead.
func (c *Client) PutInferenceEndpointRateLimits(ctx context.Context, req *PutInferenceEndpointRateLimits, opts ...api.Option) (*PutInferenceEndpointRateLimits_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/preview/serving-endpoints/%v/rate-limits", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PutInferenceEndpointRateLimits_Response{}

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

// Updates any combination of the pt endpoint's served entities, the compute
// configuration of those served entities, and the endpoint's traffic config.
// Updates are instantaneous and endpoint should be updated instantly
func (c *Client) PutProvisionedThroughputInferenceEndpointConfig(ctx context.Context, req *PutPtEndpointConfig, opts ...api.Option) (*InferenceEndpointDetailed, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/serving-endpoints/pt/%v/config", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &InferenceEndpointDetailed{}

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

// Updates the email and webhook notification settings for an endpoint.
func (c *Client) UpdateInferenceEndpointNotifications(ctx context.Context, req *UpdateInferenceEndpointNotifications, opts ...api.Option) (*UpdateInferenceEndpointNotifications_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/serving-endpoints/%v/notifications", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UpdateInferenceEndpointNotifications_Response{}

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

// Make external services call using the credentials stored in UC Connection.
func (c *Client) HttpRequest(ctx context.Context, req *ExternalFunctionRequest, opts ...api.Option) (*ExternalFunctionResponse, error) {
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
	baseURL.Path = "/api/2.0/external-function"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ExternalFunctionResponse{}

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
