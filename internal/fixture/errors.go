package fixture

import (
	"fmt"
	"strings"
)

// MismatchError is returned when a request does not match the expected fixture.
type MismatchError struct {
	Request  Request
	Expected Fixture
	Index    int
	Cause    error
}

func (e *MismatchError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("fixture mismatch at index %d:\n", e.Index))
	sb.WriteString(fmt.Sprintf("  received: %s %s\n", e.Request.Method, e.Request.Path))
	sb.WriteString(fmt.Sprintf("  expected: %s %s\n", e.Expected.Method, e.Expected.Path))
	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf("  cause: %v\n", e.Cause))
	}
	return sb.String()
}

func (e *MismatchError) Unwrap() error {
	return e.Cause
}

// UnexpectedRequestError is returned when a request is received after all
// fixtures have been consumed.
type UnexpectedRequestError struct {
	Request  Request
	Message  string
	Expected []Fixture
	Received []Request
}

func (e *UnexpectedRequestError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("unexpected request: %s %s\n", e.Request.Method, e.Request.Path))
	sb.WriteString(fmt.Sprintf("  %s\n", e.Message))
	sb.WriteString(fmt.Sprintf("  expected %d fixtures, received %d requests\n",
		len(e.Expected), len(e.Received)))
	return sb.String()
}

// UnconsumedFixturesError is returned when Verify() is called but not all
// fixtures were consumed.
type UnconsumedFixturesError struct {
	Fixtures  []Fixture
	Consumed  []bool
	NextIndex int
}

func (e *UnconsumedFixturesError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("unconsumed fixtures: %d of %d fixtures were not matched\n",
		len(e.Fixtures)-e.NextIndex, len(e.Fixtures)))
	sb.WriteString("  unconsumed:\n")
	for i := e.NextIndex; i < len(e.Fixtures); i++ {
		f := e.Fixtures[i]
		sb.WriteString(fmt.Sprintf("    [%d] %s %s\n", i, f.Method, f.Path))
	}
	return sb.String()
}
