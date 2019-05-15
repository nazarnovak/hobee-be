package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	"hobee-be/pkg/herrors2"
	"hobee-be/pkg/log"
	"hobee-be/pkg/socket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func GOT(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loggedIn, err := isLoggedIn(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if !loggedIn {
			log.Critical(ctx, herrors.New("Attempting to access websockets without being logged in"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}


		uuid, err := uuid.NewV4()
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}
fmt.Println("New socket connected at:", time.Now().UTC().String())
		s := socket.New(uuid, c)
// Creates a new socket with pkg/socket
// Run reader, which when received "search" command should somehow add the socket to search. That means
// we will import pkg/matcher
// pkg/matcher has references to pkg/socket.Socket as a part of Add/Remove functions
		go s.Reader(ctx)
		s.Writer(ctx)
responseJSONSuccess(ctx, w)
return
		ch := make(chan string)

		// Reading messages
		go func(chOut chan<- string) {
			//for {
				//var msg *message.WSMessage
				//err := c.ReadJSON(&msg)
				//// TODO:2017/04/19 22:44:18 main3.go:135: Read error: websocket: close 1006 (abnormal closure): unexpected EOF
				//if err != nil {
				//	println("ws:", err.Error())
				//	return
				//}
				//println(msg.Type, msg.Text)
				//switch msg.Type {
				//case hconst.SYSTEM_MESSAGE:
				//	handleSystemMessage(u, msg.Message)
				//case hconst.OWN_MESSAGE:
				//	if u.RoomID == "" {
				//		continue
				//	}
				//
				//	err = rooms.Broadcast(u, msg.Message)
				//	if err != nil {
				//		log.Error(r.Context(), herrors.Wrap(err, "ctx", "Write error"))
				//		continue
				//	}
				//}
			//}
		}(ch)
		// Send a test message in 5
		go func() {
			time.Sleep(time.Second * 2)
			out := socket.Message{
				Type: "s",
				Text: "c",
			}

			j, err := json.Marshal(out)
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
				return
			}
			c.WriteMessage(websocket.TextMessage, j)
		}()
		// Writing messages
		func(chIn <-chan string) {
			for {
				msg := <-chIn
				fmt.Println("received message", msg)
			}
		}(ch)
	}
}
