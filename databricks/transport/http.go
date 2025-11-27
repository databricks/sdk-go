package transport

import (
	"context"
	"net/http"

	"github.com/databricks/sdk-go/databricks/internal"
	"github.com/databricks/sdk-go/databricks/options"
)

func NewHTTPClient(ctx context.Context, opts ...options.ClientOption) (*http.Client, error) {
	copts := &internal.ClientOptions{}
	for _, opt := range opts {
		if err := opt(copts); err != nil {
			return nil, err
		}
	}
	_ = copts // not used for now
	return http.DefaultClient, nil
}
