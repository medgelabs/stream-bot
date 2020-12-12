package bot

import (
	"errors"
	"strings"
	"sync"
)

type Bot struct {
	sync.Mutex
	events          chan Event
	consumers       []Handler
	inboundPlugins  InboundPluginCollection
	outboundPlugins OutboundPluginCollection
	listening       bool
}

func New() Bot {
	return Bot{
		events:          make(chan Event, 0),
		consumers:       make([]Handler, 0),
		inboundPlugins:  make(InboundPluginCollection, 0),
		outboundPlugins: make(OutboundPluginCollection, 0),
		listening:       false,
	}
}

func (bot *Bot) sendEvent(evt Event) {
	for _, plugin := range bot.outboundPlugins {
		plugin.GetChannel() <- evt
	}
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
