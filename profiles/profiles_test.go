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

// clearProfileEnv clears all environment variables that could affect profile
// resolution. Call this at the top of tests that don't use WithoutEnv() to
// prevent ambient env vars from causing flaky failures.
func clearProfileEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABRICKS_CONFIG_FILE", "")
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "")
	for _, prop := range properties {
		t.Setenv(prop.envVar, "")
	}
}

func TestResolve_withFileAndProfile(t *testing.T) {
	want := &Profile{
		Host:         "https://workspace.cloud.databricks.com",
		Token:        "workspace-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("workspace"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_defaultSection(t *testing.T) {
	want := &Profile{
		Host:  "https://default.cloud.databricks.com",
		Token: "default-token",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_missingExplicitFile(t *testing.T) {
	got, err := Resolve(WithFile("testdata/nonexistent"), WithoutEnv())

	if !errors.Is(err, ErrConfigFileNotFound) {
		t.Errorf("Resolve() error = %v, want ErrConfigFileNotFound", err)
	}
	if got != nil {
		t.Errorf("Resolve() = %v, want nil", got)
	}
}

func TestResolve_missingEnvConfigFile(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_CONFIG_FILE", "testdata/nonexistent")

	got, err := Resolve()

	if !errors.Is(err, ErrConfigFileNotFound) {
		t.Errorf("Resolve() error = %v, want ErrConfigFileNotFound", err)
	}
	if got != nil {
		t.Errorf("Resolve() = %v, want nil", got)
	}
}

func TestResolve_missingDefaultFileSkipped(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("DATABRICKS_HOST", "https://env.cloud.databricks.com")
	want := &Profile{Host: "https://env.cloud.databricks.com"}

	got, err := Resolve()

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_missingExplicitProfile(t *testing.T) {
	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("nonexistent"),
		WithoutEnv(),
	)

	if !errors.Is(err, ErrProfileNotFound) {
		t.Errorf("Resolve() error = %v, want ErrProfileNotFound", err)
	}
	if got != nil {
		t.Errorf("Resolve() = %v, want nil", got)
	}
}

func TestResolve_missingDefaultSection(t *testing.T) {
	want := &Profile{}

	got, err := Resolve(
		WithFile("testdata/databrickscfg_no_default"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_hashInValues(t *testing.T) {
	want := &Profile{
		Host:         "https://hash.cloud.databricks.com",
		Token:        "abc#def#ghi",
		ClientSecret: "secret#with#hashes",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("hash-in-value"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_azure(t *testing.T) {
	want := &Profile{
		Host:              "https://adb-123.azuredatabricks.net",
		AzureClientID:     "az-client-id",
		AzureClientSecret: "az-client-secret",
		AzureTenantID:     "az-tenant-id",
		AzureResourceID:   "/subscriptions/sub-id/resourceGroups/rg/providers/Microsoft.Databricks/workspaces/ws",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("azure"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_emptySection(t *testing.T) {
	want := &Profile{}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("empty"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_envOnly(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_HOST", "https://env.cloud.databricks.com")
	t.Setenv("DATABRICKS_TOKEN", "env-token")
	t.Setenv("DATABRICKS_CLIENT_ID", "env-client-id")
	want := &Profile{
		Host:     "https://env.cloud.databricks.com",
		Token:    "env-token",
		ClientID: "env-client-id",
	}

	got, err := Resolve(WithFile("testdata/databrickscfg_no_default"))

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_envOverridesFile(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_HOST", "https://env-override.cloud.databricks.com")
	t.Setenv("DATABRICKS_TOKEN", "env-override-token")
	want := &Profile{
		Host:         "https://env-override.cloud.databricks.com",
		Token:        "env-override-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("workspace"),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_noOptions(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_CONFIG_FILE", "testdata/databrickscfg")
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "workspace")
	want := &Profile{
		Host:         "https://workspace.cloud.databricks.com",
		Token:        "workspace-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve()

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_noOptionsEnvOverridesFile(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_CONFIG_FILE", "testdata/databrickscfg")
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "workspace")
	t.Setenv("DATABRICKS_HOST", "https://override.cloud.databricks.com")
	want := &Profile{
		Host:         "https://override.cloud.databricks.com",
		Token:        "workspace-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve()

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_withoutEnv(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_HOST", "https://should-be-ignored.cloud.databricks.com")
	want := &Profile{
		Host:         "https://workspace.cloud.databricks.com",
		Token:        "workspace-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("workspace"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestResolve_extraKeys(t *testing.T) {
	want := &Profile{
		Host: "https://extra.cloud.databricks.com",
		Extra: map[string]string{
			"custom_key":  "custom-value",
			"another_key": "another-value",
		},
	}

	got, err := Resolve(
		WithFile("testdata/databrickscfg"),
		WithProfile("extra-keys"),
		WithoutEnv(),
	)

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
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
		fileVal := strings.Repeat("f", i+1) // "f", "ff", "fff", ...
		envVal := strings.Repeat("e", i+1)  // "e", "ee", "eee", ...
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
		want := strings.Repeat("f", i+1)
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
		want := strings.Repeat("e", i+1)
		got := prop.get(envProfile)
		if got != want {
			t.Errorf("env: property %q (envVar=%q) = %q, want %q", prop.iniKey, prop.envVar, got, want)
		}
	}
}

// TestResolve_allPropertiesCovered verifies that every exported string field on
// Profile (except Extra) has a corresponding entry in the properties table.
func TestResolve_allPropertiesCovered(t *testing.T) {
	typ := reflect.TypeOf(Profile{})
	covered := make(map[string]bool)
	for _, prop := range properties {
		p := &Profile{}
		prop.set(p, "sentinel")
		v := reflect.ValueOf(*p)
		for i := range v.NumField() {
			f := v.Field(i)
			if f.Kind() == reflect.String && f.String() == "sentinel" {
				covered[typ.Field(i).Name] = true
			}
		}
	}

	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.Type.Kind() != reflect.String {
			continue // skip Extra (map)
		}
		if !covered[f.Name] {
			t.Errorf("Profile.%s has no entry in properties table", f.Name)
		}
	}
}

// TestProfile_String_redactsAllSensitiveFields verifies that every property
// marked as sensitive is redacted in String() output.
func TestProfile_String_redactsAllSensitiveFields(t *testing.T) {
	for _, prop := range properties {
		if !prop.sensitive {
			continue
		}
		p := &Profile{}
		sentinel := "LEAKED-" + prop.iniKey
		prop.set(p, sentinel)

		got := p.String()

		if strings.Contains(got, sentinel) {
			t.Errorf("String() leaked sensitive field %q (iniKey=%q)", prop.envVar, prop.iniKey)
		}
		if !strings.Contains(got, prop.iniKey+"=********") {
			t.Errorf("String() should show redacted %q, got:\n%s", prop.iniKey, got)
		}
	}
}

func TestProfile_String_showsNonSensitiveFields(t *testing.T) {
	p := &Profile{
		Host:     "https://test.cloud.databricks.com",
		ClientID: "visible-client-id",
	}

	got := p.String()

	if !strings.Contains(got, "host=https://test.cloud.databricks.com") {
		t.Errorf("String() should contain host value, got:\n%s", got)
	}
	if !strings.Contains(got, "client_id=visible-client-id") {
		t.Errorf("String() should contain client_id value, got:\n%s", got)
	}
}

func TestProfile_String_omitsExtra(t *testing.T) {
	p := &Profile{
		Host: "https://test.cloud.databricks.com",
		Extra: map[string]string{
			"secret_key": "secret_value",
		},
	}

	got := p.String()

	if strings.Contains(got, "secret_key") {
		t.Errorf("String() should omit Extra keys, got:\n%s", got)
	}
}

func TestProfile_String_emptyProfile(t *testing.T) {
	p := &Profile{}

	got := p.String()

	if got != "" {
		t.Errorf("String() = %q, want empty", got)
	}
}

func TestProfile_String_ordering(t *testing.T) {
	p := &Profile{
		Host:      "https://test.cloud.databricks.com",
		AccountID: "acc-123",
		Token:     "secret-token",
		ClientID:  "client-abc",
		ClusterID: "cluster-xyz",
	}
	want := "host=https://test.cloud.databricks.com\n" +
		"account_id=acc-123\n" +
		"token=********\n" +
		"client_id=client-abc\n" +
		"cluster_id=cluster-xyz\n"

	got := p.String()

	if got != want {
		t.Errorf("String() =\n%s\nwant:\n%s", got, want)
	}
}

func TestProfile_SaveToFile_newFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	want := &Profile{
		Host:     "https://saved.cloud.databricks.com",
		Token:    "saved-token",
		ClientID: "saved-client-id",
	}

	err := want.SaveToFile(path, "my-profile")

	if err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("my-profile"), WithoutEnv())
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestProfile_SaveToFile_filePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{Token: "secret"}

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

func TestProfile_SaveToFile_preservesOtherSections(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p1 := &Profile{Host: "https://one.cloud.databricks.com", Token: "token-one"}
	p2 := &Profile{Host: "https://two.cloud.databricks.com", Token: "token-two"}
	if err := p1.SaveToFile(path, "one"); err != nil {
		t.Fatalf("SaveToFile(one) error: %v", err)
	}
	if err := p2.SaveToFile(path, "two"); err != nil {
		t.Fatalf("SaveToFile(two) error: %v", err)
	}

	got1, err := Resolve(WithFile(path), WithProfile("one"), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve(one) error: %v", err)
	}
	if diff := cmp.Diff(p1, got1); diff != "" {
		t.Errorf("profile one mismatch (-want +got):\n%s", diff)
	}

	got2, err := Resolve(WithFile(path), WithProfile("two"), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve(two) error: %v", err)
	}
	if diff := cmp.Diff(p2, got2); diff != "" {
		t.Errorf("profile two mismatch (-want +got):\n%s", diff)
	}
}

func TestProfile_SaveToFile_overwritesExistingSection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p1 := &Profile{Host: "https://old.cloud.databricks.com", Token: "old-token"}
	if err := p1.SaveToFile(path, "test"); err != nil {
		t.Fatalf("SaveToFile() first: %v", err)
	}
	want := &Profile{Host: "https://new.cloud.databricks.com"}

	if err := want.SaveToFile(path, "test"); err != nil {
		t.Fatalf("SaveToFile() second: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("test"), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestProfile_SaveToFile_roundTripPreservesExtra(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	want := &Profile{
		Host:  "https://extra.cloud.databricks.com",
		Token: "extra-token",
		Extra: map[string]string{
			"custom_key":  "custom-value",
			"another_key": "another-value",
		},
	}

	if err := want.SaveToFile(path, "extras"); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("extras"), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestProfile_SaveToFile_extraKeyCollisionSkipped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{
		Host: "https://real.cloud.databricks.com",
		Extra: map[string]string{
			"host": "https://evil.cloud.databricks.com",
		},
	}
	want := "https://real.cloud.databricks.com"

	if err := p.SaveToFile(path, "test"); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("test"), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if got.Host != want {
		t.Errorf("Resolve().Host = %q, want %q", got.Host, want)
	}
}

// TestProfile_SaveToFile_allProperties verifies that every property in the
// mapping table round-trips correctly through SaveToFile and Resolve.
func TestProfile_SaveToFile_allProperties(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	p := &Profile{}
	for i, prop := range properties {
		prop.set(p, strings.Repeat("s", i+1))
	}

	if err := p.SaveToFile(path, "all-fields"); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithProfile("all-fields"), WithoutEnv())
	if err != nil {
		t.Fatalf("Resolve() error: %v", err)
	}

	for i, prop := range properties {
		want := strings.Repeat("s", i+1)
		gotVal := prop.get(got)
		if gotVal != want {
			t.Errorf("round-trip: property %q (iniKey=%q) = %q, want %q", prop.envVar, prop.iniKey, gotVal, want)
		}
	}
}

func TestProfile_SaveToFile_defaultProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "databrickscfg")
	want := &Profile{
		Host:  "https://default.cloud.databricks.com",
		Token: "default-token",
	}

	if err := want.SaveToFile(path, "DEFAULT"); err != nil {
		t.Fatalf("SaveToFile() error: %v", err)
	}
	got, err := Resolve(WithFile(path), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
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

func TestWithFile_emptyStringIsNoOp(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_CONFIG_FILE", "testdata/databrickscfg")
	want := &Profile{
		Host:  "https://default.cloud.databricks.com",
		Token: "default-token",
	}

	got, err := Resolve(WithFile(""), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestWithProfile_emptyStringIsNoOp(t *testing.T) {
	clearProfileEnv(t)
	t.Setenv("DATABRICKS_CONFIG_PROFILE", "workspace")
	want := &Profile{
		Host:         "https://workspace.cloud.databricks.com",
		Token:        "workspace-token",
		AccountID:    "acc-123",
		ClientID:     "client-abc",
		ClientSecret: "secret-xyz",
		ClusterID:    "0123-456789-abcdef",
		WarehouseID:  "abc123def456",
	}

	got, err := Resolve(WithFile("testdata/databrickscfg"), WithProfile(""), WithoutEnv())

	if err != nil {
		t.Errorf("Resolve() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
	}
}

func TestListProfiles(t *testing.T) {
	want := []string{"DEFAULT", "workspace", "azure", "hash-in-value", "extra-keys", "empty"}

	got, err := ListProfiles("testdata/databrickscfg")

	if err != nil {
		t.Errorf("ListProfiles() error: %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ListProfiles() mismatch (-want +got):\n%s", diff)
	}
}

func TestListProfiles_noDefault(t *testing.T) {
	got, err := ListProfiles("testdata/databrickscfg_no_default")

	if err != nil {
		t.Errorf("ListProfiles() error: %v", err)
	}
	for _, name := range got {
		if name == "DEFAULT" {
			t.Error("ListProfiles() includes DEFAULT section with no keys")
		}
	}
}

func TestListProfiles_missingFile(t *testing.T) {
	got, err := ListProfiles("testdata/nonexistent")

	if err == nil {
		t.Error("ListProfiles() error = nil, want error for missing file")
	}
	if got != nil {
		t.Errorf("ListProfiles() = %v, want nil", got)
	}
}

func TestListProfiles_emptyPath(t *testing.T) {
	got, err := ListProfiles("")

	if err == nil {
		t.Error("ListProfiles() error = nil, want error for empty path")
	}
	if got != nil {
		t.Errorf("ListProfiles() = %v, want nil", got)
	}
}
