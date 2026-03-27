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

type DataSourceFormat string

const (
	DataSourceFormat_Delta                     DataSourceFormat = "DELTA"
	DataSourceFormat_Csv                       DataSourceFormat = "CSV"
	DataSourceFormat_Json                      DataSourceFormat = "JSON"
	DataSourceFormat_Avro                      DataSourceFormat = "AVRO"
	DataSourceFormat_Parquet                   DataSourceFormat = "PARQUET"
	DataSourceFormat_Orc                       DataSourceFormat = "ORC"
	DataSourceFormat_Text                      DataSourceFormat = "TEXT"
	DataSourceFormat_UnityCatalog              DataSourceFormat = "UNITY_CATALOG"
	DataSourceFormat_Deltasharing              DataSourceFormat = "DELTASHARING"
	DataSourceFormat_DatabricksFormat          DataSourceFormat = "DATABRICKS_FORMAT"
	DataSourceFormat_MysqlFormat               DataSourceFormat = "MYSQL_FORMAT"
	DataSourceFormat_OracleFormat              DataSourceFormat = "ORACLE_FORMAT"
	DataSourceFormat_PostgresqlFormat          DataSourceFormat = "POSTGRESQL_FORMAT"
	DataSourceFormat_RedshiftFormat            DataSourceFormat = "REDSHIFT_FORMAT"
	DataSourceFormat_SnowflakeFormat           DataSourceFormat = "SNOWFLAKE_FORMAT"
	DataSourceFormat_SqldwFormat               DataSourceFormat = "SQLDW_FORMAT"
	DataSourceFormat_SqlserverFormat           DataSourceFormat = "SQLSERVER_FORMAT"
	DataSourceFormat_SalesforceFormat          DataSourceFormat = "SALESFORCE_FORMAT"
	DataSourceFormat_SalesforceDataCloudFormat DataSourceFormat = "SALESFORCE_DATA_CLOUD_FORMAT"
	DataSourceFormat_TeradataFormat            DataSourceFormat = "TERADATA_FORMAT"
	DataSourceFormat_BigqueryFormat            DataSourceFormat = "BIGQUERY_FORMAT"
	DataSourceFormat_NetsuiteFormat            DataSourceFormat = "NETSUITE_FORMAT"
	DataSourceFormat_WorkdayRaasFormat         DataSourceFormat = "WORKDAY_RAAS_FORMAT"
	DataSourceFormat_MongodbFormat             DataSourceFormat = "MONGODB_FORMAT"
	DataSourceFormat_Hive                      DataSourceFormat = "HIVE"
	DataSourceFormat_VectorIndexFormat         DataSourceFormat = "VECTOR_INDEX_FORMAT"
	DataSourceFormat_DatabricksRowStoreFormat  DataSourceFormat = "DATABRICKS_ROW_STORE_FORMAT"
	DataSourceFormat_DeltaUniformHudi          DataSourceFormat = "DELTA_UNIFORM_HUDI"
	DataSourceFormat_DeltaUniformIceberg       DataSourceFormat = "DELTA_UNIFORM_ICEBERG"
	DataSourceFormat_Iceberg                   DataSourceFormat = "ICEBERG"
)

// String representation for [fmt.Print].
func (f *DataSourceFormat) String() string {
	return string(*f)
}

type SecurableKind string

