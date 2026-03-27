package v1

import "encoding/json"

type Behavior string

const (
	Behavior_BehaviorUnspecified Behavior = "BEHAVIOR_UNSPECIFIED"
	Behavior_None                Behavior = "NONE"
	Behavior_Block               Behavior = "BLOCK"
	Behavior_Mask                Behavior = "MASK"
)

// String representation for [fmt.Print].
func (f *Behavior) String() string {
	return string(*f)
}

type FoundationModelProduct string

const (
	FoundationModelProduct_FoundationModelProductUnspecified FoundationModelProduct = "FOUNDATION_MODEL_PRODUCT_UNSPECIFIED"
	FoundationModelProduct_Ppt                               FoundationModelProduct = "PPT"
	FoundationModelProduct_Pt                                FoundationModelProduct = "PT"
	FoundationModelProduct_Batch                             FoundationModelProduct = "BATCH"
)

// String representation for [fmt.Print].
func (f *FoundationModelProduct) String() string {
	return string(*f)
}

type ServedModelDeploymentState string

const (
	ServedModelDeploymentState_DeploymentUnknown    ServedModelDeploymentState = "DEPLOYMENT_UNKNOWN"
	ServedModelDeploymentState_DeploymentCreating   ServedModelDeploymentState = "DEPLOYMENT_CREATING"
	ServedModelDeploymentState_DeploymentRecovering ServedModelDeploymentState = "DEPLOYMENT_RECOVERING"
	ServedModelDeploymentState_DeploymentReady      ServedModelDeploymentState = "DEPLOYMENT_READY"
	ServedModelDeploymentState_DeploymentFailed     ServedModelDeploymentState = "DEPLOYMENT_FAILED"
	ServedModelDeploymentState_DeploymentAborted    ServedModelDeploymentState = "DEPLOYMENT_ABORTED"
	ServedModelDeploymentState_DeploymentStopped    ServedModelDeploymentState = "DEPLOYMENT_STOPPED"
)

// String representation for [fmt.Print].
func (f *ServedModelDeploymentState) String() string {
	return string(*f)
}

type ServingEndpointDetailedPermissionLevel string

const (
	ServingEndpointDetailedPermissionLevel_CanManage ServingEndpointDetailedPermissionLevel = "CAN_MANAGE"
	ServingEndpointDetailedPermissionLevel_CanQuery  ServingEndpointDetailedPermissionLevel = "CAN_QUERY"
	ServingEndpointDetailedPermissionLevel_CanView   ServingEndpointDetailedPermissionLevel = "CAN_VIEW"
)

// String representation for [fmt.Print].
func (f *ServingEndpointDetailedPermissionLevel) String() string {
	return string(*f)
}

type ExternalFunctionRequest_HttpMethod string

const (
	ExternalFunctionRequest_HttpMethod_HttpMethodUnspecified ExternalFunctionRequest_HttpMethod = "HTTP_METHOD_UNSPECIFIED"
	ExternalFunctionRequest_HttpMethod_Get                   ExternalFunctionRequest_HttpMethod = "GET"
	ExternalFunctionRequest_HttpMethod_Post                  ExternalFunctionRequest_HttpMethod = "POST"
	ExternalFunctionRequest_HttpMethod_Put                   ExternalFunctionRequest_HttpMethod = "PUT"
	ExternalFunctionRequest_HttpMethod_Delete                ExternalFunctionRequest_HttpMethod = "DELETE"
	ExternalFunctionRequest_HttpMethod_Patch                 ExternalFunctionRequest_HttpMethod = "PATCH"
)

// String representation for [fmt.Print].
func (f *ExternalFunctionRequest_HttpMethod) String() string {
	return string(*f)
}

type InferenceEndpointState_ConfigUpdateState string

