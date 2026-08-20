package types_test

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/databricks/sdk-go/core/types"
	"github.com/google/go-cmp/cmp"
)

func TestDuration_MarshalJSON_usesCanonicalFractionWidth(t *testing.T) {
	testCases := []struct {
		name  string
		value types.Duration
		json  string
	}{
		{name: "zero", value: types.Duration{}, json: `"0s"`},
		{name: "positive fractional", value: types.Duration{Seconds: 12, Nanos: 345_000_000}, json: `"12.345s"`},
		{name: "microsecond precision", value: types.Duration{Seconds: 12, Nanos: 345_678_000}, json: `"12.345678s"`},
		{name: "nanosecond precision", value: types.Duration{Seconds: 12, Nanos: 345_678_901}, json: `"12.345678901s"`},
		{name: "millisecond precision uses three digits", value: types.Duration{Nanos: 1_000_000}, json: `"0.001s"`},
		{name: "microsecond precision uses six digits", value: types.Duration{Nanos: 1_000}, json: `"0.000001s"`},
		{name: "nanosecond precision uses nine digits", value: types.Duration{Nanos: 1}, json: `"0.000000001s"`},
		{name: "negative subsecond", value: types.Duration{Nanos: -500_000_000}, json: `"-0.500s"`},
		{name: "negative seconds and nanos", value: types.Duration{Seconds: -12, Nanos: -345_000_000}, json: `"-12.345s"`},
		{name: "protobuf maximum", value: types.Duration{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos}, json: `"315576000000.999999999s"`},
		{name: "protobuf minimum", value: types.Duration{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos}, json: `"-315576000000.999999999s"`},
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

		})
	}
}

