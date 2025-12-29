
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
	
	LastRunId *string `json:"lastRunId"`
	
	LatestRunStatus *AnomalyDetectionRunStatus `json:"latestRunStatus"`
	
	JobType *AnomalyDetectionJobType `json:"jobType"`
	
	ExcludedTableFullNames []string `json:"excludedTableFullNames"`
	
	CustomCheckConfigurations []CustomCheckConfiguration `json:"customCheckConfigurations"`
	
}



type ColumnMatcher struct {
	
	VariableName *string `json:"variableName"`
	
	ColumnNames []string `json:"columnNames"`
	
}



type CreateQualityMonitorRequest struct {
	
	QualityMonitor *QualityMonitor `json:"qualityMonitor"`
	
}



type CustomCheckConfiguration struct {
	
	ScalarCheck *CustomScalarCheck `json:"scalarCheck"`
	
}



type CustomCheckThresholds struct {
	
	LowerBound *Threshold `json:"lowerBound"`
	
	UpperBound *Threshold `json:"upperBound"`
	
}



type CustomScalarCheck struct {
	
	CheckName *string `json:"checkName"`
	
	SqlQuery *string `json:"sqlQuery"`
	
	ColumnMatchers []ColumnMatcher `json:"columnMatchers"`
	
	Thresholds *CustomCheckThresholds `json:"thresholds"`
	
}



type DeleteQualityMonitorRequest struct {
	
	ObjectType *string `json:"objectType"`
	
	ObjectId *string `json:"objectId"`
	
}



type GetQualityMonitorRequest struct {
	
	ObjectType *string `json:"objectType"`
	
	ObjectId *string `json:"objectId"`
	
}



type ListQualityMonitorRequest struct {
	
	PageToken *string `json:"pageToken"`
	
	PageSize *int `json:"pageSize"`
	
}



type ListQualityMonitorResponse struct {
	
	QualityMonitors []QualityMonitor `json:"qualityMonitors"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type QualityMonitor struct {
	
	ObjectType *string `json:"objectType"`
	
	ObjectId *string `json:"objectId"`
	
	AnomalyDetectionConfig *AnomalyDetectionConfig `json:"anomalyDetectionConfig"`
	
}



type Threshold struct {
	
	BoundValue *int64 `json:"boundValue"`
	
	ThresholdType *ThresholdType `json:"thresholdType"`
	
}



type UpdateQualityMonitorRequest struct {
	
	ObjectType *string `json:"objectType"`
	
	ObjectId *string `json:"objectId"`
	
	QualityMonitor *QualityMonitor `json:"qualityMonitor"`
	
}


