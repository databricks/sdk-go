# Databricks Core

[![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/core.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/core)

> [!WARNING]
> **Preview: not for production use.** This SDK is in active development. APIs are
> experimental and breaking changes may occur at any time. For production use
> cases, use the current [Databricks SDK for Go](https://github.com/databricks/databricks-sdk-go).

Internal core of the [Databricks Modular Go SDK](https://github.com/databricks/sdk-go), providing foundational primitives for error handling, operation execution with retry and rate limiting, configuration profile resolution, and client metadata collection.

## Installation

```
go get github.com/databricks/sdk-go/core
```

## Packages

| Package | Description |
| --- | --- |
| [apierr](https://pkg.go.dev/github.com/databricks/sdk-go/core/apierr) | Transport-agnostic API error types with canonical error codes and structured error details. |
| [clientinfo](https://pkg.go.dev/github.com/databricks/sdk-go/core/clientinfo) | Client and environment metadata for User-Agent headers, with auto-detection of CI/CD providers and runtimes. |
| [ops](https://pkg.go.dev/github.com/databricks/sdk-go/core/ops) | Operation execution with configurable retry, timeout, and rate limiting. |
| [profiles](https://pkg.go.dev/github.com/databricks/sdk-go/core/profiles) | Resolution of Databricks configuration profiles from `~/.databrickscfg` files and environment variables. |

## Go Version Support

See [go.work](https://github.com/databricks/sdk-go/blob/main/go.work) for the
minimum supported Go version.

## License

See [LICENSE](https://github.com/databricks/sdk-go/blob/main/LICENSE).
