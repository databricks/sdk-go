package v1

type AggregationGranularity string

const (
	AggregationGranularityAggregationGranularityUnspecified AggregationGranularity = "AGGREGATION_GRANULARITY_UNSPECIFIED"
	AggregationGranularityAggregationGranularity5Minutes    AggregationGranularity = "AGGREGATION_GRANULARITY_5_MINUTES"
	AggregationGranularityAggregationGranularity30Minutes   AggregationGranularity = "AGGREGATION_GRANULARITY_30_MINUTES"
	AggregationGranularityAggregationGranularity1Hour       AggregationGranularity = "AGGREGATION_GRANULARITY_1_HOUR"
	AggregationGranularityAggregationGranularity1Day        AggregationGranularity = "AGGREGATION_GRANULARITY_1_DAY"
	AggregationGranularityAggregationGranularity1Week       AggregationGranularity = "AGGREGATION_GRANULARITY_1_WEEK"
	AggregationGranularityAggregationGranularity2Weeks      AggregationGranularity = "AGGREGATION_GRANULARITY_2_WEEKS"
	AggregationGranularityAggregationGranularity3Weeks      AggregationGranularity = "AGGREGATION_GRANULARITY_3_WEEKS"
	AggregationGranularityAggregationGranularity4Weeks      AggregationGranularity = "AGGREGATION_GRANULARITY_4_WEEKS"
	AggregationGranularityAggregationGranularity1Month      AggregationGranularity = "AGGREGATION_GRANULARITY_1_MONTH"
	AggregationGranularityAggregationGranularity1Year       AggregationGranularity = "AGGREGATION_GRANULARITY_1_YEAR"
)

// String representation for [fmt.Print].
func (f *AggregationGranularity) String() string {
	return string(*f)
}

type AnomalyDetectionJobType string

const (
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeUnspecified    AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED"
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeNormal         AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_NORMAL"
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeInternalHidden AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN"
)

// String representation for [fmt.Print].
func (f *AnomalyDetectionJobType) String() string {
	return string(*f)
}

type CronSchedulePauseStatus string

const (
	CronSchedulePauseStatusCronSchedulePauseStatusUnspecified CronSchedulePauseStatus = "CRON_SCHEDULE_PAUSE_STATUS_UNSPECIFIED"
	CronSchedulePauseStatusCronSchedulePauseStatusUnpaused    CronSchedulePauseStatus = "CRON_SCHEDULE_PAUSE_STATUS_UNPAUSED"
	CronSchedulePauseStatusCronSchedulePauseStatusPaused      CronSchedulePauseStatus = "CRON_SCHEDULE_PAUSE_STATUS_PAUSED"
)

// String representation for [fmt.Print].
func (f *CronSchedulePauseStatus) String() string {
	return string(*f)
}

type DataProfilingCustomMetricType string

const (
	DataProfilingCustomMetricTypeDataProfilingCustomMetricTypeUnspecified DataProfilingCustomMetricType = "DATA_PROFILING_CUSTOM_METRIC_TYPE_UNSPECIFIED"
	DataProfilingCustomMetricTypeDataProfilingCustomMetricTypeAggregate   DataProfilingCustomMetricType = "DATA_PROFILING_CUSTOM_METRIC_TYPE_AGGREGATE"
	DataProfilingCustomMetricTypeDataProfilingCustomMetricTypeDerived     DataProfilingCustomMetricType = "DATA_PROFILING_CUSTOM_METRIC_TYPE_DERIVED"
	DataProfilingCustomMetricTypeDataProfilingCustomMetricTypeDrift       DataProfilingCustomMetricType = "DATA_PROFILING_CUSTOM_METRIC_TYPE_DRIFT"
)

// String representation for [fmt.Print].
func (f *DataProfilingCustomMetricType) String() string {
	return string(*f)
}

type DataProfilingStatus string

const (
	DataProfilingStatusDataProfilingStatusUnspecified   DataProfilingStatus = "DATA_PROFILING_STATUS_UNSPECIFIED"
	DataProfilingStatusDataProfilingStatusActive        DataProfilingStatus = "DATA_PROFILING_STATUS_ACTIVE"
	DataProfilingStatusDataProfilingStatusPending       DataProfilingStatus = "DATA_PROFILING_STATUS_PENDING"
	DataProfilingStatusDataProfilingStatusDeletePending DataProfilingStatus = "DATA_PROFILING_STATUS_DELETE_PENDING"
	DataProfilingStatusDataProfilingStatusError         DataProfilingStatus = "DATA_PROFILING_STATUS_ERROR"
	DataProfilingStatusDataProfilingStatusFailed        DataProfilingStatus = "DATA_PROFILING_STATUS_FAILED"
)

// String representation for [fmt.Print].
func (f *DataProfilingStatus) String() string {
	return string(*f)
}

type InferenceProblemType string

