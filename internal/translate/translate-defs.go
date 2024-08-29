package translate

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"github.com/snivilised/li18ngo/nfs"
	"golang.org/x/text/language"
)

// 📚 package: translate contains internal li18ngo definitions that
// client does not need direct access to, unless explicitly exposed
// by li18ngo-api.go.

const (
	// Li18ngoSourceID the id that represents this module. If client want
	// to provides translations for languages that li18ngo does not, then
	// the localizer the create created for this purpose should use this
	// SourceID. So whenever the Text function is used on templates defined
	// inside this module, the translation process is directed to use the
	// correct i18n.Localizer (identified by the SourceID). The Source is
	// statically defined for all templates defined in li18ngo.
	Li18ngoSourceID = "github.com/snivilised/li18ngo"
)

var (
	ErrSafePanicWarning = errors.New("please ensure li18ngo.Use is invoked")
)

type (
	SupportedLanguages []language.Tag

	Localisable interface {
		Message() *i18n.Message
		SourceID() string
	}

	TranslationSource struct {
		// Name of dependency's translation file
		Name string
		Path string
	}

	// TranslationFiles maps a source id to a TranslationSource
	TranslationFiles map[string]TranslationSource

	// LoadFrom denotes where to load the translation file from
	LoadFrom struct {
		// Path denoting where to load language file from, defaults to exe location
		//
		Path string

		// Sources are the translation files that need to be loaded. They represent
		// the client app/library dependencies.
		//
		// The source id would typically be the name of a package that is the source
		// of string messages that are to be translated. Actually, we could use
		// the top level url of the package by convention, as that is unique.
		// So li18ngo would use "github.com/snivilised/li18ngo" but clients
		// are free to use whatever naming scheme they want to use for their own
		// dependencies.
		//
		Sources TranslationFiles
	}

	// LocalizerCreatorFn represents the signature of the function that can
	// optionally be provided to override how an i18n Localizer is created.
	LocalizerCreatorFn func(li *LanguageInfo, sourceID string,
		dirFS nfs.MkDirAllFS,
	) (*i18n.Localizer, error)

	// UseOptionFn functional options function required by Use.
	UseOptionFn func(*UseOptions)

	// UseOptions the options provided to the Use function
	UseOptions struct {
		// Tag sets the language to use
		//
		Tag language.Tag

		// From denotes where to load the translation file from
		//
		From LoadFrom

		// DefaultIsAcceptable controls whether an error is returned if the
		// request language is not available. By default DefaultIsAcceptable
		// is true so that the application continues in the default language
		// even if the requested language is not available.
		//
		DefaultIsAcceptable bool

		// Create allows the client to  override the default function to create
		// the i18n Localizer(s) (1 per language).
		//
		Create LocalizerCreatorFn

		// Custom set-able by the client for what ever purpose is required.
		//
		Custom any

		// FS is a file system from where translations are loaded from. This
		// does not have to performed explicitly asa it will be created using
		// the From field if not specified.
		FS fs.StatFS
	}

	// LanguageInfo information pertaining to setting language. Auto detection
	// is not supported. Any executable that supports i18n, should perform
	// auto detection and then invoke Use, with the detected language tag
	LanguageInfo struct {
		UseOptions

		// Default language reflects the base language. If all else fails, messages will
		// be in this language. It is fixed at BritishEnglish reflecting the language this
		// package is written in.
		//
		Default language.Tag

		// Supported indicates the list of languages for which translations are available.
		//
		Supported SupportedLanguages
	}

	LocalizerInfo struct {
		// Localizer by default created internally, but can be overridden by
		// the client if they provide a create function to the Translator Factory
		//
		Localizer *i18n.Localizer

		SourceID string
	}

	Translator interface {
		Localise(data Localisable) string
		LanguageInfo() *LanguageInfo
		negotiate(other Translator) (Translator, error)
		add(info *LocalizerInfo, source *TranslationSource)
	}

	localizerContainer map[string]*i18n.Localizer
)

// AddSource adds a translation source
func (lf *LoadFrom) AddSource(sourceID string, source *TranslationSource) {
	if _, found := lf.Sources[sourceID]; !found {
		lf.Sources[sourceID] = *source
	}
}

var (
	tx              Translator
	DefaultLanguage = language.BritishEnglish
)
