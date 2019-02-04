package user

type User struct {
	ID     int64
	Group  int64
	RoomID string
	Paired chan bool
}
