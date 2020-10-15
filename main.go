package main

import (
	"fmt"
	"log"
	"medgebot/bot"
	persistance "medgebot/internal/pkg/ledger"
	"medgebot/internal/pkg/storage"
	"os"
	"time"
)

const DefaultSecretPath = "secret/data/twitchToken"

func main() {
	channel := "#medgelabs"
	nick := "medgelabs"

	// Ledger for the auto greeter
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	redisEngine, err := persistance.NewRedis(redisHost, redisPort)
	if err != nil {
		log.Fatalf("FATAL - connect to Redis - %s", err)
	}
	ledger := persistance.New(redisEngine)

	// pre-seed names we don't want greeted
	err = ledger.Add("tmi.twitch.tv")
	err = ledger.Add("streamlabs")
	err = ledger.Add(nick)
	err = ledger.Add(nick + "@tmi.twitch.tv")
	if err != nil {
		fmt.Printf("error adding username to ledger - %v", err)
	}

	// Initialize Secrets Store
	vaultUrl := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	vaultEngine, err := storage.NewVault(vaultUrl, vaultToken, DefaultSecretPath)
	if err != nil {
		log.Fatalf("FATAL: Vault connect - %v", err)
	}
	store := storage.New(vaultEngine)

	password, err := store.GetString("token")
	if err != nil {
		log.Fatalf("FATAL: Get Twitch Token from store - %v", err)
	}

	// Initialize desired state for the bot
	chatBot := bot.New()
	chatBot.RegisterPingPong()
	chatBot.RegisterReadLogger()
	chatBot.HandleCommands()
	chatBot.RegisterGreeter(ledger)

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
