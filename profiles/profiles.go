// Package profiles resolves Databricks configuration from INI-style config
// files (~/.databrickscfg) and environment variables.
//
// A [Profile] is a read-only snapshot of configuration values. Use [Resolve]
// to build a profile. By default, values are read from the config file first,
// then environment variables are overlaid on top. Use [ResolveOption] values
// to control which file, profile, or environment variables are used.
//
//	// Use defaults (DATABRICKS_CONFIG_FILE, DATABRICKS_CONFIG_PROFILE, env overwrite).
//	p, err := profiles.Resolve()
//
//	// Use a specific file and profile, without env overlay.
//	p, err := profiles.Resolve(
//	    profiles.WithFile("/path/to/databrickscfg"),
//	    profiles.WithProfile("staging"),
//	    profiles.WithoutEnv(),
//	)
package profiles

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

// ErrConfigFileNotFound is returned when an explicitly requested config file
// (via [WithFile] or DATABRICKS_CONFIG_FILE) does not exist.
var ErrConfigFileNotFound = errors.New("config file not found")

// ErrProfileNotFound is returned when the requested profile (INI section) does
// not exist in the config file.
var ErrProfileNotFound = errors.New("profile not found")

// Profile holds configuration values resolved from a databrickscfg file and/or
// environment variables. The zero value is meaningful: all empty strings means
// nothing was configured.
type Profile struct {
	Host        string
	WorkspaceID string
	AccountID   string
	Token       string
	Username    string
	Password    string
	AuthType    string

	// OAuth M2M
	ClientID     string
	ClientSecret string

	// Databricks CLI
	DatabricksCLIPath string

	// Metadata Service
	MetadataServiceURL string

	// OIDC
	ActionsIDTokenRequestURL   string
	ActionsIDTokenRequestToken string
	OIDCTokenEnv               string
	OIDCTokenFilePath          string
	TokenAudience              string
	DiscoveryURL               string

	// Azure
	AzureClientID     string
	AzureClientSecret string
	AzureTenantID     string
	AzureResourceID   string
	AzureEnvironment  string
	AzureLoginAppID   string
	AzureUseMSI       string

	// GCP
	GoogleCredentials    string
	GoogleServiceAccount string

	// Operational
	ClusterID           string
	WarehouseID         string
	ServerlessComputeID string

	// Extra holds INI keys that are not mapped to a known field. This is
	// only populated by [WithFile]; environment variables do not contribute
	// to Extra.
	Extra map[string]string
}

// String returns a human-readable representation of the profile. Sensitive
// fields are redacted. [Extra] values are omitted since their sensitivity is
// unknown.
func (p *Profile) String() string {
	var b strings.Builder
	for _, prop := range properties {
		v := prop.get(p)
		if v == "" {
			continue
		}
		if prop.sensitive {
			v = "********"
		}
		fmt.Fprintf(&b, "%s=%s\n", prop.iniKey, v)
	}
	return b.String()
}

// ResolveOption configures the behavior of [Resolve].
type ResolveOption func(o *options)

type options struct {
	filePath  *string
	profile   *string
	ignoreEnv bool
}

// WithFile overrides the path to the databrickscfg file. If not set, Resolve
// reads DATABRICKS_CONFIG_FILE from the environment, falling back to
// ~/.databrickscfg. An empty string is treated as if WithFile was not called.
func WithFile(path string) ResolveOption {
	return func(o *options) {
		if path != "" {
			o.filePath = &path
		}
	}
}

// WithProfile overrides the profile (INI section) to load. If not set,
// Resolve reads DATABRICKS_CONFIG_PROFILE from the environment, falling back
// to the DEFAULT section. An empty string is treated as if WithProfile was not
// called.
func WithProfile(profile string) ResolveOption {
	return func(o *options) {
		if profile != "" {
			o.profile = &profile
		}
	}
}

// WithoutEnv disables overlaying profile field environment variables (e.g.
// DATABRICKS_HOST, DATABRICKS_TOKEN) on top of the config file values. By
// default, these environment variables take precedence over file values.
//
// This does not affect DATABRICKS_CONFIG_FILE or DATABRICKS_CONFIG_PROFILE,
// which control where the config file is located. Use [WithFile] and
// [WithProfile] to override those.
func WithoutEnv() ResolveOption {
	return func(o *options) { o.ignoreEnv = true }
}

