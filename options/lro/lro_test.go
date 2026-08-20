package lro

import (
	"testing"
	"time"

	"github.com/databricks/sdk-go/options/internaloptions"
)

func TestWithTimeout(t *testing.T) {
	testCases := []struct {
		name        string
		timeouts    []time.Duration
		wantTimeout time.Duration
	}{
		{
			name: "unset",
		},
		{
			name:        "positive",
			timeouts:    []time.Duration{3 * time.Second},
			wantTimeout: 3 * time.Second,
		},
		{
			name:        "later positive overrides earlier timeout",
			timeouts:    []time.Duration{3 * time.Second, 5 * time.Second},
			wantTimeout: 5 * time.Second,
		},
		{
			name:     "zero clears earlier timeout",
			timeouts: []time.Duration{3 * time.Second, 0},
		},
		{
			name:        "negative",
			timeouts:    []time.Duration{-time.Second},
			wantTimeout: -time.Second,
		},
		{
			name:        "positive overrides earlier negative timeout",
			timeouts:    []time.Duration{-time.Second, 5 * time.Second},
			wantTimeout: 5 * time.Second,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			options := internaloptions.LROOptions{}
			for _, timeout := range testCase.timeouts {
				if err := WithTimeout(timeout)(&options); err != nil {
					t.Fatalf("WithTimeout(): %v", err)
				}
			}

			if options.Timeout != testCase.wantTimeout {
				t.Errorf("Timeout = %v, want %v", options.Timeout, testCase.wantTimeout)
			}
		})
	}
}
