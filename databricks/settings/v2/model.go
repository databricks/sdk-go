
package v2

import (
	
	"fmt"
	"io"
)



type AccessPolicyType string

const (
	AccessPolicyTypeAccessPolicyTypeUnspecified AccessPolicyType = "ACCESS_POLICY_TYPE_UNSPECIFIED"
	AccessPolicyTypeAllowAllDomains AccessPolicyType = "ALLOW_ALL_DOMAINS"
	AccessPolicyTypeAllowApprovedDomains AccessPolicyType = "ALLOW_APPROVED_DOMAINS"
	AccessPolicyTypeDenyAllDomains AccessPolicyType = "DENY_ALL_DOMAINS"
)

// String representation for [fmt.Print].
func (f *AccessPolicyType) String() string {
	return string(*f)
}



type DayOfWeek string

const (
	DayOfWeekDayOfWeekUnspecified DayOfWeek = "DAY_OF_WEEK_UNSPECIFIED"
	DayOfWeekMonday DayOfWeek = "MONDAY"
	DayOfWeekTuesday DayOfWeek = "TUESDAY"
	DayOfWeekWednesday DayOfWeek = "WEDNESDAY"
	DayOfWeekThursday DayOfWeek = "THURSDAY"
	DayOfWeekFriday DayOfWeek = "FRIDAY"
	DayOfWeekSaturday DayOfWeek = "SATURDAY"
	DayOfWeekSunday DayOfWeek = "SUNDAY"
)

// String representation for [fmt.Print].
func (f *DayOfWeek) String() string {
	return string(*f)
}



type PersonalComputeMessageEnum string

const (
	PersonalComputeMessageEnumPersonalComputeMessageEnumUnspecified PersonalComputeMessageEnum = "PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED"
	PersonalComputeMessageEnumOn PersonalComputeMessageEnum = "ON"
	PersonalComputeMessageEnumDelegate PersonalComputeMessageEnum = "DELEGATE"
)

// String representation for [fmt.Print].
func (f *PersonalComputeMessageEnum) String() string {
	return string(*f)
}



type Status string

const (
	StatusStatusUnspecified Status = "STATUS_UNSPECIFIED"
	StatusAllowAll Status = "ALLOW_ALL"
	StatusRestrictTokensAndJobRunAs Status = "RESTRICT_TOKENS_AND_JOB_RUN_AS"
)

// String representation for [fmt.Print].
func (f *Status) String() string {
	return string(*f)
}



type WeekDayFrequency string

const (
	WeekDayFrequencyWeekDayFrequencyUnspecified WeekDayFrequency = "WEEK_DAY_FREQUENCY_UNSPECIFIED"
	WeekDayFrequencyFirstOfMonth WeekDayFrequency = "FIRST_OF_MONTH"
	WeekDayFrequencySecondOfMonth WeekDayFrequency = "SECOND_OF_MONTH"
	WeekDayFrequencyThirdOfMonth WeekDayFrequency = "THIRD_OF_MONTH"
	WeekDayFrequencyFourthOfMonth WeekDayFrequency = "FOURTH_OF_MONTH"
	WeekDayFrequencyFirstAndThirdOfMonth WeekDayFrequency = "FIRST_AND_THIRD_OF_MONTH"
	WeekDayFrequencySecondAndFourthOfMonth WeekDayFrequency = "SECOND_AND_FOURTH_OF_MONTH"
	WeekDayFrequencyEveryWeek WeekDayFrequency = "EVERY_WEEK"
)

// String representation for [fmt.Print].
func (f *WeekDayFrequency) String() string {
	return string(*f)
}





type AibiDashboardEmbeddingAccessPolicy struct {
	
	AccessPolicyType *AccessPolicyType `json:"accessPolicyType"`
	
}



type AibiDashboardEmbeddingApprovedDomains struct {
	
	ApprovedDomains []string `json:"approvedDomains"`
	
}



type BooleanMessage struct {
	
	Value *bool `json:"value"`
	
}



type ClusterAutoRestartMessage struct {
	
	Enabled *bool `json:"enabled"`
	
	CanToggle *bool `json:"canToggle"`
	
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow"`
	
	EnablementDetails *EnablementDetails `json:"enablementDetails"`
	
	RestartEvenIfNoUpdatesAvailable *bool `json:"restartEvenIfNoUpdatesAvailable"`
	
}



type EnablementDetails struct {
	
	UnavailableForNonEnterpriseTier *bool `json:"unavailableForNonEnterpriseTier"`
	
	UnavailableForDisabledEntitlement *bool `json:"unavailableForDisabledEntitlement"`
	
	ForcedForComplianceMode *bool `json:"forcedForComplianceMode"`
	
}



type GetPublicAccountSettingRequest struct {
	
	AccountId *string `json:"accountId"`
	
	Name *string `json:"name"`
	
}



type GetPublicAccountUserPreferenceRequest struct {
	
	AccountId *string `json:"accountId"`
	
	UserId *string `json:"userId"`
	
	Name *string `json:"name"`
	
}



