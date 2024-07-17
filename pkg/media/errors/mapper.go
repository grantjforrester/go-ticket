package errors

// ErrorMapper describes a common pattern for converting Go errors into error codes and error objects
// suitable for output to users.
type ErrorMapper interface {

	// MapError takes an Go error and returns an error code and an error object for output.
	MapError(err error) (int, any)
}
