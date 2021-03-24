package server

import (
	"medgebot/bot"
	"net/http"
)

// DebugClient is a dummy client for the Bot that we use to send messages
type DebugClient struct {
	events chan<- bot.Event
}

// SetDestination from bot.Client
func (c *DebugClient) SetDestination(events chan<- bot.Event) {
	c.events = events
}

// SendSub sends a mock Subscription event to the Bot for testing
func (c *DebugClient) SendSub() {
	evt := bot.NewSubEvent()
	evt.Sender = "saltymoth"
	evt.Amount = 5
	evt.Message = "I am a willing test subject"

	c.events <- evt
}

func (s *Server) debugSub(client *DebugClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendSub()
	}
}

// SendBit sends a mock Bit event to the Bot for testing
func (c *DebugClient) SendBit() {
	evt := bot.NewBitsEvent()
	evt.Sender = "Przemko9856"
	evt.Amount = 100
	evt.Message = "I am a willing test subject"

	c.events <- evt
}

func (s *Server) debugBit(client *DebugClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client.SendBit()
	}
}
