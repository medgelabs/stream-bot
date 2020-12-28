package bot

import (
	"medgebot/bot/bottest"
	"testing"
)

// Raid handler through the Bot
func TestRaidHandler(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	tmpl := bottest.MakeTemplate("testRaid", "Welcome {{.Sender}}'s {{.Amount}} raiders!")
	bot.RegisterRaidHandler(HandlerTemplate{
		templ: tmpl,
	}, 1)

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Valid bits event
	evt := NewRaidEvent()
	evt.Sender = "shito86"
	evt.Amount = 5
	bot.events <- evt

	response := <-checker.events
	if response.Message != "Welcome shito86's 5 raiders!" {
		t.Fatalf("Invalid raid response: %+v", response)
	}

}

func TestRaidHandlerIgnoresInvalidEvents(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	tmpl := bottest.MakeTemplate("testRaid", "Welcome {{.Sender}}'s {{.Amount}} raiders!")
	bot.RegisterRaidHandler(HandlerTemplate{
		templ: tmpl,
	}, 0)

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Invalid event
	evt := NewBitsEvent()
	evt.Amount = 100
	bot.events <- evt

	select {
	case resp := <-checker.events:
		t.Fatalf("Received message from RaidHandler for invalid message - %+v", resp)
	default:
		// If we don't receive a response, the Bot didn't erroneously parse the wrong message
	}
}
