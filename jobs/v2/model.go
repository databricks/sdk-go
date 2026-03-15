package v2

type AuthenticationMethod string

const (
	AuthenticationMethod_Oauth AuthenticationMethod = "OAUTH"
	AuthenticationMethod_Pat   AuthenticationMethod = "PAT"
)

// String representation for [fmt.Print].
func (f *AuthenticationMethod) String() string {
	return string(*f)
}

type AwsAvailability string

const (
	AwsAvailability_Spot             AwsAvailability = "SPOT"
	AwsAvailability_OnDemand         AwsAvailability = "ON_DEMAND"
	AwsAvailability_SpotWithFallback AwsAvailability = "SPOT_WITH_FALLBACK"
)

// String representation for [fmt.Print].
func (f *AwsAvailability) String() string {
	return string(*f)
}

type AzureAvailability string

const (
	AzureAvailability_SpotAzure             AzureAvailability = "SPOT_AZURE"
	AzureAvailability_OnDemandAzure         AzureAvailability = "ON_DEMAND_AZURE"
	AzureAvailability_SpotWithFallbackAzure AzureAvailability = "SPOT_WITH_FALLBACK_AZURE"
)

// String representation for [fmt.Print].
func (f *AzureAvailability) String() string {
	return string(*f)
}

type ComputeKind string

const (
	ComputeKind_ComputeKindUnspecified ComputeKind = "COMPUTE_KIND_UNSPECIFIED"
	ComputeKind_ClassicPreview         ComputeKind = "CLASSIC_PREVIEW"
)

// String representation for [fmt.Print].
func (f *ComputeKind) String() string {
	return string(*f)
}

type ConfidentialComputeType string

const (
	ConfidentialComputeType_SevSnp ConfidentialComputeType = "SEV_SNP"
)

// String representation for [fmt.Print].
func (f *ConfidentialComputeType) String() string {
	return string(*f)
}

type DataSecurityMode string

const (
	DataSecurityMode_None                      DataSecurityMode = "NONE"
	DataSecurityMode_SingleUser                DataSecurityMode = "SINGLE_USER"
	DataSecurityMode_UserIsolation             DataSecurityMode = "USER_ISOLATION"
	DataSecurityMode_LegacyTableAcl            DataSecurityMode = "LEGACY_TABLE_ACL"
	DataSecurityMode_LegacyPassthrough         DataSecurityMode = "LEGACY_PASSTHROUGH"
	DataSecurityMode_LegacySingleUser          DataSecurityMode = "LEGACY_SINGLE_USER"
	DataSecurityMode_LegacySingleUserStandard  DataSecurityMode = "LEGACY_SINGLE_USER_STANDARD"
	DataSecurityMode_DataSecurityModeStandard  DataSecurityMode = "DATA_SECURITY_MODE_STANDARD"
	DataSecurityMode_DataSecurityModeDedicated DataSecurityMode = "DATA_SECURITY_MODE_DEDICATED"
	DataSecurityMode_DataSecurityModeAuto      DataSecurityMode = "DATA_SECURITY_MODE_AUTO"
)

// String representation for [fmt.Print].
func (f *DataSecurityMode) String() string {
	return string(*f)
}

type DbtPlatformRunStatus string

const (
	DbtPlatformRunStatus_DbtPlatformRunStatusUnspecified DbtPlatformRunStatus = "DBT_PLATFORM_RUN_STATUS_UNSPECIFIED"
	DbtPlatformRunStatus_Queued                          DbtPlatformRunStatus = "QUEUED"
	DbtPlatformRunStatus_Starting                        DbtPlatformRunStatus = "STARTING"
	DbtPlatformRunStatus_Running                         DbtPlatformRunStatus = "RUNNING"
	DbtPlatformRunStatus_Success                         DbtPlatformRunStatus = "SUCCESS"
	DbtPlatformRunStatus_Error                           DbtPlatformRunStatus = "ERROR"
	DbtPlatformRunStatus_Cancelled                       DbtPlatformRunStatus = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *DbtPlatformRunStatus) String() string {
	return string(*f)
}

type EbsVolumeType string

const (
	EbsVolumeType_GeneralPurposeSsd      EbsVolumeType = "GENERAL_PURPOSE_SSD"
	EbsVolumeType_ThroughputOptimizedHdd EbsVolumeType = "THROUGHPUT_OPTIMIZED_HDD"
)

// String representation for [fmt.Print].
func (f *EbsVolumeType) String() string {
	return string(*f)
}

type Format string

const (
	Format_SingleTask Format = "SINGLE_TASK"
	Format_MultiTask  Format = "MULTI_TASK"
)

// String representation for [fmt.Print].
func (f *Format) String() string {
	return string(*f)
}

type GcpAvailability string

const (
	GcpAvailability_PreemptibleGcp             GcpAvailability = "PREEMPTIBLE_GCP"
	GcpAvailability_OnDemandGcp                GcpAvailability = "ON_DEMAND_GCP"
	GcpAvailability_PreemptibleWithFallbackGcp GcpAvailability = "PREEMPTIBLE_WITH_FALLBACK_GCP"
)

// String representation for [fmt.Print].
func (f *GcpAvailability) String() string {
	return string(*f)
}

type HardwareAcceleratorType string

const (
	HardwareAcceleratorType_Gpu1xA10  HardwareAcceleratorType = "GPU_1xA10"
	HardwareAcceleratorType_Gpu8xH100 HardwareAcceleratorType = "GPU_8xH100"
)

// String representation for [fmt.Print].
func (f *HardwareAcceleratorType) String() string {
	return string(*f)
}

type JobEditMode string

const (
	JobEditMode_UiLocked JobEditMode = "UI_LOCKED"
	JobEditMode_Editable JobEditMode = "EDITABLE"
)

// String representation for [fmt.Print].
func (f *JobEditMode) String() string {
	return string(*f)
}

type JobsHealthMetric string

const (
	JobsHealthMetric_RunDurationSeconds      JobsHealthMetric = "RUN_DURATION_SECONDS"
	JobsHealthMetric_StreamingBacklogBytes   JobsHealthMetric = "STREAMING_BACKLOG_BYTES"
	JobsHealthMetric_StreamingBacklogRecords JobsHealthMetric = "STREAMING_BACKLOG_RECORDS"
	JobsHealthMetric_StreamingBacklogSeconds JobsHealthMetric = "STREAMING_BACKLOG_SECONDS"
	JobsHealthMetric_StreamingBacklogFiles   JobsHealthMetric = "STREAMING_BACKLOG_FILES"
)

// String representation for [fmt.Print].
func (f *JobsHealthMetric) String() string {
	return string(*f)
}

type JobsHealthOperator string

const (
	JobsHealthOperator_GreaterThan JobsHealthOperator = "GREATER_THAN"
)

// String representation for [fmt.Print].
func (f *JobsHealthOperator) String() string {
	return string(*f)
}

type RepairType string

const (
	RepairType_Original RepairType = "ORIGINAL"
	RepairType_Repair   RepairType = "REPAIR"
)

// String representation for [fmt.Print].
func (f *RepairType) String() string {
	return string(*f)
}

type RunType string

const (
	RunType_JobRun      RunType = "JOB_RUN"
	RunType_WorkflowRun RunType = "WORKFLOW_RUN"
	RunType_SubmitRun   RunType = "SUBMIT_RUN"
)

// String representation for [fmt.Print].
func (f *RunType) String() string {
	return string(*f)
}

type RuntimeEngine string

const (
	RuntimeEngine_Null     RuntimeEngine = "NULL"
	RuntimeEngine_Standard RuntimeEngine = "STANDARD"
	RuntimeEngine_Photon   RuntimeEngine = "PHOTON"
)

// String representation for [fmt.Print].
func (f *RuntimeEngine) String() string {
	return string(*f)
}

type SchedulePauseStatus string

const (
	SchedulePauseStatus_Unpaused SchedulePauseStatus = "UNPAUSED"
	SchedulePauseStatus_Paused   SchedulePauseStatus = "PAUSED"
)

// String representation for [fmt.Print].
func (f *SchedulePauseStatus) String() string {
	return string(*f)
}

type Source string

const (
	Source_Workspace Source = "WORKSPACE"
	Source_Git       Source = "GIT"
)

// String representation for [fmt.Print].
func (f *Source) String() string {
	return string(*f)
}

type StorageMode string

const (
	StorageMode_DirectQuery StorageMode = "DIRECT_QUERY"
	StorageMode_Import      StorageMode = "IMPORT"
	StorageMode_Dual        StorageMode = "DUAL"
)

// String representation for [fmt.Print].
func (f *StorageMode) String() string {
	return string(*f)
}

type TaskDependencyType string

const (
	TaskDependencyType_AllSuccess        TaskDependencyType = "ALL_SUCCESS"
	TaskDependencyType_AllDone           TaskDependencyType = "ALL_DONE"
	TaskDependencyType_NoneFailed        TaskDependencyType = "NONE_FAILED"
	TaskDependencyType_AtLeastOneSuccess TaskDependencyType = "AT_LEAST_ONE_SUCCESS"
	TaskDependencyType_AllFailed         TaskDependencyType = "ALL_FAILED"
	TaskDependencyType_AtLeastOneFailed  TaskDependencyType = "AT_LEAST_ONE_FAILED"
)

// String representation for [fmt.Print].
func (f *TaskDependencyType) String() string {
	return string(*f)
}

type TaskRetryMode string

const (
	TaskRetryMode_Never     TaskRetryMode = "NEVER"
	TaskRetryMode_OnFailure TaskRetryMode = "ON_FAILURE"
)

// String representation for [fmt.Print].
func (f *TaskRetryMode) String() string {
	return string(*f)
}

type TriggerType string

const (
	TriggerType_Periodic          TriggerType = "PERIODIC"
	TriggerType_OneTime           TriggerType = "ONE_TIME"
	TriggerType_Retry             TriggerType = "RETRY"
	TriggerType_RunJobTask        TriggerType = "RUN_JOB_TASK"
	TriggerType_FileArrival       TriggerType = "FILE_ARRIVAL"
	TriggerType_Continuous        TriggerType = "CONTINUOUS"
	TriggerType_Table             TriggerType = "TABLE"
	TriggerType_ContinuousRestart TriggerType = "CONTINUOUS_RESTART"
)

// String representation for [fmt.Print].
func (f *TriggerType) String() string {
	return string(*f)
}

type ViewType string

const (
	ViewType_Notebook  ViewType = "NOTEBOOK"
	ViewType_Dashboard ViewType = "DASHBOARD"
)

// String representation for [fmt.Print].
func (f *ViewType) String() string {
	return string(*f)
}

type ViewsToExport string

const (
	ViewsToExport_Code       ViewsToExport = "CODE"
	ViewsToExport_Dashboards ViewsToExport = "DASHBOARDS"
	ViewsToExport_All        ViewsToExport = "ALL"
)

// String representation for [fmt.Print].
func (f *ViewsToExport) String() string {
	return string(*f)
}

type AccessControlRequest_JobPermission string

const (
	AccessControlRequest_JobPermission_CanView      AccessControlRequest_JobPermission = "CAN_VIEW"
	AccessControlRequest_JobPermission_CanManageRun AccessControlRequest_JobPermission = "CAN_MANAGE_RUN"
	AccessControlRequest_JobPermission_IsOwner      AccessControlRequest_JobPermission = "IS_OWNER"
	AccessControlRequest_JobPermission_CanManage    AccessControlRequest_JobPermission = "CAN_MANAGE"
)

// String representation for [fmt.Print].
func (f *AccessControlRequest_JobPermission) String() string {
	return string(*f)
}

type AlertEvaluationState_AlertEvaluationState string

const (
	AlertEvaluationState_AlertEvaluationState_AlertEvaluationStateUnspecified AlertEvaluationState_AlertEvaluationState = "ALERT_EVALUATION_STATE_UNSPECIFIED"
	AlertEvaluationState_AlertEvaluationState_Unknown                         AlertEvaluationState_AlertEvaluationState = "UNKNOWN"
	AlertEvaluationState_AlertEvaluationState_Triggered                       AlertEvaluationState_AlertEvaluationState = "TRIGGERED"
	AlertEvaluationState_AlertEvaluationState_Ok                              AlertEvaluationState_AlertEvaluationState = "OK"
	AlertEvaluationState_AlertEvaluationState_Error                           AlertEvaluationState_AlertEvaluationState = "ERROR"
)

