package i18n

// CLIENT-TODO: Should be updated to use url of the implementing project,
// so should not be left as astrolib. (this should be set by auto-check)
const Li18ngoSourceID = "github.com/snivilised/li18ngo"

type li18ngoTemplData struct{}

func (td li18ngoTemplData) SourceID() string {
	return Li18ngoSourceID
}
