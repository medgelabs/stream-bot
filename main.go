package main

import (
	"flag"
	"fmt"
	"medgebot/bot"
	"medgebot/cache"
	"medgebot/config"
	"medgebot/irc"
	log "medgebot/logger"
	"medgebot/pubsub"
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
		log.Fatal("init config", err)
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
		log.Fatal("Create secret store", err)
	}

	password, err := store.TwitchToken()
	if err != nil {
		log.Fatal("Get Twitch Token from store", err)
	}

	// IRC
	nick := conf.Nick()
	if nick == "" {
		log.Fatal("config key - nick not found / empty", nil)
	}

	ircConfig := irc.Config{
		Nick:     nick,
		Password: fmt.Sprintf("oauth:%s", password),
		Channel:  channel,
	}

	ircWs := ws.NewWebSocket("wss", "irc-ws.chat.twitch.tv:443")
	err = ircWs.Connect()
	if err != nil {
		log.Fatal("irc ws connect", err)
	}

	ircClient := irc.NewClient(ircWs)
	ircWs.SetPostReconnectFunc(func() error {
		return ircClient.Start(ircConfig)
	})
	defer ircClient.Close()

	err = ircClient.Start(ircConfig)
	if err != nil {
		log.Fatal("start IRC", err)
	}

	// IRC is both a Client and a ChatClient
	chatBot.RegisterClient(ircClient)
	chatBot.SetChatClient(ircClient)

	// TODO pubsub is only used for ChannelPoints at this time.
	// If we use pubsub for other features, it wouldn't make sense to
	// guard pubsub creation behind this feature flag
	if conf.ChannelPointsEnabled() || enableAll {
		pubSubWs := ws.NewWebSocket("wss", "pubsub-edge.twitch.tv")
		err = pubSubWs.Connect()
		if err != nil {
			log.Fatal("pubsub ws connect", err)
		}

		pubsub := pubsub.NewClient(pubSubWs, conf.ChannelID(), password)
		pubSubWs.SetPostReconnectFunc(pubsub.Start)
		pubsub.Start()
		chatBot.RegisterClient(pubsub)
	}

	// Feature Toggles
	if conf.CommandsEnabled() || enableAll {
		cmds := conf.KnownCommands()
		var commands []bot.Command

		for _, cmd := range cmds {
			cmdTemplate, err := template.New(cmd.Prefix).Parse(cmd.Message)
			if err != nil {
				log.Fatal(fmt.Sprintf("parse known Command [%+v]", cmd), err)
			}

			cmd := bot.Command{
				Prefix:          cmd.Prefix,
				MessageTemplate: bot.NewHandlerTemplate(cmdTemplate),
			}

			commands = append(commands, cmd)
		}

		chatBot.HandleCommands(commands)
	}

	// if config.greeterEnabled() {
	if conf.GreeterEnabled() || enableAll {
		// Cache for the auto greeter
		cache, err := cache.New(conf)
		if err != nil {
			log.Fatal("create cache", err)
		}

		// pre-seed names we want ignored
		cache.Put("streamlabs", "")
		cache.Put("nightbot", "")
		cache.Put("ranaebot", "")
		cache.Put("soundalerts", "")
		cache.Put("jtv", "")
		cache.Put(strings.TrimPrefix(channel, "#"), "") // Prevent greeting the broadcaster

		// Greeter config
		greetMessageFormat := conf.GreetMessageFormat()
		greetTempl, err := template.New("greeter").Parse(greetMessageFormat)
		if err != nil {
			log.Fatal("invalid Greeter message in config", err)
		}

		chatBot.RegisterGreeter(cache, bot.NewHandlerTemplate(greetTempl))
	}

	if conf.RaidsEnabled() || enableAll {
		raidMessageFormat := conf.RaidsMessageFormat()
		raidDelay := conf.RaidDelay()

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Fatal("invalid raid message in config", err)
		}
		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if conf.BitsEnabled() || enableAll {
		bitsMessageFormat := conf.BitsMessageFormat()

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Fatal("invalid bits message in config", err)
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if conf.SubsEnabled() || enableAll {
		subsMessageFormat := conf.SubsMessageFormat()
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Fatal("invalid subs message in config", err)
		}

		giftSubsMessageFormat := conf.GiftSubsMessageFormat()
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Fatal("invalid subs message in config", err)
		}

		chatBot.RegisterSubsHandler(
			bot.NewHandlerTemplate(subsTempl), bot.NewHandlerTemplate(giftSubsTempl))
	}

	// Alerts link between the Bot and the Web API
	ws := bot.WriteOnlyUnsafeWebSocket{}
	if conf.AlertsEnabled() || enableAll {
		chatBot.RegisterAlertHandler(&ws)
	}

	// Start the Bot only after all handlers are loaded
	if err := chatBot.Start(); err != nil {
		log.Fatal("bot connect", err)
	}

	// TODO should we have a setter for WS if AlertsEnabled?
	// Start HTTP server
	metricsCache, _ := cache.InMemory(0)
	srv := server.New(&metricsCache, &ws)
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal("start HTTP server", err)
	}
}
