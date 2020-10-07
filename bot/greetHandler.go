package bot

import (
	"fmt"
	"log"
	"medgebot/irc"
	"time"
)

func (bot *Bot) RegisterGreeter( /*ledger Ledger */ ) {
	initLedger()

	bot.RegisterHandler(func(msg irc.Message) {
		username := msg.User
		if absent(username) {
			log.Printf("Never seen %s before", username)
			add(username)

			msg := fmt.Sprintf("Welcome %s!", username)
			time.Sleep(2 * time.Second)
			bot.SendMessage(msg)
		}
	})
}

var ledger map[string]int

func initLedger() {
	ledger = make(map[string]int)
}

func absent(key string) bool {
	_, ok := ledger[key]
	return !ok
}

func add(key string) error {
	ledger[key] = 1
	return nil
}
