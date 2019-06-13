package roomusersocketshierarchytest

import (
	"fmt"
)

type Socket struct {
	ID int
}

type User struct {
	UUID string
	Sockets []Socket
}

type Room struct {
	Users []*User
}

func main() {
	s1, s2 := Socket{ID: 1}, Socket{ID: 2}

	u1, u2 := &User{UUID: "1", Sockets: []Socket{s1}}, &User{UUID: "2", Sockets: []Socket{s2}}

	r := Room{Users: []*User{u1, u2}}

	for _, u := range r.Users {
		for _, s := range u.Sockets {
			fmt.Printf("User UUID: %s -> Socket ID %d\n", u.UUID, s.ID			)
		}
	}

	s3 := Socket{ID: 3}
	u1.Sockets = append(u1.Sockets, s3)

	for _, u := range r.Users {
		for _, s := range u.Sockets {
			fmt.Printf("User UUID: %s -> Socket ID %d\n", u.UUID, s.ID			)
		}
	}
}
