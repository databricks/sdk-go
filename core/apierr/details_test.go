package apierr

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseErrorDetails(t *testing.T) {
	testCases := []struct {
		name       string
		rawDetails []json.RawMessage
		want       ErrorDetails
	}{
		{
			name:       "empty details",
			rawDetails: []json.RawMessage{},
			want:       ErrorDetails{},
		},
		{
			name: "all known error details types",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.ErrorInfo",
					"reason": "reason",
					"domain": "domain",
					"metadata": {"k1": "v1", "k2": "v2"}
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.RequestInfo",
					"request_id": "req42",
					"serving_data": "data"
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.RetryInfo",
					"retry_delay": "42.000000001s"
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.DebugInfo",
					"stack_entries": ["entry1", "entry2"],
					"detail": "detail"
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.QuotaFailure",
					"violations": [{"subject": "subject", "description": "description"}]
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.PreconditionFailure",
					"violations": [{"type": "type", "subject": "subject", "description": "description"}]
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.BadRequest",
					"field_violations": [{"field": "field", "description": "description"}]
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.ResourceInfo",
					"resource_type": "resource_type",
					"resource_name": "resource_name",
					"owner": "owner",
					"description": "description"
				}`),
				json.RawMessage(`{
					"@type": "type.googleapis.com/google.rpc.Help",
					"links": [{"description": "description", "url": "url"}]
				}`),
			},
			want: ErrorDetails{
				ErrorInfo: &ErrorInfo{
					Reason:   "reason",
					Domain:   "domain",
					Metadata: map[string]string{"k1": "v1", "k2": "v2"},
				},
				RequestInfo: &RequestInfo{
					RequestID:   "req42",
					ServingData: "data",
				},
				RetryInfo: &RetryInfo{
					RetryDelay: 42*time.Second + time.Nanosecond,
				},
				DebugInfo: &DebugInfo{
					StackEntries: []string{"entry1", "entry2"},
					Detail:       "detail",
				},
				QuotaFailure: &QuotaFailure{
					Violations: []QuotaFailureViolation{{Subject: "subject", Description: "description"}},
				},
				PreconditionFailure: &PreconditionFailure{
					Violations: []PreconditionFailureViolation{{Type: "type", Subject: "subject", Description: "description"}},
				},
				BadRequest: &BadRequest{
					FieldViolations: []BadRequestFieldViolation{{Field: "field", Description: "description"}},
				},
				ResourceInfo: &ResourceInfo{
					ResourceType: "resource_type",
					ResourceName: "resource_name",
					Owner:        "owner",
					Description:  "description",
				},
				Help: &Help{
					Links: []HelpLink{{Description: "description", URL: "url"}},
				},
			},
		},
		{
			name: "unknown error details type",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{"@type": "foo", "reason": "reason"}`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					map[string]any{
						"@type":  "foo",
						"reason": "reason",
					},
				},
			},
		},
		{
			name: "invalid error details - not valid JSON",
			rawDetails: []json.RawMessage{
				json.RawMessage(`not valid json`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					[]byte(`not valid json`),
				},
			},
		},
		{
			name: "invalid error details - not an object",
			rawDetails: []json.RawMessage{
				json.RawMessage(`42`),
				json.RawMessage(`"foobar"`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					42.0,
					"foobar",
				},
			},
		},
		{
			name: "invalid error details - object without @type",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{"foo": "bar"}`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					map[string]any{
						"foo": "bar",
					},
				},
			},
		},
		{
			name: "invalid error details - known type with invalid fields",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.RequestInfo", "request_id": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retry_delay": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.DebugInfo", "stack_entries": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.QuotaFailure", "violations": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.PreconditionFailure", "violations": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.BadRequest", "field_violations": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.ResourceInfo", "resource_type": 0}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.Help", "links": 0}`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					map[string]any{
						"@type":  "type.googleapis.com/google.rpc.ErrorInfo",
						"reason": 0.0,
					},
					map[string]any{
						"@type":      "type.googleapis.com/google.rpc.RequestInfo",
						"request_id": 0.0,
					},
					map[string]any{
						"@type":       "type.googleapis.com/google.rpc.RetryInfo",
						"retry_delay": 0.0,
					},
					map[string]any{
						"@type":         "type.googleapis.com/google.rpc.DebugInfo",
						"stack_entries": 0.0,
					},
					map[string]any{
						"@type":      "type.googleapis.com/google.rpc.QuotaFailure",
						"violations": 0.0,
					},
					map[string]any{
						"@type":      "type.googleapis.com/google.rpc.PreconditionFailure",
						"violations": 0.0,
					},
					map[string]any{
						"@type":            "type.googleapis.com/google.rpc.BadRequest",
						"field_violations": 0.0,
					},
					map[string]any{
						"@type":         "type.googleapis.com/google.rpc.ResourceInfo",
						"resource_type": 0.0,
					},
					map[string]any{
						"@type": "type.googleapis.com/google.rpc.Help",
						"links": 0.0,
					},
				},
			},
		},
		{
			name: "only keep the last error details of a type",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": "first"}`),
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": "second"}`),
			},
			want: ErrorDetails{
				ErrorInfo: &ErrorInfo{
					Reason: "second",
				},
			},
		},
		{
			name: "RetryInfo with invalid duration format",
			rawDetails: []json.RawMessage{
				json.RawMessage(`{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retry_delay": "invalid"}`),
			},
			want: ErrorDetails{
				UnknownDetails: []any{
					map[string]any{
						"@type":       "type.googleapis.com/google.rpc.RetryInfo",
						"retry_delay": "invalid",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseErrorDetails(tc.rawDetails)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected error details (-want +got):\n%s", diff)
			}
		})
	}
}
