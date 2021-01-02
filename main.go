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
	"text/template"
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
	var enableSubsMessage bool

	flag.StringVar(&channel, "channel", "", "Channel name, without the #, to join")
	flag.StringVar(&nick, "nick", "", "Nickname to join Chat with")
	flag.StringVar(&ledgerType, "ledger", ledger.REDIS, fmt.Sprintf("Ledger type string. Options: %s, %s, %s", ledger.REDIS, ledger.FILE, ledger.MEM))
	flag.StringVar(&secretStoreType, "store", secret.VAULT, fmt.Sprintf("Secret Store type string. Options: %s, %s", secret.VAULT, secret.ENV))

	flag.BoolVar(&enableAll, "all", false, "Enable all features")
	flag.BoolVar(&enableGreeter, "greeter", false, "Enable the auto-greeter")
	flag.BoolVar(&enableCommands, "commands", false, "Enable Command processing")
	flag.BoolVar(&enableRaidMessage, "raids", false, "Enable Raid auto-shoutout")
	flag.BoolVar(&enableBitsMessage, "bits", false, "Enable Bits auto-thanks")
	flag.BoolVar(&enableSubsMessage, "subs", false, "Enable Subs auto-thanks")

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

	// Initialize Secrets Store
	store, err := secret.NewSecretStore(secretStoreType)
	if err != nil {
		log.Fatalf("FATAL: Create secret store - %v", err)
	}

	password, err := store.TwitchToken()
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

	// IRC is both a Client and a ChatClient
	chatBot.RegisterClient(irc)
	chatBot.SetChatClient(irc)

	// Feature Toggles
	if enableCommands || enableAll {
		chatBot.HandleCommands()
	}

	if enableGreeter || enableAll {
		// Ledger for the auto greeter
		// TODO get from config
		var expirationTime int64 = 1000 * 60 * 60 * 12 // 12 hours
		ledger, err := ledger.NewLedger(ledgerType, expirationTime)
		if err != nil {
			log.Fatalf("FATAL: create ledger - %v", err)
		}

		// pre-seed names we don't want greeted
		ledger.Add("streamlabs")
		ledger.Add("nightbot")
		ledger.Add("ranaebot")
		ledger.Add("soundalerts")
		ledger.Add(strings.TrimPrefix(channel, "#")) // Prevent greeting the broadcaster

		// Greeter config
		greetKey := fmt.Sprintf("%s.greeter.messageFormat", strings.Trim(channel, "#"))
		greetMessageFormat := viper.GetString(greetKey)
		greetBot := greeter.New(ledger)
		greetTempl, err := template.New("greeter").Parse(greetMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid Greeter message in config - %v", err)
		}

		chatBot.RegisterGreeter(greetBot, bot.NewHandlerTemplate(greetTempl))
	}

	if enableRaidMessage || enableAll {
		raidKey := fmt.Sprintf("%s.raid.messageFormat", strings.Trim(channel, "#"))
		raidMessageFormat := viper.GetString(raidKey)

		viper.SetDefault("%s.raid.delaySeconds", 0)
		raidDelayKey := fmt.Sprintf("%s.raid.delaySeconds", strings.Trim(channel, "#"))
		raidDelay := viper.GetInt(raidDelayKey)

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid raid message in config - %v", err)
		}
		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if enableBitsMessage || enableAll {
		bitsKey := fmt.Sprintf("%s.bits.messageFormat", strings.Trim(channel, "#"))
		bitsMessageFormat := viper.GetString(bitsKey)

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid bits message in config - %v", err)
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if enableSubsMessage || enableAll {
		subsKey := fmt.Sprintf("%s.subs.messageFormat", strings.Trim(channel, "#"))
		subsMessageFormat := viper.GetString(subsKey)
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		giftSubsKey := fmt.Sprintf("%s.giftsubs.messageFormat", strings.Trim(channel, "#"))
		giftSubsMessageFormat := viper.GetString(giftSubsKey)
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		chatBot.RegisterSubsHandler(
			bot.NewHandlerTemplate(subsTempl), bot.NewHandlerTemplate(giftSubsTempl))
	}

	if err := chatBot.Start(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}

	// Keep the process alive
	for {
		time.Sleep(time.Second)
	}
}
