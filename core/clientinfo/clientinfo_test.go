package clientinfo

import (
	"errors"
	"os"
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
			desc: "multiple agents report the multiple sentinel",
			env:  map[string]string{"CLAUDECODE": "1", "CURSOR_AGENT": "1"},
			want: prefix + " agent/multiple",
		},
		{
			desc: "AGENT fallback",
			env:  map[string]string{"AGENT": "goose"},
			want: prefix + " agent/goose",
		},
		{
			desc: "AI_AGENT fallback",
			env:  map[string]string{"AI_AGENT": "cursor"},
			want: prefix + " agent/cursor",
		},
		{
			desc: "omnigent meta-harness",
			env:  map[string]string{"OMNIGENT": "1"},
			want: prefix + " meta-harness/omnigent",
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
			desc: "upstream omitted when product is empty",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0",
			},
			want: prefix,
		},
		{
			desc: "upstream omitted when version is empty",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "",
			},
			want: prefix,
		},
		{
			desc: "upstream sanitized when rendered",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform provider",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0/dev",
			},
			want: prefix + " upstream/terraform-provider upstream-version/1.5.0-dev",
		},
		{
			desc: "all env detection combined",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0",
				"GITHUB_ACTIONS":                  "true",
				"DATABRICKS_RUNTIME_VERSION":      "15.5",
				"CLAUDECODE":                      "1",
				"OMNIGENT":                        "1",
			},
			want: prefix + " upstream/terraform upstream-version/1.5.0 cicd/github runtime/15.5 agent/claude-code meta-harness/omnigent",
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

func TestDetectUpstream(t *testing.T) {
	testCases := []struct {
		desc        string
		env         map[string]string
		wantProduct string
		wantVersion string
	}{
		{desc: "unset"},
		{desc: "only product", env: map[string]string{"DATABRICKS_SDK_UPSTREAM": "terraform"}},
		{desc: "only version", env: map[string]string{"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0"}},
		{desc: "empty product", env: map[string]string{"DATABRICKS_SDK_UPSTREAM": "", "DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0"}},
		{desc: "empty version", env: map[string]string{"DATABRICKS_SDK_UPSTREAM": "terraform", "DATABRICKS_SDK_UPSTREAM_VERSION": ""}},
		{
			desc: "valid values unchanged",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform-provider",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0-dev+build.1",
			},
			wantProduct: "terraform-provider",
			wantVersion: "1.5.0-dev+build.1",
		},
		{
			desc: "malformed values sanitized",
			env: map[string]string{
				"DATABRICKS_SDK_UPSTREAM":         "terraform provider/beta",
				"DATABRICKS_SDK_UPSTREAM_VERSION": "1.5.0/dev\r\nnext",
			},
			wantProduct: "terraform-provider-beta",
			wantVersion: "1.5.0-dev--next",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotProduct, gotVersion := detectUpstream(mockEnv(tc.env))
			if gotProduct != tc.wantProduct || gotVersion != tc.wantVersion {
				t.Errorf("detectUpstream() = (%q, %q), want (%q, %q)", gotProduct, gotVersion, tc.wantProduct, tc.wantVersion)
			}
		})
	}
}