// String representation for [fmt.Print].
func (f *AlertEvaluationState_AlertEvaluationState) String() string {
	return string(*f)
}

type CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState string

const (
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_RunLifeCycleStateUnspecified CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "RUN_LIFE_CYCLE_STATE_UNSPECIFIED"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Pending                      CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "PENDING"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Running                      CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "RUNNING"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Terminating                  CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "TERMINATING"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Terminated                   CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "TERMINATED"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Skipped                      CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "SKIPPED"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_InternalError                CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "INTERNAL_ERROR"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Blocked                      CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "BLOCKED"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_WaitingForRetry              CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "WAITING_FOR_RETRY"
	CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState_Queued                       CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState = "QUEUED"
)

// String representation for [fmt.Print].
func (f *CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState) String() string {
	return string(*f)
}

type CleanRoomTaskRunResultState_CleanRoomTaskRunResultState string

const (
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_RunResultStateUnspecified    CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "RUN_RESULT_STATE_UNSPECIFIED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Success                      CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "SUCCESS"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Failed                       CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "FAILED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Timedout                     CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "TIMEDOUT"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Canceled                     CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "CANCELED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_MaximumConcurrentRunsReached CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "MAXIMUM_CONCURRENT_RUNS_REACHED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_UpstreamCanceled             CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "UPSTREAM_CANCELED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_UpstreamFailed               CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "UPSTREAM_FAILED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Excluded                     CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "EXCLUDED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Evicted                      CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "EVICTED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_SuccessWithFailures          CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "SUCCESS_WITH_FAILURES"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_UpstreamEvicted              CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "UPSTREAM_EVICTED"
	CleanRoomTaskRunResultState_CleanRoomTaskRunResultState_Disabled                     CleanRoomTaskRunResultState_CleanRoomTaskRunResultState = "DISABLED"
)

// String representation for [fmt.Print].
func (f *CleanRoomTaskRunResultState_CleanRoomTaskRunResultState) String() string {
	return string(*f)
}

type ConditionTask_ConditionTaskOperator string

const (
	ConditionTask_ConditionTaskOperator_EqualTo            ConditionTask_ConditionTaskOperator = "EQUAL_TO"
	ConditionTask_ConditionTaskOperator_GreaterThan        ConditionTask_ConditionTaskOperator = "GREATER_THAN"
	ConditionTask_ConditionTaskOperator_GreaterThanOrEqual ConditionTask_ConditionTaskOperator = "GREATER_THAN_OR_EQUAL"
	ConditionTask_ConditionTaskOperator_LessThan           ConditionTask_ConditionTaskOperator = "LESS_THAN"
	ConditionTask_ConditionTaskOperator_LessThanOrEqual    ConditionTask_ConditionTaskOperator = "LESS_THAN_OR_EQUAL"
	ConditionTask_ConditionTaskOperator_NotEqual           ConditionTask_ConditionTaskOperator = "NOT_EQUAL"
)

// String representation for [fmt.Print].
func (f *ConditionTask_ConditionTaskOperator) String() string {
	return string(*f)
}

type JobDeployment_DeploymentKind string

const (
	JobDeployment_DeploymentKind_Bundle        JobDeployment_DeploymentKind = "BUNDLE"
	JobDeployment_DeploymentKind_SystemManaged JobDeployment_DeploymentKind = "SYSTEM_MANAGED"
)

// String representation for [fmt.Print].
func (f *JobDeployment_DeploymentKind) String() string {
	return string(*f)
}

type JobSource_DirtyState string

const (
	JobSource_DirtyState_NotSynced    JobSource_DirtyState = "NOT_SYNCED"
	JobSource_DirtyState_Disconnected JobSource_DirtyState = "DISCONNECTED"
)

// String representation for [fmt.Print].
func (f *JobSource_DirtyState) String() string {
	return string(*f)
}

type ModelTriggerConfiguration_ModelTriggerCondition string

const (
	ModelTriggerConfiguration_ModelTriggerCondition_ConditionUnspecified ModelTriggerConfiguration_ModelTriggerCondition = "CONDITION_UNSPECIFIED"
	ModelTriggerConfiguration_ModelTriggerCondition_ModelCreated         ModelTriggerConfiguration_ModelTriggerCondition = "MODEL_CREATED"
	ModelTriggerConfiguration_ModelTriggerCondition_ModelVersionReady    ModelTriggerConfiguration_ModelTriggerCondition = "MODEL_VERSION_READY"
	ModelTriggerConfiguration_ModelTriggerCondition_ModelAliasSet        ModelTriggerConfiguration_ModelTriggerCondition = "MODEL_ALIAS_SET"
)

// String representation for [fmt.Print].
func (f *ModelTriggerConfiguration_ModelTriggerCondition) String() string {
	return string(*f)
}

type PerformanceTarget_PerformanceTarget string

const (
	PerformanceTarget_PerformanceTarget_PerformanceTargetUnspecified PerformanceTarget_PerformanceTarget = "PERFORMANCE_TARGET_UNSPECIFIED"
	PerformanceTarget_PerformanceTarget_PerformanceOptimized         PerformanceTarget_PerformanceTarget = "PERFORMANCE_OPTIMIZED"
	PerformanceTarget_PerformanceTarget_Standard                     PerformanceTarget_PerformanceTarget = "STANDARD"
)

// String representation for [fmt.Print].
func (f *PerformanceTarget_PerformanceTarget) String() string {
	return string(*f)
}

type PeriodicTriggerConfiguration_TimeUnit string

const (
	PeriodicTriggerConfiguration_TimeUnit_TimeUnitUnspecified PeriodicTriggerConfiguration_TimeUnit = "TIME_UNIT_UNSPECIFIED"
	PeriodicTriggerConfiguration_TimeUnit_Hours               PeriodicTriggerConfiguration_TimeUnit = "HOURS"
	PeriodicTriggerConfiguration_TimeUnit_Days                PeriodicTriggerConfiguration_TimeUnit = "DAYS"
	PeriodicTriggerConfiguration_TimeUnit_Weeks               PeriodicTriggerConfiguration_TimeUnit = "WEEKS"
)

// String representation for [fmt.Print].
func (f *PeriodicTriggerConfiguration_TimeUnit) String() string {
	return string(*f)
}

type QueueDetailsCode_Code string

const (
	QueueDetailsCode_Code_ActiveRunsLimitReached        QueueDetailsCode_Code = "ACTIVE_RUNS_LIMIT_REACHED"
	QueueDetailsCode_Code_MaxConcurrentRunsReached      QueueDetailsCode_Code = "MAX_CONCURRENT_RUNS_REACHED"
	QueueDetailsCode_Code_ActiveRunJobTasksLimitReached QueueDetailsCode_Code = "ACTIVE_RUN_JOB_TASKS_LIMIT_REACHED"
)

// String representation for [fmt.Print].
func (f *QueueDetailsCode_Code) String() string {
	return string(*f)
}

type RunLifeCycleState_RunLifeCycleState string

const (
	RunLifeCycleState_RunLifeCycleState_Pending         RunLifeCycleState_RunLifeCycleState = "PENDING"
	RunLifeCycleState_RunLifeCycleState_Running         RunLifeCycleState_RunLifeCycleState = "RUNNING"
	RunLifeCycleState_RunLifeCycleState_Terminating     RunLifeCycleState_RunLifeCycleState = "TERMINATING"
	RunLifeCycleState_RunLifeCycleState_Terminated      RunLifeCycleState_RunLifeCycleState = "TERMINATED"
	RunLifeCycleState_RunLifeCycleState_Skipped         RunLifeCycleState_RunLifeCycleState = "SKIPPED"
	RunLifeCycleState_RunLifeCycleState_InternalError   RunLifeCycleState_RunLifeCycleState = "INTERNAL_ERROR"
	RunLifeCycleState_RunLifeCycleState_Blocked         RunLifeCycleState_RunLifeCycleState = "BLOCKED"
	RunLifeCycleState_RunLifeCycleState_WaitingForRetry RunLifeCycleState_RunLifeCycleState = "WAITING_FOR_RETRY"
	RunLifeCycleState_RunLifeCycleState_Queued          RunLifeCycleState_RunLifeCycleState = "QUEUED"
)

// String representation for [fmt.Print].
func (f *RunLifeCycleState_RunLifeCycleState) String() string {
	return string(*f)
}

type RunLifecycleStateV2_State string

const (
	RunLifecycleStateV2_State_Blocked     RunLifecycleStateV2_State = "BLOCKED"
	RunLifecycleStateV2_State_Pending     RunLifecycleStateV2_State = "PENDING"
	RunLifecycleStateV2_State_Queued      RunLifecycleStateV2_State = "QUEUED"
	RunLifecycleStateV2_State_Running     RunLifecycleStateV2_State = "RUNNING"
	RunLifecycleStateV2_State_Terminating RunLifecycleStateV2_State = "TERMINATING"
	RunLifecycleStateV2_State_Terminated  RunLifecycleStateV2_State = "TERMINATED"
	RunLifecycleStateV2_State_Waiting     RunLifecycleStateV2_State = "WAITING"
)

// String representation for [fmt.Print].
func (f *RunLifecycleStateV2_State) String() string {
	return string(*f)
}

type RunResultState_RunResultState string

const (
	RunResultState_RunResultState_Success                      RunResultState_RunResultState = "SUCCESS"
	RunResultState_RunResultState_Failed                       RunResultState_RunResultState = "FAILED"
	RunResultState_RunResultState_Timedout                     RunResultState_RunResultState = "TIMEDOUT"
	RunResultState_RunResultState_Canceled                     RunResultState_RunResultState = "CANCELED"
	RunResultState_RunResultState_MaximumConcurrentRunsReached RunResultState_RunResultState = "MAXIMUM_CONCURRENT_RUNS_REACHED"
	RunResultState_RunResultState_UpstreamCanceled             RunResultState_RunResultState = "UPSTREAM_CANCELED"
	RunResultState_RunResultState_UpstreamFailed               RunResultState_RunResultState = "UPSTREAM_FAILED"
	RunResultState_RunResultState_Excluded                     RunResultState_RunResultState = "EXCLUDED"
	RunResultState_RunResultState_SuccessWithFailures          RunResultState_RunResultState = "SUCCESS_WITH_FAILURES"
	RunResultState_RunResultState_Disabled                     RunResultState_RunResultState = "DISABLED"
)

// String representation for [fmt.Print].
func (f *RunResultState_RunResultState) String() string {
	return string(*f)
}

type SqlAlertState_SqlAlertState string

const (
	SqlAlertState_SqlAlertState_Unknown   SqlAlertState_SqlAlertState = "UNKNOWN"
	SqlAlertState_SqlAlertState_Ok        SqlAlertState_SqlAlertState = "OK"
	SqlAlertState_SqlAlertState_Triggered SqlAlertState_SqlAlertState = "TRIGGERED"
)

// String representation for [fmt.Print].
func (f *SqlAlertState_SqlAlertState) String() string {
	return string(*f)
}

type SqlTask_SqlTaskQueryStatus string

const (
	SqlTask_SqlTaskQueryStatus_Pending   SqlTask_SqlTaskQueryStatus = "PENDING"
	SqlTask_SqlTaskQueryStatus_Running   SqlTask_SqlTaskQueryStatus = "RUNNING"
	SqlTask_SqlTaskQueryStatus_Success   SqlTask_SqlTaskQueryStatus = "SUCCESS"
	SqlTask_SqlTaskQueryStatus_Failed    SqlTask_SqlTaskQueryStatus = "FAILED"
	SqlTask_SqlTaskQueryStatus_Cancelled SqlTask_SqlTaskQueryStatus = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *SqlTask_SqlTaskQueryStatus) String() string {
	return string(*f)
}

type TableTriggerConfiguration_Condition string

