package types

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// MaxDurationSeconds is the largest valid Duration seconds value, inclusive.
	MaxDurationSeconds = int64(315_576_000_000)
	// MinDurationSeconds is the smallest valid Duration seconds value, inclusive.
	MinDurationSeconds = -MaxDurationSeconds
	// MaxDurationNanos is the largest valid Duration nanoseconds value, inclusive.
	MaxDurationNanos = int32(999_999_999)
	// MinDurationNanos is the smallest valid Duration nanoseconds value, inclusive.
	MinDurationNanos = -MaxDurationNanos

	nanosecondsPerSecond = int64(time.Second)
)

// Duration is a signed, fixed-length span of time with nanosecond precision.
// It models the google.protobuf.Duration well-known type without the range
// loss of [time.Duration].
//
// The zero value represents a duration of zero. Seconds must be between
// -315,576,000,000 and +315,576,000,000, inclusive. Nanos must be between
// -999,999,999 and +999,999,999, inclusive, and must have the same sign as
// Seconds when both are non-zero. [Duration.CheckValid] checks these invariants.
//
// Its representation, range, and normalization rules follow the
// [google.protobuf.Duration definition].
//
// [google.protobuf.Duration definition]: https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/duration.proto
type Duration struct {
	// Seconds is the signed whole-second component.
	Seconds int64
	// Nanos is the signed fractional-second component in nanoseconds.
	Nanos int32
}

// NewFromDuration constructs a valid Duration from a standard-library duration.
func NewFromDuration(duration time.Duration) *Duration {
	nanos := duration.Nanoseconds()
	seconds := nanos / nanosecondsPerSecond
	nanos -= seconds * nanosecondsPerSecond
	return &Duration{Seconds: seconds, Nanos: int32(nanos)}
}

// AsDuration converts d to a standard-library duration. It returns the nearest
// [time.Duration] boundary when d is outside the standard library's range. A
// nil receiver converts to zero.
func (d *Duration) AsDuration() time.Duration {
	if d == nil {
		return 0
	}
	const (
		minSeconds = int64(math.MinInt64 / nanosecondsPerSecond)
		maxSeconds = int64(math.MaxInt64 / nanosecondsPerSecond)
		minNanos   = int32(math.MinInt64 - minSeconds*nanosecondsPerSecond)
		maxNanos   = int32(math.MaxInt64 - maxSeconds*nanosecondsPerSecond)
	)
	nanos := int64(d.Nanos)
	secondsCarry := nanos / nanosecondsPerSecond
	if secondsCarry > 0 && d.Seconds > math.MaxInt64-secondsCarry {
		return time.Duration(math.MaxInt64)
	}
	if secondsCarry < 0 && d.Seconds < math.MinInt64-secondsCarry {
		return time.Duration(math.MinInt64)
	}
	seconds := d.Seconds + secondsCarry
	nanos %= nanosecondsPerSecond
	if seconds < minSeconds || seconds == minSeconds && nanos < int64(minNanos) {
		return time.Duration(math.MinInt64)
	}
	if seconds > maxSeconds || seconds == maxSeconds && nanos > int64(maxNanos) {
		return time.Duration(math.MaxInt64)
	}
	return time.Duration(seconds*nanosecondsPerSecond + nanos)
}

// Add returns the sum of d and other with the nanoseconds normalized. Add
// expects each non-nil operand to satisfy [Duration.CheckValid] and does not
// validate either operand. It does not clamp the result to the Duration range;
// callers can continue arithmetic and use CheckValid on the final result. Add
// returns nil when d is nil. A nil other is treated as zero.
func (d *Duration) Add(other *Duration) *Duration {
	if d == nil {
		return nil
	}
	if other == nil {
		return &Duration{Seconds: d.Seconds, Nanos: d.Nanos}
	}
	result := normalizeDuration(d.Seconds+other.Seconds, int64(d.Nanos)+int64(other.Nanos))
	return &result
}

func normalizeDuration(seconds, nanos int64) Duration {
	seconds += nanos / nanosecondsPerSecond
	nanos %= nanosecondsPerSecond
	if seconds > 0 && nanos < 0 {
		seconds--
		nanos += nanosecondsPerSecond
	} else if seconds < 0 && nanos > 0 {
		seconds++
		nanos -= nanosecondsPerSecond
	}
	return Duration{Seconds: seconds, Nanos: int32(nanos)}
}

func formatDuration(d *Duration) string {
	negative := d.Seconds < 0 || d.Nanos < 0
	seconds := d.Seconds
	nanos := int64(d.Nanos)
	if seconds < 0 {
		seconds = -seconds
	}
	if nanos < 0 {
		nanos = -nanos
	}
	out := strconv.FormatInt(seconds, 10)
	if nanos != 0 {
		fraction := fmt.Sprintf("%09d", nanos)
		switch {
		case fraction[3:] == "000000":
			fraction = fraction[:3]
		case fraction[6:] == "000":
			fraction = fraction[:6]
		}
		out += "." + fraction
	}
	if negative {
		out = "-" + out
	}
	return out + "s"
}

