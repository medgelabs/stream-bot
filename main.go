package main

import (
	"fmt"
	"log"
	"medgebot/bot"
	"medgebot/greeter"
	"medgebot/ledger"
	"medgebot/secret"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {

	// Initialize configuration and read from config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("FATAL: read config.yaml - %v", err)
	}

	channel := "#medgelabs"
	nick := "medgelabs"

	// Ledger for the auto greeter
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	ledger, err := ledger.NewRedisLedger(redisHost, redisPort)
	if err != nil {
		log.Fatalf("FATAL - connect to Redis - %s", err)
	}

	// pre-seed names we don't want greeted
	ledger.Add("tmi.twitch.tv")
	ledger.Add("streamlabs")
	ledger.Add("nightbot")
	ledger.Add(nick)
	ledger.Add(strings.TrimPrefix(channel, "#")) // Prevent greeting the broadcaster
	ledger.Add(nick + "@tmi.twitch.tv")

	// Initialize Secrets Store
	vaultUrl := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	store := secret.NewVaultStore("secret/data/twitchToken")
	if err := store.Connect(vaultUrl, vaultToken); err != nil {
		log.Fatalf("FATAL: Vault connect - %v", err)
	}

	password, err := store.GetTwitchToken()
	if err != nil {
		log.Fatalf("FATAL: Get Twitch Token from store - %v", err)
	}

	// Greeter config
	var greetConfig greeter.Config
	confKey := fmt.Sprintf("greeter.%s", strings.Trim(channel, "#"))
	confSub := viper.Sub(confKey)
	if confSub == nil {
		log.Fatalf("FATAL: key %s not found in config", confKey)
	}

	confSub.Unmarshal(&greetConfig)
	greetBot := greeter.New(greetConfig, &ledger)

	/*

		confKey := fmt.Sprintf("greeter.%s", strings.Trim(channel, "#"))
		greetConfig := greeter.Config{
			MessageFormat: viper.GetString(confKey),
		}

		greetBot := greeter.New(greetConfig, &ledger)
	*/

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterPingPong()
	chatBot.RegisterReadLogger()
	chatBot.HandleCommands()
	chatBot.RegisterGreeter(greetBot)

	if err := chatBot.Connect(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}
	defer chatBot.Close()

	if err := chatBot.Authenticate(nick, password); err != nil {
		log.Fatalf("FATAL: bot authentication failure - %s", err)
	}

	if err := chatBot.Join(channel); err != nil {
		log.Fatalf("FATAL: bot join channel failed: %s", err)
	}

	// TODO _no_
	for {
		time.Sleep(time.Second)
	}
}
