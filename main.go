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
	"os"
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
		log.Fatal(err, "init config")
	}

	// Channel should be prefixed with # by default. Add it if missing
	if !strings.HasPrefix(channel, "#") {
		channel = fmt.Sprintf("#%s", channel)
	}

	// Cache for various stream metrics
	metricsCache := mustCreateFileCache("metrics.txt", 0)

	// Initialize desired state for the bot
	chatBot := bot.New(metricsCache)
	chatBot.RegisterReadLogger()

	// Initialize Secrets Store
	store, err := secret.NewSecretStore(conf)
	if err != nil {
		log.Fatal(err, "Create secret store")
	}

	password, err := store.TwitchToken()
	if err != nil {
		log.Fatal(err, "Get Twitch Token from store")
	}

	// IRC
	nick := conf.Nick()
	if nick == "" {
		log.Fatal(nil, "config key - nick not found / empty")
	}

	ircConfig := irc.Config{
		Nick:     nick,
		Password: fmt.Sprintf("oauth:%s", password),
		Channel:  channel,
	}

	ircWs := ws.NewWebSocket("wss", "irc-ws.chat.twitch.tv:443")
	err = ircWs.Connect()
	if err != nil {
		log.Fatal(err, "irc ws connect")
	}

	ircClient := irc.NewClient(ircWs)
	ircWs.SetPostReconnectFunc(func() error {
		return ircClient.Start(ircConfig)
	})
	defer ircClient.Close()

	err = ircClient.Start(ircConfig)
	if err != nil {
		log.Fatal(err, "start IRC")
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
			log.Fatal(err, "pubsub ws connect")
		}

		pubsub := pubsub.NewClient(pubSubWs, conf.ChannelID(), password)
		pubSubWs.SetPostReconnectFunc(pubsub.Start)
		pubsub.Start()
		chatBot.RegisterClient(pubsub)
	}

	// Shoutout Command
	chatBot.HandleShoutoutCommand()

	// Feature Toggles
	if conf.CommandsEnabled() || enableAll {
		cmds := conf.KnownCommands()
		var commands []bot.Command

		for _, cmd := range cmds {
			cmdTemplate, err := template.New(cmd.Prefix).Parse(cmd.Message)
			if err != nil {
				log.Fatal(err, "parse known Command [%+v]", cmd)
			}

			cmd := bot.Command{
				Prefix:          cmd.Prefix,
				IsAlias:         cmd.AliasFor != "",
				AliasFor:        cmd.AliasFor,
				MessageTemplate: bot.NewHandlerTemplate(cmdTemplate),
			}

			commands = append(commands, cmd)
		}

		chatBot.HandleCommands(commands)
	}

	// if config.greeterEnabled() {
	if conf.GreeterEnabled() || enableAll {
		// Cache for the auto greeter
		greeterCache := mustCreateFileCache("greeter.txt", conf.CacheExpirationTime())

		// pre-seed names we want ignored
		greeterCache.Put("streamlabs", "")
		greeterCache.Put("nightbot", "")
		greeterCache.Put("ranaebot", "")
		greeterCache.Put("soundalerts", "")
		greeterCache.Put("jtv", "")
		greeterCache.Put(strings.TrimPrefix(channel, "#"), "") // Prevent greeting the broadcaster

		// Greeter config
		greetMessageFormat := conf.GreetMessageFormat()
		greetTempl, err := template.New("greeter").Parse(greetMessageFormat)
		if err != nil {
			log.Fatal(err, "invalid Greeter message in config")
		}

		chatBot.RegisterGreeter(greeterCache, bot.NewHandlerTemplate(greetTempl))
	}

	if conf.RaidsEnabled() || enableAll {
		raidMessageFormat := conf.RaidsMessageFormat()
		raidDelay := conf.RaidDelay()

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Fatal(err, "invalid raid message in config")
		}
		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if conf.BitsEnabled() || enableAll {
		bitsMessageFormat := conf.BitsMessageFormat()

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Fatal(err, "invalid bits message in config")
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if conf.SubsEnabled() || enableAll {
		subsMessageFormat := conf.SubsMessageFormat()
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Fatal(err, "invalid subs message in config")
		}

		giftSubsMessageFormat := conf.GiftSubsMessageFormat()
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Fatal(err, "invalid subs message in config")
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
		log.Fatal(err, "bot connect")
	}

	// TODO should we have a setter for WS if AlertsEnabled?
	// Start HTTP server
	debugClient := server.DebugClient{}
	chatBot.RegisterClient(&debugClient)
	srv := server.New(metricsCache, &ws, &debugClient)
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err, "start HTTP server")
	}
}

// Create a PersistableCache backed by a file. Panics if it cannot
func mustCreateFileCache(filepath string, keyExpirationSeconds int64) *cache.PersistableCache {
	cacheFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err, "create cache file")
	}

	cache, err := cache.FilePersisted(cacheFile, keyExpirationSeconds)
	if err != nil {
		log.Fatal(err, "read greeter cache file")
	}

	return &cache
}
