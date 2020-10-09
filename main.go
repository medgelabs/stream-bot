package main

import (
	"log"
	"medgebot/bot"
	"os"
	"time"
)

func main() {
	channel := "#medgelabs"
	nick := "medgelabs"
	password := os.Getenv("TWITCH_TOKEN")

	chatBot := bot.New()
	chatBot.RegisterPingPong()
	chatBot.RegisterReadLogger()
	chatBot.HandleCommands()

	ledger, err := bot.NewRedisLedger("127.0.0.1", "6379")
	if err != nil {
		log.Fatalf("FATAL - connect to Redis - %s", err)
	}

	// pre-seed names we don't want greeted
	ledger.Add("tmi.twitch.tv")
	ledger.Add("streamlabs")
	ledger.Add(nick)
	ledger.Add(nick + "@tmi.twitch.tv")

	chatBot.RegisterGreeter(&ledger)

	if err := chatBot.Connect(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}
	defer chatBot.Close()

	if err := chatBot.Authenticate(nick, password); err != nil {
		log.Fatalf("FATAL: bot authentication failure - %s", err)
	}

	if err := chatBot.Join(channel); err != nil {
		log.Fatalf("FATAL: bot join channel failed: %s", err)
	}

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
