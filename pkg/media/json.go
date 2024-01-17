package media

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/grantjforrester/go-ticket/pkg/media/errors"
)

// JSONHandler is a Handler implementation for reading and writing HTTP requests and responses containing JSON.
type JSONHandler struct {
	ErrorMap errors.ErrorMapper
}

var _ Handler = (*JSONHandler)(nil)

// Decodes the request body as JSON into resource.
// Returns MediaError if request body is not valid JSON or cannot be marshalled into given resource struct.
func (j JSONHandler) ReadResource(req *http.Request, resource any) error {
	jsonBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read resource: %w", err)
	}

	err = json.Unmarshal(jsonBytes, resource)
	switch e := err.(type) {
	case *json.SyntaxError:
		return MediaError{Message: "invalid json"}
	case *json.UnmarshalTypeError:
		return MediaError{Message: fmt.Sprintf("invalid type for field: %s", e.Field)}
	}

	return nil
}

// Encodes the resource into JSON and writes into response body.  Sets Content-Type to "application/json" and
// status code on the response.
// Panics if the given resource cannot be encoded to JSON.
func (j JSONHandler) WriteResponse(resp http.ResponseWriter, status int, resource any) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	if resource != nil {
		err := json.NewEncoder(resp).Encode(resource)
		if err != nil {
			log.Panicf("Could not encode resource to JSON: %v", err)
		}
	}
}

// Encodes the error into JSON and writes into response body. Sets Content-type to "application/json".
// The formatting of the JSON and the status code returned are retrieved from the handler's error map. See
// ErrorMap.MapError.
func (j JSONHandler) WriteError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/json")
	statusCode, errorResource := j.ErrorMap.MapError(err)
	j.WriteResponse(resp, statusCode, errorResource)
}
