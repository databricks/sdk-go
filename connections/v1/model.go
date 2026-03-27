package v1

type ConnectionType string

const (
	ConnectionType_UnknownConnectionType ConnectionType = "UNKNOWN_CONNECTION_TYPE"
	ConnectionType_Mysql                 ConnectionType = "MYSQL"
	ConnectionType_Postgresql            ConnectionType = "POSTGRESQL"
	ConnectionType_Snowflake             ConnectionType = "SNOWFLAKE"
	ConnectionType_Redshift              ConnectionType = "REDSHIFT"
	ConnectionType_Sqldw                 ConnectionType = "SQLDW"
	ConnectionType_Sqlserver             ConnectionType = "SQLSERVER"
	ConnectionType_Databricks            ConnectionType = "DATABRICKS"
	ConnectionType_Salesforce            ConnectionType = "SALESFORCE"
	ConnectionType_Bigquery              ConnectionType = "BIGQUERY"
	ConnectionType_WorkdayRaas           ConnectionType = "WORKDAY_RAAS"
	ConnectionType_HiveMetastore         ConnectionType = "HIVE_METASTORE"
	ConnectionType_Ga4RawData            ConnectionType = "GA4_RAW_DATA"
	ConnectionType_Servicenow            ConnectionType = "SERVICENOW"
	ConnectionType_SalesforceDataCloud   ConnectionType = "SALESFORCE_DATA_CLOUD"
	ConnectionType_Glue                  ConnectionType = "GLUE"
	ConnectionType_Oracle                ConnectionType = "ORACLE"
	ConnectionType_Palantir              ConnectionType = "PALANTIR"
	ConnectionType_Teradata              ConnectionType = "TERADATA"
	ConnectionType_Http                  ConnectionType = "HTTP"
	ConnectionType_PowerBi               ConnectionType = "POWER_BI"
)

// String representation for [fmt.Print].
func (f *ConnectionType) String() string {
	return string(*f)
}

type CredentialType string

const (
	CredentialType_UnknownCredentialType      CredentialType = "UNKNOWN_CREDENTIAL_TYPE"
	CredentialType_UsernamePassword           CredentialType = "USERNAME_PASSWORD"
	CredentialType_OauthU2m                   CredentialType = "OAUTH_U2M"
	CredentialType_OauthM2m                   CredentialType = "OAUTH_M2M"
	CredentialType_OauthRefreshToken          CredentialType = "OAUTH_REFRESH_TOKEN"
	CredentialType_OauthAccessToken           CredentialType = "OAUTH_ACCESS_TOKEN"
	CredentialType_OauthResourceOwnerPassword CredentialType = "OAUTH_RESOURCE_OWNER_PASSWORD"
	CredentialType_ServiceCredential          CredentialType = "SERVICE_CREDENTIAL"
	CredentialType_BearerToken                CredentialType = "BEARER_TOKEN"
	CredentialType_OidcToken                  CredentialType = "OIDC_TOKEN"
	CredentialType_PemPrivateKey              CredentialType = "PEM_PRIVATE_KEY"
	CredentialType_OauthU2mMapping            CredentialType = "OAUTH_U2M_MAPPING"
	CredentialType_AnyStaticCredential        CredentialType = "ANY_STATIC_CREDENTIAL"
	CredentialType_OauthMtls                  CredentialType = "OAUTH_MTLS"
	CredentialType_SswsToken                  CredentialType = "SSWS_TOKEN"
	CredentialType_EdgegridAkamai             CredentialType = "EDGEGRID_AKAMAI"
	CredentialType_OauthGoogleServiceAccount  CredentialType = "OAUTH_GOOGLE_SERVICE_ACCOUNT"
)

