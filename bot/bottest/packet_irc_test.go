package bottest

import (
	"testing"
)

func TestMakeRaidMessage(t *testing.T) {
	raider := "medgelabs"
	result := MakeRaidMessage("medgelabs", 1, "medgelabs")

	if !HasTag(result, "msg-id", "raid") {
		t.Fatalf("Missing tag msg-id. Got %s", result)
	}

	if !HasTag(result, "display-name", raider) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-displayName", raider) {
		t.Fatalf("Missing tag msg-param-displayName. Got %s", result)
	}

	if !HasTag(result, "msg-param-viewerCount", "1") {
		t.Fatalf("Missing tag msg-param-viewerCount. Got %s", result)
	}

	if !HasCommand(result, "USERNOTICE") {
		t.Fatalf("Missing command USERNOTICE. Got %s", result)
	}
}

func TestMakeIrcMessage(t *testing.T) {
	tags := make(map[string]string)
	tags["display-name"] = "medgelabs"
	tags["emotes"] = ""
	tags["subscriber"] = "1"

	result := MakeIrcMessage("!hello", "medgelabs", "PRIVMSG", "medgelabs", tags)

	if !HasTag(result, "display-name", "medgelabs") {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "emotes", "") {
		t.Fatalf("Missing tag emotes. Got %s", result)
	}

	if !HasTag(result, "subscriber", "1") {
		t.Fatalf("Missing tag subscriber. Got %s", result)
	}

	if !HasCommand(result, "PRIVMSG") {
		t.Fatalf("Missing command PRIVMSG. Got %s", result)
	}
}
