package v1

type ChannelName string

const (
	ChannelName_ChannelNameUnspecified ChannelName = "CHANNEL_NAME_UNSPECIFIED"
	ChannelName_ChannelNamePreview     ChannelName = "CHANNEL_NAME_PREVIEW"
	ChannelName_ChannelNameCurrent     ChannelName = "CHANNEL_NAME_CURRENT"
	ChannelName_ChannelNamePrevious    ChannelName = "CHANNEL_NAME_PREVIOUS"
	ChannelName_ChannelNameCustom      ChannelName = "CHANNEL_NAME_CUSTOM"
)

// String representation for [fmt.Print].
func (f *ChannelName) String() string {
	return string(*f)
}

type PlansState string

const (
	PlansState_IgnoredSmallDuration  PlansState = "IGNORED_SMALL_DURATION"
	PlansState_IgnoredLargePlansSize PlansState = "IGNORED_LARGE_PLANS_SIZE"
	PlansState_Exists                PlansState = "EXISTS"
	PlansState_Unknown               PlansState = "UNKNOWN"
	PlansState_Empty                 PlansState = "EMPTY"
	PlansState_IgnoredSparkPlanType  PlansState = "IGNORED_SPARK_PLAN_TYPE"
)

// String representation for [fmt.Print].
func (f *PlansState) String() string {
	return string(*f)
}

type QueryStatementType string

const (
	QueryStatementType_Other    QueryStatementType = "OTHER"
	QueryStatementType_Alter    QueryStatementType = "ALTER"
	QueryStatementType_Analyze  QueryStatementType = "ANALYZE"
	QueryStatementType_Copy     QueryStatementType = "COPY"
	QueryStatementType_Create   QueryStatementType = "CREATE"
	QueryStatementType_Delete   QueryStatementType = "DELETE"
	QueryStatementType_Describe QueryStatementType = "DESCRIBE"
	QueryStatementType_Drop     QueryStatementType = "DROP"
	QueryStatementType_Explain  QueryStatementType = "EXPLAIN"
	QueryStatementType_Grant    QueryStatementType = "GRANT"
	QueryStatementType_Insert   QueryStatementType = "INSERT"
	QueryStatementType_Merge    QueryStatementType = "MERGE"
	QueryStatementType_Optimize QueryStatementType = "OPTIMIZE"
	QueryStatementType_Refresh  QueryStatementType = "REFRESH"
	QueryStatementType_Replace  QueryStatementType = "REPLACE"
	QueryStatementType_Revoke   QueryStatementType = "REVOKE"
	QueryStatementType_Select   QueryStatementType = "SELECT"
	QueryStatementType_Set      QueryStatementType = "SET"
	QueryStatementType_Show     QueryStatementType = "SHOW"
	QueryStatementType_Truncate QueryStatementType = "TRUNCATE"
	QueryStatementType_Update   QueryStatementType = "UPDATE"
	QueryStatementType_Use      QueryStatementType = "USE"
	QueryStatementType_Call     QueryStatementType = "CALL"
)

// String representation for [fmt.Print].
func (f *QueryStatementType) String() string {
	return string(*f)
}

type QueryStatus string

const (
	QueryStatus_Queued    QueryStatus = "QUEUED"
	QueryStatus_Started   QueryStatus = "STARTED"
	QueryStatus_Compiling QueryStatus = "COMPILING"
	QueryStatus_Compiled  QueryStatus = "COMPILED"
	QueryStatus_Running   QueryStatus = "RUNNING"
	QueryStatus_Canceled  QueryStatus = "CANCELED"
	QueryStatus_Failed    QueryStatus = "FAILED"
	QueryStatus_Finished  QueryStatus = "FINISHED"
)

// String representation for [fmt.Print].
func (f *QueryStatus) String() string {
	return string(*f)
}

// Details about a Channel..
type ChannelInfo struct {
	Name         *ChannelName `json:"name"`
	DbsqlVersion *string      `json:"dbsql_version"`
}

type ExternalQuerySource struct {
	DashboardId       *string                      `json:"dashboard_id"`
	LegacyDashboardId *string                      `json:"legacy_dashboard_id"`
	AlertId           *string                      `json:"alert_id"`
	NotebookId        *string                      `json:"notebook_id"`
	SqlQueryId        *string                      `json:"sql_query_id"`
	JobInfo           *ExternalQuerySource_JobInfo `json:"job_info"`
	GenieSpaceId      *string                      `json:"genie_space_id"`
}

type ExternalQuerySource_JobInfo struct {
	JobId        *string `json:"job_id"`
	JobRunId     *string `json:"job_run_id"`
	JobTaskRunId *string `json:"job_task_run_id"`
}

// Fetches a list of queries conforming to the provided set of query filters.
//
// If the the number of queries to return takes > 10 seconds, the request will
// timeout. In that case, please reduce the time range to ensure ListQueries
// conforms to the 10 second max query time limit..
type ListQueries struct {
	FilterBy       *QueryFilter `json:"filter_by"`
	MaxResults     *int         `json:"max_results"`
	PageToken      *string      `json:"page_token"`
	IncludeMetrics *bool        `json:"include_metrics"`
}

type ListQueries_Response struct {
	NextPageToken *string     `json:"next_page_token"`
	HasNextPage   *bool       `json:"has_next_page"`
	Res           []QueryInfo `json:"res"`
}

