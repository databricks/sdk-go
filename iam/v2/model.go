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

// An identity rule that controls which principals can access an account..
type AccountAccessIdentityRule struct {
	Action        *AccountAccessRuleAction `json:"action"`
	ExternalId    *string                  `json:"external_id"`
	DisplayName   *string                  `json:"display_name"`
	PrincipalType *PrincipalType           `json:"principal_type"`
}

// Request message for creating an account access identity rule..
type CreateAccountAccessIdentityRuleRequest struct {
	AccountId                 *string                    `json:"account_id"`
	AccountAccessIdentityRule *AccountAccessIdentityRule `json:"account_access_identity_rule"`
}

// TODO: Write description later when this method is implemented.
type CreateGroupProxyRequest struct {
	Group *Group `json:"group"`
}

// TODO: Write description later when this method is implemented.
type CreateGroupRequest struct {
	AccountId *string `json:"account_id"`
	Group     *Group  `json:"group"`
}

// TODO: Write description later when this method is implemented.
type CreateServicePrincipalProxyRequest struct {
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

// TODO: Write description later when this method is implemented.
type CreateServicePrincipalRequest struct {
	AccountId        *string           `json:"account_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

// TODO: Write description later when this method is implemented.
type CreateUserProxyRequest struct {
	User *User `json:"user"`
}

// TODO: Write description later when this method is implemented.
type CreateUserRequest struct {
	AccountId *string `json:"account_id"`
	User      *User   `json:"user"`
}

// TODO: Write description later when this method is implemented.
type CreateWorkspaceAccessDetailLocalRequest struct {
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
}

// TODO: Write description later when this method is implemented.
type CreateWorkspaceAccessDetailRequest struct {
	AccountId             *string                `json:"account_id"`
	Parent                *string                `json:"parent"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
}

// Request message for deleting an account access identity rule..
type DeleteAccountAccessIdentityRuleRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteGroupProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteServicePrincipalProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteUserProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteWorkspaceAccessDetailLocalRequest struct {
	PrincipalId *int64 `json:"principal_id"`
}

// TODO: Write description later when this method is implemented.
type DeleteWorkspaceAccessDetailRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PrincipalId *int64  `json:"principal_id"`
}

// Request message for getting an account access identity rule..
type GetAccountAccessIdentityRuleRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

// TODO: Write description later when this method is implemented.
type GetGroupProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type GetGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type GetServicePrincipalProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type GetServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type GetUserProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// TODO: Write description later when this method is implemented.
type GetUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
}

// Request message for getting the access details for a principal in the current
// workspace..
type GetWorkspaceAccessDetailLocalRequest struct {
	PrincipalId *int64                     `json:"principal_id"`
	View        *WorkspaceAccessDetailView `json:"view"`
}

// Request message for getting the access details for a principal in a
// workspace..
type GetWorkspaceAccessDetailRequest struct {
	AccountId   *string                    `json:"account_id"`
	WorkspaceId *int64                     `json:"workspace_id"`
	PrincipalId *int64                     `json:"principal_id"`
	View        *WorkspaceAccessDetailView `json:"view"`
}

// The details of a Group resource..
type Group struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	ExternalId *string `json:"external_id"`
	GroupName  *string `json:"group_name"`
}

// Request message for listing account access identity rules..
type ListAccountAccessIdentityRulesRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// Response message for listing account access identity rules..
type ListAccountAccessIdentityRulesResponse struct {
	AccountAccessIdentityRules []AccountAccessIdentityRule `json:"account_access_identity_rules"`
	NextPageToken              *string                     `json:"next_page_token"`
}

// TODO: Write description later when this method is implemented.
type ListGroupsProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListGroupsRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListGroupsResponse struct {
	Groups        []Group `json:"groups"`
	NextPageToken *string `json:"next_page_token"`
}

// TODO: Write description later when this method is implemented.
type ListServicePrincipalsProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListServicePrincipalsRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListServicePrincipalsResponse struct {
	ServicePrincipals []ServicePrincipal `json:"service_principals"`
	NextPageToken     *string            `json:"next_page_token"`
}

