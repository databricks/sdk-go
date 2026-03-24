package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/databricks/options"
	dq "github.com/databricks/sdk-go/dataquality/v1"
)

func ptr[T any](v T) *T { return &v }

func main() {
	host := os.Getenv("DATABRICKS_HOST")
	token := os.Getenv("DATABRICKS_TOKEN")

	if host == "" || token == "" {
		server := newMockServer()
		defer server.Close()
		host = server.URL
		token = "mock-token"
		fmt.Println("=== Running with mock server (set DATABRICKS_HOST and DATABRICKS_TOKEN for a real server) ===")
	} else {
		fmt.Printf("=== Running against %s ===\n", host)
	}

	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	creds := auth.NewTokenCredentials(auth.TokenProviderFn(
		func(context.Context) (*auth.Token, error) {
			return &auth.Token{Value: token}, nil
		},
	))

	client, err := dq.NewClient(ctx,
		options.WithHost(host),
		options.WithCredentials(creds),
		options.WithLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 1. Create a table monitor with snapshot profiling.
	fmt.Println("\n--- Create Monitor ---")
	monitor, err := client.CreateMonitor(ctx, &dq.CreateMonitorRequest{
		Monitor: &dq.Monitor{
			ObjectType: ptr("table"),
			ObjectId:   ptr("main.default.my_table"),
			DataProfilingConfig: &dq.DataProfilingConfig{
				OutputSchemaId: ptr("main.default"),
				Snapshot:       &dq.SnapshotConfig{},
			},
		},
	})
	if err != nil {
		log.Fatalf("CreateMonitor: %v", err)
	}
	fmt.Printf("Created monitor: type=%s id=%s\n", *monitor.ObjectType, *monitor.ObjectId)

	// 2. Read the monitor back.
	fmt.Println("\n--- Get Monitor ---")
	monitor, err = client.GetMonitor(ctx, &dq.GetMonitorRequest{
		ObjectType: ptr("table"),
		ObjectId:   ptr("main.default.my_table"),
	})
	if err != nil {
		log.Fatalf("GetMonitor: %v", err)
	}
	fmt.Printf("Got monitor: type=%s id=%s\n", *monitor.ObjectType, *monitor.ObjectId)

	// 3. Update the monitor to add a daily schedule.
	fmt.Println("\n--- Update Monitor ---")
	monitor, err = client.UpdateMonitor(ctx, &dq.UpdateMonitorRequest{
		ObjectType: ptr("table"),
		ObjectId:   ptr("main.default.my_table"),
		Monitor: &dq.Monitor{
			ObjectType: ptr("table"),
			ObjectId:   ptr("main.default.my_table"),
			DataProfilingConfig: &dq.DataProfilingConfig{
				OutputSchemaId: ptr("main.default"),
				Snapshot:       &dq.SnapshotConfig{},
				Schedule: &dq.CronSchedule{
					QuartzCronExpression: ptr("0 0 12 * * ?"),
					TimezoneId:           ptr("UTC"),
				},
			},
		},
		UpdateMask: ptr("data_profiling_config.schedule"),
	})
	if err != nil {
		log.Fatalf("UpdateMonitor: %v", err)
	}
	fmt.Printf("Updated monitor with schedule: %s\n", *monitor.DataProfilingConfig.Schedule.QuartzCronExpression)

	// 4. Trigger a manual refresh.
	fmt.Println("\n--- Create Refresh ---")
	refresh, err := client.CreateRefresh(ctx, &dq.CreateRefreshRequest{
		Refresh: &dq.Refresh{
			ObjectType: ptr("table"),
			ObjectId:   ptr("main.default.my_table"),
		},
	})
	if err != nil {
		log.Fatalf("CreateRefresh: %v", err)
	}
	fmt.Printf("Created refresh: id=%d state=%s\n", *refresh.RefreshId, *refresh.State)

	// 5. List all refreshes using automatic pagination.
	fmt.Println("\n--- List Refreshes (paginated) ---")
	count := 0
	for r, err := range client.ListRefreshIter(ctx, &dq.ListRefreshRequest{
		ObjectType: ptr("table"),
		ObjectId:   ptr("main.default.my_table"),
	}) {
		if err != nil {
			log.Fatalf("ListRefreshIter: %v", err)
		}
		count++
		fmt.Printf("  refresh %d: state=%s\n", *r.RefreshId, *r.State)
	}
	fmt.Printf("Total refreshes: %d\n", count)

	// 6. Delete the monitor.
	fmt.Println("\n--- Delete Monitor ---")
	err = client.DeleteMonitor(ctx, &dq.DeleteMonitorRequest{
		ObjectType: ptr("table"),
		ObjectId:   ptr("main.default.my_table"),
	})
	if err != nil {
		log.Fatalf("DeleteMonitor: %v", err)
	}
	fmt.Println("Deleted monitor for main.default.my_table")

	fmt.Println("\n=== Done ===")
}

// newMockServer returns an httptest.Server that simulates the Data Quality API.
func newMockServer() *httptest.Server {
	var nextRefreshID atomic.Int64
	nextRefreshID.Store(1)

	// Store created refreshes so the list endpoint can return them.
	var refreshes []dq.Refresh

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path

		switch {
		case r.Method == http.MethodPost && path == "/api/data-quality/v1/monitors":
			var body dq.Monitor
			json.NewDecoder(r.Body).Decode(&body)
			body.DataProfilingConfig.Status = ptr(dq.DataProfilingStatus_DataProfilingStatusPending)
			json.NewEncoder(w).Encode(body)

		case r.Method == http.MethodGet && strings.HasPrefix(path, "/api/data-quality/v1/monitors/") && !strings.Contains(path, "/refreshes"):
			json.NewEncoder(w).Encode(dq.Monitor{
				ObjectType: ptr("table"),
				ObjectId:   ptr("main.default.my_table"),
				DataProfilingConfig: &dq.DataProfilingConfig{
					OutputSchemaId: ptr("main.default"),
					Status:         ptr(dq.DataProfilingStatus_DataProfilingStatusActive),
					Snapshot:       &dq.SnapshotConfig{},
				},
			})

		case r.Method == http.MethodPatch && strings.HasPrefix(path, "/api/data-quality/v1/monitors/"):
			var body dq.Monitor
			json.NewDecoder(r.Body).Decode(&body)
			body.DataProfilingConfig.Status = ptr(dq.DataProfilingStatus_DataProfilingStatusActive)
			json.NewEncoder(w).Encode(body)

		case r.Method == http.MethodDelete && strings.HasPrefix(path, "/api/data-quality/v1/monitors/"):
			w.WriteHeader(http.StatusOK)

		case r.Method == http.MethodPost && strings.Contains(path, "/refreshes") && !strings.Contains(path, "/cancel"):
			id := nextRefreshID.Add(1) - 1
			ref := dq.Refresh{
				ObjectType: ptr("table"),
				ObjectId:   ptr("main.default.my_table"),
				RefreshId:  ptr(id),
				State:      ptr(dq.RefreshState_MonitorRefreshStatePending),
				Trigger:    ptr(dq.RefreshTrigger_MonitorRefreshTriggerManual),
			}
			refreshes = append(refreshes, ref)
			json.NewEncoder(w).Encode(ref)

		case r.Method == http.MethodGet && strings.Contains(path, "/refreshes"):
			json.NewEncoder(w).Encode(dq.ListRefreshResponse{
				Refreshes: refreshes,
			})

		default:
			http.Error(w, fmt.Sprintf(`{"error":"not found: %s %s"}`, r.Method, path), http.StatusNotFound)
		}
	}))
}
