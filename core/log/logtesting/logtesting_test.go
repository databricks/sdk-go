package logtesting

import (
	"context"
	"reflect"
	"testing"

	"github.com/databricks/sdk-go/core/log"
)

func TestLogger_Log(t *testing.T) {
	testContext := context.WithValue(context.Background(), "test-key", "test-value")

	testCases := []struct {
		desc        string
		wantRecords []Record
	}{
		{
			desc: "debug level with no args",
			wantRecords: []Record{
				{
					Context: testContext,
					Level:   log.LevelDebug,
					Message: "simple message",
					Args:    nil,
				},
			},
		},
		{
			desc: "info level with format args",
			wantRecords: []Record{
				{
					Context: testContext,
					Level:   log.LevelInfo,
					Message: "user %s logged in",
					Args:    []any{"alice"},
				},
			},
		},
		{
			desc: "error level with multiple args",
			wantRecords: []Record{
				{
					Context: testContext,
					Level:   log.LevelDebug,
					Message: "simple message",
					Args:    nil,
				},
			},
		},
		{
			desc: "multiple log records",
			wantRecords: []Record{
				{
					Context: testContext,
					Level:   log.LevelDebug,
					Message: "message 1",
				},
				{
					Context: testContext,
					Level:   log.LevelInfo,
					Message: "message 2",
				},
				{
					Context: testContext,
					Level:   log.LevelError,
					Message: "message 3",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			logger := &Logger{}

			for _, record := range tc.wantRecords {
				logger.Log(record.Context, record.Level, record.Message, record.Args...)
			}

			gotRecords := logger.Records()
			if !reflect.DeepEqual(gotRecords, tc.wantRecords) {
				t.Errorf("records: got %v, want %v", gotRecords, tc.wantRecords)
			}
		})
	}
}

func TestLogger_Records_ReturnsCopy(t *testing.T) {
	logger := &Logger{}
	logger.Log(context.Background(), log.LevelInfo, "original")

	records1 := logger.Records()
	records1[0].Message = "modified"

	records2 := logger.Records()
	if got := records2[0].Message; got != "original" {
		t.Errorf("message: got %q, want %q (Records should return a copy)", got, "original")
	}
}
