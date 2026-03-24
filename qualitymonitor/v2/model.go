package v2

type AnomalyDetectionJobType string

const (
	AnomalyDetectionJobType_AnomalyDetectionJobTypeUnspecified    AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED"
	AnomalyDetectionJobType_AnomalyDetectionJobTypeNormal         AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_NORMAL"
	AnomalyDetectionJobType_AnomalyDetectionJobTypeInternalHidden AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN"
)

// String representation for [fmt.Print].
func (f *AnomalyDetectionJobType) String() string {
	return string(*f)
}

type AnomalyDetectionRunStatus string

const (
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusUnknown                AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_UNKNOWN"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusRunning                AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_RUNNING"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusCanceled               AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_CANCELED"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusPending                AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_PENDING"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusSuccess                AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_SUCCESS"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusFailed                 AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_FAILED"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusJobDeleted             AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_JOB_DELETED"
	AnomalyDetectionRunStatus_AnomalyDetectionRunStatusWorkspaceMismatchError AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_WORKSPACE_MISMATCH_ERROR"
)

// String representation for [fmt.Print].
func (f *AnomalyDetectionRunStatus) String() string {
	return string(*f)
}

type ThresholdType string

const (
	ThresholdType_ThresholdTypeUnspecified ThresholdType = "THRESHOLD_TYPE_UNSPECIFIED"
	ThresholdType_ThresholdTypeAuto        ThresholdType = "THRESHOLD_TYPE_AUTO"
	ThresholdType_ThresholdTypeUnbounded   ThresholdType = "THRESHOLD_TYPE_UNBOUNDED"
	ThresholdType_ThresholdTypeManual      ThresholdType = "THRESHOLD_TYPE_MANUAL"
)

// String representation for [fmt.Print].
func (f *ThresholdType) String() string {
	return string(*f)
}

type AnomalyDetectionConfig struct {
	LastRunId                   *string                      `json:"last_run_id"`
	LatestRunStatus             *AnomalyDetectionRunStatus   `json:"latest_run_status"`
	JobType                     *AnomalyDetectionJobType     `json:"job_type"`
	ExcludedTableFullNames      []string                     `json:"excluded_table_full_names"`
	CustomCheckConfigurations   []CustomCheckConfiguration   `json:"custom_check_configurations"`
	ValidityCheckConfigurations []ValidityCheckConfiguration `json:"validity_check_configurations"`
}

type ColumnMatcher struct {
	VariableName *string  `json:"variable_name"`
	ColumnNames  []string `json:"column_names"`
}

type CreateQualityMonitorRequest struct {
	QualityMonitor *QualityMonitor `json:"quality_monitor"`
}

type CustomCheckConfiguration struct {
	ScalarCheck *CustomScalarCheck `json:"scalar_check"`
}

type CustomCheckThresholds struct {
	LowerBound *Threshold `json:"lower_bound"`
	UpperBound *Threshold `json:"upper_bound"`
}

type CustomScalarCheck struct {
	CheckName      *string                `json:"check_name"`
	SqlQuery       *string                `json:"sql_query"`
	ColumnMatchers []ColumnMatcher        `json:"column_matchers"`
	Thresholds     *CustomCheckThresholds `json:"thresholds"`
}

type DeleteQualityMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
}

type GetQualityMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId   *string `json:"object_id"`
}

type ListQualityMonitorRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListQualityMonitorResponse struct {
	QualityMonitors []QualityMonitor `json:"quality_monitors"`
	NextPageToken   *string          `json:"next_page_token"`
}

type PercentNullValidityCheck struct {
	ColumnNames []string `json:"column_names"`
	UpperBound  *float64 `json:"upper_bound"`
}

type QualityMonitor struct {
	ObjectType                  *string                      `json:"object_type"`
	ObjectId                    *string                      `json:"object_id"`
	AnomalyDetectionConfig      *AnomalyDetectionConfig      `json:"anomaly_detection_config"`
	ValidityCheckConfigurations []ValidityCheckConfiguration `json:"validity_check_configurations"`
}

type RangeValidityCheck struct {
	ColumnNames []string `json:"column_names"`
	LowerBound  *float64 `json:"lower_bound"`
	UpperBound  *float64 `json:"upper_bound"`
}

type Threshold struct {
	BoundValue    *int64         `json:"bound_value"`
	ThresholdType *ThresholdType `json:"threshold_type"`
}

type UniquenessValidityCheck struct {
	ColumnNames []string `json:"column_names"`
}

type UpdateQualityMonitorRequest struct {
	ObjectType     *string         `json:"object_type"`
	ObjectId       *string         `json:"object_id"`
	QualityMonitor *QualityMonitor `json:"quality_monitor"`
}

type ValidityCheckConfiguration struct {
	Name                     *string                   `json:"name"`
	PercentNullValidityCheck *PercentNullValidityCheck `json:"percent_null_validity_check"`
	RangeValidityCheck       *RangeValidityCheck       `json:"range_validity_check"`
	UniquenessValidityCheck  *UniquenessValidityCheck  `json:"uniqueness_validity_check"`
}
