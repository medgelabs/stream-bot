package irc

import (
	"fmt"
	"io"
	"log"
	"medgebot/bot"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	MAX_MSG_SIZE = 1024 // bytes

	MSG_RAID = iota
	MSG_SUB
	MSG_GIFTSUB
	MSG_BITS
	MSG_CHAT
	MSG_SYSTEM // mostly for messages like PING/PONG
)

// Irc client
type Irc struct {
	sync.Mutex
	conn           io.ReadWriteCloser
	inboundEvents  chan bot.Event
	outboundEvents chan<- bot.Event
}

type Config struct {
	Nick     string
	Password string
	Channel  string
}

// Message represents a line of text from the IRC stream
type Message struct {
	Tags    map[string]string
	User    string
	Command string
	Params  string // aka message content
}

// Tag returns a tag on the message, or an empty string and a bool indicating if the
// tag was found
func (msg Message) Tag(tag string) string {
	return msg.Tags[tag]
}

func (msg *Message) AddTag(tag, value string) {
	msg.Tags[tag] = value
}

func NewClient(conn io.ReadWriteCloser) *Irc {
	return &Irc{
		conn:          conn,
		inboundEvents: make(chan bot.Event),
	}
}

func (irc *Irc) Start(config Config) error {
	if err := irc.Authenticate(config.Nick, config.Password); err != nil {
		return errors.Errorf("FATAL: irc authentication failure - %s", err)
	}

	if err := irc.Join(config.Channel); err != nil {
		return errors.Errorf("FATAL: irc join channel failed: %s", err)
	}

	// Command Capability Request for UserNotices (raids, subs, etc)
	if err := irc.CapReq("commands"); err != nil {
		return errors.Errorf("FATAL: irc CapReq COMMANDS failed: %s", err)
	}

	// Tags for bits/subs/raids metadata
	if err := irc.CapReq("tags"); err != nil {
		return errors.Errorf("FATAL: irc CapReq TAGS failed: %s", err)
	}

	// Read loop for receiving messages from IRC
	go func() {
		for {
			irc.read()
		}
	}()

	// Read loop for receiving messages from the bot
	go func() {
		for recv := range irc.inboundEvents {
			irc.PrivMsg(config.Channel, recv.Message)
		}
	}()

	return nil
}

// Authenticate connects to the IRC stream with the given nick and password
func (irc *Irc) Authenticate(nick, password string) error {
	if err := irc.sendPass(password); err != nil {
		log.Printf("ERROR: send PASS failed: %s", err)
		return err
	}
	log.Println("< PASS ***")

	if err := irc.sendNick(nick); err != nil {
		log.Printf("ERROR: send NICK failed: %s", err)
		return err
	}

	return nil
}

// Join the given IRC channel. Must be called AFTER PASS and NICK
func (irc *Irc) Join(channel string) error {
	joinCmd := Message{
		Command: "JOIN",
		Params:  channel,
	}

	return irc.write(joinCmd)
}

// Capability Request for IRC. DO NOT include the twitch.tv/ prefix
func (irc *Irc) CapReq(capability string) error {
	msg := Message{
		Command: "CAP REQ",
		Params:  ":twitch.tv/" + capability,
	}

	return irc.write(msg)
}

// PrivMsg sends a "private message" to the IRC, no prefix attached
func (irc *Irc) PrivMsg(channel, message string) error {
	msg := Message{
		Command: "PRIVMSG",
		Params:  channel + " :" + message,
	}

	return irc.write(msg)
}

func (msg Message) String() string {
	return fmt.Sprintf("%s %s %s", msg.User, msg.Command, msg.Params)
}

// SendPong reponds to the Ping heartbeat with the given body
func (irc *Irc) sendPong(body string) {
	msg := Message{
		Command: "PONG",
		Params:  body,
	}

	if err := irc.write(msg); err != nil {
		log.Printf("ERROR: irc.send PONG - %v", err)
	}
}

func (irc *Irc) Close() error {
	return irc.conn.Close()
}

// SendPass sends the PASS command to the IRC
func (irc *Irc) sendPass(token string) error {
	passCmd := Message{
		Command: "PASS",
		Params:  token,
	}

	return irc.write(passCmd)
}

// SendNick sends the NICK command to the IRC
func (irc *Irc) sendNick(nick string) error {
	nickCmd := Message{
		Command: "NICK",
		Params:  nick,
	}

	return irc.write(nickCmd)
}

