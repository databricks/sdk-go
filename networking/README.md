# github.com/databricks/sdk-go/networking

> [!NOTE]
>
> ## Beta
>
> **This SDK is in Beta and is supported for production use cases.** Interfaces might still change slightly before GA (e.g. name standardization and minor ergonomic tweaks). We are keen to hear feedback from early adopters — please [file issues](https://github.com/databricks/sdk-go/issues), and we will address them.

## Installation

```bash
go get github.com/databricks/sdk-go/networking@latest
```

## Usage

```go
import "github.com/databricks/sdk-go/networking/v1"

client, err := networking.NewClient(ctx)
if err != nil {
	return err
}
```

For a full getting-started guide, see the [root README](../README.md).
