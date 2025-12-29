
package sdk

import (
	
	"fmt"
	"io"
)



type LaunchStage string

const (
	LaunchStageLaunchStageUnspecified LaunchStage = "LAUNCH_STAGE_UNSPECIFIED"
	LaunchStageDevelopment LaunchStage = "DEVELOPMENT"
	LaunchStagePrivatePreview LaunchStage = "PRIVATE_PREVIEW"
	LaunchStagePublicBeta LaunchStage = "PUBLIC_BETA"
	LaunchStagePublicPreview LaunchStage = "PUBLIC_PREVIEW"
	LaunchStageGa LaunchStage = "GA"
)

// String representation for [fmt.Print].
func (f *LaunchStage) String() string {
	return string(*f)
}





type Binding struct {
	
	Pairs []BindingPair `json:"pairs"`
	
}



type BindingPair struct {
	
	PollMethodField *string `json:"pollMethodField"`
	
	RequestField *string `json:"requestField"`
	
	ResponseField *string `json:"responseField"`
	
}



type LongRunningOperation struct {
	
	OperationInfo *OperationInfo `json:"operationInfo"`
	
	OperationMethods *OperationMethods `json:"operationMethods"`
	
}



type OffsetInfo struct {
	
	Offset *string `json:"offset"`
	
	MaxResults *string `json:"maxResults"`
	
	DefaultMaxResults *int `json:"defaultMaxResults"`
	
}



type OperationInfo struct {
	
	ResponseType *string `json:"responseType"`
	
	MetadataType *string `json:"metadataType"`
	
}



type OperationMethods struct {
	
	Get *string `json:"get"`
	
	List *string `json:"list"`
	
	Wait *string `json:"wait"`
	
	Delete *string `json:"delete"`
	
	Cancel *string `json:"cancel"`
	
}



type PageTokenInfo struct {
	
	Request *string `json:"request"`
	
	Response *string `json:"response"`
	
	MaxResults *string `json:"maxResults"`
	
	DefaultMaxResults *int `json:"defaultMaxResults"`
	
}



type Pagination struct {
	
	OffsetInfo *OffsetInfo `json:"offsetInfo"`
	
	TokenInfo *PageTokenInfo `json:"tokenInfo"`
	
	Results *string `json:"results"`
	
}



type StateInfo struct {
	
	StatePath []string `json:"statePath"`
	
	SuccessStates []string `json:"successStates"`
	
	FailureStates []string `json:"failureStates"`
	
	MessagePath []string `json:"messagePath"`
	
}



type WaitForState struct {
	
	MethodToPoll *string `json:"methodToPoll"`
	
	Binding *Binding `json:"binding"`
	
	StateInfo *StateInfo `json:"stateInfo"`
	
}


