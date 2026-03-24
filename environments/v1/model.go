package v1

import (
	"encoding/json"
	"time"
)

type BaseEnvironmentType string

const (
	BaseEnvironmentType_BaseEnvironmentTypeUnspecified BaseEnvironmentType = "BASE_ENVIRONMENT_TYPE_UNSPECIFIED"
	BaseEnvironmentType_Cpu                            BaseEnvironmentType = "CPU"
	BaseEnvironmentType_Gpu                            BaseEnvironmentType = "GPU"
)

// String representation for [fmt.Print].
func (f *BaseEnvironmentType) String() string {
	return string(*f)
}

type ErrorCode string

const (
	ErrorCode_Unknown                          ErrorCode = "UNKNOWN"
	ErrorCode_InternalError                    ErrorCode = "INTERNAL_ERROR"
	ErrorCode_TemporarilyUnavailable           ErrorCode = "TEMPORARILY_UNAVAILABLE"
	ErrorCode_IoError                          ErrorCode = "IO_ERROR"
	ErrorCode_BadRequest                       ErrorCode = "BAD_REQUEST"
	ErrorCode_ServiceUnderMaintenance          ErrorCode = "SERVICE_UNDER_MAINTENANCE"
	ErrorCode_WorkspaceTemporarilyUnavailable  ErrorCode = "WORKSPACE_TEMPORARILY_UNAVAILABLE"
	ErrorCode_DeadlineExceeded                 ErrorCode = "DEADLINE_EXCEEDED"
	ErrorCode_Cancelled                        ErrorCode = "CANCELLED"
	ErrorCode_ResourceExhausted                ErrorCode = "RESOURCE_EXHAUSTED"
	ErrorCode_Aborted                          ErrorCode = "ABORTED"
	ErrorCode_NotFound                         ErrorCode = "NOT_FOUND"
	ErrorCode_AlreadyExists                    ErrorCode = "ALREADY_EXISTS"
	ErrorCode_Unauthenticated                  ErrorCode = "UNAUTHENTICATED"
	ErrorCode_Unavailable                      ErrorCode = "UNAVAILABLE"
	ErrorCode_InvalidParameterValue            ErrorCode = "INVALID_PARAMETER_VALUE"
	ErrorCode_EndpointNotFound                 ErrorCode = "ENDPOINT_NOT_FOUND"
	ErrorCode_MalformedRequest                 ErrorCode = "MALFORMED_REQUEST"
	ErrorCode_InvalidState                     ErrorCode = "INVALID_STATE"
	ErrorCode_PermissionDenied                 ErrorCode = "PERMISSION_DENIED"
	ErrorCode_FeatureDisabled                  ErrorCode = "FEATURE_DISABLED"
	ErrorCode_CustomerUnauthorized             ErrorCode = "CUSTOMER_UNAUTHORIZED"
	ErrorCode_RequestLimitExceeded             ErrorCode = "REQUEST_LIMIT_EXCEEDED"
	ErrorCode_ResourceConflict                 ErrorCode = "RESOURCE_CONFLICT"
	ErrorCode_UnparseableHttpError             ErrorCode = "UNPARSEABLE_HTTP_ERROR"
	ErrorCode_NotImplemented                   ErrorCode = "NOT_IMPLEMENTED"
	ErrorCode_DataLoss                         ErrorCode = "DATA_LOSS"
	ErrorCode_InvalidStateTransition           ErrorCode = "INVALID_STATE_TRANSITION"
	ErrorCode_CouldNotAcquireLock              ErrorCode = "COULD_NOT_ACQUIRE_LOCK"
	ErrorCode_ResourceAlreadyExists            ErrorCode = "RESOURCE_ALREADY_EXISTS"
	ErrorCode_ResourceDoesNotExist             ErrorCode = "RESOURCE_DOES_NOT_EXIST"
	ErrorCode_QuotaExceeded                    ErrorCode = "QUOTA_EXCEEDED"
	ErrorCode_MaxBlockSizeExceeded             ErrorCode = "MAX_BLOCK_SIZE_EXCEEDED"
	ErrorCode_MaxReadSizeExceeded              ErrorCode = "MAX_READ_SIZE_EXCEEDED"
	ErrorCode_PartialDelete                    ErrorCode = "PARTIAL_DELETE"
	ErrorCode_MaxListSizeExceeded              ErrorCode = "MAX_LIST_SIZE_EXCEEDED"
	ErrorCode_DryRunFailed                     ErrorCode = "DRY_RUN_FAILED"
	ErrorCode_ResourceLimitExceeded            ErrorCode = "RESOURCE_LIMIT_EXCEEDED"
	ErrorCode_DirectoryNotEmpty                ErrorCode = "DIRECTORY_NOT_EMPTY"
	ErrorCode_DirectoryProtected               ErrorCode = "DIRECTORY_PROTECTED"
	ErrorCode_MaxNotebookSizeExceeded          ErrorCode = "MAX_NOTEBOOK_SIZE_EXCEEDED"
	ErrorCode_MaxChildNodeSizeExceeded         ErrorCode = "MAX_CHILD_NODE_SIZE_EXCEEDED"
	ErrorCode_SearchQueryTooLong               ErrorCode = "SEARCH_QUERY_TOO_LONG"
	ErrorCode_SearchQueryTooShort              ErrorCode = "SEARCH_QUERY_TOO_SHORT"
	ErrorCode_ManagedResourceGroupDoesNotExist ErrorCode = "MANAGED_RESOURCE_GROUP_DOES_NOT_EXIST"
	ErrorCode_PermissionNotPropagated          ErrorCode = "PERMISSION_NOT_PROPAGATED"
	ErrorCode_DeploymentTimeout                ErrorCode = "DEPLOYMENT_TIMEOUT"
	ErrorCode_GitConflict                      ErrorCode = "GIT_CONFLICT"
	ErrorCode_GitUnknownRef                    ErrorCode = "GIT_UNKNOWN_REF"
	ErrorCode_GitSensitiveTokenDetected        ErrorCode = "GIT_SENSITIVE_TOKEN_DETECTED"
	ErrorCode_GitUrlNotOnAllowList             ErrorCode = "GIT_URL_NOT_ON_ALLOW_LIST"
	ErrorCode_GitRemoteError                   ErrorCode = "GIT_REMOTE_ERROR"
	ErrorCode_ProjectsOperationTimeout         ErrorCode = "PROJECTS_OPERATION_TIMEOUT"
	ErrorCode_IpynbFileInRepo                  ErrorCode = "IPYNB_FILE_IN_REPO"
	ErrorCode_InsecurePartnerResponse          ErrorCode = "INSECURE_PARTNER_RESPONSE"
	ErrorCode_MalformedPartnerResponse         ErrorCode = "MALFORMED_PARTNER_RESPONSE"
	ErrorCode_MetastoreDoesNotExist            ErrorCode = "METASTORE_DOES_NOT_EXIST"
	ErrorCode_DacDoesNotExist                  ErrorCode = "DAC_DOES_NOT_EXIST"
	ErrorCode_CatalogDoesNotExist              ErrorCode = "CATALOG_DOES_NOT_EXIST"
	ErrorCode_SchemaDoesNotExist               ErrorCode = "SCHEMA_DOES_NOT_EXIST"
	ErrorCode_TableDoesNotExist                ErrorCode = "TABLE_DOES_NOT_EXIST"
	ErrorCode_ShareDoesNotExist                ErrorCode = "SHARE_DOES_NOT_EXIST"
	ErrorCode_RecipientDoesNotExist            ErrorCode = "RECIPIENT_DOES_NOT_EXIST"
	ErrorCode_StorageCredentialDoesNotExist    ErrorCode = "STORAGE_CREDENTIAL_DOES_NOT_EXIST"
	ErrorCode_ExternalLocationDoesNotExist     ErrorCode = "EXTERNAL_LOCATION_DOES_NOT_EXIST"
	ErrorCode_PrincipalDoesNotExist            ErrorCode = "PRINCIPAL_DOES_NOT_EXIST"
	ErrorCode_ProviderDoesNotExist             ErrorCode = "PROVIDER_DOES_NOT_EXIST"
	ErrorCode_MetastoreAlreadyExists           ErrorCode = "METASTORE_ALREADY_EXISTS"
	ErrorCode_DacAlreadyExists                 ErrorCode = "DAC_ALREADY_EXISTS"
	ErrorCode_CatalogAlreadyExists             ErrorCode = "CATALOG_ALREADY_EXISTS"
	ErrorCode_SchemaAlreadyExists              ErrorCode = "SCHEMA_ALREADY_EXISTS"
	ErrorCode_TableAlreadyExists               ErrorCode = "TABLE_ALREADY_EXISTS"
	ErrorCode_ShareAlreadyExists               ErrorCode = "SHARE_ALREADY_EXISTS"
	ErrorCode_RecipientAlreadyExists           ErrorCode = "RECIPIENT_ALREADY_EXISTS"
	ErrorCode_StorageCredentialAlreadyExists   ErrorCode = "STORAGE_CREDENTIAL_ALREADY_EXISTS"
	ErrorCode_ExternalLocationAlreadyExists    ErrorCode = "EXTERNAL_LOCATION_ALREADY_EXISTS"
	ErrorCode_ProviderAlreadyExists            ErrorCode = "PROVIDER_ALREADY_EXISTS"
	ErrorCode_CatalogNotEmpty                  ErrorCode = "CATALOG_NOT_EMPTY"
	ErrorCode_SchemaNotEmpty                   ErrorCode = "SCHEMA_NOT_EMPTY"
	ErrorCode_MetastoreNotEmpty                ErrorCode = "METASTORE_NOT_EMPTY"
	ErrorCode_ProviderShareNotAccessible       ErrorCode = "PROVIDER_SHARE_NOT_ACCESSIBLE"
)

