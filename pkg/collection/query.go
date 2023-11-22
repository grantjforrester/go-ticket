package collection

type Query struct {
	Filters []FilterSpec
	Sorts	[]SortSpec
	Page	PageSpec
}