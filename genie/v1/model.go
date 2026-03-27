package v1

type ColumnTypeName string

const (
	ColumnTypeName_Boolean         ColumnTypeName = "BOOLEAN"
	ColumnTypeName_Byte            ColumnTypeName = "BYTE"
	ColumnTypeName_Short           ColumnTypeName = "SHORT"
	ColumnTypeName_Int             ColumnTypeName = "INT"
	ColumnTypeName_Long            ColumnTypeName = "LONG"
	ColumnTypeName_Float           ColumnTypeName = "FLOAT"
	ColumnTypeName_Double          ColumnTypeName = "DOUBLE"
	ColumnTypeName_Date            ColumnTypeName = "DATE"
	ColumnTypeName_Timestamp       ColumnTypeName = "TIMESTAMP"
	ColumnTypeName_String          ColumnTypeName = "STRING"
	ColumnTypeName_Binary          ColumnTypeName = "BINARY"
	ColumnTypeName_Decimal         ColumnTypeName = "DECIMAL"
	ColumnTypeName_Interval        ColumnTypeName = "INTERVAL"
	ColumnTypeName_Array           ColumnTypeName = "ARRAY"
	ColumnTypeName_Struct          ColumnTypeName = "STRUCT"
	ColumnTypeName_Map             ColumnTypeName = "MAP"
	ColumnTypeName_Char            ColumnTypeName = "CHAR"
	ColumnTypeName_Null            ColumnTypeName = "NULL"
	ColumnTypeName_UserDefinedType ColumnTypeName = "USER_DEFINED_TYPE"
	ColumnTypeName_TimestampNtz    ColumnTypeName = "TIMESTAMP_NTZ"
	ColumnTypeName_Variant         ColumnTypeName = "VARIANT"
	ColumnTypeName_Geometry        ColumnTypeName = "GEOMETRY"
	ColumnTypeName_Geography       ColumnTypeName = "GEOGRAPHY"
	ColumnTypeName_Time            ColumnTypeName = "TIME"
	ColumnTypeName_TableType       ColumnTypeName = "TABLE_TYPE"
)

// String representation for [fmt.Print].
func (f *ColumnTypeName) String() string {
	return string(*f)
}

type ErrorCode string

