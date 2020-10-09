package bot

import (
	"fmt"
	"log"
	"medgebot/irc"
	"time"
)

func (bot *Bot) RegisterGreeter(ledger Ledger) {
	bot.RegisterHandler(func(msg irc.Message) {
		username := msg.User
		if ledger.Absent(username) {
			log.Printf("Never seen %s before", username)
			ledger.Add(username)

			msg := fmt.Sprintf("Welcome @%s!", username)
			time.Sleep(5 * time.Second)
			bot.SendMessage(msg)
		}
	})
}
