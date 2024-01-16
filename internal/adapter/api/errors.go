package api

type PathNotFoundError struct {
	Message string
}

func (ve *PathNotFoundError) Error() string {
	return ve.Message
}
