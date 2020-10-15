package ledger

type persistence interface {
	absent(string) bool
	add(string) error
}
