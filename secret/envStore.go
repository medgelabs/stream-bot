package secret

import (
	"fmt"
)

// EnvStore is an in-memory secrets store that grabs secrets from
// environment variables
type EnvStore struct {
	twitchToken string
	clientID    string
}

// NewEnvStore is an in-memory Secret store driven by ENV variables
func NewEnvStore(twitchToken, clientID string) EnvStore {
	return EnvStore{
		twitchToken: twitchToken,
		clientID:    clientID,
	}
}

// TwitchToken returns the Twitch Access Token for the Bot
func (s EnvStore) TwitchToken() (string, error) {
	if s.twitchToken == "" {
		return "", fmt.Errorf("ERROR: twitch token not in env")
	}

	return s.twitchToken, nil
}

// ClientID returns the Twitch OAuth ClientID used in API calls
func (s EnvStore) ClientID() (string, error) {
	if s.clientID == "" {
		return "", fmt.Errorf("ERROR: client ID not in env")
	}

	return s.clientID, nil
}
