package v1

type AutoTaggingConfig_AutoTaggingMode string

const (
	AutoTaggingConfig_AutoTaggingMode_AutoTaggingModeUnspecified AutoTaggingConfig_AutoTaggingMode = "AUTO_TAGGING_MODE_UNSPECIFIED"
	AutoTaggingConfig_AutoTaggingMode_AutoTaggingDisabled        AutoTaggingConfig_AutoTaggingMode = "AUTO_TAGGING_DISABLED"
	AutoTaggingConfig_AutoTaggingMode_AutoTaggingEnabled         AutoTaggingConfig_AutoTaggingMode = "AUTO_TAGGING_ENABLED"
)

// String representation for [fmt.Print].
func (f *AutoTaggingConfig_AutoTaggingMode) String() string {
	return string(*f)
}

// Auto-tagging configuration for a classification tag. When enabled, detected
// columns are automatically tagged with Unity Catalog tags..
type AutoTaggingConfig struct {
	ClassificationTag *string                            `json:"classification_tag"`
	AutoTaggingMode   *AutoTaggingConfig_AutoTaggingMode `json:"auto_tagging_mode"`
}

// Data Classification configuration for a Unity Catalog catalog. This message
// follows the "At Most One Resource" pattern: at most one CatalogConfig exists
// per catalog. - Full CRUD operations are supported: Create enables Data
// Classification, Delete disables it - It has no unique identifier of its own
// and uses its parent catalog's identifier (catalog_name).
type CatalogConfig struct {
	Name                    *string                    `json:"name"`
	IncludedSchemas         *CatalogConfig_SchemaNames `json:"included_schemas"`
	AutoTagConfigs          []AutoTaggingConfig        `json:"auto_tag_configs"`
	EffectiveAutoTagConfigs []AutoTaggingConfig        `json:"effective_auto_tag_configs"`
}

// Wrapper message for a list of schema names..
type CatalogConfig_SchemaNames struct {
	Names []string `json:"names"`
}

// Create Data Classification configuration for a catalog. Creating a config
// enables Data Classification for the catalog..
type CreateCatalogConfigRequest struct {
	Parent        *string        `json:"parent"`
	CatalogConfig *CatalogConfig `json:"catalog_config"`
}

// Delete Data Classification configuration for a catalog. Deleting the config
// disables Data Classification for the catalog..
type DeleteCatalogConfigRequest struct {
	Name *string `json:"name"`
}

// Get Data Classification configuration for a catalog..
type GetCatalogConfigRequest struct {
	Name *string `json:"name"`
}

// Request to update the Data Classification configuration for a catalog.
//
// Uses field mask to support partial updates of the configuration. Only the
// fields specified in the update_mask will be modified..
type UpdateCatalogConfigRequest struct {
	CatalogConfig *CatalogConfig `json:"catalog_config"`
	UpdateMask    *string        `json:"update_mask"`
}
