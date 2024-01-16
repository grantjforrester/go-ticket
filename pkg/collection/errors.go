package collection

// QueryError is returned when a query is invalid.
type QueryError struct {
	Message string
}

func (qe QueryError) Error() string {
	return qe.Message
}