const (
	ErrorCode_Unknown                          ErrorCode = "UNKNOWN"
	ErrorCode_InternalError                    ErrorCode = "INTERNAL_ERROR"
	ErrorCode_TemporarilyUnavailable           ErrorCode = "TEMPORARILY_UNAVAILABLE"
	ErrorCode_IoError                          ErrorCode = "IO_ERROR"
	ErrorCode_BadRequest                       ErrorCode = "BAD_REQUEST"
	ErrorCode_ServiceUnderMaintenance          ErrorCode = "SERVICE_UNDER_MAINTENANCE"
	ErrorCode_WorkspaceTemporarilyUnavailable  ErrorCode = "WORKSPACE_TEMPORARILY_UNAVAILABLE"
	ErrorCode_DeadlineExceeded                 ErrorCode = "DEADLINE_EXCEEDED"
	ErrorCode_Cancelled                        ErrorCode = "CANCELLED"
	ErrorCode_ResourceExhausted                ErrorCode = "RESOURCE_EXHAUSTED"
	ErrorCode_Aborted                          ErrorCode = "ABORTED"
	ErrorCode_NotFound                         ErrorCode = "NOT_FOUND"
	ErrorCode_AlreadyExists                    ErrorCode = "ALREADY_EXISTS"
	ErrorCode_Unauthenticated                  ErrorCode = "UNAUTHENTICATED"
	ErrorCode_Unavailable                      ErrorCode = "UNAVAILABLE"
	ErrorCode_InvalidParameterValue            ErrorCode = "INVALID_PARAMETER_VALUE"
	ErrorCode_EndpointNotFound                 ErrorCode = "ENDPOINT_NOT_FOUND"
	ErrorCode_MalformedRequest                 ErrorCode = "MALFORMED_REQUEST"
	ErrorCode_InvalidState                     ErrorCode = "INVALID_STATE"
	ErrorCode_PermissionDenied                 ErrorCode = "PERMISSION_DENIED"
	ErrorCode_FeatureDisabled                  ErrorCode = "FEATURE_DISABLED"
	ErrorCode_CustomerUnauthorized             ErrorCode = "CUSTOMER_UNAUTHORIZED"
	ErrorCode_RequestLimitExceeded             ErrorCode = "REQUEST_LIMIT_EXCEEDED"
	ErrorCode_ResourceConflict                 ErrorCode = "RESOURCE_CONFLICT"
	ErrorCode_UnparseableHttpError             ErrorCode = "UNPARSEABLE_HTTP_ERROR"
	ErrorCode_NotImplemented                   ErrorCode = "NOT_IMPLEMENTED"
	ErrorCode_DataLoss                         ErrorCode = "DATA_LOSS"
	ErrorCode_InvalidStateTransition           ErrorCode = "INVALID_STATE_TRANSITION"
	ErrorCode_CouldNotAcquireLock              ErrorCode = "COULD_NOT_ACQUIRE_LOCK"
	ErrorCode_ResourceAlreadyExists            ErrorCode = "RESOURCE_ALREADY_EXISTS"
	ErrorCode_ResourceDoesNotExist             ErrorCode = "RESOURCE_DOES_NOT_EXIST"
	ErrorCode_QuotaExceeded                    ErrorCode = "QUOTA_EXCEEDED"
	ErrorCode_MaxBlockSizeExceeded             ErrorCode = "MAX_BLOCK_SIZE_EXCEEDED"
	ErrorCode_MaxReadSizeExceeded              ErrorCode = "MAX_READ_SIZE_EXCEEDED"
	ErrorCode_PartialDelete                    ErrorCode = "PARTIAL_DELETE"
	ErrorCode_MaxListSizeExceeded              ErrorCode = "MAX_LIST_SIZE_EXCEEDED"
	ErrorCode_DryRunFailed                     ErrorCode = "DRY_RUN_FAILED"
	ErrorCode_ResourceLimitExceeded            ErrorCode = "RESOURCE_LIMIT_EXCEEDED"
	ErrorCode_DirectoryNotEmpty                ErrorCode = "DIRECTORY_NOT_EMPTY"
	ErrorCode_DirectoryProtected               ErrorCode = "DIRECTORY_PROTECTED"
	ErrorCode_MaxNotebookSizeExceeded          ErrorCode = "MAX_NOTEBOOK_SIZE_EXCEEDED"
	ErrorCode_MaxChildNodeSizeExceeded         ErrorCode = "MAX_CHILD_NODE_SIZE_EXCEEDED"
	ErrorCode_SearchQueryTooLong               ErrorCode = "SEARCH_QUERY_TOO_LONG"
	ErrorCode_SearchQueryTooShort              ErrorCode = "SEARCH_QUERY_TOO_SHORT"
	ErrorCode_ManagedResourceGroupDoesNotExist ErrorCode = "MANAGED_RESOURCE_GROUP_DOES_NOT_EXIST"
	ErrorCode_PermissionNotPropagated          ErrorCode = "PERMISSION_NOT_PROPAGATED"
	ErrorCode_DeploymentTimeout                ErrorCode = "DEPLOYMENT_TIMEOUT"
	ErrorCode_GitConflict                      ErrorCode = "GIT_CONFLICT"
	ErrorCode_GitUnknownRef                    ErrorCode = "GIT_UNKNOWN_REF"
	ErrorCode_GitSensitiveTokenDetected        ErrorCode = "GIT_SENSITIVE_TOKEN_DETECTED"
	ErrorCode_GitUrlNotOnAllowList             ErrorCode = "GIT_URL_NOT_ON_ALLOW_LIST"
	ErrorCode_GitRemoteError                   ErrorCode = "GIT_REMOTE_ERROR"
	ErrorCode_ProjectsOperationTimeout         ErrorCode = "PROJECTS_OPERATION_TIMEOUT"
	ErrorCode_IpynbFileInRepo                  ErrorCode = "IPYNB_FILE_IN_REPO"
	ErrorCode_InsecurePartnerResponse          ErrorCode = "INSECURE_PARTNER_RESPONSE"
	ErrorCode_MalformedPartnerResponse         ErrorCode = "MALFORMED_PARTNER_RESPONSE"
	ErrorCode_MetastoreDoesNotExist            ErrorCode = "METASTORE_DOES_NOT_EXIST"
	ErrorCode_DacDoesNotExist                  ErrorCode = "DAC_DOES_NOT_EXIST"
	ErrorCode_CatalogDoesNotExist              ErrorCode = "CATALOG_DOES_NOT_EXIST"
	ErrorCode_SchemaDoesNotExist               ErrorCode = "SCHEMA_DOES_NOT_EXIST"
	ErrorCode_TableDoesNotExist                ErrorCode = "TABLE_DOES_NOT_EXIST"
	ErrorCode_ShareDoesNotExist                ErrorCode = "SHARE_DOES_NOT_EXIST"
	ErrorCode_RecipientDoesNotExist            ErrorCode = "RECIPIENT_DOES_NOT_EXIST"
	ErrorCode_StorageCredentialDoesNotExist    ErrorCode = "STORAGE_CREDENTIAL_DOES_NOT_EXIST"
	ErrorCode_ExternalLocationDoesNotExist     ErrorCode = "EXTERNAL_LOCATION_DOES_NOT_EXIST"
	ErrorCode_PrincipalDoesNotExist            ErrorCode = "PRINCIPAL_DOES_NOT_EXIST"
	ErrorCode_ProviderDoesNotExist             ErrorCode = "PROVIDER_DOES_NOT_EXIST"
	ErrorCode_MetastoreAlreadyExists           ErrorCode = "METASTORE_ALREADY_EXISTS"
	ErrorCode_DacAlreadyExists                 ErrorCode = "DAC_ALREADY_EXISTS"
	ErrorCode_CatalogAlreadyExists             ErrorCode = "CATALOG_ALREADY_EXISTS"
	ErrorCode_SchemaAlreadyExists              ErrorCode = "SCHEMA_ALREADY_EXISTS"
	ErrorCode_TableAlreadyExists               ErrorCode = "TABLE_ALREADY_EXISTS"
	ErrorCode_ShareAlreadyExists               ErrorCode = "SHARE_ALREADY_EXISTS"
	ErrorCode_RecipientAlreadyExists           ErrorCode = "RECIPIENT_ALREADY_EXISTS"
	ErrorCode_StorageCredentialAlreadyExists   ErrorCode = "STORAGE_CREDENTIAL_ALREADY_EXISTS"
	ErrorCode_ExternalLocationAlreadyExists    ErrorCode = "EXTERNAL_LOCATION_ALREADY_EXISTS"
	ErrorCode_ProviderAlreadyExists            ErrorCode = "PROVIDER_ALREADY_EXISTS"
	ErrorCode_CatalogNotEmpty                  ErrorCode = "CATALOG_NOT_EMPTY"
	ErrorCode_SchemaNotEmpty                   ErrorCode = "SCHEMA_NOT_EMPTY"
	ErrorCode_MetastoreNotEmpty                ErrorCode = "METASTORE_NOT_EMPTY"
	ErrorCode_ProviderShareNotAccessible       ErrorCode = "PROVIDER_SHARE_NOT_ACCESSIBLE"
)

// String representation for [fmt.Print].
func (f *ErrorCode) String() string {
	return string(*f)
}

type EvaluationStatusType string

const (
	EvaluationStatusType_EvaluationStatusTypeUnspecified EvaluationStatusType = "EVALUATION_STATUS_TYPE_UNSPECIFIED"
	EvaluationStatusType_Running                         EvaluationStatusType = "RUNNING"
	EvaluationStatusType_Done                            EvaluationStatusType = "DONE"
	EvaluationStatusType_NotStarted                      EvaluationStatusType = "NOT_STARTED"
	EvaluationStatusType_EvaluationFailed                EvaluationStatusType = "EVALUATION_FAILED"
	EvaluationStatusType_EvaluationCancelled             EvaluationStatusType = "EVALUATION_CANCELLED"
	EvaluationStatusType_EvaluationTimeout               EvaluationStatusType = "EVALUATION_TIMEOUT"
)

// String representation for [fmt.Print].
func (f *EvaluationStatusType) String() string {
	return string(*f)
}

type Format string

const (
	Format_FormatUnspecified Format = "FORMAT_UNSPECIFIED"
	Format_JsonArray         Format = "JSON_ARRAY"
	Format_ArrowStream       Format = "ARROW_STREAM"
	Format_Csv               Format = "CSV"
)

// String representation for [fmt.Print].
func (f *Format) String() string {
	return string(*f)
}

type GenieEvalAssessment string

const (
	GenieEvalAssessment_GenieEvalAssessmentUnspecified GenieEvalAssessment = "GENIE_EVAL_ASSESSMENT_UNSPECIFIED"
	GenieEvalAssessment_Good                           GenieEvalAssessment = "GOOD"
	GenieEvalAssessment_Bad                            GenieEvalAssessment = "BAD"
	GenieEvalAssessment_NeedsReview                    GenieEvalAssessment = "NEEDS_REVIEW"
)

