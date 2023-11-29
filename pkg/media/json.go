package media

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JSONHandler struct {
	ErrorMap ErrorMapper
}

func (j JSONHandler) ReadResource(req *http.Request, resource any) error {
	jsonBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("json_media_handler: failed to read resource: %w", err)
	}

	return json.Unmarshal(jsonBytes, resource)
}

func (j JSONHandler) WriteResponse(resp http.ResponseWriter, status int, resource any) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	if resource != nil {
		json.NewEncoder(resp).Encode(resource)
	}
}

func (j JSONHandler) WriteError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/json")
	statusCode, errorResource := j.ErrorMap.MapError(err)
	j.WriteResponse(resp, statusCode, errorResource)
}