const (
	SecurableKind_TableStandard                                        SecurableKind = "TABLE_STANDARD"
	SecurableKind_TableExternal                                        SecurableKind = "TABLE_EXTERNAL"
	SecurableKind_TableDelta                                           SecurableKind = "TABLE_DELTA"
	SecurableKind_TableDeltaExternal                                   SecurableKind = "TABLE_DELTA_EXTERNAL"
	SecurableKind_TableView                                            SecurableKind = "TABLE_VIEW"
	SecurableKind_TableMetricView                                      SecurableKind = "TABLE_METRIC_VIEW"
	SecurableKind_TableDeltasharing                                    SecurableKind = "TABLE_DELTASHARING"
	SecurableKind_TableDeltasharingMutable                             SecurableKind = "TABLE_DELTASHARING_MUTABLE"
	SecurableKind_TableViewDeltasharing                                SecurableKind = "TABLE_VIEW_DELTASHARING"
	SecurableKind_TableMetricViewDeltasharing                          SecurableKind = "TABLE_METRIC_VIEW_DELTASHARING"
	SecurableKind_TableMaterializedViewDeltasharing                    SecurableKind = "TABLE_MATERIALIZED_VIEW_DELTASHARING"
	SecurableKind_TableStreamingLiveTableDeltasharing                  SecurableKind = "TABLE_STREAMING_LIVE_TABLE_DELTASHARING"
	SecurableKind_TableForeignDeltasharing                             SecurableKind = "TABLE_FOREIGN_DELTASHARING"
	SecurableKind_TableDeltaIcebergDeltasharing                        SecurableKind = "TABLE_DELTA_ICEBERG_DELTASHARING"
	SecurableKind_TableDeltasharingOpenDirBased                        SecurableKind = "TABLE_DELTASHARING_OPEN_DIR_BASED"
	SecurableKind_TableFeatureStore                                    SecurableKind = "TABLE_FEATURE_STORE"
	SecurableKind_TableFeatureStoreExternal                            SecurableKind = "TABLE_FEATURE_STORE_EXTERNAL"
	SecurableKind_TableStreamingLiveTable                              SecurableKind = "TABLE_STREAMING_LIVE_TABLE"
	SecurableKind_TableSystem                                          SecurableKind = "TABLE_SYSTEM"
	SecurableKind_TableSystemDeltasharing                              SecurableKind = "TABLE_SYSTEM_DELTASHARING"
	SecurableKind_TableMaterializedView                                SecurableKind = "TABLE_MATERIALIZED_VIEW"
	SecurableKind_TableInternal                                        SecurableKind = "TABLE_INTERNAL"
	SecurableKind_TableForeignBigquery                                 SecurableKind = "TABLE_FOREIGN_BIGQUERY"
	SecurableKind_TableForeignMysql                                    SecurableKind = "TABLE_FOREIGN_MYSQL"
	SecurableKind_TableForeignOracle                                   SecurableKind = "TABLE_FOREIGN_ORACLE"
	SecurableKind_TableForeignPalantir                                 SecurableKind = "TABLE_FOREIGN_PALANTIR"
	SecurableKind_TableForeignPostgresql                               SecurableKind = "TABLE_FOREIGN_POSTGRESQL"
	SecurableKind_TableForeignSqldw                                    SecurableKind = "TABLE_FOREIGN_SQLDW"
	SecurableKind_TableForeignRedshift                                 SecurableKind = "TABLE_FOREIGN_REDSHIFT"
	SecurableKind_TableForeignSnowflake                                SecurableKind = "TABLE_FOREIGN_SNOWFLAKE"
	SecurableKind_TableForeignSqlserver                                SecurableKind = "TABLE_FOREIGN_SQLSERVER"
	SecurableKind_TableForeignSalesforce                               SecurableKind = "TABLE_FOREIGN_SALESFORCE"
	SecurableKind_TableForeignSalesforceDataCloud                      SecurableKind = "TABLE_FOREIGN_SALESFORCE_DATA_CLOUD"
	SecurableKind_TableForeignSalesforceDataCloudFileSharing           SecurableKind = "TABLE_FOREIGN_SALESFORCE_DATA_CLOUD_FILE_SHARING"
	SecurableKind_TableForeignSalesforceDataCloudFileSharingView       SecurableKind = "TABLE_FOREIGN_SALESFORCE_DATA_CLOUD_FILE_SHARING_VIEW"
	SecurableKind_TableForeignTeradata                                 SecurableKind = "TABLE_FOREIGN_TERADATA"
	SecurableKind_TableForeignNetsuite                                 SecurableKind = "TABLE_FOREIGN_NETSUITE"
	SecurableKind_TableForeignDatabricks                               SecurableKind = "TABLE_FOREIGN_DATABRICKS"
	SecurableKind_TableForeignWorkdayRaas                              SecurableKind = "TABLE_FOREIGN_WORKDAY_RAAS"
	SecurableKind_TableForeignHiveMetastore                            SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE"
	SecurableKind_TableForeignHiveMetastoreManaged                     SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_MANAGED"
	SecurableKind_TableForeignHiveMetastoreDbfsManaged                 SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_DBFS_MANAGED"
	SecurableKind_TableForeignHiveMetastoreExternal                    SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_EXTERNAL"
	SecurableKind_TableForeignHiveMetastoreDbfsExternal                SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_DBFS_EXTERNAL"
	SecurableKind_TableForeignHiveMetastoreView                        SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_VIEW"
	SecurableKind_TableForeignHiveMetastoreDbfsView                    SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_DBFS_VIEW"
	SecurableKind_TableForeignHiveMetastoreShallowCloneManaged         SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_SHALLOW_CLONE_MANAGED"
	SecurableKind_TableForeignHiveMetastoreDbfsShallowCloneManaged     SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_DBFS_SHALLOW_CLONE_MANAGED"
	SecurableKind_TableForeignHiveMetastoreShallowCloneExternal        SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_SHALLOW_CLONE_EXTERNAL"
	SecurableKind_TableForeignHiveMetastoreDbfsShallowCloneExternal    SecurableKind = "TABLE_FOREIGN_HIVE_METASTORE_DBFS_SHALLOW_CLONE_EXTERNAL"
	SecurableKind_TableForeignMongodb                                  SecurableKind = "TABLE_FOREIGN_MONGODB"
	SecurableKind_TableDeltaUniformHudiExternal                        SecurableKind = "TABLE_DELTA_UNIFORM_HUDI_EXTERNAL"
	SecurableKind_TableDeltaUniformIcebergExternal                     SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_EXTERNAL"
	SecurableKind_TableDeltaUniformIcebergForeignHiveMetastoreExternal SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_FOREIGN_HIVE_METASTORE_EXTERNAL"
	SecurableKind_TableDeltaUniformIcebergForeignHiveMetastoreManaged  SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_FOREIGN_HIVE_METASTORE_MANAGED"
	SecurableKind_TableDeltaUniformIcebergForeignSnowflake             SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_FOREIGN_SNOWFLAKE"
	SecurableKind_TableDeltaUniformIcebergForeignDeltasharing          SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_FOREIGN_DELTASHARING"
	SecurableKind_TableDeltaUniformIcebergExternalDeltasharing         SecurableKind = "TABLE_DELTA_UNIFORM_ICEBERG_EXTERNAL_DELTASHARING"
	SecurableKind_TableIcebergUniformManaged                           SecurableKind = "TABLE_ICEBERG_UNIFORM_MANAGED"
	SecurableKind_TableDeltaIcebergManaged                             SecurableKind = "TABLE_DELTA_ICEBERG_MANAGED"
	SecurableKind_TableOnlineVectorIndexReplica                        SecurableKind = "TABLE_ONLINE_VECTOR_INDEX_REPLICA"
	SecurableKind_TableOnlineVectorIndexDirect                         SecurableKind = "TABLE_ONLINE_VECTOR_INDEX_DIRECT"
	SecurableKind_TableOnlineView                                      SecurableKind = "TABLE_ONLINE_VIEW"
	SecurableKind_TableDbStorage                                       SecurableKind = "TABLE_DB_STORAGE"
	SecurableKind_TableManagedPostgresql                               SecurableKind = "TABLE_MANAGED_POSTGRESQL"
)

