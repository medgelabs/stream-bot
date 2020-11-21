package bot

import (
	"errors"
	"fmt"
	"sync"
)

type Bot struct {
	sync.Mutex
	events          chan Event
	channel         string
	consumers       []Handler
	inboundPlugins  InboundPluginCollection
	outboundPlugins OutboundPluginCollection
}

func New() Bot {
	return Bot{
		events:          make(chan Event, 0),
		channel:         "",
		consumers:       make([]Handler, 0),
		inboundPlugins:  make(InboundPluginCollection, 0),
		outboundPlugins: make(OutboundPluginCollection, 0),
	}
}

func (bot *Bot) sendEvent(evt Event) {
	for _, plugin := range bot.outboundPlugins {
		plugin.GetInboundChannel() <- evt
	}
}

func (bot *Bot) sendEventToPlugin(id string, evt Event) error {
	plugin, ok := bot.outboundPlugins[id]
	if !ok {
		return errors.New("outbound plugin not registered")
	}

	plugin.GetInboundChannel() <- evt

	return nil
}

// Start the bot and listen for incoming events
func (bot *Bot) Start() error {
	// Ensure single concurrent reader, per doc requirements
	go bot.listen()

	return nil
}

// PrivMsg sends a message to the given channel, without prefix
func (bot *Bot) SendMessage(message string) {
	evt := NewChatEvent()
	evt.Message = message

	go bot.sendEvent(evt)
}

// RegisterHandler registers a function that will be called concurrently when a message is received
func (bot *Bot) RegisterHandler(consumer Handler) {
	bot.Mutex.Lock()
	defer bot.Mutex.Unlock()

	consumers := append(bot.consumers, consumer)
	bot.consumers = consumers
}

// Start listening for Events on the inbound channel and broadcast out
// to the Handlers
func (bot *Bot) listen() {
	// Spawn goroutines for Handlers
	for _, consumer := range bot.consumers {
		go consumer.Listen()
	}

	for {
		select {
		case evt := <-bot.events:
			bot.Mutex.Lock()
			fmt.Printf("%+v", evt)
			for _, consumer := range bot.consumers {
				consumer.Receive(evt)
			}
			bot.Mutex.Unlock()
		default:
			// TODO QUIT message
		}
	}
}
