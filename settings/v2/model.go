package v2

type AccessPolicyType string

const (
	AccessPolicyTypeAccessPolicyTypeUnspecified AccessPolicyType = "ACCESS_POLICY_TYPE_UNSPECIFIED"
	AccessPolicyTypeAllowAllDomains             AccessPolicyType = "ALLOW_ALL_DOMAINS"
	AccessPolicyTypeAllowApprovedDomains        AccessPolicyType = "ALLOW_APPROVED_DOMAINS"
	AccessPolicyTypeDenyAllDomains              AccessPolicyType = "DENY_ALL_DOMAINS"
)

// String representation for [fmt.Print].
func (f *AccessPolicyType) String() string {
	return string(*f)
}

type DayOfWeek string

const (
	DayOfWeekDayOfWeekUnspecified DayOfWeek = "DAY_OF_WEEK_UNSPECIFIED"
	DayOfWeekMonday               DayOfWeek = "MONDAY"
	DayOfWeekTuesday              DayOfWeek = "TUESDAY"
	DayOfWeekWednesday            DayOfWeek = "WEDNESDAY"
	DayOfWeekThursday             DayOfWeek = "THURSDAY"
	DayOfWeekFriday               DayOfWeek = "FRIDAY"
	DayOfWeekSaturday             DayOfWeek = "SATURDAY"
	DayOfWeekSunday               DayOfWeek = "SUNDAY"
)

// String representation for [fmt.Print].
func (f *DayOfWeek) String() string {
	return string(*f)
}

type PersonalComputeMessageEnum string

const (
	PersonalComputeMessageEnumPersonalComputeMessageEnumUnspecified PersonalComputeMessageEnum = "PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED"
	PersonalComputeMessageEnumOn                                    PersonalComputeMessageEnum = "ON"
	PersonalComputeMessageEnumDelegate                              PersonalComputeMessageEnum = "DELEGATE"
)

// String representation for [fmt.Print].
func (f *PersonalComputeMessageEnum) String() string {
	return string(*f)
}

type Status string

const (
	StatusStatusUnspecified         Status = "STATUS_UNSPECIFIED"
	StatusAllowAll                  Status = "ALLOW_ALL"
	StatusRestrictTokensAndJobRunAs Status = "RESTRICT_TOKENS_AND_JOB_RUN_AS"
)

// String representation for [fmt.Print].
func (f *Status) String() string {
	return string(*f)
}

type WeekDayFrequency string

const (
	WeekDayFrequencyWeekDayFrequencyUnspecified WeekDayFrequency = "WEEK_DAY_FREQUENCY_UNSPECIFIED"
	WeekDayFrequencyFirstOfMonth                WeekDayFrequency = "FIRST_OF_MONTH"
	WeekDayFrequencySecondOfMonth               WeekDayFrequency = "SECOND_OF_MONTH"
	WeekDayFrequencyThirdOfMonth                WeekDayFrequency = "THIRD_OF_MONTH"
	WeekDayFrequencyFourthOfMonth               WeekDayFrequency = "FOURTH_OF_MONTH"
	WeekDayFrequencyFirstAndThirdOfMonth        WeekDayFrequency = "FIRST_AND_THIRD_OF_MONTH"
	WeekDayFrequencySecondAndFourthOfMonth      WeekDayFrequency = "SECOND_AND_FOURTH_OF_MONTH"
	WeekDayFrequencyEveryWeek                   WeekDayFrequency = "EVERY_WEEK"
)

// String representation for [fmt.Print].
func (f *WeekDayFrequency) String() string {
	return string(*f)
}

type AibiDashboardEmbeddingAccessPolicy struct {
	AccessPolicyType *AccessPolicyType `json:"access_policy_type"`
}

type AibiDashboardEmbeddingApprovedDomains struct {
	ApprovedDomains []string `json:"approved_domains"`
}

type BooleanMessage struct {
	Value *bool `json:"value"`
}

type ClusterAutoRestartMessage struct {
	Enabled                         *bool              `json:"enabled"`
	CanToggle                       *bool              `json:"can_toggle"`
	MaintenanceWindow               *MaintenanceWindow `json:"maintenance_window"`
	EnablementDetails               *EnablementDetails `json:"enablement_details"`
	RestartEvenIfNoUpdatesAvailable *bool              `json:"restart_even_if_no_updates_available"`
}

type EnablementDetails struct {
	UnavailableForNonEnterpriseTier   *bool `json:"unavailable_for_non_enterprise_tier"`
	UnavailableForDisabledEntitlement *bool `json:"unavailable_for_disabled_entitlement"`
	ForcedForComplianceMode           *bool `json:"forced_for_compliance_mode"`
}

type GetPublicAccountSettingRequest struct {
	AccountId *string `json:"account_id"`
	Name      *string `json:"name"`
}

type GetPublicAccountUserPreferenceRequest struct {
	AccountId *string `json:"account_id"`
	UserId    *string `json:"user_id"`
	Name      *string `json:"name"`
}

