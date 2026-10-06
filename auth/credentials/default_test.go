package credentials

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/profiles"
	"github.com/google/go-cmp/cmp"
)

const testHost = "https://workspace.example"

func isolateOIDCEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv(defaultOIDCTokenEnv, "")
	t.Setenv("TEST_OIDC_TOKEN", "")
}

// configuredStrategy returns a strategy that always builds credentials whose
// single auth header identifies the strategy by label.
func configuredStrategy(label string) strategy {
	return strategy{
		name:                    label,
		supportsGroupAssumption: true,
		configure: func(profiles.Profile) (auth.Credentials, error) {
			return auth.NewTokenCredentials(label, auth.TokenProviderFn(
				func(context.Context) (*auth.Token, error) {
					return &auth.Token{Value: label}, nil
				},
			)), nil
		},
	}
}

// unconfiguredStrategy returns a strategy that never applies.
func unconfiguredStrategy(label string) strategy {
	return strategy{
		name:      label,
		configure: func(profiles.Profile) (auth.Credentials, error) { return nil, nil },
	}
}

func loaderFor(p profiles.Profile) func() (profiles.Profile, error) {
	return func() (profiles.Profile, error) { return p, nil }
}

func newTestChain(strategies []strategy, p profiles.Profile) *defaultCredentials {
	return &defaultCredentials{loadProfile: loaderFor(p), strategies: strategies}
}

func TestDefaultStrategies_Order(t *testing.T) {
	strategies := defaultStrategies()
	got := make([]string, len(strategies))
	for i, strategy := range strategies {
		got[i] = strategy.name
	}
	want := []string{
		"pat",
		"oauth-m2m",
		"databricks-cli",
		"env-oidc",
		"file-oidc",
	}
	if !slices.Equal(got, want) {
		t.Errorf("default strategy order = %v, want %v", got, want)
	}
}

func TestDefaultCredentials_Resolution(t *testing.T) {
	testCases := []struct {
		name        string
		strategies  []strategy
		profile     profiles.Profile
		wantName    string
		wantHeaders []auth.Header
	}{
		{
			name:        "returns the first strategy that provides headers",
			strategies:  []strategy{configuredStrategy("first"), configuredStrategy("second")},
			wantName:    "first",
			wantHeaders: []auth.Header{{Key: "Authorization", Value: "Bearer first"}},
		},
		{
			name: "falls through after a configuration error",
			strategies: []strategy{
				{
					name: "first",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return nil, errors.New("configuration failed")
					},
				},
				configuredStrategy("second"),
			},
			wantName:    "second",
			wantHeaders: []auth.Header{{Key: "Authorization", Value: "Bearer second"}},
		},
		{
			name:       "falls through when configuration does not apply",
			strategies: []strategy{unconfiguredStrategy("first"), configuredStrategy("second")},
			wantName:   "second",
			wantHeaders: []auth.Header{
				{Key: "Authorization", Value: "Bearer second"},
			},
		},
		{
			name: "falls through after token acquisition fails",
			strategies: []strategy{
				{
					name: "first",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return auth.NewTokenCredentials("first", auth.TokenProviderFn(
							func(context.Context) (*auth.Token, error) {
								return nil, errors.New("token acquisition failed")
							},
						)), nil
					},
				},
				configuredStrategy("second"),
			},
			wantName:    "second",
			wantHeaders: []auth.Header{{Key: "Authorization", Value: "Bearer second"}},
		},
		{
			name: "explicit auth type pins selected strategy",
			strategies: []strategy{
				configuredStrategy("first"),
				configuredStrategy("second"),
			},
			profile:     profiles.Profile{AuthType: "second"},
			wantName:    "second",
			wantHeaders: []auth.Header{{Key: "Authorization", Value: "Bearer second"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			credentials := newTestChain(tc.strategies, tc.profile)
			headers, err := credentials.AuthHeaders(context.Background())
			if err != nil {
				t.Fatalf("AuthHeaders() error = %v", err)
			}
			if diff := cmp.Diff(tc.wantHeaders, headers); diff != "" {
				t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
			}
			if got := credentials.Name(); got != tc.wantName {
				t.Errorf("Name() = %q, want %q", got, tc.wantName)
			}
		})
	}
}

