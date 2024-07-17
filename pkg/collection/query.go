package collection

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// QuerySpec describes a common pattern for querying RESTful resource collections.
type QuerySpec struct {
	Filters []FilterExpr
	Sorts   []SortExpr
	Page    uint64 // 0 is not set
	Size    uint64 // 0 is not set
}

// FilterExpr describes a criteria for matching resources in a collection.
type FilterExpr struct {
	Field    string
	Operator Operator
	Value    any
}

// Operator describes a mathemtical comparison to be performed with 1 or more values .
type Operator string

// SortExpr describes an order for returning matching resources in a collection.
type SortExpr struct {
	Field     string
	Direction Direction
}

type Direction string

// Page describes a subset of query results from a collection of resources.
type Page[T any] struct {
	Results []T    `json:"results"`
	Page    uint64 `json:"page"`
	Size    uint64 `json:"size"`
}

// FieldCapability describes a field of a resource and how it may be used in a collection
// query.
type FieldCapability struct {

	// Filter describes whether the field can be used in a filter expression of a query.
	Filter bool

	// FilterOps describes the operators that may be used with this field in the filter.
	// expression of a query.
	FilterOps []Operator

	// Sort describes whether the field can be used in a sort expression.
	Sort bool
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
