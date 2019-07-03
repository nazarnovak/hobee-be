package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/socket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Prevents requests from random sites directly to handshake with WS
var allowedOrigins = []string{
	// Local
	"http://localhost:8080",
	// Heroku
	"https://hobee.herokuapp.com",
	// Testing with ngrok
	"https://b518cf15.ngrok.io",
}

func GOT(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if err := checkOrigin(r); err != nil {
			log.Critical(ctx, err)
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		uuidStr, err := getCookieUUID(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if uuidStr == "" {
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

fmt.Println("New socket connected at:", time.Now().UTC().String())
		s := socket.New(c)

		user := socket.AttachSocketToUser(uuidStr, s)
// Creates a new socket with pkg/socket
// Run reader, which when received "search" command should somehow add the socket to search. That means
// we will import pkg/matcher
// pkg/matcher has references to pkg/socket.Socket as a part of Add/Remove functions
		go user.Reader(ctx, s)
		go user.Writer(ctx, s)

		// Initial message - to notify if user is already connected to a room and the socket should join it, or there is no
		// existing room and the user should search for a match
		msg := socket.Message{
			Type: socket.MessageTypeSystem,
			Text: socket.SystemSearch,
			AuthorUUID: uuidStr,
			Timestamp:  time.Now().UTC(),
		}

		roomUUID := socket.UserInARoomUUID(uuidStr)
		// If we don't find user in existing rooms - we notify FE about it and "allow" it to go into search mode
		if roomUUID != "" {
			active, err := socket.IsRoomActive(roomUUID)
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
				return
			}

			msg.Type = socket.MessageTypeActivity
			msg.Text = socket.ActivityRoomActive
			if !active {
				msg.Text = socket.ActivityRoomInactive
			}
		}

		o, err := json.Marshal(msg)
		if err != nil {
			log.Critical(ctx, err)
			return
		}

		s.Send <- o
return
	//	ch := make(chan string)
	//
	//	// Reading messages
	//	go func(chOut chan<- string) {
	//		//for {
	//			//var msg *message.WSMessage
	//			//err := c.ReadJSON(&msg)
	//			//// TODO:2017/04/19 22:44:18 main3.go:135: Read error: websocket: close 1006 (abnormal closure): unexpected EOF
	//			//if err != nil {
	//			//	println("ws:", err.Error())
	//			//	return
	//			//}
	//			//println(msg.Type, msg.Text)
	//			//switch msg.Type {
	//			//case hconst.SYSTEM_MESSAGE:
	//			//	handleSystemMessage(u, msg.Message)
	//			//case hconst.OWN_MESSAGE:
	//			//	if u.RoomID == "" {
	//			//		continue
	//			//	}
	//			//
	//			//	err = rooms.Broadcast(u, msg.Message)
	//			//	if err != nil {
	//			//		log.Error(r.Context(), herrors.Wrap(err, "ctx", "Write error"))
	//			//		continue
	//			//	}
	//			//}
	//		//}
	//	}(ch)
	//	// Send a test message in 5
	//	go func() {
	//		time.Sleep(time.Second * 2)
	//		out := socket.Message{
	//			Type: "s",
	//			Text: "c",
	//		}
	//
	//		j, err := json.Marshal(out)
	//		if err != nil {
	//			log.Critical(ctx, herrors.Wrap(err))
	//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
	//			return
	//		}
	//		c.WriteMessage(websocket.TextMessage, j)
	//	}()
	//	// Writing messages
	//	func(chIn <-chan string) {
	//		for {
	//			msg := <-chIn
	//			fmt.Println("received message", msg)
	//		}
	//	}(ch)
	}
}

func checkOrigin(r *http.Request) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return herrors.New("Empty origin")
	}

	found := false
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			found = true
		}
	}

	if !found {
		return herrors.New("Origin not allowed", "origin", origin)
	}

	return nil
}
