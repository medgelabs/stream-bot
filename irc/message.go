package irc

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	MSG_RAID = iota
	MSG_SUB
	MSG_GIFTSUB
	MSG_BITS
	MSG_CHAT
)

// Message represents a line of text from the IRC stream
type Message struct {
	Tags     map[string]string
	User     string
	Command  string
	Channel  string
	Contents string
}

func NewMessage() Message {
	return Message{
		Tags: make(map[string]string),
	}
}

// Tag returns a tag on the message, or an empty string and a bool indicating if the
// tag was found
func (msg Message) Tag(tag string) string {
	return msg.Tags[tag]
}

// AddTag populates the given key/value pair
func (msg *Message) AddTag(tag, value string) {
	msg.Tags[tag] = value
}

// Attempt to parse a line from IRC to a Message
func parseIrcLine(message string) Message {
	if strings.TrimSpace(message) == "" {
		log.Println("WARN: attempt to parse empty IRC line")
		return NewMessage()
	}

	// TrimSpace to get rid of /r/n
	msgStr := strings.TrimSpace(string(message))
	tokens := strings.Split(msgStr, " ")

	// If tags are present, parse them out
	msg := NewMessage()
	cursor := 0

	// Tags
	if strings.HasPrefix(tokens[cursor], "@") {
		tagsSlice := strings.Split(strings.TrimLeft(tokens[cursor], "@"), ";")
		for _, tag := range tagsSlice {
			parts := strings.Split(tag, "=")

			// Random chance of empty tag? No panicking
			if len(parts) < 2 {
				continue
			}

			msg.AddTag(parts[0], parts[1])
		}
		cursor++
	}

	// Parse Username, if present
	if strings.HasPrefix(tokens[cursor], ":") {

		nameTag := msg.Tag("display-name")
		if nameTag != "" {
			msg.User = nameTag
		} else {
			// Parse from the prefix before the command (not ideal)
			rawUsername := strings.Split(tokens[cursor], ":")[1]
			username := strings.Split(rawUsername, "!")[0]
			msg.User = username
		}

		cursor++
	}

	// Next cursor point should be Command
	msg.Command = tokens[cursor]
	cursor++

	// TODO this caused a panic on some kind of message
	// Then, Channel
	msg.Channel = strings.TrimPrefix(tokens[cursor], "#")
	cursor++

	// The rest should be the combined String beyond the ":"
	combinedContents := strings.Join(tokens[cursor:], " ")
	msg.Contents = strings.TrimPrefix(combinedContents, ":")
	return msg
}

// Parse a msgType from Tags on a USERNOTICE to one of our iota constants, or MSG_SYSTEM if
// unknown
func (msg Message) parseUserNoticeMessageType() int {
	msgType := msg.Tag("msg-id")

	switch msgType {
	case "raid":
		return MSG_RAID
	case "sub", "resub":
		return MSG_SUB
	case "subgift":
		return MSG_GIFTSUB
	default:
		return MSG_CHAT
	}
}

// Check if message is a Raid message
func (msg Message) IsRaidMessage() bool {
	return msg.parseUserNoticeMessageType() == MSG_RAID
}

// Return the Raider for a Raid message
func (msg Message) Raider() string {
	return msg.Tag("msg-param-displayName")
}

// Return the Raid Size for a Raid message
func (msg Message) RaidSize() int {
	raidSizeStr := msg.Tag("msg-param-viewerCount")
	raidSize, err := strconv.Atoi(raidSizeStr)
	if err != nil {
		log.Printf("ERROR: invalid raid size [%s], defaulting to 0", raidSizeStr)
		return 0
	}

	return raidSize
}

// Check if message is a Bits message
func (msg Message) IsBitsMessage() bool {
	_, err := strconv.Atoi(msg.Tag("bits"))
	return err == nil
}

// Return the User that donated bits in a Bits message
func (msg Message) BitsSender() string {
	return msg.Tag("display-name")
}

// Return the amount of bits donated in a Bits message
func (msg Message) BitsAmount() int {
	bitsStr := msg.Tag("bits")
	amount, err := strconv.Atoi(bitsStr)
	if err != nil {
		log.Printf("ERROR: invalid bits [%s], defaulting to 0", bitsStr)
		return 0
	}

	return amount
}

// Check if message is a Sub/Resub message
func (msg Message) IsSubscriptionMessage() bool {
	return msg.parseUserNoticeMessageType() == MSG_SUB
}

// Return the Subscriber for a Subscription message
func (msg Message) Subscriber() string {
	return msg.Tag("display-name") // because Twitch is inconsistent
}

// Return the number of months subscribed for a Subscription message
func (msg Message) SubMonths() int {
	monthsStr := msg.Tag("msg-param-cumulative-months")
	months, err := strconv.Atoi(monthsStr)
	if err != nil {
		log.Printf("ERROR: invalid months [%s], defaulting to 0", monthsStr)
		return 0
	}

	return months
}

// Check if message is a Sub/Resub message
func (msg Message) IsGiftSubscriptionMessage() bool {
	return msg.parseUserNoticeMessageType() == MSG_GIFTSUB
}

// Return the Recipient of a Gift Subscription
func (msg Message) GiftRecipient() string {
	return msg.Tag("msg-param-recipient-display-name")
}

// Return the Sender of a Gift Subscription
func (msg Message) GiftSender() string {
	return msg.Tag("display-name")
}

func (msg Message) String() string {
	return fmt.Sprintf("%s %s %s #%s :%s", msg.Tags, msg.User, msg.Command, msg.Channel, msg.Contents)
}
