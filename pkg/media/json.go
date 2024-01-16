package media

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/grantjforrester/go-ticket/pkg/media/errors"
)

type JSONHandler struct {
	ErrorMap errors.ErrorMapper
}

var _ Handler = (*JSONHandler)(nil)

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

func (j JSONHandler) WriteError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/json")
	statusCode, errorResource := j.ErrorMap.MapError(err)
	j.WriteResponse(resp, statusCode, errorResource)
}
