package ticket

import (
	"fmt"
	"strings"
)

// Metadata holds information common to all domain entities.
type Metadata struct {
	// ID is a unique identifier for a domain entity.
	ID string `json:"id"`

	// Version is an identifier that changes as the domain entity's properties change.
	Version string `json:"version"`
}

// Validate checks metadata properties are valid. Returns error if validation fails.
func (m Metadata) Validate() error {
	errs := []string{}

	if m.ID == "" {
		errs = append(errs, "missing field: id")
	}

	if m.Version == "" {
		errs = append(errs, "missing field: version")
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, ","))
	}

	return nil
}
