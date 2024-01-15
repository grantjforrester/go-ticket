package repository

import (
	"context"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

// Repository describes a common pattern for CRUD operations on persistent entities.
type Repository[T any] interface {

	// Create creates a new entity in the repository using the given transaction.
	// Returns the new entity, or error.
	Create(Tx, T) (T, error)

	// Finds an entity by its unique id using the given transaction.
	// Returns the found entity, or error.
	Read(Tx, string) (T, error)

	// Updates the entity using the given transaction.
	// Returns the updated entity, or error.
	Update(Tx, T) (T, error)

	// Deletes an entity with the unique id using the given transaction.
	Delete(Tx, string) error

	// Finds entities based on the criteria in the query using the given transaction.
	// Returns a page of matching entities, or error.
	Query(Tx, Query) (collection.Page[T], error)

	// Starts a new transaction in the repository. Returns the transaction, or error.
	StartTx(context.Context, bool) (Tx, error)
}
