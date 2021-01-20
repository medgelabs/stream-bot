package ledger

// Ledger describes a persistent store of key/value pairs
type Ledger interface {
	Absent(key string) bool
	Get(key string) (string, error)
	Put(key string, value string) error
}
