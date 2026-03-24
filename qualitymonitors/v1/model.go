package v1

type CustomMetricType string

const (
	CustomMetricType_CustomMetricTypeUnspecified CustomMetricType = "CUSTOM_METRIC_TYPE_UNSPECIFIED"
	CustomMetricType_CustomMetricTypeAggregate   CustomMetricType = "CUSTOM_METRIC_TYPE_AGGREGATE"
	CustomMetricType_CustomMetricTypeDerived     CustomMetricType = "CUSTOM_METRIC_TYPE_DERIVED"
	CustomMetricType_CustomMetricTypeDrift       CustomMetricType = "CUSTOM_METRIC_TYPE_DRIFT"
)

// String representation for [fmt.Print].
func (f *CustomMetricType) String() string {
	return string(*f)
}

type MonitorStatus string

const (
	MonitorStatus_MonitorStatusUnspecified   MonitorStatus = "MONITOR_STATUS_UNSPECIFIED"
	MonitorStatus_MonitorStatusActive        MonitorStatus = "MONITOR_STATUS_ACTIVE"
	MonitorStatus_MonitorStatusPending       MonitorStatus = "MONITOR_STATUS_PENDING"
	MonitorStatus_MonitorStatusDeletePending MonitorStatus = "MONITOR_STATUS_DELETE_PENDING"
	MonitorStatus_MonitorStatusError         MonitorStatus = "MONITOR_STATUS_ERROR"
	MonitorStatus_MonitorStatusFailed        MonitorStatus = "MONITOR_STATUS_FAILED"
)

// String representation for [fmt.Print].
func (f *MonitorStatus) String() string {
	return string(*f)
}

type ProblemType string

const (
	ProblemType_ProblemTypeUnspecified    ProblemType = "PROBLEM_TYPE_UNSPECIFIED"
	ProblemType_ProblemTypeClassification ProblemType = "PROBLEM_TYPE_CLASSIFICATION"
	ProblemType_ProblemTypeRegression     ProblemType = "PROBLEM_TYPE_REGRESSION"
)

// String representation for [fmt.Print].
func (f *ProblemType) String() string {
	return string(*f)
}

type RefreshState string

const (
	RefreshState_Unknown  RefreshState = "UNKNOWN"
	RefreshState_Pending  RefreshState = "PENDING"
	RefreshState_Running  RefreshState = "RUNNING"
	RefreshState_Success  RefreshState = "SUCCESS"
	RefreshState_Failed   RefreshState = "FAILED"
	RefreshState_Canceled RefreshState = "CANCELED"
)

// String representation for [fmt.Print].
func (f *RefreshState) String() string {
	return string(*f)
}

type RefreshTrigger string

const (
	RefreshTrigger_UnknownTrigger RefreshTrigger = "UNKNOWN_TRIGGER"
	RefreshTrigger_Schedule       RefreshTrigger = "SCHEDULE"
	RefreshTrigger_Manual         RefreshTrigger = "MANUAL"
)

// String representation for [fmt.Print].
func (f *RefreshTrigger) String() string {
	return string(*f)
}

type SchedulePauseStatus string

const (
	SchedulePauseStatus_Unspecified SchedulePauseStatus = "UNSPECIFIED"
	SchedulePauseStatus_Unpaused    SchedulePauseStatus = "UNPAUSED"
	SchedulePauseStatus_Paused      SchedulePauseStatus = "PAUSED"
)

// String representation for [fmt.Print].
func (f *SchedulePauseStatus) String() string {
	return string(*f)
}

type CancelRefresh struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
	RefreshId        *int64  `json:"refresh_id"`
}

type CancelRefresh_Response struct {
}

type CreateMonitor struct {
	FullTableNameArg         *string                     `json:"full_table_name_arg"`
	SkipBuiltinDashboard     *bool                       `json:"skip_builtin_dashboard"`
	WarehouseId              *string                     `json:"warehouse_id"`
	OutputSchemaName         *string                     `json:"output_schema_name"`
	AssetsDir                *string                     `json:"assets_dir"`
	InferenceLog             *InferenceLogAnalysisConfig `json:"inference_log"`
	TimeSeries               *TimeSeriesAnalysisConfig   `json:"time_series"`
	Snapshot                 *SnapshotAnalysisConfig     `json:"snapshot"`
	SlicingExprs             []string                    `json:"slicing_exprs"`
	CustomMetrics            []CustomMetric              `json:"custom_metrics"`
	BaselineTableName        *string                     `json:"baseline_table_name"`
	Schedule                 *MonitorCronSchedule        `json:"schedule"`
	Notifications            *Notifications              `json:"notifications"`
	DataClassificationConfig *DataClassificationConfig   `json:"data_classification_config"`
	TableName                *string                     `json:"table_name"`
	Status                   *MonitorStatus              `json:"status"`
	LatestMonitorFailureMsg  *string                     `json:"latest_monitor_failure_msg"`
	ProfileMetricsTableName  *string                     `json:"profile_metrics_table_name"`
	DriftMetricsTableName    *string                     `json:"drift_metrics_table_name"`
	DashboardId              *string                     `json:"dashboard_id"`
	MonitorVersion           *int64                      `json:"monitor_version"`
}

