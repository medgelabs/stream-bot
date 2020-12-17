package bot

import (
	"log"
	"medgebot/greeter"
	"strings"
	"time"
)

func (bot *Bot) RegisterGreeter(greeter greeter.Greeter, messageTemplate HandlerTemplate) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			username := strings.ToLower(evt.Sender)
			if strings.TrimSpace(username) == "" {
				log.Printf("Empty username for: %+v", evt)
				return
			}

			if greeter.HasNotGreeted(username) {
				log.Printf("Never seen %s before", username)
				time.Sleep(3 * time.Second)

				bot.SendMessage(messageTemplate.Parse(evt))
				greeter.RecordGreeting(username)
			}
		}),
	)
}
