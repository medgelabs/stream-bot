package bot

import (
	"fmt"
	"log"
)

func (bot *Bot) RegisterBitsHandler(messageFormat string) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Printf("> %s cheered %d bits!", evt.Sender, evt.Amount)
				log.Println(fmt.Sprintf(messageFormat, evt.Sender))
				// bot.SendMessage(fmt.Sprintf(messageFormat, evt.Sender))
			}
		}),
	)
}
