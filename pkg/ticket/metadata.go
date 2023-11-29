package ticket

type Metadata struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

func NewMetadata(id string, version string) Metadata {
	return Metadata{ID: id, Version: version}
}
