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

type CreateTaskRequest struct {
	Name   *string     `json:"name"`
	Config *TaskConfig `json:"config"`
}

type GetTaskRequest struct {
	TaskId *string `json:"task_id"`
}

type ParametersEntry struct {
}

type Task struct {
	TaskId    *string     `json:"task_id"`
	Name      *string     `json:"name"`
	Status    *TaskStatus `json:"status"`
	Config    *TaskConfig `json:"config"`
	CreatedAt *int64      `json:"created_at"`
	UpdatedAt *int64      `json:"updated_at"`
}

type TaskConfig struct {
	TimeoutSeconds *int              `json:"timeout_seconds"`
	MaxRetries     *int              `json:"max_retries"`
	Parameters     map[string]string `json:"parameters"`
}

type TaskStatus struct {
	State    *TaskState `json:"state"`
	Message  *string    `json:"message"`
	Progress *int       `json:"progress"`
}
