package protobuf

import (
	"fmt"
	"io"
)

type CType string

const (
	CTypeString CType = "STRING"
	CTypeCord CType = "CORD"
	CTypeStringPiece CType = "STRING_PIECE"
)

// String representation for [fmt.Print].
func (f *CType) String() string {
	return string(*f)
}

type DefaultSymbolVisibility string

const (
	DefaultSymbolVisibilityDefaultSymbolVisibilityUnknown DefaultSymbolVisibility = "DEFAULT_SYMBOL_VISIBILITY_UNKNOWN"
	DefaultSymbolVisibilityExportAll DefaultSymbolVisibility = "EXPORT_ALL"
	DefaultSymbolVisibilityExportTopLevel DefaultSymbolVisibility = "EXPORT_TOP_LEVEL"
	DefaultSymbolVisibilityLocalAll DefaultSymbolVisibility = "LOCAL_ALL"
	DefaultSymbolVisibilityStrict DefaultSymbolVisibility = "STRICT"
)

// String representation for [fmt.Print].
func (f *DefaultSymbolVisibility) String() string {
	return string(*f)
}

type Edition string

const (
	EditionEditionUnknown Edition = "EDITION_UNKNOWN"
	EditionEditionLegacy Edition = "EDITION_LEGACY"
	EditionEditionProto2 Edition = "EDITION_PROTO2"
	EditionEditionProto3 Edition = "EDITION_PROTO3"
	EditionEdition2023 Edition = "EDITION_2023"
	EditionEdition2024 Edition = "EDITION_2024"
	EditionEdition1TestOnly Edition = "EDITION_1_TEST_ONLY"
	EditionEdition2TestOnly Edition = "EDITION_2_TEST_ONLY"
	EditionEdition99997TestOnly Edition = "EDITION_99997_TEST_ONLY"
	EditionEdition99998TestOnly Edition = "EDITION_99998_TEST_ONLY"
	EditionEdition99999TestOnly Edition = "EDITION_99999_TEST_ONLY"
	EditionEditionMax Edition = "EDITION_MAX"
)

// String representation for [fmt.Print].
func (f *Edition) String() string {
	return string(*f)
}

type EnforceNamingStyle string

const (
	EnforceNamingStyleEnforceNamingStyleUnknown EnforceNamingStyle = "ENFORCE_NAMING_STYLE_UNKNOWN"
	EnforceNamingStyleStyle2024 EnforceNamingStyle = "STYLE2024"
	EnforceNamingStyleStyleLegacy EnforceNamingStyle = "STYLE_LEGACY"
)

// String representation for [fmt.Print].
func (f *EnforceNamingStyle) String() string {
	return string(*f)
}

type EnumType string

const (
	EnumTypeEnumTypeUnknown EnumType = "ENUM_TYPE_UNKNOWN"
	EnumTypeOpen EnumType = "OPEN"
	EnumTypeClosed EnumType = "CLOSED"
)

// String representation for [fmt.Print].
func (f *EnumType) String() string {
	return string(*f)
}

type FieldPresence string

const (
	FieldPresenceFieldPresenceUnknown FieldPresence = "FIELD_PRESENCE_UNKNOWN"
	FieldPresenceExplicit FieldPresence = "EXPLICIT"
	FieldPresenceImplicit FieldPresence = "IMPLICIT"
	FieldPresenceLegacyRequired FieldPresence = "LEGACY_REQUIRED"
)

// String representation for [fmt.Print].
func (f *FieldPresence) String() string {
	return string(*f)
}

type IdempotencyLevel string

const (
	IdempotencyLevelIdempotencyUnknown IdempotencyLevel = "IDEMPOTENCY_UNKNOWN"
	IdempotencyLevelNoSideEffects IdempotencyLevel = "NO_SIDE_EFFECTS"
	IdempotencyLevelIdempotent IdempotencyLevel = "IDEMPOTENT"
)

// String representation for [fmt.Print].
func (f *IdempotencyLevel) String() string {
	return string(*f)
}

type JsType string

const (
	JsTypeJsNormal JsType = "JS_NORMAL"
	JsTypeJsString JsType = "JS_STRING"
	JsTypeJsNumber JsType = "JS_NUMBER"
)

// String representation for [fmt.Print].
func (f *JsType) String() string {
	return string(*f)
}

type JsonFormat string

const (
	JsonFormatJsonFormatUnknown JsonFormat = "JSON_FORMAT_UNKNOWN"
	JsonFormatAllow JsonFormat = "ALLOW"
	JsonFormatLegacyBestEffort JsonFormat = "LEGACY_BEST_EFFORT"
)

// String representation for [fmt.Print].
func (f *JsonFormat) String() string {
	return string(*f)
}

type Label string

const (
	LabelLabelOptional Label = "LABEL_OPTIONAL"
	LabelLabelRepeated Label = "LABEL_REPEATED"
	LabelLabelRequired Label = "LABEL_REQUIRED"
)

// String representation for [fmt.Print].
func (f *Label) String() string {
	return string(*f)
}

type MessageEncoding string

const (
	MessageEncodingMessageEncodingUnknown MessageEncoding = "MESSAGE_ENCODING_UNKNOWN"
	MessageEncodingLengthPrefixed MessageEncoding = "LENGTH_PREFIXED"
	MessageEncodingDelimited MessageEncoding = "DELIMITED"
)

