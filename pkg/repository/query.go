package repository

// Query is an abstraction of criteria for searching a persistent store for matching entities.
// Each implementation of repository will provide its own implementation of Query.
type Query interface{}
