package models

import (
	"github.com/gorilla/websocket"

	"hobee-be/pkg/user"
)

type WSUser struct {
	user.User
	Socket *websocket.Conn
}
