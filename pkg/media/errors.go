package media

// MediaError is returned when a received resource could not be correctly parsed.
type MediaError struct {
	Message string
}

func (me MediaError) Error() string {
	return me.Message
}
