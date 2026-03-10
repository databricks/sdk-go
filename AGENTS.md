# Databricks SDK for Go

New Go SDK for the Databricks platform. Multi-module workspace with separate
Go modules for each service area, plus hand-written authentication and core
infrastructure.

For scoped coding rules, skills, and prompt templates, see the `.agent/` directory.

### Rules (`.agent/rules/`)
- **api-design.md** — API design principles for the Databricks SDK.
- **error-handling.md** — Error handling conventions for the Databricks SDK.
- **style-guide.md** — Go code style and naming conventions for the Databricks SDK.
- **testing.md** — Testing standards and anti-patterns for the Databricks SDK.

### Skills (`.agent/skills/`)
- **write-pr-description** — Write or improve a GitHub pull request description.

### Prompts (`.agent/prompts/`)
- **code-review.md** — Review PR or code changes for backward compatibility, correctness, and style.

## Development

Prerequisites: Go 1.25+ (see `go.work`).

```bash
go test ./...          # Run tests for the current module
go vet ./...           # Vet the current module
gofmt -l -w .          # Format all Go files
golangci-lint run ./... # Lint the current module
```

CI runs tests across all modules discovered via `go.work`, with `go test -v -race`.
Linting uses `golangci-lint` and formatting is checked via `gofmt`.

## Architecture

This is a **Go workspace** (`go.work`) with multiple modules:

```
├── auth/            # Authentication infrastructure (hand-written)
│   ├── credentials/ # Credential providers (CLI, PAT, etc.)
│   └── ...          # Token types, interfaces
├── core/            # Core utilities shared across modules
├── databricks/      # Databricks workspace APIs
│   ├── jobs/        # Jobs API client
│   ├── serving/     # Model serving API client
│   └── ...
├── dataquality/     # Data quality APIs
├── iam/             # Identity and access management APIs
├── qualitymonitor/  # Quality monitor APIs
├── sdk/             # High-level SDK entry point
├── settings/        # Settings APIs
└── examples/        # Usage examples
```

## Common Mistakes

- Do NOT use `context.Background()` in production code (only in tests)
- Do NOT add dependencies without checking license compatibility (Apache 2.0 required)
- Do NOT break backwards compatibility of exported APIs without discussion

## Where to Put New Code

1. New auth method? → `auth/credentials/` (implement the appropriate interface)
2. Auth infrastructure? → `auth/`
3. New service API client? → appropriate service module (`databricks/`, `iam/`, etc.)
4. Shared utilities? → `core/`
5. Usage examples? → `examples/`

## Pull Requests

PR template requires: Changes (what + why) and Tests sections.
Run tests and lint before submitting.