const (
	TableTriggerConfiguration_Condition_AnyUpdated TableTriggerConfiguration_Condition = "ANY_UPDATED"
	TableTriggerConfiguration_Condition_AllUpdated TableTriggerConfiguration_Condition = "ALL_UPDATED"
)

// String representation for [fmt.Print].
func (f *TableTriggerConfiguration_Condition) String() string {
	return string(*f)
}

type TerminationCode_Code string

const (
	TerminationCode_Code_Success                     TerminationCode_Code = "SUCCESS"
	TerminationCode_Code_Canceled                    TerminationCode_Code = "CANCELED"
	TerminationCode_Code_DriverError                 TerminationCode_Code = "DRIVER_ERROR"
	TerminationCode_Code_ClusterError                TerminationCode_Code = "CLUSTER_ERROR"
	TerminationCode_Code_RepositoryCheckoutFailed    TerminationCode_Code = "REPOSITORY_CHECKOUT_FAILED"
	TerminationCode_Code_InvalidClusterRequest       TerminationCode_Code = "INVALID_CLUSTER_REQUEST"
	TerminationCode_Code_WorkspaceRunLimitExceeded   TerminationCode_Code = "WORKSPACE_RUN_LIMIT_EXCEEDED"
	TerminationCode_Code_FeatureDisabled             TerminationCode_Code = "FEATURE_DISABLED"
	TerminationCode_Code_ClusterRequestLimitExceeded TerminationCode_Code = "CLUSTER_REQUEST_LIMIT_EXCEEDED"
	TerminationCode_Code_StorageAccessError          TerminationCode_Code = "STORAGE_ACCESS_ERROR"
	TerminationCode_Code_RunExecutionError           TerminationCode_Code = "RUN_EXECUTION_ERROR"
	TerminationCode_Code_UnauthorizedError           TerminationCode_Code = "UNAUTHORIZED_ERROR"
	TerminationCode_Code_LibraryInstallationError    TerminationCode_Code = "LIBRARY_INSTALLATION_ERROR"
	TerminationCode_Code_MaxConcurrentRunsExceeded   TerminationCode_Code = "MAX_CONCURRENT_RUNS_EXCEEDED"
	TerminationCode_Code_MaxSparkContextsExceeded    TerminationCode_Code = "MAX_SPARK_CONTEXTS_EXCEEDED"
	TerminationCode_Code_ResourceNotFound            TerminationCode_Code = "RESOURCE_NOT_FOUND"
	TerminationCode_Code_InvalidRunConfiguration     TerminationCode_Code = "INVALID_RUN_CONFIGURATION"
	TerminationCode_Code_InternalError               TerminationCode_Code = "INTERNAL_ERROR"
	TerminationCode_Code_CloudFailure                TerminationCode_Code = "CLOUD_FAILURE"
	TerminationCode_Code_MaxJobQueueSizeExceeded     TerminationCode_Code = "MAX_JOB_QUEUE_SIZE_EXCEEDED"
	TerminationCode_Code_Skipped                     TerminationCode_Code = "SKIPPED"
	TerminationCode_Code_UserCanceled                TerminationCode_Code = "USER_CANCELED"
	TerminationCode_Code_BudgetPolicyLimitExceeded   TerminationCode_Code = "BUDGET_POLICY_LIMIT_EXCEEDED"
	TerminationCode_Code_Disabled                    TerminationCode_Code = "DISABLED"
	TerminationCode_Code_SuccessWithFailures         TerminationCode_Code = "SUCCESS_WITH_FAILURES"
)

// String representation for [fmt.Print].
func (f *TerminationCode_Code) String() string {
	return string(*f)
}

type TerminationType_Type string

const (
	TerminationType_Type_Success       TerminationType_Type = "SUCCESS"
	TerminationType_Type_InternalError TerminationType_Type = "INTERNAL_ERROR"
	TerminationType_Type_ClientError   TerminationType_Type = "CLIENT_ERROR"
	TerminationType_Type_CloudFailure  TerminationType_Type = "CLOUD_FAILURE"
)

// String representation for [fmt.Print].
func (f *TerminationType_Type) String() string {
	return string(*f)
}

type AccessControlRequest struct {
	UserName             *string                             `json:"user_name"`
	GroupName            *string                             `json:"group_name"`
	ServicePrincipalName *string                             `json:"service_principal_name"`
	PermissionLevel      *AccessControlRequest_JobPermission `json:"permission_level"`
}

// A storage location in Adls Gen2.
type Adlsgen2Info struct {
	Destination *string `json:"destination"`
}

// Defines an agentic task configuration for job-based execution. References an
// existing agent with a goal and optional output schema..
type AgenticTask struct {
	Goal             *string           `json:"goal"`
	SupervisorAgent  *SupervisorAgent  `json:"supervisor_agent"`
	TaskOutputSchema map[string]string `json:"task_output_schema"`
}

type AgenticTask_TaskOutputSchemaEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type AlertEvaluationState struct {
}

type AlertTask struct {
	AlertId       *string               `json:"alert_id"`
	WarehouseId   *string               `json:"warehouse_id"`
	WorkspacePath *string               `json:"workspace_path"`
	Subscribers   []AlertTaskSubscriber `json:"subscribers"`
}

type AlertTaskOutput struct {
	AlertState *AlertEvaluationState_AlertEvaluationState `json:"alert_state"`
}

// Represents a subscriber that will receive alert notifications. A subscriber
// can be either a user (via email) or a notification destination (via
// destination_id)..
type AlertTaskSubscriber struct {
	UserName      *string `json:"user_name"`
	DestinationId *string `json:"destination_id"`
}

type AutoScale struct {
	MinWorkers *int `json:"min_workers"`
	MaxWorkers *int `json:"max_workers"`
}

// Attributes set during cluster creation which are related to Amazon Web
// Services..
type AwsAttributes struct {
	FirstOnDemand       *int             `json:"first_on_demand"`
	Availability        *AwsAvailability `json:"availability"`
	ZoneId              *string          `json:"zone_id"`
	InstanceProfileArn  *string          `json:"instance_profile_arn"`
	SpotBidPricePercent *int             `json:"spot_bid_price_percent"`
	EbsVolumeType       *EbsVolumeType   `json:"ebs_volume_type"`
	EbsVolumeCount      *int             `json:"ebs_volume_count"`
	EbsVolumeSize       *int             `json:"ebs_volume_size"`
	EbsVolumeIops       *int             `json:"ebs_volume_iops"`
	EbsVolumeThroughput *int             `json:"ebs_volume_throughput"`
}

// Attributes set during cluster creation which are related to Microsoft Azure..
type AzureAttributes struct {
	LogAnalyticsInfo *LogAnalyticsInfo  `json:"log_analytics_info"`
	FirstOnDemand    *int               `json:"first_on_demand"`
	Availability     *AzureAvailability `json:"availability"`
	SpotBidMaxPrice  *float64           `json:"spot_bid_max_price"`
}

type BaseJob struct {
	JobId                   *int64             `json:"job_id"`
	CreatorUserName         *string            `json:"creator_user_name"`
	RunAsUserName           *string            `json:"run_as_user_name"`
	Settings                *JobSettings       `json:"settings"`
	CreatedTime             *int64             `json:"created_time"`
	TriggerState            *TriggerStateProto `json:"trigger_state"`
	HasMore                 *bool              `json:"has_more"`
	EffectiveBudgetPolicyId *string            `json:"effective_budget_policy_id"`
	EffectiveUsagePolicyId  *string            `json:"effective_usage_policy_id"`
	Path                    *string            `json:"path"`
}

type BaseRun struct {
	JobId                      *int64                               `json:"job_id"`
	RunId                      *int64                               `json:"run_id"`
	CreatorUserName            *string                              `json:"creator_user_name"`
	NumberInJob                *int64                               `json:"number_in_job"`
	OriginalAttemptRunId       *int64                               `json:"original_attempt_run_id"`
	State                      *RunState                            `json:"state"`
	Schedule                   *CronSchedule                        `json:"schedule"`
	ClusterSpec                *ClusterSpec                         `json:"cluster_spec"`
	ClusterInstance            *ClusterInstance                     `json:"cluster_instance"`
	JobParameters              []Run_JobLevelParameters             `json:"job_parameters"`
	OverridingParameters       *RunParameters                       `json:"overriding_parameters"`
	Trigger                    *TriggerType                         `json:"trigger"`
	TriggerInfo                *RunTriggerInfo                      `json:"trigger_info"`
	RunName                    *string                              `json:"run_name"`
	RunPageUrl                 *string                              `json:"run_page_url"`
	RunType                    *RunType                             `json:"run_type"`
	Tasks                      []RunTask                            `json:"tasks"`
	Description                *string                              `json:"description"`
	AttemptNumber              *int                                 `json:"attempt_number"`
	JobClusters                []JobCluster                         `json:"job_clusters"`
	GitSource                  *GitSource                           `json:"git_source"`
	RepairHistory              []Repair                             `json:"repair_history"`
	Status                     *RunStatus                           `json:"status"`
	JobRunId                   *int64                               `json:"job_run_id"`
	HasMore                    *bool                                `json:"has_more"`
	EffectivePerformanceTarget *PerformanceTarget_PerformanceTarget `json:"effective_performance_target"`
	EffectiveUsagePolicyId     *string                              `json:"effective_usage_policy_id"`
	StartTime                  *int64                               `json:"start_time"`
	SetupDuration              *int64                               `json:"setup_duration"`
	ExecutionDuration          *int64                               `json:"execution_duration"`
	CleanupDuration            *int64                               `json:"cleanup_duration"`
	EndTime                    *int64                               `json:"end_time"`
	RunDuration                *int64                               `json:"run_duration"`
	QueueDuration              *int64                               `json:"queue_duration"`
}

type CancelAllRuns struct {
	JobId         *int64 `json:"job_id"`
	AllQueuedRuns *bool  `json:"all_queued_runs"`
}

// All runs were cancelled successfully..
type CancelAllRuns_Response struct {
}

type CancelRun struct {
	RunId *int64 `json:"run_id"`
}

// Run was cancelled successfully..
type CancelRun_Response struct {
}

type CleanRoomTaskRunLifeCycleState struct {
}

type CleanRoomTaskRunResultState struct {
}

// Stores the run state of the clean rooms notebook task..
type CleanRoomTaskRunState struct {
	LifeCycleState *CleanRoomTaskRunLifeCycleState_CleanRoomTaskRunLifeCycleState `json:"life_cycle_state"`
	ResultState    *CleanRoomTaskRunResultState_CleanRoomTaskRunResultState       `json:"result_state"`
}

// Clean Rooms notebook task for V1 Clean Room service (GA). Replaces the
// deprecated CleanRoomNotebookTask (defined above) which was for V0 service..
type CleanRoomsNotebookTask struct {
	CleanRoomName          *string           `json:"clean_room_name"`
	NotebookName           *string           `json:"notebook_name"`
	Etag                   *string           `json:"etag"`
	NotebookBaseParameters map[string]string `json:"notebook_base_parameters"`
}

type CleanRoomsNotebookTask_CleanRoomsNotebookTaskOutput struct {
	CleanRoomJobRunState   *CleanRoomTaskRunState       `json:"clean_room_job_run_state"`
	NotebookOutput         *NotebookTask_NotebookOutput `json:"notebook_output"`
	OutputSchemaInfo       *OutputSchemaInfo            `json:"output_schema_info"`
	SharedOutputSchemaInfo *OutputSchemaInfo            `json:"shared_output_schema_info"`
}

// Name-based parameters for jobs running notebook tasks..
type CleanRoomsNotebookTask_NotebookBaseParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ClusterInstance struct {
	ClusterId      *string `json:"cluster_id"`
	SparkContextId *string `json:"spark_context_id"`
}

// Cluster log delivery config.
type ClusterLogConf struct {
	Dbfs    *DbfsStorageInfo    `json:"dbfs"`
	S3      *S3StorageInfo      `json:"s3"`
	Volumes *VolumesStorageInfo `json:"volumes"`
}

type ClusterSpec struct {
	ExistingClusterId *string                 `json:"existing_cluster_id"`
	NewCluster        *ClusterSpec_NewCluster `json:"new_cluster"`
	JobClusterKey     *string                 `json:"job_cluster_key"`
	Libraries         []Library               `json:"libraries"`
}

