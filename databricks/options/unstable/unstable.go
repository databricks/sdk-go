// Package unstable contains functionalities that need to be shared
// between SDK packages but are not considered part of the public API
// of the Databricks SDKs.
//
// IMPORTANT: Clients should not directly depend on functionalities in
// this package. All functionalities in this package are considered as
// UNSTABLE and may change without notice.
package unstable

import (
	"log/slog"

	"github.com/databricks/sdk-go/databricks/internal"
	"github.com/databricks/sdk-go/databricks/options"
)

// ResolvedOptions contains the resolved configuration values from ClientOptions.
type ResolvedOptions struct {
	Logger *slog.Logger
	Host   string
}

// Resolve applies all ClientOptions and returns the resolved configuration.
func Resolve(opts ...options.ClientOption) (ResolvedOptions, error) {
	o := &internal.ClientOptions{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return ResolvedOptions{}, err
		}
	}
	return ResolvedOptions{
		Logger: o.Logger,
		Host:   o.Host,
	}, nil
}
