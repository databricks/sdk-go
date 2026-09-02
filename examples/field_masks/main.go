//go:build examples

package main

import (
	"fmt"
	"log"

	"github.com/databricks/sdk-go/core/types"
	"github.com/databricks/sdk-go/dataquality/v1"
)

func main() {
	// NewFieldMask validates these API field paths against Monitor and its nested
	// generated types. Only fields in the mask are updated.
	updateMask, err := types.NewFieldMask[dataquality.Monitor](
		"data_profiling_config.output_schema_id",
		// Oneof variants use their API field name directly in the mask.
		"data_profiling_config.time_series",
		"data_profiling_config.schedule",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Schedule is intentionally unset. Because its path is present in the mask,
	// this request clears the current schedule while updating the other fields.
	request := &dataquality.UpdateMonitorRequest{
		// These fields identify the monitor. UpdateMask applies to the Monitor
		// resource below, not to these request fields.
		ObjectType: new("table"),
		ObjectId:   new("00000000-0000-0000-0000-000000000000"),
		Monitor: &dataquality.Monitor{
			DataProfilingConfig: &dataquality.DataProfilingConfig{
				OutputSchemaId: new("main.monitoring"),
				AnalysisConfig: &dataquality.DataProfilingConfig_AnalysisConfig_TimeSeries{
					TimeSeries: dataquality.TimeSeriesConfig{
						TimestampColumn: new("event_time"),
					},
				},
			},
		},
		UpdateMask: updateMask,
	}

	config := request.Monitor.DataProfilingConfig
	fmt.Printf("Update %s %s with mask %s\n", *request.ObjectType, *request.ObjectId, request.UpdateMask)
	fmt.Printf("Output schema: %s\n", *config.OutputSchemaId)
	fmt.Printf("Clear schedule: %t\n", config.Schedule == nil)
}
