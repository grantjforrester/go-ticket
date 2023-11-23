package collection_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grantjforrester/go-ticket/pkg/collection"
)

func TestShouldReturnEmpty(t *testing.T) {
	// Given
	query := MustParseQuery("")

	// When
	result, err := collection.ParseQuery(query)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Filters, 0)
	assert.Len(t, result.Sorts, 0)
	assert.Equal(t, result.Page, uint64(1))
	assert.Equal(t, result.Size, uint64(100))
}

func TestShouldReturnFilter(t *testing.T) {
	// Given
	query := MustParseQuery("filter=foo%3D%3Dbar")

	// When
	result, err := collection.ParseQuery(query)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Filters, 1)
	assert.Equal(t, result.Filters[0].Field, "foo")
	assert.Equal(t, result.Filters[0].Operator, collection.EQ)
	assert.Equal(t, result.Filters[0].Value, "bar")
}

func TestShouldReturnErrorOnInvalidFilter(t *testing.T) {
	// Given
	query := MustParseQuery("filter=foo")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "filter")
	assert.Contains(t, err.Error(), "foo")
}

func TestShouldReturnSort(t *testing.T) {
	// Given
	query := MustParseQuery("sort=foo+asc")

	// When
	result, err := collection.ParseQuery(query)

	// Then
	assert.Nil(t, err)
	assert.Len(t, result.Sorts, 1)
	assert.Equal(t, result.Sorts[0].Field, "foo")
	assert.Equal(t, result.Sorts[0].Direction, collection.ASC)
}

func TestShouldReturnErrorInInvalidSort(t *testing.T) {
	// Given
	query := MustParseQuery("sort=foo")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "sort")
	assert.Contains(t, err.Error(), "foo")
}

func TestShouldReturnPageSpecWithIndexAndNoSize(t *testing.T) {
	// Given
	query := MustParseQuery("page=1")

	// When
	result, err := collection.ParseQuery(query)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), result.Page)
	assert.Equal(t, uint64(100), result.Size)
}

func TestShouldReturnPageSpecWithIndexAndSize(t *testing.T) {
	// Given
	query := MustParseQuery("page=2&size=100")

	// When
	result, err := collection.ParseQuery(query)

	// Then
	assert.Nil(t, err)
	assert.Equal(t, uint64(2), result.Page)
	assert.Equal(t, uint64(100), result.Size)
}

func TestShouldReturnErrorOnZeroPageIndex(t *testing.T) {
	// Given
	query := MustParseQuery("page=0")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "page")
	assert.Contains(t, err.Error(), "0")
}

func TestShouldReturnErrorOnInvalidPageIndex(t *testing.T) {
	// Given
	query := MustParseQuery("page=foo")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "page")
	assert.Contains(t, err.Error(), "foo")
}

func TestShouldReturnErrorOnZeroPageSize(t *testing.T) {
	// Given
	query := MustParseQuery("page=1&size=0")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "size")
	assert.Contains(t, err.Error(), "0")
}

func TestShouldReturnErrorOnInvalidPageSize(t *testing.T) {
	// Given
	query := MustParseQuery("page=1&size=foo")

	// When
	_, err := collection.ParseQuery(query)

	// Then
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "size")
	assert.Contains(t, err.Error(), "foo")
}

func MustParseQuery(query string) url.Values {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}
	return parsed
}
