// Package profiles contains utility to resolve Databricks configuration
// profiles.
//
// A profile is a named collection of configuration values. It is typically
// stored in a file called ~/.databrickscfg. Profiles can be resolved from the
// file and/or environment variables using the Resolve function.
//
// Examples:
//
//	// Use defaults: default profile + env overlay.
//	p, err := profiles.Resolve()
//
//	// Specific profile, no env overlay.
//	p, err := profiles.Resolve(
//	    profiles.WithFile("/path/to/databrickscfg"),
//	    profiles.WithProfile("staging"),
//	)
//
//	// Specific profile with env overlay.
//	p, err := profiles.Resolve(
//	    profiles.WithProfile("staging"),
//	    profiles.WithEnv(),
//	)
//
//	// Only env vars, no config file.
//	p, err := profiles.Resolve(profiles.WithEnv())
package profiles

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"gopkg.in/ini.v1"
)

var (
	// ErrConfigFileNotFound is returned when an explicitly requested config file
	// (via [WithFile] or DATABRICKS_CONFIG_FILE) does not exist.
	ErrConfigFileNotFound = errors.New("config file not found")

	// ErrProfileNotFound is returned when the requested profile (INI section) does
	// not exist in the config file.
	ErrProfileNotFound = errors.New("profile not found")

	// ErrEmptyPath is returned when an empty path is passed.
	ErrEmptyPath = errors.New("empty path")

	// ErrEmptyProfile is returned when an empty profile name is passed.
	ErrEmptyProfile = errors.New("empty profile")

	// ErrInvalidProfileName is returned when a profile name is reserved or
	// otherwise not usable.
	ErrInvalidProfileName = errors.New("invalid profile name")
)

// settingsSection is the reserved INI section name for SDK settings. It is
// never treated as a profile. The default_profile key in this section
// determines which profile to load when no profile is explicitly requested.
const settingsSection = "__settings__"

// Secret is a string that is obfuscated in all string representations.
type Secret string

const obfuscatedSecret = "********"

func (s Secret) String() string               { return obfuscatedSecret }
func (s Secret) GoString() string             { return obfuscatedSecret }
func (s Secret) MarshalText() ([]byte, error) { return []byte(obfuscatedSecret), nil }
func (s Secret) MarshalJSON() ([]byte, error) { return []byte(obfuscatedSecret), nil }
func (s Secret) LogValue() slog.Value         { return slog.StringValue(obfuscatedSecret) }

