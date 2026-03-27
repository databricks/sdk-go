package v1

import (
	"io"
)

type CatalogIsolationMode string

const (
	CatalogIsolationMode_Open          CatalogIsolationMode = "OPEN"
	CatalogIsolationMode_Isolated      CatalogIsolationMode = "ISOLATED"
	CatalogIsolationMode_OpenInAccount CatalogIsolationMode = "OPEN_IN_ACCOUNT"
)

// String representation for [fmt.Print].
func (f *CatalogIsolationMode) String() string {
	return string(*f)
}

type CatalogType string

const (
	CatalogType_ManagedCatalog       CatalogType = "MANAGED_CATALOG"
	CatalogType_DeltasharingCatalog  CatalogType = "DELTASHARING_CATALOG"
	CatalogType_SystemCatalog        CatalogType = "SYSTEM_CATALOG"
	CatalogType_InternalCatalog      CatalogType = "INTERNAL_CATALOG"
	CatalogType_ForeignCatalog       CatalogType = "FOREIGN_CATALOG"
	CatalogType_ManagedOnlineCatalog CatalogType = "MANAGED_ONLINE_CATALOG"
)

// String representation for [fmt.Print].
func (f *CatalogType) String() string {
	return string(*f)
}

type DrReplicationStatus string

const (
	DrReplicationStatus_DrReplicationStatusUnspecified DrReplicationStatus = "DR_REPLICATION_STATUS_UNSPECIFIED"
	DrReplicationStatus_DrReplicationStatusPrimary     DrReplicationStatus = "DR_REPLICATION_STATUS_PRIMARY"
	DrReplicationStatus_DrReplicationStatusSecondary   DrReplicationStatus = "DR_REPLICATION_STATUS_SECONDARY"
)

