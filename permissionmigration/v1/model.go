package v1

type MigratePermissionsRequest struct {
	WorkspaceId            *int64  `json:"workspace_id"`
	FromWorkspaceGroupName *string `json:"from_workspace_group_name"`
	ToAccountGroupName     *string `json:"to_account_group_name"`
	Size                   *int    `json:"size"`
}

type MigratePermissionsResponse struct {
	PermissionsMigrated *int `json:"permissions_migrated"`
}
