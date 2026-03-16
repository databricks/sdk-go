package v1

import "encoding/json"

type ColumnTypeName string

const (
	ColumnTypeName_Boolean         ColumnTypeName = "BOOLEAN"
	ColumnTypeName_Byte            ColumnTypeName = "BYTE"
	ColumnTypeName_Short           ColumnTypeName = "SHORT"
	ColumnTypeName_Int             ColumnTypeName = "INT"
	ColumnTypeName_Long            ColumnTypeName = "LONG"
	ColumnTypeName_Float           ColumnTypeName = "FLOAT"
	ColumnTypeName_Double          ColumnTypeName = "DOUBLE"
	ColumnTypeName_Date            ColumnTypeName = "DATE"
	ColumnTypeName_Timestamp       ColumnTypeName = "TIMESTAMP"
	ColumnTypeName_String          ColumnTypeName = "STRING"
	ColumnTypeName_Binary          ColumnTypeName = "BINARY"
	ColumnTypeName_Decimal         ColumnTypeName = "DECIMAL"
	ColumnTypeName_Interval        ColumnTypeName = "INTERVAL"
	ColumnTypeName_Array           ColumnTypeName = "ARRAY"
	ColumnTypeName_Struct          ColumnTypeName = "STRUCT"
	ColumnTypeName_Map             ColumnTypeName = "MAP"
	ColumnTypeName_Char            ColumnTypeName = "CHAR"
	ColumnTypeName_Null            ColumnTypeName = "NULL"
	ColumnTypeName_UserDefinedType ColumnTypeName = "USER_DEFINED_TYPE"
)

// String representation for [fmt.Print].
func (f *ColumnTypeName) String() string {
	return string(*f)
}

type Disposition string

const (
	Disposition_FetchDispositionUnspecified Disposition = "FETCH_DISPOSITION_UNSPECIFIED"
	Disposition_Inline                      Disposition = "INLINE"
	Disposition_ExternalLinks               Disposition = "EXTERNAL_LINKS"
)

// String representation for [fmt.Print].
func (f *Disposition) String() string {
	return string(*f)
}

type Format string

const (
	Format_FormatUnspecified Format = "FORMAT_UNSPECIFIED"
	Format_JsonArray         Format = "JSON_ARRAY"
	Format_ArrowStream       Format = "ARROW_STREAM"
	Format_Csv               Format = "CSV"
)

// String representation for [fmt.Print].
func (f *Format) String() string {
	return string(*f)
}

type ServiceErrorCode string

const (
	ServiceErrorCode_Unknown                         ServiceErrorCode = "UNKNOWN"
	ServiceErrorCode_InternalError                   ServiceErrorCode = "INTERNAL_ERROR"
	ServiceErrorCode_TemporarilyUnavailable          ServiceErrorCode = "TEMPORARILY_UNAVAILABLE"
	ServiceErrorCode_IoError                         ServiceErrorCode = "IO_ERROR"
	ServiceErrorCode_BadRequest                      ServiceErrorCode = "BAD_REQUEST"
	ServiceErrorCode_ServiceUnderMaintenance         ServiceErrorCode = "SERVICE_UNDER_MAINTENANCE"
	ServiceErrorCode_WorkspaceTemporarilyUnavailable ServiceErrorCode = "WORKSPACE_TEMPORARILY_UNAVAILABLE"
	ServiceErrorCode_DeadlineExceeded                ServiceErrorCode = "DEADLINE_EXCEEDED"
	ServiceErrorCode_Cancelled                       ServiceErrorCode = "CANCELLED"
	ServiceErrorCode_ResourceExhausted               ServiceErrorCode = "RESOURCE_EXHAUSTED"
	ServiceErrorCode_Aborted                         ServiceErrorCode = "ABORTED"
	ServiceErrorCode_NotFound                        ServiceErrorCode = "NOT_FOUND"
	ServiceErrorCode_AlreadyExists                   ServiceErrorCode = "ALREADY_EXISTS"
	ServiceErrorCode_Unauthenticated                 ServiceErrorCode = "UNAUTHENTICATED"
)

// String representation for [fmt.Print].
func (f *ServiceErrorCode) String() string {
	return string(*f)
}

type TimeoutAction string

const (
	TimeoutAction_TimeoutActionUnspecified TimeoutAction = "TIMEOUT_ACTION_UNSPECIFIED"
	TimeoutAction_Continue                 TimeoutAction = "CONTINUE"
	TimeoutAction_Cancel                   TimeoutAction = "CANCEL"
)

// String representation for [fmt.Print].
func (f *TimeoutAction) String() string {
	return string(*f)
}

type StatementStatus_State string

const (
	StatementStatus_State_StateUnspecified StatementStatus_State = "STATE_UNSPECIFIED"
	StatementStatus_State_Pending          StatementStatus_State = "PENDING"
	StatementStatus_State_Running          StatementStatus_State = "RUNNING"
	StatementStatus_State_Succeeded        StatementStatus_State = "SUCCEEDED"
	StatementStatus_State_Failed           StatementStatus_State = "FAILED"
	StatementStatus_State_Canceled         StatementStatus_State = "CANCELED"
	StatementStatus_State_Closed           StatementStatus_State = "CLOSED"
)