// String representation for [fmt.Print].
func (f *SecurableKind) String() string {
	return string(*f)
}

type SecurableType string

const (
	SecurableType_Catalog           SecurableType = "CATALOG"
	SecurableType_Schema            SecurableType = "SCHEMA"
	SecurableType_Table             SecurableType = "TABLE"
	SecurableType_StorageCredential SecurableType = "STORAGE_CREDENTIAL"
	SecurableType_ExternalLocation  SecurableType = "EXTERNAL_LOCATION"
	SecurableType_Function          SecurableType = "FUNCTION"
	SecurableType_Share             SecurableType = "SHARE"
	SecurableType_Provider          SecurableType = "PROVIDER"
	SecurableType_Recipient         SecurableType = "RECIPIENT"
	SecurableType_CleanRoom         SecurableType = "CLEAN_ROOM"
	SecurableType_Metastore         SecurableType = "METASTORE"
	SecurableType_Pipeline          SecurableType = "PIPELINE"
	SecurableType_Volume            SecurableType = "VOLUME"
	SecurableType_Connection        SecurableType = "CONNECTION"
	SecurableType_Credential        SecurableType = "CREDENTIAL"
	SecurableType_ExternalMetadata  SecurableType = "EXTERNAL_METADATA"
	SecurableType_StagingTable      SecurableType = "STAGING_TABLE"
)

