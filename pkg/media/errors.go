package media

/*
 * Returned when unexpexted media input received.
 */
type MediaError struct {
	Message string
}

func (me MediaError) Error() string {
	return me.Message
}