// Custom metric definition..
type CustomMetric struct {
	Name           *string           `json:"name"`
	Definition     *string           `json:"definition"`
	InputColumns   []string          `json:"input_columns"`
	OutputDataType *string           `json:"output_data_type"`
	Type           *CustomMetricType `json:"type"`
}

// Data classification related configuration..
type DataClassificationConfig struct {
	Enabled *bool `json:"enabled"`
}

type DataMonitorInfo struct {
	OutputSchemaName         *string                     `json:"output_schema_name"`
	AssetsDir                *string                     `json:"assets_dir"`
	InferenceLog             *InferenceLogAnalysisConfig `json:"inference_log"`
	TimeSeries               *TimeSeriesAnalysisConfig   `json:"time_series"`
	Snapshot                 *SnapshotAnalysisConfig     `json:"snapshot"`
	SlicingExprs             []string                    `json:"slicing_exprs"`
	CustomMetrics            []CustomMetric              `json:"custom_metrics"`
	BaselineTableName        *string                     `json:"baseline_table_name"`
	Schedule                 *MonitorCronSchedule        `json:"schedule"`
	Notifications            *Notifications              `json:"notifications"`
	DataClassificationConfig *DataClassificationConfig   `json:"data_classification_config"`
	TableName                *string                     `json:"table_name"`
	Status                   *MonitorStatus              `json:"status"`
	LatestMonitorFailureMsg  *string                     `json:"latest_monitor_failure_msg"`
	ProfileMetricsTableName  *string                     `json:"profile_metrics_table_name"`
	DriftMetricsTableName    *string                     `json:"drift_metrics_table_name"`
	DashboardId              *string                     `json:"dashboard_id"`
	MonitorVersion           *int64                      `json:"monitor_version"`
}

type DeleteMonitor struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
}

type DeleteMonitor_Response struct {
}

type Destination struct {
	EmailAddresses []string `json:"email_addresses"`
}

type GetMonitor struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
}

type GetRefresh struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
	RefreshId        *int64  `json:"refresh_id"`
}

type InferenceLogAnalysisConfig struct {
	ProblemType        *ProblemType `json:"problem_type"`
	TimestampCol       *string      `json:"timestamp_col"`
	Granularities      []string     `json:"granularities"`
	PredictionCol      *string      `json:"prediction_col"`
	LabelCol           *string      `json:"label_col"`
	ModelIdCol         *string      `json:"model_id_col"`
	PredictionProbaCol *string      `json:"prediction_proba_col"`
}

type ListRefreshes struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
}

type ListRefreshes_Response struct {
	Refreshes []RefreshInfo `json:"refreshes"`
}

type MonitorCronSchedule struct {
	QuartzCronExpression *string              `json:"quartz_cron_expression"`
	TimezoneId           *string              `json:"timezone_id"`
	PauseStatus          *SchedulePauseStatus `json:"pause_status"`
}

type Notifications struct {
	OnFailure                      *Destination `json:"on_failure"`
	OnNewClassificationTagDetected *Destination `json:"on_new_classification_tag_detected"`
}

type RefreshInfo struct {
	RefreshId   *int64          `json:"refresh_id"`
	State       *RefreshState   `json:"state"`
	Message     *string         `json:"message"`
	StartTimeMs *int64          `json:"start_time_ms"`
	EndTimeMs   *int64          `json:"end_time_ms"`
	Trigger     *RefreshTrigger `json:"trigger"`
}

type RegenerateDashboard struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
	WarehouseId      *string `json:"warehouse_id"`
}

type RegenerateDashboard_Response struct {
	DashboardId  *string `json:"dashboard_id"`
	ParentFolder *string `json:"parent_folder"`
}

type RunRefresh struct {
	FullTableNameArg *string `json:"full_table_name_arg"`
}

// Snapshot analysis configuration.
type SnapshotAnalysisConfig struct {
}

// Time series analysis configuration..
type TimeSeriesAnalysisConfig struct {
	TimestampCol  *string  `json:"timestamp_col"`
	Granularities []string `json:"granularities"`
}

type UpdateMonitor struct {
	FullTableNameArg         *string                     `json:"full_table_name_arg"`
	OutputSchemaName         *string                     `json:"output_schema_name"`
	AssetsDir                *string                     `json:"assets_dir"`
	InferenceLog             *InferenceLogAnalysisConfig `json:"inference_log"`
	TimeSeries               *TimeSeriesAnalysisConfig   `json:"time_series"`
	Snapshot                 *SnapshotAnalysisConfig     `json:"snapshot"`
	SlicingExprs             []string                    `json:"slicing_exprs"`
	CustomMetrics            []CustomMetric              `json:"custom_metrics"`
	BaselineTableName        *string                     `json:"baseline_table_name"`
	Schedule                 *MonitorCronSchedule        `json:"schedule"`
	Notifications            *Notifications              `json:"notifications"`
	DataClassificationConfig *DataClassificationConfig   `json:"data_classification_config"`
	TableName                *string                     `json:"table_name"`
	Status                   *MonitorStatus              `json:"status"`
	LatestMonitorFailureMsg  *string                     `json:"latest_monitor_failure_msg"`
	ProfileMetricsTableName  *string                     `json:"profile_metrics_table_name"`
	DriftMetricsTableName    *string                     `json:"drift_metrics_table_name"`
	DashboardId              *string                     `json:"dashboard_id"`
	MonitorVersion           *int64                      `json:"monitor_version"`
}
