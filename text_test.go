package li18ngo_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/lab"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
	nef "github.com/snivilised/nefilim"
	"golang.org/x/text/language"
)

const (
	relative             = "test/data/l10n"
	TestGrafficoSourceID = "github.com/snivilised/graffico"
)

type grafficoData struct{}

func (td grafficoData) SourceID() string {
	return TestGrafficoSourceID
}

// =============================================================================
// 📨 PavementGraffitiReportGraffico
//
// A test message for pavement graffiti reporting.
// =============================================================================

// PavementGraffitiReportGrafficoTemplData Report of graffiti found on a
// pavement.
type PavementGraffitiReportGrafficoTemplData struct {
	grafficoData
	// Primary is the primary colour of the graffiti
	Primary string
}

// Message returns the i18n message for PavementGraffitiReportGrafficoTemplData.
func (td PavementGraffitiReportGrafficoTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pavement-graffiti-report.graffico.test",
		Description: "Report of graffiti found on a pavement",
		Other:       "Found graffiti on pavement; primary colour: '{{.Primary}}'",
	}
}

// NewPavementGraffitiReportGrafficoTemplData creates a new
// PavementGraffitiReportGrafficoTemplData.
func NewPavementGraffitiReportGrafficoTemplData(primary string) PavementGraffitiReportGrafficoTemplData {
	return PavementGraffitiReportGrafficoTemplData{
		grafficoData: grafficoData{},
		Primary:      primary,
	}
}

// =============================================================================
// 📨 WrongSourceIDGrafficoTest
//
// This message has the wrong source ID and should be ignored by i18n-gen.
// =============================================================================

// WrongSourceIDGrafficoTemplData Message with wrong source ID for testing
// purposes.
type WrongSourceIDGrafficoTemplData struct {
	grafficoData
}

// Message returns the i18n message for WrongSourceIDGrafficoTestTemplData.
func (td WrongSourceIDGrafficoTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "wrong-source-id.graffico.test",
		Description: "Message with wrong source ID for testing purposes",
		Other:       "This message should be ignored by i18n-gen.",
	}
}

var _ = Describe("Text", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile li18ngo.TranslationFiles
	)

	BeforeAll(func() {
		repo = lab.Repo("")
		l10nPath = lab.Path(repo, relative)
		queryFS := nef.NewMakeDirFS(nef.Rel{
			Root: repo,
		})

		Expect(queryFS.DirectoryExists(relative)).To(BeTrue(),
			fmt.Sprintf("l10n '%v' path does not exist", relative),
		)

		testTranslationFile = li18ngo.TranslationFiles{
			li18ngo.Li18ngoSourceID: li18ngo.TranslationSource{
				Name: "test",
			},
		}
	})

	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("Default Language", func() {
		BeforeEach(func() {
			if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
				o.Tag = li18ngo.DefaultLanguage
				o.From.Sources = testTranslationFile
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("given: ThirdPartyError", func() {
			It("🧪 should: contain the third party error text", func() {
				if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}

				err := locale.NewThirdPartyWrapperError(errors.New("computer says no"))
				Expect(err.Error()).To(ContainSubstring("computer says no"))
			})

			Context("Text", func() {
				Context("given: a template data instance", func() {
					It("🧪 should: evaluate translated text", func() {
						Expect(li18ngo.Text(locale.ThirdPartyWrapperErrorTemplData{
							Wrapped: "out of stock",
						})).NotTo(BeNil())
					})
				})
			})
		})
	})

	Context("Foreign Language (en-US)", func() {
		BeforeEach(func() {
			if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
				o.Tag = language.AmericanEnglish
				o.From.Path = l10nPath
				o.From.Sources = testTranslationFile
				o.DefaultIsAcceptable = false
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("Text", func() {
			Context("given: a template data instance", func() {
				It("🧪 should: evaluate translated text(internationalization)", func() {
					text := li18ngo.Text(locale.InternationalisationTemplData{})
					Expect(text).To(
						Equal("internationalization"),
					)
				})

				It("🧪 should: evaluate translated text(localization)", func() {
					Expect(li18ngo.Text(locale.LocalisationTemplData{})).To(
						Equal("localization"),
					)
				})
			})
		})
	})

	Context("Multiple Sources", func() {
		Context("Foreign Language", func() {
			It("🧪 should: translate text with the correct localizer", func() {
				if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From = li18ngo.LoadFrom{
						Path: l10nPath,
						Sources: li18ngo.TranslationFiles{
							li18ngo.Li18ngoSourceID:     li18ngo.TranslationSource{Name: "test"},
							locale.TestGrafficoSourceID: li18ngo.TranslationSource{Name: "test.graffico"},
						},
					}
					o.DefaultIsAcceptable = false
				}); err != nil {
					Fail(err.Error())
				}
				actual := li18ngo.Text(PavementGraffitiReportGrafficoTemplData{
					Primary: "Violet",
				})
				Expect(actual).To(Equal(expectUS))
			})
		})
	})

	When("extendio source not provided", func() {
		Context("Default Language", func() {
			It("🧪 should: create factory that contains the extendio source", func() {
				if err := li18ngo.Use(func(o *li18ngo.UseOptions) {
					o.Tag = language.BritishEnglish
					o.From = li18ngo.LoadFrom{
						Path: l10nPath,
						Sources: li18ngo.TranslationFiles{
							locale.TestGrafficoSourceID: li18ngo.TranslationSource{Name: "test.graffico"},
						},
					}
				}); err != nil {
					Fail(err.Error())
				}

				actual := li18ngo.Text(locale.InternationalisationTemplData{})
				Expect(actual).To(Equal("internationalisation"))
			})
		})
	})

	Context("translation source contains path", func() {
		Context("Foreign Language", func() {
			Context("given: source path exists", func() {
				It("🧪 should: create localizer from source path", func() {
					actual := testPavementGraffitiUS(&textTE{
						sourcePath:        l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: path exists", func() {
				It("🧪 should: create localizer from path", func() {
					actual := testPavementGraffitiUS(&textTE{
						path:              l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: neither path exists", func() {
				It("🧪 should: create localizer using default language", func() {
					actual := testPavementGraffitiUS(&textTE{
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectGB))
				})
			})

			Context("given: neither path exists", func() {
				XIt("🧪 should: create localizer using default language", func() {
					// Its not clear what this test is supposed to do
					//
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"could not load translations for") {
							Fail("translation file not available with exe")
						}
					}()

					_ = testPavementGraffitiUS(&textTE{
						defaultAcceptable: false,
					})
				})
			})
		})
	})
})
