# Examples

This directory contains example code that demonstrates how to use the Databricks Go SDK. These examples are **for illustration purposes only** and are not intended for production use.

## Running

By default, examples run against a built-in mock server so no credentials are needed:

```sh
go run github.com/databricks/sdk-go/examples/dataquality
```

To run against a real Databricks workspace, set the following environment variables:

```sh
DATABRICKS_HOST=https://....databricks.com \
DATABRICKS_TOKEN=dapi... \
  go run github.com/databricks/sdk-go/examples/dataquality
```