const (
	InferenceEndpointState_ConfigUpdateState_ConfigUpdateStateUnspecified InferenceEndpointState_ConfigUpdateState = "CONFIG_UPDATE_STATE_UNSPECIFIED"
	InferenceEndpointState_ConfigUpdateState_NotUpdating                  InferenceEndpointState_ConfigUpdateState = "NOT_UPDATING"
	InferenceEndpointState_ConfigUpdateState_InProgress                   InferenceEndpointState_ConfigUpdateState = "IN_PROGRESS"
	InferenceEndpointState_ConfigUpdateState_UpdateFailed                 InferenceEndpointState_ConfigUpdateState = "UPDATE_FAILED"
	InferenceEndpointState_ConfigUpdateState_UpdateCanceled               InferenceEndpointState_ConfigUpdateState = "UPDATE_CANCELED"
)

// String representation for [fmt.Print].
func (f *InferenceEndpointState_ConfigUpdateState) String() string {
	return string(*f)
}

type InferenceEndpointState_ReadyState string

const (
	InferenceEndpointState_ReadyState_ReadyStateUnspecified InferenceEndpointState_ReadyState = "READY_STATE_UNSPECIFIED"
	InferenceEndpointState_ReadyState_Ready                 InferenceEndpointState_ReadyState = "READY"
	InferenceEndpointState_ReadyState_NotReady              InferenceEndpointState_ReadyState = "NOT_READY"
)

// String representation for [fmt.Print].
func (f *InferenceEndpointState_ReadyState) String() string {
	return string(*f)
}

type Ai21LabsConfig struct {
	Ai21labsApiKey          *string `json:"ai21labs_api_key"`
	Ai21labsApiKeyPlaintext *string `json:"ai21labs_api_key_plaintext"`
}

type AiGatewayConfig struct {
	UsageTrackingConfig  *UsageTrackingConfig  `json:"usage_tracking_config"`
	InferenceTableConfig *InferenceTableConfig `json:"inference_table_config"`
	RateLimits           []AiGatewayRateLimit  `json:"rate_limits"`
	Guardrails           *AiGuardrails         `json:"guardrails"`
	FallbackConfig       *FallbackConfig       `json:"fallback_config"`
}

type AiGatewayRateLimit struct {
	Calls         *int64  `json:"calls"`
	Key           *string `json:"key"`
	RenewalPeriod *string `json:"renewal_period"`
	Principal     *string `json:"principal"`
	Tokens        *int64  `json:"tokens"`
}

type AiGuardrailParameters struct {
	Safety          *bool        `json:"safety"`
	Pii             *PiiSettings `json:"pii"`
	ValidTopics     []string     `json:"valid_topics"`
	InvalidKeywords []string     `json:"invalid_keywords"`
}

type AiGuardrails struct {
	Input  *AiGuardrailParameters `json:"input"`
	Output *AiGuardrailParameters `json:"output"`
}

type AmazonBedrockConfig struct {
	AwsRegion                   *string `json:"aws_region"`
	AwsAccessKeyId              *string `json:"aws_access_key_id"`
	AwsSecretAccessKey          *string `json:"aws_secret_access_key"`
	BedrockProvider             *string `json:"bedrock_provider"`
	AwsAccessKeyIdPlaintext     *string `json:"aws_access_key_id_plaintext"`
	AwsSecretAccessKeyPlaintext *string `json:"aws_secret_access_key_plaintext"`
	InstanceProfileArn          *string `json:"instance_profile_arn"`
	UcServiceCredentialName     *string `json:"uc_service_credential_name"`
}

type AnthropicConfig struct {
	AnthropicApiKey          *string `json:"anthropic_api_key"`
	AnthropicApiKeyPlaintext *string `json:"anthropic_api_key_plaintext"`
}

type ApiKeyAuth struct {
	Key            *string `json:"key"`
	Value          *string `json:"value"`
	ValuePlaintext *string `json:"value_plaintext"`
}

type AutoCaptureConfig struct {
	CatalogName     *string           `json:"catalog_name"`
	SchemaName      *string           `json:"schema_name"`
	TableNamePrefix *string           `json:"table_name_prefix"`
	State           *AutoCaptureState `json:"state"`
	Enabled         *bool             `json:"enabled"`
}

type AutoCaptureState struct {
	PayloadTable *PayloadTable `json:"payload_table"`
}

type BearerTokenAuth struct {
	Token          *string `json:"token"`
	TokenPlaintext *string `json:"token_plaintext"`
}

