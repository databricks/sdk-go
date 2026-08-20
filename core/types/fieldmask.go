package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// FieldMask represents normalized wire paths for fields of T. T is expected to
// be a generated SDK model. Custom types that reproduce the generator's
// reflection metadata may work but are not supported.
type FieldMask[T any] struct {
	value string
}

// NewFieldMask validates paths against the fields of T and returns their
// sorted, deduplicated, and parent-subsumed representation. The path "*"
// represents full replacement and subsumes all other valid paths.
func NewFieldMask[T any](paths ...string) (*FieldMask[T], error) {
	if err := validateFieldMaskPaths(reflect.TypeFor[T](), paths); err != nil {
		return nil, err
	}
	return &FieldMask[T]{value: normalizeFieldMaskPaths(paths)}, nil
}

// String returns the normalized wire paths separated by commas.
func (m FieldMask[T]) String() string {
	return m.value
}

// MarshalJSON encodes the normalized field mask as a JSON string.
func (m FieldMask[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.value)
}

// UnmarshalJSON decodes and validates a comma-separated field mask without
// changing the receiver when validation fails.
func (m *FieldMask[T]) UnmarshalJSON(data []byte) error {
	if m == nil {
		return fmt.Errorf("cannot unmarshal a field mask into a nil receiver")
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var paths []string
	if value != "" {
		paths = strings.Split(value, ",")
	}
	parsed, err := NewFieldMask[T](paths...)
	if err != nil {
		return err
	}
	m.value = parsed.value
	return nil
}
