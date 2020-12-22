package bot

import (
	"errors"
	"strings"
	"sync"
)

type Bot struct {
	sync.Mutex

	// Handlers of bot.Event messages
	consumers []Handler

	// Producers of data from external systems
	clients []Client

	// Where the Bot sends messages to get to Chat
	chatClient ChatClient

	// events
	events    chan Event
	listening bool
}

// Client represents data being sent TO the Bot
type Client interface {
	SetDestination(chan<- Event)
}

// ChatClient is a connector to Chat that the Bot sends messages to
type ChatClient interface {
	Channel() chan<- Event
}

func New() Bot {
	return Bot{
		consumers: make([]Handler, 0),
		clients:   make([]Client, 0),
		events:    make(chan Event, 0),
		listening: false,
	}
}

// RegisterClient links a Client that will send data TO the Bot.
// This method also set's the Client's Destination channel
func (bot *Bot) RegisterClient(client Client) {
	bot.Lock()
	defer bot.Unlock()

	client.SetDestination(bot.events)
	bot.clients = append(bot.clients, client)
}

// SetChatClient register's the client that will allow the Bot to send
// messages to Chat
func (bot *Bot) SetChatClient(client ChatClient) {
	bot.Lock()
	defer bot.Unlock()

	bot.chatClient = client
}

// sendEvent sends a Bot event to Write-enabled clients
func (bot *Bot) sendEvent(evt Event) {
	// TODO trace outbound
	bot.chatClient.Channel() <- evt
}

// Start the bot and listen for incoming events
func (bot *Bot) Start() error {
	// Ensure single concurrent reader, per doc requirements
	go bot.listen()

	return nil
}

// PrivMsg sends a message to the given channel, without prefix
func (bot *Bot) SendMessage(message string) {
	// TODO do we ever need to send empty messages?
	if strings.TrimSpace(message) == "" {
		return
	}

	evt := NewChatEvent()
	evt.Message = message

	go bot.sendEvent(evt)
}

// RegisterHandler registers a function that will be called concurrently when a message is received
func (bot *Bot) RegisterHandler(consumer Handler) error {
	if bot.listening {
		return errors.New("RegisterHandler called after bot already listening")
	}

	bot.Lock()
	defer bot.Unlock()

	consumers := append(bot.consumers, consumer)
	bot.consumers = consumers
	return nil
}

// Start listening for Events on the inbound channel and broadcast out
// to the Handlers
func (bot *Bot) listen() {
	// Spawn goroutines for Handlers
	for _, consumer := range bot.consumers {
		go consumer.Listen()
	}

	bot.listening = true

	for {
		select {
		case evt := <-bot.events:
			bot.Mutex.Lock()
			for _, consumer := range bot.consumers {
				consumer.Receive(evt)
			}
			bot.Mutex.Unlock()
		default:
			// TODO QUIT message
		}
	}
}
