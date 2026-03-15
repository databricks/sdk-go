package credentials

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/profiles"
)

func name(c auth.Credentials) string {
	if c == nil {
		return ""
	}
	return c.Name()
}

func TestResolveCredentials(t *testing.T) {
	testCases := []struct {
		name     string
		profile  *profiles.Profile
		wantName string // empty means expect nil credentials
		wantErr  error  // nil means expect no error
	}{
		{
			name:     "PAT from token",
			profile:  &profiles.Profile{Token: "my-token"},
			wantName: "pat",
		},
		{
			name:    "empty profile",
			profile: &profiles.Profile{},
			wantErr: ErrCannotResolveCredentials,
		},
		{
			name:     "auth type forces PAT",
			profile:  &profiles.Profile{AuthType: "pat", Token: "forced-token"},
			wantName: "pat",
		},
		{
			name:    "unknown auth type",
			profile: &profiles.Profile{AuthType: "nonexistent"},
			wantErr: ErrUnknownAuthType,
		},
		{
			name:    "auth type set but not configured",
			profile: &profiles.Profile{AuthType: "pat"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &resolveConfig{
				logger: slog.New(slog.DiscardHandler),
			}

			got, err := resolveCredentials(context.Background(), cfg, tc.profile)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("resolveCredentials() error = %v, want %v", err, tc.wantErr)
			}
			if gotName := name(got); gotName != tc.wantName {
				t.Errorf("resolveCredentials().Name() = %q, want %q", gotName, tc.wantName)
			}
		})
	}
}
