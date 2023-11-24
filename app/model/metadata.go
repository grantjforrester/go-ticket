package model

type Metadata struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

func NewMetadata(id string, version string) Metadata {
	return Metadata{Id: id, Version: version}
}
