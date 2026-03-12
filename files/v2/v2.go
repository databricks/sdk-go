package v2

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

// Appends a block of data to the stream specified by the input handle. If the
// handle does not exist, this call will throw an exception with
// “RESOURCE_DOES_NOT_EXIST“.
//
// If the block of data exceeds 1 MB, this call will throw an exception with
// “MAX_BLOCK_SIZE_EXCEEDED“.
func (c *Client) AddBlock(ctx context.Context, req *AddBlock, opts ...api.Option) (*AddBlock_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/add-block"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &AddBlock_Response{}

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

// Closes the stream specified by the input handle. If the handle does not
// exist, this call throws an exception with “RESOURCE_DOES_NOT_EXIST“.
func (c *Client) Close(ctx context.Context, req *Close, opts ...api.Option) (*Close_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/close"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Close_Response{}

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

// Opens a stream to write to a file and returns a handle to this stream. There
// is a 10 minute idle timeout on this handle. If a file or directory already
// exists on the given path and __overwrite__ is set to false, this call will
// throw an exception with “RESOURCE_ALREADY_EXISTS“.
//
// A typical workflow for file upload would be:
//
// 1. Issue a “create“ call and get a handle. 2. Issue one or more
// “add-block“ calls with the handle you have. 3. Issue a “close“ call with
// the handle you have.
func (c *Client) Create(ctx context.Context, req *Create, opts ...api.Option) (*Create_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/create"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Create_Response{}

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

// Delete the file or directory (optionally recursively delete all files in the
// directory). This call throws an exception with `IO_ERROR` if the path is a
// non-empty directory and `recursive` is set to `false` or on other similar
// errors.
//
// When you delete a large number of files, the delete operation is done in
// increments. The call returns a response after approximately 45 seconds with
// an error message (503 Service Unavailable) asking you to re-invoke the delete
// operation until the directory structure is fully deleted.
//
// For operations that delete more than 10K files, we discourage using the DBFS
// REST API, but advise you to perform such operations in the context of a
// cluster, using the [File system utility
// (dbutils.fs)](/dev-tools/databricks-utils.html#dbutils-fs). `dbutils.fs`
// covers the functional scope of the DBFS REST API, but from notebooks. Running
// such operations using notebooks provides better control and manageability,
// such as selective deletes, and the possibility to automate periodic delete
// jobs.
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
	baseURL.Path = "/api/2.0/dbfs/delete"
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

// Gets the file information for a file or directory. If the file or directory
// does not exist, this call throws an exception with `RESOURCE_DOES_NOT_EXIST`.
func (c *Client) GetStatus(ctx context.Context, req *GetStatus, opts ...api.Option) (*GetStatus_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/dbfs/get-status"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &GetStatus_Response{}

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

// List the contents of a directory, or details of the file. If the file or
// directory does not exist, this call throws an exception with
// `RESOURCE_DOES_NOT_EXIST`.
//
// When calling list on a large directory, the list operation will time out
// after approximately 60 seconds. We strongly recommend using list only on
// directories containing less than 10K files and discourage using the DBFS REST
// API for operations that list more than 10K files. Instead, we recommend that
// you perform such operations in the context of a cluster, using the [File
// system utility (dbutils.fs)](/dev-tools/databricks-utils.html#dbutils-fs),
// which provides the same functionality without timing out.
func (c *Client) List(ctx context.Context, req *ListStatus, opts ...api.Option) (*ListStatus_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/dbfs/list"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListStatus_Response{}

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

// Creates the given directory and necessary parent directories if they do not
// exist. If a file (not a directory) exists at any prefix of the input path,
// this call throws an exception with `RESOURCE_ALREADY_EXISTS`. **Note**: If
// this operation fails, it might have succeeded in creating some of the
// necessary parent directories.
func (c *Client) Mkdirs(ctx context.Context, req *MkDirs, opts ...api.Option) (*MkDirs_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/mkdirs"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &MkDirs_Response{}

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

// Moves a file from one location to another location within DBFS. If the source
// file does not exist, this call throws an exception with
// `RESOURCE_DOES_NOT_EXIST`. If a file already exists in the destination path,
// this call throws an exception with `RESOURCE_ALREADY_EXISTS`. If the given
// source path is a directory, this call always recursively moves all files.
func (c *Client) Move(ctx context.Context, req *Move, opts ...api.Option) (*Move_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/move"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Move_Response{}

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

// Uploads a file through the use of multipart form post. It is mainly used for
// streaming uploads, but can also be used as a convenient single call for data
// upload.
//
// Alternatively you can pass contents as base64 string.
//
// The amount of data that can be passed (when not streaming) using the
// __contents__ parameter is limited to 1 MB. `MAX_BLOCK_SIZE_EXCEEDED` will be
// thrown if this limit is exceeded.
//
// If you want to upload large files, use the streaming upload. For details, see
// :method:dbfs/create, :method:dbfs/addBlock, :method:dbfs/close.
func (c *Client) Put(ctx context.Context, req *Put, opts ...api.Option) (*Put_Response, error) {
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
	baseURL.Path = "/api/2.0/dbfs/put"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Put_Response{}

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

// Returns the contents of a file. If the file does not exist, this call throws
// an exception with `RESOURCE_DOES_NOT_EXIST`. If the path is a directory, the
// read length is negative, or if the offset is negative, this call throws an
// exception with `INVALID_PARAMETER_VALUE`. If the read length exceeds 1 MB,
// this call throws an exception with `MAX_READ_SIZE_EXCEEDED`.
//
// If `offset + length` exceeds the number of bytes in a file, it reads the
// contents until the end of file.
func (c *Client) Read(ctx context.Context, req *Read, opts ...api.Option) (*Read_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.0/dbfs/read"
	queryParams := url.Values{}
	if req.Path != nil {
		queryParams.Add("path", fmt.Sprintf("%v", *req.Path))
	}
	if req.Offset != nil {
		queryParams.Add("offset", fmt.Sprintf("%v", *req.Offset))
	}
	if req.Length != nil {
		queryParams.Add("length", fmt.Sprintf("%v", *req.Length))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &Read_Response{}

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
