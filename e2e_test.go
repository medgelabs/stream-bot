package main

import (
	"fmt"
	"medgebot/bot"
	"medgebot/bot/bottest"
	"medgebot/cache"
	"medgebot/irc"
	"medgebot/irc/irctest"
	"medgebot/ws/wstest"
	"testing"
)

const (
	USER    = "medgelabs"
	CHANNEL = "medgelabs"
)

func TestRaids(t *testing.T) {
	ws := wstest.NewWebsocket()

	ircClient := irc.NewClient(ws)
	ircConf := irc.Config{
		Nick:     USER,
		Password: "oauth:superSpookyGhostMachineTestSecret",
		Channel:  "#" + CHANNEL,
	}
	if err := ircClient.Start(ircConf); err != nil {
		t.Fatalf("Failed to start IRC client: %v", err)
	}

	cache, _ := cache.InMemory(0)
	chatBot := bot.New(&cache)
	chatBot.RegisterClient(ircClient)
	chatBot.SetChatClient(ircClient)

	raidTmpl := bot.NewHandlerTemplate(bottest.MakeTemplate("raids", "{{.Sender}} raid of {{.Amount}}"))
	chatBot.RegisterRaidHandler(raidTmpl, 1)

	// We must Start the bot AFTER the handler is registered
	chatBot.Start()

	// We should see "USER raid of RAID_SIZE" eventually come through the IRC client
	raidSize := 5
	expectedMessage := fmt.Sprintf("PRIVMSG #medgelabs :%s raid of %d", USER, raidSize)
	ws.SendAndWait(irctest.MakeRaidMessage(USER, raidSize, CHANNEL))

	if !ws.Received(expectedMessage) {
		t.Fatalf("Did not see expected Raid Message.\nWS Dump:\n%s", ws.String())
	}
}
