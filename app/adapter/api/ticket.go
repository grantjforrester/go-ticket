package api

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"

	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/ticket"
)

func (api *API) registerTicketRoutes(router *mux.Router) {
	log.Println("registering tickets")
	router.HandleFunc("/tickets", api.queryTickets).Methods("GET")
	router.HandleFunc("/tickets", api.createTicket).Methods("POST")
	router.HandleFunc("/tickets/{key}", api.readTicket).Methods("GET")
	router.HandleFunc("/tickets/{key}", api.updateTicket).Methods("PUT")
	router.HandleFunc("/tickets/{key}", api.deleteTicket).Methods("DELETE")
}

func (api *API) queryTickets(resp http.ResponseWriter, req *http.Request) {
	urlQuery, _ := url.ParseQuery(req.URL.RawQuery)
	querySpec, err := collection.ParseQuery(urlQuery)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	tickets, err := api.services.Ticket.QueryTickets(req.Context(), querySpec)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	api.mediaHandler.WriteResponse(resp, http.StatusOK, tickets)
}

func (api *API) readTicket(resp http.ResponseWriter, req *http.Request) {
	ticketID := path.Base(req.URL.Path)

	ticket, err := api.services.Ticket.ReadTicket(req.Context(), ticketID)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	api.mediaHandler.WriteResponse(resp, http.StatusOK, ticket)
}

func (api *API) createTicket(resp http.ResponseWriter, req *http.Request) {
	inTicket := ticket.TicketWithMetadata{}
	err := api.mediaHandler.ReadResource(req, &inTicket)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	createdTicket, err := api.services.Ticket.CreateTicket(req.Context(), inTicket)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	api.mediaHandler.WriteResponse(resp, http.StatusCreated, createdTicket)
}

func (api *API) updateTicket(resp http.ResponseWriter, req *http.Request) {
	ticketID := path.Base(req.URL.Path)
	inTicket := ticket.TicketWithMetadata{}
	err := api.mediaHandler.ReadResource(req, &inTicket)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}
	inTicket.ID = ticketID

	updatedTicket, err := api.services.Ticket.UpdateTicket(req.Context(), inTicket)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	api.mediaHandler.WriteResponse(resp, http.StatusOK, updatedTicket)
}

func (api *API) deleteTicket(resp http.ResponseWriter, req *http.Request) {
	ticketID := path.Base(req.URL.Path)

	err := api.services.Ticket.DeleteAlert(req.Context(), ticketID)
	if err != nil {
		api.mediaHandler.WriteError(resp, err)
		return
	}

	api.mediaHandler.WriteResponse(resp, http.StatusNoContent, nil)
}
