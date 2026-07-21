package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/databricks/sdk-go/auth"
)

var (
	errProfileOrHostRequired = errors.New("either Profile or Host is required")
	errCLINotFound           = errors.New("databricks CLI not found")
)

// U2MOptions configures a user-to-machine [auth.TokenProvider] that obtains
// tokens by shelling out to the Databricks CLI. The CLI must already be
// authenticated (via "databricks auth login"); this provider does not perform
// an interactive OAuth flow.
type U2MOptions struct {
	// Profile is the databricks CLI profile name. When set, the CLI is
	// invoked with --profile. Host, when also set, is used as a fallback
	// for older CLI versions that do not support --profile.
	Profile string

	// Host is the workspace or account URL. Required when Profile is empty.
	Host string

	// AccountID is passed as --account-id when invoking the CLI via --host
	// for account-level hosts. Optional.
	AccountID string

	// CLIPath overrides the "databricks" binary to execute. If empty, the
	// binary is looked up in PATH.
	CLIPath string
}

// NewU2MCredentials returns an [auth.TokenProvider] that obtains tokens by
// shelling out to the Databricks CLI.
//
// The returned provider does not cache tokens or retry on failure. Wrap it
// with [auth.NewCachedTokenProvider] and [retrying.NewTokenProvider] as
// needed.
func NewU2MCredentials(opts U2MOptions) (auth.TokenProvider, error) {
	if opts.Profile == "" && opts.Host == "" {
		return nil, errProfileOrHostRequired
	}
	cliPath, err := resolveCLIPath(opts.CLIPath)
	if err != nil {
		return nil, err
	}
	primary, fallback := buildCLICommands(cliPath, opts)
	return &u2mTokenProvider{cmd: primary, hostCmd: fallback}, nil
}

type u2mTokenProvider struct {
	// cmd is the primary command to execute. Uses --profile when available.
	cmd []string

	// hostCmd is the fallback command. Uses --host, and is only set when the
	// primary command uses --profile. It is tried only if the primary call
	// fails with an "unknown flag: --profile" error, indicating an older CLI.
	hostCmd []string
}

// cliTokenResponse is the JSON shape returned by "databricks auth token".
type cliTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expiry      string `json:"expiry"`
}

func (u *u2mTokenProvider) Token(ctx context.Context) (*auth.Token, error) {
	tok, err := execCLI(ctx, u.cmd)
	if err != nil && u.hostCmd != nil && isUnknownFlagError(err) {
		return execCLI(ctx, u.hostCmd)
	}
	return tok, err
}

func execCLI(ctx context.Context, args []string) (*auth.Token, error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("databricks CLI: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, fmt.Errorf("databricks CLI: %w", err)
	}
	var resp cliTokenResponse
	if err := json.Unmarshal(stdout, &resp); err != nil {
		return nil, fmt.Errorf("parsing CLI response: %w", err)
	}
	expiry, err := parseExpiry(resp.Expiry)
	if err != nil {
		return nil, fmt.Errorf("parsing token expiry: %w", err)
	}
	return &auth.Token{
		Value:  resp.AccessToken,
		Type:   resp.TokenType,
		Expiry: expiry,
	}, nil
}

// buildCLICommands constructs the CLI command(s) to fetch an auth token.
//
// When Profile is set, the primary command uses --profile and a fallback
// --host command is also returned (used if the primary command fails because
// the CLI is too old to support --profile). When Profile is empty, only the
// --host command is returned.
func buildCLICommands(cliPath string, opts U2MOptions) (primary, fallback []string) {
	if opts.Profile != "" {
		primary = []string{cliPath, "auth", "token", "--profile", opts.Profile}
		if opts.Host != "" {
			fallback = buildHostCommand(cliPath, opts)
		}
		return primary, fallback
	}
	return buildHostCommand(cliPath, opts), nil
}

func buildHostCommand(cliPath string, opts U2MOptions) []string {
	cmd := []string{cliPath, "auth", "token", "--host", opts.Host}
	if opts.AccountID != "" {
		cmd = append(cmd, "--account-id", opts.AccountID)
	}
	return cmd
}

// isUnknownFlagError reports whether err indicates the CLI does not recognize
// the --profile flag. This happens with CLI versions that predate
// profile-based token lookup.
func isUnknownFlagError(err error) bool {
	return strings.Contains(err.Error(), "unknown flag: --profile")
}

// parseExpiry accepts the subset of formats the modern Databricks CLI emits.
func parseExpiry(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unrecognized expiry format %q", s)
}

func resolveCLIPath(override string) (string, error) {
	name := "databricks"
	if override != "" {
		name = override
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%w: %v", errCLINotFound, err)
	}
	return path, nil
}
