package collection

type Direction int

const (
	ASC Direction = iota
	DESC
)

type SortSpec struct {
	Field     string
	Direction Direction
}
