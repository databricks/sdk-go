package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func mustMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal mock response: %v", err)
	}
	return b
}

func mockTaskServer(t *testing.T, createResponse *Task) *httptest.Server {
	t.Helper()

	createPayload := mustMarshalJSON(t, createResponse)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/2.0/tasks":
			if _, err := w.Write(createPayload); err != nil {
				return
			}

		default:
			http.Error(w, fmt.Sprintf(`{"error":"not found: %s %s"}`, r.Method, r.URL.Path), http.StatusNotFound)
		}
	}))
}

func mockWaitServer(t *testing.T, pollResponses []Task) *httptest.Server {
	t.Helper()

	var pollCount int

	pollPayloads := make([][]byte, len(pollResponses))
	for i, response := range pollResponses {
		pollPayloads[i] = mustMarshalJSON(t, response)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/2.0/tasks/"):
			if len(pollPayloads) == 0 {
				http.Error(w, `{"error":"no poll responses configured"}`, http.StatusInternalServerError)
				return
			}
			idx := pollCount
			pollCount++
			if idx >= len(pollPayloads) {
				http.Error(w, fmt.Sprintf(`{"error":"unexpected poll request %d"}`, idx), http.StatusInternalServerError)
				return
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
		createResponse  *Task
		wantRawResponse *Task
		wantErr         string
	}{
		{
			name: "success",
			createResponse: &Task{
				TaskId: ptr(successCreateTaskID),
				Status: &TaskStatus{State: ptr(TaskStatePending)},
			},
			wantRawResponse: &Task{
				TaskId: ptr(successCreateTaskID),
				Status: &TaskStatus{State: ptr(TaskStatePending)},
			},
		},
		{
			name:           "nil TaskId in response",
			createResponse: &Task{Status: &TaskStatus{State: ptr(TaskStatePending)}},
			wantErr:        "nil TaskId in response",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockTaskServer(t, tc.createResponse)
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

func TestDone(t *testing.T) {
	testCases := []struct {
		name         string
		pollResponse Task
		wantDone     bool
		wantErr      bool
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
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockWaitServer(t, []Task{tc.pollResponse})
			defer server.Close()

			waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

			done, err := waiter.Done(context.Background())

			if err != nil {
				if !tc.wantErr {
					t.Fatalf("Done: %v", err)
				}
				if done != tc.wantDone {
					t.Errorf("expected done=%v, got %v", tc.wantDone, done)
				}
				return
			}
			if tc.wantErr {
				t.Fatal("expected error, got nil")
			}
			if done != tc.wantDone {
				t.Errorf("expected done=%v, got %v", tc.wantDone, done)
			}
		})
	}
}
