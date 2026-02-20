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
	"testing"

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
	createStatus   int
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

	createStatus := cfg.createStatus
	if createStatus == 0 {
		createStatus = http.StatusOK
	}
	createResponse := cfg.createResponse
	if createResponse == nil {
		createResponse = &Task{
			TaskId: ptr(defaultCreateTaskID),
			Status: &TaskStatus{State: ptr(TaskStatePending)},
		}
	}

	createPayload := mustMarshalJSON(t, createResponse)
	createErrorPayload := mustMarshalJSON(t, map[string]string{"message": "mock create error"})

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/2.0/tasks":
			if createStatus < 200 || createStatus >= 300 {
				w.WriteHeader(createStatus)
				if _, err := w.Write(createErrorPayload); err != nil {
					return
				}
				return
			}
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
		serverCfg       mockTaskServerConfig
		wantRawResponse *Task
		wantErr         string
		wantErrStatus   int
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
			wantErr: "response field TaskId required for polling is nil",
		},
		{
			name: "create task http error",
			serverCfg: mockTaskServerConfig{
				createStatus: http.StatusInternalServerError,
			},
			wantErr:       "databricks-api error: mock create error",
			wantErrStatus: http.StatusInternalServerError,
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
				if tc.wantErrStatus != 0 {
					var apiErr *apierr.APIError
					if !errors.As(err, &apiErr) {
						t.Fatalf("expected *apierr.APIError, got %T (%v)", err, err)
					}
					if got := apiErr.HTTPStatusCode(); got != tc.wantErrStatus {
						t.Errorf("expected HTTP status %d, got %d", tc.wantErrStatus, got)
					}
				}
				return
			}
			if tc.wantErr != "" {
				t.Fatalf("expected error %q, got nil", tc.wantErr)
			}
			if waiter == nil {
				t.Fatal("expected waiter, got nil")
			}
			if waiter.service != client {
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
