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
	"github.com/databricks/sdk-go/databricks/apierr/codes"
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
			server := mockWaitServer(t, tc.pollResponses)
			defer server.Close()

			client := newTestClient(t, server)
			waiter := newWaiter(client, defaultCreateTaskID)

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
// the timeout is stripped from inner GetTask calls (the poll loop owns
// the deadline), the retrier is overridden to drive polling, and the
// rate limiter is disabled at the poll loop level.
func TestWait_UserTimeout(t *testing.T) {
	server := mockWaitServer(t, []Task{
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
	})
	defer server.Close()

	waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

	// The timeout is too tight for the whole polling sequence (backoff sleeps
	// exceed 50ms) but each individual GetTask call succeeds instantly. This
	// verifies that WithTimeout bounds the entire Wait, not individual calls.
	//
	// A goroutine guards against Wait hanging if WithTimeout is broken.
	// A context timeout cannot be the safety net because DeadlineExceeded
	// from WithTimeout vs the parent context is indistinguishable.
	errCh := make(chan error, 1)
	go func() {
		_, err := waiter.Wait(context.Background(), api.WithTimeout(50*time.Millisecond))
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected DeadlineExceeded, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("api.WithTimeout did not fire within 5s")
	}
}

func TestWait_PollingRetrierOverridesUserRetrier(t *testing.T) {
	// If WithDisableRetry leaks into the outer polling loop, the first "still
	// running" response will be treated as a fatal error and the test will fail.
	server := mockWaitServer(t, []Task{
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}},
		{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}},
	})
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

func TestWait_UserLimiterNotUsedByOuterExecute(t *testing.T) {
	// A 429 on the first request forces an inner retry, producing 3 HTTP
	// calls across only 2 outer poll iterations. If the limiter were on the
	// outer Execute it would only fire twice.
	var requestCount int
	responses := []struct {
		code int
		body []byte
	}{
		{code: http.StatusTooManyRequests},
		{body: mustMarshalJSON(t, Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateRunning)}})},
		{body: mustMarshalJSON(t, Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}})},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idx := requestCount
		requestCount++
		if idx >= len(responses) {
			http.Error(w, fmt.Sprintf(`{"error":"unexpected request %d"}`, idx), http.StatusInternalServerError)
			return
		}
		resp := responses[idx]
		if resp.code != 0 {
			http.Error(w, `{"error":"rate limit"}`, resp.code)
			return
		}
		if _, err := w.Write(resp.body); err != nil {
			return
		}
	}))
	defer server.Close()

	waiter := newWaiter(newTestClient(t, server), defaultCreateTaskID)

	var limiterCalls int
	limiter := &countingLimiter{calls: &limiterCalls}

	result, err := waiter.Wait(context.Background(),
		api.WithLimiter(limiter),
		api.WithRetrier(func() api.Retrier {
			return api.RetryOnCodes(api.BackoffPolicy{
				Initial: time.Millisecond,
				Maximum: time.Millisecond,
			}, codes.ResourceExhausted)
		}),
	)
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}
	wantResult := &Task{TaskId: ptr(defaultCreateTaskID), Status: &TaskStatus{State: ptr(TaskStateCompleted)}}
	if diff := cmp.Diff(wantResult, result); diff != "" {
		t.Errorf("unexpected result (-want +got):\n%s", diff)
	}
	if limiterCalls != len(responses) {
		t.Errorf("expected limiter to fire once per HTTP request (%d), got %d", len(responses), limiterCalls)
	}
}

func TestDone(t *testing.T) {
	testCases := []struct {
		name         string
		pollResponse Task
		wantDone     bool
		wantErr      string
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := mockWaitServer(t, []Task{tc.pollResponse})
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
