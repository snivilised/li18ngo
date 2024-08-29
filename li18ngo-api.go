package li18ngo

import (
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/nfs"
)

var (
	// 🌐 translate

	// DefaultLanguage represents the default language of this module
	DefaultLanguage = translate.DefaultLanguage

	// ErrSafePanicWarning is the error raised as a panic if the client
	// has accidentally not called Use before working with li18ngo.
	ErrSafePanicWarning = translate.ErrSafePanicWarning

	// Li18ngoSourceID the id that represents this module. If client want
	// to provides translations for languages that li18ngo does not, then
	// the localizer the create created for this purpose should use this
	// SourceID. So whenever the Text function is used on templates defined
	// inside this module, the translation process is directed to use the
	// correct i18n.Localizer (identified by the SourceID). The Source is
	// statically defined for all templates defined in li18ngo.
	Li18ngoSourceID = translate.Li18ngoSourceID

	// Text is the function to use to obtain a string created from
	// registered Localizers. The data parameter must be a go template
	// defining the input parameters and the translatable message content.
	Text = translate.Text

	// Use, must be called before any string data can be translated.
	// If requesting the default language, then only the language Tag
	// needs to be provided. If the requested language is not the default
	// and therefore requires translation from the translation file(s), then
	// the App and Path properties must be provided indicating
	// how the i18n bundle is created.
	// If only the Default language, then Use can even be called without
	// specifying the Tag and in this case the default language will be
	// used. The client MUST call Use before using any functionality in
	// this package.
	Use = translate.Use
)

type (
	// 🌐 nfs

	// ExistsInFS provides the facility to check the existence
	// of a path in the underlying file system.
	ExistsInFS = nfs.ExistsInFS

	// MkDirAllFS is a file system with a MkDirAll method.
	MkDirAllFS = nfs.MkDirAllFS

	// 🌐 translate

	// LoadFrom denotes where to load the translation file from
	LoadFrom = translate.LoadFrom

	// LocalisableError is an error that is translate-able (Localisable)
	LocalisableError = translate.LocalisableError

	// LanguageInfo information pertaining to setting language. Auto detection
	// is not supported. Any executable that supports i18n, should perform
	// auto detection and then invoke Use, with the detected language tag
	LanguageInfo = translate.LanguageInfo

	// LocalizerCreatorFn represents the signature of the function that can
	// optionally be provided to override how an i18n Localizer is created.
	LocalizerCreatorFn = translate.LocalizerCreatorFn

	// SupportedLanguages is a collection of the language Tags that a module
	// can define to express what languages it contains translations for.
	SupportedLanguages = translate.SupportedLanguages

	// TranslationSource
	// Name: core name of dependency's translation file. The actual file
	// is derived from this name in the form: <name>.active.<lang>.json;
	// eg li18ngo.active.en-GB.json.
	// Path: file system path to the translation file. If missing, then
	// it will default to the location of the executable file.
	TranslationSource = translate.TranslationSource

	// TranslationFiles maps a source id to a TranslationSource
	TranslationFiles = translate.TranslationFiles

	// UseOptions the options provided to the Use function
	UseOptions = translate.UseOptions

	// UseOptionFn functional options function required by Use.
	UseOptionFn = translate.UseOptionFn
)
