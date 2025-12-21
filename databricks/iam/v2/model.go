
package v2

import (
	
	"fmt"
	"io"
)



type AccessType string

const AccessTypeAccessTypeUnspecified AccessType = `ACCESS_TYPE_UNSPECIFIED`
const AccessTypeDirect AccessType = `DIRECT`
const AccessTypeIndirect AccessType = `INDIRECT`

// String representation for [fmt.Print]
func (f *AccessType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *AccessType) Set(v string) error {
	switch v {
	case `ACCESS_TYPE_UNSPECIFIED`, `DIRECT`, `INDIRECT`:
		*f = AccessType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ACCESS_TYPE_UNSPECIFIED", "DIRECT", "INDIRECT"`, v)
	}
}

// Values returns all possible values for AccessType.
//
// There is no guarantee on the order of the values in the slice.
func (f *AccessType) Values() []AccessType {
	return []AccessType{
		AccessTypeAccessTypeUnspecified,
		AccessTypeDirect,
		AccessTypeIndirect,
	}
}

// Type always returns AccessType to satisfy [pflag.Value] interface
func (f *AccessType) Type() string {
	return "AccessType"
}



type AccountAccessRuleAction string

const AccountAccessRuleActionAccountAccessRuleActionUnspecified AccountAccessRuleAction = `ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED`
const AccountAccessRuleActionAllow AccountAccessRuleAction = `ALLOW`
const AccountAccessRuleActionDeny AccountAccessRuleAction = `DENY`

// String representation for [fmt.Print]
func (f *AccountAccessRuleAction) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *AccountAccessRuleAction) Set(v string) error {
	switch v {
	case `ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED`, `ALLOW`, `DENY`:
		*f = AccountAccessRuleAction(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED", "ALLOW", "DENY"`, v)
	}
}

// Values returns all possible values for AccountAccessRuleAction.
//
// There is no guarantee on the order of the values in the slice.
func (f *AccountAccessRuleAction) Values() []AccountAccessRuleAction {
	return []AccountAccessRuleAction{
		AccountAccessRuleActionAccountAccessRuleActionUnspecified,
		AccountAccessRuleActionAllow,
		AccountAccessRuleActionDeny,
	}
}

// Type always returns AccountAccessRuleAction to satisfy [pflag.Value] interface
func (f *AccountAccessRuleAction) Type() string {
	return "AccountAccessRuleAction"
}



type PrincipalType string

const PrincipalTypePrincipalTypeUnspecified PrincipalType = `PRINCIPAL_TYPE_UNSPECIFIED`
const PrincipalTypeUser PrincipalType = `USER`
const PrincipalTypeServicePrincipal PrincipalType = `SERVICE_PRINCIPAL`
const PrincipalTypeGroup PrincipalType = `GROUP`

// String representation for [fmt.Print]
func (f *PrincipalType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *PrincipalType) Set(v string) error {
	switch v {
	case `PRINCIPAL_TYPE_UNSPECIFIED`, `USER`, `SERVICE_PRINCIPAL`, `GROUP`:
		*f = PrincipalType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "PRINCIPAL_TYPE_UNSPECIFIED", "USER", "SERVICE_PRINCIPAL", "GROUP"`, v)
	}
}

// Values returns all possible values for PrincipalType.
//
// There is no guarantee on the order of the values in the slice.
func (f *PrincipalType) Values() []PrincipalType {
	return []PrincipalType{
		PrincipalTypePrincipalTypeUnspecified,
		PrincipalTypeUser,
		PrincipalTypeServicePrincipal,
		PrincipalTypeGroup,
	}
}

// Type always returns PrincipalType to satisfy [pflag.Value] interface
func (f *PrincipalType) Type() string {
	return "PrincipalType"
}



type State string

const StateStateUnspecified State = `STATE_UNSPECIFIED`
const StateActive State = `ACTIVE`
const StateInactive State = `INACTIVE`

// String representation for [fmt.Print]
func (f *State) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *State) Set(v string) error {
	switch v {
	case `STATE_UNSPECIFIED`, `ACTIVE`, `INACTIVE`:
		*f = State(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "STATE_UNSPECIFIED", "ACTIVE", "INACTIVE"`, v)
	}
}

// Values returns all possible values for State.
//
// There is no guarantee on the order of the values in the slice.
func (f *State) Values() []State {
	return []State{
		StateStateUnspecified,
		StateActive,
		StateInactive,
	}
}

// Type always returns State to satisfy [pflag.Value] interface
func (f *State) Type() string {
	return "State"
}



type WorkspaceAccessDetailView string

const WorkspaceAccessDetailViewWorkspaceAccessDetailViewUnspecified WorkspaceAccessDetailView = `WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED`
const WorkspaceAccessDetailViewBasic WorkspaceAccessDetailView = `BASIC`
const WorkspaceAccessDetailViewFull WorkspaceAccessDetailView = `FULL`

