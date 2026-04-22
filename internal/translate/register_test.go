package translate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
)

var _ = Describe("Register", func() {
	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("library bootstrap", func() {
		When("Register is called without options", func() {
			It("🧪 should: succeed and activate the translator", func() {
				// Register has the same mechanics as Use; a library calling
				// Register before the host calls Use should not cause a panic
				// when Text or Describe is subsequently invoked.
				Expect(li18ngo.Register()).To(Succeed())
			})
		})

		When("Register is called before Use", func() {
			It("🧪 should: allow Describe to return localised text", func() {
				Expect(li18ngo.Register()).To(Succeed())

				result := li18ngo.Render(locale.LocalisationTemplData{})
				Expect(result).To(Equal("localisation"))
			})
		})

		When("Register is called before Use", func() {
			It("🧪 should: not cause Text to panic when host later calls Use", func() {
				// Simulates: library calls Register at init time, host calls
				// Use at startup - the negotiation path should merge sources.
				Expect(li18ngo.Register()).To(Succeed())
				Expect(li18ngo.Use()).To(Succeed())

				result := li18ngo.Text(locale.LocalisationTemplData{})
				Expect(result).To(Equal("localisation"))
			})
		})

		When("Use is called before Register", func() {
			It("🧪 should: not cause Text to panic when library later calls Register", func() {
				// Simulates: host bootstraps first, library registers late -
				// the more common ordering in practice.
				Expect(li18ngo.Use()).To(Succeed())
				Expect(li18ngo.Register()).To(Succeed())

				result := li18ngo.Text(locale.LocalisationTemplData{})
				Expect(result).To(Equal("localisation"))
			})
		})
	})
})
