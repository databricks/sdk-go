package v2

import "time"

type Aggregation string

const (
	Aggregation_Sum           Aggregation = "SUM"
	Aggregation_Count         Aggregation = "COUNT"
	Aggregation_CountDistinct Aggregation = "COUNT_DISTINCT"
	Aggregation_Avg           Aggregation = "AVG"
	Aggregation_Median        Aggregation = "MEDIAN"
	Aggregation_Min           Aggregation = "MIN"
	Aggregation_Max           Aggregation = "MAX"
	Aggregation_Stddev        Aggregation = "STDDEV"
)

// String representation for [fmt.Print].
func (f *Aggregation) String() string {
	return string(*f)
}

type AlertEvaluationState string

const (
	AlertEvaluationState_Unknown   AlertEvaluationState = "UNKNOWN"
	AlertEvaluationState_Triggered AlertEvaluationState = "TRIGGERED"
	AlertEvaluationState_Ok        AlertEvaluationState = "OK"
	AlertEvaluationState_Error     AlertEvaluationState = "ERROR"
)

// String representation for [fmt.Print].
func (f *AlertEvaluationState) String() string {
	return string(*f)
}

type AlertLifecycleState string

const (
	AlertLifecycleState_Active  AlertLifecycleState = "ACTIVE"
	AlertLifecycleState_Deleted AlertLifecycleState = "DELETED"
)

// String representation for [fmt.Print].
func (f *AlertLifecycleState) String() string {
	return string(*f)
}

type ComparisonOperator string

const (
	ComparisonOperator_LessThan           ComparisonOperator = "LESS_THAN"
	ComparisonOperator_GreaterThan        ComparisonOperator = "GREATER_THAN"
	ComparisonOperator_Equal              ComparisonOperator = "EQUAL"
	ComparisonOperator_NotEqual           ComparisonOperator = "NOT_EQUAL"
	ComparisonOperator_GreaterThanOrEqual ComparisonOperator = "GREATER_THAN_OR_EQUAL"
	ComparisonOperator_LessThanOrEqual    ComparisonOperator = "LESS_THAN_OR_EQUAL"
	ComparisonOperator_IsNull             ComparisonOperator = "IS_NULL"
	ComparisonOperator_IsNotNull          ComparisonOperator = "IS_NOT_NULL"
)

// String representation for [fmt.Print].
func (f *ComparisonOperator) String() string {
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

type Alert struct {
	Id                  *string              `json:"id"`
	DisplayName         *string              `json:"display_name"`
	OwnerUserName       *string              `json:"owner_user_name"`
	CreateTime          *time.Time           `json:"create_time"`
	UpdateTime          *time.Time           `json:"update_time"`
	ParentPath          *string              `json:"parent_path"`
	QueryText           *string              `json:"query_text"`
	WarehouseId         *string              `json:"warehouse_id"`
	RunAsUserName       *string              `json:"run_as_user_name"`
	Evaluation          *AlertEvaluation     `json:"evaluation"`
	Schedule            *CronSchedule        `json:"schedule"`
	LifecycleState      *AlertLifecycleState `json:"lifecycle_state"`
	CustomSummary       *string              `json:"custom_summary"`
	CustomDescription   *string              `json:"custom_description"`
	RunAs               *AlertRunAs          `json:"run_as"`
	EffectiveRunAs      *AlertRunAs          `json:"effective_run_as"`
	EffectiveParentPath *string              `json:"effective_parent_path"`
}

type AlertEvaluation struct {
	Source             *AlertOperandColumn   `json:"source"`
	ComparisonOperator *ComparisonOperator   `json:"comparison_operator"`
	Threshold          *AlertOperand         `json:"threshold"`
	Notification       *AlertNotification    `json:"notification"`
	State              *AlertEvaluationState `json:"state"`
	LastEvaluatedAt    *time.Time            `json:"last_evaluated_at"`
	EmptyResultState   *AlertEvaluationState `json:"empty_result_state"`
}

type AlertNotification struct {
	Subscriptions    []AlertSubscription `json:"subscriptions"`
	RetriggerSeconds *int                `json:"retrigger_seconds"`
	NotifyOnOk       *bool               `json:"notify_on_ok"`
}

type AlertOperand struct {
	Column *AlertOperandColumn `json:"column"`
	Value  *AlertOperandValue  `json:"value"`
}

type AlertOperandColumn struct {
	Name        *string      `json:"name"`
	Display     *string      `json:"display"`
	Aggregation *Aggregation `json:"aggregation"`
}

type AlertOperandValue struct {
	StringValue *string  `json:"string_value"`
	DoubleValue *float64 `json:"double_value"`
	BoolValue   *bool    `json:"bool_value"`
}

type AlertRunAs struct {
	UserName             *string `json:"user_name"`
	ServicePrincipalName *string `json:"service_principal_name"`
}

type AlertSubscription struct {
	UserEmail     *string `json:"user_email"`
	DestinationId *string `json:"destination_id"`
}

type CreateAlertRequest struct {
	Alert *Alert `json:"alert"`
}

type CronSchedule struct {
	QuartzCronSchedule   *string              `json:"quartz_cron_schedule"`
	TimezoneId           *string              `json:"timezone_id"`
	PauseStatus          *SchedulePauseStatus `json:"pause_status"`
	EffectivePauseStatus *SchedulePauseStatus `json:"effective_pause_status"`
}

// Represents an empty message, similar to google.protobuf.Empty, which is not
// available in the firm right now..
type Empty struct {
}

type GetAlertRequest struct {
	Id *string `json:"id"`
}

type ListAlertsRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListAlertsResponse struct {
	Alerts        []Alert `json:"alerts"`
	NextPageToken *string `json:"next_page_token"`
}

type TrashAlertRequest struct {
	Id    *string `json:"id"`
	Purge *bool   `json:"purge"`
}

type UpdateAlertRequest struct {
	Alert      *Alert  `json:"alert"`
	UpdateMask *string `json:"update_mask"`
}
