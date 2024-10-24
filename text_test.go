package li18ngo_test

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/lab"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
	nef "github.com/snivilised/nefilim"
	"golang.org/x/text/language"
)

const (
	relative = "test/data/l10n"
)

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

				err := li18ngo.NewThirdPartyErr(errors.New("computer says no"))
				Expect(err.Error()).To(ContainSubstring("computer says no"))
			})

			Context("Text", func() {
				Context("given: a template data instance", func() {
					It("🧪 should: evaluate translated text", func() {
						Expect(li18ngo.Text(li18ngo.ThirdPartyErrorTemplData{
							Error: errors.New("out of stock"),
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
					text := li18ngo.Text(InternationalisationTemplData{})
					Expect(text).To(
						Equal("internationalization"),
					)
				})

				It("🧪 should: evaluate translated text(localization)", func() {
					Expect(li18ngo.Text(LocalisationTemplData{})).To(
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
				actual := li18ngo.Text(PavementGraffitiReportTemplData{
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

				actual := li18ngo.Text(InternationalisationTemplData{})
				Expect(actual).To(Equal("internationalisation"))
			})
		})
	})

	Context("translation source contains path", func() {
		Context("Foreign Language", func() {
			Context("given: source path exists", func() {
				It("🧪 should: create localizer from source path", func() {
					actual := testTranslationPath(&textTE{
						sourcePath:        l10nPath,
						defaultAcceptable: true,
					})

					Expect(actual).To(Equal(expectUS))
				})
			})

			Context("given: path exists", func() {
				It("🧪 should: create localizer from path", func() {
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
					// Its not clear what this test is supposed to do
					//
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
				})
			})
		})
	})
})