func TestDetectAgent(t *testing.T) {
	testCases := []struct {
		desc string
		env  map[string]string
		want string
	}{
		{desc: "no agent", want: ""},
		{desc: "amp", env: map[string]string{"AMP_CURRENT_THREAD_ID": "thread"}, want: "amp"},
		{desc: "antigravity", env: map[string]string{"ANTIGRAVITY_AGENT": "1"}, want: "antigravity"},
		{desc: "augment", env: map[string]string{"AUGMENT_AGENT": "1"}, want: "augment"},
		{desc: "claude code", env: map[string]string{"CLAUDECODE": "1"}, want: "claude-code"},
		{desc: "cline", env: map[string]string{"CLINE_ACTIVE": "1"}, want: "cline"},
		{desc: "codex", env: map[string]string{"CODEX_CI": "1"}, want: "codex"},
		{desc: "copilot CLI", env: map[string]string{"COPILOT_CLI": "1"}, want: "copilot-cli"},
		{desc: "cursor", env: map[string]string{"CURSOR_AGENT": "1"}, want: "cursor"},
		{desc: "gemini CLI", env: map[string]string{"GEMINI_CLI": "1"}, want: "gemini-cli"},
		{desc: "goose", env: map[string]string{"GOOSE_TERMINAL": "1"}, want: "goose"},
		{desc: "kiro", env: map[string]string{"KIRO": "1"}, want: "kiro"},
		{desc: "openclaw", env: map[string]string{"OPENCLAW_SHELL": "1"}, want: "openclaw"},
		{desc: "opencode", env: map[string]string{"OPENCODE": "1"}, want: "opencode"},
		{desc: "VS Code agent", env: map[string]string{"VSCODE_AGENT": "1"}, want: "vscode-agent"},
		{desc: "windsurf", env: map[string]string{"WINDSURF_AGENT": "1"}, want: "windsurf"},
		{desc: "empty explicit value counts", env: map[string]string{"CLAUDECODE": ""}, want: "claude-code"},
		{desc: "multiple explicit agents", env: map[string]string{"CLAUDECODE": "1", "CURSOR_AGENT": "1"}, want: "multiple"},
		{desc: "AGENT fallback", env: map[string]string{"AGENT": "goose"}, want: "goose"},
		{desc: "AGENT sanitized", env: map[string]string{"AGENT": "claude code/agent"}, want: "claude-code-agent"},
		{desc: "AGENT length capped", env: map[string]string{"AGENT": strings.Repeat("a", 100)}, want: strings.Repeat("a", 64)},
		{desc: "empty AGENT falls through", env: map[string]string{"AGENT": "", "AI_AGENT": "cursor"}, want: "cursor"},
		{desc: "AGENT wins over AI_AGENT", env: map[string]string{"AGENT": "claude-code", "AI_AGENT": "cursor"}, want: "claude-code"},
		{desc: "explicit matcher wins over fallback", env: map[string]string{"CLAUDECODE": "1", "AGENT": "goose"}, want: "claude-code"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := detectAgent(mockEnv(tc.env)); got != tc.want {
				t.Errorf("detectAgent() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDetectRuntime(t *testing.T) {
	testCases := []struct {
		desc string
		env  map[string]string
		want string
	}{
		{desc: "unset", want: ""},
		{desc: "empty", env: map[string]string{"DATABRICKS_RUNTIME_VERSION": ""}, want: ""},
		{desc: "version", env: map[string]string{"DATABRICKS_RUNTIME_VERSION": "15.5"}, want: "15.5"},
		{desc: "sanitized", env: map[string]string{"DATABRICKS_RUNTIME_VERSION": "15.5 beta/2"}, want: "15.5-beta-2"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := detectRuntimeVersion(mockEnv(tc.env)); got != tc.want {
				t.Errorf("detectRuntimeVersion() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDetectCICD(t *testing.T) {
	testCases := []struct {
		desc string
		env  map[string]string
		want string
	}{
		{desc: "unset", want: ""},
		{desc: "github", env: map[string]string{"GITHUB_ACTIONS": "true"}, want: "github"},
		{desc: "gitlab", env: map[string]string{"GITLAB_CI": "true"}, want: "gitlab"},
		{desc: "jenkins", env: map[string]string{"JENKINS_URL": ""}, want: "jenkins"},
		{desc: "azure devops", env: map[string]string{"TF_BUILD": "True"}, want: "azure-devops"},
		{desc: "circle", env: map[string]string{"CIRCLECI": "true"}, want: "circle"},
		{desc: "travis", env: map[string]string{"TRAVIS": "true"}, want: "travis"},
		{desc: "bitbucket", env: map[string]string{"BITBUCKET_BUILD_NUMBER": ""}, want: "bitbucket"},
		{desc: "google cloud build", env: map[string]string{"PROJECT_ID": "project", "BUILD_ID": "build", "PROJECT_NUMBER": "123", "LOCATION": "us-central1"}, want: "google-cloud-build"},
		{desc: "aws codebuild", env: map[string]string{"CODEBUILD_BUILD_ARN": ""}, want: "aws-code-build"},
		{desc: "terraform cloud", env: map[string]string{"TFC_RUN_ID": ""}, want: "tf-cloud"},
		{desc: "github exact value", env: map[string]string{"GITHUB_ACTIONS": "True"}, want: ""},
		{desc: "gitlab exact value", env: map[string]string{"GITLAB_CI": "1"}, want: ""},
		{desc: "azure devops exact value", env: map[string]string{"TF_BUILD": "true"}, want: ""},
		{desc: "circle exact value", env: map[string]string{"CIRCLECI": "1"}, want: ""},
		{desc: "travis exact value", env: map[string]string{"TRAVIS": "1"}, want: ""},
		{desc: "google cloud build missing project", env: map[string]string{"BUILD_ID": "build", "PROJECT_NUMBER": "123", "LOCATION": "us-central1"}, want: ""},
		{desc: "google cloud build missing build", env: map[string]string{"PROJECT_ID": "project", "PROJECT_NUMBER": "123", "LOCATION": "us-central1"}, want: ""},
		{desc: "google cloud build missing project number", env: map[string]string{"PROJECT_ID": "project", "BUILD_ID": "build", "LOCATION": "us-central1"}, want: ""},
		{desc: "google cloud build missing location", env: map[string]string{"PROJECT_ID": "project", "BUILD_ID": "build", "PROJECT_NUMBER": "123"}, want: ""},
		{desc: "first match wins", env: map[string]string{"GITHUB_ACTIONS": "true", "GITLAB_CI": "true"}, want: "github"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := detectCICD(mockEnv(tc.env)); got != tc.want {
				t.Errorf("detectCICD() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDetectMetaHarness(t *testing.T) {
	testCases := []struct {
		desc string
		env  map[string]string
		want string
	}{
		{desc: "unset", want: ""},
		{desc: "present", env: map[string]string{"OMNIGENT": "1"}, want: "omnigent"},
		{desc: "empty value counts", env: map[string]string{"OMNIGENT": ""}, want: "omnigent"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if got := detectMetaHarness(mockEnv(tc.env)); got != tc.want {
				t.Errorf("detectMetaHarness() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDetectMetaHarness_multiple(t *testing.T) {
	metaHarnesses := []environmentProductDef{
		{envVar: "FIRST_HARNESS", product: "first"},
		{envVar: "SECOND_HARNESS", product: "second"},
	}
	env := mockEnv(map[string]string{"FIRST_HARNESS": "1", "SECOND_HARNESS": "1"})

	if got, want := detectEnvironmentProduct(env, metaHarnesses), "multiple"; got != want {
		t.Errorf("detectEnvironmentProduct() = %q, want %q", got, want)
	}
}

func TestDetectFunctions(t *testing.T) {
	isolateDetectionEnvironment(t)
	t.Setenv("DATABRICKS_SDK_UPSTREAM", "terraform provider")
	t.Setenv("DATABRICKS_SDK_UPSTREAM_VERSION", "1.5.0/dev")
	t.Setenv("DATABRICKS_RUNTIME_VERSION", "15.5 beta")
	t.Setenv("GITHUB_ACTIONS", "true")
	t.Setenv("AGENT", "cursor")
	t.Setenv("OMNIGENT", "1")

	if product, version := DetectUpstream(); product != "terraform-provider" || version != "1.5.0-dev" {
		t.Errorf("DetectUpstream() = (%q, %q), want (%q, %q)", product, version, "terraform-provider", "1.5.0-dev")
	}
	if got, want := DetectRuntimeVersion(), "15.5-beta"; got != want {
		t.Errorf("DetectRuntimeVersion() = %q, want %q", got, want)
	}
	if got, want := DetectCICDProvider(), "github"; got != want {
		t.Errorf("DetectCICDProvider() = %q, want %q", got, want)
	}
	if got, want := DetectAgentProvider(), "cursor"; got != want {
		t.Errorf("DetectAgentProvider() = %q, want %q", got, want)
	}
	if got, want := DetectMetaHarnessProvider(), "omnigent"; got != want {
		t.Errorf("DetectMetaHarnessProvider() = %q, want %q", got, want)
	}

	t.Setenv("AGENT", "claude-code")
	if got, want := DetectAgentProvider(), "claude-code"; got != want {
		t.Errorf("second DetectAgentProvider() = %q, want %q", got, want)
	}
}

func isolateDetectionEnvironment(t *testing.T) {
	t.Helper()
	names := make([]string, 0, len(knownAgents)+len(cicdProviders)+len(knownMetaHarnesses)+5)
	for _, agent := range knownAgents {
		names = append(names, agent.envVar)
	}
	for _, provider := range cicdProviders {
		for _, envVar := range provider.envVars {
			names = append(names, envVar.name)
		}
	}
	for _, metaHarness := range knownMetaHarnesses {
		names = append(names, metaHarness.envVar)
	}
	names = append(names, "AGENT", "AI_AGENT", "DATABRICKS_RUNTIME_VERSION", "DATABRICKS_SDK_UPSTREAM", "DATABRICKS_SDK_UPSTREAM_VERSION")
	for _, name := range names {
		value, wasSet := os.LookupEnv(name)
		if err := os.Unsetenv(name); err != nil {
			t.Fatalf("unset %s: %v", name, err)
		}
		t.Cleanup(func() {
			if wasSet {
				if err := os.Setenv(name, value); err != nil {
					t.Errorf("restore %s: %v", name, err)
				}
				return
			}
			if err := os.Unsetenv(name); err != nil {
				t.Errorf("unset %s: %v", name, err)
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
