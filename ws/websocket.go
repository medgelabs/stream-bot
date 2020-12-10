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
}

func NewWebsocket() *Connection {
	return &Connection{}
}

func (ws *Connection) Connect(scheme, server string) error {
	u := url.URL{Scheme: scheme, Host: server, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	ws.conn = conn
	return nil
}

func (ws *Connection) Read(dest []byte) (int, error) {
	_, message, err := ws.conn.ReadMessage()
	if err != nil {
		log.Printf("ERROR: read ws - %v", err)
		return 0, err
	}

	return copy(dest, message), nil
}

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

func (ws *Connection) Close() error {
	return ws.conn.Close()
}
