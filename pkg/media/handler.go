package media

import (
	"net/http"
)
type Handler interface {
	ReadResource(r *http.Request, v any) error
	WriteResponse(w http.ResponseWriter, statusCode int, resource any)
	WriteError(w http.ResponseWriter, err error)
}