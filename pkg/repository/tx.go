package repository

type Tx interface {
	Rollback() error
	Commit() error
}