const (
	InferenceProblemTypeInferenceProblemTypeUnspecified    InferenceProblemType = "INFERENCE_PROBLEM_TYPE_UNSPECIFIED"
	InferenceProblemTypeInferenceProblemTypeClassification InferenceProblemType = "INFERENCE_PROBLEM_TYPE_CLASSIFICATION"
	InferenceProblemTypeInferenceProblemTypeRegression     InferenceProblemType = "INFERENCE_PROBLEM_TYPE_REGRESSION"
)

// String representation for [fmt.Print].
func (f *InferenceProblemType) String() string {
	return string(*f)
}

type RefreshState string

const (
	RefreshStateMonitorRefreshStateUnknown  RefreshState = "MONITOR_REFRESH_STATE_UNKNOWN"
	RefreshStateMonitorRefreshStatePending  RefreshState = "MONITOR_REFRESH_STATE_PENDING"
	RefreshStateMonitorRefreshStateRunning  RefreshState = "MONITOR_REFRESH_STATE_RUNNING"
	RefreshStateMonitorRefreshStateSuccess  RefreshState = "MONITOR_REFRESH_STATE_SUCCESS"
	RefreshStateMonitorRefreshStateFailed   RefreshState = "MONITOR_REFRESH_STATE_FAILED"
	RefreshStateMonitorRefreshStateCanceled RefreshState = "MONITOR_REFRESH_STATE_CANCELED"
)

// String representation for [fmt.Print].
func (f *RefreshState) String() string {
	return string(*f)
}

type RefreshTrigger string

const (
	RefreshTriggerMonitorRefreshTriggerUnknown    RefreshTrigger = "MONITOR_REFRESH_TRIGGER_UNKNOWN"
	RefreshTriggerMonitorRefreshTriggerManual     RefreshTrigger = "MONITOR_REFRESH_TRIGGER_MANUAL"
	RefreshTriggerMonitorRefreshTriggerSchedule   RefreshTrigger = "MONITOR_REFRESH_TRIGGER_SCHEDULE"
	RefreshTriggerMonitorRefreshTriggerDataChange RefreshTrigger = "MONITOR_REFRESH_TRIGGER_DATA_CHANGE"
)

// String representation for [fmt.Print].
func (f *RefreshTrigger) String() string {
	return string(*f)
}

// Anomaly Detection Configurations..
type AnomalyDetectionConfig struct {
	AnomalyDetectionWorkflowId  *int64                       `json:"anomaly_detection_workflow_id"`
	PublishHealthIndicator      *bool                        `json:"publish_health_indicator"`
	JobType                     *AnomalyDetectionJobType     `json:"job_type"`
	ExcludedTableFullNames      []string                     `json:"excluded_table_full_names"`
	ValidityCheckConfigurations []ValidityCheckConfiguration `json:"validity_check_configurations"`
}

// Request to cancel a refresh..
type CancelRefreshRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
	RefreshId  *int64  `json:"refresh_id"`
}

// Response to cancelling a refresh..
type CancelRefreshResponse struct {
	Refresh *Refresh `json:"refresh"`
}

// Request to create a Monitor..
type CreateMonitorRequest struct {
	Monitor *Monitor `json:"monitor"`
}

// Request to create a refresh..
type CreateRefreshRequest struct {
	Refresh *Refresh `json:"refresh"`
}

// The data quality monitoring workflow cron schedule..
type CronSchedule struct {
	QuartzCronExpression *string                  `json:"quartz_cron_expression"`
	TimezoneId           *string                  `json:"timezone_id"`
	PauseStatus          *CronSchedulePauseStatus `json:"pause_status"`
}

// Data Profiling Configurations..
type DataProfilingConfig struct {
	OutputSchemaId              *string                     `json:"output_schema_id"`
	AssetsDir                   *string                     `json:"assets_dir"`
	InferenceLog                *InferenceLogConfig         `json:"inference_log"`
	TimeSeries                  *TimeSeriesConfig           `json:"time_series"`
	Snapshot                    *SnapshotConfig             `json:"snapshot"`
	SlicingExprs                []string                    `json:"slicing_exprs"`
	CustomMetrics               []DataProfilingCustomMetric `json:"custom_metrics"`
	BaselineTableName           *string                     `json:"baseline_table_name"`
	Schedule                    *CronSchedule               `json:"schedule"`
	NotificationSettings        *NotificationSettings       `json:"notification_settings"`
	SkipBuiltinDashboard        *bool                       `json:"skip_builtin_dashboard"`
	WarehouseId                 *string                     `json:"warehouse_id"`
	MonitoredTableName          *string                     `json:"monitored_table_name"`
	Status                      *DataProfilingStatus        `json:"status"`
	LatestMonitorFailureMessage *string                     `json:"latest_monitor_failure_message"`
	ProfileMetricsTableName     *string                     `json:"profile_metrics_table_name"`
	DriftMetricsTableName       *string                     `json:"drift_metrics_table_name"`
	DashboardId                 *string                     `json:"dashboard_id"`
	MonitorVersion              *int64                      `json:"monitor_version"`
	EffectiveWarehouseId        *string                     `json:"effective_warehouse_id"`
}

