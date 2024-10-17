package li18ngo

import (
	"github.com/snivilised/li18ngo/internal/translate"
)

type TemplData struct{}

func (td TemplData) SourceID() string {
	return translate.Li18ngoSourceID
}
