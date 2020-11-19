package bot

import (
	"strings"
)

// Register a handler to handle PING/PONG message exchange
func (bot *Bot) RegisterPingPong() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if strings.HasPrefix(evt.Message, "PING") {
				bot.client.SendPong([]string{evt.Message})
			}
		}),
	)
}
