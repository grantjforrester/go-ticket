package api

type PathNotFoundError struct {
	Message string
}

func (ve *PathNotFoundError) Error() string {
	return ve.Message
}

type InvalidQueryError struct {
	Message string
}

func (ve *InvalidQueryError) Error() string {
	return ve.Message
}