package files

// Downloads to a local file: sequentially, or in ranges written in parallel,
// through the same endpoint walk as a streaming download.

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// downloadGOOS is runtime.GOOS, a variable so that tests can take the Windows
// path.
var downloadGOOS = runtime.GOOS

// downloadTo downloads the file to destination. A new destination is written
// to a temp file next to it, moved into place once the download finishes, so a
// failed download leaves nothing behind; an existing one is written in place,
// as opening it for writing would: through a symlink, keeping its permissions,
// and into a pipe or a device.
func (e *engine) downloadTo(ctx context.Context, path, destination string, cfg downloadConfig) error {
	// The check is best effort: nothing stops a file being created at the
	// destination after it.
	if info, err := os.Stat(destination); err == nil {
		if cfg.noOverwrite {
			return &fs.PathError{Op: "download", Path: destination, Err: fs.ErrExist}
		}
		if info.IsDir() {
			return fmt.Errorf("download to %s: the destination is a directory", destination)
		}
	}
	parallel := cfg.parallel && takesRanges(destination)
	if parallel && downloadGOOS == "windows" {
		// On Windows, the file is downloaded sequentially instead.
		e.c.logger.WarnContext(ctx, "parallel download is not supported on Windows, so the file is downloaded sequentially")
		parallel = false
	}
	ranges := []byteRange{{start: 0, end: -1}}
	size := int64(-1)
	modifiedAt := ""
	parallelism := 1
	if parallel {
		n := cfg.resolvedParallelism(e.tun)
		if n < 1 || n > 64 {
			return fmt.Errorf("parallelism must be between 1 and 64, got %d", n)
		}
		info, err := e.c.GetFileMetadata(ctx, GetFileMetadataRequest{FilePath: new(path)})
		if err != nil {
			return err
		}
		if info.LastModified != nil {
			modifiedAt = *info.LastModified
		}
		// A file that fits one range is downloaded whole: a range covering the
		// whole file may be answered with a 200 rather than a 206.
		if info.ContentLength != nil && *info.ContentLength > max(e.tun.downloadMinParallelSize, e.tun.downloadPartSize) {
			size = *info.ContentLength
			ranges = splitIntoRanges(size, e.tun.downloadPartSize)
			parallelism = n
		}
	}
	dest, err := newDownloadDestination(destination, size)
	if err != nil {
		return err
	}
	if err := e.walkDownloadTo(ctx, cfg, path, ranges, modifiedAt, dest, parallelism); err != nil {
		dest.discard()
		return err
	}
	return dest.finish()
}

// walkDownloadTo downloads every range from each endpoint in turn, moving on
// only while no byte has arrived.
func (e *engine) walkDownloadTo(ctx context.Context, cfg downloadConfig, path string, ranges []byteRange, modifiedAt string, dest *downloadDestination, parallelism int) error {
	health := e.endpointHealth(false)
	order := health.order(e.clock(), e.transferEndpoints(ctx))
	var walkedPast []endpoint
	for i, ep := range order.endpoints {
		dc := e.newDownloadContext(cfg, path, ep, i == len(order.endpoints)-1, parallelism)
		err := e.downloadRanges(ctx, dc, ranges, modifiedAt, dest, parallelism)
		if err == nil {
			e.transferDone(ctx, health, order, walkedPast, order.endpoints[:i+1])
			return nil
		}
		if !mayWalkDownload(ctx, dc, err) {
			return err
		}
		e.c.logger.InfoContext(ctx, "download failed before any data arrived, trying the next endpoint", "endpoint", ep.String(), "error", err)
		if dc.countsAgainst(ctx, err) {
			walkedPast = append(walkedPast, ep)
		}
	}
	panic("unreachable: the last endpoint returns")
}

// downloadRanges downloads the ranges from one endpoint into dest, on up to
// parallelism goroutines. The first range's first bytes arrive before any
// other range starts, so that a failure before the endpoint is committed is
// the only request in flight, and the walk can tell where it stands.
func (e *engine) downloadRanges(ctx context.Context, dc *downloadContext, ranges []byteRange, modifiedAt string, dest *downloadDestination, parallelism int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	first, err := e.openDownload(ctx, dc, ranges[0], modifiedAt)
	if err != nil {
		return err
	}
	buf := make([]byte, 32<<10)
	n, err := first.Read(buf)
	if n > 0 {
		dc.committed.Store(true)
		if werr := dest.writeAt(ranges[0].start, buf[:n]); werr != nil {
			_ = first.Close()
			return werr
		}
	}
	if err != nil && err != io.EOF {
		_ = first.Close()
		return err
	}

	var mu sync.Mutex
	var firstErr error
	fail := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
		cancel()
	}
	// Every job runs, so that the first range is closed and a cancellation
	// between jobs is recorded rather than leaving ranges unwritten.
	work := make(chan func(), len(ranges))
	work <- func() {
		defer first.Close()
		if err == io.EOF {
			return
		}
		if err := copyRange(dest, first, ranges[0].start+int64(n)); err != nil {
			fail(err)
		}
	}
	for _, rng := range ranges[1:] {
		work <- func() {
			if err := ctx.Err(); err != nil {
				fail(err)
				return
			}
			d, err := e.openDownload(ctx, dc, rng, modifiedAt)
			if err != nil {
				fail(err)
				return
			}
			defer d.Close()
			if err := copyRange(dest, d, rng.start); err != nil {
				fail(err)
			}
		}
	}
	close(work)
	var wg sync.WaitGroup
	for range min(parallelism, len(ranges)) {
		wg.Go(func() {
			for job := range work {
				job()
			}
		})
	}
	wg.Wait()
	return firstErr
}

