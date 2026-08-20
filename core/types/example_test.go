package types_test

import (
	"fmt"
	"time"

	"github.com/databricks/sdk-go/core/types"
)

func ExampleNewFromDuration() {
	requestTimeout := types.NewFromDuration(15 * time.Minute)
	fmt.Println(requestTimeout)
	// Output:
	// 900s
}

func ExampleNewFromTime() {
	requestStartTime := types.NewFromTime(time.Date(
		2024, time.January, 15, 11, 30, 0, 0,
		time.FixedZone("UTC+1", 60*60),
	))
	fmt.Println(requestStartTime)
	// Output:
	// 2024-01-15T10:30:00Z
}

func ExampleTime_Add() {
	createdAt := types.NewFromTime(time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC))
	ttl := &types.Duration{Seconds: 3_600, Nanos: 500_000_000}
	expiresAt := createdAt.Add(ttl)

	fmt.Println(expiresAt)
	// Output:
	// 2024-01-15T11:30:00.500Z
}

func ExampleTime_comparisonUsingStandardLibrary() {
	updatedAt := types.NewFromTime(time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC))
	expiresAt := types.NewFromTime(time.Date(2024, time.January, 15, 11, 30, 0, 0, time.UTC))

	fmt.Println(updatedAt.AsTime().Before(expiresAt.AsTime()))
	// Output:
	// true
}

func ExampleTime_elapsedUsingStandardLibrary() {
	startedAt := types.NewFromTime(time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC))
	finishedAt := types.NewFromTime(time.Date(2024, time.January, 15, 10, 32, 30, 0, time.UTC))

	fmt.Println(finishedAt.AsTime().Sub(startedAt.AsTime()))
	// Output:
	// 2m30s
}
