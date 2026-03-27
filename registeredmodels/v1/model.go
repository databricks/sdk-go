package v1

type ModelVersionStatus string

const (
	ModelVersionStatus_ModelVersionStatusUnknown ModelVersionStatus = "MODEL_VERSION_STATUS_UNKNOWN"
	ModelVersionStatus_PendingRegistration       ModelVersionStatus = "PENDING_REGISTRATION"
	ModelVersionStatus_FailedRegistration        ModelVersionStatus = "FAILED_REGISTRATION"
	ModelVersionStatus_Ready                     ModelVersionStatus = "READY"
)

// String representation for [fmt.Print].
func (f *ModelVersionStatus) String() string {
	return string(*f)
}

// A connection that is dependent on a SQL object..
type ConnectionDependency struct {
	ConnectionName *string `json:"connection_name"`
}

type CreateRegisteredModel struct {
	Name            *string                    `json:"name"`
	CatalogName     *string                    `json:"catalog_name"`
	SchemaName      *string                    `json:"schema_name"`
	Owner           *string                    `json:"owner"`
	Comment         *string                    `json:"comment"`
	StorageLocation *string                    `json:"storage_location"`
	MetastoreId     *string                    `json:"metastore_id"`
	FullName        *string                    `json:"full_name"`
	CreatedAt       *int64                     `json:"created_at"`
	CreatedBy       *string                    `json:"created_by"`
	UpdatedAt       *int64                     `json:"updated_at"`
	UpdatedBy       *string                    `json:"updated_by"`
	Aliases         []RegisteredModelAliasInfo `json:"aliases"`
	BrowseOnly      *bool                      `json:"browse_only"`
}

// A credential that is dependent on a SQL object..
type CredentialDependency struct {
	CredentialName *string `json:"credential_name"`
}

type DeleteModelVersion struct {
	FullNameArg *string `json:"full_name_arg"`
	VersionArg  *int64  `json:"version_arg"`
}

type DeleteModelVersion_Response struct {
}

type DeleteRegisteredModel struct {
	FullNameArg *string `json:"full_name_arg"`
}

type DeleteRegisteredModel_Response struct {
}

type DeleteRegisteredModelAlias struct {
	FullNameArg *string `json:"full_name_arg"`
	AliasArg    *string `json:"alias_arg"`
}

type DeleteRegisteredModelAlias_Response struct {
}

// A dependency of a SQL object. One of the following fields must be defined:
// __table__, __function__, __connection__, or __credential__..
type Dependency struct {
	Table      *TableDependency      `json:"table"`
	Function   *FunctionDependency   `json:"function"`
	Connection *ConnectionDependency `json:"connection"`
	Credential *CredentialDependency `json:"credential"`
}

// A list of dependencies..
type DependencyList struct {
	Dependencies []Dependency `json:"dependencies"`
}

// A function that is dependent on a SQL object..
type FunctionDependency struct {
	FunctionFullName *string `json:"function_full_name"`
}

type GetModelVersion struct {
	FullNameArg    *string `json:"full_name_arg"`
	VersionArg     *int64  `json:"version_arg"`
	IncludeAliases *bool   `json:"include_aliases"`
	IncludeBrowse  *bool   `json:"include_browse"`
}

type GetModelVersionByAlias struct {
	FullNameArg    *string `json:"full_name_arg"`
	AliasArg       *string `json:"alias_arg"`
	IncludeAliases *bool   `json:"include_aliases"`
}

type GetRegisteredModel struct {
	FullNameArg    *string `json:"full_name_arg"`
	IncludeAliases *bool   `json:"include_aliases"`
	IncludeBrowse  *bool   `json:"include_browse"`
}

