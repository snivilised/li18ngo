package translate

import (
	"io/fs"

	"github.com/nicksnyder/go-i18n/v2/i18n"
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
}

func (mc *multiContainer) localise(data Localisable) (string, error) {
	localizer, err := mc.find(data.SourceID())

	if err != nil {
		return "", err
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
