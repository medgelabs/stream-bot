package bot

import (
	"log"
	"medgebot/internal/pkg/irc"
)

// Prints messages to the console
func (bot *Bot) RegisterReadLogger() {
	bot.RegisterHandler(func(msg irc.Message) {
		log.Printf("> %s", msg.String())
	})
}
