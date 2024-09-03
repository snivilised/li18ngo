package translate_test

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // ginkgo ok
)

func TestTranslate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translate Suite")
}

type ForeignTemplData struct{}

func (td ForeignTemplData) SourceID() string {
	return "github.com/snivilised/foreign"
}

type ParadiseLostTemplData struct {
	ForeignTemplData
}

func (td ParadiseLostTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "paradise-lost.message",
		Description: "paradise lost",
		Other:       "paradise lost",
	}
}

type AtTheGatesOfSilentMemoryTemplData struct {
	ForeignTemplData
}

func (td AtTheGatesOfSilentMemoryTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "at-the-gates-of-silent-memory.message",
		Description: "at the gates of silent memory",
		Other:       "at the gates of silent memory",
	}
}
