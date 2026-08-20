package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	// MinTimestampSeconds is 0001-01-01T00:00:00Z in Unix seconds, the
	// smallest valid Time seconds value.
	MinTimestampSeconds = int64(-62_135_596_800)
	// MaxTimestampSeconds is 9999-12-31T23:59:59Z in Unix seconds, the
	// largest valid Time seconds value.
	MaxTimestampSeconds = int64(253_402_300_799)
)

// Time is an instant with nanosecond precision. It models the
// google.protobuf.Timestamp well-known type.
//
// The zero value represents the Unix epoch, 1970-01-01T00:00:00Z. Seconds is
// limited to instants from 0001-01-01T00:00:00Z through
// 9999-12-31T23:59:59Z, inclusive. Nanos must be between 0 and 999,999,999,
// inclusive. [Time.CheckValid] checks these invariants.
//
// Its representation and range follow the [google.protobuf.Timestamp definition].
//
// [google.protobuf.Timestamp definition]: https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/timestamp.proto
type Time struct {
	// Seconds is the number of whole seconds since the Unix epoch.
	Seconds int64
	// Nanos is the non-negative fractional-second component in nanoseconds.
	Nanos int32
}

// NewFromTime constructs a Time from a standard-library time. Go can represent
// years outside the Protocol Buffers Timestamp range, so callers converting
// such values must use [Time.CheckValid] before sending them to an API.
func NewFromTime(value time.Time) *Time {
	return &Time{Seconds: value.Unix(), Nanos: int32(value.Nanosecond())}
}

// AsTime converts t to a standard-library time in UTC. It uses [time.Unix]
// normalization for invalid field combinations rather than validating them. A
// nil receiver converts to the Unix epoch.
func (t *Time) AsTime() time.Time {
	if t == nil {
		return time.Unix(0, 0).UTC()
	}
	return time.Unix(t.Seconds, int64(t.Nanos)).UTC()
}

// Add returns t with d added and the nanoseconds normalized. Add expects each
// non-nil operand to satisfy its CheckValid method and does not validate either
// operand. It does not clamp the result to the Time range; callers can continue
// arithmetic and use [Time.CheckValid] on the final result. Add returns nil
// when t is nil. A nil d is treated as zero.
func (t *Time) Add(d *Duration) *Time {
	if t == nil {
		return nil
	}
	if d == nil {
		return &Time{Seconds: t.Seconds, Nanos: t.Nanos}
	}
	result := normalizeTime(t.Seconds+d.Seconds, int64(t.Nanos)+int64(d.Nanos))
	return &result
}

func formatTime(t *Time) string {
	formatted := t.AsTime().Format("2006-01-02T15:04:05.000000000")
	formatted = strings.TrimSuffix(formatted, "000")
	formatted = strings.TrimSuffix(formatted, "000")
	formatted = strings.TrimSuffix(formatted, ".000")
	return formatted + "Z"
}

// String returns an RFC 3339 representation of t in UTC. It returns "<nil>"
// for a nil receiver. Invalid field combinations are rendered as their
// component values.
func (t *Time) String() string {
	if t == nil {
		return "<nil>"
	}
	if err := t.CheckValid(); err != nil {
		return fmt.Sprintf("Time{Seconds: %d, Nanos: %d}", t.Seconds, t.Nanos)
	}
	return formatTime(t)
}

// MarshalJSON encodes t as a ProtoJSON Timestamp in UTC with a Z suffix,
// using zero, three, six, or nine fractional digits as needed to represent its
// nanoseconds exactly. It returns an error when t does not satisfy
// [Time.CheckValid].
//
// See the [ProtoJSON well-known type mapping].
//
// [ProtoJSON well-known type mapping]: https://protobuf.dev/programming-guides/json/#format-description
func (t Time) MarshalJSON() ([]byte, error) {
	if err := t.CheckValid(); err != nil {
		return nil, err
	}
	return json.Marshal(formatTime(&t))
}

