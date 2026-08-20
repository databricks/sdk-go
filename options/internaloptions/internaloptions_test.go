package internaloptions

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/databricks/sdk-go/auth"
)

// stubCredentials satisfies the credentials requirement of Resolve for tests
// that exercise other resolution behavior (e.g. logger defaulting).
type stubCredentials struct{}

func (stubCredentials) Name() string { return "stub" }

func (stubCredentials) AuthHeaders(context.Context) ([]auth.Header, error) { return nil, nil }

// isolateProfileEnv points profile resolution at a temp config file and clears
// the environment variables that could otherwise leak the developer's real
// credentials into the test.
func isolateProfileEnv(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	for _, v := range []string{
		"DATABRICKS_CONFIG_FILE", "DATABRICKS_CONFIG_PROFILE", "DATABRICKS_HOST",
		"DATABRICKS_TOKEN", "DATABRICKS_CLIENT_ID", "DATABRICKS_CLIENT_SECRET",
		"DATABRICKS_AUTH_TYPE",
	} {
		t.Setenv(v, "")
	}
}

func writeConfigFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "databrickscfg")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("writing config file: %v", err)
	}
	return path
}

func TestClientOptionsResolve_DefaultLogger(t *testing.T) {
	c := &ClientOptions{Credentials: stubCredentials{}}
	if err := c.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if c.Logger == nil {
		t.Fatal("expected default logger to be set")
	}
	if c.Logger.Enabled(context.Background(), slog.LevelError) {
		t.Fatal("expected default logger to be disabled")
	}
}

func TestClientOptionsResolve_PreservesProvidedLogger(t *testing.T) {
	provided := slog.Default()
	c := &ClientOptions{Logger: provided, Credentials: stubCredentials{}}
	if err := c.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if c.Logger != provided {
		t.Fatal("expected provided logger to be preserved")
	}
}

func TestClientOptionsResolve_DefaultCredentialsFromProfile(t *testing.T) {
	isolateProfileEnv(t)
	path := writeConfigFile(t, "[DEFAULT]\nhost = https://workspace.example\ntoken = dapi-abc\n")

	c := &ClientOptions{ProfileFile: path}
	if err := c.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if c.Credentials == nil {
		t.Fatal("expected credentials to be resolved from the profile")
	}
	// The PAT strategy should have won, and its header should carry the token.
	headers, err := c.Credentials.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders: %v", err)
	}
	want := []auth.Header{{Key: "Authorization", Value: "Bearer dapi-abc"}}
	if len(headers) != 1 || headers[0] != want[0] {
		t.Errorf("AuthHeaders() = %v, want %v", headers, want)
	}
}

func TestClientOptionsResolve_PreservesProvidedCredentials(t *testing.T) {
	isolateProfileEnv(t)
	path := writeConfigFile(t, "[DEFAULT]\nhost = https://workspace.example\ntoken = dapi-abc\n")

	provided := stubCredentials{}
	c := &ClientOptions{ProfileFile: path, Credentials: provided}
	if err := c.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if c.Credentials != auth.Credentials(provided) {
		t.Error("expected explicitly provided credentials to be preserved")
	}
}

func TestClientOptionsResolve_NoCredentialsWithoutProfileResolution(t *testing.T) {
	isolateProfileEnv(t)

	c := &ClientOptions{DisableProfileResolution: true}
	err := c.Resolve()
	if err == nil {
		t.Fatal("expected an error when no credentials are available and resolution is disabled")
	}
}

func TestClientOptionsResolve_ProfilePrecedence(t *testing.T) {
	testCases := []struct {
		name               string
		configFileContents string
		envHost            string
		envAccountID       string
		envWorkspaceID     string
		clientOptions      ClientOptions
		wantHost           string
		wantAccountID      string
		wantWorkspaceID    string
	}{
		{
			name:            "environment fills unset options",
			envHost:         "https://env.example.com",
			envAccountID:    "env-account",
			envWorkspaceID:  "env-workspace",
			wantHost:        "https://env.example.com",
			wantAccountID:   "env-account",
			wantWorkspaceID: "env-workspace",
		},
		{
			name:               "environment overrides config file values",
			configFileContents: "[DEFAULT]\nhost = https://profile.example.com\naccount_id = profile-account\nworkspace_id = profile-workspace\n",
			envHost:            "https://env.example.com",
			envAccountID:       "env-account",
			envWorkspaceID:     "env-workspace",
			wantHost:           "https://env.example.com",
			wantAccountID:      "env-account",
			wantWorkspaceID:    "env-workspace",
		},
		{
			name:           "explicit options override environment",
			envHost:        "https://env.example.com",
			envAccountID:   "env-account",
			envWorkspaceID: "env-workspace",
			clientOptions: ClientOptions{
				Host:        "https://explicit.example.com",
				AccountID:   "explicit-account",
				WorkspaceID: "explicit-workspace",
			},
			wantHost:        "https://explicit.example.com",
			wantAccountID:   "explicit-account",
			wantWorkspaceID: "explicit-workspace",
		},
		{
			name:           "disabled profile resolution ignores environment",
			envHost:        "https://env.example.com",
			envAccountID:   "env-account",
			envWorkspaceID: "env-workspace",
			clientOptions: ClientOptions{
				DisableProfileResolution: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())
			t.Setenv("DATABRICKS_CONFIG_FILE", "")
			t.Setenv("DATABRICKS_CONFIG_PROFILE", "")
			t.Setenv("DATABRICKS_HOST", tc.envHost)
			t.Setenv("DATABRICKS_ACCOUNT_ID", tc.envAccountID)
			t.Setenv("DATABRICKS_WORKSPACE_ID", tc.envWorkspaceID)

			options := tc.clientOptions
			options.Credentials = stubCredentials{}
			if tc.configFileContents != "" {
				options.ProfileFile = writeConfigFile(t, tc.configFileContents)
			}
			if err := options.Resolve(); err != nil {
				t.Fatalf("Resolve: %v", err)
			}
			if options.Host != tc.wantHost {
				t.Errorf("Host = %q, want %q", options.Host, tc.wantHost)
			}
			if options.AccountID != tc.wantAccountID {
				t.Errorf("AccountID = %q, want %q", options.AccountID, tc.wantAccountID)
			}
			if options.WorkspaceID != tc.wantWorkspaceID {
				t.Errorf("WorkspaceID = %q, want %q", options.WorkspaceID, tc.wantWorkspaceID)
			}
		})
	}
}
