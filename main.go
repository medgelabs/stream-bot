package main

import (
	"flag"
	"fmt"
	"log"
	"medgebot/bot"
	"medgebot/config"
	"medgebot/irc"
	"medgebot/ledger"
	"medgebot/secret"
	"medgebot/server"
	"medgebot/ws"
	"net/http"
	"strings"
	"text/template"
)

func main() {

	// CLI argument processing
	var channel string
	var configPath string
	var enableAll bool

	flag.StringVar(&channel, "channel", "", "Channel name, without the #, to join")
	flag.StringVar(&configPath, "config", ".", "Path to the config.yaml file. Default: .")
	flag.BoolVar(&enableAll, "all", false, "Enable all features")
	flag.Parse()

	conf, err := config.New(channel, configPath)
	if err != nil {
		log.Fatalf("FATAL: init config - %v", err)
	}

	// Channel should be prefixed with # by default. Add it if missing
	if !strings.HasPrefix(channel, "#") {
		channel = fmt.Sprintf("#%s", channel)
	}

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterReadLogger()

	// Initialize Secrets Store
	store, err := secret.NewSecretStore(conf)
	if err != nil {
		log.Fatalf("FATAL: Create secret store - %v", err)
	}

	password, err := store.TwitchToken()
	if err != nil {
		log.Fatalf("FATAL: Get Twitch Token from store - %v", err)
	}

	// IRC
	nick := conf.Nick()
	if nick == "" {
		log.Fatalf("FATAL: config key - nick not found / empty")
	}

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
	if conf.CommandsEnabled() || enableAll {
		chatBot.HandleCommands()
	}

	// if config.greeterEnabled() {
	if conf.GreeterEnabled() || enableAll {
		// Ledger for the auto greeter
		ledger, err := ledger.NewLedger(conf)
		if err != nil {
			log.Fatalf("FATAL: create ledger - %v", err)
		}

		// pre-seed names we don't want greeted
		ledger.Put("streamlabs", "")
		ledger.Put("nightbot", "")
		ledger.Put("ranaebot", "")
		ledger.Put("soundalerts", "")
		ledger.Put(strings.TrimPrefix(channel, "#"), "") // Prevent greeting the broadcaster

		// Greeter config
		greetMessageFormat := conf.GreetMessageFormat()
		greetTempl, err := template.New("greeter").Parse(greetMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid Greeter message in config - %v", err)
		}

		chatBot.RegisterGreeter(ledger, bot.NewHandlerTemplate(greetTempl))
	}

	if conf.RaidsEnabled() || enableAll {
		raidMessageFormat := conf.RaidsMessageFormat()
		raidDelay := conf.RaidDelay()

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid raid message in config - %v", err)
		}
		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if conf.BitsEnabled() || enableAll {
		bitsMessageFormat := conf.BitsMessageFormat()

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid bits message in config - %v", err)
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if conf.SubsEnabled() || enableAll {
		subsMessageFormat := conf.SubsMessageFormat()
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		giftSubsMessageFormat := conf.GiftSubsMessageFormat()
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		chatBot.RegisterSubsHandler(
			bot.NewHandlerTemplate(subsTempl), bot.NewHandlerTemplate(giftSubsTempl))
	}

	// Start the Bot only after all handlers are loaded
	if err := chatBot.Start(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}

	// Start HTTP server
	srv := server.New()
	log.Fatal(http.ListenAndServe(":8080", srv))
}
