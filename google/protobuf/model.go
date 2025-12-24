
package protobuf

import (
	
	"fmt"
	"io"
)



type CType string

const CTypeString CType = `STRING`
const CTypeCord CType = `CORD`
const CTypeStringPiece CType = `STRING_PIECE`

// String representation for [fmt.Print]
func (f *CType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *CType) Set(v string) error {
	switch v {
	case `STRING`, `CORD`, `STRING_PIECE`:
		*f = CType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "STRING", "CORD", "STRING_PIECE"`, v)
	}
}

// Values returns all possible values for CType.
//
// There is no guarantee on the order of the values in the slice.
func (f *CType) Values() []CType {
	return []CType{
		CTypeString,
		CTypeCord,
		CTypeStringPiece,
	}
}

// Type always returns CType to satisfy [pflag.Value] interface
func (f *CType) Type() string {
	return "CType"
}



type DefaultSymbolVisibility string

const DefaultSymbolVisibilityDefaultSymbolVisibilityUnknown DefaultSymbolVisibility = `DEFAULT_SYMBOL_VISIBILITY_UNKNOWN`
const DefaultSymbolVisibilityExportAll DefaultSymbolVisibility = `EXPORT_ALL`
const DefaultSymbolVisibilityExportTopLevel DefaultSymbolVisibility = `EXPORT_TOP_LEVEL`
const DefaultSymbolVisibilityLocalAll DefaultSymbolVisibility = `LOCAL_ALL`
const DefaultSymbolVisibilityStrict DefaultSymbolVisibility = `STRICT`

// String representation for [fmt.Print]
func (f *DefaultSymbolVisibility) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *DefaultSymbolVisibility) Set(v string) error {
	switch v {
	case `DEFAULT_SYMBOL_VISIBILITY_UNKNOWN`, `EXPORT_ALL`, `EXPORT_TOP_LEVEL`, `LOCAL_ALL`, `STRICT`:
		*f = DefaultSymbolVisibility(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "DEFAULT_SYMBOL_VISIBILITY_UNKNOWN", "EXPORT_ALL", "EXPORT_TOP_LEVEL", "LOCAL_ALL", "STRICT"`, v)
	}
}

// Values returns all possible values for DefaultSymbolVisibility.
//
// There is no guarantee on the order of the values in the slice.
func (f *DefaultSymbolVisibility) Values() []DefaultSymbolVisibility {
	return []DefaultSymbolVisibility{
		DefaultSymbolVisibilityDefaultSymbolVisibilityUnknown,
		DefaultSymbolVisibilityExportAll,
		DefaultSymbolVisibilityExportTopLevel,
		DefaultSymbolVisibilityLocalAll,
		DefaultSymbolVisibilityStrict,
	}
}

// Type always returns DefaultSymbolVisibility to satisfy [pflag.Value] interface
func (f *DefaultSymbolVisibility) Type() string {
	return "DefaultSymbolVisibility"
}



type Edition string

const EditionEditionUnknown Edition = `EDITION_UNKNOWN`
const EditionEditionLegacy Edition = `EDITION_LEGACY`
const EditionEditionProto2 Edition = `EDITION_PROTO2`
const EditionEditionProto3 Edition = `EDITION_PROTO3`
const EditionEdition2023 Edition = `EDITION_2023`
const EditionEdition2024 Edition = `EDITION_2024`
const EditionEdition1TestOnly Edition = `EDITION_1_TEST_ONLY`
const EditionEdition2TestOnly Edition = `EDITION_2_TEST_ONLY`
const EditionEdition99997TestOnly Edition = `EDITION_99997_TEST_ONLY`
const EditionEdition99998TestOnly Edition = `EDITION_99998_TEST_ONLY`
const EditionEdition99999TestOnly Edition = `EDITION_99999_TEST_ONLY`
const EditionEditionMax Edition = `EDITION_MAX`

