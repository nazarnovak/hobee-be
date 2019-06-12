package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
	"hobee-be/pkg/message"
)

type Room struct {
	ID        uuid.UUID
	Messages  []message.Message
	Broadcast chan Broadcast
	Sockets   [2]*Socket
}

type Broadcast struct {
	Socket *Socket
	Type   MessageType
	Text   []byte
}

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	mutex       = &sync.Mutex{}
	rooms       = map[uuid.UUID]*Room{}
)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Rooms(matchedSockets <-chan [2]*Socket) {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case sockets := <-matchedSockets:
				roomID, err := getUniqueRoomID()
				if err != nil {
					log.Critical(ctx, herrors.Wrap(err))
					return
				}

				mutex.Lock()

				bc := make(chan Broadcast)

				sockets[0].Broadcast, sockets[1].Broadcast = bc, bc
				// Should roomID be also added to sockets for reference when they close connection so you need to stop room from existing?
				room := &Room{
					ID:        roomID,
					Messages:  []message.Message{},
					Broadcast: bc,
					Sockets:   [2]*Socket{sockets[0], sockets[1]},
				}

				rooms[roomID] = room

				mutex.Unlock()

fmt.Printf("Got 2 sockets in room: %s\n", roomID)

				msg := Message{MessageTypeSystem, SystemConnected}
				o, err := json.Marshal(msg)
				if err != nil {
					log.Critical(ctx, err)
					continue
				}

				go room.Broadcaster()

				sockets[0].Send <- o
				sockets[1].Send <- o

			}
		}
	}()
}

func getUniqueRoomID() (uuid.UUID, error) {
	u := uuid.NewV4()

	if _, ok := rooms[u]; !ok {
		return u, nil
	}

	return getUniqueRoomID()
}

func (r *Room) Broadcaster() {
	ctx := context.Background()

	for {
		select {
		case b := <-r.Broadcast:
			for _, socket := range r.Sockets {
				switch {
				case b.Type == MessageTypeChatting:
					t := MessageTypeOwn
					if b.Socket != socket {
						t = MessageTypeBuddy
					}

					msg := Message{
						Type: t,
						Text: string(b.Text),
					}

					o, err := json.Marshal(msg)
					if err != nil {
						log.Critical(ctx, err)
						continue
					}

					socket.Send <- o
				case b.Type == MessageTypeSystem:
					// Maybe this will error twice? Since we're ranging through all the sockets in a room
					if string(b.Text) != SystemDisconnected {
						log.Critical(ctx, herrors.New("Unknown system message text", "text", string(b.Text)))
						continue
					}

					msg := Message{
						Type: MessageTypeSystem,
						Text: SystemDisconnected,
					}

					o, err := json.Marshal(msg)
					if err != nil {
						log.Critical(ctx, err)
						continue
					}

					socket.Send <- o
				default:
					log.Critical(ctx, herrors.New("Unknown type passed to the broadcaster"))
				}

			}
		}
	}
}

//func Close(id string) {
//	// TODO: Add actual context here?
//	ctx := context.Background()
//
//	r, ok := rooms[id]
//	if !ok {
//		return
//	}
//
//	for _, wsu := range r.Users {
//		wsu.RoomID = ""
//
//		msg := models.WSMessage{hconst.SYSTEM_MESSAGE, hconst.SYS_DC}
//		j, err := json.Marshal(msg)
//		if err != nil {
//			log.Error(ctx, err)
//			continue
//		}
//
//		err = wsu.Socket.WriteMessage(websocket.TextMessage, j)
//		if err != nil {
//			if !websocket.IsCloseError(err, websocket.CloseGoingAway) {
//				// TODO: If client clicks "Disconnect" - np. If closes tab, log will have "websocket: close sent", cus
//				// websocket will be closed here at that point
//				log.Error(ctx, herrors.Wrap(err, "ctx", "Read error"))
//			}
//			continue
//		}
//	}
//
//	// TODO: Save messages to DB
//	mutex.Lock()
//	delete(rooms, id)
//	mutex.Unlock()
//}
//
//func Broadcast(u *models.WSUser, msg string) error {
//	// TODO: pass actual context here?
//	// TODO: I don't actually get to this part?
//	if _, ok := rooms[u.RoomID]; !ok {
//		return herrors.New("Unknown room", "room", u.RoomID)
//	}
//
//	paired := rooms[u.RoomID].Users
//
//	for _, pu := range paired {
//		out := models.WSMessage{Type: hconst.OWN_MESSAGE}
//		if u.Socket != pu.Socket {
//			out.Type = hconst.BUDDY_MESSAGE
//		}
//
//		out.Message = msg
//		j, err := json.Marshal(out)
//		if err != nil {
//			return herrors.Wrap(err)
//		}
//
//		err = pu.Socket.WriteMessage(websocket.TextMessage, j)
//		if err != nil {
//			return herrors.Wrap(err)
//		}
//	}
//
//	// Add a message to the chatlog as well
//	//a := fmt.Sprintf("%p", author)
//	//chatlog[r.Id] = append(chatlog[r.Id], LogMessage{Author: a, Message: msg})
//
//	return nil
//}
//
//// This function returns a roomname if user is already in a chat
//func GetAlreadyChattingUserRoom(userID int64) string {
//	for roomName, room := range rooms {
//		for _, user := range room.Users {
//			if userID == user.ID {
//				return roomName
//			}
//		}
//	}
//
//	return ""
//}
//
//func JoinExistingRoom(userID int64) (bool, error) {
//	existingRoom := GetAlreadyChattingUserRoom(userID)
//	if existingRoom == "" {
//		return false, nil
//	}
//
//	room, ok := rooms[existingRoom]
//	if !ok {
//		return false, herrors.New("Tried to connect to a room that doesn't exist anymore", "existingroom", existingRoom, "userID", userID)
//	}
//
//	mutex.Lock()
//
//	mutex.Unlock()
//}
