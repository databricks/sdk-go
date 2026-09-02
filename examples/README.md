# Examples

Each subdirectory contains a complete program for one Databricks Go SDK task.
Choose the example closest to what you want to build, then copy it or run it
from this directory.

All programs use the `examples` build tag. Include it when running one:

```sh
go run -tags=examples ./time_types
```

## Choose an example

| Program | What it demonstrates | Required setup |
| --- | --- | --- |
| `authentication/resolution` | Create a client with automatic credential resolution | A profile or authentication environment variables (see below) |
| `authentication/m2m` | Construct and pass OAuth M2M credentials explicitly | `DATABRICKS_HOST`, `DATABRICKS_CLIENT_ID`, and `DATABRICKS_CLIENT_SECRET` |
| `authentication/custom` | Implement `auth.Credentials` for an external token provider | Set `workspaceHost` in the program and provide `EXTERNAL_ACCESS_TOKEN` |
| `client_options` | Set a profile, logger, and request timeout | A Databricks profile |
| `pagination` | Iterate across pages and stop after enough results | A profile or authentication environment variables |
| `error_handling` | Inspect the code, HTTP status, and message of an API error | A profile or authentication environment variables |
| `field_masks` | Build an update with a typed field mask, including a oneof path and a field to clear | None |
| `oneofs` | Construct and inspect mutually exclusive oneof variants | None |
| `time_types` | Convert SDK timestamp and duration values | None |

## Authenticate to a workspace

For examples that use a workspace, install the
[Databricks CLI](https://docs.databricks.com/aws/en/dev-tools/cli/install), then
log in and select the profile explicitly:

```sh
databricks auth login --host https://<workspace-host> --profile SDK_EXAMPLES
DATABRICKS_CONFIG_PROFILE=SDK_EXAMPLES go run -tags=examples ./pagination
```

Alternatively, provide credentials through environment variables:

```sh
DATABRICKS_HOST=https://<workspace-host> \
DATABRICKS_CLIENT_ID=<service-principal-client-id> \
DATABRICKS_CLIENT_SECRET=<service-principal-client-secret> \
  go run -tags=examples ./pagination
```

The `client_options` example passes the profile directly as a client option:

```sh
go run -tags=examples ./client_options -profile SDK_EXAMPLES
```
