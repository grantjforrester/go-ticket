package api

import (
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/grantjforrester/go-ticket/app/model"
	"github.com/grantjforrester/go-ticket/pkg/collection"
)

func (a *Api) registerTickets() {
	log.Println("registering tickets")
	a.router.HandleFunc("/tickets", a.queryTickets).Methods("GET")
	a.router.HandleFunc("/tickets", a.createTicket).Methods("POST")
	a.router.HandleFunc("/tickets/{key}", a.readTicket).Methods("GET")
	a.router.HandleFunc("/tickets/{key}", a.updateTicket).Methods("PUT")
	a.router.HandleFunc("/tickets/{key}", a.deleteTicket).Methods("DELETE")
}

func (a *Api) queryTickets(w http.ResponseWriter, r *http.Request) {
	rqlQuery, _ := url.ParseQuery(r.URL.RawQuery)
	query, err := collection.ParseQuery(rqlQuery)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	alerts, err := a.service.QueryTickets(r.Context(), query)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	a.mediaHandler.WriteResponse(w, http.StatusOK, alerts)
}

func (a *Api) readTicket(w http.ResponseWriter, r *http.Request) {
	alertId := path.Base(r.URL.Path)

	alert, err := a.service.ReadTicket(r.Context(), alertId)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	a.mediaHandler.WriteResponse(w, http.StatusOK, alert)
}

func (a *Api) createTicket(w http.ResponseWriter, r *http.Request) {
	inAlert := model.TicketWithMetadata{}
	err := a.mediaHandler.ReadResource(r, &inAlert)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	createdAlert, err := a.service.CreateTicket(r.Context(), inAlert)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	a.mediaHandler.WriteResponse(w, http.StatusCreated, createdAlert)
}

func (a *Api) updateTicket(w http.ResponseWriter, r *http.Request) {
	alertId := path.Base(r.URL.Path)
	inAlert := model.TicketWithMetadata{}
	err := a.mediaHandler.ReadResource(r, &inAlert)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}
	inAlert.Id = alertId

	updatedAlert, err := a.service.UpdateTicket(r.Context(), inAlert)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	a.mediaHandler.WriteResponse(w, http.StatusOK, updatedAlert)
}

func (a *Api) deleteTicket(w http.ResponseWriter, r *http.Request) {
	alertId := path.Base(r.URL.Path)

	err := a.service.DeleteAlert(r.Context(), alertId)
	if err != nil {
		a.mediaHandler.WriteError(w, err)
		return
	}

	a.mediaHandler.WriteResponse(w, http.StatusNoContent, nil)
}
