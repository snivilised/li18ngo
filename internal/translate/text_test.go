package translate_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
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
