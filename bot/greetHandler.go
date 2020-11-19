package bot

import (
	"log"
	"medgebot/greeter"
	"strings"
	"time"
)

func (bot *Bot) RegisterGreeter(greeter greeter.Greeter) {
	bot.RegisterHandler(
		NewHandler(func(msg Event) {
			username := msg.Sender
			if strings.TrimSpace(username) == "" {
				return
			}

			if greeter.HasNotGreeted(username) {
				log.Printf("Never seen %s before", username)
				time.Sleep(3 * time.Second)

				msg := greeter.Greet(username)
				bot.SendMessage(msg)
				greeter.RecordGreeting(username)
			}
		}),
	)
}
