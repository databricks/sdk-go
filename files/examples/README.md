# files/v2 examples

Runnable programs that exercise the `files/v2` client end to end. Each
subdirectory is its own `main` package that runs against a real Databricks
workspace, reading its host and credentials from the default Databricks profile.

## Programs

- `upload/` — upload a payload to a Unity Catalog volume, then download it back
  and verify the round trip.

## Running

    bazel run //deco/oss/repos/sdk-go/files/examples/upload

Configure the default Databricks profile before running (see the [Databricks
Auth] documentation), and set `remotePath` in the program to a volume you own.

[Databricks Auth]: https://pkg.go.dev/github.com/databricks/sdk-go/auth
