package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/grantjforrester/go-ticket/pkg/authz"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/repository"
	"github.com/grantjforrester/go-ticket/pkg/ticket"

	"github.com/google/uuid"
)

var QueryDefaults = struct {
	Page uint64
	Size uint64
}{
	Page: uint64(1),
	Size: uint64(100),
}

var ticketCapabilities = map[string]collection.FieldCapability{
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

	applyDefaults(&query)
	if err := query.Validate(ticketCapabilities); err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, err
	}

	tx, err := svc.repository.StartTx(context, true)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("could not start tx: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	tickets, err := svc.repository.Query(tx, query)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("query ticket from repository failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("cound not commit tx: %w", err)
	}

	return tickets, nil
}

func (svc TicketService) ReadTicket(context context.Context, ticketID string) (ticket.TicketWithMetadata, error) {
	if err := svc.authorizer.IsAuthorized(context, "ReadTicket"); err != nil {
		return ticket.TicketWithMetadata{}, err
	}

	_, err := uuid.Parse(ticketID)
	if err != nil {
		return ticket.TicketWithMetadata{}, RequestError{Message: fmt.Sprintf("invalid ticket id: %s", ticketID)}
	}

	tx, err := svc.repository.StartTx(context, true)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("could not start tx: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	t, err := svc.repository.Read(tx, ticketID)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("read ticket from repository failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("cound not commit tx: %w", err)
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
		return ticket.TicketWithMetadata{}, fmt.Errorf("could not start tx: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	newTicket, err := svc.repository.Create(tx, t)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("create ticket in repository failed: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("cound not commit tx: %w", err)
	}

	return newTicket, nil
}

func (svc TicketService) UpdateTicket(context context.Context, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	if err := svc.authorizer.IsAuthorized(context, "UpdateTicket"); err != nil {
		return ticket.TicketWithMetadata{}, err
	}

	_, err := uuid.Parse(t.Metadata.ID)
	if err != nil {
		return ticket.TicketWithMetadata{}, RequestError{Message: fmt.Sprintf("invalid ticket id: %s", t.Metadata.ID)}
	}

	if err := t.Validate(); err != nil {
		return ticket.TicketWithMetadata{}, RequestError{Message: err.Error()}
	}

	tx, err := svc.repository.StartTx(context, false)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("could not start tx: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	updatedTicket, err := svc.repository.Update(tx, t)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket in repository failed: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("cound not commit tx: %w", err)
	}

	return updatedTicket, nil
}

func (svc TicketService) DeleteTicket(context context.Context, ticketID string) error {
	if err := svc.authorizer.IsAuthorized(context, "DeleteTicket"); err != nil {
		return err
	}

	_, err := uuid.Parse(ticketID)
	if err != nil {
		return RequestError{Message: fmt.Sprintf("invalid ticket id: %s", ticketID)}
	}

	tx, err := svc.repository.StartTx(context, false)
	if err != nil {
		return fmt.Errorf("could not start tx: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Rollback())
	}()

	err = svc.repository.Delete(tx, ticketID)
	if err != nil {
		return fmt.Errorf("delete ticket from repository: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("cound not commit tx: %w", err)
	}

	return nil
}

func applyDefaults(query *collection.QuerySpec) {
	if query.Page == 0 {
		query.Page = QueryDefaults.Page
	}

	if query.Size == 0 {
		query.Size = QueryDefaults.Size
	}
}
