package bot

import (
	log "medgebot/logger"
)

// RegisterChannelPointHandler responds to Channel Point redemption messages
func (bot *Bot) RegisterChannelPointHandler() {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			log.Info("%+v", evt)
		}),
	)
}
