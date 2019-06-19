package socket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

var (
	// Online mutex works with the usersSocketsMap (or user/sockets online map)
	onlineMutex     = &sync.Mutex{}
	usersSocketsMap = map[string]*User{}
)

type status int

const (
	// Disconnected could be either first time you log in, or when you disconnected from a room, then user will have
	// roomUUID, from which messages will be loaded
	statusDisconnected = 0
	statusSearching    = 1
	statusTalking      = 2
)

type User struct {
	UUID     string
	Sockets  []*Socket
	RoomUUID string

	Status status
}

// attachSocketToUser attaches one of the sockets to an existing user in the map (which is sort of like online), or
// creates a new user and attaches that to the online. It returns the user instance
func AttachSocketToUser(uuid string, s *Socket) *User {
	onlineMutex.Lock()
	if _, ok := usersSocketsMap[uuid]; !ok {
		u := &User{UUID: uuid, Sockets: []*Socket{}}
		usersSocketsMap[uuid] = u
	}

	// Add the socket to the newly created user, or to an existing user
	usersSocketsMap[uuid].Sockets = append(usersSocketsMap[uuid].Sockets, s)
	onlineMutex.Unlock()

	return usersSocketsMap[uuid]
}

func (u *User) Reader(ctx context.Context, s *Socket) {
	defer func() {
		u.Close(ctx, s)
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
			u.handleSystemMessage(ctx, s, msg.Text)
		case MessageTypeOwn:
			s.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeChatting, Text: []byte(msg.Text)}
		default:
			err := herrors.New("Unknown type received in the message", "msg", msg)
			log.Critical(ctx, err)
		}
	}
}

func (u *User) Writer(ctx context.Context, s *Socket) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		u.Close(ctx, s)
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
			s.conn.SetReadDeadline(time.Now().Add(writeWait))
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if err != websocket.ErrCloseSent {
					log.Error(ctx, herrors.Wrap(err))
				}
				return
			}
		}
	}
}

// If socket disconnects - we need to close the socket not to have a memory leak
func (u *User) Close(ctx context.Context, s *Socket) {
fmt.Println("Closing down disconnected websocket")
	onlineMutex.Lock()
	for k, socket := range u.Sockets {
		if socket.conn != s.conn {
			continue
		}

		u.Sockets = append(u.Sockets[:k], u.Sockets[k+1:]...)
	}

	onlineMutex.Unlock()

	// Close the actual websocket
	s.conn.Close()
}

func (u *User) handleSystemMessage(ctx context.Context, s *Socket, cmd string) {
	switch cmd {
	case SystemSearch:
		// Enter search mode for user
		searchAdd(u)
	case SystemDisconnected:
		// Disconnect from the current the conversation, but still part of a room until next search
// UpdateStatus(users[0].UUID, statusDisconnected)
// UpdateStatus(users[1].UUID, statusDisconnected)
// UpdateStatus(room[uuid], statusDisconnected)
		s.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeSystem, Text: []byte(SystemDisconnected)}
	default:
		err := herrors.New("Unknown command received on websocket conn", "cmd", cmd)
		log.Critical(ctx, err)
	}
}

func UpdateStatus(uuid string, status status) error {
	onlineMutex.Lock()
	if _, ok := usersSocketsMap[uuid]; !ok {
		return herrors.New("Could not find user with UUID in userSocketsMap", "uuid", uuid)
	}

	usersSocketsMap[uuid].Status = status
	onlineMutex.Unlock()

	return nil
}