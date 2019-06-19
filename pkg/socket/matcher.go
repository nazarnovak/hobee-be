package socket

import (
	"fmt"
	"time"
)

var (
	searchingUsers = map[string]*User{}
)

func searchAdd(u *User) {
	fmt.Printf("Added a user: %+v\n", u.UUID)
	matcherMutex.Lock()
	if _, ok := searchingUsers[u.UUID]; !ok {
		searchingUsers[u.UUID] = u
	}
	UpdateStatus(u.UUID, statusSearching)
	matcherMutex.Unlock()
}

func searchRemove(uuid string) {
	matcherMutex.Lock()
	delete(searchingUsers, uuid)
	matcherMutex.Unlock()
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
