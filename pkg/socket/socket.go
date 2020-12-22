package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeActivity MessageType = "a"
	MessageTypeChatting MessageType = "c"	
	MessageTypeSystem   MessageType = "s"
	MessageTypeResultLike   MessageType = "rl"
	MessageTypeResultReport   MessageType = "rr"

	ActivityUserActive   = "ua"
	ActivityUserInactive = "ui"
	ActivityOwnTyping    = "t"
	ActivityRoomActive   = "ra"
	ActivityRoomInactive = "ri"

	MessageTypeOwn   MessageType = "o"
	MessageTypeBuddy MessageType = "b"

	SystemSearch       = "s"
	SystemConnected    = "c"
	SystemDisconnected = "d"
	SystemCloseRoom    = "cr"

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

	ip        string
	userAgent string
}

type Message struct {
	//UUID       string      `json:"uuid"`
	AuthorUUID string      `json:"authoruuid"`
	Type       MessageType `json:"type"`
	Text       string      `json:"text"`
	Timestamp  time.Time   `json:"timestamp"`
}

func New(conn *websocket.Conn, userAgent string) *Socket {
	return &Socket{
		conn:      conn,
		Send:      make(chan []byte),
		ip:        conn.RemoteAddr().String(),
		userAgent: userAgent,
	}
}

func isValidLikeReasons(reasons []LikeReason) bool {
	for _, reason := range reasons {
		found := false

		for _, likeReason := range allLikeReasons {
			if reason == likeReason {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func isValidReportReasons(reasons []ReportReason) bool {
	for _, reason := range reasons {
		found := false

		for _, reportReason := range allReportReasons {
			if reason == reportReason {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}