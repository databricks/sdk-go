
package v2

import (
	
	"fmt"
	"io"
)



type AccessPolicyType string

const AccessPolicyTypeAccessPolicyTypeUnspecified AccessPolicyType = `ACCESS_POLICY_TYPE_UNSPECIFIED`
const AccessPolicyTypeAllowAllDomains AccessPolicyType = `ALLOW_ALL_DOMAINS`
const AccessPolicyTypeAllowApprovedDomains AccessPolicyType = `ALLOW_APPROVED_DOMAINS`
const AccessPolicyTypeDenyAllDomains AccessPolicyType = `DENY_ALL_DOMAINS`

// String representation for [fmt.Print]
func (f *AccessPolicyType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *AccessPolicyType) Set(v string) error {
	switch v {
	case `ACCESS_POLICY_TYPE_UNSPECIFIED`, `ALLOW_ALL_DOMAINS`, `ALLOW_APPROVED_DOMAINS`, `DENY_ALL_DOMAINS`:
		*f = AccessPolicyType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ACCESS_POLICY_TYPE_UNSPECIFIED", "ALLOW_ALL_DOMAINS", "ALLOW_APPROVED_DOMAINS", "DENY_ALL_DOMAINS"`, v)
	}
}

// Values returns all possible values for AccessPolicyType.
//
// There is no guarantee on the order of the values in the slice.
func (f *AccessPolicyType) Values() []AccessPolicyType {
	return []AccessPolicyType{
		AccessPolicyTypeAccessPolicyTypeUnspecified,
		AccessPolicyTypeAllowAllDomains,
		AccessPolicyTypeAllowApprovedDomains,
		AccessPolicyTypeDenyAllDomains,
	}
}

// Type always returns AccessPolicyType to satisfy [pflag.Value] interface
func (f *AccessPolicyType) Type() string {
	return "AccessPolicyType"
}



type DayOfWeek string

const DayOfWeekDayOfWeekUnspecified DayOfWeek = `DAY_OF_WEEK_UNSPECIFIED`
const DayOfWeekMonday DayOfWeek = `MONDAY`
const DayOfWeekTuesday DayOfWeek = `TUESDAY`
const DayOfWeekWednesday DayOfWeek = `WEDNESDAY`
const DayOfWeekThursday DayOfWeek = `THURSDAY`
const DayOfWeekFriday DayOfWeek = `FRIDAY`
const DayOfWeekSaturday DayOfWeek = `SATURDAY`
const DayOfWeekSunday DayOfWeek = `SUNDAY`

// String representation for [fmt.Print]
func (f *DayOfWeek) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *DayOfWeek) Set(v string) error {
	switch v {
	case `DAY_OF_WEEK_UNSPECIFIED`, `MONDAY`, `TUESDAY`, `WEDNESDAY`, `THURSDAY`, `FRIDAY`, `SATURDAY`, `SUNDAY`:
		*f = DayOfWeek(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "DAY_OF_WEEK_UNSPECIFIED", "MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"`, v)
	}
}

// Values returns all possible values for DayOfWeek.
//
// There is no guarantee on the order of the values in the slice.
func (f *DayOfWeek) Values() []DayOfWeek {
	return []DayOfWeek{
		DayOfWeekDayOfWeekUnspecified,
		DayOfWeekMonday,
		DayOfWeekTuesday,
		DayOfWeekWednesday,
		DayOfWeekThursday,
		DayOfWeekFriday,
		DayOfWeekSaturday,
		DayOfWeekSunday,
	}
}

// Type always returns DayOfWeek to satisfy [pflag.Value] interface
func (f *DayOfWeek) Type() string {
	return "DayOfWeek"
}



type PersonalComputeMessageEnum string

const PersonalComputeMessageEnumPersonalComputeMessageEnumUnspecified PersonalComputeMessageEnum = `PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED`
const PersonalComputeMessageEnumOn PersonalComputeMessageEnum = `ON`
const PersonalComputeMessageEnumDelegate PersonalComputeMessageEnum = `DELEGATE`

// String representation for [fmt.Print]
func (f *PersonalComputeMessageEnum) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *PersonalComputeMessageEnum) Set(v string) error {
	switch v {
	case `PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED`, `ON`, `DELEGATE`:
		*f = PersonalComputeMessageEnum(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "PERSONAL_COMPUTE_MESSAGE_ENUM_UNSPECIFIED", "ON", "DELEGATE"`, v)
	}
}

