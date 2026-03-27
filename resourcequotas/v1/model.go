package v1

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

type GetQuota struct {
	ParentSecurableType *string `json:"parent_securable_type"`
	ParentFullName      *string `json:"parent_full_name"`
	QuotaName           *string `json:"quota_name"`
}

type GetQuota_Response struct {
	QuotaInfo *QuotaInfo `json:"quota_info"`
}

type ListQuotas struct {
	MaxResults *int    `json:"max_results"`
	PageToken  *string `json:"page_token"`
}

type ListQuotas_Response struct {
	Quotas        []QuotaInfo `json:"quotas"`
	NextPageToken *string     `json:"next_page_token"`
}

type QuotaInfo struct {
	ParentSecurableType *SecurableType `json:"parent_securable_type"`
	ParentFullName      *string        `json:"parent_full_name"`
	QuotaName           *string        `json:"quota_name"`
	QuotaCount          *int           `json:"quota_count"`
	QuotaLimit          *int           `json:"quota_limit"`
	LastRefreshedAt     *int64         `json:"last_refreshed_at"`
}