// String representation for [fmt.Print].
func (f *GenieEvalAssessment) String() string {
	return string(*f)
}

type GenieEvalResponseType string

const (
	GenieEvalResponseType_GenieEvalResponseTypeUnspecified GenieEvalResponseType = "GENIE_EVAL_RESPONSE_TYPE_UNSPECIFIED"
	GenieEvalResponseType_Text                             GenieEvalResponseType = "TEXT"
	GenieEvalResponseType_Sql                              GenieEvalResponseType = "SQL"
)

// String representation for [fmt.Print].
func (f *GenieEvalResponseType) String() string {
	return string(*f)
}

type GenieFeedbackRating string

const (
	GenieFeedbackRating_GenieFeedbackRatingUnspecified GenieFeedbackRating = "GENIE_FEEDBACK_RATING_UNSPECIFIED"
	GenieFeedbackRating_Positive                       GenieFeedbackRating = "POSITIVE"
	GenieFeedbackRating_Negative                       GenieFeedbackRating = "NEGATIVE"
	GenieFeedbackRating_None                           GenieFeedbackRating = "NONE"
)

// String representation for [fmt.Print].
func (f *GenieFeedbackRating) String() string {
	return string(*f)
}

type NullValue string

const (
	NullValue_NullValue NullValue = "NULL_VALUE"
)

// String representation for [fmt.Print].
func (f *NullValue) String() string {
	return string(*f)
}

type ResponsePhase string

const (
	ResponsePhase_ResponsePhaseUnspecified ResponsePhase = "RESPONSE_PHASE_UNSPECIFIED"
	ResponsePhase_ResponsePhaseThinking    ResponsePhase = "RESPONSE_PHASE_THINKING"
	ResponsePhase_ResponsePhaseVerifying   ResponsePhase = "RESPONSE_PHASE_VERIFYING"
)

// String representation for [fmt.Print].
func (f *ResponsePhase) String() string {
	return string(*f)
}

type ScoreReason string

const (
	ScoreReason_ScoreReasonUnspecified                              ScoreReason = "SCORE_REASON_UNSPECIFIED"
	ScoreReason_EmptyResult                                         ScoreReason = "EMPTY_RESULT"
	ScoreReason_ResultMissingRows                                   ScoreReason = "RESULT_MISSING_ROWS"
	ScoreReason_ResultExtraRows                                     ScoreReason = "RESULT_EXTRA_ROWS"
	ScoreReason_ResultMissingColumns                                ScoreReason = "RESULT_MISSING_COLUMNS"
	ScoreReason_ResultExtraColumns                                  ScoreReason = "RESULT_EXTRA_COLUMNS"
	ScoreReason_SingleCellDifference                                ScoreReason = "SINGLE_CELL_DIFFERENCE"
	ScoreReason_EmptyGoodSql                                        ScoreReason = "EMPTY_GOOD_SQL"
	ScoreReason_ColumnTypeDifference                                ScoreReason = "COLUMN_TYPE_DIFFERENCE"
	ScoreReason_LlmJudgeMissingJoin                                 ScoreReason = "LLM_JUDGE_MISSING_JOIN"
	ScoreReason_LlmJudgeWrongFilter                                 ScoreReason = "LLM_JUDGE_WRONG_FILTER"
	ScoreReason_LlmJudgeWrongAggregation                            ScoreReason = "LLM_JUDGE_WRONG_AGGREGATION"
	ScoreReason_LlmJudgeWrongColumns                                ScoreReason = "LLM_JUDGE_WRONG_COLUMNS"
	ScoreReason_LlmJudgeSyntaxError                                 ScoreReason = "LLM_JUDGE_SYNTAX_ERROR"
	ScoreReason_LlmJudgeSemanticError                               ScoreReason = "LLM_JUDGE_SEMANTIC_ERROR"
	ScoreReason_LlmJudgeOther                                       ScoreReason = "LLM_JUDGE_OTHER"
	ScoreReason_LlmJudgeMissingOrIncorrectFilter                    ScoreReason = "LLM_JUDGE_MISSING_OR_INCORRECT_FILTER"
	ScoreReason_LlmJudgeIncompleteOrPartialOutput                   ScoreReason = "LLM_JUDGE_INCOMPLETE_OR_PARTIAL_OUTPUT"
	ScoreReason_LlmJudgeMisinterpretationOfUserRequest              ScoreReason = "LLM_JUDGE_MISINTERPRETATION_OF_USER_REQUEST"
	ScoreReason_LlmJudgeInstructionComplianceOrMissingBusinessLogic ScoreReason = "LLM_JUDGE_INSTRUCTION_COMPLIANCE_OR_MISSING_BUSINESS_LOGIC"
	ScoreReason_LlmJudgeIncorrectMetricCalculation                  ScoreReason = "LLM_JUDGE_INCORRECT_METRIC_CALCULATION"
	ScoreReason_LlmJudgeIncorrectTableOrFieldUsage                  ScoreReason = "LLM_JUDGE_INCORRECT_TABLE_OR_FIELD_USAGE"
	ScoreReason_LlmJudgeIncorrectFunctionUsage                      ScoreReason = "LLM_JUDGE_INCORRECT_FUNCTION_USAGE"
	ScoreReason_LlmJudgeMissingOrIncorrectJoin                      ScoreReason = "LLM_JUDGE_MISSING_OR_INCORRECT_JOIN"
	ScoreReason_LlmJudgeMissingOrIncorrectAggregation               ScoreReason = "LLM_JUDGE_MISSING_OR_INCORRECT_AGGREGATION"
	ScoreReason_LlmJudgeFormattingError                             ScoreReason = "LLM_JUDGE_FORMATTING_ERROR"
)

// String representation for [fmt.Print].
func (f *ScoreReason) String() string {
	return string(*f)
}

type TextAttachmentPurpose string

const (
	TextAttachmentPurpose_TextAttachmentPurposeUnspecified TextAttachmentPurpose = "TEXT_ATTACHMENT_PURPOSE_UNSPECIFIED"
	TextAttachmentPurpose_FollowUpQuestion                 TextAttachmentPurpose = "FOLLOW_UP_QUESTION"
)

// String representation for [fmt.Print].
func (f *TextAttachmentPurpose) String() string {
	return string(*f)
}

type ThoughtType string

