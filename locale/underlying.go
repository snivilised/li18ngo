package locale

import "github.com/snivilised/li18ngo/locale/enums"

type UnderlyingField struct {
	// Note must match a {{.<Note>}} token in the Other string exactly.
	Note string

	// GoType must be a valid native Go type (e.g. "string", "int", "uint",
	// "error").
	GoType string

	// Tale is the doc comment emitted for this field in the generated struct.
	// If Tale is empty a 🔥 TODO reminder is emitted instead.
	Tale string
}

// UnderlyingTemplData is the descriptor for a single i18n message.
// Populate one entry per message in the Underliers map below, then
// run go generate to produce the auto files.
type UnderlyingTemplData struct {
	// MessageID is the go-i18n message ID. Must be unique across all entries.
	MessageID string

	// Seed is the PascalCase base name used to derive all generated
	// identifiers: XxxTemplData, XxxError, ErrXxx, NewXxxError.
	Seed string

	// Type controls which code is generated. See the guide above.
	TypeName enums.UnderlyingType

	// Description is the go-i18n message description (a short human
	// summary) and is also used as the struct-level doc comment in
	// generated code. If empty, a 🔥 TODO reminder is emitted instead.
	Description string

	// Story is inserted into the generated banner comment block as the
	// overall narrative for the message. Long stories are word-wrapped
	// at 80 characters automatically. If empty, a 🔥 TODO reminder is
	// emitted instead.
	Story string

	// Other is the go-i18n Other translation string. May contain
	// {{.FieldName}} tokens for dynamic messages. Each token must have
	// a matching Fields entry.
	Other string

	// Fields lists the variable fields for dynamic messages. Must be
	// empty for static types and non-empty for dynamic types. For
	// wrapper types exactly one entry must have GoType "error" and
	// Name "Wrapped".
	Fields []UnderlyingField
}

// Underliers is the map type read by i18n-gen at code-generation time.
// The map key must equal the MessageID field of the value.
type Underliers map[string]UnderlyingTemplData

// underliers is the single source of truth for all i18n messages in this
// package. Edit this map and run go generate to regenerate the auto files.
var underliers = Underliers{
	//
}

func init() {
	// Prevent the underliers variable from being flagged as unused by the
	// Go compiler. i18n-gen reads this variable from the AST at generate
	// time; it is not referenced at runtime.
	_ = underliers
}
