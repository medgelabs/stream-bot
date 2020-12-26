package bot

import (
	"log"
)

// Prints Chat messages to the console
func (bot *Bot) RegisterReadLogger() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			// Prefer IRC client tracing instead
			if evt.IsChatEvent() {
				return
			}

			log.Printf("%+v", evt)
		}),
	)
}
