module github.com/databricks/sdk-go/examples

go 1.26.0

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/core => ../core

replace github.com/databricks/sdk-go/dataquality => ../dataquality

replace github.com/databricks/sdk-go/options => ../options

require (
	github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/core v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/dataquality v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/options v0.0.0-00010101000000-000000000000
)

require github.com/google/go-cmp v0.7.0

require gopkg.in/ini.v1 v1.67.0 // indirect