// String representation for [fmt.Print].
func (f *ErrorCode) String() string {
	return string(*f)
}

type WorkspaceBaseEnvironmentProvider string

const (
	WorkspaceBaseEnvironmentProvider_WorkspaceBaseEnvironmentProviderUnspecified WorkspaceBaseEnvironmentProvider = "WORKSPACE_BASE_ENVIRONMENT_PROVIDER_UNSPECIFIED"
	WorkspaceBaseEnvironmentProvider_Admin                                       WorkspaceBaseEnvironmentProvider = "ADMIN"
	WorkspaceBaseEnvironmentProvider_Databricks                                  WorkspaceBaseEnvironmentProvider = "DATABRICKS"
)

// String representation for [fmt.Print].
func (f *WorkspaceBaseEnvironmentProvider) String() string {
	return string(*f)
}

type WorkspaceBaseEnvironmentCache_Status string

const (
	WorkspaceBaseEnvironmentCache_Status_StatusUnspecified WorkspaceBaseEnvironmentCache_Status = "STATUS_UNSPECIFIED"
	WorkspaceBaseEnvironmentCache_Status_Pending           WorkspaceBaseEnvironmentCache_Status = "PENDING"
	WorkspaceBaseEnvironmentCache_Status_Created           WorkspaceBaseEnvironmentCache_Status = "CREATED"
	WorkspaceBaseEnvironmentCache_Status_Failed            WorkspaceBaseEnvironmentCache_Status = "FAILED"
	WorkspaceBaseEnvironmentCache_Status_Expired           WorkspaceBaseEnvironmentCache_Status = "EXPIRED"
	WorkspaceBaseEnvironmentCache_Status_Invalid           WorkspaceBaseEnvironmentCache_Status = "INVALID"
	WorkspaceBaseEnvironmentCache_Status_Refreshing        WorkspaceBaseEnvironmentCache_Status = "REFRESHING"
)

