package media_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	jsonMedia "github.com/grantjforrester/go-ticket/pkg/media"
)

type validStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestShouldUnmarshallValidResource(t *testing.T) {
	// Given
	handler := jsonMedia.JsonHandler{}
	request := mockRequest(validStruct{Foo: "mock foo", Bar: 1})
	resource := validStruct{}

	// When
	err := handler.ReadResource(request, &resource)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, validStruct{Foo: "mock foo", Bar: 1}, resource)
}

func mockRequest(body any) *http.Request {
	jsonBytes, _ := json.Marshal(body)
	request, _ := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(jsonBytes))
	return request
}
