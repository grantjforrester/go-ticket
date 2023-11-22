package app

import (
	_ "github.com/lib/pq"

	"github.com/grantjforrester/go-ticket/app/adapter/api"
	"github.com/grantjforrester/go-ticket/app/adapter/repository"
	"github.com/grantjforrester/go-ticket/app/service"
	"github.com/grantjforrester/go-ticket/pkg/config"
	"github.com/grantjforrester/go-ticket/pkg/media"
)

type App interface {
	Start()
	Stop()
}

func NewApp(config config.Provider) App {
	// secondary adapters
	connectionPool := repository.NewSqlConnectionPool(config)
	repository := repository.NewSqlTicketRepository(connectionPool)

	// services
	service := service.NewTicketService(repository)
	
	// primary adapters
	mediaHandler := media.JsonHandler{ ErrorMap: api.NewErrorMapper() }
	
	api := api.NewApi(config, service, mediaHandler)
	
	return api
}