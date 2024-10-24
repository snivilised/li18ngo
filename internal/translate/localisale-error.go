package translate

// LocalisableError is an error that is translate-able (Localisable)
type LocalisableError struct {
	Data Localisable
}

func (le LocalisableError) Error() string {
	return Text(le.Data)
}
