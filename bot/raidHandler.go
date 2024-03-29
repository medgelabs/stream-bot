package bot

import (
	"fmt"
	"medgebot/bot/viewer"
	log "medgebot/logger"
	"time"
)

// RegisterRaidHandler registers the Raid Auto-Thank feature with the Bot
func (bot *Bot) RegisterRaidHandler(messageTemplate HandlerTemplate, delaySeconds int) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsRaidEvent() {
				log.Info(fmt.Sprintf("%s is raiding with %d raiders!", evt.Sender, evt.Amount))

				if delaySeconds != 0 {
					time.Sleep(time.Duration(delaySeconds) * time.Second)
				}

				bot.SendMessage(messageTemplate.Parse(evt))

				metric := viewer.Metric{
					Name:   evt.Sender,
					Amount: evt.Amount,
				}
				bot.dataStore.Put(viewer.LastRaider, metric.String())
			}
		}),
	)
}
