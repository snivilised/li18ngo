package translate

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	nef "github.com/snivilised/nefilim"
)

type multiplexor struct {
}

func (mx *multiplexor) invoke(localizer *i18n.Localizer, data Localisable) (string, error) {
	return localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: data.Message(),
		TemplateData:   data,
	})
}

type multiContainer struct {
	multiplexor
	localizers localizerContainer
	queryFS    fs.StatFS
	fS         nef.MakeDirFS
	create     LocalizerCreatorFn
}

func (mc *multiContainer) localise(data Localisable) (string, error) {
	id := data.SourceID()
	localizer, err := mc.find(id)

	if err != nil {
		localizer, err = mc.mitigate(id)

		if err != nil {
			return "", err
		}

		mc.add(&LocalizerInfo{
			Localizer: localizer,
			SourceID:  id,
		})
	}

	return mc.invoke(localizer, data)
}

func (mc *multiContainer) add(info *LocalizerInfo) {
	if _, found := mc.localizers[info.SourceID]; found {
		return
	}

	mc.localizers[info.SourceID] = info.Localizer
}

func (mc *multiContainer) find(id string) (*i18n.Localizer, error) {
	if loc, found := mc.localizers[id]; found {
		return loc, nil
	}

	return nil, NewCouldNotFindLocalizerNativeError(id)
}

func (mc *multiContainer) mitigate(id string) (*i18n.Localizer, error) {
	return mc.create(&LanguageInfo{
		UseOptions: UseOptions{
			Tag:                 DefaultLanguage,
			DefaultIsAcceptable: true,
			Create:              mc.create,
			FS:                  mc.queryFS,
		},
		Default: DefaultLanguage,
		Supported: SupportedLanguages{
			DefaultLanguage,
		},
	}, id, mc.fS)
}
