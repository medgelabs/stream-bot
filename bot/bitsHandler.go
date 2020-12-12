package bot

import (
	"log"
)

func (bot *Bot) RegisterBitsHandler(messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Printf("> %s cheered %d bits!", evt.Sender, evt.Amount)
				bot.SendMessage(messageTemplate.Parse(evt))
			}
		}),
	)
}