// String representation for [fmt.Print]
func (f *Edition) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Edition) Set(v string) error {
	switch v {
	case `EDITION_UNKNOWN`, `EDITION_LEGACY`, `EDITION_PROTO2`, `EDITION_PROTO3`, `EDITION_2023`, `EDITION_2024`, `EDITION_1_TEST_ONLY`, `EDITION_2_TEST_ONLY`, `EDITION_99997_TEST_ONLY`, `EDITION_99998_TEST_ONLY`, `EDITION_99999_TEST_ONLY`, `EDITION_MAX`:
		*f = Edition(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "EDITION_UNKNOWN", "EDITION_LEGACY", "EDITION_PROTO2", "EDITION_PROTO3", "EDITION_2023", "EDITION_2024", "EDITION_1_TEST_ONLY", "EDITION_2_TEST_ONLY", "EDITION_99997_TEST_ONLY", "EDITION_99998_TEST_ONLY", "EDITION_99999_TEST_ONLY", "EDITION_MAX"`, v)
	}
}

// Values returns all possible values for Edition.
//
// There is no guarantee on the order of the values in the slice.
func (f *Edition) Values() []Edition {
	return []Edition{
		EditionEditionUnknown,
		EditionEditionLegacy,
		EditionEditionProto2,
		EditionEditionProto3,
		EditionEdition2023,
		EditionEdition2024,
		EditionEdition1TestOnly,
		EditionEdition2TestOnly,
		EditionEdition99997TestOnly,
		EditionEdition99998TestOnly,
		EditionEdition99999TestOnly,
		EditionEditionMax,
	}
}

// Type always returns Edition to satisfy [pflag.Value] interface
func (f *Edition) Type() string {
	return "Edition"
}



type EnforceNamingStyle string

const EnforceNamingStyleEnforceNamingStyleUnknown EnforceNamingStyle = `ENFORCE_NAMING_STYLE_UNKNOWN`
const EnforceNamingStyleStyle2024 EnforceNamingStyle = `STYLE2024`
const EnforceNamingStyleStyleLegacy EnforceNamingStyle = `STYLE_LEGACY`

// String representation for [fmt.Print]
func (f *EnforceNamingStyle) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *EnforceNamingStyle) Set(v string) error {
	switch v {
	case `ENFORCE_NAMING_STYLE_UNKNOWN`, `STYLE2024`, `STYLE_LEGACY`:
		*f = EnforceNamingStyle(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ENFORCE_NAMING_STYLE_UNKNOWN", "STYLE2024", "STYLE_LEGACY"`, v)
	}
}

// Values returns all possible values for EnforceNamingStyle.
//
// There is no guarantee on the order of the values in the slice.
func (f *EnforceNamingStyle) Values() []EnforceNamingStyle {
	return []EnforceNamingStyle{
		EnforceNamingStyleEnforceNamingStyleUnknown,
		EnforceNamingStyleStyle2024,
		EnforceNamingStyleStyleLegacy,
	}
}

// Type always returns EnforceNamingStyle to satisfy [pflag.Value] interface
func (f *EnforceNamingStyle) Type() string {
	return "EnforceNamingStyle"
}



type EnumType string

const EnumTypeEnumTypeUnknown EnumType = `ENUM_TYPE_UNKNOWN`
const EnumTypeOpen EnumType = `OPEN`
const EnumTypeClosed EnumType = `CLOSED`

// String representation for [fmt.Print]
func (f *EnumType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *EnumType) Set(v string) error {
	switch v {
	case `ENUM_TYPE_UNKNOWN`, `OPEN`, `CLOSED`:
		*f = EnumType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "ENUM_TYPE_UNKNOWN", "OPEN", "CLOSED"`, v)
	}
}

// Values returns all possible values for EnumType.
//
// There is no guarantee on the order of the values in the slice.
func (f *EnumType) Values() []EnumType {
	return []EnumType{
		EnumTypeEnumTypeUnknown,
		EnumTypeOpen,
		EnumTypeClosed,
	}
}

// Type always returns EnumType to satisfy [pflag.Value] interface
func (f *EnumType) Type() string {
	return "EnumType"
}



type FieldPresence string

const FieldPresenceFieldPresenceUnknown FieldPresence = `FIELD_PRESENCE_UNKNOWN`
const FieldPresenceExplicit FieldPresence = `EXPLICIT`
const FieldPresenceImplicit FieldPresence = `IMPLICIT`
const FieldPresenceLegacyRequired FieldPresence = `LEGACY_REQUIRED`

// String representation for [fmt.Print]
func (f *FieldPresence) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *FieldPresence) Set(v string) error {
	switch v {
	case `FIELD_PRESENCE_UNKNOWN`, `EXPLICIT`, `IMPLICIT`, `LEGACY_REQUIRED`:
		*f = FieldPresence(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "FIELD_PRESENCE_UNKNOWN", "EXPLICIT", "IMPLICIT", "LEGACY_REQUIRED"`, v)
	}
}

// Values returns all possible values for FieldPresence.
//
// There is no guarantee on the order of the values in the slice.
func (f *FieldPresence) Values() []FieldPresence {
	return []FieldPresence{
		FieldPresenceFieldPresenceUnknown,
		FieldPresenceExplicit,
		FieldPresenceImplicit,
		FieldPresenceLegacyRequired,
	}
}

// Type always returns FieldPresence to satisfy [pflag.Value] interface
func (f *FieldPresence) Type() string {
	return "FieldPresence"
}



type IdempotencyLevel string

const IdempotencyLevelIdempotencyUnknown IdempotencyLevel = `IDEMPOTENCY_UNKNOWN`
const IdempotencyLevelNoSideEffects IdempotencyLevel = `NO_SIDE_EFFECTS`
const IdempotencyLevelIdempotent IdempotencyLevel = `IDEMPOTENT`

