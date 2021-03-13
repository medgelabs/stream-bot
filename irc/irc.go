package irc

import (
	"fmt"
	"io"
	"medgebot/bot"
	log "medgebot/logger"
	"sync"

	"github.com/pkg/errors"
)

const (
	// MaxMessageSize defines the maximum size of an IRC message that can be received.
	// This is used to size the message buffer slice for the io.Read method
	MaxMessageSize = 1024 // bytes
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
			if err := irc.read(); err != nil {
				log.Error(err, "irc read")
				break
			}
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
		log.Error(err, "send PASS failed")
		return err
	}
	log.Info("< PASS ***")

	if err := irc.sendNick(nick); err != nil {
		log.Error(err, "send NICK failed")
		return err
	}

	return nil
}

// Join the given IRC channel. Must be called AFTER PASS and NICK
func (irc *Irc) Join(channel string) error {
	joinCmd := Message{
		Command:  "JOIN",
		Contents: channel,
	}

	return irc.write(joinCmd)
}

// CapReq for IRC. DO NOT include the twitch.tv/ prefix
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

// Close closes the IRC connection
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
		log.Error(err, "irc.send PONG")
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
	buff := make([]byte, MaxMessageSize)
	len, err := irc.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "ERROR: read irc")
	}

	if len == 0 {
		log.Warn("Empty message buffer")
		return errors.New("Empty message buffer")
	}

	// trace inbound IRC message

	str := string(buff)
	log.Info(str)
	msg := parseIrcLine(str)

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
			log.Warn("Unknown USERNOTICE: " + msg.String())
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
		log.Info("< %s", message)
	}
	return nil
}

// bot.ChatClient

// Channel returns a channel for bot.Event used to send events to the IRC client
func (irc *Irc) Channel() chan<- bot.Event {
	return irc.inboundEvents
}

// bot.Client

// SetDestination sets the outbound channel for bot.Events the IRC client will send to
func (irc *Irc) SetDestination(outbound chan<- bot.Event) {
	irc.outboundEvents = outbound
}

// sendEvent abstracts the process to send events to the bot
func (irc *Irc) sendEvent(evt bot.Event) {
	irc.outboundEvents <- evt
}
