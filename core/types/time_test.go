package types_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/types"
	"github.com/google/go-cmp/cmp"
)

func TestTime_JSONRoundTrip(t *testing.T) {
	testCases := []struct {
		name  string
		value *types.Time
		json  string
	}{
		{name: "whole second", value: &types.Time{Seconds: 1_705_315_800}, json: `"2024-01-15T10:50:00Z"`},
		{name: "millisecond precision", value: &types.Time{Seconds: 1_705_315_800, Nanos: 123_000_000}, json: `"2024-01-15T10:50:00.123Z"`},
		{name: "microsecond precision", value: &types.Time{Seconds: 1_705_315_800, Nanos: 123_456_000}, json: `"2024-01-15T10:50:00.123456Z"`},
		{name: "nanosecond precision", value: &types.Time{Seconds: 1_705_315_800, Nanos: 123_456_789}, json: `"2024-01-15T10:50:00.123456789Z"`},
		{name: "minimum", value: &types.Time{Seconds: types.MinTimestampSeconds}, json: `"0001-01-01T00:00:00Z"`},
		{name: "maximum", value: &types.Time{Seconds: types.MaxTimestampSeconds, Nanos: 999_999_999}, json: `"9999-12-31T23:59:59.999999999Z"`},
		{name: "fractional pre epoch", value: &types.Time{Seconds: -1, Nanos: 500_000_000}, json: `"1969-12-31T23:59:59.500Z"`},
		{name: "one nanosecond", value: &types.Time{Nanos: 1}, json: `"1970-01-01T00:00:00.000000001Z"`},
		{name: "one microsecond", value: &types.Time{Nanos: 1_000}, json: `"1970-01-01T00:00:00.000001Z"`},
		{name: "millisecond plus nanosecond", value: &types.Time{Nanos: 1_000_001}, json: `"1970-01-01T00:00:00.001000001Z"`},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if got := string(encoded); got != tc.json {
				t.Errorf("Marshal() = %q, want %q", got, tc.json)
			}

			var got types.Time
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if diff := cmp.Diff(tc.value, &got); diff != "" {
				t.Errorf("round trip mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTime_String(t *testing.T) {
	testCases := []struct {
		name  string
		value *types.Time
		want  string
	}{
		{name: "nil", want: "<nil>"},
		{name: "zero", value: &types.Time{}, want: "1970-01-01T00:00:00Z"},
		{name: "fractional", value: &types.Time{Seconds: -1, Nanos: 500_000_000}, want: "1969-12-31T23:59:59.500Z"},
		{name: "invalid", value: &types.Time{Nanos: -1}, want: "Time{Seconds: 0, Nanos: -1}"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTime_UnmarshalJSON_acceptsRFC3339(t *testing.T) {
	testCases := []struct {
		input string
		want  *types.Time
	}{
		{input: `"2024-01-15T11:50:00.123+01:00"`, want: &types.Time{Seconds: 1_705_315_800, Nanos: 123_000_000}},
		{input: `"1970-01-01T23:59:00+23:59"`, want: &types.Time{}},
		{input: `"1969-12-31T00:01:00-23:59"`, want: &types.Time{}},
		{input: `"1970-01-01T00:00:00.1Z"`, want: &types.Time{Nanos: 100_000_000}},
		{input: `"2000-02-29T00:00:00Z"`, want: &types.Time{Seconds: 951_782_400}},
		{input: `"1970-01-01T00:00:00.000000001Z"`, want: &types.Time{Nanos: 1}},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var got types.Time
			if err := json.Unmarshal([]byte(tc.input), &got); err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if diff := cmp.Diff(tc.want, &got); diff != "" {
				t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTime_StandardLibraryConversionAndAdd(t *testing.T) {
	stdlib := time.Date(2024, 1, 15, 10, 30, 0, 0, time.FixedZone("offset", 3600))
	custom := types.NewFromTime(stdlib)
	if got := custom.AsTime(); !got.Equal(stdlib) {
		t.Errorf("AsTime() = %v, want instant %v", got, stdlib)
	}

	got := custom.Add(&types.Duration{Seconds: 3_600, Nanos: 500_000_000})
	want := types.NewFromTime(stdlib.Add(time.Hour + 500*time.Millisecond))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Add() mismatch (-want +got):\n%s", diff)
	}
}

func TestTime_nilBehavior(t *testing.T) {
	var value *types.Time
	if got := value.AsTime(); !got.Equal(time.Unix(0, 0)) {
		t.Errorf("nil.AsTime() = %v, want Unix epoch", got)
	}
	if got := value.Add(&types.Duration{Seconds: 1}); got != nil {
		t.Errorf("nil.Add() = %+v, want nil", got)
	}

	value = &types.Time{Seconds: 1, Nanos: 2}
	got := value.Add(nil)
	if got == value || *got != *value {
		t.Errorf("Add(nil) = %+v, want an equal copy", got)
	}
}

func TestTime_Add_normalizesWithoutClamping(t *testing.T) {
	testCases := []struct {
		name     string
		time     *types.Time
		duration *types.Duration
		want     *types.Time
		valid    bool
	}{
		{name: "carry", time: &types.Time{Seconds: 10, Nanos: 900_000_000}, duration: &types.Duration{Nanos: 200_000_000}, want: &types.Time{Seconds: 11, Nanos: 100_000_000}, valid: true},
		{name: "borrow", time: &types.Time{Seconds: 10, Nanos: 100_000_000}, duration: &types.Duration{Nanos: -200_000_000}, want: &types.Time{Seconds: 9, Nanos: 900_000_000}, valid: true},
		{name: "maximum plus one nanosecond", time: &types.Time{Seconds: types.MaxTimestampSeconds, Nanos: 999_999_999}, duration: &types.Duration{Nanos: 1}, want: &types.Time{Seconds: types.MaxTimestampSeconds + 1}, valid: false},
		{name: "minimum minus one nanosecond", time: &types.Time{Seconds: types.MinTimestampSeconds}, duration: &types.Duration{Nanos: -1}, want: &types.Time{Seconds: types.MinTimestampSeconds - 1, Nanos: 999_999_999}, valid: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalTime, originalDuration := *tc.time, *tc.duration
			got := tc.time.Add(tc.duration)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Add() mismatch (-want +got):\n%s", diff)
			}
			if (*got).CheckValid() == nil != tc.valid {
				t.Errorf("Add().CheckValid() valid = %t, want %t", got.CheckValid() == nil, tc.valid)
			}
			if *tc.time != originalTime || *tc.duration != originalDuration {
				t.Error("Add() modified an operand")
			}
			if !tc.valid {
				if _, err := json.Marshal(got); err == nil {
					t.Error("Marshal(Add()) error = nil for out-of-range result")
				}
			}
		})
	}
}

func TestTime_UnmarshalJSON_rejectsInvalidValue(t *testing.T) {
	testCases := []string{
		`"not-a-time"`,
		`"2024-01-15T10:50:00Z "`,
		`"2024-01-15t10:50:00Z"`,
		`"2024-01-15T10:50:00z"`,
		`"2024-01-15T10:50:00"`,
		`"2024-01-15T10:50:00.Z"`,
		`"2024-01-15T10:50:00,1Z"`,
		`"2024-01-15T10:50:00.1234567890Z"`,
		`"2024-01-15T10:50:00+24:00"`,
		`"2024-01-15T10:50:00+23:60"`,
		`"2023-02-29T00:00:00Z"`,
		`"1900-02-29T00:00:00Z"`,
		`"2024-01-15T10:50:60Z"`,
		`"0000-01-01T00:00:00Z"`,
		`"0001-01-01T00:00:00+00:01"`,
		`"9999-12-31T23:59:59-00:01"`,
		`1`,
		`null`,
	}
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			got := types.Time{Seconds: 7, Nanos: 8}
			if err := json.Unmarshal([]byte(input), &got); err == nil {
				t.Fatalf("Unmarshal(%s) error = nil", input)
			}
			if got != (types.Time{Seconds: 7, Nanos: 8}) {
				t.Errorf("Unmarshal(%s) modified receiver to %+v", input, got)
			}
		})
	}
}

func TestTime_MarshalJSON_rejectsInvalidValue(t *testing.T) {
	for _, value := range []types.Time{
		{Seconds: types.MinTimestampSeconds - 1},
		{Seconds: types.MaxTimestampSeconds + 1},
		{Nanos: -1},
		{Nanos: 1_000_000_000},
	} {
		if _, err := json.Marshal(value); err == nil {
			t.Errorf("Marshal(%+v) error = nil", value)
		}
	}
}

func TestTime_CheckValid(t *testing.T) {
	testCases := []struct {
		name  string
		value *types.Time
		valid bool
	}{
		{name: "nil", value: nil},
		{name: "zero", value: &types.Time{}, valid: true},
		{name: "minimum", value: &types.Time{Seconds: types.MinTimestampSeconds}, valid: true},
		{name: "maximum", value: &types.Time{Seconds: types.MaxTimestampSeconds, Nanos: 999_999_999}, valid: true},
		{name: "seconds underflow", value: &types.Time{Seconds: types.MinTimestampSeconds - 1}},
		{name: "seconds overflow", value: &types.Time{Seconds: types.MaxTimestampSeconds + 1}},
		{name: "negative nanos", value: &types.Time{Nanos: -1}},
		{name: "nanos overflow", value: &types.Time{Nanos: 1_000_000_000}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.value.CheckValid()
			if (err == nil) != tc.valid {
				t.Errorf("CheckValid() error = %v, valid = %t", err, tc.valid)
			}
		})
	}
}
