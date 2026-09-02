//go:build examples

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/databricks/sdk-go/core/types"
)

func main() {
	// Use these constructors when assigning time values to SDK request fields;
	// use AsTime and AsDuration when consuming values from SDK responses.
	startedAt := types.NewFromTime(time.Date(2026, time.January, 15, 10, 30, 0, 0, time.UTC))
	timeout := types.NewFromDuration(15 * time.Minute)
	finishedAt := startedAt.Add(timeout)

	// Arithmetic does not validate the result; check it before sending it to an API.
	if err := finishedAt.CheckValid(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Started: %s\n", startedAt)
	fmt.Printf("Finished: %s\n", finishedAt)
	fmt.Printf("Elapsed: %s\n", finishedAt.AsTime().Sub(startedAt.AsTime()))
	fmt.Printf("Timeout as time.Duration: %s\n", timeout.AsDuration())
}
