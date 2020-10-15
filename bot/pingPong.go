package bot

import "medgebot/internal/pkg/irc"

// Register a handler to handle PING/PONG message exchange
func (bot *Bot) RegisterPingPong() {
	bot.RegisterHandler(func(msg irc.Message) {
		if msg.Command == "PING" {
			bot.client.SendPong(msg.Params)
		}
	})
}
