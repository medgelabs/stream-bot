package ledger

// Ledger describes a persistent store of users in the chat
type Ledger interface {
	Absent(string) bool
	Add(string) error
}
