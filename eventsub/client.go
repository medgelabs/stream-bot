package eventsub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"medgebot/logger"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	// JSON is a constant for the JSON ContentType header value
	JSON = "application/json"

	// ChallengeKey is the cache key for the challenge string received on Start
	ChallengeKey = "challengeKey"
)

// Client is a client to the Twitch Client API that also handles message parsing for events
type Client struct {
	client *http.Client
	secret string
	Config
}

// Config is a input struct for the New() constructor
// serverURL is the URL to the Twitch EventSub API
// callbackURL is the local URL for the Bot which Twitch will send events to
type Config struct {
	ServerURL     string
	CallbackURL   string
	ClientID      string
	AccessToken   string
	BroadcasterID string
}

// New constructs a new Client to EventSub
func New(config Config) Client {
	secret := uuid.NewString()

	return Client{
		client: http.DefaultClient,
		secret: secret,
		Config: config,
	}
}

// SubscriptionRequest represents a request to Twitch to create an EventSub subscription
type SubscriptionRequest struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
}

// Condition is the condition block of a SubscriptionRequest
type Condition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

// Transport is the transport block of a SubscriptionRequest
type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

// Start creates subscriptions to the desired events
func (c *Client) Start() error {
	if c.ServerURL == "" {
		return errors.Errorf("serverURL must not be empty")
	}
	if c.CallbackURL == "" {
		return errors.Errorf("callbackURL must not be empty")
	}

	jsonBody, err := json.Marshal(SubscriptionRequest{
		Type:    "channel.channel_points_custom_reward_redemption.add",
		Version: "1",
		Condition: Condition{
			BroadcasterUserID: c.BroadcasterID,
		},
		Transport: Transport{
			Method:   "webhook",
			Callback: c.CallbackURL,
			Secret:   c.secret,
		},
	})
	if err != nil {
		return errors.Wrap(err, "Marshal subscription request failed")
	}

	url := fmt.Sprintf("%s/subscriptions", c.ServerURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return errors.Wrap(err, "POST to create subscriptions failed")
	}
	req.Header.Add("Client-ID", c.ClientID)
	req.Header.Add("Authorization", "Bearer "+c.AccessToken)
	req.Header.Add("Content-Type", JSON)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "POST to create subscriptions failed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Read subscriptions response body")
	}
	defer resp.Body.Close()

	challenge, err := jsonparser.GetString(body, "challenge")
	if err != nil {
		return errors.Wrap(err, "Get challenge key from body")
	}

	logger.Info("EventSub challenge key received: %s", challenge)

	return nil
}

// Secret returns the generated secret used for Subscription initialization
func (c *Client) Secret() string {
	return c.secret
}

// TODO message parsing
