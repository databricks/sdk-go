package v1

import "encoding/json"

type ChatMessageRole string

const (
	ChatMessageRole_ChatMessageRoleUnspecified ChatMessageRole = "CHAT_MESSAGE_ROLE_UNSPECIFIED"
	ChatMessageRole_System                     ChatMessageRole = "SYSTEM"
	ChatMessageRole_User                       ChatMessageRole = "USER"
	ChatMessageRole_Assistant                  ChatMessageRole = "ASSISTANT"
)

// String representation for [fmt.Print].
func (f *ChatMessageRole) String() string {
	return string(*f)
}

type EmbeddingsV1ResponseEmbeddingElementObject string

const (
	EmbeddingsV1ResponseEmbeddingElementObject_EmbeddingsV1ResponseEmbeddingElementObjectUnspecified EmbeddingsV1ResponseEmbeddingElementObject = "EMBEDDINGS_V1_RESPONSE_EMBEDDING_ELEMENT_OBJECT_UNSPECIFIED"
	EmbeddingsV1ResponseEmbeddingElementObject_Embedding                                             EmbeddingsV1ResponseEmbeddingElementObject = "EMBEDDING"
)

// String representation for [fmt.Print].
func (f *EmbeddingsV1ResponseEmbeddingElementObject) String() string {
	return string(*f)
}

type QueryEndpointResponseObject string

const (
	QueryEndpointResponseObject_QueryEndpointResponseObjectUnspecified QueryEndpointResponseObject = "QUERY_ENDPOINT_RESPONSE_OBJECT_UNSPECIFIED"
	QueryEndpointResponseObject_TextCompletion                         QueryEndpointResponseObject = "TEXT_COMPLETION"
	QueryEndpointResponseObject_ChatCompletion                         QueryEndpointResponseObject = "CHAT_COMPLETION"
	QueryEndpointResponseObject_List                                   QueryEndpointResponseObject = "LIST"
)

// String representation for [fmt.Print].
func (f *QueryEndpointResponseObject) String() string {
	return string(*f)
}

type ChatMessage struct {
	Role    *ChatMessageRole `json:"role"`
	Content *string          `json:"content"`
}

type DataframeSplitInput struct {
	Index   []int             `json:"index"`
	Columns []json.RawMessage `json:"columns"`
	Data    []json.RawMessage `json:"data"`
}

type EmbeddingsV1ResponseEmbeddingElement struct {
	Embedding []float64                                   `json:"embedding"`
	Index     *int                                        `json:"index"`
	Object    *EmbeddingsV1ResponseEmbeddingElementObject `json:"object"`
}

type ExternalModelUsageElement struct {
	PromptTokens     *int `json:"prompt_tokens"`
	CompletionTokens *int `json:"completion_tokens"`
	TotalTokens      *int `json:"total_tokens"`
}

type QueryEndpointInput struct {
	Name             *string              `json:"name"`
	Prompt           *json.RawMessage     `json:"prompt"`
	Input            *json.RawMessage     `json:"input"`
	Messages         []ChatMessage        `json:"messages"`
	Temperature      *float64             `json:"temperature"`
	Stop             []string             `json:"stop"`
	MaxTokens        *int                 `json:"max_tokens"`
	N                *int                 `json:"n"`
	Stream           *bool                `json:"stream"`
	ExtraParams      map[string]string    `json:"extra_params"`
	DataframeRecords []json.RawMessage    `json:"dataframe_records"`
	DataframeSplit   *DataframeSplitInput `json:"dataframe_split"`
	Instances        []json.RawMessage    `json:"instances"`
	Inputs           *json.RawMessage     `json:"inputs"`
	ClientRequestId  *string              `json:"client_request_id"`
	UsageContext     map[string]string    `json:"usage_context"`
}

type QueryEndpointInput_ExtraParamsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type QueryEndpointInput_UsageContextEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type QueryEndpointResponse struct {
	Choices         []V1ResponseChoiceElement              `json:"choices"`
	Data            []EmbeddingsV1ResponseEmbeddingElement `json:"data"`
	Model           *string                                `json:"model"`
	Usage           *ExternalModelUsageElement             `json:"usage"`
	Id              *string                                `json:"id"`
	Created         *int64                                 `json:"created"`
	Object          *QueryEndpointResponseObject           `json:"object"`
	Predictions     []json.RawMessage                      `json:"predictions"`
	Outputs         []json.RawMessage                      `json:"outputs"`
	ServedModelName *string                                `json:"served-model-name"`
}

type V1ResponseChoiceElement struct {
	Text         *string      `json:"text"`
	Message      *ChatMessage `json:"message"`
	Index        *int         `json:"index"`
	FinishReason *string      `json:"finishReason"`
	Logprobs     *int         `json:"logprobs"`
}