const (
	ThoughtType_ThoughtTypeUnspecified   ThoughtType = "THOUGHT_TYPE_UNSPECIFIED"
	ThoughtType_ThoughtTypeDescription   ThoughtType = "THOUGHT_TYPE_DESCRIPTION"
	ThoughtType_ThoughtTypeUnderstanding ThoughtType = "THOUGHT_TYPE_UNDERSTANDING"
	ThoughtType_ThoughtTypeDataSourcing  ThoughtType = "THOUGHT_TYPE_DATA_SOURCING"
	ThoughtType_ThoughtTypeInstructions  ThoughtType = "THOUGHT_TYPE_INSTRUCTIONS"
	ThoughtType_ThoughtTypeSteps         ThoughtType = "THOUGHT_TYPE_STEPS"
)

// String representation for [fmt.Print].
func (f *ThoughtType) String() string {
	return string(*f)
}

type VerificationSection string

const (
	VerificationSection_VerificationSectionUnspecified           VerificationSection = "VERIFICATION_SECTION_UNSPECIFIED"
	VerificationSection_VerificationSectionSqlExamplesValidation VerificationSection = "VERIFICATION_SECTION_SQL_EXAMPLES_VALIDATION"
	VerificationSection_VerificationSectionVerificationQueries   VerificationSection = "VERIFICATION_SECTION_VERIFICATION_QUERIES"
	VerificationSection_VerificationSectionProposedImprovement   VerificationSection = "VERIFICATION_SECTION_PROPOSED_IMPROVEMENT"
	VerificationSection_VerificationSectionFinalDecision         VerificationSection = "VERIFICATION_SECTION_FINAL_DECISION"
)

// String representation for [fmt.Print].
func (f *VerificationSection) String() string {
	return string(*f)
}

type MessageError_Type string

const (
	MessageError_Type_TypeUnspecified                                  MessageError_Type = "TYPE_UNSPECIFIED"
	MessageError_Type_UnexpectedReplyProcessException                  MessageError_Type = "UNEXPECTED_REPLY_PROCESS_EXCEPTION"
	MessageError_Type_GenericChatCompletionException                   MessageError_Type = "GENERIC_CHAT_COMPLETION_EXCEPTION"
	MessageError_Type_ContextExceededException                         MessageError_Type = "CONTEXT_EXCEEDED_EXCEPTION"
	MessageError_Type_DeploymentNotFoundException                      MessageError_Type = "DEPLOYMENT_NOT_FOUND_EXCEPTION"
	MessageError_Type_FunctionsNotAvailableException                   MessageError_Type = "FUNCTIONS_NOT_AVAILABLE_EXCEPTION"
	MessageError_Type_InvalidCompletionRequestException                MessageError_Type = "INVALID_COMPLETION_REQUEST_EXCEPTION"
	MessageError_Type_ContentFilterException                           MessageError_Type = "CONTENT_FILTER_EXCEPTION"
	MessageError_Type_FunctionArgumentsInvalidJsonException            MessageError_Type = "FUNCTION_ARGUMENTS_INVALID_JSON_EXCEPTION"
	MessageError_Type_RetryableProcessingException                     MessageError_Type = "RETRYABLE_PROCESSING_EXCEPTION"
	MessageError_Type_InvalidFunctionCallException                     MessageError_Type = "INVALID_FUNCTION_CALL_EXCEPTION"
	MessageError_Type_LocalContextExceededException                    MessageError_Type = "LOCAL_CONTEXT_EXCEEDED_EXCEPTION"
	MessageError_Type_ChatCompletionNetworkException                   MessageError_Type = "CHAT_COMPLETION_NETWORK_EXCEPTION"
	MessageError_Type_InvalidChatCompletionJsonException               MessageError_Type = "INVALID_CHAT_COMPLETION_JSON_EXCEPTION"
	MessageError_Type_GenericChatCompletionServiceException            MessageError_Type = "GENERIC_CHAT_COMPLETION_SERVICE_EXCEPTION"
	MessageError_Type_WarehouseAccessMissingException                  MessageError_Type = "WAREHOUSE_ACCESS_MISSING_EXCEPTION"
	MessageError_Type_WarehouseNotFoundException                       MessageError_Type = "WAREHOUSE_NOT_FOUND_EXCEPTION"
	MessageError_Type_NoTablesToQueryException                         MessageError_Type = "NO_TABLES_TO_QUERY_EXCEPTION"
	MessageError_Type_SqlExecutionException                            MessageError_Type = "SQL_EXECUTION_EXCEPTION"
	MessageError_Type_ReplyProcessTimeoutException                     MessageError_Type = "REPLY_PROCESS_TIMEOUT_EXCEPTION"
	MessageError_Type_CouldNotGetUcSchemaException                     MessageError_Type = "COULD_NOT_GET_UC_SCHEMA_EXCEPTION"
	MessageError_Type_InvalidTableIdentifierException                  MessageError_Type = "INVALID_TABLE_IDENTIFIER_EXCEPTION"
	MessageError_Type_TooManyTablesException                           MessageError_Type = "TOO_MANY_TABLES_EXCEPTION"
	MessageError_Type_FunctionArgumentsInvalidException                MessageError_Type = "FUNCTION_ARGUMENTS_INVALID_EXCEPTION"
	MessageError_Type_GenericSqlExecApiCallException                   MessageError_Type = "GENERIC_SQL_EXEC_API_CALL_EXCEPTION"
	MessageError_Type_ChatCompletionClientException                    MessageError_Type = "CHAT_COMPLETION_CLIENT_EXCEPTION"
	MessageError_Type_ChatCompletionClientTimeoutException             MessageError_Type = "CHAT_COMPLETION_CLIENT_TIMEOUT_EXCEPTION"
	MessageError_Type_UnknownAiModel                                   MessageError_Type = "UNKNOWN_AI_MODEL"
	MessageError_Type_TablesMissingException                           MessageError_Type = "TABLES_MISSING_EXCEPTION"
	MessageError_Type_MessageDeletedWhileExecutingException            MessageError_Type = "MESSAGE_DELETED_WHILE_EXECUTING_EXCEPTION"
	MessageError_Type_MessageUpdatedWhileExecutingException            MessageError_Type = "MESSAGE_UPDATED_WHILE_EXECUTING_EXCEPTION"
	MessageError_Type_BlockMultipleExecutionsException                 MessageError_Type = "BLOCK_MULTIPLE_EXECUTIONS_EXCEPTION"
	MessageError_Type_InvalidCertifiedAnswerIdentifierException        MessageError_Type = "INVALID_CERTIFIED_ANSWER_IDENTIFIER_EXCEPTION"
	MessageError_Type_TooManyCertifiedAnswersException                 MessageError_Type = "TOO_MANY_CERTIFIED_ANSWERS_EXCEPTION"
	MessageError_Type_RateLimitExceededGenericException                MessageError_Type = "RATE_LIMIT_EXCEEDED_GENERIC_EXCEPTION"
	MessageError_Type_RateLimitExceededSpecifiedWaitException          MessageError_Type = "RATE_LIMIT_EXCEEDED_SPECIFIED_WAIT_EXCEPTION"
	MessageError_Type_FunctionCallMissingParameterException            MessageError_Type = "FUNCTION_CALL_MISSING_PARAMETER_EXCEPTION"
	MessageError_Type_InvalidCertifiedAnswerFunctionException          MessageError_Type = "INVALID_CERTIFIED_ANSWER_FUNCTION_EXCEPTION"
	MessageError_Type_IllegalParameterDefinitionException              MessageError_Type = "ILLEGAL_PARAMETER_DEFINITION_EXCEPTION"
	MessageError_Type_NoQueryToVisualizeException                      MessageError_Type = "NO_QUERY_TO_VISUALIZE_EXCEPTION"
	MessageError_Type_NoDeploymentsAvailableToWorkspace                MessageError_Type = "NO_DEPLOYMENTS_AVAILABLE_TO_WORKSPACE"
	MessageError_Type_StopProcessDueToAutoRegenerate                   MessageError_Type = "STOP_PROCESS_DUE_TO_AUTO_REGENERATE"
	MessageError_Type_FunctionArgumentsInvalidTypeException            MessageError_Type = "FUNCTION_ARGUMENTS_INVALID_TYPE_EXCEPTION"
	MessageError_Type_MessageCancelledWhileExecutingException          MessageError_Type = "MESSAGE_CANCELLED_WHILE_EXECUTING_EXCEPTION"
	MessageError_Type_CouldNotGetModelDeploymentsException             MessageError_Type = "COULD_NOT_GET_MODEL_DEPLOYMENTS_EXCEPTION"
	MessageError_Type_GeneratedSqlQueryTooLongException                MessageError_Type = "GENERATED_SQL_QUERY_TOO_LONG_EXCEPTION"
	MessageError_Type_MissingSqlQueryException                         MessageError_Type = "MISSING_SQL_QUERY_EXCEPTION"
	MessageError_Type_DescribeQueryUnexpectedFailure                   MessageError_Type = "DESCRIBE_QUERY_UNEXPECTED_FAILURE"
	MessageError_Type_DescribeQueryTimeout                             MessageError_Type = "DESCRIBE_QUERY_TIMEOUT"
	MessageError_Type_DescribeQueryInvalidSqlError                     MessageError_Type = "DESCRIBE_QUERY_INVALID_SQL_ERROR"
	MessageError_Type_InvalidSqlUnknownTableException                  MessageError_Type = "INVALID_SQL_UNKNOWN_TABLE_EXCEPTION"
	MessageError_Type_InvalidSqlMultipleStatementsException            MessageError_Type = "INVALID_SQL_MULTIPLE_STATEMENTS_EXCEPTION"
	MessageError_Type_InvalidSqlMultipleDatasetReferencesException     MessageError_Type = "INVALID_SQL_MULTIPLE_DATASET_REFERENCES_EXCEPTION"
	MessageError_Type_InvalidChatCompletionArgumentsJsonException      MessageError_Type = "INVALID_CHAT_COMPLETION_ARGUMENTS_JSON_EXCEPTION"
	MessageError_Type_MessageAttachmentTooLongError                    MessageError_Type = "MESSAGE_ATTACHMENT_TOO_LONG_ERROR"
	MessageError_Type_InternalCatalogPathOverlapException              MessageError_Type = "INTERNAL_CATALOG_PATH_OVERLAP_EXCEPTION"
	MessageError_Type_InternalCatalogMissingUcPathException            MessageError_Type = "INTERNAL_CATALOG_MISSING_UC_PATH_EXCEPTION"
	MessageError_Type_ExceededMaxTokenLengthException                  MessageError_Type = "EXCEEDED_MAX_TOKEN_LENGTH_EXCEPTION"
	MessageError_Type_InternalCatalogAssetCreationOngoingException     MessageError_Type = "INTERNAL_CATALOG_ASSET_CREATION_ONGOING_EXCEPTION"
	MessageError_Type_InternalCatalogAssetCreationFailedException      MessageError_Type = "INTERNAL_CATALOG_ASSET_CREATION_FAILED_EXCEPTION"
	MessageError_Type_InternalCatalogAssetCreationUnsupportedException MessageError_Type = "INTERNAL_CATALOG_ASSET_CREATION_UNSUPPORTED_EXCEPTION"
	MessageError_Type_UnsupportedConversationTypeException             MessageError_Type = "UNSUPPORTED_CONVERSATION_TYPE_EXCEPTION"
	MessageError_Type_CouldNotGetDashboardSchemaException              MessageError_Type = "COULD_NOT_GET_DASHBOARD_SCHEMA_EXCEPTION"
)

