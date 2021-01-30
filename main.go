package main

import (
	"flag"
	"fmt"
	"medgebot/bot"
	"medgebot/config"
	"medgebot/irc"
	"medgebot/ledger"
	log "medgebot/logger"
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
		log.Panic("init config", err)
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
		log.Panic("Create secret store", err)
	}

	password, err := store.TwitchToken()
	if err != nil {
		log.Panic("Get Twitch Token from store", err)
	}

	// IRC
	nick := conf.Nick()
	if nick == "" {
		log.Panic("config key - nick not found / empty", nil)
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
		log.Panic("start IRC", err)
	}

	// IRC is both a Client and a ChatClient
	chatBot.RegisterClient(irc)
	chatBot.SetChatClient(irc)

	// Feature Toggles
	if conf.CommandsEnabled() || enableAll {
		cmds := conf.KnownCommands()
		var commands []bot.Command

		for _, cmd := range cmds {
			cmdTemplate, err := template.New(cmd.Prefix).Parse(cmd.Message)
			if err != nil {
				log.Panic(fmt.Sprintf("parse known Command [%+v]", cmd), err)
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
		// Ledger for the auto greeter
		ledger, err := ledger.NewLedger(conf)
		if err != nil {
			log.Panic("create ledger", err)
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
			log.Panic("invalid Greeter message in config", err)
		}

		chatBot.RegisterGreeter(ledger, bot.NewHandlerTemplate(greetTempl))
	}

	if conf.RaidsEnabled() || enableAll {
		raidMessageFormat := conf.RaidsMessageFormat()
		raidDelay := conf.RaidDelay()

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Panic("invalid raid message in config", err)
		}
		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if conf.BitsEnabled() || enableAll {
		bitsMessageFormat := conf.BitsMessageFormat()

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Panic("invalid bits message in config", err)
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if conf.SubsEnabled() || enableAll {
		subsMessageFormat := conf.SubsMessageFormat()
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Panic("invalid subs message in config", err)
		}

		giftSubsMessageFormat := conf.GiftSubsMessageFormat()
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Panic("invalid subs message in config", err)
		}

		chatBot.RegisterSubsHandler(
			bot.NewHandlerTemplate(subsTempl), bot.NewHandlerTemplate(giftSubsTempl))
	}

	// Start the Bot only after all handlers are loaded
	if err := chatBot.Start(); err != nil {
		log.Panic("bot connect", err)
	}

	// Start HTTP server
	srv := server.New()
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Panic("start HTTP server", err)
	}
}
