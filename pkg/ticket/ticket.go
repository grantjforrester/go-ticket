package ticket

import (
	"fmt"
	"strings"
)

// Ticket represents a reminder of work to be done in a typical ITSM.
type Ticket struct {

	// Summary is a brief single sentence describing the work to be done.
	Summary string `json:"summary"`

	// Description is a full and detailed description of the work to be done.
	Description string `json:"description"`

	// Status describes whether the work has been completed.
	Status string `json:"status"`
}

// Validate checks the mandatory ticket properties are valid. Returns error if validation fails.
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

// TicketWithMetadata merges the types Ticket and Metadata.
type TicketWithMetadata struct {

	// Metadata identifies the ticket.
	Metadata

	// Ticket holds the ticket details.
	Ticket
}

// Validate checks the ticket and metadata properties are valid. Returns error if validation fails.
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