type CohereConfig struct {
	CohereApiKey          *string `json:"cohere_api_key"`
	CohereApiKeyPlaintext *string `json:"cohere_api_key_plaintext"`
	CohereApiBase         *string `json:"cohere_api_base"`
}

type Content struct {
	Schema *Schema `json:"schema"`
}

type CreateInferenceEndpoint struct {
	Name               *string             `json:"name"`
	Config             *EndpointCoreConfig `json:"config"`
	Tags               []EndpointTag       `json:"tags"`
	RouteOptimized     *bool               `json:"route_optimized"`
	RateLimits         []RateLimit         `json:"rate_limits"`
	AiGateway          *AiGatewayConfig    `json:"ai_gateway"`
	BudgetPolicyId     *string             `json:"budget_policy_id"`
	EmailNotifications *EmailNotifications `json:"email_notifications"`
	Description        *string             `json:"description"`
}

type CreatePtEndpoint struct {
	Name               *string               `json:"name"`
	Config             *PtEndpointCoreConfig `json:"config"`
	Tags               []EndpointTag         `json:"tags"`
	AiGateway          *AiGatewayConfig      `json:"ai_gateway"`
	BudgetPolicyId     *string               `json:"budget_policy_id"`
	EmailNotifications *EmailNotifications   `json:"email_notifications"`
}

// Configs needed to create a custom provider model route..
type CustomProviderConfig struct {
	CustomProviderUrl *string          `json:"custom_provider_url"`
	BearerTokenAuth   *BearerTokenAuth `json:"bearer_token_auth"`
	ApiKeyAuth        *ApiKeyAuth      `json:"api_key_auth"`
}

// Details necessary to query this object's API through the DataPlane APIs..
type DataPlaneInfo struct {
	EndpointUrl          *string `json:"endpoint_url"`
	AuthorizationDetails *string `json:"authorization_details"`
}

type DatabricksModelServingConfig struct {
	DatabricksApiToken          *string `json:"databricks_api_token"`
	DatabricksWorkspaceUrl      *string `json:"databricks_workspace_url"`
	DatabricksApiTokenPlaintext *string `json:"databricks_api_token_plaintext"`
}

type DeleteInferenceEndpoint struct {
	Name *string `json:"name"`
}

type DeleteInferenceEndpoint_Response struct {
}

type EmailNotifications struct {
	OnUpdateSuccess []string `json:"on_update_success"`
	OnUpdateFailure []string `json:"on_update_failure"`
}

type EndpointCoreConfig struct {
	ServedEntities    []ServedModel      `json:"served_entities"`
	ServedModels      []ServedModel      `json:"served_models"`
	TrafficConfig     *TrafficConfig     `json:"traffic_config"`
	AutoCaptureConfig *AutoCaptureConfig `json:"auto_capture_config"`
}

type EndpointCoreConfigOutput struct {
	ConfigVersion     *int64             `json:"config_version"`
	ServedEntities    []ServedModel      `json:"served_entities"`
	ServedModels      []ServedModel      `json:"served_models"`
	TrafficConfig     *TrafficConfig     `json:"traffic_config"`
	AutoCaptureConfig *AutoCaptureConfig `json:"auto_capture_config"`
}

type EndpointCoreConfigSummary struct {
	ServedEntities []ServedModelLite `json:"served_entities"`
	ServedModels   []ServedModelLite `json:"served_models"`
}

type EndpointTag struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Simple Proto message for testing.
type ExternalFunctionRequest struct {
	ConnectionName *string                             `json:"connection_name"`
	Method         *ExternalFunctionRequest_HttpMethod `json:"method"`
	Path           *string                             `json:"path"`
	Json           *string                             `json:"json"`
	Headers        *string                             `json:"headers"`
	Params         *string                             `json:"params"`
	SubDomain      *string                             `json:"sub_domain"`
}

type ExternalFunctionResponse struct {
}

