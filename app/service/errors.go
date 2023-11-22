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

/*
 * The request failed due to a hardware or software failure.
 */
type SystemError struct {
	Message string
}

func (se SystemError) Error() string {
	return se.Message
}