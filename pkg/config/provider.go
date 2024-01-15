package config

// Provider describes a common pattern for passing configuration values.
type Provider interface {
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
}
