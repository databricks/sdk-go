module github.com/databricks/sdk-go/iam

go 1.25.4

replace github.com/databricks/sdk-go/sdk => ../sdk

replace github.com/databricks/sdk-go/internal => ../internal

replace github.com/databricks/sdk-go/jobs => ../jobs

replace github.com/databricks/sdk-go/qualitymonitor => ../qualitymonitor

replace github.com/databricks/sdk-go/settings => ../settings

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/databricks => ../databricks

replace github.com/databricks/sdk-go/google => ../google

require github.com/databricks/sdk-go/databricks v0.0.0-00010101000000-000000000000

require github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000 // indirect
