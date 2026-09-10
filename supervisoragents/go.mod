module github.com/databricks/sdk-go/supervisoragents

go 1.26.0

replace github.com/databricks/sdk-go/auth => ../auth

replace github.com/databricks/sdk-go/core => ../core

replace github.com/databricks/sdk-go/options => ../options

require (
	github.com/databricks/sdk-go/auth v0.0.1-dev.8
	github.com/databricks/sdk-go/core v0.0.1-dev.8
	github.com/databricks/sdk-go/options v0.0.1-dev.8
)

require (
	github.com/databricks/databricks-sdk-go v0.92.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/oauth2 v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)
