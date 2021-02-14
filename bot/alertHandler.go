package bot

import (
	"sync"

	"github.com/gorilla/websocket"
)

/*
  AlertHandler connects the Bot event broadcast to a (potentially) connected
  client listening over a WebSocket connection. Because we cannot guarantee that
  this connection will be available on start, we consider the WS to be "unsafe".

  We only ever _write_ to this WebSocket. We should never read from it.

  Note that this is, effectively, global mutable state and, as such, we should
  be as defensive about its use as that typically calls for.

  Assumption of control:
	  main.go - creates the reference and passes it around
		  ws := WriteOnlyUnsafeWebSocket{}
		  startServer(ws)
		  RegisterAlertHandler(ws)

	  server.go(alertWs WriteOnlyUnsafeWebSocket): - handles actual Connection references
		WS Connect -> alertsWs.Open(conn)
*/

// RegisterAlertHandler publishes Alert-related events to be displayed by
// on-screen consumers
func (bot *Bot) RegisterAlertHandler(ws *WriteOnlyUnsafeWebSocket) {
	bot.RegisterHandler(
		NewHandler(func(evt Event) {
			if !ws.Connected() {
				// TODO retry to account for disconnects?
				return
			}

			// Send JSON payload to consumer for rendering
		}),
	)
}

// WriteOnlyUnsafeWebSocket represents a write-only WebSocket that should be
// used defensively, as there is no guarantee the connection was initialized
// or is currently connected.
type WriteOnlyUnsafeWebSocket struct {
	sync.Mutex
	conn *websocket.Conn // Write-only connection
}

// Connected checks if we have a WebSocket connection at all
func (ws *WriteOnlyUnsafeWebSocket) Connected() bool {
	return ws.conn != nil
}

// SetConnection (re)assigns the reference to the open Alerts WebSocket via a mutex-locked
// assignment
func (ws *WriteOnlyUnsafeWebSocket) SetConnection(conn *websocket.Conn) {
	ws.Lock()
	defer ws.Unlock()

	ws.conn = conn
}
