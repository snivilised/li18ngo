package translate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo/translate"
	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translate Suite")
}

type textTE struct {
	message           string
	path              string
	sourcePath        string
	name              string
	defaultAcceptable bool
}

const (
	expectUS = "Found graffiti on sidewalk; primary color: 'Violet'"
	expectGB = "Found graffiti on pavement; primary colour: 'Violet'"
)

func testTranslationPath(entry *textTE) string {
	// this test form required, because DescribeTable can't be used
	// due to not being able to setup state correctly, eg l10nPath
	//
	if err := translate.Use(func(o *translate.UseOptions) {
		o.Tag = language.AmericanEnglish
		o.DefaultIsAcceptable = entry.defaultAcceptable
		o.From = translate.LoadFrom{
			Path: entry.path,
			Sources: translate.TranslationFiles{
				GrafficoSourceID: translate.TranslationSource{
					Path: entry.sourcePath,
					Name: "test.graffico",
				},
			},
		}
	}); err != nil {
		Fail(err.Error())
	}

	return translate.Text(PavementGraffitiReportTemplData{
		Primary: "Violet",
	})
}
