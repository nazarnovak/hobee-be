package api

import (
	"net/http"
)

func WS(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//ctx := r.Context()
		//
		//userIdStr := LoggedInUserId(r, secret)
		//if userIdStr == "" {
		//	log.Error(ctx, herrors.New("Not logged in"))
		//	ResponseJSONError(ctx, w, "Not logged in", http.StatusBadRequest)
		//	return
		//}
		//
		//c, err := upgrader.Upgrade(w, r, nil)
		//if err != nil {
		//	log.Critical(ctx, herrors.Wrap(err, "userid", userIdStr))
		//	return
		//}
		//
		//userID, err := strconv.ParseInt(userIdStr, 10, 64)
		//if err != nil {
		//	log.Critical(ctx, herrors.Wrap(err, "userid", userIdStr))
		//}
		//
		//s := socket.New(userID, c)
		//
		//go s.ReadFrom(ctx)
		//go s.WriteTo(ctx)
		//
		//status := socket.Status(userID)
		//switch status {
		//case user.Idle:
		//	// User is idle, do nothing
		//case user.Searching:
		//	// User is chatting, send that back to the socket so UI can switch to search mode
		//case user.Chatting:
		//	// User is chatting, send that back to the current socket so UI can switch to chat mode
		//case user.Disconnected:
		//	// User is chatting, send that back to the current socket so UI can switch to ended chat mode
		//case user.Unknown:
		//	err := herrors.New("Something went wrong when getting status", "userid", userID)
		//	log.Critical(ctx, err)
		//default:
		//	err := herrors.New("Something went super wrong when getting status", "userid", userID)
		//	log.Critical(ctx, err)
		//}
		// Try to join in the search pool if already there

		// Otherwise do nothing and rely on the socket to handle operations

// 1) Check if user already in search - add this socket to search as well under the same user
// 2) If they're a part of a conversation that is ongoing - connect to it
// 3) Search for a new match
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
