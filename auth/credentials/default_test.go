package credentials

import (
	"context"
	"errors"
	"os"
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
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable() error = %v", err)
	}

	testCases := []struct {
		name       string
		strategies []strategy
		profile    profiles.Profile
		setEnv     func(*testing.T)
		wantName   string
	}{
		{
			name:       "returns the first configured strategy",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc"},
			wantName:   "pat",
		},
		{
			name:       "falls through to the next strategy when earlier ones are unconfigured",
			strategies: []strategy{unconfiguredStrategy("pat"), configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost},
			wantName:   "oauth-m2m",
		},
		{
			name:       "auth_type pins a strategy over an earlier configured strategy",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc", AuthType: "oauth-m2m"},
			wantName:   "oauth-m2m",
		},
		{
			name:    "auth_type pins environment OIDC when PAT is configured",
			profile: profiles.Profile{Host: testHost, AuthType: "env-oidc", Token: "earlier-pat"},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "env-oidc",
		},
		{
			name:     "grouped file OIDC selected by auth_type",
			profile:  profiles.Profile{Host: testHost, AuthType: "file-oidc", OIDCTokenFilePath: "/missing", GroupID: "group-id"},
			wantName: "file-oidc",
		},
		{
			name:    "grouped environment OIDC selected without auth_type",
			profile: profiles.Profile{Host: testHost, GroupID: "group-id"},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "env-oidc",
		},
		{
			name:    "custom environment OIDC selected without auth_type",
			profile: profiles.Profile{Host: testHost, OIDCTokenEnv: "TEST_OIDC_TOKEN"},
			setEnv: func(t *testing.T) {
				t.Setenv("TEST_OIDC_TOKEN", "id-token")
			},
			wantName: "env-oidc",
		},
		{
			name:     "file OIDC selected without auth_type",
			profile:  profiles.Profile{Host: testHost, OIDCTokenFilePath: "/missing"},
			wantName: "file-oidc",
		},
		{
			name:    "PAT precedes environment OIDC without auth_type",
			profile: profiles.Profile{Host: testHost, Token: "earlier-pat"},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "pat",
		},
		{
			name: "OAuth M2M precedes environment OIDC without auth_type",
			profile: profiles.Profile{
				Host:         testHost,
				ClientID:     "client-id",
				ClientSecret: "client-secret",
			},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "oauth-m2m",
		},
		{
			name: "CLI precedes environment OIDC without auth_type",
			profile: profiles.Profile{
				Name:              "configured-profile",
				Host:              testHost,
				DatabricksCLIPath: executable,
			},
			setEnv: func(t *testing.T) {
				t.Setenv(defaultOIDCTokenEnv, "id-token")
			},
			wantName: "databricks-cli",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isolateOIDCEnvironment(t)
			if tc.setEnv != nil {
				tc.setEnv(t)
			}
			strategies := tc.strategies
			if strategies == nil {
				strategies = defaultStrategies()
			}
			credentials, err := newTestChain(strategies, tc.profile).resolveChain()
			if err != nil {
				t.Fatalf("resolveChain() error = %v", err)
			}
			if got := credentials.Name(); got != tc.wantName {
				t.Errorf("Name() = %q, want %q", got, tc.wantName)
			}
		})
	}
}

func TestDefaultCredentials_CachesResolvedStrategy(t *testing.T) {
	buildCount := 0
	counting := strategy{
		name: "counting",
		configure: func(profiles.Profile) (auth.Credentials, error) {
			buildCount++
			return auth.NewTokenCredentials("counting", auth.TokenProviderFn(
				func(context.Context) (*auth.Token, error) {
					return &auth.Token{Value: "x"}, nil
				},
			)), nil
		},
	}
	creds := newTestChain([]strategy{counting}, profiles.Profile{})
	if _, err := creds.AuthHeaders(context.Background()); err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	if _, err := creds.AuthHeaders(context.Background()); err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	if buildCount != 1 {
		t.Errorf("configure called %d times, want 1", buildCount)
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
		desc       string
		strategies []strategy
		profile    profiles.Profile
		wantErr    error
	}{
		{
			desc:       "no strategy is configured",
			strategies: []strategy{{name: "pat", configure: configurePAT}},
			profile:    profiles.Profile{Host: testHost},
			wantErr:    ErrNoAuthConfigured,
		},
		{
			desc:       "no strategy matches auth_type",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc", AuthType: "made-up"},
			wantErr:    ErrAuthTypeNotFound,
		},
		{
			desc:       "the strategy named by auth_type is not configured",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, AuthType: "pat"},
			wantErr:    ErrNoAuthConfigured,
		},
		{
			desc:    "environment OIDC token is missing",
			profile: profiles.Profile{Host: testHost, AuthType: "env-oidc"},
			wantErr: ErrNoAuthConfigured,
		},
		{
			desc:    "file OIDC path is missing",
			profile: profiles.Profile{Host: testHost, AuthType: "file-oidc"},
			wantErr: ErrNoAuthConfigured,
		},
		{
			desc: "file OIDC host is missing",
			profile: profiles.Profile{
				AuthType:          "file-oidc",
				OIDCTokenFilePath: "/token",
			},
			wantErr: ErrNoAuthConfigured,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			isolateOIDCEnvironment(t)
			strategies := tc.strategies
			if strategies == nil {
				strategies = defaultStrategies()
			}
			creds := newTestChain(strategies, tc.profile)
			_, err := creds.AuthHeaders(context.Background())
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("AuthHeaders() err = %v, want %v", err, tc.wantErr)
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

// TestDefaultCredentials_DoesNotFallBackAfterSelectedProviderFails verifies
// that a token failure from the selected provider is returned without trying a
// later credential strategy.
func TestDefaultCredentials_DoesNotFallBackAfterSelectedProviderFails(t *testing.T) {
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

	_, err := creds.AuthHeaders(context.Background())
	if !errors.Is(err, providerErr) {
		t.Fatalf("AuthHeaders() error = %v, want %v", err, providerErr)
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

// TestNewDefaultCredentials_PATFromProfile exercises the public constructor
// end-to-end with a real strategy (PAT needs no network) via an explicit
// profile, confirming the default chain wires up correctly.
func TestNewDefaultCredentials_PATFromProfile(t *testing.T) {
	creds := NewDefaultCredentials(DefaultCredentialsOptions{
		Profile: &profiles.Profile{Host: testHost, Token: "dapi-xyz"},
	})
	headers, err := creds.AuthHeaders(context.Background())
	if err != nil {
		t.Fatalf("AuthHeaders() error = %v", err)
	}
	want := []auth.Header{{Key: "Authorization", Value: "Bearer dapi-xyz"}}
	if diff := cmp.Diff(want, headers); diff != "" {
		t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
	}
	if got := creds.Name(); got != "pat" {
		t.Errorf("Name() = %q, want %q", got, "pat")
	}
}
