package collection

type Page[T any] struct {
	Results		[]T
	This 		PageSpec
	Count		uint64
}

type PageSpec struct {
	Idx uint64
	Size uint64
}