package secret

import (
	"fmt"
	"medgebot/config"
	"strings"

	"github.com/pkg/errors"
)

const (
	VAULT = "vault"
	ENV   = "env"
)

// NewSecretStore returns a SecretStore for the given type, or an error if an invalid type is passed
func NewSecretStore(config config.Config) (Store, error) {
	storeType := config.SecretStore()
	if storeType == "" {
		return nil, errors.New("config key secretStore not found / empty")
	}

	switch strings.ToLower(storeType) {
	case VAULT:
		vaultUrl := config.VaultAddress()
		vaultToken := config.VaultToken()
		store := NewVaultStore("secret/data/twitchToken")
		err := store.Connect(vaultUrl, vaultToken)
		return &store, err
	case ENV:
		twitchToken := config.TwitchToken()
		store := NewEnvStore(twitchToken)
		return &store, nil
	default:
		return nil, fmt.Errorf("Invalid storeType - %s. Valid values are: %s, %s", storeType, VAULT, ENV)
	}
}
