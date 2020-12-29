package bot

func (bot *Bot) RegisterSubsHandler(subsTemplate, giftSubsTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsSubEvent() {
				bot.SendMessage(subsTemplate.Parse(evt))
			} else if evt.IsGiftSubEvent() {
				bot.SendMessage(giftSubsTemplate.Parse(evt))
			} else {
				return // no messaging otherwise
			}
		}),
	)
}
