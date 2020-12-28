package bot

// ChatClient is a connector to Chat that the Bot sends messages to
type ChatClient interface {
	Channel() chan<- Event
}

// TestChatClient emulates a ChatClient receiving messages from the bot
// Only use in test!
type TestChatClient struct {
	events chan Event
}

func NewTestChatClient() TestChatClient {
	return TestChatClient{
		events: make(chan Event),
	}
}

// bot.ChatClient
func (c TestChatClient) Channel() chan<- Event {
	return c.events
}
