package v1

type TaskState string

const (
	TaskStateTaskStateUnspecified TaskState = "TASK_STATE_UNSPECIFIED"
	TaskStatePending              TaskState = "PENDING"
	TaskStateRunning              TaskState = "RUNNING"
	TaskStateCompleted            TaskState = "COMPLETED"
	TaskStateFailed               TaskState = "FAILED"
	TaskStateCancelled            TaskState = "CANCELLED"
	TaskStateInternalError        TaskState = "INTERNAL_ERROR"
)

// String representation for [fmt.Print].
func (f *TaskState) String() string {
	return string(*f)
}

// Request to cancel a task..
type CancelTaskRequest struct {
	TaskId *string `json:"task_id"`
	Reason *string `json:"reason"`
}

// Request to create a new task..
type CreateTaskRequest struct {
	Name   *string     `json:"name"`
	Config *TaskConfig `json:"config"`
}

// Request to get task status..
type GetTaskRequest struct {
	TaskId *string `json:"task_id"`
}

type ParametersEntry struct {
}

// Represents a task entity..
type Task struct {
	TaskId    *string     `json:"task_id"`
	Name      *string     `json:"name"`
	Status    *TaskStatus `json:"status"`
	Config    *TaskConfig `json:"config"`
	CreatedAt *int64      `json:"created_at"`
	UpdatedAt *int64      `json:"updated_at"`
}

// Configuration for a task..
type TaskConfig struct {
	TimeoutSeconds *int              `json:"timeout_seconds"`
	MaxRetries     *int              `json:"max_retries"`
	Parameters     map[string]string `json:"parameters"`
}

// Status information for a task..
type TaskStatus struct {
	State    *TaskState `json:"state"`
	Message  *string    `json:"message"`
	Progress *int       `json:"progress"`
}