// Read reads from the IRC stream, one line at a time
func (irc *Irc) read() {
	buff := make([]byte, MAX_MSG_SIZE)
	len, err := irc.conn.Read(buff)
	if err != nil {
		log.Printf("ERROR: read irc - %v", err)
		return
	}

	if len == 0 {
		log.Println("Empty message buffer")
		return
	}

	msg := parseIrcLine(string(buff))

	// Intercept for PING/PONG
	if msg.Command == "PING" {
		// log.Printf("> PING %s", contents)
		irc.sendPong(msg.Params)
		return
	}

	// Now, convert to bot.Event
	switch msg.Command {

	// PRIVMSG is almost always a chat message
	case "PRIVMSG":
		irc.outboundEvents <- bot.Event{
			Type:    bot.CHAT_MSG,
			Sender:  msg.User,
			Message: msg.Params,
		}

	// USERNOTICE forms most event type messages to be parsed
	case "USERNOTICE":
		msgType := parseMessageType(msg)

		switch msgType {
		case MSG_RAID:
			irc.outboundEvents <- parseRaidMessage(msg)
		case MSG_BITS:
			irc.outboundEvents <- parseBitsMessage(msg)
		default:
			log.Printf("Unknown USERNOTICE: %s", msg.String())
		}

	default:
		log.Printf("<<< %s", msg.String())
	}
}

// Attempt to parse a line from IRC to a Message
func parseIrcLine(message string) Message {
	// TrimSpace to get rid of /r/n
	msgStr := strings.TrimSpace(string(message))
	tokens := strings.Split(msgStr, " ")

	// If tags are present, parse them out
	msg := Message{}
	cursor := 0

	// Tags
	if strings.HasPrefix(tokens[cursor], "@") {
		tagsSlice := strings.Split(strings.TrimLeft(tokens[cursor], "@"), ";")
		for _, tag := range tagsSlice {
			parts := strings.Split(tag, "=")
			msg.AddTag(parts[0], parts[1])
		}

		log.Println(msg.Tags)
		cursor++
	}

	// Prefix, therefore parse username
	if strings.HasPrefix(tokens[cursor], ":") {
		rawUsername := strings.Split(tokens[cursor], ":")[1]
		username := strings.Split(rawUsername, "!")[0]
		msg.User = username
		cursor++
	}

	// Remaining cursor points should be Command and Params
	msg.Command = tokens[cursor]

	// The combined String beyond the ":"
	msg.Params = strings.TrimPrefix(strings.Join(tokens[cursor+1:], " "), ":")
	return msg
}

// Parse a msgType from Tags to one of our iota constants, or MSG_SYSTEM if
// unknown
func parseMessageType(msg Message) int {

	// Bits doesn't use msg-id, so we must check for that explicitly
	if bits := msg.Tag("bits"); bits != "" {
		return MSG_BITS
	}

	msgType := msg.Tag("msg-id")
	switch msgType {
	case "raid":
		return MSG_RAID
	default:
		return MSG_SYSTEM
	}
}

func parseRaidMessage(msg Message) bot.Event {
	raider := msg.Tag("msg-param-displayName")
	raidSizeStr := msg.Tag("msg-param-viewerCount")
	raidSize, err := strconv.Atoi(raidSizeStr)
	if err != nil {
		log.Printf("ERROR: invalid raid size [%s], defaulting to 0", raidSizeStr)
		raidSize = 0
	}

	return bot.Event{
		Type:   bot.RAID,
		Sender: raider,
		Amount: raidSize,
	}
}

func parseBitsMessage(msg Message) bot.Event {
	cheerer := msg.Tag("display-name")
	bitsStr := msg.Tag("bits")
	amount, err := strconv.Atoi(bitsStr)
	if err != nil {
		log.Printf("ERROR: invalid bits [%s], defaulting to 0", bitsStr)
	}

	return bot.Event{
		Type:   bot.BITS,
		Sender: cheerer,
		Amount: amount,
	}
}

// Write writes a message to the IRC stream
func (irc *Irc) write(message Message) error {
	msgStr := fmt.Sprintf("%s %s", message.Command, message.Params)

	// Lock since WriteMessage requires only one concurrent execution
	if _, err := irc.conn.Write([]byte(msgStr)); err != nil {
		return err
	}

	if message.Command != "PASS" && message.Command != "PONG" {
		log.Printf("< %s", message)
	}
	return nil
}

// Pluggable
func (irc *Irc) GetChannel() chan<- bot.Event {
	return irc.inboundEvents
}

func (irc *Irc) SetChannel(outbound chan<- bot.Event) {
	irc.outboundEvents = outbound
}