// String representation for [fmt.Print].
func (f *WorkspaceBaseEnvironmentCache_Status) String() string {
	return string(*f)
}

// Request message for CreateWorkspaceBaseEnvironment..
type CreateWorkspaceBaseEnvironmentRequest struct {
	WorkspaceBaseEnvironment   *WorkspaceBaseEnvironment `json:"workspace_base_environment"`
	WorkspaceBaseEnvironmentId *string                   `json:"workspace_base_environment_id"`
	RequestId                  *string                   `json:"request_id"`
}

// Databricks Error that is returned by all Databricks APIs..
type DatabricksServiceExceptionWithDetailsProto struct {
	ErrorCode  *ErrorCode        `json:"error_code"`
	Message    *string           `json:"message"`
	StackTrace *string           `json:"stack_trace"`
	Details    []json.RawMessage `json:"details"`
}

// A singleton resource representing the default workspace base environment
// configuration. This resource contains the workspace base environments that
// are used as defaults for serverless notebooks and jobs in the workspace, for
// both CPU and GPU compute types..
type DefaultWorkspaceBaseEnvironment struct {
	Name                        *string `json:"name"`
	CpuWorkspaceBaseEnvironment *string `json:"cpu_workspace_base_environment"`
	GpuWorkspaceBaseEnvironment *string `json:"gpu_workspace_base_environment"`
}

