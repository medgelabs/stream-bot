package storage

type storable interface {
	put(key string, value interface{}) error
	get(key string) (interface{}, error)
}
