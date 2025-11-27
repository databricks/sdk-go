// Package log provides the Databricks Logger interface and a default logger
// implementation.
//
// The Logger interface is designed to log messages at different levels of
// severity. When creating a logger, you can set a minimum log level, meaning
// any messages below this level will be ignored.
package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Level represents the importance or severity of a log event. The higher the
// level, the more important or severe that event is.
type Level int

const (
	LevelTrace Level = -8
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is anything that can log a message at a given level.
type Logger interface {
	Log(ctx context.Context, level Level, format string, a ...any)
}

// NewLogger returns a new default logger that logs on [os.Stderr].
func NewLogger() Logger {
	return &logger{out: os.Stderr}
}

type logger struct {
	out io.Writer
}

func (l *logger) Log(_ context.Context, level Level, format string, a ...any) {
	fmt.Fprintf(l.out, "[%s] %s\n", level.String(), fmt.Sprintf(format, a...))
}