type ClusterSpec_NewCluster struct {
	ApplyPolicyDefaultValues   *bool                `json:"apply_policy_default_values"`
	ClusterName                *string              `json:"cluster_name"`
	SparkVersion               *string              `json:"spark_version"`
	SparkConf                  map[string]string    `json:"spark_conf"`
	AwsAttributes              *AwsAttributes       `json:"aws_attributes"`
	AzureAttributes            *AzureAttributes     `json:"azure_attributes"`
	GcpAttributes              *GcpAttributes       `json:"gcp_attributes"`
	NodeTypeId                 *string              `json:"node_type_id"`
	DriverNodeTypeId           *string              `json:"driver_node_type_id"`
	WorkerNodeTypeFlexibility  *NodeTypeFlexibility `json:"worker_node_type_flexibility"`
	DriverNodeTypeFlexibility  *NodeTypeFlexibility `json:"driver_node_type_flexibility"`
	SshPublicKeys              []string             `json:"ssh_public_keys"`
	CustomTags                 map[string]string    `json:"custom_tags"`
	ClusterLogConf             *ClusterLogConf      `json:"cluster_log_conf"`
	SparkEnvVars               map[string]string    `json:"spark_env_vars"`
	AutoterminationMinutes     *int                 `json:"autotermination_minutes"`
	EnableElasticDisk          *bool                `json:"enable_elastic_disk"`
	InitScripts                []InitScriptInfo     `json:"init_scripts"`
	DockerImage                *DockerImage         `json:"docker_image"`
	InstancePoolId             *string              `json:"instance_pool_id"`
	SingleUserName             *string              `json:"single_user_name"`
	PolicyId                   *string              `json:"policy_id"`
	EnableLocalDiskEncryption  *bool                `json:"enable_local_disk_encryption"`
	DriverInstancePoolId       *string              `json:"driver_instance_pool_id"`
	WorkloadType               *WorkloadType        `json:"workload_type"`
	DataSecurityMode           *DataSecurityMode    `json:"data_security_mode"`
	RuntimeEngine              *RuntimeEngine       `json:"runtime_engine"`
	Kind                       *ComputeKind         `json:"kind"`
	UseMlRuntime               *bool                `json:"use_ml_runtime"`
	IsSingleNode               *bool                `json:"is_single_node"`
	RemoteDiskThroughput       *int                 `json:"remote_disk_throughput"`
	TotalInitialRemoteDiskSize *int                 `json:"total_initial_remote_disk_size"`
	NumWorkers                 *int                 `json:"num_workers"`
	Autoscale                  *AutoScale           `json:"autoscale"`
}

type ClusterSpec_NewCluster_CustomTagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark configuration key-value pairs.
type ClusterSpec_NewCluster_SparkConfEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Spark environment variable key-value pairs.
type ClusterSpec_NewCluster_SparkEnvVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type Compute struct {
	HardwareAccelerator *HardwareAcceleratorType `json:"hardware_accelerator"`
}

type ComputeConfig struct {
	NumGpus       *int    `json:"num_gpus"`
	GpuNodePoolId *string `json:"gpu_node_pool_id"`
	GpuType       *string `json:"gpu_type"`
}

type ConditionTask struct {
	Op      *ConditionTask_ConditionTaskOperator `json:"op"`
	Left    *string                              `json:"left"`
	Right   *string                              `json:"right"`
	Outcome *string                              `json:"outcome"`
}

type ContinuousSettings struct {
	PauseStatus   *SchedulePauseStatus `json:"pause_status"`
	TaskRetryMode *TaskRetryMode       `json:"task_retry_mode"`
}

type CreateJob struct {
	AccessControlList       []AccessControlRequest               `json:"access_control_list"`
	Name                    *string                              `json:"name"`
	Description             *string                              `json:"description"`
	EmailNotifications      *JobEmailNotifications               `json:"email_notifications"`
	WebhookNotifications    *WebhookNotifications                `json:"webhook_notifications"`
	NotificationSettings    *NotificationSettings                `json:"notification_settings"`
	TimeoutSeconds          *int                                 `json:"timeout_seconds"`
	Health                  *JobsHealthRules                     `json:"health"`
	Schedule                *CronSchedule                        `json:"schedule"`
	Trigger                 *TriggerSettings                     `json:"trigger"`
	Continuous              *ContinuousSettings                  `json:"continuous"`
	MaxConcurrentRuns       *int                                 `json:"max_concurrent_runs"`
	Tasks                   []TaskSettings                       `json:"tasks"`
	JobClusters             []JobCluster                         `json:"job_clusters"`
	GitSource               *GitSource                           `json:"git_source"`
	Tags                    map[string]string                    `json:"tags"`
	Format                  *Format                              `json:"format"`
	Queue                   *QueueSettings                       `json:"queue"`
	Parameters              []JobLevelParameter                  `json:"parameters"`
	RunAs                   *JobRunAs                            `json:"run_as"`
	EditMode                *JobEditMode                         `json:"edit_mode"`
	Deployment              *JobDeployment                       `json:"deployment"`
	Environments            []JobEnvironment                     `json:"environments"`
	BudgetPolicyId          *string                              `json:"budget_policy_id"`
	UsagePolicyId           *string                              `json:"usage_policy_id"`
	PerformanceTarget       *PerformanceTarget_PerformanceTarget `json:"performance_target"`
	ParentPath              *string                              `json:"parent_path"`
	MaxRetries              *int                                 `json:"max_retries"`
	MinRetryIntervalMillis  *int                                 `json:"min_retry_interval_millis"`
	RetryOnTimeout          *bool                                `json:"retry_on_timeout"`
	DisableAutoOptimization *bool                                `json:"disable_auto_optimization"`
}

// Job was created successfully.
type CreateJob_Response struct {
	JobId *int64 `json:"job_id"`
}

// <Databricks> proto compiler is too old and does not support map. This is wire
// compatible with map<string,string>. See
// https://developers.google.com/protocol-buffers/docs/proto#backwards_compatibility..
type CreateJob_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CronSchedule struct {
	QuartzCronExpression *string                    `json:"quartz_cron_expression"`
	TimezoneId           *string                    `json:"timezone_id"`
	PauseStatus          *SchedulePauseStatus       `json:"pause_status"`
	SqlCondition         *SqlConditionConfiguration `json:"sql_condition"`
}

type DashboardPageSnapshot struct {
	PageDisplayName    *string             `json:"page_display_name"`
	WidgetErrorDetails []WidgetErrorDetail `json:"widget_error_details"`
}

// Configures the Lakeview Dashboard job task type..
type DashboardTask struct {
	Subscription *Subscription     `json:"subscription"`
	WarehouseId  *string           `json:"warehouse_id"`
	DashboardId  *string           `json:"dashboard_id"`
	Filters      map[string]string `json:"filters"`
}

type DashboardTask_FiltersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type DashboardTaskOutput struct {
	PageSnapshots []DashboardPageSnapshot `json:"page_snapshots"`
}

// A storage location in DBFS.
type DbfsStorageInfo struct {
	Destination *string `json:"destination"`
}

// Format of response retrieved from dbt Cloud, for inclusion in output
// Deprecated in favor of DbtPlatformJobRunStep.
type DbtCloudJobRunStep struct {
	Index  *int                  `json:"index"`
	Name   *string               `json:"name"`
	Status *DbtPlatformRunStatus `json:"status"`
	Logs   *string               `json:"logs"`
}

// Deprecated in favor of DbtPlatformTask.
type DbtCloudTask struct {
	DbtCloudJobId          *int64  `json:"dbt_cloud_job_id"`
	ConnectionResourceName *string `json:"connection_resource_name"`
}

// Deprecated in favor of DbtPlatformTaskOutput.
type DbtCloudTaskOutput struct {
	DbtCloudJobRunId     *int64               `json:"dbt_cloud_job_run_id"`
	DbtCloudJobRunUrl    *string              `json:"dbt_cloud_job_run_url"`
	DbtCloudJobRunOutput []DbtCloudJobRunStep `json:"dbt_cloud_job_run_output"`
}

// Format of response retrieved from dbt platform, for inclusion in output.
type DbtPlatformJobRunStep struct {
	Index         *int                  `json:"index"`
	Name          *string               `json:"name"`
	Status        *DbtPlatformRunStatus `json:"status"`
	Logs          *string               `json:"logs"`
	NameTruncated *bool                 `json:"name_truncated"`
	LogsTruncated *bool                 `json:"logs_truncated"`
}

type DbtPlatformTask struct {
	DbtPlatformJobId       *string `json:"dbt_platform_job_id"`
	ConnectionResourceName *string `json:"connection_resource_name"`
}

type DbtPlatformTaskOutput struct {
	DbtPlatformJobRunId     *string                 `json:"dbt_platform_job_run_id"`
	DbtPlatformJobRunUrl    *string                 `json:"dbt_platform_job_run_url"`
	DbtPlatformJobRunOutput []DbtPlatformJobRunStep `json:"dbt_platform_job_run_output"`
	StepsTruncated          *bool                   `json:"steps_truncated"`
}

type DbtTask struct {
	ProjectDirectory  *string  `json:"project_directory"`
	Commands          []string `json:"commands"`
	Schema            *string  `json:"schema"`
	WarehouseId       *string  `json:"warehouse_id"`
	ProfilesDirectory *string  `json:"profiles_directory"`
	Catalog           *string  `json:"catalog"`
	Source            *Source  `json:"source"`
}

type DbtTask_DbtTaskOutput struct {
	ArtifactsLink    *string           `json:"artifacts_link"`
	ArtifactsHeaders map[string]string `json:"artifacts_headers"`
}

type DbtTask_DbtTaskOutput_ArtifactsHeadersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type DeleteJob struct {
	JobId *int64 `json:"job_id"`
}

// Job was deleted successfully..
type DeleteJob_Response struct {
}

type DeleteRun struct {
	RunId *int64 `json:"run_id"`
}

// Run was deleted successfully..
type DeleteRun_Response struct {
}

type DockerBasicAuth struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type DockerImage struct {
	Url       *string          `json:"url"`
	BasicAuth *DockerBasicAuth `json:"basic_auth"`
}

// The environment entity used to preserve serverless environment side panel,
// jobs' environment for non-notebook task, and DLT's environment for classic
// and serverless pipelines. In this minimal environment spec, only pip
// dependencies are supported..
type Environment struct {
	Client             *string  `json:"client"`
	Dependencies       []string `json:"dependencies"`
	BaseEnvironment    *string  `json:"base_environment"`
	EnvironmentVersion *string  `json:"environment_version"`
	JavaDependencies   []string `json:"java_dependencies"`
}

// Retrieves the export of a job run task..
type ExportRun struct {
	RunId         *int64         `json:"run_id"`
	ViewsToExport *ViewsToExport `json:"views_to_export"`
}

// Run was exported successfully..
type ExportRun_Response struct {
	Views []ViewItem `json:"views"`
}

type FileArrivalTriggerConfiguration struct {
	Url                           *string `json:"url"`
	MinTimeBetweenTriggersSeconds *int    `json:"min_time_between_triggers_seconds"`
	WaitAfterLastChangeSeconds    *int    `json:"wait_after_last_change_seconds"`
}

type FileArrivalTriggerState struct {
	UsingFileEvents *bool `json:"using_file_events"`
}

type ForEachTask struct {
	Inputs      *string       `json:"inputs"`
	Concurrency *int          `json:"concurrency"`
	Task        *TaskSettings `json:"task"`
}

// Attributes set during cluster creation which are related to GCP..
type GcpAttributes struct {
	UsePreemptibleExecutors *bool                    `json:"use_preemptible_executors"`
	GoogleServiceAccount    *string                  `json:"google_service_account"`
	BootDiskSize            *int                     `json:"boot_disk_size"`
	Availability            *GcpAvailability         `json:"availability"`
	ZoneId                  *string                  `json:"zone_id"`
	LocalSsdCount           *int                     `json:"local_ssd_count"`
	FirstOnDemand           *int                     `json:"first_on_demand"`
	ConfidentialComputeType *ConfidentialComputeType `json:"confidential_compute_type"`
}

