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

	// File is an optional output-file prefix. When empty, the message is written
	// to the default output file for its kind:
	//   - messages-cobra-auto.go
	//   - messages-general-auto.go
	//   - messages-errors-auto.go
	//
	// When File is set, it is used as a prefix instead of "messages", producing
	// a custom output file of the same kind, for example:
	//   - File: "system-automation" + cobra kind  -> system-automation-cobra-auto.go
	//   - File: "system-automation" + general kind -> system-automation-general-auto.go
	//   - File: "system-automation" + error kind   -> system-automation-errors-auto.go
	//
	// Only letters, digits, underscores and dashes are permitted in File. A
	// trailing dash is silently stripped to prevent double-dash sequences in the
	// resulting filename. Any other invalid character is a terminating error.
	File string
}

// Underliers is the map type read by i18n-gen at code-generation time.
// The map key must equal the MessageID field of the value.
type Underliers map[string]UnderlyingTemplData

// underliers is the single source of truth for all i18n messages in this
// package. Edit this map and run go generate to regenerate the auto files.
var _ = Underliers{
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

	// TODO: generating silly function name, eg: locale.NewThirdPartyErrorWrapperError
	// The source ID for this message is wrong. It is currently li18ngo, instead of
	// being graffico and that is because PavementGraffitiReportGrafficoTemplData is
	// defined with Li18ngoTemplData instead of GrafficoData, which has been erroneously
	// removed.
	"third-party.error-wrapper-msg": {
		MessageID:   "third-party.error-wrapper-msg",
		Seed:        "ThirdPartyWrapper",
		TypeName:    enums.UnderlyingTypeStaticErrorWrapperMsg,
		Description: "Wrapper for third-party errors",
		Story:       "ThirdPartyErrorWrapper is used to wrap errors from third-party libraries.",
		Other:       "Third party error occurred: '{{.Wrapped}}'",
		Fields: []UnderlyingField{
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "is the original error from the third-party library that is being wrapped",
			},
		},
	},
}
