package errors

import (
	"errors"
	"log"
	"reflect"
)

// RFC7807Mapper is an implementation of ErrorMapper that, given a Go Error,
// returns an HTTP status code and RFC7807Error.
type RFC7807Mapper struct {
	defaultError RFC7807Error
	errorMap     map[reflect.Type]RFC7807Error
}

var _ ErrorMapper = (*RFC7807Mapper)(nil)

// RFC7807Error represents an error in JSON RFC7807 format.
type RFC7807Error struct {
	TypeURI string `json:"type"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Detail  string `json:"detail"`
}

// NewRFC7807ErrorMapper creates a new RFC7807ErrorMapper that returns the given
// default error if no matching error was found.
func NewRFC7807ErrorMapper(defaultError RFC7807Error) RFC7807Mapper {
	errorMap := make(map[reflect.Type]RFC7807Error)
	return RFC7807Mapper{
		defaultError: defaultError,
		errorMap:     errorMap,
	}
}

// MapError will, given a Go Error, return an appropriate
// HTTP status code and RFC7807Error according to registered mapping rules.
// See RegisterError.
//
// Matching is performed by comparing the TypeOf the error against the TypeOf the errors in
// registered mappings. If not matched the error is unwrapped using Unwrap and wrapped errors compared.
//
// If the error is not matched to any registered mapping rule then the defaultError and its status is
// returned.
// Unmatched errors will always also be logged using err.error().
func (m *RFC7807Mapper) MapError(err error) (int, any) {
	if match, unwrappedErr, ok := m.matchError(err); ok {
		// return specific error
		return match.Status, m.formatError(unwrappedErr, match)
	} else {
		// return default error
		log.Println("Error: ", err.Error())
		return m.defaultError.Status, m.formatError(errors.New(""), m.defaultError)
	}
}

// RegisterError allows a rule to be added to his mapper that describes how a Go error
// should be handled.
func (m *RFC7807Mapper) RegisterError(err error, mapping RFC7807Error) {
	errorType := reflect.TypeOf(err).Elem()
	m.errorMap[errorType] = mapping
}

// matchError compares type of error (and any wrapped errors) with mappings.
// Returns mapping and actual (possibly unwrapped) error matched of false if no match found.
func (m *RFC7807Mapper) matchError(err error) (RFC7807Error, error, bool) {
	for e := err; e != nil; {
		errorType := reflect.TypeOf(e)
		match, ok := m.errorMap[errorType]
		if ok {
			return match, e, ok
		}
		e = errors.Unwrap(e)
	}
	return RFC7807Error{}, nil, false
}

// formatError formats a Go error into an RFC7807Error according to a match.
func (m *RFC7807Mapper) formatError(err error, match RFC7807Error) RFC7807Error {
	detail := err.Error()
	if match.Detail != "" {
		detail = match.Detail
	}
	return RFC7807Error{
		TypeURI: match.TypeURI,
		Title:   match.Title,
		Status:  match.Status,
		Detail:  detail,
	}
}
