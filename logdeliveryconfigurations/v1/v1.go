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

// Creates a new <Databricks> log delivery configuration to enable delivery of
// the specified type of logs to your storage location. This requires that you
// already created a [credential object](:method:Credentials/Create) (which
// encapsulates a cross-account service IAM role) and a [storage configuration
// object](:method:Storage/Create) (which encapsulates an S3 bucket). For full
// details, including the required IAM role policies and bucket policies, see
// [Deliver and access billable usage logs] or [Configure audit logging].
// **Note**: There is a limit on the number of log delivery configurations
// available per account (each limit applies separately to each log type
// including billable usage and audit logs). You can create a maximum of two
// enabled account-level delivery configurations (configurations without a
// workspace filter) per type. Additionally, you can create two enabled
// workspace-level delivery configurations per workspace for each log type,
// which means that the same workspace ID can occur in the workspace filter for
// no more than two delivery configurations per log type.
//
// You cannot delete a log delivery configuration, but you can disable it (see
// [Enable or disable log delivery
// configuration](:method:LogDelivery/PatchStatus)).
//
// [Configure audit logging]: https://docs.databricks.com/administration-guide/account-settings/audit-logs.html
// [Deliver and access billable usage logs]: https://docs.databricks.com/administration-guide/account-settings/billable-usage-delivery.html
func (c *Client) CreateLogDeliveryConfiguration(ctx context.Context, req *CreateLogDeliveryConfiguration, opts ...api.Option) (*CreateLogDeliveryConfiguration_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/accounts/%v/log-delivery", *req.LogDeliveryConfiguration.AccountId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CreateLogDeliveryConfiguration_Response{}

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

// Gets a <Databricks> log delivery configuration object for an account, both
// specified by ID.
func (c *Client) GetLogDeliveryConfiguration(ctx context.Context, req *GetLogDeliveryConfiguration, opts ...api.Option) (*GetLogDeliveryConfiguration_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.0/accounts//log-delivery/%v", *req.ConfigId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetLogDeliveryConfiguration_Response{}

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

// Gets all <Databricks> log delivery configurations associated with an account
// specified by ID.
func (c *Client) ListLogDeliveryConfiguration(ctx context.Context, req *ListLogDeliveryConfiguration, opts ...api.Option) (*ListLogDeliveryConfiguration_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/accounts/{account_id}/log-delivery"
	queryParams := url.Values{}
	if req.CredentialsId != nil {
		queryParams.Add("credentials_id", fmt.Sprintf("%v", *req.CredentialsId))
	}
	if req.StorageConfigurationId != nil {
		queryParams.Add("storage_configuration_id", fmt.Sprintf("%v", *req.StorageConfigurationId))
	}
	if req.Status != nil {
		queryParams.Add("status", fmt.Sprintf("%v", *req.Status))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListLogDeliveryConfiguration_Response{}

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

func (c *Client) ListLogDeliveryConfigurationIter(ctx context.Context, req *ListLogDeliveryConfiguration, opts ...api.Option) iter.Seq2[*LogDeliveryConfiguration, error] {
	return func(yield func(*LogDeliveryConfiguration, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListLogDeliveryConfiguration{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListLogDeliveryConfiguration(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.LogDeliveryConfigurations {
				if !yield(&resp.LogDeliveryConfigurations[i], nil) {
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

// Enables or disables a log delivery configuration. Deletion of delivery
// configurations is not supported, so disable log delivery configurations that
// are no longer needed. Note that you can't re-enable a delivery configuration
// if this would violate the delivery configuration limits described under
// [Create log delivery](:method:LogDelivery/Create).
func (c *Client) UpdateLogDeliveryConfiguration(ctx context.Context, req *UpdateLogDeliveryConfiguration, opts ...api.Option) (*UpdateLogDeliveryConfiguration_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.0/accounts//log-delivery/%v", *req.ConfigId)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UpdateLogDeliveryConfiguration_Response{}

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
