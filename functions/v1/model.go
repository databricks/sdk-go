package v1

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
	ColumnTypeName_TimestampNtz    ColumnTypeName = "TIMESTAMP_NTZ"
	ColumnTypeName_Variant         ColumnTypeName = "VARIANT"
	ColumnTypeName_Geometry        ColumnTypeName = "GEOMETRY"
	ColumnTypeName_Geography       ColumnTypeName = "GEOGRAPHY"
	ColumnTypeName_Time            ColumnTypeName = "TIME"
	ColumnTypeName_TableType       ColumnTypeName = "TABLE_TYPE"
)

// String representation for [fmt.Print].
func (f *ColumnTypeName) String() string {
	return string(*f)
}

type FunctionParameterMode string

const (
	FunctionParameterMode_In FunctionParameterMode = "IN"
)

// String representation for [fmt.Print].
func (f *FunctionParameterMode) String() string {
	return string(*f)
}

type FunctionParameterType string

const (
	FunctionParameterType_Param  FunctionParameterType = "PARAM"
	FunctionParameterType_Column FunctionParameterType = "COLUMN"
)

// String representation for [fmt.Print].
func (f *FunctionParameterType) String() string {
	return string(*f)
}

type FunctionInfo_ParameterStyle string

const (
	FunctionInfo_ParameterStyle_S FunctionInfo_ParameterStyle = "S"
)

// String representation for [fmt.Print].
func (f *FunctionInfo_ParameterStyle) String() string {
	return string(*f)
}

type FunctionInfo_RoutineBody string

const (
	FunctionInfo_RoutineBody_Sql      FunctionInfo_RoutineBody = "SQL"
	FunctionInfo_RoutineBody_External FunctionInfo_RoutineBody = "EXTERNAL"
)

// String representation for [fmt.Print].
func (f *FunctionInfo_RoutineBody) String() string {
	return string(*f)
}

type FunctionInfo_SecurityType string

const (
	FunctionInfo_SecurityType_Definer FunctionInfo_SecurityType = "DEFINER"
)

// String representation for [fmt.Print].
func (f *FunctionInfo_SecurityType) String() string {
	return string(*f)
}

type FunctionInfo_SqlDataAccess string

const (
	FunctionInfo_SqlDataAccess_ContainsSql  FunctionInfo_SqlDataAccess = "CONTAINS_SQL"
	FunctionInfo_SqlDataAccess_ReadsSqlData FunctionInfo_SqlDataAccess = "READS_SQL_DATA"
	FunctionInfo_SqlDataAccess_NoSql        FunctionInfo_SqlDataAccess = "NO_SQL"
)

// String representation for [fmt.Print].
func (f *FunctionInfo_SqlDataAccess) String() string {
	return string(*f)
}

// A connection that is dependent on a SQL object..
type ConnectionDependency struct {
	ConnectionName *string `json:"connection_name"`
}

type CreateFunction struct {
	Name                *string                      `json:"name"`
	CatalogName         *string                      `json:"catalog_name"`
	SchemaName          *string                      `json:"schema_name"`
	InputParams         *FunctionParameterInfos      `json:"input_params"`
	DataType            *ColumnTypeName              `json:"data_type"`
	FullDataType        *string                      `json:"full_data_type"`
	RoutineBody         *FunctionInfo_RoutineBody    `json:"routine_body"`
	RoutineDefinition   *string                      `json:"routine_definition"`
	ParameterStyle      *FunctionInfo_ParameterStyle `json:"parameter_style"`
	IsDeterministic     *bool                        `json:"is_deterministic"`
	SqlDataAccess       *FunctionInfo_SqlDataAccess  `json:"sql_data_access"`
	IsNullCall          *bool                        `json:"is_null_call"`
	SecurityType        *FunctionInfo_SecurityType   `json:"security_type"`
	SpecificName        *string                      `json:"specific_name"`
	ReturnParams        *FunctionParameterInfos      `json:"return_params"`
	ExternalName        *string                      `json:"external_name"`
	ExternalLanguage    *string                      `json:"external_language"`
	SqlPath             *string                      `json:"sql_path"`
	Owner               *string                      `json:"owner"`
	Comment             *string                      `json:"comment"`
	Properties          *string                      `json:"properties"`
	RoutineDependencies *DependencyList              `json:"routine_dependencies"`
	MetastoreId         *string                      `json:"metastore_id"`
	FullName            *string                      `json:"full_name"`
	CreatedAt           *int64                       `json:"created_at"`
	CreatedBy           *string                      `json:"created_by"`
	UpdatedAt           *int64                       `json:"updated_at"`
	UpdatedBy           *string                      `json:"updated_by"`
	FunctionId          *string                      `json:"function_id"`
	BrowseOnly          *bool                        `json:"browse_only"`
}

type CreateFunctionRequest struct {
	FunctionInfo *CreateFunction `json:"function_info"`
}

// A credential that is dependent on a SQL object..
type CredentialDependency struct {
	CredentialName *string `json:"credential_name"`
}

type DeleteFunction struct {
	FullNameArg *string `json:"full_name_arg"`
	Force       *bool   `json:"force"`
}

type DeleteFunction_Response struct {
}

// A dependency of a SQL object. One of the following fields must be defined:
// __table__, __function__, __connection__, or __credential__..
type Dependency struct {
	Table      *TableDependency      `json:"table"`
	Function   *FunctionDependency   `json:"function"`
	Connection *ConnectionDependency `json:"connection"`
	Credential *CredentialDependency `json:"credential"`
}

