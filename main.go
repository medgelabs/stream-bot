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

type Irc struct {
	conn *websocket.Conn
}

func NewIrc() *Irc {
	return &Irc{
		conn: nil,
	}
}

func (irc *Irc) Connect(scheme, server string) error {
	u := url.URL{Scheme: scheme, Host: server, Path: "/"}
	log.Println("connecting to " + u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	irc.conn = conn
	return nil
}

func (irc *Irc) Read() (string, error) {
	_, message, err := irc.conn.ReadMessage()
	if err != nil {
		// TODO check if conn is open. If not - reconnect?
		return "", err
	}
	msgStr := string(message)
	log.Printf("> %s", msgStr)
	return msgStr, nil
}

func (irc *Irc) Write(message string) error {
	if irc.conn == nil {
		return fmt.Errorf("Irc.conn is nil. Did you forget to call Irc.Connect()?")
	}

	if err := irc.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		return err
	}

	if !strings.HasPrefix(message, "PASS") {
		log.Printf("< %s", message)
	}
	return nil
}

func (irc *Irc) Close() {
	irc.conn.Close()
	log.Println("INFO: connection closed")
}

func main() {
	irc := NewIrc()
	if err := irc.Connect("wss", "irc-ws.chat.twitch.tv:443"); err != nil {
		log.Fatalf("FATAL: connect - %s", err)
	}
	defer irc.Close()

	token := os.Getenv("TWITCH_TOKEN")
	passCmd := fmt.Sprintf("PASS %s", token)
	if err := irc.Write(passCmd); err != nil {
		log.Fatalf("FATAL: send PASS failed: %s", err)
	}
	log.Println("< PASS ***")

	nickCmd := fmt.Sprintf("NICK medgelabs")
	if err := irc.Write(nickCmd); err != nil {
		log.Fatalf("FATAL: send NICK failed: %s", err)
	}

	go func() {
		for {
			msgStr, err := irc.Read()
			if err != nil {
				log.Println("ERROR: read from connection - " + err.Error())
				break
			}

			// PING / PONG must be honored...or we get YEETd
			// PING :tmi.twitch.tv -> []string{PING, :tmi.twitch.tv}
			tokens := strings.Split(msgStr, " ")
			switch tokens[0] {
			case "PING":
				if err := irc.Write("PONG " + tokens[1]); err != nil {
					log.Printf("ERROR: send PONG failed: %s", err)
				}
			}
		}
	}()

	time.Sleep(time.Second * 2)
	joinCmd := fmt.Sprintf("JOIN #medgelabs")
	// TODO parameterize
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
