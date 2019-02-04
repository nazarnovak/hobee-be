package models

type WSMessage struct {
	Type    int       `json:"type"`
	Message string    `json:"msg"`
}
