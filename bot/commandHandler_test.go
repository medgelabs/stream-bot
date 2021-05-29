package bot

import (
	"medgebot/bot/bottest"
	"medgebot/cache"
	"strings"
	"testing"
)

// Bits handler through the Bot
func TestCommandHandler(t *testing.T) {
	// Initialize Bot
	cache, _ := cache.InMemory(0)
	bot := New(&cache)
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

func TestCoinThrow(t *testing.T) {
	// Initialize Bot
	cache, _ := cache.InMemory(0)
	bot := New(&cache)
	checker := NewTestChatClient()
	bot.SetChatClient(checker)
	bot.HandleCommands([]Command{})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Fire off !coin commands many times
	var responses []Event
	for i := 0; i < 100; i++ {
		evt := NewChatEvent()
		evt.Sender = "medgelabs"
		evt.Message = "!coin"
		bot.events <- evt

		resp := <-checker.events
		responses = append(responses, resp)
	}

	heads := 0
	tails := 0
	for _, event := range responses {
		if strings.Contains(event.Message, "heads") {
			heads++
		} else if strings.Contains(event.Message, "tails") {
			tails++
		} else {
			t.Logf("Unknown coin flip result - %s", event.Message)
		}
	}

	t.Logf("Heads: %d | Tails: %d", heads, tails)

	if heads < 10 {
		t.Fatalf("Oddly distributed amount of heads flips: %d", heads)
	}

	if tails < 10 {
		t.Fatalf("Oddly distributed amount of tails flips: %d", heads)
	}
}

func TestCommandHandlerIgnoresRegularChatMessages(t *testing.T) {
	// Initialize Bot
	cache, _ := cache.InMemory(0)
	bot := New(&cache)
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
	cache, _ := cache.InMemory(0)
	bot := New(&cache)
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
