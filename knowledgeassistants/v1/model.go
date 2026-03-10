package v1

import "time"

type KnowledgeAssistant_State string

const (
	KnowledgeAssistant_State_StateUnspecified KnowledgeAssistant_State = "STATE_UNSPECIFIED"
	KnowledgeAssistant_State_Creating         KnowledgeAssistant_State = "CREATING"
	KnowledgeAssistant_State_Active           KnowledgeAssistant_State = "ACTIVE"
	KnowledgeAssistant_State_Failed           KnowledgeAssistant_State = "FAILED"
)

// String representation for [fmt.Print].
func (f *KnowledgeAssistant_State) String() string {
	return string(*f)
}

type KnowledgeSource_State string

const (
	KnowledgeSource_State_StateUnspecified KnowledgeSource_State = "STATE_UNSPECIFIED"
	KnowledgeSource_State_Updating         KnowledgeSource_State = "UPDATING"
	KnowledgeSource_State_Updated          KnowledgeSource_State = "UPDATED"
	KnowledgeSource_State_FailedUpdate     KnowledgeSource_State = "FAILED_UPDATE"
)

// String representation for [fmt.Print].
func (f *KnowledgeSource_State) String() string {
	return string(*f)
}

type CreateKnowledgeAssistantRequest struct {
	KnowledgeAssistant *KnowledgeAssistant `json:"knowledge_assistant"`
}

type CreateKnowledgeSourceRequest struct {
	Parent          *string          `json:"parent"`
	KnowledgeSource *KnowledgeSource `json:"knowledge_source"`
}

// A request to delete a Knowledge Assistant..
type DeleteKnowledgeAssistantRequest struct {
	Name *string `json:"name"`
}

type DeleteKnowledgeSourceRequest struct {
	Name *string `json:"name"`
}

// FileTableSpec specifies a file table source configuration..
type FileTableSpec struct {
	TableName *string `json:"table_name"`
	FileCol   *string `json:"file_col"`
}

// FilesSpec specifies a files source configuration..
type FilesSpec struct {
	Path *string `json:"path"`
}

// A request to retrieve a Knowledge Assistant..
type GetKnowledgeAssistantRequest struct {
	Name *string `json:"name"`
}

type GetKnowledgeSourceRequest struct {
	Name *string `json:"name"`
}

// IndexSpec specifies a vector search index source configuration..
type IndexSpec struct {
	IndexName *string `json:"index_name"`
	TextCol   *string `json:"text_col"`
	DocUriCol *string `json:"doc_uri_col"`
}

// Entity message that represents a knowledge assistant. Note: REQUIRED
// annotations below represent create-time requirements. For updates, required
// fields are determined by the update mask..
type KnowledgeAssistant struct {
	Name         *string                   `json:"name"`
	State        *KnowledgeAssistant_State `json:"state"`
	Id           *string                   `json:"id"`
	DisplayName  *string                   `json:"display_name"`
	Description  *string                   `json:"description"`
	Instructions *string                   `json:"instructions"`
	Creator      *string                   `json:"creator"`
	CreateTime   *time.Time                `json:"create_time"`
	EndpointName *string                   `json:"endpoint_name"`
	ExperimentId *string                   `json:"experiment_id"`
	ErrorInfo    *string                   `json:"error_info"`
}

// KnowledgeSource represents a source of knowledge for the KnowledgeAssistant.
// Used in create/update requests and returned in Get/List responses. Note:
// REQUIRED annotations below represent create-time requirements. For updates,
// required fields are determined by the update mask..
type KnowledgeSource struct {
	Name                *string                `json:"name"`
	DisplayName         *string                `json:"display_name"`
	Description         *string                `json:"description"`
	SourceType          *string                `json:"source_type"`
	Index               *IndexSpec             `json:"index"`
	Files               *FilesSpec             `json:"files"`
	FileTable           *FileTableSpec         `json:"file_table"`
	State               *KnowledgeSource_State `json:"state"`
	Id                  *string                `json:"id"`
	KnowledgeCutoffTime *time.Time             `json:"knowledge_cutoff_time"`
	CreateTime          *time.Time             `json:"create_time"`
}

// A request to list Knowledge Assistants..
type ListKnowledgeAssistantsRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

// A list of Knowledge Assistants..
type ListKnowledgeAssistantsResponse struct {
	KnowledgeAssistants []KnowledgeAssistant `json:"knowledge_assistants"`
	NextPageToken       *string              `json:"next_page_token"`
}

type ListKnowledgeSourcesRequest struct {
	Parent    *string `json:"parent"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type ListKnowledgeSourcesResponse struct {
	KnowledgeSources []KnowledgeSource `json:"knowledge_sources"`
	NextPageToken    *string           `json:"next_page_token"`
}

type SyncKnowledgeSourcesRequest struct {
	Name *string `json:"name"`
}

type UpdateKnowledgeAssistantRequest struct {
	KnowledgeAssistant *KnowledgeAssistant `json:"knowledge_assistant"`
	UpdateMask         *string             `json:"update_mask"`
}

type UpdateKnowledgeSourceRequest struct {
	Name            *string          `json:"name"`
	KnowledgeSource *KnowledgeSource `json:"knowledge_source"`
	UpdateMask      *string          `json:"update_mask"`
}
