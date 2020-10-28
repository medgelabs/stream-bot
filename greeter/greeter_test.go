package greeter

import (
	"medgebot/ledger"
	"testing"
)

func TestGreet(t *testing.T) {
	ledger := ledger.NewInMemoryLedger()
	config := Config{
		MessageFormat: "Welcome %s",
	}

	expected := "Welcome Barry"
	g := New(config, &ledger)
	if greeting := g.Greet("Barry"); greeting != expected {
		t.Errorf("Incorrect greeting. Got %s, expecting %s", greeting, expected)
	}
}
