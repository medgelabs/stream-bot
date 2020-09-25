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

	bot.ReadStreamFunc(func(msg string) {
		log.Printf("> %s", msg)
	})

	// Read goroutine for the main chat stream
	// go func() {
	// // PING / PONG must be honored...or we get disconnected
	// if msg.Command == "PING" {
	// pong := irc.Message{
	// Command: "PONG",
	// Params:  msg.Params,
	// }

	// if err := client.Write(pong); err != nil {
	// log.Printf("ERROR: send PONG failed: %s", err)
	// }
	// }

	// // Otherwise - handle PRIVMSG
	// if msg.Command == "PRIVMSG" {
	// channel := msg.Params[0]
	// contents := strings.TrimPrefix(strings.Join(msg.Params[1:], " "), ":")

	// // Command processing
	// // TODO make better
	// if strings.HasPrefix(contents, "!hello") {
	// if err := client.PrivMsg(channel, "WORLD"); err != nil {
	// // if err := client.Write(msg); err != nil {
	// // if err := client.Write("PRIVMSG " + channel + " :WORLD!"); err != nil {
	// log.Printf("ERROR: send failed: %s", err)
	// }
	// }

	// if strings.HasPrefix(contents, "!sorcery") {
	// if err := client.PrivMsg(channel, "!so @SorceryAndSarcasm"); err != nil {
	// log.Printf("ERROR: send failed: %s", err)
	// }
	// }
	// }
	// }()

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