// Profile holds configuration values resolved from a databrickscfg file and/or
// environment variables. The zero value is meaningful: all empty strings means
// nothing was configured.
//
// TODO: Consider filtering properties that are not relevant in SDK-Mod.
type Profile struct {
	// Name is the profile name (INI section name). Set automatically by
	// [Resolve]; must be set manually when constructing a Profile for use
	// with [SaveToFile].
	Name string

	// Host is the Databricks workspace or Accounts API endpoint URL.
	Host string

	// WorkspaceID is the Databricks Workspace ID, used with unified hosts.
	WorkspaceID string

	// AccountID is the Databricks Account ID for Accounts API.
	AccountID string

	// Token is the personal access token for PAT authentication.
	Token Secret

	// Username is the username for basic authentication.
	Username string

	// Password is the password for basic authentication.
	Password Secret

	// AuthType selects a specific auth method when multiple credentials are
	// available in the environment.
	AuthType string

	// ClientID is the OAuth client ID for M2M authentication.
	ClientID string

	// ClientSecret is the OAuth client secret for M2M authentication.
	ClientSecret Secret

	// DatabricksCLIPath is the path to the Databricks CLI binary
	// (version >= 0.100.0) used for CLI-based authentication.
	DatabricksCLIPath string

	// MetadataServiceURL is the URL of a metadata service that provides
	// authentication credentials (e.g. Databricks VSCode extension).
	MetadataServiceURL Secret

	// ActionsIDTokenRequestURL is the GitHub Actions URL for requesting an
	// OIDC token. Set automatically by GitHub Actions runners.
	ActionsIDTokenRequestURL string

	// ActionsIDTokenRequestToken is the bearer token used to authenticate
	// with the GitHub Actions OIDC provider.
	ActionsIDTokenRequestToken Secret

	// OIDCTokenEnv is the name of an environment variable containing an
	// OIDC ID token.
	OIDCTokenEnv string

	// OIDCTokenFilePath is the path to a file containing an OIDC ID token.
	OIDCTokenFilePath string

	// TokenAudience is the audience to specify when fetching an ID token
	// for Workload Identity Federation.
	TokenAudience string

	// DiscoveryURL is an OpenID Connect discovery URL. When set, OIDC
	// endpoints are fetched from this URL instead of the default
	// host-based well-known endpoint.
	DiscoveryURL string

	// AzureClientID is the Azure AD application (client) ID for Azure
	// client-secret or MSI authentication.
	AzureClientID string

	// AzureClientSecret is the Azure AD client secret for Azure
	// client-secret authentication.
	AzureClientSecret Secret

	// AzureTenantID is the Azure AD tenant ID.
	AzureTenantID string

	// AzureResourceID is the Azure Resource Manager ID for an Azure
	// Databricks workspace, exchanged for a host URL.
	AzureResourceID string

	// AzureEnvironment selects the Azure cloud environment (PUBLIC,
	// USGOVERNMENT, CHINA). When empty, determined from the workspace host.
	AzureEnvironment string

	// AzureLoginAppID is the Azure Login Application ID. Deprecated: this
	// field no longer has any effect and will be removed in the future.
	AzureLoginAppID string

	// AzureUseMSI enables Azure Managed Service Identity authentication.
	AzureUseMSI *bool

	// GoogleCredentials holds GCP service account credentials JSON for
	// Google credentials-based authentication.
	GoogleCredentials Secret

	// GoogleServiceAccount is the GCP service account email for
	// Google ID-based authentication.
	GoogleServiceAccount string

	// ClusterID is the default Databricks cluster ID for compute operations.
	ClusterID string

	// WarehouseID is the default Databricks SQL warehouse ID.
	WarehouseID string

	// ServerlessComputeID is the default serverless compute resource ID.
	ServerlessComputeID string

	// Extra holds INI keys that are not mapped to a known field. This is
	// only populated by [WithFile]; environment variables do not contribute
	// to Extra.
	Extra map[string]string
}

// ResolveOption configures the behavior of [Resolve].
type ResolveOption func(o *options) error

type options struct {
	withFile bool
	withEnv  bool
	filePath string
	profile  string
}

// WithFile overrides the path to the databrickscfg file. If not set, Resolve
// reads DATABRICKS_CONFIG_FILE from the environment, falling back to
// ~/.databrickscfg. An empty string returns [ErrEmptyPath].
func WithFile(path string) ResolveOption {
	return func(o *options) error {
		if path == "" {
			return ErrEmptyPath
		}
		o.filePath = path
		o.withFile = true
		return nil
	}
}

// WithDefaultProfile loads the default profile from the config file. The
// profile name is resolved in order from: DATABRICKS_CONFIG_PROFILE
// environment variable, the default_profile key in the [__settings__] section
// of the config file, and finally the DEFAULT section.
//
// Note: DATABRICKS_CONFIG_PROFILE controls which profile NAME to load. This
// is independent of [WithEnv], which controls whether profile values from
// environment variables are overlaid on top of the profile.
func WithDefaultProfile() ResolveOption {
	return func(o *options) error {
		o.withFile = true
		return nil
	}
}

// WithProfile loads a specific named profile (INI section) from the config
// file. If the profile does not exist, [Resolve] returns [ErrProfileNotFound].
// An empty string returns [ErrEmptyProfile].
func WithProfile(profile string) ResolveOption {
	return func(o *options) error {
		if profile == "" {
			return ErrEmptyProfile
		}
		o.profile = profile
		o.withFile = true
		return nil
	}
}