// String representation for [fmt.Print].
func (f *MessageError_Type) String() string {
	return string(*f)
}

type MessageStatus_MessageStatus string

const (
	MessageStatus_MessageStatus_FetchingMetadata   MessageStatus_MessageStatus = "FETCHING_METADATA"
	MessageStatus_MessageStatus_FilteringContext   MessageStatus_MessageStatus = "FILTERING_CONTEXT"
	MessageStatus_MessageStatus_AskingAi           MessageStatus_MessageStatus = "ASKING_AI"
	MessageStatus_MessageStatus_PendingWarehouse   MessageStatus_MessageStatus = "PENDING_WAREHOUSE"
	MessageStatus_MessageStatus_ExecutingQuery     MessageStatus_MessageStatus = "EXECUTING_QUERY"
	MessageStatus_MessageStatus_Failed             MessageStatus_MessageStatus = "FAILED"
	MessageStatus_MessageStatus_Completed          MessageStatus_MessageStatus = "COMPLETED"
	MessageStatus_MessageStatus_Submitted          MessageStatus_MessageStatus = "SUBMITTED"
	MessageStatus_MessageStatus_QueryResultExpired MessageStatus_MessageStatus = "QUERY_RESULT_EXPIRED"
	MessageStatus_MessageStatus_Cancelled          MessageStatus_MessageStatus = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *MessageStatus_MessageStatus) String() string {
	return string(*f)
}

type StatementStatus_State string

const (
	StatementStatus_State_StateUnspecified StatementStatus_State = "STATE_UNSPECIFIED"
	StatementStatus_State_Pending          StatementStatus_State = "PENDING"
	StatementStatus_State_Running          StatementStatus_State = "RUNNING"
	StatementStatus_State_Succeeded        StatementStatus_State = "SUCCEEDED"
	StatementStatus_State_Failed           StatementStatus_State = "FAILED"
	StatementStatus_State_Canceled         StatementStatus_State = "CANCELED"
	StatementStatus_State_Closed           StatementStatus_State = "CLOSED"
)

