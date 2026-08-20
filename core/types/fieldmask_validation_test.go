package types

import (
	"fmt"
	"reflect"
	"testing"
)

type fieldMaskRecursiveTestModel struct {
	Related *fieldMaskRecursiveTestModel `fieldmask:"related"`
	Title   string                       `fieldmask:"title"`
}

type fieldMaskOneofTest interface {
	fieldMaskOneofTest()
}

type fieldMaskPrintedEditionTest struct {
	Printer string `fieldmask:"printer"`
}

type fieldMaskEditionPrintedTest struct {
	PrintedEdition fieldMaskPrintedEditionTest `fieldmask:"printed_edition"`
}

func (*fieldMaskEditionPrintedTest) fieldMaskOneofTest() {}

type fieldMaskEditionDigitalTest struct {
	DigitalURI string `fieldmask:"digital_uri"`
}

func (*fieldMaskEditionDigitalTest) fieldMaskOneofTest() {}

type fieldMaskEditionMetadataTest struct {
	*fieldMaskEditionPrintedTest
	*fieldMaskEditionDigitalTest
}

type fieldMaskUntaggedWrapperTest struct {
	Value string
}

func (*fieldMaskUntaggedWrapperTest) fieldMaskOneofTest() {}

type fieldMaskMalformedWrapperTest struct {
	First  string `fieldmask:"first"`
	Second string `fieldmask:"second"`
}

func (*fieldMaskMalformedWrapperTest) fieldMaskOneofTest() {}

type fieldMaskPartialEditionMetadataTest struct {
	*fieldMaskEditionPrintedTest
	*fieldMaskUntaggedWrapperTest
	*fieldMaskMalformedWrapperTest
}