// String representation for [fmt.Print].
func (f *SecurableType) String() string {
	return string(*f)
}

type SseEncryptionAlgorithm string

const (
	SseEncryptionAlgorithm_SseEncryptionAlgorithmUnspecified SseEncryptionAlgorithm = "SSE_ENCRYPTION_ALGORITHM_UNSPECIFIED"
	SseEncryptionAlgorithm_AwsSseS3                          SseEncryptionAlgorithm = "AWS_SSE_S3"
	SseEncryptionAlgorithm_AwsSseKms                         SseEncryptionAlgorithm = "AWS_SSE_KMS"
)

// String representation for [fmt.Print].
func (f *SseEncryptionAlgorithm) String() string {
	return string(*f)
}

type TableType string

const (
	TableType_Managed              TableType = "MANAGED"
	TableType_External             TableType = "EXTERNAL"
	TableType_View                 TableType = "VIEW"
	TableType_MaterializedView     TableType = "MATERIALIZED_VIEW"
	TableType_StreamingTable       TableType = "STREAMING_TABLE"
	TableType_ManagedShallowClone  TableType = "MANAGED_SHALLOW_CLONE"
	TableType_Foreign              TableType = "FOREIGN"
	TableType_ExternalShallowClone TableType = "EXTERNAL_SHALLOW_CLONE"
	TableType_MetricView           TableType = "METRIC_VIEW"
)

// String representation for [fmt.Print].
func (f *TableType) String() string {
	return string(*f)
}

type OptionSpec_OauthStage string

const (
	OptionSpec_OauthStage_OauthStageUnspecified   OptionSpec_OauthStage = "OAUTH_STAGE_UNSPECIFIED"
	OptionSpec_OauthStage_BeforeAuthorizationCode OptionSpec_OauthStage = "BEFORE_AUTHORIZATION_CODE"
	OptionSpec_OauthStage_BeforeAccessToken       OptionSpec_OauthStage = "BEFORE_ACCESS_TOKEN"
)

// String representation for [fmt.Print].
func (f *OptionSpec_OauthStage) String() string {
	return string(*f)
}

type OptionSpec_OptionType string

const (
	OptionSpec_OptionType_OptionTypeUnspecified   OptionSpec_OptionType = "OPTION_TYPE_UNSPECIFIED"
	OptionSpec_OptionType_OptionBoolean           OptionSpec_OptionType = "OPTION_BOOLEAN"
	OptionSpec_OptionType_OptionNumber            OptionSpec_OptionType = "OPTION_NUMBER"
	OptionSpec_OptionType_OptionBigint            OptionSpec_OptionType = "OPTION_BIGINT"
	OptionSpec_OptionType_OptionString            OptionSpec_OptionType = "OPTION_STRING"
	OptionSpec_OptionType_OptionEnum              OptionSpec_OptionType = "OPTION_ENUM"
	OptionSpec_OptionType_OptionServiceCredential OptionSpec_OptionType = "OPTION_SERVICE_CREDENTIAL"
	OptionSpec_OptionType_OptionMultilineString   OptionSpec_OptionType = "OPTION_MULTILINE_STRING"
	OptionSpec_OptionType_OptionStorageCredential OptionSpec_OptionType = "OPTION_STORAGE_CREDENTIAL"
)

// String representation for [fmt.Print].
func (f *OptionSpec_OptionType) String() string {
	return string(*f)
}

type ColumnInfo struct {
	Name             *string         `json:"name"`
	TypeText         *string         `json:"type_text"`
	TypeName         *ColumnTypeName `json:"type_name"`
	Position         *int            `json:"position"`
	TypePrecision    *int            `json:"type_precision"`
	TypeScale        *int            `json:"type_scale"`
	TypeIntervalType *string         `json:"type_interval_type"`
	TypeJson         *string         `json:"type_json"`
	Comment          *string         `json:"comment"`
	Nullable         *bool           `json:"nullable"`
	PartitionIndex   *int            `json:"partition_index"`
	Mask             *ColumnMask     `json:"mask"`
}