// String representation for [fmt.Print].
func (f *StatementStatus_State) String() string {
	return string(*f)
}

type ChunkInfo struct {
	ChunkIndex            *int    `json:"chunk_index"`
	RowOffset             *int64  `json:"row_offset"`
	RowCount              *int64  `json:"row_count"`
	ByteCount             *int64  `json:"byte_count"`
	NextChunkIndex        *int    `json:"next_chunk_index"`
	NextChunkInternalLink *string `json:"next_chunk_internal_link"`
}

type ColumnInfo struct {
	Name             *string         `json:"name"`
	TypeText         *string         `json:"type_text"`
	TypeName         *ColumnTypeName `json:"type_name"`
	Position         *int            `json:"position"`
	TypePrecision    *int            `json:"type_precision"`
	TypeScale        *int            `json:"type_scale"`
	TypeIntervalType *string         `json:"type_interval_type"`
	TypeJson         *string         `json:"type_json"`
	Comment          *string         `json:"comment"`
	Nullable         *bool           `json:"nullable"`
	PartitionIndex   *int            `json:"partition_index"`
	Mask             *ColumnMask     `json:"mask"`
}

type ColumnMask struct {
	FunctionName     *string                  `json:"function_name"`
	UsingColumnNames []string                 `json:"using_column_names"`
	UsingArguments   []PolicyFunctionArgument `json:"using_arguments"`
}

// Serialization format for DatabricksServiceException. Note the definition of
// this message should be in sync with
// DatabricksServiceExceptionWithDetailsProto defined in
// /api-base/proto/exception_with_details.proto except the later one has an
// extra error details field defined..
type DatabricksServiceExceptionProto struct {
	ErrorCode  *ErrorCode `json:"error_code"`
	Message    *string    `json:"message"`
	StackTrace *string    `json:"stack_trace"`
}

type ExternalLink struct {
	ExternalLink          *string           `json:"external_link"`
	Expiration            *string           `json:"expiration"`
	HttpHeaders           map[string]string `json:"http_headers"`
	ChunkIndex            *int              `json:"chunk_index"`
	RowOffset             *int64            `json:"row_offset"`
	RowCount              *int64            `json:"row_count"`
	ByteCount             *int64            `json:"byte_count"`
	NextChunkIndex        *int              `json:"next_chunk_index"`
	NextChunkInternalLink *string           `json:"next_chunk_internal_link"`
}

type ExternalLink_HttpHeadersEntry struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

// Genie AI Response.
type GenieAttachment struct {
	Text               *TextAttachment                    `json:"text"`
	Query              *GenieQueryAttachment              `json:"query"`
	SuggestedQuestions *GenieSuggestedQuestionsAttachment `json:"suggested_questions"`
	AttachmentId       *string                            `json:"attachment_id"`
}

type GenieConversation struct {
	Id                   *string `json:"id"`
	SpaceId              *string `json:"space_id"`
	UserId               *int64  `json:"user_id"`
	CreatedTimestamp     *int64  `json:"created_timestamp"`
	LastUpdatedTimestamp *int64  `json:"last_updated_timestamp"`
	Title                *string `json:"title"`
	ConversationId       *string `json:"conversation_id"`
}

type GenieConversationSummary struct {
	ConversationId   *string `json:"conversation_id"`
	Title            *string `json:"title"`
	CreatedTimestamp *int64  `json:"created_timestamp"`
}

type GenieCreateConversationMessageRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	Content        *string `json:"content"`
}

type GenieCreateEvalRunRequest struct {
	SpaceId              *string  `json:"space_id"`
	BenchmarkQuestionIds []string `json:"benchmark_question_ids"`
}

type GenieCreateMessageCommentRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	MessageId      *string `json:"message_id"`
	Content        *string `json:"content"`
}

type GenieCreateSpaceRequest struct {
	WarehouseId     *string `json:"warehouse_id"`
	ParentPath      *string `json:"parent_path"`
	SerializedSpace *string `json:"serialized_space"`
	Title           *string `json:"title"`
	Description     *string `json:"description"`
}

type GenieDeleteConversationMessageRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	MessageId      *string `json:"message_id"`
}

type GenieDeleteConversationRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
}

type GenieEvalResponse struct {
	Response           *string                `json:"response"`
	SqlExecutionResult *StatementResponse     `json:"sql_execution_result"`
	ResponseType       *GenieEvalResponseType `json:"response_type"`
}

// Shows summary information for an evaluation result. For detailed information
// including SQL execution results, actual/expected responses, and assessment
// scores, use GenieGetEvalResultDetails..
type GenieEvalResult struct {
	ResultId            *string               `json:"result_id"`
	SpaceId             *string               `json:"space_id"`
	BenchmarkQuestionId *string               `json:"benchmark_question_id"`
	Status              *EvaluationStatusType `json:"status"`
	Question            *string               `json:"question"`
	BenchmarkAnswer     *string               `json:"benchmark_answer"`
	CreatedByUser       *int64                `json:"created_by_user"`
}

// Shows detailed information for an evaluation result..
type GenieEvalResultDetails struct {
	ResultId            *string               `json:"result_id"`
	SpaceId             *string               `json:"space_id"`
	BenchmarkQuestionId *string               `json:"benchmark_question_id"`
	EvalRunStatus       *EvaluationStatusType `json:"eval_run_status"`
	Assessment          *GenieEvalAssessment  `json:"assessment"`
	ManualAssessment    *bool                 `json:"manual_assessment"`
	AssessmentReasons   []ScoreReason         `json:"assessment_reasons"`
	ActualResponse      []GenieEvalResponse   `json:"actual_response"`
	ExpectedResponse    []GenieEvalResponse   `json:"expected_response"`
}

type GenieEvalRunResponse struct {
	EvalRunId            *string               `json:"eval_run_id"`
	EvalRunStatus        *EvaluationStatusType `json:"eval_run_status"`
	RunByUser            *int64                `json:"run_by_user"`
	CreatedTimestamp     *int64                `json:"created_timestamp"`
	NumQuestions         *int64                `json:"num_questions"`
	NumCorrect           *int64                `json:"num_correct"`
	NumNeedsReview       *int64                `json:"num_needs_review"`
	NumDone              *int64                `json:"num_done"`
	LastUpdatedTimestamp *int64                `json:"last_updated_timestamp"`
}

type GenieExecuteMessageAttachmentQueryRequest struct {
	MessageId      *string `json:"message_id"`
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	AttachmentId   *string `json:"attachment_id"`
}

