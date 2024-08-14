package translate

import "io/fs"

type multiplexor struct {
}

func (mx *multiplexor) invoke(localizer *Localizer, data Localisable) string {
	return localizer.MustLocalize(&LocalizeConfig{
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

	return mc.invoke(localizer, data), err
}

func (mc *multiContainer) add(info *LocalizerInfo) {
	if _, found := mc.localizers[info.sourceID]; found {
		return
	}

	mc.localizers[info.sourceID] = info.Localizer
}

func (mc *multiContainer) find(id string) (*Localizer, error) {
	if loc, found := mc.localizers[id]; found {
		return loc, nil
	}

	return nil, NewCouldNotFindLocalizerNativeError(id)
}
