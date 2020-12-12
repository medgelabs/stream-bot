package greeter

import (
	"medgebot/ledger"
	"testing"
)

func TestGreet(t *testing.T) {
	ledger := ledger.NewInMemoryLedger()
	New(&ledger)

	// TODO test HasNotGreeted() and RecordGreeting()
}
