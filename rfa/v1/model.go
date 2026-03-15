package v1

type DestinationType string

const (
	DestinationType_DestinationTypeUnspecified DestinationType = "DESTINATION_TYPE_UNSPECIFIED"
	DestinationType_Email                      DestinationType = "EMAIL"
	DestinationType_Slack                      DestinationType = "SLACK"
	DestinationType_GenericWebhook             DestinationType = "GENERIC_WEBHOOK"
	DestinationType_MicrosoftTeams             DestinationType = "MICROSOFT_TEAMS"
	DestinationType_Url                        DestinationType = "URL"
)

// String representation for [fmt.Print].
func (f *DestinationType) String() string {
	return string(*f)
}

type PrincipalType string

const (
	PrincipalType_PrincipalTypeUnspecified PrincipalType = "PRINCIPAL_TYPE_UNSPECIFIED"
	PrincipalType_UserPrincipal            PrincipalType = "USER_PRINCIPAL"
	PrincipalType_GroupPrincipal           PrincipalType = "GROUP_PRINCIPAL"
	PrincipalType_ServicePrincipal         PrincipalType = "SERVICE_PRINCIPAL"
)

// String representation for [fmt.Print].
func (f *PrincipalType) String() string {
	return string(*f)
}

type SecurableType string

const (
	SecurableType_Catalog           SecurableType = "CATALOG"
	SecurableType_Schema            SecurableType = "SCHEMA"
	SecurableType_Table             SecurableType = "TABLE"
	SecurableType_StorageCredential SecurableType = "STORAGE_CREDENTIAL"
	SecurableType_ExternalLocation  SecurableType = "EXTERNAL_LOCATION"
	SecurableType_Function          SecurableType = "FUNCTION"
	SecurableType_Share             SecurableType = "SHARE"
	SecurableType_Provider          SecurableType = "PROVIDER"
	SecurableType_Recipient         SecurableType = "RECIPIENT"
	SecurableType_CleanRoom         SecurableType = "CLEAN_ROOM"
	SecurableType_Metastore         SecurableType = "METASTORE"
	SecurableType_Pipeline          SecurableType = "PIPELINE"
	SecurableType_Volume            SecurableType = "VOLUME"
	SecurableType_Connection        SecurableType = "CONNECTION"
	SecurableType_Credential        SecurableType = "CREDENTIAL"
	SecurableType_ExternalMetadata  SecurableType = "EXTERNAL_METADATA"
	SecurableType_StagingTable      SecurableType = "STAGING_TABLE"
)

// String representation for [fmt.Print].
func (f *SecurableType) String() string {
	return string(*f)
}

type SpecialDestination string

const (
	SpecialDestination_SpecialDestinationUnspecified           SpecialDestination = "SPECIAL_DESTINATION_UNSPECIFIED"
	SpecialDestination_SpecialDestinationCatalogOwner          SpecialDestination = "SPECIAL_DESTINATION_CATALOG_OWNER"
	SpecialDestination_SpecialDestinationExternalLocationOwner SpecialDestination = "SPECIAL_DESTINATION_EXTERNAL_LOCATION_OWNER"
	SpecialDestination_SpecialDestinationConnectionOwner       SpecialDestination = "SPECIAL_DESTINATION_CONNECTION_OWNER"
	SpecialDestination_SpecialDestinationCredentialOwner       SpecialDestination = "SPECIAL_DESTINATION_CREDENTIAL_OWNER"
	SpecialDestination_SpecialDestinationMetastoreOwner        SpecialDestination = "SPECIAL_DESTINATION_METASTORE_OWNER"
)

// String representation for [fmt.Print].
func (f *SpecialDestination) String() string {
	return string(*f)
}

type AccessRequestDestinations struct {
	Destinations               []NotificationDestination `json:"destinations"`
	Securable                  *Securable                `json:"securable"`
	AreAnyDestinationsHidden   *bool                     `json:"are_any_destinations_hidden"`
	DestinationSourceSecurable *Securable                `json:"destination_source_securable"`
	SecurableType              *string                   `json:"securable_type"`
	FullName                   *string                   `json:"full_name"`
}

type BatchCreateAccessRequestsRequest struct {
	Requests []CreateAccessRequest `json:"requests"`
}

type BatchCreateAccessRequestsResponse struct {
	Responses []CreateAccessRequestResponse `json:"responses"`
}

type CreateAccessRequest struct {
	BehalfOf             *Principal             `json:"behalf_of"`
	Comment              *string                `json:"comment"`
	SecurablePermissions []SecurablePermissions `json:"securable_permissions"`
}

type CreateAccessRequestResponse struct {
	BehalfOf            *Principal                  `json:"behalf_of"`
	RequestDestinations []AccessRequestDestinations `json:"request_destinations"`
}

type GetAccessRequestDestinationsRequest struct {
	SecurableType *string `json:"securable_type"`
	FullName      *string `json:"full_name"`
}

type NotificationDestination struct {
	DestinationId      *string             `json:"destination_id"`
	DestinationType    *DestinationType    `json:"destination_type"`
	SpecialDestination *SpecialDestination `json:"special_destination"`
}

type Principal struct {
	Id            *string        `json:"id"`
	PrincipalType *PrincipalType `json:"principal_type"`
}

// Generic definition of a securable, which is uniquely defined in a metastore
// by its type and full name..
type Securable struct {
	Type          *SecurableType `json:"type"`
	FullName      *string        `json:"full_name"`
	ProviderShare *string        `json:"provider_share"`
}

type SecurablePermissions struct {
	Securable   *Securable `json:"securable"`
	Permissions []string   `json:"permissions"`
}

type UpdateAccessRequestDestinationsRequest struct {
	AccessRequestDestinations *AccessRequestDestinations `json:"access_request_destinations"`
	UpdateMask                *string                    `json:"update_mask"`
}