type ExternalModel struct {
	Provider                     *string                       `json:"provider"`
	Name                         *string                       `json:"name"`
	Task                         *string                       `json:"task"`
	Ai21labsConfig               *Ai21LabsConfig               `json:"ai21labs_config"`
	AnthropicConfig              *AnthropicConfig              `json:"anthropic_config"`
	AmazonBedrockConfig          *AmazonBedrockConfig          `json:"amazon_bedrock_config"`
	CohereConfig                 *CohereConfig                 `json:"cohere_config"`
	GoogleCloudVertexAiConfig    *GoogleCloudVertexAiConfig    `json:"google_cloud_vertex_ai_config"`
	DatabricksModelServingConfig *DatabricksModelServingConfig `json:"databricks_model_serving_config"`
	OpenaiConfig                 *OpenAiConfig                 `json:"openai_config"`
	PalmConfig                   *PaLmConfig                   `json:"palm_config"`
	CustomProviderConfig         *CustomProviderConfig         `json:"custom_provider_config"`
}

type FallbackConfig struct {
	Enabled *bool `json:"enabled"`
}

// All fields are not sensitive as they are hard-coded in the system and made
// available to customers..
type FoundationModel struct {
	Name        *string `json:"name"`
	DisplayName *string `json:"display_name"`
	Docs        *string `json:"docs"`
	Description *string `json:"description"`
}

type GetInferenceEndpoint struct {
	Name *string `json:"name"`
}

type GetInferenceEndpointSchema struct {
	Name *string `json:"name"`
}

// The top level proto message that represents an OpenAPI 3.0 document..
type GetOpenApiResponse struct {
	Openapi *string             `json:"openapi"`
	Info    *Info               `json:"info"`
	Servers []Server            `json:"servers"`
	Paths   map[string]PathItem `json:"paths"`
}

type GetOpenApiResponse_PathsEntry struct {
	Key   *string   `json:"key"`
	Value *PathItem `json:"value"`
}

type GoogleCloudVertexAiConfig struct {
	PrivateKey          *string `json:"private_key"`
	ProjectId           *string `json:"project_id"`
	Region              *string `json:"region"`
	PrivateKeyPlaintext *string `json:"private_key_plaintext"`
}

type InferenceEndpoint struct {
	Name                 *string                    `json:"name"`
	Creator              *string                    `json:"creator"`
	CreationTimestamp    *int64                     `json:"creation_timestamp"`
	LastUpdatedTimestamp *int64                     `json:"last_updated_timestamp"`
	State                *InferenceEndpointState    `json:"state"`
	Config               *EndpointCoreConfigSummary `json:"config"`
	Tags                 []EndpointTag              `json:"tags"`
	Id                   *string                    `json:"id"`
	Task                 *string                    `json:"task"`
	AiGateway            *AiGatewayConfig           `json:"ai_gateway"`
	BudgetPolicyId       *string                    `json:"budget_policy_id"`
	Description          *string                    `json:"description"`
	UsagePolicyId        *string                    `json:"usage_policy_id"`
}

type InferenceEndpointDetailed struct {
	Name                 *string                                 `json:"name"`
	Creator              *string                                 `json:"creator"`
	CreationTimestamp    *int64                                  `json:"creation_timestamp"`
	LastUpdatedTimestamp *int64                                  `json:"last_updated_timestamp"`
	State                *InferenceEndpointState                 `json:"state"`
	Config               *EndpointCoreConfigOutput               `json:"config"`
	PendingConfig        *PendingConfig                          `json:"pending_config"`
	Id                   *string                                 `json:"id"`
	PermissionLevel      *ServingEndpointDetailedPermissionLevel `json:"permission_level"`
	Tags                 []EndpointTag                           `json:"tags"`
	Task                 *string                                 `json:"task"`
	RouteOptimized       *bool                                   `json:"route_optimized"`
	EndpointUrl          *string                                 `json:"endpoint_url"`
	DataPlaneInfo        *ModelDataPlaneInfo                     `json:"data_plane_info"`
	AiGateway            *AiGatewayConfig                        `json:"ai_gateway"`
	BudgetPolicyId       *string                                 `json:"budget_policy_id"`
	EmailNotifications   *EmailNotifications                     `json:"email_notifications"`
	Description          *string                                 `json:"description"`
}

