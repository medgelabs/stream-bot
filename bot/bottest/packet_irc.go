package bottest

import (
	"fmt"
	"strconv"
	"strings"
)

func MakeRaidMessage(raider string, raidSize int, channel string) string {
	tags := make(map[string]string)

	tags["msg-id"] = "raid"
	tags["display-name"] = raider
	tags["msg-param-displayName"] = raider
	tags["msg-param-viewerCount"] = strconv.Itoa(raidSize)

	return MakeIrcMessage("", "", "USERNOTICE", channel, tags)
}

// Helper for creating an IRC message to be used in Integration tests
func MakeIrcMessage(body, sender, command, channel string, tags map[string]string) string {
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
