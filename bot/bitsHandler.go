package bot

import (
	"fmt"
	log "medgebot/logger"
)

func (bot *Bot) RegisterBitsHandler(messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Info(fmt.Sprintf("> %s cheered %d bits!", evt.Sender, evt.Amount))
				bot.SendMessage(messageTemplate.Parse(evt))
			}
		}),
	)
}
