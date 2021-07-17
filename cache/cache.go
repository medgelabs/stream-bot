package cache

// Acceptable types of caches
const (
	FILE = "file"
	MEM  = "mem"
)

// Cache describes a store of key/value pairs
type Cache interface {
	Absent(key string) bool
	Get(key string) (string, error)
	GetOrDefault(key, defaultValue string) (string, error)
	Put(key, value string) error
	Append(key, separator, value string) error
	Clear(key string)
}
