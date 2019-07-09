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
	Active    bool
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
	rooms             = map[string]*Room{}
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
					u.Broadcast = bc
					u.RoomUUID = roomID.String()
				}

				room := &Room{
					ID:        roomID,
					Messages:  []Message{},
					Broadcast: bc,
					Users:     [2]*User{users[0], users[1]},
					Active:    true,
				}

				rooms[roomID.String()] = room

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

	if _, ok := rooms[u.String()]; !ok {
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

			// We received a signal that we will close this room since it's not used anymore - stop the looping now
			if b.Type == MessageTypeSystem && string(b.Text) == SystemCloseRoom {
				break
			}

			// Removing messageTypeSystem will break the flow (investigate), need to add like/report type here too for
			// results
			if !r.Active && b.Type != MessageTypeSystem {
				continue
			}

			roomMessagesMutex.Lock()
			msg := Message{
				Type:       b.Type,
				Text:       string(b.Text),
				AuthorUUID: b.UUID,
				Timestamp:  time.Now().UTC(),
			}

			// Do not add "typing" event to the room, Spammy McSpammerson
			if b.Type != MessageTypeActivity && string(b.Text) != ActivityOwnTyping {
				r.Messages = append(r.Messages, msg)
			}

			// "Wipe" the author after adding it to the room, so it doesn't get exposed to FE (not like it matters,
			// but yeah)
			msg.AuthorUUID = ""

			roomMessagesMutex.Unlock()

			for _, user := range r.Users {
				// User might be nil when they already went searching for a new conversation
				if user == nil {
					continue
				}

				for _, socket := range user.Sockets {
					switch {
					case b.Type == MessageTypeChatting:
						t := MessageTypeOwn
						if b.UUID != user.UUID {
							t = MessageTypeBuddy
						}

						msg.AuthorUUID = string(t)

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

						msg.AuthorUUID = string(MessageTypeSystem)

						if string(b.Text) == SystemDisconnected {
							msg.AuthorUUID = string(MessageTypeOwn)
							if b.UUID != user.UUID {
								msg.AuthorUUID = string(MessageTypeBuddy)
							}
						}

						o, err := json.Marshal(msg)
						if err != nil {
							log.Critical(ctx, err)
							continue
						}

						socket.Send <- o
					case b.Type == MessageTypeActivity:
						t := MessageTypeOwn
						if b.UUID != user.UUID {
							t = MessageTypeBuddy
						}

						// Don't send "typing" event to the user, who emitted it
						if string(b.Text) == ActivityOwnTyping && t == MessageTypeOwn {
							continue
						}

						msg.AuthorUUID = string(t)

						o, err := json.Marshal(msg)
						if err != nil {
							log.Critical(ctx, err)
							continue
						}

						socket.Send <- o
					default:
						log.Critical(ctx, herrors.New("Unknown type passed to the broadcaster", "type", b.Type))
					}
				}
			}
		}
	}
}

// CloseRoom closes the rooms broadcast channel, since there is no need for that anymore.
func CloseRoom(uuid string) error {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[uuid]
	if !ok {
		return herrors.New("Failed to find a room", "uuid", uuid)
	}

	room.Broadcast <- Broadcast{UUID: "", Type: MessageTypeSystem, Text: []byte(SystemCloseRoom)}

	close(room.Broadcast)
	delete(rooms, uuid)

	return nil
}

func RoomMessages(uuid string) ([]Message, error) {
	matcherMutex.Lock()

	room, ok := rooms[uuid]
	if !ok {
		return nil, herrors.New("Failed to find a room", "uuid", uuid)
	}

	matcherMutex.Unlock()

	// We need to make a separate copy of the messages, since if we don't want to change anything in the original struct
	// to meddle with original data
	msgsCopy := []Message{}

	for _, msg := range room.Messages {
		msgCopy := Message{
			AuthorUUID: msg.AuthorUUID,
			Type:       msg.Type,
			Text:       msg.Text,
			Timestamp:  msg.Timestamp,
		}

		msgsCopy = append(msgsCopy, msgCopy)
	}

	return msgsCopy, nil
}

func IsRoomActive(uuid string) (bool, error) {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[uuid]
	if !ok {
		return false, herrors.New("Failed to find a room, maybe it wasn't cleaned up but user roomUUID was changed",
			"uuid", uuid)
	}

	return room.Active, nil
}

func GetRoomBroadcastChannel(uuid string) (chan<- Broadcast, error) {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[uuid]
	if !ok {
		return nil, herrors.New("Failed to find a room", "roomuuid", uuid)
	}

	return room.Broadcast, nil
}

func RoomRemoveUser(roomuuid, useruuid string) error {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[roomuuid]
	if !ok {
		return herrors.New("Failed to find a room", "roomuuid", roomuuid)
	}

	for k, u := range room.Users {
		if u == nil {
			continue
		}

		if u.UUID != useruuid {
			continue
		}

		room.Users[k] = nil
	}

	return nil
}

func IsAllRoomUsersDisconnected(uuid string) (bool, error) {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[uuid]
	if !ok {
		return false, herrors.New("Failed to find a room", "roomuuid", uuid)
	}

	disconnected := true

	for _, u := range room.Users {
		if u != nil {
			disconnected = false
			break
		}
	}

	return disconnected, nil
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
