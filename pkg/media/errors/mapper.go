package errors

// ErrorMapper describes a common pattern for converting Go errors into an error object suitable
// for output to user.
type ErrorMapper interface {

	// MapError takes an Go error and returns an error code and an error object for output.
	MapError(err error) (int, any)
}
