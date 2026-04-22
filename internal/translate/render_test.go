package translate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
)

var _ = Describe("Render", func() {
	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("before Use has been called", func() {
		When("tx is nil", func() {
			It("🧪 should: return the canonical Other string without panicking", func() {
				// This is the core library-safety guarantee: Render must never
				// panic regardless of initialisation state.
				result := li18ngo.Render(localisableErrorFixtureTemplData{})
				Expect(result).To(Equal("something went wrong in the fixture"))
			})
		})
	})

	Context("after Use has been called", func() {
		When("a native package message is described", func() {
			It("🧪 should: return the localised string via the active translator", func() {
				Expect(li18ngo.Use()).To(Succeed())

				result := li18ngo.Render(locale.LocalisationTemplData{})
				Expect(result).To(Equal("localisation"))
			})
		})
	})
})
