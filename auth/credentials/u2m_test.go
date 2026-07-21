package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/google/go-cmp/cmp"
)

func TestNewU2MCredentials_validation(t *testing.T) {
	_, err := NewU2MCredentials(U2MOptions{})
	if !errors.Is(err, errProfileOrHostRequired) {
		t.Errorf("NewU2MCredentials() err = %v, want %v", err, errProfileOrHostRequired)
	}
}

func TestNewU2MCredentials_CLINotFound(t *testing.T) {
	_, err := NewU2MCredentials(U2MOptions{
		Host:    "https://example.cloud.databricks.com",
		CLIPath: "/path/that/does/not/exist/databricks",
	})
	if !errors.Is(err, errCLINotFound) {
		t.Errorf("NewU2MCredentials() err = %v, want %v", err, errCLINotFound)
	}
}

func TestBuildCLICommands(t *testing.T) {
	testCases := []struct {
		desc         string
		opts         U2MOptions
		wantPrimary  []string
		wantFallback []string
	}{
		{
			desc:        "profile only",
			opts:        U2MOptions{Profile: "DEFAULT"},
			wantPrimary: []string{"/bin/databricks", "auth", "token", "--profile", "DEFAULT"},
		},
		{
			desc: "profile with host fallback",
			opts: U2MOptions{Profile: "DEFAULT", Host: "https://example.com"},
			wantPrimary: []string{
				"/bin/databricks", "auth", "token", "--profile", "DEFAULT",
			},
			wantFallback: []string{
				"/bin/databricks", "auth", "token", "--host", "https://example.com",
			},
		},
		{
			desc: "host only",
			opts: U2MOptions{Host: "https://example.com"},
			wantPrimary: []string{
				"/bin/databricks", "auth", "token", "--host", "https://example.com",
			},
		},
		{
			desc: "host with account ID",
			opts: U2MOptions{Host: "https://accounts.example.com", AccountID: "acct-123"},
			wantPrimary: []string{
				"/bin/databricks", "auth", "token",
				"--host", "https://accounts.example.com",
				"--account-id", "acct-123",
			},
		},
		{
			desc: "profile with host and account ID fallback",
			opts: U2MOptions{Profile: "prod", Host: "https://accounts.example.com", AccountID: "acct-123"},
			wantPrimary: []string{
				"/bin/databricks", "auth", "token", "--profile", "prod",
			},
			wantFallback: []string{
				"/bin/databricks", "auth", "token",
				"--host", "https://accounts.example.com",
				"--account-id", "acct-123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotPrimary, gotFallback := buildCLICommands("/bin/databricks", tc.opts)
			if diff := cmp.Diff(tc.wantPrimary, gotPrimary); diff != "" {
				t.Errorf("primary command mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantFallback, gotFallback); diff != "" {
				t.Errorf("fallback command mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseExpiry(t *testing.T) {
	testCases := []struct {
		desc    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			desc:  "RFC3339Nano",
			input: "2024-03-20T10:30:00.123456789Z",
			want:  time.Date(2024, 3, 20, 10, 30, 0, 123456789, time.UTC),
		},
		{
			desc:  "RFC3339",
			input: "2024-03-20T10:30:00Z",
			want:  time.Date(2024, 3, 20, 10, 30, 0, 0, time.UTC),
		},
		{
			desc:    "unrecognized",
			input:   "not a date",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := parseExpiry(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("parseExpiry(%q) want error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseExpiry(%q) err = %v, want nil", tc.input, err)
			}
			if !got.Equal(tc.want) {
				t.Errorf("parseExpiry(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// Tests below drive the full u2mTokenProvider.Token path by asking the test
// binary to re-exec itself as a fake CLI. Behavior is selected via env vars
// set on the exec.Cmd.

const (
	envFakeMode   = "TEST_FAKE_CLI_MODE"
	envFakeExpiry = "TEST_FAKE_CLI_EXPIRY"
	envFakeToken  = "TEST_FAKE_CLI_TOKEN"
)

// TestMain intercepts invocations where the test binary is being used as a
// stand-in for the databricks CLI.
func TestMain(m *testing.M) {
	mode := os.Getenv(envFakeMode)
	if mode == "" {
		os.Exit(m.Run())
	}
	switch mode {
	case "ok":
		expiry := os.Getenv(envFakeExpiry)
		if expiry == "" {
			expiry = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
		}
		token := os.Getenv(envFakeToken)
		if token == "" {
			token = "fake-access-token"
		}
		fmt.Fprintf(os.Stdout,
			`{"access_token":%q,"token_type":"Bearer","expiry":%q}`,
			token, expiry)
		os.Exit(0)
	case "unknown-profile-flag":
		// Behave like an old CLI that does not understand --profile.
		for _, a := range os.Args[1:] {
			if a == "--profile" {
				fmt.Fprintln(os.Stderr, "Error: unknown flag: --profile")
				os.Exit(1)
			}
		}
		// Fall through and behave like "ok" for the --host invocation.
		fmt.Fprintf(os.Stdout,
			`{"access_token":"fallback-token","token_type":"Bearer","expiry":%q}`,
			time.Now().Add(1*time.Hour).UTC().Format(time.RFC3339))
		os.Exit(0)
	case "error":
		fmt.Fprintln(os.Stderr, "Error: something broke")
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown fake CLI mode: %s\n", mode)
		os.Exit(2)
	}
}

// execAsFakeCLI configures the test binary to re-exec itself as the fake
// databricks CLI by wiring env vars through exec.Command. Since u2mTokenProvider
// uses exec.CommandContext directly, we rely on the inherited process env.
// We therefore set the env var on the caller before invoking Token.
func setFakeCLIEnv(t *testing.T, mode string) {
	t.Helper()
	t.Setenv(envFakeMode, mode)
}

func TestU2M_Token_Success(t *testing.T) {
	setFakeCLIEnv(t, "ok")
	t.Setenv(envFakeToken, "hello")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	tp, err := NewU2MCredentials(U2MOptions{
		Host:    "https://example.cloud.databricks.com",
		CLIPath: exe,
	})
	if err != nil {
		t.Fatalf("NewU2MCredentials() error = %v", err)
	}

	got, err := tp.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if got.Value != "hello" {
		t.Errorf("Token().Value = %q, want %q", got.Value, "hello")
	}
	if got.Type != "Bearer" {
		t.Errorf("Token().Type = %q, want %q", got.Type, "Bearer")
	}
	if got.Expiry.IsZero() {
		t.Errorf("Token().Expiry is zero, want non-zero")
	}
}

func TestU2M_Token_Error(t *testing.T) {
	setFakeCLIEnv(t, "error")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	tp, err := NewU2MCredentials(U2MOptions{
		Host:    "https://example.cloud.databricks.com",
		CLIPath: exe,
	})
	if err != nil {
		t.Fatalf("NewU2MCredentials() error = %v", err)
	}

	_, err = tp.Token(context.Background())
	if err == nil {
		t.Fatalf("Token() err = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "something broke") {
		t.Errorf("Token() err = %v, want to contain %q", err, "something broke")
	}
}

func TestU2M_Token_UnknownProfileFlagFallback(t *testing.T) {
	setFakeCLIEnv(t, "unknown-profile-flag")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	tp, err := NewU2MCredentials(U2MOptions{
		Profile: "DEFAULT",
		Host:    "https://example.cloud.databricks.com",
		CLIPath: exe,
	})
	if err != nil {
		t.Fatalf("NewU2MCredentials() error = %v", err)
	}

	got, err := tp.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if got.Value != "fallback-token" {
		t.Errorf("Token().Value = %q, want %q (fallback path)", got.Value, "fallback-token")
	}
}

func TestU2M_Token_UnknownProfileFlagNoFallback(t *testing.T) {
	setFakeCLIEnv(t, "unknown-profile-flag")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	tp, err := NewU2MCredentials(U2MOptions{
		Profile: "DEFAULT",
		CLIPath: exe,
	})
	if err != nil {
		t.Fatalf("NewU2MCredentials() error = %v", err)
	}

	_, err = tp.Token(context.Background())
	if err == nil {
		t.Fatal("Token() err = nil, want an error because no --host fallback")
	}
	if !strings.Contains(err.Error(), "unknown flag: --profile") {
		t.Errorf("Token() err = %v, want to contain %q", err, "unknown flag: --profile")
	}
}

// We exercise one path where the CLI writes a valid token but with no Type,
// to confirm we propagate whatever the CLI says.
func TestU2M_Token_Acts_As_OpaqueTokenBearer(t *testing.T) {
	setFakeCLIEnv(t, "ok")
	t.Setenv(envFakeToken, "opaque-token")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	tp, err := NewU2MCredentials(U2MOptions{Host: "https://x", CLIPath: exe})
	if err != nil {
		t.Fatalf("NewU2MCredentials() error = %v", err)
	}

	got, err := tp.Token(context.Background())
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	want := &auth.Token{Value: "opaque-token", Type: "Bearer"}
	if got.Value != want.Value || got.Type != want.Type {
		t.Errorf("Token() = %+v, want Value=%q Type=%q", got, want.Value, want.Type)
	}
}
