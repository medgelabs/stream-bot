package irc

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type Irc struct {
	conn *websocket.Conn
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

func (irc *Irc) Read() (string, error) {
	_, message, err := irc.conn.ReadMessage()
	if err != nil {
		// TODO check if conn is open. If not - reconnect?
		return "", err
	}
	msgStr := string(message)
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
