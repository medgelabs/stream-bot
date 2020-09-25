package main

import (
	"log"
	"medgebot/irc"
	"os"
	"strings"
	"time"
)

func main() {
	channel := "#medgelabs"
	nick := "medgelabs"
	token := os.Getenv("TWITCH_TOKEN")

	client := irc.NewClient()
	if err := client.Connect("wss", "irc-ws.chat.twitch.tv:443"); err != nil {
		log.Fatalf("FATAL: connect - %s", err)
	}
	defer client.Close()

	// Authenticate with the IRC
	if err := client.SendPass(token); err != nil {
		log.Fatalf("FATAL: send PASS failed: %s", err)
	}
	log.Println("< PASS ***")

	if err := client.SendNick(nick); err != nil {
		log.Fatalf("FATAL: send NICK failed: %s", err)
	}

	// Read goroutine for the main chat stream
	go func() {
		for {
			msg, err := client.Read()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
				break
			}
			log.Printf("> %s", msg.String())

			// PING / PONG must be honored...or we get disconnected
			if msg.Command == "PING" {
				pong := irc.Message{
					Command: "PONG",
					Params:  msg.Params,
				}

				if err := client.Write(pong); err != nil {
					log.Printf("ERROR: send PONG failed: %s", err)
				}
			}

			// Otherwise - handle PRIVMSG
			if msg.Command == "PRIVMSG" {
				channel := msg.Params[0]
				contents := strings.TrimPrefix(strings.Join(msg.Params[1:], " "), ":")

				// Command processing
				// TODO make better
				if strings.HasPrefix(contents, "!hello") {
					if err := client.PrivMsg(channel, "WORLD"); err != nil {
						// if err := client.Write(msg); err != nil {
						// if err := client.Write("PRIVMSG " + channel + " :WORLD!"); err != nil {
						log.Printf("ERROR: send failed: %s", err)
					}
				}

				if strings.HasPrefix(contents, "!sorcery") {
					if err := client.PrivMsg(channel, "!so @SorceryAndSarcasm"); err != nil {
						log.Printf("ERROR: send failed: %s", err)
					}
				}
			}
		}
	}()

	time.Sleep(time.Second)

	if err := client.Join(channel); err != nil {
		log.Fatalf("FATAL: send JOIN failed: %s", err)
	}

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
