package wstest

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Websocket stubs a ws connection and keeps track of messages sent/received
// with convenience methods for Tests. Does NOT communicate on the network
type Websocket struct {
	sync.Mutex
	lines      []string
	readCursor int
}

// NewWebsocket returns a new, fresh Websocket stub
func NewWebsocket() *Websocket {
	return &Websocket{
		lines:      make([]string, 0),
		readCursor: 0,
	}
}

// Send is a convenience method over Write()
func (w *Websocket) Send(message string) {
	w.Write([]byte(message))
}

// SendAndWait is a convenience method over Write() that also waits
// for a new message to arrive
func (w *Websocket) SendAndWait(message string) {
	w.Write([]byte(message))

	w.Lock()
	current := len(w.lines)
	w.Unlock()

	for {
		if len(w.lines) != current {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// Received indicates if the given message contents was _ever_ received
// on this Websocket
func (w *Websocket) Received(contents string) bool {
	for _, line := range w.lines {
		if line == contents {
			return true
		}
	}

	// Not found
	return false
}

// String returns the current WS line buffer as a \n delimited string
func (w *Websocket) String() string {
	var sb strings.Builder
	for _, line := range w.lines {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Read Cursor at: %d", w.readCursor))

	return sb.String()
}

// io.ReadWriteCloser
func (w *Websocket) Read(dst []byte) (int, error) {
	// Block until lines available
	for {
		if len(w.lines) != 0 && w.readCursor < len(w.lines) {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	w.Lock()
	defer w.Unlock()

	head := w.lines[w.readCursor]
	copy(dst, []byte(head))
	w.readCursor++

	return len(head), nil
}

func (w *Websocket) Write(data []byte) (int, error) {
	w.Lock()
	defer w.Unlock()

	w.lines = append(w.lines, string(data))
	return len(data), nil
}

func (w *Websocket) Close() error {
	return nil
}
