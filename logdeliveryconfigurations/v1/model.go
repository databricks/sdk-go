package v1

type LogDeliveryConfigStatus string

const (
	LogDeliveryConfigStatus_Enabled  LogDeliveryConfigStatus = "ENABLED"
	LogDeliveryConfigStatus_Disabled LogDeliveryConfigStatus = "DISABLED"
)

// String representation for [fmt.Print].
func (f *LogDeliveryConfigStatus) String() string {
	return string(*f)
}

type LogDeliveryOutputFormat string

const (
	LogDeliveryOutputFormat_Csv  LogDeliveryOutputFormat = "CSV"
	LogDeliveryOutputFormat_Json LogDeliveryOutputFormat = "JSON"
)

// String representation for [fmt.Print].
func (f *LogDeliveryOutputFormat) String() string {
	return string(*f)
}

type LogDeliveryStatusEnum string

const (
	LogDeliveryStatusEnum_Created       LogDeliveryStatusEnum = "CREATED"
	LogDeliveryStatusEnum_Succeeded     LogDeliveryStatusEnum = "SUCCEEDED"
	LogDeliveryStatusEnum_UserFailure   LogDeliveryStatusEnum = "USER_FAILURE"
	LogDeliveryStatusEnum_SystemFailure LogDeliveryStatusEnum = "SYSTEM_FAILURE"
	LogDeliveryStatusEnum_NotFound      LogDeliveryStatusEnum = "NOT_FOUND"
)

// String representation for [fmt.Print].
func (f *LogDeliveryStatusEnum) String() string {
	return string(*f)
}

type LogDeliveryType string

const (
	LogDeliveryType_BillableUsage LogDeliveryType = "BILLABLE_USAGE"
	LogDeliveryType_AuditLogs     LogDeliveryType = "AUDIT_LOGS"
)

// String representation for [fmt.Print].
func (f *LogDeliveryType) String() string {
	return string(*f)
}

// * Properties of the new log delivery configuration..
type CreateLogDeliveryConfiguration struct {
	LogDeliveryConfiguration *CreateLogDeliveryConfigurationParams `json:"log_delivery_configuration"`
}

type CreateLogDeliveryConfiguration_Response struct {
	LogDeliveryConfiguration *LogDeliveryConfiguration `json:"log_delivery_configuration"`
}

// * Log Delivery Configuration.
type CreateLogDeliveryConfigurationParams struct {
	ConfigId               *string                  `json:"config_id"`
	ConfigName             *string                  `json:"config_name"`
	LogType                *LogDeliveryType         `json:"log_type"`
	OutputFormat           *LogDeliveryOutputFormat `json:"output_format"`
	AccountId              *string                  `json:"account_id"`
	CredentialsId          *string                  `json:"credentials_id"`
	StorageConfigurationId *string                  `json:"storage_configuration_id"`
	WorkspaceIdsFilter     []int64                  `json:"workspace_ids_filter"`
	DeliveryPathPrefix     *string                  `json:"delivery_path_prefix"`
	DeliveryStartTime      *string                  `json:"delivery_start_time"`
	Status                 *LogDeliveryConfigStatus `json:"status"`
	CreationTime           *int64                   `json:"creation_time"`
	UpdateTime             *int64                   `json:"update_time"`
	LogDeliveryStatus      *LogDeliveryStatus       `json:"log_delivery_status"`
}

// * Get Log Delivery Configuration.
type GetLogDeliveryConfiguration struct {
	ConfigId  *string `json:"config_id"`
	AccountId *string `json:"account_id"`
}

type GetLogDeliveryConfiguration_Response struct {
	LogDeliveryConfiguration *LogDeliveryConfiguration `json:"log_delivery_configuration"`
}

// * List Log Delivery Configuration.
type ListLogDeliveryConfiguration struct {
	AccountId              *string                  `json:"account_id"`
	CredentialsId          *string                  `json:"credentials_id"`
	StorageConfigurationId *string                  `json:"storage_configuration_id"`
	Status                 *LogDeliveryConfigStatus `json:"status"`
	PageToken              *string                  `json:"page_token"`
}

type ListLogDeliveryConfiguration_Response struct {
	LogDeliveryConfigurations []LogDeliveryConfiguration `json:"log_delivery_configurations"`
	NextPageToken             *string                    `json:"next_page_token"`
}

// * Log Delivery Configuration.
type LogDeliveryConfiguration struct {
	ConfigId               *string                  `json:"config_id"`
	ConfigName             *string                  `json:"config_name"`
	LogType                *LogDeliveryType         `json:"log_type"`
	OutputFormat           *LogDeliveryOutputFormat `json:"output_format"`
	AccountId              *string                  `json:"account_id"`
	CredentialsId          *string                  `json:"credentials_id"`
	StorageConfigurationId *string                  `json:"storage_configuration_id"`
	WorkspaceIdsFilter     []int64                  `json:"workspace_ids_filter"`
	DeliveryPathPrefix     *string                  `json:"delivery_path_prefix"`
	DeliveryStartTime      *string                  `json:"delivery_start_time"`
	Status                 *LogDeliveryConfigStatus `json:"status"`
	CreationTime           *int64                   `json:"creation_time"`
	UpdateTime             *int64                   `json:"update_time"`
	LogDeliveryStatus      *LogDeliveryStatus       `json:"log_delivery_status"`
}

type LogDeliveryStatus struct {
	Status                    *LogDeliveryStatusEnum `json:"status"`
	LastAttemptTime           *string                `json:"last_attempt_time"`
	LastSuccessfulAttemptTime *string                `json:"last_successful_attempt_time"`
	Message                   *string                `json:"message"`
}

// * Update Log Delivery Configuration.
type UpdateLogDeliveryConfiguration struct {
	ConfigId  *string                  `json:"config_id"`
	AccountId *string                  `json:"account_id"`
	Status    *LogDeliveryConfigStatus `json:"status"`
}

type UpdateLogDeliveryConfiguration_Response struct {
}
