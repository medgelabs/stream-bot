package main

import (
	"medgebot/bot"
	"medgebot/bot/bottest"
	"medgebot/irc"
	"testing"
)

const (
	USER    = "medgelabs"
	CHANNEL = "medgelabs"
)

func TestRaids(t *testing.T) {
	ws := bottest.NewTestWebsocket()
	ircClient := irc.NewClient(ws)
	ircConf := irc.Config{
		Nick:     USER,
		Password: "oauth:superSpookyGhostMachineTestSecret",
		Channel:  CHANNEL,
	}

	if err := ircClient.Start(ircConf); err != nil {
		t.Fatalf("Failed to start IRC client: %v", err)
	}

	chatBot := bot.New()
	chatBot.RegisterClient(ircClient)
	chatBot.SetChatClient(ircClient)

	raidTmpl := bot.NewHandlerTemplate(bottest.MakeTemplate("raids", "Raid {{.Sender}} of {{.Amount}}"))
	chatBot.RegisterRaidHandler(raidTmpl, 0)
}
