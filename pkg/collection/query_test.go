package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

func TestShouldReturnEmpty(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{}
	var pageIdx = ""
	var pageSize = ""

	// When
	result, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Filters, 0)
	assert.Len(t, result.Sorts, 0)
	assert.Equal(t, result.Page.Idx, uint64(0))
	assert.Equal(t, result.Page.Size, uint64(0))
}

func TestShouldReturnFilter(t *testing.T) {
	// Given
	var filters = []string{"foo==bar"}
	var sorts = []string{}
	var pageIdx = ""
	var pageSize = ""

	// When
	result, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, result.Filters[0].Field, "foo")
	assert.Equal(t, result.Filters[0].Operator, collection.EQ)
	assert.Equal(t, result.Filters[0].Value, "bar")
}

func TestShouldReturnErrorOnInvalidFilter(t *testing.T) {
	// Given
	var filters = []string{ "foo" }
	var sorts = []string{}
	var pageIdx = ""
	var pageSize = ""

	// When
	_, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "filter")
}

func TestShouldReturnSort(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{"foo asc"}
	var pageIdx = ""
	var pageSize = ""

	// When
	result, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Sorts, 1)
	assert.Equal(t, result.Sorts[0].Field, "foo")
	assert.Equal(t, result.Sorts[0].Direction, collection.ASC)
}

func TestShouldReturnErrorInInvalidSort(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{"asc"}
	var pageIdx = ""
	var pageSize = ""

	// When
	_, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "sort")
}

func TestShouldReturnPageSpecWithIndexAndNoSize(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{}
	var pageIdx = "1"
	var pageSize = ""

	// When
	result, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, collection.PageSpec{ Idx: 1 }, result.Page)
}

func TestShouldReturnPageSpecWithIndexAndSize(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{}
	var pageIdx = "2"
	var pageSize = "100"

	// When
	result, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, collection.PageSpec{ Idx: 2, Size: 100 }, result.Page)
}

func TestShouldReturnErrorOnInvalidPageIndex(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{}
	var pageIdx = "foo"
	var pageSize = ""

	// When
	_, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "page index")
}

func TestShouldReturnErrorOnInvalidPageSize(t *testing.T) {
	// Given
	var filters = []string{}
	var sorts = []string{}
	var pageIdx = "1"
	var pageSize = "foo"

	// When
	_, err := collection.ParseQuery(filters, sorts, pageIdx, pageSize)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "page size")
}