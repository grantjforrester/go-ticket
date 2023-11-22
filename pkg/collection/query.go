package collection

import (
	"fmt"
	"regexp"
	"strconv"
)

type Query struct {
	Filters []FilterSpec
	Sorts	[]SortSpec
	Page	PageSpec
}

const (
	ASC_SEQ    = "asc"
	DESC_SEQ   = "desc"
	EQ_SEQ     = "=="
	NE_SEQ     = "!="
	GT_SEQ     = ">"
	LT_SEQ     = "<"
	GE_SEQ     = ">="
	LE_SEQ     = "<="
)

var FIELD_PATTERN = `\w+`
var VALUE_PATTERN = `.+`
var OPERATOR_PATTERN = fmt.Sprintf("%s|%s|%s|%s|%s|%s", EQ_SEQ, NE_SEQ, GT_SEQ, LT_SEQ, GE_SEQ, LE_SEQ)
var FILTER_PATTERN = fmt.Sprintf("(%s)(%s)(%s)", FIELD_PATTERN, OPERATOR_PATTERN, VALUE_PATTERN)
var SORT_PATTERN = fmt.Sprintf("(%s) (%s|%s)", FIELD_PATTERN, ASC_SEQ, DESC_SEQ)

var DIRECTIONS = map[string]Direction{
	ASC_SEQ:  ASC,
	DESC_SEQ: DESC,
}

var OPERATORS = map[string]Operator{
	EQ_SEQ: EQ,
	NE_SEQ: NE,
	GT_SEQ: GT,
	LT_SEQ: LT,
	GE_SEQ: GE,
	LE_SEQ: LE,
}

func ParseQuery(filterClauses []string, sortClauses []string, pageIdx string, pageSize string ) (Query, error) {
	var sort_re = regexp.MustCompile(SORT_PATTERN) 
	var filter_re = regexp.MustCompile(FILTER_PATTERN)

	query := Query{}

	filterSpecs, err := parseFilterClauses(filterClauses, *filter_re)
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Filters = filterSpecs

	sortSpecs, err := parseSortClauses(sortClauses, *sort_re)
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Sorts = sortSpecs

	pageSpec, err := parsePageClause(pageIdx, pageSize)
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Page = pageSpec


	return query, nil
}

func parseFilterClauses(clauses []string, regexp regexp.Regexp) ([]FilterSpec, error) {
	if len(clauses) == 0 {
		return []FilterSpec{}, nil
	}

	filters := []FilterSpec{}
	for _, filter := range clauses {
		if regexp.MatchString(filter) {
			parts := regexp.FindStringSubmatch(filter)
			field := parts[1]
			operator := parts[2]
			value := parts[3]
			filters = append(filters, FilterSpec{ Field: field, Operator: OPERATORS[operator], Value: value})
		} else {
			return nil, fmt.Errorf("invalid filter: %s", filter)
		}
	}
	return filters, nil
}

func parseSortClauses(clauses []string, regexp regexp.Regexp) ([]SortSpec, error) {
	if len(clauses) == 0 {
		return []SortSpec{}, nil
	}

	sorts := []SortSpec{}
	for _, sort := range clauses {
		if regexp.MatchString(sort) {
			parts := regexp.FindStringSubmatch(sort)
			field := parts[1]
			direction := parts[2]
			sorts = append(sorts, SortSpec{Field: field, Direction: DIRECTIONS[direction]})
		} else {
			return nil, fmt.Errorf("invalid sort: %s", sort)
		}
	}
	return sorts, nil
}

func parsePageClause(pageIdx string, pageSize string) (PageSpec, error) {
	if pageIdx == "" {
		return PageSpec{}, nil
	}

	idx, err := strconv.ParseUint(pageIdx, 10, 64)
	if err != nil {
		return PageSpec{}, fmt.Errorf("invalid page index: %s", pageIdx)
	}

	var size uint64
	if pageSize != "" {
		size, err = strconv.ParseUint(pageSize, 10, 64)
		if err != nil {
			return PageSpec{}, fmt.Errorf("invalid page size: %s", pageSize)
		}
	}
	return PageSpec{Idx: idx, Size: size}, nil
}