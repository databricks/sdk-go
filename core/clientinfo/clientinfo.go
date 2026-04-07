// Package clientinfo collects information about the client and its
// environment into an immutable [ClientInfo] value.
//
// [Default] returns a [ClientInfo] pre-populated with process-level
// defaults (SDK version, Go version, OS, and detected environment).
// [ClientInfo.With] derives a new value with additional key/value
// segments; it never mutates the original.
//
// Databricks tooling (Terraform provider, CLI, partner integrations)
// should call [SetProduct], [SetPartner], or [AddToDefault] at startup
// to register global metadata before any client is created.
package clientinfo

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/databricks/sdk-go/core/internal"
)

var (
	// ErrOddKeyvals indicates an odd number of key/value arguments.
	ErrOddKeyvals = errors.New("odd number of key/value arguments")

	// ErrInvalidKey indicates a segment key with invalid characters.
	ErrInvalidKey = errors.New("invalid key")

	// ErrInvalidValue indicates a segment value with invalid characters.
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidVersion indicates a version string that is not valid semver.
	ErrInvalidVersion = errors.New("invalid version")
)

// ClientInfo is an immutable, ordered list of key/value segments. Use
// [Default] to get the process-level defaults and [ClientInfo.With] to
// derive new values with additional segments.
type ClientInfo struct {
	segments []segment
}

type segment struct {
	key, value string
}

// With returns a new [ClientInfo] with the given key/value pairs appended.
// The original is not modified. An odd number of arguments returns an error.
//
// Keys and values must contain only alphanumeric characters plus the
// characters: undescore ('_'), dot ('.'), plus ('+'), or hyphen ('-').
//
// Exact key+value duplicates are silently ignored. On error, the zero value
// is returned (all-or-nothing).
func (ci ClientInfo) With(keyvals ...string) (ClientInfo, error) {
	if len(keyvals)%2 != 0 {
		return ClientInfo{}, fmt.Errorf("%w: got %d", ErrOddKeyvals, len(keyvals))
	}
	if len(keyvals) == 0 {
		return ci, nil
	}

	newSegments := make([]segment, len(ci.segments), len(ci.segments)+len(keyvals)/2)
	copy(newSegments, ci.segments)

	for i := 0; i < len(keyvals); i += 2 {
		key, value := keyvals[i], keyvals[i+1]
		if !isValidSegment(key) {
			return ClientInfo{}, fmt.Errorf("%w: %s", ErrInvalidKey, key)
		}
		if !isValidSegment(value) {
			return ClientInfo{}, fmt.Errorf("%w for %q: %s", ErrInvalidValue, key, value)
		}
		if containsSegment(newSegments, key, value) {
			continue
		}
		newSegments = append(newSegments, segment{key, value})
	}

	return ClientInfo{segments: newSegments}, nil
}

// String returns a string representation of the client info suitable for
// inclusion in HTTP headers. Key/value pairs are formatted as "key/value"
// and joined by spaces in the order they were inserted.
func (ci ClientInfo) String() string {
	parts := make([]string, len(ci.segments))
	for i, s := range ci.segments {
		parts[i] = s.key + "/" + s.value
	}
	return strings.Join(parts, " ")
}

func containsSegment(segments []segment, key, value string) bool {
	for _, s := range segments {
		if s.key == key && s.value == value {
			return true
		}
	}
	return false
}

// base holds segments added via [AddToDefault], [SetProduct], and
// [SetPartner]. [Default] returns a copy of this with env detection
// appended.
var base ClientInfo

// Default returns a [ClientInfo] populated with SDK metadata, runtime
// information, segments registered via [AddToDefault], and automatically
// detected environment properties.
func Default() ClientInfo {
	return defaultWithEnv(os.LookupEnv)
}

// lookupFunc has the same signature as [os.LookupEnv]. It is used to
// abstract environment access for testing.
type lookupFunc func(string) (string, bool)

func defaultWithEnv(lookupEnv lookupFunc) ClientInfo {
	// 3 fixed + base segments + up to 5 env detection segments.
	s := make([]segment, 0, 3+len(base.segments)+5)
	s = append(s,
		segment{internal.ModuleName, internal.Version},
		segment{"go", cachedGoVersion},
		segment{"os", runtime.GOOS},
	)
	s = append(s, base.segments...)
	// DATABRICKS_SDK_UPSTREAM and DATABRICKS_SDK_UPSTREAM_VERSION are set
	// by tools built on top of this SDK (e.g. Terraform provider, Pulumi)
	// to identify themselves as the upstream product. Both must be present
	// for the upstream segment to be included.
	if p, ok := lookupEnv("DATABRICKS_SDK_UPSTREAM"); ok {
		if v, ok := lookupEnv("DATABRICKS_SDK_UPSTREAM_VERSION"); ok {
			s = append(s, segment{"upstream", sanitize(p)}, segment{"upstream-version", sanitize(v)})
		}
	}
	if p := detectCICD(lookupEnv); p != "" {
		s = append(s, segment{"cicd", p})
	}
	if v, ok := lookupEnv("DATABRICKS_RUNTIME_VERSION"); ok && v != "" {
		s = append(s, segment{"runtime", sanitize(v)})
	}
	if a := detectAgent(lookupEnv); a != "" {
		s = append(s, segment{"agent", a})
	}
	return ClientInfo{segments: s}
}

// SetProduct sets the product name and version globally. The version must
// be a valid semver string.
//
// Must be called before any client is created. Not safe for concurrent use.
func SetProduct(name, version string) error {
	if !isSemVer(version) {
		return fmt.Errorf("%w: %s", ErrInvalidVersion, version)
	}
	return AddToDefault(name, version)
}

