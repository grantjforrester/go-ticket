package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grantjforrester/go-ticket/internal/service"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/repository"
	"github.com/grantjforrester/go-ticket/pkg/ticket"
)

type SQLTicketRepository struct {
	connectionPool *sql.DB
}

var _ repository.Repository[ticket.TicketWithMetadata] = (*SQLTicketRepository)(nil)

func NewSQLTicketRepository(pool *sql.DB) SQLTicketRepository {
	return SQLTicketRepository{connectionPool: pool}
}

func (s SQLTicketRepository) Create(tx repository.Tx, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	var uuid string

	err := ptx.QueryRow(`INSERT INTO tickets (id, version, summary, description, status)
			  				VALUES (uuid_generate_v4(), 0, $1, $2, $3)
			  				RETURNING id`, t.Summary, t.Description, t.Status).Scan(&uuid)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("create ticket failed: %w", err)
	}

	createdTicket, err := s.Read(tx, uuid)

	return createdTicket, err
}

func (s SQLTicketRepository) Read(tx repository.Tx, ticketID string) (ticket.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	row := ptx.QueryRow(`SELECT id, version, summary, description, status 
						 	FROM tickets
							WHERE id = $1`, ticketID)

	t := ticket.TicketWithMetadata{}
	switch err := row.Scan(&t.ID, &t.Version, &t.Summary, &t.Description, &t.Status); err {
	case nil:
		return t, nil
	case sql.ErrNoRows:
		return ticket.TicketWithMetadata{}, &service.NotFoundError{Message: fmt.Sprintf("no ticket with id %s found", ticketID)}
	default:
		return ticket.TicketWithMetadata{}, err
	}
}

func (s SQLTicketRepository) Update(tx repository.Tx, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	_, err := s.Read(tx, t.Metadata.ID)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (1): %w", err)
	}

	res, err := ptx.Exec(`UPDATE tickets
							SET summary = $3, description = $4, status = $5
							WHERE id = $1
							AND version = $2`,
		t.ID, t.Version, t.Summary, t.Description, t.Status)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (2): %w", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update ticket failed (3): no rows updated")
	}
	if rowCount != 1 {
		return ticket.TicketWithMetadata{}, &service.ConflictError{Message: "update ticket failed (4): version conflict"}
	}

	return s.Read(tx, t.Metadata.ID)
}

func (s SQLTicketRepository) Delete(tx repository.Tx, ticketID string) error {
	ptx := tx.(*sql.Tx)
	_, err := ptx.Exec(`DELETE FROM tickets WHERE id = $1`, ticketID)
	if err != nil {
		return fmt.Errorf("delete ticket failed: %w", err)
	}

	return nil
}

func (s SQLTicketRepository) Query(tx repository.Tx, query repository.Query) (collection.Page[ticket.TicketWithMetadata], error) {
	ptx := tx.(*sql.Tx)
	qspec := query.(collection.QuerySpec)
	results := []ticket.TicketWithMetadata{}

	rows, err := ptx.Query(`SELECT id, version, summary, description, status 
							FROM tickets
							LIMIT $1 OFFSET $2`, query.Size, (query.Page-1)*query.Size)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{},
			fmt.Errorf("query tickets failed (1): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := ticket.TicketWithMetadata{}
		err := rows.Scan(&t.ID, &t.Version, &t.Summary, &t.Description, &t.Status)
		if err != nil {
			return collection.Page[ticket.TicketWithMetadata]{},
				fmt.Errorf("query tickets failed (2): %w", err)
		}
		results = append(results, t)
	}

	size := uint64(len(results))
	page := uint64(0)
	if size > 0 {
		page = qspec.Page
	}
	return collection.Page[ticket.TicketWithMetadata]{
		Results: results,
		Page:    page,
		Size:    size,
	}, nil
}

func (s SQLTicketRepository) StartTx(ctx context.Context, readOnly bool) (repository.Tx, error) {
	opts := sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: readOnly}
	tx, err := s.connectionPool.BeginTx(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("start tx failed: %w", err)
	}
	return tx, nil
}