type ColumnMask struct {
	FunctionName     *string                  `json:"function_name"`
	UsingColumnNames []string                 `json:"using_column_names"`
	UsingArguments   []PolicyFunctionArgument `json:"using_arguments"`
}

// Defines when an option should be hidden based on another option's value. For
// example, for pre-created OAuth connections, some options are conditionally
// hidden. This field works in conjunction with OptionSpec.is_hidden: - If
// OptionSpec.is_hidden is true, the option is always hidden regardless of
// ConditionalDisplay. - If OptionSpec.is_hidden is false (or unset),
// ConditionalDisplay determines visibility: - If depends_on_option matches any
// value in hidden_when_values, hide this option. - Otherwise, show this option..
type ConditionalDisplay struct {
	DependsOnOption  *string  `json:"depends_on_option"`
	HiddenWhenValues []string `json:"hidden_when_values"`
}

// A connection that is dependent on a SQL object..
type ConnectionDependency struct {
	ConnectionName *string `json:"connection_name"`
}

type CreateTable struct {
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	SchemaName                          *string                              `json:"schema_name"`
	TableType                           *TableType                           `json:"table_type"`
	DataSourceFormat                    *DataSourceFormat                    `json:"data_source_format"`
	StorageLocation                     *string                              `json:"storage_location"`
	ViewDefinition                      *string                              `json:"view_definition"`
	ViewDependencies                    *DependencyList                      `json:"view_dependencies"`
	SqlPath                             *string                              `json:"sql_path"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageCredentialName               *string                              `json:"storage_credential_name"`
	TableConstraints                    []TableConstraint                    `json:"table_constraints"`
	RowFilter                           *RowFilter                           `json:"row_filter"`
	PipelineId                          *string                              `json:"pipeline_id"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	DataAccessConfigurationId           *string                              `json:"data_access_configuration_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	TableId                             *string                              `json:"table_id"`
	DeltaRuntimePropertiesKvpairs       *DeltaRuntimePropertiesKvPairs       `json:"delta_runtime_properties_kvpairs"`
	DeletedAt                           *int64                               `json:"deleted_at"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	AccessPoint                         *string                              `json:"access_point"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	EncryptionDetails                   *EncryptionDetails                   `json:"encryption_details"`
	SecurableKindManifest               *SecurableKindManifest               `json:"securable_kind_manifest"`
	Columns                             []ColumnInfo                         `json:"columns"`
	Properties                          map[string]string                    `json:"properties"`
}

type CreateTable_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateTableConstraint struct {
	FullNameArg *string          `json:"full_name_arg"`
	Constraint  *TableConstraint `json:"constraint"`
}

// A credential that is dependent on a SQL object..
type CredentialDependency struct {
	CredentialName *string `json:"credential_name"`
}

type DeleteTable struct {
	FullNameArg *string `json:"full_name_arg"`
}

type DeleteTable_Response struct {
}

type DeleteTableConstraint struct {
	FullNameArg    *string `json:"full_name_arg"`
	ConstraintName *string `json:"constraint_name"`
	Cascade        *bool   `json:"cascade"`
}

type DeleteTableConstraint_Response struct {
}

// Properties pertaining to the current state of the delta table as given by the
// commit server. This does not contain **delta.*** (input) properties in
// __TableInfo.properties__..
type DeltaRuntimePropertiesKvPairs struct {
	DeltaRuntimeProperties map[string]string `json:"delta_runtime_properties"`
}

type DeltaRuntimePropertiesKvPairs_DeltaRuntimePropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
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

type EffectivePredictiveOptimizationFlag struct {
	Value             *string `json:"value"`
	InheritedFromType *string `json:"inherited_from_type"`
	InheritedFromName *string `json:"inherited_from_name"`
}

// Encryption options that apply to clients connecting to cloud storage..
type EncryptionDetails struct {
	SseEncryptionDetails *SseEncryptionDetails `json:"sse_encryption_details"`
}

type ForeignKeyConstraint struct {
	Name          *string  `json:"name"`
	ChildColumns  []string `json:"child_columns"`
	ParentTable   *string  `json:"parent_table"`
	ParentColumns []string `json:"parent_columns"`
	Rely          *bool    `json:"rely"`
}

// A function that is dependent on a SQL object..
type FunctionDependency struct {
	FunctionFullName *string `json:"function_full_name"`
}

type GetTable struct {
	FullNameArg                 *string `json:"full_name_arg"`
	IncludeDeltaMetadata        *bool   `json:"include_delta_metadata"`
	IncludeBrowse               *bool   `json:"include_browse"`
	IncludeManifestCapabilities *bool   `json:"include_manifest_capabilities"`
}

type ListTableSummaries struct {
	CatalogName                 *string `json:"catalog_name"`
	SchemaNamePattern           *string `json:"schema_name_pattern"`
	TableNamePattern            *string `json:"table_name_pattern"`
	MaxResults                  *int    `json:"max_results"`
	PageToken                   *string `json:"page_token"`
	IncludeManifestCapabilities *bool   `json:"include_manifest_capabilities"`
}

type ListTableSummaries_Response struct {
	Tables        []TableSummary `json:"tables"`
	NextPageToken *string        `json:"next_page_token"`
}

type ListTables struct {
	CatalogName                 *string `json:"catalog_name"`
	SchemaName                  *string `json:"schema_name"`
	MaxResults                  *int    `json:"max_results"`
	PageToken                   *string `json:"page_token"`
	OmitColumns                 *bool   `json:"omit_columns"`
	OmitProperties              *bool   `json:"omit_properties"`
	OmitUsername                *bool   `json:"omit_username"`
	IncludeBrowse               *bool   `json:"include_browse"`
	IncludeManifestCapabilities *bool   `json:"include_manifest_capabilities"`
}

type ListTables_Response struct {
	Tables        []TableInfo `json:"tables"`
	NextPageToken *string     `json:"next_page_token"`
}

type NamedTableConstraint struct {
	Name *string `json:"name"`
}

// Spec of an allowed option on a securable kind and its attributes. This is
// mostly used by UI to provide user friendly hints and descriptions in order to
// facilitate the securable creation process..
type OptionSpec struct {
	Name               *string                `json:"name"`
	Type               *OptionSpec_OptionType `json:"type"`
	DefaultValue       *string                `json:"default_value"`
	AllowedValues      []string               `json:"allowed_values"`
	Hint               *string                `json:"hint"`
	Description        *string                `json:"description"`
	IsRequired         *bool                  `json:"is_required"`
	IsSecret           *bool                  `json:"is_secret"`
	IsHidden           *bool                  `json:"is_hidden"`
	IsUpdatable        *bool                  `json:"is_updatable"`
	OauthStage         *OptionSpec_OauthStage `json:"oauth_stage"`
	IsLoggable         *bool                  `json:"is_loggable"`
	IsCreatable        *bool                  `json:"is_creatable"`
	IsCopiable         *bool                  `json:"is_copiable"`
	ConditionalDisplay *ConditionalDisplay    `json:"conditional_display"`
}

// A positional argument passed to a row filter or column mask function.
// Distinguishes between column references and literals..
type PolicyFunctionArgument struct {
	Column   *string `json:"column"`
	Constant *string `json:"constant"`
}

type PrimaryKeyConstraint struct {
	Name              *string  `json:"name"`
	ChildColumns      []string `json:"child_columns"`
	TimeseriesColumns []string `json:"timeseries_columns"`
	Rely              *bool    `json:"rely"`
}

type RowFilter struct {
	FunctionName     *string                  `json:"function_name"`
	InputColumnNames []string                 `json:"input_column_names"`
	InputArguments   []PolicyFunctionArgument `json:"input_arguments"`
}

// Manifest of a specific securable kind..
type SecurableKindManifest struct {
	SecurableType        *SecurableType `json:"securable_type"`
	SecurableKind        *SecurableKind `json:"securable_kind"`
	AssignablePrivileges []string       `json:"assignable_privileges"`
	Options              []OptionSpec   `json:"options"`
	Capabilities         []string       `json:"capabilities"`
}

// Server-Side Encryption properties for clients communicating with AWS s3..
type SseEncryptionDetails struct {
	Algorithm    *SseEncryptionAlgorithm `json:"algorithm"`
	AwsKmsKeyArn *string                 `json:"aws_kms_key_arn"`
}

// A table constraint, as defined by *one* of the following fields being set:
// __primary_key_constraint__, __foreign_key_constraint__,
// __named_table_constraint__..
type TableConstraint struct {
	PrimaryKeyConstraint *PrimaryKeyConstraint `json:"primary_key_constraint"`
	ForeignKeyConstraint *ForeignKeyConstraint `json:"foreign_key_constraint"`
	NamedTableConstraint *NamedTableConstraint `json:"named_table_constraint"`
}

// A table that is dependent on a SQL object..
type TableDependency struct {
	TableFullName *string `json:"table_full_name"`
}

type TableExists struct {
	FullNameArg *string `json:"full_name_arg"`
}

type TableExists_Response struct {
	TableExists *bool `json:"table_exists"`
}

type TableInfo struct {
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	SchemaName                          *string                              `json:"schema_name"`
	TableType                           *TableType                           `json:"table_type"`
	DataSourceFormat                    *DataSourceFormat                    `json:"data_source_format"`
	StorageLocation                     *string                              `json:"storage_location"`
	ViewDefinition                      *string                              `json:"view_definition"`
	ViewDependencies                    *DependencyList                      `json:"view_dependencies"`
	SqlPath                             *string                              `json:"sql_path"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageCredentialName               *string                              `json:"storage_credential_name"`
	TableConstraints                    []TableConstraint                    `json:"table_constraints"`
	RowFilter                           *RowFilter                           `json:"row_filter"`
	PipelineId                          *string                              `json:"pipeline_id"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	DataAccessConfigurationId           *string                              `json:"data_access_configuration_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	TableId                             *string                              `json:"table_id"`
	DeltaRuntimePropertiesKvpairs       *DeltaRuntimePropertiesKvPairs       `json:"delta_runtime_properties_kvpairs"`
	DeletedAt                           *int64                               `json:"deleted_at"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	AccessPoint                         *string                              `json:"access_point"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	EncryptionDetails                   *EncryptionDetails                   `json:"encryption_details"`
	SecurableKindManifest               *SecurableKindManifest               `json:"securable_kind_manifest"`
	Columns                             []ColumnInfo                         `json:"columns"`
	Properties                          map[string]string                    `json:"properties"`
}