type GetPublicWorkspaceSettingRequest struct {
	
	Name *string `json:"name"`
	
}



type IntegerMessage struct {
	
	Value *int `json:"value"`
	
}



type ListAccountSettingsMetadataRequest struct {
	
	AccountId *string `json:"accountId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
}



type ListAccountSettingsMetadataResponse struct {
	
	SettingsMetadata []SettingsMetadata `json:"settingsMetadata"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListAccountUserPreferencesMetadataRequest struct {
	
	AccountId *string `json:"accountId"`
	
	UserId *string `json:"userId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
}



type ListAccountUserPreferencesMetadataResponse struct {
	
	SettingsMetadata []SettingsMetadata `json:"settingsMetadata"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListWorkspaceSettingsMetadataRequest struct {
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
}



type ListWorkspaceSettingsMetadataResponse struct {
	
	SettingsMetadata []SettingsMetadata `json:"settingsMetadata"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type MaintenanceWindow struct {
	
	WeekDayBasedSchedule *WeekDayBasedSchedule `json:"weekDayBasedSchedule"`
	
}



type PatchPublicAccountSettingRequest struct {
	
	AccountId *string `json:"accountId"`
	
	Name *string `json:"name"`
	
	Setting *Setting `json:"setting"`
	
}



type PatchPublicAccountUserPreferenceRequest struct {
	
	AccountId *string `json:"accountId"`
	
	UserId *string `json:"userId"`
	
	Name *string `json:"name"`
	
	Setting *UserPreference `json:"setting"`
	
}



type PatchPublicWorkspaceSettingRequest struct {
	
	Name *string `json:"name"`
	
	Setting *Setting `json:"setting"`
	
}



type PersonalComputeMessage struct {
	
	Value *PersonalComputeMessageEnum `json:"value"`
	
}



type RestrictWorkspaceAdminsMessage struct {
	
	Status *Status `json:"status"`
	
}



type Setting struct {
	
	Name *string `json:"name"`
	
	BooleanVal *BooleanMessage `json:"booleanVal"`
	
	StringVal *StringMessage `json:"stringVal"`
	
	IntegerVal *IntegerMessage `json:"integerVal"`
	
	AutomaticClusterUpdateWorkspace *ClusterAutoRestartMessage `json:"automaticClusterUpdateWorkspace"`
	
	AibiDashboardEmbeddingApprovedDomains *AibiDashboardEmbeddingApprovedDomains `json:"aibiDashboardEmbeddingApprovedDomains"`
	
	AibiDashboardEmbeddingAccessPolicy *AibiDashboardEmbeddingAccessPolicy `json:"aibiDashboardEmbeddingAccessPolicy"`
	
	RestrictWorkspaceAdmins *RestrictWorkspaceAdminsMessage `json:"restrictWorkspaceAdmins"`
	
	PersonalCompute *PersonalComputeMessage `json:"personalCompute"`
	
	EffectiveBooleanVal *BooleanMessage `json:"effectiveBooleanVal"`
	
	EffectiveStringVal *StringMessage `json:"effectiveStringVal"`
	
	EffectiveIntegerVal *IntegerMessage `json:"effectiveIntegerVal"`
	
	EffectiveAutomaticClusterUpdateWorkspace *ClusterAutoRestartMessage `json:"effectiveAutomaticClusterUpdateWorkspace"`
	
	EffectiveAibiDashboardEmbeddingApprovedDomains *AibiDashboardEmbeddingApprovedDomains `json:"effectiveAibiDashboardEmbeddingApprovedDomains"`
	
	EffectiveAibiDashboardEmbeddingAccessPolicy *AibiDashboardEmbeddingAccessPolicy `json:"effectiveAibiDashboardEmbeddingAccessPolicy"`
	
	EffectiveRestrictWorkspaceAdmins *RestrictWorkspaceAdminsMessage `json:"effectiveRestrictWorkspaceAdmins"`
	
	EffectivePersonalCompute *PersonalComputeMessage `json:"effectivePersonalCompute"`
	
}



type SettingsMetadata struct {
	
	Name *string `json:"name"`
	
	Description *string `json:"description"`
	
	Type *string `json:"type"`
	
	DocsLink *string `json:"docsLink"`
	
}



type StringMessage struct {
	
	Value *string `json:"value"`
	
}



type UserPreference struct {
	
	Name *string `json:"name"`
	
	UserId *string `json:"userId"`
	
	BooleanVal *BooleanMessage `json:"booleanVal"`
	
	StringVal *StringMessage `json:"stringVal"`
	
	EffectiveBooleanVal *BooleanMessage `json:"effectiveBooleanVal"`
	
	EffectiveStringVal *StringMessage `json:"effectiveStringVal"`
	
}



type WeekDayBasedSchedule struct {
	
	Frequency *WeekDayFrequency `json:"frequency"`
	
	DayOfWeek *DayOfWeek `json:"dayOfWeek"`
	
	WindowStartTime *WindowStartTime `json:"windowStartTime"`
	
}



type WindowStartTime struct {
	
	Hours *int `json:"hours"`
	
	Minutes *int `json:"minutes"`
	
}


