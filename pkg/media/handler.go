package media

import (
	"net/http"
)

// Handler describes a common pattern for reading and writing resources from http requests and responses.
type Handler interface {

	// Read a resource from a request into the pointer resource.
	// Expected resource format determined by handler implementation.
	ReadResource(r *http.Request, resource any) error

	// Write a resource to the response with the given status code.
	// Resource format determined by handler implementation.
	WriteResponse(w http.ResponseWriter, statusCode int, resource any)

	// Write an error to the response.
	// Status code and error format determined by handler implementation.
	WriteError(w http.ResponseWriter, err error)
}
