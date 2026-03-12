package v1

type PermissionLevel string

const (
	PermissionLevel_CanManage                   PermissionLevel = "CAN_MANAGE"
	PermissionLevel_CanRestart                  PermissionLevel = "CAN_RESTART"
	PermissionLevel_CanAttachTo                 PermissionLevel = "CAN_ATTACH_TO"
	PermissionLevel_IsOwner                     PermissionLevel = "IS_OWNER"
	PermissionLevel_CanManageRun                PermissionLevel = "CAN_MANAGE_RUN"
	PermissionLevel_CanView                     PermissionLevel = "CAN_VIEW"
	PermissionLevel_CanRead                     PermissionLevel = "CAN_READ"
	PermissionLevel_CanRun                      PermissionLevel = "CAN_RUN"
	PermissionLevel_CanEdit                     PermissionLevel = "CAN_EDIT"
	PermissionLevel_CanUse                      PermissionLevel = "CAN_USE"
	PermissionLevel_CanManageStagingVersions    PermissionLevel = "CAN_MANAGE_STAGING_VERSIONS"
	PermissionLevel_CanManageProductionVersions PermissionLevel = "CAN_MANAGE_PRODUCTION_VERSIONS"
	PermissionLevel_CanEditMetadata             PermissionLevel = "CAN_EDIT_METADATA"
	PermissionLevel_CanViewMetadata             PermissionLevel = "CAN_VIEW_METADATA"
	PermissionLevel_CanBind                     PermissionLevel = "CAN_BIND"
	PermissionLevel_CanQuery                    PermissionLevel = "CAN_QUERY"
	PermissionLevel_CanMonitor                  PermissionLevel = "CAN_MONITOR"
	PermissionLevel_CanCreate                   PermissionLevel = "CAN_CREATE"
	PermissionLevel_CanMonitorOnly              PermissionLevel = "CAN_MONITOR_ONLY"
	PermissionLevel_CanCreateApp                PermissionLevel = "CAN_CREATE_APP"
)

// String representation for [fmt.Print].
func (f *PermissionLevel) String() string {
	return string(*f)
}

type AccessControlRequest struct {
	UserName             *string          `json:"user_name"`
	GroupName            *string          `json:"group_name"`
	ServicePrincipalName *string          `json:"service_principal_name"`
	PermissionLevel      *PermissionLevel `json:"permission_level"`
}

type AccessControlResponse struct {
	UserName             *string      `json:"user_name"`
	GroupName            *string      `json:"group_name"`
	ServicePrincipalName *string      `json:"service_principal_name"`
	DisplayName          *string      `json:"display_name"`
	AllPermissions       []Permission `json:"all_permissions"`
}

type GetObjectPermissions struct {
	RequestObjectType *string `json:"request_object_type"`
	RequestObjectId   *string `json:"request_object_id"`
}

type GetPermissionLevels struct {
	RequestObjectType *string `json:"request_object_type"`
	RequestObjectId   *string `json:"request_object_id"`
}

type GetPermissionLevels_Response struct {
	PermissionLevels []PermissionsDescription `json:"permission_levels"`
}

type Permission struct {
	PermissionLevel     *PermissionLevel `json:"permission_level"`
	Inherited           *bool            `json:"inherited"`
	InheritedFromObject []string         `json:"inherited_from_object"`
}

type PermissionsDescription struct {
	PermissionLevel *PermissionLevel `json:"permission_level"`
	Description     *string          `json:"description"`
}

type PermissionsResponse struct {
	ObjectId          *string                 `json:"object_id"`
	ObjectType        *string                 `json:"object_type"`
	AccessControlList []AccessControlResponse `json:"access_control_list"`
}

type SetObjectPermissions struct {
	RequestObjectType *string                `json:"request_object_type"`
	RequestObjectId   *string                `json:"request_object_id"`
	AccessControlList []AccessControlRequest `json:"access_control_list"`
}

type UpdateObjectPermissions struct {
	RequestObjectType *string                `json:"request_object_type"`
	RequestObjectId   *string                `json:"request_object_id"`
	AccessControlList []AccessControlRequest `json:"access_control_list"`
}