func TestDuration_String(t *testing.T) {
	testCases := []struct {
		name  string
		value *types.Duration
		want  string
	}{
		{name: "nil", want: "<nil>"},
		{name: "zero", value: &types.Duration{}, want: "0s"},
		{name: "whole seconds", value: &types.Duration{Seconds: 90}, want: "90s"},
		{name: "fractional", value: &types.Duration{Seconds: -12, Nanos: -345_000_000}, want: "-12.345s"},
		{name: "invalid", value: &types.Duration{Seconds: 1, Nanos: -1}, want: "Duration{Seconds: 1, Nanos: -1}"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestDuration_JSONRoundTrip(t *testing.T) {
	testCases := []types.Duration{
		{},
		{Seconds: 12},
		{Seconds: 12, Nanos: 345_678_901},
		{Seconds: -12, Nanos: -345_678_901},
		{Nanos: 1},
		{Nanos: -1},
		{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos},
		{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos},
	}
	for _, want := range testCases {
		encoded, err := json.Marshal(want)
		if err != nil {
			t.Fatalf("Marshal(%+v) error = %v", want, err)
		}
		var got types.Duration
		if err := json.Unmarshal(encoded, &got); err != nil {
			t.Fatalf("Unmarshal(%s) error = %v", encoded, err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("round trip mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestDuration_UnmarshalJSON_acceptsFixedPointSeconds(t *testing.T) {
	testCases := []struct {
		input string
		want  types.Duration
	}{
		{input: `"1s"`, want: types.Duration{Seconds: 1}},
		{input: `"-0s"`, want: types.Duration{}},
		{input: `"0.1s"`, want: types.Duration{Nanos: 100_000_000}},
		{input: `"0.12s"`, want: types.Duration{Nanos: 120_000_000}},
		{input: `"0.123s"`, want: types.Duration{Nanos: 123_000_000}},
		{input: `"0.1234s"`, want: types.Duration{Nanos: 123_400_000}},
		{input: `"0.12345s"`, want: types.Duration{Nanos: 123_450_000}},
		{input: `"0.123456s"`, want: types.Duration{Nanos: 123_456_000}},
		{input: `"0.1234567s"`, want: types.Duration{Nanos: 123_456_700}},
		{input: `"0.12345678s"`, want: types.Duration{Nanos: 123_456_780}},
		{input: `"-1s"`, want: types.Duration{Seconds: -1}},
		{input: `"-0.1s"`, want: types.Duration{Nanos: -100_000_000}},
		{input: `"0.00000001s"`, want: types.Duration{Nanos: 10}},
		{input: `"1.123456789s"`, want: types.Duration{Seconds: 1, Nanos: 123_456_789}},
		{input: `"315576000000.999999999s"`, want: types.Duration{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos}},
		{input: `"-315576000000.999999999s"`, want: types.Duration{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos}},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			var got types.Duration
			if err := json.Unmarshal([]byte(tc.input), &got); err != nil {
				t.Fatalf("Unmarshal(%s) error = %v", tc.input, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Unmarshal(%s) mismatch (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestDuration_UnmarshalJSON_rejectsInvalidValues(t *testing.T) {
	testCases := []string{
		`""`,
		`"a"`,
		`"-s"`,
		`"s"`,
		`"1"`,
		`"1S"`,
		`"+1s"`,
		`".1s"`,
		`"-.1s"`,
		`"1.s"`,
		`"01s"`,
		`"1,1s"`,
		`"1e3s"`,
		`" 1s"`,
		`"1.1234567890s"`,
		`"1.as"`,
		`"1.2.3s"`,
		`"315576000001s"`,
		`"-315576000001s"`,
		`"--1s"`,
		`1`,
		`null`,
	}
	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			got := types.Duration{Seconds: 7, Nanos: 8}
			if err := json.Unmarshal([]byte(input), &got); err == nil {
				t.Fatalf("Unmarshal(%s) error = nil", input)
			}
			if got != (types.Duration{Seconds: 7, Nanos: 8}) {
				t.Errorf("Unmarshal(%s) modified receiver to %+v", input, got)
			}
		})
	}
}

func TestDuration_TimeConversion(t *testing.T) {
	var nilDuration *types.Duration
	if got := nilDuration.AsDuration(); got != 0 {
		t.Errorf("nil.AsDuration() = %v, want zero", got)
	}

	testCases := []time.Duration{
		0,
		2*time.Hour + 3*time.Millisecond,
		-2*time.Hour - 3*time.Millisecond,
		time.Duration(1<<63 - 1),
		time.Duration(-1 << 63),
	}
	for _, want := range testCases {
		if got := types.NewFromDuration(want).AsDuration(); got != want {
			t.Errorf("NewFromDuration(%v).AsDuration() = %v", want, got)
		}
	}

	const nanosPerSecond = int64(time.Second)
	maxSeconds := int64(math.MaxInt64 / nanosPerSecond)
	maxNanos := int32(math.MaxInt64 - maxSeconds*nanosPerSecond)
	minSeconds := int64(math.MinInt64 / nanosPerSecond)
	minNanos := int32(math.MinInt64 - minSeconds*nanosPerSecond)
	clampCases := []struct {
		name  string
		value types.Duration
		want  time.Duration
	}{
		{name: "maximum exact", value: types.Duration{Seconds: maxSeconds, Nanos: maxNanos}, want: time.Duration(math.MaxInt64)},
		{name: "one above maximum", value: types.Duration{Seconds: maxSeconds, Nanos: maxNanos + 1}, want: time.Duration(math.MaxInt64)},
		{name: "minimum exact", value: types.Duration{Seconds: minSeconds, Nanos: minNanos}, want: time.Duration(math.MinInt64)},
		{name: "one below minimum", value: types.Duration{Seconds: minSeconds, Nanos: minNanos - 1}, want: time.Duration(math.MinInt64)},
		{name: "invalid nanos overflow below boundary second", value: types.Duration{Seconds: maxSeconds - 1, Nanos: math.MaxInt32}, want: time.Duration(math.MaxInt64)},
		{name: "invalid nanos underflow above boundary second", value: types.Duration{Seconds: minSeconds + 1, Nanos: math.MinInt32}, want: time.Duration(math.MinInt64)},
		{name: "protobuf maximum", value: types.Duration{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos}, want: time.Duration(math.MaxInt64)},
		{name: "protobuf minimum", value: types.Duration{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos}, want: time.Duration(math.MinInt64)},
		{name: "arbitrary maximum fields", value: types.Duration{Seconds: math.MaxInt64, Nanos: math.MaxInt32}, want: time.Duration(math.MaxInt64)},
		{name: "arbitrary minimum fields", value: types.Duration{Seconds: math.MinInt64, Nanos: math.MinInt32}, want: time.Duration(math.MinInt64)},
	}
	for _, tc := range clampCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.AsDuration(); got != tc.want {
				t.Errorf("AsDuration() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDuration_Add_normalizesWithoutClamping(t *testing.T) {
	testCases := []struct {
		name  string
		left  *types.Duration
		right *types.Duration
		want  *types.Duration
		valid bool
	}{
		{name: "zero", left: &types.Duration{}, right: &types.Duration{}, want: &types.Duration{}, valid: true},
		{name: "carry", left: &types.Duration{Seconds: 10, Nanos: 900_000_000}, right: &types.Duration{Nanos: 200_000_000}, want: &types.Duration{Seconds: 11, Nanos: 100_000_000}, valid: true},
		{name: "positive borrow", left: &types.Duration{Seconds: 10, Nanos: 100_000_000}, right: &types.Duration{Nanos: -200_000_000}, want: &types.Duration{Seconds: 9, Nanos: 900_000_000}, valid: true},
		{name: "negative borrow", left: &types.Duration{Seconds: -10, Nanos: -100_000_000}, right: &types.Duration{Nanos: 200_000_000}, want: &types.Duration{Seconds: -9, Nanos: -900_000_000}, valid: true},
		{name: "cancel", left: &types.Duration{Seconds: 10, Nanos: 100}, right: &types.Duration{Seconds: -10, Nanos: -100}, want: &types.Duration{}, valid: true},
		{name: "maximum plus one nanosecond", left: &types.Duration{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos}, right: &types.Duration{Nanos: 1}, want: &types.Duration{Seconds: types.MaxDurationSeconds + 1}, valid: false},
		{name: "minimum minus one nanosecond", left: &types.Duration{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos}, right: &types.Duration{Nanos: -1}, want: &types.Duration{Seconds: types.MinDurationSeconds - 1}, valid: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			left, right := *tc.left, *tc.right
			got := tc.left.Add(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Add() mismatch (-want +got):\n%s", diff)
			}
			if (got.CheckValid() == nil) != tc.valid {
				t.Errorf("Add().CheckValid() valid = %t, want %t", got.CheckValid() == nil, tc.valid)
			}
			if *tc.left != left || *tc.right != right {
				t.Error("Add() modified an operand")
			}
		})
	}
}

func TestDuration_Add_nil(t *testing.T) {
	var duration *types.Duration
	if got := duration.Add(&types.Duration{Seconds: 1}); got != nil {
		t.Errorf("nil.Add() = %+v, want nil", got)
	}

	duration = &types.Duration{Seconds: 1, Nanos: 2}
	got := duration.Add(nil)
	if got == duration || *got != *duration {
		t.Errorf("Add(nil) = %+v, want an equal copy", got)
	}
}

func TestDuration_MarshalJSON_rejectsInvalidValue(t *testing.T) {
	for _, value := range []types.Duration{
		{Seconds: types.MaxDurationSeconds + 1},
		{Seconds: types.MinDurationSeconds - 1},
		{Nanos: types.MaxDurationNanos + 1},
		{Nanos: types.MinDurationNanos - 1},
		{Seconds: 1, Nanos: -1},
		{Seconds: -1, Nanos: 1},
	} {
		if _, err := json.Marshal(value); err == nil {
			t.Errorf("Marshal(%+v) error = nil", value)
		}
	}
}

func TestDuration_CheckValid(t *testing.T) {
	testCases := []struct {
		name  string
		value *types.Duration
		valid bool
	}{
		{name: "nil", value: nil},
		{name: "zero", value: &types.Duration{}, valid: true},
		{name: "maximum", value: &types.Duration{Seconds: types.MaxDurationSeconds, Nanos: types.MaxDurationNanos}, valid: true},
		{name: "minimum", value: &types.Duration{Seconds: types.MinDurationSeconds, Nanos: types.MinDurationNanos}, valid: true},
		{name: "seconds overflow", value: &types.Duration{Seconds: types.MaxDurationSeconds + 1}},
		{name: "seconds underflow", value: &types.Duration{Seconds: types.MinDurationSeconds - 1}},
		{name: "nanos overflow", value: &types.Duration{Nanos: types.MaxDurationNanos + 1}},
		{name: "nanos underflow", value: &types.Duration{Nanos: types.MinDurationNanos - 1}},
		{name: "positive seconds negative nanos", value: &types.Duration{Seconds: 1, Nanos: -1}},
		{name: "negative seconds positive nanos", value: &types.Duration{Seconds: -1, Nanos: 1}},
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
