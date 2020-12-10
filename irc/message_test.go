package irc

import (
	"strconv"
	"testing"
)

const (
	// User that should be used in below IRC messages
	assistant = "assistant1"

	CHAT_MSG = "@badge-info=subscriber/4;badges=subscriber/3;color=#8A2BE2;display-name=assistant1;emotes=;flags=;id=a6db5e1c-6969-41da-8447-29e347fa84d3;mod=0;room-id=000000001;subscriber=1;tmi-sent-ts=1607601409666;turbo=0;user-id=000000001;user-type= :assistant1!assistant1@assistant1.tmi.twitch.tv PRIVMSG #medgelabs :Yes, we can test"

	BITS_MSG = "@badge-info=subscriber/10;badges=vip/1,subscriber/9,bits-leader/3;bits=1;color=#008000;display-name=assistant1;emotes=;flags=;id=f1b5a719-b65f-466d-b38a-5d48681b870e;mod=0;room-id=000000001;subscriber=1;tmi-sent-ts=1607601525765;turbo=0;user-id=000000002;user-type= :assistant1!assistant1@assistant1.tmi.twitch.tv PRIVMSG #medgelabs :Cheer100"

	RAID_MSG = `@badge-info=subscriber/6;badges=moderator/1,subscriber/6,bits/100000;color=#2E8B57;display-name=assistant1;emotes=;flags=;id=9c764aba-ed5b-4fc1-9b43-cad3b5952829;login=assistant1;mod=0;msg-id=raid;msg-param-displayName=assistant1;msg-param-login=assistant1;msg-param-profileImageURL=https://static-cdn.jtvnw.net/jtv_user_pictures/noop.png;msg-param-viewerCount=1;room-id=000000001;subscriber=1;system-msg=1\sassistant1s\sfrom\sassistant1\shave\sjoined!;tmi-sent-ts=1607601582443;user-id=000000003;user-type=mod :tmi.twitch.tv USERNOTICE #medgelabs`
)

func TestParseIrcLine(t *testing.T) {
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
	}

	for _, test := range tests {
		result := parseIrcLine(test.input).parseUserNoticeMessageType()
		if result != test.expected {
			t.Fatalf("Expected %d message type, but got %d for %s", test.expected, result, test.description)
		}
	}
}
