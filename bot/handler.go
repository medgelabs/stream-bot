package bot

type Handler struct {
	consumer func(Event)
	msgChan  chan Event
}

func NewHandler(consumer func(Event)) Handler {
	return Handler{
		consumer: consumer,
		msgChan:  make(chan Event, 10),
	}
}

func (h Handler) Listen() {
	for msg := range h.msgChan {
		h.consumer(msg)
	}
}

func (h Handler) Receive(msg Event) {
	h.msgChan <- msg
}
