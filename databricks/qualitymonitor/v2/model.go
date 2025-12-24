
package v2

import (
	
	"fmt"
	"io"
)



type AnomalyDetectionJobType string

const AnomalyDetectionJobTypeAnomalyDetectionJobTypeUnspecified AnomalyDetectionJobType = `ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED`
const AnomalyDetectionJobTypeAnomalyDetectionJobTypeNormal AnomalyDetectionJobType = `ANOMALY_DETECTION_JOB_TYPE_NORMAL`
const AnomalyDetectionJobTypeAnomalyDetectionJobTypeInternalHidden AnomalyDetectionJobType = `ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN`

// String representation for [fmt.Print]
func (f *AnomalyDetectionJobType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *AnomalyDetectionJobType) Set(v string) error {
	switch v {
	case `ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED`, `ANOMALY_DETECTION_JOB_TYPE_NORMAL`, `ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN`:
		*f = AnomalyDetectionJobType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ANOMALY_DETECTION_JOB_TYPE_UNSPECIFIED", "ANOMALY_DETECTION_JOB_TYPE_NORMAL", "ANOMALY_DETECTION_JOB_TYPE_INTERNAL_HIDDEN"`, v)
	}
}

// Values returns all possible values for AnomalyDetectionJobType.
//
// There is no guarantee on the order of the values in the slice.
func (f *AnomalyDetectionJobType) Values() []AnomalyDetectionJobType {
	return []AnomalyDetectionJobType{
		AnomalyDetectionJobTypeAnomalyDetectionJobTypeUnspecified,
		AnomalyDetectionJobTypeAnomalyDetectionJobTypeNormal,
		AnomalyDetectionJobTypeAnomalyDetectionJobTypeInternalHidden,
	}
}

// Type always returns AnomalyDetectionJobType to satisfy [pflag.Value] interface
func (f *AnomalyDetectionJobType) Type() string {
	return "AnomalyDetectionJobType"
}



type AnomalyDetectionRunStatus string

const AnomalyDetectionRunStatusAnomalyDetectionRunStatusUnknown AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_UNKNOWN`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusRunning AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_RUNNING`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusCanceled AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_CANCELED`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusPending AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_PENDING`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusSuccess AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_SUCCESS`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusFailed AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_FAILED`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusJobDeleted AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_JOB_DELETED`
const AnomalyDetectionRunStatusAnomalyDetectionRunStatusWorkspaceMismatchError AnomalyDetectionRunStatus = `ANOMALY_DETECTION_RUN_STATUS_WORKSPACE_MISMATCH_ERROR`

// String representation for [fmt.Print]
func (f *AnomalyDetectionRunStatus) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *AnomalyDetectionRunStatus) Set(v string) error {
	switch v {
	case `ANOMALY_DETECTION_RUN_STATUS_UNKNOWN`, `ANOMALY_DETECTION_RUN_STATUS_RUNNING`, `ANOMALY_DETECTION_RUN_STATUS_CANCELED`, `ANOMALY_DETECTION_RUN_STATUS_PENDING`, `ANOMALY_DETECTION_RUN_STATUS_SUCCESS`, `ANOMALY_DETECTION_RUN_STATUS_FAILED`, `ANOMALY_DETECTION_RUN_STATUS_JOB_DELETED`, `ANOMALY_DETECTION_RUN_STATUS_WORKSPACE_MISMATCH_ERROR`:
		*f = AnomalyDetectionRunStatus(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ANOMALY_DETECTION_RUN_STATUS_UNKNOWN", "ANOMALY_DETECTION_RUN_STATUS_RUNNING", "ANOMALY_DETECTION_RUN_STATUS_CANCELED", "ANOMALY_DETECTION_RUN_STATUS_PENDING", "ANOMALY_DETECTION_RUN_STATUS_SUCCESS", "ANOMALY_DETECTION_RUN_STATUS_FAILED", "ANOMALY_DETECTION_RUN_STATUS_JOB_DELETED", "ANOMALY_DETECTION_RUN_STATUS_WORKSPACE_MISMATCH_ERROR"`, v)
	}
}

// Values returns all possible values for AnomalyDetectionRunStatus.
//
// There is no guarantee on the order of the values in the slice.
func (f *AnomalyDetectionRunStatus) Values() []AnomalyDetectionRunStatus {
	return []AnomalyDetectionRunStatus{
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusUnknown,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusRunning,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusCanceled,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusPending,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusSuccess,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusFailed,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusJobDeleted,
		AnomalyDetectionRunStatusAnomalyDetectionRunStatusWorkspaceMismatchError,
	}
}

// Type always returns AnomalyDetectionRunStatus to satisfy [pflag.Value] interface
func (f *AnomalyDetectionRunStatus) Type() string {
	return "AnomalyDetectionRunStatus"
}



type ThresholdType string

const ThresholdTypeThresholdTypeUnspecified ThresholdType = `THRESHOLD_TYPE_UNSPECIFIED`
const ThresholdTypeThresholdTypeAuto ThresholdType = `THRESHOLD_TYPE_AUTO`
const ThresholdTypeThresholdTypeUnbounded ThresholdType = `THRESHOLD_TYPE_UNBOUNDED`
const ThresholdTypeThresholdTypeManual ThresholdType = `THRESHOLD_TYPE_MANUAL`

// String representation for [fmt.Print]
func (f *ThresholdType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *ThresholdType) Set(v string) error {
	switch v {
	case `THRESHOLD_TYPE_UNSPECIFIED`, `THRESHOLD_TYPE_AUTO`, `THRESHOLD_TYPE_UNBOUNDED`, `THRESHOLD_TYPE_MANUAL`:
		*f = ThresholdType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "THRESHOLD_TYPE_UNSPECIFIED", "THRESHOLD_TYPE_AUTO", "THRESHOLD_TYPE_UNBOUNDED", "THRESHOLD_TYPE_MANUAL"`, v)
	}
}

// Values returns all possible values for ThresholdType.
//
// There is no guarantee on the order of the values in the slice.
func (f *ThresholdType) Values() []ThresholdType {
	return []ThresholdType{
		ThresholdTypeThresholdTypeUnspecified,
		ThresholdTypeThresholdTypeAuto,
		ThresholdTypeThresholdTypeUnbounded,
		ThresholdTypeThresholdTypeManual,
	}
}

// Type always returns ThresholdType to satisfy [pflag.Value] interface
func (f *ThresholdType) Type() string {
	return "ThresholdType"
}





type AnomalyDetectionConfig struct {
	
	LastRunId *string `json:"lastRunId"`
	
	LatestRunStatus *AnomalyDetectionRunStatus `json:"latestRunStatus"`
	
	JobType *AnomalyDetectionJobType `json:"jobType"`
	
	ExcludedTableFullNames []string `json:"excludedTableFullNames"`
	
	CustomCheckConfiguration *CustomCheckConfiguration `json:"customCheckConfiguration"`
	
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


