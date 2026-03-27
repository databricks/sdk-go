package v2

type AccountAccessRuleAction string

const (
	AccountAccessRuleAction_AccountAccessRuleActionUnspecified AccountAccessRuleAction = "ACCOUNT_ACCESS_RULE_ACTION_UNSPECIFIED"
	AccountAccessRuleAction_Deny                               AccountAccessRuleAction = "DENY"
)

// String representation for [fmt.Print].
func (f *AccountAccessRuleAction) String() string {
	return string(*f)
}

type Entitlement string

const (
	Entitlement_EntitlementUnspecified  Entitlement = "ENTITLEMENT_UNSPECIFIED"
	Entitlement_WorkspaceAccess         Entitlement = "WORKSPACE_ACCESS"
	Entitlement_WorkspaceConsume        Entitlement = "WORKSPACE_CONSUME"
	Entitlement_DatabricksSqlAccess     Entitlement = "DATABRICKS_SQL_ACCESS"
	Entitlement_WorkspaceAdmin          Entitlement = "WORKSPACE_ADMIN"
	Entitlement_AllowClusterCreate      Entitlement = "ALLOW_CLUSTER_CREATE"
	Entitlement_AllowInstancePoolCreate Entitlement = "ALLOW_INSTANCE_POOL_CREATE"
)

// String representation for [fmt.Print].
func (f *Entitlement) String() string {
	return string(*f)
}

type GroupMembershipSource string

const (
	GroupMembershipSource_GroupMembershipSourceUnspecified GroupMembershipSource = "GROUP_MEMBERSHIP_SOURCE_UNSPECIFIED"
	GroupMembershipSource_Internal                         GroupMembershipSource = "INTERNAL"
	GroupMembershipSource_IdentityProvider                 GroupMembershipSource = "IDENTITY_PROVIDER"
)

// String representation for [fmt.Print].
func (f *GroupMembershipSource) String() string {
	return string(*f)
}

type PrincipalType string

const (
	PrincipalType_PrincipalTypeUnspecified PrincipalType = "PRINCIPAL_TYPE_UNSPECIFIED"
	PrincipalType_User                     PrincipalType = "USER"
	PrincipalType_ServicePrincipal         PrincipalType = "SERVICE_PRINCIPAL"
	PrincipalType_Group                    PrincipalType = "GROUP"
)

// String representation for [fmt.Print].
func (f *PrincipalType) String() string {
	return string(*f)
}

type State string

const (
	State_StateUnspecified State = "STATE_UNSPECIFIED"
	State_Active           State = "ACTIVE"
	State_Inactive         State = "INACTIVE"
)

// String representation for [fmt.Print].
func (f *State) String() string {
	return string(*f)
}

type WorkspaceAccessDetailView string

const (
	WorkspaceAccessDetailView_WorkspaceAccessDetailViewUnspecified WorkspaceAccessDetailView = "WORKSPACE_ACCESS_DETAIL_VIEW_UNSPECIFIED"
	WorkspaceAccessDetailView_Basic                                WorkspaceAccessDetailView = "BASIC"
	WorkspaceAccessDetailView_Full                                 WorkspaceAccessDetailView = "FULL"
)

// String representation for [fmt.Print].
func (f *WorkspaceAccessDetailView) String() string {
	return string(*f)
}

type WorkspacePermission string

const (
	WorkspacePermission_WorkspacePermissionUnspecified WorkspacePermission = "WORKSPACE_PERMISSION_UNSPECIFIED"
	WorkspacePermission_UserPermission                 WorkspacePermission = "USER_PERMISSION"
	WorkspacePermission_AdminPermission                WorkspacePermission = "ADMIN_PERMISSION"
)

// String representation for [fmt.Print].
func (f *WorkspacePermission) String() string {
	return string(*f)
}

type WorkspaceAccessDetail_AccessType string

