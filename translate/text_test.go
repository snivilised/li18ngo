package translate_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	"github.com/snivilised/li18ngo/internal/helpers"
	"github.com/snivilised/li18ngo/translate"

	"golang.org/x/text/language"
)

var _ = Describe("Text", Ordered, func() {
	var (
		repo                string
		l10nPath            string
		testTranslationFile translate.TranslationFiles
	)

	BeforeAll(func() {
		repo = helpers.Repo("")
		l10nPath = helpers.Path(repo, "test/data/l10n")
		// Expect(utils.FolderExists(l10nPath)).To(BeTrue())

		testTranslationFile = translate.TranslationFiles{
			translate.Li18ngoSourceID: translate.TranslationSource{Name: "test"},
		}
	})

	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("Default Language", func() {
		BeforeEach(func() {
			if err := translate.Use(func(o *translate.UseOptions) {
				o.Tag = translate.DefaultLanguage
				o.From.Sources = testTranslationFile
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("given: ThirdPartyError", func() {
			It("🧪 should: contain the third party error text", func() {
				if err := translate.Use(func(o *translate.UseOptions) {
					o.Tag = language.BritishEnglish
				}); err != nil {
					Fail(err.Error())
				}

				err := translate.NewThirdPartyErr(errors.New("computer says no"))
				Expect(err.Error()).To(ContainSubstring("computer says no"))
			})

			Context("Text", func() {
				Context("given: a template data instance", func() {
					It("🧪 should: evaluate translated text", func() {
						Expect(translate.Text(translate.ThirdPartyErrorTemplData{
							Error: errors.New("out of stock"),
						})).NotTo(BeNil())
					})
				})
			})
		})
	})

	Context("Foreign Language", func() {
		BeforeEach(func() {
			if err := translate.Use(func(o *translate.UseOptions) {
				o.Tag = language.AmericanEnglish
				o.From.Path = l10nPath
				o.From.Sources = testTranslationFile
			}); err != nil {
				Fail(err.Error())
			}
		})

		Context("Text", func() {
			Context("given: a template data instance", func() {
				XIt("🧪 should: evaluate translated text(internationalization)", func() {
					Expect(translate.Text(translate.InternationalisationTemplData{})).To(
						Equal("internationalization"),
					)
				})

				XIt("🧪 should: evaluate translated text(localization)", func() {
					Expect(translate.Text(translate.LocalisationTemplData{})).To(
						Equal("localization"),
					)
				})
			})
		})
	})

	Context("Multiple Sources", func() {
		Context("Foreign Language", func() {
			XIt("🧪 should: translate text with the correct localizer", func() {
				if err := translate.Use(func(o *translate.UseOptions) {
					o.Tag = language.AmericanEnglish
					o.From = translate.LoadFrom{
						Path: l10nPath,
						Sources: translate.TranslationFiles{
							translate.Li18ngoSourceID: translate.TranslationSource{Name: "test"},
							GrafficoSourceID:          translate.TranslationSource{Name: "test.graffico"},
						},
					}
				}); err != nil {
					Fail(err.Error())
				}
				actual := translate.Text(PavementGraffitiReportTemplData{
					Primary: "Violet",
				})
				Expect(actual).To(Equal(expectUS))
			})
		})
	})

	When("extendio source not provided", func() {
		Context("Default Language", func() {
			It("🧪 should: create factory that contains the extendio source", func() {
				if err := translate.Use(func(o *translate.UseOptions) {
					o.Tag = language.BritishEnglish
					o.From = translate.LoadFrom{
						Path: l10nPath,
						Sources: translate.TranslationFiles{
							GrafficoSourceID: translate.TranslationSource{Name: "test.graffico"},
						},
					}
				}); err != nil {
					Fail(err.Error())
				}

				actual := translate.Text(translate.InternationalisationTemplData{})
				Expect(actual).To(Equal("internationalisation"))
			})
		})
	})

	Context("translation source contains path", func() {
		Context("Foreign Language", func() {
			Context("given: source path exists", func() {
				XIt("🧪 should: create localizer from source path", func() {
					actual := testTranslationPath(&textTE{
						sourcePath:        l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: path exists", func() {
				XIt("🧪 should: create localizer from path", func() {
					actual := testTranslationPath(&textTE{
						path:              l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: neither path exists", func() {
				It("🧪 should: create localizer using default language", func() {
					actual := testTranslationPath(&textTE{
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectGB))
				})
			})

			Context("given: neither path exists", func() {
				XIt("🧪 should: create localizer using default language", func() {
					defer func() {
						pe := recover()
						if err, ok := pe.(error); !ok || !strings.Contains(err.Error(),
							"could not load translations for") {
							Fail("translation file not available with exe")
						}
					}()

					_ = testTranslationPath(&textTE{
						defaultAcceptable: false,
					})

					Fail("❌ expected panic due translation file not available with exe")
				})
			})
		})
	})
})