// String representation for [fmt.Print]
func (f *IdempotencyLevel) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *IdempotencyLevel) Set(v string) error {
	switch v {
	case `IDEMPOTENCY_UNKNOWN`, `NO_SIDE_EFFECTS`, `IDEMPOTENT`:
		*f = IdempotencyLevel(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "IDEMPOTENCY_UNKNOWN", "NO_SIDE_EFFECTS", "IDEMPOTENT"`, v)
	}
}

// Values returns all possible values for IdempotencyLevel.
//
// There is no guarantee on the order of the values in the slice.
func (f *IdempotencyLevel) Values() []IdempotencyLevel {
	return []IdempotencyLevel{
		IdempotencyLevelIdempotencyUnknown,
		IdempotencyLevelNoSideEffects,
		IdempotencyLevelIdempotent,
	}
}

// Type always returns IdempotencyLevel to satisfy [pflag.Value] interface
func (f *IdempotencyLevel) Type() string {
	return "IdempotencyLevel"
}



type JsType string

const JsTypeJsNormal JsType = `JS_NORMAL`
const JsTypeJsString JsType = `JS_STRING`
const JsTypeJsNumber JsType = `JS_NUMBER`

// String representation for [fmt.Print]
func (f *JsType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *JsType) Set(v string) error {
	switch v {
	case `JS_NORMAL`, `JS_STRING`, `JS_NUMBER`:
		*f = JsType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "JS_NORMAL", "JS_STRING", "JS_NUMBER"`, v)
	}
}

// Values returns all possible values for JsType.
//
// There is no guarantee on the order of the values in the slice.
func (f *JsType) Values() []JsType {
	return []JsType{
		JsTypeJsNormal,
		JsTypeJsString,
		JsTypeJsNumber,
	}
}

// Type always returns JsType to satisfy [pflag.Value] interface
func (f *JsType) Type() string {
	return "JsType"
}



type JsonFormat string

const JsonFormatJsonFormatUnknown JsonFormat = `JSON_FORMAT_UNKNOWN`
const JsonFormatAllow JsonFormat = `ALLOW`
const JsonFormatLegacyBestEffort JsonFormat = `LEGACY_BEST_EFFORT`

// String representation for [fmt.Print]
func (f *JsonFormat) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *JsonFormat) Set(v string) error {
	switch v {
	case `JSON_FORMAT_UNKNOWN`, `ALLOW`, `LEGACY_BEST_EFFORT`:
		*f = JsonFormat(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "JSON_FORMAT_UNKNOWN", "ALLOW", "LEGACY_BEST_EFFORT"`, v)
	}
}

// Values returns all possible values for JsonFormat.
//
// There is no guarantee on the order of the values in the slice.
func (f *JsonFormat) Values() []JsonFormat {
	return []JsonFormat{
		JsonFormatJsonFormatUnknown,
		JsonFormatAllow,
		JsonFormatLegacyBestEffort,
	}
}

// Type always returns JsonFormat to satisfy [pflag.Value] interface
func (f *JsonFormat) Type() string {
	return "JsonFormat"
}



type Label string

const LabelLabelOptional Label = `LABEL_OPTIONAL`
const LabelLabelRepeated Label = `LABEL_REPEATED`
const LabelLabelRequired Label = `LABEL_REQUIRED`

// String representation for [fmt.Print]
func (f *Label) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Label) Set(v string) error {
	switch v {
	case `LABEL_OPTIONAL`, `LABEL_REPEATED`, `LABEL_REQUIRED`:
		*f = Label(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "LABEL_OPTIONAL", "LABEL_REPEATED", "LABEL_REQUIRED"`, v)
	}
}

// Values returns all possible values for Label.
//
// There is no guarantee on the order of the values in the slice.
func (f *Label) Values() []Label {
	return []Label{
		LabelLabelOptional,
		LabelLabelRepeated,
		LabelLabelRequired,
	}
}

// Type always returns Label to satisfy [pflag.Value] interface
func (f *Label) Type() string {
	return "Label"
}



type MessageEncoding string

const MessageEncodingMessageEncodingUnknown MessageEncoding = `MESSAGE_ENCODING_UNKNOWN`
const MessageEncodingLengthPrefixed MessageEncoding = `LENGTH_PREFIXED`
const MessageEncodingDelimited MessageEncoding = `DELIMITED`

// String representation for [fmt.Print]
func (f *MessageEncoding) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *MessageEncoding) Set(v string) error {
	switch v {
	case `MESSAGE_ENCODING_UNKNOWN`, `LENGTH_PREFIXED`, `DELIMITED`:
		*f = MessageEncoding(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "MESSAGE_ENCODING_UNKNOWN", "LENGTH_PREFIXED", "DELIMITED"`, v)
	}
}

// Values returns all possible values for MessageEncoding.
//
// There is no guarantee on the order of the values in the slice.
func (f *MessageEncoding) Values() []MessageEncoding {
	return []MessageEncoding{
		MessageEncodingMessageEncodingUnknown,
		MessageEncodingLengthPrefixed,
		MessageEncodingDelimited,
	}
}

// Type always returns MessageEncoding to satisfy [pflag.Value] interface
func (f *MessageEncoding) Type() string {
	return "MessageEncoding"
}



type OptimizeMode string

const OptimizeModeSpeed OptimizeMode = `SPEED`
const OptimizeModeCodeSize OptimizeMode = `CODE_SIZE`
const OptimizeModeLiteRuntime OptimizeMode = `LITE_RUNTIME`

// String representation for [fmt.Print]
func (f *OptimizeMode) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *OptimizeMode) Set(v string) error {
	switch v {
	case `SPEED`, `CODE_SIZE`, `LITE_RUNTIME`:
		*f = OptimizeMode(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "SPEED", "CODE_SIZE", "LITE_RUNTIME"`, v)
	}
}

