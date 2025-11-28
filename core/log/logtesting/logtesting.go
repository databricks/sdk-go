// Package logtesting provides a thread-safe logger that captures log records
// for testing. It is intended for internal use only.
package logtesting

import (
	"context"
	"sync"

	"github.com/databricks/sdk-go/core/log"
)

// Record represents a log record captured by the logger.
type Record struct {
	Context context.Context
	Level   log.Level
	Message string
	Args    []any
}

// Logger is a thread-safe logger that captures log records for testing.
type Logger struct {
	records []Record
	mu      sync.Mutex
}

func (l *Logger) Log(ctx context.Context, level log.Level, format string, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, Record{
		Context: ctx,
		Level:   level,
		Message: format,
		Args:    a,
	})
}

// Records returns a copy of all log records captured by the logger.
func (l *Logger) Records() []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]Record{}, l.records...)
}
