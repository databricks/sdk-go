package v1

type IsolationMode string

const (
	IsolationMode_IsolationModeUnspecified   IsolationMode = "ISOLATION_MODE_UNSPECIFIED"
	IsolationMode_IsolationModeOpen          IsolationMode = "ISOLATION_MODE_OPEN"
	IsolationMode_IsolationModeIsolated      IsolationMode = "ISOLATION_MODE_ISOLATED"
	IsolationMode_IsolationModeOpenInAccount IsolationMode = "ISOLATION_MODE_OPEN_IN_ACCOUNT"
)

// String representation for [fmt.Print].
func (f *IsolationMode) String() string {
	return string(*f)
}

type PathOperation string

const (
	PathOperation_PathRead        PathOperation = "PATH_READ"
	PathOperation_PathReadWrite   PathOperation = "PATH_READ_WRITE"
	PathOperation_PathCreateTable PathOperation = "PATH_CREATE_TABLE"
)

// String representation for [fmt.Print].
func (f *PathOperation) String() string {
	return string(*f)
}

type TableOperation string

const (
	TableOperation_Read      TableOperation = "READ"
	TableOperation_ReadWrite TableOperation = "READ_WRITE"
)

// String representation for [fmt.Print].
func (f *TableOperation) String() string {
	return string(*f)
}

type ValidateCredential_Result string

const (
	ValidateCredential_Result_Pass ValidateCredential_Result = "PASS"
	ValidateCredential_Result_Fail ValidateCredential_Result = "FAIL"
	ValidateCredential_Result_Skip ValidateCredential_Result = "SKIP"
)

// String representation for [fmt.Print].
func (f *ValidateCredential_Result) String() string {
	return string(*f)
}

type ValidateStorageCredential_FileOperation string

const (
	ValidateStorageCredential_FileOperation_List       ValidateStorageCredential_FileOperation = "LIST"
	ValidateStorageCredential_FileOperation_Read       ValidateStorageCredential_FileOperation = "READ"
	ValidateStorageCredential_FileOperation_Write      ValidateStorageCredential_FileOperation = "WRITE"
	ValidateStorageCredential_FileOperation_Delete     ValidateStorageCredential_FileOperation = "DELETE"
	ValidateStorageCredential_FileOperation_PathExists ValidateStorageCredential_FileOperation = "PATH_EXISTS"
)

// String representation for [fmt.Print].
func (f *ValidateStorageCredential_FileOperation) String() string {
	return string(*f)
}

type ValidateStorageCredential_Result string

const (
	ValidateStorageCredential_Result_Pass ValidateStorageCredential_Result = "PASS"
	ValidateStorageCredential_Result_Fail ValidateStorageCredential_Result = "FAIL"
	ValidateStorageCredential_Result_Skip ValidateStorageCredential_Result = "SKIP"
)

// String representation for [fmt.Print].
func (f *ValidateStorageCredential_Result) String() string {
	return string(*f)
}

// AWS temporary credentials for API authentication. Read more at
// https://docs.aws.amazon.com/STS/latest/APIReference/API_Credentials.html..
type AwsCredentials struct {
	AccessKeyId     *string `json:"access_key_id"`
	SecretAccessKey *string `json:"secret_access_key"`
	SessionToken    *string `json:"session_token"`
	AccessPoint     *string `json:"access_point"`
}

// The AWS IAM role configuration.
type AwsIamRole struct {
	RoleArn            *string `json:"role_arn"`
	UnityCatalogIamArn *string `json:"unity_catalog_iam_arn"`
	ExternalId         *string `json:"external_id"`
}

// Azure Active Directory token, essentially the Oauth token for Azure Service
// Principal or Managed Identity. Read more at
// https://learn.microsoft.com/en-us/azure/databricks/dev-tools/api/latest/aad/service-prin-aad-token.
type AzureActiveDirectoryToken struct {
	AadToken *string `json:"aad_token"`
}

// The Azure managed identity configuration..
type AzureManagedIdentity struct {
	AccessConnectorId *string `json:"access_connector_id"`
	ManagedIdentityId *string `json:"managed_identity_id"`
	CredentialId      *string `json:"credential_id"`
}

// The Azure service principal configuration. Only applicable when purpose is
// **STORAGE**..
type AzureServicePrincipal struct {
	DirectoryId   *string `json:"directory_id"`
	ApplicationId *string `json:"application_id"`
	ClientSecret  *string `json:"client_secret"`
}

// Azure temporary credentials for API authentication. Read more at
// https://docs.microsoft.com/en-us/rest/api/storageservices/create-user-delegation-sas.
type AzureUserDelegationSas struct {
	SasToken *string `json:"sas_token"`
}

// The Cloudflare API token configuration. Read more at
// https://developers.cloudflare.com/r2/api/s3/tokens/.
type CloudflareApiToken struct {
	AccessKeyId     *string `json:"access_key_id"`
	SecretAccessKey *string `json:"secret_access_key"`
	AccountId       *string `json:"account_id"`
}

