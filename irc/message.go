package irc

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

// AddTag populates the given key/value pair
func (msg *Message) AddTag(tag, value string) {
	if msg.Tags == nil {
		msg.Tags = make(map[string]string)
	}

	msg.Tags[tag] = value
}
