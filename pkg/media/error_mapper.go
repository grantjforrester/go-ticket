package media

type ErrorMapper interface {
	MapError(err error) (int, any)
}

