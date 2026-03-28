package profiles

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// resetEnv clears all environment variables that could affect profile
// resolution and sets HOME to a temp directory so the real ~/.databrickscfg is
// never consulted. Call this at the top of every test.
func resetEnv(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("DATABRICKS_CONFIG_FILE", "")
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "")
	for _, prop := range properties {
		t.Setenv(prop.envVar, "")
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		opts    []ResolveOption
		env     map[string]string
		want    *Profile
		wantErr error
	}{
		{
			name: "withFileAndProfile",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("workspace"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:         "workspace",
				Host:         "https://workspace.cloud.databricks.com",
				Token:        Secret("workspace-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
		},
		{
			name: "defaultSection",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:  "DEFAULT",
				Host:  "https://default.cloud.databricks.com",
				Token: Secret("default-token"),
			},
		},
		{
			name:    "missingExplicitFile",
			opts:    []ResolveOption{WithFile("testdata/nonexistent"), WithoutEnv()},
			wantErr: ErrConfigFileNotFound,
		},
		{
			name:    "missingEnvConfigFile",
			env:     map[string]string{"DATABRICKS_CONFIG_FILE": "testdata/nonexistent"},
			wantErr: ErrConfigFileNotFound,
		},
		{
			name: "missingDefaultFileSkipped",
			env:  map[string]string{"DATABRICKS_HOST": "https://env.cloud.databricks.com"},
			want: &Profile{Host: "https://env.cloud.databricks.com"},
		},
		{
			name: "missingExplicitProfile",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("nonexistent"),
				WithoutEnv(),
			},
			wantErr: ErrProfileNotFound,
		},
		{
			name: "missingDefaultSection",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_no_default"),
				WithoutEnv(),
			},
			want: &Profile{},
		},
		{
			name: "hashInValues",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("hash-in-value"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:         "hash-in-value",
				Host:         "https://hash.cloud.databricks.com",
				Token:        Secret("abc#def#ghi"),
				ClientSecret: Secret("secret#with#hashes"),
			},
		},
		{
			name: "azure",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("azure"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:              "azure",
				Host:              "https://adb-123.azuredatabricks.net",
				AzureClientID:     "az-client-id",
				AzureClientSecret: Secret("az-client-secret"),
				AzureTenantID:     "az-tenant-id",
				AzureResourceID:   "/subscriptions/sub-id/resourceGroups/rg/providers/Microsoft.Databricks/workspaces/ws",
			},
		},
		{
			name: "emptySection",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("empty"),
				WithoutEnv(),
			},
			want: &Profile{Name: "empty"},
		},
		{
			name: "envOnly",
			opts: []ResolveOption{WithFile("testdata/databrickscfg_no_default")},
			env: map[string]string{
				"DATABRICKS_HOST":      "https://env.cloud.databricks.com",
				"DATABRICKS_TOKEN":     "env-token",
				"DATABRICKS_CLIENT_ID": "env-client-id",
			},
			want: &Profile{
				Host:     "https://env.cloud.databricks.com",
				Token:    Secret("env-token"),
				ClientID: "env-client-id",
			},
		},
		{
			name: "envOverridesFile",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("workspace"),
			},
			env: map[string]string{
				"DATABRICKS_HOST":  "https://env-override.cloud.databricks.com",
				"DATABRICKS_TOKEN": "env-override-token",
			},
			want: &Profile{
				Name:         "workspace",
				Host:         "https://env-override.cloud.databricks.com",
				Token:        Secret("env-override-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
		},
		{
			name: "noOptions",
			env: map[string]string{
				"DATABRICKS_CONFIG_FILE":    "testdata/databrickscfg",
				"DATABRICKS_CONFIG_PROFILE": "workspace",
			},
			want: &Profile{
				Name:         "workspace",
				Host:         "https://workspace.cloud.databricks.com",
				Token:        Secret("workspace-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
		},
		{
			name: "noOptionsEnvOverridesFile",
			env: map[string]string{
				"DATABRICKS_CONFIG_FILE":    "testdata/databrickscfg",
				"DATABRICKS_CONFIG_PROFILE": "workspace",
				"DATABRICKS_HOST":           "https://override.cloud.databricks.com",
			},
			want: &Profile{
				Name:         "workspace",
				Host:         "https://override.cloud.databricks.com",
				Token:        Secret("workspace-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
		},
		{
			name: "withoutEnv",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("workspace"),
				WithoutEnv(),
			},
			env: map[string]string{
				"DATABRICKS_HOST": "https://should-be-ignored.cloud.databricks.com",
			},
			want: &Profile{
				Name:         "workspace",
				Host:         "https://workspace.cloud.databricks.com",
				Token:        Secret("workspace-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
		},
		{
			name: "extraKeys",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile("extra-keys"),
				WithoutEnv(),
			},
			want: &Profile{
				Name: "extra-keys",
				Host: "https://extra.cloud.databricks.com",
				Extra: map[string]string{
					"custom_key":  "custom-value",
					"another_key": "another-value",
				},
			},
		},
		{
			name:    "withFileEmptyStringIsError",
			opts:    []ResolveOption{WithFile("")},
			wantErr: ErrEmptyPath,
		},
		{
			name:    "withProfileEmptyStringIsError",
			opts:    []ResolveOption{WithProfile("")},
			wantErr: ErrEmptyProfile,
		},
		{
			name: "defaultProfileResolves",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:  "my-workspace",
				Host:  "https://my-workspace.cloud.databricks.com",
				Token: Secret("my-workspace-token"),
			},
		},
		{
			name: "defaultProfileFallbackEmptyKey",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings_empty"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:  "DEFAULT",
				Host:  "https://default.cloud.databricks.com",
				Token: Secret("default-token"),
			},
		},
		{
			name: "defaultProfileOverriddenByExplicitProfile",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings"),
				WithProfile("DEFAULT"),
				WithoutEnv(),
			},
			want: &Profile{
				Name:  "DEFAULT",
				Host:  "https://default.cloud.databricks.com",
				Token: Secret("default-token"),
			},
		},
		{
			name: "defaultProfileOverriddenByEnvProfile",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings"),
				WithoutEnv(),
			},
			env: map[string]string{
				"DATABRICKS_CONFIG_PROFILE": "DEFAULT",
			},
			want: &Profile{
				Name:  "DEFAULT",
				Host:  "https://default.cloud.databricks.com",
				Token: Secret("default-token"),
			},
		},
		{
			name: "defaultProfileSelfReference",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings_self_ref"),
				WithoutEnv(),
			},
			wantErr: ErrInvalidProfileName,
		},
		{
			name: "defaultProfileNonexistentSection",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings_nonexistent"),
				WithoutEnv(),
			},
			wantErr: ErrProfileNotFound,
		},
		{
			name: "explicitSettingsProfileRejected",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg_settings"),
				WithProfile("__settings__"),
				WithoutEnv(),
			},
			wantErr: ErrInvalidProfileName,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			got, gotErr := Resolve(tc.opts...)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Resolve() error = %v, want %v", gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestResolve_allProperties verifies that every property in the mapping table