type GetPublicWorkspaceSettingRequest struct {
	Name *string `json:"name"`
}

type IntegerMessage struct {
	Value *int `json:"value"`
}

type ListAccountSettingsMetadataRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListAccountSettingsMetadataResponse struct {
	SettingsMetadata []SettingsMetadata `json:"settings_metadata"`
	NextPageToken    *string            `json:"next_page_token"`
}

type ListAccountUserPreferencesMetadataRequest struct {
	AccountId *string `json:"account_id"`
	UserId    *string `json:"user_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListAccountUserPreferencesMetadataResponse struct {
	SettingsMetadata []SettingsMetadata `json:"settings_metadata"`
	NextPageToken    *string            `json:"next_page_token"`
}

type ListWorkspaceSettingsMetadataRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListWorkspaceSettingsMetadataResponse struct {
	SettingsMetadata []SettingsMetadata `json:"settings_metadata"`
	NextPageToken    *string            `json:"next_page_token"`
}

type MaintenanceWindow struct {
	WeekDayBasedSchedule *WeekDayBasedSchedule `json:"week_day_based_schedule"`
}

type PatchPublicAccountSettingRequest struct {
	AccountId *string  `json:"account_id"`
	Name      *string  `json:"name"`
	Setting   *Setting `json:"setting"`
}

type PatchPublicAccountUserPreferenceRequest struct {
	AccountId *string         `json:"account_id"`
	UserId    *string         `json:"user_id"`
	Name      *string         `json:"name"`
	Setting   *UserPreference `json:"setting"`
}

type PatchPublicWorkspaceSettingRequest struct {
	Name    *string  `json:"name"`
	Setting *Setting `json:"setting"`
}

type PersonalComputeMessage struct {
	Value *PersonalComputeMessageEnum `json:"value"`
}

type RestrictWorkspaceAdminsMessage struct {
	Status *Status `json:"status"`
}

type Setting struct {
	Name                                           *string                                `json:"name"`
	BooleanVal                                     *BooleanMessage                        `json:"boolean_val"`
	StringVal                                      *StringMessage                         `json:"string_val"`
	IntegerVal                                     *IntegerMessage                        `json:"integer_val"`
	AutomaticClusterUpdateWorkspace                *ClusterAutoRestartMessage             `json:"automatic_cluster_update_workspace"`
	AibiDashboardEmbeddingApprovedDomains          *AibiDashboardEmbeddingApprovedDomains `json:"aibi_dashboard_embedding_approved_domains"`
	AibiDashboardEmbeddingAccessPolicy             *AibiDashboardEmbeddingAccessPolicy    `json:"aibi_dashboard_embedding_access_policy"`
	RestrictWorkspaceAdmins                        *RestrictWorkspaceAdminsMessage        `json:"restrict_workspace_admins"`
	PersonalCompute                                *PersonalComputeMessage                `json:"personal_compute"`
	EffectiveBooleanVal                            *BooleanMessage                        `json:"effective_boolean_val"`
	EffectiveStringVal                             *StringMessage                         `json:"effective_string_val"`
	EffectiveIntegerVal                            *IntegerMessage                        `json:"effective_integer_val"`
	EffectiveAutomaticClusterUpdateWorkspace       *ClusterAutoRestartMessage             `json:"effective_automatic_cluster_update_workspace"`
	EffectiveAibiDashboardEmbeddingApprovedDomains *AibiDashboardEmbeddingApprovedDomains `json:"effective_aibi_dashboard_embedding_approved_domains"`
	EffectiveAibiDashboardEmbeddingAccessPolicy    *AibiDashboardEmbeddingAccessPolicy    `json:"effective_aibi_dashboard_embedding_access_policy"`
	EffectiveRestrictWorkspaceAdmins               *RestrictWorkspaceAdminsMessage        `json:"effective_restrict_workspace_admins"`
	EffectivePersonalCompute                       *PersonalComputeMessage                `json:"effective_personal_compute"`
}

type SettingsMetadata struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Type        *string `json:"type"`
	DocsLink    *string `json:"docs_link"`
}

type StringMessage struct {
	Value *string `json:"value"`
}

type UserPreference struct {
	Name                *string         `json:"name"`
	UserId              *string         `json:"user_id"`
	BooleanVal          *BooleanMessage `json:"boolean_val"`
	StringVal           *StringMessage  `json:"string_val"`
	EffectiveBooleanVal *BooleanMessage `json:"effective_boolean_val"`
	EffectiveStringVal  *StringMessage  `json:"effective_string_val"`
}

type WeekDayBasedSchedule struct {
	Frequency       *WeekDayFrequency `json:"frequency"`
	DayOfWeek       *DayOfWeek        `json:"day_of_week"`
	WindowStartTime *WindowStartTime  `json:"window_start_time"`
}

type WindowStartTime struct {
	Hours   *int `json:"hours"`
	Minutes *int `json:"minutes"`
}
