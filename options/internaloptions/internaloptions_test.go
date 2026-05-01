package internaloptions

import (
	"context"
	"log/slog"
	"testing"
)

func TestClientOptionsInitialize_DefaultLogger(t *testing.T) {
	c := &ClientOptions{}
	if err := c.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	if c.Logger == nil {
		t.Fatal("expected default logger to be set")
	}
	if c.Logger.Enabled(context.Background(), slog.LevelError) {
		t.Fatal("expected default logger to be disabled")
	}
}

func TestClientOptionsInitialize_PreservesProvidedLogger(t *testing.T) {
	provided := slog.Default()
	c := &ClientOptions{Logger: provided}
	if err := c.Initialize(); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	if c.Logger != provided {
		t.Fatal("expected provided logger to be preserved")
	}
}