// A storage location in Google Cloud Platform's GCS.
type GcsStorageInfo struct {
	Destination *string `json:"destination"`
}

type GenAiComputeTask struct {
	DlRuntimeImage         *string        `json:"dl_runtime_image"`
	Compute                *ComputeConfig `json:"compute"`
	Command                *string        `json:"command"`
	Source                 *Source        `json:"source"`
	TrainingScriptPath     *string        `json:"training_script_path"`
	YamlParametersFilePath *string        `json:"yaml_parameters_file_path"`
	YamlParameters         *string        `json:"yaml_parameters"`
	MlflowExperimentName   *string        `json:"mlflow_experiment_name"`
	DockerImageUrl         *string        `json:"docker_image_url"`
}

// Retrieves information about a single job..
type GetJob struct {
	JobId     *int64  `json:"job_id"`
	PageToken *string `json:"page_token"`
}

// Job was retrieved successfully..
type GetJob_Response struct {
	NextPageToken           *string            `json:"next_page_token"`
	JobId                   *int64             `json:"job_id"`
	CreatorUserName         *string            `json:"creator_user_name"`
	RunAsUserName           *string            `json:"run_as_user_name"`
	Settings                *JobSettings       `json:"settings"`
	CreatedTime             *int64             `json:"created_time"`
	TriggerState            *TriggerStateProto `json:"trigger_state"`
	HasMore                 *bool              `json:"has_more"`
	EffectiveBudgetPolicyId *string            `json:"effective_budget_policy_id"`
	EffectiveUsagePolicyId  *string            `json:"effective_usage_policy_id"`
	Path                    *string            `json:"path"`
}

type GetRun struct {
	RunId                 *int64  `json:"run_id"`
	IncludeHistory        *bool   `json:"include_history"`
	IncludeResolvedValues *bool   `json:"include_resolved_values"`
	PageToken             *string `json:"page_token"`
}

// Run was retrieved successfully.
type GetRun_Response struct {
	NextPageToken              *string                              `json:"next_page_token"`
	JobId                      *int64                               `json:"job_id"`
	RunId                      *int64                               `json:"run_id"`
	CreatorUserName            *string                              `json:"creator_user_name"`
	NumberInJob                *int64                               `json:"number_in_job"`
	OriginalAttemptRunId       *int64                               `json:"original_attempt_run_id"`
	State                      *RunState                            `json:"state"`
	Schedule                   *CronSchedule                        `json:"schedule"`
	ClusterSpec                *ClusterSpec                         `json:"cluster_spec"`
	ClusterInstance            *ClusterInstance                     `json:"cluster_instance"`
	JobParameters              []Run_JobLevelParameters             `json:"job_parameters"`
	OverridingParameters       *RunParameters                       `json:"overriding_parameters"`
	Trigger                    *TriggerType                         `json:"trigger"`
	TriggerInfo                *RunTriggerInfo                      `json:"trigger_info"`
	RunName                    *string                              `json:"run_name"`
	RunPageUrl                 *string                              `json:"run_page_url"`
	RunType                    *RunType                             `json:"run_type"`
	Tasks                      []RunTask                            `json:"tasks"`
	Description                *string                              `json:"description"`
	AttemptNumber              *int                                 `json:"attempt_number"`
	JobClusters                []JobCluster                         `json:"job_clusters"`
	GitSource                  *GitSource                           `json:"git_source"`
	RepairHistory              []Repair                             `json:"repair_history"`
	Status                     *RunStatus                           `json:"status"`
	JobRunId                   *int64                               `json:"job_run_id"`
	HasMore                    *bool                                `json:"has_more"`
	EffectivePerformanceTarget *PerformanceTarget_PerformanceTarget `json:"effective_performance_target"`
	EffectiveUsagePolicyId     *string                              `json:"effective_usage_policy_id"`
	StartTime                  *int64                               `json:"start_time"`
	SetupDuration              *int64                               `json:"setup_duration"`
	ExecutionDuration          *int64                               `json:"execution_duration"`
	CleanupDuration            *int64                               `json:"cleanup_duration"`
	EndTime                    *int64                               `json:"end_time"`
	RunDuration                *int64                               `json:"run_duration"`
	QueueDuration              *int64                               `json:"queue_duration"`
}

// Retrieves both the output and the metadata of a run..
type GetRunOutput struct {
	RunId *int64 `json:"run_id"`
}

// Run output was retrieved successfully..
type GetRunOutput_Response struct {
	Metadata                 *Run                                                 `json:"metadata"`
	Error                    *string                                              `json:"error"`
	Info                     *string                                              `json:"info"`
	NotebookOutput           *NotebookTask_NotebookOutput                         `json:"notebook_output"`
	SqlOutput                *SqlTask_SqlOutput                                   `json:"sql_output"`
	DbtOutput                *DbtTask_DbtTaskOutput                               `json:"dbt_output"`
	RunJobOutput             *RunJobTask_RunJobTaskOutput                         `json:"run_job_output"`
	CleanRoomsNotebookOutput *CleanRoomsNotebookTask_CleanRoomsNotebookTaskOutput `json:"clean_rooms_notebook_output"`
	DashboardOutput          *DashboardTaskOutput                                 `json:"dashboard_output"`
	DbtCloudOutput           *DbtCloudTaskOutput                                  `json:"dbt_cloud_output"`
	DbtPlatformOutput        *DbtPlatformTaskOutput                               `json:"dbt_platform_output"`
	AlertOutput              *AlertTaskOutput                                     `json:"alert_output"`
	Logs                     *string                                              `json:"logs"`
	LogsTruncated            *bool                                                `json:"logs_truncated"`
	ErrorTrace               *string                                              `json:"error_trace"`
}

// Read-only state of the remote repository at the time the job was run. This
// field is only included on job runs..
type GitMetadataSnapshot struct {
	UsedCommit *string `json:"used_commit"`
}

// An optional specification for a remote Git repository containing the source
// code used by tasks. Version-controlled source code is supported by notebook,
// dbt, Python script, and SQL File tasks.
//
// If `git_source` is set, these tasks retrieve the file from the remote
// repository by default. However, this behavior can be overridden by setting
// `source` to `WORKSPACE` on the task.
//
// Note: dbt and SQL File tasks support only version-controlled sources. If dbt
// or SQL File tasks are used, `git_source` must be defined on the job..
type GitSource struct {
	GitUrl         *string              `json:"git_url"`
	GitProvider    *string              `json:"git_provider"`
	GitBranch      *string              `json:"git_branch"`
	GitTag         *string              `json:"git_tag"`
	GitCommit      *string              `json:"git_commit"`
	GitSnapshot    *GitMetadataSnapshot `json:"git_snapshot"`
	JobSource      *JobSource           `json:"job_source"`
	SparseCheckout *SparseCheckout      `json:"sparse_checkout"`
}

// Config for an individual init script Next ID: 11.
type InitScriptInfo struct {
	Dbfs      *DbfsStorageInfo      `json:"dbfs"`
	S3        *S3StorageInfo        `json:"s3"`
	File      *LocalFileInfo        `json:"file"`
	Gcs       *GcsStorageInfo       `json:"gcs"`
	Abfss     *Adlsgen2Info         `json:"abfss"`
	Workspace *WorkspaceStorageInfo `json:"workspace"`
	Volumes   *VolumesStorageInfo   `json:"volumes"`
}

type JobCluster struct {
	JobClusterKey *string                 `json:"job_cluster_key"`
	NewCluster    *ClusterSpec_NewCluster `json:"new_cluster"`
}

type JobDeployment struct {
	Kind             *JobDeployment_DeploymentKind `json:"kind"`
	MetadataFilePath *string                       `json:"metadata_file_path"`
}

type JobEmailNotifications struct {
	OnStart                            []string `json:"on_start"`
	OnSuccess                          []string `json:"on_success"`
	OnFailure                          []string `json:"on_failure"`
	OnDurationWarningThresholdExceeded []string `json:"on_duration_warning_threshold_exceeded"`
	OnStreamingBacklogExceeded         []string `json:"on_streaming_backlog_exceeded"`
	NoAlertForSkippedRuns              *bool    `json:"no_alert_for_skipped_runs"`
}

type JobEnvironment struct {
	EnvironmentKey *string      `json:"environment_key"`
	Spec           *Environment `json:"spec"`
}

type JobLevelParameter struct {
	Name    *string `json:"name"`
	Default *string `json:"default"`
}

// Write-only setting. Specifies the user or service principal that the job runs
// as. If not specified, the job runs as the user who created the job.
//
// Either `user_name` or `service_principal_name` should be specified. If not,
// an error is thrown..
type JobRunAs struct {
	UserName             *string `json:"user_name"`
	ServicePrincipalName *string `json:"service_principal_name"`
	GroupName            *string `json:"group_name"`
}

type JobSettings struct {
	Name                    *string                              `json:"name"`
	Description             *string                              `json:"description"`
	EmailNotifications      *JobEmailNotifications               `json:"email_notifications"`
	WebhookNotifications    *WebhookNotifications                `json:"webhook_notifications"`
	NotificationSettings    *NotificationSettings                `json:"notification_settings"`
	TimeoutSeconds          *int                                 `json:"timeout_seconds"`
	Health                  *JobsHealthRules                     `json:"health"`
	Schedule                *CronSchedule                        `json:"schedule"`
	Trigger                 *TriggerSettings                     `json:"trigger"`
	Continuous              *ContinuousSettings                  `json:"continuous"`
	MaxConcurrentRuns       *int                                 `json:"max_concurrent_runs"`
	Tasks                   []TaskSettings                       `json:"tasks"`
	JobClusters             []JobCluster                         `json:"job_clusters"`
	GitSource               *GitSource                           `json:"git_source"`
	Tags                    map[string]string                    `json:"tags"`
	Format                  *Format                              `json:"format"`
	Queue                   *QueueSettings                       `json:"queue"`
	Parameters              []JobLevelParameter                  `json:"parameters"`
	RunAs                   *JobRunAs                            `json:"run_as"`
	EditMode                *JobEditMode                         `json:"edit_mode"`
	Deployment              *JobDeployment                       `json:"deployment"`
	Environments            []JobEnvironment                     `json:"environments"`
	BudgetPolicyId          *string                              `json:"budget_policy_id"`
	UsagePolicyId           *string                              `json:"usage_policy_id"`
	PerformanceTarget       *PerformanceTarget_PerformanceTarget `json:"performance_target"`
	ParentPath              *string                              `json:"parent_path"`
	MaxRetries              *int                                 `json:"max_retries"`
	MinRetryIntervalMillis  *int                                 `json:"min_retry_interval_millis"`
	RetryOnTimeout          *bool                                `json:"retry_on_timeout"`
	DisableAutoOptimization *bool                                `json:"disable_auto_optimization"`
}

// <Databricks> proto compiler is too old and does not support map. This is wire
// compatible with map<string,string>. See
// https://developers.google.com/protocol-buffers/docs/proto#backwards_compatibility..
type JobSettings_TagsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// The source of the job specification in the remote repository when the job is
// source controlled..
type JobSource struct {
	JobConfigPath       *string               `json:"job_config_path"`
	ImportFromGitBranch *string               `json:"import_from_git_branch"`
	DirtyState          *JobSource_DirtyState `json:"dirty_state"`
}

type JobsHealthRule struct {
	Metric *JobsHealthMetric   `json:"metric"`
	Op     *JobsHealthOperator `json:"op"`
	Value  *int64              `json:"value"`
}

// An optional set of health rules that can be defined for this job..
type JobsHealthRules struct {
	Rules []JobsHealthRule `json:"rules"`
}

type Library struct {
	Jar          *string            `json:"jar"`
	Egg          *string            `json:"egg"`
	Pypi         *PythonPyPiLibrary `json:"pypi"`
	Maven        *MavenLibrary      `json:"maven"`
	Cran         *RCranLibrary      `json:"cran"`
	Whl          *string            `json:"whl"`
	Requirements *string            `json:"requirements"`
}

// Lists all jobs..
type ListJobs struct {
	Offset      *int    `json:"offset"`
	Limit       *int    `json:"limit"`
	ExpandTasks *bool   `json:"expand_tasks"`
	Name        *string `json:"name"`
	PageToken   *string `json:"page_token"`
}