// Request message for DeleteWorkspaceBaseEnvironment..
type DeleteWorkspaceBaseEnvironmentRequest struct {
	Name *string `json:"name"`
}

// Request message for GetDefaultWorkspaceBaseEnvironment..
type GetDefaultWorkspaceBaseEnvironmentRequest struct {
	Name *string `json:"name"`
}

// The request message for `GetOperation` method..
type GetOperationRequest struct {
	Name *string `json:"name"`
}

// Request message for GetWorkspaceBaseEnvironment..
type GetWorkspaceBaseEnvironmentRequest struct {
	Name *string `json:"name"`
}

// Request message for ListWorkspaceBaseEnvironments..
type ListWorkspaceBaseEnvironmentsRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// Response message for ListWorkspaceBaseEnvironments..
type ListWorkspaceBaseEnvironmentsResponse struct {
	WorkspaceBaseEnvironments []WorkspaceBaseEnvironment `json:"workspace_base_environments"`
	NextPageToken             *string                    `json:"next_page_token"`
}

// This resource represents a long-running operation that is the result of a
// network API call..
type Operation struct {
	Name     *string                                     `json:"name"`
	Metadata *json.RawMessage                            `json:"metadata"`
	Done     *bool                                       `json:"done"`
	Error    *DatabricksServiceExceptionWithDetailsProto `json:"error"`
	Response *json.RawMessage                            `json:"response"`
}

// Request message for RefreshWorkspaceBaseEnvironments..
type RefreshWorkspaceBaseEnvironmentRequest struct {
	Name *string `json:"name"`
}

// Request message for UpdateDefaultWorkspaceBaseEnvironment..
type UpdateDefaultWorkspaceBaseEnvironmentRequest struct {
	DefaultWorkspaceBaseEnvironment *DefaultWorkspaceBaseEnvironment `json:"default_workspace_base_environment"`
	UpdateMask                      *string                          `json:"update_mask"`
}

// Request message for UpdateWorkspaceBaseEnvironment..
type UpdateWorkspaceBaseEnvironmentRequest struct {
	Name                     *string                   `json:"name"`
	WorkspaceBaseEnvironment *WorkspaceBaseEnvironment `json:"workspace_base_environment"`
}

// A WorkspaceBaseEnvironment defines a workspace-level environment
// configuration consisting of an environment version and a list of
// dependencies..
type WorkspaceBaseEnvironment struct {
	Name                    *string                               `json:"name"`
	DisplayName             *string                               `json:"display_name"`
	Filepath                *string                               `json:"filepath"`
	CreatorUserId           *string                               `json:"creator_user_id"`
	CreateTime              *time.Time                            `json:"create_time"`
	LastUpdatedUserId       *string                               `json:"last_updated_user_id"`
	UpdateTime              *time.Time                            `json:"update_time"`
	Status                  *WorkspaceBaseEnvironmentCache_Status `json:"status"`
	Message                 *string                               `json:"message"`
	IsDefault               *bool                                 `json:"is_default"`
	BaseEnvironmentType     *BaseEnvironmentType                  `json:"base_environment_type"`
	BaseEnvironmentProvider *WorkspaceBaseEnvironmentProvider     `json:"base_environment_provider"`
}

// Materialized environment information for a WorkspaceBaseEnvironment..
type WorkspaceBaseEnvironmentCache struct {
}

// Metadata for the WorkspaceBaseEnvironment long-running operations. This
// message tracks the progress of the workspace base environment long-running
// process..
type WorkspaceBaseEnvironmentOperationMetadata struct {
}
