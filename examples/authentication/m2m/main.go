//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/auth/credentials"
	"github.com/databricks/sdk-go/jobs/v2"
	"github.com/databricks/sdk-go/options/client"
)

func main() {
	host := os.Getenv("DATABRICKS_HOST")
	clientID := os.Getenv("DATABRICKS_CLIENT_ID")
	clientSecret := os.Getenv("DATABRICKS_CLIENT_SECRET")
	if host == "" || clientID == "" || clientSecret == "" {
		log.Fatal("DATABRICKS_HOST, DATABRICKS_CLIENT_ID, and DATABRICKS_CLIENT_SECRET must be set")
	}

	// Passing these credentials to the client below ensures it uses M2M, even
	// if other authentication settings are available.
	tokenProvider, err := credentials.NewM2MCredentials(credentials.M2MOptions{
		Host:         host,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Fatalf("creating M2M credentials: %v", err)
	}
	// Cache the provider so access tokens are reused until they near expiry.
	creds := auth.NewTokenCredentials("oauth-m2m", auth.NewCachedTokenProvider(tokenProvider))

	ctx := context.Background()
	c, err := jobs.NewClient(ctx, client.WithHost(host), client.WithCredentials(creds))
	if err != nil {
		log.Fatalf("creating Jobs client: %v", err)
	}
	if _, err := c.ListJobs(ctx, jobs.ListJobsRequest{Limit: new(1)}); err != nil {
		log.Fatalf("listing jobs: %v", err)
	}

	fmt.Println("Authenticated successfully with explicit M2M credentials")
}
