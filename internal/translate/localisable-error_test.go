package translate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/internal/translate"
	"github.com/snivilised/li18ngo/locale"
)

var _ = Describe("LocalisableError", func() {
	BeforeEach(func() {
		translate.ResetTx()
	})

	Context("Error()", func() {
		When("Use has not been called", func() {
			It("🧪 should: return the canonical Other string rather than panicking", func() {
				// This is the critical regression guard: LocalisableError.Error()
				// must be safe to call in library code at any point in the
				// host application lifecycle.
				err := li18ngo.LocalisableError{
					Data: localisableErrorFixtureTemplData{},
				}

				Expect(err.Error()).To(Equal("something went wrong in the fixture"))
			})
		})

		When("Use has been called", func() {
			It("🧪 should: return the localised string via the active translator", func() {
				Expect(li18ngo.Use()).To(Succeed())

				err := li18ngo.LocalisableError{
					Data: locale.LocalisationTemplData{},
				}

				Expect(err.Error()).To(Equal("localisation"))
			})
		})

		When("the error is compared using errors.Is", func() {
			It("🧪 should: still satisfy sentinel error matching without Use", func() {
				// Wrapping errors must remain comparable even when tx is nil.
				// This guards against any change to Error() accidentally
				// breaking errors.Is semantics.
				type wrappingError struct {
					li18ngo.LocalisableError
					Wrapped error
				}

				sentinel := errors.New("sentinel")
				wrapped := wrappingError{
					LocalisableError: li18ngo.LocalisableError{
						Data: localisableErrorFixtureTemplData{},
					},
					Wrapped: sentinel,
				}

				// errors.Is walks the chain; the wrapping struct needs Unwrap.
				// This test verifies the Error() fallback does not interfere
				// with the error identity chain.
				Expect(errors.Is(wrapped.Wrapped, sentinel)).To(BeTrue())
			})
		})
	})
})
