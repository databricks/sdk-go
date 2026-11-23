# github.com/databricks/sdk-go/files

> [!NOTE]
>
> ## Beta
>
> **This SDK is in Beta and is supported for production use cases.** Interfaces might still change slightly before GA (e.g. name standardization and minor ergonomic tweaks). We are keen to hear feedback from early adopters — please [file issues](https://github.com/databricks/sdk-go/issues), and we will address them.

## Installation

```bash
go get github.com/databricks/sdk-go/files@latest
```

## Usage

```go
import "github.com/databricks/sdk-go/files/v2"

client, err := files.NewClient(ctx)
if err != nil {
	return err
}
```

For a full getting-started guide, see the [root README](../README.md).

## Uploading files

`Upload` sends the contents of an `io.Reader` to a file in a Unity Catalog volume, and `UploadFrom` sends a local file:

```go
_, err := client.Upload(ctx, "/Volumes/main/default/my_volume/report.csv", strings.NewReader(report),
	files.WithOverwrite(false),
)
if errors.Is(err, files.ErrAlreadyExists) {
	// The file exists, and WithOverwrite(false) asked not to replace it.
}

_, err = client.UploadFrom(ctx, "/Volumes/main/default/my_volume/data.parquet", "data.parquet",
	files.WithOverwrite(true),
	files.WithProgress(func(p files.Progress) { log.Printf("%d of %d bytes", p.Transferred, p.Total) }),
)
```

How an upload is sent:

- **Endpoints.** The SDK asks the workspace which hosts to use for file transfers, and asks again every 30 minutes. A transfer tries them in the order the workspace gives. If the workspace gives no usable list, transfers use presigned URLs and then the workspace's Files API.
- **Size.** Under 50 MiB, the file goes in one request, so it passes over presigned URLs, which take only parts. From 50 MiB, it goes in parts to an endpoint that takes parts, such as cloud storage over presigned URLs, or in one request to one that doesn't, such as the workspace, which takes at most 5 GiB. A reader of unknown size is read ahead up to 50 MiB to choose.
- **Fallback.** If the upload fails on an endpoint before any part has landed, for example because a firewall blocks the storage account, it is sent again through the next endpoint, such as the workspace. It isn't sent again if the file already exists, if the context is done, or if no endpoint left can take it, as for an upload known to be over 5 GiB, the most the workspace takes in one request. To send it again, the upload rewinds a seekable reader, rereads a local file, or replays the bytes it has kept from a stream.
- **Skipping.** After several uploads in a row fail on an endpoint and then complete on a later one, the client skips that endpoint for a while and tries it again from time to time. An upload only the skipped endpoint can take, such as one over 5 GiB or of unknown size when no other endpoint takes parts, still goes there.
- **Retries.** Failed requests are retried, except a one-request upload of a reader that can't be rewound. A one-request upload is cut off only when no bytes move for 60 seconds. The context bounds the whole upload, retries included, so set a deadline on it to cap the upload time.
- **Memory.** With more than one part in flight, the default, parts are read from an `io.ReaderAt`, such as an `*os.File`, as they are sent. Otherwise each in-flight part is held in memory: by default up to 256 MiB of parts for a reader without random access. A reader of unknown size also holds up to 50 MiB it read ahead until its first parts are cut from it.

Options:

| Option | Effect |
|---|---|
| `WithOverwrite(bool)` | Whether to replace an existing file. Without it, the server's default applies. |
| `WithPartSize(int64)` | Set the part size. By default it is chosen from the content length. |
| `WithParallelism(int)` | Set how many parts are sent at once. |
| `WithProgress(ProgressFunc)` | Report the bytes uploaded so far. |
| `WithLimiter(Limiter)` | Cap concurrent transfers across every upload given the same `Limiter`, such as `NewLimiter(n)`. |
| `WithTransferClient(*http.Client)` | Send parts to presigned URLs through this client, for example for a proxy or custom CA. It must not add Databricks credentials. Requests to Databricks, including the fallback through the workspace, use the client's own transport. |

## Downloading files

`Download` streams a file from a Unity Catalog volume, and `DownloadTo` writes one to a local path:

```go
resp, err := client.Download(ctx, "/Volumes/main/default/my_volume/report.csv")
if err != nil {
	return err
}
defer resp.Contents.Close()
_, err = io.Copy(os.Stdout, resp.Contents)

_, err = client.DownloadTo(ctx, "/Volumes/main/default/my_volume/data.parquet", "data.parquet",
	files.WithDownloadParallel(true),
)
```

`Download` returns the same response as the generated `DownloadFile`, which is unchanged.

How a download is sent:

- **Endpoints.** The file is read from the first endpoint in the same list. If one fails before the first byte arrives, for example because a firewall blocks the storage account, the next one is tried. After several downloads in a row fail on an endpoint and then complete on a later one, the client skips that endpoint for a while and tries it again from time to time.
- **Resuming.** A connection that breaks, or moves no bytes for 60 seconds, after bytes have arrived is resumed where it left off. If the file changed since the download began, the download fails instead.
- **Local file.** `DownloadTo` writes a new local file to a temporary file next to it and moves it into place when the download completes, so a failed download leaves nothing behind. An existing file is written in place.
- **Parallel ranges.** With `WithDownloadParallel(true)`, `DownloadTo` downloads a file over 50 MiB in ranges of 10 MiB, 10 at a time by default. Smaller files, local paths that aren't regular files, and downloads on Windows go sequentially.
- **Retries.** Failed requests are retried. The context bounds the whole download, retries and reads of `Contents` included, so keep it alive until `Contents` is closed.

Options:

| Option | Effect |
|---|---|
| `WithDownloadOverwrite(bool)` | Whether `DownloadTo` replaces an existing local file. Defaults to true. With false, it fails with an error matching `fs.ErrExist`. |
| `WithDownloadParallel(bool)` | Whether `DownloadTo` downloads a large file in parallel ranges. Defaults to false. |
| `WithDownloadParallelism(int)` | How many ranges a parallel `DownloadTo` downloads at once, from 1 to 64. Defaults to 10, and is ignored without `WithDownloadParallel(true)`. |
| `WithDownloadTransferClient(*http.Client)` | Download from presigned URLs through this client, for example for a proxy or custom CA. It must not add Databricks credentials. Requests to Databricks use the client's own transport. |

`Download` uses only `WithDownloadTransferClient`. The other options apply to `DownloadTo`.
