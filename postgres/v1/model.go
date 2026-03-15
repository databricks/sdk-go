package v1

import (
	"encoding/json"
	"time"
)

type EndpointType string

const (
	EndpointType_EndpointTypeUnspecified EndpointType = "ENDPOINT_TYPE_UNSPECIFIED"
	EndpointType_EndpointTypeReadWrite   EndpointType = "ENDPOINT_TYPE_READ_WRITE"
	EndpointType_EndpointTypeReadOnly    EndpointType = "ENDPOINT_TYPE_READ_ONLY"
)

// String representation for [fmt.Print].
func (f *EndpointType) String() string {
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

type ProvisioningPhase string

const (
	ProvisioningPhase_ProvisioningPhaseUnspecified ProvisioningPhase = "PROVISIONING_PHASE_UNSPECIFIED"
	ProvisioningPhase_ProvisioningPhaseMain        ProvisioningPhase = "PROVISIONING_PHASE_MAIN"
	ProvisioningPhase_ProvisioningPhaseIndexScan   ProvisioningPhase = "PROVISIONING_PHASE_INDEX_SCAN"
	ProvisioningPhase_ProvisioningPhaseIndexSort   ProvisioningPhase = "PROVISIONING_PHASE_INDEX_SORT"
)

// String representation for [fmt.Print].
func (f *ProvisioningPhase) String() string {
	return string(*f)
}

type SyncedTableState string

const (
	SyncedTableState_SyncedTableStateUnspecified                SyncedTableState = "SYNCED_TABLE_STATE_UNSPECIFIED"
	SyncedTableState_SyncedTableProvisioning                    SyncedTableState = "SYNCED_TABLE_PROVISIONING"
	SyncedTableState_SyncedTableProvisioningPipelineResources   SyncedTableState = "SYNCED_TABLE_PROVISIONING_PIPELINE_RESOURCES"
	SyncedTableState_SyncedTableProvisioningInitialSnapshot     SyncedTableState = "SYNCED_TABLE_PROVISIONING_INITIAL_SNAPSHOT"
	SyncedTableState_SyncedTableOnline                          SyncedTableState = "SYNCED_TABLE_ONLINE"
	SyncedTableState_SyncedTableOnlineContinuousUpdate          SyncedTableState = "SYNCED_TABLE_ONLINE_CONTINUOUS_UPDATE"
	SyncedTableState_SyncedTableOnlineTriggeredUpdate           SyncedTableState = "SYNCED_TABLE_ONLINE_TRIGGERED_UPDATE"
	SyncedTableState_SyncedTableOnlineNoPendingUpdate           SyncedTableState = "SYNCED_TABLE_ONLINE_NO_PENDING_UPDATE"
	SyncedTableState_SyncedTabledOffline                        SyncedTableState = "SYNCED_TABLED_OFFLINE"
	SyncedTableState_SyncedTableOfflineFailed                   SyncedTableState = "SYNCED_TABLE_OFFLINE_FAILED"
	SyncedTableState_SyncedTableOnlinePipelineFailed            SyncedTableState = "SYNCED_TABLE_ONLINE_PIPELINE_FAILED"
	SyncedTableState_SyncedTableOnlineUpdatingPipelineResources SyncedTableState = "SYNCED_TABLE_ONLINE_UPDATING_PIPELINE_RESOURCES"
)

// String representation for [fmt.Print].
func (f *SyncedTableState) String() string {
	return string(*f)
}

type BranchStatus_State string

const (
	BranchStatus_State_StateUnspecified BranchStatus_State = "STATE_UNSPECIFIED"
	BranchStatus_State_Init             BranchStatus_State = "INIT"
	BranchStatus_State_Importing        BranchStatus_State = "IMPORTING"
	BranchStatus_State_Resetting        BranchStatus_State = "RESETTING"
	BranchStatus_State_Ready            BranchStatus_State = "READY"
	BranchStatus_State_Archived         BranchStatus_State = "ARCHIVED"
)

// String representation for [fmt.Print].
func (f *BranchStatus_State) String() string {
	return string(*f)
}

type ComputeInstance_ComputeState string

const (
	ComputeInstance_ComputeState_ComputeStateUnspecified ComputeInstance_ComputeState = "COMPUTE_STATE_UNSPECIFIED"
	ComputeInstance_ComputeState_Init                    ComputeInstance_ComputeState = "INIT"
	ComputeInstance_ComputeState_Idle                    ComputeInstance_ComputeState = "IDLE"
	ComputeInstance_ComputeState_Active                  ComputeInstance_ComputeState = "ACTIVE"
)

// String representation for [fmt.Print].
func (f *ComputeInstance_ComputeState) String() string {
	return string(*f)
}

type ComputeInstance_ComputeType string

const (
	ComputeInstance_ComputeType_ComputeTypeUnspecified ComputeInstance_ComputeType = "COMPUTE_TYPE_UNSPECIFIED"
	ComputeInstance_ComputeType_ReadWrite              ComputeInstance_ComputeType = "READ_WRITE"
	ComputeInstance_ComputeType_ReadOnly               ComputeInstance_ComputeType = "READ_ONLY"
	ComputeInstance_ComputeType_HotStandby             ComputeInstance_ComputeType = "HOT_STANDBY"
)

// String representation for [fmt.Print].
func (f *ComputeInstance_ComputeType) String() string {
	return string(*f)
}

type EndpointStatus_State string

const (
	EndpointStatus_State_StateUnspecified EndpointStatus_State = "STATE_UNSPECIFIED"
	EndpointStatus_State_Init             EndpointStatus_State = "INIT"
	EndpointStatus_State_Active           EndpointStatus_State = "ACTIVE"
	EndpointStatus_State_Idle             EndpointStatus_State = "IDLE"
	EndpointStatus_State_Degraded         EndpointStatus_State = "DEGRADED"
)

// String representation for [fmt.Print].
func (f *EndpointStatus_State) String() string {
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

type RequestedClaims_PermissionSet string

const (
	RequestedClaims_PermissionSet_PermissionSetUnspecified RequestedClaims_PermissionSet = "PERMISSION_SET_UNSPECIFIED"
	RequestedClaims_PermissionSet_ReadOnly                 RequestedClaims_PermissionSet = "READ_ONLY"
)

// String representation for [fmt.Print].
func (f *RequestedClaims_PermissionSet) String() string {
	return string(*f)
}

type Role_AuthMethod string

const (
	Role_AuthMethod_AuthMethodUnspecified Role_AuthMethod = "AUTH_METHOD_UNSPECIFIED"
	Role_AuthMethod_NoLogin               Role_AuthMethod = "NO_LOGIN"
	Role_AuthMethod_PgPasswordScramSha256 Role_AuthMethod = "PG_PASSWORD_SCRAM_SHA_256"
	Role_AuthMethod_LakebaseOauthV1       Role_AuthMethod = "LAKEBASE_OAUTH_V1"
)

// String representation for [fmt.Print].
func (f *Role_AuthMethod) String() string {
	return string(*f)
}

type Role_IdentityType string

const (
	Role_IdentityType_IdentityTypeUnspecified Role_IdentityType = "IDENTITY_TYPE_UNSPECIFIED"
	Role_IdentityType_User                    Role_IdentityType = "USER"
	Role_IdentityType_ServicePrincipal        Role_IdentityType = "SERVICE_PRINCIPAL"
	Role_IdentityType_Group                   Role_IdentityType = "GROUP"
)

// String representation for [fmt.Print].
func (f *Role_IdentityType) String() string {
	return string(*f)
}

type Role_MembershipRole string

const (
	Role_MembershipRole_MembershipRoleUnspecified Role_MembershipRole = "MEMBERSHIP_ROLE_UNSPECIFIED"
	Role_MembershipRole_DatabricksSuperuser       Role_MembershipRole = "DATABRICKS_SUPERUSER"
)

// String representation for [fmt.Print].
func (f *Role_MembershipRole) String() string {
	return string(*f)
}

type SyncedTableSpec_SyncedTableSchedulingPolicy string

const (
	SyncedTableSpec_SyncedTableSchedulingPolicy_SyncedTableSchedulingPolicyUnspecified SyncedTableSpec_SyncedTableSchedulingPolicy = "SYNCED_TABLE_SCHEDULING_POLICY_UNSPECIFIED"
	SyncedTableSpec_SyncedTableSchedulingPolicy_Continuous                             SyncedTableSpec_SyncedTableSchedulingPolicy = "CONTINUOUS"
	SyncedTableSpec_SyncedTableSchedulingPolicy_Triggered                              SyncedTableSpec_SyncedTableSchedulingPolicy = "TRIGGERED"
	SyncedTableSpec_SyncedTableSchedulingPolicy_Snapshot                               SyncedTableSpec_SyncedTableSchedulingPolicy = "SNAPSHOT"
)

// String representation for [fmt.Print].
func (f *SyncedTableSpec_SyncedTableSchedulingPolicy) String() string {
	return string(*f)
}

type Branch struct {
	Name       *string       `json:"name"`
	Uid        *string       `json:"uid"`
	Parent     *string       `json:"parent"`
	CreateTime *time.Time    `json:"create_time"`
	UpdateTime *time.Time    `json:"update_time"`
	Spec       *BranchSpec   `json:"spec"`
	Status     *BranchStatus `json:"status"`
}

type BranchSpec struct {
	SourceBranch     *string        `json:"source_branch"`
	SourceBranchLsn  *string        `json:"source_branch_lsn"`
	SourceBranchTime *time.Time     `json:"source_branch_time"`
	IsProtected      *bool          `json:"is_protected"`
	ExpireTime       *time.Time     `json:"expire_time"`
	Ttl              *time.Duration `json:"ttl"`
	NoExpiry         *bool          `json:"no_expiry"`
}

type BranchStatus struct {
	SourceBranch     *string             `json:"source_branch"`
	SourceBranchLsn  *string             `json:"source_branch_lsn"`
	SourceBranchTime *time.Time          `json:"source_branch_time"`
	Default          *bool               `json:"default"`
	IsProtected      *bool               `json:"is_protected"`
	CurrentState     *BranchStatus_State `json:"current_state"`
	PendingState     *BranchStatus_State `json:"pending_state"`
	StateChangeTime  *time.Time          `json:"state_change_time"`
	LogicalSizeBytes *int64              `json:"logical_size_bytes"`
	ExpireTime       *time.Time          `json:"expire_time"`
	BranchId         *string             `json:"branch_id"`
}

type Catalog struct {
	Name                      *string `json:"name"`
	Database                  *string `json:"database"`
	Uid                       *string `json:"uid"`
	CreateDatabaseIfNotExists *bool   `json:"create_database_if_not_exists"`
	Project                   *string `json:"project"`
	Branch                    *string `json:"branch"`
}

type ComputeInstance struct {
	Name              *string                       `json:"name"`
	ComputeInstanceId *string                       `json:"compute_instance_id"`
	CurrentState      *ComputeInstance_ComputeState `json:"current_state"`
	PendingState      *ComputeInstance_ComputeState `json:"pending_state"`
	Role              *ComputeInstance_ComputeType  `json:"role"`
	ComputeHost       *string                       `json:"compute_host"`
	StartTime         *time.Time                    `json:"start_time"`
	SuspendTime       *time.Time                    `json:"suspend_time"`
}

type CreateBranchRequest struct {
	Parent   *string `json:"parent"`
	BranchId *string `json:"branch_id"`
	Branch   *Branch `json:"branch"`
}

type CreateCatalogRequest struct {
	Catalog *Catalog `json:"catalog"`
}

type CreateDatabaseRequest struct {
	Parent     *string   `json:"parent"`
	DatabaseId *string   `json:"database_id"`
	Database   *Database `json:"database"`
}

type CreateEndpointRequest struct {
	Parent     *string   `json:"parent"`
	EndpointId *string   `json:"endpoint_id"`
	Endpoint   *Endpoint `json:"endpoint"`
}

type CreateProjectRequest struct {
	ProjectId *string  `json:"project_id"`
	Project   *Project `json:"project"`
}

type CreateRoleRequest struct {
	Parent *string `json:"parent"`
	RoleId *string `json:"role_id"`
	Role   *Role   `json:"role"`
}

type CreateSyncedTableRequest struct {
	SyncedTable *SyncedTable `json:"synced_table"`
}

type CreateTableRequest struct {
	Table *Table `json:"table"`
}

// Database represents a Postgres database within a Branch..
type Database struct {
	Name       *string                  `json:"name"`
	Parent     *string                  `json:"parent"`
	CreateTime *time.Time               `json:"create_time"`
	UpdateTime *time.Time               `json:"update_time"`
	Spec       *Database_DatabaseSpec   `json:"spec"`
	Status     *Database_DatabaseStatus `json:"status"`
}

type Database_DatabaseSpec struct {
	Role             *string `json:"role"`
	PostgresDatabase *string `json:"postgres_database"`
}

type Database_DatabaseStatus struct {
	Role             *string `json:"role"`
	PostgresDatabase *string `json:"postgres_database"`
}

type DatabaseCredential struct {
	Token      *string    `json:"token"`
	ExpireTime *time.Time `json:"expire_time"`
}

// Databricks Error that is returned by all Databricks APIs..
type DatabricksServiceExceptionWithDetailsProto struct {
	ErrorCode  *ErrorCode        `json:"error_code"`
	Message    *string           `json:"message"`
	StackTrace *string           `json:"stack_trace"`
	Details    []json.RawMessage `json:"details"`
}

type DeleteBranchRequest struct {
	Name *string `json:"name"`
}

type DeleteCatalogRequest struct {
	Name *string `json:"name"`
}

type DeleteDatabaseRequest struct {
	Name *string `json:"name"`
}

type DeleteEndpointRequest struct {
	Name *string `json:"name"`
}

type DeleteProjectRequest struct {
	Name *string `json:"name"`
}

type DeleteRoleRequest struct {
	Name            *string `json:"name"`
	ReassignOwnedTo *string `json:"reassign_owned_to"`
}

type DeleteSyncedTableRequest struct {
	Name *string `json:"name"`
}

type DeleteTableRequest struct {
	Name *string `json:"name"`
}

// Copied from database_table_statuses.proto to decouple SDK packages..
type DeltaTableSyncInfo struct {
	DeltaCommitVersion   *int64     `json:"delta_commit_version"`
	DeltaCommitTimestamp *time.Time `json:"delta_commit_timestamp"`
}

// Request to disable Forward ETL.
type DisableForwardEtlRequest struct {
	Parent        *string `json:"parent"`
	TenantId      *string `json:"tenant_id"`
	TimelineId    *string `json:"timeline_id"`
	PgDatabaseOid *int64  `json:"pg_database_oid"`
	PgSchemaOid   *int64  `json:"pg_schema_oid"`
}

// Response to disable Forward ETL.
type DisableForwardEtlResponse struct {
	Disabled *bool `json:"disabled"`
}

type Endpoint struct {
	Name       *string         `json:"name"`
	Uid        *string         `json:"uid"`
	Parent     *string         `json:"parent"`
	CreateTime *time.Time      `json:"create_time"`
	UpdateTime *time.Time      `json:"update_time"`
	Spec       *EndpointSpec   `json:"spec"`
	Status     *EndpointStatus `json:"status"`
}

type EndpointGroupSpec struct {
	Min                       *int  `json:"min"`
	Max                       *int  `json:"max"`
	EnableReadableSecondaries *bool `json:"enable_readable_secondaries"`
}

type EndpointGroupStatus struct {
	Min                       *int  `json:"min"`
	Max                       *int  `json:"max"`
	EnableReadableSecondaries *bool `json:"enable_readable_secondaries"`
}

// Encapsulates various hostnames (r/w or r/o, pooled or not) for an endpoint..
type EndpointHosts struct {
	Host                *string `json:"host"`
	ReadOnlyHost        *string `json:"read_only_host"`
	ReadWritePooledHost *string `json:"read_write_pooled_host"`
	ReadOnlyPooledHost  *string `json:"read_only_pooled_host"`
}

// A collection of settings for a compute endpoint..
type EndpointSettings struct {
	PgSettings map[string]string `json:"pg_settings"`
}

type EndpointSettings_PgSettingsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type EndpointSpec struct {
	EndpointType           *EndpointType      `json:"endpoint_type"`
	AutoscalingLimitMinCu  *float64           `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu  *float64           `json:"autoscaling_limit_max_cu"`
	Disabled               *bool              `json:"disabled"`
	SuspendTimeoutDuration *time.Duration     `json:"suspend_timeout_duration"`
	NoSuspension           *bool              `json:"no_suspension"`
	Settings               *EndpointSettings  `json:"settings"`
	Group                  *EndpointGroupSpec `json:"group"`
}

type EndpointStatus struct {
	EndpointType           *EndpointType         `json:"endpoint_type"`
	Hosts                  *EndpointHosts        `json:"hosts"`
	LastActiveTime         *time.Time            `json:"last_active_time"`
	AutoscalingLimitMinCu  *float64              `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu  *float64              `json:"autoscaling_limit_max_cu"`
	CurrentState           *EndpointStatus_State `json:"current_state"`
	PendingState           *EndpointStatus_State `json:"pending_state"`
	Disabled               *bool                 `json:"disabled"`
	SuspendTimeoutDuration *time.Duration        `json:"suspend_timeout_duration"`
	Settings               *EndpointSettings     `json:"settings"`
	Group                  *EndpointGroupStatus  `json:"group"`
}

// Forward ETL configuration.
type ForwardEtlConfig struct {
	WorkspaceId      *int64  `json:"workspace_id"`
	TenantId         *string `json:"tenant_id"`
	TimelineId       *string `json:"timeline_id"`
	PgDatabaseOid    *int64  `json:"pg_database_oid"`
	PgSchemaOid      *int64  `json:"pg_schema_oid"`
	UcCatalogId      *string `json:"uc_catalog_id"`
	UcSchemaId       *string `json:"uc_schema_id"`
	Enabled          *bool   `json:"enabled"`
	CreateTimeMillis *int64  `json:"create_time_millis"`
	UpdateTimeMillis *int64  `json:"update_time_millis"`
}

// Database metadata.
type ForwardEtlDatabase struct {
	Name *string `json:"name"`
	Oid  *int64  `json:"oid"`
}

// Forward ETL metadata response.
type ForwardEtlMetadata struct {
	Databases []ForwardEtlDatabase `json:"databases"`
	Schemas   []ForwardEtlSchema   `json:"schemas"`
}

// Schema metadata.
type ForwardEtlSchema struct {
	Name *string `json:"name"`
	Oid  *int64  `json:"oid"`
}

// Forward ETL status response.
type ForwardEtlStatus struct {
	Configurations []ForwardEtlConfig       `json:"configurations"`
	TableMappings  []ForwardEtlTableMapping `json:"table_mappings"`
}

// Per-table replication mapping.
type ForwardEtlTableMapping struct {
	PgTableOid    *int64  `json:"pg_table_oid"`
	UcTableId     *string `json:"uc_table_id"`
	LastSyncedLsn *string `json:"last_synced_lsn"`
	PgTableName   *string `json:"pg_table_name"`
	UcTableName   *string `json:"uc_table_name"`
	Enabled       *bool   `json:"enabled"`
}

type GenerateDatabaseCredentialRequest struct {
	Claims     []RequestedClaims `json:"claims"`
	Endpoint   *string           `json:"endpoint"`
	GroupName  *string           `json:"group_name"`
	Ttl        *time.Duration    `json:"ttl"`
	ExpireTime *time.Time        `json:"expire_time"`
}

type GetBranchRequest struct {
	Name *string `json:"name"`
}

type GetCatalogRequest struct {
	Name *string `json:"name"`
}

type GetComputeInstanceRequest struct {
	Name *string `json:"name"`
}

type GetDatabaseRequest struct {
	Name *string `json:"name"`
}

type GetEndpointRequest struct {
	Name *string `json:"name"`
}

// Request to get Forward ETL metadata.
type GetForwardEtlMetadataRequest struct {
	Parent     *string `json:"parent"`
	TenantId   *string `json:"tenant_id"`
	TimelineId *string `json:"timeline_id"`
}

// Request to get Forward ETL status.
type GetForwardEtlStatusRequest struct {
	Parent     *string `json:"parent"`
	TenantId   *string `json:"tenant_id"`
	TimelineId *string `json:"timeline_id"`
}

// The request message for `GetOperation` method..
type GetOperationRequest struct {
	Name *string `json:"name"`
}

type GetProjectRequest struct {
	Name *string `json:"name"`
}

type GetRoleRequest struct {
	Name *string `json:"name"`
}

type GetSyncedTableRequest struct {
	Name *string `json:"name"`
}

type GetTableRequest struct {
	Name *string `json:"name"`
}

type InitialEndpointSpec struct {
	Group *EndpointGroupSpec `json:"group"`
}

type ListBranchesRequest struct {
	Parent    *string `json:"parent"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListBranchesResponse struct {
	Branches      []Branch `json:"branches"`
	NextPageToken *string  `json:"next_page_token"`
}

type ListComputeInstancesRequest struct {
	Parent    *string `json:"parent"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListComputeInstancesResponse struct {
	ComputeInstances []ComputeInstance `json:"compute_instances"`
	NextPageToken    *string           `json:"next_page_token"`
}

// List Databases..
type ListDatabasesRequest struct {
	Parent    *string `json:"parent"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListDatabasesResponse struct {
	Databases     []Database `json:"databases"`
	NextPageToken *string    `json:"next_page_token"`
}

type ListEndpointsRequest struct {
	Parent    *string `json:"parent"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListEndpointsResponse struct {
	Endpoints     []Endpoint `json:"endpoints"`
	NextPageToken *string    `json:"next_page_token"`
}

type ListProjectsRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListProjectsResponse struct {
	Projects      []Project `json:"projects"`
	NextPageToken *string   `json:"next_page_token"`
}

type ListRolesRequest struct {
	Parent    *string `json:"parent"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListRolesResponse struct {
	Roles         []Role  `json:"roles"`
	NextPageToken *string `json:"next_page_token"`
}

type NewPipelineSpec struct {
	StorageCatalog *string `json:"storage_catalog"`
	StorageSchema  *string `json:"storage_schema"`
	BudgetPolicyId *string `json:"budget_policy_id"`
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

type Project struct {
	Name                *string              `json:"name"`
	Uid                 *string              `json:"uid"`
	CreateTime          *time.Time           `json:"create_time"`
	UpdateTime          *time.Time           `json:"update_time"`
	Spec                *ProjectSpec         `json:"spec"`
	Status              *ProjectStatus       `json:"status"`
	InitialEndpointSpec *InitialEndpointSpec `json:"initial_endpoint_spec"`
}

type ProjectCustomTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// A collection of settings for a compute endpoint..
type ProjectDefaultEndpointSettings struct {
	AutoscalingLimitMinCu  *float64          `json:"autoscaling_limit_min_cu"`
	AutoscalingLimitMaxCu  *float64          `json:"autoscaling_limit_max_cu"`
	SuspendTimeoutDuration *time.Duration    `json:"suspend_timeout_duration"`
	NoSuspension           *bool             `json:"no_suspension"`
	PgSettings             map[string]string `json:"pg_settings"`
}

type ProjectDefaultEndpointSettings_PgSettingsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ProjectSpec struct {
	DisplayName              *string                         `json:"display_name"`
	PgVersion                *int                            `json:"pg_version"`
	HistoryRetentionDuration *time.Duration                  `json:"history_retention_duration"`
	DefaultEndpointSettings  *ProjectDefaultEndpointSettings `json:"default_endpoint_settings"`
	BudgetPolicyId           *string                         `json:"budget_policy_id"`
	CustomTags               []ProjectCustomTag              `json:"custom_tags"`
	WorkspaceKeyEncrypted    *bool                           `json:"workspace_key_encrypted"`
	EnablePgNativeLogin      *bool                           `json:"enable_pg_native_login"`
	DefaultBranch            *string                         `json:"default_branch"`
}

type ProjectStatus struct {
	DisplayName                 *string                         `json:"display_name"`
	PgVersion                   *int                            `json:"pg_version"`
	HistoryRetentionDuration    *time.Duration                  `json:"history_retention_duration"`
	DefaultEndpointSettings     *ProjectDefaultEndpointSettings `json:"default_endpoint_settings"`
	BranchLogicalSizeLimitBytes *int64                          `json:"branch_logical_size_limit_bytes"`
	SyntheticStorageSizeBytes   *int64                          `json:"synthetic_storage_size_bytes"`
	ComputeLastActiveTime       *time.Time                      `json:"compute_last_active_time"`
	BudgetPolicyId              *string                         `json:"budget_policy_id"`
	CustomTags                  []ProjectCustomTag              `json:"custom_tags"`
	Owner                       *string                         `json:"owner"`
	EnablePgNativeLogin         *bool                           `json:"enable_pg_native_login"`
	DefaultBranch               *string                         `json:"default_branch"`
}

// Copied over from managed-catalog/api/messages/common.proto to decouple SDK
// packages. xref go/unified-api-packages-dd.
type ProvisioningInfo struct {
}

type RequestedClaims struct {
	PermissionSet *RequestedClaims_PermissionSet `json:"permission_set"`
	Resources     []RequestedResource            `json:"resources"`
}

type RequestedResource struct {
	UnspecifiedResourceName *string `json:"unspecified_resource_name"`
	TableName               *string `json:"table_name"`
}

// Role represents a Postgres role within a Branch..
type Role struct {
	Name       *string          `json:"name"`
	Parent     *string          `json:"parent"`
	CreateTime *time.Time       `json:"create_time"`
	UpdateTime *time.Time       `json:"update_time"`
	Spec       *Role_RoleSpec   `json:"spec"`
	Status     *Role_RoleStatus `json:"status"`
}

// Attributes that can be granted to a Postgres role. We are only implementing a
// subset for now, see xref:
// https://www.postgresql.org/docs/16/sql-createrole.html The values follow
// Postgres keyword naming e.g. CREATEDB, BYPASSRLS, etc. which is why they
// don't include typical underscores between words..
type Role_Attributes struct {
	Createdb   *bool `json:"createdb"`
	Createrole *bool `json:"createrole"`
	Bypassrls  *bool `json:"bypassrls"`
}

type Role_RoleSpec struct {
	MembershipRoles []Role_MembershipRole `json:"membership_roles"`
	IdentityType    *Role_IdentityType    `json:"identity_type"`
	Attributes      *Role_Attributes      `json:"attributes"`
	AuthMethod      *Role_AuthMethod      `json:"auth_method"`
	PostgresRole    *string               `json:"postgres_role"`
}

type Role_RoleStatus struct {
	MembershipRoles []Role_MembershipRole `json:"membership_roles"`
	IdentityType    *Role_IdentityType    `json:"identity_type"`
	Attributes      *Role_Attributes      `json:"attributes"`
	AuthMethod      *Role_AuthMethod      `json:"auth_method"`
	PostgresRole    *string               `json:"postgres_role"`
}

type SyncedTable struct {
	Name                          *string                 `json:"name"`
	Database                      *string                 `json:"database"`
	Project                       *string                 `json:"project"`
	Branch                        *string                 `json:"branch"`
	Spec                          *SyncedTableSpec        `json:"spec"`
	UnityCatalogProvisioningState *ProvisioningInfo_State `json:"unity_catalog_provisioning_state"`
	DataSynchronizationStatus     *SyncedTableStatus      `json:"data_synchronization_status"`
	TableServingUrl               *string                 `json:"table_serving_url"`
}

// Detailed status of a synced table. Shown if the synced table is in the
// SYNCED_CONTINUOUS_UPDATE or the SYNCED_UPDATING_PIPELINE_RESOURCES state.
// Copied from database_table_statuses.proto to decouple SDK packages..
type SyncedTableContinuousUpdateStatus struct {
	LastProcessedCommitVersion  *int64                       `json:"last_processed_commit_version"`
	Timestamp                   *time.Time                   `json:"timestamp"`
	InitialPipelineSyncProgress *SyncedTablePipelineProgress `json:"initial_pipeline_sync_progress"`
}

// Detailed status of a synced table. Shown if the synced table is in the
// OFFLINE_FAILED or the SYNCED_PIPELINE_FAILED state. Copied from
// database_table_statuses.proto to decouple SDK packages..
type SyncedTableFailedStatus struct {
	LastProcessedCommitVersion *int64     `json:"last_processed_commit_version"`
	Timestamp                  *time.Time `json:"timestamp"`
}

// Progress information of the Synced Table data synchronization pipeline.
// Copied from database_table_statuses.proto to decouple SDK packages..
type SyncedTablePipelineProgress struct {
	LatestVersionCurrentlyProcessing *int64             `json:"latest_version_currently_processing"`
	SyncedRowCount                   *int64             `json:"synced_row_count"`
	TotalRowCount                    *int64             `json:"total_row_count"`
	SyncProgressCompletion           *float64           `json:"sync_progress_completion"`
	EstimatedCompletionTimeSeconds   *float64           `json:"estimated_completion_time_seconds"`
	ProvisioningPhase                *ProvisioningPhase `json:"provisioning_phase"`
}

// Copied from database_table_statuses.proto to decouple SDK packages..
type SyncedTablePosition struct {
	SyncStartTimestamp *time.Time          `json:"sync_start_timestamp"`
	SyncEndTimestamp   *time.Time          `json:"sync_end_timestamp"`
	DeltaTableSyncInfo *DeltaTableSyncInfo `json:"delta_table_sync_info"`
}

// Detailed status of a synced table. Shown if the synced table is in the
// PROVISIONING_PIPELINE_RESOURCES or the PROVISIONING_INITIAL_SNAPSHOT state.
// Copied from database_table_statuses.proto to decouple SDK packages..
type SyncedTableProvisioningStatus struct {
	InitialPipelineSyncProgress *SyncedTablePipelineProgress `json:"initial_pipeline_sync_progress"`
}

type SyncedTableSpec struct {
	SchedulingPolicy               *SyncedTableSpec_SyncedTableSchedulingPolicy `json:"scheduling_policy"`
	SourceTableFullName            *string                                      `json:"source_table_full_name"`
	PrimaryKeyColumns              []string                                     `json:"primary_key_columns"`
	TimeseriesKey                  *string                                      `json:"timeseries_key"`
	ExistingPipelineId             *string                                      `json:"existing_pipeline_id"`
	CreateDatabaseObjectsIfMissing *bool                                        `json:"create_database_objects_if_missing"`
	NewPipelineSpec                *NewPipelineSpec                             `json:"new_pipeline_spec"`
	AcceleratedSync                *bool                                        `json:"accelerated_sync"`
}

type SyncedTableStatus struct {
	Message                *string                            `json:"message"`
	PipelineId             *string                            `json:"pipeline_id"`
	DetailedState          *SyncedTableState                  `json:"detailed_state"`
	ProvisioningStatus     *SyncedTableProvisioningStatus     `json:"provisioning_status"`
	ContinuousUpdateStatus *SyncedTableContinuousUpdateStatus `json:"continuous_update_status"`
	TriggeredUpdateStatus  *SyncedTableTriggeredUpdateStatus  `json:"triggered_update_status"`
	FailedStatus           *SyncedTableFailedStatus           `json:"failed_status"`
	LastSync               *SyncedTablePosition               `json:"last_sync"`
}

// Detailed status of a synced table. Shown if the synced table is in the
// SYNCED_TRIGGERED_UPDATE or the SYNCED_NO_PENDING_UPDATE state. Copied from
// database_table_statuses.proto to decouple SDK packages..
type SyncedTableTriggeredUpdateStatus struct {
	LastProcessedCommitVersion *int64                       `json:"last_processed_commit_version"`
	Timestamp                  *time.Time                   `json:"timestamp"`
	TriggeredUpdateProgress    *SyncedTablePipelineProgress `json:"triggered_update_progress"`
}

// Table represents a non-synced database table in a Lakebase project. Unlike
// SyncedTable, this does not have a data synchronization pipeline..
type Table struct {
	Name            *string `json:"name"`
	Database        *string `json:"database"`
	Project         *string `json:"project"`
	Branch          *string `json:"branch"`
	TableServingUrl *string `json:"table_serving_url"`
}

type UpdateBranchRequest struct {
	Branch     *Branch `json:"branch"`
	UpdateMask *string `json:"update_mask"`
}

type UpdateDatabaseRequest struct {
	Database   *Database `json:"database"`
	UpdateMask *string   `json:"update_mask"`
}

type UpdateEndpointRequest struct {
	Endpoint   *Endpoint `json:"endpoint"`
	UpdateMask *string   `json:"update_mask"`
}

type UpdateProjectRequest struct {
	Project    *Project `json:"project"`
	UpdateMask *string  `json:"update_mask"`
}

type UpdateRoleRequest struct {
	Role       *Role   `json:"role"`
	UpdateMask *string `json:"update_mask"`
}
