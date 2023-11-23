package collection

import (
	"fmt"
)

type Query struct {
	Filters []FilterSpec
	Sorts	[]SortSpec
	Page	uint64
	Size	uint64
}

type FieldCapability struct {
	Filter  bool
	Sort	bool
}

func (q Query) Validate(fieldCapabilities map[string]FieldCapability) error {
	for _, field := range q.Filters {
		if f, ok := fieldCapabilities[field.Field]; !ok || !f.Filter {
			return fmt.Errorf("invalid filter: %s", field.Field)
		}
	}

	for _, field := range q.Sorts {
		if f, ok := fieldCapabilities[field.Field]; !ok || !f.Sort {
			return fmt.Errorf("invalid sort: %s", field.Field)
		}
	}
	
	return nil
}