// round-trips correctly through INI file loading and env var loading. This
// catches copy-paste bugs in the property set/get closures.
func TestResolve_allProperties(t *testing.T) {
	// Build an INI file with every known key set to a unique value.
	var iniLines []string
	iniLines = append(iniLines, "[all-fields]")
	envValues := make(map[string]string)
	for i, prop := range properties {
		var fileVal, envVal string
		if prop.iniKey == "azure_use_msi" {
			fileVal = "true"
			envVal = "true"
		} else {
			fileVal = strings.Repeat("f", i+1) // "f", "ff", "fff", ...
			envVal = strings.Repeat("e", i+1)  // "e", "ee", "eee", ...
		}
		iniLines = append(iniLines, prop.iniKey+" = "+fileVal)
		envValues[prop.envVar] = envVal
	}
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "databrickscfg")
	if err := os.WriteFile(cfgPath, []byte(strings.Join(iniLines, "\n")+"\n"), 0600); err != nil {
		t.Fatalf("writing test config: %v", err)
	}

	// Test file loading: every field should be set to its file value.
	fileProfile, err := Resolve(WithFile(cfgPath), WithProfile("all-fields"), WithoutEnv())
	if err != nil {
		t.Fatalf("Resolve() file error: %v", err)
	}
	for i, prop := range properties {
		var want string
		if prop.iniKey == "azure_use_msi" {
			want = "true"
		} else {
			want = strings.Repeat("f", i+1)
		}
		got := prop.get(fileProfile)
		if got != want {
			t.Errorf("file: property %q (iniKey=%q) = %q, want %q", prop.envVar, prop.iniKey, got, want)
		}
	}

	// Test env loading: set all env vars and resolve with an empty file.
	for envVar, val := range envValues {
		t.Setenv(envVar, val)
	}
	emptyCfg := filepath.Join(dir, "empty_databrickscfg")
	if err := os.WriteFile(emptyCfg, []byte("[DEFAULT]\n"), 0600); err != nil {
		t.Fatalf("writing empty config: %v", err)
	}
	envProfile, err := Resolve(WithFile(emptyCfg))
	if err != nil {
		t.Fatalf("Resolve() env error: %v", err)
	}
	for i, prop := range properties {
		var want string
		if prop.iniKey == "azure_use_msi" {
			want = "true"
		} else {
			want = strings.Repeat("e", i+1)
		}
		got := prop.get(envProfile)
		if got != want {
			t.Errorf("env: property %q (envVar=%q) = %q, want %q", prop.iniKey, prop.envVar, got, want)
		}
	}
}