func TestDefaultCredentials_CachesResolution(t *testing.T) {
	testCases := []struct {
		name         string
		credentials  auth.Credentials
		configureErr error
		wantErr      error
	}{
		{
			name: "successful resolution",
			credentials: auth.NewTokenCredentials("counting", auth.TokenProviderFn(
				func(context.Context) (*auth.Token, error) {
					return &auth.Token{Value: "x"}, nil
				},
			)),
		},
		{
			name:         "failed resolution",
			configureErr: errors.New("configuration failed"),
			wantErr:      ErrNoAuthConfigured,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configureCalls := 0
			counting := strategy{
				name: "counting",
				configure: func(profiles.Profile) (auth.Credentials, error) {
					configureCalls++
					return tc.credentials, tc.configureErr
				},
			}
			creds := newTestChain([]strategy{counting}, profiles.Profile{})

			for range 2 {
				_, err := creds.AuthHeaders(context.Background())
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("AuthHeaders() error = %v, want %v", err, tc.wantErr)
				}
			}
			if configureCalls != 1 {
				t.Errorf("configure calls = %d, want 1", configureCalls)
			}
		})
	}
}

func TestDefaultCredentials_InvokesLoaderExactlyOnce(t *testing.T) {
	loaderCalls := 0
	loader := func() (profiles.Profile, error) {
		loaderCalls++
		return profiles.Profile{Host: testHost, Token: "dapi-abc"}, nil
	}
	creds := &defaultCredentials{
		loadProfile: loader,
		strategies:  []strategy{{name: "pat", configure: configurePAT}},
	}
	if _, err := creds.AuthHeaders(context.Background()); err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	if _, err := creds.AuthHeaders(context.Background()); err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	if loaderCalls != 1 {
		t.Errorf("loader called %d times, want 1", loaderCalls)
	}
}

func TestDefaultCredentials_Errors(t *testing.T) {
	testCases := []struct {
		name        string
		strategies  []strategy
		profile     profiles.Profile
		wantErr     error
		wantMessage string
	}{
		{
			name: "automatic exhaustion returns generic error",
			strategies: []strategy{
				{
					name: "configuration-error",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return nil, errors.New("configuration failed")
					},
				},
				{
					name: "token-error",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return auth.NewTokenCredentials("token-error", auth.TokenProviderFn(
							func(context.Context) (*auth.Token, error) {
								return nil, errors.New("token acquisition failed")
							},
						)), nil
					},
				},
			},
			wantErr: ErrNoAuthConfigured,
			wantMessage: "cannot configure default credentials, please check " + authDocURL +
				" to configure credentials for your preferred authentication method",
		},
		{
			name: "selected strategy configuration error",
			strategies: []strategy{
				{
					name: "selected",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return nil, errTokenRequired
					},
				},
				configuredStrategy("fallback"),
			},
			profile: profiles.Profile{AuthType: "selected"},
			wantErr: errTokenRequired,
		},
		{
			name: "selected strategy initial header error",
			strategies: []strategy{
				{
					name: "selected",
					configure: func(profiles.Profile) (auth.Credentials, error) {
						return auth.NewTokenCredentials("selected", auth.TokenProviderFn(
							func(context.Context) (*auth.Token, error) {
								return nil, errTokenRequired
							},
						)), nil
					},
				},
				configuredStrategy("fallback"),
			},
			profile: profiles.Profile{AuthType: "selected"},
			wantErr: errTokenRequired,
		},
		{
			name:       "oauth m2m missing client secret",
			strategies: defaultStrategies(),
			profile: profiles.Profile{
				Host:     testHost,
				ClientID: "client-id",
				AuthType: "oauth-m2m",
			},
			wantErr: errClientSecretRequired,
		},
		{
			name:       "unknown auth type",
			strategies: defaultStrategies(),
			profile:    profiles.Profile{AuthType: "made-up"},
			wantErr:    ErrAuthTypeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creds := newTestChain(tc.strategies, tc.profile)
			_, err := creds.AuthHeaders(context.Background())
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("AuthHeaders() error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantMessage != "" && err.Error() != tc.wantMessage {
				t.Errorf("AuthHeaders() error = %q, want %q", err, tc.wantMessage)
			}
		})
	}
}

// TestDefaultCredentials_GroupedResolutionSkipsUnsupportedStrategies verifies
// that automatic resolution bypasses strategies that cannot assume a group.
func TestDefaultCredentials_GroupedResolutionSkipsUnsupportedStrategies(t *testing.T) {
	unsupportedCalls := 0
	unsupported := strategy{
		name: "pat",
		configure: func(profiles.Profile) (auth.Credentials, error) {
			unsupportedCalls++
			return NewPATCredentials("dapi-abc")
		},
	}

	creds := newTestChain(
		[]strategy{unsupported, configuredStrategy("oauth-m2m")},
		profiles.Profile{GroupID: "group-id"},
	)

	if _, err := creds.AuthHeaders(context.Background()); err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}

	if unsupportedCalls != 0 {
		t.Errorf("unsupported configure calls = %d, want 0", unsupportedCalls)
	}

	if got, want := creds.Name(), "oauth-m2m"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