const (
	WorkspaceAccessDetail_AccessType_AccessTypeUnspecified WorkspaceAccessDetail_AccessType = "ACCESS_TYPE_UNSPECIFIED"
	WorkspaceAccessDetail_AccessType_Direct                WorkspaceAccessDetail_AccessType = "DIRECT"
	WorkspaceAccessDetail_AccessType_Indirect              WorkspaceAccessDetail_AccessType = "INDIRECT"
)

// String representation for [fmt.Print].
func (f *WorkspaceAccessDetail_AccessType) String() string {
	return string(*f)
}

// An identity rule that controls which principals can access an account..
type AccountAccessIdentityRule struct {
	Action              *AccountAccessRuleAction `json:"action"`
	ExternalPrincipalId *string                  `json:"external_principal_id"`
	DisplayName         *string                  `json:"display_name"`
	PrincipalType       *PrincipalType           `json:"principal_type"`
	Name                *string                  `json:"name"`
}

// Request message for creating an account access identity rule..
type CreateAccountAccessIdentityRuleRequest struct {
	Parent                    *string                    `json:"parent"`
	ExternalPrincipalId       *string                    `json:"external_principal_id"`
	AccountAccessIdentityRule *AccountAccessIdentityRule `json:"account_access_identity_rule"`
}

// Request message for creating a group membership (assigning a principal to a
// group)..
type CreateDirectGroupMemberProxyRequest struct {
	GroupId           *int64             `json:"group_id"`
	DirectGroupMember *DirectGroupMember `json:"direct_group_member"`
}