// Values returns all possible values for OptimizeMode.
//
// There is no guarantee on the order of the values in the slice.
func (f *OptimizeMode) Values() []OptimizeMode {
	return []OptimizeMode{
		OptimizeModeSpeed,
		OptimizeModeCodeSize,
		OptimizeModeLiteRuntime,
	}
}

// Type always returns OptimizeMode to satisfy [pflag.Value] interface
func (f *OptimizeMode) Type() string {
	return "OptimizeMode"
}



type OptionRetention string

const OptionRetentionRetentionUnknown OptionRetention = `RETENTION_UNKNOWN`
const OptionRetentionRetentionRuntime OptionRetention = `RETENTION_RUNTIME`
const OptionRetentionRetentionSource OptionRetention = `RETENTION_SOURCE`

// String representation for [fmt.Print]
func (f *OptionRetention) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *OptionRetention) Set(v string) error {
	switch v {
	case `RETENTION_UNKNOWN`, `RETENTION_RUNTIME`, `RETENTION_SOURCE`:
		*f = OptionRetention(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "RETENTION_UNKNOWN", "RETENTION_RUNTIME", "RETENTION_SOURCE"`, v)
	}
}

// Values returns all possible values for OptionRetention.
//
// There is no guarantee on the order of the values in the slice.
func (f *OptionRetention) Values() []OptionRetention {
	return []OptionRetention{
		OptionRetentionRetentionUnknown,
		OptionRetentionRetentionRuntime,
		OptionRetentionRetentionSource,
	}
}

// Type always returns OptionRetention to satisfy [pflag.Value] interface
func (f *OptionRetention) Type() string {
	return "OptionRetention"
}



type OptionTargetType string

const OptionTargetTypeTargetTypeUnknown OptionTargetType = `TARGET_TYPE_UNKNOWN`
const OptionTargetTypeTargetTypeFile OptionTargetType = `TARGET_TYPE_FILE`
const OptionTargetTypeTargetTypeExtensionRange OptionTargetType = `TARGET_TYPE_EXTENSION_RANGE`
const OptionTargetTypeTargetTypeMessage OptionTargetType = `TARGET_TYPE_MESSAGE`
const OptionTargetTypeTargetTypeField OptionTargetType = `TARGET_TYPE_FIELD`
const OptionTargetTypeTargetTypeOneof OptionTargetType = `TARGET_TYPE_ONEOF`
const OptionTargetTypeTargetTypeEnum OptionTargetType = `TARGET_TYPE_ENUM`
const OptionTargetTypeTargetTypeEnumEntry OptionTargetType = `TARGET_TYPE_ENUM_ENTRY`
const OptionTargetTypeTargetTypeService OptionTargetType = `TARGET_TYPE_SERVICE`
const OptionTargetTypeTargetTypeMethod OptionTargetType = `TARGET_TYPE_METHOD`

// String representation for [fmt.Print]
func (f *OptionTargetType) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *OptionTargetType) Set(v string) error {
	switch v {
	case `TARGET_TYPE_UNKNOWN`, `TARGET_TYPE_FILE`, `TARGET_TYPE_EXTENSION_RANGE`, `TARGET_TYPE_MESSAGE`, `TARGET_TYPE_FIELD`, `TARGET_TYPE_ONEOF`, `TARGET_TYPE_ENUM`, `TARGET_TYPE_ENUM_ENTRY`, `TARGET_TYPE_SERVICE`, `TARGET_TYPE_METHOD`:
		*f = OptionTargetType(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "TARGET_TYPE_UNKNOWN", "TARGET_TYPE_FILE", "TARGET_TYPE_EXTENSION_RANGE", "TARGET_TYPE_MESSAGE", "TARGET_TYPE_FIELD", "TARGET_TYPE_ONEOF", "TARGET_TYPE_ENUM", "TARGET_TYPE_ENUM_ENTRY", "TARGET_TYPE_SERVICE", "TARGET_TYPE_METHOD"`, v)
	}
}

// Values returns all possible values for OptionTargetType.
//
// There is no guarantee on the order of the values in the slice.
func (f *OptionTargetType) Values() []OptionTargetType {
	return []OptionTargetType{
		OptionTargetTypeTargetTypeUnknown,
		OptionTargetTypeTargetTypeFile,
		OptionTargetTypeTargetTypeExtensionRange,
		OptionTargetTypeTargetTypeMessage,
		OptionTargetTypeTargetTypeField,
		OptionTargetTypeTargetTypeOneof,
		OptionTargetTypeTargetTypeEnum,
		OptionTargetTypeTargetTypeEnumEntry,
		OptionTargetTypeTargetTypeService,
		OptionTargetTypeTargetTypeMethod,
	}
}