// UnmarshalJSON decodes a ProtoJSON Timestamp. It accepts UTC or numeric
// timezone offsets and between zero and nine fractional digits. In accordance
// with the protobuf Timestamp JSON form, fractional seconds use a period; the
// comma separator permitted by the broader RFC 3339 grammar is rejected. The
// JSON value must be a string; null is rejected. If decoding fails, t is
// unchanged.
//
// See the [ProtoJSON well-known type mapping].
//
// [ProtoJSON well-known type mapping]: https://protobuf.dev/programming-guides/json/#format-description
func (t *Time) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if err := validateTimestampSyntax(value); err != nil {
		return fmt.Errorf("parse timestamp %q: %w", value, err)
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return fmt.Errorf("parse timestamp %q: %w", value, err)
	}
	result := NewFromTime(parsed)
	if err := result.CheckValid(); err != nil {
		return fmt.Errorf("parse timestamp %q: %w", value, err)
	}
	*t = *result
	return nil
}

func validateTimestampSyntax(value string) error {
	if len(value) < len("0001-01-01T00:00:00Z") ||
		value[4] != '-' || value[7] != '-' || value[10] != 'T' ||
		value[13] != ':' || value[16] != ':' {
		return fmt.Errorf("value must have the form YYYY-MM-DDTHH:MM:SS[.fffffffff](Z|+HH:MM|-HH:MM)")
	}
	for _, index := range []int{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 14, 15, 17, 18} {
		if value[index] < '0' || value[index] > '9' {
			return fmt.Errorf("date and time components must contain only decimal digits")
		}
	}

	rest := value[19:]
	if len(rest) > 0 && rest[0] == '.' {
		rest = rest[1:]
		digits := 0
		for digits < len(rest) && rest[digits] >= '0' && rest[digits] <= '9' {
			digits++
		}
		if digits == 0 || digits > 9 {
			return fmt.Errorf("fraction must contain between 1 and 9 digits")
		}
		rest = rest[digits:]
	}
	if rest == "Z" {
		return nil
	}
	if len(rest) != 6 || rest[0] != '+' && rest[0] != '-' || rest[3] != ':' {
		return fmt.Errorf("timezone must be Z or a numeric offset in the form +HH:MM or -HH:MM")
	}
	for _, index := range []int{1, 2, 4, 5} {
		if rest[index] < '0' || rest[index] > '9' {
			return fmt.Errorf("timezone offset must contain only decimal digits")
		}
	}
	hours := int(rest[1]-'0')*10 + int(rest[2]-'0')
	minutes := int(rest[4]-'0')*10 + int(rest[5]-'0')
	if hours > 23 || minutes > 59 {
		return fmt.Errorf("timezone offset is outside the RFC 3339 range")
	}
	return nil
}

// CheckValid returns an error unless t satisfies the Protocol Buffers
// Timestamp range and normalization rules. A nil Time is invalid.
func (t *Time) CheckValid() error {
	if t == nil {
		return fmt.Errorf("invalid nil Time")
	}
	if t.Seconds < MinTimestampSeconds {
		return fmt.Errorf("timestamp seconds %d before 0001-01-01", t.Seconds)
	}
	if t.Seconds > MaxTimestampSeconds {
		return fmt.Errorf("timestamp seconds %d after 9999-12-31", t.Seconds)
	}
	if t.Nanos < 0 || t.Nanos >= int32(time.Second) {
		return fmt.Errorf("timestamp nanoseconds %d outside protobuf range", t.Nanos)
	}
	return nil
}

func normalizeTime(seconds, nanos int64) Time {
	seconds += nanos / nanosecondsPerSecond
	nanos %= nanosecondsPerSecond
	if nanos < 0 {
		seconds--
		nanos += nanosecondsPerSecond
	}
	return Time{Seconds: seconds, Nanos: int32(nanos)}
}
