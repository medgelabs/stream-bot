package bot

import (
	"log"
)

func (bot *Bot) RegisterRaidHandler(messageFormat string) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsRaidEvent() {
				// TODO shout out the raider
				// bot.SendMessage(fmt.Sprintf(messageFormat, evt.Sender))
				log.Printf("> Raid of %d from %s!", evt.Amount, evt.Sender)
			}
		}),
	)
}