// String representation for [fmt.Print].
func (f *MessageEncoding) String() string {
	return string(*f)
}

type OptimizeMode string

const (
	OptimizeModeSpeed OptimizeMode = "SPEED"
	OptimizeModeCodeSize OptimizeMode = "CODE_SIZE"
	OptimizeModeLiteRuntime OptimizeMode = "LITE_RUNTIME"
)

// String representation for [fmt.Print].
func (f *OptimizeMode) String() string {
	return string(*f)
}

type OptionRetention string

const (
	OptionRetentionRetentionUnknown OptionRetention = "RETENTION_UNKNOWN"
	OptionRetentionRetentionRuntime OptionRetention = "RETENTION_RUNTIME"
	OptionRetentionRetentionSource OptionRetention = "RETENTION_SOURCE"
)

// String representation for [fmt.Print].
func (f *OptionRetention) String() string {
	return string(*f)
}

type OptionTargetType string

const (
	OptionTargetTypeTargetTypeUnknown OptionTargetType = "TARGET_TYPE_UNKNOWN"
	OptionTargetTypeTargetTypeFile OptionTargetType = "TARGET_TYPE_FILE"
	OptionTargetTypeTargetTypeExtensionRange OptionTargetType = "TARGET_TYPE_EXTENSION_RANGE"
	OptionTargetTypeTargetTypeMessage OptionTargetType = "TARGET_TYPE_MESSAGE"
	OptionTargetTypeTargetTypeField OptionTargetType = "TARGET_TYPE_FIELD"
	OptionTargetTypeTargetTypeOneof OptionTargetType = "TARGET_TYPE_ONEOF"
	OptionTargetTypeTargetTypeEnum OptionTargetType = "TARGET_TYPE_ENUM"
	OptionTargetTypeTargetTypeEnumEntry OptionTargetType = "TARGET_TYPE_ENUM_ENTRY"
	OptionTargetTypeTargetTypeService OptionTargetType = "TARGET_TYPE_SERVICE"
	OptionTargetTypeTargetTypeMethod OptionTargetType = "TARGET_TYPE_METHOD"
)

// String representation for [fmt.Print].
func (f *OptionTargetType) String() string {
	return string(*f)
}

type RepeatedFieldEncoding string

const (
	RepeatedFieldEncodingRepeatedFieldEncodingUnknown RepeatedFieldEncoding = "REPEATED_FIELD_ENCODING_UNKNOWN"
	RepeatedFieldEncodingPacked RepeatedFieldEncoding = "PACKED"
	RepeatedFieldEncodingExpanded RepeatedFieldEncoding = "EXPANDED"
)

// String representation for [fmt.Print].
func (f *RepeatedFieldEncoding) String() string {
	return string(*f)
}

type Semantic string

const (
	SemanticNone Semantic = "NONE"
	SemanticSet Semantic = "SET"
	SemanticAlias Semantic = "ALIAS"
)

// String representation for [fmt.Print].
func (f *Semantic) String() string {
	return string(*f)
}

type SymbolVisibility string

const (
	SymbolVisibilityVisibilityUnset SymbolVisibility = "VISIBILITY_UNSET"
	SymbolVisibilityVisibilityLocal SymbolVisibility = "VISIBILITY_LOCAL"
	SymbolVisibilityVisibilityExport SymbolVisibility = "VISIBILITY_EXPORT"
)

// String representation for [fmt.Print].
func (f *SymbolVisibility) String() string {
	return string(*f)
}

type Type string

const (
	TypeTypeDouble Type = "TYPE_DOUBLE"
	TypeTypeFloat Type = "TYPE_FLOAT"
	TypeTypeInt64 Type = "TYPE_INT64"
	TypeTypeUint64 Type = "TYPE_UINT64"
	TypeTypeInt32 Type = "TYPE_INT32"
	TypeTypeFixed64 Type = "TYPE_FIXED64"
	TypeTypeFixed32 Type = "TYPE_FIXED32"
	TypeTypeBool Type = "TYPE_BOOL"
	TypeTypeString Type = "TYPE_STRING"
	TypeTypeGroup Type = "TYPE_GROUP"
	TypeTypeMessage Type = "TYPE_MESSAGE"
	TypeTypeBytes Type = "TYPE_BYTES"
	TypeTypeUint32 Type = "TYPE_UINT32"
	TypeTypeEnum Type = "TYPE_ENUM"
	TypeTypeSfixed32 Type = "TYPE_SFIXED32"
	TypeTypeSfixed64 Type = "TYPE_SFIXED64"
	TypeTypeSint32 Type = "TYPE_SINT32"
	TypeTypeSint64 Type = "TYPE_SINT64"
)

// String representation for [fmt.Print].
func (f *Type) String() string {
	return string(*f)
}

type Utf8Validation string

const (
	Utf8ValidationUtf8ValidationUnknown Utf8Validation = "UTF8_VALIDATION_UNKNOWN"
	Utf8ValidationVerify Utf8Validation = "VERIFY"
	Utf8ValidationNone Utf8Validation = "NONE"
)

// String representation for [fmt.Print].
func (f *Utf8Validation) String() string {
	return string(*f)
}

type VerificationState string

const (
	VerificationStateDeclaration VerificationState = "DECLARATION"
	VerificationStateUnverified VerificationState = "UNVERIFIED"
)

// String representation for [fmt.Print].
func (f *VerificationState) String() string {
	return string(*f)
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
