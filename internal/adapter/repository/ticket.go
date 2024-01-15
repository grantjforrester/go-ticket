package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grantjforrester/go-ticket/pkg/collection"
	sq "github.com/grantjforrester/go-ticket/pkg/collection/sql"
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
		return ticket.TicketWithMetadata{}, fmt.Errorf("insert statement failed: %w", err)
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
		return ticket.TicketWithMetadata{}, NotFoundError{Message: fmt.Sprintf("no ticket with id %s found", ticketID)}
	default:
		return ticket.TicketWithMetadata{}, err
	}
}

func (s SQLTicketRepository) Update(tx repository.Tx, t ticket.TicketWithMetadata) (ticket.TicketWithMetadata, error) {
	ptx := tx.(*sql.Tx)
	_, err := s.Read(tx, t.Metadata.ID)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("read ticket failed: %w", err)
	}

	res, err := ptx.Exec(`UPDATE tickets
							SET summary = $3, description = $4, status = $5
							WHERE id = $1
							AND version = $2`,
		t.ID, t.Version, t.Summary, t.Description, t.Status)
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("update statement failed: %w", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return ticket.TicketWithMetadata{}, fmt.Errorf("count of updated rows failed: %w", err)
	}
	if rowCount != 1 {
		return ticket.TicketWithMetadata{}, ConflictError{Message: "version conflict"}
	}

	return s.Read(tx, t.Metadata.ID)
}

func (s SQLTicketRepository) Delete(tx repository.Tx, ticketID string) error {
	ptx := tx.(*sql.Tx)
	_, err := ptx.Exec(`DELETE FROM tickets WHERE id = $1`, ticketID)
	if err != nil {
		return fmt.Errorf("delete statement failed: %w", err)
	}

	return nil
}

func (s SQLTicketRepository) Query(tx repository.Tx, query repository.Query) (collection.Page[ticket.TicketWithMetadata], error) {
	ptx := tx.(*sql.Tx)
	qspec := query.(collection.QuerySpec)
	results := []ticket.TicketWithMetadata{}
	qry, args, err := sq.SQLQuery{
		Fields: []string{"id", "version", "summary", "description", "status"},
		Table:  "tickets",
		Query:  qspec,
	}.ToSQL()
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{}, fmt.Errorf("building sql query failed): %w", err)
	}

	rows, err := ptx.Query(qry, args...)
	if err != nil {
		return collection.Page[ticket.TicketWithMetadata]{},
			fmt.Errorf("executing query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		t := ticket.TicketWithMetadata{}
		err := rows.Scan(&t.ID, &t.Version, &t.Summary, &t.Description, &t.Status)
		if err != nil {
			return collection.Page[ticket.TicketWithMetadata]{},
				fmt.Errorf("error reading row: %w", err)
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
