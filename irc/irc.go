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

/*

(Re)Sub:

Tags: [badge-info=founder/4 badges=founder/0 color=#2E8B57 display-name=medgelabs emotes= flags= id=c1839900-5cf5-4660-8249-944ed3349173 login=medgelabs mod=0 msg-id=resub msg-param-cumulative-months=4 msg-param-months=0 msg-param-multimonth-duration=0 msg-param-multimonth-tenure=0 msg-param-should-share-streak=0 msg-param-sub-plan-name=Channel\sSubscription\s(pk_ninten) msg-param-sub-plan=1000 msg-param-was-gifted=false room-id=425558714 subscriber=1 system-msg=medgelabs\ssubscribed\sat\sTier\s1.\sThey've\ssubscribed\sfor\s4\smonths! tmi-sent-ts=1607433305720 user-id=62232210 user-type=]
Mesage: tmi.twitch.tv USERNOTICE #nintenbrews :Oh snap it's a ninten


Bits:

Tags: [badge-info=founder/4 badges=founder/0 bits=100 color=#2E8B57 display-name=medgelabs emotes= flags= id=8e275b11-22d8-4cc9-a75f-c773b5b43a3f mod=0 room-id=425558714 subscriber=0 tmi-sent-ts=1607433450079 turbo=0 user-id=62232210 user-type=]
Message: medgelabs: cheer100 HAPPY BIRTHDAY MRS. NINTEN

Raid:

Tags: [badge-info=founder/4 badges=founder/0 color=#2E8B57 display-name=medgelabs emotes= flags= id=bf27e226-96e3-4494-992d-d33cc0b2d39f login=medgelabs mod=0 msg-id=raid msg-param-displayName=medgelabs msg-param-login=medgelabs msg-param-profileImageURL=https://static-cdn.jtvnw.net/jtv_user_pictures/67a0147d-adec-4121-b5a3-a195f20e0fd0-profile_image-70x70.png msg-param-viewerCount=1 room-id=425558714 subscriber=1 system-msg=1\sraiders\sfrom\smedgelabs\shave\sjoined! tmi-sent-ts=1607438402227 user-id=62232210 user-type=]
Message: tmi.twitch.tv USERNOTICE #nintenbrews

*/

// Irc client
type Irc struct {
	sync.Mutex
	conn           io.ReadWriteCloser
	inboundEvents  chan bot.Event
	outboundEvents chan<- bot.Event
}

type Config struct {
	Scheme   string
	Host     string
	Nick     string
	Password string
	Channel  string
}

// Message represents a line of text from the IRC stream
type Message struct {
	Tags    []string
	User    string
	Command string
	Params  []string
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
		Params:  []string{channel},
	}

	return irc.write(joinCmd)
}

// Capability Request for IRC. DO NOT include the twitch.tv/ prefix
func (irc *Irc) CapReq(capability string) error {
	msg := Message{
		Command: "CAP REQ",
		Params:  []string{":twitch.tv/" + capability},
	}

	return irc.write(msg)
}

// PrivMsg sends a "private message" to the IRC, no prefix attached
func (irc *Irc) PrivMsg(channel, message string) error {
	msg := Message{
		Command: "PRIVMSG",
		Params:  []string{channel, ":" + message},
	}

	return irc.write(msg)
}

func (msg Message) String() string {
	return fmt.Sprintf("%s %s %s", msg.User, msg.Command, strings.Join(msg.Params, " "))
}

// SendPong reponds to the Ping heartbeat with the given body
func (irc *Irc) sendPong(body []string) {
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
		Params:  []string{token},
	}

	return irc.write(passCmd)
}

// SendNick sends the NICK command to the IRC
func (irc *Irc) sendNick(nick string) error {
	nickCmd := Message{
		Command: "NICK",
		Params:  []string{nick},
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

	message := string(buff)

	// TrimSpace to get rid of /r/n
	msgStr := strings.TrimSpace(string(message))
	tokens := strings.Split(msgStr, " ")

	// If tags are present, parse them out
	msg := Message{}
	cursor := 0

	// Tags
	if strings.HasPrefix(tokens[cursor], "@") {
		msg.Tags = strings.Split(strings.TrimLeft(tokens[cursor], "@"), ";")
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
	msg.Params = tokens[cursor+1:]

	// Finally, strip excess for Message
	contents := strings.TrimPrefix(strings.Join(msg.Params[1:], " "), ":")

	// Intercept for PING/PONG
	if msg.Command == "PING" {
		// log.Printf("> PING %s", contents)
		irc.sendPong(msg.Params)
		return
	}

	// TODO remove this when tested
	if msg.Command == "PRIVMSG" {
		irc.outboundEvents <- bot.Event{
			Type:    bot.CHAT_MSG,
			Sender:  msg.User,
			Message: contents,
		}
	} else {
		log.Println("<<< " + msg.String())
	}
}

// Write writes a message to the IRC stream
func (irc *Irc) write(message Message) error {
	msgStr := fmt.Sprintf("%s %s", message.Command, strings.Join(message.Params, " "))

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
