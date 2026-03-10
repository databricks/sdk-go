---
description: "Testing standards and anti-patterns for the Databricks SDK."
globs: "**/*_test.go"
alwaysApply: false
---

# Testing

Use Go's standard `testing` package with `google/go-cmp` for complex comparisons.
Do NOT introduce `testify` — this SDK minimizes third-party dependencies.

Table-driven tests are preferred:

```go
for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        // ...
        if diff := cmp.Diff(tc.want, got); diff != "" {
            t.Errorf("mismatch (-want +got):\n%s", diff)
        }
    })
}
```

- Unit tests: `*_test.go` alongside source files
- Call `t.Helper()` at the start of test helper functions (correct line numbers in failures)
- For configurable options, always test three cases: (1) explicit values, (2) nil/unset default, (3) empty value
- Inject `time.Now` as `func() time.Time` for time-dependent logic — makes tests deterministic
- Always check errors — never use `_` for error returns in tests
- Assert on exact expected values, not partial matches (e.g., full header value, not just prefix)
- Extract repeated test values to shared `const`/`var` to convey intent and avoid silent drift
- Add explicit timeouts to channel/goroutine-based tests to prevent hanging
- Use `defer close(ch)` for channel cleanup to handle error paths safely

## Naming and Style

- Use `want` for expected values and `got` for actual values — not `expected`/`actual`
- Use **named struct fields** in test mock initialization: `&IDToken{Value: tok}` not `&IDToken{tok}`
- Use existing test utilities rather than creating new one-off mock types

## Test Design Philosophy

- **Pragmatic about test value**: don't write tests that can't catch real regressions.
  If removing the tested code wouldn't be caught by the test, the test is not useful
- **Prefer blackbox tests over whitebox tests**: verify observable behavior, not
  implementation details. When structural constraints make blackbox testing impossible,
  consider refactoring for testability — but not at the cost of breaking backward compatibility
- **Test at the right abstraction level**: test behavior through the public interface of
  a component, not through its internals
- **Mock at the right level**: mock the nearest interface boundary, not deep internals
- **Prefer isolated tests by design**: use dependency injection and mocks rather than
  relying on environment cleanup utilities

## Anti-Patterns (never do these)

- Deleting or skipping flaky tests instead of fixing the root cause
- Tests depending on execution order or shared mutable state
- Non-descriptive test names (`TestFoo`, `TestBasic`) — follow the same pattern as Go
  `Example` functions: `TestFunctionName_description` (e.g., `TestHostType_missingScheme`,
  `TestConfig_HostType_missingScheme` for a method on a type)
- `time.Now` in tests — use a mock clock
- Exporting methods solely for test access — use `cmp.Diff` to compare structs,
  or write internal (`_test.go` in same package) tests instead
- Unnecessary concurrency in tests — if goroutines/channels aren't clearly needed, don't use them
