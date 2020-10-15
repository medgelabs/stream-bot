package ledger

type InMemory struct {
	ledger map[string]int
}

func NewInMemory() InMemory {
	return InMemory{
		ledger: make(map[string]int),
	}
}

func (l *InMemory) absent(key string) bool {
	_, ok := l.ledger[key]
	return !ok
}

func (l *InMemory) add(key string) error {
	l.ledger[key] = 1
	return nil
}
