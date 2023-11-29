package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grantjforrester/go-ticket/pkg/authz"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/repository"
	"github.com/grantjforrester/go-ticket/pkg/ticket"
)

var ticketQueries = map[string]collection.FieldCapability{
	"summary": {Filter: true, Sort: true},
}

type TicketService struct {
	authorizer authz.Authorizer
	repository TicketRepository
}

type TicketRepository repository.Repository[ticket.TicketWithMetadata]

func NewTicketService(r TicketRepository, a authz.Authorizer) TicketService {
	return TicketService{repository: r, authorizer: a}
}

func (svc TicketService) QueryTickets(context context.Context, query collection.QuerySpec) (collection.Page[ticket.TicketWithMetadata], error) {
	if err := svc.authorizer.IsAuthorized(context, "QueryTickets"); err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, err
	}

	if err := query.Validate(ticketQueries); err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, err
	}

	tx, err := svc.repository.StartTx(context, true)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("query tickets failed (1): %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	tickets, err := svc.repository.Query(tx, query)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("query tickets failed (2): %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("query tickets failed (3): %w", err)
	}

	return tickets, nil
}

func (svc TicketService) ReadTicket(context context.Context, ticketID string) (ticket.TicketWithMetadata, error) {
	if err := svc.authorizer.IsAuthorized(context, "ReadTicket"); err != nil {
		return ticket.TicketWithMetadata{}, err
	}

	tx, err := svc.repository.StartTx(context, true)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("read ticket failed (1): %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	t, err := svc.repository.Read(tx, ticketID)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("read ticket failed (2): %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("read ticket failed (3): %w", err)
	}

	return t, nil
}

func (svc TicketService) CreateTicket(context context.Context, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	if err := svc.authorizer.IsAuthorized(context, "CreateTicket"); err != nil {
		return ticket.TicketWithMetadata{}, err
	}

	if err := t.Validate(); err != nil {
		return ticket.TicketWithMetadata{}, &RequestError{Message: err.Error()}
	}

	tx, err := svc.repository.StartTx(context, false)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("create ticket failed (1): %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	newTicket, err := svc.repository.Create(tx, t)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("create ticket failed (2): %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("create ticket failed (3): %w", err)
	}

	return newTicket, nil
}

func (svc TicketService) UpdateTicket(context context.Context, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	if err := svc.authorizer.IsAuthorized(context, "UpdateTicket"); err != nil {
		return ticket.TicketWithMetadata{}, err
	}

	if err := t.Validate(); err != nil {
		return ticket.TicketWithMetadata{}, &RequestError{Message: err.Error()}
	}

	tx, err := svc.repository.StartTx(context, false)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (1): %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	updatedTicket, err := svc.repository.Update(tx, t)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (3): %w", err)
	}

	return updatedTicket, nil
}

func (svc TicketService) DeleteAlert(context context.Context, ticketID string) error {
	if err := svc.authorizer.IsAuthorized(context, "DeleteTicket"); err != nil {
		return err
	}

	tx, err := svc.repository.StartTx(context, false)
	if err != nil {
		return fmt.Errorf("delete ticket failed (1): %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	err = svc.repository.Delete(tx, ticketID)
	if err != nil {
		return fmt.Errorf("delete ticket failed (2): %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("delete ticket failed (3): %w", err)
	}

	return nil
}