// Values returns all possible values for PersonalComputeMessageEnum.
//
// There is no guarantee on the order of the values in the slice.
func (f *PersonalComputeMessageEnum) Values() []PersonalComputeMessageEnum {
	return []PersonalComputeMessageEnum{
		PersonalComputeMessageEnumPersonalComputeMessageEnumUnspecified,
		PersonalComputeMessageEnumOn,
		PersonalComputeMessageEnumDelegate,
	}
}

// Type always returns PersonalComputeMessageEnum to satisfy [pflag.Value] interface
func (f *PersonalComputeMessageEnum) Type() string {
	return "PersonalComputeMessageEnum"
}



type Status string

const StatusStatusUnspecified Status = `STATUS_UNSPECIFIED`
const StatusAllowAll Status = `ALLOW_ALL`
const StatusRestrictTokensAndJobRunAs Status = `RESTRICT_TOKENS_AND_JOB_RUN_AS`

// String representation for [fmt.Print]
func (f *Status) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Status) Set(v string) error {
	switch v {
	case `STATUS_UNSPECIFIED`, `ALLOW_ALL`, `RESTRICT_TOKENS_AND_JOB_RUN_AS`:
		*f = Status(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "STATUS_UNSPECIFIED", "ALLOW_ALL", "RESTRICT_TOKENS_AND_JOB_RUN_AS"`, v)
	}
}

// Values returns all possible values for Status.
//
// There is no guarantee on the order of the values in the slice.
func (f *Status) Values() []Status {
	return []Status{
		StatusStatusUnspecified,
		StatusAllowAll,
		StatusRestrictTokensAndJobRunAs,
	}
}

// Type always returns Status to satisfy [pflag.Value] interface
func (f *Status) Type() string {
	return "Status"
}



type WeekDayFrequency string

const WeekDayFrequencyWeekDayFrequencyUnspecified WeekDayFrequency = `WEEK_DAY_FREQUENCY_UNSPECIFIED`
const WeekDayFrequencyFirstOfMonth WeekDayFrequency = `FIRST_OF_MONTH`
const WeekDayFrequencySecondOfMonth WeekDayFrequency = `SECOND_OF_MONTH`
const WeekDayFrequencyThirdOfMonth WeekDayFrequency = `THIRD_OF_MONTH`
const WeekDayFrequencyFourthOfMonth WeekDayFrequency = `FOURTH_OF_MONTH`
const WeekDayFrequencyFirstAndThirdOfMonth WeekDayFrequency = `FIRST_AND_THIRD_OF_MONTH`
const WeekDayFrequencySecondAndFourthOfMonth WeekDayFrequency = `SECOND_AND_FOURTH_OF_MONTH`
const WeekDayFrequencyEveryWeek WeekDayFrequency = `EVERY_WEEK`

// String representation for [fmt.Print]
func (f *WeekDayFrequency) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *WeekDayFrequency) Set(v string) error {
	switch v {
	case `WEEK_DAY_FREQUENCY_UNSPECIFIED`, `FIRST_OF_MONTH`, `SECOND_OF_MONTH`, `THIRD_OF_MONTH`, `FOURTH_OF_MONTH`, `FIRST_AND_THIRD_OF_MONTH`, `SECOND_AND_FOURTH_OF_MONTH`, `EVERY_WEEK`:
		*f = WeekDayFrequency(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "WEEK_DAY_FREQUENCY_UNSPECIFIED", "FIRST_OF_MONTH", "SECOND_OF_MONTH", "THIRD_OF_MONTH", "FOURTH_OF_MONTH", "FIRST_AND_THIRD_OF_MONTH", "SECOND_AND_FOURTH_OF_MONTH", "EVERY_WEEK"`, v)
	}
}

// Values returns all possible values for WeekDayFrequency.
//
// There is no guarantee on the order of the values in the slice.
func (f *WeekDayFrequency) Values() []WeekDayFrequency {
	return []WeekDayFrequency{
		WeekDayFrequencyWeekDayFrequencyUnspecified,
		WeekDayFrequencyFirstOfMonth,
		WeekDayFrequencySecondOfMonth,
		WeekDayFrequencyThirdOfMonth,
		WeekDayFrequencyFourthOfMonth,
		WeekDayFrequencyFirstAndThirdOfMonth,
		WeekDayFrequencySecondAndFourthOfMonth,
		WeekDayFrequencyEveryWeek,
	}
}

// Type always returns WeekDayFrequency to satisfy [pflag.Value] interface
func (f *WeekDayFrequency) Type() string {
	return "WeekDayFrequency"
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


