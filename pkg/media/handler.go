package media

import (
	"net/http"
)

// Handler describes a common pattern for reading and writing resources from HTTP requests and responses.
type Handler interface {

	// Read a resource from a request into the given struct.
	// The expected resource format is determined by the handler implementation.
	// If the expected resource cannot be parsed correctly a MediaError is returned.
	ReadResource(r *http.Request, resource any) error

	// Writes the given resource to the response writer with the given status code.
	// The resource format is determined by the handler implementation.
	WriteResponse(w http.ResponseWriter, statusCode int, resource any)

	// Writes the given error to the response.
	// Status code and error format is determined by the handler implementation.
	WriteError(w http.ResponseWriter, err error)
}
