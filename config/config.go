package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	conf *viper.Viper
}

func LoadConfig() Config {
	conf := viper.New()

	// Initialize configuration and read from config.yaml
	conf.SetConfigName("config")
	conf.SetConfigType("yaml")
	conf.AddConfigPath(".")
	if err := conf.ReadInConfig(); err != nil {
		log.Fatalf("FATAL: read config.yaml - %v", err)
	}

	return Config{
		conf: conf,
	}
}

// GetString returns a config value, cast as a string. Empty string if not present
func (c *Config) GetString(key string) string {
	if !c.conf.IsSet(key) {
		log.Printf("WARN: key %s not found in config", key)
		return ""
	}

	return c.conf.GetString(key)
}

// GetInt64 returns a config value, cast as an int64. 0 if not present
func (c *Config) GetIntOrDefault(key string, defaultVal int) int {
	if !c.conf.IsSet(key) {
		log.Printf("WARN: key %s not found in config", key)
		return defaultVal
	}

	return c.conf.GetInt(key)
}

// GetInt64 returns a config value, cast as an int64. 0 if not present
func (c *Config) GetInt64(key string) int64 {
	if !c.conf.IsSet(key) {
		log.Printf("WARN: key %s not found in config", key)
		return 0
	}

	return c.conf.GetInt64(key)
}

// GetList returns a config list as a slice. Empty slice if not present
func (c *Config) GetList(key string) []string {
	if !c.conf.IsSet(key) {
		log.Printf("WARN: key %s not found in config", key)
		return []string{}
	}

	return c.conf.GetStringSlice(key)
}

// FeatureEnabled checks the value of featureName.enabled in config. false if not present in config
func (c *Config) FeatureEnabled(featureName string) bool {
	key := fmt.Sprintf("%s.enabled", featureName)
	if !c.conf.IsSet(key) {
		log.Printf("WARN: key %s not found in config", key)
		return false
	}

	return c.conf.GetBool(key)
}
