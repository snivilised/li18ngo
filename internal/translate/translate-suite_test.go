package translate_test

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/snivilised/li18ngo/locale"
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

// ---------------------------------------------------------------------------
// Test fixtures
//
// Two minimal Localisable implementations used only in this file, keeping
// each describe block self-contained and avoiding cross-file fixture coupling.
// ---------------------------------------------------------------------------

// localisableErrorFixtureTemplData is a static non-error message used to
// exercise LocalisableError.Error() in isolation.
type localisableErrorFixtureTemplData struct {
	locale.Li18ngoTemplData
}

func (td localisableErrorFixtureTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "localisable-error-fixture",
		Description: "fixture message for LocalisableError tests",
		Other:       "something went wrong in the fixture",
	}
}
