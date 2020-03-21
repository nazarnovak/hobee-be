package socket

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"github.com/nazarnovak/hobee-be/pkg/db"
	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type ReportReason string

var (
	reasonSpam      ReportReason = "rsp"
	reasonHarassing ReportReason = "rha"
	reasonRacism    ReportReason = "rra"
	reasonSex       ReportReason = "rse"
	reasonOther     ReportReason = "rot"

	//allReasons = []ReportReason{reasonSpam, reasonHarassing ...}
)

type Room struct {
	ID        uuid.UUID
	Messages  []Message
	Broadcast chan Broadcast
	Users     [2]*User
	Active    bool
	Results   [2]*Result
}

type Result struct {
	AuthorUUID string       `json:"-"`
	Liked      bool         `json:"liked"`
	Reported   ReportReason `json:"reported"`
}

// TODO: Doesn't have to be attached to a Room, could just have RoomID as a field instead?

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

func (r *Room) SaveMessages(secret string) error {
	q := `INSERT INTO chats(user1, user2, room, chat, started, finished)
		VALUES($1, $2, $3, $4, $5, $6);`

	chatBytes, err := json.Marshal(r.Messages)
	if err != nil {
		return herrors.Wrap(err)
	}

	encryptedMsgs, err := EncryptMessages(chatBytes, secret)
	if err != nil {
		return herrors.Wrap(err)
	}

	if _, err := db.Instance.Exec(q, r.Users[0].UUID, r.Users[1].UUID, r.ID, encryptedMsgs, r.Messages[0].Timestamp, time.Now().UTC()); err != nil {
		return herrors.Wrap(err)
	}

	return nil
	//filename := fmt.Sprintf("%s.%s", r.ID.String(), "csv")
	//
	//if exists := FileExists(fmt.Sprintf("%s/%s", "chats", filename)); exists {
	//	return herrors.New("Filename already exists", "roomuuid", r.ID.String())
	//}
	//
	//file, err := os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_CREATE|os.O_WRONLY, 0777)
	//if err != nil {
	//	return herrors.Wrap(err)
	//}
	//defer file.Close()
	//
	//rows := make([][]string, 0, len(r.Messages)+1)
	//
	//// Headers
	//rows = append(rows, []string{"timestamp", "authoruuid", "type", "text"})
	//
	//for _, msg := range r.Messages {
	//	row := []string{msg.Timestamp.Format(time.RFC3339), msg.AuthorUUID, string(msg.Type), msg.Text}
	//
	//	rows = append(rows, row)
	//}
	//
	//wr := csv.NewWriter(file)
	//wr.Comma = ';'
	//
	//if err := wr.WriteAll(rows); err != nil {
	//	return herrors.Wrap(err)
	//}
	//wr.Flush()

	return nil
}

func EncryptMessages(messages []byte, secret string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return nil, herrors.Wrap(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, herrors.Wrap(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		return nil, herrors.Wrap(err)
	}

	return gcm.Seal(nonce, nonce, messages, nil), nil
}

func saveResultLikeCSV(roomuuid, useruuid string, liked bool) error {
	filename := fmt.Sprintf("%s:%s.%s", roomuuid, useruuid, "csv")

	exists := FileExists(fmt.Sprintf("%s/%s", "chats", filename))
	if !exists {
		file, err := os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return herrors.Wrap(err)
		}
		defer file.Close()

		// Always going to be 2 rows: headers + like/report
		rows := make([][]string, 0, 2)

		// Headers
		rows = append(rows, []string{"liked", "reported"})

		// The default value for reported is not reported - empty string
		reported := ""

		row := []string{strconv.FormatBool(liked), reported}
		rows = append(rows, row)

		wr := csv.NewWriter(file)
		wr.Comma = ';'

		if err := wr.WriteAll(rows); err != nil {
			return herrors.Wrap(err)
		}
		wr.Flush()

		return nil
	}

	// If file already exists - we'll only edit the field we need to change
	// First we only read the values from the file
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_RDONLY, 0777)
	if err != nil {
		return herrors.Wrap(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		return herrors.Wrap(err)
	}

	if len(rows) != 2 {
		return herrors.New("Expecting 2 rows in the csv")
	}

	// Liked will be on the second row, first column
	rows[1][0] = strconv.FormatBool(liked)

	// Now we can truncate the file and write the new values
	file, err = os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return herrors.Wrap(err)
	}

	wr := csv.NewWriter(file)
	wr.Comma = ';'

	if err := wr.WriteAll(rows); err != nil {
		return herrors.Wrap(err)
	}
	wr.Flush()

	return nil
}

