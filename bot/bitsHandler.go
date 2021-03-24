package bot

import (
	"fmt"
	"medgebot/bot/viewer"
	"medgebot/cache"
	log "medgebot/logger"
)

// RegisterBitsHandler adds the Bits handler logic to the Bot
func (bot *Bot) RegisterBitsHandler(messageTemplate HandlerTemplate, metricsCache cache.Cache) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsBitsEvent() {
				log.Info(fmt.Sprintf("> %s cheered %d bits!", evt.Sender, evt.Amount))
				bot.SendMessage(messageTemplate.Parse(evt))

				// TODO if evt.isDebug() { return }
				metric := viewer.Metric{
					Name:   evt.Sender,
					Amount: evt.Amount,
				}
				metricsCache.Put("lastBits", metric.String())
			}
		}),
	)
}
