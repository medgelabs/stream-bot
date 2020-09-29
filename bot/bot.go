package bot

import (
	"fmt"
	"log"
	"medgebot/irc"
	"sync"
)

type Bot struct {
	sync.Mutex
	client    *irc.Irc
	consumers []func(irc.Message)
}

func New() Bot {
	client := irc.NewClient()

	return Bot{
		client:    client,
		consumers: []func(irc.Message){},
	}
}

// Connect to the bot client
func (bot *Bot) Connect() error {
	if err := bot.client.Connect("wss", "irc-ws.chat.twitch.tv:443"); err != nil {
		log.Printf("ERROR: bot connect - %s", err)
		return err
	}

	go bot.readChat()
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
	return err
}

// PrivMsg sends a message to the given channel, without prefix
func (bot *Bot) SendMessage(channel, message string) string {
	return fmt.Sprintf("PRIVMSG #%s %s", channel, message)
}

// RegisterHandler registers a function that will be called concurrently when a message is received
func (bot *Bot) RegisterHandler(consumer func(irc.Message)) {
	bot.Mutex.Lock()
	defer bot.Mutex.Unlock()

	consumers := append(bot.consumers, consumer)
	bot.consumers = consumers
}

// readChat Reads from the client and passes the parsed messages to the stream channel
func (bot *Bot) readChat() {
	for {
		msg, err := bot.client.Read()
		if err != nil {
			log.Println("ERROR: read - " + err.Error())
			break
		}

		bot.Mutex.Lock()
		for _, consumer := range bot.consumers {
			go consumer(msg)
		}
		bot.Mutex.Unlock()
	}
}
