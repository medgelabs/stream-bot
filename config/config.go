package config

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// TODO required key functionality?

// Config serves as the centralized place for configuration coming
// from various sources, i.e config.yaml, ENV vars, etc
type Config struct {
	channel string
	config  *viper.Viper
}

// Initialize configuration and read from config.yaml
func New(channel string, configPath string) (Config, error) {
	conf := viper.New()
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(configPath)

	if channel == "" {
		return Config{}, errors.New("ERROR: channel cannot be empty in config.New()")
	}
	if err := conf.ReadInConfig(); err != nil {
		return Config{}, errors.Errorf("FATAL: read %s/config.yaml - %v", configPath, err)
	}

	return Config{
		channel: channel,
		config:  conf,
	}, nil
}

// Nick returns the nickname to join IRC with
func (c *Config) Nick() string {
	nick := c.config.GetString(c.key("nick"))
	if nick == "" {
		log.Fatalf("FATAL: config key - nick not found / empty")
	}

	return nick
}

// Ledger returns the desired ledger type, which should match the ledger/ledgerFactory.go enum
func (c *Config) Ledger() string {
	ledger := c.config.GetString(c.key("ledger"))
	if ledger == "" {
		log.Fatalf("FATAL: config key - ledger not found / empty")
	}

	return ledger
}

// SecretStore returns the desired secret store type, which should match the
// secret/storeFactory.go enum
func (c *Config) SecretStore() string {
	secretStore := c.config.GetString(c.key("secretStore"))
	if secretStore == "" {
		log.Fatalf("FATAL: config key - secretStore not found / empty")
	}

	return secretStore
}

// Feature Flags - built as opt-in

// GreeterEnabled checks the Greeter feature flag
func (c *Config) GreeterEnabled() bool {
	flagValue := c.config.GetBool(c.key("greeter.enabled"))
	return flagValue
}

// GreeterExpirationTime grabs the expiration time for the Greeter ledger
func (c *Config) GreeterExpirationTime() int64 {
	value := c.config.GetInt64(c.key("greeter.ledger.expirationTime"))
	return value
}

// GreetMessageFormat returns the text/template formatted String for Greet messages
func (c *Config) GreetMessageFormat() string {
	msgFormat := c.config.GetString(c.key("greeter.messageFormat"))
	return msgFormat
}

// CommandsEnabled checks the Commands feature flag
func (c *Config) CommandsEnabled() bool {
	flagValue := c.config.GetBool(c.key("commands.enabled"))
	return flagValue
}

// RaidsEnabled checks the Raids feature flag
func (c *Config) RaidsEnabled() bool {
	flagValue := c.config.GetBool(c.key("raids.enabled"))
	return flagValue
}

// RaidDelay returns the delay in seconds between a Raid and the Raid Message being sent
func (c *Config) RaidDelay() int {
	delay := c.config.GetInt(c.key("raids.delaySeconds"))
	return delay
}

// RaidsMessageFormat returns the text/template formatted String for Raids messages
func (c *Config) RaidsMessageFormat() string {
	msgFormat := c.config.GetString(c.key("raids.messageFormat"))
	return msgFormat
}

// BitsEnabled checks the Bits feature flag
func (c *Config) BitsEnabled() bool {
	flagValue := c.config.GetBool(c.key("bits.enabled"))
	return flagValue
}

// BitsMessageFormat returns the text/template formatted String for Bits messages
func (c *Config) BitsMessageFormat() string {
	msgFormat := c.config.GetString(c.key("bits.messageFormat"))
	return msgFormat
}

// SubsEnabled checks the Subs feature flag
func (c *Config) SubsEnabled() bool {
	flagValue := c.config.GetBool(c.key("subs.enabled"))
	return flagValue
}

// SubsMessageFormat returns the text/template formatted String for Subs messages
func (c *Config) SubsMessageFormat() string {
	msgFormat := c.config.GetString(c.key("subs.messageFormat"))
	return msgFormat
}

// GiftSubsMessageFormat returns the text/template formatted String for Gift Subs messages
func (c *Config) GiftSubsMessageFormat() string {
	msgFormat := c.config.GetString(c.key("giftsubs.messageFormat"))
	return msgFormat
}

// key constructs valid channel-based config keys for Viper lookups
func (c *Config) key(path string) string {
	return fmt.Sprintf("%s.%s", c.channel, path)
}
