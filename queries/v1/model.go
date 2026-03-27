package v1

import "time"

type DatePrecision string

const (
	DatePrecision_DayPrecision    DatePrecision = "DAY_PRECISION"
	DatePrecision_MinutePrecision DatePrecision = "MINUTE_PRECISION"
	DatePrecision_SecondPrecision DatePrecision = "SECOND_PRECISION"
)

// String representation for [fmt.Print].
func (f *DatePrecision) String() string {
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

type RunAsMode string

const (
	RunAsMode_Owner  RunAsMode = "OWNER"
	RunAsMode_Viewer RunAsMode = "VIEWER"
)

// String representation for [fmt.Print].
func (f *RunAsMode) String() string {
	return string(*f)
}

type DateRangeValue_DynamicDateRange string

const (
	DateRangeValue_DynamicDateRange_Today        DateRangeValue_DynamicDateRange = "TODAY"
	DateRangeValue_DynamicDateRange_Yesterday    DateRangeValue_DynamicDateRange = "YESTERDAY"
	DateRangeValue_DynamicDateRange_ThisWeek     DateRangeValue_DynamicDateRange = "THIS_WEEK"
	DateRangeValue_DynamicDateRange_ThisMonth    DateRangeValue_DynamicDateRange = "THIS_MONTH"
	DateRangeValue_DynamicDateRange_ThisYear     DateRangeValue_DynamicDateRange = "THIS_YEAR"
	DateRangeValue_DynamicDateRange_LastWeek     DateRangeValue_DynamicDateRange = "LAST_WEEK"
	DateRangeValue_DynamicDateRange_LastMonth    DateRangeValue_DynamicDateRange = "LAST_MONTH"
	DateRangeValue_DynamicDateRange_LastYear     DateRangeValue_DynamicDateRange = "LAST_YEAR"
	DateRangeValue_DynamicDateRange_LastHour     DateRangeValue_DynamicDateRange = "LAST_HOUR"
	DateRangeValue_DynamicDateRange_Last8Hours   DateRangeValue_DynamicDateRange = "LAST_8_HOURS"
	DateRangeValue_DynamicDateRange_Last24Hours  DateRangeValue_DynamicDateRange = "LAST_24_HOURS"
	DateRangeValue_DynamicDateRange_Last7Days    DateRangeValue_DynamicDateRange = "LAST_7_DAYS"
	DateRangeValue_DynamicDateRange_Last14Days   DateRangeValue_DynamicDateRange = "LAST_14_DAYS"
	DateRangeValue_DynamicDateRange_Last30Days   DateRangeValue_DynamicDateRange = "LAST_30_DAYS"
	DateRangeValue_DynamicDateRange_Last60Days   DateRangeValue_DynamicDateRange = "LAST_60_DAYS"
	DateRangeValue_DynamicDateRange_Last90Days   DateRangeValue_DynamicDateRange = "LAST_90_DAYS"
	DateRangeValue_DynamicDateRange_Last12Months DateRangeValue_DynamicDateRange = "LAST_12_MONTHS"
)

// String representation for [fmt.Print].
func (f *DateRangeValue_DynamicDateRange) String() string {
	return string(*f)
}

type DateValue_DynamicDate string

const (
	DateValue_DynamicDate_Now       DateValue_DynamicDate = "NOW"
	DateValue_DynamicDate_Yesterday DateValue_DynamicDate = "YESTERDAY"
)

// String representation for [fmt.Print].
func (f *DateValue_DynamicDate) String() string {
	return string(*f)
}

type CreateQueryRequest struct {
	Query                  *CreateQueryRequestQuery `json:"query"`
	AutoResolveDisplayName *bool                    `json:"auto_resolve_display_name"`
}

type CreateQueryRequestQuery struct {
	Id                   *string          `json:"id"`
	DisplayName          *string          `json:"display_name"`
	Description          *string          `json:"description"`
	OwnerUserName        *string          `json:"owner_user_name"`
	WarehouseId          *string          `json:"warehouse_id"`
	QueryText            *string          `json:"query_text"`
	RunAsMode            *RunAsMode       `json:"run_as_mode"`
	LifecycleState       *LifecycleState  `json:"lifecycle_state"`
	LastModifierUserName *string          `json:"last_modifier_user_name"`
	ParentPath           *string          `json:"parent_path"`
	Tags                 []string         `json:"tags"`
	CreateTime           *time.Time       `json:"create_time"`
	UpdateTime           *time.Time       `json:"update_time"`
	Parameters           []QueryParameter `json:"parameters"`
	ApplyAutoLimit       *bool            `json:"apply_auto_limit"`
	Catalog              *string          `json:"catalog"`
	Schema               *string          `json:"schema"`
}

type DateRange struct {
	Start *string `json:"start"`
	End   *string `json:"end"`
}

type DateRangeValue struct {
	DynamicDateRangeValue *DateRangeValue_DynamicDateRange `json:"dynamic_date_range_value"`
	DateRangeValue        *DateRange                       `json:"date_range_value"`
	Precision             *DatePrecision                   `json:"precision"`
	StartDayOfWeek        *int                             `json:"start_day_of_week"`
}

type DateValue struct {
	DynamicDateValue *DateValue_DynamicDate `json:"dynamic_date_value"`
	DateValue        *string                `json:"date_value"`
	Precision        *DatePrecision         `json:"precision"`
}

// Represents an empty message, similar to google.protobuf.Empty, which is not
// available in the firm right now..
type Empty struct {
}

type EnumValue struct {
	Values             []string            `json:"values"`
	EnumOptions        *string             `json:"enum_options"`
	MultiValuesOptions *MultiValuesOptions `json:"multi_values_options"`
}

type GetQueryRequest struct {
	Id *string `json:"id"`
}

type ListQueriesRequest struct {
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListQueriesResponse struct {
	Results       []ListQueryObjectsResponseQuery `json:"results"`
	NextPageToken *string                         `json:"next_page_token"`
}

type ListQueryObjectsResponseQuery struct {
	Id                   *string          `json:"id"`
	DisplayName          *string          `json:"display_name"`
	Description          *string          `json:"description"`
	OwnerUserName        *string          `json:"owner_user_name"`
	WarehouseId          *string          `json:"warehouse_id"`
	QueryText            *string          `json:"query_text"`
	RunAsMode            *RunAsMode       `json:"run_as_mode"`
	LifecycleState       *LifecycleState  `json:"lifecycle_state"`
	LastModifierUserName *string          `json:"last_modifier_user_name"`
	ParentPath           *string          `json:"parent_path"`
	Tags                 []string         `json:"tags"`
	CreateTime           *time.Time       `json:"create_time"`
	UpdateTime           *time.Time       `json:"update_time"`
	Parameters           []QueryParameter `json:"parameters"`
	ApplyAutoLimit       *bool            `json:"apply_auto_limit"`
	Catalog              *string          `json:"catalog"`
	Schema               *string          `json:"schema"`
}

type ListVisualizationsForQueryRequest struct {
	Id        *string `json:"id"`
	PageToken *string `json:"page_token"`
	PageSize  *int    `json:"page_size"`
}

type ListVisualizationsForQueryResponse struct {
	Results       []Visualization `json:"results"`
	NextPageToken *string         `json:"next_page_token"`
}

type MultiValuesOptions struct {
	Prefix    *string `json:"prefix"`
	Separator *string `json:"separator"`
	Suffix    *string `json:"suffix"`
}

type NumericValue struct {
	Value *float64 `json:"value"`
}

type Query struct {
	Id                   *string          `json:"id"`
	DisplayName          *string          `json:"display_name"`
	Description          *string          `json:"description"`
	OwnerUserName        *string          `json:"owner_user_name"`
	WarehouseId          *string          `json:"warehouse_id"`
	QueryText            *string          `json:"query_text"`
	RunAsMode            *RunAsMode       `json:"run_as_mode"`
	LifecycleState       *LifecycleState  `json:"lifecycle_state"`
	LastModifierUserName *string          `json:"last_modifier_user_name"`
	ParentPath           *string          `json:"parent_path"`
	Tags                 []string         `json:"tags"`
	CreateTime           *time.Time       `json:"create_time"`
	UpdateTime           *time.Time       `json:"update_time"`
	Parameters           []QueryParameter `json:"parameters"`
	ApplyAutoLimit       *bool            `json:"apply_auto_limit"`
	Catalog              *string          `json:"catalog"`
	Schema               *string          `json:"schema"`
}

type QueryBackedValue struct {
	Values             []string            `json:"values"`
	QueryId            *string             `json:"query_id"`
	MultiValuesOptions *MultiValuesOptions `json:"multi_values_options"`
}

type QueryParameter struct {
	Title            *string           `json:"title"`
	Name             *string           `json:"name"`
	TextValue        *TextValue        `json:"text_value"`
	NumericValue     *NumericValue     `json:"numeric_value"`
	EnumValue        *EnumValue        `json:"enum_value"`
	DateValue        *DateValue        `json:"date_value"`
	DateRangeValue   *DateRangeValue   `json:"date_range_value"`
	QueryBackedValue *QueryBackedValue `json:"query_backed_value"`
}

type TextValue struct {
	Value *string `json:"value"`
}

type TrashQueryRequest struct {
	Id *string `json:"id"`
}

type UpdateQueryRequest struct {
	Query                  *UpdateQueryRequestQuery `json:"query"`
	UpdateMask             *string                  `json:"update_mask"`
	Id                     *string                  `json:"id"`
	AutoResolveDisplayName *bool                    `json:"auto_resolve_display_name"`
}

type UpdateQueryRequestQuery struct {
	Id                   *string          `json:"id"`
	DisplayName          *string          `json:"display_name"`
	Description          *string          `json:"description"`
	OwnerUserName        *string          `json:"owner_user_name"`
	WarehouseId          *string          `json:"warehouse_id"`
	QueryText            *string          `json:"query_text"`
	RunAsMode            *RunAsMode       `json:"run_as_mode"`
	LifecycleState       *LifecycleState  `json:"lifecycle_state"`
	LastModifierUserName *string          `json:"last_modifier_user_name"`
	ParentPath           *string          `json:"parent_path"`
	Tags                 []string         `json:"tags"`
	CreateTime           *time.Time       `json:"create_time"`
	UpdateTime           *time.Time       `json:"update_time"`
	Parameters           []QueryParameter `json:"parameters"`
	ApplyAutoLimit       *bool            `json:"apply_auto_limit"`
	Catalog              *string          `json:"catalog"`
	Schema               *string          `json:"schema"`
}

type Visualization struct {
	Id                  *string    `json:"id"`
	DisplayName         *string    `json:"display_name"`
	Type                *string    `json:"type"`
	CreateTime          *time.Time `json:"create_time"`
	UpdateTime          *time.Time `json:"update_time"`
	SerializedQueryPlan *string    `json:"serialized_query_plan"`
	SerializedOptions   *string    `json:"serialized_options"`
	QueryId             *string    `json:"query_id"`
}
