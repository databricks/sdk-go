package files

// Shared buffering and cloud-leg header helpers for the large-file uploader,
// plus the best-effort abort glue between the orchestration and the two
// transports. The cloud-transport pieces (part bodies, retry classification,
// presigned-URL-expiry detection) live in the cloudstorage subpackage.

import (
	"context"
	"io"
	"maps"
)

// --- Buffer helpers ---

// fillBuffer reads from r until buf holds at least minSize bytes or the stream
// ends. A short read at end-of-stream is not an error.
func fillBuffer(buf []byte, minSize int64, r io.Reader) ([]byte, error) {
	need := minSize - int64(len(buf))
	if need <= 0 {
		return buf, nil
	}
	tmp := make([]byte, need)
	n, err := io.ReadFull(r, tmp)
	buf = append(buf, tmp[:n]...)
	if err == nil || err == io.EOF || err == io.ErrUnexpectedEOF {
		return buf, nil
	}
	return buf, err
}

// readUpTo reads up to n bytes from r, returning fewer only at end-of-stream.
func readUpTo(r io.Reader, n int64) ([]byte, error) {
	buf := make([]byte, n)
	read, err := io.ReadFull(r, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return buf[:read], nil
}

// --- Header helpers ---

func mergeHeaders(base, override map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(override))
	maps.Copy(out, base)
	maps.Copy(out, override)
	return out
}

// octetStreamHeaders returns the headers for a cloud-storage request: the binary
// content type plus the presigned URL's own headers (which win on conflict). The
// returned map is freshly allocated, so callers may add request-specific headers
// such as Content-Range.
func octetStreamHeaders(presigned map[string]string) map[string]string {
	return mergeHeaders(map[string]string{"Content-Type": "application/octet-stream"}, presigned)
}

// --- Single-shot delegation and best-effort abort glue ---

// singleShotUpload sends the whole body, of size bytes or -1 if unknown, in one
// PUT to the current endpoint's Files API. Progress is reported as the body
// streams.
func (e *engine) singleShotUpload(ctx context.Context, uc *uploadContext, body io.Reader, size int64) error {
	if uc.progress != nil {
		body = withProgress(body, uc.progress)
	}
	// A single-shot upload is one transfer; hold one slot for its duration so it
	// draws from the same concurrency budget as multipart parts.
	if err := uc.limiter.Acquire(ctx); err != nil {
		return err
	}
	defer uc.limiter.Release()
	return e.uploadSingleShot(ctx, uc, body, size)
}

// abortMultipartBestEffort aborts the multipart session on a context detached
// from the caller's (see cleanupContext), and logs rather than returns a
// failure.
func (e *engine) abortMultipartBestEffort(ctx context.Context, uc *uploadContext, token string) {
	ctx, cancel := e.cleanupContext(ctx)
	defer cancel()
	if err := e.abortMultipart(ctx, uc, token); err != nil {
		// Best-effort cleanup: the server expires abandoned sessions, so a failed
		// abort is logged but doesn't change the upload's outcome.
		e.c.logger.WarnContext(ctx, "failed to abort multipart upload", "error", err)
	}
}

// cleanupContext returns a context for a best-effort cleanup (aborting a partial
// upload) that is detached from the caller's cancellation and deadline but
// bounded by its own timeout. Cleanup most needs to run exactly when the upload
// context has already been cancelled or has expired; reusing it would make the
// abort fail instantly and leak the partial upload on the storage provider.
func (e *engine) cleanupContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), e.tun.cleanupTimeout)
}
