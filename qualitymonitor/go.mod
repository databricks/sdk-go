module github.com/databricks/sdk-go/qualitymonitor

go 1.26.0

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/core => ../core

replace github.com/databricks/sdk-go/databricks => ../databricks

require (
	github.com/databricks/sdk-go/core v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/databricks v0.0.0-00010101000000-000000000000
)

require github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000 // indirect
