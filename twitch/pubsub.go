package twitch

import (
	"encoding/json"
	"fmt"
	"log"
	"medgebot/bot"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Number of times to retry reconnects
	MAX_RETRIES = 5

	CHANNEL_POINT_TOPIC = "channel-points-channel-v1"
	SUBS_TOPIC          = "channel-subscribe-events-v1"
	BITS_TOPIC          = "channel-bits-events-v2"
)

type PubSubClient struct {
	sync.Mutex
	conn *websocket.Conn

	// To account for reconnects, we need to store connection details
	channelId      string
	authToken      string
	serverHost     string
	serverScheme   string
	outboundEvents chan<- bot.Event
}

// event from a listened topic
type event struct {
	Type  string `json:"type"`
	Error string `json:"error,omitempty"`
	Data  struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	} `json:"data,omitempty"`
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
	topicStr := fmt.Sprintf("%s.%s", CHANNEL_POINT_TOPIC, client.channelId)
	return client.listen(topicStr)
}

// Start the listener as an infinite loop
func (client *PubSubClient) Start() {
	go client.sendPing()

	for {
		evt, err := client.read()
		if err != nil {
			log.Printf("ERROR: pubsub read - %v", err)
			break
		}

		// Now, we figure out what the message is and, if valid, parse and
		// send to the outbound consumers
		if strings.HasPrefix(evt.Data.Topic, CHANNEL_POINT_TOPIC) {
			botEvt, err := parseChannelPointRedeemV1(evt)
			if err != nil {
				log.Printf("ERROR: pubsub parse - %v", err)
				continue
			}
			client.outboundEvents <- botEvt
		}
	}
}

// Attempt to parse a Point Redemption V1 event
func parseChannelPointRedeemV1(message event) (bot.Event, error) {
	var redeem struct {
		Type string `json:"type"`
		Data struct {
			Redemption struct {
				User struct {
					Name string `json:"display_name"`
				} `json:"user"`
				Reward struct {
					Title     string `json:"title"`
					Cost      int    `json:"cost"`
					UserInput string `json:"user_input"`
				} `json:"reward"`
			} `json:"redemption"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(message.Data.Message), &redeem); err != nil {
		log.Printf("PointRedemption parse json - %v", err)
		return bot.Event{}, errors.Errorf("PointRedemption parse failed - %v", err)
	}

	evt := bot.NewPointsEvent()
	data := redeem.Data.Redemption
	evt.Amount = data.Reward.Cost
	evt.Sender = data.User.Name
	evt.Title = data.Reward.Title
	evt.Message = data.Reward.UserInput

	return evt, nil
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
	type listenEvent struct {
		Type string `json:"type"`
		Data struct {
			Topics    []string `json:"topics"`
			AuthToken string   `json:"auth_token"`
		}
	}

	return client.write(listenEvent{
		Type: "LISTEN",
		Data: struct {
			Topics    []string `json:"topics"`
			AuthToken string   `json:"auth_token"`
		}{
			Topics:    []string{topic},
			AuthToken: client.authToken,
		},
	})
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
				log.Printf("ERROR: pubsub PING failed - %v", err)
				break
			}
		}
	}
}

// Read reads from the PubSub stream, one event at a time
func (client *PubSubClient) read() (event, error) {
	_, message, err := client.conn.ReadMessage()
	if err != nil {
		// TODO should check if conn is closed
		if recErr := client.reconnect(); recErr != nil {
			return event{}, err
		}
	}

	var evt event
	if err := json.Unmarshal([]byte(message), &evt); err != nil {
		return event{}, err
	}

	return evt, nil
}

// Write writes a message to the PubSub stream
func (client *PubSubClient) write(message interface{}) error {
	if client.conn == nil {
		return fmt.Errorf("PubSub.conn is nil. Did you forget to call PubSub.Connect()?")
	}

	// Lock since we must only allow one concurrent write
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	if err := client.conn.WriteJSON(message); err != nil {
		// TODO should check if conn is closed
		if recErr := client.reconnect(); recErr != nil {
			return err
		}
	}

	return nil
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
	client.outboundEvents = outbound
}
