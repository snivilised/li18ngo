package locale

const (
	Li18ngoSourceID = "github.com/snivilised/li18ngo"
)

type Li18ngoTemplData struct{}

func (td Li18ngoTemplData) SourceID() string {
	return Li18ngoSourceID
}
