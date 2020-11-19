package bot

import (
	"log"
)

// Prints messages to the console
func (bot *Bot) RegisterReadLogger() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			log.Printf("> %s", evt)
		}),
	)
}
