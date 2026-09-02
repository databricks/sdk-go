package credentials

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/databricks/databricks-sdk-go/credentials/u2m"
	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/auth/oidc"
	"github.com/databricks/sdk-go/core/profiles"
)

const authDocURL = "https://docs.databricks.com/aws/en/dev-tools/auth/index"

const defaultOIDCTokenEnv = "DATABRICKS_OIDC_TOKEN"

var (
	// ErrNoAuthConfigured is returned when no strategy in the default chain
	// could be configured from the resolved profile and environment.
	ErrNoAuthConfigured = errors.New("cannot configure default credentials")

	// ErrAuthTypeNotFound is returned when the profile requests an auth_type
	// that does not match any strategy in the default chain.
	ErrAuthTypeNotFound = errors.New("auth type not found")
)

// strategy is one entry in the default credential chain: a name plus a way to
// build credentials from a profile. Configure returns nil (without an error)
// when the strategy does not apply to the given profile, so the chain can move
// on to the next strategy.
type strategy struct {
	name                    string
	supportsGroupAssumption bool
	configure               func(profiles.Profile) (auth.Credentials, error)
}

// defaultStrategies returns the strategies tried by [NewDefaultCredentials], in
// priority order. The order mirrors the other Databricks SDKs and must not
// change without considering environments compatible with more than one
// strategy.
func defaultStrategies() []strategy {
	return []strategy{
		{
			name:      "pat",
			configure: configurePAT,
		},
		{
			name:                    "oauth-m2m",
			supportsGroupAssumption: true,
			configure:               configureM2M,
		},
		{
			name:      "databricks-cli",
			configure: configureU2M,
		},
		{
			name:                    "env-oidc",
			supportsGroupAssumption: true,
			configure:               configureEnvOIDC,
		},
		{
			name:                    "file-oidc",
			supportsGroupAssumption: true,
			configure:               configureFileOIDC,
		},
	}
}

// DefaultCredentialsOptions configures [NewDefaultCredentials].
type DefaultCredentialsOptions struct {
	// Profile is a pre-resolved profile to use. When nil, the profile is
	// resolved on first use from the default config file (~/.databrickscfg)
	// and DATABRICKS_* environment variables.
	Profile *profiles.Profile
}

// NewDefaultCredentials returns [auth.Credentials] that resolve to the first
// configured authentication strategy on first use.
//
// Strategies are tried in this order:
//  1. PAT (pat).
//  2. OAuth M2M (oauth-m2m).
//  3. Databricks CLI (databricks-cli).
//  4. Environment OIDC (env-oidc).
//  5. File OIDC (file-oidc).
//
// If the profile sets auth_type, only the strategy with that name is tried.
// Resolution is deferred until the first [auth.Credentials.AuthHeaders] call
// and then memoized, so profile resolution and any network discovery happen
// lazily and at most once.
func NewDefaultCredentials(opts DefaultCredentialsOptions) auth.Credentials {
	explicit := opts.Profile
	loadProfile := func() (profiles.Profile, error) {
		if explicit != nil {
			return *explicit, nil
		}
		p, err := profiles.Resolve()
		if err != nil {
			return profiles.Profile{}, fmt.Errorf("resolving default profile: %w", err)
		}
		return *p, nil
	}
	return &defaultCredentials{
		loadProfile: loadProfile,
		strategies:  defaultStrategies(),
	}
}

// defaultCredentials lazily resolves a profile and selects a strategy the first
// time AuthHeaders is called. It is safe for concurrent use: resolution runs
// once under [sync.Once], and the selected credentials are published through an
// [atomic.Pointer] so Name can read them without blocking on AuthHeaders.
type defaultCredentials struct {
	loadProfile func() (profiles.Profile, error)
	strategies  []strategy

	once       sync.Once
	resolved   atomic.Pointer[auth.Credentials]
	resolveErr error
}

// Name returns "default" until a strategy has been selected, then the name of
// the selected strategy so callers (logging, telemetry) can tell which
// authentication method won.
func (c *defaultCredentials) Name() string {
	if resolved := c.resolved.Load(); resolved != nil {
		return (*resolved).Name()
	}
	return "default"
}

func (c *defaultCredentials) AuthHeaders(ctx context.Context) ([]auth.Header, error) {
	c.once.Do(func() {
		creds, err := c.resolveChain()
		if err != nil {
			c.resolveErr = err
			return
		}
		c.resolved.Store(&creds)
	})
	if c.resolveErr != nil {
		return nil, c.resolveErr
	}
	return (*c.resolved.Load()).AuthHeaders(ctx)
}

func (c *defaultCredentials) resolveChain() (auth.Credentials, error) {
	profile, err := c.loadProfile()
	if err != nil {
		return nil, err
	}

	if profile.AuthType != "" {
		return c.resolveByAuthType(profile, profile.AuthType)
	}

	for _, s := range c.strategies {
		if profile.GroupID != "" && !s.supportsGroupAssumption {
			continue
		}
		creds, err := s.configure(profile)
		if err != nil {
			return nil, err
		}
		if creds != nil {
			return creds, nil
		}
	}
	return nil, fmt.Errorf("%w, please check %s to configure credentials for your preferred authentication method", ErrNoAuthConfigured, authDocURL)
}

