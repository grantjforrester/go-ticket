package model

import (
	"errors"
	"strings"
)

type Ticket struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func (a Ticket) Validate() error {
	errs := []string{}

	if a.Summary == "" {
		errs = append(errs, "invalid ticket: missing field: summary")
	}
	if a.Description == "" {
		errs = append(errs, "invalid ticket: missing field: description")
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