// List of jobs was retrieved successfully..
type ListJobs_Response struct {
	Jobs          []BaseJob `json:"jobs"`
	HasMore       *bool     `json:"has_more"`
	NextPageToken *string   `json:"next_page_token"`
	PrevPageToken *string   `json:"prev_page_token"`
}

// Lists runs from most recently started to least..
type ListRuns struct {
	JobId         *int64   `json:"job_id"`
	ActiveOnly    *bool    `json:"active_only"`
	CompletedOnly *bool    `json:"completed_only"`
	Offset        *int     `json:"offset"`
	Limit         *int     `json:"limit"`
	RunType       *RunType `json:"run_type"`
	ExpandTasks   *bool    `json:"expand_tasks"`
	StartTimeFrom *int64   `json:"start_time_from"`
	StartTimeTo   *int64   `json:"start_time_to"`
	PageToken     *string  `json:"page_token"`
}

// List of runs was retrieved successfully..
type ListRuns_Response struct {
	Runs          []BaseRun `json:"runs"`
	HasMore       *bool     `json:"has_more"`
	NextPageToken *string   `json:"next_page_token"`
	PrevPageToken *string   `json:"prev_page_token"`
}

type LocalFileInfo struct {
	Destination *string `json:"destination"`
}

type LogAnalyticsInfo struct {
	LogAnalyticsWorkspaceId *string `json:"log_analytics_workspace_id"`
	LogAnalyticsPrimaryKey  *string `json:"log_analytics_primary_key"`
}

type MavenLibrary struct {
	Coordinates *string  `json:"coordinates"`
	Repo        *string  `json:"repo"`
	Exclusions  []string `json:"exclusions"`
}

type ModelTriggerConfiguration struct {
	SecurableName                 *string                                          `json:"securable_name"`
	Aliases                       []string                                         `json:"aliases"`
	Condition                     *ModelTriggerConfiguration_ModelTriggerCondition `json:"condition"`
	MinTimeBetweenTriggersSeconds *int                                             `json:"min_time_between_triggers_seconds"`
	WaitAfterLastChangeSeconds    *int                                             `json:"wait_after_last_change_seconds"`
}

// Configuration for flexible node types, allowing fallback to alternate node
// types during cluster launch and upscale..
type NodeTypeFlexibility struct {
	AlternateNodeTypeIds []string `json:"alternate_node_type_ids"`
}

type NotebookTask struct {
	NotebookPath   *string           `json:"notebook_path"`
	BaseParameters map[string]string `json:"base_parameters"`
	Source         *Source           `json:"source"`
	WarehouseId    *string           `json:"warehouse_id"`
}

// Name-based parameters for jobs running notebook tasks..
type NotebookTask_BaseParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type NotebookTask_NotebookOutput struct {
	Result    *string `json:"result"`
	Truncated *bool   `json:"truncated"`
}

type NotificationSettings struct {
	NoAlertForSkippedRuns  *bool `json:"no_alert_for_skipped_runs"`
	NoAlertForCanceledRuns *bool `json:"no_alert_for_canceled_runs"`
	AlertOnLastAttempt     *bool `json:"alert_on_last_attempt"`
}

// Stores the catalog name, schema name, and the output schema expiration time
// for the clean room run..
type OutputSchemaInfo struct {
	CatalogName    *string `json:"catalog_name"`
	SchemaName     *string `json:"schema_name"`
	ExpirationTime *int64  `json:"expiration_time"`
}

type PerformanceTarget struct {
}

type PeriodicTriggerConfiguration struct {
	Interval *int                                   `json:"interval"`
	Unit     *PeriodicTriggerConfiguration_TimeUnit `json:"unit"`
}

type PipelineParameters struct {
	FullRefresh *bool `json:"full_refresh"`
}

type PipelineTask struct {
	PipelineId  *string `json:"pipeline_id"`
	FullRefresh *bool   `json:"full_refresh"`
}

type PowerBiModel struct {
	WorkspaceName        *string               `json:"workspace_name"`
	ModelName            *string               `json:"model_name"`
	StorageMode          *StorageMode          `json:"storage_mode"`
	AuthenticationMethod *AuthenticationMethod `json:"authentication_method"`
	OverwriteExisting    *bool                 `json:"overwrite_existing"`
}

type PowerBiTable struct {
	Name        *string      `json:"name"`
	Catalog     *string      `json:"catalog"`
	Schema      *string      `json:"schema"`
	StorageMode *StorageMode `json:"storage_mode"`
}

type PowerBiTask struct {
	Tables                 []PowerBiTable `json:"tables"`
	WarehouseId            *string        `json:"warehouse_id"`
	PowerBiModel           *PowerBiModel  `json:"power_bi_model"`
	ConnectionResourceName *string        `json:"connection_resource_name"`
	RefreshAfterUpdate     *bool          `json:"refresh_after_update"`
}

type PythonPyPiLibrary struct {
	Package *string `json:"package"`
	Repo    *string `json:"repo"`
}

