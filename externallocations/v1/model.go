package v1

type IsolationMode string

const (
	IsolationMode_IsolationModeUnspecified   IsolationMode = "ISOLATION_MODE_UNSPECIFIED"
	IsolationMode_IsolationModeOpen          IsolationMode = "ISOLATION_MODE_OPEN"
	IsolationMode_IsolationModeIsolated      IsolationMode = "ISOLATION_MODE_ISOLATED"
	IsolationMode_IsolationModeOpenInAccount IsolationMode = "ISOLATION_MODE_OPEN_IN_ACCOUNT"
)

// String representation for [fmt.Print].
func (f *IsolationMode) String() string {
	return string(*f)
}

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

type AwsSqsQueue struct {
	QueueUrl          *string `json:"queue_url"`
	ManagedResourceId *string `json:"managed_resource_id"`
}

type AzureQueueStorage struct {
	QueueUrl          *string `json:"queue_url"`
	SubscriptionId    *string `json:"subscription_id"`
	ResourceGroup     *string `json:"resource_group"`
	ManagedResourceId *string `json:"managed_resource_id"`
}

type CreateExternalLocation struct {
	SkipValidation            *bool              `json:"skip_validation"`
	Name                      *string            `json:"name"`
	Url                       *string            `json:"url"`
	CredentialName            *string            `json:"credential_name"`
	ReadOnly                  *bool              `json:"read_only"`
	Comment                   *string            `json:"comment"`
	EnableFileEvents          *bool              `json:"enable_file_events"`
	FileEventQueue            *FileEventQueue    `json:"file_event_queue"`
	Owner                     *string            `json:"owner"`
	EncryptionDetails         *EncryptionDetails `json:"encryption_details"`
	MetastoreId               *string            `json:"metastore_id"`
	CredentialId              *string            `json:"credential_id"`
	CreatedAt                 *int64             `json:"created_at"`
	CreatedBy                 *string            `json:"created_by"`
	UpdatedAt                 *int64             `json:"updated_at"`
	UpdatedBy                 *string            `json:"updated_by"`
	BrowseOnly                *bool              `json:"browse_only"`
	IsolationMode             *IsolationMode     `json:"isolation_mode"`
	Fallback                  *bool              `json:"fallback"`
	EffectiveEnableFileEvents *bool              `json:"effective_enable_file_events"`
	EffectiveFileEventQueue   *FileEventQueue    `json:"effective_file_event_queue"`
}

type DeleteExternalLocation struct {
	NameArg *string `json:"name_arg"`
	Force   *bool   `json:"force"`
}

type DeleteExternalLocation_Response struct {
}

// Encryption options that apply to clients connecting to cloud storage..
type EncryptionDetails struct {
	SseEncryptionDetails *SseEncryptionDetails `json:"sse_encryption_details"`
}

type ExternalLocationInfo struct {
	Name                      *string            `json:"name"`
	Url                       *string            `json:"url"`
	CredentialName            *string            `json:"credential_name"`
	ReadOnly                  *bool              `json:"read_only"`
	Comment                   *string            `json:"comment"`
	EnableFileEvents          *bool              `json:"enable_file_events"`
	FileEventQueue            *FileEventQueue    `json:"file_event_queue"`
	Owner                     *string            `json:"owner"`
	EncryptionDetails         *EncryptionDetails `json:"encryption_details"`
	MetastoreId               *string            `json:"metastore_id"`
	CredentialId              *string            `json:"credential_id"`
	CreatedAt                 *int64             `json:"created_at"`
	CreatedBy                 *string            `json:"created_by"`
	UpdatedAt                 *int64             `json:"updated_at"`
	UpdatedBy                 *string            `json:"updated_by"`
	BrowseOnly                *bool              `json:"browse_only"`
	IsolationMode             *IsolationMode     `json:"isolation_mode"`
	Fallback                  *bool              `json:"fallback"`
	EffectiveEnableFileEvents *bool              `json:"effective_enable_file_events"`
	EffectiveFileEventQueue   *FileEventQueue    `json:"effective_file_event_queue"`
}

type FileEventQueue struct {
	ProvidedAqs    *AzureQueueStorage `json:"provided_aqs"`
	ProvidedSqs    *AwsSqsQueue       `json:"provided_sqs"`
	ProvidedPubsub *GcpPubsub         `json:"provided_pubsub"`
	ManagedAqs     *AzureQueueStorage `json:"managed_aqs"`
	ManagedSqs     *AwsSqsQueue       `json:"managed_sqs"`
	ManagedPubsub  *GcpPubsub         `json:"managed_pubsub"`
}

type GcpPubsub struct {
	SubscriptionName  *string `json:"subscription_name"`
	ManagedResourceId *string `json:"managed_resource_id"`
}

type GetExternalLocation struct {
	NameArg       *string `json:"name_arg"`
	IncludeBrowse *bool   `json:"include_browse"`
}

type ListExternalLocations struct {
	IncludeBrowse  *bool   `json:"include_browse"`
	MaxResults     *int    `json:"max_results"`
	PageToken      *string `json:"page_token"`
	IncludeUnbound *bool   `json:"include_unbound"`
}

type ListExternalLocations_Response struct {
	ExternalLocations []ExternalLocationInfo `json:"external_locations"`
	NextPageToken     *string                `json:"next_page_token"`
}

// Server-Side Encryption properties for clients communicating with AWS s3..
type SseEncryptionDetails struct {
	Algorithm    *SseEncryptionAlgorithm `json:"algorithm"`
	AwsKmsKeyArn *string                 `json:"aws_kms_key_arn"`
}

type UpdateExternalLocation struct {
	NameArg                   *string            `json:"name_arg"`
	NewName                   *string            `json:"new_name"`
	Force                     *bool              `json:"force"`
	SkipValidation            *bool              `json:"skip_validation"`
	Name                      *string            `json:"name"`
	Url                       *string            `json:"url"`
	CredentialName            *string            `json:"credential_name"`
	ReadOnly                  *bool              `json:"read_only"`
	Comment                   *string            `json:"comment"`
	EnableFileEvents          *bool              `json:"enable_file_events"`
	FileEventQueue            *FileEventQueue    `json:"file_event_queue"`
	Owner                     *string            `json:"owner"`
	EncryptionDetails         *EncryptionDetails `json:"encryption_details"`
	MetastoreId               *string            `json:"metastore_id"`
	CredentialId              *string            `json:"credential_id"`
	CreatedAt                 *int64             `json:"created_at"`
	CreatedBy                 *string            `json:"created_by"`
	UpdatedAt                 *int64             `json:"updated_at"`
	UpdatedBy                 *string            `json:"updated_by"`
	BrowseOnly                *bool              `json:"browse_only"`
	IsolationMode             *IsolationMode     `json:"isolation_mode"`
	Fallback                  *bool              `json:"fallback"`
	EffectiveEnableFileEvents *bool              `json:"effective_enable_file_events"`
	EffectiveFileEventQueue   *FileEventQueue    `json:"effective_file_event_queue"`
}
