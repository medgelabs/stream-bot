package main

import (
	"bufio"
	"fmt"
	"log"
	"medgebot/irc"
	"os"
	"strings"
	"time"
)

type Message struct {
	Prefix  string
	Command string
	Params  []string
}

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
			msgStr, err := irc.Read()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
				break
			}
			log.Printf("> %s", msgStr)

			trimmed := strings.TrimSpace(msgStr)
			tokens := strings.Split(trimmed, " ")

			var msg Message
			if strings.HasPrefix(tokens[0], ":") {
				msg = Message{
					// TODO this will break when prefix > 1 token
					// Need to add processing for space-delimited prefix as well
					Prefix:  tokens[0],
					Command: tokens[1],
					Params:  tokens[2:],
				}
			} else {
				msg = Message{
					Prefix:  "",
					Command: tokens[0],
					Params:  tokens[1:], // TODO are there any commands we need to handle that have no params?
				}
			}

			// PING / PONG must be honored...or we get disconnected
			if msg.Command == "PING" {
				if err := irc.Write("PONG " + tokens[1]); err != nil {
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

	time.Sleep(time.Second * 2)
	joinCmd := fmt.Sprintf("JOIN %s", channel)
	if err := irc.Write(joinCmd); err != nil {
		log.Fatalf("FATAL: send JOIN failed: %s", err)
	}

	// Allow sending IRC commands from stdin
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("ERROR: read stdin - %s", err)
			}

			if err := irc.Write(cmd); err != nil {
				log.Fatalf("FATAL: send %s failed: %s", cmd, err)
			}
		}
	}()

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
