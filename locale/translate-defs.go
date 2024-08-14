package locale

const (
	Li18ngoSourceID = "github.com/snivilised/li18ngo"
)

type li18ngoTemplData struct{}

func (td li18ngoTemplData) SourceID() string {
	return Li18ngoSourceID
}
