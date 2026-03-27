package v1

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

type CreateSchema struct {
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	StorageLocation                     *string                              `json:"storage_location"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	SchemaId                            *string                              `json:"schema_id"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type CreateSchema_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateSchema_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type DeleteSchema struct {
	FullNameArg *string `json:"full_name_arg"`
	Force       *bool   `json:"force"`
}

type DeleteSchema_Response struct {
}

type EffectivePredictiveOptimizationFlag struct {
	Value             *string `json:"value"`
	InheritedFromType *string `json:"inherited_from_type"`
	InheritedFromName *string `json:"inherited_from_name"`
}

type GetSchema struct {
	FullNameArg   *string `json:"full_name_arg"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListSchemas struct {
	CatalogName   *string `json:"catalog_name"`
	MaxResults    *int    `json:"max_results"`
	PageToken     *string `json:"page_token"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListSchemas_Response struct {
	Schemas       []SchemaInfo `json:"schemas"`
	NextPageToken *string      `json:"next_page_token"`
}

// Next ID: 45.
type SchemaInfo struct {
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	StorageLocation                     *string                              `json:"storage_location"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	SchemaId                            *string                              `json:"schema_id"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type SchemaInfo_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type SchemaInfo_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateSchema struct {
	FullNameArg                         *string                              `json:"full_name_arg"`
	NewName                             *string                              `json:"new_name"`
	Name                                *string                              `json:"name"`
	CatalogName                         *string                              `json:"catalog_name"`
	Owner                               *string                              `json:"owner"`
	Comment                             *string                              `json:"comment"`
	StorageRoot                         *string                              `json:"storage_root"`
	EnablePredictiveOptimization        *string                              `json:"enable_predictive_optimization"`
	MetastoreId                         *string                              `json:"metastore_id"`
	FullName                            *string                              `json:"full_name"`
	CreatedAt                           *int64                               `json:"created_at"`
	CreatedBy                           *string                              `json:"created_by"`
	UpdatedAt                           *int64                               `json:"updated_at"`
	UpdatedBy                           *string                              `json:"updated_by"`
	CatalogType                         *CatalogType                         `json:"catalog_type"`
	StorageLocation                     *string                              `json:"storage_location"`
	EffectivePredictiveOptimizationFlag *EffectivePredictiveOptimizationFlag `json:"effective_predictive_optimization_flag"`
	SchemaId                            *string                              `json:"schema_id"`
	BrowseOnly                          *bool                                `json:"browse_only"`
	Properties                          map[string]string                    `json:"properties"`
	Options                             map[string]string                    `json:"options"`
}

type UpdateSchema_OptionsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type UpdateSchema_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}
