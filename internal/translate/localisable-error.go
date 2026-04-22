package translate

// LocalisableError is an error that is translate-able (Localisable)
type LocalisableError struct {
	Data Localisable
}

// Error uses Describe rather than Text so that it is safe to use in library
// code before the host application has called Use.
func (le LocalisableError) Error() string {
	return Render(le.Data)
}
