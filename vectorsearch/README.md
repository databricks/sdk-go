# github.com/databricks/sdk-go/vectorsearch

> [!WARNING]
>
> ## ⚠️ PREVIEW - NOT FOR PRODUCTION USE
>
> **This SDK is in active development and is subject to change without notice.**
>
> - ❌ **Do NOT use in production environments**
> - ⚠️ **Breaking changes may occur at any time**
> - 🔬 **APIs are experimental and unstable**

## Installation

```bash
go get github.com/databricks/sdk-go/vectorsearch@latest
```

## Usage

```go
import "github.com/databricks/sdk-go/vectorsearch/v1"

client, err := vectorsearch.NewClient(ctx)
if err != nil {
	return err
}
```

For a full getting-started guide, see the [root README](../README.md).
