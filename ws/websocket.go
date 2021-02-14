package ws

import (
	"errors"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Connection implements io.ReadWriteCloser with gorilla/websocket
type Connection struct {
	sync.Mutex
	conn *websocket.Conn

	// Reconnection params
	connURL             url.URL
	autoReconnect       bool
	maxReconnectRetries int
}

// NewWebSocket creates a WS container, without connecting. Must call Connect() when
// ready to connect to the network. By default: auto-reconnect logic is disabled.
func NewWebSocket() *Connection {
	return &Connection{
		autoReconnect:       false,
		maxReconnectRetries: 0,
	}
}

// EnableAutoReconnect turns on auto-reconnect logic on a read/write error
func (ws *Connection) EnableAutoReconnect() {
	ws.autoReconnect = true
}

// SetMaxReconnects sets the maximum number of reconnect attempts made on a
// read/write error
func (ws *Connection) SetMaxReconnects(maxReconnectRetries int) {
	ws.maxReconnectRetries = maxReconnectRetries
}

// Connect actually connects the underlying WS connection
func (ws *Connection) Connect(scheme, server string) error {
	ws.connURL = url.URL{Scheme: scheme, Host: server, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(ws.connURL.String(), nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	return nil
}

// Read from the underlying connection.
// If ws.autoReconnect is true, on a read error it will attempt to reconnect the WS.
// On ws.maxReconnectRetries, the error will be returned
func (ws *Connection) Read(dest []byte) (int, error) {
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		log.Printf("ERROR: read ws - %v", err)
		return 0, err
	}

	return copy(dest, message), nil
}

// Write to the underlying connection.
// If ws.autoReconnect is true, on a write error it will attempt to reconnect the WS.
// On ws.maxReconnectRetries, the error will be returned
func (ws *Connection) Write(data []byte) (int, error) {
	if ws.conn == nil {
		return 0, errors.New("ws.conn is nil. Did you forget to call ws.Connect()?")
	}

	// Lock since WriteMessage requires only one concurrent execution
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	if err := ws.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return 0, err
	}

	return len(data), nil
}

// Close the underlying connection
func (ws *Connection) Close() error {
	return ws.conn.Close()
}
