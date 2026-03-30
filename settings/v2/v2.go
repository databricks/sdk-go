package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/databricks/sdk-go/core/ops"
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

// Get a setting value at account level. See
// :method:settingsv2/listaccountsettingsmetadata for list of setting available
// via public APIs at account level.
func (c *Client) GetPublicAccountSetting(ctx context.Context, req *GetPublicAccountSettingRequest, opts ...ops.Option) (*Setting, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/accounts//settings/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Setting{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get a user preference for a specific user. User preferences are personal
// settings that allow individual customization without affecting other users.
// See :method:settingsv2/listaccountuserpreferencesmetadata for list of user
// preferences available via public APIs.
func (c *Client) GetPublicAccountUserPreference(ctx context.Context, req *GetPublicAccountUserPreferenceRequest, opts ...ops.Option) (*UserPreference, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/accounts//users/%v/settings/%v", *req.UserId, *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UserPreference{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Get a setting value at workspace level. See
// :method:settingsv2/listworkspacesettingsmetadata for list of setting
// available via public APIs.
func (c *Client) GetPublicWorkspaceSetting(ctx context.Context, req *GetPublicWorkspaceSettingRequest, opts ...ops.Option) (*Setting, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/settings/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Setting{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// List valid setting keys and metadata. These settings are available to be
// referenced via GET :method:settingsv2/getpublicaccountsetting and PATCH
// :method:settingsv2/patchpublicaccountsetting APIs
func (c *Client) ListAccountSettingsMetadata(ctx context.Context, req *ListAccountSettingsMetadataRequest, opts ...ops.Option) (*ListAccountSettingsMetadataResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/accounts/{account_id}/settings-metadata"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListAccountSettingsMetadataResponse{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListAccountSettingsMetadataIter(ctx context.Context, req *ListAccountSettingsMetadataRequest, opts ...ops.Option) iter.Seq2[*SettingsMetadata, error] {
	return func(yield func(*SettingsMetadata, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListAccountSettingsMetadataRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListAccountSettingsMetadata(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.SettingsMetadata {
				if !yield(&resp.SettingsMetadata[i], nil) {
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

// List valid user preferences and their metadata for a specific user. User
// preferences are personal settings that allow individual customization without
// affecting other users. These settings are available to be referenced via GET
// :method:settingsv2/getpublicaccountuserpreference and PATCH
// :method:settingsv2/patchpublicaccountuserpreference APIs
func (c *Client) ListAccountUserPreferencesMetadata(ctx context.Context, req *ListAccountUserPreferencesMetadataRequest, opts ...ops.Option) (*ListAccountUserPreferencesMetadataResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/accounts//users/%v/settings-metadata", *req.UserId)
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListAccountUserPreferencesMetadataResponse{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListAccountUserPreferencesMetadataIter(ctx context.Context, req *ListAccountUserPreferencesMetadataRequest, opts ...ops.Option) iter.Seq2[*SettingsMetadata, error] {
	return func(yield func(*SettingsMetadata, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListAccountUserPreferencesMetadataRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListAccountUserPreferencesMetadata(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.SettingsMetadata {
				if !yield(&resp.SettingsMetadata[i], nil) {
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

// List valid setting keys and metadata. These settings are available to be
// referenced via GET :method:settingsv2/getpublicworkspacesetting and PATCH
// :method:settingsv2/patchpublicworkspacesetting APIs
func (c *Client) ListWorkspaceSettingsMetadata(ctx context.Context, req *ListWorkspaceSettingsMetadataRequest, opts ...ops.Option) (*ListWorkspaceSettingsMetadataResponse, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/settings-metadata"
	queryParams := url.Values{}
	if req.PageSize != nil {
		queryParams.Add("page_size", fmt.Sprintf("%v", *req.PageSize))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListWorkspaceSettingsMetadataResponse{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListWorkspaceSettingsMetadataIter(ctx context.Context, req *ListWorkspaceSettingsMetadataRequest, opts ...ops.Option) iter.Seq2[*SettingsMetadata, error] {
	return func(yield func(*SettingsMetadata, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListWorkspaceSettingsMetadataRequest{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListWorkspaceSettingsMetadata(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.SettingsMetadata {
				if !yield(&resp.SettingsMetadata[i], nil) {
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

// Patch a setting value at account level. See
// :method:settingsv2/listaccountsettingsmetadata for list of setting available
// via public APIs at account level. To determine the correct field to include
// in a patch request, refer to the type field of the setting returned in the
// :method:settingsv2/listaccountsettingsmetadata response.
//
// Note: Page refresh is required for changes to take effect in UI.
func (c *Client) PatchPublicAccountSetting(ctx context.Context, req *PatchPublicAccountSettingRequest, opts ...ops.Option) (*Setting, error) {
	body, err := json.Marshal(req.Setting)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/accounts//settings/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Setting{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Update a user preference for a specific user. User preferences are personal
// settings that allow individual customization without affecting other users.
// See :method:settingsv2/listaccountuserpreferencesmetadata for list of user
// preferences available via public APIs.
//
// Note: Page refresh is required for changes to take effect in UI.
func (c *Client) PatchPublicAccountUserPreference(ctx context.Context, req *PatchPublicAccountUserPreferenceRequest, opts ...ops.Option) (*UserPreference, error) {
	body, err := json.Marshal(req.Setting)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/accounts//users/%v/settings/%v", *req.UserId, *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UserPreference{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

// Patch a setting value at workspace level. See
// :method:settingsv2/listworkspacesettingsmetadata for list of setting
// available via public APIs at workspace level. To determine the correct field
// to include in a patch request, refer to the type field of the setting
// returned in the :method:settingsv2/listworkspacesettingsmetadata response.
//
// Note: Page refresh is required for changes to take effect in UI.
func (c *Client) PatchPublicWorkspaceSetting(ctx context.Context, req *PatchPublicWorkspaceSettingRequest, opts ...ops.Option) (*Setting, error) {
	body, err := json.Marshal(req.Setting)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/settings/%v", *req.Name)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Setting{}

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

	if err := ops.Execute(ctx, call, opts...); err != nil {
		return nil, err
	}
	return resp, nil
}
