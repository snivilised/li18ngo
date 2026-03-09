package locale_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"   
)

func TestI18n(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Locale Suite")
}
