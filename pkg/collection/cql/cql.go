package cql

import (
	"fmt"

	"net/url"
	"regexp"
	"strconv"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

// symbols
const (
	ParamPage   string               = "page"
	ParamSize   string               = "size"
	ParamSort   string               = "sort"
	ParamFilter string               = "filter"
	SortAsc     collection.Direction = "asc"
	SortDesc    collection.Direction = "desc"
	OpEq        collection.Operator  = "=="
	OpNe        collection.Operator  = "!="
	OpGt        collection.Operator  = ">"
	OpLt        collection.Operator  = "<"
	OpGe        collection.Operator  = ">="
	OpLe        collection.Operator  = "<="
)

var StringOps = []collection.Operator{OpEq, OpNe}
var BoolOps = StringOps
var NumberOps = []collection.Operator{OpEq, OpNe, OpGt, OpLt, OpGe, OpLe}

var fieldPattern = `\w+`
var valuePattern = `.+`
var operatorPattern = fmt.Sprintf("%s|%s|%s|%s|%s|%s", OpEq, OpNe, OpGt, OpLt, OpGe, OpLe)
var filterPattern = fmt.Sprintf("(%s)(%s)(%s)", fieldPattern, operatorPattern, valuePattern)
var sortPattern = fmt.Sprintf("(%s) (%s|%s)", fieldPattern, SortAsc, SortDesc)

func ParseQuery(urlQuery url.Values) (collection.QuerySpec, error) {
	var (
		frx = regexp.MustCompile(filterPattern)
		srx = regexp.MustCompile(sortPattern)
	)

	q := collection.QuerySpec{}

	fs, err := parseFilters(urlQuery[ParamFilter], *frx)
	if err != nil {
		return collection.QuerySpec{}, err
	}
	q.Filters = fs

	ss, err := parseSorts(urlQuery[ParamSort], *srx)
	if err != nil {
		return collection.QuerySpec{}, err
	}
	q.Sorts = ss

	pg, err := parsePage(urlQuery.Get(ParamPage))
	if err != nil {
		return collection.QuerySpec{}, err
	}
	q.Page = pg

	sz, err := parseSize(urlQuery.Get(ParamSize))
	if err != nil {
		return collection.QuerySpec{}, err
	}
	q.Size = sz

	return q, nil
}

func parseFilters(filters []string, regex regexp.Regexp) ([]collection.FilterExpr, error) {
	fltrs := []collection.FilterExpr{}

	for _, f := range filters {
		if regex.MatchString(f) {
			parts := regex.FindStringSubmatch(f)
			fltrs = append(fltrs, collection.FilterExpr{Field: parts[1], Operator: collection.Operator(parts[2]), Value: parts[3]})
		} else {
			return nil, collection.QueryError{Message: fmt.Sprintf("invalid filter: %s", f)}
		}
	}

	return fltrs, nil
}

func parseSorts(clauses []string, regexp regexp.Regexp) ([]collection.SortExpr, error) {
	srts := []collection.SortExpr{}

	for _, s := range clauses {
		if regexp.MatchString(s) {
			parts := regexp.FindStringSubmatch(s)
			srts = append(srts, collection.SortExpr{Field: parts[1], Direction: collection.Direction(parts[2])})
		} else {
			return nil, collection.QueryError{Message: fmt.Sprintf("invalid sort: %s", s)}
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
		return 0, collection.QueryError{Message: fmt.Sprintf("invalid page: %s", page)}
	}

	return pg, nil
}

func parseSize(size string) (uint64, error) {
	if size == "" {
		return 0, nil
	}

	sz, err := strconv.ParseUint(size, 10, 64)
	if err != nil || sz == 0 {
		return 0, collection.QueryError{Message: fmt.Sprintf("invalid size: %s", size)}
	}

	return sz, nil
}
