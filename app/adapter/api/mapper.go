package api

import (
	"github.com/grantjforrester/go-ticket/pkg/media"

	"github.com/grantjforrester/go-ticket/app/service"
)

func NewErrorMapper() media.ErrorMapper {
	errorMapper := media.NewRFC7807ErrorMapper("induction:go:err:", media.RFC7807ErrorMapping{
		Status: 500,
		Title:  "Internal Server Error",
	})
	errorMapper.RegisterError((*PathNotFoundError)(nil), media.RFC7807ErrorMapping{
		Status: 404,
		Title:  "Not Found",
	})
	errorMapper.RegisterError((*service.RequestError)(nil), media.RFC7807ErrorMapping{
		Status: 400,
		Title:  "Bad Request",
	})
	errorMapper.RegisterError((*service.NotFoundError)(nil), media.RFC7807ErrorMapping{
		Status: 404,
		Title:  "Not Found",
	})
	errorMapper.RegisterError((*service.ConflictError)(nil), media.RFC7807ErrorMapping{
		Status: 409,
		Title:  "Conflict",
	})

	return &errorMapper
}
