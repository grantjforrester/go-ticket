package collection

import (
	"fmt"
)

type QuerySpec struct {
	Filters []FilterSpec
	Sorts   []SortSpec
	Page    uint64 // 0 is not set
	Size    uint64 // 0 is not set
}

type FilterSpec struct {
	Field    string
	Operator Operator
	Value    any
}

type Operator string

type SortSpec struct {
	Field     string
	Direction Direction
}

type Direction string

type Page[T any] struct {
	Results []T    `json:"results"`
	Page    uint64 `json:"page"`
	Size    uint64 `json:"size"`
}

type FieldCapability struct {
	Filter bool
	Sort   bool
}

// Validates the query against a set of field capabilities.
func (q QuerySpec) Validate(fieldCapabilities map[string]FieldCapability) error {
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