// TODO: Write description later when this method is implemented.
type ListUsersProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListUsersRequest struct {
	AccountId *string `json:"account_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// TODO: Write description later when this method is implemented.
type ListUsersResponse struct {
	Users         []User  `json:"users"`
	NextPageToken *string `json:"next_page_token"`
}

// TODO: Write description later when this method is implemented.
type ListWorkspaceAccessDetailsLocalRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// TODO: Write description later when this method is implemented.
type ListWorkspaceAccessDetailsRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

// TODO: Write description later when this method is implemented.
type ListWorkspaceAccessDetailsResponse struct {
	WorkspaceAccessDetails []WorkspaceAccessDetail `json:"workspace_access_details"`
	NextPageToken          *string                 `json:"next_page_token"`
}

type Name struct {
}

// Request message for resolving a group with the given external ID from the
// customer's IdP into <Databricks>. Will resolve metadata such as the group's
// groupname, and inherited parent groups..
type ResolveGroupProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

// Request message for resolving a group with the given external ID from the
// customer's IdP into <Databricks>. Will resolve metadata such as the group's
// groupname, and inherited parent groups..
type ResolveGroupRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveGroupResponse struct {
	Group *Group `json:"group"`
}

// Request message for resolving a service principal with the given external ID
// from the customer's IdP into <Databricks>. Will resolve metadata such as the
// service principal's displayname, status, and inherited parent groups..
type ResolveServicePrincipalProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

// Request message for resolving a service principal with the given external ID
// from the customer's IdP into <Databricks>. Will resolve metadata such as the
// service principal's displayname, status, and inherited parent groups..
type ResolveServicePrincipalRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveServicePrincipalResponse struct {
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
}

// Request message for resolving a user with the given external ID from the
// customer's IdP into <Databricks>. Will resolve metadata such as the user's
// displayname, status, and inherited parent groups..
type ResolveUserProxyRequest struct {
	ExternalId *string `json:"external_id"`
}

// Request message for resolving a user with the given external ID from the
// customer's IdP into <Databricks>. Will resolve metadata such as the user's
// displayname, status, and inherited parent groups..
type ResolveUserRequest struct {
	AccountId  *string `json:"account_id"`
	ExternalId *string `json:"external_id"`
}

type ResolveUserResponse struct {
	User *User `json:"user"`
}

// The details of a ServicePrincipal resource..
type ServicePrincipal struct {
	AccountId       *string `json:"account_id"`
	InternalId      *int64  `json:"internal_id"`
	ExternalId      *string `json:"external_id"`
	ApplicationId   *string `json:"application_id"`
	DisplayName     *string `json:"display_name"`
	AccountSpStatus *State  `json:"account_sp_status"`
}

// TODO: Write description later when this method is implemented.
type UpdateGroupProxyRequest struct {
	InternalId *int64  `json:"internal_id"`
	Group      *Group  `json:"group"`
	UpdateMask *string `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateGroupRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	Group      *Group  `json:"group"`
	UpdateMask *string `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateServicePrincipalProxyRequest struct {
	InternalId       *int64            `json:"internal_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
	UpdateMask       *string           `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateServicePrincipalRequest struct {
	AccountId        *string           `json:"account_id"`
	InternalId       *int64            `json:"internal_id"`
	ServicePrincipal *ServicePrincipal `json:"service_principal"`
	UpdateMask       *string           `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateUserProxyRequest struct {
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateWorkspaceAccessDetailLocalRequest struct {
	PrincipalId           *int64                 `json:"principal_id"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
	UpdateMask            *string                `json:"update_mask"`
}

// TODO: Write description later when this method is implemented.
type UpdateWorkspaceAccessDetailRequest struct {
	AccountId             *string                `json:"account_id"`
	WorkspaceId           *int64                 `json:"workspace_id"`
	PrincipalId           *int64                 `json:"principal_id"`
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspace_access_detail"`
	UpdateMask            *string                `json:"update_mask"`
}

// The details of a User resource..
type User struct {
	AccountId         *string `json:"account_id"`
	InternalId        *int64  `json:"internal_id"`
	ExternalId        *string `json:"external_id"`
	Username          *string `json:"username"`
	Name              *Name   `json:"name"`
	AccountUserStatus *State  `json:"account_user_status"`
}

// The details of a principal's access to a workspace..
type WorkspaceAccessDetail struct {
	PrincipalId   *int64                `json:"principal_id"`
	WorkspaceId   *int64                `json:"workspace_id"`
	AccountId     *string               `json:"account_id"`
	PrincipalType *PrincipalType        `json:"principal_type"`
	AccessType    *AccessType           `json:"access_type"`
	Status        *State                `json:"status"`
	Permissions   []WorkspacePermission `json:"permissions"`
}
