package bot

import (
	"medgebot/bot/viewer"
	"medgebot/cache"
)

// RegisterSubsHandler adds the Subscription handler logic to the Bot
func (bot *Bot) RegisterSubsHandler(subsTemplate, giftSubsTemplate HandlerTemplate, metricsCache cache.Cache) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsSubEvent() {
				bot.SendMessage(subsTemplate.Parse(evt))

				// TODO if evt.isDebug() { return }
				metric := viewer.Metric{
					Name:   evt.Sender,
					Amount: evt.Amount,
				}
				metricsCache.Put(viewer.LastSub, metric.String())
			} else if evt.IsGiftSubEvent() {
				bot.SendMessage(giftSubsTemplate.Parse(evt))

				metric := viewer.Metric{
					Name:      evt.Sender,
					Recipient: evt.Recipient,
					Amount:    evt.Amount,
				}
				metricsCache.Put(viewer.LastGiftSub, metric.String())
			} else {
				return // no messaging otherwise
			}
		}),
	)
}
