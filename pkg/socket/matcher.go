package socket

import (
	"fmt"
	"time"
)

var (
	sockets []*Socket
)

func searchAdd(s *Socket) {
	fmt.Printf("Added a socket: %+v\n", s.UUID)
	mutex.Lock()
	sockets = append(sockets, s)
	mutex.Unlock()
}

func searchRemove(s *Socket) {
	mutex.Lock()
	for k, socket := range sockets {
		if socket != s {
			continue
		}

		sockets = append(sockets[:k], sockets[k+1:]...)
		break
	}
	mutex.Unlock()
}

func getMatchingSockets(sp chan<- [2]*Socket) {
	if len(sockets) < 2 {
		return
	}
fmt.Println("Matching 2 sockets")
	sp <- [2]*Socket{sockets[0], sockets[1]}

	searchRemove(sockets[1])
	searchRemove(sockets[0])
}

func Matcher(sp chan<- [2]*Socket) {
	for {
		getMatchingSockets(sp)
		time.Sleep(1 * time.Second)
	}
}