type CreateCredential struct {
	SkipValidation              *bool                        `json:"skip_validation"`
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

type CreateStorageCredential struct {
	SkipValidation              *bool                        `json:"skip_validation"`
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

type CredentialInfo struct {
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

// GCP long-lived credential. <Databricks>-created Google Cloud Storage service
// account..
type DatabricksGcpServiceAccount struct {
	Email        *string `json:"email"`
	PrivateKeyId *string `json:"private_key_id"`
	CredentialId *string `json:"credential_id"`
}

type DeleteCredential struct {
	NameArg *string `json:"name_arg"`
	Force   *bool   `json:"force"`
}

type DeleteCredential_Response struct {
}

type DeleteStorageCredential struct {
	NameArg *string `json:"name_arg"`
	Force   *bool   `json:"force"`
}

type DeleteStorageCredential_Response struct {
}

// GCP temporary credentials for API authentication. Read more at
// https://developers.google.com/identity/protocols/oauth2/service-account.
type GcpOauthToken struct {
	OauthToken *string `json:"oauth_token"`
}

// GCP long-lived credential. GCP Service Account..
type GcpServiceAccountKey struct {
	Email        *string `json:"email"`
	PrivateKeyId *string `json:"private_key_id"`
	PrivateKey   *string `json:"private_key"`
}

type GenerateTemporaryPathCredential struct {
	Url       *string        `json:"url"`
	Operation *PathOperation `json:"operation"`
	DryRun    *bool          `json:"dry_run"`
}

type GenerateTemporaryPathCredential_Response struct {
	AwsTempCredentials     *AwsCredentials            `json:"aws_temp_credentials"`
	AzureUserDelegationSas *AzureUserDelegationSas    `json:"azure_user_delegation_sas"`
	GcpOauthToken          *GcpOauthToken             `json:"gcp_oauth_token"`
	AzureAad               *AzureActiveDirectoryToken `json:"azure_aad"`
	R2TempCredentials      *R2Credentials             `json:"r2_temp_credentials"`
	UcEncryptedToken       *UcEncryptedToken          `json:"uc_encrypted_token"`
	ExpirationTime         *int64                     `json:"expiration_time"`
	Url                    *string                    `json:"url"`
}

type GenerateTemporaryServiceCredential struct {
	CredentialName *string                                          `json:"credential_name"`
	AzureOptions   *GenerateTemporaryServiceCredential_AzureOptions `json:"azure_options"`
	GcpOptions     *GenerateTemporaryServiceCredential_GcpOptions   `json:"gcp_options"`
}

// The Azure cloud options to customize the requested temporary credential.
type GenerateTemporaryServiceCredential_AzureOptions struct {
	Resources []string `json:"resources"`
}

// The GCP cloud options to customize the requested temporary credential.
type GenerateTemporaryServiceCredential_GcpOptions struct {
	Scopes []string `json:"scopes"`
}

type GenerateTemporaryTableCredential struct {
	TableId   *string         `json:"table_id"`
	Operation *TableOperation `json:"operation"`
}

type GenerateTemporaryTableCredential_Response struct {
	AwsTempCredentials     *AwsCredentials            `json:"aws_temp_credentials"`
	AzureUserDelegationSas *AzureUserDelegationSas    `json:"azure_user_delegation_sas"`
	GcpOauthToken          *GcpOauthToken             `json:"gcp_oauth_token"`
	AzureAad               *AzureActiveDirectoryToken `json:"azure_aad"`
	R2TempCredentials      *R2Credentials             `json:"r2_temp_credentials"`
	UcEncryptedToken       *UcEncryptedToken          `json:"uc_encrypted_token"`
	ExpirationTime         *int64                     `json:"expiration_time"`
	Url                    *string                    `json:"url"`
}

type GetCredential struct {
	NameArg *string `json:"name_arg"`
}

// TODO(UC-1710): The legacy /storage-credentials API is being deprecated.
// Please use the new consolidated /credentials API instead. See
// https://github.com/databricks-eng/universe/pull/857047#discussion_r1924779791
// for an example of a case when that wasn't possible..
type GetStorageCredential struct {
	NameArg *string `json:"name_arg"`
}

// ListCredentialsRequest is used to list credentials in the metastore. Returns
// an array of credentials (as CredentialInfo objects). The array is limited to
// the credentials that the caller has permission to access. If the caller is a
// metastore admin, retrieval of credentials is unrestricted.
//
// There is no guarantee of a specific ordering of the elements in the array..
type ListCredentials struct {
	IncludeUnbound *bool   `json:"include_unbound"`
	MaxResults     *int    `json:"max_results"`
	PageToken      *string `json:"page_token"`
}

type ListCredentials_Response struct {
	Credentials   []CredentialInfo `json:"credentials"`
	NextPageToken *string          `json:"next_page_token"`
}

type ListStorageCredentials struct {
	IncludeUnbound *bool   `json:"include_unbound"`
	MaxResults     *int    `json:"max_results"`
	PageToken      *string `json:"page_token"`
}

type ListStorageCredentials_Response struct {
	StorageCredentials []StorageCredentialInfo `json:"storage_credentials"`
	NextPageToken      *string                 `json:"next_page_token"`
}

// R2 temporary credentials for API authentication. Read more at
// https://developers.cloudflare.com/r2/api/s3/tokens/..
type R2Credentials struct {
	AccessKeyId     *string `json:"access_key_id"`
	SecretAccessKey *string `json:"secret_access_key"`
	SessionToken    *string `json:"session_token"`
}

type StorageCredentialInfo struct {
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

type TemporaryCredentials struct {
	AwsTempCredentials     *AwsCredentials            `json:"aws_temp_credentials"`
	AzureUserDelegationSas *AzureUserDelegationSas    `json:"azure_user_delegation_sas"`
	GcpOauthToken          *GcpOauthToken             `json:"gcp_oauth_token"`
	AzureAad               *AzureActiveDirectoryToken `json:"azure_aad"`
	R2TempCredentials      *R2Credentials             `json:"r2_temp_credentials"`
	UcEncryptedToken       *UcEncryptedToken          `json:"uc_encrypted_token"`
	ExpirationTime         *int64                     `json:"expiration_time"`
	Url                    *string                    `json:"url"`
}

// Encrypted token used when we cannot downscope the cloud provider token
// appropriately See:
// https://docs.google.com/document/d/1hEKDnSckuU5PIS798CtfqBElrMR6OJuR2wgz_BjhMSY.
type UcEncryptedToken struct {
	EncryptedPayload *string `json:"encrypted_payload"`
}

type UpdateCredential struct {
	NameArg                     *string                      `json:"name_arg"`
	NewName                     *string                      `json:"new_name"`
	SkipValidation              *bool                        `json:"skip_validation"`
	Force                       *bool                        `json:"force"`
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

type UpdateStorageCredential struct {
	NameArg                     *string                      `json:"name_arg"`
	NewName                     *string                      `json:"new_name"`
	SkipValidation              *bool                        `json:"skip_validation"`
	Force                       *bool                        `json:"force"`
	Name                        *string                      `json:"name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	GcpServiceAccountKey        *GcpServiceAccountKey        `json:"gcp_service_account_key"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	Comment                     *string                      `json:"comment"`
	ReadOnly                    *bool                        `json:"read_only"`
	Owner                       *string                      `json:"owner"`
	Id                          *string                      `json:"id"`
	MetastoreId                 *string                      `json:"metastore_id"`
	CreatedAt                   *int64                       `json:"created_at"`
	CreatedBy                   *string                      `json:"created_by"`
	UpdatedAt                   *int64                       `json:"updated_at"`
	UpdatedBy                   *string                      `json:"updated_by"`
	UsedForManagedStorage       *bool                        `json:"used_for_managed_storage"`
	FullName                    *string                      `json:"full_name"`
	IsolationMode               *IsolationMode               `json:"isolation_mode"`
}

// Next ID: 18.
type ValidateCredential struct {
	CredentialName              *string                      `json:"credential_name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	ExternalLocationName        *string                      `json:"external_location_name"`
	Url                         *string                      `json:"url"`
	ReadOnly                    *bool                        `json:"read_only"`
}

type ValidateCredential_Response struct {
	Results []ValidateCredential_ValidationResult `json:"results"`
	IsDir   *bool                                 `json:"isDir"`
}

type ValidateCredential_ValidationResult struct {
	Result  *ValidateCredential_Result `json:"result"`
	Message *string                    `json:"message"`
}

type ValidateStorageCredential struct {
	StorageCredentialName       *string                      `json:"storage_credential_name"`
	AwsIamRole                  *AwsIamRole                  `json:"aws_iam_role"`
	AzureServicePrincipal       *AzureServicePrincipal       `json:"azure_service_principal"`
	AzureManagedIdentity        *AzureManagedIdentity        `json:"azure_managed_identity"`
	DatabricksGcpServiceAccount *DatabricksGcpServiceAccount `json:"databricks_gcp_service_account"`
	CloudflareApiToken          *CloudflareApiToken          `json:"cloudflare_api_token"`
	ExternalLocationName        *string                      `json:"external_location_name"`
	Url                         *string                      `json:"url"`
	ReadOnly                    *bool                        `json:"read_only"`
}

type ValidateStorageCredential_Response struct {
	IsDir   *bool                                        `json:"isDir"`
	Results []ValidateStorageCredential_ValidationResult `json:"results"`
}

type ValidateStorageCredential_ValidationResult struct {
	Operation *ValidateStorageCredential_FileOperation `json:"operation"`
	Result    *ValidateStorageCredential_Result        `json:"result"`
	Message   *string                                  `json:"message"`
}