// Request message for creating a group membership (assigning a principal to a
// group)..
type CreateDirectGroupMemberRequest struct {
	AccountId         *string            `json:"account_id"`
	GroupId           *int64             `json:"group_id"`
	DirectGroupMember *DirectGroupMember `json:"direct_group_member"`
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

// Creates a user in Databricks and provisions it at the account level. Behavior
// depends on whether Account Identity Management (AIM) is enabled: - When AIM
// is enabled: The user is provisioned with an internalId. If an externalId is
// provided, the identity provider is treated as the source of truth for user
// metadata, and customer-supplied field values may be overridden. - When AIM is
// disabled: The user is provisioned with an internalId only, and
// customer-supplied metadata is used as-is..
type CreateUserProxyRequest struct {
	User *User `json:"user"`
}

// Creates a user in Databricks and provisions it at the account level. Behavior
// depends on whether Account Identity Management (AIM) is enabled: - When AIM
// is enabled: The user is provisioned with an internalId. If an externalId is
// provided, the identity provider is treated as the source of truth for user
// metadata, and customer-supplied field values may be overridden. - When AIM is
// disabled: The user is provisioned with an internalId only, and
// customer-supplied metadata is used as-is..
type CreateUserRequest struct {
	AccountId *string `json:"account_id"`
	User      *User   `json:"user"`
}

// Assign an identity directly to a workspace with the specified permissions and
// workspace-level status..
type CreateWorkspaceAssignmentDetailProxyRequest struct {
	WorkspaceAssignmentDetail *WorkspaceAssignmentDetail `json:"workspace_assignment_detail"`
}

// Assign an identity directly to a workspace with the specified permissions and
// workspace-level status..
type CreateWorkspaceAssignmentDetailRequest struct {
	AccountId                 *string                    `json:"account_id"`
	WorkspaceId               *int64                     `json:"workspace_id"`
	WorkspaceAssignmentDetail *WorkspaceAssignmentDetail `json:"workspace_assignment_detail"`
}

// Request message for deleting an account access identity rule..
type DeleteAccountAccessIdentityRuleRequest struct {
	Parent              *string `json:"parent"`
	ExternalPrincipalId *string `json:"external_principal_id"`
}

// Request message for deleting a group membership (unassigning a principal from
// a group)..
type DeleteDirectGroupMemberProxyRequest struct {
	GroupId     *int64 `json:"group_id"`
	PrincipalId *int64 `json:"principal_id"`
}

// Request message for deleting a group membership (unassigning a principal from
// a group)..
type DeleteDirectGroupMemberRequest struct {
	AccountId   *string `json:"account_id"`
	GroupId     *int64  `json:"group_id"`
	PrincipalId *int64  `json:"principal_id"`
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

// Proxy request for deleting a workspace assignment detail for a principal..
type DeleteWorkspaceAssignmentDetailProxyRequest struct {
	PrincipalId *int64 `json:"principal_id"`
}

// If the identity is directly assigned to the workspace, remove its assignment
// from the workspace.
type DeleteWorkspaceAssignmentDetailRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PrincipalId *int64  `json:"principal_id"`
}

// Represents a principal that is a direct member of a group, with its source of
// membership..
type DirectGroupMember struct {
	PrincipalId      *int64                 `json:"principal_id"`
	PrincipalType    *PrincipalType         `json:"principal_type"`
	MembershipSource *GroupMembershipSource `json:"membership_source"`
	DisplayName      *string                `json:"display_name"`
	ExternalId       *string                `json:"external_id"`
}

// Request message for getting an account access identity rule..
type GetAccountAccessIdentityRuleRequest struct {
	Parent              *string `json:"parent"`
	ExternalPrincipalId *string `json:"external_principal_id"`
}

// Request message for getting a provisioned direct group member..
type GetDirectGroupMemberProxyRequest struct {
	GroupId     *int64 `json:"group_id"`
	PrincipalId *int64 `json:"principal_id"`
}

// Request message for getting a provisioned direct group member..
type GetDirectGroupMemberRequest struct {
	AccountId   *string `json:"account_id"`
	GroupId     *int64  `json:"group_id"`
	PrincipalId *int64  `json:"principal_id"`
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

// Creates a user in Databricks and returns the resulting User resource.
// Readability of the created user depends on Account Identity Management (AIM)
// and the configured Boundary Enforcement mode: - When AIM is enabled and
// Boundary Enforcement is set to RULES_ONLY: - MVP: Any user with an internalId
// is readable, including users with an externalId populated. - Phase 2:
// Behavior to be defined. - When AIM is enabled and Boundary Enforcement is set
// to ALLOW_ALL: - Any user with an internalId is readable, including users with
// an externalId populated. - When AIM is disabled: - Returns the User resource
// corresponding to the given internalId..
type GetUserProxyRequest struct {
	InternalId *int64 `json:"internal_id"`
}

// Creates a user in Databricks and returns the resulting User resource.
// Readability of the created user depends on Account Identity Management (AIM)
// and the configured Boundary Enforcement mode: - When AIM is enabled and
// Boundary Enforcement is set to RULES_ONLY: - MVP: Any user with an internalId
// is readable, including users with an externalId populated. - Phase 2:
// Behavior to be defined. - When AIM is enabled and Boundary Enforcement is set
// to ALLOW_ALL: - Any user with an internalId is readable, including users with
// an externalId populated. - When AIM is disabled: - Returns the User resource
// corresponding to the given internalId..
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

// Proxy request for getting workspace assignment details for a principal in a
// workspace..
type GetWorkspaceAssignmentDetailProxyRequest struct {
	PrincipalId *int64 `json:"principal_id"`
}

// Get the workspace assignment details of a principal that is provisioned in
// the account and directly assigned to a workspace.
type GetWorkspaceAssignmentDetailRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PrincipalId *int64  `json:"principal_id"`
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
	Parent    *string `json:"parent"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// Response message for listing account access identity rules..
type ListAccountAccessIdentityRulesResponse struct {
	AccountAccessIdentityRules []AccountAccessIdentityRule `json:"account_access_identity_rules"`
	NextPageToken              *string                     `json:"next_page_token"`
}

// Request message for listing provisioned direct group members..
type ListDirectGroupMembersProxyRequest struct {
	GroupId   *int64  `json:"group_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// Request message for listing provisioned direct group members..
type ListDirectGroupMembersRequest struct {
	AccountId *string `json:"account_id"`
	GroupId   *int64  `json:"group_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// Response message for listing direct group members..
type ListDirectGroupMembersResponse struct {
	DirectGroupMembers []DirectGroupMember `json:"direct_group_members"`
	NextPageToken      *string             `json:"next_page_token"`
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

// Request message for listing all transitive parent groups of a principal..
type ListTransitiveParentGroupsProxyRequest struct {
	PrincipalId *int64  `json:"principal_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

// Request message for listing all transitive parent groups of a principal..
type ListTransitiveParentGroupsRequest struct {
	AccountId   *string `json:"account_id"`
	PrincipalId *int64  `json:"principal_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

// Response message for listing all transitive parent groups of a principal..
type ListTransitiveParentGroupsResponse struct {
	TransitiveParentGroups []TransitiveParentGroup `json:"transitive_parent_groups"`
	NextPageToken          *string                 `json:"next_page_token"`
}

// Returns a paginated list of account-level users. Behavior depends on whether
// Account Identity Management (AIM) is enabled: - When AIM is enabled: - The
// "externalId eq" filter only evaluates provisioned Databricks users that have
// an internalId. - The "username eq" filter only evaluates provisioned
// Databricks users that have an internalId. - Listing without filters returns
// all provisioned Databricks users. - AIM Boundary Enforcement Phase 2:
// Behavior to be defined. - When AIM is disabled: - The "externalId eq" filter
// only evaluates provisioned Databricks users that have an internalId. - The
// "username eq" filter only evaluates provisioned Databricks users that have an
// internalId. - Listing without filters returns all provisioned Databricks
// users..
type ListUsersProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
	Filter    *string `json:"filter"`
}

// Returns a paginated list of account-level users. Behavior depends on whether
// Account Identity Management (AIM) is enabled: - When AIM is enabled: - The
// "externalId eq" filter only evaluates provisioned Databricks users that have
// an internalId. - The "username eq" filter only evaluates provisioned
// Databricks users that have an internalId. - Listing without filters returns
// all provisioned Databricks users. - AIM Boundary Enforcement Phase 2:
// Behavior to be defined. - When AIM is disabled: - The "externalId eq" filter
// only evaluates provisioned Databricks users that have an internalId. - The
// "username eq" filter only evaluates provisioned Databricks users that have an
// internalId. - Listing without filters returns all provisioned Databricks
// users..
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

// Proxy request for listing workspace assignment details for a workspace..
type ListWorkspaceAssignmentDetailsProxyRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// Returns a paginated list of direct assignments to the workspace..
type ListWorkspaceAssignmentDetailsRequest struct {
	AccountId   *string `json:"account_id"`
	WorkspaceId *int64  `json:"workspace_id"`
	PageSize    *int    `json:"page_size"`
	PageToken   *string `json:"page_token"`
}

// Response message for listing workspace assignment details..
type ListWorkspaceAssignmentDetailsResponse struct {
	WorkspaceAssignmentDetails []WorkspaceAssignmentDetail `json:"workspace_assignment_details"`
	NextPageToken              *string                     `json:"next_page_token"`
}

// Request message for matching a group against the IDP. This will perform a
// sync by group_id before performing analysis to update local data which is
// safe to fix..
type MatchGroupWithIdpRequest struct {
	AccountId *string `json:"account_id"`
	GroupId   *int64  `json:"group_id"`
}

// Response message for matching a group against the IDP..
type MatchGroupWithIdpResponse struct {
	DatabricksGroup         *Group              `json:"databricks_group"`
	IdpMatchesByGroupName   []Group             `json:"idp_matches_by_group_name"`
	IdpMatchByExternalId    *Group              `json:"idp_match_by_external_id"`
	LocalMembersNotInIdp    []DirectGroupMember `json:"local_members_not_in_idp"`
	ExternalMembersNotInIdp []DirectGroupMember `json:"external_members_not_in_idp"`
}

// Request message for matching a service principal against the IDP. This will
// perform a sync by service_principal_id before performing analysis to update
// local data which is safe to fix..
type MatchServicePrincipalWithIdpRequest struct {
	AccountId          *string `json:"account_id"`
	ServicePrincipalId *int64  `json:"service_principal_id"`
}

// Response message for matching a service principal against the IDP..
type MatchServicePrincipalWithIdpResponse struct {
	DatabricksServicePrincipal *ServicePrincipal `json:"databricks_service_principal"`
	IdpMatchByAppId            *ServicePrincipal `json:"idp_match_by_app_id"`
	IdpMatchByExternalId       *ServicePrincipal `json:"idp_match_by_external_id"`
}

// Request message for matching a user against the IDP. This will perform a sync
// by user_id before performing analysis to update local data which is safe to
// fix..
type MatchUserWithIdpRequest struct {
	AccountId *string `json:"account_id"`
	UserId    *int64  `json:"user_id"`
}

// Response message for matching a user against the IDP..
type MatchUserWithIdpResponse struct {
	DatabricksUser       *User `json:"databricks_user"`
	IdpMatchByUsername   *User `json:"idp_match_by_username"`
	IdpMatchByExternalId *User `json:"idp_match_by_external_id"`
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

// Represents a group that is a transitive parent of a principal..
type TransitiveParentGroup struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	ExternalId *string `json:"external_id"`
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

// Updates an existing user in Databricks. The behavior is consistent regardless
// of whether Account Identity Management (AIM) is enabled or disabled. The
// following fields are updatable: - name.familyName - name.givenName - status -
// externalId.
type UpdateUserProxyRequest struct {
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

// Updates an existing user in Databricks. The behavior is consistent regardless
// of whether Account Identity Management (AIM) is enabled or disabled. The
// following fields are updatable: - name.familyName - name.givenName - status -
// externalId.
type UpdateUserRequest struct {
	AccountId  *string `json:"account_id"`
	InternalId *int64  `json:"internal_id"`
	User       *User   `json:"user"`
	UpdateMask *string `json:"update_mask"`
}

// Proxy request for updating a workspace assignment detail for a principal..
type UpdateWorkspaceAssignmentDetailProxyRequest struct {
	PrincipalId               *int64                     `json:"principal_id"`
	WorkspaceAssignmentDetail *WorkspaceAssignmentDetail `json:"workspace_assignment_detail"`
	UpdateMask                *string                    `json:"update_mask"`
}

// TBD since the only updatable field is permissions.
type UpdateWorkspaceAssignmentDetailRequest struct {
	AccountId                 *string                    `json:"account_id"`
	WorkspaceId               *int64                     `json:"workspace_id"`
	PrincipalId               *int64                     `json:"principal_id"`
	WorkspaceAssignmentDetail *WorkspaceAssignmentDetail `json:"workspace_assignment_detail"`
	UpdateMask                *string                    `json:"update_mask"`
}

// The details of a User resource..
type User struct {
	AccountId         *string    `json:"account_id"`
	InternalId        *int64     `json:"internal_id"`
	ExternalId        *string    `json:"external_id"`
	Username          *string    `json:"username"`
	Name              *User_Name `json:"name"`
	AccountUserStatus *State     `json:"account_user_status"`
}

type User_Name struct {
	GivenName  *string `json:"given_name"`
	FamilyName *string `json:"family_name"`
}

// The details of a principal's access to a workspace..
type WorkspaceAccessDetail struct {
	PrincipalId   *int64                            `json:"principal_id"`
	WorkspaceId   *int64                            `json:"workspace_id"`
	AccountId     *string                           `json:"account_id"`
	PrincipalType *PrincipalType                    `json:"principal_type"`
	AccessType    *WorkspaceAccessDetail_AccessType `json:"access_type"`
	Status        *State                            `json:"status"`
	Permissions   []WorkspacePermission             `json:"permissions"`
}

// The details of a principal's assignment to a workspace..
type WorkspaceAssignmentDetail struct {
	PrincipalId   *int64         `json:"principal_id"`
	WorkspaceId   *int64         `json:"workspace_id"`
	AccountId     *string        `json:"account_id"`
	PrincipalType *PrincipalType `json:"principal_type"`
	Entitlements  []Entitlement  `json:"entitlements"`
}