// WithEnv enables overlaying environment variables (e.g. DATABRICKS_HOST) on
// top of config file values. When enabled, environment variables take
// precedence over file values.
//
// This does not affect DATABRICKS_CONFIG_FILE or DATABRICKS_CONFIG_PROFILE,
// which control file and profile selection. Use [WithFile] and [WithProfile]
// to override those.
func WithEnv() ResolveOption {
	return func(o *options) error {
		o.withEnv = true
		return nil
	}
}

var defaultResolveOptions = []ResolveOption{
	WithDefaultProfile(),
	WithEnv(),
}

// Resolve creates a [Profile] from a databrickscfg file and/or environment
// variables.
//
// If no options are provided, the defaults are used:
//
//	Resolve() // equivalent to Resolve(WithDefaultProfile(), WithEnv())
func Resolve(opts ...ResolveOption) (*Profile, error) {
	if len(opts) == 0 {
		opts = defaultResolveOptions
	}

	o := options{}
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}

	p := &Profile{}
	if o.withFile {
		if err := loadFile(p, o.filePath, o.profile); err != nil {
			return nil, err
		}
	}
	if o.withEnv {
		if err := loadEnv(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// ListProfiles returns the names of all profiles (INI sections) in the given
// config file. The path must be non-empty; use [DefaultConfigFile] to obtain
// the conventional default path. The DEFAULT section is included if it has any
// keys.
func ListProfiles(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: empty path", ErrConfigFileNotFound)
	}

	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrConfigFileNotFound, path)
	}

	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return nil, fmt.Errorf("loading config file %q: %w", path, err)
	}

	var names []string
	for _, section := range f.Sections() {
		name := section.Name()
		if name == settingsSection {
			continue // the settings section is not a profile
		}
		if isPhantomDefault(section) {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// property maps a Profile field to its environment variable name and INI key
// name, with getter and setter functions. This avoids reflection while keeping
// the mapping declarative.
type property struct {
	envVar string
	iniKey string
	set    func(p *Profile, v string) error
	get    func(p *Profile) string
}

var properties = []property{
	{
		envVar: "DATABRICKS_HOST",
		iniKey: "host",
		set:    func(p *Profile, v string) error { p.Host = v; return nil },
		get:    func(p *Profile) string { return p.Host },
	},
	{
		envVar: "DATABRICKS_WORKSPACE_ID",
		iniKey: "workspace_id",
		set:    func(p *Profile, v string) error { p.WorkspaceID = v; return nil },
		get:    func(p *Profile) string { return p.WorkspaceID },
	},
	{
		envVar: "DATABRICKS_ACCOUNT_ID",
		iniKey: "account_id",
		set:    func(p *Profile, v string) error { p.AccountID = v; return nil },
		get:    func(p *Profile) string { return p.AccountID },
	},
	{
		envVar: "DATABRICKS_TOKEN",
		iniKey: "token",
		set:    func(p *Profile, v string) error { p.Token = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.Token) },
	},
	{
		envVar: "DATABRICKS_USERNAME",
		iniKey: "username",
		set:    func(p *Profile, v string) error { p.Username = v; return nil },
		get:    func(p *Profile) string { return p.Username },
	},
	{
		envVar: "DATABRICKS_PASSWORD",
		iniKey: "password",
		set:    func(p *Profile, v string) error { p.Password = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.Password) },
	},
	{
		envVar: "DATABRICKS_AUTH_TYPE",
		iniKey: "auth_type",
		set:    func(p *Profile, v string) error { p.AuthType = v; return nil },
		get:    func(p *Profile) string { return p.AuthType },
	},
	{
		envVar: "DATABRICKS_CLIENT_ID",
		iniKey: "client_id",
		set:    func(p *Profile, v string) error { p.ClientID = v; return nil },
		get:    func(p *Profile) string { return p.ClientID },
	},
	{
		envVar: "DATABRICKS_CLIENT_SECRET",
		iniKey: "client_secret",
		set:    func(p *Profile, v string) error { p.ClientSecret = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.ClientSecret) },
	},
	{
		envVar: "DATABRICKS_CLI_PATH",
		iniKey: "databricks_cli_path",
		set:    func(p *Profile, v string) error { p.DatabricksCLIPath = v; return nil },
		get:    func(p *Profile) string { return p.DatabricksCLIPath },
	},
	{
		envVar: "DATABRICKS_METADATA_SERVICE_URL",
		iniKey: "metadata_service_url",
		set:    func(p *Profile, v string) error { p.MetadataServiceURL = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.MetadataServiceURL) },
	},
	{
		envVar: "ACTIONS_ID_TOKEN_REQUEST_URL",
		iniKey: "actions_id_token_request_url",
		set:    func(p *Profile, v string) error { p.ActionsIDTokenRequestURL = v; return nil },
		get:    func(p *Profile) string { return p.ActionsIDTokenRequestURL },
	},
	{
		envVar: "ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		iniKey: "actions_id_token_request_token",
		set:    func(p *Profile, v string) error { p.ActionsIDTokenRequestToken = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.ActionsIDTokenRequestToken) },
	},
	{
		envVar: "DATABRICKS_OIDC_TOKEN_ENV",
		iniKey: "oidc_token_env",
		set:    func(p *Profile, v string) error { p.OIDCTokenEnv = v; return nil },
		get:    func(p *Profile) string { return p.OIDCTokenEnv },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_OIDC_TOKEN_FILEPATH",
		iniKey: "databricks_id_token_filepath",
		set:    func(p *Profile, v string) error { p.OIDCTokenFilePath = v; return nil },
		get:    func(p *Profile) string { return p.OIDCTokenFilePath },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_TOKEN_AUDIENCE",
		iniKey: "audience",
		set:    func(p *Profile, v string) error { p.TokenAudience = v; return nil },
		get:    func(p *Profile) string { return p.TokenAudience },
	},
	{
		envVar: "DATABRICKS_DISCOVERY_URL",
		iniKey: "discovery_url",
		set:    func(p *Profile, v string) error { p.DiscoveryURL = v; return nil },
		get:    func(p *Profile) string { return p.DiscoveryURL },
	},
	{
		envVar: "ARM_CLIENT_ID",
		iniKey: "azure_client_id",
		set:    func(p *Profile, v string) error { p.AzureClientID = v; return nil },
		get:    func(p *Profile) string { return p.AzureClientID },
	},
	{
		envVar: "ARM_CLIENT_SECRET",
		iniKey: "azure_client_secret",
		set:    func(p *Profile, v string) error { p.AzureClientSecret = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.AzureClientSecret) },
	},
	{
		envVar: "ARM_TENANT_ID",
		iniKey: "azure_tenant_id",
		set:    func(p *Profile, v string) error { p.AzureTenantID = v; return nil },
		get:    func(p *Profile) string { return p.AzureTenantID },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_AZURE_RESOURCE_ID",
		iniKey: "azure_workspace_resource_id",
		set:    func(p *Profile, v string) error { p.AzureResourceID = v; return nil },
		get:    func(p *Profile) string { return p.AzureResourceID },
	},
	{
		envVar: "ARM_ENVIRONMENT",
		iniKey: "azure_environment",
		set:    func(p *Profile, v string) error { p.AzureEnvironment = v; return nil },
		get:    func(p *Profile) string { return p.AzureEnvironment },
	},
	{
		envVar: "DATABRICKS_AZURE_LOGIN_APP_ID",
		iniKey: "azure_login_app_id",
		set:    func(p *Profile, v string) error { p.AzureLoginAppID = v; return nil },
		get:    func(p *Profile) string { return p.AzureLoginAppID },
	},
	{
		envVar: "ARM_USE_MSI",
		iniKey: "azure_use_msi",
		set: func(p *Profile, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return err
			}
			p.AzureUseMSI = &b
			return nil
		},
		get: func(p *Profile) string {
			if p.AzureUseMSI == nil {
				return ""
			}
			return strconv.FormatBool(*p.AzureUseMSI)
		},
	},
	{
		envVar: "GOOGLE_CREDENTIALS",
		iniKey: "google_credentials",
		set:    func(p *Profile, v string) error { p.GoogleCredentials = Secret(v); return nil },
		get:    func(p *Profile) string { return string(p.GoogleCredentials) },
	},
	{
		envVar: "DATABRICKS_GOOGLE_SERVICE_ACCOUNT",
		iniKey: "google_service_account",
		set:    func(p *Profile, v string) error { p.GoogleServiceAccount = v; return nil },
		get:    func(p *Profile) string { return p.GoogleServiceAccount },
	},
	{
		envVar: "DATABRICKS_CLUSTER_ID",
		iniKey: "cluster_id",
		set:    func(p *Profile, v string) error { p.ClusterID = v; return nil },
		get:    func(p *Profile) string { return p.ClusterID },
	},
	{
		envVar: "DATABRICKS_WAREHOUSE_ID",
		iniKey: "warehouse_id",
		set:    func(p *Profile, v string) error { p.WarehouseID = v; return nil },
		get:    func(p *Profile) string { return p.WarehouseID },
	},
	{
		envVar: "DATABRICKS_SERVERLESS_COMPUTE_ID",
		iniKey: "serverless_compute_id",
		set:    func(p *Profile, v string) error { p.ServerlessComputeID = v; return nil },
		get:    func(p *Profile) string { return p.ServerlessComputeID },
	},
}