// String representation for [fmt.Print].
func (f *DrReplicationStatus) String() string {
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

type ConversionInfo_State string

const (
	ConversionInfo_State_StateUnspecified ConversionInfo_State = "STATE_UNSPECIFIED"
	ConversionInfo_State_InProgress       ConversionInfo_State = "IN_PROGRESS"
	ConversionInfo_State_Completed        ConversionInfo_State = "COMPLETED"
)

// String representation for [fmt.Print].
func (f *ConversionInfo_State) String() string {
	return string(*f)
}

type ProvisioningInfo_State string

const (
	ProvisioningInfo_State_StateUnspecified ProvisioningInfo_State = "STATE_UNSPECIFIED"
	ProvisioningInfo_State_Provisioning     ProvisioningInfo_State = "PROVISIONING"
	ProvisioningInfo_State_Active           ProvisioningInfo_State = "ACTIVE"
	ProvisioningInfo_State_Failed           ProvisioningInfo_State = "FAILED"
	ProvisioningInfo_State_Deleting         ProvisioningInfo_State = "DELETING"
	ProvisioningInfo_State_Updating         ProvisioningInfo_State = "UPDATING"
	ProvisioningInfo_State_Degraded         ProvisioningInfo_State = "DEGRADED"
)

// String representation for [fmt.Print].
func (f *ProvisioningInfo_State) String() string {
	return string(*f)
}

type AzureEncryptionSettings struct {
	AzureTenantId             *string `json:"azure_tenant_id"`
	AzureCmkAccessConnectorId *string `json:"azure_cmk_access_connector_id"`
	AzureCmkManagedIdentityId *string `json:"azure_cmk_managed_identity_id"`
}

type CatalogInfo struct {
	Name                                *string                              `json:"name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	ProviderName                        *string                              `json:"provider_name"`
	ShareName                           *string                              `json:"share_name"`
	ConnectionName                      *string                              `json:"connection_name"`
	MetastoreId                         *string                              `json:"metastore_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	StorageLocation                     *string                              `json:"storage_location"`
	IsolationMode                       *CatalogIsolationMode                `json:"isolation_mode"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	ProvisioningInfo                    *ProvisioningInfo                    `json:"provisioning_info"`
	FullName                            *string                              `json:"full_name"`
	SecurableType                       *SecurableType                       `json:"securable_type"`
	ConversionInfo                      *ConversionInfo                      `json:"conversion_info"`
	DrReplicationInfo                   *DrReplicationInfo                   `json:"dr_replication_info"`
	ManagedEncryptionSettings           *EncryptionSettings                  `json:"managed_encryption_settings"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type CatalogInfo_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CatalogInfo_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Status of conversion of FOREIGN entity into UC Native entity..
type ConversionInfo struct {
	State *ConversionInfo_State `json:"state"`
}

type CreateCatalog struct {
	Name                                *string                              `json:"name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	ProviderName                        *string                              `json:"provider_name"`
	ShareName                           *string                              `json:"share_name"`
	ConnectionName                      *string                              `json:"connection_name"`
	MetastoreId                         *string                              `json:"metastore_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	StorageLocation                     *string                              `json:"storage_location"`
	IsolationMode                       *CatalogIsolationMode                `json:"isolation_mode"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	ProvisioningInfo                    *ProvisioningInfo                    `json:"provisioning_info"`
	FullName                            *string                              `json:"full_name"`
	SecurableType                       *SecurableType                       `json:"securable_type"`
	ConversionInfo                      *ConversionInfo                      `json:"conversion_info"`
	DrReplicationInfo                   *DrReplicationInfo                   `json:"dr_replication_info"`
	ManagedEncryptionSettings           *EncryptionSettings                  `json:"managed_encryption_settings"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type CreateCatalog_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateCatalog_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type DeleteCatalog struct {
	NameArg *string `json:"name_arg"`
	Force   *bool   `json:"force"`
}

type DeleteCatalog_Response struct {
}

// Metadata related to Disaster Recovery..
type DrReplicationInfo struct {
	Status             *DrReplicationStatus `json:"status"`
	ReplicatedEntities *io.ReadCloser       `json:"replicated_entities"`
}

type EffectivePredictiveOptimizationFlag struct {
	Value             *string `json:"value"`
	InheritedFromType *string `json:"inherited_from_type"`
	InheritedFromName *string `json:"inherited_from_name"`
}

// Encryption Settings are used to carry metadata for securable encryption at
// rest. Currently used for catalogs, we can use the information supplied here
// to interact with a CMK..
type EncryptionSettings struct {
	CustomerManagedKeyId    *string                  `json:"customer_managed_key_id"`
	AzureKeyVaultKeyId      *string                  `json:"azure_key_vault_key_id"`
	AzureEncryptionSettings *AzureEncryptionSettings `json:"azure_encryption_settings"`
}

type GetCatalog struct {
	NameArg       *string `json:"name_arg"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListCatalogs struct {
	IncludeBrowse  *bool   `json:"include_browse"`
	MaxResults     *int    `json:"max_results"`
	PageToken      *string `json:"page_token"`
	IncludeUnbound *bool   `json:"include_unbound"`
}

type ListCatalogs_Response struct {
	Catalogs      []CatalogInfo `json:"catalogs"`
	NextPageToken *string       `json:"next_page_token"`
}

// Status of an asynchronously provisioned resource..
type ProvisioningInfo struct {
	State *ProvisioningInfo_State `json:"state"`
}

type UpdateCatalog struct {
	NameArg                             *string                              `json:"name_arg"`
	NewName                             *string                              `json:"new_name"`
	Name                                *string                              `json:"name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	ProviderName                        *string                              `json:"provider_name"`
	ShareName                           *string                              `json:"share_name"`
	ConnectionName                      *string                              `json:"connection_name"`
	MetastoreId                         *string                              `json:"metastore_id"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	StorageLocation                     *string                              `json:"storage_location"`
	IsolationMode                       *CatalogIsolationMode                `json:"isolation_mode"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	ProvisioningInfo                    *ProvisioningInfo                    `json:"provisioning_info"`
	FullName                            *string                              `json:"full_name"`
	SecurableType                       *SecurableType                       `json:"securable_type"`
	ConversionInfo                      *ConversionInfo                      `json:"conversion_info"`
	DrReplicationInfo                   *DrReplicationInfo                   `json:"dr_replication_info"`
	ManagedEncryptionSettings           *EncryptionSettings                  `json:"managed_encryption_settings"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type UpdateCatalog_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateCatalog_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
