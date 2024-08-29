package translate_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/li18ngo/locale"
)

var _ = Describe("Text", func() {
	Context("failure", func() {
		When("Use has not been called", func() {
			It("should: raise the correct panic", func() {
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
