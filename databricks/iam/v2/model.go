
package v2

import (
	
	"fmt"
	"io"
)



type AccessType string

const (
	AccessTypeAccessTypeUnspecified AccessType = "ACCESS_TYPE_UNSPECIFIED"
	AccessTypeDirect AccessType = "DIRECT"
	AccessTypeIndirect AccessType = "INDIRECT"
)

// String representation for [fmt.Print].
func (f *AccessType) String() string {
	return string(*f)
}



type AccountAccessRuleAction string

const (
	AccountAccessRuleActionAccountAccessRuleActionUnspecified AccountAccessRuleAction = "ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED"
	AccountAccessRuleActionAllow AccountAccessRuleAction = "ALLOW"
	AccountAccessRuleActionDeny AccountAccessRuleAction = "DENY"
)

// String representation for [fmt.Print].
func (f *AccountAccessRuleAction) String() string {
	return string(*f)
}



type PrincipalType string

const (
	PrincipalTypePrincipalTypeUnspecified PrincipalType = "PRINCIPAL_TYPE_UNSPECIFIED"
	PrincipalTypeUser PrincipalType = "USER"
	PrincipalTypeServicePrincipal PrincipalType = "SERVICE_PRINCIPAL"
	PrincipalTypeGroup PrincipalType = "GROUP"
)

// String representation for [fmt.Print].
func (f *PrincipalType) String() string {
	return string(*f)
}



type State string

const (
	StateStateUnspecified State = "STATE_UNSPECIFIED"
	StateActive State = "ACTIVE"
	StateInactive State = "INACTIVE"
)

// String representation for [fmt.Print].
func (f *State) String() string {
	return string(*f)
}



type WorkspaceAccessDetailView string

const (
	WorkspaceAccessDetailViewWorkspaceAccessDetailViewUnspecified WorkspaceAccessDetailView = "WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED"
	WorkspaceAccessDetailViewBasic WorkspaceAccessDetailView = "BASIC"
	WorkspaceAccessDetailViewFull WorkspaceAccessDetailView = "FULL"
)

// String representation for [fmt.Print].
func (f *WorkspaceAccessDetailView) String() string {
	return string(*f)
}



type WorkspacePermission string

const (
	WorkspacePermissionWorkspacePermissionUnspecified WorkspacePermission = "WORKSPACE_PERMISSION_UNSPECIFIED"
	WorkspacePermissionUserPermission WorkspacePermission = "USER_PERMISSION"
	WorkspacePermissionAdminPermission WorkspacePermission = "ADMIN_PERMISSION"
)

// String representation for [fmt.Print].
func (f *WorkspacePermission) String() string {
	return string(*f)
}





type AccountAccessIdentityRule struct {
	
	Action *AccountAccessRuleAction `json:"action"`
	
	ExternalId *string `json:"externalId"`
	
	DisplayName *string `json:"displayName"`
	
	PrincipalType *PrincipalType `json:"principalType"`
	
}



type CreateAccountAccessIdentityRuleRequest struct {
	
	AccountId *string `json:"accountId"`
	
	AccountAccessIdentityRule *AccountAccessIdentityRule `json:"accountAccessIdentityRule"`
	
}



type CreateGroupProxyRequest struct {
	
	Group *Group `json:"group"`
	
}



type CreateGroupRequest struct {
	
	AccountId *string `json:"accountId"`
	
	Group *Group `json:"group"`
	
}



type CreateServicePrincipalProxyRequest struct {
	
	ServicePrincipal *ServicePrincipal `json:"servicePrincipal"`
	
}



type CreateServicePrincipalRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ServicePrincipal *ServicePrincipal `json:"servicePrincipal"`
	
}



type CreateUserProxyRequest struct {
	
	User *User `json:"user"`
	
}



type CreateUserRequest struct {
	
	AccountId *string `json:"accountId"`
	
	User *User `json:"user"`
	
}



type CreateWorkspaceAccessDetailLocalRequest struct {
	
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
}



type CreateWorkspaceAccessDetailRequest struct {
	
	AccountId *string `json:"accountId"`
	
	Parent *string `json:"parent"`
	
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
}



type DeleteAccountAccessIdentityRuleRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ExternalId *string `json:"externalId"`
	
}



type DeleteGroupProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteGroupRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteServicePrincipalProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteServicePrincipalRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteUserProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteUserRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type DeleteWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId *int64 `json:"principalId"`
	
}



type DeleteWorkspaceAccessDetailRequest struct {
	
	AccountId *string `json:"accountId"`
	
	WorkspaceId *int64 `json:"workspaceId"`
	
	PrincipalId *int64 `json:"principalId"`
	
}



type GetAccountAccessIdentityRuleRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ExternalId *string `json:"externalId"`
	
}



type GetGroupProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type GetGroupRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type GetServicePrincipalProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type GetServicePrincipalRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type GetUserProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
}



type GetUserRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
}



type GetWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId *int64 `json:"principalId"`
	
	View *WorkspaceAccessDetailView `json:"view"`
	
}



type GetWorkspaceAccessDetailRequest struct {
	
	AccountId *string `json:"accountId"`
	
	WorkspaceId *int64 `json:"workspaceId"`
	
	PrincipalId *int64 `json:"principalId"`
	
	View *WorkspaceAccessDetailView `json:"view"`
	
}



type Group struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	ExternalId *string `json:"externalId"`
	
	GroupName *string `json:"groupName"`
	
}



type ListAccountAccessIdentityRulesRequest struct {
	
	AccountId *string `json:"accountId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListAccountAccessIdentityRulesResponse struct {
	
	AccountAccessIdentityRules []AccountAccessIdentityRule `json:"accountAccessIdentityRules"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListGroupsProxyRequest struct {
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListGroupsRequest struct {
	
	AccountId *string `json:"accountId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListGroupsResponse struct {
	
	Groups []Group `json:"groups"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListServicePrincipalsProxyRequest struct {
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListServicePrincipalsRequest struct {
	
	AccountId *string `json:"accountId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListServicePrincipalsResponse struct {
	
	ServicePrincipals []ServicePrincipal `json:"servicePrincipals"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListUsersProxyRequest struct {
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListUsersRequest struct {
	
	AccountId *string `json:"accountId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
	Filter *string `json:"filter"`
	
}



type ListUsersResponse struct {
	
	Users []User `json:"users"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type ListWorkspaceAccessDetailsLocalRequest struct {
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
}



type ListWorkspaceAccessDetailsRequest struct {
	
	AccountId *string `json:"accountId"`
	
	WorkspaceId *int64 `json:"workspaceId"`
	
	PageSize *int `json:"pageSize"`
	
	PageToken *string `json:"pageToken"`
	
}



type ListWorkspaceAccessDetailsResponse struct {
	
	WorkspaceAccessDetails []WorkspaceAccessDetail `json:"workspaceAccessDetails"`
	
	NextPageToken *string `json:"nextPageToken"`
	
}



type Name struct {
	
	GivenName *string `json:"givenName"`
	
	FamilyName *string `json:"familyName"`
	
}



type ResolveGroupProxyRequest struct {
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveGroupRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveGroupResponse struct {
	
	Group *Group `json:"group"`
	
}



type ResolveServicePrincipalProxyRequest struct {
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveServicePrincipalRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveServicePrincipalResponse struct {
	
	ServicePrincipal *ServicePrincipal `json:"servicePrincipal"`
	
}



type ResolveUserProxyRequest struct {
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveUserRequest struct {
	
	AccountId *string `json:"accountId"`
	
	ExternalId *string `json:"externalId"`
	
}



type ResolveUserResponse struct {
	
	User *User `json:"user"`
	
}



type ServicePrincipal struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	ExternalId *string `json:"externalId"`
	
	ApplicationId *string `json:"applicationId"`
	
	DisplayName *string `json:"displayName"`
	
	AccountSpStatus *State `json:"accountSpStatus"`
	
}



type UpdateGroupProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
	Group *Group `json:"group"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateGroupRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	Group *Group `json:"group"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateServicePrincipalProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
	ServicePrincipal *ServicePrincipal `json:"servicePrincipal"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateServicePrincipalRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	ServicePrincipal *ServicePrincipal `json:"servicePrincipal"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateUserProxyRequest struct {
	
	InternalId *int64 `json:"internalId"`
	
	User *User `json:"user"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateUserRequest struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	User *User `json:"user"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId *int64 `json:"principalId"`
	
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type UpdateWorkspaceAccessDetailRequest struct {
	
	AccountId *string `json:"accountId"`
	
	WorkspaceId *int64 `json:"workspaceId"`
	
	PrincipalId *int64 `json:"principalId"`
	
	WorkspaceAccessDetail *WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
	UpdateMask *string `json:"updateMask"`
	
}



type User struct {
	
	AccountId *string `json:"accountId"`
	
	InternalId *int64 `json:"internalId"`
	
	ExternalId *string `json:"externalId"`
	
	Username *string `json:"username"`
	
	Name *Name `json:"name"`
	
	AccountUserStatus *State `json:"accountUserStatus"`
	
}



type WorkspaceAccessDetail struct {
	
	PrincipalId *int64 `json:"principalId"`
	
	WorkspaceId *int64 `json:"workspaceId"`
	
	AccountId *string `json:"accountId"`
	
	PrincipalType *PrincipalType `json:"principalType"`
	
	AccessType *AccessType `json:"accessType"`
	
	Status *State `json:"status"`
	
	Permissions []WorkspacePermission `json:"permissions"`
	
}