// SetPartner adds a partner identifier globally. Partner attribution is
// a first-class concept used for support ticket routing.
//
// Must be called before any client is created. Not safe for concurrent use.
func SetPartner(partner string) error {
	return AddToDefault("partner", partner)
}

// AddToDefault adds a global key/value segment that will be included in
// every [Default] call. Same key with different values is allowed (e.g.,
// multiple partners). Exact key+value duplicates are silently ignored.
//
// Must be called before any client is created. Not safe for concurrent use.
func AddToDefault(key, value string) error {
	ci, err := base.With(key, value)
	if err != nil {
		return err
	}
	base = ci
	return nil
}

const (
	semVerCore          = `(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)`
	semVerPrerelease    = `(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?`
	semVerBuildmetadata = `(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
)

var (
	regexpSemVer             = regexp.MustCompile(`^` + semVerCore + semVerPrerelease + semVerBuildmetadata + `$`)
	regexpValidSegment       = regexp.MustCompile(`^[0-9A-Za-z_.\+-]+$`)
	regexpInvalidSegmentChar = regexp.MustCompile(`[^0-9A-Za-z_.\+-]`)
)

func isSemVer(s string) bool       { return regexpSemVer.MatchString(s) }
func isValidSegment(s string) bool { return regexpValidSegment.MatchString(s) }

// sanitize replaces characters that are not valid in segment values with
// hyphens. Used for environment-sourced values (runtime version, upstream)
// that we do not control and cannot reject.
func sanitize(s string) string {
	return regexpInvalidSegmentChar.ReplaceAllString(s, "-")
}

type agentDef struct {
	envVar  string
	product string
}

var knownAgents = []agentDef{
	{"ANTIGRAVITY_AGENT", "antigravity"},
	{"CLAUDECODE", "claude-code"},
	{"CLINE_ACTIVE", "cline"},
	{"CODEX_CI", "codex"},
	{"COPILOT_CLI", "copilot-cli"},
	{"CURSOR_AGENT", "cursor"},
	{"GEMINI_CLI", "gemini-cli"},
	{"OPENCODE", "opencode"},
	{"OPENCLAW_SHELL", "openclaw"},
}

type envCheck struct {
	name          string
	expectedValue string
}

type cicdDef struct {
	name    string
	envVars []envCheck
}

var cicdProviders = []cicdDef{
	{"github", []envCheck{{"GITHUB_ACTIONS", "true"}}},
	{"gitlab", []envCheck{{"GITLAB_CI", "true"}}},
	{"jenkins", []envCheck{{"JENKINS_URL", ""}}},
	{"azure-devops", []envCheck{{"TF_BUILD", "True"}}},
	{"circle", []envCheck{{"CIRCLECI", "true"}}},
	{"travis", []envCheck{{"TRAVIS", "true"}}},
	{"bitbucket", []envCheck{{"BITBUCKET_BUILD_NUMBER", ""}}},
	{"google-cloud-build", []envCheck{{"PROJECT_ID", ""}, {"BUILD_ID", ""}, {"PROJECT_NUMBER", ""}, {"LOCATION", ""}}},
	{"aws-code-build", []envCheck{{"CODEBUILD_BUILD_ARN", ""}}},
	{"tf-cloud", []envCheck{{"TFC_RUN_ID", ""}}},
}

// detectAgent returns the name of a single detected AI coding agent, or
// empty if zero or more than one agent is detected. When multiple agents
// are present (e.g. Claude from within Cursor), we cannot reliably
// determine which one initiated the request, so we omit the segment.
//
// TODO: support reporting multiple concurrent agents.
func detectAgent(lookupEnv lookupFunc) string {
	var detected string
	count := 0
	for _, a := range knownAgents {
		if _, ok := lookupEnv(a.envVar); ok {
			detected = a.product
			count++
			if count > 1 {
				break
			}
		}
	}
	if count == 1 {
		return detected
	}
	return ""
}

func detectCICD(lookupEnv lookupFunc) string {
	for _, p := range cicdProviders {
		detected := true
		for _, ev := range p.envVars {
			v, ok := lookupEnv(ev.name)
			if !ok || (ev.expectedValue != "" && v != ev.expectedValue) {
				detected = false
				break
			}
		}
		if detected {
			return p.name
		}
	}
	return ""
}

// cachedGoVersion is computed once at package init because
// runtime.Version() never changes during a process lifetime.
var cachedGoVersion = normalizeGoVersion(runtime.Version())

// normalizeGoVersion converts a Go version string (e.g., "go1.26",
// "go1.26.0", "go1.26rc1") into a semver-compliant three-part version
// string (e.g., "1.26.0", "1.26.0-rc1").
func normalizeGoVersion(raw string) string {
	// Development builds return "devel +<hash> <date>" or similar.
	// These don't follow the "goX.Y.Z" convention, so return a fixed
	// dev version rather than producing a malformed string.
	if !strings.HasPrefix(raw, "go") {
		return "0.0.0-dev"
	}
	raw = strings.TrimPrefix(raw, "go")

	// Separate numeric prefix (e.g., "1.26.0") from pre-release suffix
	// (e.g., "rc1", "beta2"). The suffix starts at the first character
	// that is not a digit or dot.
	var suffix string
	prefix := raw
	dotCount := 0
	for i, c := range raw {
		if c == '.' {
			dotCount++
			continue
		}
		if c < '0' || c > '9' {
			suffix = raw[i:]
			prefix = raw[:i]
			break
		}
	}

	// Pad prefix.
	switch dotCount {
	case 0:
		prefix += ".0.0"
	case 1:
		prefix += ".0"
	}

	if suffix != "" {
		return fmt.Sprintf("%s-%s", prefix, suffix)
	}
	return prefix
}
