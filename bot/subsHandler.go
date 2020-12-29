package bot

func (bot *Bot) RegisterSubsHandler(subsTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsSubEvent() {
				bot.SendMessage(subsTemplate.Parse(evt))
			}
		}),
	)
}

func (bot *Bot) RegisterGiftSubsHandler(giftSubsTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if evt.IsGiftSubEvent() {
				bot.SendMessage(giftSubsTemplate.Parse(evt))
			}
		}),
	)
}
