module github.com/databricks/sdk-go/modelservingquery

go 1.26.0

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/core => ../core

replace github.com/databricks/sdk-go/options => ../options

require (
	github.com/databricks/sdk-go/auth v0.0.1-dev.1
	github.com/databricks/sdk-go/core v0.0.1-dev.1
	github.com/databricks/sdk-go/options v0.0.1-dev.1
	github.com/google/go-cmp v0.7.0
)

require gopkg.in/ini.v1 v1.67.0 // indirect
