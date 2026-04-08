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
	// -------------------------------------------------------------------------
	// General messages
	// -------------------------------------------------------------------------
	"using-config-file": {
		MessageID:   "using-config-file",
		Seed:        "UsingConfigFile",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Message to indicate which config is being used",
		Story: "UsingConfigFile is printed on startup to indicate" +
			" which configuration file has been loaded.",
		Other: "Using config file: '{{.ConfigFileName}}'",
		Fields: []UnderlyingField{
			{
				Note:   "ConfigFileName",
				GoType: "string",
				Tale:   "is the name of the config file being used",
			},
		},
	},

	"localisation.test": {
		MessageID:   "localisation.test",
		Seed:        "Localisation",
		TypeName:    enums.UnderlyingTypeStaticGeneral,
		Description: "Localisation",
		Story:       "A test message for localisation.",
		Other:       "localisation",
	},

	"internationalisation.test": {
		MessageID:   "internationalisation.test",
		Seed:        "Internationalisation",
		TypeName:    enums.UnderlyingTypeStaticGeneral,
		Description: "Internationalisation",
		Story:       "A test message for internationalisation.",
		Other:       "internationalisation",
	},

	"pavement-graffiti-report.graffico.test": {
		MessageID:   "pavement-graffiti-report.graffico.test",
		Seed:        "PavementGraffitiReportGrafficoTest",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Report of graffiti found on a pavement",
		Story:       "A test message for pavement graffiti reporting.",
		Other:       "Found graffiti on pavement; primary colour: '{{.Primary}}'",
		Fields: []UnderlyingField{
			{
				Note:   "Primary",
				GoType: "string",
				Tale:   "is the primary colour of the graffiti",
			},
		},
	},

	"wrong-source-id.graffico.test": {
		MessageID:   "wrong-source-id.graffico.test",
		Seed:        "WrongSourceIDGrafficoTest",
		TypeName:    enums.UnderlyingTypeStaticGeneral,
		Description: "Message with wrong source ID for testing purposes",
		Story:       "This message has the wrong source ID and should be ignored by i18n-gen.",
		Other:       "This message should be ignored by i18n-gen.",
	},

	// -------------------------------------------------------------------------
	// Error messages
	// -------------------------------------------------------------------------

	"path-not-found.dynamic-error": {
		MessageID:   "path-not-found.dynamic-error",
		Seed:        "PathNotFound",
		TypeName:    enums.UnderlyingTypeDynamicError,
		Description: "Directory or file path does not exist",
		Story:       "PathNotFoundError is used when a directory or file path does not exist.",
		Other:       "{{.Name}} path not found ({{.Path}})",
		Fields: []UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "is the name of the path that was not found (e.g. 'Config')",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the actual path that was not found (e.g. '/etc/config.yaml')",
			},
		},
	},

	"not-a-directory.dynamic-error": {
		MessageID:   "not-a-directory.dynamic-error",
		Seed:        "NotADirectory",
		TypeName:    enums.UnderlyingTypeDynamicError,
		Description: "File system path is not a directory",
		Story:       "NotADirectoryError is used when a file system path is expected to be a directory but is not.",
		Other:       "file system path '{{.Path}}', is not a directory",
		Fields: []UnderlyingField{
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the file system path that was expected to be a directory but is not (e.g. '/etc/config.yaml')",
			},
		},
	},

	"third-party.error-wrapper": {
		MessageID:   "third-party.error-wrapper",
		Seed:        "ThirdPartyErrorWrapper",
		TypeName:    enums.UnderlyingTypeStaticErrorWrapper,
		Description: "Wrapper for third-party errors",
		Story:       "ThirdPartyErrorWrapper is used to wrap errors from third-party libraries.",
		Other:       "An error occurred in a third-party library: {{.Wrapped}}",
		Fields: []UnderlyingField{
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "is the original error from the third-party library that is being wrapped",
			},
		},
	},
}

func init() {
	// Prevent the underliers variable from being flagged as unused by the
	// Go compiler. i18n-gen reads this variable from the AST at generate
	// time; it is not referenced at runtime.
	_ = underliers
}
