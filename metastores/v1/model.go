package v1

type DeltaSharingScope_Enum string

const (
	DeltaSharingScope_Enum_Internal            DeltaSharingScope_Enum = "INTERNAL"
	DeltaSharingScope_Enum_InternalAndExternal DeltaSharingScope_Enum = "INTERNAL_AND_EXTERNAL"
)

// String representation for [fmt.Print].
func (f *DeltaSharingScope_Enum) String() string {
	return string(*f)
}

type CreateMetastore struct {
	Name                                        *string                 `json:"name"`
	StorageRoot                                 *string                 `json:"storage_root"`
	DefaultDataAccessConfigId                   *string                 `json:"default_data_access_config_id"`
	StorageRootCredentialId                     *string                 `json:"storage_root_credential_id"`
	DeltaSharingScope                           *DeltaSharingScope_Enum `json:"delta_sharing_scope"`
	DeltaSharingRecipientTokenLifetimeInSeconds *int64                  `json:"delta_sharing_recipient_token_lifetime_in_seconds"`
	DeltaSharingOrganizationName                *string                 `json:"delta_sharing_organization_name"`
	Owner                                       *string                 `json:"owner"`
	PrivilegeModelVersion                       *string                 `json:"privilege_model_version"`
	Region                                      *string                 `json:"region"`
	MetastoreId                                 *string                 `json:"metastore_id"`
	CreatedAt                                   *int64                  `json:"created_at"`
	CreatedBy                                   *string                 `json:"created_by"`
	UpdatedAt                                   *int64                  `json:"updated_at"`
	UpdatedBy                                   *string                 `json:"updated_by"`
	StorageRootCredentialName                   *string                 `json:"storage_root_credential_name"`
	Cloud                                       *string                 `json:"cloud"`
	GlobalMetastoreId                           *string                 `json:"global_metastore_id"`
	ExternalAccessEnabled                       *bool                   `json:"external_access_enabled"`
}

type CreateMetastoreAssignment struct {
	WorkspaceId        *int64  `json:"workspace_id"`
	MetastoreId        *string `json:"metastore_id"`
	DefaultCatalogName *string `json:"default_catalog_name"`
}

type CreateMetastoreAssignment_Response struct {
}

type DeleteMetastore struct {
	Id    *string `json:"id"`
	Force *bool   `json:"force"`
}

type DeleteMetastore_Response struct {
}

type DeleteMetastoreAssignment struct {
	WorkspaceId *int64  `json:"workspace_id"`
	MetastoreId *string `json:"metastore_id"`
}

type DeleteMetastoreAssignment_Response struct {
}

type DeltaSharingScope struct {
}

type GetCurrentMetastoreAssignment struct {
}

type GetMetastore struct {
	Id *string `json:"id"`
}

type GetMetastoreSummary struct {
}

type GetMetastoreSummary_Response struct {
	MetastoreId                                 *string                 `json:"metastore_id"`
	Name                                        *string                 `json:"name"`
	DefaultDataAccessConfigId                   *string                 `json:"default_data_access_config_id"`
	StorageRootCredentialId                     *string                 `json:"storage_root_credential_id"`
	Cloud                                       *string                 `json:"cloud"`
	Region                                      *string                 `json:"region"`
	GlobalMetastoreId                           *string                 `json:"global_metastore_id"`
	StorageRootCredentialName                   *string                 `json:"storage_root_credential_name"`
	PrivilegeModelVersion                       *string                 `json:"privilege_model_version"`
	DeltaSharingScope                           *DeltaSharingScope_Enum `json:"delta_sharing_scope"`
	DeltaSharingRecipientTokenLifetimeInSeconds *int64                  `json:"delta_sharing_recipient_token_lifetime_in_seconds"`
	DeltaSharingOrganizationName                *string                 `json:"delta_sharing_organization_name"`
	StorageRoot                                 *string                 `json:"storage_root"`
	Owner                                       *string                 `json:"owner"`
	CreatedAt                                   *int64                  `json:"created_at"`
	CreatedBy                                   *string                 `json:"created_by"`
	UpdatedAt                                   *int64                  `json:"updated_at"`
	UpdatedBy                                   *string                 `json:"updated_by"`
	ExternalAccessEnabled                       *bool                   `json:"external_access_enabled"`
}

