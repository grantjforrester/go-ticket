package errors

// ErrorMapper describes a common pattern for converting Go errors into an object suitable
// for output to user.
type ErrorMapper interface {

	// MapError takes an Go error and returns a status code and an object for output.
	MapError(err error) (int, any)
}
