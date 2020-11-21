package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Number of times to retry reconnects
	MAX_RETRIES = 5
)

type PubSubClient struct {
	sync.Mutex
	conn *websocket.Conn

	// To account for reconnects, we need to store connection details
	serverHost   string
	serverScheme string
}

// Event from a listened topic on PubSub
type Event struct {
	Type  string `json:"type"`
	Error string `json:"error,omitempty"`
	// Data  string `json:"data,omitempty"`
}

func NewClient() *PubSubClient {
	return &PubSubClient{
		conn: nil,
	}
}

func (client *PubSubClient) Connect(scheme, server string) error {
	client.serverHost = server
	client.serverScheme = scheme

	u := url.URL{Scheme: scheme, Host: server, Path: "/"}
	log.Println("connecting to " + u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	client.conn = conn

	// Kick off PING/PONG handler

	// LISTEN to desired topics

	return nil
}

func (client *PubSubClient) Close() {
	if client.conn == nil {
		log.Println("WARN: PubSub.Close() called on nil connection")
		return
	}

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

// PING handler to keep the connection alive
func (client *PubSubClient) SendPing() error {
	err := client.write(Event{
		Type: "PING",
	})

	// Wait for pong?

	return err
}

// Handle the need to reconnect the client
// Maybe it was an error, maybe it was a RECONNECT event
func (client *PubSubClient) reconnect() {
	log.Println("Reconnecting PubSub client..")

	client.Lock()
	defer client.Unlock()

	// Close the existing connection to ensure no resource leaks!
	client.Close()

	err := client.Connect(client.serverHost, client.serverScheme)
	for retries := 1; err == nil || retries == MAX_RETRIES; retries++ {
		log.Printf("ERROR: reconnect. Retry %d - %v", retries, err)
		time.Sleep(time.Duration(retries) * time.Second)
		err = client.Connect(client.serverHost, client.serverScheme)

		// If max retries reached
		if retries == MAX_RETRIES-1 {
			log.Println("ERROR: Max reconnect tries hit in PubSub.reconnect()")
		}
	}
}
