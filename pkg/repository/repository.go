package repository

import "context"

import "github.com/grantjforrester/go-ticket/pkg/collection"

type Repository[T any] interface {
	Create(Tx, T) (T, error)
	Read(Tx, string) (T, error)
	Update(Tx, T) (T, error)
	Delete(Tx, string) error

	Query(Tx, collection.Query) (collection.Page[T], error)

	StartTx(context.Context, bool) (Tx, error)
}