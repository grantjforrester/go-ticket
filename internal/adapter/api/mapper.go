package api

import (
	"github.com/grantjforrester/go-ticket/internal/adapter/repository"
	"github.com/grantjforrester/go-ticket/internal/service"
	"github.com/grantjforrester/go-ticket/pkg/media/errors"
)

func NewErrorMapper() errors.ErrorMapper {
	errorMapper := errors.NewRFC7807Mapper("induction:go:err:", errors.RFC7807Mapping{
		Status: 500,
		Title:  "Internal Server Error",
	})
	errorMapper.RegisterError((*PathNotFoundError)(nil), errors.RFC7807Mapping{
		Status: 404,
		Title:  "Not Found",
	})
	errorMapper.RegisterError((*service.RequestError)(nil), errors.RFC7807Mapping{
		Status: 400,
		Title:  "Bad Request",
	})
	errorMapper.RegisterError((*repository.NotFoundError)(nil), errors.RFC7807Mapping{
		Status: 404,
		Title:  "Not Found",
	})
	errorMapper.RegisterError((*repository.ConflictError)(nil), errors.RFC7807Mapping{
		Status: 409,
		Title:  "Conflict",
	})

	return &errorMapper
}
