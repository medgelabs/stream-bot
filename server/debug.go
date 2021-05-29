package server

import (
	"medgebot/bot"
	"net/http"
)

func (s *Server) debugSub(client *DebugClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendSub()
	}
}

func (s *Server) debugGift(client *DebugClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendGiftSub()
	}
}

func (s *Server) debugBit(client *DebugClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendBit()
	}
}

// DebugClient is a dummy client for the Bot that we use to send messages
type DebugClient struct {
	events chan<- bot.Event
}

// SetDestination from bot.Client
func (c *DebugClient) SetDestination(events chan<- bot.Event) {
	c.events = events
}

// SendBit sends a mock Bit event to the Bot for testing
func (c *DebugClient) SendBit() {
	evt := bot.NewBitsEvent()
	evt.Sender = "Przemko9856"
	evt.Amount = 100
	evt.Message = "I am a willing test subject"

	c.events <- evt
}

// SendSub sends a mock Subscription event to the Bot for testing
func (c *DebugClient) SendSub() {
	evt := bot.NewSubEvent()
	evt.Sender = "saltymoth"
	evt.Amount = 5
	evt.Message = "I am a willing test subject"

	c.events <- evt
}

// SendGiftSub sends a mock Gift Sub event to the Bot for testing
func (c *DebugClient) SendGiftSub() {
	evt := bot.NewGiftSubEvent()
	evt.Sender = "srycantthnkof1"
	evt.Recipient = "SpookyGhostMachine"

	c.events <- evt
}