// A list of dependencies..
type DependencyList struct {
	Dependencies []Dependency `json:"dependencies"`
}

// A function that is dependent on a SQL object..
type FunctionDependency struct {
	FunctionFullName *string `json:"function_full_name"`
}

type FunctionInfo struct {
	Name                *string                      `json:"name"`
	CatalogName         *string                      `json:"catalog_name"`
	SchemaName          *string                      `json:"schema_name"`
	InputParams         *FunctionParameterInfos      `json:"input_params"`
	DataType            *ColumnTypeName              `json:"data_type"`
	FullDataType        *string                      `json:"full_data_type"`
	RoutineBody         *FunctionInfo_RoutineBody    `json:"routine_body"`
	RoutineDefinition   *string                      `json:"routine_definition"`
	ParameterStyle      *FunctionInfo_ParameterStyle `json:"parameter_style"`
	IsDeterministic     *bool                        `json:"is_deterministic"`
	SqlDataAccess       *FunctionInfo_SqlDataAccess  `json:"sql_data_access"`
	IsNullCall          *bool                        `json:"is_null_call"`
	SecurityType        *FunctionInfo_SecurityType   `json:"security_type"`
	SpecificName        *string                      `json:"specific_name"`
	ReturnParams        *FunctionParameterInfos      `json:"return_params"`
	ExternalName        *string                      `json:"external_name"`
	ExternalLanguage    *string                      `json:"external_language"`
	SqlPath             *string                      `json:"sql_path"`
	Owner               *string                      `json:"owner"`
	Comment             *string                      `json:"comment"`
	Properties          *string                      `json:"properties"`
	RoutineDependencies *DependencyList              `json:"routine_dependencies"`
	MetastoreId         *string                      `json:"metastore_id"`
	FullName            *string                      `json:"full_name"`
	CreatedAt           *int64                       `json:"created_at"`
	CreatedBy           *string                      `json:"created_by"`
	UpdatedAt           *int64                       `json:"updated_at"`
	UpdatedBy           *string                      `json:"updated_by"`
	FunctionId          *string                      `json:"function_id"`
	BrowseOnly          *bool                        `json:"browse_only"`
}

type FunctionParameterInfo struct {
	Name             *string                `json:"name"`
	TypeText         *string                `json:"type_text"`
	TypeJson         *string                `json:"type_json"`
	TypeName         *ColumnTypeName        `json:"type_name"`
	TypePrecision    *int                   `json:"type_precision"`
	TypeScale        *int                   `json:"type_scale"`
	TypeIntervalType *string                `json:"type_interval_type"`
	Position         *int                   `json:"position"`
	ParameterMode    *FunctionParameterMode `json:"parameter_mode"`
	ParameterType    *FunctionParameterType `json:"parameter_type"`
	ParameterDefault *string                `json:"parameter_default"`
	Comment          *string                `json:"comment"`
}

type FunctionParameterInfos struct {
	Parameters []FunctionParameterInfo `json:"parameters"`
}

type GetFunction struct {
	FullNameArg   *string `json:"full_name_arg"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListFunctions struct {
	CatalogName   *string `json:"catalog_name"`
	SchemaName    *string `json:"schema_name"`
	IncludeBrowse *bool   `json:"include_browse"`
	MaxResults    *int    `json:"max_results"`
	PageToken     *string `json:"page_token"`
}

type ListFunctions_Response struct {
	Functions     []FunctionInfo `json:"functions"`
	NextPageToken *string        `json:"next_page_token"`
}

// A table that is dependent on a SQL object..
type TableDependency struct {
	TableFullName *string `json:"table_full_name"`
}

type UpdateFunction struct {
	FullNameArg         *string                      `json:"full_name_arg"`
	Name                *string                      `json:"name"`
	CatalogName         *string                      `json:"catalog_name"`
	SchemaName          *string                      `json:"schema_name"`
	InputParams         *FunctionParameterInfos      `json:"input_params"`
	DataType            *ColumnTypeName              `json:"data_type"`
	FullDataType        *string                      `json:"full_data_type"`
	RoutineBody         *FunctionInfo_RoutineBody    `json:"routine_body"`
	RoutineDefinition   *string                      `json:"routine_definition"`
	ParameterStyle      *FunctionInfo_ParameterStyle `json:"parameter_style"`
	IsDeterministic     *bool                        `json:"is_deterministic"`
	SqlDataAccess       *FunctionInfo_SqlDataAccess  `json:"sql_data_access"`
	IsNullCall          *bool                        `json:"is_null_call"`
	SecurityType        *FunctionInfo_SecurityType   `json:"security_type"`
	SpecificName        *string                      `json:"specific_name"`
	ReturnParams        *FunctionParameterInfos      `json:"return_params"`
	ExternalName        *string                      `json:"external_name"`
	ExternalLanguage    *string                      `json:"external_language"`
	SqlPath             *string                      `json:"sql_path"`
	Owner               *string                      `json:"owner"`
	Comment             *string                      `json:"comment"`
	Properties          *string                      `json:"properties"`
	RoutineDependencies *DependencyList              `json:"routine_dependencies"`
	MetastoreId         *string                      `json:"metastore_id"`
	FullName            *string                      `json:"full_name"`
	CreatedAt           *int64                       `json:"created_at"`
	CreatedBy           *string                      `json:"created_by"`
	UpdatedAt           *int64                       `json:"updated_at"`
	UpdatedBy           *string                      `json:"updated_by"`
	FunctionId          *string                      `json:"function_id"`
	BrowseOnly          *bool                        `json:"browse_only"`
}
