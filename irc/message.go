package irc

import (
	"log"
	"strconv"
)

const (
	MSG_RAID = iota
	MSG_SUB
	MSG_GIFTSUB
	MSG_BITS
	MSG_CHAT
	MSG_SYSTEM // mostly for messages like PING/PONG
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
	if msg.Tags == nil {
		msg.Tags = make(map[string]string)
	}

	msg.Tags[tag] = value
}

// Parse a msgType from Tags on a USERNOTICE to one of our iota constants, or MSG_SYSTEM if
// unknown
func (msg Message) parseMessageType() int {
	msgType := msg.Tag("msg-id")

	switch msgType {
	case "raid":
		return MSG_RAID
	default:
		return MSG_SYSTEM
	}
}

// Check if message is a Raid message
func (msg Message) IsRaidMessage() bool {
	return msg.parseMessageType() == MSG_RAID
}

func (msg Message) Raider() string {
	return msg.Tag("msg-param-displayName")
}

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

func (msg Message) BitsSender() string {
	return msg.Tag("display-name")
}

func (msg Message) BitsAmount() int {
	bitsStr := msg.Tag("bits")
	amount, err := strconv.Atoi(bitsStr)
	if err != nil {
		log.Printf("ERROR: invalid bits [%s], defaulting to 0", bitsStr)
		return 0
	}

	return amount
}
