package api

import (
	"fmt"
	"io"
)

type FieldBehavior string

const (
	FieldBehaviorFieldBehaviorUnspecified FieldBehavior = "FIELD_BEHAVIOR_UNSPECIFIED"
	FieldBehaviorOptional FieldBehavior = "OPTIONAL"
	FieldBehaviorRequired FieldBehavior = "REQUIRED"
	FieldBehaviorOutputOnly FieldBehavior = "OUTPUT_ONLY"
	FieldBehaviorInputOnly FieldBehavior = "INPUT_ONLY"
	FieldBehaviorImmutable FieldBehavior = "IMMUTABLE"
	FieldBehaviorUnorderedList FieldBehavior = "UNORDERED_LIST"
	FieldBehaviorNonEmptyDefault FieldBehavior = "NON_EMPTY_DEFAULT"
)

// String representation for [fmt.Print].
func (f *FieldBehavior) String() string {
	return string(*f)
}

type CustomHttpPattern struct {
	Kind *string `json:"kind"`
	Path *string `json:"path"`
}

type Http struct {
	Rules []HttpRule `json:"rules"`
	FullyDecodeReservedExpansion *bool `json:"fullyDecodeReservedExpansion"`
}

type HttpRule struct {
	Selector *string `json:"selector"`
	Get *string `json:"get"`
	Put *string `json:"put"`
	Post *string `json:"post"`
	Delete *string `json:"delete"`
	Patch *string `json:"patch"`
	Custom *CustomHttpPattern `json:"custom"`
	Body *string `json:"body"`
	ResponseBody *string `json:"responseBody"`
	AdditionalBindings []HttpRule `json:"additionalBindings"`
}