// String representation for [fmt.Print].
func (f *CredentialType) String() string {
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

type ProvisioningInfo_State string

const (
	ProvisioningInfo_State_StateUnspecified ProvisioningInfo_State = "STATE_UNSPECIFIED"
	ProvisioningInfo_State_Provisioning     ProvisioningInfo_State = "PROVISIONING"
	ProvisioningInfo_State_Active           ProvisioningInfo_State = "ACTIVE"
	ProvisioningInfo_State_Failed           ProvisioningInfo_State = "FAILED"
	ProvisioningInfo_State_Deleting         ProvisioningInfo_State = "DELETING"
	ProvisioningInfo_State_Updating         ProvisioningInfo_State = "UPDATING"
	ProvisioningInfo_State_Degraded         ProvisioningInfo_State = "DEGRADED"
)

// String representation for [fmt.Print].
func (f *ProvisioningInfo_State) String() string {
	return string(*f)
}

// Next ID: 24.
type ConnectionInfo struct {
	Name                *string              `json:"name"`
	ConnectionType      *ConnectionType      `json:"connection_type"`
	Owner               *string              `json:"owner"`
	ReadOnly            *bool                `json:"read_only"`
	Comment             *string              `json:"comment"`
	EnvironmentSettings *EnvironmentSettings `json:"environment_settings"`
	FullName            *string              `json:"full_name"`
	Url                 *string              `json:"url"`
	CredentialType      *CredentialType      `json:"credential_type"`
	ConnectionId        *string              `json:"connection_id"`
	MetastoreId         *string              `json:"metastore_id"`
	CreatedAt           *int64               `json:"created_at"`
	CreatedBy           *string              `json:"created_by"`
	UpdatedAt           *int64               `json:"updated_at"`
	UpdatedBy           *string              `json:"updated_by"`
	SecurableType       *SecurableType       `json:"securable_type"`
	ProvisioningInfo    *ProvisioningInfo    `json:"provisioning_info"`
	Options             map[string]string    `json:"options"`
	Properties          map[string]string    `json:"properties"`
}

type ConnectionInfo_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ConnectionInfo_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ConnectionUserCredential struct {
	UserIdentity     *string           `json:"user_identity"`
	OptionsKvpairs   *OptionsKvPairs   `json:"options_kvpairs"`
	ProvisioningInfo *ProvisioningInfo `json:"provisioning_info"`
}

type CreateConnection struct {
	Name                *string              `json:"name"`
	ConnectionType      *ConnectionType      `json:"connection_type"`
	Owner               *string              `json:"owner"`
	ReadOnly            *bool                `json:"read_only"`
	Comment             *string              `json:"comment"`
	EnvironmentSettings *EnvironmentSettings `json:"environment_settings"`
	FullName            *string              `json:"full_name"`
	Url                 *string              `json:"url"`
	CredentialType      *CredentialType      `json:"credential_type"`
	ConnectionId        *string              `json:"connection_id"`
	MetastoreId         *string              `json:"metastore_id"`
	CreatedAt           *int64               `json:"created_at"`
	CreatedBy           *string              `json:"created_by"`
	UpdatedAt           *int64               `json:"updated_at"`
	UpdatedBy           *string              `json:"updated_by"`
	SecurableType       *SecurableType       `json:"securable_type"`
	ProvisioningInfo    *ProvisioningInfo    `json:"provisioning_info"`
	Options             map[string]string    `json:"options"`
	Properties          map[string]string    `json:"properties"`
}

type CreateConnection_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateConnection_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateUserMappedCredential struct {
	NameArg                  *string                   `json:"name_arg"`
	ConnectionUserCredential *ConnectionUserCredential `json:"connection_user_credential"`
}

type CreateUserMappedCredential_Response struct {
	ConnectionUserCredential *ConnectionUserCredential `json:"connection_user_credential"`
}

type DeleteConnection struct {
	NameArg *string `json:"name_arg"`
}

type DeleteConnection_Response struct {
}

type DeleteUserMappedCredential struct {
	NameArg         *string `json:"name_arg"`
	UserIdentityArg *string `json:"user_identity_arg"`
}

type DeleteUserMappedCredential_Response struct {
}

type EnvironmentSettings struct {
	JavaDependencies   []string `json:"java_dependencies"`
	EnvironmentVersion *string  `json:"environment_version"`
}

type GetConnection struct {
	NameArg *string `json:"name_arg"`
}

type GetUserMappedCredential struct {
	NameArg         *string `json:"name_arg"`
	UserIdentityArg *string `json:"user_identity_arg"`
}

type GetUserMappedCredential_Response struct {
	ConnectionUserCredential *ConnectionUserCredential `json:"connection_user_credential"`
}

type ListConnections struct {
	MaxResults *int    `json:"max_results"`
	PageToken  *string `json:"page_token"`
}

type ListConnections_Response struct {
	Connections   []ConnectionInfo `json:"connections"`
	NextPageToken *string          `json:"next_page_token"`
}

type ListUserMappedCredentials struct {
	NameArg      *string `json:"name_arg"`
	UserIdentity *string `json:"user_identity"`
	MaxResults   *int    `json:"max_results"`
	PageToken    *string `json:"page_token"`
}

type ListUserMappedCredentials_Response struct {
	ConnectionUserCredential []ConnectionUserCredential `json:"connection_user_credential"`
	NextPageToken            *string                    `json:"next_page_token"`
}

type OptionsKvPairs struct {
	Options map[string]string `json:"options"`
}

type OptionsKvPairs_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Status of an asynchronously provisioned resource..
type ProvisioningInfo struct {
	State *ProvisioningInfo_State `json:"state"`
}

type UpdateConnection struct {
	NameArg             *string              `json:"name_arg"`
	NewName             *string              `json:"new_name"`
	Name                *string              `json:"name"`
	ConnectionType      *ConnectionType      `json:"connection_type"`
	Owner               *string              `json:"owner"`
	ReadOnly            *bool                `json:"read_only"`
	Comment             *string              `json:"comment"`
	EnvironmentSettings *EnvironmentSettings `json:"environment_settings"`
	FullName            *string              `json:"full_name"`
	Url                 *string              `json:"url"`
	CredentialType      *CredentialType      `json:"credential_type"`
	ConnectionId        *string              `json:"connection_id"`
	MetastoreId         *string              `json:"metastore_id"`
	CreatedAt           *int64               `json:"created_at"`
	CreatedBy           *string              `json:"created_by"`
	UpdatedAt           *int64               `json:"updated_at"`
	UpdatedBy           *string              `json:"updated_by"`
	SecurableType       *SecurableType       `json:"securable_type"`
	ProvisioningInfo    *ProvisioningInfo    `json:"provisioning_info"`
	Options             map[string]string    `json:"options"`
	Properties          map[string]string    `json:"properties"`
}

type UpdateConnection_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateConnection_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateUserMappedCredential struct {
	NameArg                  *string                   `json:"name_arg"`
	ConnectionUserCredential *ConnectionUserCredential `json:"connection_user_credential"`
}

type UpdateUserMappedCredential_Response struct {
	ConnectionUserCredential *ConnectionUserCredential `json:"connection_user_credential"`
}