// TestDefaultCredentials_GroupedExplicitAuthRejectsUnsupportedStrategy verifies
// that an explicitly selected unsupported strategy is rejected before it is
// configured.
func TestDefaultCredentials_GroupedExplicitAuthRejectsUnsupportedStrategy(t *testing.T) {
	authTypes := []string{"pat", "databricks-cli"}

	for _, authType := range authTypes {
		t.Run(authType, func(t *testing.T) {
			configureCalls := 0
			s := strategy{
				name: authType,
				configure: func(profiles.Profile) (auth.Credentials, error) {
					configureCalls++
					return nil, errors.New("must not be called")
				},
			}

			creds := newTestChain([]strategy{s}, profiles.Profile{GroupID: "group-id", AuthType: authType})

			_, err := creds.AuthHeaders(context.Background())
			if err == nil || !strings.Contains(err.Error(), "does not support group role assumption") {
				t.Fatalf("AuthHeaders() error = %v", err)
			}

			if configureCalls != 0 {
				t.Errorf("configure calls = %d, want 0", configureCalls)
			}
		})
	}
}

// TestDefaultCredentials_GroupedChainExhaustionReturnsNoAuthConfigured verifies
// that skipping every unsupported strategy preserves the generic chain error.
func TestDefaultCredentials_GroupedChainExhaustionReturnsNoAuthConfigured(t *testing.T) {
	unsupported := strategy{
		name: "pat",
		configure: func(profiles.Profile) (auth.Credentials, error) {
			return NewPATCredentials("dapi-abc")
		},
	}

	creds := newTestChain([]strategy{unsupported}, profiles.Profile{GroupID: "group-id"})

	_, err := creds.AuthHeaders(context.Background())
	if !errors.Is(err, ErrNoAuthConfigured) {
		t.Fatalf("AuthHeaders() error = %v, want %v", err, ErrNoAuthConfigured)
	}
}

// TestDefaultCredentials_BuiltInUnsupportedStrategiesRejectGroup verifies that
// the built-in PAT and CLI strategies reject explicit grouped authentication.
func TestDefaultCredentials_BuiltInUnsupportedStrategiesRejectGroup(t *testing.T) {
	testCases := []struct {
		name     string
		authType string
		profile  profiles.Profile
	}{
		{
			name:     "PAT",
			authType: "pat",
			profile: profiles.Profile{
				Host:  testHost,
				Token: "dapi-normal-access",
			},
		},
		{
			name:     "Databricks CLI",
			authType: "databricks-cli",
			profile: profiles.Profile{
				Name:              "DEFAULT",
				Host:              testHost,
				DatabricksCLIPath: "/path/that/must/not/be/invoked",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profile := tc.profile
			profile.AuthType = tc.authType
			profile.GroupID = "group-id"

			creds := NewDefaultCredentials(DefaultCredentialsOptions{Profile: &profile})

			_, err := creds.AuthHeaders(context.Background())
			if err == nil || !strings.Contains(err.Error(), "does not support group role assumption") {
				t.Fatalf("AuthHeaders() error = %v", err)
			}
		})
	}
}

