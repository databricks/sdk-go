# Version changelog

## Release v0.9.0 (2026-11-23)

### New Features and Improvements
* Added `Download`, which streams a file, and `DownloadTo`, which writes it to a local path. With `WithDownloadParallel`, `DownloadTo` downloads a file over 50 MiB in parallel ranges, 10 at a time by default (set with `WithDownloadParallelism`), except on Windows. Both resume a transfer that breaks partway, try the next endpoint when one fails before the first byte, and skip an endpoint that keeps failing, as uploads do. `DownloadTo` overwrites an existing local file unless `WithDownloadOverwrite(false)` is set. `DownloadFile` is unchanged.
* `Upload`, `UploadFrom`, `Download` and `DownloadTo` now ask the workspace which hosts to use for file transfers (`GET /api/2.0/fs/endpoints`), and ask again every 30 minutes. A transfer tries them in the order the workspace gives. If the workspace gives no usable list, transfers use presigned URLs and then the workspace's Files API, as before. A request that starts or finishes a transfer now times out after 60 seconds without an answer, or after the HTTP client's own timeout if it sets one. A host that never answers then gives way to the next host instead of holding the transfer for 5 minutes.
* After repeated failures, `Upload` and `UploadFrom` skip an endpoint for a while and try it again from time to time. An endpoint is skipped after three uploads in a row fail on it and then complete on another, for two minutes at first and up to 30 minutes while it keeps failing. Throttling, unavailable answers, dropped connections and stalled transfers don't count, but a host that doesn't answer a request that starts or finishes a transfer within 60 seconds, or within the HTTP client's own timeout if it sets one, does. An upload that only the skipped endpoint can take still uses it.
* `Upload` and `UploadFrom` send the file again through the next endpoint, such as the workspace in one request, when the upload fails before the first part reaches cloud storage. That now includes a network error or timeout reaching cloud storage and a failed `initiate-upload`, which didn't fall back before. They don't send it again when the file already exists, the context is done, or no endpoint left can take it, as for an upload known to be over 5 GiB, the most the workspace takes in one request. Part uploads to cloud storage are retried once before the first part lands, while another endpoint is still left to try, and twice otherwise. Before, they were retried four times.
* `Upload` sends a reader of unknown size that ends before 50 MiB in one request at any parallelism, instead of opening a multipart upload. Before, only `WithParallelism(1)` did. A longer one also holds the 50 MiB it read ahead until its first parts are cut from it.

### Bug Fixes
* A failed multipart upload is aborted with the Files API's `abort-upload` action instead of a presigned DELETE, and a failed abort is logged as a warning.
* A single-request upload is no longer cut off after 5 minutes while its bytes keep moving. The 5-minute retry budget now only stops new attempts. An attempt that moves no bytes for 60 seconds fails and is retried like any other transient failure. `WithProgress` no longer counts bytes twice when a request is retried or the upload moves on to the next endpoint.
* A single-request upload of a seekable reader or a local file now retries with the whole body and sends its `Content-Length`. Before, a retry of an `*os.File` failed with "file already closed", and `Upload` closed the file the caller passed in.
* After a failed single-request upload, moving on to another endpoint now waits until that request has stopped reading the caller's reader before rewinding it. Before, a reader that can seek but has no `ReadAt` could send the other endpoint the wrong bytes. Uploads with the default settings weren't affected before this release, when the workspace was the only endpoint they sent single requests to.

### Internal Changes
* Removed unused code for resumable uploads, which the Files API no longer offers. The Files API requests that `Upload` and `UploadFrom` build themselves now send the SDK's User-Agent.

## Release v0.8.0 (2026-11-17)

## Release v0.7.0 (2026-11-16)

## Release v0.6.0 (2026-11-12)

### New Features and Improvements
* Structs in this package implement `json.Marshaler` and `json.Unmarshaler`, so they can be stored as JSON and read back.

### Breaking Changes
* `json.Marshal` and `json.Unmarshal` on structs in this package use a different JSON format than before.

## Release v0.5.0 (2026-11-11)

## Release v0.4.0 (2026-11-10)

## Release v0.3.0 (2026-10-13)

## Release v0.2.0 (2026-10-06)

## Release v0.1.0 (2026-10-02)

### New Features and Improvements
* Beta release.

## Release v0.0.1 (2026-10-01)

## Release v0.0.1-dev.11 (2026-09-30)

## Release v0.0.1-dev.10 (2026-09-23)

## Release v0.0.1-dev.9 (2026-09-15)

## Release v0.0.1-dev.8 (2026-09-10)

## Release v0.0.1-dev.7 (2026-09-08)

## Release v0.0.1-dev.6 (2026-09-07)

## Release v0.0.1-dev.5 (2026-09-06)

## Release v0.0.1-dev.4 (2026-09-05)

## Release v0.0.1-dev.3 (2026-09-03)

## Release v0.0.1-dev.2 (2026-09-02)

## Release v0.0.1-dev.1 (2026-08-20)