// Custom metric definition..
type DataProfilingCustomMetric struct {
	Name           *string                        `json:"name"`
	Definition     *string                        `json:"definition"`
	InputColumns   []string                       `json:"input_columns"`
	OutputDataType *string                        `json:"output_data_type"`
	Type           *DataProfilingCustomMetricType `json:"type"`
}

// Request to delete a Monitor..
type DeleteMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
}

// Request to delete a ronitor..
type DeleteRefreshRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
	RefreshId  *int64  `json:"refresh_id"`
}

// Request to get a Monitor..
type GetMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
}

// Request to get a refresh..
type GetRefreshRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
	RefreshId  *int64  `json:"refresh_id"`
}

// Inference log configuration..
type InferenceLogConfig struct {
	ProblemType                 *InferenceProblemType    `json:"problem_type"`
	TimestampColumn             *string                  `json:"timestamp_column"`
	Granularities               []AggregationGranularity `json:"granularities"`
	PredictionColumn            *string                  `json:"prediction_column"`
	LabelColumn                 *string                  `json:"label_column"`
	ModelIdColumn               *string                  `json:"model_id_column"`
	PredictionProbabilityColumn *string                  `json:"prediction_probability_column"`
}

// Request to list Monitors..
type ListMonitorRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

// Response for listing Monitors..
type ListMonitorResponse struct {
	Monitors      []Monitor `json:"monitors"`
	NextPageToken *string   `json:"next_page_token"`
}

// Request to list refreshes..
type ListRefreshRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
	PageToken  *string `json:"page_token"`
	PageSize   *int    `json:"page_size"`
}

// Response for listing refreshes..
type ListRefreshResponse struct {
	Refreshes     []Refresh `json:"refreshes"`
	NextPageToken *string   `json:"next_page_token"`
}

// Monitor for the data quality of unity catalog entities such as schema or
// table..
type Monitor struct {
	ObjectType             *string                 `json:"object_type"`
	ObjectId               *string                 `json:"object_id"`
	AnomalyDetectionConfig *AnomalyDetectionConfig `json:"anomaly_detection_config"`
	DataProfilingConfig    *DataProfilingConfig    `json:"data_profiling_config"`
}

// Destination of the data quality monitoring notification..
type NotificationDestination struct {
	EmailAddresses []string `json:"email_addresses"`
}

// Settings for sending notifications on the data quality monitoring..
type NotificationSettings struct {
	OnFailure *NotificationDestination `json:"on_failure"`
}

type PercentNullValidityCheck struct {
	ColumnNames []string `json:"column_names"`
	UpperBound  *float64 `json:"upper_bound"`
}

type RangeValidityCheck struct {
	ColumnNames []string `json:"column_names"`
	LowerBound  *float64 `json:"lower_bound"`
	UpperBound  *float64 `json:"upper_bound"`
}

// The Refresh object gives information on a refresh of the data quality
// monitoring pipeline..
type Refresh struct {
	ObjectType  *string         `json:"object_type"`
	ObjectId    *string         `json:"object_id"`
	RefreshId   *int64          `json:"refresh_id"`
	State       *RefreshState   `json:"state"`
	Message     *string         `json:"message"`
	StartTimeMs *int64          `json:"start_time_ms"`
	EndTimeMs   *int64          `json:"end_time_ms"`
	Trigger     *RefreshTrigger `json:"trigger"`
}

// Snapshot analysis configuration..
type SnapshotConfig struct {
}

// Time series analysis configuration..
type TimeSeriesConfig struct {
	TimestampColumn *string                  `json:"timestamp_column"`
	Granularities   []AggregationGranularity `json:"granularities"`
}

type UniquenessValidityCheck struct {
	ColumnNames []string `json:"column_names"`
}

// Request to update a Monitor..
type UpdateMonitorRequest struct {
	ObjectType *string  `json:"object_type"`
	ObjectId   *string  `json:"object_id"`
	Monitor    *Monitor `json:"monitor"`
	UpdateMask *string  `json:"update_mask"`
}

// Request to update a refresh..
type UpdateRefreshRequest struct {
	ObjectType *string  `json:"object_type"`
	ObjectId   *string  `json:"object_id"`
	RefreshId  *int64   `json:"refresh_id"`
	Refresh    *Refresh `json:"refresh"`
	UpdateMask *string  `json:"update_mask"`
}

type ValidityCheckConfiguration struct {
	Name                     *string                   `json:"name"`
	PercentNullValidityCheck *PercentNullValidityCheck `json:"percent_null_validity_check"`
	RangeValidityCheck       *RangeValidityCheck       `json:"range_validity_check"`
	UniquenessValidityCheck  *UniquenessValidityCheck  `json:"uniqueness_validity_check"`
}
