//go:build examples

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/databricks/sdk-go/jobs/v2"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	// DATABRICKS_CONFIG_PROFILE selects a profile from ~/.databrickscfg. If it
	// is unset, the SDK uses default_profile from [__settings__], then falls
	// back to [DEFAULT]. DATABRICKS_* values override values from the profile.
	c, err := jobs.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating Jobs client: %w", err)
	}
	// Use a small read request to confirm that the resolved credentials work.
	if _, err := c.ListJobs(ctx, jobs.ListJobsRequest{Limit: new(1)}); err != nil {
		return fmt.Errorf("listing jobs: %w", err)
	}

	fmt.Println("Authenticated successfully")
	return nil
}
