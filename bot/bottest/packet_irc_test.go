package bottest

import "testing"

func TestMakeIrcMessage(t *testing.T) {
	tags := make(map[string]string)
	tags["display-name"] = "medgelabs"
	tags["emotes"] = ""
	tags["subscriber"] = "1"

	expected := "@display-name=medgelabs;emotes=;subscriber=1; :medgelabs!medgelabs@medgelabs.tmi.twitch.tv PRIVMSG #medgelabs :!hello"
	result := MakeIrcMessage("!hello", "medgelabs", "PRIVMSG", "medgelabs", tags)

	if result != expected {
		t.Fatalf("Got wrong IRC message. Expected\n%s\nGot:\n%s", expected, result)
	}
}
