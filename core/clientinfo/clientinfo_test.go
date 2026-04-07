package clientinfo

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/databricks/sdk-go/core/internal"
)

func TestClientInfo_WithAndString(t *testing.T) {
	testCases := []struct {
		desc       string
		base       ClientInfo
		kvs        []string
		wantString string
		wantErr    error
	}{
		{
			desc:       "empty base empty args",
			wantString: "",
		},
		{
			desc:       "single pair on empty base",
			base:       ClientInfo{},
			kvs:        []string{"auth", "pat"},
			wantString: "auth/pat",
		},
		{
			desc:       "multiple pairs on empty base",
			base:       ClientInfo{},
			kvs:        []string{"dataquality", "0.0.1", "auth", "pat"},
			wantString: "dataquality/0.0.1 auth/pat",
		},
		{
			desc:       "appends to existing segments",
			base:       ClientInfo{segments: []segment{{"sdk", "1.0.0"}}},
			kvs:        []string{"auth", "pat"},
			wantString: "sdk/1.0.0 auth/pat",
		},
		{
			desc:       "no args returns same value",
			base:       ClientInfo{segments: []segment{{"sdk", "1.0.0"}}},
			wantString: "sdk/1.0.0",
		},
		{
			desc:       "preserves insertion order",
			base:       ClientInfo{},
			kvs:        []string{"zzz", "3", "aaa", "1", "mmm", "2"},
			wantString: "zzz/3 aaa/1 mmm/2",
		},
		{
			desc:       "exact duplicate silently ignored",
			base:       ClientInfo{segments: []segment{{"key", "value"}}},
			kvs:        []string{"key", "value"},
			wantString: "key/value",
		},
		{
			desc:       "duplicate within batch silently ignored",
			base:       ClientInfo{},
			kvs:        []string{"key", "value", "key", "value"},
			wantString: "key/value",
		},
		{
			desc:       "same key different value allowed",
			base:       ClientInfo{segments: []segment{{"partner", "acme"}}},
			kvs:        []string{"partner", "contoso"},
			wantString: "partner/acme partner/contoso",
		},
		{
			desc:    "odd number of arguments",
			base:    ClientInfo{},
			kvs:     []string{"key"},
			wantErr: ErrOddKeyvals,
		},
		{
			desc:    "invalid key with space",
			base:    ClientInfo{},
			kvs:     []string{"bad key", "value"},
			wantErr: ErrInvalidKey,
		},
		{
			desc:    "invalid key with slash",
			base:    ClientInfo{},
			kvs:     []string{"bad/key", "value"},
			wantErr: ErrInvalidKey,
		},
		{
			desc:    "invalid value with space",
			base:    ClientInfo{},
			kvs:     []string{"key", "bad value"},
			wantErr: ErrInvalidValue,
		},
		{
			desc:    "invalid value with special chars",
			base:    ClientInfo{},
			kvs:     []string{"key", "bad!value"},
			wantErr: ErrInvalidValue,
		},
		{
			desc:    "error on first invalid pair returns zero value",
			base:    ClientInfo{segments: []segment{{"existing", "value"}}},
			kvs:     []string{"bad key", "value"},
			wantErr: ErrInvalidKey,
		},
		{
			desc:    "error on second pair leaves base unchanged",
			base:    ClientInfo{},
			kvs:     []string{"good", "value", "bad key", "value"},
			wantErr: ErrInvalidKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, gotErr := tc.base.With(tc.kvs...)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("With() error = %v, want %v", gotErr, tc.wantErr)
			}
			if gotString := got.String(); gotString != tc.wantString {
				t.Errorf("With() = %q, want %q", gotString, tc.wantString)
			}
		})
	}
}

// TestClientInfo_With_noMutation verifies that With does not mutate the
// original ClientInfo.
func TestClientInfo_With_noMutation(t *testing.T) {
	base, _ := ClientInfo{}.With("foo", "bar")
	want := base.String()

	_, _ = base.With("extra", "value")

	if got := base.String(); got != want {
		t.Errorf("original was mutated: got %q, want %q", got, want)
	}
}

