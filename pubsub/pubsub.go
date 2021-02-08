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

	// TODO determine if Ping/Pong or other before parsing
	var envelope struct {
		Type string `json:"type"`
		Data struct {
			Topic   string          `json:"topic"`
			Message json.RawMessage // to delay processing
		} `json:"data"`
	}

	err = json.Unmarshal(buff, &envelope)
	if err != nil {
		return errors.Wrap(err, "pubsub message parse")
	}

	// ChannelPoints
	if strings.HasPrefix(envelope.Data.Topic, ChannelPointTopic) {
		var channelPoints ChannelPointRedemption
		err = json.Unmarshal(envelope.Data.Message, &channelPoints)

		client.handleChannelPointRedemption(channelPoints)
	}

	return nil
}

func (client *PubSub) handleChannelPointRedemption(msg ChannelPointRedemption) {
	log.Info(fmt.Sprintf("%+v", msg))
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
