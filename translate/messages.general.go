package translate

type Li18ngoTemplData struct{}

func (td Li18ngoTemplData) SourceID() string {
	return Li18ngoSourceID
}

// 🧊 Internationalisation

// Internationalisation
type InternationalisationTemplData struct {
	Li18ngoTemplData
}

func (td InternationalisationTemplData) Message() *Message {
	return &Message{
		ID:          "internationalisation.general",
		Description: "Internationalisation",
		Other:       "internationalisation",
	}
}

// 🧊 Localisation

// Internationalisation
type LocalisationTemplData struct {
	Li18ngoTemplData
}

func (td LocalisationTemplData) Message() *Message {
	return &Message{
		ID:          "localisation.general",
		Description: "Localisation",
		Other:       "localisation",
	}
}
