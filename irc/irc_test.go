package irc

import (
	"bytes"
	"medgebot/bot"
	"strings"
	"sync"
	"testing"
	"time"
)

type TestConnection struct {
	bytes.Buffer
	sync.Mutex
}

// Helper method for checking if the given string was sent to the TestConnection
func (t *TestConnection) Received(str string) bool {
	t.Lock()
	defer t.Unlock()
	return strings.Contains(t.String(), str)
}

func (t *TestConnection) Clear() {
	t.Lock()
	defer t.Unlock()
	t.Reset()
}

func (t *TestConnection) Close() error {
	return nil
}

func TestStart(t *testing.T) {
	// Ensure PASS, NICK, JOIN, and CAP REQ commands are sent
	conn := &TestConnection{}

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
	conn := &TestConnection{}
	config := Config{
		Nick:     "medgelabs",
		Password: "oauth:secret",
		Channel:  "#medgelabs",
	}
	irc := NewClient(conn)

	testBot := make(chan bot.Event)
	irc.SetDestination(testBot)

	irc.Start(config)
	conn.Clear()

	conn.WriteString(CHAT_MSG_TAGS)

	// Wait for message on
	select {
	case evt := <-testBot:
		if evt.Type != bot.CHAT_MSG {
			t.Fatalf("Did not receive a Chat message. Got %+v", evt)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("Failed to receive expected message")
	}
}

// func TestSendMessageToServer(t *testing.T) {
// t.Fatalf("Not Implemented")
// }
