package v1

type PolicyType string

const (
	PolicyType_PolicyTypeUnspecified PolicyType = "POLICY_TYPE_UNSPECIFIED"
	PolicyType_PolicyTypeRowFilter   PolicyType = "POLICY_TYPE_ROW_FILTER"
	PolicyType_PolicyTypeColumnMask  PolicyType = "POLICY_TYPE_COLUMN_MASK"
	PolicyType_PolicyTypeDeny        PolicyType = "POLICY_TYPE_DENY"
)

// String representation for [fmt.Print].
func (f *PolicyType) String() string {
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

type ColumnMaskOptions struct {
	FunctionName *string            `json:"function_name"`
	OnColumn     *string            `json:"on_column"`
	Using        []FunctionArgument `json:"using"`
}

type CreatePolicy struct {
	PolicyInfo *PolicyInfo `json:"policy_info"`
}

type DeletePolicy struct {
	OnSecurableType     *string `json:"on_securable_type"`
	OnSecurableFullname *string `json:"on_securable_fullname"`
	Name                *string `json:"name"`
}

type DeletePolicy_Response struct {
}

type DenyOptions struct {
	Privileges []string `json:"privileges"`
}

type FunctionArgument struct {
	Alias    *string `json:"alias"`
	Constant *string `json:"constant"`
}

type GetPolicy struct {
	OnSecurableType     *string `json:"on_securable_type"`
	OnSecurableFullname *string `json:"on_securable_fullname"`
	Name                *string `json:"name"`
}

type ListPolicies struct {
	OnSecurableType     *string `json:"on_securable_type"`
	OnSecurableFullname *string `json:"on_securable_fullname"`
	IncludeInherited    *bool   `json:"include_inherited"`
	MaxResults          *int    `json:"max_results"`
	PageToken           *string `json:"page_token"`
}

type ListPolicies_Response struct {
	Policies      []PolicyInfo `json:"policies"`
	NextPageToken *string      `json:"next_page_token"`
}

type MatchColumn struct {
	Condition *string `json:"condition"`
	Alias     *string `json:"alias"`
}

type PolicyInfo struct {
	Id                  *string            `json:"id"`
	OnSecurableType     *SecurableType     `json:"on_securable_type"`
	OnSecurableFullname *string            `json:"on_securable_fullname"`
	Name                *string            `json:"name"`
	Comment             *string            `json:"comment"`
	ToPrincipals        []string           `json:"to_principals"`
	ExceptPrincipals    []string           `json:"except_principals"`
	ForSecurableType    *SecurableType     `json:"for_securable_type"`
	WhenCondition       *string            `json:"when_condition"`
	PolicyType          *PolicyType        `json:"policy_type"`
	RowFilter           *RowFilterOptions  `json:"row_filter"`
	ColumnMask          *ColumnMaskOptions `json:"column_mask"`
	Deny                *DenyOptions       `json:"deny"`
	MatchColumns        []MatchColumn      `json:"match_columns"`
	CreatedAt           *int64             `json:"created_at"`
	CreatedBy           *string            `json:"created_by"`
	UpdatedAt           *int64             `json:"updated_at"`
	UpdatedBy           *string            `json:"updated_by"`
	UseSessionIdentity  *bool              `json:"use_session_identity"`
}

type RowFilterOptions struct {
	FunctionName *string            `json:"function_name"`
	Using        []FunctionArgument `json:"using"`
}

type UpdatePolicy struct {
	OnSecurableType     *string     `json:"on_securable_type"`
	OnSecurableFullname *string     `json:"on_securable_fullname"`
	Name                *string     `json:"name"`
	PolicyInfo          *PolicyInfo `json:"policy_info"`
	UpdateMask          *string     `json:"update_mask"`
}
