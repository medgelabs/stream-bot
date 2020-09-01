package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
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
	passCmd := fmt.Sprintf("PASS oauth:%s", token)
	err = conn.WriteMessage(websocket.TextMessage, []byte(passCmd))
	if err != nil {
		log.Fatalf("FATAL: send PASS failed: %s", err)
	}
	log.Println("< PASS oauth:****")

	// TODO make NICK a variable
	nickCmd := fmt.Sprintf("NICK medgelabs")
	err = conn.WriteMessage(websocket.TextMessage, []byte(nickCmd))
	if err != nil {
		log.Fatalf("FATAL: send NICK failed: %s", err)
	}
	log.Println("< NICK medgelabs")

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
			}
			log.Printf("> %s", message)
		}
	}()

	time.Sleep(time.Second)
	joinCmd := fmt.Sprintf("JOIN #medgelabs")
	// TODO parameterize
	err = conn.WriteMessage(websocket.TextMessage, []byte(joinCmd))
	if err != nil {
		log.Fatalf("FATAL: send JOIN failed: %s", err)
	}
	log.Println("< JOIN #medgelabs")

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
