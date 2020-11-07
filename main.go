package main

import (
	"flag"
	"fmt"
	"log"
	"medgebot/bot"
	"medgebot/greeter"
	"medgebot/ledger"
	"medgebot/secret"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {

	// channel and nick come from the CLI
	var channel string
	var nick string
	flag.StringVar(&channel, "channel", "", "Channel name, without the #, to join")
	flag.StringVar(&nick, "nick", "", "Nickname to join Chat with")

	// Inputs for factories
	var ledgerType string
	flag.StringVar(&ledgerType, "ledger", ledger.REDIS, fmt.Sprintf("Ledger type string. Options: %s, %s, %s", ledger.REDIS, ledger.FILE, ledger.MEM))

	var secretStoreType string
	flag.StringVar(&secretStoreType, "store", secret.VAULT, fmt.Sprintf("Secret Store type string. Options: %s, %s", secret.VAULT, secret.ENV))

	flag.Parse()

	// Flag error handling
	if strings.HasPrefix(channel, "#") {
		log.Fatalln("FATAL: channel cannot start with a #")
	}
	channel = fmt.Sprintf("#%s", channel)

	if nick == "" {
		log.Fatalln("FATAL: nick empty")
	}

	// Initialize configuration and read from config.yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("FATAL: read config.yaml - %v", err)
	}

	// Ledger for the auto greeter
	ledger, err := ledger.NewLedger(ledgerType)
	if err != nil {
		log.Fatalf("FATAL: create ledger - %v", err)
	}

	// pre-seed names we don't want greeted
	// TODO make variadic
	ledger.Add("tmi.twitch.tv")
	ledger.Add("streamlabs")
	ledger.Add("nightbot")
	ledger.Add(nick)
	ledger.Add(strings.TrimPrefix(channel, "#")) // Prevent greeting the broadcaster
	ledger.Add(nick + ".tmi.twitch.tv")
	ledger.Add(nick + "@tmi.twitch.tv")

	// Initialize Secrets Store
	store, err := secret.NewSecretStore(secretStoreType)
	if err != nil {
		log.Fatalf("FATAL: Create secret store - %v", err)
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
	greetBot := greeter.New(greetConfig, ledger)

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterPingPong()
	chatBot.RegisterReadLogger()
	chatBot.HandleCommands()
	chatBot.RegisterGreeter(greetBot)
	// chatBot.RegisterEmoteCounter()
	// chatBot.RegisterRaidHander()
	// chatbot.RegisterFollowTracker()
	// chatbot.RegisterSubscriberTracker()

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
