package storage

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

type Vault struct {
	client     *api.Client
	secretPath string
}

func NewVault(url, token, secretPath string) (*Vault, error) {
	client, err := api.NewClient(&api.Config{
		Address: url,
	})
	if err != nil {
		return nil, err
	}
	client.SetToken(token)

	return &Vault{
		client:     client,
		secretPath: secretPath,
	}, nil
}

func (v Vault) put(key string, value interface{}) error {
	_, err := v.client.Logical().Write(v.secretPath, map[string]interface{}{
		key: value,
	})

	if err != nil {
		return err
	}

	return nil
}

func (v Vault) get(key string) (interface{}, error) {
	secret, err := v.client.Logical().Read(v.secretPath)
	if err != nil {
		return "", fmt.Errorf("fetch value from store - %v", err)
	}
	// This can happen from the underlying api
	// @see github.com/hashicorp/vault/api@v1.0.4/logical.go:81
	if secret == nil {
		return "", fmt.Errorf("ERROR: strange things happening - no error and no ptr to secret")
	}

	dataMap := secret.Data["data"].(map[string]interface{})
	value, ok := dataMap[key]
	if !ok {
		return "", fmt.Errorf("key not found in vault")
	}

	return value, nil
}
