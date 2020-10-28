package greeter

import (
	"fmt"
	"medgebot/ledger"
)

type Greeter struct {
	config Config
	ledger ledger.Ledger
}

func New(config Config, ledger ledger.Ledger) Greeter {
	return Greeter{
		config: config,
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

// Greet generates a greeting for the given target and
func (g Greeter) Greet(target string) string {
	return fmt.Sprintf(g.config.MessageFormat, target)
}