// TestResolve_allPropertiesCovered verifies that every exported field on
// Profile (except Extra) has a corresponding entry in the properties table.
func TestResolve_allPropertiesCovered(t *testing.T) {
	// Build the set of fields that are touched by at least one property.
	covered := make(map[string]bool)
	for _, prop := range properties {
		before := reflect.ValueOf(Profile{})
		p := &Profile{}
		if err := prop.set(p, "true"); err != nil {
			t.Fatalf("set %q: %v", prop.iniKey, err)
		}
		after := reflect.ValueOf(*p)
		for i := range before.NumField() {
			if !reflect.DeepEqual(before.Field(i).Interface(), after.Field(i).Interface()) {
				covered[reflect.TypeOf(Profile{}).Field(i).Name] = true
			}
		}
	}

	typ := reflect.TypeOf(Profile{})
	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.Name == "Name" || f.Name == "Extra" {
			continue
		}
		if !covered[f.Name] {
			t.Errorf("Profile.%s has no entry in properties table", f.Name)
		}
	}
}

func TestProfile_SaveToFile(t *testing.T) {
	tests := []struct {
		desc     string
		existing *string // initial ini content; nil means no file
		profile  *Profile
		want     string
		wantErr  error
	}{
		{
			desc:    "known fields are written to the ini file",
			profile: &Profile{Name: "test", Host: "https://saved.cloud.databricks.com", Token: Secret("saved-token"), ClientID: "saved-client-id"},
			want: `[test]
host      = https://saved.cloud.databricks.com
token     = saved-token
client_id = saved-client-id
`,
		},
		{
			desc:    "extra keys are written to the ini file",
			profile: &Profile{Name: "test", Host: "https://extra.cloud.databricks.com", Extra: map[string]string{"custom_key": "custom-value", "another_key": "another-value"}},
			want: `[test]
host        = https://extra.cloud.databricks.com
another_key = another-value
custom_key  = custom-value
`,
		},
		{
			desc:    "extra keys that collide with known fields are skipped",
			profile: &Profile{Name: "test", Host: "https://real.cloud.databricks.com", Extra: map[string]string{"host": "https://evil.cloud.databricks.com"}},
			want: `[test]
host = https://real.cloud.databricks.com
`,
		},
		{
			desc:    "empty fields are omitted",
			profile: &Profile{Name: "my-profile", Host: "https://host.cloud.databricks.com"},
			want: `[my-profile]
host = https://host.cloud.databricks.com
`,
		},
		{
			desc: "other sections in the file are preserved",
			existing: new(`[other]
host  = https://other.cloud.databricks.com
token = other-token
`),
			profile: &Profile{Name: "new", Host: "https://new.cloud.databricks.com"},
			want: `[other]
host  = https://other.cloud.databricks.com
token = other-token

[new]
host = https://new.cloud.databricks.com
`,
		},
		{
			desc: "existing section is replaced entirely",
			existing: new(`[test]
host  = https://old.cloud.databricks.com
token = old-token
`),
			profile: &Profile{Name: "test", Host: "https://new.cloud.databricks.com"},
			want: `[test]
host = https://new.cloud.databricks.com
`,
		},
		{
			desc:    "empty name returns an error",
			profile: &Profile{Host: "https://test.cloud.databricks.com"},
			wantErr: ErrEmptyProfile,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "databrickscfg")
			if tc.existing != nil {
				if err := os.WriteFile(path, []byte(*tc.existing), 0600); err != nil {
					t.Fatalf("WriteFile() error: %v", err)
				}
			}

			gotErr := tc.profile.SaveToFile(path)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("SaveToFile() error = %v, want %v", gotErr, tc.wantErr)
			}
			if tc.wantErr == nil {
				got, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("ReadFile() error: %v", err)
				}
				if diff := cmp.Diff(tc.want, string(got)); diff != "" {
					t.Errorf("file content mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestProfile_SaveToFile_filePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{Name: "test", Token: Secret("secret")}

	err := p.SaveToFile(path)

	if err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("file permissions = %o, want 0600", perm)
	}
}

func TestProfile_SaveToFile_emptyPath(t *testing.T) {
	p := &Profile{Name: "my-profile", Host: "https://test.cloud.databricks.com"}

	err := p.SaveToFile("")

	if !errors.Is(err, ErrEmptyPath) {
		t.Errorf("SaveToFile() error = %v, want %v", err, ErrEmptyPath)
	}
}

func TestProfile_SaveToFile_roundTrip(t *testing.T) {
	resetEnv(t)
	src := "testdata/databrickscfg"
	p, err := Resolve(WithFile(src), WithProfile("workspace"), WithoutEnv())
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}
	if p.Name != "workspace" {
		t.Fatalf("Name = %q, want %q", p.Name, "workspace")
	}

	dst := filepath.Join(t.TempDir(), "databrickscfg")
	if err := p.SaveToFile(dst); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}

	got, err := Resolve(WithFile(dst), WithProfile("workspace"), WithoutEnv())
	if err != nil {
		t.Fatalf("re-Resolve() error: %v", err)
	}
	if diff := cmp.Diff(p, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestProfile_SaveToFile_settingsSectionRejected(t *testing.T) {
	p := &Profile{Name: "__settings__", Host: "https://test.cloud.databricks.com"}

	err := p.SaveToFile(filepath.Join(t.TempDir(), "databrickscfg"))

	if !errors.Is(err, ErrInvalidProfileName) {
		t.Errorf("SaveToFile() error = %v, want %v", err, ErrInvalidProfileName)
	}
}

func TestDefaultConfigFile_fromEnv(t *testing.T) {
	t.Setenv("DATABRICKS_CONFIG_FILE", "/custom/path/databrickscfg")
	want := "/custom/path/databrickscfg"

	got := DefaultConfigFile()

	if got != want {
		t.Errorf("DefaultConfigFile() = %q, want %q", got, want)
	}
}

func TestDefaultConfigFile_fallsBackToHome(t *testing.T) {
	t.Setenv("DATABRICKS_CONFIG_FILE", "")
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("UserHomeDir() error: %v", err)
	}
	want := filepath.Join(home, ".databrickscfg")

	got := DefaultConfigFile()

	if got != want {
		t.Errorf("DefaultConfigFile() = %q, want %q", got, want)
	}
}

func TestListProfiles(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []string
		wantErr error
	}{
		{
			name: "allProfiles",
			path: "testdata/databrickscfg",
			want: []string{"DEFAULT", "workspace", "azure", "hash-in-value", "extra-keys", "empty"},
		},
		{
			name: "noDefault",
			path: "testdata/databrickscfg_no_default",
			want: []string{"workspace"},
		},
		{
			name: "settingsSectionExcluded",
			path: "testdata/databrickscfg_settings",
			want: []string{"DEFAULT", "my-workspace"},
		},
		{
			name:    "missingFile",
			path:    "testdata/nonexistent",
			wantErr: ErrConfigFileNotFound,
		},
		{
			name:    "emptyPath",
			path:    "",
			wantErr: ErrConfigFileNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, gotErr := ListProfiles(tc.path)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("ListProfiles() error = %v, want %v", gotErr, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ListProfiles() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