func TestDefaultCredentials_FallsBackWhenProviderCannotAcquireInitialToken(t *testing.T) {
	providerErr := errors.New("group assumption rejected")
	fallbackCalls := 0
	failed := strategy{
		name:                    "oauth-m2m",
		supportsGroupAssumption: true,
		configure: func(profiles.Profile) (auth.Credentials, error) {
			return auth.NewTokenCredentials("oauth-m2m", auth.TokenProviderFn(
				func(context.Context) (*auth.Token, error) { return nil, providerErr },
			)), nil
		},
	}

	fallback := strategy{
		name:                    "fallback",
		supportsGroupAssumption: true,
		configure: func(profiles.Profile) (auth.Credentials, error) {
			fallbackCalls++
			return NewPATCredentials("dapi-fallback")
		},
	}

	creds := newTestChain([]strategy{failed, fallback}, profiles.Profile{GroupID: "group-id"})

	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	want := []auth.Header{{Key: "Authorization", Value: "Bearer dapi-fallback"}}
	if diff := cmp.Diff(want, headers); diff != "" {
		t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
	}
	if fallbackCalls != 1 {
		t.Errorf("fallback configure calls = %d, want 1", fallbackCalls)
	}
	if got, want := creds.Name(), "pat"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestDefaultCredentials_DoesNotFallBackAfterSelectedProviderFails(t *testing.T) {
	providerErr := errors.New("token refresh failed")
	failSelected := false
	fallbackCalls := 0
	selected := strategy{
		name:                    "oauth-m2m",
		supportsGroupAssumption: true,
		configure: func(profiles.Profile) (auth.Credentials, error) {
			return auth.NewTokenCredentials("oauth-m2m", auth.TokenProviderFn(
				func(context.Context) (*auth.Token, error) {
					if failSelected {
						return nil, providerErr
					}
					return &auth.Token{Value: "initial-token"}, nil
				},
			)), nil
		},
	}
	fallback := strategy{
		name:                    "fallback",
		supportsGroupAssumption: true,
		configure: func(profiles.Profile) (auth.Credentials, error) {
			fallbackCalls++
			return NewPATCredentials("dapi-fallback")
		},
	}
	creds := newTestChain([]strategy{selected, fallback}, profiles.Profile{GroupID: "group-id"})

	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("first AuthHeaders() error = %v", err)
	}
	wantHeaders := []auth.Header{{Key: "Authorization", Value: "Bearer initial-token"}}
	if diff := cmp.Diff(wantHeaders, headers); diff != "" {
		t.Errorf("first AuthHeaders() mismatch (-want +got):\n%s", diff)
	}

	failSelected = true
	_, err = creds.AuthHeaders(context.Background())
	if !errors.Is(err, providerErr) {
		t.Fatalf("second AuthHeaders() error = %v, want %v", err, providerErr)
	}
	if fallbackCalls != 0 {
		t.Errorf("fallback configure calls = %d, want 0", fallbackCalls)
	}
}

func TestDefaultCredentials_GroupedStrategiesForwardGroup(t *testing.T) {
	testCases := []struct {
		name     string
		profile  profiles.Profile
		setEnv   func(*testing.T)
		wantName string
	}{
		{
			name: "OAuth M2M wins over PAT",
			profile: profiles.Profile{
				Token:        "dapi-normal-access",
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				GroupID:      "group-id",
			},
			wantName: "oauth-m2m",
		},
		{
			name: "environment OIDC",
			profile: profiles.Profile{
				ClientID: "client-id",
				GroupID:  "group-id",
				Token:    "dapi-normal-access",
			},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "env-oidc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateOIDCEnvironment(t)
			if tc.setEnv != nil {
				tc.setEnv(t)
			}

			s := newFakeOIDCServer()
			defer s.Close()
			profile := tc.profile
			profile.Host = s.server.URL
			creds := NewDefaultCredentials(DefaultCredentialsOptions{Profile: &profile})

			headers, err := creds.AuthHeaders(context.Background())
			if err != nil {
				t.Fatalf("AuthHeaders() error = %v", err)
			}

			want := []auth.Header{{Key: "Authorization", Value: "Bearer access-token"}}
			if diff := cmp.Diff(want, headers); diff != "" {
				t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
			}

			if got := creds.Name(); got != tc.wantName {
				t.Errorf("Name() = %q, want %q", got, tc.wantName)
			}

			if got, want := s.lastRequest.Get("assume_group"), "group-id"; got != want {
				t.Errorf("assume_group = %q, want %q", got, want)
			}
		})
	}
}

func TestDefaultCredentials_Name(t *testing.T) {
	testCases := []struct {
		desc       string
		strategies []strategy
		profile    profiles.Profile
		wantName   string
	}{
		{
			desc:       "reports the first configured strategy when auth_type is not set",
			strategies: []strategy{{name: "pat", configure: configurePAT}},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc"},
			wantName:   "pat",
		},
		{
			desc:       "reports the strategy selected by auth_type",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc", AuthType: "oauth-m2m"},
			wantName:   "oauth-m2m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			creds := newTestChain(tc.strategies, tc.profile)
			if got := creds.Name(); got != "default" {
				t.Errorf("Name() before resolution = %q, want %q", got, "default")
			}
			if _, err := creds.AuthHeaders(context.Background()); err != nil {
				t.Fatalf("AuthHeaders() error = %v", err)
			}
			if got := creds.Name(); got != tc.wantName {
				t.Errorf("Name() after resolution = %q, want %q", got, tc.wantName)
			}
		})
	}
}

// TestDefaultCredentials_NameConcurrentWithAuthHeaders calls Name concurrently
// with AuthHeaders to guard against a data race on the resolved credentials.
// Meaningful under `go test -race`, which this repo's CI runs.
func TestDefaultCredentials_NameConcurrentWithAuthHeaders(t *testing.T) {
	creds := newTestChain(
		[]strategy{{name: "pat", configure: configurePAT}},
		profiles.Profile{Host: testHost, Token: "dapi-abc"},
	)
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(2)
		go func() { defer wg.Done(); _, _ = creds.AuthHeaders(context.Background()) }()
		go func() { defer wg.Done(); _ = creds.Name() }()
	}
	wg.Wait()
}

