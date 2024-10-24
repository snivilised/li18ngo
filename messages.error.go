package li18ngo

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
)

// see https://dave.cheney.net/2014/12/24/inspecting-errors
//
// The requirement for localisation runs counter to that
// explained in this article, in particular the definition
// of typed errors increasing the api surface of a package
// and therefore makes the api more brittle. This issue is
// due-ly noted, but if translations are important, then we
// have to live with this problem unless another approach
// is available. Its not really recommended to provide foreign
// translations for external packages as this creates an
// undesirable coupling, but the option is there just in case.
// To ameliorate api surface area issue, limit error definitions
// to those errors that are intended to be displayed to
// the end user. Internal errors that can be handled, should not
// have translations templates defined for them as the user
// won't see them.
//
// As is presented in the article, clients are better off
// asserting errors for behaviour, not type, but this aspect
// should not be at cross purposes with the requirement for
// localisation.
//
//  In summary then, for ...
//
// * package authors: provide predicate interface definitions
// for errors that can be handled, eg "Timeout() bool". Also,
// use errors.Wrap to add context to another error.
// * package users: don't check an error's type, query for the
// declared interface, and invoke the provided predicates
// to identify an actual error.
//
// An alternative to providing foreign translations is just
// to handle the 3rd party error and Wrapping it up with a
// local error in the desired language. Sure, the inner error
// will be defined in the library's default language, but that
// can be wrapped (errors.Wrap), providing content in the
// required but library un-supported language.
//
// There does NOT need to be a translation file for the default language
// as the default language is what's implemented in code, here in
// message files (messages.error.nav.go). Having said that, we
// still need to create a file for the default language as that file
// is used to create translations. This default file will not be
// be part of the installation set.
// ===> checked in as locale/default/active.en-GB.json
//
// 1) This file is automatically processed to create the translation
// files, currently only 'active.en-US.json' by running:
// $ goi18n extract -format json -sourceLanguage "en-GB" -out ./out
// ---> creates locale/out/active.en-GB.json (locale/default/active.en-GB.json)
// ===> implemented as task: extract
//
// 2) ... Create an empty message file for the language that you want
// to add (e.g. translate.en-US.json).
// ---> when performing updates, you don't need to create the empty file, use the existing one
// ---> check-in the translation file
// ===> this has been implemented in the extract task
//
// 3) goi18n merge -format json active.en.json translate.en-US.json -outdir <dir>
// (goi18n merge -format <json|toml> <default-language-file> <existing-active-file>)
//
// existing-active-file: when starting out, this file is blank, but must exist first.
// When updating existing translations, this file will be the one that's already
// checked-in and the result of manual translation (ie we re-named the translation file
// to be active file)
//
// current dir: ./extendio/i18n/
// $ goi18n merge -format json -sourceLanguage "en-GB" -outdir ./out ./out/active.en-GB.json ./out/l10n/translate.en-US.json
//
// ---> creates the translate.en-US.json in the current directory, this is the real one
// with the content including the hashes, ready to be translated. It also
// creates an empty active version (active.en-US.json)
//
// ---> so the go merge command needs the translate file to pre-exist
//
// 4) translate the translate.en-US.json and copy the contents to the active
// file (active.en-US.json)
//
// 5) the translated file should be renamed to 'active' version
// ---> so 'active' files denotes the file that is used in production (loaded into bundle)
// ---> check-in the active file

// ====================================================================

// ❌ PathNotFound indicates path does not exist

// PathNotFoundTemplData is
type PathNotFoundTemplData struct {
	locale.Li18ngoTemplData
	Name string
	Path string
}

func (td PathNotFoundTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "path-not-found.error",
		Description: "Directory or file path does not exist",
		Other:       "{{.Name}} path not found ({{.Path}})",
	}
}

// PathNotFoundErrorBehaviourQuery used to query if an error is:
// "File system foo is ..."
type PathNotFoundErrorBehaviourQuery interface {
	IsPathNotFound() bool
}

type PathNotFoundError struct {
	translate.LocalisableError
}

// PathNotFound enables the client to check if error is PathNotFoundError
// via QueryPathNotFoundError
func (e PathNotFoundError) IsPathNotFound() bool {
	return true
}

// NewPathNotFoundError creates a PathNotFoundError
func NewPathNotFoundError(name, path string) PathNotFoundError {
	return PathNotFoundError{
		LocalisableError: translate.LocalisableError{
			Data: PathNotFoundTemplData{
				Name: name,
				Path: path,
			},
		},
	}
}

// ❌ Not A Directory

// NotADirectoryTemplData path is not a directory
type NotADirectoryTemplData struct {
	locale.Li18ngoTemplData
	Path string
}

func (td NotADirectoryTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "not-a-directory.error",
		Description: "File system path is not a directory",
		Other:       "file system path '{{.Path}}', is not a directory",
	}
}

// NotADirectoryErrorBehaviourQuery used to query if an error is:
// "File system path is not a directory"
type NotADirectoryErrorBehaviourQuery interface {
	NotADirectory() bool
}

type NotADirectoryError struct {
	translate.LocalisableError
}

// NotADirectory enables the client to check if error is NotADirectoryError
// via QueryNotADirectoryError
func (e NotADirectoryError) NotADirectory() bool {
	return true
}

// NewNotADirectoryError creates a NotADirectoryError
func NewNotADirectoryError(path string) NotADirectoryError {
	return NotADirectoryError{
		LocalisableError: translate.LocalisableError{
			Data: NotADirectoryTemplData{
				Path: path,
			},
		},
	}
}

// ❌ Third Party Error

// ThirdPartyErrorTemplData third party un-translated error
type ThirdPartyErrorTemplData struct {
	locale.Li18ngoTemplData

	Error error
}

func (td ThirdPartyErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "third-party.error",
		Description: "These errors are generated by dependencies that don't support localisation",
		Other:       "third party error: '{{.Error}}'",
	}
}

// ThirdPartyError represents an error received by a dependency that does
// not support i18n.
type ThirdPartyError struct {
	translate.LocalisableError
}

// NewThirdPartyErr creates a ThirdPartyErr
func NewThirdPartyErr(err error) ThirdPartyError {
	return ThirdPartyError{
		LocalisableError: translate.LocalisableError{
			Data: ThirdPartyErrorTemplData{
				Error: err,
			},
		},
	}
}
