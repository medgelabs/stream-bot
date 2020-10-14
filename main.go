package main

import (
	"log"
	"medgebot/bot"
	"medgebot/secret"
	"os"
	"time"
)

func main() {
	channel := "#medgelabs"
	nick := "medgelabs"

	// Ledger for the auto greeter
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	ledger, err := bot.NewRedisLedger(redisHost, redisPort)
	if err != nil {
		log.Fatalf("FATAL - connect to Redis - %s", err)
	}

	// pre-seed names we don't want greeted
	ledger.Add("tmi.twitch.tv")
	ledger.Add("streamlabs")
	ledger.Add(nick)
	ledger.Add(nick + "@tmi.twitch.tv")

	// Initialize Secrets Store
	vaultUrl := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	store := secret.NewVaultStore("secret/twitchToken")
	if err := store.Connect(vaultUrl, vaultToken); err != nil {
		log.Fatalf("FATAL: Vault connect - %v", err)
	}

	password, err := store.GetTwitchToken()
	if err != nil {
		log.Fatalf("FATAL: Get Twitch Token from store - %v", err)
	}

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterPingPong()
	chatBot.RegisterReadLogger()
	chatBot.HandleCommands()
	chatBot.RegisterGreeter(&ledger)

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
