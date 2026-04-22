# 🌐 li18ngo: ___i18n helper___

[![A B](https://img.shields.io/badge/branching-commonflow-informational?style=flat)](https://commonflow.org)
[![A B](https://img.shields.io/badge/merge-rebase-informational?style=flat)](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
[![A B](https://img.shields.io/badge/branch%20history-linear-blue?style=flat)](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule)
[![Go Reference](https://pkg.go.dev/badge/github.com/snivilised/li18ngo.svg)](https://pkg.go.dev/github.com/snivilised/li18ngo)
[![Go report](https://goreportcard.com/badge/github.com/snivilised/li18ngo)](https://goreportcard.com/report/github.com/snivilised/li18ngo)
[![Coverage Status](https://coveralls.io/repos/github/snivilised/li18ngo/badge.svg?branch=master)](https://coveralls.io/github/snivilised/li18ngo?branch=master&kill_cache=1)
[![Li18ngo Continuous Integration](https://github.com/snivilised/li18ngo/actions/workflows/ci-workflow.yml/badge.svg)](https://github.com/snivilised/li18ngo/actions/workflows/ci-workflow.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![A B](https://img.shields.io/badge/commit-conventional-commits?style=flat)](https://www.conventionalcommits.org/)

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

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" alt="go dev" /></a>
</p>

## 🔰 Introduction

This project contains helpers to aid in the development of libraries and programs that require i18n support. It doesn't perform i18n itself, that is delegated to [go-i18n](https://github.com/nicksnyder/go-i18n). Rather it aims to provide functionality that makes using go-i18n easier. For example, implementing localised error messages can be a bit tedious and error prone so included within this module is a cli app, ___lingo___, that can generate all error related code.

## 📚 Usage

## 🎀 Features

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" alt="ginkgo" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" alt="gomega" /></a>
</p>

+ unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)
+ i18n with [go-i18n](https://github.com/nicksnyder/go-i18n)
+ linting configuration and pre-commit hooks, (see: [linting-golang](https://freshman.tech/linting-golang/)).

### 🌐 l10n Translations

This template has been setup to support localisation. The default language is `en-GB` with support for `en-US`. There is a translation file for `en-US` defined as __src/i18n/deploy/arcadia.active.en-US.json__. This is the initial translation for `en-US` that should be deployed with the app.

Make sure that the go-i18n package has been installed so that it can be invoked as cli, see [go-i18n](https://github.com/nicksnyder/go-i18n) for installation instructions.

To maintain localisation of the application, the user must take care to implement all steps to ensure translate-ability of all user facing messages. Whenever there is a need to add/change user facing messages including error messages, to maintain this state, the user must:

+ define template struct (__xxxTemplData__) in __src/i18n/messages.go__ and corresponding __Message()__ method. All messages are defined here in the same location, simplifying the message extraction process as all extractable strings occur at the same place. Please see [go-i18n](https://github.com/nicksnyder/go-i18n) for all translation/pluralisation options and other regional sensitive content.

For more detailed workflow instructions relating to i18n, please see [i18n README](./resources/doc/i18n-README.md). For details on how defining translate-able content can be achieved consistently using code generation, see [Lingo](./resources/doc/LINGO.md)

## 🚀 Coding Guidelines

### Using li18ngo in Your Application

`li18ngo` wraps `go-i18n` and adds a localisation lifecycle on top of it. The
library distinguishes between two phases: __bootstrap__ (registering languages
and activating a language) and __runtime__ (translating messages via `Text`).
You must complete the bootstrap phase before making any `Text` calls, or the
library will not behave correctly.

---

### Bootstrap Phase

Before your application renders any localised string, call `Use` to activate
the library and register your supported languages:

```go
err := li18ngo.Use(func(o *li18ngo.UseOptions) {
    o.Tag = language.BritishEnglish
    o.From = li18ngo.LoadFrom{
        Path: "path/to/translations",
        Sources: li18ngo.TranslationFiles{
            li18ngo.Li18ngoSourceID: li18ngo.TranslationSource{
                Name: "li18ngo",
            },
            <YourPackage>.SourceID: li18ngo.TranslationSource{
                Name: "<your-app>",
            },
        },
    }
})
```

`Use` must be called exactly once per process lifetime, typically at the very
start of your `main` function or, for CLI applications built on Cobra/Mamba,
inside the bootstrap function that runs before any command executes.

If you ship a library that builds on `li18ngo`, call `Register` instead of
`Use` - `Register` is the per-library entry point that tells `li18ngo` about
your translation sources so that a host application can load them:

```go
func init() {
    li18ngo.Register(func(o *li18ngo.RegisterOptions) {
        o.SourceID = SourceID
        o.DefaultFS = &translationsFS
    })
}
```

The host application then includes your `SourceID` in its `Use` call, as shown
above.

😮‍💨 It needs to be acknowledged that building i18n compliant applications and libraries
can be quite onerous and a real pain in the you know whats. For this reason, it is
not mandatory to participate in the i18n infrastructure, just because one of your
dependencies does. i18n is not mandatory, but invoking Use/Register is. If you do not want
to participate in i18n, then just invoke these functions without passing in a registration
function, ie you can just do this:

for applications:

```go
  li18ngo.Use()
```

for libraries:

```go
  li18ngo.Register()
```

When no registration function is passed into Use/Register, then `Text` will just
return the default text as created by the library/application author; ie the text
you see in the `Other` field of the `TemplData` struct's `Message` method.

---

### Translating Messages at Runtime

Once bootstrap is complete, translate a message by constructing its template
data struct and passing it to `li18ngo.Text`:

```go
msg := li18ngo.Text(MyMessageTemplData{})
```

For messages that carry dynamic content, populate the exported fields of the
template data struct before passing it:

```go
msg := li18ngo.Text(FileNotFoundTemplData{
    Name: path,
})
```

`Text` returns a plain `string`. It never returns an error - if something goes
wrong (missing translation, uninitialised library), it falls back gracefully,
so you do not need to check a return value beyond the string itself.

When you define messages with dynamic content, then a helper function is provided
to relieve the author from knowing the full way to compose the appropriate structs,
eg:

```go
  msg := li18ngo.Text(locale.NewRootCmdConfigFileUsageTemplData(
    viper.ConfigFileUsed()),
  )
```

---

### Defining Messages

All user-facing strings are represented as template data structs typically in the
`locale` package (but this can be overridden by using the `--locale` flag on `lingo`). There are distinct message patterns, summarised below. (For the full enumeration definition, see [UnderlyingType](./locale/enums/underlying-type-en.go))

| Type | Emoji | Enum(UnderlyingType) | When to use |
| --- | --- | --- | --- |
| Cobra static | 🧊 | StaticCobra | Plain informational string, no variable content |
| Cobra dynamic | 🧊 | DynamicCobra | Informational string with one or more variable tokens |
| General static | 📨 | StaticGeneral | Plain informational string, no variable content |
| General dynamic | 📨 | DynamicGeneral | Informational string with one or more variable tokens |
| Static error | ❌ | StaticError | Fixed error message, no variable content |
| Static sentinel error | ❌ | SentinelError | Fixed error message; produces an exported sentinel `var Err...`. Designed to be wrapped |
| Static error (wrapper) | ❌ | StaticErrorWrapper | Error with no variable context that wraps an underlying error, wrapped error does nopt appear in translated output |
| Static error (wrapper) | ❌ | StaticErrorWrapperMsg | Error with no variable context that wraps an underlying error, wrapped error appears in translated output |
| Dynamic error (no wrap) | ❌ | DynamicError | Error carrying variable context, does not wrap another error |
| Dynamic error (wrapper) | ❌ | DynamicErrorWrapper | Error with variable context that wraps an underlying error |

Message IDs follow a strict convention:

| Message kind | ID format |
| --- | --- |
| Non-error | `kebab-slug` |
| Static error | `kebab-slug.static-error` |
| Dynamic error | `kebab-slug.dynamic-error` |

Never omit the `.static-error` / `.dynamic-error` suffix from error message
IDs and never add those suffixes to non-error messages. This is just a suggested convention; the author is free to use whatever scheme they wish.

---

### Code Generation

The `lingo` code generation tool manages the `messages-general-auto.go`,
`messages-errors-auto.go` and `messages-cobra-auto.go` description files automatically. You
should not hand-edit those generated files. The source of truth is
the `Underliers` map - add or modify message
descriptors there and re-run `lingo` to regenerate the output files. Full
documentation for [Lingo](./resources/doc/LINGO.md) is covered separately.

---

### Translation Files

Translation files are JSON and follow the `go-i18n` format. The active
language is selected by the `Tag` passed to `Use`. If no translation file
exists for the requested language, `go-i18n` falls back to the default
language (typically `en-GB`, depending on how your `Use` call is
configured).

Place translation files in a directory that is either embedded via `go:embed`
or readable from disk. Pass the corresponding `fs.FS` or path to `Use` via
`LoadFrom`.

---

### Error Handling Conventions

Static errors expose a package-level sentinel variable (`Err<Name>`) that
callers can match using the standard `errors.Is` / `errors.As` functions - no
generated helper wrappers are provided or needed:

```go
if errors.Is(err, locale.ErrFilterMissingType) {
    // handle it
}
```

Dynamic errors that wrap an underlying error implement `Unwrap`, so
`errors.Is` and `errors.As` chains work correctly without any extra
boilerplate.

---

### Troubleshooting

| Symptom | Likely cause | Fix |
| --- | --- | --- |
| All messages return the fallback / English string despite setting a language | `Use` was not called before the first `Text` call | Move your `Use` call earlier in bootstrap, before any command or handler executes |
| `Text` returns an empty string or panics | The message ID in the template data struct does not match the key in the translation JSON file | Check that `lingo` has been re-run after any change to `Underliers`, and that the translation file has been updated |
| A library's messages are not translated | The library's `SourceID` is missing from the `Sources` map in `Use` | Add the library's `SourceID` and `TranslationSource` to the host application's `Use` call |
| `Register` has no effect | `Register` was called after `Use` | `Register` must be called before `Use` - put it in a library `init` function |
| Sentinel error does not match via `errors.Is` | Caller wrapped the sentinel in a new error without preserving the chain | Use `fmt.Errorf("...: %w", locale.ErrFoo)` to wrap, not `fmt.Errorf("...: %v", ...)` |
| Generated files contain stale or missing messages | `lingo` has not been run after editing `Underliers` | Run `lingo` and commit the regenerated files |
