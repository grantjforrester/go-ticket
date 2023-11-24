package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grantjforrester/go-ticket/app/model"
	"github.com/grantjforrester/go-ticket/app/service"
	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/repository"
)

type SqlTicketRepository struct {
	connectionPool *sql.DB
}

func NewSqlTicketRepository(pool *sql.DB) SqlTicketRepository {
	return SqlTicketRepository{connectionPool: pool}
}

func (s SqlTicketRepository) Create(tx repository.Tx, ticket model.TicketWithMetadata) (model.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	var uuid string

	err := ptx.QueryRow(`INSERT INTO tickets (id, version, summary, description)
			  				VALUES (uuid_generate_v4(), 0, $1, $2)
			  				RETURNING id`, ticket.Summary, ticket.Description).Scan(&uuid)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("create ticket failed: %w", err)
	}

	createdTicket, err := s.Read(tx, uuid)

	return createdTicket, err
}

func (s SqlTicketRepository) Read(tx repository.Tx, ticketId string) (model.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	row := ptx.QueryRow(`SELECT id, version, summary, description 
						 	FROM tickets
							WHERE id = $1`, ticketId)

	ticket := model.TicketWithMetadata{}
	switch err := row.Scan(&ticket.Id, &ticket.Version, &ticket.Summary, &ticket.Description); err {
	case nil:
		return ticket, nil
	case sql.ErrNoRows:
		return model.TicketWithMetadata{}, &service.NotFoundError{Message: fmt.Sprintf("no ticket with id %s found", ticketId)}
	default:
		return model.TicketWithMetadata{}, err
	}
}

func (s SqlTicketRepository) Update(tx repository.Tx, ticket model.TicketWithMetadata) (model.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	_, err := s.Read(tx, ticket.Metadata.Id)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("update ticket failed (1): %w", err)
	}

	res, err := ptx.Exec(`UPDATE tickets
							SET summary = $3, description = $4
							WHERE id = $1
							AND version = $2`,
		ticket.Id, ticket.Version, ticket.Summary, ticket.Description)
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("update ticket failed (2): %w", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return model.TicketWithMetadata{}, fmt.Errorf("update ticket failed (3): no rows updated")
	}
	if rowCount != 1 {
		return model.TicketWithMetadata{}, &service.ConflictError{Message: "update ticket failed (4): version conflict"}
	}

	return s.Read(tx, ticket.Metadata.Id)
}

func (s SqlTicketRepository) Delete(tx repository.Tx, ticketId string) error {
	ptx := tx.(*sql.Tx)
	_, err := ptx.Exec(`DELETE FROM tickets WHERE id = $1`, ticketId)
	if err != nil {
		return fmt.Errorf("delete ticket failed: %w", err)
	}

	return nil
}

func (s SqlTicketRepository) Query(tx repository.Tx, query collection.Query) (collection.Page[model.TicketWithMetadata], error) {
	ptx := tx.(*sql.Tx)
	results := []model.TicketWithMetadata{}

	rows, err := ptx.Query(`SELECT id, version, summary, description FROM tickets`)
	if err != nil {
		return collection.Page[model.TicketWithMetadata]{},
			fmt.Errorf("query tickets failed (1): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		ticket := model.TicketWithMetadata{}
		err := rows.Scan(&ticket.Id, &ticket.Version, &ticket.Summary, &ticket.Description)
		if err != nil {
			return collection.Page[model.TicketWithMetadata]{},
				fmt.Errorf("query tickets failed (2): %w", err)
		}
		results = append(results, ticket)
	}

	sz := uint64(len(results))
	pg := uint64(0)
	if sz > 0 {
		pg = query.Page
	}
	return collection.Page[model.TicketWithMetadata]{
		Results: results,
		Page:    pg,
		Size:    sz,
	}, nil
}

func (s SqlTicketRepository) StartTx(ctx context.Context, readOnly bool) (repository.Tx, error) {
	opts := sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: readOnly}
	tx, err := s.connectionPool.BeginTx(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("start tx failed: %w", err)
	}
	return tx, nil
}
