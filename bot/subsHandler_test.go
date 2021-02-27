package bot

import (
	"medgebot/bot/bottest"
	"testing"
)

var subsTmpl = bottest.MakeTemplate("testSubs", "{{.Sender}} subbed for {{.Amount}} months!")
var giftSubTmpl = bottest.MakeTemplate("testGiftSubs", "{{.Sender}} gifted a sub to {{.Recipient}}!")

// Subs handler through the Bot
func TestSubHandler(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.RegisterSubsHandler(
		HandlerTemplate{template: subsTmpl},
		HandlerTemplate{template: giftSubTmpl})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Valid bits event
	evt := NewSubEvent()
	evt.Sender = "srycantthnkof1"
	evt.Amount = 5
	bot.events <- evt

	response := <-checker.events
	if response.Message != "srycantthnkof1 subbed for 5 months!" {
		t.Fatalf("Invalid subscription response: %+v", response)
	}

}

func TestGiftSubHandler(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.RegisterSubsHandler(
		HandlerTemplate{template: subsTmpl},
		HandlerTemplate{template: giftSubTmpl})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Valid bits event
	evt := NewGiftSubEvent()
	evt.Sender = "BlackMarvel"
	evt.Recipient = "nojoy"
	bot.events <- evt

	response := <-checker.events
	if response.Message != "BlackMarvel gifted a sub to nojoy!" {
		t.Fatalf("Invalid gift subscription response: %+v", response)
	}

}

func TestSubHandlerIgnoresInvalidEvents(t *testing.T) {
	// Initialize Bot
	bot := New()
	checker := NewTestChatClient()
	bot.SetChatClient(checker)

	// Initialize Handler
	bot.RegisterSubsHandler(
		HandlerTemplate{template: subsTmpl},
		HandlerTemplate{template: giftSubTmpl})

	// This must happen after Handler registration, else data race occurs
	bot.Start()

	// Invalid event
	evt := NewBitsEvent()
	evt.Amount = 100
	bot.events <- evt

	select {
	case resp := <-checker.events:
		t.Fatalf("Received message from SubHandler for invalid message - %+v", resp)
	default:
		// If we don't receive a response, the Bot didn't erroneously parse the wrong message
	}
}
