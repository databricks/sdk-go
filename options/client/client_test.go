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
		WithConfigFile("databrickscfg"),
		WithProfile("workspace"),
		WithoutEnv(),
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
	if cfg.ConfigFile != "databrickscfg" {
		t.Errorf("ConfigFile = %q", cfg.ConfigFile)
	}
	if cfg.ProfileName != "workspace" {
		t.Errorf("ProfileName = %q", cfg.ProfileName)
	}
	if !cfg.DisableEnv {
		t.Error("DisableEnv = false, want true")
	}

	withoutConfigFileCfg := internaloptions.ClientOptions{}
	if err := WithoutConfigFile()(&withoutConfigFileCfg); err != nil {
		t.Fatalf("WithoutConfigFile: %v", err)
	}
	if !withoutConfigFileCfg.DisableConfigFile {
		t.Error("DisableConfigFile = false, want true")
	}
}

type stubCredentials struct{}

func (stubCredentials) Name() string { return "stub" }

func (stubCredentials) AuthHeaders(_ context.Context) ([]auth.Header, error) {
	return nil, nil
}
