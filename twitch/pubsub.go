package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type PubSubClient struct {
	sync.Mutex
	conn *websocket.Conn
}

type Event struct {
	Type  string `json:"type"`
	Data  string `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewClient() *PubSubClient {
	return &PubSubClient{
		conn: nil,
	}
}

func (client *PubSubClient) Connect(scheme, server string) error {
	u := url.URL{Scheme: scheme, Host: server, Path: "/"}
	log.Println("connecting to " + u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	client.conn = conn
	return nil
}

func (client *PubSubClient) Close() {
	client.conn.Close()
	log.Println("INFO: connection closed")
}

// Read reads from the PubSub stream, one event at a time
func (client *PubSubClient) Read() (Event, error) {
	_, message, err := client.conn.ReadMessage()
	if err != nil {
		return Event{}, err
	}

	var event Event
	if err := json.Unmarshal([]byte(message), &event); err != nil {
		return Event{}, err
	}

	// TODO Now, we figure out what the message is

	return Event{}, nil
}

func (client *PubSubClient) SendPing() error {
	err := client.write(Event{
		Type: "PING",
	})

	return err
}

// Write writes a message to the PubSub stream
func (client *PubSubClient) write(message Event) error {
	if client.conn == nil {
		return fmt.Errorf("PubSub.conn is nil. Did you forget to call PubSub.Connect()?")
	}

	// Lock since we must only allow one concurrent write
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	if err := client.conn.WriteJSON(message); err != nil {
		return err
	}

	return nil
}