type GenieExecuteMessageQueryRequest struct {
	MessageId      *string `json:"message_id"`
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
}

// Feedback containing rating and optional comment.
type GenieFeedback struct {
	Rating  *GenieFeedbackRating `json:"rating"`
	Comment *string              `json:"comment"`
}

type GenieGenerateDownloadFullQueryResultRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	MessageId      *string `json:"message_id"`
	AttachmentId   *string `json:"attachment_id"`
}

type GenieGenerateDownloadFullQueryResultResponse struct {
	DownloadId          *string `json:"download_id"`
	DownloadIdSignature *string `json:"download_id_signature"`
}

type GenieGetConversationMessageRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	MessageId      *string `json:"message_id"`
}

type GenieGetDownloadFullQueryResultRequest struct {
	SpaceId             *string `json:"space_id"`
	ConversationId      *string `json:"conversation_id"`
	MessageId           *string `json:"message_id"`
	AttachmentId        *string `json:"attachment_id"`
	DownloadId          *string `json:"download_id"`
	DownloadIdSignature *string `json:"download_id_signature"`
}

type GenieGetDownloadFullQueryResultResponse struct {
	StatementResponse *StatementResponse `json:"statement_response"`
}

type GenieGetEvalResultDetailsRequest struct {
	SpaceId   *string `json:"space_id"`
	EvalRunId *string `json:"eval_run_id"`
	ResultId  *string `json:"result_id"`
}

type GenieGetEvalRunRequest struct {
	SpaceId   *string `json:"space_id"`
	EvalRunId *string `json:"eval_run_id"`
}

type GenieGetMessageAttachmentQueryResultRequest struct {
	MessageId      *string `json:"message_id"`
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	AttachmentId   *string `json:"attachment_id"`
}

type GenieGetMessageQueryResultRequest struct {
	MessageId      *string `json:"message_id"`
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
}

type GenieGetMessageQueryResultResponse struct {
	StatementResponse *StatementResponse `json:"statement_response"`
}

type GenieGetQueryResultByAttachmentRequest struct {
	MessageId      *string `json:"message_id"`
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	AttachmentId   *string `json:"attachment_id"`
}

type GenieGetSpaceRequest struct {
	SpaceId                *string `json:"space_id"`
	IncludeSerializedSpace *bool   `json:"include_serialized_space"`
}

type GenieListConversationCommentsRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	PageSize       *int    `json:"page_size"`
	PageToken      *string `json:"page_token"`
}

type GenieListConversationCommentsResponse struct {
	Comments      []GenieMessageComment `json:"comments"`
	NextPageToken *string               `json:"next_page_token"`
}

type GenieListConversationMessagesRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	PageSize       *int    `json:"page_size"`
	PageToken      *string `json:"page_token"`
}

type GenieListConversationMessagesResponse struct {
	Messages      []GenieMessage `json:"messages"`
	NextPageToken *string        `json:"next_page_token"`
}

type GenieListConversationsRequest struct {
	SpaceId    *string `json:"space_id"`
	PageSize   *int    `json:"page_size"`
	PageToken  *string `json:"page_token"`
	IncludeAll *bool   `json:"include_all"`
}

type GenieListConversationsResponse struct {
	Conversations []GenieConversationSummary `json:"conversations"`
	NextPageToken *string                    `json:"next_page_token"`
}

