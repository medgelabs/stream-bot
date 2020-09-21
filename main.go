package main

import (
	"fmt"
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

	irc := irc.NewClient()
	if err := irc.Connect("wss", "irc-ws.chat.twitch.tv:443"); err != nil {
		log.Fatalf("FATAL: connect - %s", err)
	}
	defer irc.Close()

	passCmd := fmt.Sprintf("PASS %s", token)
	if err := irc.Write(passCmd); err != nil {
		log.Fatalf("FATAL: send PASS failed: %s", err)
	}
	log.Println("< PASS ***")

	nickCmd := fmt.Sprintf("NICK %s", nick)
	if err := irc.Write(nickCmd); err != nil {
		log.Fatalf("FATAL: send NICK failed: %s", err)
	}

	// Read goroutine for the main chat stream
	go func() {
		for {
			msg, err := irc.Read()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
				break
			}
			log.Printf("> %s", msg.String())

			// PING / PONG must be honored...or we get disconnected
			if msg.Command == "PING" {
				if err := irc.Write("PONG " + strings.Join(msg.Params, " ")); err != nil {
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
					if err := irc.Write("PRIVMSG " + channel + " :WORLD!"); err != nil {
						log.Printf("ERROR: send failed: %s", err)
					}
				}

				if strings.HasPrefix(contents, "!sorcery") {
					if err := irc.Write("PRIVMSG " + channel + " :!so @SorceryAndSarcasm"); err != nil {
						log.Printf("ERROR: send failed: %s", err)
					}
				}
			}
		}
	}()

	time.Sleep(time.Second)
	joinCmd := fmt.Sprintf("JOIN %s", channel)
	if err := irc.Write(joinCmd); err != nil {
		log.Fatalf("FATAL: send JOIN failed: %s", err)
	}

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
