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
	UUID    string
	Sockets []*Socket

	// Broadcast is used when socket receives a message and wants to broadcast it to everyone in the room, ending
	// up in Send
	Broadcast chan<- Broadcast

	RoomUUID string

	Status status

	// The latest users uuids that this user had a conversation with
	//UserHistory []string

	// The latest room uuids that this user was a part of
	RoomHistory []string
}

// attachSocketToUser attaches one of the sockets to an existing user in the map (which is sort of like online), or
// creates a new user and attaches that to the online. It returns the user instance
func AttachSocketToUser(uuid string, s *Socket) *User {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()

	if _, ok := usersSocketsMap[uuid]; !ok {
		u := &User{UUID: uuid, Sockets: []*Socket{}}
		usersSocketsMap[uuid] = u
	}

	u := usersSocketsMap[uuid]

	// If user was in a room earlier and reconnected - notify the room that the user became active again
	if len(u.Sockets) == 0 && u.RoomUUID != "" {
		active, err := IsRoomActive(u.RoomUUID)
		if err != nil {
			log.Critical(context.Background(), herrors.Wrap(err))
			return nil
		}
		if active {
			u.Broadcast <- Broadcast{UUID: uuid, Type: MessageTypeActivity, Text: []byte(ActivityUserActive)}
		}
	}

	// Add the socket to the newly created user, or to an existing user
	usersSocketsMap[uuid].Sockets = append(usersSocketsMap[uuid].Sockets, s)

	return usersSocketsMap[uuid]
}

func UserInARoomUUID(userUUID string) string {
	// TODO: Not sure if you need locks for read only?
	onlineMutex.Lock()
	defer onlineMutex.Unlock()

	roomUUID := ""

	if u, ok := usersSocketsMap[userUUID]; ok {
		if u.RoomUUID != "" {
			roomUUID = u.RoomUUID
		}
	}

	return roomUUID
}

func (u *User) Reader(ctx context.Context, s *Socket, secret string) {
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
			u.handleSystemMessage(ctx, s, msg.Text, secret)
		case MessageTypeOwn:
			u.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeChatting, Text: []byte(msg.Text)}
		case MessageTypeActivity:
			if msg.Text != ActivityOwnTyping {
				log.Critical(ctx, herrors.New("Unexpected activity message", "msg", msg))
				continue
			}

			u.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeActivity, Text: []byte(msg.Text)}
		case MessageTypeResult:
			// Likes
			if msg.Text == ResultLike || msg.Text == ResultDislike {
				liked := true

				if msg.Text == ResultDislike {
					liked = false
				}

				if err := SetRoomLike(u.RoomUUID, u.UUID, liked); err != nil {
					log.Critical(ctx, herrors.Wrap(err, "user", u, "msg", msg))
					continue
				}

				//if err := saveResultLikeCSV(u.RoomUUID, u.UUID, liked); err != nil {
				//	log.Critical(ctx, herrors.Wrap(err, "user", u, "msg", msg))
				//	continue
				//}
				continue
			}

			if isReportOption(msg.Text) {
				if err := SetRoomReport(u.RoomUUID, u.UUID, ReportReason(msg.Text)); err != nil {
					log.Critical(ctx, herrors.New("Couldn't set a like on a room", "user", u, "msg", msg))
					continue
				}

				//if err := saveResultReportedCSV(u.RoomUUID, u.UUID, ReportReason(msg.Text)); err != nil {
				//	log.Critical(ctx, herrors.Wrap(err, "user", u, "msg", msg))
				//	continue
				//}
				continue
			}

			log.Critical(ctx, herrors.New("Unexpected result message", "msg", msg))
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
		// No need to close the websocket connection because it's already done by the reader?
		//u.Close(ctx, s)
	}()
	for {
		select {
		case message, ok := <-s.Send:
			if !ok {
				log.Critical(ctx, herrors.New("Send channel is closed", "useruuid", u.UUID))
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := s.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				return
			}

			if _, err := w.Write(message); err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				return
			}

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
	defer onlineMutex.Unlock()

	for k, socket := range u.Sockets {
		if socket.conn != s.conn {
			continue
		}

		u.Sockets = append(u.Sockets[:k], u.Sockets[k+1:]...)
	}

	// If this is the last socket of the user - set a user inactive event in the room
	if len(u.Sockets) == 0 {
		// TODO: This removes user when they go into search and then close all tabs. Maybe worth leaving for now
		//u.Status = statusDisconnected
		//searchRemove(u.UUID)
		if u.RoomUUID != "" {
			u.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeActivity, Text: []byte(ActivityUserInactive)}
		}
	}

	// Close the actual websocket
	s.conn.Close()
}

func (u *User) handleSystemMessage(ctx context.Context, s *Socket, cmd, secret string) {
	switch cmd {
	case SystemSearch:
		// Enter search mode for user
		if err := cleanUpRoom(u); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
		}

		searchAdd(u)
	case SystemDisconnected:
		// Disconnect from the current the conversation, but still part of a room until next search
		// UpdateStatus(users[0].UUID, statusDisconnected)
		// UpdateStatus(users[1].UUID, statusDisconnected)
		// UpdateStatus(room[uuid], statusDisconnected)
		u.Broadcast <- Broadcast{UUID: u.UUID, Type: MessageTypeSystem, Text: []byte(SystemDisconnected)}

		if u.RoomUUID == "" {
			log.Critical(ctx, herrors.New("User tried to disconnect without having a room set", "useruuid", u.UUID))
			return
		}

		matcherMutex.Lock()
		room, ok := rooms[u.RoomUUID]
		if !ok {
			matcherMutex.Unlock()
			log.Critical(ctx, herrors.New("Failed to find a room", "roomuuid", u.RoomUUID))
			return
		}

		room.Active = false
		matcherMutex.Unlock()

		// Save the chat into a file
		if err := room.SaveMessages(secret); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			return
		}
	default:
		err := herrors.New("Unknown command received on websocket conn", "cmd", cmd)
		log.Critical(ctx, err)
	}
}

func UpdateStatus(uuid string, status status) error {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()

	if _, ok := usersSocketsMap[uuid]; !ok {
		return herrors.New("Could not find user with UUID in userSocketsMap", "uuid", uuid)
	}

	usersSocketsMap[uuid].Status = status

	return nil
}

func UserRoomHistory(uuid string) ([]string, error) {
	onlineMutex.Lock()
	defer onlineMutex.Unlock()

	if _, ok := usersSocketsMap[uuid]; !ok {
		return nil, herrors.New("Could not find user with UUID in userSocketsMap", "uuid", uuid)
	}

	return usersSocketsMap[uuid].RoomHistory, nil
}

func cleanUpRoom(u *User) error {
	// User was previously in a conversation. If they are the last to leave - we can close the old rooms broadcast
	// channel already and remove the room from the room map
	if u.RoomUUID == "" {
		return nil
	}

	if err := RoomRemoveUser(u.RoomUUID, u.UUID); err != nil {
		return herrors.Wrap(err)
	}

	allDced, err := IsAllRoomUsersDisconnected(u.RoomUUID)
	if err != nil {
		return herrors.Wrap(err)
	}

	// If there are still users in that room - we don't want to close it yet
	if !allDced {
		return nil
	}

	if err := CloseRoom(u.RoomUUID); err != nil {
		return herrors.Wrap(err)
	}

	return nil
}
