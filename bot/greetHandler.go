package bot

import (
	"fmt"
	"log"
	"medgebot/internal/pkg/irc"
	"medgebot/internal/pkg/ledger"
	"time"
)

func (bot *Bot) RegisterGreeter(ledger *ledger.Ledger) {
	bot.RegisterHandler(func(msg irc.Message) {
		username := msg.User
		if ledger.Absent(username) {
			log.Printf("Never seen %s before", username)
			err := ledger.Add(username)
			if err != nil {
				fmt.Printf("ERROR: couldn't add username to ledger - %v\n", err)
			}

			msg := fmt.Sprintf("Welcome @%s!", username)
			time.Sleep(5 * time.Second)
			err = bot.SendMessage(msg)
			if err != nil {
				fmt.Printf("ERROR: couldn't send message to user - %v\n", err)
			}
		}
	})
}
