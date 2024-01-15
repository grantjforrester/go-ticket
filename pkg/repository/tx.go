package repository

// Tx is an atomic database transaction.  All operations
// in the same transaction will all pass or all fail.
type Tx interface {
	// Commit commits the transaction.  Returns error on failure.
	Commit() error

	// Rollback aborts the transaction. Returns error on failure
	Rollback() error
}
