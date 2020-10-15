package ledger

type Ledger struct {
	Engine persistence
}

func New(engine persistence) *Ledger {
	return &Ledger{
		Engine: engine,
	}
}

func (s *Ledger) Absent(value string) bool {
	return s.Engine.absent(value)
}

func (s *Ledger) Add(key string) error {
	return s.Engine.add(key)
}
