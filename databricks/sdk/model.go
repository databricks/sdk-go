package sdk

type LaunchStage string

const (
	LaunchStageLaunchStageUnspecified LaunchStage = "LAUNCH_STAGE_UNSPECIFIED"
	LaunchStageDevelopment            LaunchStage = "DEVELOPMENT"
	LaunchStagePrivatePreview         LaunchStage = "PRIVATE_PREVIEW"
	LaunchStagePublicBeta             LaunchStage = "PUBLIC_BETA"
	LaunchStagePublicPreview          LaunchStage = "PUBLIC_PREVIEW"
	LaunchStageGa                     LaunchStage = "GA"
)

// String representation for [fmt.Print].
func (f *LaunchStage) String() string {
	return string(*f)
}

type Binding struct {
	Pairs []BindingPair `json:"pairs"`
}

type BindingPair struct {
	PollMethodField *string `json:"poll_method_field"`
	RequestField    *string `json:"request_field"`
	ResponseField   *string `json:"response_field"`
}

type LongRunningOperation struct {
	OperationInfo    *OperationInfo    `json:"operation_info"`
	OperationMethods *OperationMethods `json:"operation_methods"`
}

type OffsetInfo struct {
	Offset            *string `json:"offset"`
	MaxResults        *string `json:"max_results"`
	DefaultMaxResults *int    `json:"default_max_results"`
}

type OperationInfo struct {
	ResponseType *string `json:"response_type"`
	MetadataType *string `json:"metadata_type"`
}

type OperationMethods struct {
	Get    *string `json:"get"`
	List   *string `json:"list"`
	Wait   *string `json:"wait"`
	Delete *string `json:"delete"`
	Cancel *string `json:"cancel"`
}

type PageTokenInfo struct {
	Request           *string `json:"request"`
	Response          *string `json:"response"`
	MaxResults        *string `json:"max_results"`
	DefaultMaxResults *int    `json:"default_max_results"`
}

type Pagination struct {
	OffsetInfo *OffsetInfo    `json:"offset_info"`
	TokenInfo  *PageTokenInfo `json:"token_info"`
	Results    *string        `json:"results"`
}

type StateInfo struct {
	StatePath     []string `json:"state_path"`
	SuccessStates []string `json:"success_states"`
	FailureStates []string `json:"failure_states"`
	MessagePath   []string `json:"message_path"`
}

type WaitForState struct {
	MethodToPoll *string    `json:"method_to_poll"`
	Binding      *Binding   `json:"binding"`
	StateInfo    *StateInfo `json:"state_info"`
}
