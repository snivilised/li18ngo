# 📒 li18ngo: ___Defining Translate-able Content___

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD014/commands-show-output: Dollar signs used before commands without showing output mark down lint -->
<!-- MarkDownLint-disable MD014 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<!-- MD028/no-blanks-blockquote: Blank line inside blockquote -->
<!-- MarkDownLint-disable MD028 -->

Making an application work across different locales, is definitely no trivial task and after having dipped one's toe into using a library like __go-i18n__, it quickly becomes apparent that this can become a tiresome task. In particular, defining custom i18n Messages for translate-able content can easily become inconsistent if due care is not taken.

This document describes different scenarios and some standards that can be employed to help keep definitions consistent. Of course these are not mandatory, third parties are free to define i18n.Messages and data templates as they see fit, but this is being documented for the purposes of keeping __snivilised__ projects inline.

The categories covered are:

* single word (simple)
* phrase (multiple words)
* key/value field
* static error
* core error (wrapped error)
* dynamic error (error wrapper)

For each description of a message, there will be a definition of a data template of the form ___xxxTemplData___, eg 'InternationaliseTemplData', and this struct will embed an a fictional type 'heliosTemplData', which in reality would be replaced by a project specific type. So for example, the 'traverse' project will use type ___traverseTemplData___. The purpose of this embedded data template type is simply to define the ___SourceID___ required to work with li18ngo (if using either of the 2 __snivilised__, template projects, __arcadia__ or __astrolib__, this type will automatically be defined). For each message definition, there will be an example and then a generalised form, which the reader can copy and paste to create custom definitions. They can also be used to defined code snippets as an alternative way of speeding up implementation of new messages (see [go i18n snippets](./.vscode/snippets/go.my-i18n.json); they can be activated by copying them over to your vscode global snippets file. Adjust them according to your own needs).

## Non Error Content

### 📬 Single Word

* 📨 Message ID: "xxx.word"

```go
  type InternationaliseTemplData struct {
    heliosTemplData
  }

  // Message
  func (td InternationaliseTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "Internationalise.word",
      Description: "Internationalise",
      Other:       "internationalise",
    }
  }
```

* ⭕ Generalised form:

```go
  type FooTemplData struct {
    heliosTemplData
  }

  // Message
  func (td FooTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "---.word",
      Description: "---",
      Other:       "",
    }
  }
```

* 🍉 snippet prefix: "w18", ___Single i18n Word___

| variable | description |
|----------|-------------|
| $1 | word |
| $2 | repo name |

### 📬 Phrase

* 📨 Message ID: "xxx.phrase"

Essentially, this is no different to a `word`, except that the ID contains a different suffix.

```go
  type IHaveACunningPlanTemplData struct {
    heliosTemplData
  }

  // Message
  func (td IHaveACunningPlanTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "i-have-a-cunning-plan.phrase",
      Description: "i have a cunning plan (my lord!)",
      Other:       "i have a cunning plan",
    }
  }
```

* ⭕ Generalised form:

```go
  type FooTemplData struct {
    heliosTemplData
  }

  // Message
  func (td FoolData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "---.phrase",
      Description: "---",
      Other:       "---",
    }
  }
```

* 🍑 snippet prefix: "w18", ___i18n Phrase___

| variable | description |
|----------|-------------|
| $1 | phrase |
| $2 | repo name |

📍 __Note__ : spaces and dashes will need to be inserted into the `ID`, `Description` and `Other` fields.

### 📬 Key Value Field

* 📨 Message ID: "xxx.field"

Just to clarify, the key/value pair being addressed here is the case where the field is constant content and the value is dynamically created at runtime and thus not subject to translation, eg:

> name: marina

'name' is the translate-able field name and 'marina' is the un-translate-able value. Given that this field is going to be called 'Greeting', we could define as follows:

* 📕 TemplData:

```go
  type GreetingTemplData struct {
    heliosTemplData
    Name string
  }

  // Message
  func (td GreetingFieldTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "greeting.field",
      Description: "greeting displayed to user on application start up",
      Other:       "Name: {{.Name}}",
    }
  }
```

* ⭕ Generalised form:

