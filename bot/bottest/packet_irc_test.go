package bottest

import (
	"fmt"
	"strings"
	"testing"
)

func TestMakeRaidMessage(t *testing.T) {
	raider := "medgelabs"
	result := MakeRaidMessage("medgelabs", 1, "medgelabs")

	if !hasTag(result, "msg-id", "raid") {
		t.Fatalf("Missing tag msg-id in Raid message. Got %s", result)
	}

	if !hasTag(result, "display-name", raider) {
		t.Fatalf("Missing tag display-name in Raid message. Got %s", result)
	}

	if !hasTag(result, "msg-param-displayName", raider) {
		t.Fatalf("Missing tag msg-param-displayName in Raid message. Got %s", result)
	}

	if !hasTag(result, "msg-param-viewerCount", "1") {
		t.Fatalf("Missing tag msg-param-viewerCount in Raid message. Got %s", result)
	}

	if !hasCommand(result, "USERNOTICE") {
		t.Fatalf("Missing command USERNOTICE. Got %s", result)
	}
}

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

func hasTag(msg, tag, tagValue string) bool {
	return strings.Contains(msg, fmt.Sprintf("%s=%s", tag, tagValue))
}

func hasCommand(msg, command string) bool {
	return strings.Contains(msg, command)
}
