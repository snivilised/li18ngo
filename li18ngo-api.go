package li18ngo

import (
	"github.com/snivilised/li18ngo/internal/translate"
)

var (
	// DefaultLanguage represents the default language of this module
	DefaultLanguage = translate.DefaultLanguage

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
	// LoadFrom denotes where to load the translation file from
	LoadFrom = translate.LoadFrom

	// TranslationSource
	TranslationSource = translate.TranslationSource

	// TranslationFiles maps a source id to a TranslationSource
	TranslationFiles = translate.TranslationFiles

	// UseOptions the options provided to the Use function
	UseOptions = translate.UseOptions

	// UseOptionFn functional options function required by Use.
	UseOptionFn = translate.UseOptionFn
)
