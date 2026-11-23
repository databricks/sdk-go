package files

import (
	"bytes"
	"context"
	"io"
	"sync"
)

// uploadData is the source of one upload. The walk sends it from the start to
// each endpoint it tries.
type uploadData struct {
	// r reads the data in order, from the start of the current endpoint's
	// attempt.
	r io.Reader
	// seeker rewinds r for the next endpoint. It is set when the size was
	// learned by seeking.
	seeker io.Seeker
	// ra reads parts with positioned reads, without buffering them. It is set
	// for a source of known size that has random access.
	ra io.ReaderAt
	// stream is set for a one-pass stream, which r reads through.
	stream *replayStream
	// knownSize is the size of a seekable source.
	knownSize int64
	// sent, when set, is closed once net/http has closed the last request body
	// that read r itself; restart waits for it.
	sent <-chan struct{}
}

// newSeekableData returns the data of r, of the given size, which r was sized
// by seeking and is positioned at the start of.
func newSeekableData(r io.ReadSeeker, size int64) *uploadData {
	d := &uploadData{r: r, seeker: r, knownSize: size}
	if ra, ok := r.(io.ReaderAt); ok {
		d.ra = ra
	}
	return d
}

// newStreamData returns the data of a one-pass stream.
func newStreamData(r io.Reader) *uploadData {
	s := &replayStream{src: r, keeping: true}
	return &uploadData{r: s, stream: s, knownSize: -1}
}

// size is the size of the data, or -1 while a stream has not ended.
func (d *uploadData) size() int64 {
	if d.stream != nil {
		return d.stream.size()
	}
	return d.knownSize
}

// readAhead reads up to n bytes of a stream ahead, so that a shorter one
// learns its size. A seekable source knows it already.
func (d *uploadData) readAhead(n int64) error {
	if d.stream == nil {
		return nil
	}
	return d.stream.readAhead(n)
}

// canSendAgain reports whether the data can be sent from the start to another
// endpoint.
func (d *uploadData) canSendAgain() bool {
	return d.stream == nil || d.stream.canRestart()
}

// sending records that a request body reading r itself may be read until done
// is closed. RoundTrip always closes the body, but may do so on another
// goroutine after it returns, so the bytes it still sends would race the next
// endpoint's reads of r. It is a no-op on nil data.
func (d *uploadData) sending(done <-chan struct{}) {
	if d != nil {
		d.sent = done
	}
}

// restart positions the data at the start for the next endpoint, once no
// request of the last endpoint reads r any more.
func (d *uploadData) restart(ctx context.Context) error {
	if d.sent != nil {
		select {
		case <-d.sent:
		case <-ctx.Done():
			return context.Cause(ctx)
		}
		d.sent = nil
	}
	if d.stream != nil {
		d.stream.restart()
		return nil
	}
	_, err := d.seeker.Seek(0, io.SeekStart)
	return err
}

// whole returns the data for a single-shot PUT. A seekable source is the body
// itself, so the PUT can rewind it to retry; a stream is a bytes.Reader only
// if it ended within the bytes it kept.
func (d *uploadData) whole() io.Reader {
	if d.stream != nil {
		return d.stream.whole()
	}
	return d.r
}

// commit is called when the first part of the upload lands: it will not be
// sent again from the start, so a stream stops keeping what it reads.
func (d *uploadData) commit() {
	if d.stream != nil {
		d.stream.commit()
	}
}

// replayStream reads a one-pass stream and keeps every byte it reads, from the
// start, until the upload commits, so that the walk can send the stream again
// to the next endpoint. A multipart upload commits on its first part, so what
// it keeps is bounded by the read-ahead and one part.
type replayStream struct {
	src io.Reader

	mu sync.Mutex
	// kept holds the bytes read from src since the start. After commit it holds
	// only bytes read ahead but not yet handed out.
	kept []byte
	// pos is how many bytes of kept Read has handed out.
	pos     int
	total   int64 // bytes read from src
	keeping bool
	eof     bool
	// lost records that bytes were handed out without being kept, so the
	// stream can't start again.
	lost bool
}

func (s *replayStream) Read(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pos < len(s.kept) {
		n := copy(p, s.kept[s.pos:])
		s.pos += n
		if !s.keeping && s.pos == len(s.kept) {
			s.kept, s.pos = nil, 0
		}
		return n, nil
	}
	if s.eof {
		return 0, io.EOF
	}
	n, err := s.src.Read(p)
	s.total += int64(n)
	if n > 0 {
		if s.keeping {
			s.kept = append(s.kept, p[:n]...)
			s.pos += n
		} else {
			s.lost = true
		}
	}
	if err == io.EOF {
		s.eof = true
	}
	return n, err
}

func (s *replayStream) readAhead(n int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.keeping || s.eof || int64(len(s.kept)) >= n {
		return nil
	}
	buf := make([]byte, n-int64(len(s.kept)))
	m, err := io.ReadFull(s.src, buf)
	s.total += int64(m)
	s.kept = append(s.kept, buf[:m]...)
	switch err {
	case nil:
		return nil
	case io.EOF, io.ErrUnexpectedEOF:
		s.eof = true
		return nil
	}
	return err
}

func (s *replayStream) size() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.eof {
		return s.total
	}
	return -1
}

func (s *replayStream) canRestart() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.keeping && !s.lost
}

func (s *replayStream) restart() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pos = 0
}

func (s *replayStream) whole() io.Reader {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.eof && s.keeping {
		return bytes.NewReader(s.kept)
	}
	// The rest is streamed and can't be read again; keeping it would hold the
	// whole stream in memory.
	s.keeping, s.lost = false, true
	rest := io.MultiReader(bytes.NewReader(s.kept[s.pos:]), s.src)
	s.kept, s.pos = nil, 0
	return rest
}

func (s *replayStream) commit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.keeping {
		return
	}
	s.keeping = false
	s.kept = s.kept[s.pos:]
	s.pos = 0
}
