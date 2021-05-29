package bot

import (
	"medgebot/cache"
	log "medgebot/logger"
	"strings"
	"time"
)

// RegisterGreeter creates and registers the greeter module with the Bot
// Note: this is a different cache from the bot.metricsCache, so we still expect one as an arg
func (bot *Bot) RegisterGreeter(cache cache.Cache, messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			username := strings.ToLower(evt.Sender)
			if strings.TrimSpace(username) == "" {
				log.Info("Empty username for: %+v", evt)
				return
			}

			if cache.Absent(username) {
				log.Info("Never seen %s before", username)
				time.Sleep(3 * time.Second)

				bot.SendMessage(messageTemplate.Parse(evt))
				cache.Put(username, "")
			}
		}),
	)
}
