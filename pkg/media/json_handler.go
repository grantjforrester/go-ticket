package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type JsonHandler struct {
	ErrorMap ErrorMapper
}

func (j JsonHandler) ReadResource(r *http.Request, v any) error {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("json_media_handler: failed to read resource: %w", err)
	}

	return json.Unmarshal(jsonBytes, v)
}

func (j JsonHandler) WriteResponse(w http.ResponseWriter, statusCode int, resource any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if resource != nil {
		json.NewEncoder(w).Encode(resource)
	}
}

func (j JsonHandler) WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	statusCode, errorResource := j.ErrorMap.MapError(err)
	j.WriteResponse(w, statusCode, errorResource)
}
