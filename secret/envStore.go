package secret

import (
	"fmt"
)

// EnvStore is an in-memory secrets store that grabs secrets from
// environment variables
type EnvStore struct {
	twitchToken string
}

func NewEnvStore(twitchToken string) EnvStore {
	return EnvStore{
		twitchToken: twitchToken,
	}
}

func (s EnvStore) TwitchToken() (string, error) {
	if s.twitchToken == "" {
		return "", fmt.Errorf("ERROR: twitch token not in env")
	}

	return s.twitchToken, nil
}