// String representation for [fmt.Print]
func (f *WorkspaceAccessDetailView) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *WorkspaceAccessDetailView) Set(v string) error {
	switch v {
	case `WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED`, `BASIC`, `FULL`:
		*f = WorkspaceAccessDetailView(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED", "BASIC", "FULL"`, v)
	}
}

// Values returns all possible values for WorkspaceAccessDetailView.
//
// There is no guarantee on the order of the values in the slice.
func (f *WorkspaceAccessDetailView) Values() []WorkspaceAccessDetailView {
	return []WorkspaceAccessDetailView{
		WorkspaceAccessDetailViewWorkspaceAccessDetailViewUnspecified,
		WorkspaceAccessDetailViewBasic,
		WorkspaceAccessDetailViewFull,
	}
}

// Type always returns WorkspaceAccessDetailView to satisfy [pflag.Value] interface
func (f *WorkspaceAccessDetailView) Type() string {
	return "WorkspaceAccessDetailView"
}



type WorkspacePermission string

const WorkspacePermissionWorkspacePermissionUnspecified WorkspacePermission = `WORKSPACE_PERMISSION_UNSPECIFIED`
const WorkspacePermissionUserPermission WorkspacePermission = `USER_PERMISSION`
const WorkspacePermissionAdminPermission WorkspacePermission = `ADMIN_PERMISSION`

// String representation for [fmt.Print]
func (f *WorkspacePermission) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *WorkspacePermission) Set(v string) error {
	switch v {
	case `WORKSPACE_PERMISSION_UNSPECIFIED`, `USER_PERMISSION`, `ADMIN_PERMISSION`:
		*f = WorkspacePermission(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "WORKSPACE_PERMISSION_UNSPECIFIED", "USER_PERMISSION", "ADMIN_PERMISSION"`, v)
	}
}

// Values returns all possible values for WorkspacePermission.
//
// There is no guarantee on the order of the values in the slice.
func (f *WorkspacePermission) Values() []WorkspacePermission {
	return []WorkspacePermission{
		WorkspacePermissionWorkspacePermissionUnspecified,
		WorkspacePermissionUserPermission,
		WorkspacePermissionAdminPermission,
	}
}

// Type always returns WorkspacePermission to satisfy [pflag.Value] interface
func (f *WorkspacePermission) Type() string {
	return "WorkspacePermission"
}





type AccountAccessIdentityRule struct {
	
	Action AccountAccessRuleAction `json:"action"`
	
	ExternalId string `json:"externalId"`
	
	DisplayName string `json:"displayName"`
	
	PrincipalType PrincipalType `json:"principalType"`
	
}



type CreateAccountAccessIdentityRuleRequest struct {
	
	AccountId string `json:"accountId"`
	
	AccountAccessIdentityRule AccountAccessIdentityRule `json:"accountAccessIdentityRule"`
	
}



type CreateGroupProxyRequest struct {
	
	Group Group `json:"group"`
	
}



type CreateGroupRequest struct {
	
	AccountId string `json:"accountId"`
	
	Group Group `json:"group"`
	
}



type CreateServicePrincipalProxyRequest struct {
	
	ServicePrincipal ServicePrincipal `json:"servicePrincipal"`
	
}



type CreateServicePrincipalRequest struct {
	
	AccountId string `json:"accountId"`
	
	ServicePrincipal ServicePrincipal `json:"servicePrincipal"`
	
}



type CreateUserProxyRequest struct {
	
	User User `json:"user"`
	
}



type CreateUserRequest struct {
	
	AccountId string `json:"accountId"`
	
	User User `json:"user"`
	
}



type CreateWorkspaceAccessDetailLocalRequest struct {
	
	WorkspaceAccessDetail WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
}



type CreateWorkspaceAccessDetailRequest struct {
	
	AccountId string `json:"accountId"`
	
	Parent string `json:"parent"`
	
	WorkspaceAccessDetail WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
}



type DeleteAccountAccessIdentityRuleRequest struct {
	
	AccountId string `json:"accountId"`
	
	ExternalId string `json:"externalId"`
	
}



type DeleteGroupProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteGroupRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteServicePrincipalProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteServicePrincipalRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteUserProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteUserRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type DeleteWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId int64 `json:"principalId"`
	
}



type DeleteWorkspaceAccessDetailRequest struct {
	
	AccountId string `json:"accountId"`
	
	WorkspaceId int64 `json:"workspaceId"`
	
	PrincipalId int64 `json:"principalId"`
	
}



type GetGroupProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type GetGroupRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type GetServicePrincipalProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type GetServicePrincipalRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type GetUserProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
}



type GetUserRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
}



type GetWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId int64 `json:"principalId"`
	
	View WorkspaceAccessDetailView `json:"view"`
	
}



type GetWorkspaceAccessDetailRequest struct {
	
	AccountId string `json:"accountId"`
	
	WorkspaceId int64 `json:"workspaceId"`
	
	PrincipalId int64 `json:"principalId"`
	
	View WorkspaceAccessDetailView `json:"view"`
	
}



type Group struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	ExternalId string `json:"externalId"`
	
	GroupName string `json:"groupName"`
	
}



type ListAccountAccessIdentityRulesRequest struct {
	
	AccountId string `json:"accountId"`
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListAccountAccessIdentityRulesResponse struct {
	
	AccountAccessIdentityRules []AccountAccessIdentityRule `json:"accountAccessIdentityRules"`
	
	NextPageToken string `json:"nextPageToken"`
	
}



type ListGroupsProxyRequest struct {
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListGroupsRequest struct {
	
	AccountId string `json:"accountId"`
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListGroupsResponse struct {
	
	Groups []Group `json:"groups"`
	
	NextPageToken string `json:"nextPageToken"`
	
}



type ListServicePrincipalsProxyRequest struct {
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListServicePrincipalsRequest struct {
	
	AccountId string `json:"accountId"`
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListServicePrincipalsResponse struct {
	
	ServicePrincipals []ServicePrincipal `json:"servicePrincipals"`
	
	NextPageToken string `json:"nextPageToken"`
	
}



type ListUsersProxyRequest struct {
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListUsersRequest struct {
	
	AccountId string `json:"accountId"`
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
	Filter string `json:"filter"`
	
}



type ListUsersResponse struct {
	
	Users []User `json:"users"`
	
	NextPageToken string `json:"nextPageToken"`
	
}



type ListWorkspaceAccessDetailsLocalRequest struct {
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
}



type ListWorkspaceAccessDetailsRequest struct {
	
	AccountId string `json:"accountId"`
	
	WorkspaceId int64 `json:"workspaceId"`
	
	PageSize int `json:"pageSize"`
	
	PageToken string `json:"pageToken"`
	
}



type ListWorkspaceAccessDetailsResponse struct {
	
	WorkspaceAccessDetails []WorkspaceAccessDetail `json:"workspaceAccessDetails"`
	
	NextPageToken string `json:"nextPageToken"`
	
}



type Name struct {
	
	GivenName string `json:"givenName"`
	
	FamilyName string `json:"familyName"`
	
}



type ResolveGroupProxyRequest struct {
	
	ExternalId string `json:"externalId"`
	
}



type ResolveGroupRequest struct {
	
	AccountId string `json:"accountId"`
	
	ExternalId string `json:"externalId"`
	
}



type ResolveGroupResponse struct {
	
	Group Group `json:"group"`
	
}



type ResolveServicePrincipalProxyRequest struct {
	
	ExternalId string `json:"externalId"`
	
}



type ResolveServicePrincipalRequest struct {
	
	AccountId string `json:"accountId"`
	
	ExternalId string `json:"externalId"`
	
}



type ResolveServicePrincipalResponse struct {
	
	ServicePrincipal ServicePrincipal `json:"servicePrincipal"`
	
}



type ResolveUserProxyRequest struct {
	
	ExternalId string `json:"externalId"`
	
}



type ResolveUserRequest struct {
	
	AccountId string `json:"accountId"`
	
	ExternalId string `json:"externalId"`
	
}



type ResolveUserResponse struct {
	
	User User `json:"user"`
	
}



type ServicePrincipal struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	ExternalId string `json:"externalId"`
	
	ApplicationId string `json:"applicationId"`
	
	DisplayName string `json:"displayName"`
	
	AccountSpStatus State `json:"accountSpStatus"`
	
}



type UpdateGroupProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
	Group Group `json:"group"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateGroupRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	Group Group `json:"group"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateServicePrincipalProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
	ServicePrincipal ServicePrincipal `json:"servicePrincipal"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateServicePrincipalRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	ServicePrincipal ServicePrincipal `json:"servicePrincipal"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateUserProxyRequest struct {
	
	InternalId int64 `json:"internalId"`
	
	User User `json:"user"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateUserRequest struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	User User `json:"user"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateWorkspaceAccessDetailLocalRequest struct {
	
	PrincipalId int64 `json:"principalId"`
	
	WorkspaceAccessDetail WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
	UpdateMask string `json:"updateMask"`
	
}



type UpdateWorkspaceAccessDetailRequest struct {
	
	AccountId string `json:"accountId"`
	
	WorkspaceId int64 `json:"workspaceId"`
	
	PrincipalId int64 `json:"principalId"`
	
	WorkspaceAccessDetail WorkspaceAccessDetail `json:"workspaceAccessDetail"`
	
	UpdateMask string `json:"updateMask"`
	
}



type User struct {
	
	AccountId string `json:"accountId"`
	
	InternalId int64 `json:"internalId"`
	
	ExternalId string `json:"externalId"`
	
	Username string `json:"username"`
	
	Name Name `json:"name"`
	
	AccountUserStatus State `json:"accountUserStatus"`
	
}



type WorkspaceAccessDetail struct {
	
	PrincipalId int64 `json:"principalId"`
	
	WorkspaceId int64 `json:"workspaceId"`
	
	AccountId string `json:"accountId"`
	
	PrincipalType PrincipalType `json:"principalType"`
	
	AccessType AccessType `json:"accessType"`
	
	Status State `json:"status"`
	
	Permissions []WorkspacePermission `json:"permissions"`
	
}


