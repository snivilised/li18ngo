package translate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // ginkgo ok
)

func TestTranslate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translate Suite")
}