```go
  type FooTemplData struct {
    heliosTemplData
    Key string
  }

  // Message
  func (td FooFieldTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "---.field",
      Description: "---",
      Other:       "Name: {{.Key}}",
    }
  }
```

* 🥝 snippet prefix: "f18", ___i18n Field___

| variable | description |
|----------|-------------|
| $1 | name of the key/value field |
| $2 | repo name |
| $3 | key |

`key` is assumed to be of type string, adjust this type as appropriate.

📍 __Note__ : spaces and dashes will need to be inserted into the `ID`, `Description` and `Other` fields.

📍 __Note__ : $3/Key is a pseudo field because VSCode will not prompt for content for this parameter due to the way that snippets are implemented. After the snippet has been inserted, the key will remain as 'Key string' and will require updating as appropriate, including inside the template reference within `Other`.

## Error Content

Errors require more definitions, because we need to implement standard Go error handling features, including error wrapping to support ___errors.Is/As___. The general pattern for errors builds upon the template data definitions we have already seen in previous sections. We also need:

* definition of the error type, which uses a template data instance to define its content
* a NewXXX error constructor function, to hide away the complexity of composing the error type with the template data
* a tester function in the form __IsXXXError__ which can be used to test whether a particular error is of this error type
* optionally, a core error to support error unwrapping

### LocalisableError

___li18ngo.LocalisableError___ is the struct that can be embedded into a custom error. Doing so, bestows upon the error, the built in invocation of the ___li18ngo.Text___ function when its ___Error___ method is invoked.

---

### 📬 Static

* 📨 Message ID: "xxx.static"

This is a static error message, that has no dynamic content. The data template will not contain any extra members and the ___Other___ field will be a straight forward string containing no break out references to template fields (eg {{.Foo}}). This is similar to a `phrase`, except we have to add in the error plumbing.

The following describes an 'out of memory' error:

* 📕 TemplData:

```go
  type OutOfMemoryTemplData struct {
    heliosTemplData
  }

  func (td OutOfMemoryTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "out-of-memory.error",
      Description: "System has unable to allocate new memory",
      Other:       "out of memory",
    }
  }
```

* 💥 Error type

```go
type OutOfMemoryError struct {
  li18ngo.LocalisableError
}
```

* ⭕ Generalised form:

```go
  type FooTemplData struct {
    heliosTemplData
  }

  func (td FooTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "foo.error",
      Description: "---",
      Other:       "---",
    }
  }

  type FooError struct {
    li18ngo.LocalisableError
  }
```

* 🍓 snippet prefix: "e18s", ___static i18n error___

| variable | description |
|----------|-------------|
| $1 | name of the error |
| $2 | repo name |

* 🎁 Create:

There's no mystery about this scenario, no constructor required so just create the error:

```go
  err := OutOfMemoryError{}
```

* 🎯 Identify:

Create an instance, this can be a package level global and then use ___errors.Is___ directly

```go
var ErrOutOfMemoryError = OutOfMemoryError{}

...
  if err := operation(); != nil && errors.Is(err, ErrOutOfMemoryError) {
    //...
  }
```

To maintain consistency with other error types, the global error instance could remain un-exported and a new function ___IsOutOfMemoryError___ defined instead:

```go
  func IsOutOfMemoryError(err error) bool {
    return errors.Is(err, errOutOfMemoryError)
  }
```

---

### 📬 Dynamic

* 📨 Message ID: "xxx.dynamic"

This message optionally contains a static part with un-translate-able dynamic content so this is similar to 'Key Value Field' previously described. The following example illustrates just a single text item, but of course messages can be more complex than this, with multiple fields of various types.

> path not found: /system/app/configs/foo

'path not found' is the translate-able part and '/system/app/configs/foo' is the un-translate-able value. Now to enable identification of this error to work consistently using ___errors.Is___, we need to split the static part from the dynamic part. The static part is handled by defining a core error (see the following section). We define a wrapper around the core which constitutes the wrapper. The most concise way we can do this is as follows:

* 📕 TemplData:

