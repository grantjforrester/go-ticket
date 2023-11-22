package rql

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

// symbols
const (
	PAGE     = "page"
	SIZE     = "size"
	SORT     = "sort"
	FILTER   = "filter"
	ASC_SEQ  = "asc"
	DESC_SEQ = "desc"
	EQ_SEQ   = "=="
	NE_SEQ   = "!="
	GT_SEQ   = ">"
	LT_SEQ   = "<"
	GE_SEQ   = ">="
	LE_SEQ   = "<="
)

var DIRECTIONS = map[string]collection.Direction{
	ASC_SEQ:  collection.ASC,
	DESC_SEQ: collection.DESC,
}

var OPERATORS = map[string]collection.Operator{
	EQ_SEQ: collection.EQ,
	NE_SEQ: collection.NE,
	GT_SEQ: collection.GT,
	LT_SEQ: collection.LT,
	GE_SEQ: collection.GE,
	LE_SEQ: collection.LE,
}

var FIELD_PATTERN = `\w+`
var VALUE_PATTERN = `.+`
var OPERATOR_PATTERN = fmt.Sprintf("%s|%s|%s|%s|%s|%s", EQ_SEQ, NE_SEQ, GT_SEQ, LT_SEQ, GE_SEQ, LE_SEQ)
var FILTER_PATTERN = fmt.Sprintf("(%s)(%s)(%s)", FIELD_PATTERN, OPERATOR_PATTERN, VALUE_PATTERN)
var SORT_PATTERN = fmt.Sprintf("(%s) (%s|%s)", FIELD_PATTERN, ASC_SEQ, DESC_SEQ)

func Parse(urlQuery url.Values) (collection.Query, error) {
	var sort_re = regexp.MustCompile(SORT_PATTERN)
	var filter_re = regexp.MustCompile(FILTER_PATTERN)

	query := collection.Query{}

	filterSpecs, err := parseFilters(urlQuery[FILTER], *filter_re)
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Filters = filterSpecs

	sortSpecs, err := parseSorts(urlQuery[SORT], *sort_re)
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Sorts = sortSpecs

	pageSpec, err := parsePage(urlQuery.Get(PAGE), urlQuery.Get(SIZE))
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	query.Page = pageSpec

	return query, nil
}

func parseFilters(clauses []string, regexp regexp.Regexp) ([]collection.FilterSpec, error) {
	if len(clauses) == 0 {
		return []collection.FilterSpec{}, nil
	}

	filters := []collection.FilterSpec{}
	for _, filter := range clauses {
		if regexp.MatchString(filter) {
			parts := regexp.FindStringSubmatch(filter)
			field := parts[1]
			operator := parts[2]
			value := parts[3]
			filters = append(filters, collection.FilterSpec{Field: field, Operator: OPERATORS[operator], Value: value})
		} else {
			return nil, fmt.Errorf("invalid filter: %s", filter)
		}
	}
	return filters, nil
}

func parseSorts(clauses []string, regexp regexp.Regexp) ([]collection.SortSpec, error) {
	if len(clauses) == 0 {
		return []collection.SortSpec{}, nil
	}

	sorts := []collection.SortSpec{}
	for _, sort := range clauses {
		if regexp.MatchString(sort) {
			parts := regexp.FindStringSubmatch(sort)
			field := parts[1]
			direction := parts[2]
			sorts = append(sorts, collection.SortSpec{Field: field, Direction: DIRECTIONS[direction]})
		} else {
			return nil, fmt.Errorf("invalid sort: %s", sort)
		}
	}
	return sorts, nil
}

func parsePage(pageIdx string, pageSize string) (collection.PageSpec, error) {
	if pageIdx == "" {
		return collection.PageSpec{}, nil
	}

	idx, err := strconv.ParseUint(pageIdx, 10, 64)
	if err != nil || idx == 0 {
		return collection.PageSpec{}, fmt.Errorf("invalid page index: %s", pageIdx)
	}

	var size uint64
	if pageSize != "" {
		size, err = strconv.ParseUint(pageSize, 10, 64)
		if err != nil || size == 0 {
			return collection.PageSpec{}, fmt.Errorf("invalid page size: %s", pageSize)
		}
	}
	return collection.PageSpec{Idx: idx, Size: size}, nil
}