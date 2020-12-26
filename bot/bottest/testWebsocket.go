package bottest

import (
	"sync"
)

type Websocket struct {
	sync.Mutex
	lines []string
}

func NewTestWebsocket() *Websocket {
	return &Websocket{
		lines: make([]string, 5),
	}
}

// Send is a convenience method over Write()
func (w *Websocket) Send(message string) {
	w.Write([]byte(message))
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

// io.ReadWriteCloser
func (w *Websocket) Read(dst []byte) (int, error) {
	w.Lock()
	defer w.Unlock()

	if len(w.lines) == 0 {
		// EOF?
		return 0, nil
	}

	head := w.lines[0]
	copy(dst, []byte(head))
	w.lines = w.lines[1:]
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
