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

// Creates a new secret scope.
//
// The scope name must consist of alphanumeric characters, dashes, underscores,
// and periods, and may not exceed 128 characters.
//
// Example request:
//
// .. code::
//
// { "scope": "my-simple-databricks-scope", "initial_manage_principal": "users"
// "scope_backend_type": "databricks|azure_keyvault", # below is only required
// if scope type is azure_keyvault "backend_azure_keyvault": { "resource_id":
// "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/xxxx/providers/Microsoft.KeyVault/vaults/xxxx",
// "tenant_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "dns_name":
// "https://xxxx.vault.azure.net/", } }
//
// If “initial_manage_principal“ is specified, the initial ACL applied to the
// scope is applied to the supplied principal (user or group) with “MANAGE“
// permissions. The only supported principal for this option is the group
// “users“, which contains all users in the workspace. If
// “initial_manage_principal“ is not specified, the initial ACL with
// “MANAGE“ permission applied to the scope is assigned to the API request
// issuer's user identity.
//
// If “scope_backend_type“ is “azure_keyvault“, a secret scope is created
// with secrets from a given Azure KeyVault. The caller must provide the
// keyvault_resource_id and the tenant_id for the key vault. If
// “scope_backend_type“ is “databricks“ or is unspecified, an empty secret
// scope is created and stored in <Databricks>'s own storage.
//
// Throws “RESOURCE_ALREADY_EXISTS“ if a scope with the given name already
// exists. Throws “RESOURCE_LIMIT_EXCEEDED“ if maximum number of scopes in the
// workspace is exceeded. Throws “INVALID_PARAMETER_VALUE“ if the scope name
// is invalid. Throws “BAD_REQUEST“ if request violated constraints. Throws
// “CUSTOMER_UNAUTHORIZED“ if normal user attempts to create a scope with name
// reserved for databricks internal usage. Throws “UNAUTHENTICATED“ if unable
// to verify user access permission on Azure KeyVault
func (c *Client) CreateScope(ctx context.Context, req *CreateScope, opts ...api.Option) (*CreateScope_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/scopes/create"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &CreateScope_Response{}

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

// Deletes the given ACL on the given scope.
//
// Users must have the “MANAGE“ permission to invoke this API.
//
// Example request:
//
// .. code::
//
// { "scope": "my-secret-scope", "principal": "data-scientists" }
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope, principal, or ACL
// exists. Throws “PERMISSION_DENIED“ if the user does not have permission to
// make this API call. Throws “INVALID_PARAMETER_VALUE“ if the permission or
// principal is invalid.
func (c *Client) DeleteAcl(ctx context.Context, req *DeleteAcl, opts ...api.Option) (*DeleteAcl_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/acls/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteAcl_Response{}

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

// Deletes a secret scope.
//
// Example request:
//
// .. code::
//
// { "scope": "my-secret-scope" }
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if the scope does not exist. Throws
// “PERMISSION_DENIED“ if the user does not have permission to make this API
// call. Throws “BAD_REQUEST“ if system user attempts to delete internal
// secret scope.
func (c *Client) DeleteScope(ctx context.Context, req *DeleteScope, opts ...api.Option) (*DeleteScope_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/scopes/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteScope_Response{}

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

// Deletes the secret stored in this secret scope. You must have “WRITE“ or
// “MANAGE“ permission on the Secret Scope.
//
// Example request:
//
// .. code::
//
// { "scope": "my-secret-scope", "key": "my-secret-key" }
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope or secret exists.
// Throws “PERMISSION_DENIED“ if the user does not have permission to make
// this API call. Throws “BAD_REQUEST“ if system user attempts to delete an
// internal secret, or request is made against Azure KeyVault backed scope.
func (c *Client) DeleteSecret(ctx context.Context, req *DeleteSecret, opts ...api.Option) (*DeleteSecret_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteSecret_Response{}

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

// Describes the details about the given ACL, such as the group and permission.
//
// Users must have the “MANAGE“ permission to invoke this API.
//
// Example response:
//
// .. code::
//
// { "principal": "data-scientists", "permission": "READ" }
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope exists. Throws
// “PERMISSION_DENIED“ if the user does not have permission to make this API
// call. Throws “INVALID_PARAMETER_VALUE“ if the permission or principal is
// invalid.
func (c *Client) GetAcl(ctx context.Context, req *GetAcl, opts ...api.Option) (*AclItem, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/secrets/acls/get"
	queryParams := url.Values{}
	if req.Scope != nil {
		queryParams.Add("scope", fmt.Sprintf("%v", *req.Scope))
	}
	if req.Principal != nil {
		queryParams.Add("principal", fmt.Sprintf("%v", *req.Principal))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AclItem{}

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

// Gets a secret for a given key and scope. This API can only be called from the
// DBUtils interface. Users need the READ permission to make this call.
//
// Example response:
//
// .. code::
//
// { "key": "my-string-key", "value": <bytes of the secret value> }
//
// Note that the secret value returned is in bytes. The interpretation of the
// bytes is determined by the caller in DBUtils and the type the data is decoded
// into.
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret or secret scope exists.
// Throws “PERMISSION_DENIED“ if the user does not have permission to make
// this API call.
//
// Note: This is explicitly an undocumented API. It also doesn't need to be
// supported for the /preview prefix, because it's not a customer-facing API
// (i.e. only used for DBUtils SecretUtils to fetch secrets).
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope or secret exists.
// Throws “BAD_REQUEST“ if normal user calls get secret outside of a notebook.
// AKV specific errors: Throws “INVALID_PARAMETER_VALUE“ if secret name is not
// alphanumeric or too long. Throws “PERMISSION_DENIED“ if secret manager
// cannot access AKV with 403 error Throws “MALFORMED_REQUEST“ if secret
// manager cannot access AKV with any other 4xx error
func (c *Client) GetSecret(ctx context.Context, req *GetSecret, opts ...api.Option) (*GetSecret_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/secrets/get"
	queryParams := url.Values{}
	if req.Scope != nil {
		queryParams.Add("scope", fmt.Sprintf("%v", *req.Scope))
	}
	if req.Key != nil {
		queryParams.Add("key", fmt.Sprintf("%v", *req.Key))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetSecret_Response{}

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

// Lists the ACLs set on the given scope.
//
// Users must have the “MANAGE“ permission to invoke this API.
//
// Example response:
//
// .. code::
//
// { "acls": [{ "principal": "admins", "permission": "MANAGE" },{ "principal":
// "data-scientists", "permission": "READ" }] }
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope exists. Throws
// “PERMISSION_DENIED“ if the user does not have permission to make this API
// call.
func (c *Client) ListAcls(ctx context.Context, req *ListAcls, opts ...api.Option) (*ListAcls_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/secrets/acls/list"
	queryParams := url.Values{}
	if req.Scope != nil {
		queryParams.Add("scope", fmt.Sprintf("%v", *req.Scope))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListAcls_Response{}

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

// Lists all secret scopes available in the workspace.
//
// Example response:
//
// .. code::
//
// { "scopes": [{ "name": "my-databricks-scope", "backend_type": "DATABRICKS"
// },{ "name": "mount-points", "backend_type": "DATABRICKS" }] }
//
// Throws “PERMISSION_DENIED“ if the user does not have permission to make
// this API call.
func (c *Client) ListScopes(ctx context.Context, req *ListScopes, opts ...api.Option) (*ListScopes_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/secrets/scopes/list"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListScopes_Response{}

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

// Lists the secret keys that are stored at this scope. This is a metadata-only
// operation; secret data cannot be retrieved using this API. Users need the
// READ permission to make this call.
//
// Example response:
//
// .. code::
//
// { "secrets": [ { "key": "my-string-key"", "last_updated_timestamp":
// "1520467595000" }, { "key": "my-byte-key", "last_updated_timestamp":
// "1520467595000" }, ] }
//
// The lastUpdatedTimestamp returned is in milliseconds since epoch.
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope exists. Throws
// “PERMISSION_DENIED“ if the user does not have permission to make this API
// call.
func (c *Client) ListSecrets(ctx context.Context, req *ListSecrets, opts ...api.Option) (*ListSecrets_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/secrets/list"
	queryParams := url.Values{}
	if req.Scope != nil {
		queryParams.Add("scope", fmt.Sprintf("%v", *req.Scope))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListSecrets_Response{}

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

// Creates or overwrites the ACL associated with the given principal (user or
// group) on the specified scope point. In general, a user or group will use the
// most powerful permission available to them, and permissions are ordered as
// follows:
//
// * “MANAGE“ - Allowed to change ACLs, and read and write to this secret
// scope. * “WRITE“ - Allowed to read and write to this secret scope. *
// “READ“ - Allowed to read this secret scope and list what secrets are
// available.
//
// Note that in general, secret values can only be read from within a command on
// a cluster (for example, through a notebook). There is no API to read the
// actual secret value material outside of a cluster. However, the user's
// permission will be applied based on who is executing the command, and they
// must have at least READ permission.
//
// Users must have the “MANAGE“ permission to invoke this API.
//
// Example request:
//
// .. code::
//
// { "scope": "my-secret-scope", "principal": "data-scientists", "permission":
// "READ" }
//
// The principal is a user or group name corresponding to an existing
// <Databricks> principal to be granted or revoked access.
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope exists. Throws
// “RESOURCE_ALREADY_EXISTS“ if a permission for the principal already exists.
// Throws “INVALID_PARAMETER_VALUE“ if the permission or principal is invalid.
// Throws “PERMISSION_DENIED“ if the user does not have permission to make
// this API call.
func (c *Client) PutAcl(ctx context.Context, req *PutAcl, opts ...api.Option) (*PutAcl_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/acls/put"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PutAcl_Response{}

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

// Inserts a secret under the provided scope with the given name. If a secret
// already exists with the same name, this command overwrites the existing
// secret's value. The server encrypts the secret using the secret scope's
// encryption settings before storing it. You must have “WRITE“ or “MANAGE“
// permission on the secret scope.
//
// The secret key must consist of alphanumeric characters, dashes, underscores,
// and periods, and cannot exceed 128 characters. The maximum allowed secret
// value size is 128 KB. The maximum number of secrets in a given scope is 1000.
//
// Example request:
//
// .. code::
//
// { "scope": "my-databricks-scope", "key": "my-string-key", "string_value":
// "foobar" }
//
// The input fields "string_value" or "bytes_value" specify the type of the
// secret, which will determine the value returned when the secret value is
// requested. Exactly one must be specified.
//
// Throws “RESOURCE_DOES_NOT_EXIST“ if no such secret scope exists. Throws
// “RESOURCE_LIMIT_EXCEEDED“ if maximum number of secrets in scope is
// exceeded. Throws “INVALID_PARAMETER_VALUE“ if the request parameters are
// invalid. Throws “PERMISSION_DENIED“ if the user does not have permission to
// make this API call. Throws “MALFORMED_REQUEST“ if request is incorrectly
// formatted or conflicting. Throws “BAD_REQUEST“ if request is made against
// Azure KeyVault backed scope.
func (c *Client) PutSecret(ctx context.Context, req *PutSecret, opts ...api.Option) (*PutSecret_Response, error) {
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
	baseURL.Path = "/api/2.0/secrets/put"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &PutSecret_Response{}

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