type InferenceEndpointState struct {
	Ready        *InferenceEndpointState_ReadyState        `json:"ready"`
	ConfigUpdate *InferenceEndpointState_ConfigUpdateState `json:"config_update"`
}

type InferenceTableConfig struct {
	CatalogName     *string `json:"catalog_name"`
	SchemaName      *string `json:"schema_name"`
	TableNamePrefix *string `json:"table_name_prefix"`
	Enabled         *bool   `json:"enabled"`
}

// The object provides metadata about the API. The metadata MAY be used by the
// clients if needed, and MAY be presented in editing or documentation
// generation tools for convenience..
type Info struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Version     *string `json:"version"`
	Summary     *string `json:"summary"`
}

type ListFoundationModelEndpoints struct {
	Product *FoundationModelProduct `json:"product"`
}

type ListFoundationModelEndpoints_Response struct {
	Endpoints []InferenceEndpoint `json:"endpoints"`
}

type ListInferenceEndpoints struct {
}

type ListInferenceEndpoints_Response struct {
	Endpoints []InferenceEndpoint `json:"endpoints"`
}

// A representation of all DataPlaneInfo for operations that can be done on a
// model through Data Plane APIs..
type ModelDataPlaneInfo struct {
	QueryInfo *DataPlaneInfo `json:"query_info"`
}

// Configs needed to create an OpenAI model route..
type OpenAiConfig struct {
	OpenaiApiKey                        *string `json:"openai_api_key"`
	OpenaiApiType                       *string `json:"openai_api_type"`
	OpenaiApiBase                       *string `json:"openai_api_base"`
	OpenaiApiVersion                    *string `json:"openai_api_version"`
	OpenaiDeploymentName                *string `json:"openai_deployment_name"`
	OpenaiOrganization                  *string `json:"openai_organization"`
	MicrosoftEntraTenantId              *string `json:"microsoft_entra_tenant_id"`
	MicrosoftEntraClientId              *string `json:"microsoft_entra_client_id"`
	MicrosoftEntraClientSecret          *string `json:"microsoft_entra_client_secret"`
	OpenaiApiKeyPlaintext               *string `json:"openai_api_key_plaintext"`
	MicrosoftEntraClientSecretPlaintext *string `json:"microsoft_entra_client_secret_plaintext"`
}

// Describes a single API operation on a path..
type Operation struct {
	Tags        []string            `json:"tags"`
	Summary     *string             `json:"summary"`
	Description *string             `json:"description"`
	RequestBody *RequestBody        `json:"requestBody"`
	Responses   map[string]Response `json:"responses"`
}

type Operation_ResponsesEntry struct {
	Key   *string   `json:"key"`
	Value *Response `json:"value"`
}

type PaLmConfig struct {
	PalmApiKey          *string `json:"palm_api_key"`
	PalmApiKeyPlaintext *string `json:"palm_api_key_plaintext"`
}

type PatchInferenceEndpointBudgetPolicy struct {
	Name           *string `json:"name"`
	BudgetPolicyId *string `json:"budget_policy_id"`
}

type PatchInferenceEndpointBudgetPolicy_Response struct {
	BudgetPolicyId          *string `json:"budget_policy_id"`
	EffectiveBudgetPolicyId *string `json:"effective_budget_policy_id"`
}

type PatchInferenceEndpointDescription struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type PatchInferenceEndpointDescription_Response struct {
	Description *string `json:"description"`
}

type PatchInferenceEndpointTags struct {
	Name       *string       `json:"name"`
	AddTags    []EndpointTag `json:"add_tags"`
	DeleteTags []string      `json:"delete_tags"`
}

type PatchInferenceEndpointTags_Response struct {
	Tags []EndpointTag `json:"tags"`
}

type PatchInferenceEndpointUsagePolicy struct {
	Name          *string `json:"name"`
	UsagePolicyId *string `json:"usage_policy_id"`
}

type PatchInferenceEndpointUsagePolicy_Response struct {
	UsagePolicyId          *string `json:"usage_policy_id"`
	EffectiveUsagePolicyId *string `json:"effective_usage_policy_id"`
}

