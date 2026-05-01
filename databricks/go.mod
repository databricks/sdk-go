module github.com/databricks/sdk-go/databricks

go 1.26.0

replace (
	github.com/databricks/sdk-go/auth => ../auth
	github.com/databricks/sdk-go/core => ../core
	github.com/databricks/sdk-go/options => ../options
)

require (
	github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/options v0.0.0-00010101000000-000000000000
)

require github.com/databricks/sdk-go/core v0.0.0-00010101000000-000000000000 // indirect