func TestNewDefaultCredentials_BuiltInStrategies(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	const (
		accountHost = "https://accounts.example.com"
		accountID   = "account-id"
	)
	testCases := []struct {
		name      string
		setup     func(*testing.T) profiles.Profile
		wantName  string
		wantToken string
	}{
		{
			name: "PAT",
			setup: func(*testing.T) profiles.Profile {
				return profiles.Profile{
					Host:     testHost,
					Token:    "dapi-xyz",
					AuthType: "pat",
				}
			},
			wantName:  "pat",
			wantToken: "dapi-xyz",
		},
		{
			name: "OAuth M2M",
			setup: func(t *testing.T) profiles.Profile {
				server := newFakeOIDCServer()
				t.Cleanup(server.Close)
				return profiles.Profile{
					Host:         server.server.URL,
					ClientID:     "client-id",
					ClientSecret: "client-secret",
					AuthType:     "oauth-m2m",
				}
			},
			wantName:  "oauth-m2m",
			wantToken: "access-token",
		},
		{
			name: "Databricks CLI",
			setup: func(t *testing.T) profiles.Profile {
				setFakeCLIEnv(t, "ok")
				t.Setenv(envFakeCLIArgs, "auth token --host "+accountHost+" --account-id "+accountID)
				return profiles.Profile{
					Host:              accountHost,
					AccountID:         accountID,
					AuthType:          "databricks-cli",
					DatabricksCLIPath: executable,
				}
			},
			wantName:  "databricks-cli",
			wantToken: "fake-access-token",
		},
		{
			name: "environment OIDC with custom variable",
			setup: func(t *testing.T) profiles.Profile {
				isolateOIDCEnvironment(t)
				t.Setenv("TEST_OIDC_TOKEN", "id-token")
				server := newFakeOIDCServer()
				t.Cleanup(server.Close)
				return profiles.Profile{
					Host:         server.server.URL,
					ClientID:     "client-id",
					AuthType:     "env-oidc",
					OIDCTokenEnv: "TEST_OIDC_TOKEN",
				}
			},
			wantName:  "env-oidc",
			wantToken: "access-token",
		},
		{
			name: "file OIDC",
			setup: func(t *testing.T) profiles.Profile {
				tokenFile := filepath.Join(t.TempDir(), "oidc-token")
				if err := os.WriteFile(tokenFile, []byte("id-token"), 0o600); err != nil {
					t.Fatalf("os.WriteFile() error = %v", err)
				}
				server := newFakeOIDCServer()
				t.Cleanup(server.Close)
				return profiles.Profile{
					Host:              server.server.URL,
					ClientID:          "client-id",
					AuthType:          "file-oidc",
					OIDCTokenFilePath: tokenFile,
				}
			},
			wantName:  "file-oidc",
			wantToken: "access-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profile := tc.setup(t)
			creds := NewDefaultCredentials(DefaultCredentialsOptions{Profile: &profile})

			headers, err := creds.AuthHeaders(context.Background())
			if err != nil {
				t.Fatalf("AuthHeaders() error = %v", err)
			}
			want := []auth.Header{{Key: "Authorization", Value: "Bearer " + tc.wantToken}}
			if diff := cmp.Diff(want, headers); diff != "" {
				t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
			}
			if got, want := creds.Name(), tc.wantName; got != want {
				t.Errorf("Name() = %q, want %q", got, want)
			}
		})
	}
}

func TestNewDefaultCredentials_U2MProfileFallsBackToHostWithOldCLI(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	const (
		host      = "https://accounts.example.com"
		accountID = "account-id"
	)
	setFakeCLIEnv(t, "unknown-profile-flag")
	t.Setenv(envFakeCLIArgs, "auth token --host "+host+" --account-id "+accountID)
	profile := profiles.Profile{
		Name:              "ACCOUNT",
		Host:              host,
		AccountID:         accountID,
		AuthType:          "databricks-cli",
		DatabricksCLIPath: executable,
	}
	creds := NewDefaultCredentials(DefaultCredentialsOptions{Profile: &profile})

	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	wantHeaders := []auth.Header{{Key: "Authorization", Value: "Bearer fallback-token"}}
	if diff := cmp.Diff(wantHeaders, headers); diff != "" {
		t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
	}
	if got, want := creds.Name(), "databricks-cli"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}