// Type always returns OptionTargetType to satisfy [pflag.Value] interface
func (f *OptionTargetType) Type() string {
	return "OptionTargetType"
}



type RepeatedFieldEncoding string

const RepeatedFieldEncodingRepeatedFieldEncodingUnknown RepeatedFieldEncoding = `REPEATED_FIELD_ENCODING_UNKNOWN`
const RepeatedFieldEncodingPacked RepeatedFieldEncoding = `PACKED`
const RepeatedFieldEncodingExpanded RepeatedFieldEncoding = `EXPANDED`

// String representation for [fmt.Print]
func (f *RepeatedFieldEncoding) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *RepeatedFieldEncoding) Set(v string) error {
	switch v {
	case `REPEATED_FIELD_ENCODING_UNKNOWN`, `PACKED`, `EXPANDED`:
		*f = RepeatedFieldEncoding(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "REPEATED_FIELD_ENCODING_UNKNOWN", "PACKED", "EXPANDED"`, v)
	}
}

// Values returns all possible values for RepeatedFieldEncoding.
//
// There is no guarantee on the order of the values in the slice.
func (f *RepeatedFieldEncoding) Values() []RepeatedFieldEncoding {
	return []RepeatedFieldEncoding{
		RepeatedFieldEncodingRepeatedFieldEncodingUnknown,
		RepeatedFieldEncodingPacked,
		RepeatedFieldEncodingExpanded,
	}
}

// Type always returns RepeatedFieldEncoding to satisfy [pflag.Value] interface
func (f *RepeatedFieldEncoding) Type() string {
	return "RepeatedFieldEncoding"
}



type Semantic string

const SemanticNone Semantic = `NONE`
const SemanticSet Semantic = `SET`
const SemanticAlias Semantic = `ALIAS`

// String representation for [fmt.Print]
func (f *Semantic) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Semantic) Set(v string) error {
	switch v {
	case `NONE`, `SET`, `ALIAS`:
		*f = Semantic(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "NONE", "SET", "ALIAS"`, v)
	}
}

// Values returns all possible values for Semantic.
//
// There is no guarantee on the order of the values in the slice.
func (f *Semantic) Values() []Semantic {
	return []Semantic{
		SemanticNone,
		SemanticSet,
		SemanticAlias,
	}
}

// Type always returns Semantic to satisfy [pflag.Value] interface
func (f *Semantic) Type() string {
	return "Semantic"
}



type SymbolVisibility string

const SymbolVisibilityVisibilityUnset SymbolVisibility = `VISIBILITY_UNSET`
const SymbolVisibilityVisibilityLocal SymbolVisibility = `VISIBILITY_LOCAL`
const SymbolVisibilityVisibilityExport SymbolVisibility = `VISIBILITY_EXPORT`

// String representation for [fmt.Print]
func (f *SymbolVisibility) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *SymbolVisibility) Set(v string) error {
	switch v {
	case `VISIBILITY_UNSET`, `VISIBILITY_LOCAL`, `VISIBILITY_EXPORT`:
		*f = SymbolVisibility(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "VISIBILITY_UNSET", "VISIBILITY_LOCAL", "VISIBILITY_EXPORT"`, v)
	}
}

// Values returns all possible values for SymbolVisibility.
//
// There is no guarantee on the order of the values in the slice.
func (f *SymbolVisibility) Values() []SymbolVisibility {
	return []SymbolVisibility{
		SymbolVisibilityVisibilityUnset,
		SymbolVisibilityVisibilityLocal,
		SymbolVisibilityVisibilityExport,
	}
}

// Type always returns SymbolVisibility to satisfy [pflag.Value] interface
func (f *SymbolVisibility) Type() string {
	return "SymbolVisibility"
}



type Type string

const TypeTypeDouble Type = `TYPE_DOUBLE`
const TypeTypeFloat Type = `TYPE_FLOAT`
const TypeTypeInt64 Type = `TYPE_INT64`
const TypeTypeUint64 Type = `TYPE_UINT64`
const TypeTypeInt32 Type = `TYPE_INT32`
const TypeTypeFixed64 Type = `TYPE_FIXED64`
const TypeTypeFixed32 Type = `TYPE_FIXED32`
const TypeTypeBool Type = `TYPE_BOOL`
const TypeTypeString Type = `TYPE_STRING`
const TypeTypeGroup Type = `TYPE_GROUP`
const TypeTypeMessage Type = `TYPE_MESSAGE`
const TypeTypeBytes Type = `TYPE_BYTES`
const TypeTypeUint32 Type = `TYPE_UINT32`
const TypeTypeEnum Type = `TYPE_ENUM`
const TypeTypeSfixed32 Type = `TYPE_SFIXED32`
const TypeTypeSfixed64 Type = `TYPE_SFIXED64`
const TypeTypeSint32 Type = `TYPE_SINT32`
const TypeTypeSint64 Type = `TYPE_SINT64`

// String representation for [fmt.Print]
func (f *Type) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Type) Set(v string) error {
	switch v {
	case `TYPE_DOUBLE`, `TYPE_FLOAT`, `TYPE_INT64`, `TYPE_UINT64`, `TYPE_INT32`, `TYPE_FIXED64`, `TYPE_FIXED32`, `TYPE_BOOL`, `TYPE_STRING`, `TYPE_GROUP`, `TYPE_MESSAGE`, `TYPE_BYTES`, `TYPE_UINT32`, `TYPE_ENUM`, `TYPE_SFIXED32`, `TYPE_SFIXED64`, `TYPE_SINT32`, `TYPE_SINT64`:
		*f = Type(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "TYPE_DOUBLE", "TYPE_FLOAT", "TYPE_INT64", "TYPE_UINT64", "TYPE_INT32", "TYPE_FIXED64", "TYPE_FIXED32", "TYPE_BOOL", "TYPE_STRING", "TYPE_GROUP", "TYPE_MESSAGE", "TYPE_BYTES", "TYPE_UINT32", "TYPE_ENUM", "TYPE_SFIXED32", "TYPE_SFIXED64", "TYPE_SINT32", "TYPE_SINT64"`, v)
	}
}

