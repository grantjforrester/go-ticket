package ticket

import (
	"strings"
)

// A Ticket represents a ticket (a reminder of work to be done) in a typical ITSM.
type Ticket struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Validates the ticket properties. Returns ValidationError if validation fails.
func (a Ticket) Validate() error {
	errs := []string{}

	if a.Summary == "" {
		errs = append(errs, "invalid ticket: missing field: summary")
	}

	if a.Status == "" {
		errs = append(errs, "invalid ticket: missing field: status")
	}

	if len(errs) > 0 {
		return TicketError{Message: strings.Join(errs, ",")}
	}

	return nil
}

// A TicketWithMetadata merges the types Ticket and Metadata.
type TicketWithMetadata struct {
	Metadata
	Ticket
}
