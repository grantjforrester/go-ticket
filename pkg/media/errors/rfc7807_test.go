package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grantjforrester/go-ticket/pkg/media/errors"
)

type MockError1 struct {
}

func (r MockError1) Error() string {
	return "mock error 1"
}

type MockError2 struct {
}

func (r MockError2) Error() string {
	return "mock error 2"
}

func TestShouldReturnDefaultError(t *testing.T) {
	// given
	defaultError := errors.RFC7807Error{TypeURI: "test:err:internalservererror", Status: 500, Title: "Internal Server Error"}
	errorMapper := errors.NewRFC7807ErrorMapper(defaultError)
	mockErr := MockError1{}

	// when
	statusCode, errorResponse := errorMapper.MapError(mockErr)

	// then
	assert.Equal(t, 500, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, errors.RFC7807Error{
		TypeURI: "test:err:internalservererror",
		Title:   "Internal Server Error",
		Status:  500,
		Detail:  "",
	}, errorResponse.(errors.RFC7807Error))
}

func TestShouldMatchErrorAndReturnRfc7807Error(t *testing.T) {
	// given
	defaultError := errors.RFC7807Error{TypeURI: "test:err:internalservererror", Status: 500, Title: "Internal Server Error"}
	errorMapper := errors.NewRFC7807ErrorMapper(defaultError)
	mockErr := MockError1{}

	// when
	errorMapper.RegisterError((*MockError1)(nil), errors.RFC7807Error{TypeURI: "test:err:notfound", Status: 404, Title: "Not Found"})
	errorMapper.RegisterError((*MockError2)(nil), errors.RFC7807Error{TypeURI: "test:err:badrequest", Status: 400, Title: "Bad Request"})
	statusCode, errorResponse := errorMapper.MapError(mockErr)

	// then
	assert.Equal(t, 404, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, errors.RFC7807Error{
		TypeURI: "test:err:notfound",
		Title:   "Not Found",
		Status:  404,
		Detail:  "mock error 1",
	}, errorResponse.(errors.RFC7807Error))
}

func TestShouldReturnDefaultErrorIfNoMatch(t *testing.T) {
	// given
	defaultError := errors.RFC7807Error{TypeURI: "test:err:internalservererror", Status: 500, Title: "Internal Server Error"}
	errorMapper := errors.NewRFC7807ErrorMapper(defaultError)
	mockErr := MockError2{}

	// when
	errorMapper.RegisterError((*MockError1)(nil), errors.RFC7807Error{TypeURI: "test:err:notfound", Status: 404, Title: "Not Found"})
	statusCode, errorResponse := errorMapper.MapError(mockErr)

	// then
	assert.Equal(t, 500, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, errors.RFC7807Error{
		TypeURI: "test:err:internalservererror",
		Title:   "Internal Server Error",
		Status:  500,
		Detail:  "",
	}, errorResponse.(errors.RFC7807Error))
}

func TestShouldMatchWrappedError(t *testing.T) {
	// given
	defaultError := errors.RFC7807Error{TypeURI: "test:err:internalservererror", Status: 500, Title: "Internal Server Error"}
	errorMapper := errors.NewRFC7807ErrorMapper(defaultError)
	mockErr := fmt.Errorf("wrapper: %w", MockError1{})

	// when
	errorMapper.RegisterError((*MockError1)(nil), errors.RFC7807Error{TypeURI: "test:err:notfound", Status: 404, Title: "Not Found"})
	statusCode, errorResponse := errorMapper.MapError(mockErr)

	// then
	assert.Equal(t, 404, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, errors.RFC7807Error{
		TypeURI: "test:err:notfound",
		Title:   "Not Found",
		Status:  404,
		Detail:  "mock error 1",
	}, errorResponse.(errors.RFC7807Error))
}
