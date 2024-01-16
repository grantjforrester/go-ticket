package ticket

import (
	"fmt"
	"strings"
)

// A Ticket represents a ticket (a reminder of work to be done) in a typical ITSM.
type Ticket struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Validates the ticket properties. Returns error if validation fails.
func (t Ticket) Validate() error {
	errs := []string{}

	if t.Summary == "" {
		errs = append(errs, "missing field: summary")
	}

	if t.Status == "" {
		errs = append(errs, "missing field: status")
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ","))
	}

	return nil
}

// A TicketWithMetadata merges the types Ticket and Metadata.
type TicketWithMetadata struct {
	Metadata
	Ticket
}

// Validates the ticket with metadata properties. Returns error if validation fails.
func (t TicketWithMetadata) Validate() error {
	errs := []string{}

	if err := t.Metadata.Validate(); err != nil {
		errs = append(errs, err.Error())
	}

	if err := t.Ticket.Validate(); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ","))
	}

	return nil
}
