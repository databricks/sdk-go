package v2

type AibiDashboardEmbeddingAccessPolicy_AccessPolicyType string

const (
	AibiDashboardEmbeddingAccessPolicy_AccessPolicyType_AccessPolicyTypeUnspecified AibiDashboardEmbeddingAccessPolicy_AccessPolicyType = "ACCESS_POLICY_TYPE_UNSPECIFIED"
	AibiDashboardEmbeddingAccessPolicy_AccessPolicyType_AllowAllDomains             AibiDashboardEmbeddingAccessPolicy_AccessPolicyType = "ALLOW_ALL_DOMAINS"
	AibiDashboardEmbeddingAccessPolicy_AccessPolicyType_AllowApprovedDomains        AibiDashboardEmbeddingAccessPolicy_AccessPolicyType = "ALLOW_APPROVED_DOMAINS"
	AibiDashboardEmbeddingAccessPolicy_AccessPolicyType_DenyAllDomains              AibiDashboardEmbeddingAccessPolicy_AccessPolicyType = "DENY_ALL_DOMAINS"
)

// String representation for [fmt.Print].
func (f *AibiDashboardEmbeddingAccessPolicy_AccessPolicyType) String() string {
	return string(*f)
}

type ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek string

const (
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_DayOfWeekUnspecified ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "DAY_OF_WEEK_UNSPECIFIED"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Monday               ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "MONDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Tuesday              ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "TUESDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Wednesday            ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "WEDNESDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Thursday             ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "THURSDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Friday               ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "FRIDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Saturday             ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "SATURDAY"
	ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek_Sunday               ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek = "SUNDAY"
)

// String representation for [fmt.Print].
func (f *ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek) String() string {
	return string(*f)
}

type ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency string

const (
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_WeekDayFrequencyUnspecified ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "WEEK_DAY_FREQUENCY_UNSPECIFIED"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_FirstOfMonth                ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "FIRST_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_SecondOfMonth               ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "SECOND_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_ThirdOfMonth                ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "THIRD_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_FourthOfMonth               ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "FOURTH_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_FirstAndThirdOfMonth        ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "FIRST_AND_THIRD_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_SecondAndFourthOfMonth      ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "SECOND_AND_FOURTH_OF_MONTH"
	ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency_EveryWeek                   ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency = "EVERY_WEEK"
)

// String representation for [fmt.Print].
func (f *ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency) String() string {
	return string(*f)
}

type PersonalComputeMessage_PersonalComputeMessageEnum string

const (
	PersonalComputeMessage_PersonalComputeMessageEnum_PersonalComputeMessageEnumUnspecified PersonalComputeMessage_PersonalComputeMessageEnum = "PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED"
	PersonalComputeMessage_PersonalComputeMessageEnum_On                                    PersonalComputeMessage_PersonalComputeMessageEnum = "ON"
	PersonalComputeMessage_PersonalComputeMessageEnum_Delegate                              PersonalComputeMessage_PersonalComputeMessageEnum = "DELEGATE"
)

// String representation for [fmt.Print].
func (f *PersonalComputeMessage_PersonalComputeMessageEnum) String() string {
	return string(*f)
}

type RestrictWorkspaceAdminsMessage_Status string

const (
	RestrictWorkspaceAdminsMessage_Status_StatusUnspecified         RestrictWorkspaceAdminsMessage_Status = "STATUS_UNSPECIFIED"
	RestrictWorkspaceAdminsMessage_Status_AllowAll                  RestrictWorkspaceAdminsMessage_Status = "ALLOW_ALL"
	RestrictWorkspaceAdminsMessage_Status_RestrictTokensAndJobRunAs RestrictWorkspaceAdminsMessage_Status = "RESTRICT_TOKENS_AND_JOB_RUN_AS"
)

// String representation for [fmt.Print].
func (f *RestrictWorkspaceAdminsMessage_Status) String() string {
	return string(*f)
}

type AibiDashboardEmbeddingAccessPolicy struct {
	AccessPolicyType *AibiDashboardEmbeddingAccessPolicy_AccessPolicyType `json:"access_policy_type"`
}

type AibiDashboardEmbeddingApprovedDomains struct {
	ApprovedDomains []string `json:"approved_domains"`
}

type BooleanMessage struct {
	Value *bool `json:"value"`
}

type ClusterAutoRestartMessage struct {
	Enabled                         *bool                                        `json:"enabled"`
	CanToggle                       *bool                                        `json:"can_toggle"`
	MaintenanceWindow               *ClusterAutoRestartMessage_MaintenanceWindow `json:"maintenance_window"`
	EnablementDetails               *ClusterAutoRestartMessage_EnablementDetails `json:"enablement_details"`
	RestartEvenIfNoUpdatesAvailable *bool                                        `json:"restart_even_if_no_updates_available"`
}

// Contains an information about the enablement status judging (e.g. whether the
// enterprise tier is enabled) This is only additional information that MUST NOT
// be used to decide whether the setting is enabled or not. This is intended to
// use only for purposes like showing an error message to the customer with the
// additional details. For example, using these details we can check why exactly
// the feature is disabled for this customer..
type ClusterAutoRestartMessage_EnablementDetails struct {
	UnavailableForNonEnterpriseTier   *bool `json:"unavailable_for_non_enterprise_tier"`
	UnavailableForDisabledEntitlement *bool `json:"unavailable_for_disabled_entitlement"`
	ForcedForComplianceMode           *bool `json:"forced_for_compliance_mode"`
}

type ClusterAutoRestartMessage_MaintenanceWindow struct {
	WeekDayBasedSchedule *ClusterAutoRestartMessage_MaintenanceWindow_WeekDayBasedSchedule `json:"week_day_based_schedule"`
}

type ClusterAutoRestartMessage_MaintenanceWindow_WeekDayBasedSchedule struct {
	Frequency       *ClusterAutoRestartMessage_MaintenanceWindow_WeekDayFrequency `json:"frequency"`
	DayOfWeek       *ClusterAutoRestartMessage_MaintenanceWindow_DayOfWeek        `json:"day_of_week"`
	WindowStartTime *ClusterAutoRestartMessage_MaintenanceWindow_WindowStartTime  `json:"window_start_time"`
}

type ClusterAutoRestartMessage_MaintenanceWindow_WindowStartTime struct {
	Hours   *int `json:"hours"`
	Minutes *int `json:"minutes"`
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
	Value *PersonalComputeMessage_PersonalComputeMessageEnum `json:"value"`
}

type RestrictWorkspaceAdminsMessage struct {
	Status                *RestrictWorkspaceAdminsMessage_Status `json:"status"`
	DisableGovTagCreation *bool                                  `json:"disable_gov_tag_creation"`
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

// User Preference represents a user-specific setting scoped to an individual
// user within an account. Unlike workspace or account settings that apply to
// all users, user preferences allow personal customization (e.g., UI theme,
// editor preferences) without affecting other users..
type UserPreference struct {
	Name                *string         `json:"name"`
	UserId              *string         `json:"user_id"`
	BooleanVal          *BooleanMessage `json:"boolean_val"`
	StringVal           *StringMessage  `json:"string_val"`
	EffectiveBooleanVal *BooleanMessage `json:"effective_boolean_val"`
	EffectiveStringVal  *StringMessage  `json:"effective_string_val"`
}
