package bot

import "fmt"

type Bot struct{}

// Authenticate connects to the IRC stream with the given nick and password
func (bot *Bot) Authenticate(nick, password string) error {
	return nil
}

// Join joins to a specific channel on the IRC
func (bot *Bot) Join(channel string) error {
	return nil
}

// PrivMsg sends a message to the given channel, without prefix
func (bot *Bot) PrivMsg(channel, message string) string {
	// PREFIX PRIVMSG #channel :message
	return fmt.Sprintf("PRIVMSG #%s %s", channel, message)
}

// Pong responds to Ping heartbeats
func (bot *Bot) Pong() error {
	return nil
}
