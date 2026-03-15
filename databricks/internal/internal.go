package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/auth/credentials"
)

type ClientOptions struct {
	Host        string
	HTTPClient  *http.Client
	Credentials auth.Credentials
	Timeout     time.Duration
	Logger      *slog.Logger

	// Profile and credentials resolution options.
	NoResolution bool   // disable profile resolution entirely
	NoEnv        bool   // disable environment variable overlay
	Profile      string // profile name to use for resolution
	ProfileFile  string // path to the profile file to use for resolution
}

// Initialize initializes the client options and validates its configuration.
func (o *ClientOptions) Initialize() error {
	if o.Logger == nil {
		o.Logger = slog.New(slog.DiscardHandler)
	}
	if err := o.resolveCredentials(); err != nil {
		return err
	}
	return nil
}

func (o *ClientOptions) resolveCredentials() error {
	if o.Credentials != nil || o.NoResolution {
		return nil
	}

	var opts []credentials.ResolveOption
	if o.Profile != "" {
		opts = append(opts, credentials.WithProfile(o.Profile))
	}
	if o.ProfileFile != "" {
		opts = append(opts, credentials.WithProfileFile(o.ProfileFile))
	}
	if o.NoEnv {
		opts = append(opts, credentials.WithoutEnv())
	}
	opts = append(opts, credentials.WithLogger(o.Logger))
	if o.HTTPClient != nil {
		opts = append(opts, credentials.WithHTTPClient(o.HTTPClient))
	}

	cred, err := credentials.ResolveCredentials(context.Background(), opts...)
	if err != nil {
		return err
	}
	o.Credentials = cred

	return nil
}
