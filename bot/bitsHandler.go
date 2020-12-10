package bot

import (
	"log"
)

func (bot *Bot) RegisterBitsHandler(messageFormat string) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				// TODO Thank the sender
				// bot.SendMessage(fmt.Sprintf(messageFormat, evt.Sender))
				log.Printf("> %s cheered %d bits!", evt.Sender, evt.Amount)
			}
		}),
	)
}
