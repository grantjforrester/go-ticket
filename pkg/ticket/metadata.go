package ticket

// Metadata holds information common to all domain entities.
type Metadata struct {
	// A unique identifier for the domain entity.
	ID string `json:"id"`

	// A version identifier that changes as the domain entity's properties change.
	Version string `json:"version"`
}

// Creates a new Metadata.
func NewMetadata(id string, version string) Metadata {
	return Metadata{ID: id, Version: version}
}