type ListGlobalMetastoresRequest struct {
	PageSize       *int    `json:"page_size"`
	PageToken      *string `json:"page_token"`
	ShardLocalOnly *bool   `json:"shard_local_only"`
}

type ListGlobalMetastoresResponse struct {
	Metastores    []MetastoreInfo `json:"metastores"`
	NextPageToken *string         `json:"next_page_token"`
}

type ListMetastores struct {
	MaxResults *int    `json:"max_results"`
	PageToken  *string `json:"page_token"`
}

type ListMetastores_Response struct {
	Metastores    []MetastoreInfo `json:"metastores"`
	NextPageToken *string         `json:"next_page_token"`
}

type MetastoreAssignment struct {
	WorkspaceId        *int64  `json:"workspace_id"`
	MetastoreId        *string `json:"metastore_id"`
	DefaultCatalogName *string `json:"default_catalog_name"`
}

type MetastoreInfo struct {
	Name                                        *string                 `json:"name"`
	StorageRoot                                 *string                 `json:"storage_root"`
	DefaultDataAccessConfigId                   *string                 `json:"default_data_access_config_id"`
	StorageRootCredentialId                     *string                 `json:"storage_root_credential_id"`
	DeltaSharingScope                           *DeltaSharingScope_Enum `json:"delta_sharing_scope"`
	DeltaSharingRecipientTokenLifetimeInSeconds *int64                  `json:"delta_sharing_recipient_token_lifetime_in_seconds"`
	DeltaSharingOrganizationName                *string                 `json:"delta_sharing_organization_name"`
	Owner                                       *string                 `json:"owner"`
	PrivilegeModelVersion                       *string                 `json:"privilege_model_version"`
	Region                                      *string                 `json:"region"`
	MetastoreId                                 *string                 `json:"metastore_id"`
	CreatedAt                                   *int64                  `json:"created_at"`
	CreatedBy                                   *string                 `json:"created_by"`
	UpdatedAt                                   *int64                  `json:"updated_at"`
	UpdatedBy                                   *string                 `json:"updated_by"`
	StorageRootCredentialName                   *string                 `json:"storage_root_credential_name"`
	Cloud                                       *string                 `json:"cloud"`
	GlobalMetastoreId                           *string                 `json:"global_metastore_id"`
	ExternalAccessEnabled                       *bool                   `json:"external_access_enabled"`
}

type UpdateMetastore struct {
	Id                                          *string                 `json:"id"`
	NewName                                     *string                 `json:"new_name"`
	Name                                        *string                 `json:"name"`
	StorageRoot                                 *string                 `json:"storage_root"`
	DefaultDataAccessConfigId                   *string                 `json:"default_data_access_config_id"`
	StorageRootCredentialId                     *string                 `json:"storage_root_credential_id"`
	DeltaSharingScope                           *DeltaSharingScope_Enum `json:"delta_sharing_scope"`
	DeltaSharingRecipientTokenLifetimeInSeconds *int64                  `json:"delta_sharing_recipient_token_lifetime_in_seconds"`
	DeltaSharingOrganizationName                *string                 `json:"delta_sharing_organization_name"`
	Owner                                       *string                 `json:"owner"`
	PrivilegeModelVersion                       *string                 `json:"privilege_model_version"`
	Region                                      *string                 `json:"region"`
	MetastoreId                                 *string                 `json:"metastore_id"`
	CreatedAt                                   *int64                  `json:"created_at"`
	CreatedBy                                   *string                 `json:"created_by"`
	UpdatedAt                                   *int64                  `json:"updated_at"`
	UpdatedBy                                   *string                 `json:"updated_by"`
	StorageRootCredentialName                   *string                 `json:"storage_root_credential_name"`
	Cloud                                       *string                 `json:"cloud"`
	GlobalMetastoreId                           *string                 `json:"global_metastore_id"`
	ExternalAccessEnabled                       *bool                   `json:"external_access_enabled"`
}

type UpdateMetastoreAssignment struct {
	WorkspaceId        *int64  `json:"workspace_id"`
	MetastoreId        *string `json:"metastore_id"`
	DefaultCatalogName *string `json:"default_catalog_name"`
}

type UpdateMetastoreAssignment_Response struct {
}
