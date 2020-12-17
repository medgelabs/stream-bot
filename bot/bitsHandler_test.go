package bot

import "testing"

// Integration Test for Bits handler through the Bot
func TestBitsHandler(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewChecker()
	bot.RegisterOutboundPlugin(checker)

	// Initialize Bits Handler
	tmpl := makeTemplate("testBits", "Thanks for the {{.Amount}} bits {{.Sender}}")
	bot.RegisterBitsHandler(HandlerTemplate{
		templ: tmpl,
	})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Valid bits event
	evt := NewBitsEvent()
	evt.Sender = "ReallyFrank"
	evt.Amount = 100
	bot.events <- evt

	response := <-checker.Events
	if response.Message != "Thanks for the 100 bits ReallyFrank" {
		t.Fatalf("Got invalid bits response: %+v", response)
	}

	// Invalid event
	evt = NewRaidEvent()
	bot.events <- evt

	select {
	case resp := <-checker.Events:
		t.Fatalf("Received message from BitsHandler for invalid message - %+v", resp)
	default:
		// If we don't receive a response, the Bot didn't erroneously parse the wrong message
	}
}
