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

// Deletes an object or a directory (and optionally recursively deletes all
// objects in the directory). * If `path` does not exist, this call returns an
// error `RESOURCE_DOES_NOT_EXIST`. * If `path` is a non-empty directory and
// `recursive` is set to `false`, this call returns an error
// `DIRECTORY_NOT_EMPTY`. Object deletion cannot be undone and deleting a
// directory recursively is not atomic.
func (c *Client) Delete(ctx context.Context, req *Delete, opts ...api.Option) (*Delete_Response, error) {
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
	baseURL.Path = "/api/2.0/workspace/delete"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Delete_Response{}

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

// Exports an object or the contents of an entire directory. If `path` does not
// exist, this call returns an error `RESOURCE_DOES_NOT_EXIST`. If the exported
// data would exceed size limit, this call returns `MAX_NOTEBOOK_SIZE_EXCEEDED`.
// Currently, this API does not support exporting a library.
func (c *Client) Export(ctx context.Context, req *Export, opts ...api.Option) (*Export_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/workspace/export"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	if req.Format != nil {
		queryParams.Add("format", fmt.Sprintf("%v", *req.Format))
	}
	if req.DirectDownload != nil {
		queryParams.Add("direct_download", fmt.Sprintf("%v", *req.DirectDownload))
	}
	if req.Outputs != nil {
		queryParams.Add("outputs", fmt.Sprintf("%v", *req.Outputs))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Export_Response{}

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

// Gets the status of an object or a directory. If `path` does not exist, this
// call returns an error `RESOURCE_DOES_NOT_EXIST`.
func (c *Client) GetStatus(ctx context.Context, req *GetStatus, opts ...api.Option) (*ObjectInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/workspace/get-status"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ObjectInfo{}

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

// Imports a workspace object (for example, a notebook or file) or the contents
// of an entire directory. If `path` already exists and `overwrite` is set to
// `false`, this call returns an error `RESOURCE_ALREADY_EXISTS`. To import a
// directory, you can use either the `DBC` format or the `SOURCE` format with
// the `language` field unset. To import a single file as `SOURCE`, you must set
// the `language` field. Zip files within directories are not supported.
func (c *Client) Import(ctx context.Context, req *Import, opts ...api.Option) (*Import_Response, error) {
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
	baseURL.Path = "/api/2.0/workspace/import"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Import_Response{}

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

// Lists the contents of a directory, or the object if it is not a directory. If
// the input path does not exist, this call returns an error
// `RESOURCE_DOES_NOT_EXIST`.
func (c *Client) List(ctx context.Context, req *List, opts ...api.Option) (*List_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/workspace/list"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	if req.NotebooksModifiedAfter != nil {
		queryParams.Add("notebooks_modified_after", fmt.Sprintf("%v", *req.NotebooksModifiedAfter))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &List_Response{}

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

// Creates the specified directory (and necessary parent directories if they do
// not exist). If there is an object (not a directory) at any prefix of the
// input path, this call returns an error `RESOURCE_ALREADY_EXISTS`. Note that
// if this operation fails it may have succeeded in creating some of the
// necessary parent directories.
func (c *Client) Mkdirs(ctx context.Context, req *Mkdirs, opts ...api.Option) (*Mkdirs_Response, error) {
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
	baseURL.Path = "/api/2.0/workspace/mkdirs"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Mkdirs_Response{}

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