func resetBase(t *testing.T) {
	t.Helper()
	orig := base
	t.Cleanup(func() { base = orig })
}

// mockEnv returns a lookupFunc backed by a map. Keys not in the map
// are treated as unset.
func mockEnv(env map[string]string) lookupFunc {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

// TestDefault_passLookupFunc focuses on testing the plumbing between Default
// and defaultWithEnv. It verifies that Default passes os.LookupEnv.
func TestDefault_passLookupFunc(t *testing.T) {
	t.Setenv("DATABRICKS_SDK_UPSTREAM", "test-foo")
	t.Setenv("DATABRICKS_SDK_UPSTREAM_VERSION", "42.13.37")
	want := "upstream/test-foo upstream-version/42.13.37"

	got := Default().String()

	if !strings.Contains(got, want) {
		t.Errorf("Default() = %q, want %q", got, want)
	}
}

func TestDefault(t *testing.T) {
	prefix := internal.ModuleName + "/" + internal.Version + " go/" + cachedGoVersion + " os/" + runtime.GOOS

	testCases := []struct {
		desc string
		env  map[string]string
		want string
	}{
		{
			desc: "base segments only",
			want: prefix,
		},
		{
			desc: "github actions",
			env:  map[string]string{"GITHUB_ACTIONS": "true"},
			want: prefix + " cicd/github",
		},
		{
			desc: "gitlab ci",
			env:  map[string]string{"GITLAB_CI": "true"},
			want: prefix + " cicd/gitlab",
		},
		{
			desc: "google cloud build requires all four vars",
			env: map[string]string{
				"PROJECT_ID":     "my-project",
				"BUILD_ID":       "123",
				"PROJECT_NUMBER": "456",
				"LOCATION":       "us-central1",
			},
			want: prefix + " cicd/google-cloud-build",
		},
		{
			desc: "single agent",
			env:  map[string]string{"CLAUDECODE": "1"},
			want: prefix + " agent/claude-code",
		},
		{
			desc: "multiple agents omitted",
			env:  map[string]string{"CLAUDECODE": "1", "CURSOR_AGENT": "1"},
			want: prefix,
		},
		{
			desc: "databricks runtime",
			env:  map[string]string{"DATABRICKS_RUNTIME_VERSION": "15.5"},
			want: prefix + " runtime/15.5",
		},
		{
			desc: "databricks runtime sanitized",
			env:  map[string]string{"DATABRICKS_RUNTIME_VERSION": "15.5 beta/2"},
			want: prefix + " runtime/15.5-beta-2",
		},
		{
			desc: "upstream both vars",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0",
			},
			want: prefix + " upstream/terraform upstream-version/1.5.0",
		},
		{
			desc: "upstream omitted when only product set",
			env:  map[string]string{"DATABRICKS_SDK_UPSTREAM": "terraform"},
			want: prefix,
		},
		{
			desc: "all env detection combined",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0",
				"GITHUB_ACTIONS":                  "true",
				"DATABRICKS_RUNTIME_VERSION":      "15.5",
				"CLAUDECODE":                      "1",
			},
			want: prefix + " upstream/terraform upstream-version/1.5.0 cicd/github runtime/15.5 agent/claude-code",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resetBase(t)

			got := defaultWithEnv(mockEnv(tc.env)).String()

			if got != tc.want {
				t.Errorf("defaultWithEnv() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSetProduct(t *testing.T) {
	testCases := []struct {
		desc       string
		name       string
		version    string
		wantErr    error
		wantInBase string
	}{
		{
			desc:       "valid product",
			name:       "my-app",
			version:    "1.2.3",
			wantInBase: "my-app/1.2.3",
		},
		{
			desc:    "invalid name with space",
			name:    "invalid name",
			version: "1.0.0",
			wantErr: ErrInvalidKey,
		},
		{
			desc:    "invalid version not semver",
			name:    "valid-name",
			version: "not-semver",
			wantErr: ErrInvalidVersion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resetBase(t)

			gotErr := SetProduct(tc.name, tc.version)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("SetProduct() error = %v, want %v", gotErr, tc.wantErr)
			}
			if got := base.String(); got != tc.wantInBase {
				t.Errorf("base = %q, want %q", got, tc.wantInBase)
			}
		})
	}
}

func TestSetPartner(t *testing.T) {
	testCases := []struct {
		desc       string
		partner    string
		wantErr    error
		wantInBase string
	}{
		{
			desc:       "valid partner",
			partner:    "acme",
			wantInBase: "partner/acme",
		},
		{
			desc:    "invalid partner with space",
			partner: "bad partner",
			wantErr: ErrInvalidValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resetBase(t)

			gotErr := SetPartner(tc.partner)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("SetPartner() error = %v, want %v", gotErr, tc.wantErr)
			}
			if got := base.String(); got != tc.wantInBase {
				t.Errorf("base = %q, want %q", got, tc.wantInBase)
			}
		})
	}
}

