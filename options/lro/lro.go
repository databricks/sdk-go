// Package lro defines options used to wait for long-running operations.
package lro

import (
	"time"

	"github.com/databricks/sdk-go/options/internaloptions"
)

// Option configures a wait for a long-running operation.
type Option func(*internaloptions.LROOptions) error

// WithTimeout returns an Option that limits the complete wait, including every
// polling attempt. When the context already has a deadline, the earlier
// deadline applies. A zero duration removes a timeout set by an earlier Option.
// After all options are applied, a final negative duration returns an error
// before polling begins.
func WithTimeout(timeout time.Duration) Option {
	return func(options *internaloptions.LROOptions) error {
		options.Timeout = timeout
		return nil
	}
}
