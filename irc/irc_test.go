package irc

import (
	"medgebot/bot"
	"medgebot/irc/irctest"
	"medgebot/ws/wstest"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	// Ensure PASS, NICK, JOIN, and CAP REQ commands are sent
	conn := wstest.NewWebsocket()

	config := Config{
		Nick:     "medgelabs",
		Password: "oauth:secret",
		Channel:  "#medgelabs",
	}

	irc := NewClient(conn)
	irc.Start(config)

	output := conn.String()
	if !conn.Received("PASS " + config.Password) {
		t.Fatalf("PASS command not sent to connection. Sent: %s", output)
	}

	if !conn.Received("NICK " + config.Nick) {
		t.Fatalf("NICK command not sent to connection. Sent: %s", output)
	}

	if !conn.Received("JOIN " + config.Channel) {
		t.Fatalf("JOIN command not sent to connection. Sent: %s", output)
	}
}

func TestMessageReceivedFromServer(t *testing.T) {
	conn := wstest.NewWebsocket()
	config := Config{
		Nick:     "medgelabs",
		Password: "oauth:secret",
		Channel:  "#medgelabs",
	}
	irc := NewClient(conn)

	testBot := make(chan bot.Event)
	irc.SetDestination(testBot)
	irc.Start(config)

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
