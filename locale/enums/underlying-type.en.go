package enums

//go:generate stringer -type=UnderlyingType -linecomment -trimprefix=UnderlyingType -output underlying-type-en-auto.go

// UnderlyingType identifies the kind of message and controls code generation.
type UnderlyingType uint

// =============================================================================
//
// # UnderlyingType validation rules
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
	// Every {{.Token}} in Other must have a matching Fields entry and
	// vice versa.
	// Fields must be non-empty.
	// Generates:
	// - NewXxxTemplData constructor generated.
	UnderlyingTypeDynamicCobra // DynamicCobra

	// UnderlyingTypeStaticGeneral is a static non-error user-facing message.
	UnderlyingTypeStaticGeneral // StaticGeneral

	// UnderlyingTypeDynamicGeneral is a dynamic non-error user-facing message.
	// Every {{.Token}} in Other must have a matching Fields entry and
	// vice versa.
	// Fields must be non-empty.
	// Generates:
	// - NewXxxTemplData constructor generated.
	UnderlyingTypeDynamicGeneral // DynamicGeneral

	// UnderlyingTypeStaticError is a static error with no variable content.
	UnderlyingTypeStaticError // StaticError

	// UnderlyingTypeSentinelError is a static sentinel error designed to be
	// wrapped by outer errors.
	// Generates:
	// - XxxError
	// ErrXxx sentinel.
	UnderlyingTypeSentinelError // SentinelError

	// UnderlyingTypeStaticErrorWrapper is a static error that wraps another
	// error for Go's error chain (errors.Is/errors.As) only. The localised
	// message text is fully fixed; the wrapped error's text does not appear
	// in the translated output. Use UnderlyingTypeStaticErrorWrapperMsg when
	// you want {{.Wrapped}} to appear inside the Other string.
	// Generates:
	// - XxxError
	UnderlyingTypeStaticErrorWrapper // StaticErrorWrapper

	// UnderlyingTypeStaticErrorWrapperMsg is a static error that wraps another
	// error and includes the wrapped error's message text directly in the
	// translated output via {{.Wrapped}} in Other. If the message text is
	// fully fixed and you only need the wrapped error for the error chain,
	// use UnderlyingTypeStaticErrorWrapper instead.
	// Generates:
	// - NewXxxError(wrapped error)
	// - Error() string
	// - Unwrap() error
	UnderlyingTypeStaticErrorWrapperMsg // StaticErrorWrapperMsg

	// UnderlyingTypeDynamicError is a dynamic error with no wrapping.
	// Fields must be non-empty. No Wrapped field permitted in Fields.
	// Generates:
	// - NewXxxError
	UnderlyingTypeDynamicError // DynamicError

	// UnderlyingTypeDynamicErrorWrapper is a dynamic error that wraps
	// another error.
	// Fields must be non-empty.
	// Generates:
	// - NewXxxError(wrapped error)
	// - Error() string
	// - Unwrap() error
	UnderlyingTypeDynamicErrorWrapper // DynamicErrorWrapper
)
