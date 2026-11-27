# Version changelog

## Release v0.2.0 (2026-11-27)

### New Features and Improvements
* `ops.Execute` and `ops.ExecuteWithResult` now keep the error of the last attempt when the context ends, or the limiter fails, after a failed attempt. The message of the returned error names that error, and the new `ops.LastError` returns it. `errors.Is` and `errors.As` still match only the context or limiter error, so error handling that reads the API code or status does not change.

## Release v0.1.0 (2026-10-02)

### New Features and Improvements
* Beta release.

## Release v0.0.1 (2026-10-01)

## Release v0.0.1-dev.8 (2026-09-10)

## Release v0.0.1-dev.7 (2026-09-08)

## Release v0.0.1-dev.6 (2026-09-07)

## Release v0.0.1-dev.5 (2026-09-06)

## Release v0.0.1-dev.4 (2026-09-05)

## Release v0.0.1-dev.3 (2026-09-03)

## Release v0.0.1-dev.2 (2026-09-02)

## Release v0.0.1-dev.1 (2026-08-20)
