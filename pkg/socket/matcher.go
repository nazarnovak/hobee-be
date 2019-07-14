package socket

import (
	"context"
	"fmt"
	"time"

	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

var (
	searchingUsers = map[string]*User{}
)

func searchAdd(u *User) {
	fmt.Printf("Added a user: %+v\n", u.UUID)

	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	// Reset the roomUUID if user was connected in a room before
	u.RoomUUID = ""

	if _, ok := searchingUsers[u.UUID]; !ok {
		searchingUsers[u.UUID] = u
	}
	UpdateStatus(u.UUID, statusSearching)
}

func searchRemove(uuid string) {
	matcherMutex.Lock()
	defer matcherMutex.Unlock()

	delete(searchingUsers, uuid)
}

func getMatchingUsers(sp chan<- [2]*User) {
	if len(searchingUsers) < 2 {
		return
	}

	matchedUsers := [2]*User{}
	matched := 0

	for _, searchingUser := range searchingUsers {
		if matched > 1 {
			break
		}

		if searchingUser.Status != statusSearching {
			log.Critical(context.Background(), herrors.New("Expecting user status to be searching", "status",
				searchingUser.Status))
			continue
		}

		matchedUsers[matched] = searchingUser
		matched++
	}

	sp <- matchedUsers

	searchRemove(matchedUsers[1].UUID)
	searchRemove(matchedUsers[0].UUID)
}

func Matcher(sp chan<- [2]*User) {
	for {
		getMatchingUsers(sp)
		time.Sleep(1 * time.Second)
	}
}
