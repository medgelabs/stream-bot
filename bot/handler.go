package bot

import (
	"medgebot/irc"
)

type Handler struct {
	consumer func(irc.Message)
	msgChan  chan irc.Message
}

func NewHandler(consumer func(irc.Message)) Handler {
	return Handler{
		consumer: consumer,
		msgChan:  make(chan irc.Message, 10),
	}
}

func (h Handler) Listen() {
	for msg := range h.msgChan {
		h.consumer(msg)
	}
}

func (h Handler) Receive(msg irc.Message) {
	h.msgChan <- msg
}