// Values returns all possible values for Type.
//
// There is no guarantee on the order of the values in the slice.
func (f *Type) Values() []Type {
	return []Type{
		TypeTypeDouble,
		TypeTypeFloat,
		TypeTypeInt64,
		TypeTypeUint64,
		TypeTypeInt32,
		TypeTypeFixed64,
		TypeTypeFixed32,
		TypeTypeBool,
		TypeTypeString,
		TypeTypeGroup,
		TypeTypeMessage,
		TypeTypeBytes,
		TypeTypeUint32,
		TypeTypeEnum,
		TypeTypeSfixed32,
		TypeTypeSfixed64,
		TypeTypeSint32,
		TypeTypeSint64,
	}
}

// Type always returns Type to satisfy [pflag.Value] interface
func (f *Type) Type() string {
	return "Type"
}



type Utf8Validation string

const Utf8ValidationUtf8ValidationUnknown Utf8Validation = `UTF8_VALIDATION_UNKNOWN`
const Utf8ValidationVerify Utf8Validation = `VERIFY`
const Utf8ValidationNone Utf8Validation = `NONE`

// String representation for [fmt.Print]
func (f *Utf8Validation) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *Utf8Validation) Set(v string) error {
	switch v {
	case `UTF8_VALIDATION_UNKNOWN`, `VERIFY`, `NONE`:
		*f = Utf8Validation(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "UTF8_VALIDATION_UNKNOWN", "VERIFY", "NONE"`, v)
	}
}

// Values returns all possible values for Utf8Validation.
//
// There is no guarantee on the order of the values in the slice.
func (f *Utf8Validation) Values() []Utf8Validation {
	return []Utf8Validation{
		Utf8ValidationUtf8ValidationUnknown,
		Utf8ValidationVerify,
		Utf8ValidationNone,
	}
}

// Type always returns Utf8Validation to satisfy [pflag.Value] interface
func (f *Utf8Validation) Type() string {
	return "Utf8Validation"
}



type VerificationState string

const VerificationStateDeclaration VerificationState = `DECLARATION`
const VerificationStateUnverified VerificationState = `UNVERIFIED`

// String representation for [fmt.Print]
func (f *VerificationState) String() string {
	return string(*f)
}

// Set raw string value and validate it against allowed values
func (f *VerificationState) Set(v string) error {
	switch v {
	case `DECLARATION`, `UNVERIFIED`:
		*f = VerificationState(v)
		return nil
	default:
		return fmt.Errorf(`value "%s" is not one of "DECLARATION", "UNVERIFIED"`, v)
	}
}

// Values returns all possible values for VerificationState.
//
// There is no guarantee on the order of the values in the slice.
func (f *VerificationState) Values() []VerificationState {
	return []VerificationState{
		VerificationStateDeclaration,
		VerificationStateUnverified,
	}
}

// Type always returns VerificationState to satisfy [pflag.Value] interface
func (f *VerificationState) Type() string {
	return "VerificationState"
}





type Annotation struct {
	
	Path []int `json:"path"`
	
	SourceFile *string `json:"sourceFile"`
	
	Begin *int `json:"begin"`
	
	End *int `json:"end"`
	
	Semantic *Semantic `json:"semantic"`
	
}



type Declaration struct {
	
	Number *int `json:"number"`
	
	FullName *string `json:"fullName"`
	
	Type *string `json:"type"`
	
	Reserved *bool `json:"reserved"`
	
	Repeated *bool `json:"repeated"`
	
}



type DescriptorProto struct {
	
	Name *string `json:"name"`
	
	Field []FieldDescriptorProto `json:"field"`
	
	Extension []FieldDescriptorProto `json:"extension"`
	
	NestedType []DescriptorProto `json:"nestedType"`
	
	EnumType []EnumDescriptorProto `json:"enumType"`
	
	ExtensionRange []ExtensionRange `json:"extensionRange"`
	
	OneofDecl []OneofDescriptorProto `json:"oneofDecl"`
	
	Options *MessageOptions `json:"options"`
	
	ReservedRange []ReservedRange `json:"reservedRange"`
	
	ReservedName []string `json:"reservedName"`
	
	Visibility *SymbolVisibility `json:"visibility"`
	
}



type EditionDefault struct {
	
	Edition *Edition `json:"edition"`
	
	Value *string `json:"value"`
	
}



type Empty struct {
	
}



type EnumDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Value []EnumValueDescriptorProto `json:"value"`
	
	Options *EnumOptions `json:"options"`
	
	ReservedRange []EnumReservedRange `json:"reservedRange"`
	
	ReservedName []string `json:"reservedName"`
	
	Visibility *SymbolVisibility `json:"visibility"`
	
}



type EnumOptions struct {
	
	AllowAlias *bool `json:"allowAlias"`
	
	Deprecated *bool `json:"deprecated"`
	
	DeprecatedLegacyJsonFieldConflicts *bool `json:"deprecatedLegacyJsonFieldConflicts"`
	
	Features *FeatureSet `json:"features"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type EnumReservedRange struct {
	
	Start *int `json:"start"`
	
	End *int `json:"end"`
	
}



type EnumValueDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Number *int `json:"number"`
	
	Options *EnumValueOptions `json:"options"`
	
}



type EnumValueOptions struct {
	
	Deprecated *bool `json:"deprecated"`
	
	Features *FeatureSet `json:"features"`
	
	DebugRedact *bool `json:"debugRedact"`
	
	FeatureSupport *FeatureSupport `json:"featureSupport"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type ExtensionRange struct {
	
	Start *int `json:"start"`
	
	End *int `json:"end"`
	
	Options *ExtensionRangeOptions `json:"options"`
	
}



type ExtensionRangeOptions struct {
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
	Declaration []Declaration `json:"declaration"`
	
	Features *FeatureSet `json:"features"`
	
	Verification *VerificationState `json:"verification"`
	
}



type FeatureSet struct {
	
	FieldPresence *FieldPresence `json:"fieldPresence"`
	
	EnumType *EnumType `json:"enumType"`
	
	RepeatedFieldEncoding *RepeatedFieldEncoding `json:"repeatedFieldEncoding"`
	
	Utf8Validation *Utf8Validation `json:"utf8Validation"`
	
	MessageEncoding *MessageEncoding `json:"messageEncoding"`
	
	JsonFormat *JsonFormat `json:"jsonFormat"`
	
	EnforceNamingStyle *EnforceNamingStyle `json:"enforceNamingStyle"`
	
	DefaultSymbolVisibility *DefaultSymbolVisibility `json:"defaultSymbolVisibility"`
	
}



type FeatureSetDefaults struct {
	
	Defaults []FeatureSetEditionDefault `json:"defaults"`
	
	MinimumEdition *Edition `json:"minimumEdition"`
	
	MaximumEdition *Edition `json:"maximumEdition"`
	
}



type FeatureSetEditionDefault struct {
	
	Edition *Edition `json:"edition"`
	
	OverridableFeatures *FeatureSet `json:"overridableFeatures"`
	
	FixedFeatures *FeatureSet `json:"fixedFeatures"`
	
}



type FeatureSupport struct {
	
	EditionIntroduced *Edition `json:"editionIntroduced"`
	
	EditionDeprecated *Edition `json:"editionDeprecated"`
	
	DeprecationWarning *string `json:"deprecationWarning"`
	
	EditionRemoved *Edition `json:"editionRemoved"`
	
}



type FieldDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Number *int `json:"number"`
	
	Label *Label `json:"label"`
	
	Type *Type `json:"type"`
	
	TypeName *string `json:"typeName"`
	
	Extendee *string `json:"extendee"`
	
	DefaultValue *string `json:"defaultValue"`
	
	OneofIndex *int `json:"oneofIndex"`
	
	JsonName *string `json:"jsonName"`
	
	Options *FieldOptions `json:"options"`
	
	Proto3Optional *bool `json:"proto3Optional"`
	
}



type FieldMask struct {
	
	Paths []string `json:"paths"`
	
}



type FieldOptions struct {
	
	Ctype *CType `json:"ctype"`
	
	Packed *bool `json:"packed"`
	
	Jstype *JsType `json:"jstype"`
	
	Lazy *bool `json:"lazy"`
	
	UnverifiedLazy *bool `json:"unverifiedLazy"`
	
	Deprecated *bool `json:"deprecated"`
	
	Weak *bool `json:"weak"`
	
	DebugRedact *bool `json:"debugRedact"`
	
	Retention *OptionRetention `json:"retention"`
	
	Targets []OptionTargetType `json:"targets"`
	
	EditionDefaults []EditionDefault `json:"editionDefaults"`
	
	Features *FeatureSet `json:"features"`
	
	FeatureSupport *FeatureSupport `json:"featureSupport"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type FileDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Package *string `json:"package"`
	
	Dependency []string `json:"dependency"`
	
	PublicDependency []int `json:"publicDependency"`
	
	WeakDependency []int `json:"weakDependency"`
	
	OptionDependency []string `json:"optionDependency"`
	
	MessageType []DescriptorProto `json:"messageType"`
	
	EnumType []EnumDescriptorProto `json:"enumType"`
	
	Service []ServiceDescriptorProto `json:"service"`
	
	Extension []FieldDescriptorProto `json:"extension"`
	
	Options *FileOptions `json:"options"`
	
	SourceCodeInfo *SourceCodeInfo `json:"sourceCodeInfo"`
	
	Syntax *string `json:"syntax"`
	
	Edition *Edition `json:"edition"`
	
}



type FileDescriptorSet struct {
	
	File []FileDescriptorProto `json:"file"`
	
}



type FileOptions struct {
	
	JavaPackage *string `json:"javaPackage"`
	
	JavaOuterClassname *string `json:"javaOuterClassname"`
	
	JavaMultipleFiles *bool `json:"javaMultipleFiles"`
	
	JavaGenerateEqualsAndHash *bool `json:"javaGenerateEqualsAndHash"`
	
	JavaStringCheckUtf8 *bool `json:"javaStringCheckUtf8"`
	
	OptimizeFor *OptimizeMode `json:"optimizeFor"`
	
	GoPackage *string `json:"goPackage"`
	
	CcGenericServices *bool `json:"ccGenericServices"`
	
	JavaGenericServices *bool `json:"javaGenericServices"`
	
	PyGenericServices *bool `json:"pyGenericServices"`
	
	Deprecated *bool `json:"deprecated"`
	
	CcEnableArenas *bool `json:"ccEnableArenas"`
	
	ObjcClassPrefix *string `json:"objcClassPrefix"`
	
	CsharpNamespace *string `json:"csharpNamespace"`
	
	SwiftPrefix *string `json:"swiftPrefix"`
	
	PhpClassPrefix *string `json:"phpClassPrefix"`
	
	PhpNamespace *string `json:"phpNamespace"`
	
	PhpMetadataNamespace *string `json:"phpMetadataNamespace"`
	
	RubyPackage *string `json:"rubyPackage"`
	
	Features *FeatureSet `json:"features"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type GeneratedCodeInfo struct {
	
	Annotation []Annotation `json:"annotation"`
	
}



type Location struct {
	
	Path []int `json:"path"`
	
	Span []int `json:"span"`
	
	LeadingComments *string `json:"leadingComments"`
	
	TrailingComments *string `json:"trailingComments"`
	
	LeadingDetachedComments []string `json:"leadingDetachedComments"`
	
}



type MessageOptions struct {
	
	MessageSetWireFormat *bool `json:"messageSetWireFormat"`
	
	NoStandardDescriptorAccessor *bool `json:"noStandardDescriptorAccessor"`
	
	Deprecated *bool `json:"deprecated"`
	
	MapEntry *bool `json:"mapEntry"`
	
	DeprecatedLegacyJsonFieldConflicts *bool `json:"deprecatedLegacyJsonFieldConflicts"`
	
	Features *FeatureSet `json:"features"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type MethodDescriptorProto struct {
	
	Name *string `json:"name"`
	
	InputType *string `json:"inputType"`
	
	OutputType *string `json:"outputType"`
	
	Options *MethodOptions `json:"options"`
	
	ClientStreaming *bool `json:"clientStreaming"`
	
	ServerStreaming *bool `json:"serverStreaming"`
	
}



type MethodOptions struct {
	
	Deprecated *bool `json:"deprecated"`
	
	IdempotencyLevel *IdempotencyLevel `json:"idempotencyLevel"`
	
	Features *FeatureSet `json:"features"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type NamePart struct {
	
	NamePart *string `json:"namePart"`
	
	IsExtension *bool `json:"isExtension"`
	
}



type OneofDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Options *OneofOptions `json:"options"`
	
}



type OneofOptions struct {
	
	Features *FeatureSet `json:"features"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type ReservedRange struct {
	
	Start *int `json:"start"`
	
	End *int `json:"end"`
	
}



type ServiceDescriptorProto struct {
	
	Name *string `json:"name"`
	
	Method []MethodDescriptorProto `json:"method"`
	
	Options *ServiceOptions `json:"options"`
	
}



type ServiceOptions struct {
	
	Features *FeatureSet `json:"features"`
	
	Deprecated *bool `json:"deprecated"`
	
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
	
}



type SourceCodeInfo struct {
	
	Location []Location `json:"location"`
	
}



type UninterpretedOption struct {
	
	Name []NamePart `json:"name"`
	
	IdentifierValue *string `json:"identifierValue"`
	
	PositiveIntValue *int64 `json:"positiveIntValue"`
	
	NegativeIntValue *int64 `json:"negativeIntValue"`
	
	DoubleValue *float64 `json:"doubleValue"`
	
	StringValue *io.ReadCloser `json:"stringValue"`
	
	AggregateValue *string `json:"aggregateValue"`
	
}



type VisibilityFeature struct {
	
}


