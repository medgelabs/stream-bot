package bot

import (
	"log"
	"medgebot/irc"
	"sync"
)

type Bot struct {
	sync.Mutex
	inboundEvents  <-chan Event
	outboundEvents chan<- Event
	channel        string
	consumers      []Handler

	client *irc.Irc
}

func New() Bot {
	client := irc.NewClient()

	return Bot{
		client:    client,
		channel:   "",
		consumers: make([]Handler, 0),
	}
}

// Connect to the bot client
func (bot *Bot) Connect() error {
	if err := bot.client.Connect("wss", "irc-ws.chat.twitch.tv:443"); err != nil {
		log.Printf("ERROR: bot connect - %s", err)
		return err
	}

	// Capabilities required
	// if err := bot.client.CapReq("twitch.tv/tags"); err != nil {
	// log.Printf("ERROR: send CAP REQ failed: %s", err)
	// return err
	// }

	// Ensure single concurrent reader, per doc requirements
	go bot.listen()

	return nil
}

// Close the connection to the client
func (bot *Bot) Close() {
	bot.client.Close()
}

// Authenticate connects to the IRC stream with the given nick and password
func (bot *Bot) Authenticate(nick, password string) error {
	if err := bot.client.SendPass(password); err != nil {
		log.Printf("ERROR: send PASS failed: %s", err)
		return err
	}
	log.Println("< PASS ***")

	if err := bot.client.SendNick(nick); err != nil {
		log.Printf("ERROR: send NICK failed: %s", err)
		return err
	}

	return nil
}

// Join joins to a specific channel on the IRC
func (bot *Bot) Join(channel string) error {
	err := bot.client.Join(channel)
	bot.channel = channel

	return err
}

// PrivMsg sends a message to the given channel, without prefix
func (bot *Bot) SendMessage(message string) {
	evt := NewChatEvent()
	evt.Message = message

	bot.outboundEvents <- evt
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
		case evt := <-bot.inboundEvents:
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
