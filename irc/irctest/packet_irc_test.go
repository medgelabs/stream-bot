package irctest

import (
	"testing"
)

const (
	user    = "medgelabs"
	channel = "medgelabs"
)

func TestMakeRaidMessage(t *testing.T) {
	result := MakeRaidMessage("medgelabs", 1, "medgelabs")

	if !HasTag(result, "msg-id", "raid") {
		t.Fatalf("Missing tag msg-id. Got %s", result)
	}

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-displayName", user) {
		t.Fatalf("Missing tag msg-param-displayName. Got %s", result)
	}

	if !HasTag(result, "msg-param-viewerCount", "1") {
		t.Fatalf("Missing tag msg-param-viewerCount. Got %s", result)
	}

	if !HasCommand(result, "USERNOTICE") {
		t.Fatalf("Missing command USERNOTICE. Got %s", result)
	}
}

func TestMakeChatMessage(t *testing.T) {
	result := MakeChatMessage(user, "Hello", channel)

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasCommand(result, "PRIVMSG") {
		t.Fatalf("Missing command PRIVMSG. Got %s", result)
	}
}

func TestMakeBitsMessage(t *testing.T) {
	result := MakeBitsMessage(user, 100, channel)

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "bits", "100") {
		t.Fatalf("Missing tag bits. Got %s", result)
	}

	if !HasCommand(result, "PRIVMSG") {
		t.Fatalf("Missing command PRIVMSG. Got %s", result)
	}
}

func TestMakeSubMessage(t *testing.T) {
	result := MakeSubMessage(user, 1, channel)

	if !HasTag(result, "msg-id", "sub") {
		t.Fatalf("Missing tag msg-id. Got %s", result)
	}

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-cumulative-months", "1") {
		t.Fatalf("Missing tag msg-param-cumulative-months. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan-name", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan-name. Got %s", result)
	}

	if !HasCommand(result, "USERNOTICE") {
		t.Fatalf("Missing command USERNOTICE. Got %s", result)
	}
}

func TestMakeReSubMessage(t *testing.T) {
	result := MakeResubMessage(user, 1, channel)

	if !HasTag(result, "msg-id", "resub") {
		t.Fatalf("Missing tag msg-id. Got %s", result)
	}

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-cumulative-months", "1") {
		t.Fatalf("Missing tag msg-param-cumulative-months. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan-name", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan-name. Got %s", result)
	}

	if !HasCommand(result, "USERNOTICE") {
		t.Fatalf("Missing command USERNOTICE. Got %s", result)
	}
}

func TestMakeGiftSubMessage(t *testing.T) {
	result := MakeGiftSubMessage(user, "recipient", channel)

	if !HasTag(result, "msg-id", "subgift") {
		t.Fatalf("Missing tag msg-id. Got %s", result)
	}

	if !HasTag(result, "display-name", user) {
		t.Fatalf("Missing tag display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-recipient-display-name", "recipient") {
		t.Fatalf("Missing tag msg-param-recipient-display-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-recipient-user-name", "recipient") {
		t.Fatalf("Missing tag msg-param-recipient-user-name. Got %s", result)
	}

	if !HasTag(result, "msg-param-gift-months", "1") {
		t.Fatalf("Missing tag msg-param-gift-months. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan. Got %s", result)
	}

	if !HasTag(result, "msg-param-sub-plan-name", "Tier1") {
		t.Fatalf("Missing tag msg-param-sub-plan-name. Got %s", result)
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

	result := makeIrcMessage("!hello", "medgelabs", "PRIVMSG", "medgelabs", tags)

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
