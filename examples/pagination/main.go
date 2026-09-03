//go:build examples

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/databricks/sdk-go/jobs/v2"
)

const (
	pageSize   = 5
	jobsToList = 12
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	c, err := jobs.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating Jobs client: %w", err)
	}

	count := 0
	// Limit controls the page size. The iterator fetches more results as needed.
	for job, err := range c.ListJobsIter(ctx, jobs.ListJobsRequest{Limit: new(pageSize)}) {
		if err != nil {
			return fmt.Errorf("listing jobs: %w", err)
		}
		count++
		if job.JobId == nil {
			return fmt.Errorf("job %d omitted its ID", count)
		}
		if job.Settings != nil && job.Settings.Name != nil {
			fmt.Printf("%d: %s\n", *job.JobId, *job.Settings.Name)
		} else {
			fmt.Println(*job.JobId)
		}
		if count == jobsToList {
			break
		}
	}

	fmt.Printf("Listed %d jobs with a page size of %d.\n", count, pageSize)
	return nil
}
