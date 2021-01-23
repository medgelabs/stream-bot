package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// TODO required key functionality?

// Config serve as the centralized place for configuration coming
// from various sources, i.e config.yaml, ENV vars, etc
type Config struct {
	channel string
	config  *viper.Viper
}

// Initialize configuration and read from config.yaml
func New(channel string, configPath string) Config {
	conf := viper.New()
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(configPath)
	if err := conf.ReadInConfig(); err != nil {
		log.Fatalf("FATAL: read %s/config.yaml - %v", configPath, err)
	}

	return Config{
		channel: channel,
		config:  conf,
	}
}

// Nick returns the nickname to join IRC with
func (c *Config) Nick() string {
	nick := c.config.GetString(c.key("nick"))
	if nick == "" {
		log.Fatalf("FATAL: config key - nick not found / empty")
	}

	return nick
}

// Ledger returns the desired ledger type, which should match the ledger/ledger.go enum
func (c *Config) Ledger() string {
	ledger := c.config.GetString(c.key("ledger"))
	if ledger == "" {
		log.Fatalf("FATAL: config key - ledger not found / empty")
	}

	return ledger
}

// SecretStore returns the desired secret store type, which should match the
// secret/store.go enum
func (c *Config) SecretStore() string {
	secretStore := c.config.GetString(c.key("secretStore"))
	if secretStore == "" {
		log.Fatalf("FATAL: config key - secretStore not found / empty")
	}

	return secretStore
}

// Construct valid channel-based config keys like: channel.feature.option
func (c *Config) key(path string) string {
	return fmt.Sprintf("%s.%s", c.channel, path)
}