type QueryFilter struct {
	QueryStartTimeRange *TimeRange    `json:"query_start_time_range"`
	UserIds             []int64       `json:"user_ids"`
	Statuses            []QueryStatus `json:"statuses"`
	WarehouseIds        []string      `json:"warehouse_ids"`
	StatementIds        []string      `json:"statement_ids"`
}

type QueryInfo struct {
	QueryId            *string              `json:"query_id"`
	Status             *QueryStatus         `json:"status"`
	QueryText          *string              `json:"query_text"`
	QueryStartTimeMs   *int64               `json:"query_start_time_ms"`
	ExecutionEndTimeMs *int64               `json:"execution_end_time_ms"`
	QueryEndTimeMs     *int64               `json:"query_end_time_ms"`
	UserId             *int64               `json:"user_id"`
	UserName           *string              `json:"user_name"`
	SparkUiUrl         *string              `json:"spark_ui_url"`
	EndpointId         *string              `json:"endpoint_id"`
	RowsProduced       *int64               `json:"rows_produced"`
	ErrorMessage       *string              `json:"error_message"`
	LookupKey          *string              `json:"lookup_key"`
	Metrics            *QueryMetrics        `json:"metrics"`
	ExecutedAsUserId   *int64               `json:"executed_as_user_id"`
	ExecutedAsUserName *string              `json:"executed_as_user_name"`
	SessionId          *string              `json:"session_id"`
	IsFinal            *bool                `json:"is_final"`
	ChannelUsed        *ChannelInfo         `json:"channel_used"`
	PlansState         *PlansState          `json:"plans_state"`
	StatementType      *QueryStatementType  `json:"statement_type"`
	WarehouseId        *string              `json:"warehouse_id"`
	Duration           *int64               `json:"duration"`
	ClientApplication  *string              `json:"client_application"`
	QuerySource        *ExternalQuerySource `json:"query_source"`
	CacheQueryId       *string              `json:"cache_query_id"`
	QueryTags          []QueryTag           `json:"query_tags"`
}

// A query metric that encapsulates a set of measurements for a single query.
// Metrics come from the driver and are stored in the history service database..
type QueryMetrics struct {
	TotalTimeMs                       *int64             `json:"total_time_ms"`
	ReadBytes                         *int64             `json:"read_bytes"`
	RowsProducedCount                 *int64             `json:"rows_produced_count"`
	CompilationTimeMs                 *int64             `json:"compilation_time_ms"`
	ExecutionTimeMs                   *int64             `json:"execution_time_ms"`
	ReadRemoteBytes                   *int64             `json:"read_remote_bytes"`
	WriteRemoteBytes                  *int64             `json:"write_remote_bytes"`
	ReadCacheBytes                    *int64             `json:"read_cache_bytes"`
	SpillToDiskBytes                  *int64             `json:"spill_to_disk_bytes"`
	TaskTotalTimeMs                   *int64             `json:"task_total_time_ms"`
	ReadFilesCount                    *int64             `json:"read_files_count"`
	ReadPartitionsCount               *int64             `json:"read_partitions_count"`
	PhotonTotalTimeMs                 *int64             `json:"photon_total_time_ms"`
	RowsReadCount                     *int64             `json:"rows_read_count"`
	ResultFetchTimeMs                 *int64             `json:"result_fetch_time_ms"`
	NetworkSentBytes                  *int64             `json:"network_sent_bytes"`
	ResultFromCache                   *bool              `json:"result_from_cache"`
	PrunedBytes                       *int64             `json:"pruned_bytes"`
	PrunedFilesCount                  *int64             `json:"pruned_files_count"`
	ProvisioningQueueStartTimestamp   *int64             `json:"provisioning_queue_start_timestamp"`
	OverloadingQueueStartTimestamp    *int64             `json:"overloading_queue_start_timestamp"`
	QueryCompilationStartTimestamp    *int64             `json:"query_compilation_start_timestamp"`
	TaskTimeOverTimeRange             *TaskTimeOverRange `json:"task_time_over_time_range"`
	WorkToBeDone                      *int64             `json:"work_to_be_done"`
	RunnableTasks                     *int64             `json:"runnable_tasks"`
	ProjectedRemainingTaskTotalTimeMs *int64             `json:"projected_remaining_task_total_time_ms"`
	RemainingTaskCount                *int64             `json:"remaining_task_count"`
	ProjectedRemainingWallclockTimeMs *int64             `json:"projected_remaining_wallclock_time_ms"`
	ReadFilesBytes                    *int64             `json:"read_files_bytes"`
}

// * A query execution can be annotated with an optional key-value pair to allow
// users to attribute the executions by key and optional value to filter by.
// QueryTag is the user-facing representation..
type QueryTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type TaskTimeOverRange struct {
	Entries  []TaskTimeOverRangeEntry `json:"entries"`
	Interval *int64                   `json:"interval"`
}

type TaskTimeOverRangeEntry struct {
	TaskCompletedTimeMs *int64 `json:"task_completed_time_ms"`
}

type TimeRange struct {
	StartTimeMs *int64 `json:"start_time_ms"`
	EndTimeMs   *int64 `json:"end_time_ms"`
}
