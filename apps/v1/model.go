package v1

import (
	"encoding/json"
	"time"
)

type ComputeSize string

const (
	ComputeSize_Medium ComputeSize = "MEDIUM"
	ComputeSize_Large  ComputeSize = "LARGE"
	ComputeSize_Liquid ComputeSize = "LIQUID"
)

// String representation for [fmt.Print].
func (f *ComputeSize) String() string {
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

type AppDeployment_Mode string

const (
	AppDeployment_Mode_ModeUnspecified AppDeployment_Mode = "MODE_UNSPECIFIED"
	AppDeployment_Mode_Snapshot        AppDeployment_Mode = "SNAPSHOT"
	AppDeployment_Mode_AutoSync        AppDeployment_Mode = "AUTO_SYNC"
)

// String representation for [fmt.Print].
func (f *AppDeployment_Mode) String() string {
	return string(*f)
}

type AppDeployment_State string

const (
	AppDeployment_State_StateUnspecified AppDeployment_State = "STATE_UNSPECIFIED"
	AppDeployment_State_Succeeded        AppDeployment_State = "SUCCEEDED"
	AppDeployment_State_Failed           AppDeployment_State = "FAILED"
	AppDeployment_State_InProgress       AppDeployment_State = "IN_PROGRESS"
	AppDeployment_State_Cancelled        AppDeployment_State = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *AppDeployment_State) String() string {
	return string(*f)
}

type AppManifest_AppResourceAppSpec_AppPermission string

const (
	AppManifest_AppResourceAppSpec_AppPermission_CanUse AppManifest_AppResourceAppSpec_AppPermission = "CAN_USE"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceAppSpec_AppPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceExperimentSpec_ExperimentPermission string

const (
	AppManifest_AppResourceExperimentSpec_ExperimentPermission_CanManage AppManifest_AppResourceExperimentSpec_ExperimentPermission = "CAN_MANAGE"
	AppManifest_AppResourceExperimentSpec_ExperimentPermission_CanEdit   AppManifest_AppResourceExperimentSpec_ExperimentPermission = "CAN_EDIT"
	AppManifest_AppResourceExperimentSpec_ExperimentPermission_CanRead   AppManifest_AppResourceExperimentSpec_ExperimentPermission = "CAN_READ"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceExperimentSpec_ExperimentPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceJobSpec_JobPermission string

const (
	AppManifest_AppResourceJobSpec_JobPermission_CanManage    AppManifest_AppResourceJobSpec_JobPermission = "CAN_MANAGE"
	AppManifest_AppResourceJobSpec_JobPermission_IsOwner      AppManifest_AppResourceJobSpec_JobPermission = "IS_OWNER"
	AppManifest_AppResourceJobSpec_JobPermission_CanManageRun AppManifest_AppResourceJobSpec_JobPermission = "CAN_MANAGE_RUN"
	AppManifest_AppResourceJobSpec_JobPermission_CanView      AppManifest_AppResourceJobSpec_JobPermission = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceJobSpec_JobPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourcePostgresSpec_PostgresPermission string

const (
	AppManifest_AppResourcePostgresSpec_PostgresPermission_CanConnectAndCreate AppManifest_AppResourcePostgresSpec_PostgresPermission = "CAN_CONNECT_AND_CREATE"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourcePostgresSpec_PostgresPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceSecretSpec_SecretPermission string

const (
	AppManifest_AppResourceSecretSpec_SecretPermission_Read   AppManifest_AppResourceSecretSpec_SecretPermission = "READ"
	AppManifest_AppResourceSecretSpec_SecretPermission_Write  AppManifest_AppResourceSecretSpec_SecretPermission = "WRITE"
	AppManifest_AppResourceSecretSpec_SecretPermission_Manage AppManifest_AppResourceSecretSpec_SecretPermission = "MANAGE"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceSecretSpec_SecretPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission string

const (
	AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission_CanManage AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission = "CAN_MANAGE"
	AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission_CanQuery  AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission = "CAN_QUERY"
	AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission_CanView   AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission string

const (
	AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission_CanManage AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission = "CAN_MANAGE"
	AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission_CanUse    AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission = "CAN_USE"
	AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission_IsOwner   AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission = "IS_OWNER"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission string

const (
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_ReadVolume    AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "READ_VOLUME"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_WriteVolume   AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "WRITE_VOLUME"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_Manage        AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "MANAGE"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_Select        AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "SELECT"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_Execute       AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "EXECUTE"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_UseConnection AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "USE_CONNECTION"
	AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission_Modify        AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission = "MODIFY"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission) String() string {
	return string(*f)
}

type AppManifest_AppResourceUcSecurableSpec_UcSecurableType string

const (
	AppManifest_AppResourceUcSecurableSpec_UcSecurableType_Volume     AppManifest_AppResourceUcSecurableSpec_UcSecurableType = "VOLUME"
	AppManifest_AppResourceUcSecurableSpec_UcSecurableType_Table      AppManifest_AppResourceUcSecurableSpec_UcSecurableType = "TABLE"
	AppManifest_AppResourceUcSecurableSpec_UcSecurableType_Function   AppManifest_AppResourceUcSecurableSpec_UcSecurableType = "FUNCTION"
	AppManifest_AppResourceUcSecurableSpec_UcSecurableType_Connection AppManifest_AppResourceUcSecurableSpec_UcSecurableType = "CONNECTION"
)

// String representation for [fmt.Print].
func (f *AppManifest_AppResourceUcSecurableSpec_UcSecurableType) String() string {
	return string(*f)
}

type AppResourceApp_AppPermission string

const (
	AppResourceApp_AppPermission_CanUse AppResourceApp_AppPermission = "CAN_USE"
)

// String representation for [fmt.Print].
func (f *AppResourceApp_AppPermission) String() string {
	return string(*f)
}

type AppResourceDatabase_DatabasePermission string

const (
	AppResourceDatabase_DatabasePermission_CanConnectAndCreate AppResourceDatabase_DatabasePermission = "CAN_CONNECT_AND_CREATE"
)

// String representation for [fmt.Print].
func (f *AppResourceDatabase_DatabasePermission) String() string {
	return string(*f)
}

type AppResourceExperiment_ExperimentPermission string

const (
	AppResourceExperiment_ExperimentPermission_CanManage AppResourceExperiment_ExperimentPermission = "CAN_MANAGE"
	AppResourceExperiment_ExperimentPermission_CanEdit   AppResourceExperiment_ExperimentPermission = "CAN_EDIT"
	AppResourceExperiment_ExperimentPermission_CanRead   AppResourceExperiment_ExperimentPermission = "CAN_READ"
)

// String representation for [fmt.Print].
func (f *AppResourceExperiment_ExperimentPermission) String() string {
	return string(*f)
}

type AppResourceGenieSpace_GenieSpacePermission string

const (
	AppResourceGenieSpace_GenieSpacePermission_CanManage AppResourceGenieSpace_GenieSpacePermission = "CAN_MANAGE"
	AppResourceGenieSpace_GenieSpacePermission_CanEdit   AppResourceGenieSpace_GenieSpacePermission = "CAN_EDIT"
	AppResourceGenieSpace_GenieSpacePermission_CanRun    AppResourceGenieSpace_GenieSpacePermission = "CAN_RUN"
	AppResourceGenieSpace_GenieSpacePermission_CanView   AppResourceGenieSpace_GenieSpacePermission = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *AppResourceGenieSpace_GenieSpacePermission) String() string {
	return string(*f)
}

type AppResourceJob_JobPermission string

const (
	AppResourceJob_JobPermission_CanManage    AppResourceJob_JobPermission = "CAN_MANAGE"
	AppResourceJob_JobPermission_IsOwner      AppResourceJob_JobPermission = "IS_OWNER"
	AppResourceJob_JobPermission_CanManageRun AppResourceJob_JobPermission = "CAN_MANAGE_RUN"
	AppResourceJob_JobPermission_CanView      AppResourceJob_JobPermission = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *AppResourceJob_JobPermission) String() string {
	return string(*f)
}

type AppResourcePostgres_PostgresPermission string

const (
	AppResourcePostgres_PostgresPermission_CanConnectAndCreate AppResourcePostgres_PostgresPermission = "CAN_CONNECT_AND_CREATE"
)

// String representation for [fmt.Print].
func (f *AppResourcePostgres_PostgresPermission) String() string {
	return string(*f)
}

type AppResourceSecret_SecretPermission string

const (
	AppResourceSecret_SecretPermission_Read   AppResourceSecret_SecretPermission = "READ"
	AppResourceSecret_SecretPermission_Write  AppResourceSecret_SecretPermission = "WRITE"
	AppResourceSecret_SecretPermission_Manage AppResourceSecret_SecretPermission = "MANAGE"
)

// String representation for [fmt.Print].
func (f *AppResourceSecret_SecretPermission) String() string {
	return string(*f)
}

type AppResourceServingEndpoint_ServingEndpointPermission string

const (
	AppResourceServingEndpoint_ServingEndpointPermission_CanManage AppResourceServingEndpoint_ServingEndpointPermission = "CAN_MANAGE"
	AppResourceServingEndpoint_ServingEndpointPermission_CanQuery  AppResourceServingEndpoint_ServingEndpointPermission = "CAN_QUERY"
	AppResourceServingEndpoint_ServingEndpointPermission_CanView   AppResourceServingEndpoint_ServingEndpointPermission = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *AppResourceServingEndpoint_ServingEndpointPermission) String() string {
	return string(*f)
}

type AppResourceSqlWarehouse_SqlWarehousePermission string

const (
	AppResourceSqlWarehouse_SqlWarehousePermission_CanManage AppResourceSqlWarehouse_SqlWarehousePermission = "CAN_MANAGE"
	AppResourceSqlWarehouse_SqlWarehousePermission_CanUse    AppResourceSqlWarehouse_SqlWarehousePermission = "CAN_USE"
	AppResourceSqlWarehouse_SqlWarehousePermission_IsOwner   AppResourceSqlWarehouse_SqlWarehousePermission = "IS_OWNER"
)

// String representation for [fmt.Print].
func (f *AppResourceSqlWarehouse_SqlWarehousePermission) String() string {
	return string(*f)
}

type AppResourceUcSecurable_UcSecurablePermission string

const (
	AppResourceUcSecurable_UcSecurablePermission_ReadVolume    AppResourceUcSecurable_UcSecurablePermission = "READ_VOLUME"
	AppResourceUcSecurable_UcSecurablePermission_WriteVolume   AppResourceUcSecurable_UcSecurablePermission = "WRITE_VOLUME"
	AppResourceUcSecurable_UcSecurablePermission_Select        AppResourceUcSecurable_UcSecurablePermission = "SELECT"
	AppResourceUcSecurable_UcSecurablePermission_Execute       AppResourceUcSecurable_UcSecurablePermission = "EXECUTE"
	AppResourceUcSecurable_UcSecurablePermission_UseConnection AppResourceUcSecurable_UcSecurablePermission = "USE_CONNECTION"
	AppResourceUcSecurable_UcSecurablePermission_Modify        AppResourceUcSecurable_UcSecurablePermission = "MODIFY"
)

// String representation for [fmt.Print].
func (f *AppResourceUcSecurable_UcSecurablePermission) String() string {
	return string(*f)
}

type AppResourceUcSecurable_UcSecurableType string

const (
	AppResourceUcSecurable_UcSecurableType_Volume     AppResourceUcSecurable_UcSecurableType = "VOLUME"
	AppResourceUcSecurable_UcSecurableType_Table      AppResourceUcSecurable_UcSecurableType = "TABLE"
	AppResourceUcSecurable_UcSecurableType_Function   AppResourceUcSecurable_UcSecurableType = "FUNCTION"
	AppResourceUcSecurable_UcSecurableType_Connection AppResourceUcSecurable_UcSecurableType = "CONNECTION"
)

// String representation for [fmt.Print].
func (f *AppResourceUcSecurable_UcSecurableType) String() string {
	return string(*f)
}

type AppUpdate_UpdateStatus_UpdateState string

const (
	AppUpdate_UpdateStatus_UpdateState_UpdateStateUnspecified AppUpdate_UpdateStatus_UpdateState = "UPDATE_STATE_UNSPECIFIED"
	AppUpdate_UpdateStatus_UpdateState_NotUpdated             AppUpdate_UpdateStatus_UpdateState = "NOT_UPDATED"
	AppUpdate_UpdateStatus_UpdateState_InProgress             AppUpdate_UpdateStatus_UpdateState = "IN_PROGRESS"
	AppUpdate_UpdateStatus_UpdateState_Succeeded              AppUpdate_UpdateStatus_UpdateState = "SUCCEEDED"
	AppUpdate_UpdateStatus_UpdateState_Failed                 AppUpdate_UpdateStatus_UpdateState = "FAILED"
)

// String representation for [fmt.Print].
func (f *AppUpdate_UpdateStatus_UpdateState) String() string {
	return string(*f)
}

type ApplicationStatus_ApplicationState string

const (
	ApplicationStatus_ApplicationState_ApplicationStateUnspecified ApplicationStatus_ApplicationState = "APPLICATION_STATE_UNSPECIFIED"
	ApplicationStatus_ApplicationState_Deploying                   ApplicationStatus_ApplicationState = "DEPLOYING"
	ApplicationStatus_ApplicationState_Running                     ApplicationStatus_ApplicationState = "RUNNING"
	ApplicationStatus_ApplicationState_Crashed                     ApplicationStatus_ApplicationState = "CRASHED"
	ApplicationStatus_ApplicationState_Unavailable                 ApplicationStatus_ApplicationState = "UNAVAILABLE"
)

// String representation for [fmt.Print].
func (f *ApplicationStatus_ApplicationState) String() string {
	return string(*f)
}

type ComputeStatus_ComputeState string

const (
	ComputeStatus_ComputeState_ComputeStateUnspecified ComputeStatus_ComputeState = "COMPUTE_STATE_UNSPECIFIED"
	ComputeStatus_ComputeState_Error                   ComputeStatus_ComputeState = "ERROR"
	ComputeStatus_ComputeState_Deleting                ComputeStatus_ComputeState = "DELETING"
	ComputeStatus_ComputeState_Starting                ComputeStatus_ComputeState = "STARTING"
	ComputeStatus_ComputeState_Stopping                ComputeStatus_ComputeState = "STOPPING"
	ComputeStatus_ComputeState_Updating                ComputeStatus_ComputeState = "UPDATING"
	ComputeStatus_ComputeState_Stopped                 ComputeStatus_ComputeState = "STOPPED"
	ComputeStatus_ComputeState_Active                  ComputeStatus_ComputeState = "ACTIVE"
)

// String representation for [fmt.Print].
func (f *ComputeStatus_ComputeState) String() string {
	return string(*f)
}

type SpaceStatus_SpaceState string

const (
	SpaceStatus_SpaceState_SpaceStateUnspecified SpaceStatus_SpaceState = "SPACE_STATE_UNSPECIFIED"
	SpaceStatus_SpaceState_SpaceCreating         SpaceStatus_SpaceState = "SPACE_CREATING"
	SpaceStatus_SpaceState_SpaceActive           SpaceStatus_SpaceState = "SPACE_ACTIVE"
	SpaceStatus_SpaceState_SpaceError            SpaceStatus_SpaceState = "SPACE_ERROR"
	SpaceStatus_SpaceState_SpaceDeleting         SpaceStatus_SpaceState = "SPACE_DELETING"
	SpaceStatus_SpaceState_SpaceDeleted          SpaceStatus_SpaceState = "SPACE_DELETED"
	SpaceStatus_SpaceState_SpaceUpdating         SpaceStatus_SpaceState = "SPACE_UPDATING"
)

// String representation for [fmt.Print].
func (f *SpaceStatus_SpaceState) String() string {
	return string(*f)
}

type App struct {
	Name                        *string                      `json:"name"`
	Description                 *string                      `json:"description"`
	ComputeStatus               *ComputeStatus               `json:"compute_status"`
	AppStatus                   *ApplicationStatus           `json:"app_status"`
	Url                         *string                      `json:"url"`
	ActiveDeployment            *AppDeployment               `json:"active_deployment"`
	CreateTime                  *time.Time                   `json:"create_time"`
	Creator                     *string                      `json:"creator"`
	UpdateTime                  *time.Time                   `json:"update_time"`
	Updater                     *string                      `json:"updater"`
	PendingDeployment           *AppDeployment               `json:"pending_deployment"`
	Resources                   []AppResource                `json:"resources"`
	ServicePrincipalId          *int64                       `json:"service_principal_id"`
	ServicePrincipalName        *string                      `json:"service_principal_name"`
	DefaultSourceCodePath       *string                      `json:"default_source_code_path"`
	DefaultGitSource            *GitSource                   `json:"default_git_source"`
	BudgetPolicyId              *string                      `json:"budget_policy_id"`
	EffectiveBudgetPolicyId     *string                      `json:"effective_budget_policy_id"`
	ServicePrincipalClientId    *string                      `json:"service_principal_client_id"`
	UserApiScopes               []string                     `json:"user_api_scopes"`
	Id                          *string                      `json:"id"`
	EffectiveUserApiScopes      []string                     `json:"effective_user_api_scopes"`
	Oauth2AppIntegrationId      *string                      `json:"oauth2_app_integration_id"`
	Oauth2AppClientId           *string                      `json:"oauth2_app_client_id"`
	ComputeSize                 *ComputeSize                 `json:"compute_size"`
	UsagePolicyId               *string                      `json:"usage_policy_id"`
	EffectiveUsagePolicyId      *string                      `json:"effective_usage_policy_id"`
	GitRepository               *GitRepository               `json:"git_repository"`
	TelemetryExportDestinations []TelemetryExportDestination `json:"telemetry_export_destinations"`
	ThumbnailUrl                *string                      `json:"thumbnail_url"`
	Space                       *string                      `json:"space"`
	LastDeploymentId            *string                      `json:"last_deployment_id"`
	SpaceId                     *string                      `json:"space_id"`
}

type AppDeployment struct {
	DeploymentId        *string                 `json:"deployment_id"`
	SourceCodePath      *string                 `json:"source_code_path"`
	GitSource           *GitSource              `json:"git_source"`
	Mode                *AppDeployment_Mode     `json:"mode"`
	DeploymentArtifacts *AppDeploymentArtifacts `json:"deployment_artifacts"`
	Status              *AppDeploymentStatus    `json:"status"`
	CreateTime          *time.Time              `json:"create_time"`
	Creator             *string                 `json:"creator"`
	UpdateTime          *time.Time              `json:"update_time"`
	Command             []string                `json:"command"`
	EnvVars             []EnvVar                `json:"env_vars"`
}

type AppDeploymentArtifacts struct {
	SourceCodePath *string `json:"source_code_path"`
}

type AppDeploymentStatus struct {
	State   *AppDeployment_State `json:"state"`
	Message *string              `json:"message"`
}

// App manifest definition.
type AppManifest struct {
	Version       *int                          `json:"version"`
	Name          *string                       `json:"name"`
	Description   *string                       `json:"description"`
	ResourceSpecs []AppManifest_AppResourceSpec `json:"resource_specs"`
}

type AppManifest_AppResourceAppSpec struct {
	Name       *string                                       `json:"name"`
	Permission *AppManifest_AppResourceAppSpec_AppPermission `json:"permission"`
}

type AppManifest_AppResourceExperimentSpec struct {
	Permission *AppManifest_AppResourceExperimentSpec_ExperimentPermission `json:"permission"`
}

type AppManifest_AppResourceJobSpec struct {
	Permission *AppManifest_AppResourceJobSpec_JobPermission `json:"permission"`
}

type AppManifest_AppResourcePostgresSpec struct {
	Branch     *string                                                 `json:"branch"`
	Database   *string                                                 `json:"database"`
	Permission *AppManifest_AppResourcePostgresSpec_PostgresPermission `json:"permission"`
}

type AppManifest_AppResourceSecretSpec struct {
	Permission *AppManifest_AppResourceSecretSpec_SecretPermission `json:"permission"`
}

type AppManifest_AppResourceServingEndpointSpec struct {
	Permission *AppManifest_AppResourceServingEndpointSpec_ServingEndpointPermission `json:"permission"`
}

// AppResource related fields are copied from app.proto but excludes resource
// identifiers (e.g. name, id, key, scope, etc.).
type AppManifest_AppResourceSpec struct {
	Name                *string                                     `json:"name"`
	Description         *string                                     `json:"description"`
	SecretSpec          *AppManifest_AppResourceSecretSpec          `json:"secret_spec"`
	SqlWarehouseSpec    *AppManifest_AppResourceSqlWarehouseSpec    `json:"sql_warehouse_spec"`
	ServingEndpointSpec *AppManifest_AppResourceServingEndpointSpec `json:"serving_endpoint_spec"`
	JobSpec             *AppManifest_AppResourceJobSpec             `json:"job_spec"`
	UcSecurableSpec     *AppManifest_AppResourceUcSecurableSpec     `json:"uc_securable_spec"`
	ExperimentSpec      *AppManifest_AppResourceExperimentSpec      `json:"experiment_spec"`
	AppSpec             *AppManifest_AppResourceAppSpec             `json:"app_spec"`
	PostgresSpec        *AppManifest_AppResourcePostgresSpec        `json:"postgres_spec"`
}

type AppManifest_AppResourceSqlWarehouseSpec struct {
	Permission *AppManifest_AppResourceSqlWarehouseSpec_SqlWarehousePermission `json:"permission"`
}

type AppManifest_AppResourceUcSecurableSpec struct {
	SecurableType *AppManifest_AppResourceUcSecurableSpec_UcSecurableType       `json:"securable_type"`
	Permission    *AppManifest_AppResourceUcSecurableSpec_UcSecurablePermission `json:"permission"`
}

type AppResource struct {
	Name            *string                     `json:"name"`
	Description     *string                     `json:"description"`
	Secret          *AppResourceSecret          `json:"secret"`
	SqlWarehouse    *AppResourceSqlWarehouse    `json:"sql_warehouse"`
	ServingEndpoint *AppResourceServingEndpoint `json:"serving_endpoint"`
	Job             *AppResourceJob             `json:"job"`
	UcSecurable     *AppResourceUcSecurable     `json:"uc_securable"`
	Database        *AppResourceDatabase        `json:"database"`
	GenieSpace      *AppResourceGenieSpace      `json:"genie_space"`
	Experiment      *AppResourceExperiment      `json:"experiment"`
	App             *AppResourceApp             `json:"app"`
	Postgres        *AppResourcePostgres        `json:"postgres"`
}

type AppResourceApp struct {
	Name       *string                       `json:"name"`
	Permission *AppResourceApp_AppPermission `json:"permission"`
}

type AppResourceDatabase struct {
	InstanceName *string                                 `json:"instance_name"`
	DatabaseName *string                                 `json:"database_name"`
	Permission   *AppResourceDatabase_DatabasePermission `json:"permission"`
}

type AppResourceExperiment struct {
	ExperimentId *string                                     `json:"experiment_id"`
	Permission   *AppResourceExperiment_ExperimentPermission `json:"permission"`
}

type AppResourceGenieSpace struct {
	Name       *string                                     `json:"name"`
	SpaceId    *string                                     `json:"space_id"`
	Permission *AppResourceGenieSpace_GenieSpacePermission `json:"permission"`
}

type AppResourceJob struct {
	Id         *string                       `json:"id"`
	Permission *AppResourceJob_JobPermission `json:"permission"`
}

type AppResourcePostgres struct {
	Branch     *string                                 `json:"branch"`
	Database   *string                                 `json:"database"`
	Permission *AppResourcePostgres_PostgresPermission `json:"permission"`
}

type AppResourceSecret struct {
	Scope      *string                             `json:"scope"`
	Key        *string                             `json:"key"`
	Permission *AppResourceSecret_SecretPermission `json:"permission"`
}

type AppResourceServingEndpoint struct {
	Name       *string                                               `json:"name"`
	Permission *AppResourceServingEndpoint_ServingEndpointPermission `json:"permission"`
}

type AppResourceSqlWarehouse struct {
	Id         *string                                         `json:"id"`
	Permission *AppResourceSqlWarehouse_SqlWarehousePermission `json:"permission"`
}

type AppResourceUcSecurable struct {
	SecurableFullName *string                                       `json:"securable_full_name"`
	SecurableType     *AppResourceUcSecurable_UcSecurableType       `json:"securable_type"`
	Permission        *AppResourceUcSecurable_UcSecurablePermission `json:"permission"`
	SecurableKind     *string                                       `json:"securable_kind"`
}

type AppUpdate struct {
	Status         *AppUpdate_UpdateStatus `json:"status"`
	Description    *string                 `json:"description"`
	BudgetPolicyId *string                 `json:"budget_policy_id"`
	Resources      []AppResource           `json:"resources"`
	UserApiScopes  []string                `json:"user_api_scopes"`
	ComputeSize    *ComputeSize            `json:"compute_size"`
	UsagePolicyId  *string                 `json:"usage_policy_id"`
	GitRepository  *GitRepository          `json:"git_repository"`
}

type AppUpdate_UpdateStatus struct {
	State   *AppUpdate_UpdateStatus_UpdateState `json:"state"`
	Message *string                             `json:"message"`
}

type ApplicationStatus struct {
	State            *ApplicationStatus_ApplicationState `json:"state"`
	Message          *string                             `json:"message"`
	RunningInstances *int                                `json:"running_instances"`
}

type AsyncUpdateAppRequest struct {
	App        *App    `json:"app"`
	UpdateMask *string `json:"update_mask"`
	AppName    *string `json:"app_name"`
}

type ComputeStatus struct {
	State           *ComputeStatus_ComputeState `json:"state"`
	Message         *string                     `json:"message"`
	ActiveInstances *int                        `json:"active_instances"`
}

type CreateAppDeploymentRequest struct {
	AppName       *string        `json:"app_name"`
	AppDeployment *AppDeployment `json:"app_deployment"`
}

type CreateAppRequest struct {
	App       *App  `json:"app"`
	NoCompute *bool `json:"no_compute"`
}

type CreateCustomTemplateRequest struct {
	Template *CustomTemplate `json:"template"`
}

type CreateSpaceRequest struct {
	Space *Space `json:"space"`
}

type CustomTemplate struct {
	Name        *string      `json:"name"`
	Description *string      `json:"description"`
	GitRepo     *string      `json:"git_repo"`
	Path        *string      `json:"path"`
	Manifest    *AppManifest `json:"manifest"`
	GitProvider *string      `json:"git_provider"`
	Creator     *string      `json:"creator"`
}

// Databricks Error that is returned by all Databricks APIs..
type DatabricksServiceExceptionWithDetailsProto struct {
	ErrorCode  *ErrorCode        `json:"error_code"`
	Message    *string           `json:"message"`
	StackTrace *string           `json:"stack_trace"`
	Details    []json.RawMessage `json:"details"`
}

type DeleteAppRequest struct {
	Name *string `json:"name"`
}

type DeleteAppThumbnailRequest struct {
	Name *string `json:"name"`
}

type DeleteCustomTemplateRequest struct {
	Name *string `json:"name"`
}

type DeleteSpaceRequest struct {
	Name *string `json:"name"`
}

type EnvVar struct {
	Name      *string `json:"name"`
	Value     *string `json:"value"`
	ValueFrom *string `json:"value_from"`
}

type GetAppDeploymentRequest struct {
	AppName      *string `json:"app_name"`
	DeploymentId *string `json:"deployment_id"`
}

type GetAppRequest struct {
	Name *string `json:"name"`
}

type GetAppUpdateRequest struct {
	AppName *string `json:"app_name"`
}

type GetCustomTemplateRequest struct {
	Name *string `json:"name"`
}

// The request message for `GetOperation` method..
type GetOperationRequest struct {
	Name *string `json:"name"`
}

type GetSpaceRequest struct {
	Name *string `json:"name"`
}

// Git repository configuration specifying the location of the repository..
type GitRepository struct {
	Url      *string `json:"url"`
	Provider *string `json:"provider"`
}

// Complete git source specification including repository location and
// reference..
type GitSource struct {
	GitRepository  *GitRepository `json:"git_repository"`
	Branch         *string        `json:"branch"`
	Tag            *string        `json:"tag"`
	Commit         *string        `json:"commit"`
	SourceCodePath *string        `json:"source_code_path"`
	ResolvedCommit *string        `json:"resolved_commit"`
}

type ListAppDeploymentsRequest struct {
	AppName   *string `json:"app_name"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListAppDeploymentsResponse struct {
	AppDeployments []AppDeployment `json:"app_deployments"`
	NextPageToken  *string         `json:"next_page_token"`
}

// Request to list all apps deployed in the workspace.
type ListAppsRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
	Space     *string `json:"space"`
}

type ListAppsResponse struct {
	Apps          []App   `json:"apps"`
	NextPageToken *string `json:"next_page_token"`
}

type ListCustomTemplatesRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListCustomTemplatesResponse struct {
	Templates     []CustomTemplate `json:"templates"`
	NextPageToken *string          `json:"next_page_token"`
}

type ListSpacesRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListSpacesResponse struct {
	Spaces        []Space `json:"spaces"`
	NextPageToken *string `json:"next_page_token"`
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

type Space struct {
	Name                     *string       `json:"name"`
	Description              *string       `json:"description"`
	Status                   *SpaceStatus  `json:"status"`
	Id                       *string       `json:"id"`
	CreateTime               *time.Time    `json:"create_time"`
	Creator                  *string       `json:"creator"`
	UpdateTime               *time.Time    `json:"update_time"`
	Updater                  *string       `json:"updater"`
	Resources                []AppResource `json:"resources"`
	UserApiScopes            []string      `json:"user_api_scopes"`
	EffectiveUserApiScopes   []string      `json:"effective_user_api_scopes"`
	ServicePrincipalId       *int64        `json:"service_principal_id"`
	ServicePrincipalName     *string       `json:"service_principal_name"`
	ServicePrincipalClientId *string       `json:"service_principal_client_id"`
	UsagePolicyId            *string       `json:"usage_policy_id"`
	EffectiveUsagePolicyId   *string       `json:"effective_usage_policy_id"`
}

type SpaceStatus struct {
	State   *SpaceStatus_SpaceState `json:"state"`
	Message *string                 `json:"message"`
}

type StartAppRequest struct {
	Name *string `json:"name"`
}

type StopAppRequest struct {
	Name *string `json:"name"`
}

// A single telemetry export destination with its configuration and status..
type TelemetryExportDestination struct {
	UnityCatalog *UnityCatalog `json:"unity_catalog"`
}

// Unity Catalog Destinations for OTEL telemetry export..
type UnityCatalog struct {
	LogsTable    *string `json:"logs_table"`
	MetricsTable *string `json:"metrics_table"`
	TracesTable  *string `json:"traces_table"`
}

type UpdateAppRequest struct {
	App *App `json:"app"`
}

type UpdateCustomTemplateRequest struct {
	Template *CustomTemplate `json:"template"`
}

type UpdateSpaceRequest struct {
	Space      *Space  `json:"space"`
	UpdateMask *string `json:"update_mask"`
}

type UploadAppThumbnailRequest struct {
	Name             *string `json:"name"`
	EncodedThumbnail *string `json:"encoded_thumbnail"`
}