func TestValidateFieldMaskPaths(t *testing.T) {
	pointerRoot := reflect.TypeOf(&struct{}{})
	testCases := []struct {
		name    string
		root    reflect.Type
		paths   []string
		wantErr string
	}{
		{
			name: "empty",
			root: reflect.TypeOf(struct{}{}),
		},
		{
			name: "flat",
			root: reflect.TypeOf(struct {
				Name string `fieldmask:"name"`
			}{}),
			paths: []string{"name"},
		},
		{
			name: "nested",
			root: reflect.TypeOf(struct {
				Child *struct {
					Name string `fieldmask:"name"`
				} `fieldmask:"child"`
			}{}),
			paths: []string{"child.name"},
		},
		{
			name:  "recursive",
			root:  reflect.TypeFor[fieldMaskRecursiveTestModel](),
			paths: []string{"related.related.title"},
		},
		{
			name: "oneof",
			root: reflect.TypeOf(struct {
				Edition fieldMaskOneofTest
				_       [0]fieldMaskEditionMetadataTest `fieldmask_oneof:"Edition"`
			}{}),
			paths: []string{"printed_edition.printer", "digital_uri"},
		},
		{
			name: "terminal fields",
			root: reflect.TypeOf(struct {
				Scalar string            `fieldmask:"scalar"`
				Array  []string          `fieldmask:"array"`
				Map    map[string]string `fieldmask:"map"`
				Custom struct {
					Value string
				} `fieldmask:"custom"`
			}{}),
			paths: []string{"scalar", "array", "map", "custom"},
		},
		{
			name: "non-snake tag",
			root: reflect.TypeOf(struct {
				Legacy string `fieldmask:"legacyName"`
			}{}),
			paths: []string{"legacyName"},
		},
		{
			name: "duplicate tag keeps first field",
			root: reflect.TypeOf(struct {
				First  string `fieldmask:"duplicate"`
				Second struct {
					Child string `fieldmask:"child"`
				} `fieldmask:"duplicate"`
			}{}),
			paths:   []string{"duplicate.child"},
			wantErr: `invalid field mask path "duplicate.child": field "duplicate" may only appear at the end of a path`,
		},
		{
			name: "missing field tag",
			root: reflect.TypeOf(struct {
				Name  string `fieldmask:"name"`
				Added string
			}{}),
			paths:   []string{"added"},
			wantErr: `invalid field mask path "added": field "added" does not exist`,
		},
		{
			name: "field with conflicting metadata tags",
			root: reflect.TypeOf(struct {
				Choice fieldMaskOneofTest
				Field  string `fieldmask:"field" fieldmask_oneof:"Choice"`
			}{}),
			paths:   []string{"field"},
			wantErr: `invalid field mask path "field": field "field" does not exist`,
		},
		{
			name: "malformed oneof carrier",
			root: reflect.TypeOf(struct {
				Edition  fieldMaskOneofTest
				Metadata [0]fieldMaskEditionMetadataTest `fieldmask_oneof:"Edition"`
			}{}),
			paths:   []string{"digital_uri"},
			wantErr: `invalid field mask path "digital_uri": field "digital_uri" does not exist`,
		},
		{
			name: "oneof carrier references non-interface field",
			root: reflect.TypeOf(struct {
				Edition string
				_       [0]fieldMaskEditionMetadataTest `fieldmask_oneof:"Edition"`
			}{}),
			paths:   []string{"digital_uri"},
			wantErr: `invalid field mask path "digital_uri": field "digital_uri" does not exist`,
		},
		{
			name: "malformed oneof wrappers",
			root: reflect.TypeOf(struct {
				Edition fieldMaskOneofTest
				_       [0]struct {
					*fieldMaskUntaggedWrapperTest
					*fieldMaskMalformedWrapperTest
				} `fieldmask_oneof:"Edition"`
			}{}),
			paths:   []string{"value"},
			wantErr: `invalid field mask path "value": field "value" does not exist`,
		},
		{
			name: "partially usable oneof",
			root: reflect.TypeOf(struct {
				Edition fieldMaskOneofTest
				_       [0]fieldMaskPartialEditionMetadataTest `fieldmask_oneof:"Edition"`
			}{}),
			paths: []string{"printed_edition.printer"},
		},
		{
			name: "partially usable oneof omits malformed wrapper",
			root: reflect.TypeOf(struct {
				Edition fieldMaskOneofTest
				_       [0]fieldMaskPartialEditionMetadataTest `fieldmask_oneof:"Edition"`
			}{}),
			paths:   []string{"value"},
			wantErr: `invalid field mask path "value": field "value" does not exist`,
		},
		{
			name: "wildcard",
			root: reflect.TypeOf(struct {
				Name string `fieldmask:"name"`
			}{}),
			paths: []string{"*"},
		},
		{
			name: "wildcard path segment",
			root: reflect.TypeOf(struct {
				Child struct {
					Name string `fieldmask:"name"`
				} `fieldmask:"child"`
			}{}),
			paths:   []string{"child.*"},
			wantErr: `invalid field mask path "child.*": wildcard "*" must be the entire path`,
		},
		{
			name: "missing field",
			root: reflect.TypeOf(struct {
				Name string `fieldmask:"name"`
			}{}),
			paths:   []string{"missing"},
			wantErr: `invalid field mask path "missing": field "missing" does not exist`,
		},
		{
			name: "scalar traversal",
			root: reflect.TypeOf(struct {
				Scalar string `fieldmask:"scalar"`
			}{}),
			paths:   []string{"scalar.child"},
			wantErr: `invalid field mask path "scalar.child": field "scalar" may only appear at the end of a path`,
		},
		{
			name: "array traversal",
			root: reflect.TypeOf(struct {
				Array []struct {
					Child string `fieldmask:"child"`
				} `fieldmask:"array"`
			}{}),
			paths:   []string{"array.child"},
			wantErr: `invalid field mask path "array.child": field "array" may only appear at the end of a path`,
		},
		{
			name: "map traversal",
			root: reflect.TypeOf(struct {
				Map map[string]struct {
					Child string `fieldmask:"child"`
				} `fieldmask:"map"`
			}{}),
			paths:   []string{"map.child"},
			wantErr: `invalid field mask path "map.child": field "map" may only appear at the end of a path`,
		},
		{
			name: "untagged struct traversal",
			root: reflect.TypeOf(struct {
				Custom struct {
					Child string
				} `fieldmask:"custom"`
			}{}),
			paths:   []string{"custom.child"},
			wantErr: `invalid field mask path "custom.child": field "custom" may only appear at the end of a path`,
		},
		{
			name: "empty path",
			root: reflect.TypeOf(struct {
				Name string `fieldmask:"name"`
			}{}),
			paths:   []string{""},
			wantErr: `invalid field mask path "": path is empty`,
		},
		{
			name:    "non-struct root without paths",
			root:    reflect.TypeFor[string](),
			wantErr: "field mask root string must be a struct",
		},
		{
			name:    "pointer root without paths",
			root:    pointerRoot,
			wantErr: fmt.Sprintf("field mask root %s must be a struct", pointerRoot),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateFieldMaskPaths(testCase.root, testCase.paths)
			if testCase.wantErr == "" {
				if err != nil {
					t.Fatalf("validateFieldMaskPaths() returned error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("validateFieldMaskPaths() returned nil error")
			}
			if got := err.Error(); got != testCase.wantErr {
				t.Errorf("validateFieldMaskPaths() error = %q, want %q", got, testCase.wantErr)
			}
		})
	}
}
