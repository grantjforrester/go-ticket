package media

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

type RFC7807ErrorMapping struct {
	Status int
	Title  string
}

type RFC7807Error struct {
	TypeURI string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
}

type RFC7807ErrorMapper struct {
	uriPrefix    string
	defaultError RFC7807ErrorMapping
	errorMap     map[string]RFC7807ErrorMapping
}

func NewRFC7807ErrorMapper(uriPrefix string, defaultError RFC7807ErrorMapping) RFC7807ErrorMapper {
	errorMap := make(map[string]RFC7807ErrorMapping)
	return RFC7807ErrorMapper{
		uriPrefix:    uriPrefix,
		defaultError: defaultError,
		errorMap:     errorMap,
	}
}

func (e *RFC7807ErrorMapper) MapError(err error) (int, any) {
	if errorMapping, ok := e.matchError(err); ok {
		// return specific error
		return errorMapping.Status, e.formatError(err, errorMapping)
	} else {
		// return default error
		log.Println("Error: ", err.Error())
		return e.defaultError.Status, e.formatError(err, e.defaultError)
	}
}

func (e *RFC7807ErrorMapper) RegisterError(err error, mapping RFC7807ErrorMapping) {
	errorType := reflect.TypeOf(err).Elem().Name()
	e.errorMap[errorType] = mapping
}

func (e *RFC7807ErrorMapper) matchError(err error) (RFC7807ErrorMapping, bool) {
	for e1 := err; e1 != nil; {
		errorType := reflect.TypeOf(e1).Elem().Name()
		errorMapping, ok := e.errorMap[errorType]
		if ok {
			return errorMapping, ok
		}
		e1 = errors.Unwrap(e1)
	}
	return RFC7807ErrorMapping{}, false
}

func (e *RFC7807ErrorMapper) formatError(err error, mapping RFC7807ErrorMapping) RFC7807Error {
	return RFC7807Error{
		TypeURI: e.formatURI(mapping.Title),
		Title:   mapping.Title,
		Status:  mapping.Status,
		Detail:  err.Error(),
	}
}

func (e *RFC7807ErrorMapper) formatURI(title string) string {
	return strings.ReplaceAll(strings.ToLower(e.uriPrefix+title), " ", "")
}