// copyRange writes what d reads into dest from offset on.
func copyRange(dest *downloadDestination, d io.Reader, offset int64) error {
	buf := make([]byte, 32<<10)
	for {
		n, err := d.Read(buf)
		if n > 0 {
			if werr := dest.writeAt(offset, buf[:n]); werr != nil {
				return werr
			}
			offset += int64(n)
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func splitIntoRanges(size, partSize int64) []byteRange {
	var ranges []byteRange
	for start := int64(0); start < size; start += partSize {
		ranges = append(ranges, byteRange{start: start, end: min(start+partSize, size) - 1})
	}
	return ranges
}

// takesRanges reports whether ranges can be written at any offset of
// destination: a new file, or a regular one after following symlinks.
func takesRanges(destination string) bool {
	info, err := os.Stat(destination)
	if errors.Is(err, fs.ErrNotExist) {
		return true
	}
	return err == nil && info.Mode().IsRegular()
}

// downloadDestination is where downloadTo writes.
type downloadDestination struct {
	destination string
	// temp is the temp file a new destination is written to, or "".
	temp string
	// size is the length an existing destination is given on its first write,
	// so ranges can be written at any offset, or -1.
	size int64

	mu sync.Mutex
	f  *os.File
	// position is the file's own offset, which only an in-order Write moves.
	position int64
}

func newDownloadDestination(destination string, size int64) (*downloadDestination, error) {
	// Lstat, so that a broken symlink is written through rather than replaced.
	if _, err := os.Lstat(destination); err == nil {
		return &downloadDestination{destination: destination, size: size}, nil
	}
	dir := filepath.Dir(destination)
	for {
		var suffix [8]byte
		_, _ = rand.Read(suffix[:])
		temp := filepath.Join(dir, ".download-"+hex.EncodeToString(suffix[:]))
		// Mode 0o666, so the process umask applies, as for any new file.
		f, err := os.OpenFile(temp, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
		if errors.Is(err, fs.ErrExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		d := &downloadDestination{destination: destination, temp: temp, size: -1, f: f}
		if size >= 0 {
			if err := f.Truncate(size); err != nil {
				d.discard()
				return nil, err
			}
		}
		return d, nil
	}
}

// writeAt writes p at offset. It is safe for concurrent use.
func (d *downloadDestination) writeAt(offset int64, p []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.f == nil {
		if err := d.open(); err != nil {
			return err
		}
	}
	// Only ranges write out of order, so a pipe is only appended to.
	if offset != d.position {
		_, err := d.f.WriteAt(p, offset)
		return err
	}
	_, err := d.f.Write(p)
	d.position += int64(len(p))
	return err
}

// open opens an existing destination on the first byte, so a download that
// fails before any data leaves it as it was.
func (d *downloadDestination) open() error {
	f, err := os.OpenFile(d.destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	if d.size >= 0 {
		if err := f.Truncate(d.size); err != nil {
			_ = f.Close()
			return err
		}
	}
	d.f = f
	return nil
}

// finish closes the destination, moving a temp file into place. An empty
// download still empties an existing file.
func (d *downloadDestination) finish() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.f == nil {
		if err := d.open(); err != nil {
			return err
		}
	}
	if err := d.f.Close(); err != nil {
		d.f = nil
		d.removeTemp()
		return err
	}
	d.f = nil
	if d.temp != "" {
		// Replaces a file that appeared meanwhile, and fails for a directory.
		if err := os.Rename(d.temp, d.destination); err != nil {
			d.removeTemp()
			return err
		}
	}
	return nil
}

// discard closes the destination after a failure: a temp file is removed,
// and an existing destination keeps what was written.
func (d *downloadDestination) discard() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.f != nil {
		_ = d.f.Close()
		d.f = nil
	}
	d.removeTemp()
}

func (d *downloadDestination) removeTemp() {
	if d.temp != "" {
		_ = os.Remove(d.temp)
	}
}
