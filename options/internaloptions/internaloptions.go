// Package internaloptions contains the resolved options types produced
// by applying [github.com/databricks/sdk-go/options/client.Option] and
// [github.com/databricks/sdk-go/options/call.Option] values.
//
// IMPORTANT: This package is NOT part of the public API of the Databricks
// SDK. Its contents may change at any time without notice. Clients should
// not directly depend on functionalities in this package.
package internaloptions

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/ops"
)

// ClientOptions is the resolved client configuration produced by applying
// client.Option values.
type ClientOptions struct {
	Host        string
	HTTPClient  *http.Client
	Credentials auth.Credentials
	Timeout     time.Duration
	Logger      *slog.Logger
}

// Initialize fills in defaults and validates the resolved client options.
func (c *ClientOptions) Initialize() error {
	if c.Logger == nil {
		c.Logger = slog.New(slog.DiscardHandler)
	}
	return nil
}

// CallOptions is the resolved per-call configuration produced by applying
// call.Option values.
type CallOptions struct {
	Retrier     func() ops.Retrier
	RateLimiter ops.Limiter
	Timeout     time.Duration
}
