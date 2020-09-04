package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "wss", Host: "irc-ws.chat.twitch.tv:443", Path: "/"}
	log.Println("connecting to " + u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("FATAL: connect failed: %s", err)
	}
	defer conn.Close()

	token := os.Getenv("TWITCH_TOKEN")
	passCmd := fmt.Sprintf("PASS %s", token)
	err = conn.WriteMessage(websocket.TextMessage, []byte(passCmd))
	if err != nil {
		log.Fatalf("FATAL: send PASS failed: %s", err)
	}
	log.Println("< PASS ***")

	nickCmd := fmt.Sprintf("NICK medgelabs")
	err = Write(conn, nickCmd)
	if err != nil {
		log.Fatalf("FATAL: send NICK failed: %s", err)
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
			}
			msgStr := string(message)
			log.Printf("> %s", msgStr)

			// PING / PONG must be honored...or we get YEETd
			// PING :tmi.twitch.tv -> []string{PING, :tmi.twitch.tv}
			tokens := strings.Split(msgStr, " ")
			switch tokens[0] {
			case "PING":
				if err := Write(conn, "PONG "+tokens[1]); err != nil {
					log.Printf("ERROR: send PONG failed: %s", err)
				}
			}

		}
	}()

	time.Sleep(time.Second * 2)
	joinCmd := fmt.Sprintf("JOIN #medgelabs")
	// TODO parameterize
	if err = Write(conn, joinCmd); err != nil {
		log.Fatalf("FATAL: send JOIN failed: %s", err)
	}

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("ERROR: read stdin - %s", err)
			}

			if err = Write(conn, cmd); err != nil {
				log.Fatalf("FATAL: send %s failed: %s", cmd, err)
			}
		}
	}()

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}

func Write(conn *websocket.Conn, message string) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}

	log.Printf("< %s", message)
	return nil
}
