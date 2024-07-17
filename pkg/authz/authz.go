package authz

import "context"

// Authorizer describes a common pattern for authorizing operations.
// Each implementation of Authorizer determines the authorization mechanisms used.
type Authorizer interface {

	// IsAuthorized performs the relevant authorization checks for the operation
	// using the request context. If not authorized then AuthorizationError is returned.
	IsAuthorized(context.Context, Operation) error
}

// Operation represents a distinct function of the system that must be authorized before execution.
type Operation string

// AlwaysAuthorize is an implementation of Authorizer that always authorizes the operation.
type AlwaysAuthorize struct {
}

var _ Authorizer = (*AlwaysAuthorize)(nil)

// Always returns nil i.e. the operation is authorized.
func (a AlwaysAuthorize) IsAuthorized(_ context.Context, operation Operation) error {
	return nil
}
