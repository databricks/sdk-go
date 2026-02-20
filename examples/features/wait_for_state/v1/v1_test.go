package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/databricks/sdk-go/databricks/api"
	"github.com/databricks/sdk-go/databricks/apierr"
	"github.com/databricks/sdk-go/databricks/options"
	"github.com/google/go-cmp/cmp"
)

const (
	testTaskName        = "test-task"
	defaultCreateTaskID = "task-123"
	successCreateTaskID = "task-456"
)

func ptr[T any](v T) *T { return &v }

func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	client, err := NewClient(context.Background(),
		options.WithHost(server.URL),
		options.WithHTTPClient(server.Client()),
		options.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

type mockTaskServerConfig struct {
	createResponse *Task
	pollResponses  []Task
	pollStatus     int
}

func mustMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal mock response: %v", err)
	}
	return b
}

func mockTaskServer(t *testing.T, cfg mockTaskServerConfig) *httptest.Server {
	t.Helper()

	var pollCount int

	pollStatus := cfg.pollStatus
	if pollStatus == 0 {
		pollStatus = http.StatusOK
	}

	createPayload := mustMarshalJSON(t, cfg.createResponse)
	pollErrorPayload := mustMarshalJSON(t, map[string]string{"message": "mock get error"})
	pollPayloads := make([][]byte, len(cfg.pollResponses))
	for i, response := range cfg.pollResponses {
		pollPayloads[i] = mustMarshalJSON(t, response)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/2.0/tasks":
			if _, err := w.Write(createPayload); err != nil {
				return
			}

		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/2.0/tasks/"):
			if pollStatus < 200 || pollStatus >= 300 {
				w.WriteHeader(pollStatus)
				if _, err := w.Write(pollErrorPayload); err != nil {
					return
				}
				return
			}
			if len(pollPayloads) == 0 {
				http.Error(w, `{"error":"no poll responses configured"}`, http.StatusInternalServerError)
				return
			}
			idx := pollCount
			pollCount++
			// Clamp to the last entry so tests that poll more times than
			// there are responses (e.g. timeout tests) don't panic.
			if idx >= len(pollPayloads) {
				idx = len(pollPayloads) - 1
			}
			if _, err := w.Write(pollPayloads[idx]); err != nil {
				return
			}

		default:
			http.Error(w, fmt.Sprintf(`{"error":"not found: %s %s"}`, r.Method, r.URL.Path), http.StatusNotFound)
		}
	}))
}

func newWaiter(client *Client, taskId string) *CreateTaskWaiter {
	return &CreateTaskWaiter{
		rawResponse: &Task{TaskId: ptr(taskId)},
		client:      client,
		taskId:      ptr(taskId),
	}
}

func TestCreateTaskWaiter(t *testing.T) {
	testCases := []struct {
		name            string
		serverCfg       mockTaskServerConfig
		wantRawResponse *Task
		wantErr         string
	}{
		{
			name: "success",
			serverCfg: mockTaskServerConfig{
				createResponse: &Task{
					TaskId: ptr(successCreateTaskID),
					Status: &TaskStatus{State: ptr(TaskStatePending)},
				},
			},
			wantRawResponse: &Task{
				TaskId: ptr(successCreateTaskID),
				Status: &TaskStatus{State: ptr(TaskStatePending)},
			},
		},
		{
			name: "nil TaskId in response",
			serverCfg: mockTaskServerConfig{
				createResponse: &Task{Status: &TaskStatus{State: ptr(TaskStatePending)}},
			},
			wantErr: "nil TaskId in response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockTaskServer(t, tc.serverCfg)
			defer server.Close()

			client := newTestClient(t, server)
			waiter, err := client.CreateTaskWaiter(context.Background(), &CreateTaskRequest{
				Name: ptr(testTaskName),
			})

			if err != nil {
				if tc.wantErr == "" {
					t.Fatalf("CreateTaskWaiter: %v", err)
				}
				if got := err.Error(); got != tc.wantErr {
					t.Errorf("expected error %q, got %q", tc.wantErr, got)
				}
				return
			}
			if tc.wantErr != "" {
				t.Fatalf("expected error %q, got nil", tc.wantErr)
			}
			if waiter == nil {
				t.Fatal("expected waiter, got nil")
			}
			if waiter.client != client {
				t.Error("waiter should keep reference to the originating client")
			}
			if diff := cmp.Diff(tc.wantRawResponse, waiter.rawResponse); diff != "" {
				t.Errorf("unexpected rawResponse (-want +got):\n%s", diff)
			}
			if waiter.taskId != waiter.rawResponse.TaskId {
				t.Error("waiter.taskId should reference rawResponse.TaskId")
			}
		})
	}
}

