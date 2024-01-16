package authz

/*
 * Authorization failed.
 */
type AuthorizationError struct {
	Message string
}

func (ae AuthorizationError) Error() string {
	return ae.Message
}
