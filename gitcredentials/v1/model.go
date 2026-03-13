package v1

type CreateCredentials struct {
	GitProvider          *string `json:"git_provider"`
	GitUsername          *string `json:"git_username"`
	PersonalAccessToken  *string `json:"personal_access_token"`
	PrincipalId          *int64  `json:"principal_id"`
	Name                 *string `json:"name"`
	IsDefaultForProvider *bool   `json:"is_default_for_provider"`
	GitEmail             *string `json:"git_email"`
}

type CreateCredentials_Response struct {
	CredentialId         *int64  `json:"credential_id"`
	GitProvider          *string `json:"git_provider"`
	GitUsername          *string `json:"git_username"`
	Name                 *string `json:"name"`
	IsDefaultForProvider *bool   `json:"is_default_for_provider"`
	GitEmail             *string `json:"git_email"`
}

type Credential struct {
	CredentialId         *int64  `json:"credential_id"`
	GitProvider          *string `json:"git_provider"`
	GitUsername          *string `json:"git_username"`
	Name                 *string `json:"name"`
	IsDefaultForProvider *bool   `json:"is_default_for_provider"`
	GitEmail             *string `json:"git_email"`
}

type DeleteCredentials struct {
	Id          *int64 `json:"id"`
	PrincipalId *int64 `json:"principal_id"`
}

type DeleteCredentials_Response struct {
}

type GetCredentials struct {
	Id          *int64 `json:"id"`
	PrincipalId *int64 `json:"principal_id"`
}

type GetCredentials_Response struct {
	CredentialId         *int64  `json:"credential_id"`
	GitProvider          *string `json:"git_provider"`
	GitUsername          *string `json:"git_username"`
	Name                 *string `json:"name"`
	IsDefaultForProvider *bool   `json:"is_default_for_provider"`
	GitEmail             *string `json:"git_email"`
}

type ListCredentials struct {
	PrincipalId *int64 `json:"principal_id"`
}

type ListCredentials_Response struct {
	Credentials []Credential `json:"credentials"`
}

type UpdateCredentials struct {
	Id                   *int64  `json:"id"`
	PersonalAccessToken  *string `json:"personal_access_token"`
	GitProvider          *string `json:"git_provider"`
	GitUsername          *string `json:"git_username"`
	PrincipalId          *int64  `json:"principal_id"`
	Name                 *string `json:"name"`
	IsDefaultForProvider *bool   `json:"is_default_for_provider"`
	GitEmail             *string `json:"git_email"`
}

type UpdateCredentials_Response struct {
}
