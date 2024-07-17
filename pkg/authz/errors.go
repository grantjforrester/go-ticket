package authz

/*
 * AuthorizationError is returned when authorization of an operation failed.
 */
type AuthorizationError struct {
	Message string
}

func (ae AuthorizationError) Error() string {
	return ae.Message
}
