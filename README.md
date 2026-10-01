# Databricks Modular SDKs for Go (Beta)

> [!NOTE]
>
> ## Beta
>
> **This SDK is in Beta and is supported for production use cases.** Interfaces might still change slightly before GA (e.g. name standardization and minor ergonomic tweaks). We are keen to hear feedback from early adopters — please [file issues](https://github.com/databricks/sdk-go/issues), and we will address them.

The Databricks SDKs for Go provide typed clients for the Databricks REST API. They have a **modular architecture**, with a separate Go module for each API (for example, `github.com/databricks/sdk-go/postgres`).

## Table of Contents

- [Installation](#installation)
- [Authentication](#authentication)
- [Example](#example)
- [Modules](#modules)
  - [Shared modules](#shared-modules)
- [License](#license)

## Installation

The SDK requires Go 1.26 or later. Install the module for each API you need. For example, to work with Postgres:

```bash
go get github.com/databricks/sdk-go/postgres@latest
```

## Authentication

By default, a client reads its host and credentials from a Databricks configuration profile (`~/.databrickscfg`) and `DATABRICKS_*` environment variables. With those set, no credentials need to be passed in code:

```go
import "github.com/databricks/sdk-go/postgres/v1"

postgresClient, err := postgres.NewClient(ctx)
if err != nil {
	return err
}
```

To configure credentials explicitly, create credentials with `github.com/databricks/sdk-go/auth/credentials` and pass client options from `github.com/databricks/sdk-go/options/client`:

```go
import (
	"github.com/databricks/sdk-go/auth/credentials"
	clientopts "github.com/databricks/sdk-go/options/client"
	"github.com/databricks/sdk-go/postgres/v1"
)

pat, err := credentials.NewPATCredentials("<personal-access-token>")
if err != nil {
	return err
}
postgresClient, err := postgres.NewClient(
	ctx,
	clientopts.WithHost("https://example.cloud.databricks.com"),
	clientopts.WithCredentials(pat),
)
if err != nil {
	return err
}
```

## Example

The following lists the Postgres projects you can access, using the default authentication described above. Iterator methods page through results transparently:

```go
import (
	"context"
	"fmt"

	"github.com/databricks/sdk-go/postgres/v1"
)

func listProjects(ctx context.Context) error {
	postgresClient, err := postgres.NewClient(ctx)
	if err != nil {
		return err
	}
	for project, err := range postgresClient.ListProjectsIter(ctx, postgres.ListProjectsRequest{}) {
		if err != nil {
			return err
		}
		fmt.Println(project)
	}
	return nil
}
```

Additional runnable examples are available in [`examples`](examples) and package-specific example directories such as [`files/examples`](files/examples).

## Modules

Each Databricks API is published as a separate module named `github.com/databricks/sdk-go/<api>`. Import its client from the module's versioned package — for example, `github.com/databricks/sdk-go/postgres/v1` exports `postgres.Client`.

| Module | Documentation |
| --- | --- |
| [`github.com/databricks/sdk-go/accessmanagement`](https://pkg.go.dev/github.com/databricks/sdk-go/accessmanagement) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/accessmanagement.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/accessmanagement) |
| [`github.com/databricks/sdk-go/aifunctions`](https://pkg.go.dev/github.com/databricks/sdk-go/aifunctions) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/aifunctions.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/aifunctions) |
| [`github.com/databricks/sdk-go/aigateway`](https://pkg.go.dev/github.com/databricks/sdk-go/aigateway) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/aigateway.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/aigateway) |
| [`github.com/databricks/sdk-go/aisearch`](https://pkg.go.dev/github.com/databricks/sdk-go/aisearch) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/aisearch.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/aisearch) |
| [`github.com/databricks/sdk-go/alerts`](https://pkg.go.dev/github.com/databricks/sdk-go/alerts) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/alerts.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/alerts) |
| [`github.com/databricks/sdk-go/apps`](https://pkg.go.dev/github.com/databricks/sdk-go/apps) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/apps.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/apps) |
| [`github.com/databricks/sdk-go/authentication`](https://pkg.go.dev/github.com/databricks/sdk-go/authentication) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/authentication.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/authentication) |
| [`github.com/databricks/sdk-go/budgetpolicy`](https://pkg.go.dev/github.com/databricks/sdk-go/budgetpolicy) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/budgetpolicy.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/budgetpolicy) |
| [`github.com/databricks/sdk-go/budgets`](https://pkg.go.dev/github.com/databricks/sdk-go/budgets) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/budgets.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/budgets) |
| [`github.com/databricks/sdk-go/cleanrooms`](https://pkg.go.dev/github.com/databricks/sdk-go/cleanrooms) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/cleanrooms.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/cleanrooms) |
| [`github.com/databricks/sdk-go/clusterlibraries`](https://pkg.go.dev/github.com/databricks/sdk-go/clusterlibraries) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/clusterlibraries.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/clusterlibraries) |
| [`github.com/databricks/sdk-go/clusterpolicies`](https://pkg.go.dev/github.com/databricks/sdk-go/clusterpolicies) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/clusterpolicies.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/clusterpolicies) |
| [`github.com/databricks/sdk-go/clusters`](https://pkg.go.dev/github.com/databricks/sdk-go/clusters) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/clusters.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/clusters) |
| [`github.com/databricks/sdk-go/commandexecution`](https://pkg.go.dev/github.com/databricks/sdk-go/commandexecution) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/commandexecution.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/commandexecution) |
| [`github.com/databricks/sdk-go/customllms`](https://pkg.go.dev/github.com/databricks/sdk-go/customllms) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/customllms.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/customllms) |
| [`github.com/databricks/sdk-go/database`](https://pkg.go.dev/github.com/databricks/sdk-go/database) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/database.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/database) |
| [`github.com/databricks/sdk-go/dataclassification`](https://pkg.go.dev/github.com/databricks/sdk-go/dataclassification) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/dataclassification.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/dataclassification) |
| [`github.com/databricks/sdk-go/dataquality`](https://pkg.go.dev/github.com/databricks/sdk-go/dataquality) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/dataquality.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/dataquality) |
| [`github.com/databricks/sdk-go/disasterrecovery`](https://pkg.go.dev/github.com/databricks/sdk-go/disasterrecovery) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/disasterrecovery.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/disasterrecovery) |
| [`github.com/databricks/sdk-go/domains`](https://pkg.go.dev/github.com/databricks/sdk-go/domains) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/domains.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/domains) |
| [`github.com/databricks/sdk-go/environments`](https://pkg.go.dev/github.com/databricks/sdk-go/environments) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/environments.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/environments) |
| [`github.com/databricks/sdk-go/experiments`](https://pkg.go.dev/github.com/databricks/sdk-go/experiments) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/experiments.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/experiments) |
| [`github.com/databricks/sdk-go/features`](https://pkg.go.dev/github.com/databricks/sdk-go/features) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/features.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/features) |
| [`github.com/databricks/sdk-go/featurestore`](https://pkg.go.dev/github.com/databricks/sdk-go/featurestore) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/featurestore.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/featurestore) |
| [`github.com/databricks/sdk-go/files`](https://pkg.go.dev/github.com/databricks/sdk-go/files) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/files.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/files) |
| [`github.com/databricks/sdk-go/forecasting`](https://pkg.go.dev/github.com/databricks/sdk-go/forecasting) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/forecasting.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/forecasting) |
| [`github.com/databricks/sdk-go/genie`](https://pkg.go.dev/github.com/databricks/sdk-go/genie) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/genie.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/genie) |
| [`github.com/databricks/sdk-go/gitcredentials`](https://pkg.go.dev/github.com/databricks/sdk-go/gitcredentials) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/gitcredentials.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/gitcredentials) |
| [`github.com/databricks/sdk-go/globalinitscripts`](https://pkg.go.dev/github.com/databricks/sdk-go/globalinitscripts) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/globalinitscripts.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/globalinitscripts) |
| [`github.com/databricks/sdk-go/instancepools`](https://pkg.go.dev/github.com/databricks/sdk-go/instancepools) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/instancepools.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/instancepools) |
| [`github.com/databricks/sdk-go/instanceprofiles`](https://pkg.go.dev/github.com/databricks/sdk-go/instanceprofiles) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/instanceprofiles.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/instanceprofiles) |
| [`github.com/databricks/sdk-go/jobs`](https://pkg.go.dev/github.com/databricks/sdk-go/jobs) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/jobs.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/jobs) |
| [`github.com/databricks/sdk-go/keyconfigurations`](https://pkg.go.dev/github.com/databricks/sdk-go/keyconfigurations) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/keyconfigurations.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/keyconfigurations) |
| [`github.com/databricks/sdk-go/knowledgeassistants`](https://pkg.go.dev/github.com/databricks/sdk-go/knowledgeassistants) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/knowledgeassistants.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/knowledgeassistants) |
| [`github.com/databricks/sdk-go/lakeview`](https://pkg.go.dev/github.com/databricks/sdk-go/lakeview) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/lakeview.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/lakeview) |
| [`github.com/databricks/sdk-go/logdelivery`](https://pkg.go.dev/github.com/databricks/sdk-go/logdelivery) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/logdelivery.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/logdelivery) |
| [`github.com/databricks/sdk-go/marketplaces`](https://pkg.go.dev/github.com/databricks/sdk-go/marketplaces) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/marketplaces.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/marketplaces) |
| [`github.com/databricks/sdk-go/modelregistry`](https://pkg.go.dev/github.com/databricks/sdk-go/modelregistry) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/modelregistry.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/modelregistry) |
| [`github.com/databricks/sdk-go/modelserving`](https://pkg.go.dev/github.com/databricks/sdk-go/modelserving) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/modelserving.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/modelserving) |
| [`github.com/databricks/sdk-go/modelservingquery`](https://pkg.go.dev/github.com/databricks/sdk-go/modelservingquery) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/modelservingquery.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/modelservingquery) |
| [`github.com/databricks/sdk-go/networking`](https://pkg.go.dev/github.com/databricks/sdk-go/networking) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/networking.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/networking) |
| [`github.com/databricks/sdk-go/notificationdestinations`](https://pkg.go.dev/github.com/databricks/sdk-go/notificationdestinations) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/notificationdestinations.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/notificationdestinations) |
| [`github.com/databricks/sdk-go/oauth`](https://pkg.go.dev/github.com/databricks/sdk-go/oauth) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/oauth.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/oauth) |
| [`github.com/databricks/sdk-go/pipelines`](https://pkg.go.dev/github.com/databricks/sdk-go/pipelines) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/pipelines.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/pipelines) |
| [`github.com/databricks/sdk-go/policyfamilies`](https://pkg.go.dev/github.com/databricks/sdk-go/policyfamilies) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/policyfamilies.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/policyfamilies) |
| [`github.com/databricks/sdk-go/postgres`](https://pkg.go.dev/github.com/databricks/sdk-go/postgres) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/postgres.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/postgres) |
| [`github.com/databricks/sdk-go/queries`](https://pkg.go.dev/github.com/databricks/sdk-go/queries) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/queries.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/queries) |
| [`github.com/databricks/sdk-go/queryhistory`](https://pkg.go.dev/github.com/databricks/sdk-go/queryhistory) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/queryhistory.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/queryhistory) |
| [`github.com/databricks/sdk-go/repos`](https://pkg.go.dev/github.com/databricks/sdk-go/repos) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/repos.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/repos) |
| [`github.com/databricks/sdk-go/sandbox`](https://pkg.go.dev/github.com/databricks/sdk-go/sandbox) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/sandbox.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/sandbox) |
| [`github.com/databricks/sdk-go/scim`](https://pkg.go.dev/github.com/databricks/sdk-go/scim) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/scim.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/scim) |
| [`github.com/databricks/sdk-go/secrets`](https://pkg.go.dev/github.com/databricks/sdk-go/secrets) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/secrets.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/secrets) |
| [`github.com/databricks/sdk-go/settings`](https://pkg.go.dev/github.com/databricks/sdk-go/settings) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/settings.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/settings) |
| [`github.com/databricks/sdk-go/sharing`](https://pkg.go.dev/github.com/databricks/sdk-go/sharing) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/sharing.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/sharing) |
| [`github.com/databricks/sdk-go/statementexecution`](https://pkg.go.dev/github.com/databricks/sdk-go/statementexecution) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/statementexecution.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/statementexecution) |
| [`github.com/databricks/sdk-go/storageconfigurations`](https://pkg.go.dev/github.com/databricks/sdk-go/storageconfigurations) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/storageconfigurations.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/storageconfigurations) |
| [`github.com/databricks/sdk-go/supervisoragents`](https://pkg.go.dev/github.com/databricks/sdk-go/supervisoragents) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/supervisoragents.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/supervisoragents) |
| [`github.com/databricks/sdk-go/tagassignments`](https://pkg.go.dev/github.com/databricks/sdk-go/tagassignments) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/tagassignments.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/tagassignments) |
| [`github.com/databricks/sdk-go/tagpolicies`](https://pkg.go.dev/github.com/databricks/sdk-go/tagpolicies) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/tagpolicies.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/tagpolicies) |
| [`github.com/databricks/sdk-go/tokenmanagement`](https://pkg.go.dev/github.com/databricks/sdk-go/tokenmanagement) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/tokenmanagement.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/tokenmanagement) |
| [`github.com/databricks/sdk-go/tokens`](https://pkg.go.dev/github.com/databricks/sdk-go/tokens) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/tokens.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/tokens) |
| [`github.com/databricks/sdk-go/uc/abacpolicies`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/abacpolicies) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/abacpolicies.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/abacpolicies) |
| [`github.com/databricks/sdk-go/uc/artifactallowlists`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/artifactallowlists) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/artifactallowlists.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/artifactallowlists) |
| [`github.com/databricks/sdk-go/uc/catalogs`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/catalogs) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/catalogs.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/catalogs) |
| [`github.com/databricks/sdk-go/uc/connections`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/connections) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/connections.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/connections) |
| [`github.com/databricks/sdk-go/uc/credentials`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/credentials) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/credentials.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/credentials) |
| [`github.com/databricks/sdk-go/uc/entitytagassignments`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/entitytagassignments) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/entitytagassignments.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/entitytagassignments) |
| [`github.com/databricks/sdk-go/uc/externallineage`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externallineage) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/externallineage.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externallineage) |
| [`github.com/databricks/sdk-go/uc/externallocations`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externallocations) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/externallocations.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externallocations) |
| [`github.com/databricks/sdk-go/uc/externalmetadata`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externalmetadata) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/externalmetadata.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/externalmetadata) |
| [`github.com/databricks/sdk-go/uc/functions`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/functions) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/functions.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/functions) |
| [`github.com/databricks/sdk-go/uc/grants`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/grants) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/grants.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/grants) |
| [`github.com/databricks/sdk-go/uc/metastores`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/metastores) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/metastores.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/metastores) |
| [`github.com/databricks/sdk-go/uc/onlinetables`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/onlinetables) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/onlinetables.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/onlinetables) |
| [`github.com/databricks/sdk-go/uc/registeredmodels`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/registeredmodels) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/registeredmodels.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/registeredmodels) |
| [`github.com/databricks/sdk-go/uc/resourcequotas`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/resourcequotas) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/resourcequotas.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/resourcequotas) |
| [`github.com/databricks/sdk-go/uc/rfa`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/rfa) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/rfa.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/rfa) |
| [`github.com/databricks/sdk-go/uc/schemas`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/schemas) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/schemas.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/schemas) |
| [`github.com/databricks/sdk-go/uc/secrets`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/secrets) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/secrets.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/secrets) |
| [`github.com/databricks/sdk-go/uc/systemschemas`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/systemschemas) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/systemschemas.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/systemschemas) |
| [`github.com/databricks/sdk-go/uc/tables`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/tables) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/tables.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/tables) |
| [`github.com/databricks/sdk-go/uc/volumes`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/volumes) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/volumes.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/volumes) |
| [`github.com/databricks/sdk-go/uc/workspacebindings`](https://pkg.go.dev/github.com/databricks/sdk-go/uc/workspacebindings) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/uc/workspacebindings.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/uc/workspacebindings) |
| [`github.com/databricks/sdk-go/usagedashboards`](https://pkg.go.dev/github.com/databricks/sdk-go/usagedashboards) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/usagedashboards.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/usagedashboards) |
| [`github.com/databricks/sdk-go/vectorsearch`](https://pkg.go.dev/github.com/databricks/sdk-go/vectorsearch) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/vectorsearch.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/vectorsearch) |
| [`github.com/databricks/sdk-go/warehouses`](https://pkg.go.dev/github.com/databricks/sdk-go/warehouses) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/warehouses.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/warehouses) |
| [`github.com/databricks/sdk-go/workspaces`](https://pkg.go.dev/github.com/databricks/sdk-go/workspaces) | [![Go Reference](https://pkg.go.dev/badge/github.com/databricks/sdk-go/workspaces.svg)](https://pkg.go.dev/github.com/databricks/sdk-go/workspaces) |

### Shared modules

Three modules are shared by every API client and provide the pieces you import directly:

- [`github.com/databricks/sdk-go/core`](https://pkg.go.dev/github.com/databricks/sdk-go/core) — configuration-profile resolution, API errors, retry and rate-limiting primitives, client metadata, and shared types.
- [`github.com/databricks/sdk-go/auth`](https://pkg.go.dev/github.com/databricks/sdk-go/auth) — credential providers and the default credential chain.
- [`github.com/databricks/sdk-go/options`](https://pkg.go.dev/github.com/databricks/sdk-go/options) — options for clients, calls, and long-running operations.

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
