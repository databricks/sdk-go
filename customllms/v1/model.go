package v1

import "time"

type State string

const (
	State_StateUnspecified State = "STATE_UNSPECIFIED"
	State_Created          State = "CREATED"
	State_Running          State = "RUNNING"
	State_Completed        State = "COMPLETED"
	State_Failed           State = "FAILED"
	State_Pending          State = "PENDING"
	State_Cancelled        State = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *State) String() string {
	return string(*f)
}

type CancelCustomLlmOptimizationRunRequest struct {
	Id *string `json:"id"`
}

type CreateCustomLlmRequest struct {
	Name              *string   `json:"name"`
	Instructions      *string   `json:"instructions"`
	Datasets          []Dataset `json:"datasets"`
	Guidelines        []string  `json:"guidelines"`
	AgentArtifactPath *string   `json:"agent_artifact_path"`
}

type CustomLlm struct {
	Id                *string    `json:"id"`
	Name              *string    `json:"name"`
	EndpointName      *string    `json:"endpoint_name"`
	Instructions      *string    `json:"instructions"`
	Datasets          []Dataset  `json:"datasets"`
	Guidelines        []string   `json:"guidelines"`
	OptimizationState *State     `json:"optimization_state"`
	Creator           *string    `json:"creator"`
	CreationTime      *time.Time `json:"creation_time"`
	AgentArtifactPath *string    `json:"agent_artifact_path"`
}

type Dataset struct {
	Table *Table `json:"table"`
}

type DeleteCustomLlmRequest struct {
	Id *string `json:"id"`
}

type GetCustomLlmRequest struct {
	Id *string `json:"id"`
}

type StartCustomLlmOptimizationRunRequest struct {
	Id *string `json:"id"`
}

type Table struct {
	TablePath   *string `json:"table_path"`
	RequestCol  *string `json:"request_col"`
	ResponseCol *string `json:"response_col"`
}

type UpdateCustomLlmRequest struct {
	Id         *string    `json:"id"`
	CustomLlm  *CustomLlm `json:"custom_llm"`
	UpdateMask *string    `json:"update_mask"`
}