// Describes the operations available on a single path. A Path Item MAY be
// empty, due to ACL constraints. The path itself is still exposed to the
// documentation viewer but they will not know which operations and parameters
// are available..
type PathItem struct {
	Summary     *string    `json:"summary"`
	Description *string    `json:"description"`
	Get         *Operation `json:"get"`
	Put         *Operation `json:"put"`
	Post        *Operation `json:"post"`
	Delete      *Operation `json:"delete"`
	Options     *Operation `json:"options"`
	Head        *Operation `json:"head"`
	Patch       *Operation `json:"patch"`
	Trace       *Operation `json:"trace"`
}

type PayloadTable struct {
	Name          *string `json:"name"`
	Status        *string `json:"status"`
	StatusMessage *string `json:"status_message"`
}

type PendingConfig struct {
	ServedEntities    []ServedModel      `json:"served_entities"`
	ServedModels      []ServedModel      `json:"served_models"`
	TrafficConfig     *TrafficConfig     `json:"traffic_config"`
	ConfigVersion     *int               `json:"config_version"`
	StartTime         *int64             `json:"start_time"`
	AutoCaptureConfig *AutoCaptureConfig `json:"auto_capture_config"`
}

type PiiSettings struct {
	Behavior *Behavior `json:"behavior"`
}

type PtEndpointCoreConfig struct {
	ServedEntities []PtServedModel `json:"served_entities"`
	TrafficConfig  *TrafficConfig  `json:"traffic_config"`
}

type PtServedModel struct {
	Name                  *string `json:"name"`
	EntityName            *string `json:"entity_name"`
	EntityVersion         *string `json:"entity_version"`
	ProvisionedModelUnits *int64  `json:"provisioned_model_units"`
	BurstScalingEnabled   *bool   `json:"burst_scaling_enabled"`
}

type PutInferenceEndpointAiGateway struct {
	Name                 *string               `json:"name"`
	UsageTrackingConfig  *UsageTrackingConfig  `json:"usage_tracking_config"`
	InferenceTableConfig *InferenceTableConfig `json:"inference_table_config"`
	RateLimits           []AiGatewayRateLimit  `json:"rate_limits"`
	Guardrails           *AiGuardrails         `json:"guardrails"`
	FallbackConfig       *FallbackConfig       `json:"fallback_config"`
}

type PutInferenceEndpointAiGateway_Response struct {
	UsageTrackingConfig  *UsageTrackingConfig  `json:"usage_tracking_config"`
	InferenceTableConfig *InferenceTableConfig `json:"inference_table_config"`
	RateLimits           []AiGatewayRateLimit  `json:"rate_limits"`
	Guardrails           *AiGuardrails         `json:"guardrails"`
	FallbackConfig       *FallbackConfig       `json:"fallback_config"`
}

type PutInferenceEndpointConfig struct {
	Name              *string            `json:"name"`
	ServedEntities    []ServedModel      `json:"served_entities"`
	ServedModels      []ServedModel      `json:"served_models"`
	TrafficConfig     *TrafficConfig     `json:"traffic_config"`
	AutoCaptureConfig *AutoCaptureConfig `json:"auto_capture_config"`
}

type PutInferenceEndpointRateLimits struct {
	Name       *string     `json:"name"`
	RateLimits []RateLimit `json:"rate_limits"`
}

type PutInferenceEndpointRateLimits_Response struct {
	RateLimits []RateLimit `json:"rate_limits"`
}

type PutPtEndpointConfig struct {
	Name   *string               `json:"name"`
	Config *PtEndpointCoreConfig `json:"config"`
}

type RateLimit struct {
	Calls         *int64  `json:"calls"`
	Key           *string `json:"key"`
	RenewalPeriod *string `json:"renewal_period"`
}

// Describes a single request body..
type RequestBody struct {
	Description *string            `json:"description"`
	Content     map[string]Content `json:"content"`
}

type RequestBody_ContentEntry struct {
	Key   *string  `json:"key"`
	Value *Content `json:"value"`
}

