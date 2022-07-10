package compiler

import "github.com/NYTimes/openapi2proto/internal/option"

const (
	optkeyAnnotation         = "annotation"
	optkeyAllOfFieldPrefix   = "allof-prefix"
	optkeySkipRpcs           = "skip-rpcs"
	optKeySkipDeprecatedRpcs = "skip-deprecated-rpcs"
	optkeyPrefixEnums        = "namespace-enums"
	optkeyWrapPrimitives     = "wrap-primitives"
)

// WithAnnotation creates a new Option to specify if we should add
// google.api.http annotation to the compiled Protocol Buffers structure
func WithAnnotation(b bool) Option {
	return option.New(optkeyAnnotation, b)
}

// WithAllOfFieldPrefix creates a new Option to specify the prefix to use when create fields name of
// types declared as allOf
func WithAllOfFieldPrefix(prefix string) Option {
	return option.New(optkeyAllOfFieldPrefix, prefix)
}

// WithSkipRpcs creates a new Option to specify if we should
// generate services and rpcs in addition to messages
func WithSkipRpcs(b bool) Option {
	return option.New(optkeySkipRpcs, b)
}

// WithSkipDeprecatedRpcs creates a new Option to specify if we should
// skip generating rpcs for endpoints marked as deprecated
func WithSkipDeprecatedRpcs(b bool) Option {
	return option.New(optKeySkipDeprecatedRpcs, b)
}

// prefix enum values with their enum name to prevent protobuf namespacing issues
func WithPrefixEnums(b bool) Option {
	return option.New(optkeyPrefixEnums, b)
}

// wrap primitive types with their wrapper message types
// see https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto
// and https://developers.google.com/protocol-buffers/docs/proto3#default
func WithWrapPrimitives(b bool) Option {
	return option.New(optkeyWrapPrimitives, b)
}
