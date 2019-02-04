package rooms

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"hobee-be/models"
	"hobee-be/pkg/hconst"
	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
)

type Room struct {
	Users    []*models.WSUser
	Messages []models.Message
}

const (
	roomIDLength = 15
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	mutex       = &sync.Mutex{}
	rooms       = map[string]Room{}
)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// TODO: Some way to interact with the room from the WS - to close it if needed
func Init(userPool <-chan [2]*models.WSUser) {
// TODO: Can you provide actual context here?
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			select {
			case paired := <-userPool:
				id := getUniqueRoomID()

				mutex.Lock()
				rooms[id] = Room{Users: []*models.WSUser{paired[0], paired[1]}}
				mutex.Unlock()

				fmt.Printf("Got 2 users: %d + %d\n", paired[0].ID, paired[1].ID)

				for _, u := range paired {
					msg := models.WSMessage{hconst.SYSTEM_MESSAGE, hconst.SYS_C}
					j, err := json.Marshal(msg)
					if err != nil {
						log.Error(ctx, err)
						continue
					}

					err = u.Socket.WriteMessage(websocket.TextMessage, j)
					// TODO: Research binary more. Use it for system messages?
					//err := cl.Connection.WriteMessage(websocket.BinaryMessage, []byte(SYS_C))
					if err != nil {
						log.Error(ctx, err)
						continue
					}

					u.RoomID = id
					u.Paired <- true
				}
			}
		}
	}()
}

func Close(id string) {
	// TODO: Add actual context here?
	ctx := context.Background()

	r, ok := rooms[id]
	if !ok {
		return
	}

	for _, wsu := range r.Users {
		wsu.RoomID = ""

		msg := models.WSMessage{hconst.SYSTEM_MESSAGE, hconst.SYS_DC}
		j, err := json.Marshal(msg)
		if err != nil {
			log.Error(ctx, err)
			continue
		}

		err = wsu.Socket.WriteMessage(websocket.TextMessage, j)
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway) {
				// TODO: If client clicks "Disconnect" - np. If closes tab, log will have "websocket: close sent", cus
				// websocket will be closed here at that point
				log.Error(ctx, herrors.Wrap(err, "ctx", "Read error"))
			}
			continue
		}
	}

// TODO: Save messages to DB
	mutex.Lock()
	delete(rooms, id)
	mutex.Unlock()
}

func getUniqueRoomID() string {
// TODO: Add a counter here to prevent infinite loop if it generates common rooms? What is the chance of that?
	randID := randStringRunes(roomIDLength)

	if _, ok := rooms[randID]; !ok {
		return randID
	}

	return getUniqueRoomID()
}

func Broadcast(u *models.WSUser, msg string) error {
// TODO: pass actual context here?
// TODO: I don't actually get to this part?
	if _ , ok := rooms[u.RoomID]; !ok {
		return herrors.New("Unknown room", "room", u.RoomID)
	}

	paired := rooms[u.RoomID].Users

	for _, pu := range paired {
		out := models.WSMessage{Type: hconst.OWN_MESSAGE}
		if u.Socket != pu.Socket {
			out.Type = hconst.BUDDY_MESSAGE
		}

		out.Message = msg
		j, err := json.Marshal(out)
		if err != nil {
			return herrors.Wrap(err)
		}

		err = pu.Socket.WriteMessage(websocket.TextMessage, j)
		if err != nil {
			return herrors.Wrap(err)
		}
	}

	// Add a message to the chatlog as well
	//a := fmt.Sprintf("%p", author)
	//chatlog[r.Id] = append(chatlog[r.Id], LogMessage{Author: a, Message: msg})

	return nil
}
