package irc

import (
	"medgebot/bot"
	"medgebot/irc/irctest"
	"medgebot/ws/wstest"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	pass := "oauth:secret"
	nick := "medgelabs"
	channel := "medgelabs"

	// Ensure PASS, NICK, JOIN, and CAP REQ commands are sent
	conn := wstest.NewWebsocket()
	irc := NewClient(conn, channel)

	irc.Start(nick, pass)
	output := conn.String()
	if !conn.Received("PASS " + pass) {
		t.Fatalf("PASS command not sent to connection. Sent: %s", output)
	}

	if !conn.Received("NICK " + nick) {
		t.Fatalf("NICK command not sent to connection. Sent: %s", output)
	}

	if !conn.Received("JOIN " + channel) {
		t.Fatalf("JOIN command not sent to connection. Sent: %s", output)
	}
}

func TestMessageReceivedFromServer(t *testing.T) {
	pass := "oauth:secret"
	nick := "medgelabs"
	channel := "medgelabs"

	conn := wstest.NewWebsocket()
	irc := NewClient(conn, channel)

	testBot := make(chan bot.Event)
	irc.SetDestination(testBot)

	irc.Start(nick, pass)

	conn.Send(irctest.MakeChatMessage("testuser", "Chat!", "medgelabs"))

	// Wait for message on bot Event channel
	select {
	case evt := <-testBot:
		if evt.Type != bot.CHAT_MSG {
			t.Fatalf("Did not receive a Chat message. Got %+v", evt)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Failed to receive expected message")
	}
}
