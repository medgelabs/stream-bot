package bot

import (
	"medgebot/bot/bottest"
	"testing"
)

// Bits handler through the Bot
func TestCommandHandler(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.HandleCommands([]Command{
		{
			Prefix:          "!hello",
			MessageTemplate: NewHandlerTemplate(bottest.MakeTemplate("hello", "WORLD")),
		},
	})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Valid bits event
	evt := NewChatEvent()
	evt.Sender = "medgelabs"
	evt.Message = "!hello"
	bot.events <- evt

	response := <-checker.events
	if response.Message != "WORLD" {
		t.Fatalf("Got invalid !hello command response: %+v", response)
	}

}

func TestCommandHandlerIgnoresRegularChatMessages(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.HandleCommands([]Command{
		{
			Prefix:          "!hello",
			MessageTemplate: NewHandlerTemplate(bottest.MakeTemplate("hello", "WORLD")),
		},
	})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Invalid event
	evt := NewChatEvent()
	evt.Sender = "medgelabs"
	evt.Message = "hello" // no ! in front
	bot.events <- evt

	select {
	case resp := <-checker.events:
		t.Fatalf("Received message from CommandHandler for regular chat message - %+v", resp)
	default:
		// If we don't receive a response, the Bot didn't erroneously parse the wrong message
	}
}

func TestCommandHandlerIgnoresInvalidEvents(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.HandleCommands([]Command{})
	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Invalid event
	evt := NewRaidEvent()
	bot.events <- evt

	select {
	case resp := <-checker.events:
		t.Fatalf("Received message from CommandHandler for invalid message - %+v", resp)
	default:
		// If we don't receive a response, the Bot didn't erroneously parse the wrong message
	}
}