func (c *defaultCredentials) resolveByAuthType(profile profiles.Profile, authType string) (auth.Credentials, error) {
	for _, s := range c.strategies {
		if s.name != authType {
			continue
		}
		if profile.GroupID != "" && !s.supportsGroupAssumption {
			return nil, fmt.Errorf("auth type %q does not support group role assumption. Use OAuth M2M or Workload Identity Federation", authType)
		}
		creds, err := s.configure(profile)
		if err != nil {
			return nil, err
		}
		if creds == nil {
			return nil, fmt.Errorf("%w, please check %s to configure credentials for your preferred authentication method", ErrNoAuthConfigured, authDocURL)
		}
		return creds, nil
	}
	return nil, fmt.Errorf("%w: %q, please check %s for a list of supported auth types", ErrAuthTypeNotFound, authType, authDocURL)
}

// configurePAT selects PAT credentials when the profile has a host and a token.
func configurePAT(p profiles.Profile) (auth.Credentials, error) {
	if p.Host == "" || p.Token == "" {
		return nil, nil
	}
	return NewPATCredentials(string(p.Token))
}

// configureM2M selects OAuth M2M credentials when the profile has a host plus a
// client ID and client secret. The underlying token provider is wrapped in a
// cache so a token is reused until it nears expiry rather than re-minted on
// every request.
func configureM2M(p profiles.Profile) (auth.Credentials, error) {
	if p.Host == "" || p.ClientID == "" || p.ClientSecret == "" {
		return nil, nil
	}
	provider, err := NewM2MCredentials(M2MOptions{
		Host:         p.Host,
		ClientID:     p.ClientID,
		ClientSecret: string(p.ClientSecret),
		GroupID:      p.GroupID,
	})
	if err != nil {
		return nil, err
	}
	return auth.NewTokenCredentials("oauth-m2m", auth.NewCachedTokenProvider(provider)), nil
}

// configureU2M selects Databricks CLI (U2M) credentials when the profile was
// loaded from the config file (so its section name is known) and has a host.
// The CLI must have been logged in ahead of time via "databricks auth login".
// The underlying token provider is wrapped in a cache to avoid shelling out to
// the CLI on every request.
func configureU2M(p profiles.Profile) (auth.Credentials, error) {
	if p.Host == "" || p.Name == "" {
		return nil, nil
	}
	provider, err := NewU2MCredentials(U2MOptions{
		Profile: p.Name,
		CLIPath: p.DatabricksCLIPath,
	})
	if err != nil {
		return nil, err
	}
	return auth.NewTokenCredentials("databricks-cli", auth.NewCachedTokenProvider(provider)), nil
}

func configureEnvOIDC(profile profiles.Profile) (auth.Credentials, error) {
	name := profile.OIDCTokenEnv
	if name == "" {
		name = defaultOIDCTokenEnv
	}
	if profile.Host == "" || os.Getenv(name) == "" {
		return nil, nil
	}
	return newOIDCCredentials(profile, "env-oidc", oidc.NewEnvIDTokenProvider(name)), nil
}

func configureFileOIDC(profile profiles.Profile) (auth.Credentials, error) {
	if profile.Host == "" || profile.OIDCTokenFilePath == "" {
		return nil, nil
	}
	return newOIDCCredentials(profile, "file-oidc", oidc.NewFileTokenProvider(profile.OIDCTokenFilePath)), nil
}

func newOIDCCredentials(profile profiles.Profile, name string, idTokenProvider oidc.IDTokenProvider) auth.Credentials {
	tokenProvider := oidc.NewDatabricksOIDCTokenProvider(oidc.DatabricksOIDCTokenProviderConfig{
		ClientID:              profile.ClientID,
		AccountID:             profile.AccountID,
		Host:                  profile.Host,
		GroupID:               profile.GroupID,
		TokenEndpointProvider: oidcTokenEndpointProvider(profile),
		Audience:              profile.TokenAudience,
		IDTokenProvider:       idTokenProvider,
	})
	return auth.NewTokenCredentials(name, auth.NewCachedTokenProvider(tokenProvider))
}

func oidcTokenEndpointProvider(profile profiles.Profile) func(context.Context) (*u2m.OAuthAuthorizationServer, error) {
	return func(ctx context.Context) (*u2m.OAuthAuthorizationServer, error) {
		server, err := discoverAuthorizationServer(ctx, nil, profile.Host)
		if err != nil {
			return nil, err
		}
		return &u2m.OAuthAuthorizationServer{
			AuthorizationEndpoint: server.AuthorizationEndpoint,
			TokenEndpoint:         server.TokenEndpoint,
		}, nil
	}
}
