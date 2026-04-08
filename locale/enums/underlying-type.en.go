package enums

//go:generate stringer -type=UnderlyingType -linecomment -trimprefix=UnderlyingType -output underlying-type-en-auto.go

// UnderlyingType identifies the kind of message and controls code generation.
type UnderlyingType uint

// =============================================================================
//
// # UnderlyingType guide
//
// Each entry in the Underliers map must set a TypeName field. The type controls
// which code is generated for that message. The rules are:
//
//	UnderlyingTypeCobraStatic
//	  A short description string for a Cobra command or flag.
//	  Fields must be empty. No constructor generated.
//
//	UnderlyingTypeCobraDynamic
//	  A long description string for a Cobra command or flag.
//	  Fields must be non-empty. NewXxxTemplData constructor generated.
//	  Every {{.Token}} in Other must have a matching Fields entry and
//	  vice versa.
//
//	UnderlyingTypeGeneralStatic
//	  A non-error user-facing message with no variable content.
//	  Fields must be empty. No constructor generated.
//
//	UnderlyingTypeGeneralDynamic
//	  A non-error user-facing message with variable content.
//	  Fields must be non-empty. NewXxxTemplData constructor generated.
//	  Every {{.Token}} in Other must have a matching Fields entry and
//	  vice versa.
//
//	UnderlyingTypeErrorStatic
//	  An error with no variable content.
//	  Fields must be empty.
//	  Generates: XxxErrorTemplData, XxxError, ErrXxx sentinel.
//
//	UnderlyingTypeErrorCore
//	  A static sentinel error designed to be wrapped by outer errors.
//	  Fields must be empty.
//	  Generates: XxxErrorTemplData, XxxError, ErrXxx (exported sentinel).
//	  Callers use errors.Is(err, locale.ErrXxx) directly.
//
//	UnderlyingTypeErrorStaticWrapper
//	  A static error that wraps another error.
//	  Fields must be empty (Wrapped is implicit, not declared in Fields).
//	  Generates: XxxErrorTemplData, XxxError{Wrapped error},
//	             NewXxxError(wrapped error), AsXxxError helper,
//	             Error() string, Unwrap() error.
//
//	UnderlyingTypeErrorDynamic
//	  An error with variable content but no wrapping.
//	  Fields must be non-empty. No Wrapped field permitted in Fields.
//	  Generates: XxxTemplData, XxxError, NewXxxError(fields...),
//	             AsXxxError helper.
//
//	UnderlyingTypeErrorDynamicWrapper
//	  An error with variable content that also wraps another error.
//	  Fields must be non-empty and must contain exactly one entry with
//	  GoType "error" and Name "Wrapped". The Wrapped field may appear as
//	  {{.Wrapped}} in Other to control placement in the error message.
//	  The constructor always takes wrapped error as its first parameter.
//	  Generates: XxxTemplData (Wrapped as string), XxxError (Wrapped as
//	             error), NewXxxError(wrapped error, fields...),
//	             AsXxxError helper, Error() string, Unwrap() error.
//
// # Validation
//
// lingo validates the entire Underliers map before generating any files.
// Generation only proceeds when zero errors are found. Detected errors:
//
//   - Fields non-empty when TypeName declares static
//   - Fields empty when TypeName declares dynamic
//   - {{.Token}} in Other with no matching Fields entry
//   - Fields entry with no matching {{.Token}} in Other
//   - More than one Fields entry with GoType "error"
//   - Fields entry with GoType "error" and Name != "Wrapped"
//   - Fields entry with GoType "error" on a non-wrapper TypeName
//   - Fields non-empty on UnderlyingTypeErrorStaticWrapper
//   - {{.Wrapped}} in Other on a non-wrapper TypeName
//   - Duplicate MessageID across the map
//
// =============================================================================
const (
	// UnderlyingTypeUndefined is the zero value; always an error if seen.
	UnderlyingTypeUndefined UnderlyingType = iota // Undefined

	// UnderlyingTypeStaticCobra is a static Cobra command/flag
	// description with no variable content.
	UnderlyingTypeStaticCobra // StaticCobra

	// UnderlyingTypeDynamicCobra is a dynamic Cobra command/flag
	// description with variable content.
	UnderlyingTypeDynamicCobra // DynamicCobra

	// UnderlyingTypeStaticGeneral is a static non-error user-facing message.
	UnderlyingTypeStaticGeneral // StaticGeneral

	// UnderlyingTypeDynamicGeneral is a dynamic non-error user-facing message.
	UnderlyingTypeDynamicGeneral // DynamicGeneral

	// UnderlyingTypeStaticError is a static error with no variable content.
	UnderlyingTypeStaticError // StaticError

	// UnderlyingTypeSentinelError is a static sentinel error designed to be
	// wrapped by outer errors.
	UnderlyingTypeSentinelError // SentinelError

	// UnderlyingTypeStaticErrorWrapper is a static error that wraps
	// another error.
	UnderlyingTypeStaticErrorWrapper // StaticErrorWrapper

	// UnderlyingTypeDynamicError is a dynamic error with no wrapping.
	UnderlyingTypeDynamicError // DynamicError

	// UnderlyingTypeDynamicErrorWrapper is a dynamic error that wraps
	// another error.
	UnderlyingTypeDynamicErrorWrapper // DynamicErrorWrapper
)
