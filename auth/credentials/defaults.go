package credentials

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/profiles"
)

// ErrCannotResolveCredentials indicates that no credential strategy in the
// chain was able to provide valid credentials.
var ErrCannotResolveCredentials = errors.New("cannot resolve default credentials, please check " + docURL)

const docURL = "https://docs.databricks.com/en/dev-tools/auth.html#databricks-client-unified-authentication"

// ResolveOption configures credential resolution.
type ResolveOption func(*resolveConfig)

type resolveConfig struct {
	profileName string
	profileFile string
	withoutEnv  bool
	httpClient  *http.Client
}

// WithProfile sets the databrickscfg profile name.
func WithProfile(name string) ResolveOption {
	return func(c *resolveConfig) { c.profileName = name }
}

// WithProfileFile overrides the path to the profile file.
func WithProfileFile(path string) ResolveOption {
	return func(c *resolveConfig) { c.profileFile = path }
}

// WithoutEnv disables the environment variable overlay during profile
// resolution. See [profiles.WithoutEnv].
func WithoutEnv() ResolveOption {
	return func(c *resolveConfig) { c.withoutEnv = true }
}

// WithHTTPClient sets the HTTP client used by strategies that make HTTP calls.
// If nil, [http.DefaultClient] is used.
func WithHTTPClient(client *http.Client) ResolveOption {
	return func(c *resolveConfig) { c.httpClient = client }
}

// ResolveCredentials resolves Databricks credentials by loading configuration
// from the databrickscfg file and environment variables, then trying each
// known authentication method in priority order.
func ResolveCredentials(ctx context.Context, opts ...ResolveOption) (auth.Credentials, error) {
	cfg := &resolveConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Build profile resolution options.
	var profileOpts []profiles.ResolveOption
	if cfg.profileName != "" {
		profileOpts = append(profileOpts, profiles.WithProfile(cfg.profileName))
	}
	if cfg.profileFile != "" {
		profileOpts = append(profileOpts, profiles.WithFile(cfg.profileFile))
	}
	if cfg.withoutEnv {
		profileOpts = append(profileOpts, profiles.WithoutEnv())
	}

	// Resolve profile from config file and environment.
	profile, err := profiles.Resolve(profileOpts...)
	if err != nil {
		return nil, err
	}

	// Run the credential chain.
	return resolveCredentials(ctx, profile, cfg)
}

// resolveCredentials tries each known credential strategy in priority order
// and returns the first one that succeeds.
func resolveCredentials(ctx context.Context, p *profiles.Profile, cfg *resolveConfig) (auth.Credentials, error) {
	type entry struct {
		name string
		make func() (auth.Credentials, error)
	}

	entries := []entry{
		{name: "pat", make: func() (auth.Credentials, error) { return resolvePAT(p) }},
	}

	// If AuthType is set, only attempt the named strategy.
	if p.AuthType != "" {
		for _, e := range entries {
			if e.name == p.AuthType {
				creds, err := e.make()
				if err != nil {
					return nil, fmt.Errorf("%s auth: %w", e.name, err)
				}
				if creds == nil {
					return nil, fmt.Errorf("%s auth: not configured", e.name)
				}
				return creds, nil
			}
		}
		return nil, fmt.Errorf("auth type %q not found", p.AuthType)
	}

	// Auto-detect: try strategies in order.
	for _, e := range entries {
		creds, err := e.make()
		if err != nil {
			slog.DebugContext(ctx, "auth strategy failed", "strategy", e.name, "error", err)
			continue
		}
		if creds == nil {
			continue
		}
		slog.DebugContext(ctx, "configured auth strategy", "strategy", e.name)
		return creds, nil
	}

	return nil, ErrCannotResolveCredentials
}

func resolvePAT(p *profiles.Profile) (auth.Credentials, error) {
	if p.Token == "" {
		return nil, nil
	}
	return NewPATCredentials(p.Token)
}
