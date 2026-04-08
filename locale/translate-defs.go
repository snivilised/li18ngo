package locale

const (
	// Li18ngoSourceID defines the ID (by convention the repo URL) required
	// for i18n translation purposes.
	Li18ngoSourceID = "github.com/snivilised/li18ngo"

	// TestGrafficoSourceID test source ID
	TestGrafficoSourceID = "github.com/snivilised/graffico"
)

// Li18ngoTemplData is the base template data for all messages in this repo.
type Li18ngoTemplData struct{}

// SourceID returns the source ID for this message.
func (td Li18ngoTemplData) SourceID() string {
	return Li18ngoSourceID
}
