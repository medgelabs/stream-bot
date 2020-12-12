package irc

import (
	"reflect"
	"strconv"
	"testing"
)

const (
	// User that should be used in below IRC messages
	assistant = "assistant1"

	CHAT_MSG    = "@display-name=assistant1;subscriber=1 :assistant1!assistant1@assistant1.tmi.twitch.tv PRIVMSG #medgelabs :Yes, we can test"
	BITS_MSG    = "@bits=1;display-name=assistant1 :assistant1!assistant1@assistant1.tmi.twitch.tv PRIVMSG #medgelabs :Cheer100"
	RAID_MSG    = `@display-name=assistant1;msg-id=raid;msg-param-displayName=assistant1;msg-param-viewerCount=1 :tmi.twitch.tv USERNOTICE #medgelabs`
	SUB_MSG     = `@display-name=assistant1;msg-id=sub;msg-param-cumulative-months=1;msg-param-sub-plan=Tier1;msg-param-sub-plan-name=Tier1 :tmi.twitch.tv USERNOTICE #medgelabs :Moar testing!`
	RESUB_MSG   = `@display-name=assistant1;msg-id=resub;msg-param-cumulative-months=2;msg-param-sub-plan=Tier1;msg-param-sub-plan-name=Tier1 :tmi.twitch.tv USERNOTICE #medgelabs :Moar testing!`
	GIFTSUB_MSG = `@display-name=ReallyFrank;msg-id=subgift;msg-param-gift-months=1;msg-param-recipient-display-name=Fjoell;msg-param-recipient-user-name=Fjoell;msg-param-sub-plan=1000;msg-param-sub-plan-name=Tier1 :tmi.twitch.tv USERNOTICE #medgelabs :`
)

func TestParseIrcLine(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    Message
	}{
		{description: "Bits Message", input: BITS_MSG, expected: Message{
			Tags: map[string]string{
				"display-name": assistant,
				"bits":         "1",
			},
			User:     assistant,
			Command:  "PRIVMSG",
			Channel:  "medgelabs",
			Contents: "Cheer100",
		}},
	}

	for _, test := range tests {
		t.Run(test.description, func(tt *testing.T) {
			result := parseIrcLine(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				tt.Fatalf("Parsed Message invalid. Expected:\n %+v \nGot:\n %+v", test.expected, result)
			}
		})
	}
}

func TestChatMessageParse(t *testing.T) {
	parsed := parseIrcLine(CHAT_MSG)

	if parsed.Command != "PRIVMSG" {
		t.Fatalf("Expected PRIVMSG as Command but got %s", parsed.Command)
	}

	if parsed.User != assistant {
		t.Fatalf("Expected assistant1 as User but got %s", parsed.User)
	}

	if parsed.Channel != "medgelabs" {
		t.Fatalf("Expected medgelabs as Channel but got %s", parsed.Channel)
	}

	if parsed.Contents != "Yes, we can test" {
		t.Fatalf("Expected 'Yes, we can test' as Contents but got %s", parsed.Contents)
	}

	// Check for a couple tags to ensure parsing happened properly
	name := parsed.Tag("display-name")
	if name != assistant {
		t.Fatalf("Expected assistant1 as display-name Tag but got %s", name)
	}

	subscriber := parsed.Tag("subscriber")
	subInt, err := strconv.Atoi(subscriber)
	if err != nil || subInt != 1 {
		t.Fatalf("Expected 1 subscriber Tag but got %s", subscriber)
	}

}

func TestRaidMessageParsing(t *testing.T) {
	parsed := parseIrcLine(RAID_MSG)

	if !parsed.IsRaidMessage() {
		t.Fatalf("RAID_MSG not recognized as a Raid")
	}

	if parsed.Raider() != assistant {
		t.Fatalf("Raider should be %s, but got %s", assistant, parsed.Raider())
	}

	if parsed.RaidSize() != 1 {
		t.Fatalf("RaidSize should be 1, but got %d", parsed.RaidSize())
	}
}

func TestBitsMessageParsing(t *testing.T) {
	parsed := parseIrcLine(BITS_MSG)

	if !parsed.IsBitsMessage() {
		t.Fatalf("BITS_MSG not recognized as Bits cheering")
	}

	if parsed.BitsSender() != assistant {
		t.Fatalf("BitsSender should be %s, but got %s", assistant, parsed.BitsSender())
	}

	if parsed.BitsAmount() != 1 {
		t.Fatalf("BitsAmount should be 1, but got %d", parsed.BitsAmount())
	}
}

func TestParseUserNoticeMessageType(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    int
	}{
		{description: "Raid Message", input: RAID_MSG, expected: MSG_RAID},
		{description: "Sub Message", input: SUB_MSG, expected: MSG_SUB},
		{description: "ReSub Message", input: RESUB_MSG, expected: MSG_SUB},
		{description: "GiftSub Message", input: GIFTSUB_MSG, expected: MSG_GIFTSUB},
	}

	for _, test := range tests {
		t.Run(test.description, func(tt *testing.T) {
			result := parseIrcLine(test.input).parseUserNoticeMessageType()
			if result != test.expected {
				tt.Fatalf("Expected %d message type, but got %d for %s", test.expected, result, test.description)
			}
		})
	}
}
