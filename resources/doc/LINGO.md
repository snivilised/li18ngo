# lingo - i18n Code Generator For go-i18n Integration

<!-- MD029/ol-prefix: Ordered list item prefix -->
<!-- MarkDownLint-disable MD029 -->

`lingo` provides a helper framework and code generator that assists in defining the Go structs required to integrate with the [`go-i18n`](https://github.com/nicksnyder/go-i18n) package for internationalisation (i18n) support.  

This tool ensures consistent, type-safe, and automatically validated message structures that power localised strings throughout your application.

Each message definition, referred to as an **underlier**, is represented in the `Underliers` map. This map's keys are message identifiers, and each entry specifies metadata used by `lingo` to automatically generate Go source files compatible with `go-i18n`.

Example entry:

```go
"root-command-config-file-usage": {
  MessageID:   "root-command-config-file-usage",
  Seed:        "RootCmdConfigFileUsage",
  TypeName:    UnderlyingTypeDynamicCobra,
  Description: "root command config flag usage",
  Story: "RootCmdConfigFileUsage is the usage string for the" +
    " config file flag on the root command.",
  Other: "config file (default is $HOME/{{.ConfigFileName}}.yml)",
  Fields: []UnderlyingField{
    {
      Note:   "ConfigFileName",
      GoType: "string",
      Tale:   "is the base name of the config file without extension",
    },
  },
},
```

For a description of these types (`UnderlyingTemplData` and `UnderlyingField`) and how to define them, please refer to the source: [`underlying.go`](./locale/underlying.go)

---

## Getting Started

First, the `lingo` code generation tool needs to be installed as follows:

> go install github.com/snivilised/li18ngo/cmd/lingo@latest

Before running `lingo`, the author needs to define entries inside a specific map. The definition the author need to provide is as follows:

```go
using (
  lingo "github.com/snivilised/li18ngo/locale"
)

var _ = lingo.Underliers{
  // entries go here...
  "message-id-1": {
    MessageID: "message-id-1",
  },
  // ...
}
```

The keys to this map is a string message id value which is also populated with the same value in the `MessageID` member. This looks like duplication and it is, but doing it this way in a map helps the user by not allowing them to accidentally define multiple entries with the same key value.

You will notice that an alias has been used to refer to the `locale` package inside `li18ngo`. This is of course optional and it up to the author how they wish to do this, but when lingo is run, it is expecting to find the `Underliers` definitions via the Go Ast ("go/ast"). By default it is expecting to find this in a `locale` package found in the root of the repo, or `./src/locale`. Alternatively, the author can specify an alternative path using the `--locale` flag. Either way, there must be a populated `Underliers` map. If the author chooses to go with the default use of a `locale` package, then not providing an alternative alias for li18no's version of `locale` will not be possible from inside the native `locale` package due to a package name clash.

The `Underliers` map may be defined in any file of the user's choice as long as it is in the `locale` package or the override specified by the `--locale` flag.

Before trying to generate the code for all messages, it is advisable to use the dry run mode by using the `--dry-run` flag. Doing so invokes verification and shows the user if all entries are valid (see the section [Validation Rules](#validation-rules)
below). As an aid, the user can use the `--verbose` flag to see extra output during validation phase.

When all messages have been define, generate the i18n code, merely using (from the repo root):

> $ lingo

## Generated Files

When `lingo` is executed, the code generator produces specific Go files for different message categories:

- 🧊 **Cobra messages:**  
  `messages-cobra-auto.go` - used for CLI command and flag text.
  
- ❌ **Error messages:**  
  `messages-errors-auto.go` - defines typed, localised error messages compatible with `errors.Is()` and `unwrap` operations.

- 📨 **General messages:**  
  `messages-general-auto.go` - general, non-error user-facing messages with or without dynamic content.

These files are **automatically generated** and should **never** be hand-edited. To modify a message, adjust the `Underliers` map and re-run.

---

## UnderlyingType Reference

Every underlier declares a `TypeName` that controls how code is generated and how the message is interpreted.  
Below is a summary table of all supported types:

|Type Name|Description|
|---|---|
|`UnderlyingTypeCobraStatic`|Short description for a Cobra command or flag. No dynamic fields.|
|`UnderlyingTypeCobraDynamic`|Long, parameterised description for a Cobra command or flag. Has dynamic fields.|
|`UnderlyingTypeGeneralStatic`|Static, user-facing message without variable content.|
|`UnderlyingTypeGeneralDynamic`|Dynamic, user-facing message with variable tokens.|
|`UnderlyingTypeErrorStatic`|Error with fixed content. Generates a static sentinel and error constructor.|
|`UnderlyingTypeErrorCore`|Core sentinel error meant to be wrapped. Used directly with `errors.Is()`.|
|`UnderlyingTypeErrorStaticWrapper`|Static error that wraps another error. Automatically includes a `Wrapped` field.|
|`UnderlyingTypeStaticErrorWrapperMsg`|Static error that wraps another error and includes wrapped error.|
|`UnderlyingTypeErrorDynamic`|Dynamic error with variable fields but no wrapping.|
|`UnderlyingTypeErrorDynamicWrapper`|Dynamic error that wraps another error and may interpolate `{{.Wrapped}}` in message text.|

---

## Validation Rules

Before code generation, all entries in the `Underliers` map are validated.  
Generation proceeds **only if no errors** are detected.  

The validation ensures structural and semantic correctness using the following rules:

- Fields must be empty for *static* types.  
- Fields must be non-empty for *dynamic* types.  
- Every `{{.Token}}` in `Other` must correspond to a field in `Fields`.  
- Every field must be used once in `Other`.  
- Only one field may have `GoType: "error"`.  
- Error fields must be named `Wrapped` when present.  
- Non-wrapper types may not define or use `Wrapped`.  
- Static wrapper errors must not define `Fields`.  
- `{{.Wrapped}}` tokens are only valid on wrapper types.  
- Duplicate `MessageID`s across the map are not allowed.

This ensures that `lingo` produces coherent, fully type-safe output for all translation templates.

---

## UnderlyingField and UnderlyingTemplData

Each entry in an `Underlier` may define an array of **`UnderlyingField`** structures, representing the parameters injected into the message template.  
An `UnderlyingField` typically defines:

- `Note` - Descriptive name or usage hint for the field.  
- `GoType` - The Go type of the field (`string`, `int`, `error`, etc.).  
- `Tale` - Additional context or documentation describing its role.

For dynamic messages, every `{{.Token}}` in the template (found in the `Other` field) must correspond to a `UnderlyingField` entry.

The code generator uses these fields to produce a strongly typed **`UnderlyingTemplData`** struct for each message, such as `RootCmdConfigFileUsageTemplData`.  
This generated struct contains the exact field definitions required to populate the message's dynamic content and is used to construct localised instances during runtime.

For every dynamic message, a corresponding `NewXxxTemplData` constructor is generated, allowing you to instantiate the template data easily.

---

## Type Descriptions

### UnderlyingTypeCobraStatic

Used for short, static descriptions of Cobra commands or flags.  
No dynamic fields are permitted and no constructor is generated.

### UnderlyingTypeCobraDynamic

Used for long, parameterised descriptions in Cobra commands or flags.  
Fields must be non-empty and correspond to template tokens.  
Generates a `NewXxxTemplData` constructor.

### UnderlyingTypeGeneralStatic

Represents a static, non-error user message with no dynamic content.  
Produces no constructor or data struct.

### UnderlyingTypeGeneralDynamic

Represents a dynamic, non-error user message with one or more variable fields.  
Each field maps to a template token.  
A constructor `NewXxxTemplData` is generated.

### UnderlyingTypeErrorStatic

Defines an error message with fixed text and no fields.  
Produces:

- `XxxErrorTemplData`
- `XxxError`
- `ErrXxx` sentinel.

### UnderlyingTypeErrorCore

Defines a static core sentinel error intended for wrapping.  
Produces:

- `XxxErrorTemplData`
- `XxxError`
- Exported `ErrXxx` sentinel usable with `errors.Is()`.

### UnderlyingTypeErrorStaticWrapper

Defines a static error that wraps another error implicitly.  
No `Fields` are declared.  
Generates:

- `XxxErrorTemplData`
- `XxxError{Wrapped error}`
- `NewXxxError(wrapped error)`
- `Error()` and `Unwrap()` methods.

### UnderlyingTypeStaticErrorWrapperMsg

Is a static error that wraps
another error and includes the wrapped error's message in the
translated output via `{{.Wrapped}}`. Use this instead of
`UnderlyingTypeStaticErrorWrapper` when `Other` contains `{{.Wrapped}}`.

- `XxxErrorTemplData`
- `XxxError{Wrapped error}`
- `NewXxxError(wrapped error)`
- `Error()` and `Unwrap()` methods.

### UnderlyingTypeErrorDynamic

Defines dynamic errors with variable content but no wrapping.  
Fields must be present, but no `Wrapped` field.  
Generates:

- `XxxTemplData`
- `XxxError`
- `NewXxxError(fields...)`

### UnderlyingTypeErrorDynamicWrapper

Defines dynamic errors that also wrap another error.  
Fields must include exactly one `error` type named `Wrapped`.  
The constructor takes the wrapped error first, followed by other parameters.  
Generates:

- `XxxTemplData` (stringified `Wrapped`)
- `XxxError` (containing actual error)
- `NewXxxError(wrapped error, fields...)`

---

## Example: Full Generation Flow

Here's how a typical workflow looks end-to-end:

1. **Define your message in the Underliers map**

   ```go
   var _ = map[string]Underlying{
    "root-command-config-file-usage": {
      MessageID:   "root-command-config-file-usage",
      Seed:        "RootCmdConfigFileUsage",
      TypeName:    UnderlyingTypeCobraDynamic,
      Description: "Usage string for the root command config flag.",
      Story:       "RootCmdConfigFileUsage provides the help text for the config file flag.",
      Other:       "config file (default is $HOME/{{.ConfigFileName}}.yml)",
      Fields: []UnderlyingField{
          {
            Note: "ConfigFileName",
            GoType: "string",
            Tale: "Base name of the config file without extension.",
          },
      },
    },
   }
   ```

2. **Run the generator**

> $ lingo

3. **Generated output (excerpt)**

   ```go
   // Code generated by lingo; DO NOT EDIT.

   type RootCmdConfigFileUsageTemplData struct {
       ConfigFileName string
   }

   func NewRootCmdConfigFileUsageTemplData(configFileName string) RootCmdConfigFileUsageTemplData {
       return RootCmdConfigFileUsageTemplData{ConfigFileName: configFileName}
   }

   var RootCmdConfigFileUsage = i18n.Message{
       ID:    "root-command-config-file-usage",
       Other: "config file (default is $HOME/{{.ConfigFileName}}.yml)",
   }
   ```

4. **Usage at runtime**

   ```go
   data := NewRootCmdConfigFileUsageTemplData("my-config")
   localised := localizer.MustLocalizeMessage(&RootCmdConfigFileUsage, data)
   fmt.Println(localised)
   ```

---
