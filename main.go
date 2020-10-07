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

	bot := bot.New()
	bot.RegisterPingPong()
	bot.RegisterReadLogger()
	bot.HandleCommands()
	bot.RegisterGreeter()

	if err := bot.Connect(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}
	defer bot.Close()

	if err := bot.Authenticate(nick, password); err != nil {
		log.Fatalf("FATAL: bot authentication failure - %s", err)
	}

	if err := bot.Join(channel); err != nil {
		log.Fatalf("FATAL: bot join channel failed: %s", err)
	}

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
