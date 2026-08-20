package types

import (
	"fmt"
	"reflect"
	"strings"
)

func validateFieldMaskPaths(root reflect.Type, paths []string) error {
	if root.Kind() != reflect.Struct {
		return fmt.Errorf("field mask root %s must be a struct", root)
	}
	for _, path := range paths {
		if err := validateFieldMaskPath(root, path); err != nil {
			return fmt.Errorf("invalid field mask path %q: %w", path, err)
		}
	}
	return nil
}

func validateFieldMaskPath(root reflect.Type, path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	if path == "*" {
		return nil
	}
	current := root
	segments := strings.Split(path, ".")
	for i, segment := range segments {
		if segment == "*" {
			return fmt.Errorf("wildcard %q must be the entire path", segment)
		}
		field, ok := findFieldMaskField(current, segment)
		if !ok {
			return fmt.Errorf("field %q does not exist", segment)
		}
		if i < len(segments)-1 {
			child, ok := fieldMaskMessageType(field.Type)
			if !ok {
				return fmt.Errorf("field %q may only appear at the end of a path", segment)
			}
			current = child
		}
	}
	return nil
}

func findFieldMaskField(root reflect.Type, path string) (reflect.StructField, bool) {
	for i := range root.NumField() {
		field := root.Field(i)
		fieldPath, hasFieldPath := field.Tag.Lookup("fieldmask")
		oneofGroup, hasOneofGroup := field.Tag.Lookup("fieldmask_oneof")
		switch {
		case hasFieldPath && hasOneofGroup:
			continue
		case hasFieldPath:
			if field.PkgPath == "" && fieldPath == path {
				return field, true
			}
		case hasOneofGroup:
			if oneofGroup == "" {
				continue
			}
			if field, ok := findFieldMaskOneofField(root, field, oneofGroup, path); ok {
				return field, true
			}
		}
	}
	return reflect.StructField{}, false
}

func findFieldMaskOneofField(root reflect.Type, carrier reflect.StructField, group, path string) (reflect.StructField, bool) {
	if carrier.Name != "_" || carrier.Type.Kind() != reflect.Array || carrier.Type.Len() != 0 {
		return reflect.StructField{}, false
	}
	oneofField, ok := root.FieldByName(group)
	if !ok || oneofField.PkgPath != "" || oneofField.Type.Kind() != reflect.Interface {
		return reflect.StructField{}, false
	}

	metadataType := carrier.Type.Elem()
	if metadataType.Kind() != reflect.Struct || metadataType.NumField() == 0 {
		return reflect.StructField{}, false
	}
	for i := range metadataType.NumField() {
		embedded := metadataType.Field(i)
		if !embedded.Anonymous {
			continue
		}
		wrapperType := embedded.Type
		if wrapperType.Kind() == reflect.Pointer {
			wrapperType = wrapperType.Elem()
		}
		if wrapperType.Kind() != reflect.Struct || wrapperType.NumField() != 1 {
			continue
		}
		payload := wrapperType.Field(0)
		fieldPath, ok := payload.Tag.Lookup("fieldmask")
		if ok && payload.PkgPath == "" && fieldPath == path {
			return payload, true
		}
	}
	return reflect.StructField{}, false
}

func fieldMaskMessageType(typ reflect.Type) (reflect.Type, bool) {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, false
	}
	for i := range typ.NumField() {
		field := typ.Field(i)
		if _, ok := field.Tag.Lookup("fieldmask"); ok {
			return typ, true
		}
		if _, ok := field.Tag.Lookup("fieldmask_oneof"); ok {
			return typ, true
		}
	}
	return nil, false
}
