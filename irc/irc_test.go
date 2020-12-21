package irc

import (
	"bytes"
	"strings"
	"testing"
)

type TestConnection struct {
	bytes.Buffer
}

// Helper method for checking if the given string was sent to the TestConnection
func (t *TestConnection) Received(str string) bool {
	return strings.Contains(t.String(), str)
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

func TestRead(t *testing.T) {}

func TestWrite(t *testing.T) {}
