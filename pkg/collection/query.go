package collection

import (
	"fmt"
	"golang.org/x/exp/slices"
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
	Filter    bool
	FilterOps []Operator
	Sort      bool
}

// Validates the query against a set of field capabilities.
func (q QuerySpec) Validate(fieldCapabilities map[string]FieldCapability) error {
	for _, filter := range q.Filters {
		if f, ok := fieldCapabilities[filter.Field]; !ok || !f.Filter {
			return QueryError{Message: fmt.Sprintf("invalid filter field: %s", filter.Field)}
		}
		if !slices.Contains(fieldCapabilities[filter.Field].FilterOps, filter.Operator) {
			return QueryError{Message: fmt.Sprintf("invalid filter operator: %s", filter.Operator)}
		}
	}

	for _, sort := range q.Sorts {
		if f, ok := fieldCapabilities[sort.Field]; !ok || !f.Sort {
			return QueryError{Message: fmt.Sprintf("invalid sort field: %s", sort.Field)}
		}
	}

	return nil
}
