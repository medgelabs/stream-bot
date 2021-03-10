package cache

// Cache describes a store of key/value pairs
type Cache interface {
	Absent(key string) bool
	Get(key string) (string, error)
	Put(key string, value string) error
}
