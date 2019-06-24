package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type Room struct {
	ID        uuid.UUID
	Messages  []Message
	Broadcast chan Broadcast
	Users     [2]*User
	//Summaries [2]*Summary
}

// TODO: Doesn't have to be attached to a Room, could just have RoomID as a field instead?

//type reportReason int
//
//var (
//	reasonSpam reportReason = 0
//	reasonHarassing reportReason = 1
//	reasonRacism reportReason = 2
//	reasonSex reportReason = 3
//	reasonOther reportReason = 4
//
//	allReasons = []reportReason{reasonSpam, reasonHarassing ...}
//
//)
//
//type Summary struct {
//	AuthorUUID string
//	Liked bool
//	Reported reportReason
//}

type Broadcast struct {
	UUID string
	Type MessageType
	Text []byte
}

var (
	letterRunes       = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	matcherMutex      = &sync.Mutex{}
	roomMessagesMutex = &sync.Mutex{}
	rooms             = map[uuid.UUID]*Room{}
)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Rooms(matchedUsers <-chan [2]*User) {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case users := <-matchedUsers:
				roomID, err := getUniqueRoomID()
				if err != nil {
					log.Critical(ctx, herrors.Wrap(err))
					return
				}

				matcherMutex.Lock()

				// Broadcast is a shared channel between room and each socket, that way if you send something to either
				// room's broadcast or socket broadcast channel - it will be sent to everyone who's joined in that room
				bc := make(chan Broadcast)

				for _, u := range users {
					for _, s := range u.Sockets {
						s.Broadcast = bc
					}

					u.RoomUUID = roomID.String()
				}

				// Should roomID be also added to sockets for reference when they close connection so you need to stop room from existing?
				room := &Room{
					ID:        roomID,
					Messages:  []Message{},
					Broadcast: bc,
					Users:     [2]*User{users[0], users[1]},
				}

				rooms[roomID] = room

				UpdateStatus(users[0].UUID, statusTalking)
				UpdateStatus(users[1].UUID, statusTalking)

				matcherMutex.Unlock()

				fmt.Printf("Got 2 sockets in room: %s\n", roomID)

				go room.Broadcaster()

				room.Broadcast <- Broadcast{UUID: "", Type: MessageTypeSystem, Text: []byte(SystemConnected)}
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
			// TODO: Even tho the messages are already saved here, if there's an error happening it might cause
			// inconsistencies, so for example it will save an unknown type of a message, or someone might send a message,
			// and when trying to broadcast the message back to them there might be an error, so they'll try again and
			// that will duplicate the message
			roomMessagesMutex.Lock()
			msg := Message{
				Type: b.Type,
				Text: string(b.Text),
				AuthorUUID: b.UUID,
				Timestamp:  time.Now().UTC(),
			}

			r.Messages = append(r.Messages, msg)

			// "Wipe" the author after adding it to the room, so it doesn't get exposed to FE (not like it matters,
			// but yeah)
			msg.AuthorUUID = ""

			roomMessagesMutex.Unlock()

			for _, user := range r.Users {
				for _, socket := range user.Sockets {
					switch {
					case b.Type == MessageTypeChatting:
						t := MessageTypeOwn
						if b.UUID != user.UUID {
							t = MessageTypeBuddy
						}

						msg.Type = t

						o, err := json.Marshal(msg)
						if err != nil {
							log.Critical(ctx, err)
							continue
						}

						socket.Send <- o
					case b.Type == MessageTypeSystem:
						// Maybe this will error twice? Since we're ranging through all the sockets in a room
						if string(b.Text) != SystemConnected && string(b.Text) != SystemDisconnected {
							log.Critical(ctx, herrors.New("Unknown system message text", "text", string(b.Text)))
							continue
						}

						// If someone disconnected - we don't have to have broadcast channel alive anymore - we clean it
						// up
						if string(b.Text) == SystemDisconnected {
							r.Close()
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
}

// Close closes the rooms broadcast channel, since there is no need for that anymore.
func (r *Room) Close() {
	close(r.Broadcast)
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
//	matcherMutex.Lock()
//	delete(rooms, id)
//	matcherMutex.Unlock()
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
//}