func saveResultReportedCSV(roomuuid, useruuid string, reported ReportReason) error {
	filename := fmt.Sprintf("%s:%s.%s", roomuuid, useruuid, "csv")

	exists := FileExists(fmt.Sprintf("%s/%s", "chats", filename))
	if !exists {
		file, err := os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return herrors.Wrap(err)
		}
		defer file.Close()

		// Always going to be 2 rows: headers + like/report
		rows := make([][]string, 0, 2)

		// Headers
		rows = append(rows, []string{"liked", "reported"})

		// The default value for liked is not liked - false
		liked := false

		row := []string{strconv.FormatBool(liked), string(reported)}
		rows = append(rows, row)

		wr := csv.NewWriter(file)
		wr.Comma = ';'

		if err := wr.WriteAll(rows); err != nil {
			return herrors.Wrap(err)
		}
		wr.Flush()

		return nil
	}

	// If file already exists - we'll only edit the field we need to change
	// First we only read the values from the file
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_RDONLY, 0777)
	if err != nil {
		return herrors.Wrap(err)
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		return herrors.Wrap(err)
	}

	if len(rows) != 2 {
		return herrors.New("Expecting 2 rows in the csv")
	}

	// Reported will be on the second row, second column
	rows[1][1] = string(reported)

	// Now we can truncate the file and write the new values
	file, err = os.OpenFile(fmt.Sprintf("%s/%s", "chats", filename), os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return herrors.Wrap(err)
	}

	wr := csv.NewWriter(file)
	wr.Comma = ';'

	if err := wr.WriteAll(rows); err != nil {
		return herrors.Wrap(err)
	}
	wr.Flush()

	return nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

//func randStringRunes(n int) string {
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letterRunes[rand.Intn(len(letterRunes))]
//	}
//	return string(b)
//}

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
					Results: [2]*Result{
						{AuthorUUID: users[0].UUID},
						{AuthorUUID: users[1].UUID},
					},
				}

				rooms[roomID.String()] = room

				UpdateStatus(users[0].UUID, statusTalking)
				UpdateStatus(users[1].UUID, statusTalking)

				matcherMutex.Unlock()

				// TODO: Do not match with the last peerson the user had a conversation with
				// We add the user to the other users history
				//users[0].UserHistory, users[1].UserHistory = []string{users[1].UUID}, []string{users[0].UUID}

				addRoomToUserRoomHistory(users[0], roomID.String())
				addRoomToUserRoomHistory(users[1], roomID.String())

				fmt.Printf("Got 2 sockets in room: %s\n", roomID)

				go room.Broadcaster()

				room.Broadcast <- Broadcast{UUID: "", Type: MessageTypeSystem, Text: []byte(SystemConnected)}

				// If someone went into search mode and closed the tab - inform the other user that the buddy is offline
				if len(users[0].Sockets) == 0 {
					room.Broadcast <- Broadcast{UUID: users[0].UUID, Type: MessageTypeActivity, Text: []byte(ActivityUserInactive)}
				}

				if len(users[1].Sockets) == 0 {
					room.Broadcast <- Broadcast{UUID: users[1].UUID, Type: MessageTypeActivity, Text: []byte(ActivityUserInactive)}
				}
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
			if !(b.Type == MessageTypeActivity && string(b.Text) == ActivityOwnTyping) {
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
	defer matcherMutex.Unlock()

	room, ok := rooms[uuid]
	if !ok {
		return nil, herrors.New("Failed to find a room", "uuid", uuid)
	}

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

func SetRoomLike(roomuuid, useruuid string, liked bool) error {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[roomuuid]
	if !ok {
		return herrors.New("Failed to find a room", "roomuuid", roomuuid)
	}

	for k, r := range room.Results {
		if r == nil {
			continue
		}

		if r.AuthorUUID != useruuid {
			continue
		}

		room.Results[k].Liked = liked
	}

	// Set like in DB
	user := "user1"
	if useruuid == room.Users[1].UUID {
		user = "user2"
	}

	q := fmt.Sprintf(`UPDATE chats SET %[1]s_liked = $1 WHERE room = $2;`, user)
	if _, err := db.Instance.Exec(q, liked, roomuuid); err != nil {
		return herrors.Wrap(err)
	}

	return nil
}

func SetRoomReport(roomuuid, useruuid string, reason ReportReason) error {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[roomuuid]
	if !ok {
		return herrors.New("Failed to find a room", "roomuuid", roomuuid)
	}

	for k, r := range room.Results {
		if r == nil {
			continue
		}

		if r.AuthorUUID != useruuid {
			continue
		}

		room.Results[k].Reported = reason
	}

	// Set report in DB
	user := "user1"
	if useruuid == room.Users[1].UUID {
		user = "user2"
	}

	q := fmt.Sprintf(`UPDATE chats SET %[1]s_reported = $1 WHERE room = $2;`, user)
	if _, err := db.Instance.Exec(q, reason, roomuuid); err != nil {
		return herrors.Wrap(err)
	}
	return nil
}

func addRoomToUserRoomHistory(user *User, roomuuid string) {
	user.RoomHistory = append(user.RoomHistory, roomuuid)

	// We only want to start with the latest 3 conversations you had, so truncate other rooms
	if len(user.RoomHistory) > 3 {
		// Remove the oldest room from history
		user.RoomHistory = append(user.RoomHistory[:0], user.RoomHistory[1:]...)
	}

	return
}

func GetRoomUserResults(roomuuid, useruuid string) (*Result, error) {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	room, ok := rooms[roomuuid]
	if !ok {
		return nil, herrors.New("Failed to find a room", "roomuuid", roomuuid)
	}

	for _, result := range room.Results {
		if result.AuthorUUID != useruuid {
			continue
		}

		return result, nil
	}

	return nil, herrors.New("Could not find user results in a room", "roomuuid", roomuuid,
		"useruuid", useruuid)
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
