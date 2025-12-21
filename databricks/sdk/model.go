
package sdk

import (
	
	"fmt"
	"io"
)



type LaunchStage string

const LaunchStageLaunchStageUnspecified LaunchStage = `LAUNCH_STAGE_UNSPECIFIED`
const LaunchStageDevelopment LaunchStage = `DEVELOPMENT`
const LaunchStagePrivatePreview LaunchStage = `PRIVATE_PREVIEW`
const LaunchStagePublicBeta LaunchStage = `PUBLIC_BETA`
const LaunchStagePublicPreview LaunchStage = `PUBLIC_PREVIEW`
const LaunchStageGa LaunchStage = `GA`

// String representation for [fmt.Print]
func (f *LaunchStage) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *LaunchStage) Set(v string) error {
	switch v {
	case `LAUNCH_STAGE_UNSPECIFIED`, `DEVELOPMENT`, `PRIVATE_PREVIEW`, `PUBLIC_BETA`, `PUBLIC_PREVIEW`, `GA`:
		*f = LaunchStage(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "LAUNCH_STAGE_UNSPECIFIED", "DEVELOPMENT", "PRIVATE_PREVIEW", "PUBLIC_BETA", "PUBLIC_PREVIEW", "GA"`, v)
	}
}

// Values returns all possible values for LaunchStage.
//
// There is no guarantee on the order of the values in the slice.
func (f *LaunchStage) Values() []LaunchStage {
	return []LaunchStage{
		LaunchStageLaunchStageUnspecified,
		LaunchStageDevelopment,
		LaunchStagePrivatePreview,
		LaunchStagePublicBeta,
		LaunchStagePublicPreview,
		LaunchStageGa,
	}
}

// Type always returns LaunchStage to satisfy [pflag.Value] interface
func (f *LaunchStage) Type() string {
	return "LaunchStage"
}





type Binding struct {
	
	Pairs []BindingPair `json:"pairs"`
	
}



type BindingPair struct {
	
	PollMethodField string `json:"pollMethodField"`
	
	RequestField string `json:"requestField"`
	
	ResponseField string `json:"responseField"`
	
}



type LongRunningOperation struct {
	
	OperationInfo OperationInfo `json:"operationInfo"`
	
	OperationMethods OperationMethods `json:"operationMethods"`
	
}



type OffsetInfo struct {
	
	Offset string `json:"offset"`
	
	MaxResults string `json:"maxResults"`
	
	DefaultMaxResults int `json:"defaultMaxResults"`
	
}



type OperationInfo struct {
	
	ResponseType string `json:"responseType"`
	
	MetadataType string `json:"metadataType"`
	
}



type OperationMethods struct {
	
	Get string `json:"get"`
	
	List string `json:"list"`
	
	Wait string `json:"wait"`
	
	Delete string `json:"delete"`
	
	Cancel string `json:"cancel"`
	
}



type PageTokenInfo struct {
	
	Request string `json:"request"`
	
	Response string `json:"response"`
	
	MaxResults string `json:"maxResults"`
	
	DefaultMaxResults int `json:"defaultMaxResults"`
	
}



type Pagination struct {
	
	OffsetInfo OffsetInfo `json:"offsetInfo"`
	
	TokenInfo PageTokenInfo `json:"tokenInfo"`
	
	Results string `json:"results"`
	
}



type StateInfo struct {
	
	StatePath []string `json:"statePath"`
	
	SuccessStates []string `json:"successStates"`
	
	FailureStates []string `json:"failureStates"`
	
	MessagePath []string `json:"messagePath"`
	
}



type WaitForState struct {
	
	MethodToPoll string `json:"methodToPoll"`
	
	Binding Binding `json:"binding"`
	
	StateInfo StateInfo `json:"stateInfo"`
	
}