type GenieListEvalResultsRequest struct {
	SpaceId   *string `json:"space_id"`
	EvalRunId *string `json:"eval_run_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type GenieListEvalResultsResponse struct {
	EvalResults   []GenieEvalResult `json:"eval_results"`
	NextPageToken *string           `json:"next_page_token"`
}

type GenieListEvalRunsRequest struct {
	SpaceId   *string `json:"space_id"`
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type GenieListEvalRunsResponse struct {
	EvalRuns      []GenieEvalRunResponse `json:"eval_runs"`
	NextPageToken *string                `json:"next_page_token"`
}

type GenieListMessageCommentsRequest struct {
	SpaceId        *string `json:"space_id"`
	ConversationId *string `json:"conversation_id"`
	MessageId      *string `json:"message_id"`
	PageSize       *int    `json:"page_size"`
	PageToken      *string `json:"page_token"`
}

type GenieListMessageCommentsResponse struct {
	Comments      []GenieMessageComment `json:"comments"`
	NextPageToken *string               `json:"next_page_token"`
}

type GenieListSpacesRequest struct {
	PageSize  *int    `json:"page_size"`
	PageToken *string `json:"page_token"`
}

type GenieListSpacesResponse struct {
	Spaces        []GenieSpace `json:"spaces"`
	NextPageToken *string      `json:"next_page_token"`
}

type GenieMessage struct {
	Id                   *string                      `json:"id"`
	SpaceId              *string                      `json:"space_id"`
	ConversationId       *string                      `json:"conversation_id"`
	UserId               *int64                       `json:"user_id"`
	CreatedTimestamp     *int64                       `json:"created_timestamp"`
	LastUpdatedTimestamp *int64                       `json:"last_updated_timestamp"`
	Status               *MessageStatus_MessageStatus `json:"status"`
	Content              *string                      `json:"content"`
	Attachments          []GenieAttachment            `json:"attachments"`
	QueryResult          *Result                      `json:"query_result"`
	Error                *MessageError                `json:"error"`
	MessageId            *string                      `json:"message_id"`
	Feedback             *GenieFeedback               `json:"feedback"`
}

// A comment on a Genie conversation message..
type GenieMessageComment struct {
	SpaceId          *string `json:"space_id"`
	ConversationId   *string `json:"conversation_id"`
	MessageId        *string `json:"message_id"`
	MessageCommentId *string `json:"message_comment_id"`
	UserId           *int64  `json:"user_id"`
	Content          *string `json:"content"`
	CreatedTimestamp *int64  `json:"created_timestamp"`
}

type GenieQueryAttachment struct {
	Title                *string                    `json:"title"`
	Query                *string                    `json:"query"`
	Description          *string                    `json:"description"`
	LastUpdatedTimestamp *int64                     `json:"last_updated_timestamp"`
	Parameters           []QueryAttachmentParameter `json:"parameters"`
	Id                   *string                    `json:"id"`
	StatementId          *string                    `json:"statement_id"`
	QueryResultMetadata  *GenieResultMetadata       `json:"query_result_metadata"`
	Thoughts             []Thought                  `json:"thoughts"`
}

type GenieResultMetadata struct {
	RowCount    *int64 `json:"row_count"`
	IsTruncated *bool  `json:"is_truncated"`
}

type GenieSendMessageFeedbackRequest struct {
	SpaceId        *string              `json:"space_id"`
	ConversationId *string              `json:"conversation_id"`
	MessageId      *string              `json:"message_id"`
	Rating         *GenieFeedbackRating `json:"rating"`
	Comment        *string              `json:"comment"`
}

type GenieSpace struct {
	SpaceId         *string `json:"space_id"`
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	WarehouseId     *string `json:"warehouse_id"`
	ParentPath      *string `json:"parent_path"`
	SerializedSpace *string `json:"serialized_space"`
}

type GenieStartConversationMessageRequest struct {
	SpaceId *string `json:"space_id"`
	Content *string `json:"content"`
}

type GenieStartConversationResponse struct {
	MessageId      *string            `json:"message_id"`
	Message        *GenieMessage      `json:"message"`
	ConversationId *string            `json:"conversation_id"`
	Conversation   *GenieConversation `json:"conversation"`
}

// Follow-up questions suggested by Genie.
type GenieSuggestedQuestionsAttachment struct {
	Questions []string `json:"questions"`
}

type GenieTrashSpaceRequest struct {
	SpaceId *string `json:"space_id"`
}

type GenieUpdateSpaceRequest struct {
	SpaceId         *string `json:"space_id"`
	SerializedSpace *string `json:"serialized_space"`
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	WarehouseId     *string `json:"warehouse_id"`
}

// copied from proto3 / Google Well Known Types, source:
// https://github.com/protocolbuffers/protobuf/blob/450d24ca820750c5db5112a6f0b0c2efb9758021/src/google/protobuf/struct.proto
// `ListValue` is a wrapper around a repeated field of values.
//
// The JSON representation for `ListValue` is JSON array..
type ListValue struct {
	Values []Value `json:"values"`
}

// <Databricks> proto compiler is too old and does not support map. This is wire
// compatible with map<string,string>. See
// https://developers.google.com/protocol-buffers/docs/proto#backwards_compatibility..
type MapStringValueEntry struct {
	Key   *string `json:"key"`
	Value *Value  `json:"value"`
}

type MessageError struct {
	Error *string            `json:"error"`
	Type  *MessageError_Type `json:"type"`
}

type MessageStatus struct {
}

// A positional argument passed to a row filter or column mask function.
// Distinguishes between column references and literals..
type PolicyFunctionArgument struct {
	Column   *string `json:"column"`
	Constant *string `json:"constant"`
}

type QueryAttachmentParameter struct {
	Keyword *string `json:"keyword"`
	Value   *string `json:"value"`
	SqlType *string `json:"sql_type"`
}

type Result struct {
	StatementId          *string `json:"statement_id"`
	RowCount             *int64  `json:"row_count"`
	IsTruncated          *bool   `json:"is_truncated"`
	StatementIdSignature *string `json:"statement_id_signature"`
}

// Contains the result data of a single chunk when using `INLINE` disposition.
// When using `EXTERNAL_LINKS` disposition, the array `external_links` is used
// instead to provide URLs to the result data in cloud storage. Exactly one of
// these alternatives is used. (While the `external_links` array prepares the
// API to return multiple links in a single response. Currently only a single
// link is returned.).
type ResultData struct {
	ExternalLinks         []ExternalLink `json:"external_links"`
	DataArray             []ListValue    `json:"data_array"`
	ChunkIndex            *int           `json:"chunk_index"`
	RowOffset             *int64         `json:"row_offset"`
	RowCount              *int64         `json:"row_count"`
	ByteCount             *int64         `json:"byte_count"`
	NextChunkIndex        *int           `json:"next_chunk_index"`
	NextChunkInternalLink *string        `json:"next_chunk_internal_link"`
}

// The result manifest provides schema and metadata for the result set..
type ResultManifest struct {
	Format          *Format     `json:"format"`
	Schema          *Schema     `json:"schema"`
	TotalChunkCount *int        `json:"total_chunk_count"`
	Chunks          []ChunkInfo `json:"chunks"`
	TotalRowCount   *int64      `json:"total_row_count"`
	TotalByteCount  *int64      `json:"total_byte_count"`
	Truncated       *bool       `json:"truncated"`
}

type Schema struct {
	ColumnCount *int         `json:"column_count"`
	Columns     []ColumnInfo `json:"columns"`
}

type StatementResponse struct {
	StatementId *string          `json:"statement_id"`
	Status      *StatementStatus `json:"status"`
	Manifest    *ResultManifest  `json:"manifest"`
	Result      *ResultData      `json:"result"`
}

// The status response includes execution state and if relevant, error
// information..
type StatementStatus struct {
	State *StatementStatus_State           `json:"state"`
	Error *DatabricksServiceExceptionProto `json:"error"`
}

// copied from proto3 / Google Well Known Types, source:
// https://github.com/protocolbuffers/protobuf/blob/450d24ca820750c5db5112a6f0b0c2efb9758021/src/google/protobuf/struct.proto
// `Struct` represents a structured data value, consisting of fields which map
// to dynamically typed values. In some languages, `Struct` might be supported
// by a native representation. For example, in scripting languages like JS a
// struct is represented as an object. The details of that representation are
// described together with the proto support for the language.
//
// The JSON representation for `Struct` is JSON object..
type Struct struct {
	Fields []MapStringValueEntry `json:"fields"`
}

type TextAttachment struct {
	Content              *string                `json:"content"`
	Id                   *string                `json:"id"`
	Phase                *ResponsePhase         `json:"phase"`
	VerificationMetadata *VerificationMetadata  `json:"verification_metadata"`
	Purpose              *TextAttachmentPurpose `json:"purpose"`
}

// A single thought in the AI's reasoning process for a query..
type Thought struct {
	ThoughtType *ThoughtType `json:"thought_type"`
	Content     []string     `json:"content"`
}

// copied from proto3 / Google Well Known Types, source:
// https://github.com/protocolbuffers/protobuf/blob/450d24ca820750c5db5112a6f0b0c2efb9758021/src/google/protobuf/struct.proto
// `Value` represents a dynamically typed value which can be either null, a
// number, a string, a boolean, a recursive struct value, or a list of values. A
// producer of value is expected to set one of these variants. Absence of any
// variant indicates an error.
//
// The JSON representation for `Value` is JSON value..
type Value struct {
	NullValue   *NullValue `json:"null_value"`
	NumberValue *float64   `json:"number_value"`
	StringValue *string    `json:"string_value"`
	BoolValue   *bool      `json:"bool_value"`
	StructValue *Struct    `json:"struct_value"`
	ListValue   *ListValue `json:"list_value"`
}

// Metadata for verification phase attachments.
type VerificationMetadata struct {
	Section *VerificationSection `json:"section"`
	Index   *int                 `json:"index"`
}
