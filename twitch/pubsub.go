package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"medgebot/bot"
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
	channelId      string
	authToken      string
	serverHost     string
	serverScheme   string
	outboundevents chan<- bot.Event
}

// event from a listened topic on PubSub
type event struct {
	Type  string      `json:"type"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func NewPubSubClient(channelId, authToken string) *PubSubClient {
	return &PubSubClient{
		conn:      nil,
		channelId: channelId,
		authToken: authToken,
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

	// LISTEN to desired topics
	topicStr := fmt.Sprintf("%s.%s", "channel-points-channel-v1", client.channelId)
	return client.listen(topicStr)
}

// Start the listener as an infinite loop
func (client *PubSubClient) Start() {
	go client.sendPing()

	for {
		_, err := client.read()
		if err != nil {
			log.Printf("ERROR: pubsub read - %v", err)
			break
		}
	}
}

func (client *PubSubClient) Close() {
	if client.conn == nil {
		log.Println("WARN: PubSub.Close() called on nil connection")
		return
	}

	client.conn.Close()
}

// Listen to the given topic
func (client *PubSubClient) listen(topic string) error {
	evt := event{
		Type: "LISTEN",
		Data: struct {
			Topics    []string `json:"topics"`
			AuthToken string   `json:"auth_token"`
		}{
			[]string{topic},
			client.authToken,
		},
	}

	return client.write(evt)
}

// Read reads from the PubSub stream, one event at a time
func (client *PubSubClient) read() (event, error) {
	_, message, err := client.conn.ReadMessage()
	if err != nil {
		return event{}, err
	}

	var evt event
	if err := json.Unmarshal([]byte(message), &evt); err != nil {
		return event{}, err
	}
	log.Printf("%+v", evt)

	// TODO Now, we figure out what the message is

	// TODO how to validate we received a PONG in time?

	return event{}, nil
}

// Write writes a message to the PubSub stream
func (client *PubSubClient) write(message event) error {
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
func (client *PubSubClient) sendPing() {
	pingTick := time.NewTicker(4 * time.Minute)
	for {
		select {
		case <-pingTick.C:
			err := client.write(event{
				Type: "PING",
			})
			if err != nil {
				break
			}

		}
	}
}

// Handle the need to reconnect the client
// Maybe it was an error, maybe it was a RECONNECT event
func (client *PubSubClient) reconnect() error {
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
			return err
		}
	}

	return nil
}

func (client *PubSubClient) SetChannel(outbound chan<- bot.Event) {
	client.outboundevents = outbound
}
