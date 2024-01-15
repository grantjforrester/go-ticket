package service

/*
 * The request cannot be proccessed due to a caller error.
 */
type RequestError struct {
	Message string
}

func (ve RequestError) Error() string {
	return ve.Message
}
