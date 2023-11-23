package collection

type Page[T any] struct {
	Results		[]T		`json:"results"`
	Page 		uint64	`json:"page"`
	Size		uint64	`json:"size"`
}