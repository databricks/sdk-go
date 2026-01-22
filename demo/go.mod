module github.com/databricks/sdk-go/demo

go 1.25.4

replace github.com/databricks/sdk-go/qualitymonitor => ../qualitymonitor

replace github.com/databricks/sdk-go/dataquality => ../dataquality

replace github.com/databricks/sdk-go/databricks => ../databricks

replace github.com/databricks/sdk-go/auth => ../auth

require (
	github.com/databricks/sdk-go/dataquality v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/databricks v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/qualitymonitor v0.0.0-00010101000000-000000000000
)
