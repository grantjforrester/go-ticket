package model

import (
	"errors"
	"strings"
)

type Ticket struct {
	Summary string		`json:"summary" query:"filter,sort"`
	Description string	`json:"description" query:"-"`
}

func(a Ticket) Validate() error {
	errs := []string{}
	
	if a.Summary == "" {
		errs = append(errs, "invalid alert: missing field: summary")
	}
	if a.Description == "" {
		errs = append(errs, "invalid alert: missing field: description")
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
