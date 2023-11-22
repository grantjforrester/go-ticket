package service

import (
	"context"
	"fmt"

	"github.com/grantjforrester/go-ticket/app/model"
	"github.com/grantjforrester/go-ticket/pkg/repository"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/authz"
)

type TicketRepository repository.Repository[model.TicketWithMetadata]

type TicketService struct {
	authorizer authz.Authorizer
	repository TicketRepository
}

func NewTicketService(r TicketRepository) TicketService {
	return TicketService{repository: r}
}

func (as TicketService) QueryTickets(context context.Context, query collection.Query) (collection.Page[model.TicketWithMetadata], error) {
	if err := as.authorizer.IsAuthorized(context, "QueryTickets"); err != nil {
		return collection.Page[model.TicketWithMetadata]{}, err
	}

	tx, err := as.repository.StartTx(context, true)
	if err != nil {
		return collection.Page[model.TicketWithMetadata]{}, fmt.Errorf("alertservice: get alerts failed (1): %w", err)
	}

	alerts, err := as.repository.Query(tx, query)

	if err != nil {
		tx.Rollback()
		return collection.Page[model.TicketWithMetadata]{}, fmt.Errorf("ticketservice: query tickets failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return collection.Page[model.TicketWithMetadata]{}, fmt.Errorf("ticketservice: query tickets failed (3): %w", err)
	}

	return alerts, nil
}

func (as TicketService) ReadTicket(context context.Context, alertId string) (model.TicketWithMetadata, error) {
	if err := as.authorizer.IsAuthorized(context, "ReadTicket"); err != nil {
		return model.TicketWithMetadata{}, err
	}
	
	tx, err := as.repository.StartTx(context, true)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: read ticket failed (1): %w", err)
	}

	ticket, err := as.repository.Read(tx, alertId)

	if err != nil {
		tx.Rollback()
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: read ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: read ticket failed (3): %w", err)
	}

	return ticket, nil
}

func (as TicketService) CreateTicket(context context.Context, ticket model.TicketWithMetadata) (model.TicketWithMetadata, error) {
	if err := as.authorizer.IsAuthorized(context, "CreateTicket"); err != nil {
		return model.TicketWithMetadata{}, err
	}
	
	if err := ticket.Validate(); err != nil {
		return model.TicketWithMetadata{}, &RequestError{Message: err.Error()}
	}

	tx, err := as.repository.StartTx(context, false)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: create ticket failed (1): %w", err)
	}

	newTicket, err := as.repository.Create(tx, ticket)

	if err != nil {
		tx.Rollback()
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: create ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: create ticket failed (3): %w", err)
	}

	return newTicket, nil
}

func (as TicketService) UpdateTicket(context context.Context, ticket model.TicketWithMetadata) (model.TicketWithMetadata, error) {
	if err := as.authorizer.IsAuthorized(context, "UpdateTicket"); err != nil {
		return model.TicketWithMetadata{}, err
	}
	
	if err := ticket.Validate(); err != nil {
		return model.TicketWithMetadata{}, &RequestError{Message: err.Error()}
	}

	tx, err := as.repository.StartTx(context, false)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: update ticket failed (1): %w", err)
	}

	updatedTicket, err := as.repository.Update(tx, ticket)

	if err != nil {
		tx.Rollback()
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: update ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("ticketservice: update ticket failed (3): %w", err)
	}

	return updatedTicket, nil
}

func (as TicketService) DeleteAlert(context context.Context, alertId string) error {
	if err := as.authorizer.IsAuthorized(context, "DeleteTicket"); err != nil {
		return err
	}

	tx, err := as.repository.StartTx(context, false)
	if err != nil {
		return fmt.Errorf("ticketservice: delete ticket failed (1): %w", err)
	}

	err = as.repository.Delete(tx, alertId)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ticketservice: delete ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ticketservice: delete ticket failed (3): %w", err)
	}

	return nil
}