package collection

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

// symbols
const (
	PAGE   string    = "page"
	SIZE   string    = "size"
	SORT   string    = "sort"
	FILTER string    = "filter"
	ASC    Direction = "asc"
	DESC   Direction = "desc"
	EQ     Operator  = "=="
	NE     Operator  = "!="
	GT     Operator  = ">"
	LT     Operator  = "<"
	GE     Operator  = ">="
	LE     Operator  = "<="

	DEFAULT_PAGE = uint64(1)
	DEFAULT_SIZE = uint64(100)
)

var FIELD_PATTERN = `\w+`
var VALUE_PATTERN = `.+`
var OPERATOR_PATTERN = fmt.Sprintf("%s|%s|%s|%s|%s|%s", EQ, NE, GT, LT, GE, LE)
var FILTER_PATTERN = fmt.Sprintf("(%s)(%s)(%s)", FIELD_PATTERN, OPERATOR_PATTERN, VALUE_PATTERN)
var SORT_PATTERN = fmt.Sprintf("(%s) (%s|%s)", FIELD_PATTERN, ASC, DESC)

func ParseQuery(urlQuery url.Values) (Query, error) {
	var frx = regexp.MustCompile(FILTER_PATTERN)
	var srx = regexp.MustCompile(SORT_PATTERN)

	q := Query{}

	fs, err := parseFilters(urlQuery[FILTER], *frx)
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Filters = fs

	ss, err := parseSorts(urlQuery[SORT], *srx)
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Sorts = ss

	pg, sz, err := parsePage(urlQuery.Get(PAGE), urlQuery.Get(SIZE))
	if err != nil {
		return Query{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Page = pg
	q.Size = sz

	return q, nil
}

func parseFilters(filters []string, regex regexp.Regexp) ([]FilterSpec, error) {
	fltrs := []FilterSpec{}

	for _, f := range filters {
		if regex.MatchString(f) {
			parts := regex.FindStringSubmatch(f)
			fltrs = append(fltrs, FilterSpec{Field: parts[1], Operator: Operator(parts[2]), Value: parts[3]})
		} else {
			return nil, fmt.Errorf("invalid filter: %s", f)
		}
	}

	return fltrs, nil
}

func parseSorts(clauses []string, regexp regexp.Regexp) ([]SortSpec, error) {
	srts := []SortSpec{}

	for _, s := range clauses {
		if regexp.MatchString(s) {
			parts := regexp.FindStringSubmatch(s)
			srts = append(srts, SortSpec{Field: parts[1], Direction: Direction(parts[2])})
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