type Response struct {
	Description *string            `json:"description"`
	Content     map[string]Content `json:"content"`
}

type Response_ContentEntry struct {
	Key   *string  `json:"key"`
	Value *Content `json:"value"`
}

type Route struct {
	ServedModelName   *string `json:"served_model_name"`
	TrafficPercentage *int    `json:"traffic_percentage"`
	ServedEntityName  *string `json:"served_entity_name"`
}

// An object representing the request or response body..
type Schema struct {
	Description          *string           `json:"description"`
	Required             []string          `json:"required"`
	Type                 *string           `json:"type"`
	Format               *string           `json:"format"`
	Example              *json.RawMessage  `json:"example"`
	AllOf                []Schema          `json:"allOf"`
	OneOf                []Schema          `json:"oneOf"`
	AnyOf                []Schema          `json:"anyOf"`
	Items                *Schema           `json:"items"`
	Properties           map[string]Schema `json:"properties"`
	AdditionalProperties *Schema           `json:"additionalProperties"`
	Enum                 []string          `json:"enum"`
	PrefixItems          []Schema          `json:"prefixItems"`
	Nullable             *bool             `json:"nullable"`
	Default              *string           `json:"default"`
	MinItems             *int64            `json:"minItems"`
	MaxItems             *int64            `json:"maxItems"`
	Examples             []json.RawMessage `json:"examples"`
}

type Schema_PropertiesEntry struct {
	Key   *string `json:"key"`
	Value *Schema `json:"value"`
}

type ServedModel struct {
	Name                      *string           `json:"name"`
	ExternalModel             *ExternalModel    `json:"external_model"`
	EntityName                *string           `json:"entity_name"`
	EntityVersion             *string           `json:"entity_version"`
	MinProvisionedThroughput  *int              `json:"min_provisioned_throughput"`
	MaxProvisionedThroughput  *int              `json:"max_provisioned_throughput"`
	MinProvisionedConcurrency *int              `json:"min_provisioned_concurrency"`
	MaxProvisionedConcurrency *int              `json:"max_provisioned_concurrency"`
	WorkloadSize              *string           `json:"workload_size"`
	ProvisionedModelUnits     *int64            `json:"provisioned_model_units"`
	BurstScalingEnabled       *bool             `json:"burst_scaling_enabled"`
	ScaleToZeroEnabled        *bool             `json:"scale_to_zero_enabled"`
	ModelName                 *string           `json:"model_name"`
	ModelVersion              *string           `json:"model_version"`
	EnvironmentVars           map[string]string `json:"environment_vars"`
	InstanceProfileArn        *string           `json:"instance_profile_arn"`
	FoundationModel           *FoundationModel  `json:"foundation_model"`
	State                     *ServedModelState `json:"state"`
	Creator                   *string           `json:"creator"`
	CreationTimestamp         *int64            `json:"creation_timestamp"`
}

type ServedModel_EnvironmentVarsEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type ServedModelLite struct {
	Name            *string          `json:"name"`
	ModelName       *string          `json:"model_name"`
	EntityName      *string          `json:"entity_name"`
	ModelVersion    *string          `json:"model_version"`
	EntityVersion   *string          `json:"entity_version"`
	ExternalModel   *ExternalModel   `json:"external_model"`
	FoundationModel *FoundationModel `json:"foundation_model"`
}

type ServedModelState struct {
	Deployment             *ServedModelDeploymentState `json:"deployment"`
	DeploymentStateMessage *string                     `json:"deployment_state_message"`
}

// An object representing a Server..
type Server struct {
	Url         *string `json:"url"`
	Description *string `json:"description"`
}

type TrafficConfig struct {
	Routes []Route `json:"routes"`
}

type UpdateInferenceEndpointNotifications struct {
	Name               *string             `json:"name"`
	EmailNotifications *EmailNotifications `json:"email_notifications"`
}

type UpdateInferenceEndpointNotifications_Response struct {
	Name               *string             `json:"name"`
	EmailNotifications *EmailNotifications `json:"email_notifications"`
}

type UsageTrackingConfig struct {
	Enabled *bool `json:"enabled"`
}