type ListModelVersions struct {
	FullNameArg   *string `json:"full_name_arg"`
	MaxResults    *int64  `json:"max_results"`
	PageToken     *string `json:"page_token"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListModelVersions_Response struct {
	ModelVersions []ModelVersionInfo `json:"model_versions"`
	NextPageToken *string            `json:"next_page_token"`
}

type ListRegisteredModels struct {
	CatalogName   *string `json:"catalog_name"`
	SchemaName    *string `json:"schema_name"`
	IncludeBrowse *bool   `json:"include_browse"`
	MaxResults    *int64  `json:"max_results"`
	PageToken     *string `json:"page_token"`
}

type ListRegisteredModels_Response struct {
	RegisteredModels []RegisteredModelInfo `json:"registered_models"`
	NextPageToken    *string               `json:"next_page_token"`
}

type ModelVersionInfo struct {
	ModelName                *string                    `json:"model_name"`
	CatalogName              *string                    `json:"catalog_name"`
	SchemaName               *string                    `json:"schema_name"`
	Source                   *string                    `json:"source"`
	Comment                  *string                    `json:"comment"`
	RunId                    *string                    `json:"run_id"`
	RunWorkspaceId           *int64                     `json:"run_workspace_id"`
	ModelVersionDependencies *DependencyList            `json:"model_version_dependencies"`
	Status                   *ModelVersionStatus        `json:"status"`
	Version                  *int64                     `json:"version"`
	StorageLocation          *string                    `json:"storage_location"`
	MetastoreId              *string                    `json:"metastore_id"`
	CreatedAt                *int64                     `json:"created_at"`
	CreatedBy                *string                    `json:"created_by"`
	UpdatedAt                *int64                     `json:"updated_at"`
	UpdatedBy                *string                    `json:"updated_by"`
	Id                       *string                    `json:"id"`
	Aliases                  []RegisteredModelAliasInfo `json:"aliases"`
}

type RegisteredModelAliasInfo struct {
	AliasName   *string `json:"alias_name"`
	VersionNum  *int64  `json:"version_num"`
	Id          *string `json:"id"`
	ModelName   *string `json:"model_name"`
	CatalogName *string `json:"catalog_name"`
	SchemaName  *string `json:"schema_name"`
}

type RegisteredModelInfo struct {
	Name            *string                    `json:"name"`
	CatalogName     *string                    `json:"catalog_name"`
	SchemaName      *string                    `json:"schema_name"`
	Owner           *string                    `json:"owner"`
	Comment         *string                    `json:"comment"`
	StorageLocation *string                    `json:"storage_location"`
	MetastoreId     *string                    `json:"metastore_id"`
	FullName        *string                    `json:"full_name"`
	CreatedAt       *int64                     `json:"created_at"`
	CreatedBy       *string                    `json:"created_by"`
	UpdatedAt       *int64                     `json:"updated_at"`
	UpdatedBy       *string                    `json:"updated_by"`
	Aliases         []RegisteredModelAliasInfo `json:"aliases"`
	BrowseOnly      *bool                      `json:"browse_only"`
}

type SetRegisteredModelAlias struct {
	FullNameArg *string `json:"full_name_arg"`
	AliasArg    *string `json:"alias_arg"`
	VersionNum  *int64  `json:"version_num"`
}

// A table that is dependent on a SQL object..
type TableDependency struct {
	TableFullName *string `json:"table_full_name"`
}

type UpdateModelVersion struct {
	FullNameArg              *string                    `json:"full_name_arg"`
	VersionArg               *int64                     `json:"version_arg"`
	ModelName                *string                    `json:"model_name"`
	CatalogName              *string                    `json:"catalog_name"`
	SchemaName               *string                    `json:"schema_name"`
	Source                   *string                    `json:"source"`
	Comment                  *string                    `json:"comment"`
	RunId                    *string                    `json:"run_id"`
	RunWorkspaceId           *int64                     `json:"run_workspace_id"`
	ModelVersionDependencies *DependencyList            `json:"model_version_dependencies"`
	Status                   *ModelVersionStatus        `json:"status"`
	Version                  *int64                     `json:"version"`
	StorageLocation          *string                    `json:"storage_location"`
	MetastoreId              *string                    `json:"metastore_id"`
	CreatedAt                *int64                     `json:"created_at"`
	CreatedBy                *string                    `json:"created_by"`
	UpdatedAt                *int64                     `json:"updated_at"`
	UpdatedBy                *string                    `json:"updated_by"`
	Id                       *string                    `json:"id"`
	Aliases                  []RegisteredModelAliasInfo `json:"aliases"`
}

type UpdateRegisteredModel struct {
	FullNameArg     *string                    `json:"full_name_arg"`
	NewName         *string                    `json:"new_name"`
	Name            *string                    `json:"name"`
	CatalogName     *string                    `json:"catalog_name"`
	SchemaName      *string                    `json:"schema_name"`
	Owner           *string                    `json:"owner"`
	Comment         *string                    `json:"comment"`
	StorageLocation *string                    `json:"storage_location"`
	MetastoreId     *string                    `json:"metastore_id"`
	FullName        *string                    `json:"full_name"`
	CreatedAt       *int64                     `json:"created_at"`
	CreatedBy       *string                    `json:"created_by"`
	UpdatedAt       *int64                     `json:"updated_at"`
	UpdatedBy       *string                    `json:"updated_by"`
	Aliases         []RegisteredModelAliasInfo `json:"aliases"`
	BrowseOnly      *bool                      `json:"browse_only"`
}
