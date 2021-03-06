package secret

import (
	"fmt"
	"medgebot/config"
	"strings"

	"github.com/pkg/errors"
)

const (
	ENV = "env"
)

// NewSecretStore returns a SecretStore for the given type, or an error if an invalid type is passed
func NewSecretStore(config config.Config) (Store, error) {
	storeType := config.SecretStore()
	if storeType == "" {
		return nil, errors.New("config key secretStore not found / empty")
	}

	switch strings.ToLower(storeType) {
	case ENV:
		twitchToken := config.TwitchToken()
		store := NewEnvStore(twitchToken)
		return &store, nil
	default:
		return nil, fmt.Errorf("Invalid storeType - %s. Valid values are: %v", storeType, []string{ENV})
	}
}
