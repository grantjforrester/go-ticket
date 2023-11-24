package authz

import "context"

type Operation string

type Authorizer interface {
	IsAuthorized(context.Context, Operation) error
}

type AlwaysAuthorize struct {
}

var _ Authorizer = (*AlwaysAuthorize)(nil)

func (a AlwaysAuthorize) IsAuthorized(c context.Context, operation Operation) error {
	return nil
}
