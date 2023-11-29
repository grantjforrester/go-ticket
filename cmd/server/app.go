package main

import (
	_ "github.com/lib/pq"

	"github.com/grantjforrester/go-ticket/internal/adapter/api"
	"github.com/grantjforrester/go-ticket/internal/adapter/repository"
	"github.com/grantjforrester/go-ticket/internal/service"
	"github.com/grantjforrester/go-ticket/pkg/authz"
	"github.com/grantjforrester/go-ticket/pkg/config"
	"github.com/grantjforrester/go-ticket/pkg/media"
)

type App interface {
	Start()
	Stop()
}

func NewApp(config config.Provider) App {
	// secondary adapters
	connectionPool := repository.NewSQLConnectionPool(config)
	repository := repository.NewSQLTicketRepository(connectionPool)

	// services
	authorizer := authz.AlwaysAuthorize{}
	service := service.NewTicketService(repository, authorizer)

	// primary adapters
	mediaHandler := media.JSONHandler{ErrorMap: api.NewErrorMapper()}
	api := api.NewAPI(config, api.Services{Ticket: service}, mediaHandler)

	return api
}
