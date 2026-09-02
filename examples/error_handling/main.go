//go:build examples

package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/jobs/v2"
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

	// Job IDs are positive, so this request deterministically returns an API error.
	_, err = c.GetJob(ctx, jobs.GetJobRequest{JobId: new(int64(-1))})
	if err == nil {
		return errors.New("expected the invalid job lookup to fail")
	}

	// APIError preserves the structured Databricks error response, and
	// errors.As also finds it through wrapped errors.
	var apiError *apierr.APIError
	if !errors.As(err, &apiError) {
		return fmt.Errorf("getting invalid job returned an unexpected error type: %w", err)
	}
	fmt.Printf("Code: %s\n", apiError.Code())
	fmt.Printf("HTTP status: %d\n", apiError.HTTPStatusCode())
	fmt.Printf("Message: %s\n", apiError.Message())
	return nil
}
