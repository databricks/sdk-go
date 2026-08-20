package credentials

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/profiles"
	"github.com/google/go-cmp/cmp"
)

const testHost = "https://workspace.example"

// configuredStrategy returns a strategy that always builds credentials whose
// single auth header identifies the strategy by label.
func configuredStrategy(label string) strategy {
	return strategy{
		name: label,
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

func TestDefaultCredentials_Resolution(t *testing.T) {
	testCases := []struct {
		desc       string
		strategies []strategy
		profile    profiles.Profile
		wantValue  string // expected bearer token value in the Authorization header
	}{
		{
			desc:       "returns the first configured strategy",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc"},
			wantValue:  "dapi-abc",
		},
		{
			desc:       "falls through to the next strategy when earlier ones are unconfigured",
			strategies: []strategy{unconfiguredStrategy("pat"), configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost},
			wantValue:  "oauth-m2m",
		},
		{
			// PAT is configured and comes first, but auth_type pins oauth-m2m.
			desc:       "selects the strategy named by auth_type over an earlier configured strategy",
			strategies: []strategy{{name: "pat", configure: configurePAT}, configuredStrategy("oauth-m2m")},
			profile:    profiles.Profile{Host: testHost, Token: "dapi-abc", AuthType: "oauth-m2m"},
			wantValue:  "oauth-m2m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			creds := newTestChain(tc.strategies, tc.profile)
			headers, err := creds.AuthHeaders(context.Background())
			if err != nil {
				t.Fatalf("AuthHeaders() error = %v", err)
			}
			want := []auth.Header{{Key: "Authorization", Value: "Bearer " + tc.wantValue}}
			if diff := cmp.Diff(want, headers); diff != "" {
				t.Errorf("AuthHeaders() mismatch (-want +got):\n%s", diff)
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
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			creds := newTestChain(tc.strategies, tc.profile)
			_, err := creds.AuthHeaders(context.Background())
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("AuthHeaders() err = %v, want %v", err, tc.wantErr)
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
