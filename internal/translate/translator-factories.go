package translate

import (
	"io/fs"

	"github.com/snivilised/li18ngo/internal/lo"
	"github.com/snivilised/li18ngo/internal/nfs"
)

type translatorFactory struct {
	Create LocalizerCreatorFn
	legacy Translator
}

func (f *translatorFactory) setup(lang *LanguageInfo) {
	verifyLanguage(lang)

	if f.Create == nil {
		f.Create = createLocalizer
	}
}

// multiTranslatorFactory creates a translator instance from the provided
// Localizers.
//
// Note, in the case where a source client wants to provide a localizer
// for a language that one of ite dependencies does not support, then
// the translator should create the localizer based on its own default
// language, but we load the client provided translation file at the same
// name as the dependency would have created it for, then this file will
// be loaded as per usual.
type multiTranslatorFactory struct {
	translatorFactory
}

func (f *multiTranslatorFactory) New(lang *LanguageInfo) (Translator, error) {
	f.setup(lang)

	queryFS := lo.TernaryF(lang.FS != nil,
		func() fs.StatFS {
			return lang.FS
		},
		func() fs.StatFS {
			native := nfs.NewReadDirFS(lang.From.Path)
			return nfs.StatFSFromFS(native)
		},
	)

	dirFS := lo.TernaryF(lang.FS != nil,
		func() nfs.MkDirAllFS {
			return nfs.DirFSFromFS(lang.FS)
		},
		func() nfs.MkDirAllFS {
			return nfs.DirFSFromFS(queryFS)
		},
	)

	multi := &multiContainer{
		localizers: make(localizerContainer),
		queryFS:    queryFS,
	}

	for id := range lang.From.Sources {
		localizer, err := f.Create(lang, id, dirFS)

		if err != nil {
			return nil, err
		}

		multi.add(&LocalizerInfo{
			SourceID:  id,
			Localizer: localizer,
		})
	}

	return &i18nTranslator{
		mx:           multi,
		languageInfo: lang,
	}, nil
}