func TestWait(t *testing.T) {
	testCases := []struct {
		name          string
		pollResponses []Task
		wantResult    *Task
		wantErr       string
	}{
		{
			name: "completed after polling",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
			},
			wantResult: &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
		},
		{
			name: "completed immediately",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
			},
			wantResult: &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
		},
		{
			name: "cancelled",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCancelled)}},
			},
			wantResult: &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCancelled)}},
		},
		{
			name: "failed with message",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{
					State:   ptr(TaskStateFailed),
					Message: ptr("something went wrong"),
				}},
			},
			wantErr: "terminal state FAILED: something went wrong",
		},
		{
			name: "failed without message",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{
					State: ptr(TaskStateFailed),
				}},
			},
			wantErr: "terminal state FAILED: (no message)",
		},
		{
			name: "internal error",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{
					State: ptr(TaskStateInternalError),
				}},
			},
			wantErr: "terminal state INTERNAL_ERROR: (no message)",
		},
		{
			name: "missing status",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID)},
			},
			wantErr: errMissingStatus.Error(),
		},
		{
			name: "nil state in status",
			pollResponses: []Task{
				{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{}},
			},
			wantErr: errMissingStatus.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockTaskServer(t, mockTaskServerConfig{pollResponses: tc.pollResponses})
			defer server.Close()

			client := newTestClient(t, server)
			waiter, err := client.CreateTaskWaiter(context.Background(), &CreateTaskRequest{
				Name: ptr(testTaskName),
			})
			if err != nil {
				t.Fatalf("CreateTaskWaiter: %v", err)
			}

			result, err := waiter.Wait(context.Background())

			if err != nil {
				if tc.wantErr == "" {
					t.Fatalf("Wait: %v", err)
				}
				if got := err.Error(); got != tc.wantErr {
					t.Errorf("expected error %q, got %q", tc.wantErr, got)
				}
				return
			}
			if tc.wantErr != "" {
				t.Fatalf("expected error %q, got nil", tc.wantErr)
			}
			if diff := cmp.Diff(tc.wantResult, result); diff != "" {
				t.Errorf("unexpected result (-want +got):\n%s", diff)
			}
		})
	}
}

// Option override tests verify the two-level option semantics in Wait:
// user opts - retrier and limiter - apply to inner GetTask calls,
// while the outer polling loop uses user timeout and its own retrier and limiter.
func TestWait_UserTimeout(t *testing.T) {
	server := mockTaskServer(t, mockTaskServerConfig{pollResponses: []Task{
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
	}})
	defer server.Close()

	waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

	start := time.Now()
	_, err := waiter.Wait(context.Background(), api.WithTimeout(50*time.Millisecond))
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got: %v", err)
	}
	if elapsed > 5*time.Second {
		t.Errorf("timeout should have fired quickly, took %v", elapsed)
	}
}

func TestWait_PollingRetrierOverridesUserRetrier(t *testing.T) {
	// If the user's WithDisableRetry applied to the polling loop, the first
	// non-terminal poll would be fatal. Polling must continue regardless.
	server := mockTaskServer(t, mockTaskServerConfig{pollResponses: []Task{
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
	}})
	defer server.Close()

	waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

	result, err := waiter.Wait(context.Background(), api.WithDisableRetry())
	if err != nil {
		t.Fatalf("Wait should succeed despite WithDisableRetry: %v", err)
	}
	wantResult := &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}}
	if diff := cmp.Diff(wantResult, result); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
}

type countingLimiter struct {
	calls *int
}

func (l *countingLimiter) Wait(context.Context) error {
	(*l.calls)++
	return nil
}

func TestWait_UserLimiterAppliesToGetTaskCalls(t *testing.T) {
	// The user's limiter is forced to nil for the polling loop, but opts
	// still flow to inner GetTask calls. Verify it is actually invoked.
	server := mockTaskServer(t, mockTaskServerConfig{pollResponses: []Task{
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
	}})
	defer server.Close()

	waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

	var calls int
	limiter := &countingLimiter{calls: &calls}

	result, err := waiter.Wait(context.Background(), api.WithLimiter(limiter))
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}
	wantResult := &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}}
	if diff := cmp.Diff(wantResult, result); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
	if got := calls; got != 2 {
		t.Errorf("expected limiter to be called exactly 2 times, got %d", got)
	}
}

func TestDone(t *testing.T) {
	testCases := []struct {
		name          string
		pollResponse  Task
		pollStatus    int
		wantDone      bool
		wantErr       string
		wantErrStatus int
	}{
		{
			name:         "completed",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
			wantDone:     true,
		},
		{
			name:         "cancelled",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCancelled)}},
			wantDone:     true,
		},
		{
			name:         "failed",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateFailed)}},
			wantDone:     true,
		},
		{
			name:         "internal error",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateInternalError)}},
			wantDone:     true,
		},
		{
			name:         "running - not done",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
			wantDone:     false,
		},
		{
			name:         "pending - not done",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStatePending)}},
			wantDone:     false,
		},
		{
			name:         "missing status",
			pollResponse: Task{TaskId: ptr(defaultCreateTaskID)},
			wantDone:     false,
			wantErr:      errMissingStatus.Error(),
		},
		{
			name:          "get task http error",
			pollStatus:    http.StatusInternalServerError,
			wantDone:      false,
			wantErr:       "databricks-api error: mock get error",
			wantErrStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := mockTaskServerConfig{pollStatus: tc.pollStatus}
			if tc.pollStatus == 0 {
				cfg.pollResponses = []Task{tc.pollResponse}
			}
			server := mockTaskServer(t, cfg)
			defer server.Close()

			waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

			done, err := waiter.Done(context.Background())

			if err != nil {
				if tc.wantErr == "" {
					t.Fatalf("Done: %v", err)
				}
				if got := err.Error(); got != tc.wantErr {
					t.Errorf("expected error %q, got %q", tc.wantErr, got)
				}
				if tc.wantErrStatus != 0 {
					var apiErr *apierr.APIError
					if !errors.As(err, &apiErr) {
						t.Fatalf("expected *apierr.APIError, got %T (%v)", err, err)
					}
					if got := apiErr.HTTPStatusCode(); got != tc.wantErrStatus {
						t.Errorf("expected HTTP status %d, got %d", tc.wantErrStatus, got)
					}
				}
				if done != tc.wantDone {
					t.Errorf("expected done=%v, got %v", tc.wantDone, done)
				}
				return
			}
			if tc.wantErr != "" {
				t.Fatalf("expected error %q, got nil", tc.wantErr)
			}
			if done != tc.wantDone {
				t.Errorf("expected done=%v, got %v", tc.wantDone, done)
			}
		})
	}
}
