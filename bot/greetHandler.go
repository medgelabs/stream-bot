package bot

import (
	"log"
	"medgebot/ledger"
	"strings"
	"time"
)

func (bot *Bot) RegisterGreeter(cache ledger.Cache, messageTemplate HandlerTemplate) {
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
