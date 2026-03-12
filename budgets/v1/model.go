package v1

type ActionConfigurationType string

const (
	ActionConfigurationType_EmailNotification ActionConfigurationType = "EMAIL_NOTIFICATION"
)

// String representation for [fmt.Print].
func (f *ActionConfigurationType) String() string {
	return string(*f)
}

type AlertConfigurationQuantityType string

const (
	AlertConfigurationQuantityType_ListPriceDollarsUsd AlertConfigurationQuantityType = "LIST_PRICE_DOLLARS_USD"
)

// String representation for [fmt.Print].
func (f *AlertConfigurationQuantityType) String() string {
	return string(*f)
}

type AlertConfigurationTimePeriod string

const (
	AlertConfigurationTimePeriod_Month AlertConfigurationTimePeriod = "MONTH"
)

// String representation for [fmt.Print].
func (f *AlertConfigurationTimePeriod) String() string {
	return string(*f)
}

type AlertConfigurationTriggerType string

const (
	AlertConfigurationTriggerType_CumulativeSpendingExceeded AlertConfigurationTriggerType = "CUMULATIVE_SPENDING_EXCEEDED"
)

// String representation for [fmt.Print].
func (f *AlertConfigurationTriggerType) String() string {
	return string(*f)
}

type BudgetConfigurationFilter_Operator string

const (
	BudgetConfigurationFilter_Operator_In BudgetConfigurationFilter_Operator = "IN"
)

// String representation for [fmt.Print].
func (f *BudgetConfigurationFilter_Operator) String() string {
	return string(*f)
}

type ActionConfiguration struct {
	ActionConfigurationId *string                  `json:"action_configuration_id"`
	ActionType            *ActionConfigurationType `json:"action_type"`
	Target                *string                  `json:"target"`
}

type AlertConfiguration struct {
	AlertConfigurationId *string                         `json:"alert_configuration_id"`
	TimePeriod           *AlertConfigurationTimePeriod   `json:"time_period"`
	TriggerType          *AlertConfigurationTriggerType  `json:"trigger_type"`
	QuantityType         *AlertConfigurationQuantityType `json:"quantity_type"`
	QuantityThreshold    *string                         `json:"quantity_threshold"`
	ActionConfigurations []ActionConfiguration           `json:"action_configurations"`
}

type BudgetConfiguration struct {
	BudgetConfigurationId *string                    `json:"budget_configuration_id"`
	AccountId             *string                    `json:"account_id"`
	CreateTime            *int64                     `json:"create_time"`
	UpdateTime            *int64                     `json:"update_time"`
	AlertConfigurations   []AlertConfiguration       `json:"alert_configurations"`
	Filter                *BudgetConfigurationFilter `json:"filter"`
	DisplayName           *string                    `json:"display_name"`
}

type BudgetConfigurationFilter struct {
	WorkspaceId *BudgetConfigurationFilter_WorkspaceIdClause `json:"workspace_id"`
	Tags        []BudgetConfigurationFilter_TagClause        `json:"tags"`
}

type BudgetConfigurationFilter_Clause struct {
	Operator *BudgetConfigurationFilter_Operator `json:"operator"`
	Values   []string                            `json:"values"`
}

type BudgetConfigurationFilter_TagClause struct {
	Key   *string                           `json:"key"`
	Value *BudgetConfigurationFilter_Clause `json:"value"`
}

type BudgetConfigurationFilter_WorkspaceIdClause struct {
	Operator *BudgetConfigurationFilter_Operator `json:"operator"`
	Values   []int64                             `json:"values"`
}

type CreateBudgetConfiguration struct {
	Budget *CreateBudgetConfigurationBudget `json:"budget"`
}

type CreateBudgetConfiguration_Response struct {
	Budget *BudgetConfiguration `json:"budget"`
}

type CreateBudgetConfigurationBudget struct {
	BudgetConfigurationId *string                    `json:"budget_configuration_id"`
	AccountId             *string                    `json:"account_id"`
	CreateTime            *int64                     `json:"create_time"`
	UpdateTime            *int64                     `json:"update_time"`
	AlertConfigurations   []AlertConfiguration       `json:"alert_configurations"`
	Filter                *BudgetConfigurationFilter `json:"filter"`
	DisplayName           *string                    `json:"display_name"`
}

// * Delete budget.
type DeleteBudgetConfiguration struct {
	BudgetId  *string `json:"budget_id"`
	AccountId *string `json:"account_id"`
}

type DeleteBudgetConfiguration_Response struct {
}

type GetBudgetConfiguration struct {
	BudgetId           *string `json:"budget_id"`
	AccountId          *string `json:"account_id"`
	IncludeSpendStatus *bool   `json:"include_spend_status"`
}

type GetBudgetConfiguration_Response struct {
	Budget *BudgetConfiguration `json:"budget"`
}

type ListBudgetConfigurations struct {
	AccountId          *string `json:"account_id"`
	PageToken          *string `json:"page_token"`
	IncludeSpendStatus *bool   `json:"include_spend_status"`
}

type ListBudgetConfigurations_Response struct {
	Budgets       []BudgetConfiguration `json:"budgets"`
	NextPageToken *string               `json:"next_page_token"`
}

type UpdateBudgetConfiguration struct {
	BudgetId *string                          `json:"budget_id"`
	Budget   *UpdateBudgetConfigurationBudget `json:"budget"`
}

type UpdateBudgetConfiguration_Response struct {
	Budget *BudgetConfiguration `json:"budget"`
}

type UpdateBudgetConfigurationBudget struct {
	BudgetConfigurationId *string                    `json:"budget_configuration_id"`
	AccountId             *string                    `json:"account_id"`
	CreateTime            *int64                     `json:"create_time"`
	UpdateTime            *int64                     `json:"update_time"`
	AlertConfigurations   []AlertConfiguration       `json:"alert_configurations"`
	Filter                *BudgetConfigurationFilter `json:"filter"`
	DisplayName           *string                    `json:"display_name"`
}
