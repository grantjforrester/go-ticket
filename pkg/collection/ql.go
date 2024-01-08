package collection

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

// symbols
const (
	ParamPage   string    = "page"
	ParamSize   string    = "size"
	ParamSort   string    = "sort"
	ParamFilter string    = "filter"
	SortAsc     Direction = "asc"
	SortDesc    Direction = "desc"
	OpEq        Operator  = "=="
	OpNe        Operator  = "!="
	OpGt        Operator  = ">"
	OpLt        Operator  = "<"
	OpGe        Operator  = ">="
	OpLe        Operator  = "<="
)

var FieldPattern = `\w+`
var ValuePattern = `.+`
var OperatorPattern = fmt.Sprintf("%s|%s|%s|%s|%s|%s", OpEq, OpNe, OpGt, OpLt, OpGe, OpLe)
var FilterPattern = fmt.Sprintf("(%s)(%s)(%s)", FieldPattern, OperatorPattern, ValuePattern)
var SortPattern = fmt.Sprintf("(%s) (%s|%s)", FieldPattern, SortAsc, SortDesc)

func ParseQuery(urlQuery url.Values) (QuerySpec, error) {
	var (
		frx = regexp.MustCompile(FilterPattern)
		srx = regexp.MustCompile(SortPattern)
	)

	q := QuerySpec{}

	fs, err := parseFilters(urlQuery[ParamFilter], *frx)
	if err != nil {
		return QuerySpec{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Filters = fs

	ss, err := parseSorts(urlQuery[ParamSort], *srx)
	if err != nil {
		return QuerySpec{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Sorts = ss

	pg, err := parsePage(urlQuery.Get(ParamPage))
	if err != nil {
		return QuerySpec{}, fmt.Errorf("invalid query: %w", err)
	}
	q.Page = pg

	sz, err := parseSize(urlQuery.Get(ParamSize))
	if err != nil {
		return QuerySpec{}, fmt.Errorf("invalid query: %w", err)
	}
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

func parsePage(page string) (uint64, error) {
	if page == "" {
		return 0, nil
	}

	pg, err := strconv.ParseUint(page, 10, 64)
	if err != nil || pg == 0 {
		return 0, fmt.Errorf("invalid page: %s", page)
	}

	return pg, nil
}

func parseSize(size string) (uint64, error) {
	if size == "" {
		return 0, nil
	}

	sz, err := strconv.ParseUint(size, 10, 64)
	if err != nil || sz == 0 {
		return 0, fmt.Errorf("invalid size: %s", size)
	}

	return sz, nil
}
