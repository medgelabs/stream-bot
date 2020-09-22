package bot

import (
	"strings"
	"testing"
)

func TestPrivMsg(t *testing.T) {
	bot := Bot{}
	result := bot.PrivMsg("medgelabs", "Foobar is a lie")

	if !strings.HasPrefix(result, "PRIVMSG") {
		t.Fatalf("Expected PRIVMSG command. Got: %s", result)
	}

	if result != "PRIVMSG #medgelabs :Foobar is a lie" {
		t.Fatalf("Received invalid PRIVMSG line: %s", result)
	}
}
