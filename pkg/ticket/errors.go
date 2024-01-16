package ticket

// TicketError is returned when an invalid ticket is encountered.
type TicketError struct {
	Message string
}

func (ve TicketError) Error() string {
	return ve.Message
}
