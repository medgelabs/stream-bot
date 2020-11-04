package secret

import (
	"fmt"
	"os"
	"strings"
)

const (
	VAULT = "vault"
	ENV   = "env"
)

func NewSecretStore(storeType string) (Store, error) {
	switch strings.ToLower(storeType) {
	case VAULT:
		vaultUrl := os.Getenv("VAULT_ADDR")
		vaultToken := os.Getenv("VAULT_TOKEN")
		store := NewVaultStore("secret/data/twitchToken")
		err := store.Connect(vaultUrl, vaultToken)
		return &store, err
	case ENV:
		store := NewEnvStore()
		return &store, nil
	default:
		return nil, fmt.Errorf("Invalid storeType " + storeType)
	}
}