// String representation for [fmt.Print].
func (f *StatementStatus_State) String() string {
	return string(*f)
}

type CancelStatementRequest struct {
	StatementId *string `json:"statement_id"`
}

type CancelStatementResponse struct {
}

type ChunkInfo struct {
	ChunkIndex            *int    `json:"chunk_index"`
	RowOffset             *int64  `json:"row_offset"`
	RowCount              *int64  `json:"row_count"`
	ByteCount             *int64  `json:"byte_count"`
	NextChunkIndex        *int    `json:"next_chunk_index"`
	NextChunkInternalLink *string `json:"next_chunk_internal_link"`
}

type ColumnInfo struct {
	Name             *string         `json:"name"`
	TypeText         *string         `json:"type_text"`
	TypeName         *ColumnTypeName `json:"type_name"`
	Position         *int            `json:"position"`
	TypePrecision    *int            `json:"type_precision"`
	TypeScale        *int            `json:"type_scale"`
	TypeIntervalType *string         `json:"type_interval_type"`
}

type ExecuteStatementRequest struct {
	Statement     *string              `json:"statement"`
	WarehouseId   *string              `json:"warehouse_id"`
	Catalog       *string              `json:"catalog"`
	Schema        *string              `json:"schema"`
	RowLimit      *int64               `json:"row_limit"`
	ByteLimit     *int64               `json:"byte_limit"`
	Format        *Format              `json:"format"`
	Disposition   *Disposition         `json:"disposition"`
	WaitTimeout   *string              `json:"wait_timeout"`
	OnWaitTimeout *TimeoutAction       `json:"on_wait_timeout"`
	Parameters    []StatementParameter `json:"parameters"`
	QueryTags     []QueryTag           `json:"query_tags"`
}

type ExternalLink struct {
	ExternalLink          *string           `json:"external_link"`
	Expiration            *string           `json:"expiration"`
	HttpHeaders           map[string]string `json:"http_headers"`
	ChunkIndex            *int              `json:"chunk_index"`
	RowOffset             *int64            `json:"row_offset"`
	RowCount              *int64            `json:"row_count"`
	ByteCount             *int64            `json:"byte_count"`
	NextChunkIndex        *int              `json:"next_chunk_index"`
	NextChunkInternalLink *string           `json:"next_chunk_internal_link"`
}

type ExternalLink_HttpHeadersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type GetResultDataRequest struct {
	StatementId *string `json:"statement_id"`
	ChunkIndex  *int    `json:"chunk_index"`
}

type GetStatementResultRequest struct {
	StatementId *string `json:"statement_id"`
}

// * A query execution can be annotated with an optional key-value pair to allow
// users to attribute the executions by key and optional value to filter by.
// QueryTag is the user-facing representation..
type QueryTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Contains the result data of a single chunk when using `INLINE` disposition.
// When using `EXTERNAL_LINKS` disposition, the array `external_links` is used
// instead to provide URLs to the result data in cloud storage. Exactly one of
// these alternatives is used. (While the `external_links` array prepares the
// API to return multiple links in a single response. Currently only a single
// link is returned.).
type ResultData struct {
	ExternalLinks         []ExternalLink    `json:"external_links"`
	DataArray             []json.RawMessage `json:"data_array"`
	ChunkIndex            *int              `json:"chunk_index"`
	RowOffset             *int64            `json:"row_offset"`
	RowCount              *int64            `json:"row_count"`
	ByteCount             *int64            `json:"byte_count"`
	NextChunkIndex        *int              `json:"next_chunk_index"`
	NextChunkInternalLink *string           `json:"next_chunk_internal_link"`
}

// The result manifest provides schema and metadata for the result set..
type ResultManifest struct {
	Format          *Format     `json:"format"`
	Schema          *Schema     `json:"schema"`
	TotalChunkCount *int        `json:"total_chunk_count"`
	Chunks          []ChunkInfo `json:"chunks"`
	TotalRowCount   *int64      `json:"total_row_count"`
	TotalByteCount  *int64      `json:"total_byte_count"`
	Truncated       *bool       `json:"truncated"`
}

// The schema is an ordered list of column descriptions..
type Schema struct {
	ColumnCount *int         `json:"column_count"`
	Columns     []ColumnInfo `json:"columns"`
}

type ServiceError struct {
	ErrorCode *ServiceErrorCode `json:"error_code"`
	Message   *string           `json:"message"`
}

type StatementParameter struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
	Type  *string `json:"type"`
}

type StatementResponse struct {
	StatementId *string          `json:"statement_id"`
	Status      *StatementStatus `json:"status"`
	Manifest    *ResultManifest  `json:"manifest"`
	Result      *ResultData      `json:"result"`
}

// The status response includes execution state and if relevant, error
// information..
type StatementStatus struct {
	State *StatementStatus_State `json:"state"`
	Error *ServiceError          `json:"error"`
}
