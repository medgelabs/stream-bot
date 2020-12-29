package main

import (
	"flag"
	"fmt"
	"log"
	"medgebot/bot"
	"medgebot/config"
	"medgebot/greeter"
	"medgebot/irc"
	"medgebot/ledger"
	"medgebot/secret"
	"medgebot/ws"
	"strings"
	"text/template"
	"time"
)

func main() {

	// CLI argument processing
	var channel string
	var nick string
	var ledgerType string
	var secretStoreType string

	flag.StringVar(&channel, "channel", "", "Channel name, without the #, to join")
	flag.StringVar(&nick, "nick", "", "Nickname to join Chat with")
	flag.StringVar(&ledgerType, "ledger", ledger.REDIS, fmt.Sprintf("Ledger type string. Options: %s, %s, %s", ledger.REDIS, ledger.FILE, ledger.MEM))
	flag.StringVar(&secretStoreType, "store", secret.VAULT, fmt.Sprintf("Secret Store type string. Options: %s, %s", secret.VAULT, secret.ENV))

	flag.Parse()

	// Flag error handling
	if !strings.HasPrefix(channel, "#") {
		channel = fmt.Sprintf("#%s", channel)
	}

	if nick == "" {
		log.Fatalln("FATAL: nick empty")
	}

	conf := config.LoadConfig()

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
	ircWs := ws.NewWebsocket()
	ircWs.Connect("wss", "irc-ws.chat.twitch.tv:443")
	irc := irc.NewClient(ircWs, channel)
	defer irc.Close()

	pass := fmt.Sprintf("oauth:%s", password)
	err = irc.Start(nick, pass)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// IRC is both a Client and a ChatClient
	chatBot.RegisterClient(irc)
	chatBot.SetChatClient(irc)

	// Feature enablement

	if conf.FeatureEnabled("commands") {
		chatBot.HandleCommands()
	}

	if conf.FeatureEnabled("greeter") {
		var expirationTime int64 = conf.GetInt64("greeter.expirationTimeSeconds")
		ledger, err := ledger.NewLedger(ledgerType, expirationTime)
		if err != nil {
			log.Fatalf("FATAL: create ledger - %v", err)
		}

		// pre-seed names we don't want greeted
		ignores := conf.GetList("greeter.ignore")
		for _, ignore := range ignores {
			ledger.Add(ignore)
		}

		// Greeter config
		greetMessageFormat := conf.GetString("greeter.messageFormat")
		greetBot := greeter.New(ledger)
		greetTempl, err := template.New("greeter").Parse(greetMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid Greeter message in config - %v", err)
		}

		chatBot.RegisterGreeter(greetBot, bot.NewHandlerTemplate(greetTempl))
	}

	if conf.FeatureEnabled("raids") {
		raidMessageFormat := conf.GetString("raids.messageFormat")
		raidDelay := conf.GetIntOrDefault("raids.delaySeconds", 0)

		raidTempl, err := template.New("raids").Parse(raidMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid raid message in config - %v", err)
		}

		chatBot.RegisterRaidHandler(
			bot.NewHandlerTemplate(raidTempl), raidDelay)
	}

	if conf.FeatureEnabled("bits") {
		bitsMessageFormat := conf.GetString("bits.messageFormat")

		bitsTempl, err := template.New("bits").Parse(bitsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid bits message in config - %v", err)
		}

		chatBot.RegisterBitsHandler(
			bot.NewHandlerTemplate(bitsTempl))
	}

	if conf.FeatureEnabled("subs") {
		subsMessageFormat := conf.GetString("subs.messageFormat")
		subsTempl, err := template.New("subs").Parse(subsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		chatBot.RegisterSubsHandler(
			bot.NewHandlerTemplate(subsTempl))
	}

	if conf.FeatureEnabled("giftsubs") {
		giftSubsMessageFormat := conf.GetString("giftsubs.messageFormat")
		giftSubsTempl, err := template.New("giftsubs").Parse(giftSubsMessageFormat)
		if err != nil {
			log.Fatalf("FATAL: invalid subs message in config - %v", err)
		}

		chatBot.RegisterGiftSubsHandler(
			bot.NewHandlerTemplate(giftSubsTempl))
	}

	if err := chatBot.Start(); err != nil {
		log.Fatalf("FATAL: bot connect - %v", err)
	}

	// Keep the process alive
	for {
		time.Sleep(time.Second)
	}
}
