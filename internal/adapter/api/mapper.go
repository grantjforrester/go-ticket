package api

import (
	"github.com/grantjforrester/go-ticket/internal/adapter/repository"
	"github.com/grantjforrester/go-ticket/internal/service"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/media"
	"github.com/grantjforrester/go-ticket/pkg/media/errors"
)

func NewErrorMapper() errors.ErrorMapper {
	errorMapper := errors.NewRFC7807ErrorMapper(errors.RFC7807Error{
		TypeURI: "ticket:err:internalservererror",
		Status:  500,
		Title:   "Internal Server Error",
	})
	errorMapper.RegisterError((*PathNotFoundError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:notfound",
		Status:  404,
		Title:   "Not Found",
	})
	errorMapper.RegisterError((*service.RequestError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:badrequest",
		Status:  400,
		Title:   "Bad Request",
	})
	errorMapper.RegisterError((*media.MediaError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:badrequest",
		Status:  400,
		Title:   "Bad Request",
	})
	errorMapper.RegisterError((*collection.QueryError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:badrequest",
		Status:  400,
		Title:   "Bad Request",
	})
	errorMapper.RegisterError((*repository.NotFoundError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:notfound",
		Status:  404,
		Title:   "Not Found",
	})
	errorMapper.RegisterError((*repository.ConflictError)(nil), errors.RFC7807Error{
		TypeURI: "ticket:err:conflict",
		Status:  409,
		Title:   "Conflict",
	})

	return &errorMapper
}
