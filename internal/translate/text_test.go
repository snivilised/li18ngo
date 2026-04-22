package translate_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
)

var _ = Describe("Text", func() {
	Context("failure", func() {
		When("Use has not been called", func() {
			It("🧪 should: raise the correct panic", func() {
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok && errors.Is(err, li18ngo.ErrSafePanicWarning) {
							return
						}
					}
					Fail("safe panic warning has not occurred")
				}()

				_ = li18ngo.Text(locale.InternationalisationTemplData{})
				Fail("expected panic to occur")
			})
		})
	})

	Context("Use", func() {
		Context("native package", func() {
			When("invoked without arguments", func() {
				It("🧪 should: get default text", func() {
					Expect(li18ngo.Use()).To(Succeed())

					text := li18ngo.Text(locale.LocalisationTemplData{})
					Expect(text).To(Equal("localisation"))
				})
			})
		})

		Context("foreign package", func() {
			When("invoked without arguments", func() {
				It("🧪 should: get default text", func() {
					Expect(li18ngo.Use()).To(Succeed())

					Expect(li18ngo.Text(ParadiseLostTemplData{})).To(
						Equal("paradise lost"),
						"new localizer should be created for unknown source id",
					)

					Expect(li18ngo.Text(AtTheGatesOfSilentMemoryTemplData{})).To(
						Equal("at the gates of silent memory"),
						"should use localizer that was created on the fly",
					)
				})
			})
		})
	})
})

var _ = Describe("Text", func() {
	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("after Register but before Use", func() {
		When("Text is called by application code", func() {
			It("🧪 should: not panic because Register activates the translator", func() {
				// Register activates tx just as Use does; application code
				// calling Text after a library has called Register should
				// therefore not encounter a panic.
				Expect(li18ngo.Register()).To(Succeed())

				Expect(func() {
					_ = li18ngo.Text(locale.LocalisationTemplData{})
				}).NotTo(Panic())
			})
		})
	})

	Context("without Use or Register", func() {
		When("Text is called by application code", func() {
			It("🧪 should: still panic with ErrSafePanicWarning", func() {
				// Regression guard: the existing panic contract for Text must
				// be preserved. Describe's nil-safety must not weaken Text.
				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok && errors.Is(err, li18ngo.ErrSafePanicWarning) {
							return
						}
					}
					Fail("safe panic warning has not occurred")
				}()

				_ = li18ngo.Text(locale.InternationalisationTemplData{})
				Fail("expected panic to occur")
			})
		})
	})
})
