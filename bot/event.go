package bot

const (
	CHAT_MSG = iota
	BITS
	SUB
	GIFTSUB
	POINT_REDEMPTION
)

// All-encompassing model for Events that the Bot understands, such as:
// Chat messages
// Bit donations
// Subscriptions / Re-subscriptions / Gifted Subscriptions
// Point redemptions
type Event struct {
	Type      int    // Identify what kind of Event we are receiving
	Sender    string // Source user, empty if not tied to a user
	Recipient string // Target user, if applicable (i.e gifted subscription)
	Title     string // title of the Channel Point redemption made
	Message   string // User-supplied message, empty if not provided
	Amount    int    // Any numerical amount tied to the message (bits, points, sub count)
}

func NewChatEvent() Event {
	return Event{
		Type: CHAT_MSG,
	}
}

func (evt Event) IsChatEvent() bool {
	return evt.Type == CHAT_MSG
}

func NewBitsEvent() Event {
	return Event{
		Type: BITS,
	}
}

func (evt Event) IsBitsEvent() bool {
	return evt.Type == BITS
}

func NewSubEvent() Event {
	return Event{
		Type: SUB,
	}
}

func (evt Event) IsSubEvent() bool {
	return evt.Type == SUB
}

func NewGiftSubEvent() Event {
	return Event{
		Type: GIFTSUB,
	}
}

func (evt Event) IsGiftSubEvent() bool {
	return evt.Type == GIFTSUB
}

func NewPointsEvent() Event {
	return Event{
		Type: POINT_REDEMPTION,
	}
}

func (evt Event) IsPointsEvent() bool {
	return evt.Type == POINT_REDEMPTION
}
