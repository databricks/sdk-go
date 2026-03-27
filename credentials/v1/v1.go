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

// Creates a new credential. The type of credential to be created is determined
// by the **purpose** field, which should be either **SERVICE** or **STORAGE**.
//
// The caller must be a metastore admin or have the metastore privilege
// **CREATE_STORAGE_CREDENTIAL** for storage credentials, or
// **CREATE_SERVICE_CREDENTIAL** for service credentials.
func (c *Client) CreateCredential(ctx context.Context, req *CreateCredential, opts ...api.Option) (*StorageCredentialInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Creates a new storage credential.
//
// The caller must be a metastore admin or have the
// **CREATE_STORAGE_CREDENTIAL** privilege on the metastore.
func (c *Client) CreateStorageCredential(ctx context.Context, req *CreateStorageCredential, opts ...api.Option) (*StorageCredentialInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/storage-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Deletes a service or storage credential from the metastore. The caller must
// be an owner of the credential.
func (c *Client) DeleteCredential(ctx context.Context, req *DeleteCredential, opts ...api.Option) (*DeleteCredential_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	if req.Force != nil {
		queryParams.Add("force", fmt.Sprintf("%v", *req.Force))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteCredential_Response{}

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

// Deletes a storage credential from the metastore. The caller must be an owner
// of the storage credential.
func (c *Client) DeleteStorageCredential(ctx context.Context, req *DeleteStorageCredential, opts ...api.Option) (*DeleteStorageCredential_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/storage-credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	if req.Force != nil {
		queryParams.Add("force", fmt.Sprintf("%v", *req.Force))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteStorageCredential_Response{}

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

// Get a short-lived credential for directly accessing cloud storage locations
// registered in <Databricks>. The Generate Temporary Path Credentials API is
// only supported for external storage paths, specifically external locations
// and external tables. Managed tables are not supported by this API. The
// metastore must have **external_access_enabled** flag set to true (default
// false). The caller must have the **EXTERNAL_USE_LOCATION** privilege on the
// external location; this privilege can only be granted by external location
// owners. For requests on existing external tables, the caller must also have
// the **EXTERNAL_USE_SCHEMA** privilege on the parent schema; this privilege
// can only be granted by catalog owners.
func (c *Client) GenerateTemporaryPathCredential(ctx context.Context, req *GenerateTemporaryPathCredential, opts ...api.Option) (*GenerateTemporaryPathCredential_Response, error) {
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
	baseURL.Path = "/api/2.0/unity-catalog/temporary-path-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GenerateTemporaryPathCredential_Response{}

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

// Returns a set of temporary credentials generated using the specified service
// credential. The caller must be a metastore admin or have the metastore
// privilege **ACCESS** on the service credential.
func (c *Client) GenerateTemporaryServiceCredential(ctx context.Context, req *GenerateTemporaryServiceCredential, opts ...api.Option) (*TemporaryCredentials, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/temporary-service-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TemporaryCredentials{}

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

// Get a short-lived credential for directly accessing the table data on cloud
// storage. The metastore must have **external_access_enabled** flag set to true
// (default false). The caller must have the **EXTERNAL_USE_SCHEMA** privilege
// on the parent schema and this privilege can only be granted by catalog
// owners.
func (c *Client) GenerateTemporaryTableCredential(ctx context.Context, req *GenerateTemporaryTableCredential, opts ...api.Option) (*GenerateTemporaryTableCredential_Response, error) {
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
	baseURL.Path = "/api/2.0/unity-catalog/temporary-table-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GenerateTemporaryTableCredential_Response{}

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

// Gets a service or storage credential from the metastore. The caller must be a
// metastore admin, the owner of the credential, or have any permission on the
// credential.
func (c *Client) GetCredential(ctx context.Context, req *GetCredential, opts ...api.Option) (*StorageCredentialInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Gets a storage credential from the metastore. The caller must be a metastore
// admin, the owner of the storage credential, or have some permission on the
// storage credential.
func (c *Client) GetStorageCredential(ctx context.Context, req *GetStorageCredential, opts ...api.Option) (*StorageCredentialInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/storage-credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Gets an array of credentials (as __CredentialInfo__ objects).
//
// The array is limited to only the credentials that the caller has permission
// to access. If the caller is a metastore admin, retrieval of credentials is
// unrestricted. There is no guarantee of a specific ordering of the elements in
// the array.
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListCredentials(ctx context.Context, req *ListCredentials, opts ...api.Option) (*ListCredentials_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/credentials"
	queryParams := url.Values{}
	if req.IncludeUnbound != nil {
		queryParams.Add("include_unbound", fmt.Sprintf("%v", *req.IncludeUnbound))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListCredentials_Response{}

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

func (c *Client) ListCredentialsIter(ctx context.Context, req *ListCredentials, opts ...api.Option) iter.Seq2[*CredentialInfo, error] {
	return func(yield func(*CredentialInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListCredentials{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListCredentials(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Credentials {
				if !yield(&resp.Credentials[i], nil) {
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

// Gets an array of storage credentials (as __StorageCredentialInfo__ objects).
// The array is limited to only those storage credentials the caller has
// permission to access. If the caller is a metastore admin, retrieval of
// credentials is unrestricted. There is no guarantee of a specific ordering of
// the elements in the array.
//
// NOTE: we recommend using max_results=0 to use the paginated version of this
// API. Unpaginated calls will be deprecated soon.
//
// PAGINATION BEHAVIOR: When using pagination (max_results >= 0), a page may
// contain zero results while still providing a next_page_token. Clients must
// continue reading pages until next_page_token is absent, which is the only
// indication that the end of results has been reached.
func (c *Client) ListStorageCredentials(ctx context.Context, req *ListStorageCredentials, opts ...api.Option) (*ListStorageCredentials_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/storage-credentials"
	queryParams := url.Values{}
	if req.IncludeUnbound != nil {
		queryParams.Add("include_unbound", fmt.Sprintf("%v", *req.IncludeUnbound))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListStorageCredentials_Response{}

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

func (c *Client) ListStorageCredentialsIter(ctx context.Context, req *ListStorageCredentials, opts ...api.Option) iter.Seq2[*StorageCredentialInfo, error] {
	return func(yield func(*StorageCredentialInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListStorageCredentials{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListStorageCredentials(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.StorageCredentials {
				if !yield(&resp.StorageCredentials[i], nil) {
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

// Updates a service or storage credential on the metastore.
//
// The caller must be the owner of the credential or a metastore admin or have
// the `MANAGE` permission. If the caller is a metastore admin, only the
// __owner__ field can be changed.
func (c *Client) UpdateCredential(ctx context.Context, req *UpdateCredential, opts ...api.Option) (*StorageCredentialInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Updates a storage credential on the metastore.
//
// The caller must be the owner of the storage credential or a metastore admin.
// If the caller is a metastore admin, only the **owner** field can be changed.
func (c *Client) UpdateStorageCredential(ctx context.Context, req *UpdateStorageCredential, opts ...api.Option) (*StorageCredentialInfo, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/storage-credentials/%v", *req.NameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &StorageCredentialInfo{}

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

// Validates a credential.
//
// For service credentials (purpose is **SERVICE**), either the
// __credential_name__ or the cloud-specific credential must be provided.
//
// For storage credentials (purpose is **STORAGE**), at least one of
// __external_location_name__ and __url__ need to be provided. If only one of
// them is provided, it will be used for validation. And if both are provided,
// the __url__ will be used for validation, and __external_location_name__ will
// be ignored when checking overlapping urls. Either the __credential_name__ or
// the cloud-specific credential must be provided.
//
// The caller must be a metastore admin or the credential owner or have the
// required permission on the metastore and the credential (e.g.,
// **CREATE_EXTERNAL_LOCATION** when purpose is **STORAGE**).
func (c *Client) ValidateCredential(ctx context.Context, req *ValidateCredential, opts ...api.Option) (*ValidateCredential_Response, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/validate-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ValidateCredential_Response{}

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

// Validates a storage credential. At least one of __external_location_name__
// and __url__ need to be provided. If only one of them is provided, it will be
// used for validation. And if both are provided, the __url__ will be used for
// validation, and __external_location_name__ will be ignored when checking
// overlapping urls.
//
// Either the __storage_credential_name__ or the cloud-specific credential must
// be provided.
//
// The caller must be a metastore admin or the storage credential owner or have
// the **CREATE_EXTERNAL_LOCATION** privilege on the metastore and the storage
// credential.
func (c *Client) ValidateStorageCredential(ctx context.Context, req *ValidateStorageCredential, opts ...api.Option) (*ValidateStorageCredential_Response, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/validate-storage-credentials"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ValidateStorageCredential_Response{}

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
