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

	DEFAULT_PAGE = uint64(1)
	DEFAULT_SIZE = uint64(100)
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
	var frx = regexp.MustCompile(FILTER_PATTERN)
	var srx = regexp.MustCompile(SORT_PATTERN)

	q := collection.Query{}

	fs, err := parseFilters(urlQuery[FILTER], *frx)
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Filters = fs

	ss, err := parseSorts(urlQuery[SORT], *srx)
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Sorts = ss

	pg, sz, err := parsePage(urlQuery.Get(PAGE), urlQuery.Get(SIZE))
	if err != nil {
		return collection.Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Page = pg
	q.Size = sz

	return q, nil
}

func parseFilters(filters []string, regex regexp.Regexp) ([]collection.FilterSpec, error) {
	fltrs := []collection.FilterSpec{}

	for _, f := range filters {
		if regex.MatchString(f) {
			parts := regex.FindStringSubmatch(f)
			fltrs = append(fltrs, collection.FilterSpec{Field: parts[1], Operator: OPERATORS[parts[2]], Value: parts[3]})
		} else {
			return nil, fmt.Errorf("invalid filter: %s", f)
		}
	}

	return fltrs, nil
}

func parseSorts(clauses []string, regexp regexp.Regexp) ([]collection.SortSpec, error) {
	srts := []collection.SortSpec{}

	for _, s := range clauses {
		if regexp.MatchString(s) {
			parts := regexp.FindStringSubmatch(s)
			srts = append(srts, collection.SortSpec{Field: parts[1], Direction: DIRECTIONS[parts[2]]})
		} else {
			return nil, fmt.Errorf("invalid sort: %s", s)
		}
	}

	return srts, nil
}

func parsePage(page string, size string) (uint64, uint64, error) {
	var pg = DEFAULT_PAGE
	var err error

	if page != "" {
		pg, err = strconv.ParseUint(page, 10, 64)
		if err != nil || pg == 0 {
			return 0, 0, fmt.Errorf("invalid page: %s", page)
		}
	}

	var sz = DEFAULT_SIZE
	if size != "" {
		sz, err = strconv.ParseUint(size, 10, 64)
		if err != nil || sz == 0 {
			return 0, 0, fmt.Errorf("invalid size: %s", size)
		}
	}

	return pg, sz, nil
}
