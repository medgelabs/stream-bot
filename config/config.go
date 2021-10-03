package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

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

// ChannelID returns the numeric ChannelID for the current broadcaster
func (c *Config) ChannelID() string {
	channelID := c.config.GetString(c.key("channelId"))
	return channelID
}

// Nick returns the nickname to join IRC with
func (c *Config) Nick() string {
	nick := c.config.GetString(c.key("nick"))
	return nick
}

// Cache returns the desired Cache type, which should match the cache/cache.go enum
func (c *Config) Cache() string {
	val := c.config.GetString(c.key("cacheType"))
	return val
}

// SecretStore returns the desired secret store type, which should match the
// secret/storeFactory.go enum
func (c *Config) SecretStore() string {
	secretStore := c.config.GetString(c.key("secretStore"))
	return secretStore
}

// TwitchToken if Store type is ENV
func (c *Config) TwitchToken() string {
	return os.Getenv("TWITCH_TOKEN")
}

// Feature Flags - built as opt-in

// GreeterEnabled checks the Greeter feature flag
func (c *Config) GreeterEnabled() bool {
	flagValue := c.config.GetBool(c.key("greeter.enabled"))
	return flagValue
}

// CacheExpirationTime grabs the expiration time for the Greeter Cache
func (c *Config) CacheExpirationTime() int64 {
	value := c.config.GetInt64(c.key("greeter.cache.expirationTime"))
	return value
}

// GreetMessageFormat returns the text/template formatted String for Greet messages
func (c *Config) GreetMessageFormat() string {
	msgFormat := c.config.GetString(c.key("greeter.messageFormat"))
	return msgFormat
}

// ChannelPointsEnabled checks the Commands feature flag
func (c *Config) ChannelPointsEnabled() bool {
	flagValue := c.config.GetBool(c.key("channelPoints.enabled"))
	return flagValue
}

// CommandsEnabled checks the Commands feature flag
func (c *Config) CommandsEnabled() bool {
	flagValue := c.config.GetBool(c.key("commands.enabled"))
	return flagValue
}

// KnownCommand returns a slice of map[prefix]message pairs, to be parsed elsewhere,
// that represent commands the Bot responds to
type KnownCommand struct {
	Prefix   string `mapstructure:"prefix"`
	Message  string `mapstructure:"message"`
	AliasFor string `mapstructure:"aliasFor"`
}

// KnownCommands returns a slice of commands the Bot knows how to respond to
func (c *Config) KnownCommands() []KnownCommand {
	var commands []KnownCommand
	c.config.UnmarshalKey(c.key("commands.known"), &commands)
	return commands
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

// PollsEnabled checks the Subs feature flag
func (c *Config) PollsEnabled() bool {
	flagValue := c.config.GetBool(c.key("polls.enabled"))
	return flagValue
}

// key constructs valid channel-based config keys for Viper lookups
func (c *Config) key(path string) string {
	return fmt.Sprintf("%s.%s", c.channel, path)
}
