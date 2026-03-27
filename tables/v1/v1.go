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

// Creates a new table in the specified catalog and schema.
//
// To create an external delta table, the caller must have the
// **EXTERNAL_USE_SCHEMA** privilege on the parent schema and the
// **EXTERNAL_USE_LOCATION** privilege on the external location. These
// privileges must always be granted explicitly, and cannot be inherited through
// ownership or **ALL_PRIVILEGES**.
//
// Standard UC permissions needed to create tables still apply: **USE_CATALOG**
// on the parent catalog (or ownership of the parent catalog), **CREATE_TABLE**
// and **USE_SCHEMA** on the parent schema (or ownership of the parent schema),
// and **CREATE_EXTERNAL_TABLE** on external location.
//
// The **columns** field needs to be in a Spark compatible format, so we
// recommend you use Spark to create these tables. The API itself does not
// validate the correctness of the column spec. If the spec is not Spark
// compatible, the tables may not be readable by Databricks Runtime.
//
// NOTE: The Create Table API for external clients only supports creating
// **external delta tables**. The values shown in the respective enums are all
// values supported by <Databricks>, however for this specific Create Table API,
// only **table_type** **EXTERNAL** and **data_source_format** **DELTA** are
// supported. Additionally, column masks are not supported when creating tables
// through this API.
func (c *Client) CreateTable(ctx context.Context, req *CreateTable, opts ...api.Option) (*TableInfo, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/tables"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TableInfo{}

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

// Creates a new table constraint.
//
// For the table constraint creation to succeed, the user must satisfy both of
// these conditions: - the user must have the **USE_CATALOG** privilege on the
// table's parent catalog, the **USE_SCHEMA** privilege on the table's parent
// schema, and be the owner of the table. - if the new constraint is a
// __ForeignKeyConstraint__, the user must have the **USE_CATALOG** privilege on
// the referenced parent table's catalog, the **USE_SCHEMA** privilege on the
// referenced parent table's schema, and be the owner of the referenced parent
// table.
func (c *Client) CreateTableConstraint(ctx context.Context, req *CreateTableConstraint, opts ...api.Option) (*TableConstraint, error) {
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
	baseURL.Path = "/api/2.1/unity-catalog/constraints"
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TableConstraint{}

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

// Deletes a table from the specified parent catalog and schema. The caller must
// be the owner of the parent catalog, have the **USE_CATALOG** privilege on the
// parent catalog and be the owner of the parent schema, or be the owner of the
// table and have the **USE_CATALOG** privilege on the parent catalog and the
// **USE_SCHEMA** privilege on the parent schema.
func (c *Client) DeleteTable(ctx context.Context, req *DeleteTable, opts ...api.Option) (*DeleteTable_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteTable_Response{}

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

// Deletes a table constraint.
//
// For the table constraint deletion to succeed, the user must satisfy both of
// these conditions: - the user must have the **USE_CATALOG** privilege on the
// table's parent catalog, the **USE_SCHEMA** privilege on the table's parent
// schema, and be the owner of the table. - if __cascade__ argument is **true**,
// the user must have the following permissions on all of the child tables: the
// **USE_CATALOG** privilege on the table's catalog, the **USE_SCHEMA**
// privilege on the table's schema, and be the owner of the table.
func (c *Client) DeleteTableConstraint(ctx context.Context, req *DeleteTableConstraint, opts ...api.Option) (*DeleteTableConstraint_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/constraints/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.ConstraintName != nil {
		queryParams.Add("constraint_name", fmt.Sprintf("%v", *req.ConstraintName))
	}
	if req.Cascade != nil {
		queryParams.Add("cascade", fmt.Sprintf("%v", *req.Cascade))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &DeleteTableConstraint_Response{}

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

// Gets a table from the metastore for a specific catalog and schema. The caller
// must satisfy one of the following requirements: * Be a metastore admin * Be
// the owner of the parent catalog * Be the owner of the parent schema and have
// the **USE_CATALOG** privilege on the parent catalog * Have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema, and either be the table owner or have the
// **SELECT** privilege on the table.
func (c *Client) GetTable(ctx context.Context, req *GetTable, opts ...api.Option) (*TableInfo, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v", *req.FullNameArg)
	queryParams := url.Values{}
	if req.IncludeDeltaMetadata != nil {
		queryParams.Add("include_delta_metadata", fmt.Sprintf("%v", *req.IncludeDeltaMetadata))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	if req.IncludeManifestCapabilities != nil {
		queryParams.Add("include_manifest_capabilities", fmt.Sprintf("%v", *req.IncludeManifestCapabilities))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TableInfo{}

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

// Gets an array of summaries for tables for a schema and catalog within the
// metastore. The table summaries returned are either:
//
// * summaries for tables (within the current metastore and parent catalog and
// schema), when the user is a metastore admin, or: * summaries for tables and
// schemas (within the current metastore and parent catalog) for which the user
// has ownership or the **SELECT** privilege on the table and ownership or
// **USE_SCHEMA** privilege on the schema, provided that the user also has
// ownership or the **USE_CATALOG** privilege on the parent catalog.
//
// There is no guarantee of a specific ordering of the elements in the array.
//
// PAGINATION BEHAVIOR: The API is by default paginated, a page may contain zero
// results while still providing a next_page_token. Clients must continue
// reading pages until next_page_token is absent, which is the only indication
// that the end of results has been reached.
func (c *Client) ListTableSummaries(ctx context.Context, req *ListTableSummaries, opts ...api.Option) (*ListTableSummaries_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/table-summaries"
	queryParams := url.Values{}
	if req.CatalogName != nil {
		queryParams.Add("catalog_name", fmt.Sprintf("%v", *req.CatalogName))
	}
	if req.SchemaNamePattern != nil {
		queryParams.Add("schema_name_pattern", fmt.Sprintf("%v", *req.SchemaNamePattern))
	}
	if req.TableNamePattern != nil {
		queryParams.Add("table_name_pattern", fmt.Sprintf("%v", *req.TableNamePattern))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.IncludeManifestCapabilities != nil {
		queryParams.Add("include_manifest_capabilities", fmt.Sprintf("%v", *req.IncludeManifestCapabilities))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListTableSummaries_Response{}

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

func (c *Client) ListTableSummariesIter(ctx context.Context, req *ListTableSummaries, opts ...api.Option) iter.Seq2[*TableSummary, error] {
	return func(yield func(*TableSummary, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListTableSummaries{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListTableSummaries(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Tables {
				if !yield(&resp.Tables[i], nil) {
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

// Gets an array of all tables for the current metastore under the parent
// catalog and schema. The caller must be a metastore admin or an owner of (or
// have the **SELECT** privilege on) the table. For the latter case, the caller
// must also be the owner or have the **USE_CATALOG** privilege on the parent
// catalog and the **USE_SCHEMA** privilege on the parent schema. There is no
// guarantee of a specific ordering of the elements in the array.
//
// NOTE: **view_dependencies** and **table_constraints** are not returned by
// ListTables queries.
//
// NOTE: we recommend using max_results=0 to use the paginated version of this
// API. Unpaginated calls will be deprecated soon.
//
// PAGINATION BEHAVIOR: When using pagination (max_results >= 0), a page may
// contain zero results while still providing a next_page_token. Clients must
// continue reading pages until next_page_token is absent, which is the only
// indication that the end of results has been reached.
func (c *Client) ListTables(ctx context.Context, req *ListTables, opts ...api.Option) (*ListTables_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = "/api/2.1/unity-catalog/tables"
	queryParams := url.Values{}
	if req.CatalogName != nil {
		queryParams.Add("catalog_name", fmt.Sprintf("%v", *req.CatalogName))
	}
	if req.SchemaName != nil {
		queryParams.Add("schema_name", fmt.Sprintf("%v", *req.SchemaName))
	}
	if req.MaxResults != nil {
		queryParams.Add("max_results", fmt.Sprintf("%v", *req.MaxResults))
	}
	if req.PageToken != nil {
		queryParams.Add("page_token", fmt.Sprintf("%v", *req.PageToken))
	}
	if req.OmitColumns != nil {
		queryParams.Add("omit_columns", fmt.Sprintf("%v", *req.OmitColumns))
	}
	if req.OmitProperties != nil {
		queryParams.Add("omit_properties", fmt.Sprintf("%v", *req.OmitProperties))
	}
	if req.OmitUsername != nil {
		queryParams.Add("omit_username", fmt.Sprintf("%v", *req.OmitUsername))
	}
	if req.IncludeBrowse != nil {
		queryParams.Add("include_browse", fmt.Sprintf("%v", *req.IncludeBrowse))
	}
	if req.IncludeManifestCapabilities != nil {
		queryParams.Add("include_manifest_capabilities", fmt.Sprintf("%v", *req.IncludeManifestCapabilities))
	}
	baseURL.RawQuery = queryParams.Encode()

	resp := &ListTables_Response{}

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

func (c *Client) ListTablesIter(ctx context.Context, req *ListTables, opts ...api.Option) iter.Seq2[*TableInfo, error] {
	return func(yield func(*TableInfo, error) bool) {
		// Deep copy the request via JSON round-trip to avoid modifying the original.
		reqBody, err := json.Marshal(req)
		if err != nil {
			yield(nil, err)
			return
		}
		pageReq := ListTables{}
		if err := json.Unmarshal(reqBody, &pageReq); err != nil {
			yield(nil, err)
			return
		}
		for {
			resp, err := c.ListTables(ctx, &pageReq, opts...)
			if err != nil {
				yield(nil, err)
				return
			}
			for i := range resp.Tables {
				if !yield(&resp.Tables[i], nil) {
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

// Gets if a table exists in the metastore for a specific catalog and schema.
// The caller must satisfy one of the following requirements: * Be a metastore
// admin * Be the owner of the parent catalog * Be the owner of the parent
// schema and have the **USE_CATALOG** privilege on the parent catalog * Have
// the **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema, and either be the table owner or have the
// **SELECT** privilege on the table. * Have **BROWSE** privilege on the parent
// catalog * Have **BROWSE** privilege on the parent schema
func (c *Client) TableExists(ctx context.Context, req *TableExists, opts ...api.Option) (*TableExists_Response, error) {

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	baseURL, err := url.Parse(c.host)
	if err != nil {
		return nil, err
	}
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v/exists", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &TableExists_Response{}

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

// Change the owner of the table. The caller must be the owner of the parent
// catalog, have the **USE_CATALOG** privilege on the parent catalog and be the
// owner of the parent schema, or be the owner of the table and have the
// **USE_CATALOG** privilege on the parent catalog and the **USE_SCHEMA**
// privilege on the parent schema.
func (c *Client) UpdateTable(ctx context.Context, req *UpdateTable, opts ...api.Option) (*UpdateTable_Response, error) {
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
	baseURL.Path = fmt.Sprintf("/api/2.1/unity-catalog/tables/%v", *req.FullNameArg)
	queryParams := url.Values{}
	baseURL.RawQuery = queryParams.Encode()

	resp := &UpdateTable_Response{}

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
