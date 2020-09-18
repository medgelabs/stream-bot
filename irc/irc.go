package irc

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Irc client
type Irc struct {
	conn *websocket.Conn
}

// Message represents a line of text from the IRC stream
type Message struct {
	Prefix  string
	Command string
	Params  []string
}

func (msg Message) String() string {
	return fmt.Sprintf("%s %s :%s", msg.Prefix, msg.Command, strings.Join(msg.Params, " "))
}

func NewClient() *Irc {
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

func (irc *Irc) Read() (Message, error) {
	_, message, err := irc.conn.ReadMessage()
	if err != nil {
		// TODO check if conn is open. If not - reconnect?
		return Message{}, err
	}

	// TrimSpace to get rid of /r/n
	msgStr := strings.TrimSpace(string(message))
	tokens := strings.Split(msgStr, " ")

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

	return msg, nil
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
