package v1

import (
	"io"
)

type AclPermission string

const (
	AclPermission_Read   AclPermission = "READ"
	AclPermission_Write  AclPermission = "WRITE"
	AclPermission_Manage AclPermission = "MANAGE"
)

// String representation for [fmt.Print].
func (f *AclPermission) String() string {
	return string(*f)
}

type ScopeBackendType string

const (
	ScopeBackendType_Databricks    ScopeBackendType = "DATABRICKS"
	ScopeBackendType_AzureKeyvault ScopeBackendType = "AZURE_KEYVAULT"
)

// String representation for [fmt.Print].
func (f *ScopeBackendType) String() string {
	return string(*f)
}

// An item representing an ACL rule applied to the given principal (user or
// group) on the associated scope point..
type AclItem struct {
	Principal  *string        `json:"principal"`
	Permission *AclPermission `json:"permission"`
}

// The metadata of the Azure KeyVault for a secret scope of type
// `AZURE_KEYVAULT`.
type AzureKeyVaultSecretScopeMetadata struct {
	ResourceId *string `json:"resource_id"`
	DnsName    *string `json:"dns_name"`
}

type CreateScope struct {
	Scope                  *string                           `json:"scope"`
	InitialManagePrincipal *string                           `json:"initial_manage_principal"`
	ScopeBackendType       *ScopeBackendType                 `json:"scope_backend_type"`
	BackendAzureKeyvault   *AzureKeyVaultSecretScopeMetadata `json:"backend_azure_keyvault"`
}

type CreateScope_Response struct {
}

type DeleteAcl struct {
	Scope     *string `json:"scope"`
	Principal *string `json:"principal"`
}

type DeleteAcl_Response struct {
}

type DeleteScope struct {
	Scope *string `json:"scope"`
}

type DeleteScope_Response struct {
}

type DeleteSecret struct {
	Scope *string `json:"scope"`
	Key   *string `json:"key"`
}

type DeleteSecret_Response struct {
}

type GetAcl struct {
	Scope     *string `json:"scope"`
	Principal *string `json:"principal"`
}

type GetSecret struct {
	Scope *string `json:"scope"`
	Key   *string `json:"key"`
}

type GetSecret_Response struct {
	Key   *string        `json:"key"`
	Value *io.ReadCloser `json:"value"`
}

type ListAcls struct {
	Scope *string `json:"scope"`
}

type ListAcls_Response struct {
	Items []AclItem `json:"items"`
}

type ListScopes struct {
}

type ListScopes_Response struct {
	Scopes []SecretScope `json:"scopes"`
}

type ListSecrets struct {
	Scope *string `json:"scope"`
}

type ListSecrets_Response struct {
	Secrets []SecretMetadata `json:"secrets"`
}

type PutAcl struct {
	Scope      *string        `json:"scope"`
	Principal  *string        `json:"principal"`
	Permission *AclPermission `json:"permission"`
}

type PutAcl_Response struct {
}

type PutSecret struct {
	Scope       *string        `json:"scope"`
	Key         *string        `json:"key"`
	StringValue *string        `json:"string_value"`
	BytesValue  *io.ReadCloser `json:"bytes_value"`
}

type PutSecret_Response struct {
}

// The metadata about a secret. Returned when listing secrets. Does not contain
// the actual secret value..
type SecretMetadata struct {
	Key                  *string `json:"key"`
	LastUpdatedTimestamp *int64  `json:"last_updated_timestamp"`
}

// An organizational resource for storing secrets. Secret scopes can be
// different types (Databricks-managed, Azure KeyVault backed, etc), and ACLs
// can be applied to control permissions for all secrets within a scope..
type SecretScope struct {
	Name             *string                           `json:"name"`
	BackendType      *ScopeBackendType                 `json:"backend_type"`
	KeyvaultMetadata *AzureKeyVaultSecretScopeMetadata `json:"keyvault_metadata"`
}
