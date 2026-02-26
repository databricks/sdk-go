package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/databricks/sdk-go/databricks/options"
	"github.com/google/go-cmp/cmp"
)

const (
	testTaskName        = "test-task"
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
