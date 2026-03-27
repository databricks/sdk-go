package v1

type EffectivePrivilege struct {
	Privilege         *string `json:"privilege"`
	InheritedFromType *string `json:"inherited_from_type"`
	InheritedFromName *string `json:"inherited_from_name"`
}

type EffectivePrivilegeAssignment struct {
	Principal  *string              `json:"principal"`
	Privileges []EffectivePrivilege `json:"privileges"`
}

type GetEffectivePermissions struct {
	SecurableType     *string `json:"securable_type"`
	SecurableFullName *string `json:"securable_full_name"`
	Principal         *string `json:"principal"`
	MaxResults        *int    `json:"max_results"`
	PageToken         *string `json:"page_token"`
}

type GetEffectivePermissions_Response struct {
	NextPageToken        *string                        `json:"next_page_token"`
	PrivilegeAssignments []EffectivePrivilegeAssignment `json:"privilege_assignments"`
}

type GetPermissions struct {
	SecurableType            *string `json:"securable_type"`
	SecurableFullName        *string `json:"securable_full_name"`
	Principal                *string `json:"principal"`
	MaxResults               *int    `json:"max_results"`
	PageToken                *string `json:"page_token"`
	IncludeDeletedPrincipals *bool   `json:"include_deleted_principals"`
}

type GetPermissions_Response struct {
	NextPageToken        *string               `json:"next_page_token"`
	PrivilegeAssignments []PrivilegeAssignment `json:"privilege_assignments"`
}

type ListPrivilegeAssignmentsRequest struct {
	SecurableType            *string `json:"securable_type"`
	FullName                 *string `json:"full_name"`
	Principal                *string `json:"principal"`
	IncludeDeletedPrincipals *bool   `json:"include_deleted_principals"`
	PageSize                 *int    `json:"page_size"`
	PageToken                *string `json:"page_token"`
}

type ListPrivilegeAssignmentsResponse struct {
	PrivilegeAssignments []PrivilegeAssignment `json:"privilege_assignments"`
	NextPageToken        *string               `json:"next_page_token"`
}

type PermissionsChange struct {
	Principal   *string  `json:"principal"`
	Add         []string `json:"add"`
	Remove      []string `json:"remove"`
	PrincipalId *int64   `json:"principal_id"`
}

type PrivilegeAssignment struct {
	Principal   *string  `json:"principal"`
	Privileges  []string `json:"privileges"`
	PrincipalId *int64   `json:"principal_id"`
}

type UpdatePermissions struct {
	SecurableType     *string             `json:"securable_type"`
	SecurableFullName *string             `json:"securable_full_name"`
	Changes           []PermissionsChange `json:"changes"`
}

type UpdatePermissions_Response struct {
	PrivilegeAssignments []PrivilegeAssignment `json:"privilege_assignments"`
}
