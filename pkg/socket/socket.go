package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeSystem MessageType = "s"

	MessageTypeChatting MessageType = "c"
	MessageTypeOwn      MessageType = "o"
	MessageTypeBuddy    MessageType = "b"

	SystemSearch       = "s"
	SystemConnected    = "c"
	SystemDisconnected = "d"
	SystemTalking      = "t"

	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

// socket represents a uuid and a websocket connection
type Socket struct {
	conn *websocket.Conn
	// Send is used to send message to websockets
	Send chan []byte
	// Send is used when socket receives a message and wants to broadcast it to everyone in the room, ending
	// up in Send
	Broadcast chan Broadcast
}

type Message struct {
	AuthorUUID string
	Type       MessageType `json:"type"`
	Text       string      `json:"text"`
	Timestamp  time.Time   `json:"timestamp"`
}

func New(conn *websocket.Conn) *Socket {
	return &Socket{conn: conn, Send: make(chan []byte)}
}
