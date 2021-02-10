package pubsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"medgebot/bot"
	log "medgebot/logger"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

const (
	// MaxMessageSize defines the maximum size of a message that can be received.
	// This is used to size the message buffer slice for the io.Read method
	MaxMessageSize = 2048 // bytes

	// PingInterval for the Ping/Pong loop
	PingInterval = 4 * time.Minute

	ChannelPointTopic = "channel-points-channel-v1"
)

// PubSub wraps the websocket connection to Twitch PubSub and acts as a
// Producer of ChannelPoint messages. We do NOT send messages to PubSub
type PubSub struct {
	sync.Mutex
	conn           io.ReadWriteCloser
	outboundEvents chan<- bot.Event

	// For reconnect purposes
	channelID    string
	authToken    string
	serverScheme string
	serverHost   string
}

// NewClient creates a non-connected Client
func NewClient(conn io.ReadWriteCloser, channelID, authToken string) *PubSub {
	return &PubSub{
		conn:      conn,
		channelID: channelID,
		authToken: authToken,
	}
}

// Start actually connects to PubSub and starts the read loops
func (client *PubSub) Start() error {
	// Read loop for receiving messages
	go func() {
		for {
			if err := client.read(); err != nil {
				log.Error("pubsub read", err)
				break
			}
		}
	}()

	// Subscribe to ChannelPointRedemption messages
	if err := client.listen("channel-points-channel-v1"); err != nil {
		return err
	}

	// Lastly, ensure we heartbeat with PubSub
	go client.ping()

	return nil
}

// Listen fires off a request to listen to the given topics
func (client *PubSub) listen(topic string) error {
	type topicRequest struct {
		Type  string `json:"type"`
		Nonce string `json:"nonce"`
		Data  struct {
			Topics    []string `json:"topics"`
			AuthToken string   `json:"auth_token"`
		} `json:"data"`
	}

	topicStr := fmt.Sprintf("%s.%s", topic, client.channelID)
	log.Info("PubSub: LISTEN " + topicStr)
	return client.write(topicRequest{
		Type: "LISTEN",
		Data: struct {
			Topics    []string `json:"topics"`
			AuthToken string   `json:"auth_token"`
		}{
			Topics:    []string{topicStr},
			AuthToken: client.authToken,
		},
	})
}

// Ping sends a PING message every 4 minutes
func (client *PubSub) ping() {
	ticker := time.NewTicker(PingInterval)

	go func() {
		for range ticker.C {
			ping := struct {
				Type string `json:"type"`
			}{
				Type: "PING",
			}

			log.Info("PubSub: PING")
			if err := client.write(ping); err != nil {
				log.Error("pubsub ping", err)
			}
		}
	}()
}

// Read reads from the PubSub stream
func (client *PubSub) read() error {
	buff := make([]byte, MaxMessageSize)
	len, err := client.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "read pubsub")
	}

	if len == 0 {
		log.Warn("Empty message buffer from PubSub")
		return errors.New("Empty message buffer")
	}

	str := string(buff)
	log.Info("PubSub: " + str)

	// First, we extract the type to see what we're receiving
	msgType, err := jsonparser.GetString(buff, "type")
	if err != nil {
		log.Warn("type field not found. Got message: " + str)
		return errors.Wrap(err, "pubsub unknown message format")
	}

	switch msgType {
	case "RESPONSE":
		// Ignore as this is a response to the LISTEN command

	case "PONG":
		// Ignore, as this is handled by the PingPong goroutine

	case "MESSAGE":
		topic, err := jsonparser.GetString(buff, "data", "topic")
		if err != nil {
			msg := "pubsub data.topic not found in: " + str
			log.Error(msg, err)
			return errors.Wrap(err, msg)
		}

		messageJSON, err := jsonparser.GetString(buff, "data", "message")
		if err != nil {
			msg := "pubsub data.message not found in: " + str
			log.Error(msg, err)
			return errors.Wrap(err, msg)
		}

		// Parse topic for the appropriate Unmarshall call
		// ChannelPoints
		if strings.HasPrefix(topic, ChannelPointTopic) {
			var channelPoints ChannelPointRedemption
			err = json.Unmarshal([]byte(messageJSON), &channelPoints)
			if err != nil {
				msg := "pubsub data.message invalid: " + str
				log.Error(msg, err)
				return errors.Wrap(err, msg)
			}

			client.handleChannelPointRedemption(channelPoints)
		}

	default:
		log.Warn("Unknown message. Skipping: " + str)
	}

	return nil
}

func (client *PubSub) handleChannelPointRedemption(msg ChannelPointRedemption) {
	evt := bot.NewPointsEvent()
	evt.Title = msg.Data.Redemption.Reward.Title
	evt.Sender = msg.Data.Redemption.User.DisplayName
	evt.Amount = msg.Data.Redemption.Reward.Cost
	evt.Message = msg.Data.Redemption.UserInput

	client.outboundEvents <- evt
}

// Write writes a message to the PubSub stream
func (client *PubSub) write(message interface{}) error {
	if client.conn == nil {
		return errors.Errorf("PubSub.conn is nil")
	}

	// Lock since we must only allow one concurrent write
	client.Mutex.Lock()
	defer client.Mutex.Unlock()

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(message)
	_, err := client.conn.Write(buf.Bytes())
	return err
}

// bot.Client

// SetDestination sets the outbound channel for bot.Events the client will send to the bot
func (client *PubSub) SetDestination(outbound chan<- bot.Event) {
	client.outboundEvents = outbound
}
