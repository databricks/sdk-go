package log

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestLevel_order(t *testing.T) {
	var wantOrder = []Level{
		LevelTrace,
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
	}

	for i := 1; i < len(wantOrder); i++ {
		if wantOrder[i-1] >= wantOrder[i] {
			t.Errorf("Level %s should be lower than %s", wantOrder[i-1].String(), wantOrder[i].String())
		}
	}
}

func TestLogger_Log(t *testing.T) {
	testCases := []struct {
		level  Level
		format string
		args   []any
		want   string
	}{
		{
			level:  LevelTrace,
			format: "test %s",
			args:   []any{"message"},
			want:   "[TRACE] test message\n",
		},
		{
			level:  LevelDebug,
			format: "test %s",
			args:   []any{"message"},
			want:   "[DEBUG] test message\n",
		},
		{
			level:  LevelInfo,
			format: "test %s",
			args:   []any{"message"},
			want:   "[INFO] test message\n",
		},
		{
			level:  LevelWarn,
			format: "test %s",
			args:   []any{"message"},
			want:   "[WARN] test message\n",
		},
		{
			level:  LevelError,
			format: "test %s",
			args:   []any{"message"},
			want:   "[ERROR] test message\n",
		},
		{
			level:  Level(42), // unknown level
			format: "test %s",
			args:   []any{"message"},
			want:   "[UNKNOWN] test message\n",
		},
	}

	for _, tc := range testCases {
		out := &bytes.Buffer{}
		l := logger{out: out}

		l.Log(context.Background(), tc.level, tc.format, tc.args...)

		if got := out.String(); got != tc.want {
			t.Errorf("Log(): got %s, want %s", got, tc.want)
		}
	}
}

func TestNew(t *testing.T) {
	l := NewLogger().(*logger)
	if l.out != os.Stderr {
		t.Errorf("NewLogger() got out %v, want %v", l.out, os.Stderr)
	}
}
