package controllers

import (
	"context"
	"fmt"
	"hobee-be/api"
	"net/http"

	"github.com/gorilla/websocket"

	//"hobee-be/models"
	//"hobee-be/pkg/hconst"
	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
	//"hobee-be/pkg/matcher"
	//"hobee-be/pkg/rooms"
	//"hobee-be/pkg/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type socketOnline struct {
	userID int64
	conn *websocket.Conn
}

func (so *socketOnline) readFrom(ctx context.Context) {
	defer func() {
		// If socket dies - remove the current socket address from the room, so clean up
		so.conn.Close()
	}()

	// TODO: What each of these do?
	//so.conn.SetReadLimit(maxMessageSize)
	//so.conn.SetReadDeadline(time.Now().Add(pongWait))
	//so.conn.SetPongHandler(func(string) error { so.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := so.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Critical(ctx, herrors.Wrap(err))
			}
			break
		}
fmt.Println("Received message", string(message))
		// TODO: Broadcast message to the room if connected to the room
		//c.hub.broadcast <- message
	}
}

func (so *socketOnline) writeTo(ctx context.Context) {
	//ticker := time.NewTicker(pingPeriod)
	//defer func() {
	//	ticker.Stop()
	//	c.conn.Close()
	//}()
	//for {
	//	select {
	//	case message, ok := <-c.send:
	//		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//		if !ok {
	//			// The hub closed the channel.
	//			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	//			return
	//		}
	//
	//		w, err := c.conn.NextWriter(websocket.TextMessage)
	//		if err != nil {
	//			return
	//		}
	//		w.Write(message)
	//
	//		// Add queued chat messages to the current websocket message.
	//		n := len(c.send)
	//		for i := 0; i < n; i++ {
	//			w.Write(newline)
	//			w.Write(<-c.send)
	//		}
	//
	//		if err := w.Close(); err != nil {
	//			return
	//		}
	//	case <-ticker.C:
	//		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	//		if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
	//			return
	//		}
	//	}
	//}
}

func WS(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userIdStr := api.LoggedInUserId(r, secret)
		if userIdStr == "" {
			log.Error(ctx, herrors.New("Not logged in"))
			api.ResponseJSONError(ctx, w, "Not logged in", http.StatusBadRequest)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// Log.Error
			//println("ws1:", err.Error())
			//if err := api.ResponseJSONError(w, "Already logged in", http.StatusInternalServerError); err != nil {
			//
			//	println("ws2:", err.Error())
			//}
			return
		}

		so := &socketOnline{userID: 1, conn: c}

		//existingUserRoom := getExistingUserRoom(userID)
		//if existingUserRoom != "" {
		//	// TODO: If you're already in a chat - just point socket to that chat and continue the conversation
		//	return
		//}

		// TODO: If you're not in any rooms yet - search for a new match

		go so.readFrom(ctx)
		go so.writeTo(ctx)
		////pinger to close the room if connection is closed
		//// TODO: How to 1) Set the room for the users once he's connected; 2) Disconnect the user from the room if either he or his
		//// buddy disconnects?
		//u := &models.WSUser{User: user.User{ID: 1, Group: 2, RoomID: "", Paired: make(chan bool)}, Socket: c}
		//
//		for {
//			var msg *models.WSMessage
//			err := c.ReadJSON(&msg)
//			// TODO:2017/04/19 22:44:18 main3.go:135: Read error: websocket: close 1006 (abnormal closure): unexpected EOF
//			if err != nil {
//				println("ws:", err.Error())
//				return
//				//// When you close the window and you were in a room
//				//if websocket.IsCloseError(err, websocket.CloseGoingAway) && u.RoomID != "" {
//				//	rooms.Close(u.RoomID)
//				//	return nil
//				//}
//				//return herrors.Wrap(err, "ctx", "read error")
//			}
//println(msg.Message)
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
//		}
		//
		//return nil
	}
}

//func handleSystemMessage(u *models.WSUser, msg string) {
//		switch msg {
//		case hconst.SYS_S:
//			// TODO: What if you send more than 1 search messages, will it break here?
//			// TODO: How to actually verify if user currently matched or not matched?
//			//if cl.Room != "" {
//			//	return
//			//}
//
//			if u.RoomID != "" {
//				return
//			}
//
//			matcher.AddUser(u)
//
//			<-u.Paired
//		case hconst.SYS_DC:
//			if u.RoomID == "" {
//				return
//			}
//
//			rooms.Close(u.RoomID)
//
//			//TODO: Disconnect from the other user
//			//r, ok := rooms[cl.Room]
//			//if !ok {
//			//	log.Println("Error: room not found")
//			//	return
//			//}
//			//
//			//r.Close()
//		}
//	}
//}

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
