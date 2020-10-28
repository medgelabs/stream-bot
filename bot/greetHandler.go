package bot

import (
	"log"
	"medgebot/greeter"
	"medgebot/irc"
	"strings"
	"time"
)

func (bot *Bot) RegisterGreeter(greeter greeter.Greeter) {
	bot.RegisterHandler(func(msg irc.Message) {
		username := msg.User
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
	})
}
