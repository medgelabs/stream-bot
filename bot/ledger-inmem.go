package bot

// InMemoryLedger provides in memory ledger capabilities
type InMemoryLedger struct {
	ledger map[string]int
}

func NewInMemoryLedger() InMemoryLedger {
	return InMemoryLedger{
		ledger: make(map[string]int),
	}
}

func (l *InMemoryLedger) Absent(key string) bool {
	_, ok := l.ledger[key]
	return !ok
}

func (l *InMemoryLedger) Add(key string) error {
	l.ledger[key] = 1
	return nil
}
