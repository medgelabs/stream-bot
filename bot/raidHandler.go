package bot

import "log"

func (bot *Bot) RegisterRaidHandler(messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsRaidEvent() {
				log.Printf("%s is raiding with %d raiders!", evt.Sender, evt.Amount)
				bot.SendMessage(messageTemplate.Parse(evt))
			}
		}),
	)
}
