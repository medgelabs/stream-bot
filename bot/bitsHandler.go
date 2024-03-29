package bot

import (
	"fmt"
	"medgebot/bot/viewer"
	log "medgebot/logger"
)

// RegisterBitsHandler adds the Bits handler logic to the Bot
func (bot *Bot) RegisterBitsHandler(messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Info(fmt.Sprintf("> %s cheered %d bits!", evt.Sender, evt.Amount))
				bot.SendMessage(messageTemplate.Parse(evt))

				metric := viewer.Metric{
					Name:   evt.Sender,
					Amount: evt.Amount,
				}
				bot.dataStore.Put(viewer.LastBits, metric.String())
			}
		}),
	)
}
