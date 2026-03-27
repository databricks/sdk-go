module github.com/databricks/sdk-go/examples

go 1.26.0

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/databricks => ../databricks

replace github.com/databricks/sdk-go/dataquality => ../dataquality

require (
	github.com/databricks/sdk-go/auth v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/databricks v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/dataquality v0.0.0-00010101000000-000000000000
)

require github.com/google/go-cmp v0.7.0
