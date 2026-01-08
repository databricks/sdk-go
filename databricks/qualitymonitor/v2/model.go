package v2

import (
	"fmt"
	"io"
)

type AnomalyDetectionJobType string

const (
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeUnspecified AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED"
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeNormal AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_NORMAL"
	AnomalyDetectionJobTypeAnomalyDetectionJobTypeInternalHidden AnomalyDetectionJobType = "ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN"
)

// String representation for [fmt.Print].
func (f *AnomalyDetectionJobType) String() string {
	return string(*f)
}

type AnomalyDetectionRunStatus string

const (
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusUnknown AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_UNKNOWN"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusRunning AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_RUNNING"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusCanceled AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_CANCELED"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusPending AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_PENDING"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusSuccess AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_SUCCESS"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusFailed AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_FAILED"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusJobDeleted AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_JOB_DELETED"
	AnomalyDetectionRunStatusAnomalyDetectionRunStatusWorkspaceMismatchError AnomalyDetectionRunStatus = "ANOMALY_DETECTION_RUN_STATUS_WORKSPACE_MISMATCH_ERROR"
)

// String representation for [fmt.Print].
func (f *AnomalyDetectionRunStatus) String() string {
	return string(*f)
}

type ThresholdType string

const (
	ThresholdTypeThresholdTypeUnspecified ThresholdType = "THRESHOLD_TYPE_UNSPECIFIED"
	ThresholdTypeThresholdTypeAuto ThresholdType = "THRESHOLD_TYPE_AUTO"
	ThresholdTypeThresholdTypeUnbounded ThresholdType = "THRESHOLD_TYPE_UNBOUNDED"
	ThresholdTypeThresholdTypeManual ThresholdType = "THRESHOLD_TYPE_MANUAL"
)

// String representation for [fmt.Print].
func (f *ThresholdType) String() string {
	return string(*f)
}

type AnomalyDetectionConfig struct {
	LastRunId *string `json:"last_run_id"`
	LatestRunStatus *AnomalyDetectionRunStatus `json:"latest_run_status"`
	JobType *AnomalyDetectionJobType `json:"job_type"`
	ExcludedTableFullNames []string `json:"excluded_table_full_names"`
	CustomCheckConfigurations []CustomCheckConfiguration `json:"custom_check_configurations"`
}

type ColumnMatcher struct {
	VariableName *string `json:"variable_name"`
	ColumnNames []string `json:"column_names"`
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
	CheckName *string `json:"check_name"`
	SqlQuery *string `json:"sql_query"`
	ColumnMatchers []ColumnMatcher `json:"column_matchers"`
	Thresholds *CustomCheckThresholds `json:"thresholds"`
}

type DeleteQualityMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId *string `json:"object_id"`
}

type GetQualityMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId *string `json:"object_id"`
}

type ListQualityMonitorRequest struct {
	PageToken *string `json:"page_token"`
	PageSize *int `json:"page_size"`
}

type ListQualityMonitorResponse struct {
	QualityMonitors []QualityMonitor `json:"quality_monitors"`
	NextPageToken *string `json:"next_page_token"`
}

type QualityMonitor struct {
	ObjectType *string `json:"object_type"`
	ObjectId *string `json:"object_id"`
	AnomalyDetectionConfig *AnomalyDetectionConfig `json:"anomaly_detection_config"`
}

type Threshold struct {
	BoundValue *int64 `json:"bound_value"`
	ThresholdType *ThresholdType `json:"threshold_type"`
}

type UpdateQualityMonitorRequest struct {
	ObjectType *string `json:"object_type"`
	ObjectId *string `json:"object_id"`
	QualityMonitor *QualityMonitor `json:"quality_monitor"`
}
