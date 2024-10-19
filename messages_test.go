package li18ngo_test

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/locale"
)

// 🧊 Internationalisation

// Internationalisation
type InternationalisationTemplData struct {
	li18ngo.TemplData
}

func (td InternationalisationTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "internationalisation.test",
		Description: "Internationalisation",
		Other:       "internationalisation",
	}
}

// 🧊 Localisation

// Internationalisation
type LocalisationTemplData struct {
	li18ngo.TemplData
}

func (td LocalisationTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "localisation.test",
		Description: "Localisation",
		Other:       "localisation",
	}
}

type GrafficoData struct{}

func (td GrafficoData) SourceID() string {
	return locale.TestGrafficoSourceID
}

// 🧊 Pavement Graffiti Report

// PavementGraffitiReportTemplData
type PavementGraffitiReportTemplData struct {
	GrafficoData
	Primary string
}

func (td PavementGraffitiReportTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pavement-graffiti-report.graffico.test",
		Description: "Report of graffiti found on a pavement",
		Other:       "Found graffiti on pavement; primary colour: '{{.Primary}}'",
	}
}

// ☢️ Wrong Source Id

// WrongSourceIDTemplData
type WrongSourceIDTemplData struct {
	GrafficoData
}

func (td WrongSourceIDTemplData) SourceID() string {
	return "FOO-BAR"
}

func (td WrongSourceIDTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "wrong-source-id.graffico.test",
		Description: "Incorrect Source ID which doesn't match the one in the localizer",
		Other:       "Message with wrong source ID",
	}
}
