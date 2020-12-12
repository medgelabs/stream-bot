package greeter

import (
	"medgebot/ledger"
)

type Greeter struct {
	ledger ledger.Ledger
}

func New(ledger ledger.Ledger) Greeter {
	return Greeter{
		ledger: ledger,
	}
}

// HasNotGreeted determines if the given target has not been greeted
func (g Greeter) HasNotGreeted(target string) bool {
	return g.ledger.Absent(target)
}

// RecordGreeting makes note that someone has been greeted
func (g Greeter) RecordGreeting(target string) error {
	return g.ledger.Add(target)
}
