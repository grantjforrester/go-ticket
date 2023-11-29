package api

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/grantjforrester/go-ticket/pkg/config"
	"github.com/grantjforrester/go-ticket/pkg/media"

	"github.com/grantjforrester/go-ticket/internal/service"
)

type API struct {
	port         int
	server       *http.Server
	services     Services
	mediaHandler media.Handler
}

type Services struct {
	Ticket service.TicketService
}

//go:embed openapi.yml
var openapi []byte

func NewAPI(config config.Provider, svcs Services, mh media.Handler) API {
	prt := config.GetInt("api_port")

	rtr := mux.NewRouter()
	srv := &http.Server{Addr: fmt.Sprintf(":%d", prt), Handler: rtr}
	api := API{port: prt, server: srv, services: svcs, mediaHandler: mh}

	// register standard endpoints
	rtr.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	rtr.HandleFunc("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write(openapi)
		if err != nil {
			mh.WriteError(w, err)
		}
	})

	// register api routes
	v1 := rtr.PathPrefix("/api/v1").Subrouter()
	api.registerTicketRoutes(v1)

	// default not found
	rtr.NotFoundHandler = http.HandlerFunc(api.PathNotFound)

	return api
}

func (api API) PathNotFound(w http.ResponseWriter, r *http.Request) {
	err := PathNotFoundError{Message: "resource not found: " + r.RequestURI}
	api.mediaHandler.WriteError(w, &err)
}

func (api API) Start() {
	go func() {
		log.Println("API started on port", api.port)
		err := api.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicln(err)
		}
	}()
}

func (api API) Stop() {
	log.Println("Stopping API on port", api.port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = api.server.Shutdown(ctx)
	log.Println("API stopped")
}
