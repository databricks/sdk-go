package credentials

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeCfg(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "databrickscfg")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing test config: %v", err)
	}
	return path
}

func clearProfileEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABRICKS_CONFIG_FILE", "")
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "")
	t.Setenv("DATABRICKS_HOST", "")
	t.Setenv("DATABRICKS_TOKEN", "")
	t.Setenv("DATABRICKS_AUTH_TYPE", "")
	t.Setenv("DATABRICKS_ACCOUNT_ID", "")
	t.Setenv("DATABRICKS_CLIENT_ID", "")
	t.Setenv("DATABRICKS_CLIENT_SECRET", "")
	t.Setenv("DATABRICKS_METADATA_SERVICE_URL", "")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "")
	t.Setenv("DATABRICKS_CLI_PATH", "")
}

func TestResolveCredentials_profileProvidesToken(t *testing.T) {
	cfg := writeCfg(t, `[DEFAULT]
host = https://test.cloud.databricks.com
token = profile-token
`)
	clearProfileEnv(t)

	creds, err := ResolveCredentials(context.Background(),
		WithProfileFile(cfg),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("ResolveCredentials() error: %v", err)
	}
	if creds.Name() != "pat" {
		t.Errorf("ResolveCredentials().Name() = %q, want %q", creds.Name(), "pat")
	}
}

func TestResolveCredentials_withProfileForwarded(t *testing.T) {
	cfg := writeCfg(t, `[DEFAULT]
host = https://default.cloud.databricks.com

[staging]
host = https://staging.cloud.databricks.com
token = staging-token
`)
	clearProfileEnv(t)

	creds, err := ResolveCredentials(context.Background(),
		WithProfileFile(cfg),
		WithProfile("staging"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("ResolveCredentials() error: %v", err)
	}
	if creds.Name() != "pat" {
		t.Errorf("ResolveCredentials().Name() = %q, want %q", creds.Name(), "pat")
	}
}

func TestResolveCredentials_profileResolutionError(t *testing.T) {
	clearProfileEnv(t)

	got, err := ResolveCredentials(context.Background(),
		WithProfileFile("/nonexistent/path/databrickscfg"),
		WithoutEnv(),
	)

	if err == nil {
		t.Error("ResolveCredentials() error = nil, want error for nonexistent config file")
	}
	if got != nil {
		t.Errorf("ResolveCredentials() = %v, want nil", got)
	}
}

func TestResolveCredentials_noCredentials(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("HOME", t.TempDir())

	got, err := ResolveCredentials(context.Background(),
		WithoutEnv(),
	)

	if !errors.Is(err, ErrCannotResolveCredentials) {
		t.Errorf("ResolveCredentials() error = %v, want ErrCannotResolveCredentials", err)
	}
	if got != nil {
		t.Errorf("ResolveCredentials() = %v, want nil", got)
	}
}

func TestResolveCredentials_envOverridesConfigFile(t *testing.T) {
	cfg := writeCfg(t, `[DEFAULT]
host = https://file.cloud.databricks.com
token = file-token
`)
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_TOKEN", "env-token")
	want := "Bearer env-token"

	creds, err := ResolveCredentials(context.Background(),
		WithProfileFile(cfg),
	)

	if err != nil {
		t.Fatalf("ResolveCredentials() error: %v", err)
	}
	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders() error: %v", err)
	}
	if len(headers) == 0 || headers[0].Value != want {
		t.Errorf("AuthHeaders()[0].Value = %q, want %q", headers[0].Value, want)
	}
}

func TestResolveCredentials_withoutEnvIgnoresEnvVars(t *testing.T) {
	cfg := writeCfg(t, `[DEFAULT]
host = https://file.cloud.databricks.com
token = file-token
`)
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_TOKEN", "env-token")
	want := "Bearer file-token"

	creds, err := ResolveCredentials(context.Background(),
		WithProfileFile(cfg),
		WithoutEnv(),
	)

	if err != nil {
		t.Fatalf("ResolveCredentials() error: %v", err)
	}
	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders() error: %v", err)
	}
	if len(headers) == 0 || headers[0].Value != want {
		t.Errorf("AuthHeaders()[0].Value = %q, want %q", headers[0].Value, want)
	}
}
