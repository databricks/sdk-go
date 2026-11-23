package files

// The exported downloads: Download streams a file and DownloadTo writes one to
// a local path, both through the endpoint walk of engine.download and
// engine.downloadTo.

import (
	"context"
	"net/http"
)

// DownloadToResult holds the result of a download to a local file. It is
// currently empty and exists for forward compatibility.
type DownloadToResult struct{}

// DownloadOption configures a Download or DownloadTo call.
type DownloadOption func(*downloadConfig)

// WithDownloadOverwrite controls whether DownloadTo replaces an existing local
// file. It defaults to true. With false, DownloadTo fails with an error that
// matches fs.ErrExist when the local path exists. DownloadTo only: Download
// ignores it.
func WithDownloadOverwrite(overwrite bool) DownloadOption {
	return func(c *downloadConfig) { c.noOverwrite = !overwrite }
}

// WithDownloadParallel controls whether DownloadTo downloads a large file in
// ranges, written in parallel. It defaults to false. A file of at most 50 MiB,
// a local path that isn't a regular file, such as a pipe, and any download on
// Windows are downloaded sequentially. DownloadTo only: Download ignores it.
func WithDownloadParallel(parallel bool) DownloadOption {
	return func(c *downloadConfig) { c.parallel = parallel }
}

// WithDownloadParallelism sets how many ranges a parallel DownloadTo downloads
// at once, from 1 to 64. It defaults to 10. It applies only with
// WithDownloadParallel(true), and is ignored without it. DownloadTo only:
// Download ignores it.
func WithDownloadParallelism(n int) DownloadOption {
	return func(c *downloadConfig) { c.parallelism = &n }
}

// WithDownloadTransferClient overrides the HTTP client that downloads from
// presigned cloud-storage URLs, replacing the one the call would otherwise
// create internally. Use it to route those transfers through a custom
// transport, proxy, or CA. Requests to Databricks hosts, including downloads
// from a Files API host, use the client's own transport.
//
// Because presigned URLs authenticate themselves, the client must not attach
// Databricks credentials; the storage provider rejects extra authentication. Do
// not set http.Client.Timeout, a whole-request deadline that would abort a
// legitimately long transfer; bound the download with the context passed to
// Download or DownloadTo instead.
func WithDownloadTransferClient(client *http.Client) DownloadOption {
	return func(c *downloadConfig) { c.transferClient = client }
}

func resolveDownloadConfig(opts []DownloadOption) downloadConfig {
	var cfg downloadConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// Download opens the remote file at filePath, an absolute path such as
// "/Volumes/catalog/schema/volume/file", for reading. It returns the same
// response as DownloadFile, whose Contents the caller must close.
//
// The file is read from the first endpoint that answers, in the client's
// endpoint order: by default presigned cloud-storage URLs, then the workspace.
// An endpoint that fails before the first byte arrives is passed over for the
// next one, and one that keeps failing is skipped for a while. A connection
// that breaks or stalls once bytes have arrived is resumed where it left off;
// a file that changed since the download began fails it instead.
//
// Only WithDownloadTransferClient applies; the other options are for
// DownloadTo. ctx bounds the whole download, including every read of Contents
// and its retries, so keep it alive until Contents is closed.
func (c *Client) Download(ctx context.Context, filePath string, opts ...DownloadOption) (*DownloadFileResponse, error) {
	return newClientEngine(ctx, c).download(ctx, filePath, resolveDownloadConfig(opts))
}

// DownloadTo downloads the remote file at filePath to the local file at
// localPath. A new localPath is written to a temporary file next to it and
// moved into place once the download completes, so a failed download leaves
// nothing behind. An existing localPath is written in place, unless
// WithDownloadOverwrite(false) is set, in which case DownloadTo fails before
// downloading. A localPath that is a directory is an error.
//
// The download walks the same endpoints as Download, and resumes in the same
// way. With WithDownloadParallel(true), a file over 50 MiB is downloaded in
// ranges of 10 MiB, on up to 10 at once by default (WithDownloadParallelism).
//
// ctx bounds the entire operation, including all retries; pass a context with a
// deadline to cap the total download time.
func (c *Client) DownloadTo(ctx context.Context, filePath, localPath string, opts ...DownloadOption) (*DownloadToResult, error) {
	if err := newClientEngine(ctx, c).downloadTo(ctx, filePath, localPath, resolveDownloadConfig(opts)); err != nil {
		return nil, err
	}
	return &DownloadToResult{}, nil
}
