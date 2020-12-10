package main

import (
	"flag"
	"fmt"
	"log"
	"medgebot/bot"
	"medgebot/greeter"
	"medgebot/irc"
	"medgebot/ledger"
	"medgebot/secret"
	"medgebot/ws"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func main() {

	// CLI argument processing
	var channel string
	var nick string
	var ledgerType string
	var secretStoreType string
	var enableAll bool
	var enableGreeter bool
	var enableCommands bool
	var enableRaidMessage bool
	var enableBitsMessage bool

	flag.StringVar(&channel, "channel", "", "Channel name, without the #, to join")
	flag.StringVar(&nick, "nick", "", "Nickname to join Chat with")
	flag.StringVar(&ledgerType, "ledger", ledger.REDIS, fmt.Sprintf("Ledger type string. Options: %s, %s, %s", ledger.REDIS, ledger.FILE, ledger.MEM))
	flag.StringVar(&secretStoreType, "store", secret.VAULT, fmt.Sprintf("Secret Store type string. Options: %s, %s", secret.VAULT, secret.ENV))

	flag.BoolVar(&enableAll, "all", false, "Enable all features")
	flag.BoolVar(&enableGreeter, "greeter", false, "Enable the auto-greeter")
	flag.BoolVar(&enableCommands, "commands", false, "Enable Command processing")
	flag.BoolVar(&enableRaidMessage, "raids", false, "Enable Raid auto-shoutout")
	flag.BoolVar(&enableBitsMessage, "bits", false, "Enable Bits auto-thanks")

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

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterReadLogger()

	/// Plugin Registration ///

	// Initialize Secrets Store
	store, err := secret.NewSecretStore(secretStoreType)
	if err != nil {
		log.Fatalf("FATAL: Create secret store - %v", err)
	}

	password, err := store.GetTwitchToken()
	if err != nil {
		log.Fatalf("FATAL: Get Twitch Token from store - %v", err)
	}

	// IRC
	ircConfig := irc.Config{
		Nick:     nick,
		Password: fmt.Sprintf("oauth:%s", password),
		Channel:  channel,
	}

	ircWs := ws.NewWebsocket()
	ircWs.Connect("wss", "irc-ws.chat.twitch.tv:443")
	irc := irc.NewClient(ircWs)
	defer irc.Close()

	err = irc.Start(ircConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := chatBot.RegisterPlugin(irc); err != nil {
		log.Fatalf("FATAL: failed to register plugin: %v", err)
	}

	// PubSub
	// channelIdKey := fmt.Sprintf("%s.channelId", strings.Trim(channel, "#"))
	// channelId := viper.GetString(channelIdKey)
	// pubSub := twitch.NewPubSubClient(channelId, password)
	// defer pubSub.Close()

	// err = pubSub.Connect("wss", "pubsub-edge.twitch.tv")
	// if err != nil {
	// log.Fatalf(err.Error())
	// }

	// go pubSub.Start()

	// if err := chatBot.RegisterInboundPlugin(pubSub); err != nil {
	// log.Fatalf("FATAL: failed to register PubSub: %v", err)
	// }
	/// Plugin Registration END ///

	if enableCommands || enableAll {
		chatBot.HandleCommands()
	}

	if enableGreeter || enableAll {
		// Ledger for the auto greeter
		ledger, err := ledger.NewLedger(ledgerType)
		if err != nil {
			log.Fatalf("FATAL: create ledger - %v", err)
		}

		// pre-seed names we don't want greeted
		ledger.Add("tmi.twitch.tv")
		ledger.Add("streamlabs")
		ledger.Add("nightbot")
		ledger.Add("ranaebot")
		ledger.Add("soundalerts")
		ledger.Add(nick)
		ledger.Add(strings.TrimPrefix(channel, "#")) // Prevent greeting the broadcaster
		ledger.Add(nick + ".tmi.twitch.tv")
		ledger.Add(nick + "@tmi.twitch.tv")

		// Greeter config
		var greetConfig greeter.Config
		confKey := fmt.Sprintf("%s.greeter", strings.Trim(channel, "#"))
		confSub := viper.Sub(confKey)
		if confSub == nil {
			log.Fatalf("FATAL: key %s not found in config", confKey)
		}

		confSub.Unmarshal(&greetConfig)
		greetBot := greeter.New(greetConfig, ledger)
		chatBot.RegisterGreeter(greetBot)
	}

	if enableRaidMessage || enableAll {
		raidKey := fmt.Sprintf("%s.raid.messageFormat", strings.Trim(channel, "#"))
		raidMessageFormat := viper.GetString(raidKey)
		chatBot.RegisterRaidHandler(raidMessageFormat)
	}

	if enableBitsMessage || enableAll {
		bitsKey := fmt.Sprintf("%s.bits.messageFormat", strings.Trim(channel, "#"))
		bitsMessageFormat := viper.GetString(bitsKey)
		chatBot.RegisterBitsHandler(bitsMessageFormat)
	}

	if err := chatBot.Start(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}

	// Keep the process alive
	for {
		time.Sleep(time.Second)
	}
}
