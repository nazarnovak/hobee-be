package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"

	"hobee-be/models"
	"hobee-be/pkg/hconst"
	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
	"hobee-be/pkg/matcher"
	"hobee-be/pkg/rooms"
	"hobee-be/pkg/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WS(w http.ResponseWriter, r *http.Request) error {
	// Check the cookie here if the user is even logged in?

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return herrors.Wrap(err)
	}

	//pinger to close the room if connection is closed
	// TODO: How to 1) Set the room for the users once he's connected; 2) Disconnect the user from the room if either he or his
	// buddy disconnects?
	u := &models.WSUser{User: user.User{ID: 1, Group: 2, RoomID: "", Paired: make(chan bool)}, Socket: c}

	for {
		var msg *models.WSMessage
		err := c.ReadJSON(&msg)
		// TODO:2017/04/19 22:44:18 main3.go:135: Read error: websocket: close 1006 (abnormal closure): unexpected EOF
		if err != nil {
			// When you close the window and you were in a room
			if websocket.IsCloseError(err, websocket.CloseGoingAway) && u.RoomID != "" {
				rooms.Close(u.RoomID)
				return nil
			}
			return herrors.Wrap(err, "ctx", "read error")
		}

		switch msg.Type {
		case hconst.SYSTEM_MESSAGE:
			handleSystemMessage(u, msg.Message)
		case hconst.OWN_MESSAGE:
			if u.RoomID == "" {
				continue
			}

			err = rooms.Broadcast(u, msg.Message)
			if err != nil {
				log.Error(r.Context(), herrors.Wrap(err, "ctx", "Write error"))
				continue
			}
		}
	}

	return nil
}

func handleSystemMessage(u *models.WSUser, msg string) {
	switch msg {
	case hconst.SYS_S:
		// TODO: What if you send more than 1 search messages, will it break here?
		// TODO: How to actually verify if user currently matched or not matched?
		//if cl.Room != "" {
		//	return
		//}

		if u.RoomID != "" {
			return
		}

		matcher.AddUser(u)

		<-u.Paired
	case hconst.SYS_DC:
		if u.RoomID == "" {
			return
		}

		rooms.Close(u.RoomID)

		//TODO: Disconnect from the other user
		//r, ok := rooms[cl.Room]
		//if !ok {
		//	log.Println("Error: room not found")
		//	return
		//}
		//
		//r.Close()
	}
}

//func pinger(cl *Client) {
//	t := time.NewTicker(5 * time.Second)
//
//	for {
//		<- t.C
//
//		err := cl.Connection.WriteMessage(websocket.PingMessage, []byte(""));
//		if err == nil {
//			continue
//		}
//
//		t.Stop()
//		cl.Close()
//		log.Println(err)
//		return
//	}
//
//}