// String returns a formatted representation of d as signed seconds with an
// "s" suffix. It returns "<nil>" for a nil receiver. Invalid field
// combinations are rendered as their component values.
func (d *Duration) String() string {
	if d == nil {
		return "<nil>"
	}
	if err := d.CheckValid(); err != nil {
		return fmt.Sprintf("Duration{Seconds: %d, Nanos: %d}", d.Seconds, d.Nanos)
	}
	return formatDuration(d)
}

// MarshalJSON encodes d as a ProtoJSON Duration, using zero, three, six,
// or nine fractional digits as needed to represent its nanoseconds exactly.
// It returns an error when d does not satisfy [Duration.CheckValid].
//
// See the [ProtoJSON well-known type mapping].
//
// [ProtoJSON well-known type mapping]: https://protobuf.dev/programming-guides/json/#format-description
func (d Duration) MarshalJSON() ([]byte, error) {
	if err := d.CheckValid(); err != nil {
		return nil, err
	}
	return json.Marshal(formatDuration(&d))
}

// UnmarshalJSON decodes the SDK's accepted subset of ProtoJSON Duration inputs.
// It accepts the form
// -?(0|[1-9][0-9]*)(\.[0-9]{1,9})?s. It intentionally rejects non-canonical
// extensions such as a leading plus sign or a missing digit on either side of
// the decimal point. The JSON value must be a string; null is rejected. If
// decoding fails, d is unchanged.
//
// See the [ProtoJSON well-known type mapping].
//
// [ProtoJSON well-known type mapping]: https://protobuf.dev/programming-guides/json/#format-description
func (d *Duration) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	seconds, nanos, err := parseDuration(value)
	if err != nil {
		return fmt.Errorf("parse duration %q: %w", value, err)
	}
	parsed := Duration{Seconds: seconds, Nanos: nanos}
	if err := parsed.CheckValid(); err != nil {
		return fmt.Errorf("parse duration %q: %w", value, err)
	}
	*d = parsed
	return nil
}

// parseDuration parses the protobuf JSON Duration form:
//
//	-?(0|[1-9][0-9]*)(\.[0-9]{1,9})?s
//
// Whole seconds are required. A decimal point must be followed by one to nine
// fractional digits. Leading plus signs, leading zeroes, and a missing "s"
// suffix are rejected.
func parseDuration(input string) (int64, int32, error) {
	if len(input) < 2 || input[len(input)-1] != 's' {
		return 0, 0, fmt.Errorf("value must end with %q", "s")
	}
	value := input[:len(input)-1]

	negative := false
	if value[0] == '-' {
		negative = true
		value = value[1:]
	}
	if len(value) == 0 {
		return 0, 0, fmt.Errorf("value must contain seconds")
	}

	integerEnd := strings.IndexByte(value, '.')
	if integerEnd < 0 {
		integerEnd = len(value)
	}
	integerPart := value[:integerEnd]
	if len(integerPart) == 0 || len(integerPart) > 1 && integerPart[0] == '0' {
		return 0, 0, fmt.Errorf("seconds must be zero or start with a non-zero digit")
	}
	for _, digit := range integerPart {
		if digit < '0' || digit > '9' {
			return 0, 0, fmt.Errorf("seconds must contain only decimal digits")
		}
	}

	fraction := ""
	if integerEnd < len(value) {
		fraction = value[integerEnd+1:]
		if len(fraction) == 0 || len(fraction) > 9 {
			return 0, 0, fmt.Errorf("fraction must contain between 1 and 9 digits")
		}
		for _, digit := range fraction {
			if digit < '0' || digit > '9' {
				return 0, 0, fmt.Errorf("fraction must contain only decimal digits")
			}
		}
	}

	seconds, err := strconv.ParseInt(integerPart, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse seconds: %w", err)
	}

	var nanos int64
	if fraction != "" {
		fraction += strings.Repeat("0", 9-len(fraction))
		nanos, err = strconv.ParseInt(fraction, 10, 32)
		if err != nil {
			return 0, 0, fmt.Errorf("parse fraction: %w", err)
		}
	}

	if negative {
		seconds = -seconds
		nanos = -nanos
	}
	return seconds, int32(nanos), nil
}

// CheckValid returns an error unless d satisfies the Protocol Buffers Duration
// range and normalization rules. A nil Duration is invalid.
func (d *Duration) CheckValid() error {
	if d == nil {
		return fmt.Errorf("invalid nil Duration")
	}
	if d.Seconds < MinDurationSeconds || d.Seconds > MaxDurationSeconds {
		return fmt.Errorf("duration seconds %d outside protobuf range", d.Seconds)
	}
	if d.Nanos < MinDurationNanos || d.Nanos > MaxDurationNanos {
		return fmt.Errorf("duration nanoseconds %d outside protobuf range", d.Nanos)
	}
	if d.Seconds > 0 && d.Nanos < 0 || d.Seconds < 0 && d.Nanos > 0 {
		return fmt.Errorf("duration seconds and nanoseconds have different signs")
	}
	return nil
}