// Resolve creates a Profile from the databrickscfg file and (optionally)
// environment variables.
//
// By default, Resolve reads DATABRICKS_CONFIG_FILE and
// DATABRICKS_CONFIG_PROFILE from the environment to locate the config file and
// profile, then overlays all profile environment variables on top. Use
// [WithFile], [WithProfile], and [WithoutEnv] to override this behavior.
//
// If the config file does not exist and no file was explicitly requested
// (neither [WithFile] nor DATABRICKS_CONFIG_FILE), the file is silently
// skipped. An explicitly provided path -- via [WithFile] or the environment
// variable -- that does not exist is an error.
func Resolve(opts ...ResolveOption) (*Profile, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	filePath := ""
	explicitFile := false
	if o.filePath != nil {
		filePath = *o.filePath
		explicitFile = true
	} else if v := os.Getenv("DATABRICKS_CONFIG_FILE"); v != "" {
		filePath = v
		explicitFile = true
	}

	profile := ""
	if o.profile != nil {
		profile = *o.profile
	} else {
		profile = os.Getenv("DATABRICKS_CONFIG_PROFILE")
	}

	p := &Profile{}
	if err := loadFile(p, filePath, profile, explicitFile); err != nil {
		return nil, err
	}
	if !o.ignoreEnv {
		loadEnv(p)
	}
	return p, nil
}

