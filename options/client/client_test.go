package client

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/options/internaloptions"
)

func TestOptionsApply_AllFields(t *testing.T) {
	httpClient := &http.Client{}
	creds := stubCredentials{}
	logger := slog.Default()

	opts := []Option{
		WithHost("https://example.cloud.databricks.com"),
		WithHTTPClient(httpClient),
		WithCredentials(creds),
		WithTimeout(7 * time.Second),
		WithLogger(logger),
	}

	cfg := internaloptions.ClientOptions{}
	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			t.Fatalf("Apply: %v", err)
		}
	}

	if cfg.Host != "https://example.cloud.databricks.com" {
		t.Errorf("Host = %q", cfg.Host)
	}
	if cfg.HTTPClient != httpClient {
		t.Error("HTTPClient mismatch")
	}
	if cfg.Credentials != auth.Credentials(creds) {
		t.Error("Credentials mismatch")
	}
	if cfg.Timeout != 7*time.Second {
		t.Errorf("Timeout = %v", cfg.Timeout)
	}
	if cfg.Logger != logger {
		t.Error("Logger mismatch")
	}
}

type stubCredentials struct{}

func (stubCredentials) AuthHeaders(_ context.Context) ([]auth.Header, error) {
	return nil, nil
}