// iniLoadOptions are the options used for parsing databrickscfg files.
var iniLoadOptions = ini.LoadOptions{
	// Passwords and tokens may contain '#'.
	SpaceBeforeInlineComment: true,
}

// isPhantomDefault reports whether section is an empty DEFAULT section that
// should be treated as non-existent.
//
// In the INI format, DEFAULT is typically treated as a special section that
// provides fallback values to other sections, not a regular section. This
// conflicts with Databricks's historical treatment of DEFAULT as a regular
// profile name.
//
// Known limitation: an intentionally empty [DEFAULT] section in the file will
// be silently ignored.
func isPhantomDefault(section *ini.Section) bool {
	return section != nil && section.Name() == "DEFAULT" && len(section.Keys()) == 0
}

// DefaultConfigFile returns the path to the default databrickscfg file. It
// reads DATABRICKS_CONFIG_FILE from the environment, falling back to
// ~/.databrickscfg. If the home directory cannot be determined, it returns an
// empty string.
func DefaultConfigFile() string {
	if v := os.Getenv("DATABRICKS_CONFIG_FILE"); v != "" {
		return v
	}
	return defaultConfigFilePath()
}

// defaultConfigFilePath returns ~/.databrickscfg without checking environment
// variables.
func defaultConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".databrickscfg")
}