func TestAddToDefault(t *testing.T) {
	testCases := []struct {
		desc       string
		setup      []segment
		key        string
		value      string
		wantErr    error
		wantInBase string
	}{
		{
			desc:       "valid pair",
			key:        "test-key",
			value:      "test-value",
			wantInBase: "test-key/test-value",
		},
		{
			desc:    "invalid key with space",
			key:     "bad key",
			value:   "value",
			wantErr: ErrInvalidKey,
		},
		{
			desc:       "exact duplicate silently ignored",
			setup:      []segment{{"dup", "value"}},
			key:        "dup",
			value:      "value",
			wantInBase: "dup/value",
		},
		{
			desc:       "same key different value allowed",
			setup:      []segment{{"partner", "acme"}},
			key:        "partner",
			value:      "contoso",
			wantInBase: "partner/acme partner/contoso",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			resetBase(t)
			for _, kv := range tc.setup {
				if err := AddToDefault(kv.key, kv.value); err != nil {
					t.Fatalf("setup AddToDefault(%q, %q) error: %v", kv.key, kv.value, err)
				}
			}

			gotErr := AddToDefault(tc.key, tc.value)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("AddToDefault() error = %v, want %v", gotErr, tc.wantErr)
			}
			if got := base.String(); got != tc.wantInBase {
				t.Errorf("base = %q, want %q", got, tc.wantInBase)
			}
		})
	}
}

func TestIsSemVer(t *testing.T) {
	testCases := []struct {
		input string
		want  bool
	}{
		{"1.2.3", true},
		{"0.0.0-dev+2e014739024a", true},
		{"1.2.3.4", false},
		{"1.2", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := isSemVer(tc.input); got != tc.want {
				t.Errorf("isSemVer(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsValidSegment(t *testing.T) {
	testCases := []struct {
		input string
		want  bool
	}{
		{"foo", true},
		{"FOO", true},
		{"FOO123", true},
		{"foo_bar", true},
		{"foo-bar", true},
		{"foo.bar", true},
		{"foo+bar", true},
		{"foo bar", false},
		{"foo/bar", false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := isValidSegment(tc.input); got != tc.want {
				t.Errorf("isValidSegment(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"foo", "foo"},
		{"FOO", "FOO"},
		{"foo_bar", "foo_bar"},
		{"foo-bar", "foo-bar"},
		{"foo+bar", "foo+bar"},
		{"foo.bar", "foo.bar"},
		{"1.2.3", "1.2.3"},
		{"foo bar", "foo-bar"},
		{"foo/bar", "foo-bar"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := sanitize(tc.input); got != tc.want {
				t.Errorf("sanitize(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestNormalizeGoVersion(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"go1.26.0", "1.26.0"},
		{"go1.26", "1.26.0"},
		{"go1.26rc1", "1.26.0-rc1"},
		{"go1.26.0rc1", "1.26.0-rc1"},
		{"go1.26beta2", "1.26.0-beta2"},
		{"go2", "2.0.0"},
		{"devel +abc123def Mon Jan 1 00:00:00 2024 +0000", "0.0.0-dev"},
		{"devel go1.23-abc123", "0.0.0-dev"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			if got := normalizeGoVersion(tc.input); got != tc.want {
				t.Errorf("normalizeGoVersion(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
