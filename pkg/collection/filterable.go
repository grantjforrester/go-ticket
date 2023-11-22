package collection

type Operator int

const (
	EQ Operator = iota
	NE
	GT
	LT
	GE
	LE
)

type FilterSpec struct {
	Field		string
	Operator	Operator
	Value		any
}