type PythonWheelTask struct {
	PackageName     *string           `json:"package_name"`
	EntryPoint      *string           `json:"entry_point"`
	Parameters      []string          `json:"parameters"`
	NamedParameters map[string]string `json:"named_parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type PythonWheelTask_NamedParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type QueueDetails struct {
	Code    *QueueDetailsCode_Code `json:"code"`
	Message *string                `json:"message"`
}

type QueueDetailsCode struct {
}

type QueueSettings struct {
	Enabled *bool `json:"enabled"`
}

type RCranLibrary struct {
	Package *string `json:"package"`
	Repo    *string `json:"repo"`
}

type Repair struct {
	Type                       *RepairType                          `json:"type"`
	StartTime                  *int64                               `json:"start_time"`
	EndTime                    *int64                               `json:"end_time"`
	State                      *RunState                            `json:"state"`
	Id                         *int64                               `json:"id"`
	TaskRunIds                 []int64                              `json:"task_run_ids"`
	Status                     *RunStatus                           `json:"status"`
	EffectivePerformanceTarget *PerformanceTarget_PerformanceTarget `json:"effective_performance_target"`
}

type RepairRun struct {
	RunId               *int64                               `json:"run_id"`
	LatestRepairId      *int64                               `json:"latest_repair_id"`
	RerunTasks          []string                             `json:"rerun_tasks"`
	JobParameters       map[string]string                    `json:"job_parameters"`
	RerunAllFailedTasks *bool                                `json:"rerun_all_failed_tasks"`
	RerunDependentTasks *bool                                `json:"rerun_dependent_tasks"`
	PerformanceTarget   *PerformanceTarget_PerformanceTarget `json:"performance_target"`
	PipelineParams      *PipelineParameters                  `json:"pipeline_params"`
	JarParams           []string                             `json:"jar_params"`
	NotebookParams      map[string]string                    `json:"notebook_params"`
	PythonParams        []string                             `json:"python_params"`
	SparkSubmitParams   []string                             `json:"spark_submit_params"`
	PythonNamedParams   map[string]string                    `json:"python_named_params"`
	SqlParams           map[string]string                    `json:"sql_params"`
	DbtCommands         []string                             `json:"dbt_commands"`
}

// Name-based parameters for jobs running notebook tasks..
type RepairRun_JobParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RepairRun_NotebookParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RepairRun_PythonNamedParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Run repair was initiated..
type RepairRun_Response struct {
	RepairId *int64 `json:"repair_id"`
}

// Name-based parameters for jobs running notebook tasks..
type RepairRun_SqlParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ResetJob struct {
	JobId       *int64       `json:"job_id"`
	NewSettings *JobSettings `json:"new_settings"`
}

// Job was overwritten successfully..
type ResetJob_Response struct {
}

type ResolvedValues struct {
	NotebookTask    *ResolvedValues_NotebookTaskResolvedValues    `json:"notebook_task"`
	SparkJarTask    *ResolvedValues_SparkJarTaskResolvedValues    `json:"spark_jar_task"`
	SparkPythonTask *ResolvedValues_SparkPythonTaskResolvedValues `json:"spark_python_task"`
	SparkSubmitTask *ResolvedValues_SparkSubmitTaskResolvedValues `json:"spark_submit_task"`
	PythonWheelTask *ResolvedValues_PythonWheelTaskResolvedValues `json:"python_wheel_task"`
	DbtTask         *ResolvedValues_DbtTaskResolvedValues         `json:"dbt_task"`
	SqlTask         *ResolvedValues_SqlTaskResolvedValues         `json:"sql_task"`
	RunJobTask      *ResolvedValues_RunJobTaskResolvedValues      `json:"run_job_task"`
	ConditionTask   *ResolvedValues_ConditionTaskResolvedValues   `json:"condition_task"`
	SimulationTask  *ResolvedValues_SimulationTaskResolvedValues  `json:"simulation_task"`
}

type ResolvedValues_ConditionTaskResolvedValues struct {
	Left  *string `json:"left"`
	Right *string `json:"right"`
}

type ResolvedValues_DbtTaskResolvedValues struct {
	Commands []string `json:"commands"`
}

type ResolvedValues_NotebookTaskResolvedValues struct {
	BaseParameters map[string]string `json:"base_parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_NotebookTaskResolvedValues_BaseParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ResolvedValues_PythonWheelTaskResolvedValues struct {
	Parameters      []string          `json:"parameters"`
	NamedParameters map[string]string `json:"named_parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_PythonWheelTaskResolvedValues_NamedParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ResolvedValues_RunJobTaskResolvedValues struct {
	Parameters    map[string]string `json:"parameters"`
	JobParameters map[string]string `json:"job_parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_RunJobTaskResolvedValues_JobParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_RunJobTaskResolvedValues_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ResolvedValues_SimulationTaskResolvedValues struct {
	Parameters map[string]string `json:"parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_SimulationTaskResolvedValues_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ResolvedValues_SparkJarTaskResolvedValues struct {
	Parameters []string `json:"parameters"`
}

type ResolvedValues_SparkPythonTaskResolvedValues struct {
}

type ResolvedValues_SparkSubmitTaskResolvedValues struct {
}

type ResolvedValues_SqlTaskResolvedValues struct {
	Parameters map[string]string `json:"parameters"`
}

// Name-based parameters for jobs running notebook tasks..
type ResolvedValues_SqlTaskResolvedValues_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type Run struct {
	JobId                      *int64                               `json:"job_id"`
	RunId                      *int64                               `json:"run_id"`
	CreatorUserName            *string                              `json:"creator_user_name"`
	NumberInJob                *int64                               `json:"number_in_job"`
	OriginalAttemptRunId       *int64                               `json:"original_attempt_run_id"`
	State                      *RunState                            `json:"state"`
	Schedule                   *CronSchedule                        `json:"schedule"`
	ClusterSpec                *ClusterSpec                         `json:"cluster_spec"`
	ClusterInstance            *ClusterInstance                     `json:"cluster_instance"`
	JobParameters              []Run_JobLevelParameters             `json:"job_parameters"`
	OverridingParameters       *RunParameters                       `json:"overriding_parameters"`
	Trigger                    *TriggerType                         `json:"trigger"`
	TriggerInfo                *RunTriggerInfo                      `json:"trigger_info"`
	RunName                    *string                              `json:"run_name"`
	RunPageUrl                 *string                              `json:"run_page_url"`
	RunType                    *RunType                             `json:"run_type"`
	Tasks                      []RunTask                            `json:"tasks"`
	Description                *string                              `json:"description"`
	AttemptNumber              *int                                 `json:"attempt_number"`
	JobClusters                []JobCluster                         `json:"job_clusters"`
	GitSource                  *GitSource                           `json:"git_source"`
	RepairHistory              []Repair                             `json:"repair_history"`
	Status                     *RunStatus                           `json:"status"`
	JobRunId                   *int64                               `json:"job_run_id"`
	HasMore                    *bool                                `json:"has_more"`
	EffectivePerformanceTarget *PerformanceTarget_PerformanceTarget `json:"effective_performance_target"`
	EffectiveUsagePolicyId     *string                              `json:"effective_usage_policy_id"`
	StartTime                  *int64                               `json:"start_time"`
	SetupDuration              *int64                               `json:"setup_duration"`
	ExecutionDuration          *int64                               `json:"execution_duration"`
	CleanupDuration            *int64                               `json:"cleanup_duration"`
	EndTime                    *int64                               `json:"end_time"`
	RunDuration                *int64                               `json:"run_duration"`
	QueueDuration              *int64                               `json:"queue_duration"`
}

type Run_JobLevelParameters struct {
	Name    *string `json:"name"`
	Default *string `json:"default"`
	Value   *string `json:"value"`
}

type RunJobTask struct {
	JobId             *int64              `json:"job_id"`
	JobParameters     map[string]string   `json:"job_parameters"`
	PipelineParams    *PipelineParameters `json:"pipeline_params"`
	JarParams         []string            `json:"jar_params"`
	NotebookParams    map[string]string   `json:"notebook_params"`
	PythonParams      []string            `json:"python_params"`
	SparkSubmitParams []string            `json:"spark_submit_params"`
	PythonNamedParams map[string]string   `json:"python_named_params"`
	SqlParams         map[string]string   `json:"sql_params"`
	DbtCommands       []string            `json:"dbt_commands"`
}

// Name-based parameters for jobs running notebook tasks..
type RunJobTask_JobParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunJobTask_NotebookParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunJobTask_PythonNamedParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type RunJobTask_RunJobTaskOutput struct {
	RunId *int64 `json:"run_id"`
}

// Name-based parameters for jobs running notebook tasks..
type RunJobTask_SqlParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type RunLifeCycleState struct {
}

type RunLifecycleStateV2 struct {
}

type RunNow struct {
	JobId             *int64                               `json:"job_id"`
	JobParameters     map[string]string                    `json:"job_parameters"`
	IdempotencyToken  *string                              `json:"idempotency_token"`
	Queue             *QueueSettings                       `json:"queue"`
	Only              []string                             `json:"only"`
	PerformanceTarget *PerformanceTarget_PerformanceTarget `json:"performance_target"`
	PipelineParams    *PipelineParameters                  `json:"pipeline_params"`
	JarParams         []string                             `json:"jar_params"`
	NotebookParams    map[string]string                    `json:"notebook_params"`
	PythonParams      []string                             `json:"python_params"`
	SparkSubmitParams []string                             `json:"spark_submit_params"`
	PythonNamedParams map[string]string                    `json:"python_named_params"`
	SqlParams         map[string]string                    `json:"sql_params"`
	DbtCommands       []string                             `json:"dbt_commands"`
}

// Name-based parameters for jobs running notebook tasks..
type RunNow_JobParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunNow_NotebookParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunNow_PythonNamedParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Run was started successfully..
type RunNow_Response struct {
	RunId       *int64 `json:"run_id"`
	NumberInJob *int64 `json:"number_in_job"`
}

// Name-based parameters for jobs running notebook tasks..
type RunNow_SqlParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type RunParameters struct {
	PipelineParams    *PipelineParameters `json:"pipeline_params"`
	JarParams         []string            `json:"jar_params"`
	NotebookParams    map[string]string   `json:"notebook_params"`
	PythonParams      []string            `json:"python_params"`
	SparkSubmitParams []string            `json:"spark_submit_params"`
	PythonNamedParams map[string]string   `json:"python_named_params"`
	SqlParams         map[string]string   `json:"sql_params"`
	DbtCommands       []string            `json:"dbt_commands"`
}

// Name-based parameters for jobs running notebook tasks..
type RunParameters_NotebookParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunParameters_PythonNamedParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Name-based parameters for jobs running notebook tasks..
type RunParameters_SqlParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type RunResultState struct {
}

// The current state of the run..
type RunState struct {
	LifeCycleState          *RunLifeCycleState_RunLifeCycleState `json:"life_cycle_state"`
	ResultState             *RunResultState_RunResultState       `json:"result_state"`
	StateMessage            *string                              `json:"state_message"`
	UserCancelledOrTimedout *bool                                `json:"user_cancelled_or_timedout"`
	QueueReason             *string                              `json:"queue_reason"`
}

// The current status of the run.
type RunStatus struct {
	State              *RunLifecycleStateV2_State `json:"state"`
	TerminationDetails *TerminationDetails        `json:"termination_details"`
	QueueDetails       *QueueDetails              `json:"queue_details"`
}

// Used when outputting a child run, in GetRun or ListRuns..
type RunTask struct {
	RunId                      *int64                               `json:"run_id"`
	State                      *RunState                            `json:"state"`
	RunPageUrl                 *string                              `json:"run_page_url"`
	ClusterInstance            *ClusterInstance                     `json:"cluster_instance"`
	AttemptNumber              *int                                 `json:"attempt_number"`
	GitSource                  *GitSource                           `json:"git_source"`
	ResolvedValues             *ResolvedValues                      `json:"resolved_values"`
	Status                     *RunStatus                           `json:"status"`
	EffectivePerformanceTarget *PerformanceTarget_PerformanceTarget `json:"effective_performance_target"`
	TaskKey                    *string                              `json:"task_key"`
	Description                *string                              `json:"description"`
	DependsOn                  []TaskDependency                     `json:"depends_on"`
	RunIf                      *TaskDependencyType                  `json:"run_if"`
	TimeoutSeconds             *int                                 `json:"timeout_seconds"`
	EmailNotifications         *JobEmailNotifications               `json:"email_notifications"`
	Health                     *JobsHealthRules                     `json:"health"`
	NotificationSettings       *NotificationSettings                `json:"notification_settings"`
	WebhookNotifications       *WebhookNotifications                `json:"webhook_notifications"`
	EnvironmentKey             *string                              `json:"environment_key"`
	Disabled                   *bool                                `json:"disabled"`
	Compute                    *Compute                             `json:"compute"`
	NotebookTask               *NotebookTask                        `json:"notebook_task"`
	SparkJarTask               *SparkJarTask                        `json:"spark_jar_task"`
	SparkPythonTask            *SparkPythonTask                     `json:"spark_python_task"`
	SparkSubmitTask            *SparkSubmitTask                     `json:"spark_submit_task"`
	PipelineTask               *PipelineTask                        `json:"pipeline_task"`
	PythonWheelTask            *PythonWheelTask                     `json:"python_wheel_task"`
	DbtTask                    *DbtTask                             `json:"dbt_task"`
	SqlTask                    *SqlTask                             `json:"sql_task"`
	RunJobTask                 *RunJobTask                          `json:"run_job_task"`
	ConditionTask              *ConditionTask                       `json:"condition_task"`
	ForEachTask                *ForEachTask                         `json:"for_each_task"`
	CleanRoomsNotebookTask     *CleanRoomsNotebookTask              `json:"clean_rooms_notebook_task"`
	GenAiComputeTask           *GenAiComputeTask                    `json:"gen_ai_compute_task"`
	AlertTask                  *AlertTask                           `json:"alert_task"`
	PowerBiTask                *PowerBiTask                         `json:"power_bi_task"`
	DashboardTask              *DashboardTask                       `json:"dashboard_task"`
	DbtCloudTask               *DbtCloudTask                        `json:"dbt_cloud_task"`
	DbtPlatformTask            *DbtPlatformTask                     `json:"dbt_platform_task"`
	AgenticTask                *AgenticTask                         `json:"agentic_task"`
	ExistingClusterId          *string                              `json:"existing_cluster_id"`
	NewCluster                 *ClusterSpec_NewCluster              `json:"new_cluster"`
	JobClusterKey              *string                              `json:"job_cluster_key"`
	Libraries                  []Library                            `json:"libraries"`
	MaxRetries                 *int                                 `json:"max_retries"`
	MinRetryIntervalMillis     *int                                 `json:"min_retry_interval_millis"`
	RetryOnTimeout             *bool                                `json:"retry_on_timeout"`
	DisableAutoOptimization    *bool                                `json:"disable_auto_optimization"`
	StartTime                  *int64                               `json:"start_time"`
	SetupDuration              *int64                               `json:"setup_duration"`
	ExecutionDuration          *int64                               `json:"execution_duration"`
	CleanupDuration            *int64                               `json:"cleanup_duration"`
	EndTime                    *int64                               `json:"end_time"`
	RunDuration                *int64                               `json:"run_duration"`
	QueueDuration              *int64                               `json:"queue_duration"`
}

type RunTaskSettings struct {
	TaskKey                 *string                 `json:"task_key"`
	Description             *string                 `json:"description"`
	DependsOn               []TaskDependency        `json:"depends_on"`
	RunIf                   *TaskDependencyType     `json:"run_if"`
	TimeoutSeconds          *int                    `json:"timeout_seconds"`
	EmailNotifications      *JobEmailNotifications  `json:"email_notifications"`
	Health                  *JobsHealthRules        `json:"health"`
	NotificationSettings    *NotificationSettings   `json:"notification_settings"`
	WebhookNotifications    *WebhookNotifications   `json:"webhook_notifications"`
	EnvironmentKey          *string                 `json:"environment_key"`
	Disabled                *bool                   `json:"disabled"`
	Compute                 *Compute                `json:"compute"`
	NotebookTask            *NotebookTask           `json:"notebook_task"`
	SparkJarTask            *SparkJarTask           `json:"spark_jar_task"`
	SparkPythonTask         *SparkPythonTask        `json:"spark_python_task"`
	SparkSubmitTask         *SparkSubmitTask        `json:"spark_submit_task"`
	PipelineTask            *PipelineTask           `json:"pipeline_task"`
	PythonWheelTask         *PythonWheelTask        `json:"python_wheel_task"`
	DbtTask                 *DbtTask                `json:"dbt_task"`
	SqlTask                 *SqlTask                `json:"sql_task"`
	RunJobTask              *RunJobTask             `json:"run_job_task"`
	ConditionTask           *ConditionTask          `json:"condition_task"`
	ForEachTask             *ForEachTask            `json:"for_each_task"`
	CleanRoomsNotebookTask  *CleanRoomsNotebookTask `json:"clean_rooms_notebook_task"`
	GenAiComputeTask        *GenAiComputeTask       `json:"gen_ai_compute_task"`
	AlertTask               *AlertTask              `json:"alert_task"`
	PowerBiTask             *PowerBiTask            `json:"power_bi_task"`
	DashboardTask           *DashboardTask          `json:"dashboard_task"`
	DbtCloudTask            *DbtCloudTask           `json:"dbt_cloud_task"`
	DbtPlatformTask         *DbtPlatformTask        `json:"dbt_platform_task"`
	AgenticTask             *AgenticTask            `json:"agentic_task"`
	ExistingClusterId       *string                 `json:"existing_cluster_id"`
	NewCluster              *ClusterSpec_NewCluster `json:"new_cluster"`
	JobClusterKey           *string                 `json:"job_cluster_key"`
	Libraries               []Library               `json:"libraries"`
	MaxRetries              *int                    `json:"max_retries"`
	MinRetryIntervalMillis  *int                    `json:"min_retry_interval_millis"`
	RetryOnTimeout          *bool                   `json:"retry_on_timeout"`
	DisableAutoOptimization *bool                   `json:"disable_auto_optimization"`
}

// Additional details about what triggered the run.
type RunTriggerInfo struct {
	SqlCondition *SqlConditionRunInfoDetails `json:"sql_condition"`
	RunId        *int64                      `json:"run_id"`
}

// A storage location in Amazon S3.
type S3StorageInfo struct {
	Destination      *string `json:"destination"`
	Region           *string `json:"region"`
	Endpoint         *string `json:"endpoint"`
	EnableEncryption *bool   `json:"enable_encryption"`
	EncryptionType   *string `json:"encryption_type"`
	KmsKey           *string `json:"kms_key"`
	CannedAcl        *string `json:"canned_acl"`
}

type SparkJarTask struct {
	JarUri        *string  `json:"jar_uri"`
	MainClassName *string  `json:"main_class_name"`
	Parameters    []string `json:"parameters"`
	RunAsRepl     *bool    `json:"run_as_repl"`
}

type SparkPythonTask struct {
	PythonFile *string  `json:"python_file"`
	Parameters []string `json:"parameters"`
	Source     *Source  `json:"source"`
}

type SparkSubmitTask struct {
	Parameters []string `json:"parameters"`
}

type SparseCheckout struct {
	Patterns []string `json:"patterns"`
}

type SqlAlertState struct {
}

type SqlConditionConfiguration struct {
	SqlQueryId  *string `json:"sql_query_id"`
	WarehouseId *string `json:"warehouse_id"`
}

// SQL condition evaluation details captured at the time the run was triggered.
type SqlConditionRunInfoDetails struct {
	ConditionEvaluationSqlStatementId *string `json:"condition_evaluation_sql_statement_id"`
	ConditionEvaluationSatisfied      *bool   `json:"condition_evaluation_satisfied"`
	ConditionEvaluationSqlSessionId   *string `json:"condition_evaluation_sql_session_id"`
}

type SqlConditionState struct {
	LatestConditionEvaluationSqlStatementId *string `json:"latest_condition_evaluation_sql_statement_id"`
	LatestConditionEvaluationSatisfied      *bool   `json:"latest_condition_evaluation_satisfied"`
	LatestConditionEvaluationSqlSessionId   *string `json:"latest_condition_evaluation_sql_session_id"`
}

type SqlTask struct {
	Parameters  map[string]string `json:"parameters"`
	Query       *SqlTaskQuery     `json:"query"`
	Dashboard   *SqlTaskDashboard `json:"dashboard"`
	Alert       *SqlTaskAlert     `json:"alert"`
	File        *SqlTaskFile      `json:"file"`
	WarehouseId *string           `json:"warehouse_id"`
}

// Name-based parameters for jobs running notebook tasks..
type SqlTask_ParametersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type SqlTask_SqlAlertOutput struct {
	QueryText     *string                      `json:"query_text"`
	SqlStatements []SqlTask_SqlStatementOutput `json:"sql_statements"`
	OutputLink    *string                      `json:"output_link"`
	WarehouseId   *string                      `json:"warehouse_id"`
	AlertState    *SqlAlertState_SqlAlertState `json:"alert_state"`
}

type SqlTask_SqlDashboardOutput struct {
	Widgets     []SqlTask_SqlDashboardWidgetOutput `json:"widgets"`
	WarehouseId *string                            `json:"warehouse_id"`
}

type SqlTask_SqlDashboardWidgetOutput struct {
	WidgetId    *string                     `json:"widget_id"`
	WidgetTitle *string                     `json:"widget_title"`
	OutputLink  *string                     `json:"output_link"`
	Status      *SqlTask_SqlTaskQueryStatus `json:"status"`
	Error       *SqlTask_SqlOutputError     `json:"error"`
	StartTime   *int64                      `json:"start_time"`
	EndTime     *int64                      `json:"end_time"`
}

type SqlTask_SqlOutput struct {
	QueryOutput     *SqlTask_SqlQueryOutput     `json:"query_output"`
	DashboardOutput *SqlTask_SqlDashboardOutput `json:"dashboard_output"`
	AlertOutput     *SqlTask_SqlAlertOutput     `json:"alert_output"`
}

type SqlTask_SqlOutputError struct {
	Message *string `json:"message"`
}

type SqlTask_SqlQueryOutput struct {
	QueryText     *string                      `json:"query_text"`
	EndpointId    *string                      `json:"endpoint_id"`
	SqlStatements []SqlTask_SqlStatementOutput `json:"sql_statements"`
	OutputLink    *string                      `json:"output_link"`
	WarehouseId   *string                      `json:"warehouse_id"`
}

type SqlTask_SqlStatementOutput struct {
	LookupKey *string `json:"lookup_key"`
}

type SqlTaskAlert struct {
	AlertId            *string               `json:"alert_id"`
	Subscriptions      []SqlTaskSubscription `json:"subscriptions"`
	PauseSubscriptions *bool                 `json:"pause_subscriptions"`
}

type SqlTaskDashboard struct {
	DashboardId        *string               `json:"dashboard_id"`
	Subscriptions      []SqlTaskSubscription `json:"subscriptions"`
	CustomSubject      *string               `json:"custom_subject"`
	PauseSubscriptions *bool                 `json:"pause_subscriptions"`
}

type SqlTaskFile struct {
	Path   *string `json:"path"`
	Source *Source `json:"source"`
}

type SqlTaskQuery struct {
	QueryId *string `json:"query_id"`
}

type SqlTaskSubscription struct {
	UserName      *string `json:"user_name"`
	DestinationId *string `json:"destination_id"`
}

type SubmitRun struct {
	AccessControlList    []AccessControlRequest `json:"access_control_list"`
	Queue                *QueueSettings         `json:"queue"`
	RunAs                *JobRunAs              `json:"run_as"`
	RunName              *string                `json:"run_name"`
	TimeoutSeconds       *int                   `json:"timeout_seconds"`
	Health               *JobsHealthRules       `json:"health"`
	IdempotencyToken     *string                `json:"idempotency_token"`
	Tasks                []RunTaskSettings      `json:"tasks"`
	GitSource            *GitSource             `json:"git_source"`
	WebhookNotifications *WebhookNotifications  `json:"webhook_notifications"`
	EmailNotifications   *JobEmailNotifications `json:"email_notifications"`
	NotificationSettings *NotificationSettings  `json:"notification_settings"`
	Environments         []JobEnvironment       `json:"environments"`
	BudgetPolicyId       *string                `json:"budget_policy_id"`
	UsagePolicyId        *string                `json:"usage_policy_id"`
}

// Run was created and started successfully..
type SubmitRun_Response struct {
	RunId *int64 `json:"run_id"`
}

type Subscription struct {
	Subscribers   []Subscription_Subscriber `json:"subscribers"`
	Paused        *bool                     `json:"paused"`
	CustomSubject *string                   `json:"custom_subject"`
}

type Subscription_Subscriber struct {
	UserName      *string `json:"user_name"`
	DestinationId *string `json:"destination_id"`
}

// Configuration for a Supervisor Agent..
type SupervisorAgent struct {
	AgentId *string `json:"agent_id"`
}

type TableState struct {
	TableName      *string `json:"table_name"`
	HasSeenUpdates *bool   `json:"has_seen_updates"`
}

type TableTriggerConfiguration struct {
	TableNames                    []string                             `json:"table_names"`
	MinTimeBetweenTriggersSeconds *int                                 `json:"min_time_between_triggers_seconds"`
	WaitAfterLastChangeSeconds    *int                                 `json:"wait_after_last_change_seconds"`
	Condition                     *TableTriggerConfiguration_Condition `json:"condition"`
}

type TableTriggerState struct {
	LastSeenTableStates     []TableState `json:"last_seen_table_states"`
	UsingScalableMonitoring *bool        `json:"using_scalable_monitoring"`
}

type TaskDependency struct {
	TaskKey *string `json:"task_key"`
	Outcome *string `json:"outcome"`
}

type TaskSettings struct {
	TaskKey                 *string                 `json:"task_key"`
	DependsOn               []TaskDependency        `json:"depends_on"`
	RunIf                   *TaskDependencyType     `json:"run_if"`
	TimeoutSeconds          *int                    `json:"timeout_seconds"`
	Health                  *JobsHealthRules        `json:"health"`
	EmailNotifications      *JobEmailNotifications  `json:"email_notifications"`
	NotificationSettings    *NotificationSettings   `json:"notification_settings"`
	WebhookNotifications    *WebhookNotifications   `json:"webhook_notifications"`
	Description             *string                 `json:"description"`
	EnvironmentKey          *string                 `json:"environment_key"`
	Disabled                *bool                   `json:"disabled"`
	Compute                 *Compute                `json:"compute"`
	NotebookTask            *NotebookTask           `json:"notebook_task"`
	SparkJarTask            *SparkJarTask           `json:"spark_jar_task"`
	SparkPythonTask         *SparkPythonTask        `json:"spark_python_task"`
	SparkSubmitTask         *SparkSubmitTask        `json:"spark_submit_task"`
	PipelineTask            *PipelineTask           `json:"pipeline_task"`
	PythonWheelTask         *PythonWheelTask        `json:"python_wheel_task"`
	DbtTask                 *DbtTask                `json:"dbt_task"`
	SqlTask                 *SqlTask                `json:"sql_task"`
	RunJobTask              *RunJobTask             `json:"run_job_task"`
	ConditionTask           *ConditionTask          `json:"condition_task"`
	ForEachTask             *ForEachTask            `json:"for_each_task"`
	CleanRoomsNotebookTask  *CleanRoomsNotebookTask `json:"clean_rooms_notebook_task"`
	GenAiComputeTask        *GenAiComputeTask       `json:"gen_ai_compute_task"`
	AlertTask               *AlertTask              `json:"alert_task"`
	PowerBiTask             *PowerBiTask            `json:"power_bi_task"`
	DashboardTask           *DashboardTask          `json:"dashboard_task"`
	DbtCloudTask            *DbtCloudTask           `json:"dbt_cloud_task"`
	DbtPlatformTask         *DbtPlatformTask        `json:"dbt_platform_task"`
	AgenticTask             *AgenticTask            `json:"agentic_task"`
	ExistingClusterId       *string                 `json:"existing_cluster_id"`
	NewCluster              *ClusterSpec_NewCluster `json:"new_cluster"`
	JobClusterKey           *string                 `json:"job_cluster_key"`
	Libraries               []Library               `json:"libraries"`
	MaxRetries              *int                    `json:"max_retries"`
	MinRetryIntervalMillis  *int                    `json:"min_retry_interval_millis"`
	RetryOnTimeout          *bool                   `json:"retry_on_timeout"`
	DisableAutoOptimization *bool                   `json:"disable_auto_optimization"`
}

type TerminationCode struct {
}

type TerminationDetails struct {
	Code    *TerminationCode_Code `json:"code"`
	Type    *TerminationType_Type `json:"type"`
	Message *string               `json:"message"`
}

type TerminationType struct {
}

type TriggerSettings struct {
	PauseStatus  *SchedulePauseStatus             `json:"pause_status"`
	FileArrival  *FileArrivalTriggerConfiguration `json:"file_arrival"`
	Periodic     *PeriodicTriggerConfiguration    `json:"periodic"`
	TableUpdate  *TableTriggerConfiguration       `json:"table_update"`
	Model        *ModelTriggerConfiguration       `json:"model"`
	SqlCondition *SqlConditionConfiguration       `json:"sql_condition"`
}

type TriggerStateProto struct {
	Table        *TableTriggerState       `json:"table"`
	FileArrival  *FileArrivalTriggerState `json:"file_arrival"`
	SqlCondition *SqlConditionState       `json:"sql_condition"`
}

type UpdateJob struct {
	JobId          *int64       `json:"job_id"`
	NewSettings    *JobSettings `json:"new_settings"`
	FieldsToRemove []string     `json:"fields_to_remove"`
}

// Job was updated successfully..
type UpdateJob_Response struct {
}

type ViewItem struct {
	Content *string   `json:"content"`
	Name    *string   `json:"name"`
	Type    *ViewType `json:"type"`
}

// A storage location back by UC Volumes..
type VolumesStorageInfo struct {
	Destination *string `json:"destination"`
}

type Webhook struct {
	Id *string `json:"id"`
}

type WebhookNotifications struct {
	OnStart                            []Webhook `json:"on_start"`
	OnSuccess                          []Webhook `json:"on_success"`
	OnFailure                          []Webhook `json:"on_failure"`
	OnDurationWarningThresholdExceeded []Webhook `json:"on_duration_warning_threshold_exceeded"`
	OnStreamingBacklogExceeded         []Webhook `json:"on_streaming_backlog_exceeded"`
}

type WidgetErrorDetail struct {
	Message *string `json:"message"`
}

// Cluster Attributes showing for clusters workload types..
type WorkloadType struct {
	Clients *WorkloadType_ClientsTypes `json:"clients"`
}

type WorkloadType_ClientsTypes struct {
	Notebooks *bool `json:"notebooks"`
	Jobs      *bool `json:"jobs"`
}

// A storage location in Workspace Filesystem (WSFS).
type WorkspaceStorageInfo struct {
	Destination *string `json:"destination"`
}
