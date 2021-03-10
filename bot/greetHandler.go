package bot

import (
	"log"
	"medgebot/cache"
	"strings"
	"time"
)

// RegisterGreeter creates and registers the greeter module with the Bot
func (bot *Bot) RegisterGreeter(cache cache.Cache, messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			username := strings.ToLower(evt.Sender)
			if strings.TrimSpace(username) == "" {
				log.Printf("Empty username for: %+v", evt)
				return
			}

			if cache.Absent(username) {
				log.Printf("Never seen %s before", username)
				time.Sleep(3 * time.Second)

				bot.SendMessage(messageTemplate.Parse(evt))
				cache.Put(username, "")
			}
		}),
	)
}
