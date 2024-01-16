package errors

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

// RFC7807Error represents an error in JSON RFC7807 format.
type RFC7807Error struct {
	TypeURI string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
}

// RFC7807Mapper is an implementation of ErrorMapper that, given a Go Error,
// reformats to a status code and RFC7807Error.
type RFC7807Mapper struct {
	uriPrefix    string
	defaultError RFC7807Mapping
	errorMap     map[string]RFC7807Mapping
}

var _ ErrorMapper = (*RFC7807Mapper)(nil)

// RFC7807Mapping describes the details of an RFC7807 error to be returned.
type RFC7807Mapping struct {
	Status int
	Title  string
}

// NewRFC7807Mapper creates a new RFC7807ErrorMapper.
func NewRFC7807Mapper(uriPrefix string, defaultError RFC7807Mapping) RFC7807Mapper {
	errorMap := make(map[string]RFC7807Mapping)
	return RFC7807Mapper{
		uriPrefix:    uriPrefix,
		defaultError: defaultError,
		errorMap:     errorMap,
	}
}

// MapError is an implementation of ErrorMapper.MapError().  Given a Go Error, an appropriate
// status code and RFC7807Error according to registered mappings. See RegisterError.
//
// Matching is performed by comparing the TypeOf the error against the TypeOf the errors in
// registered mappings. If not matched the error is unwrapped using Unwrap and wrapped errors compared.
//
// If the error is not matched to any registered mapping then the defaultError and its status is
// returned and with the Detail field deliberately left blank.
// Unmatched errors will be logged using err.error().
func (m *RFC7807Mapper) MapError(err error) (int, any) {
	if errorMapping, matchedError, ok := m.matchError(err); ok {
		// return specific error
		return errorMapping.Status, m.formatError(matchedError, errorMapping)
	} else {
		// return default error
		log.Println("Error: ", err.Error())
		return m.defaultError.Status, m.formatError(errors.New(""), m.defaultError)
	}
}

// RegisterError allows a rule to be added to mapper that describes how a matching Go error
// should be handled.
func (m *RFC7807Mapper) RegisterError(err error, mapping RFC7807Mapping) {
	errorType := reflect.TypeOf(err).Elem().Name()
	m.errorMap[errorType] = mapping
}

// matchError compares type of error (and any wrapped errors) with mappings.
// Returns mapping and actual (possibly unwrapped) error matched of false if no match found.
func (m *RFC7807Mapper) matchError(err error) (RFC7807Mapping, error, bool) {
	for e := err; e != nil; {
		errorType := reflect.TypeOf(e).Name()
		errorMapping, ok := m.errorMap[errorType]
		if ok {
			return errorMapping, e, ok
		}
		e = errors.Unwrap(e)
	}
	return RFC7807Mapping{}, nil, false
}

// formatError formats a Go error into an RFC7807Error.
func (m *RFC7807Mapper) formatError(err error, mapping RFC7807Mapping) RFC7807Error {
	return RFC7807Error{
		TypeURI: m.formatURI(mapping.Title),
		Title:   mapping.Title,
		Status:  mapping.Status,
		Detail:  err.Error(),
	}
}

// format URI creates error URI from prefix + title with spaces removed.
func (m *RFC7807Mapper) formatURI(title string) string {
	return strings.ReplaceAll(strings.ToLower(m.uriPrefix+title), " ", "")
}
