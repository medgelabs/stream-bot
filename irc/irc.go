package irc

import (
	"fmt"
	"io"
	"log"
	"medgebot/bot"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const (
	MAX_MSG_SIZE = 1024 // bytes
)

// Irc client
type Irc struct {
	sync.Mutex
	channel        string
	conn           io.ReadWriteCloser
	inboundEvents  chan bot.Event
	outboundEvents chan<- bot.Event
}

func NewClient(conn io.ReadWriteCloser, channel string) *Irc {
	return &Irc{
		channel:       channel,
		conn:          conn,
		inboundEvents: make(chan bot.Event),
	}
}

func (irc *Irc) Start(nick, pass string) error {
	if err := irc.Authenticate(nick, pass); err != nil {
		return errors.Errorf("FATAL: irc authentication failure - %s", err)
	}

	if err := irc.Join(irc.channel); err != nil {
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
			if err := irc.read(); err != nil {
				log.Printf("ERROR: irc read - %v", err)
				break
			}
		}
	}()

	// Read loop for receiving messages from the bot
	go func() {
		for recv := range irc.inboundEvents {
			irc.PrivMsg(irc.channel, recv.Message)
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
		Command:  "JOIN",
		Contents: strings.TrimPrefix(channel, "#"),
	}

	return irc.write(joinCmd)
}

// Capability Request for IRC. DO NOT include the twitch.tv/ prefix
func (irc *Irc) CapReq(capability string) error {
	msg := Message{
		Command:  "CAP REQ",
		Contents: ":twitch.tv/" + capability,
	}

	return irc.write(msg)
}

// PrivMsg sends a "private message" to the IRC, no prefix attached
func (irc *Irc) PrivMsg(channel, message string) error {
	msg := Message{
		Command:  "PRIVMSG",
		Contents: channel + " :" + message,
	}

	return irc.write(msg)
}

func (irc *Irc) Close() error {
	return irc.conn.Close()
}

// SendPong reponds to the Ping heartbeat with the given body
func (irc *Irc) sendPong(body string) {
	msg := Message{
		Command:  "PONG",
		Contents: body,
	}

	if err := irc.write(msg); err != nil {
		log.Printf("ERROR: irc.send PONG - %v", err)
	}
}

// SendPass sends the PASS command to the IRC
func (irc *Irc) sendPass(token string) error {
	passCmd := Message{
		Command:  "PASS",
		Contents: token,
	}

	return irc.write(passCmd)
}

// SendNick sends the NICK command to the IRC
func (irc *Irc) sendNick(nick string) error {
	nickCmd := Message{
		Command:  "NICK",
		Contents: nick,
	}

	return irc.write(nickCmd)
}

// Read reads from the IRC stream, one line at a time
func (irc *Irc) read() error {
	buff := make([]byte, MAX_MSG_SIZE)
	len, err := irc.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "ERROR: read irc")
	}

	if len == 0 {
		log.Println("Empty message buffer")
		return errors.New("Empty message buffer")
	}

	// trace inbound IRC message
	log.Println(string(buff))

	msg := parseIrcLine(string(buff))

	// Intercept for PING/PONG
	if msg.Command == "PING" {
		irc.sendPong(msg.Contents)
		return nil
	}

	// Now, convert to bot.Event
	switch msg.Command {

	// PRIVMSG is almost always, either, a chat message or bits cheer
	case "PRIVMSG":
		if msg.IsBitsMessage() {
			evt := bot.NewBitsEvent()
			evt.Sender = msg.BitsSender()
			evt.Amount = msg.BitsAmount()
			irc.sendEvent(evt)
		} else {
			evt := bot.NewChatEvent()
			evt.Sender = msg.User
			evt.Message = msg.Contents
			irc.sendEvent(evt)
		}

	// USERNOTICE forms most event type messages to be parsed
	case "USERNOTICE":
		switch {
		case msg.IsRaidMessage():
			evt := bot.NewRaidEvent()
			evt.Sender = msg.Raider()
			evt.Amount = msg.RaidSize()
			irc.sendEvent(evt)
		case msg.IsSubscriptionMessage():
			evt := bot.NewSubEvent()
			evt.Sender = msg.Subscriber()
			evt.Amount = msg.SubMonths()
			irc.sendEvent(evt)
		case msg.IsGiftSubscriptionMessage():
			evt := bot.NewGiftSubEvent()
			evt.Sender = msg.GiftSender()
			evt.Recipient = msg.GiftRecipient()
			irc.sendEvent(evt)
		default:
			log.Printf("Unknown USERNOTICE: %s", msg.String())
		}

	default:
		// log.Printf("<<< %s", msg.String())
	}

	return nil
}

// Write a message to the IRC stream
func (irc *Irc) write(message Message) error {
	msgStr := fmt.Sprintf("%s %s", message.Command, message.Contents)

	// Lock since WriteMessage requires only one concurrent execution
	if _, err := irc.conn.Write([]byte(msgStr)); err != nil {
		return err
	}

	if message.Command != "PASS" && message.Command != "PONG" {
		log.Printf("< %s", message)
	}
	return nil
}

// bot.ChatClient
func (irc *Irc) Channel() chan<- bot.Event {
	return irc.inboundEvents
}

// bot.Client
func (irc *Irc) SetDestination(outbound chan<- bot.Event) {
	irc.outboundEvents = outbound
}

// Send Event to the bot
func (irc *Irc) sendEvent(evt bot.Event) {
	irc.outboundEvents <- evt
}
