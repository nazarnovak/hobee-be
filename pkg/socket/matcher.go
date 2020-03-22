package socket

import (
	"time"
)

var (
	searchingUsers = map[string]*User{}
)

func searchAdd(u *User) {
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

		// TODO: This removes user when they go into search and then close all tabs. Maybe worth leaving for now
		//if searchingUser.Status != statusSearching {
		//	log.Critical(context.Background(), herrors.New("Expecting user status to be searching", "status",
		//		searchingUser.Status))
		//	continue
		//}

		// Do not match with the last user you talked to if it was 1 minute ago
		if matched > 0 && len(matchedUsers[0].UserHistory) > 0 {
			if matchedUsers[0].UserHistory[0].UserUUID == searchingUser.UUID &&
				time.Now().UTC().Sub(matchedUsers[0].UserHistory[0].LastMessage) < time.Second * 10 {
				continue
			}
		}

		matchedUsers[matched] = searchingUser
		matched++
	}

	// We couldn't pair anyone together, return
	if matched < 2 {
		return
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
