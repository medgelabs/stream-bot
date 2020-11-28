package bot

import (
	"log"
)

// Prints Chat messages to the console
func (bot *Bot) RegisterReadLogger() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			log.Printf("> %s: %s", evt.Sender, evt.Message)
		}),
	)
}
