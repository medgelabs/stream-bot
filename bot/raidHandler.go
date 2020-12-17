package bot

import (
	"log"
	"time"
)

func (bot *Bot) RegisterRaidHandler(messageTemplate HandlerTemplate, delaySeconds int) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsRaidEvent() {
				log.Printf("%s is raiding with %d raiders!", evt.Sender, evt.Amount)

				if delaySeconds != 0 {
					time.Sleep(time.Duration(delaySeconds) * time.Second)
				}

				bot.SendMessage(messageTemplate.Parse(evt))
			}
		}),
	)
}
