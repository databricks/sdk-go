module github.com/databricks/sdk-go/options

go 1.26.0

replace (
	github.com/databricks/sdk-go/auth => ../auth
	github.com/databricks/sdk-go/core => ../core
)

require (
	github.com/databricks/sdk-go/auth v0.0.1-dev.1
	github.com/databricks/sdk-go/core v0.0.1-dev.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