type TableInfo_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type TableSummary struct {
	FullName              *string                `json:"full_name"`
	TableType             *TableType             `json:"table_type"`
	SecurableKindManifest *SecurableKindManifest `json:"securable_kind_manifest"`
}

type UpdateTable struct {
	FullNameArg                         *string                              `json:"full_name_arg"`
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	SchemaName                          *string                              `json:"schema_name"`
	TableType                           *TableType                           `json:"table_type"`
	DataSourceFormat                    *DataSourceFormat                    `json:"data_source_format"`
	StorageLocation                     *string                              `json:"storage_location"`
	ViewDefinition                      *string                              `json:"view_definition"`
	ViewDependencies                    *DependencyList                      `json:"view_dependencies"`
	SqlPath                             *string                              `json:"sql_path"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageCredentialName               *string                              `json:"storage_credential_name"`
	TableConstraints                    []TableConstraint                    `json:"table_constraints"`
	RowFilter                           *RowFilter                           `json:"row_filter"`
	PipelineId                          *string                              `json:"pipeline_id"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	DataAccessConfigurationId           *string                              `json:"data_access_configuration_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	TableId                             *string                              `json:"table_id"`
	DeltaRuntimePropertiesKvpairs       *DeltaRuntimePropertiesKvPairs       `json:"delta_runtime_properties_kvpairs"`
	DeletedAt                           *int64                               `json:"deleted_at"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	AccessPoint                         *string                              `json:"access_point"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	EncryptionDetails                   *EncryptionDetails                   `json:"encryption_details"`
	SecurableKindManifest               *SecurableKindManifest               `json:"securable_kind_manifest"`
	Columns                             []ColumnInfo                         `json:"columns"`
	Properties                          map[string]string                    `json:"properties"`
}

type UpdateTable_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateTable_Response struct {
}
