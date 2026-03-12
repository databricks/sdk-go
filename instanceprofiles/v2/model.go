package v2

type AddInstanceProfile struct {
	SkipValidation        *bool   `json:"skip_validation"`
	InstanceProfileArn    *string `json:"instance_profile_arn"`
	IsMetaInstanceProfile *bool   `json:"is_meta_instance_profile"`
	IamRoleArn            *string `json:"iam_role_arn"`
}

type AddInstanceProfile_Response struct {
}

type EditInstanceProfile struct {
	InstanceProfileArn    *string `json:"instance_profile_arn"`
	IsMetaInstanceProfile *bool   `json:"is_meta_instance_profile"`
	IamRoleArn            *string `json:"iam_role_arn"`
}

type EditInstanceProfile_Response struct {
}

type InstanceProfile struct {
	InstanceProfileArn    *string `json:"instance_profile_arn"`
	IsMetaInstanceProfile *bool   `json:"is_meta_instance_profile"`
	IamRoleArn            *string `json:"iam_role_arn"`
}

type ListInstanceProfiles struct {
}

type ListInstanceProfiles_Response struct {
	InstanceProfiles []InstanceProfile `json:"instance_profiles"`
}

type RemoveInstanceProfile struct {
	InstanceProfileArn *string `json:"instance_profile_arn"`
}

type RemoveInstanceProfile_Response struct {
}
