module github.com/databricks/sdk-go/jobs

go 1.25.5

require (
	github.com/databricks/sdk-go/databricks v0.0.0-00010101000000-000000000000
	github.com/databricks/sdk-go/internal v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.7.0
)

replace github.com/databricks/sdk-go/databricks => ../databricks
replace github.com/databricks/sdk-go/internal => ../internal
