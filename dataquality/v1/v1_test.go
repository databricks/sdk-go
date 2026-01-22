package v1

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/databricks/options"
)

func ptr[T any](v T) *T {
	return &v
}

func TestListRefreshIter_Pagination(t *testing.T) {
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		pageToken := r.URL.Query().Get("page_token")

		var resp ListRefreshResponse
		switch pageToken {
		case "":
			resp = ListRefreshResponse{
				Refreshes:     []Refresh{{RefreshId: ptr(int64(1))}, {RefreshId: ptr(int64(2))}},
				NextPageToken: ptr("page2"),
			}
		case "page2":
			resp = ListRefreshResponse{
				Refreshes:     []Refresh{{RefreshId: ptr(int64(3))}},
				NextPageToken: ptr(""),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	auth := auth.NewTokenCredentials(auth.TokenProviderFn(func(ctx context.Context) (*auth.Token, error) {
		return &auth.Token{
			Value: "token",
		}, nil
	}))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	client, err := NewClient(context.Background(), options.WithHost(server.URL), options.WithCredentials(auth), options.WithLogger(logger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var gotIDs []int64
	for refresh, err := range client.ListRefreshIter(context.Background(), &ListRefreshRequest{
		ObjectType: ptr("table"),
		ObjectId:   ptr("obj1"),
	}) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		gotIDs = append(gotIDs, *refresh.RefreshId)
	}

	// Verify all items collected across pages
	if len(gotIDs) != 3 {
		t.Errorf("got %d items, want 3", len(gotIDs))
	}

	// Verify correct number of HTTP requests
	if requestCount != 2 {
		t.Errorf("got %d requests, want 2", requestCount)
	}
}
