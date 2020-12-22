package bottest

import (
	"strings"
)

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
	sb.WriteString(sender + "!" + sender + "@" + sender + ".tmi.twitch.tv")
	sb.WriteString(" ")

	// Command
	sb.WriteString(strings.ToUpper(command))
	sb.WriteString(" ")

	// Channel
	sb.WriteString("#" + channel)
	sb.WriteString(" :")

	// Message
	sb.WriteString(body)

	return sb.String()
}
