//go:build examples

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/databricks/sdk-go/jobs/v2"
	"github.com/databricks/sdk-go/options/client"
)

func main() {
	profile := flag.String("profile", "", "Databricks configuration profile")
	flag.Parse()
	if *profile == "" {
		log.Fatal("-profile must be set")
	}

	if err := run(context.Background(), *profile); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, profile string) error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	// WithProfile selects a profile in code instead of through
	// DATABRICKS_CONFIG_PROFILE. WithLogger enables SDK logs, and WithTimeout
	// sets the default timeout for API calls made by this client.
	c, err := jobs.NewClient(ctx,
		client.WithProfile(profile),
		client.WithLogger(logger),
		client.WithTimeout(30*time.Second),
	)
	if err != nil {
		return fmt.Errorf("creating Jobs client: %w", err)
	}
	if _, err := c.ListJobs(ctx, jobs.ListJobsRequest{Limit: new(1)}); err != nil {
		return fmt.Errorf("listing jobs: %w", err)
	}

	fmt.Printf("Configured client with profile %q\n", profile)
	return nil
}