// ListProfiles returns the names of all profiles (INI sections) in the given
// config file. The path must be non-empty; use [DefaultConfigFile] to obtain
// the conventional default path. The DEFAULT section is included if it has any
// keys.
func ListProfiles(path string) ([]string, error) {
	if path == "" {
		return nil, fmt.Errorf("path must not be empty")
	}

	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return nil, fmt.Errorf("loading config file %q: %w", path, err)
	}

	var names []string
	for _, section := range f.Sections() {
		name := section.Name()
		if name == "DEFAULT" && len(section.Keys()) == 0 {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// property maps a Profile field to its environment variable name and INI key
// name, with getter and setter functions. This avoids reflection while keeping
// the mapping declarative. The sensitive flag indicates whether the field
// contains secrets that should be redacted in string representations.
type property struct {
	envVar    string
	iniKey    string
	sensitive bool
	set       func(p *Profile, v string)
	get       func(p *Profile) string
}

var properties = []property{
	{
		envVar: "DATABRICKS_HOST",
		iniKey: "host",
		set:    func(p *Profile, v string) { p.Host = v },
		get:    func(p *Profile) string { return p.Host },
	},
	{
		envVar: "DATABRICKS_WORKSPACE_ID",
		iniKey: "workspace_id",
		set:    func(p *Profile, v string) { p.WorkspaceID = v },
		get:    func(p *Profile) string { return p.WorkspaceID },
	},
	{
		envVar: "DATABRICKS_ACCOUNT_ID",
		iniKey: "account_id",
		set:    func(p *Profile, v string) { p.AccountID = v },
		get:    func(p *Profile) string { return p.AccountID },
	},
	{
		envVar:    "DATABRICKS_TOKEN",
		iniKey:    "token",
		sensitive: true,
		set:       func(p *Profile, v string) { p.Token = v },
		get:       func(p *Profile) string { return p.Token },
	},
	{
		envVar: "DATABRICKS_USERNAME",
		iniKey: "username",
		set:    func(p *Profile, v string) { p.Username = v },
		get:    func(p *Profile) string { return p.Username },
	},
	{
		envVar:    "DATABRICKS_PASSWORD",
		iniKey:    "password",
		sensitive: true,
		set:       func(p *Profile, v string) { p.Password = v },
		get:       func(p *Profile) string { return p.Password },
	},
	{
		envVar: "DATABRICKS_AUTH_TYPE",
		iniKey: "auth_type",
		set:    func(p *Profile, v string) { p.AuthType = v },
		get:    func(p *Profile) string { return p.AuthType },
	},
	{
		envVar: "DATABRICKS_CLIENT_ID",
		iniKey: "client_id",
		set:    func(p *Profile, v string) { p.ClientID = v },
		get:    func(p *Profile) string { return p.ClientID },
	},
	{
		envVar:    "DATABRICKS_CLIENT_SECRET",
		iniKey:    "client_secret",
		sensitive: true,
		set:       func(p *Profile, v string) { p.ClientSecret = v },
		get:       func(p *Profile) string { return p.ClientSecret },
	},
	{
		envVar: "DATABRICKS_CLI_PATH",
		iniKey: "databricks_cli_path",
		set:    func(p *Profile, v string) { p.DatabricksCLIPath = v },
		get:    func(p *Profile) string { return p.DatabricksCLIPath },
	},
	{
		envVar:    "DATABRICKS_METADATA_SERVICE_URL",
		iniKey:    "metadata_service_url",
		sensitive: true,
		set:       func(p *Profile, v string) { p.MetadataServiceURL = v },
		get:       func(p *Profile) string { return p.MetadataServiceURL },
	},
	{
		envVar: "ACTIONS_ID_TOKEN_REQUEST_URL",
		iniKey: "actions_id_token_request_url",
		set:    func(p *Profile, v string) { p.ActionsIDTokenRequestURL = v },
		get:    func(p *Profile) string { return p.ActionsIDTokenRequestURL },
	},
	{
		envVar:    "ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		iniKey:    "actions_id_token_request_token",
		sensitive: true,
		set:       func(p *Profile, v string) { p.ActionsIDTokenRequestToken = v },
		get:       func(p *Profile) string { return p.ActionsIDTokenRequestToken },
	},
	{
		envVar: "DATABRICKS_OIDC_TOKEN_ENV",
		iniKey: "oidc_token_env",
		set:    func(p *Profile, v string) { p.OIDCTokenEnv = v },
		get:    func(p *Profile) string { return p.OIDCTokenEnv },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_OIDC_TOKEN_FILEPATH",
		iniKey: "databricks_id_token_filepath",
		set:    func(p *Profile, v string) { p.OIDCTokenFilePath = v },
		get:    func(p *Profile) string { return p.OIDCTokenFilePath },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_TOKEN_AUDIENCE",
		iniKey: "audience",
		set:    func(p *Profile, v string) { p.TokenAudience = v },
		get:    func(p *Profile) string { return p.TokenAudience },
	},
	{
		envVar: "DATABRICKS_DISCOVERY_URL",
		iniKey: "discovery_url",
		set:    func(p *Profile, v string) { p.DiscoveryURL = v },
		get:    func(p *Profile) string { return p.DiscoveryURL },
	},
	{
		envVar: "ARM_CLIENT_ID",
		iniKey: "azure_client_id",
		set:    func(p *Profile, v string) { p.AzureClientID = v },
		get:    func(p *Profile) string { return p.AzureClientID },
	},
	{
		envVar:    "ARM_CLIENT_SECRET",
		iniKey:    "azure_client_secret",
		sensitive: true,
		set:       func(p *Profile, v string) { p.AzureClientSecret = v },
		get:       func(p *Profile) string { return p.AzureClientSecret },
	},
	{
		envVar: "ARM_TENANT_ID",
		iniKey: "azure_tenant_id",
		set:    func(p *Profile, v string) { p.AzureTenantID = v },
		get:    func(p *Profile) string { return p.AzureTenantID },
	},
	{
		// Intentional mismatch between envVar and iniKey for backward compatibility.
		envVar: "DATABRICKS_AZURE_RESOURCE_ID",
		iniKey: "azure_workspace_resource_id",
		set:    func(p *Profile, v string) { p.AzureResourceID = v },
		get:    func(p *Profile) string { return p.AzureResourceID },
	},
	{
		envVar: "ARM_ENVIRONMENT",
		iniKey: "azure_environment",
		set:    func(p *Profile, v string) { p.AzureEnvironment = v },
		get:    func(p *Profile) string { return p.AzureEnvironment },
	},
	{
		envVar: "DATABRICKS_AZURE_LOGIN_APP_ID",
		iniKey: "azure_login_app_id",
		set:    func(p *Profile, v string) { p.AzureLoginAppID = v },
		get:    func(p *Profile) string { return p.AzureLoginAppID },
	},
	{
		envVar: "ARM_USE_MSI",
		iniKey: "azure_use_msi",
		set:    func(p *Profile, v string) { p.AzureUseMSI = v },
		get:    func(p *Profile) string { return p.AzureUseMSI },
	},
	{
		envVar:    "GOOGLE_CREDENTIALS",
		iniKey:    "google_credentials",
		sensitive: true,
		set:       func(p *Profile, v string) { p.GoogleCredentials = v },
		get:       func(p *Profile) string { return p.GoogleCredentials },
	},
	{
		envVar: "DATABRICKS_GOOGLE_SERVICE_ACCOUNT",
		iniKey: "google_service_account",
		set:    func(p *Profile, v string) { p.GoogleServiceAccount = v },
		get:    func(p *Profile) string { return p.GoogleServiceAccount },
	},
	{
		envVar: "DATABRICKS_CLUSTER_ID",
		iniKey: "cluster_id",
		set:    func(p *Profile, v string) { p.ClusterID = v },
		get:    func(p *Profile) string { return p.ClusterID },
	},
	{
		envVar: "DATABRICKS_WAREHOUSE_ID",
		iniKey: "warehouse_id",
		set:    func(p *Profile, v string) { p.WarehouseID = v },
		get:    func(p *Profile) string { return p.WarehouseID },
	},
	{
		envVar: "DATABRICKS_SERVERLESS_COMPUTE_ID",
		iniKey: "serverless_compute_id",
		set:    func(p *Profile, v string) { p.ServerlessComputeID = v },
		get:    func(p *Profile) string { return p.ServerlessComputeID },
	},
}

// iniLoadOptions are the options used for parsing databrickscfg files.
var iniLoadOptions = ini.LoadOptions{
	// Passwords and tokens may contain '#'.
	SpaceBeforeInlineComment: true,
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
// variables. Used by loadFile to avoid re-reading DATABRICKS_CONFIG_FILE,
// which Resolve has already handled.
func defaultConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".databrickscfg")
}

// SaveToFile writes the profile to the named section of a databrickscfg file.
// Both path and profile must be non-empty. Use [DefaultConfigFile] to obtain
// the conventional default path. If the file exists, other sections are
// preserved; the named section is replaced entirely. If the file does not
// exist, it is created with mode 0600.
//
// Both known fields and [Profile.Extra] entries are written, so a
// [Resolve]/SaveToFile round-trip never loses unknown keys. Fields set to the
// empty string are omitted from the output.
func (p *Profile) SaveToFile(path, profile string) error {
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	if profile == "" {
		return fmt.Errorf("profile must not be empty")
	}

	// Ensure the file exists with restrictive permissions before writing
	// secrets into it. O_CREATE|O_EXCL is atomic: it creates the file only
	// if it does not already exist, avoiding a Stat/WriteFile TOCTOU race.
	if f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err == nil {
		f.Close()
	}

	// Load the existing file.
	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return fmt.Errorf("loading config file %q: %w", path, err)
	}

	// Delete the section first to ensure a clean replacement.
	f.DeleteSection(profile)
	section, err := f.NewSection(profile)
	if err != nil {
		return fmt.Errorf("creating section %q: %w", profile, err)
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

	return f.SaveTo(path)
}

// loadFile populates p from the named section of a databrickscfg file. If path
// is empty, it falls back to ~/.databrickscfg. If profile is
// empty, it reads the DEFAULT section (a missing DEFAULT is not an error).
// When explicit is false and the file does not exist, loadFile returns nil
// without error.
func loadFile(p *Profile, path, profile string, explicit bool) error {
	if path == "" {
		path = defaultConfigFilePath()
	}

	// If the file doesn't exist and no explicit path was requested, skip
	// silently. This allows env-only configurations to work without a
	// config file on disk.
	if !explicit {
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			return nil
		}
	}

	if explicit {
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrConfigFileNotFound, path)
		}
	}

	f, err := ini.LoadSources(iniLoadOptions, path)
	if err != nil {
		return fmt.Errorf("loading config file %q: %w", path, err)
	}

	if profile == "" {
		profile = "DEFAULT"
	}

	section, err := f.GetSection(profile)
	if err != nil {
		// A missing DEFAULT section is not an error; leave p unchanged.
		if profile == "DEFAULT" {
			return nil
		}
		return fmt.Errorf("%w: %q in %s", ErrProfileNotFound, profile, path)
	}

	known := make(map[string]bool, len(properties))
	for _, prop := range properties {
		known[prop.iniKey] = true
		if section.HasKey(prop.iniKey) {
			prop.set(p, section.Key(prop.iniKey).String())
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
func loadEnv(p *Profile) {
	for _, prop := range properties {
		if v := os.Getenv(prop.envVar); v != "" {
			prop.set(p, v)
		}
	}
}
