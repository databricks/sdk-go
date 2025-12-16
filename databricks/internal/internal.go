package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
)

type ClientOptions struct {
	HTTPClient  *http.Client
	Credentials auth.Credentials
	Timeout     time.Duration
	Logger      *slog.Logger
}

// Initialize initializes the client options and validates its configuration.
func (o *ClientOptions) Initialize() error {
	if o.Logger == nil {
		o.Logger = slog.New(noopHandler{})
	}
	return nil
}

// nopHandler is a handler that does nothing.
type noopHandler struct{}

func (h noopHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (h noopHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h noopHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h noopHandler) WithGroup(_ string) slog.Handler {
	return h
}
