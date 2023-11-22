package media_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	em "github.com/grantjforrester/go-ticket/pkg/media"
)

type MockError1 struct {

}

func (r *MockError1) Error() string {
	return "mock error 1"
}

type MockError2 struct {

}

func (r *MockError2) Error() string {
	return "mock error 2"
}

func TestShouldReturnDefaultError(t *testing.T) {
	// given
	defaultError := em.Rfc7807ErrorMapping{Status: 500, Title: "Internal Server Error"}
	rfc7807_error_mapper := em.NewRfc7807ErrorMapper("test:err:", defaultError)
	mockErr := MockError1{}

	// when
	statusCode, errorResponse := rfc7807_error_mapper.MapError(&mockErr)

	// then
	assert.Equal(t, 500, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, em.Rfc7807Error{
		TypeUri: "test:err:internalservererror",
		Title: "Internal Server Error",
		Status: 500,
		Detail: "mock error 1",
	}, errorResponse.(em.Rfc7807Error))
}

func TestShouldMatchErrorAndReturnRfc7807Error(t *testing.T) {
	// given
	defaultError := em.Rfc7807ErrorMapping{Status: 500, Title: "Internal Server Error"}
	rfc7807_error_mapper := em.NewRfc7807ErrorMapper("test:err:", defaultError)
	mockErr := MockError1{}

	// when
	rfc7807_error_mapper.RegisterError((*MockError1)(nil), em.Rfc7807ErrorMapping{Status: 404, Title: "Not Found"})
	rfc7807_error_mapper.RegisterError((*MockError2)(nil), em.Rfc7807ErrorMapping{Status: 400, Title: "Bad Request"})
	statusCode, errorResponse := rfc7807_error_mapper.MapError(&mockErr)

	// then
	assert.Equal(t, 404, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, em.Rfc7807Error{
		TypeUri: "test:err:notfound",
		Title: "Not Found",
		Status: 404,
		Detail: "mock error 1",
	}, errorResponse.(em.Rfc7807Error))
}

func TestShouldReturnDefaultErrorIfNoMatch(t *testing.T) {
	// given
	defaultError := em.Rfc7807ErrorMapping{Status: 500, Title: "Internal Server Error"}
	rfc7807_error_mapper := em.NewRfc7807ErrorMapper("test:err:", defaultError)
	mockErr := MockError2{}

	// when
	rfc7807_error_mapper.RegisterError((*MockError1)(nil), em.Rfc7807ErrorMapping{Status: 404, Title: "Not Found"})
	statusCode, errorResponse := rfc7807_error_mapper.MapError(&mockErr)

	// then
	assert.Equal(t, 500, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, em.Rfc7807Error{
		TypeUri: "test:err:internalservererror",
		Title: "Internal Server Error",
		Status: 500,
		Detail: "mock error 2",
	}, errorResponse.(em.Rfc7807Error))
}

func TestShouldMatchWrappedError(t *testing.T) {
	// given
	defaultError := em.Rfc7807ErrorMapping{Status: 500, Title: "Internal Server Error"}
	rfc7807_error_mapper := em.NewRfc7807ErrorMapper("test:err:", defaultError)
	mockErr := fmt.Errorf("wrapper: %w", &MockError1{})

	// when
	rfc7807_error_mapper.RegisterError((*MockError1)(nil), em.Rfc7807ErrorMapping{Status: 404, Title: "Not Found"})
	statusCode, errorResponse := rfc7807_error_mapper.MapError(mockErr)

	// then
	assert.Equal(t, 404, statusCode)
	assert.NotNil(t, errorResponse)
	assert.Equal(t, em.Rfc7807Error{
		TypeUri: "test:err:notfound",
		Title: "Not Found",
		Status: 404,
		Detail: "wrapper: mock error 1",
	}, errorResponse.(em.Rfc7807Error))
}