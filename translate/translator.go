package translate

import (
	"github.com/snivilised/li18ngo/internal/lo"
	"github.com/snivilised/li18ngo/utils"
	"golang.org/x/text/language"
)

type Translator interface {
	Localise(data Localisable) string
	LanguageInfoRef() utils.RoProp[*LanguageInfo]
	negotiate(other Translator) (Translator, error)
	add(info *LocalizerInfo, source *TranslationSource)
}

var DefaultLanguage = utils.NewRoProp(language.BritishEnglish)
var tx Translator
var TxRef utils.RoProp[Translator] = utils.NewRoProp(tx)

type localizerContainer map[string]*Localizer

// Use, must be called by the client before any string data
// can be translated. If the client requests the default
// language, then only the language Tag needs to be provided.
// If the requested language is not the default and therefore
// requires translation from the translation file(s), then
// the client must provide the App and Path properties indicating
// how the l18n bundle is created.
// If the client just wishes to use the Default language, then Use
// can even be called without specifying the Tag and in this case
// the default language will be used. The client MUST call Use
// before using any functionality in this package.
func Use(options ...UseOptionFn) error {
	var err error

	o := &UseOptions{}

	o.DefaultIsAcceptable = true
	o.Tag = DefaultLanguage.Get()

	for _, fo := range options {
		fo(o)
	}

	lang := NewLanguageInfo(o)

	if !containsLanguage(lang.Supported, o.Tag) {
		if o.DefaultIsAcceptable {
			o.Tag = DefaultLanguage.Get()
			lang.Tag = o.Tag
		} else {
			err = NewFailedToCreateTranslatorNativeError(o.Tag)
		}
	}

	if err == nil {
		err = applyLanguage(lang)
	}

	return err
}

func verifyLanguage(lang *LanguageInfo) {
	if lang.From.Sources == nil {
		lang.From.Sources = make(TranslationFiles)
	}

	// By adding in the source for li18ngo, we relieve the client from having
	// to do this. After-all, it should be taken as read that if the client is
	// using li18ngo then the translations for li18ngo should be loaded,
	// otherwise li18ngo will not be able to convey these translations to the
	// client. The client app has to make sure that when their app is deployed,
	// the translations file(s) for li18ngo are named as 'li18ngo', as you
	// can see below, that that is the name assigned to the app name of the
	// source. There is little value in making this customisable as this would
	// just lead to confusion. If the client really wants to control the name
	// of the translation file for li18ngo, they can provide an override
	// 'Create' function on UseOptions.
	//
	if _, found := lang.From.Sources[Li18ngoSourceID]; !found {
		lang.From.Sources[Li18ngoSourceID] = TranslationSource{Name: "li18ngo"}
	}
}

func applyLanguage(lang *LanguageInfo) error {
	verifyLanguage(lang)
	factory := &multiTranslatorFactory{
		AbstractTranslatorFactory: AbstractTranslatorFactory{
			Create: lang.Create,
			legacy: tx,
		},
	}

	newTranslator, err := factory.New(lang)
	if err != nil {
		return err
	}

	tx, err = negotiateTranslators(tx, newTranslator)

	TxRef = utils.NewRoProp(tx)

	return err
}

func negotiateTranslators(legacyTX, incomingTX Translator) (Translator, error) {
	var (
		err error
	)

	result := incomingTX

	if legacyTX != nil {
		result, err = legacyTX.negotiate(incomingTX)
	}

	return result, err
}

// ResetTx, do not use, required for unit testing only and is
// not considered part of the public api and may be removed without
// corresponding version number change.
func ResetTx() {
	// having to do this smells a bit, but required so unit tests can
	// remain isolated (this is why package globals are bad, but sometimes
	// unavoidable). This is all because we want to be able to call the Text
	// function easily. If we defined the Text function on an object, then that
	// would require passing that state around in many places, making the code
	// much more brittle and cumbersome to maintain.
	//
	tx = nil
	TxRef = utils.NewRoProp(tx)
}

// NewLanguageInfo gets a new instance of Language info from the use options
// specified. This is specific to li18ngo. Client applications should
// provide their own version that reflects their own defaults.
func NewLanguageInfo(o *UseOptions) *LanguageInfo {
	return &LanguageInfo{
		UseOptions: *o,
		Default:    DefaultLanguage.Get(),
		Supported: SupportedLanguages{
			DefaultLanguage.Get(),
			language.AmericanEnglish,
		},
	}
}

// Text is the function to use to obtain a string created from
// registered Localizers. The data parameter must be a go template
// defining the input parameters and the translatable message content.
func Text(data Localisable) string {
	return tx.Localise(data)
}

// i18nTranslator provides the translation implementation used by the
// Text function
type i18nTranslator struct {
	mx              *multiContainer
	languageInfo    *LanguageInfo
	languageInfoRef utils.RoProp[*LanguageInfo]
}

func (t *i18nTranslator) LanguageInfoRef() utils.RoProp[*LanguageInfo] {
	return t.languageInfoRef
}

func (t *i18nTranslator) Localise(data Localisable) string {
	// TODO: the error is ignored for now, but we should enable
	// the client to register an error hander. Localise should
	// be easy to use without worrying about the need to handle
	// an error that should not really ever happen. Or we revert
	// to default. There is no reason to ever return an error here.
	//
	tx, _ := t.mx.localise(data)
	return tx
}

func containsLanguage(languages SupportedLanguages, tag language.Tag) bool {
	return lo.ContainsBy(languages, func(t language.Tag) bool {
		return t == tag
	})
}

func (t *i18nTranslator) negotiate(incomingTX Translator) (Translator, error) {
	incomingLang := incomingTX.LanguageInfoRef().Get()
	legacyLang := t.LanguageInfoRef().Get()
	incTX, ok := incomingTX.(*i18nTranslator)

	if !ok {
		return nil, ErrInvalidTranslator
	}

	legacySources := legacyLang.From.Sources
	incomingSources := incomingLang.From.Sources

	for sourceID, source := range incomingSources {
		if _, found := legacySources[sourceID]; !found {
			s := source // copy required to avoid "implicit memory aliasing in for loop"
			t.add(&LocalizerInfo{
				Localizer: incTX.mx.localizers[sourceID],
				sourceID:  sourceID,
			}, &s)
		}
	}

	return t, nil
}

func (t *i18nTranslator) add(info *LocalizerInfo, source *TranslationSource) {
	t.mx.add(info)
	t.languageInfo.From.AddSource(info.sourceID, source)
}
