package user

type Status string

const (
	Unknown      Status = "unknown"
	Idle         Status = "idle"
	Searching    Status = "searching"
	Chatting     Status = "chatting"
	Disconnected Status = "disconnected"
)