// SaveToFile writes the profile to the [Profile.Name] section of a
// databrickscfg file. Path must be non-empty; use [DefaultConfigFile] to
// obtain the conventional default path. [Profile.Name] must be non-empty. If
// the file exists, other sections are preserved; the named section is replaced
// entirely. If the file does not exist, it is created with mode 0600.
//
// Both known fields and [Profile.Extra] entries are written, so a
// [Resolve]/SaveToFile round-trip never loses unknown keys. Fields set to the
// empty string are omitted from the output.
//
// Warning: this method will save properties that might have been loaded from
// environment variables.
func (p *Profile) SaveToFile(path string) error {
	if path == "" {
		return ErrEmptyPath
	}
	if p.Name == "" {
		return ErrEmptyProfile
	}
	if p.Name == settingsSection {
		return fmt.Errorf("%w: %q is a reserved section", ErrInvalidProfileName, settingsSection)
	}

	// Ensure the file exists with restrictive permissions before writing
	// secrets into it. O_CREATE|O_EXCL is atomic: it creates the file only
	// if it does not already exist, avoiding a Stat/WriteFile TOCTOU race.
	if f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err == nil {
		f.Close()
	} else if !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("creating config file %q: %w", path, err)
	}

	// Load the existing file.
	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return fmt.Errorf("loading config file %q: %w", path, err)
	}

	// Delete the section first to ensure a clean replacement.
	f.DeleteSection(p.Name)
	section, err := f.NewSection(p.Name)
	if err != nil {
		return fmt.Errorf("creating section %q: %w", p.Name, err)
	}

	// Build a set of known INI keys so Extra entries don't collide.
	known := make(map[string]bool, len(properties))
	for _, prop := range properties {
		known[prop.iniKey] = true
		v := prop.get(p)
		if v == "" {
			continue
		}
		if _, err := section.NewKey(prop.iniKey, v); err != nil {
			return fmt.Errorf("setting key %q: %w", prop.iniKey, err)
		}
	}

	// Write extra properties in sorted order for deterministic output.
	// Skip any key that collides with a known property.
	for _, k := range slices.SortedFunc(maps.Keys(p.Extra), cmp.Compare) {
		if known[k] {
			continue
		}
		if _, err := section.NewKey(k, p.Extra[k]); err != nil {
			return fmt.Errorf("setting key %q: %w", k, err)
		}
	}

	// Re-apply restrictive permissions after writing. ini.SaveTo may
	// create a new file or truncate the existing one, potentially
	// resetting the permissions to the user's umask default.
	if err := f.SaveTo(path); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

