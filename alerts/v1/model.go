package v1

import "time"

type AlertOperator string

const (
	AlertOperator_GreaterThan        AlertOperator = "GREATER_THAN"
	AlertOperator_GreaterThanOrEqual AlertOperator = "GREATER_THAN_OR_EQUAL"
	AlertOperator_LessThan           AlertOperator = "LESS_THAN"
	AlertOperator_LessThanOrEqual    AlertOperator = "LESS_THAN_OR_EQUAL"
	AlertOperator_Equal              AlertOperator = "EQUAL"
	AlertOperator_NotEqual           AlertOperator = "NOT_EQUAL"
	AlertOperator_IsNull             AlertOperator = "IS_NULL"
)

// String representation for [fmt.Print].
func (f *AlertOperator) String() string {
	return string(*f)
}

type AlertState string

const (
	AlertState_Unknown   AlertState = "UNKNOWN"
	AlertState_Ok        AlertState = "OK"
	AlertState_Triggered AlertState = "TRIGGERED"
)

// String representation for [fmt.Print].
func (f *AlertState) String() string {
	return string(*f)
}

type LifecycleState string

const (
	LifecycleState_Active  LifecycleState = "ACTIVE"
	LifecycleState_Trashed LifecycleState = "TRASHED"
)

// String representation for [fmt.Print].
func (f *LifecycleState) String() string {
	return string(*f)
}

type Alert struct {
	Id                 *string         `json:"id"`
	DisplayName        *string         `json:"display_name"`
	QueryId            *string         `json:"query_id"`
	State              *AlertState     `json:"state"`
	SecondsToRetrigger *int            `json:"seconds_to_retrigger"`
	LifecycleState     *LifecycleState `json:"lifecycle_state"`
	TriggerTime        *time.Time      `json:"trigger_time"`
	CustomBody         *string         `json:"custom_body"`
	CustomSubject      *string         `json:"custom_subject"`
	Condition          *AlertCondition `json:"condition"`
	OwnerUserName      *string         `json:"owner_user_name"`
	ParentPath         *string         `json:"parent_path"`
	CreateTime         *time.Time      `json:"create_time"`
	UpdateTime         *time.Time      `json:"update_time"`
	NotifyOnOk         *bool           `json:"notify_on_ok"`
}

type AlertCondition struct {
	Op               *AlertOperator `json:"op"`
	Operand          *AlertOperand  `json:"operand"`
	Threshold        *AlertOperand  `json:"threshold"`
	EmptyResultState *AlertState    `json:"empty_result_state"`
}

type AlertOperand struct {
	Value  *AlertOperandValue  `json:"value"`
	Column *AlertOperandColumn `json:"column"`
}

type AlertOperandColumn struct {
	Name *string `json:"name"`
}

type AlertOperandValue struct {
	StringValue *string  `json:"string_value"`
	DoubleValue *float64 `json:"double_value"`
	BoolValue   *bool    `json:"bool_value"`
}

type CreateAlertRequest struct {
	Alert                  *CreateAlertRequestAlert `json:"alert"`
	AutoResolveDisplayName *bool                    `json:"auto_resolve_display_name"`
}

type CreateAlertRequestAlert struct {
	Id                 *string         `json:"id"`
	DisplayName        *string         `json:"display_name"`
	QueryId            *string         `json:"query_id"`
	State              *AlertState     `json:"state"`
	SecondsToRetrigger *int            `json:"seconds_to_retrigger"`
	LifecycleState     *LifecycleState `json:"lifecycle_state"`
	TriggerTime        *time.Time      `json:"trigger_time"`
	CustomBody         *string         `json:"custom_body"`
	CustomSubject      *string         `json:"custom_subject"`
	Condition          *AlertCondition `json:"condition"`
	OwnerUserName      *string         `json:"owner_user_name"`
	ParentPath         *string         `json:"parent_path"`
	CreateTime         *time.Time      `json:"create_time"`
	UpdateTime         *time.Time      `json:"update_time"`
	NotifyOnOk         *bool           `json:"notify_on_ok"`
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
	Results       []ListAlertsResponseAlert `json:"results"`
	NextPageToken *string                   `json:"next_page_token"`
}

type ListAlertsResponseAlert struct {
	Id                 *string         `json:"id"`
	DisplayName        *string         `json:"display_name"`
	QueryId            *string         `json:"query_id"`
	State              *AlertState     `json:"state"`
	SecondsToRetrigger *int            `json:"seconds_to_retrigger"`
	LifecycleState     *LifecycleState `json:"lifecycle_state"`
	TriggerTime        *time.Time      `json:"trigger_time"`
	CustomBody         *string         `json:"custom_body"`
	CustomSubject      *string         `json:"custom_subject"`
	Condition          *AlertCondition `json:"condition"`
	OwnerUserName      *string         `json:"owner_user_name"`
	ParentPath         *string         `json:"parent_path"`
	CreateTime         *time.Time      `json:"create_time"`
	UpdateTime         *time.Time      `json:"update_time"`
	NotifyOnOk         *bool           `json:"notify_on_ok"`
}

type TrashAlertRequest struct {
	Id *string `json:"id"`
}

type UpdateAlertRequest struct {
	Alert                  *UpdateAlertRequestAlert `json:"alert"`
	UpdateMask             *string                  `json:"update_mask"`
	Id                     *string                  `json:"id"`
	AutoResolveDisplayName *bool                    `json:"auto_resolve_display_name"`
}

type UpdateAlertRequestAlert struct {
	Id                 *string         `json:"id"`
	DisplayName        *string         `json:"display_name"`
	QueryId            *string         `json:"query_id"`
	State              *AlertState     `json:"state"`
	SecondsToRetrigger *int            `json:"seconds_to_retrigger"`
	LifecycleState     *LifecycleState `json:"lifecycle_state"`
	TriggerTime        *time.Time      `json:"trigger_time"`
	CustomBody         *string         `json:"custom_body"`
	CustomSubject      *string         `json:"custom_subject"`
	Condition          *AlertCondition `json:"condition"`
	OwnerUserName      *string         `json:"owner_user_name"`
	ParentPath         *string         `json:"parent_path"`
	CreateTime         *time.Time      `json:"create_time"`
	UpdateTime         *time.Time      `json:"update_time"`
	NotifyOnOk         *bool           `json:"notify_on_ok"`
}
