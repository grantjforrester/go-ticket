package cql_test

import (
	"testing"

	"github.com/grantjforrester/go-ticket/pkg/collection"
	"github.com/grantjforrester/go-ticket/pkg/collection/cql"
	"github.com/stretchr/testify/assert"
)

func TestReturnOnlySelectFrom(t *testing.T) {
	// Given
	q := cql.SQLQuery{Fields: []string{"foo"}, Table: "bar"}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar", sql)
	assert.Len(t, args, 0)
}

func TestReturnWhereEquals(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Fields: []string{"foo"},
		Table:  "bar",
		Query: collection.QuerySpec{
			Filters: []collection.FilterExpr{{
				Field:    "bam",
				Operator: cql.OpEq,
				Value:    "baz",
			}}}}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar WHERE bam = $1", sql)
	assert.Len(t, args, 1)
	assert.Equal(t, "baz", args[0])
}

func TestReturnWhereNotEquals(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Fields: []string{"foo"},
		Table:  "bar",
		Query: collection.QuerySpec{
			Filters: []collection.FilterExpr{{
				Field:    "bam",
				Operator: cql.OpNe,
				Value:    "baz",
			}}}}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar WHERE bam <> $1", sql)
	assert.Len(t, args, 1)
	assert.Equal(t, "baz", args[0])
}

func TestReturnOrderBy(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Table:  "bar",
		Fields: []string{"foo"},
		Query: collection.QuerySpec{
			Sorts: []collection.SortExpr{{
				Field:     "bar",
				Direction: cql.SortAsc,
			}}}}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar ORDER BY bar ASC", sql)
	assert.Len(t, args, 0)
}

func TestReturnLimit(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Table:  "bar",
		Fields: []string{"foo"},
		Query: collection.QuerySpec{
			Size: 10,
		},
	}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar LIMIT 10", sql)
	assert.Len(t, args, 0)
}

func TestReturnLimitAndOffset(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Table:  "bar",
		Fields: []string{"foo"},
		Query: collection.QuerySpec{
			Size: 10,
			Page: 2,
		},
	}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar LIMIT 10 OFFSET 10", sql)
	assert.Len(t, args, 0)
}

func TestReturnAllWhenPageWithNoSize(t *testing.T) {
	// Given
	q := cql.SQLQuery{
		Table:  "bar",
		Fields: []string{"foo"},
		Query: collection.QuerySpec{
			Page: 2,
		},
	}

	// When
	sql, args, err := q.ToSQL()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "SELECT foo FROM bar", sql)
	assert.Len(t, args, 0)
}
