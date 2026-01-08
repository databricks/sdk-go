package v2

type AccessType string

const (
	AccessTypeAccessTypeUnspecified AccessType = "ACCESS_TYPE_UNSPECIFIED"
	AccessTypeDirect                AccessType = "DIRECT"
	AccessTypeIndirect              AccessType = "INDIRECT"
)

// String representation for [fmt.Print].
func (f *AccessType) String() string {
	return string(*f)
}

type AccountAccessRuleAction string

const (
	AccountAccessRuleActionAccountAccessRuleActionUnspecified AccountAccessRuleAction = "ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED"
	AccountAccessRuleActionAllow                              AccountAccessRuleAction = "ALLOW"
	AccountAccessRuleActionDeny                               AccountAccessRuleAction = "DENY"
)

// String representation for [fmt.Print].
func (f *AccountAccessRuleAction) String() string {
	return string(*f)
}

type PrincipalType string

const (
	PrincipalTypePrincipalTypeUnspecified PrincipalType = "PRINCIPAL_TYPE_UNSPECIFIED"
	PrincipalTypeUser                     PrincipalType = "USER"
	PrincipalTypeServicePrincipal         PrincipalType = "SERVICE_PRINCIPAL"
	PrincipalTypeGroup                    PrincipalType = "GROUP"
)

// String representation for [fmt.Print].
func (f *PrincipalType) String() string {
	return string(*f)
}

type State string

const (
	StateStateUnspecified State = "STATE_UNSPECIFIED"
	StateActive           State = "ACTIVE"
	StateInactive         State = "INACTIVE"
)

// String representation for [fmt.Print].
func (f *State) String() string {
	return string(*f)
}

type WorkspaceAccessDetailView string

const (
	WorkspaceAccessDetailViewWorkspaceAccessDetailViewUnspecified WorkspaceAccessDetailView = "WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED"
	WorkspaceAccessDetailViewBasic                                WorkspaceAccessDetailView = "BASIC"
	WorkspaceAccessDetailViewFull                                 WorkspaceAccessDetailView = "FULL"
)

// String representation for [fmt.Print].
func (f *WorkspaceAccessDetailView) String() string {
	return string(*f)
}

type WorkspacePermission string

const (
	WorkspacePermissionWorkspacePermissionUnspecified WorkspacePermission = "WORKSPACE_PERMISSION_UNSPECIFIED"
	WorkspacePermissionUserPermission                 WorkspacePermission = "USER_PERMISSION"
	WorkspacePermissionAdminPermission                WorkspacePermission = "ADMIN_PERMISSION"
)

// String representation for [fmt.Print].
func (f *WorkspacePermission) String() string {
	return string(*f)
}

type AccountAccessIdentityRule struct {
	Action        *AccountAccessRuleAction `json:"action"`
	ExternalId    *string                  `json:"external_id"`
	DisplayName   *string                  `json:"display_name"`
	PrincipalType *PrincipalType           `json:"principal_type"`
}

type CreateAccountAccessIdentityRuleRequest struct {
	AccountId                 *string                    `json:"account_id"`
	AccountAccessIdentityRule *AccountAccessIdentityRule `json:"account_access_identity_rule"`
}

type CreateGroupProxyRequest struct {
	Group *Group `json:"group"`
}

type CreateGroupRequest struct {
	AccountId *string `json:"account_id"`
	Group     *Group  `json:"group"`
}

type CreateServicePrincipalProxyRequest struct {
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

type CreateServicePrincipalRequest struct {
	AccountId        *string           `json:"account_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

type CreateUserProxyRequest struct {
	User *User `json:"user"`
}

type CreateUserRequest struct {
	AccountId *string `json:"account_id"`
	User      *User   `json:"user"`
}

type CreateWorkspaceAccessDetailLocalRequest struct {
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
}

type CreateWorkspaceAccessDetailRequest struct {
	AccountId             *string                `json:"account_id"`
	Parent                *string                `json:"parent"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
}

type DeleteAccountAccessIdentityRuleRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type DeleteGroupProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type DeleteGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type DeleteServicePrincipalProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type DeleteServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type DeleteUserProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type DeleteUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type DeleteWorkspaceAccessDetailLocalRequest struct {
	PrincipalId *int64 `json:"principal_id"`
}

type DeleteWorkspaceAccessDetailRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PrincipalId *int64  `json:"principal_id"`
}

type GetAccountAccessIdentityRuleRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type GetGroupProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type GetGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type GetServicePrincipalProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type GetServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type GetUserProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

type GetUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

type GetWorkspaceAccessDetailLocalRequest struct {
	PrincipalId *int64                     `json:"principal_id"`
	View        *WorkspaceAccessDetailView `json:"view"`
}

type GetWorkspaceAccessDetailRequest struct {
	AccountId   *string                    `json:"account_id"`
	WorkspaceId *int64                     `json:"workspace_id"`
	PrincipalId *int64                     `json:"principal_id"`
	View        *WorkspaceAccessDetailView `json:"view"`
}

type Group struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	ExternalId *string `json:"external_id"`
	GroupName  *string `json:"group_name"`
}

type ListAccountAccessIdentityRulesRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListAccountAccessIdentityRulesResponse struct {
	AccountAccessIdentityRules []AccountAccessIdentityRule `json:"account_access_identity_rules"`
	NextPageToken              *string                     `json:"next_page_token"`
}

type ListGroupsProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListGroupsRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListGroupsResponse struct {
	Groups        []Group `json:"groups"`
	NextPageToken *string `json:"next_page_token"`
}

type ListServicePrincipalsProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListServicePrincipalsRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListServicePrincipalsResponse struct {
	ServicePrincipals []ServicePrincipal `json:"service_principals"`
	NextPageToken     *string            `json:"next_page_token"`
}

type ListUsersProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListUsersRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

type ListUsersResponse struct {
	Users         []User  `json:"users"`
	NextPageToken *string `json:"next_page_token"`
}

type ListWorkspaceAccessDetailsLocalRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListWorkspaceAccessDetailsRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

type ListWorkspaceAccessDetailsResponse struct {
	WorkspaceAccessDetails []WorkspaceAccessDetail `json:"workspace_access_details"`
	NextPageToken          *string                 `json:"next_page_token"`
}

type Name struct {
	GivenName  *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
}

type ResolveGroupProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

type ResolveGroupRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveGroupResponse struct {
	Group *Group `json:"group"`
}

type ResolveServicePrincipalProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

type ResolveServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveServicePrincipalResponse struct {
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

type ResolveUserProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

type ResolveUserRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveUserResponse struct {
	User *User `json:"user"`
}

type ServicePrincipal struct {
	AccountId       *string `json:"account_id"`
	InternalId      *int64  `json:"internal_id"`
	ExternalId      *string `json:"external_id"`
	ApplicationId   *string `json:"application_id"`
	DisplayName     *string `json:"display_name"`
	AccountSpStatus *State  `json:"account_sp_status"`
}

type UpdateGroupProxyRequest struct {
	InternalId *int64  `json:"internal_id"`
	Group      *Group  `json:"group"`
	UpdateMask *string `json:"update_mask"`
}

type UpdateGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	Group      *Group  `json:"group"`
	UpdateMask *string `json:"update_mask"`
}

type UpdateServicePrincipalProxyRequest struct {
	InternalId       *int64            `json:"internal_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
	UpdateMask       *string           `json:"update_mask"`
}

type UpdateServicePrincipalRequest struct {
	AccountId        *string           `json:"account_id"`
	InternalId       *int64            `json:"internal_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
	UpdateMask       *string           `json:"update_mask"`
}

type UpdateUserProxyRequest struct {
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

type UpdateUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

type UpdateWorkspaceAccessDetailLocalRequest struct {
	PrincipalId           *int64                 `json:"principal_id"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
	UpdateMask            *string                `json:"update_mask"`
}

type UpdateWorkspaceAccessDetailRequest struct {
	AccountId             *string                `json:"account_id"`
	WorkspaceId           *int64                 `json:"workspace_id"`
	PrincipalId           *int64                 `json:"principal_id"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
	UpdateMask            *string                `json:"update_mask"`
}

type User struct {
	AccountId         *string `json:"account_id"`
	InternalId        *int64  `json:"internal_id"`
	ExternalId        *string `json:"external_id"`
	Username          *string `json:"username"`
	Name              *Name   `json:"name"`
	AccountUserStatus *State  `json:"account_user_status"`
}

type WorkspaceAccessDetail struct {
	PrincipalId   *int64                `json:"principal_id"`
	WorkspaceId   *int64                `json:"workspace_id"`
	AccountId     *string               `json:"account_id"`
	PrincipalType *PrincipalType        `json:"principal_type"`
	AccessType    *AccessType           `json:"access_type"`
	Status        *State                `json:"status"`
	Permissions   []WorkspacePermission `json:"permissions"`
}
