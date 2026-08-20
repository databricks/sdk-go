package types_test

import (
	"encoding/json"
	"testing"

	"github.com/databricks/sdk-go/core/types"
)

type fieldMaskPrintedEdition struct {
	Printer *string `fieldmask:"printer"`
}

type fieldMaskBookEdition interface {
	fieldMaskBookEdition()
}

type fieldMaskBookEditionPrintedEdition struct {
	PrintedEdition fieldMaskPrintedEdition `fieldmask:"printed_edition"`
}

func (*fieldMaskBookEditionPrintedEdition) fieldMaskBookEdition() {}

type fieldMaskBookEditionDigitalURI struct {
	DigitalURI string `fieldmask:"digital_uri"`
}

func (*fieldMaskBookEditionDigitalURI) fieldMaskBookEdition() {}

type fieldMaskBookEditionMetadata struct {
	*fieldMaskBookEditionPrintedEdition
	*fieldMaskBookEditionDigitalURI
}

type fieldMaskBook struct {
	Title       *string                    `fieldmask:"title"`
	RelatedBook *fieldMaskBook             `fieldmask:"related_book"`
	Labels      []string                   `fieldmask:"labels"`
	Attributes  map[string]string          `fieldmask:"attributes"`
	PublishedAt *types.Time                `fieldmask:"published_at"`
	ReadingTime *types.Duration            `fieldmask:"reading_time"`
	Details     json.RawMessage            `fieldmask:"details"`
	Values      []json.RawMessage          `fieldmask:"values"`
	Properties  map[string]json.RawMessage `fieldmask:"properties"`
	LegacyName  *string                    `fieldmask:"legacyName"`
	Edition     fieldMaskBookEdition
	_           [0]fieldMaskBookEditionMetadata `fieldmask_oneof:"Edition"`
}

func TestNewFieldMask(t *testing.T) {
	testCases := []struct {
		name    string
		paths   []string
		want    string
		wantErr bool
	}{
		{name: "empty", want: ""},
		{name: "flat", paths: []string{"title"}, want: "title"},
		{name: "nested oneof message", paths: []string{"printed_edition.printer"}, want: "printed_edition.printer"},
		{name: "recursive", paths: []string{"related_book.related_book.title"}, want: "related_book.related_book.title"},
		{name: "oneof scalar", paths: []string{"digital_uri"}, want: "digital_uri"},
		{name: "collections as terminals", paths: []string{"labels", "attributes"}, want: "attributes,labels"},
		{
			name:  "well-known JSON terminals",
			paths: []string{"values", "reading_time", "properties", "published_at", "details"},
			want:  "details,properties,published_at,reading_time,values",
		},
		{name: "non-snake wire name", paths: []string{"legacyName"}, want: "legacyName"},
		{
			name:  "sort deduplicate and subsume children",
			paths: []string{"title", "printed_edition.printer", "printed_edition", "title"},
			want:  "printed_edition,title",
		},
		{name: "empty path", paths: []string{""}, wantErr: true},
		{name: "missing field", paths: []string{"missing"}, wantErr: true},
		{name: "Go field name", paths: []string{"RelatedBook"}, wantErr: true},
		{name: "wildcard", paths: []string{"*"}, want: "*"},
		{
			name:  "wildcard subsumes other paths",
			paths: []string{"title", "*", "printed_edition.printer"},
			want:  "*",
		},
		{name: "bare oneof group", paths: []string{"edition"}, wantErr: true},
		{name: "oneof group prefix", paths: []string{"edition.printed_edition"}, wantErr: true},
		{name: "scalar traversal", paths: []string{"title.value"}, wantErr: true},
		{name: "array traversal", paths: []string{"labels.value"}, wantErr: true},
		{name: "map traversal", paths: []string{"attributes.value"}, wantErr: true},
		{
			name:    "validate before parent subsumption",
			paths:   []string{"printed_edition", "printed_edition.missing"},
			wantErr: true,
		},
		{
			name:    "validate before wildcard subsumption",
			paths:   []string{"*", "missing"},
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mask, err := types.NewFieldMask[fieldMaskBook](testCase.paths...)
			if testCase.wantErr {
				if err == nil {
					t.Fatal("NewFieldMask() returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewFieldMask() returned error: %v", err)
			}
			if got := mask.String(); got != testCase.want {
				t.Errorf("String() = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestFieldMask_JSONRoundTrip(t *testing.T) {
	mask, err := types.NewFieldMask[fieldMaskBook]("title", "printed_edition.printer", "legacyName")
	if err != nil {
		t.Fatalf("NewFieldMask() returned error: %v", err)
	}
	data, err := json.Marshal(mask)
	if err != nil {
		t.Fatalf("json.Marshal() returned error: %v", err)
	}
	if got, want := string(data), `"legacyName,printed_edition.printer,title"`; got != want {
		t.Errorf("json.Marshal() = %s, want %s", got, want)
	}

	var got types.FieldMask[fieldMaskBook]
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() returned error: %v", err)
	}
	if got.String() != mask.String() {
		t.Errorf("round trip String() = %q, want %q", got.String(), mask.String())
	}
}

func TestFieldMask_UnmarshalJSON_zeroValueAndNilPointer(t *testing.T) {
	var value types.FieldMask[fieldMaskBook]
	if err := json.Unmarshal([]byte(`"title,printed_edition.printer"`), &value); err != nil {
		t.Fatalf("json.Unmarshal() into zero value returned error: %v", err)
	}
	if got, want := value.String(), "printed_edition.printer,title"; got != want {
		t.Errorf("zero value String() = %q, want %q", got, want)
	}

	var pointer *types.FieldMask[fieldMaskBook]
	if err := json.Unmarshal([]byte(`"title"`), &pointer); err != nil {
		t.Fatalf("json.Unmarshal() into nil pointer returned error: %v", err)
	}
	if pointer == nil || pointer.String() != "title" {
		t.Errorf("nil pointer unmarshal = %#v, want mask %q", pointer, "title")
	}
}

func TestFieldMask_UnmarshalJSON_preservesReceiverOnFailure(t *testing.T) {
	mask, err := types.NewFieldMask[fieldMaskBook]("title")
	if err != nil {
		t.Fatalf("NewFieldMask() returned error: %v", err)
	}
	if err := json.Unmarshal([]byte(`"missing"`), mask); err == nil {
		t.Fatal("json.Unmarshal() returned nil error")
	}
	if got := mask.String(); got != "title" {
		t.Errorf("String() after failed unmarshal = %q, want %q", got, "title")
	}
}
