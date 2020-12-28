package irctest

import (
	"fmt"
	"strconv"
	"strings"
)

// MakeRaidMessage generates a well-formed Raid IRC message
func MakeRaidMessage(raider string, raidSize int, channel string) string {
	tags := make(map[string]string)

	tags["msg-id"] = "raid"
	tags["display-name"] = raider
	tags["msg-param-displayName"] = raider
	tags["msg-param-viewerCount"] = strconv.Itoa(raidSize)

	return makeIrcMessage("", "", "USERNOTICE", channel, tags)
}

// MakeChatMessage generates a well-formed simple Chat IRC message
func MakeChatMessage(sender, content, channel string) string {
	tags := make(map[string]string)
	tags["display-name"] = sender

	return makeIrcMessage(sender, content, "PRIVMSG", channel, tags)
}

// MakeBitsMessage generates a well-formed Bits Cheer IRC message
func MakeBitsMessage(sender string, bits int, channel string) string {
	tags := make(map[string]string)
	tags["display-name"] = sender
	tags["bits"] = strconv.Itoa(bits)

	return makeIrcMessage(sender, fmt.Sprintf("Cheer%d", bits), "PRIVMSG", channel, tags)
}

// MakeSubMessage generates a well-formed Subscription event IRC message
func MakeSubMessage(subscriber string, months int, channel string) string {
	tags := make(map[string]string)
	tags["msg-id"] = "sub"
	tags["display-name"] = subscriber
	tags["msg-param-cumulative-months"] = strconv.Itoa(months)
	tags["msg-param-sub-plan"] = "Tier1"
	tags["msg-param-sub-plan-name"] = "Tier1"

	return makeIrcMessage("", "", "USERNOTICE", channel, tags)
}

// MakeResubMessage generates a well-formed Re-Subscription event IRC message
func MakeResubMessage(subscriber string, months int, channel string) string {
	tags := make(map[string]string)
	tags["msg-id"] = "resub"
	tags["display-name"] = subscriber
	tags["msg-param-cumulative-months"] = strconv.Itoa(months)
	tags["msg-param-sub-plan"] = "Tier1"
	tags["msg-param-sub-plan-name"] = "Tier1"

	return makeIrcMessage("", "", "USERNOTICE", channel, tags)
}

// MakeGiftSubMessage generates a well-formed Gift Subscription event IRC message
func MakeGiftSubMessage(sender, recipient, channel string) string {
	tags := make(map[string]string)
	tags["msg-id"] = "subgift"
	tags["display-name"] = sender
	tags["msg-param-gift-months"] = "1"
	tags["msg-param-recipient-display-name"] = recipient
	tags["msg-param-recipient-user-name"] = recipient
	tags["msg-param-sub-plan"] = "Tier1"
	tags["msg-param-sub-plan-name"] = "Tier1"

	return makeIrcMessage("", "", "USERNOTICE", channel, tags)
}

// Helper for creating an IRC message
func makeIrcMessage(sender, body, command, channel string, tags map[string]string) string {
	var sb strings.Builder

	// Tags
	sb.WriteString("@")
	for k, v := range tags {
		sb.WriteString(k + "=" + v + ";")
	}
	sb.WriteString(" :")

	// Username
	if sender != "" {
		sb.WriteString(sender + "!" + sender + "@" + sender + ".tmi.twitch.tv")
	} else {
		sb.WriteString("tmi.twitch.tv")
	}
	sb.WriteString(" ")

	// Command
	sb.WriteString(strings.ToUpper(command))
	sb.WriteString(" ")

	// Channel
	sb.WriteString("#" + channel)

	// Message
	if body != "" {
		sb.WriteString(" :")
		sb.WriteString(body)
	}

	str := sb.String()
	// Removes trailing ; from last tag
	str = strings.Replace(str, "; :", " :", 1)
	return str
}

// HasTag helper for determining if an IRC message has the given tag/tagValue
func HasTag(msg, tag, tagValue string) bool {
	return strings.Contains(msg, fmt.Sprintf("%s=%s", tag, tagValue))
}

// HasCommand helper for determining if an IRC is the given Command
func HasCommand(msg, command string) bool {
	return strings.Contains(msg, command)
}
