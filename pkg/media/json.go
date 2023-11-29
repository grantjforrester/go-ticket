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

var _ Handler = (*JSONHandler)(nil)

func (j JSONHandler) ReadResource(req *http.Request, resource any) error {
	jsonBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("failed to read resource: %w", err)
	}

	return json.Unmarshal(jsonBytes, resource)
}

func (j JSONHandler) WriteResponse(resp http.ResponseWriter, status int, resource any) error {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(status)
	if resource != nil {
		err := json.NewEncoder(resp).Encode(resource)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	}

	return nil
}

func (j JSONHandler) WriteError(resp http.ResponseWriter, err error) {
	resp.Header().Set("Content-Type", "application/json")
	statusCode, errorResource := j.ErrorMap.MapError(err)
	j.WriteResponse(resp, statusCode, errorResource)
}