```go
  type PathNotFoundTemplData struct {
    heliosTemplData
    Path string
  }

  func (td PathNotFoundTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "path-not-found.dynamic-error",
      Description: "path not found dynamic error",
      Other:       "path: {{.Path}}",
    }
  }

  type PathNotFoundError struct {
    li18ngo.LocalisableError
    Wrapped error
  }

  func (e PathNotFoundError) Error() string {
    return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
  }

  func (e PathNotFoundError) Unwrap() error {
    return e.Wrapped
  }

  func NewPathNotFoundError(path string) error {
    return &PathNotFoundError{
      LocalisableError: li18ngo.LocalisableError{
        Data: PathNotFoundTemplData{
          Path: path,
        },
      },
      Wrapped: errCorePathNotFound,
    }
  }
```

So the effect of this is that the core error is the fundamental error required for identification purposes and the wrapper with the dynamic part adds context to this fundamental.

📍 __Note__ also that we now need to define an alternative implementation of the ___Error___ method, which combines the static part with the dynamic. The ___Unwrap___ method defined is invoked for us whenever the client invokes ___errors.Is___, which will be done via the ___IsPathNotFoundError___ defined for the core.

* ⭕ Generalised form:

```go
  type FooTemplData struct {
    heliosTemplData
    Field string
  }

  func (td FooTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "---.dynamic-error",
      Description: "--- dynamic error",
      Other:       "field: {{.Field}}",
    }
  }

  type FooError struct {
    li18ngo.LocalisableError
    Wrapped error
  }

  func (e FooError) Error() string {
    return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
  }

  func (e FooError) Unwrap() error {
    return e.Wrapped
  }

  func NewFooError(field string) error {
    return &FooError{
      LocalisableError: li18ngo.LocalisableError{
        Data: FooTemplData{
          Field: field,
        },
      },
      Wrapped: errCoreFoo,
    }
  }
```

* 🍒 snippet prefix: "e18d", ___dynamic i18n error___

| variable | description |
|----------|-------------|
| $1 | name of the error |
| $2 | repo name |
| $3 | context field |

📍 __Note__ The snippet will generate code to access the core error. This core error reference will probably be wrong, unless the core error has been pre-defined correctly first. It is recommended to use the core error snippet first then this one next.

* 🎁 Create:

```go
  path := "/system/app/configs/foo"
  NewPathNotFoundError(pattern)
```

* 🎯 Identify:

```go
  func IsPathNotFoundError(err error) bool {
    return errors.Is(err, errCorePathNotFound)
  }

  ...

  if err := operation(); err != nil && IsPathNotFoundError(err) {
    // some remedial action
  }
```

---

### 📬 Core

* 📨 Message ID: "xxx.core-error"

The core error is not meant to be used in isolation, it's purpose is simply to be wrapped by a dynamic error.

```go
  type CorePathNotFoundErrorTemplData struct {
    heliosTemplData
  }

  func IsPathNotFoundError(err error) bool {
    return errors.Is(err, errCorePathNotFound)
  }

  func (td CorePathNotFoundErrorTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "path-not-found.core-error",
      Description: "path not found core error",
      Other:       "path not found",
    }
  }

  type CorePathNotFoundError struct {
    li18ngo.LocalisableError
  }

  var errCorePathNotFound = CorePathNotFoundError{
    LocalisableError: li18ngo.LocalisableError{
      Data: CorePathNotFoundErrorTemplData{},
    },
  }
```

* ⭕ Generalised form:

```go
  type CoreFooErrorTemplData struct {
    heliosTemplData
  }

  func IsFooError(err error) bool {
    return errors.Is(err, errCoreFoo)
  }

  func (td CoreFooErrorTemplData) Message() *i18n.Message {
    return &i18n.Message{
      ID:          "---.core-error",
      Description: "--- core error",
      Other:       "---",
    }
  }

  type CoreFooError struct {
    li18ngo.LocalisableError
  }

  var errCoreFoo = CoreFooError{
    LocalisableError: li18ngo.LocalisableError{
      Data: CoreFooErrorTemplData{},
    },
  }
```

* 🥥 snippet prefix: "e18c", ___core i18n error___

| variable | description |
|----------|-------------|
| $1 | name of the core error |
| $2 | repo name |
