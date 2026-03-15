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
var ErrCannotResolveCredentials = errors.New("cannot resolve default credentials, please check " + authDocURL)

// ErrUnknownAuthType indicates that the requested auth type does not match
// any known credential strategy.
var ErrUnknownAuthType = errors.New("unknown auth type")

const authDocURL = "https://docs.databricks.com/en/dev-tools/auth.html#databricks-client-unified-authentication"

// ResolveOption configures credential resolution.
type ResolveOption func(*resolveConfig)

type resolveConfig struct {
	profileName string
	profileFile string
	withoutEnv  bool

	httpClient *http.Client
	logger     *slog.Logger
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

// WithLogger sets the logger used by strategies.
// If nil, a no-op logger is used.
func WithLogger(logger *slog.Logger) ResolveOption {
	return func(c *resolveConfig) { c.logger = logger }
}

// WithHTTPClient sets the HTTP client used by strategies that make HTTP calls.
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

	if cfg.logger == nil {
		cfg.logger = slog.New(slog.DiscardHandler) // no logging by default
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
	return resolveCredentials(ctx, cfg, profile)
}

type credentialBuilder struct {
	name  string
	build func(p *profiles.Profile) (auth.Credentials, error)
}

var builders = []credentialBuilder{
	{name: "pat", build: resolvePAT},
}

// resolveCredentials tries each known credential strategy in priority order
// and returns the first one that succeeds.
func resolveCredentials(ctx context.Context, cfg *resolveConfig, p *profiles.Profile) (auth.Credentials, error) {
	// If an auth type is specified, try to configure the credentials for that
	// specific auth type. If an error is encountered, return it.
	if p.AuthType != "" {
		for _, b := range builders {
			if b.name == p.AuthType {
				cfg.logger.DebugContext(ctx, "attempting to configure auth", "strategy", b.name)
				return b.build(p)
			}
		}
		return nil, fmt.Errorf("%w: %q, please check %s for a list of supported auth types", ErrUnknownAuthType, p.AuthType, authDocURL)
	}

	// If no auth type is specified, try the strategies in order. If a strategy
	// succeeds, return the credentials provider. If a strategy fails, swallow
	// the error and try the next strategy.
	for _, b := range builders {
		cfg.logger.DebugContext(ctx, "attempting to configure auth", "strategy", b.name)
		creds, err := b.build(p)
		if err != nil || creds == nil {
			cfg.logger.DebugContext(ctx, "failed to configure auth", "strategy", b.name, "error", err)
			continue
		}
		// Many credentials providers can only be truly validated after a
		// request is made (e.g. because they need to exercise some hand-shake
		// or verify that tokens exist in the cache). We perform a dry run to
		// validate the credentials provider.
		if err := dryRun(ctx, creds); err != nil {
			cfg.logger.DebugContext(ctx, "failed to configure auth", "strategy", b.name, "error", err)
			continue
		}
		return creds, nil // success
	}

	return nil, ErrCannotResolveCredentials
}

func dryRun(ctx context.Context, c auth.Credentials) error {
	_, err := c.AuthHeaders(ctx)
	return err
}

func resolvePAT(p *profiles.Profile) (auth.Credentials, error) {
	if p.Token == "" {
		return nil, nil
	}
	return NewPATCredentials(p.Token)
}
