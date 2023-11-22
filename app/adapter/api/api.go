package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/grantjforrester/go-ticket/pkg/config"
	"github.com/grantjforrester/go-ticket/pkg/media"

	"github.com/grantjforrester/go-ticket/app/service"
)

type Api struct {
	port int
	server *http.Server
	router *mux.Router
	service service.TicketService
	mediaHandler media.Handler
}

func NewApi(config config.Provider, svc service.TicketService, mh media.Handler) Api {
	p := config.GetInt("api_port")
	r := mux.NewRouter()
	srv := &http.Server{Addr: fmt.Sprintf(":%d",p), Handler: r}
	api := Api{port: p, server: srv, router: r, service: svc, mediaHandler: mh}

	api.registerTickets()

	r.HandleFunc("foo", func(w http.ResponseWriter, r *http.Request){ w.WriteHeader(200)})

	r.NotFoundHandler = http.HandlerFunc(api.PathNotFound)

	return api 
}

func (a Api) PathNotFound (w http.ResponseWriter, r *http.Request) {
	err := PathNotFoundError{Message: "resource not found: " + r.RequestURI}
	a.mediaHandler.WriteError(w, &err)
}

func (a Api) Start() {
	go func() {
		log.Println("API started on port", a.port)
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicln(err)
		}
	}()
}

func (a Api) Stop() {
	log.Println("Stopping API on port", a.port)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
	a.server.Shutdown(ctx)
	log.Println("API stopped")
}