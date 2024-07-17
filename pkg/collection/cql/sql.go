package cql

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

// SQLQuery wraps a QuerySpec and generates SQL.
type SQLQuery struct {
	Fields []string
	Table  string
	Query  collection.QuerySpec
}

// Returns SQL and and arguments for placeholders.
func (q SQLQuery) ToSQL() (string, []any, error) {
	sql := sq.Select(q.Fields...).From(q.Table).OrderBy(q.orderBy()...)

	for _, f := range q.Query.Filters {
		foo := mapFilter(f)
		sql = sql.Where(foo)
	}

	if q.limit() > 0 {
		sql = sql.Limit(q.limit())
	}

	if q.offset() > 0 {
		sql = sql.Offset(q.offset())
	}

	return sql.
		PlaceholderFormat(sq.Dollar).
		ToSql()
}

// returns SQL ORDER BY clause from QuerySpec.
func (q SQLQuery) orderBy() []string {
	var clause = make([]string, len(q.Query.Sorts))
	for i, s := range q.Query.Sorts {
		clause[i] = fmt.Sprintf("%s %s", s.Field, mapDirection(s.Direction))
	}
	return clause
}

// returns SQL limit from QuerySpec.
func (q SQLQuery) limit() uint64 {
	return q.Query.Size
}

// returns SQL offset form QuerySpec.
func (q SQLQuery) offset() uint64 {
	if q.Query.Size == 0 || q.Query.Page == 0 {
		return 0
	}

	return (q.Query.Page - 1) * q.Query.Size
}

// returns a given Squirrel expression for a filter
func mapFilter(filter collection.FilterExpr) any {
	switch filter.Operator {
	case OpEq:
		return sq.Eq{filter.Field: filter.Value}
	case OpNe:
		return sq.NotEq{filter.Field: filter.Value}
	case OpLt:
		return sq.Lt{filter.Field: filter.Value}
	case OpLe:
		return sq.LtOrEq{filter.Field: filter.Value}
	case OpGt:
		return sq.Gt{filter.Field: filter.Value}
	case OpGe:
		return sq.GtOrEq{filter.Field: filter.Value}
	default:
		panic(fmt.Sprintf("unknown query filter operator: %s", filter.Operator))
	}
}

// returns a given Squirrel order by direction for a given direction
func mapDirection(direction collection.Direction) string {
	switch direction {
	case SortAsc:
		return "ASC"
	case SortDesc:
		return "DESC"
	default:
		panic(fmt.Sprintf("unknown query sort direction: %s", direction))
	}
}
