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
			want: &Profile{},
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
				Host: "https://extra.cloud.databricks.com",
				Extra: map[string]string{
					"custom_key":  "custom-value",
					"another_key": "another-value",
				},
			},
		},
		{
			name: "withFileEmptyStringIsNoOp",
			opts: []ResolveOption{WithFile(""), WithoutEnv()},
			env: map[string]string{
				"DATABRICKS_CONFIG_FILE": "testdata/databrickscfg",
			},
			want: &Profile{
				Host:  "https://default.cloud.databricks.com",
				Token: Secret("default-token"),
			},
		},
		{
			name: "withProfileEmptyStringIsNoOp",
			opts: []ResolveOption{
				WithFile("testdata/databrickscfg"),
				WithProfile(""),
				WithoutEnv(),
			},
			env: map[string]string{
				"DATABRICKS_CONFIG_PROFILE": "workspace",
			},
			want: &Profile{
				Host:         "https://workspace.cloud.databricks.com",
				Token:        Secret("workspace-token"),
				AccountID:    "acc-123",
				ClientID:     "client-abc",
				ClientSecret: Secret("secret-xyz"),
				ClusterID:    "0123-456789-abcdef",
				WarehouseID:  "abc123def456",
			},
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
		if f.Name == "Extra" {
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
		section  string
		want     string
	}{
		{
			desc:    "known fields are written to the ini file",
			profile: &Profile{Host: "https://saved.cloud.databricks.com", Token: Secret("saved-token"), ClientID: "saved-client-id"},
			section: "test",
			want: "[test]\n" +
				"host      = https://saved.cloud.databricks.com\n" +
				"token     = saved-token\n" +
				"client_id = saved-client-id\n",
		},
		{
			desc:    "extra keys are written to the ini file",
			profile: &Profile{Host: "https://extra.cloud.databricks.com", Extra: map[string]string{"custom_key": "custom-value", "another_key": "another-value"}},
			section: "test",
			want: "[test]\n" +
				"host        = https://extra.cloud.databricks.com\n" +
				"another_key = another-value\n" +
				"custom_key  = custom-value\n",
		},
		{
			desc:    "extra keys that collide with known fields are skipped",
			profile: &Profile{Host: "https://real.cloud.databricks.com", Extra: map[string]string{"host": "https://evil.cloud.databricks.com"}},
			section: "test",
			want: "[test]\n" +
				"host = https://real.cloud.databricks.com\n",
		},
		{
			desc:    "empty fields are omitted",
			profile: &Profile{Host: "https://host.cloud.databricks.com"},
			section: "my-profile",
			want: "[my-profile]\n" +
				"host = https://host.cloud.databricks.com\n",
		},
		{
			desc: "other sections in the file are preserved",
			existing: new("[other]\n" +
				"host  = https://other.cloud.databricks.com\n" +
				"token = other-token\n"),
			profile: &Profile{Host: "https://new.cloud.databricks.com"},
			section: "new",
			want: "[other]\n" +
				"host  = https://other.cloud.databricks.com\n" +
				"token = other-token\n" +
				"\n" +
				"[new]\n" +
				"host = https://new.cloud.databricks.com\n",
		},
		{
			desc: "existing section is replaced entirely",
			existing: new("[test]\n" +
				"host  = https://old.cloud.databricks.com\n" +
				"token = old-token\n"),
			profile: &Profile{Host: "https://new.cloud.databricks.com"},
			section: "test",
			want: "[test]\n" +
				"host = https://new.cloud.databricks.com\n",
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
			if err := tc.profile.SaveToFile(path, tc.section); err != nil {
				t.Fatalf("SaveToFile() error: %v", err)
			}

			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("ReadFile() error: %v", err)
			}
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("file content mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestProfile_SaveToFile_filePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{Token: Secret("secret")}

	err := p.SaveToFile(path, "test")

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

func TestProfile_SaveToFile_emptyPathErrors(t *testing.T) {
	p := &Profile{Host: "https://test.cloud.databricks.com"}

	err := p.SaveToFile("", "my-profile")

	if err == nil {
		t.Error("SaveToFile() error = nil, want error for empty path")
	}
}

func TestProfile_SaveToFile_emptyProfileErrors(t *testing.T) {
	p := &Profile{Host: "https://test.cloud.databricks.com"}

	err := p.SaveToFile("/tmp/test", "")

	if err == nil {
		t.Error("SaveToFile() error = nil, want error for empty profile")
	}
}

// TestProfile_SaveToFile_allProperties verifies that every property in the
// mapping table round-trips correctly through SaveToFile and Resolve.
func TestProfile_SaveToFile_allProperties(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{}
	for i, prop := range properties {
		var val string
		if prop.iniKey == "azure_use_msi" {
			val = "true"
		} else {
			val = strings.Repeat("s", i+1)
		}
		if err := prop.set(p, val); err != nil {
			t.Fatalf("set %q: %v", prop.iniKey, err)
		}
	}

	if err := p.SaveToFile(path, "all-fields"); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("all-fields"), WithoutEnv())
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}

	for i, prop := range properties {
		var want string
		if prop.iniKey == "azure_use_msi" {
			want = "true"
		} else {
			want = strings.Repeat("s", i+1)
		}
		gotVal := prop.get(got)
		if gotVal != want {
			t.Errorf("round-trip: property %q (iniKey=%q) = %q, want %q", prop.envVar, prop.iniKey, gotVal, want)
		}
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
