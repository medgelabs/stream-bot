package irc

import (
	"medgebot/irc/irctest"
	"reflect"
	"testing"
)

const (
	// User that should be used in below IRC messages
	assistant = "assistant1"
	channel   = "medgelabs"

	CHAT_MSG_BASE = ":assistant1!assistant1@assistant1.tmi.twitch.tv PRIVMSG #medgelabs :Yes, we can test"
)

func TestTags(t *testing.T) {
	msg := NewMessage()

	msg.AddTag("display-name", assistant)
	if msg.Tag("display-name") != assistant {
		t.Fatalf("Expected display-name tag to be %s. Got: %s", assistant, msg.Tag("display-name"))
	}
}

func TestParseIrcLine(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    Message
	}{
		{description: "Chat Message with display-name tag", input: irctest.MakeChatMessage(assistant, "Yes, we can test", channel), expected: Message{
			Tags: map[string]string{
				"display-name": assistant,
			},
			User:     assistant,
			Command:  "PRIVMSG",
			Channel:  channel,
			Contents: "Yes, we can test",
		}},
		{description: "Chat Message with no display-name tag should still parse", input: CHAT_MSG_BASE, expected: Message{
			Tags:     map[string]string{},
			User:     "assistant1",
			Command:  "PRIVMSG",
			Channel:  channel,
			Contents: "Yes, we can test",
		}},
		{description: "Bits Message", input: irctest.MakeBitsMessage(assistant, 1, "medgelabs"), expected: Message{
			Tags: map[string]string{
				"display-name": assistant,
				"bits":         "1",
			},
			User:     assistant,
			Command:  "PRIVMSG",
			Channel:  channel,
			Contents: "Cheer1",
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

func TestRaidMessageParsing(t *testing.T) {
	parsed := parseIrcLine(irctest.MakeRaidMessage(assistant, 1, channel))

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

func TestRaidMessageParsingInvalidRaidSize(t *testing.T) {
	// Invalid raid size value should default to 0
	invalid := parseIrcLine(irctest.MakeRaidMessage(assistant, 1, channel))
	invalid.AddTag("msg-param-viewerCount", "asdf") // Invalid raid size
	if invalid.RaidSize() != 0 {
		t.Fatalf("Invalid raid size should have defaulted to 0. Got: %d", invalid.RaidSize())
	}
}

func TestBitsMessageParsing(t *testing.T) {
	parsed := parseIrcLine(irctest.MakeBitsMessage(assistant, 1, "medgelabs"))

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

func TestSubMessageParsing(t *testing.T) {
	parsed := parseIrcLine(irctest.MakeSubMessage(assistant, 1, "medgelabs"))

	if !parsed.IsSubscriptionMessage() {
		t.Fatalf("Subscription Message not recognized as sub event")
	}

	if parsed.Subscriber() != assistant {
		t.Fatalf("Subscriber should be %s, but got %s", assistant, parsed.Subscriber())
	}

	if parsed.SubMonths() != 1 {
		t.Fatalf("SubMonths should be 1, but got %d", parsed.SubMonths())
	}
}

func TestResubMessageParsing(t *testing.T) {
	parsed := parseIrcLine(irctest.MakeResubMessage(assistant, 2, "medgelabs"))

	if !parsed.IsSubscriptionMessage() {
		t.Fatalf("Re-Subscription Message not recognized as sub event")
	}

	if parsed.Subscriber() != assistant {
		t.Fatalf("Subscriber should be %s, but got %s", assistant, parsed.Subscriber())
	}

	if parsed.SubMonths() != 2 {
		t.Fatalf("SubMonths should be 2, but got %d", parsed.SubMonths())
	}
}

func TestGiftSubMessageParsing(t *testing.T) {
	parsed := parseIrcLine(irctest.MakeGiftSubMessage("ReallyFrank", "Fjoell", "medgelabs"))

	if !parsed.IsGiftSubscriptionMessage() {
		t.Fatalf("Gift Subscription Message not recognized as subgift event")
	}

	if parsed.GiftRecipient() != "Fjoell" {
		t.Fatalf("Recipient should be Fjoell, but got %s", parsed.GiftRecipient())
	}

	if parsed.GiftSender() != "ReallyFrank" {
		t.Fatalf("Sender should be ReallyFrank, but got %s", parsed.GiftSender())
	}
}

func TestParseUserNoticeMessageType(t *testing.T) {
	tests := []struct {
		description string
		input       string
		expected    int
	}{
		{description: "Chat Message", input: irctest.MakeChatMessage(assistant, "Hello", channel), expected: MSG_CHAT},
		{description: "Raid Message", input: irctest.MakeRaidMessage(assistant, 1, channel), expected: MSG_RAID},
		{description: "Sub Message", input: irctest.MakeSubMessage(assistant, 1, channel), expected: MSG_SUB},
		{description: "ReSub Message", input: irctest.MakeResubMessage(assistant, 2, channel), expected: MSG_SUB},
		{description: "GiftSub Message", input: irctest.MakeGiftSubMessage("ReallyFrank", "Fjoell", channel), expected: MSG_GIFTSUB},
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

// Trailing semicolon at the end of the tags block would cause an empty tag to be registered, which
// caused the `msg.AddTag(parts[0], parts[1])` part to panic (index out of bounds).
// This test should not panic in such an event
func TestEmptyTagDoesntExplode(t *testing.T) {
	line := "@display-name=medgelabs; tmi.twitch.tv PRIVMSG :Trailing semicolon causes empty tag. Should not explode"
	_ = parseIrcLine(line)
}
