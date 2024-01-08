package ticket

import (
	"errors"
	"strings"
)

type Ticket struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (a Ticket) Validate() error {
	errs := []string{}

	if a.Summary == "" {
		errs = append(errs, "invalid ticket: missing field: summary")
	}

	if a.Status == "" {
		errs = append(errs, "invalid ticket: missing field: status")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ","))
	}

	return nil
}

type TicketWithMetadata struct {
	Metadata
	Ticket
}
