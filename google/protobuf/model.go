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

type Label string

const (
	LabelLabelOptional Label = "LABEL_OPTIONAL"
	LabelLabelRequired Label = "LABEL_REQUIRED"
	LabelLabelRepeated Label = "LABEL_REPEATED"
)

// String representation for [fmt.Print].
func (f *Label) String() string {
	return string(*f)
}

type NullValue string

const (
	NullValueNullValue NullValue = "NULL_VALUE"
)

// String representation for [fmt.Print].
func (f *NullValue) String() string {
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

type Annotation struct {
	Path []int `json:"path"`
	SourceFile *string `json:"sourceFile"`
	Begin *int `json:"begin"`
	End *int `json:"end"`
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
}

type Empty struct {
}

type EnumDescriptorProto struct {
	Name *string `json:"name"`
	Value []EnumValueDescriptorProto `json:"value"`
	Options *EnumOptions `json:"options"`
	ReservedRange []EnumReservedRange `json:"reservedRange"`
	ReservedName []string `json:"reservedName"`
}

type EnumOptions struct {
	AllowAlias *bool `json:"allowAlias"`
	Deprecated *bool `json:"deprecated"`
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
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
}

type ExtensionRange struct {
	Start *int `json:"start"`
	End *int `json:"end"`
	Options *ExtensionRangeOptions `json:"options"`
}

type ExtensionRangeOptions struct {
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
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
	Deprecated *bool `json:"deprecated"`
	Weak *bool `json:"weak"`
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
}

type FieldsEntry struct {
	Key *string `json:"key"`
	Value *json.RawMessage `json:"value"`
}

type FileDescriptorProto struct {
	Name *string `json:"name"`
	Package *string `json:"package"`
	Dependency []string `json:"dependency"`
	PublicDependency []int `json:"publicDependency"`
	WeakDependency []int `json:"weakDependency"`
	MessageType []DescriptorProto `json:"messageType"`
	EnumType []EnumDescriptorProto `json:"enumType"`
	Service []ServiceDescriptorProto `json:"service"`
	Extension []FieldDescriptorProto `json:"extension"`
	Options *FileOptions `json:"options"`
	SourceCodeInfo *SourceCodeInfo `json:"sourceCodeInfo"`
	Syntax *string `json:"syntax"`
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
	PhpGenericServices *bool `json:"phpGenericServices"`
	Deprecated *bool `json:"deprecated"`
	CcEnableArenas *bool `json:"ccEnableArenas"`
	ObjcClassPrefix *string `json:"objcClassPrefix"`
	CsharpNamespace *string `json:"csharpNamespace"`
	SwiftPrefix *string `json:"swiftPrefix"`
	PhpClassPrefix *string `json:"phpClassPrefix"`
	PhpNamespace *string `json:"phpNamespace"`
	PhpMetadataNamespace *string `json:"phpMetadataNamespace"`
	RubyPackage *string `json:"rubyPackage"`
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
}

type GeneratedCodeInfo struct {
	Annotation []Annotation `json:"annotation"`
}

type ListValue struct {
	Values []json.RawMessage `json:"values"`
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
	Deprecated *bool `json:"deprecated"`
	UninterpretedOption []UninterpretedOption `json:"uninterpretedOption"`
}

type SourceCodeInfo struct {
	Location []Location `json:"location"`
}

type Struct struct {
	Fields map[string]json.RawMessage `json:"fields"`
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

type Value struct {
	NullValue *NullValue `json:"nullValue"`
	NumberValue *float64 `json:"numberValue"`
	StringValue *string `json:"stringValue"`
	BoolValue *bool `json:"boolValue"`
	StructValue *json.RawMessage `json:"structValue"`
	ListValue *json.RawMessage `json:"listValue"`
}
