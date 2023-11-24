package media

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

type Rfc7807ErrorMapping struct {
	Status int
	Title  string
}

type Rfc7807Error struct {
	TypeUri string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
}

type Rfc7807ErrorMapper struct {
	uriPrefix    string
	defaultError Rfc7807ErrorMapping
	errorMap     map[string]Rfc7807ErrorMapping
}

func NewRfc7807ErrorMapper(uriPrefix string, defaultError Rfc7807ErrorMapping) Rfc7807ErrorMapper {
	errorMap := make(map[string]Rfc7807ErrorMapping)
	return Rfc7807ErrorMapper{
		uriPrefix:    uriPrefix,
		defaultError: defaultError,
		errorMap:     errorMap,
	}
}

func (e *Rfc7807ErrorMapper) MapError(err error) (int, any) {
	if errorMapping, ok := e.matchError(err); ok {
		// return specific error
		return errorMapping.Status, e.mapError(err, errorMapping)
	} else {
		// return default error
		log.Println("Error: ", err.Error())
		return e.defaultError.Status, e.mapError(err, e.defaultError)
	}
}

func (e *Rfc7807ErrorMapper) RegisterError(err error, mapping Rfc7807ErrorMapping) {
	errorType := reflect.TypeOf(err).Elem().Name()

	e.errorMap[errorType] = mapping
}

func (e *Rfc7807ErrorMapper) matchError(err error) (Rfc7807ErrorMapping, bool) {
	for e1 := err; e1 != nil; {
		errorType := reflect.TypeOf(e1).Elem().Name()
		errorMapping, ok := e.errorMap[errorType]
		if ok {
			return errorMapping, ok
		}
		e1 = errors.Unwrap(e1)
	}
	return Rfc7807ErrorMapping{}, false
}

func (e *Rfc7807ErrorMapper) mapError(err error, mapping Rfc7807ErrorMapping) Rfc7807Error {
	return Rfc7807Error{
		TypeUri: e.formatUri(mapping.Title),
		Title:   mapping.Title,
		Status:  mapping.Status,
		Detail:  err.Error(),
	}
}

func (e *Rfc7807ErrorMapper) formatUri(title string) string {
	return strings.ReplaceAll(strings.ToLower(e.uriPrefix+title), " ", "")
}
