package v1

type BindingType string

const (
	BindingType_BindingTypeUnspecified BindingType = "BINDING_TYPE_UNSPECIFIED"
	BindingType_BindingTypeReadWrite   BindingType = "BINDING_TYPE_READ_WRITE"
	BindingType_BindingTypeReadOnly    BindingType = "BINDING_TYPE_READ_ONLY"
)

// String representation for [fmt.Print].
func (f *BindingType) String() string {
	return string(*f)
}

type GetCatalogWorkspaceBindings struct {
	CatalogName *string `json:"catalog_name"`
}

type GetCatalogWorkspaceBindings_Response struct {
	Workspaces []int64 `json:"workspaces"`
}

type GetWorkspaceBindings struct {
	SecurableType     *string `json:"securable_type"`
	SecurableFullName *string `json:"securable_full_name"`
	MaxResults        *int    `json:"max_results"`
	PageToken         *string `json:"page_token"`
}

type GetWorkspaceBindings_Response struct {
	Bindings      []WorkspaceBindingInfo `json:"bindings"`
	NextPageToken *string                `json:"next_page_token"`
}

type UpdateCatalogWorkspaceBindings struct {
	CatalogName        *string `json:"catalog_name"`
	AssignWorkspaces   []int64 `json:"assign_workspaces"`
	UnassignWorkspaces []int64 `json:"unassign_workspaces"`
}

type UpdateCatalogWorkspaceBindings_Response struct {
	Workspaces []int64 `json:"workspaces"`
}

type UpdateWorkspaceBindings struct {
	SecurableType     *string                `json:"securable_type"`
	SecurableFullName *string                `json:"securable_full_name"`
	Add               []WorkspaceBindingInfo `json:"add"`
	Remove            []WorkspaceBindingInfo `json:"remove"`
}

// A list of workspace IDs that are bound to the securable.
type UpdateWorkspaceBindings_Response struct {
	Bindings []WorkspaceBindingInfo `json:"bindings"`
}

type WorkspaceBindingInfo struct {
	WorkspaceId *int64       `json:"workspace_id"`
	BindingType *BindingType `json:"binding_type"`
}