// loadFile populates p from the databrickscfg file. If path is empty, it
// checks DATABRICKS_CONFIG_FILE, falling back to ~/.databrickscfg. If profile
// is empty, it checks DATABRICKS_CONFIG_PROFILE, then the default_profile key
// in the config file, falling back to the DEFAULT section. A missing file is
// silently skipped unless an explicit path was provided. A missing DEFAULT
// section is not an error.
func loadFile(p *Profile, path, profile string) error {
	explicitFile := path != ""
	if path == "" {
		if v := os.Getenv("DATABRICKS_CONFIG_FILE"); v != "" {
			path = v
			explicitFile = true
		} else {
			path = defaultConfigFilePath()
		}
	}
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if explicitFile {
			return fmt.Errorf("%w: %s", ErrConfigFileNotFound, path)
		}
		return nil
	}

	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return fmt.Errorf("loading config file %q: %w", path, err)
	}

	// Profile name resolution chain:
	// 1. From environment variable
	// 2. From default_profile key in __settings__ section
	// 3. From DEFAULT section
	if profile == "" {
		profile = os.Getenv("DATABRICKS_CONFIG_PROFILE")
	}
	if profile == "" {
		if s, err := f.GetSection(settingsSection); err == nil {
			profile = s.Key("default_profile").String()
		}
	}
	if profile == settingsSection {
		return fmt.Errorf("%w: %q is a reserved section", ErrInvalidProfileName, settingsSection)
	}
	// If no profile was resolved, fall back to DEFAULT. A missing DEFAULT
	// section is not an error since nobody explicitly asked for it.
	explicitProfile := profile != ""
	if !explicitProfile {
		profile = "DEFAULT"
	}
	section, err := f.GetSection(profile)
	if err != nil || isPhantomDefault(section) {
		if explicitProfile {
			return fmt.Errorf("%w: %q in %s", ErrProfileNotFound, profile, path)
		}
		return nil
	}
	p.Name = profile

	known := make(map[string]bool, len(properties))
	for _, prop := range properties {
		known[prop.iniKey] = true
		if section.HasKey(prop.iniKey) {
			if err := prop.set(p, section.Key(prop.iniKey).String()); err != nil {
				return fmt.Errorf("setting %q: %w", prop.iniKey, err)
			}
		}
	}
	for _, key := range section.Keys() {
		if known[key.Name()] {
			continue
		}
		if p.Extra == nil {
			p.Extra = make(map[string]string)
		}
		p.Extra[key.Name()] = key.String()
	}
	return nil
}

// loadEnv populates p from environment variables. Empty environment variables
// are treated as unset and do not override existing values.
func loadEnv(p *Profile) error {
	for _, prop := range properties {
		if v := os.Getenv(prop.envVar); v != "" {
			if err := prop.set(p, v); err != nil {
				return fmt.Errorf("setting %q from %s: %w", prop.iniKey, prop.envVar, err)
			}
		}
	}
	return nil
}
