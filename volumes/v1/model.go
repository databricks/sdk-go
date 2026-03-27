package v1

type SseEncryptionAlgorithm string

const (
	SseEncryptionAlgorithm_SseEncryptionAlgorithmUnspecified SseEncryptionAlgorithm = "SSE_ENCRYPTION_ALGORITHM_UNSPECIFIED"
	SseEncryptionAlgorithm_AwsSseS3                          SseEncryptionAlgorithm = "AWS_SSE_S3"
	SseEncryptionAlgorithm_AwsSseKms                         SseEncryptionAlgorithm = "AWS_SSE_KMS"
)

// String representation for [fmt.Print].
func (f *SseEncryptionAlgorithm) String() string {
	return string(*f)
}

type VolumeType string

const (
	VolumeType_Managed  VolumeType = "MANAGED"
	VolumeType_External VolumeType = "EXTERNAL"
)

// String representation for [fmt.Print].
func (f *VolumeType) String() string {
	return string(*f)
}

type CreateVolume struct {
	Name              *string            `json:"name"`
	CatalogName       *string            `json:"catalog_name"`
	SchemaName        *string            `json:"schema_name"`
	VolumeType        *VolumeType        `json:"volume_type"`
	StorageLocation   *string            `json:"storage_location"`
	Owner             *string            `json:"owner"`
	Comment           *string            `json:"comment"`
	FullName          *string            `json:"full_name"`
	VolumeId          *string            `json:"volume_id"`
	MetastoreId       *string            `json:"metastore_id"`
	CreatedAt         *int64             `json:"created_at"`
	CreatedBy         *string            `json:"created_by"`
	UpdatedAt         *int64             `json:"updated_at"`
	UpdatedBy         *string            `json:"updated_by"`
	AccessPoint       *string            `json:"access_point"`
	EncryptionDetails *EncryptionDetails `json:"encryption_details"`
	BrowseOnly        *bool              `json:"browse_only"`
}

type DeleteVolume struct {
	FullNameArg *string `json:"full_name_arg"`
}

type DeleteVolume_Response struct {
}

// Encryption options that apply to clients connecting to cloud storage..
type EncryptionDetails struct {
	SseEncryptionDetails *SseEncryptionDetails `json:"sse_encryption_details"`
}

type GetVolume struct {
	FullNameArg   *string `json:"full_name_arg"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListVolumes struct {
	CatalogName   *string `json:"catalog_name"`
	SchemaName    *string `json:"schema_name"`
	IncludeBrowse *bool   `json:"include_browse"`
	MaxResults    *int    `json:"max_results"`
	PageToken     *string `json:"page_token"`
}

type ListVolumes_Response struct {
	Volumes       []VolumeInfo `json:"volumes"`
	NextPageToken *string      `json:"next_page_token"`
}

// Server-Side Encryption properties for clients communicating with AWS s3..
type SseEncryptionDetails struct {
	Algorithm    *SseEncryptionAlgorithm `json:"algorithm"`
	AwsKmsKeyArn *string                 `json:"aws_kms_key_arn"`
}

type UpdateVolume struct {
	FullNameArg       *string            `json:"full_name_arg"`
	NewName           *string            `json:"new_name"`
	Name              *string            `json:"name"`
	CatalogName       *string            `json:"catalog_name"`
	SchemaName        *string            `json:"schema_name"`
	VolumeType        *VolumeType        `json:"volume_type"`
	StorageLocation   *string            `json:"storage_location"`
	Owner             *string            `json:"owner"`
	Comment           *string            `json:"comment"`
	FullName          *string            `json:"full_name"`
	VolumeId          *string            `json:"volume_id"`
	MetastoreId       *string            `json:"metastore_id"`
	CreatedAt         *int64             `json:"created_at"`
	CreatedBy         *string            `json:"created_by"`
	UpdatedAt         *int64             `json:"updated_at"`
	UpdatedBy         *string            `json:"updated_by"`
	AccessPoint       *string            `json:"access_point"`
	EncryptionDetails *EncryptionDetails `json:"encryption_details"`
	BrowseOnly        *bool              `json:"browse_only"`
}

type VolumeInfo struct {
	Name              *string            `json:"name"`
	CatalogName       *string            `json:"catalog_name"`
	SchemaName        *string            `json:"schema_name"`
	VolumeType        *VolumeType        `json:"volume_type"`
	StorageLocation   *string            `json:"storage_location"`
	Owner             *string            `json:"owner"`
	Comment           *string            `json:"comment"`
	FullName          *string            `json:"full_name"`
	VolumeId          *string            `json:"volume_id"`
	MetastoreId       *string            `json:"metastore_id"`
	CreatedAt         *int64             `json:"created_at"`
	CreatedBy         *string            `json:"created_by"`
	UpdatedAt         *int64             `json:"updated_at"`
	UpdatedBy         *string            `json:"updated_by"`
	AccessPoint       *string            `json:"access_point"`
	EncryptionDetails *EncryptionDetails `json:"encryption_details"`
	BrowseOnly        *bool              `json:"browse_only"`
}
