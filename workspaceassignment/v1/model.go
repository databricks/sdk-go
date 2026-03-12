package v1

type Permission string

const (
	Permission_Unknown Permission = "UNKNOWN"
	Permission_User    Permission = "USER"
	Permission_Admin   Permission = "ADMIN"
)

// String representation for [fmt.Print].
func (f *Permission) String() string {
	return string(*f)
}

// Removes all permission assignments for a workspace given a principal..
type DeleteWorkspacePermissionAssignment struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PrincipalId *int64  `json:"principal_id"`
}

type DeleteWorkspacePermissionAssignment_Response struct {
}

// Gets all the permission assignments for a workspace, given an account and a
// workspace..
type GetWorkspacePermissionAssignments struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PageToken   *string `json:"page_token"`
	MaxResults  *int    `json:"max_results"`
	Filter      *string `json:"filter"`
}

type GetWorkspacePermissionAssignments_Response struct {
	PermissionAssignments []WorkspacePermissionAssignmentOutput `json:"permission_assignments"`
	NextPageToken         *string                               `json:"next_page_token"`
	PrevPageToken         *string                               `json:"prev_page_token"`
}

// List permissions for a workspace, given an account and a workspace..
type ListWorkspacePermissions struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
}

type ListWorkspacePermissions_Response struct {
	Permissions []PermissionOutput `json:"permissions"`
}

type PermissionOutput struct {
	PermissionLevel *Permission `json:"permission_level"`
	Description     *string     `json:"description"`
}

// Information about the principal assigned to the workspace..
type PrincipalOutput struct {
	UserName             *string `json:"user_name"`
	GroupName            *string `json:"group_name"`
	ServicePrincipalName *string `json:"service_principal_name"`
	PrincipalId          *int64  `json:"principal_id"`
	DisplayName          *string `json:"display_name"`
}

type UpdateWorkspacePermissionAssignment struct {
	AccountId   *string      `json:"account_id"`
	WorkspaceId *int64       `json:"workspace_id"`
	PrincipalId *int64       `json:"principal_id"`
	Permissions []Permission `json:"permissions"`
}

// The output format for existing workspace PermissionAssignment records, which
// contains some info for user consumption..
type WorkspacePermissionAssignmentOutput struct {
	Principal   *PrincipalOutput `json:"principal"`
	Permissions []Permission     `json:"permissions"`
	Error       *string          `json:"error"`
}
