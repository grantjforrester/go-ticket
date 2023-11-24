package config

type Provider interface {
	Get(key string) any
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
}
