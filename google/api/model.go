
package api

import (
	
	"fmt"
	"io"
)



type FieldBehavior string

const FieldBehaviorFieldBehaviorUnspecified FieldBehavior = `FIELD_BEHAVIOR_UNSPECIFIED`
const FieldBehaviorOptional FieldBehavior = `OPTIONAL`
const FieldBehaviorRequired FieldBehavior = `REQUIRED`
const FieldBehaviorOutputOnly FieldBehavior = `OUTPUT_ONLY`
const FieldBehaviorInputOnly FieldBehavior = `INPUT_ONLY`
const FieldBehaviorImmutable FieldBehavior = `IMMUTABLE`
const FieldBehaviorUnorderedList FieldBehavior = `UNORDERED_LIST`
const FieldBehaviorNonEmptyDefault FieldBehavior = `NON_EMPTY_DEFAULT`

// String representation for [fmt.Print]
func (f *FieldBehavior) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *FieldBehavior) Set(v string) error {
	switch v {
	case `FIELD_BEHAVIOR_UNSPECIFIED`, `OPTIONAL`, `REQUIRED`, `OUTPUT_ONLY`, `INPUT_ONLY`, `IMMUTABLE`, `UNORDERED_LIST`, `NON_EMPTY_DEFAULT`:
		*f = FieldBehavior(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "FIELD_BEHAVIOR_UNSPECIFIED", "OPTIONAL", "REQUIRED", "OUTPUT_ONLY", "INPUT_ONLY", "IMMUTABLE", "UNORDERED_LIST", "NON_EMPTY_DEFAULT"`, v)
	}
}

// Values returns all possible values for FieldBehavior.
//
// There is no guarantee on the order of the values in the slice.
func (f *FieldBehavior) Values() []FieldBehavior {
	return []FieldBehavior{
		FieldBehaviorFieldBehaviorUnspecified,
		FieldBehaviorOptional,
		FieldBehaviorRequired,
		FieldBehaviorOutputOnly,
		FieldBehaviorInputOnly,
		FieldBehaviorImmutable,
		FieldBehaviorUnorderedList,
		FieldBehaviorNonEmptyDefault,
	}
}

// Type always returns FieldBehavior to satisfy [pflag.Value] interface
func (f *FieldBehavior) Type() string {
	return "FieldBehavior"
}





type CustomHttpPattern struct {
	
	Kind string `json:"kind"`
	
	Path string `json:"path"`
	
}



type Http struct {
	
	Rules []HttpRule `json:"rules"`
	
	FullyDecodeReservedExpansion bool `json:"fullyDecodeReservedExpansion"`
	
}



type HttpRule struct {
	
	Selector string `json:"selector"`
	
	Get string `json:"get"`
	
	Put string `json:"put"`
	
	Post string `json:"post"`
	
	Delete string `json:"delete"`
	
	Patch string `json:"patch"`
	
	Custom CustomHttpPattern `json:"custom"`
	
	Body string `json:"body"`
	
	ResponseBody string `json:"responseBody"`
	
	AdditionalBindings []HttpRule `json:"additionalBindings"`
	
}


