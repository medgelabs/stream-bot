package ws

import (
	"fmt"
	log "medgebot/logger"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Connection implements io.ReadWriteCloser with gorilla/websocket,
// also allowing for automatic reconnect/retry logic on a read/write
// error
type Connection struct {
	sync.Mutex
	conn    *websocket.Conn
	connURL url.URL

	// Reconnection/retry params
	postReconnectFunc func() error
	maxRetries        int
	retryWaitDuration time.Duration
}

// NewWebSocket creates a WS container without connecting to the actual network.
// Must call Connect() when ready to connect to the network. This is to simplify reconnects.
// By default: postReconnectFunc is a no-op and should be set with SetPostReconnectFunc().
// By default: maxRetries is set to 0 and should be set with SetMaxRetries()
func NewWebSocket(scheme, server string) *Connection {
	return &Connection{
		connURL:           url.URL{Scheme: scheme, Host: server, Path: "/"},
		postReconnectFunc: func() error { return nil }, // no-op
		maxRetries:        0,
		retryWaitDuration: 50 * time.Millisecond,
	}
}

// SetPostReconnectFunc assigns the function that will be run after a successful
// reconnect. This accounts for client-specific actions like re-authentication, etc
func (ws *Connection) SetPostReconnectFunc(f func() error) {
	ws.postReconnectFunc = f
}

// SetMaxRetries sets the maximum number of retry attempts made on a
// read/write error
func (ws *Connection) SetMaxRetries(maxRetries int) {
	ws.maxRetries = maxRetries
}

// Connect actually connects the underlying WS connection
func (ws *Connection) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(ws.connURL.String(), nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	return nil
}

// Reconnect attempts Connect() and the given reconnectFunc(). If either fails,
// a wrapped error is returned describing which part failed
func (ws *Connection) Reconnect() error {
	err := ws.Connect()
	if err != nil {
		return errors.Wrap(err, "reconnect failed")
	}

	err = ws.postReconnectFunc()
	if err != nil {
		return errors.Wrap(err, "post-reconnect func failed")
	}

	return nil
}

// Read from the underlying connection.
// On a read error, if maxRetries > 0, it will attempt to reconnect the WS.
// On ws.maxReconnectRetries, the error will be returned
func (ws *Connection) Read(dest []byte) (int, error) {
	_, message, err := ws.conn.ReadMessage()
	if err == nil {
		return copy(dest, message), nil
	}

	// On error, retry
	log.Error("read ws", err)
	if ws.maxRetries > 0 {
		for i := 1; i <= ws.maxRetries; i++ {
			// Make sure we wait before retrying
			nextWait := time.Duration(i) * ws.retryWaitDuration
			time.Sleep(nextWait)

			reconnErr := ws.Reconnect()
			if reconnErr != nil {
				log.Error(fmt.Sprintf("retry %d failed. Retry after %d", i, nextWait), reconnErr)
				continue
			}

			_, message, err := ws.conn.ReadMessage()
			if err == nil {
				return copy(dest, message), nil
			}

			log.Error(fmt.Sprintf("retry %d failed. Retry after %d", i, nextWait), err)
		}
	}

	// Absolute failure, return the error
	return 0, err
}

// Write to the underlying connection.
// On a write error, if maxRetries > 0, it will attempt to reconnect the WS.
// On ws.maxReconnectRetries, the error will be returned
func (ws *Connection) Write(data []byte) (int, error) {
	// Lock since WriteMessage requires only one concurrent execution
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	err := ws.conn.WriteMessage(websocket.TextMessage, data)

	// On success, just return
	if err == nil {
		return len(data), nil
	}

	// Otherwise, if we hit an error and we attempt reconnects
	if ws.maxRetries > 0 {
		for i := 1; i <= ws.maxRetries; i++ {
			// Make sure we wait before retrying
			nextWait := time.Duration(i) * ws.retryWaitDuration
			time.Sleep(nextWait)

			reconnErr := ws.Reconnect()
			if reconnErr != nil {
				log.Error(fmt.Sprintf("retry %d failed. Retry after %d", i, nextWait), reconnErr)
				continue
			}

			err := ws.conn.WriteMessage(websocket.TextMessage, data)
			if err == nil {
				return len(data), nil
			}

			log.Error(fmt.Sprintf("retry %d failed. Retry after %d", i, nextWait), err)
		}
	}

	// Absolute failure, return the error
	return 0, err
}

// Close the underlying connection
func (ws *Connection) Close() error {
	return ws.conn.Close()
}
