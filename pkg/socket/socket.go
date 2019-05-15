package socket

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
)

type MessageType string

const (
	MessageTypeSystem MessageType = "s"
	MessageTypeOwn    MessageType = "o"
	MessageTypeBuddy  MessageType = "b"

	SystemSearch       = "s"
	SystemConnected    = "c"
	SystemDisconnected = "d"

	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//var userSocketsMap = map[int64][]*socket{}

// socket represents a uuid and a websocket connection
type Socket struct {
	UUID uuid.UUID
	conn *websocket.Conn
	// Send is used to send message to websockets
	Send chan []byte
	// Send is used when socket receives a message and wants to broadcast it to everyone in the room, ending
	// up in Send
	Broadcast chan Broadcast
}

type Message struct {
	Type MessageType `json:"type"`
	Text string      `json:"text"`
}

func New(uuid uuid.UUID, conn *websocket.Conn) *Socket {
	return &Socket{UUID: uuid, conn: conn, Send: make(chan []byte)}
}

func (s *Socket) Reader(ctx context.Context) {
	defer func() {
		s.Close(ctx)
	}()

	s.conn.SetReadLimit(maxMessageSize)

	// Set the deadline for the first message we expect to receive
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg Message
		if err := s.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Critical(ctx, herrors.Wrap(err))
			}
			break
		}

		// Extend deadline
		s.conn.SetReadDeadline(time.Now().Add(pongWait))
fmt.Printf("%+v\n", msg)
		switch msg.Type {
		case MessageTypeSystem:
			s.handleSystemMessage(ctx, msg.Text)
		case MessageTypeOwn:
fmt.Printf("Received own message: %+v\n", msg.Text)
			s.Broadcast <- Broadcast{Socket: s, Text: []byte(msg.Text)}
		default:
			err := herrors.New("Unknown type received in the message", "msg", msg)
			log.Critical(ctx, err)
		}
	}
}

func (s *Socket) Writer(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.conn.Close()
	}()
	for {
		select {
		case message, ok := <-s.Send:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Critical(ctx, herrors.New("Room no longer exists"))
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				return
			}
			w.Write(message)

			//// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Clean up if socket closes
func (s *Socket) Close(ctx context.Context) {
	s.conn.Close()
}

func (s *Socket) ReceiveMessage(msg Message) error {

	return nil
}

func (s *Socket) handleSystemMessage(ctx context.Context, cmd string) {
	switch cmd {
	case SystemSearch:
		// Enter search mode for user
		searchAdd(s)
	case SystemDisconnected:
		// Disconnect from the current room and end the conversation
	default:
		err := herrors.New("Unknown command received on websocket conn", "cmd", cmd)
		log.Critical(ctx, err)
	}
}
