package main

import (
	"fmt"
	"time"

	"github.com/databricks/sdk-go/core/types"
)

// Generated SDK models use pointer fields with these shapes. Keeping the
// example local to core/types lets it run before any particular service module
// is selected.
type apiResponse struct {
	CreateTime *types.Time
	UpdateTime *types.Time
	TTL        *types.Duration
}

type apiRequest struct {
	NotBefore *types.Time
	Retention *types.Duration
}

func main() {
	now := time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC)
	requestTimeout := types.NewFromDuration(15 * time.Minute)
	requestStartTime := types.NewFromTime(now)
	fmt.Printf("request timeout: %s, starts: %s\n", requestTimeout, requestStartTime)

	// A response value can be manipulated without converting through the
	// range-limited time.Duration type.
	response := &apiResponse{
		CreateTime: types.NewFromTime(now),
		UpdateTime: types.NewFromTime(now.Add(30 * time.Minute)),
		TTL:        &types.Duration{Seconds: 86_400, Nanos: 500_000_000},
	}
	expiresAt := response.CreateTime.Add(response.TTL)
	fmt.Printf("resource expires: %s\n", expiresAt)

	// SDK values can be forwarded directly into another generated request.
	request := &apiRequest{
		NotBefore: response.CreateTime,
		Retention: response.TTL,
	}
	fmt.Printf("forwarded request: %s, retention %s\n", request.NotBefore, request.Retention)

	// Convert to time.Time for standard-library comparison and subtraction.
	fmt.Printf("expired: %t\n", expiresAt.AsTime().Before(now))
	fmt.Printf("updated before expiry: %t\n", response.UpdateTime.AsTime().Before(expiresAt.AsTime()))
	fmt.Printf("elapsed: %s\n", expiresAt.AsTime().Sub(response.CreateTime.AsTime()))

	// AsDuration clamps values outside time.Duration's range. Round-tripping is
	// the current way to determine whether a conversion preserved the value.
	retention := &types.Duration{Seconds: types.MaxDurationSeconds}
	converted := retention.AsDuration()
	roundTripped := types.NewFromDuration(converted)
	exact := retention.Seconds == roundTripped.Seconds && retention.Nanos == roundTripped.Nanos
	fmt.Printf("stdlib conversion exact: %t\n", exact)
}
