package repository

/*
 * The target resource of the request could not be found.
 */
type NotFoundError struct {
	Message string
}

func (nfe NotFoundError) Error() string {
	return nfe.Message
}

/*
 * The resource has been modified by another party.
 */
type ConflictError struct {
	Message string
}

func (ce ConflictError) Error() string {
	return ce.Message
}
