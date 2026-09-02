//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/jobs/v2"
	"github.com/databricks/sdk-go/options/client"
)

type externalCredentials struct {
	token func(context.Context) (string, error)
}

const workspaceHost = "https://<workspace-host>"

func (externalCredentials) Name() string {
	// Name gives SDK logs and telemetry a stable label for this auth type.
	return "external"
}

func (c externalCredentials) AuthHeaders(ctx context.Context) ([]auth.Header, error) {
	// The SDK calls AuthHeaders when signing a request, so an implementation can
	// retrieve or refresh tokens from an external provider here.
	token, err := c.token(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting external token: %w", err)
	}
	return []auth.Header{{Key: "Authorization", Value: "Bearer " + token}}, nil
}

func main() {
	creds := externalCredentials{token: func(context.Context) (string, error) {
		// Replace this environment lookup with a call to your external token provider.
		token := os.Getenv("EXTERNAL_ACCESS_TOKEN")
		if token == "" {
			return "", fmt.Errorf("EXTERNAL_ACCESS_TOKEN must be set")
		}
		return token, nil
	}}

	ctx := context.Background()
	c, err := jobs.NewClient(ctx, client.WithHost(workspaceHost), client.WithCredentials(creds))
	if err != nil {
		log.Fatalf("creating Jobs client: %v", err)
	}
	// Use a small read request to confirm that the supplied credentials work.
	if _, err := c.ListJobs(ctx, jobs.ListJobsRequest{Limit: new(1)}); err != nil {
		log.Fatalf("listing jobs: %v", err)
	}

	fmt.Printf("Authenticated successfully with %s credentials\n", creds.Name())
}
