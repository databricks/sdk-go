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
}

// A BindingPair maps a single field in either the request or response of the
// annotated method to a field of the method_to_poll operation..
type BindingPair struct {
}

// Long-Running Operation (LRO) configuration for API methods.
//
// This annotation is used to mark methods that return a long-running operation
// and specify the associated operation management methods and types.
//
// Long-running operations are asynchronous operations that may take an extended
// period of time to complete. Instead of blocking until completion, these
// methods immediately return an operation handle that can be used to track
// progress and retrieve results..
type LongRunningOperation struct {
	OperationInfo    *OperationInfo    `json:"operation_info"`
	OperationMethods *OperationMethods `json:"operation_methods"`
}

// Offset-based pagination.
//
// To retrieve results: 1. Set the offset field to 0 (or omit it) for the first
// page. 2. Set max_results to control how many items per page. 3. For
// subsequent pages, increment offset by the page size.
//
// Example: To get items 0-99, set offset=0 and max_results=100. To get items
// 100-199, set offset=100 and max_results=100..
type OffsetInfo struct {
}

// Specifies the response and metadata types for the operation..
type OperationInfo struct {
}

// References to methods that manage the long-running operation.
//
// These methods allow clients to check status, wait for completion, list
// operations, and cancel or delete operations..
type OperationMethods struct {
}

// Token-based pagination.
//
// To retrieve results: 1. Make the first request without a page token (or with
// an empty token). 2. Check the response for a next page token. 3. If present,
// include that token in the next request. 4. Repeat until the response contains
// no next page token.
//
// Note: Page tokens are opaque strings and should not be parsed or constructed
// by clients. Always use the exact token returned by the API..
type PageTokenInfo struct {
}

// Indicates that an API method returns paginated results.
//
// Methods annotated with this option may not return all results in a single
// call. Multiple requests are required to retrieve the complete dataset..
type Pagination struct {
	OffsetInfo *OffsetInfo    `json:"offset_info"`
	TokenInfo  *PageTokenInfo `json:"token_info"`
	Results    *string        `json:"results"`
}

// Groups fields used to find and interpret the operation's state..
type StateInfo struct {
}

// Annotation used to describe legacy Long Running Operations. WaitForState
// contains 4 key components: * method_to_poll: the method which returns the
// status of the operations. * binding: a mapping of fields used to construct
// the request to such method. * state_info: a message describing how to find
// and interpret the operation's state.
//
// The following example demonstrates how to annotate an RPC, `Create`, which
// initiates a long-running operation.
//
// 1. A client calls `Create(CreateRequest)`, which returns an `Object` response
// immediately. 2. The `wait_for_state` annotation tells the client generator
// that this operation isn't finished and must be polled. 3. To check the
// status, the client will poll the `Get` method (specified by `method_to_poll:
// "Get"`). 4. To build the `GetRequest` for polling, the client uses the
// `binding`. This example maps the `object_id` field from the initial `Create`
// *response* to the `id` field in the `Get` *request* (as defined by
// `poll_method_field: "id"` and `response_field: "object_id"`). 5. When the
// client receives the `Object` response from polling `Get`, it inspects the
// `state` field (defined by `state_info.state_path: "state"`). If the path were
// nested (e.g., `state_path: ["status", "phase"]`), it would inspect
// `object.status.phase`. 6. If the `state` field's value is `"RUNNING"`, the
// operation is complete and successful (per `state_info.success_states:
// "RUNNING"`). 7. If the `state` is `"STOPPED"` or `"FAILED"`, the operation is
// complete and has failed (per `state_info.failure_states: ["STOPPED",
// "FAILED"]`). 8. If the state is anything else, the client continues to poll
// the `Get` method until the timeout is reached.
//
// rpc Create(CreateRequest) returns Object { option
// (databricks.sdk.wait_for_state) = { method_to_poll: "Get" binding: {
// binding_pairs: { poll_method_field: "id" response_field: "object_id" } }
// state_info: { state_path: "state" success_states: "RUNNING" failure_states:
// ["STOPPED", "FAILED"] } } } rpc Get(GetRequest) returns Object { }
//
// message Object { optional string object_id = 1; optional State state = 2; }
// message GetRequest { optional string id = 1; }.
type WaitForState struct {
	MethodToPoll *string    `json:"method_to_poll"`
	Binding      *Binding   `json:"binding"`
	StateInfo    *StateInfo `json:"state_info"`
}
